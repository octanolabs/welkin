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

func Connect(uri string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := dbclient.Ping(ctx, readpref.Primary())
	cancel()
	if err != nil {
		return false
	}
	return true
}
