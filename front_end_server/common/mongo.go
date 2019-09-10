package common

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"front_end_server/env"
)

var (
	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
)

func InitMongo() *mongo.Client {
	var err error

	log.Println("Connecting to MongoDB...")
	uri := fmt.Sprintf(
		"mongodb://%v:%v@%v:%v",
		env.MongoUser,
		env.MongoPassword,
		env.MongoHostName,
		env.MongoPort,
	)

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicln("Error setting up Mongo client:", err)
	}

	for i := uint64(0); i < env.MongoRetryTimes; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
		err = mongoClient.Connect(ctx)
		if err == nil {
			log.Println("Successfully connected to MongoDB...")
			mongoDatabase = mongoClient.Database(env.MongoDb)

			ctx, _ = context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
			err = mongoClient.Ping(ctx, readpref.Primary())
			if err != nil {
				log.Fatalln("Mongo not responsive")
			}

			setupMongoIndices()

			return mongoClient
		}
	}

	log.Panicln("Error connecting to MongoDB client:", err)

	return nil
}

func GetMongoDb() *mongo.Database {
	return mongoDatabase
}

func setupMongoIndices() {
	var err error

	log.Println("Setting up MongoDB indices...")

	model := mongo.IndexModel{
		Keys: bson.M{
			"user_name": 1,
			"email":     1,
		},
		Options: options.Index().SetUnique(true),
	}

	collection := mongoDatabase.Collection("users")

	for i := uint64(0); i < env.MongoRetryTimes; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
		opts := options.CreateIndexes().SetMaxTime(time.Duration(env.MongoRetrySeconds) * time.Second)

		_, err = collection.Indexes().CreateOne(ctx, model, opts)
		if err == nil {
			log.Println("Successfully set up MongoDB indices...")
			return
		}
	}

	log.Panicln("Failed setting up indices", err)
}
