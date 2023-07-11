package dev.emortal.api.service.matchmaker;

import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

import java.util.Collection;
import java.util.UUID;

public interface MatchmakerService {

    @Nullable PlayerQueueInfo getPlayerQueueInfo(@NotNull UUID playerId);

    @NotNull QueuePlayerResult queuePlayer(@NotNull String gameModeId, @NotNull UUID playerId, @NotNull QueueOptions options);

    default @NotNull QueuePlayerResult queuePlayer(@NotNull String gameModeId, @NotNull UUID playerId) {
        return this.queuePlayer(gameModeId, playerId, QueueOptions.DEFAULT);
    }

    @NotNull DequeuePlayerResult dequeuePlayer(@NotNull UUID playerId);

    void sendPlayerToLobby(@NotNull UUID playerId, boolean sendParty);

    void sendPlayersToLobby(@NotNull Collection<UUID> playerIds, boolean sendParties);

    void queueInitialLobby(@NotNull UUID playerId);

    @NotNull ChangeMapVoteResult changeMapVote(@NotNull UUID playerId, @NotNull String newMapId);
}
