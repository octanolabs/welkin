package store

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

// GetState ...
// Returns app state from db
func GetState() State {
	var state State
	filter := bson.D{}
	collection := getCollection("state")
	ctx, cancel := context.WithTimeout(context.Background(), twoSeconds)
	err := collection.FindOne(ctx, filter).Decode(&state)
	if err != nil {
		log.Fatal(err)
	}
	cancel()
	return state
}

// SetState ...
// Update main state
func SetState(newState State) {
	state := getCollection("state")
	ctx, cancel := context.WithTimeout(context.Background(), twoSeconds)
	filter := bson.D{}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err := state.FindOneAndUpdate(ctx, filter, newState, opts)
	if err != nil {
		log.Fatal(err)
	}
	cancel()
	return
}
