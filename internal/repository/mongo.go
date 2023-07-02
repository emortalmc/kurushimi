package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
	"kurushimi/internal/config"
	"kurushimi/internal/repository/registrytypes"
	"sync"
	"time"
)

var _ Repository = &mongoRepository{}

type mongoRepository struct {
	client   *mongo.Client
	database *mongo.Database

	queuedPlayerCollection *mongo.Collection
	ticketCollection       *mongo.Collection
	pendingMatchCollection *mongo.Collection
	backfillCollection     *mongo.Collection
}

func NewMongoRepository(ctx context.Context, logger *zap.SugaredLogger, cfg *config.MongoDBConfig) (Repository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI).SetRegistry(createCodecRegistry()))
	if err != nil {
		return nil, err
	}

	database := client.Database(databaseName)
	repo := &mongoRepository{
		client:   client,
		database: database,

		queuedPlayerCollection: database.Collection(queuedPlayerCollectionName),
		ticketCollection:       database.Collection(ticketCollectionName),
		pendingMatchCollection: database.Collection(pendingMatchCollectionName),
		backfillCollection:     database.Collection(backfillCollectionName),
	}

	repo.createIndexes(ctx)
	logger.Infow("created mongo indexes")

	return repo, nil
}

var (
	// no queued player indexes

	ticketIndexes = []mongo.IndexModel{
		{
			Keys:    bson.M{"playerIds": 1},
			Options: options.Index().SetName("playerIds"),
		},
		{
			Keys:    bson.M{"gameModeId": 1},
			Options: options.Index().SetName("gameModeId"),
		},
		{
			Keys:    bson.D{{Key: "gameModeId", Value: 1}, {Key: "inPendingMatch", Value: 1}},
			Options: options.Index().SetName("gameModeId_inPendingMatch"),
		},
		{
			Keys:    bson.D{{Key: "gameModeId", Value: 1}, {Key: "removals", Value: 1}},
			Options: options.Index().SetName("gameModeId_removals"),
		},
		{
			Keys:    bson.M{"partyId": 1},
			Options: options.Index().SetName("partyId"),
		},
	}

	pendingMatchIndexes = []mongo.IndexModel{
		{
			Keys:    bson.M{"gameModeId": 1},
			Options: options.Index().SetName("gameModeId"),
		},
		{
			Keys:    bson.M{"ticketIds": 1},
			Options: options.Index().SetName("ticketIds"),
		},
	}
)

func (m *mongoRepository) createIndexes(ctx context.Context) {
	collIndexes := map[*mongo.Collection][]mongo.IndexModel{
		m.ticketCollection:       ticketIndexes,
		m.pendingMatchCollection: pendingMatchIndexes,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(collIndexes))

	for coll, indexes := range collIndexes {
		go func(coll *mongo.Collection, indexes []mongo.IndexModel) {
			defer wg.Done()
			_, err := m.createCollIndexes(ctx, coll, indexes)
			if err != nil {
				panic(fmt.Sprintf("failed to create indexes for collection %s: %s", coll.Name(), err))
			}
		}(coll, indexes)
	}

	wg.Wait()
}

func (m *mongoRepository) createCollIndexes(ctx context.Context, coll *mongo.Collection, indexes []mongo.IndexModel) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return 0, err
	}

	return len(result), nil
}

func (m *mongoRepository) HealthCheck(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return m.client.Ping(ctx, nil)
}

func (m *mongoRepository) ExecuteTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) error) error {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOpts := options.Transaction().SetWriteConcern(wc)

	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}, txnOpts)

	return err
}

func createCodecRegistry() *bsoncodec.Registry {
	return bson.NewRegistryBuilder().
		RegisterTypeEncoder(registrytypes.UUIDType, bsoncodec.ValueEncoderFunc(registrytypes.UuidEncodeValue)).
		RegisterTypeDecoder(registrytypes.UUIDType, bsoncodec.ValueDecoderFunc(registrytypes.UuidDecodeValue)).
		Build()
}