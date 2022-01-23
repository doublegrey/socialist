package datastore

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Datastore struct {
	Mongo *mongo.Database
	Redis *redis.Client
}

var Connections = &Datastore{}

func New() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create mongo client: %v\n", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		log.Fatalf("Failed to initialize mongo client: %v\n", err)
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping mongo database: %v\n", err)
	}
	Connections.Mongo = client.Database("socialist")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	Connections.Redis = rdb
}
