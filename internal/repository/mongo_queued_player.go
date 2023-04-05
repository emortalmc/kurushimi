package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"kurushimi/internal/repository/model"
	"kurushimi/internal/repository/registrytypes"
	"time"
)

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

func (m *mongoRepository) CreateQueuedPlayers(ctx context.Context, players []*model.QueuedPlayer) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	converted := make([]interface{}, len(players))
	for i, player := range players {
		converted[i] = player
	}

	_, err := m.queuedPlayerCollection.InsertMany(ctx, converted)
	return err
}

func (m *mongoRepository) DeleteQueuedPlayer(ctx context.Context, playerId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := m.queuedPlayerCollection.DeleteOne(ctx, bson.M{"_id": playerId})
	return err
}

func (m *mongoRepository) DeleteAllQueuedPlayersById(ctx context.Context, playerIds []uuid.UUID) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := m.queuedPlayerCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": playerIds}})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (m *mongoRepository) GetQueuedPlayerById(ctx context.Context, playerId uuid.UUID) (*model.QueuedPlayer, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var player model.QueuedPlayer
	if err := m.queuedPlayerCollection.FindOne(ctx, bson.M{"_id": playerId}).Decode(&player); err != nil {
		return nil, err
	}

	return &player, nil
}

func (m *mongoRepository) GetAllQueuedPlayersByIds(ctx context.Context, playerIds []uuid.UUID) ([]*model.QueuedPlayer, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := m.queuedPlayerCollection.Find(ctx, bson.M{"_id": bson.M{"$in": playerIds}})
	if err != nil {
		return nil, err
	}

	var players []*model.QueuedPlayer
	err = cursor.All(ctx, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

func (m *mongoRepository) SetMapIdOfQueuedPlayer(ctx context.Context, playerId uuid.UUID, mapId string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := m.queuedPlayerCollection.UpdateOne(ctx, bson.M{"_id": playerId}, bson.M{"$set": bson.M{"mapId": mapId}})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func createCodecRegistry() *bsoncodec.Registry {
	return bson.NewRegistryBuilder().
		RegisterTypeEncoder(registrytypes.UUIDType, bsoncodec.ValueEncoderFunc(registrytypes.UuidEncodeValue)).
		RegisterTypeDecoder(registrytypes.UUIDType, bsoncodec.ValueDecoderFunc(registrytypes.UuidDecodeValue)).
		Build()
}
