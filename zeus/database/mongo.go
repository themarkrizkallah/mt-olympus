package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"zeus/env"
)

var client *mongo.Client

func Init() (*mongo.Client, error) {
	var err error

	uri := fmt.Sprintf(
		"mongodb://%v:%v@%v:%v",
		env.MongoUser,
		env.MongoPassword,
		env.MongoHostName,
		env.MongoPort,
	)

	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	for i := uint64(0); i < env.MongoRetryTimes; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)

		if err = client.Connect(ctx); err == nil {
			ctx, _ = context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
			if err = client.Ping(ctx, readpref.Primary()); err != nil {
				return nil, err
			}

			return client, nil
		}
	}

	return nil, err
}

func SetupIndices(cfg IndicesConfig, dbName string) []error {
	var errs []error

	db := client.Database(dbName)

	for _, indexConfig := range []IndexConfig{cfg.userIndexConfig} {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
		opts := options.CreateIndexes().SetMaxTime(time.Duration(env.MongoRetrySeconds) * time.Second)

		collection := db.Collection(indexConfig.collection)
		if _, err := collection.Indexes().CreateMany(ctx, indexConfig.indexModel, opts); err != nil {
			errs = append(errs, err)
		}

	}

	return errs
}

func GetDB(dbName string) *mongo.Database {
	return client.Database(dbName)
}
