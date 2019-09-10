package common

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"front_end_server/env"
)

var (
	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
)

func InitMongo() *mongo.Client {
	var err error

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(env.MongoUri))
	if err != nil {
		log.Panicln("Error setting up Mongo client:", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Panicln("Error connecting to Mongo client:", err)
	}

	mongoDatabase = mongoClient.Database(env.MongoDb)
	setupMongoIndices()

	return mongoClient
}

func GetMongoDb() *mongo.Database {
	return mongoDatabase
}

func setupMongoIndices() {
	const timeoutSeconds = 10

	model := mongo.IndexModel{
		Keys: bson.M{
			"user_name": 1,
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	collection := mongoDatabase.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), timeoutSeconds * time.Second)
	opts := options.CreateIndexes().SetMaxTime(timeoutSeconds * time.Second)

	_, err := collection.Indexes().CreateOne(ctx, model, opts)
	if err != nil {
		log.Panicln("Failed setting up indices", err)
	}
}