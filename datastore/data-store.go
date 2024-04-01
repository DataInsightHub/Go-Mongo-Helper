package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	DataStore struct {
		Client   *mongo.Client
		Database *mongo.Database
		Ctx      context.Context
	}
)

func NewDataStore(mongoDbUri, mongoDbName string, dataStoreOptions ...DataStoreOptions) (*DataStore, error) {
	ops := &dataStoreOption{
		timeout: 10 * time.Second,
		usePing: true,
	}

	for _, datastoreOption := range dataStoreOptions {
		datastoreOption.apply(ops)
	}

	ctx, _ := context.WithTimeout(context.Background(), ops.timeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		return nil, err
	}

	if ops.usePing {
		// Check connection
		err = client.Ping(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	db := client.Database(mongoDbName)

	store := &DataStore{
		Client:   client,
		Database: db,
		Ctx:      ctx,
	}

	return store, nil
}

func (dataStore *DataStore) Disconnect() error {
	err := dataStore.Client.Disconnect(dataStore.Ctx)
	if err != nil {
		return err
	}

	return nil
}
