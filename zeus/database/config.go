package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexConfig struct {
	collection string
	indexModel []mongo.IndexModel
}

type IndicesConfig struct {
	userIndexConfig IndexConfig
}

func DefaultIndexConfig() IndicesConfig {
	userIndexCfg := IndexConfig{
		collection: "users",
		indexModel: []mongo.IndexModel{
			{Keys: bson.M{"user_name": 1}, Options: options.Index().SetUnique(true)},
			{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
		},
	}

	return IndicesConfig{userIndexConfig: userIndexCfg}
}
