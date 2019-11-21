package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbclient *mongo.Client
)

const (
	tenSeconds = 10 * time.Second
	twoSeconds = 2 * time.Second
)

// Connect ...
// Connect to mongodb via uri
func Connect(uri string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), tenSeconds)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return false
	}
	dbclient = client
	cancel()
	return isConnected()
}

func isConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), twoSeconds)
	err := dbclient.Ping(ctx, readpref.Primary())
	cancel()
	if err != nil {
		return false
	}
	return true
}

func getCollection(name string) *mongo.Collection {
	collection := dbclient.Database("DATABASE").Collection(name)
	return collection
}

// UpdateState ...
// Update main state
func UpdateState(newState *State) {
	state := getCollection("state")
	ctx, cancel := context.WithTimeout(context.Background(), twoSeconds)
	_, _ = state.InsertOne(ctx, newState)
	cancel()
	return
}
