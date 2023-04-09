package director

import (
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	v1 "agones.dev/agones/pkg/client/clientset/versioned/typed/allocation/v1"
	"context"
	"errors"
	"fmt"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/internal/config/allocationselector"
	"kurushimi/internal/kafka"
	matchfunction2 "kurushimi/internal/matchfunction"
	"kurushimi/internal/repository"
	"kurushimi/internal/repository/model"
	"kurushimi/internal/utils/protoutils"
	"kurushimi/pkg/pb"
	"log"
	"sort"
	"sync"
	"time"
)

type Director interface {
	Start(ctx context.Context)
}

type directorImpl struct {
	logger *zap.SugaredLogger

	repo     repository.Repository
	notifier kafka.Notifier

	allocationClient v1.GameServerAllocationInterface

	configs map[string]*liveconfig.GameModeConfig
}

func New(logger *zap.SugaredLogger, repo repository.Repository, notifier kafka.Notifier, namespace string,
	agonesClient *versioned.Clientset, cfgController liveconfig.GameModeConfigController) Director {

	d := &directorImpl{
		logger: logger,

		repo:     repo,
		notifier: notifier,

		allocationClient: agonesClient.AllocationV1().GameServerAllocations(namespace),

		configs: cfgController.GetConfigs(),
	}

	cfgController.AddGlobalUpdateListener(d.onGameModeConfigUpdate)
	return d
}

func (d *directorImpl) Start(ctx context.Context) {
	for _, c := range d.configs {
		go func(config *liveconfig.GameModeConfig) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			for {
				lastRunTime := time.Now()
				d.run(ctx, config)

				// Wait for the next run
				timeSinceLastRun := time.Since(lastRunTime)
				if timeSinceLastRun < config.MatchmakerInfo.Rate {
					time.Sleep(config.MatchmakerInfo.Rate - timeSinceLastRun)
				}
			}
		}(c)
	}
}

func (d *directorImpl) run(ctx context.Context, originalConfig *liveconfig.GameModeConfig) {
	temp := *originalConfig
	config := &temp

	// account for deletions
	err := d.processDequeues(ctx, config)
	if err != nil {
		d.logger.Errorw("failed to process dequeues", "error", err)
		return
	}

	// run match function
	matches, err := d.runMatchFunction(ctx, config)
	if err != nil {
		d.logger.Errorw("failed to run match function", "error", err)
		return
	}

	if len(matches) == 0 {
		return
	}

	d.logger.Debugw("match function returned matches", "count", len(matches))
	d.logger.Debugw("matches", "matches", matches)

	// todo allocate gameservers
}

func (d *directorImpl) processDequeues(ctx context.Context, config *liveconfig.GameModeConfig) error {
	configId := config.Id

	tickets, err := d.repo.GetTicketsWithDequeueRequest(ctx, configId)
	if err != nil {
		return err
	}

	// ticketIdsToUpdate contains the ids of every ticket that isn't in the ticketIdsToDelete slice
	ticketIdsToUpdate := make([]primitive.ObjectID, 0)
	playerIdsToDelete := make([]uuid.UUID, 0)
	ticketIdsToDelete := make([]primitive.ObjectID, 0)

	for _, ticket := range tickets {
		removals := ticket.Removals

		if removals.MarkedForRemoval {
			playerIdsToDelete = append(playerIdsToDelete, ticket.PlayerIds...)
			ticketIdsToDelete = append(ticketIdsToDelete, ticket.Id)
		} else {
			playerIdsToDelete = append(playerIdsToDelete, removals.PlayersForRemoval...)
			ticketIdsToUpdate = append(ticketIdsToUpdate, ticket.Id)
		}
	}

	// Delete tickets
	if len(ticketIdsToDelete) > 0 {
		modified, err := d.repo.DeleteAllTicketsById(ctx, ticketIdsToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete tickets: %w", err)
		}

		if int(modified) != len(ticketIdsToDelete) {
			d.logger.Warnw("deleted tickets count does not match expected count", "deleted", modified, "expected", len(ticketIdsToDelete))
		}

		// Delete tickets from any PendingMatches they are in.
		// NOTE: The modified count is irrelevant as we don't know if they were in any PendingMatches
		_, err = d.repo.RemoveTicketsFromPendingMatchesById(ctx, ticketIdsToDelete)
		if err != nil {
			return fmt.Errorf("failed to remove tickets from pending matches: %w", err)
		}

		// Send Kafka notifications
		for _, ticket := range tickets {
			if ticket.Removals.MarkedForRemoval {
				if err := d.notifier.TicketDeleted(ctx, ticket.ToProto(), pb.TicketDeletedMessage_MANUAL_DEQUEUE); err != nil {
					d.logger.Errorw("failed to send ticket deleted notification", "error", err)
				}
			}
		}
	}

	if len(ticketIdsToUpdate) > 0 {
		// Update tickets removal requests.
		modified, err := d.repo.ResetAllDequeueRequestsById(ctx, ticketIdsToUpdate)
		if err != nil {
			return fmt.Errorf("failed to reset dequeue requests: %w", err)
		}

		if int(modified) != len(ticketIdsToUpdate) {
			d.logger.Warnw("updated tickets count does not match expected count", "updated", modified, "expected", len(ticketIdsToUpdate))
		}
	}

	if len(playerIdsToDelete) > 0 {
		// Delete players
		modified, err := d.repo.DeleteAllQueuedPlayersById(ctx, playerIdsToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete players: %w", err)
		}

		if int(modified) != len(playerIdsToDelete) {
			d.logger.Warnw("deleted players count does not match expected count", "deleted", modified, "expected", len(playerIdsToDelete))
		}
		return nil
	}

	return err
}

func (d *directorImpl) runMatchFunction(ctx context.Context, config *liveconfig.GameModeConfig) ([]*pb.Match, error) {
	// NOTE: these tickets are ALL the tickets for this gamemode, even ones already in a PendingMatch
	tickets, err := d.repo.GetTicketsByGameMode(ctx, config.Id)
	if err != nil {
		return nil, err
	}

	d.logger.Debugw("matchmaker running", "gamemode", config.Id, "tickets", len(tickets), "method", config.MatchmakerInfo.MatchMethod)

	// make matches
	var matches []*pb.Match
	switch config.MatchmakerInfo.MatchMethod {
	case liveconfig.MatchMethodCountdown:
		pendingMatches, err := d.repo.GetPendingMatchesByGameMode(ctx, config.Id)
		if err != nil {
			return nil, err
		}

		// Create a pending matches map
		pendingMatchesMap := make(map[primitive.ObjectID]*model.PendingMatch)
		for _, pendingMatch := range pendingMatches {
			pendingMatchesMap[pendingMatch.Id] = pendingMatch
		}

		log.Printf("pending matches: %+v", pendingMatchesMap)

		// Create a ticket map
		ticketMap := make(map[primitive.ObjectID]*model.Ticket)
		for _, ticket := range tickets {
			ticketMap[ticket.Id] = ticket
		}

		log.Printf("pending matches before: %v", pendingMatchesMap)
		// Clean up existing pending matches
		deletedPending := matchfunction2.CountdownRemoveInvalidPendingMatches(pendingMatchesMap, ticketMap, config)
		log.Printf("pending matches after: %v", pendingMatchesMap)

		// Delete matches from db and notify with reason cancelled
		if len(deletedPending) > 0 {
			//ticketUpdates := make(map[primitive.ObjectID]bool)

			for _, match := range deletedPending {
				if err := d.notifier.PendingMatchDeleted(ctx, match, pb.PendingMatchDeletedMessage_CANCELLED); err != nil {
					d.logger.Errorw("failed to send pending match deleted notification", "error", err)
				}

				for _, ticketId := range match.TicketIds {
					ticket := ticketMap[ticketId]
					if err := d.notifier.TicketUpdated(ctx, ticket); err != nil {
						d.logger.Errorw("failed to send ticket updated notification", "error", err)
					}
				}

				//for _, ticketId := range match.TicketIds {
				//	ticketUpdates[ticketId] = false
				//}
			}

			// Update the tickets in the DB
			//if _, err := d.repo.MassUpdateTicketInPendingMatch(ctx, ticketUpdates); err != nil {
			//	return nil, fmt.Errorf("failed to update tickets in pending match: %w", err)
			//}

			// Delete the pending matches from the db
			deletedIds := make([]primitive.ObjectID, 0)
			for _, match := range deletedPending {
				deletedIds = append(deletedIds, match.Id)
			}
			if err := d.repo.DeletePendingMatches(ctx, deletedIds); err != nil {
				return nil, err
			}
		}

		createdPending, updatedPending, deletedPending, createdMatches, err := matchfunction2.RunCountdown(d.logger, ticketMap, pendingMatchesMap, config)
		if err != nil {
			return nil, err
		}

		d.logger.Debugw("countdown matchmaker results",
			"gamemode", config.Id,
			"pending", len(updatedPending),
			"deleted", len(deletedPending),
			"matches", len(createdMatches),
		)

		matches = createdMatches // assign the created matches to the return value

		// handle the pending match stuff
		if len(updatedPending) > 0 {
			err = d.repo.UpdatePendingMatches(ctx, updatedPending)
			if err != nil {
				return nil, err
			}

			for _, match := range updatedPending {
				if err := d.notifier.PendingMatchUpdated(ctx, match); err != nil {
					d.logger.Errorw("failed to notify pending match updated", "error", err)
				}
			}
		}

		// handle created pending matches (just notify rn)
		if len(createdPending) > 0 {
			err = d.repo.CreatePendingMatches(ctx, createdPending)
			if err != nil {
				return nil, err
			}

			for _, match := range createdPending {

				if err := d.notifier.PendingMatchCreated(ctx, match); err != nil {
					d.logger.Errorw("failed to notify pending match created", "error", err)
				}
			}
		}

		// Delete pending matches from DB and notify with reason match created
		if len(deletedPending) > 0 {
			deletedIds := make([]primitive.ObjectID, 0)
			for _, pending := range deletedPending {
				deletedIds = append(deletedIds, pending.Id)
			}

			err = d.repo.DeletePendingMatches(ctx, deletedIds)
			if err != nil {
				return nil, err
			}

			for _, match := range deletedPending {
				if err := d.notifier.PendingMatchDeleted(ctx, match, pb.PendingMatchDeletedMessage_MATCH_CREATED); err != nil {
					d.logger.Errorw("failed to notify pending match deleted", "error", err)
				}
			}
		}

		inPendingMatchUpdates := make(map[primitive.ObjectID]bool)
		for _, ticket := range tickets {
			if ticket.InternalUpdates != nil && ticket.InternalUpdates.InPendingMatchUpdated {
				inPendingMatchUpdates[ticket.Id] = ticket.InPendingMatch
			}
		}
		if len(inPendingMatchUpdates) > 0 {
			_, err = d.repo.MassUpdateTicketInPendingMatch(ctx, inPendingMatchUpdates)
			if err != nil {
				return nil, err
			}
		}

		// (Kafka) send notifications for updated tickets
		for _, ticket := range tickets {
			if ticket.InternalUpdates != nil {
				if err := d.notifier.TicketUpdated(ctx, ticket); err != nil {
					d.logger.Errorw("failed to send ticket updated notification", "error", err)
				}
			}
		}
	case liveconfig.MatchMethodInstant:
		matches, err = matchfunction2.RunInstant(tickets, config)
	}
	if err != nil {
		return nil, err
	}

	d.logger.Debugw("matchmaker finished", "gamemode", config.Id, "matches", len(matches))

	if len(matches) == 0 {
		return nil, nil
	}

	if len(config.Maps) > 0 {
		err = d.calculateMaps(ctx, matches)
		if err != nil {
			return nil, err
		}
	}

	// Create teams for Matches
	teamInfo := config.TeamInfo
	if teamInfo.TeamCount > 0 {
		for _, match := range matches {
			teams := createTeams(match.Tickets, teamInfo.TeamCount, teamInfo.TeamSize)
			match.Teams = teams
		}
	}

	// Assign a server for each match
	errorMap := d.allocateServers(ctx, config, matches)
	// TODO let's use the errorMap to do some retry logic and not delete the Tickets and QueuedPlayers

	// TODO remove this
	_ = errorMap

	// Notify of match creation
	for _, match := range matches {
		err = d.notifier.MatchCreated(ctx, match)
		if err != nil {
			d.logger.Errorw("error notifying of match creation", "match", match.Id, "error", err)
		}
	}

	// delete all Tickets and QueuedPlayers in Matches (not PendingMatches)
	for _, match := range matches {
		ticketIds := make([]primitive.ObjectID, 0)
		playerIds := make([]uuid.UUID, 0)
		for _, pbTicket := range match.Tickets {
			var ticket *model.Ticket
			for _, t := range tickets {
				if t.Id.Hex() == pbTicket.Id {
					ticket = t
					break
				}
			}

			// (Kafka) Notify of ticket deletion
			// TODO this is actually wrong. This match may have been deleted because there were no longer enough tickets to maintain the PendingMatch
			if err := d.notifier.TicketDeleted(ctx, pbTicket, pb.TicketDeletedMessage_MATCH_CREATED); err != nil {
				d.logger.Errorw("failed to notify ticket deleted", "error", err)
			}

			ticketIds = append(ticketIds, ticket.Id)

			playerIds = append(playerIds, ticket.PlayerIds...)
		}

		// Delete Tickets
		deletedCount, err := d.repo.DeleteAllTicketsById(ctx, ticketIds)
		if err != nil {
			return nil, err
		}

		if int(deletedCount) != len(ticketIds) {
			d.logger.Warnw("deleted tickets count does not match expected count", "deleted", deletedCount, "expected", len(ticketIds))
		}

		// Delete QueuedPlayers
		deletedCount, err = d.repo.DeleteAllQueuedPlayersById(ctx, playerIds)
		if err != nil {
			return nil, err
		}

		if int(deletedCount) != len(playerIds) {
			d.logger.Warnw("deleted players count does not match expected count", "deleted", deletedCount, "expected", len(playerIds))
		}
	}

	return matches, nil
}

// calculate map retrieves the map votes for those present in a Match
// and assigns the MapId field of a pb.Match
func (d *directorImpl) calculateMaps(ctx context.Context, matches []*pb.Match) error {
	for _, match := range matches {
		playerIds := make([]uuid.UUID, 0)

		for _, ticket := range match.Tickets {
			for _, playerId := range ticket.PlayerIds {
				parsedId, err := uuid.Parse(playerId)
				if err != nil {
					// Note: We're not returning the error as map selection isn't critical
					d.logger.Errorw("failed to parse player id", "playerId", playerId)
					continue
				}

				playerIds = append(playerIds, parsedId)
			}
		}

		players, err := d.repo.GetAllQueuedPlayersByIds(ctx, playerIds)
		if err != nil {
			return err
		}

		// map of map id to number of votes
		mapVotes := make(map[string]int)
		for _, player := range players {
			if player.MapId != nil {
				mapVotes[*player.MapId]++
			}
		}

		// find the map with the most votes
		var mostVotedMapId *string
		mostVotes := 0
		for mapId, votes := range mapVotes {
			if votes > mostVotes {
				mostVotedMapId = &mapId
				mostVotes = votes
			}
		}

		match.MapId = mostVotedMapId
	}

	return nil
}

// TODO this needs to be redone.
// We should favour filling all teams or make that behaviour configurable (e.g. a min teams or prefer spread setting)
func createTeams(tickets []*pb.Ticket, teamCount int, maxTeamSize int) []*pb.Match_MatchTeam {
	if teamCount == 0 {
		return nil
	}
	if teamCount == 1 {
		// Make one team with all players
		playerIds := make([]string, 0)
		for _, ticket := range tickets {
			playerIds = append(playerIds, ticket.PlayerIds...)
		}
		return []*pb.Match_MatchTeam{{PlayerIds: playerIds}}
	}

	if teamCount > 0 {
		// A party may be over the max team size, so we need to split them up into groups of <= max team size
		groups := make([][]string, 0)
		for _, ticket := range tickets {
			playerIds := make([]string, 0)
			copy(playerIds, ticket.PlayerIds)

			for len(playerIds) > 0 {
				// Split the party into groups of size <= max team size
				group := make([]string, 0)
				for len(playerIds) > 0 && len(group) < maxTeamSize {
					group = append(group, playerIds[0])
					playerIds = playerIds[1:]
				}
			}
		}

		// Sort groups from largest to smallest.
		// This can help improve matchmaking, although it's not a fix.
		sort.Slice(groups, func(i, j int) bool {
			return len(groups[i]) > len(groups[j])
		})

		teams := make([]*pb.Match_MatchTeam, 0)
		for _, group := range groups {
			groupSize := len(group)
			foundTeam := false
			for _, team := range teams {
				teamSpace := maxTeamSize - len(team.PlayerIds)
				if teamSpace >= groupSize {
					team.PlayerIds = append(team.PlayerIds, group...)
					foundTeam = true
					break
				}
			}

			if !foundTeam {
				// We can create a new team :)
				if len(teams) < teamCount {
					teams = append(teams, &pb.Match_MatchTeam{PlayerIds: group})
				} else {
					// Split up thr group
					// NOTE: This is a naive approach for now. We know there is space for every player, but there may not be
					// space in teams for this combination of players in groups. Therefore, we're going to split the group up
				}
			}
		}
	}

	return nil
}

// allocateServers allocates servers for the given matches
// returns: map of match id to error
// NOTE: this function blocks until all matches have been allocated
// TODO let's make a system to retry failed allocations
func (d *directorImpl) allocateServers(ctx context.Context, config *liveconfig.GameModeConfig, matches []*pb.Match) map[string]error {
	wg := sync.WaitGroup{}
	wg.Add(len(matches))

	matchErrors := make(map[string]error)

	for _, match := range matches {
		go func(match *pb.Match) {
			defer wg.Done()
			var selector *allocatorv1.GameServerAllocation
			switch config.MatchmakerInfo.SelectMethod {
			case liveconfig.SelectMethodAvailable:
				selector = allocationselector.CreateAvailableSelector(config, match)
			case liveconfig.SelectMethodPlayerCount:
				selector = allocationselector.CreatePlayerBasedSelector(config, match, protoutils.GetMatchPlayerCount(match))
			}

			response, err := d.allocationClient.Create(ctx, selector, v12.CreateOptions{})
			if err != nil {
				matchErrors[match.Id] = err
				return
			}

			allocation := response.Status
			if allocation.State != allocatorv1.GameServerAllocationAllocated {
				matchErrors[match.Id] = errors.New("allocation failed due to invalid state: " + string(allocation.State))
				return
			}

			// TODO should we be getting the port like this?
			// TODO also let's try a cluster only port
			match.Assignment = &pb.Assignment{
				ServerId:      allocation.GameServerName,
				ServerAddress: allocation.Address,
				ServerPort:    uint32(allocation.Ports[0].Port),
			}
		}(match)
	}

	wg.Wait()
	return matchErrors
}

func (d *directorImpl) onGameModeConfigUpdate(update liveconfig.ConfigUpdate[liveconfig.GameModeConfig]) {
	// TODO
}
