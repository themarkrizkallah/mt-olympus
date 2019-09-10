package users

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"front_end_server/common"
	"front_end_server/env"
)

const collectionName = "users"

type User struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserName  string             `bson:"user_name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Address   string             `bson:"address"`
}

func FindUserByLoginPayload(payload LoginPayload) (User, error) {
	if len(payload.UserName) > 0 {
		return FindUserByUserName(payload.UserName)
	} else if len(payload.Email) > 0 {
		return FindUserByEmail(payload.Email)
	}

	return User{}, errors.New("user_name and email are empty")
}

func FindUserById(hexId string) (User, error) {
	var (
		err  error
		user User
	)

	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		log.Fatalln(err)
	}

	filter := bson.M{"_id": id}
	collection := common.GetMongoDb().Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return user, collection.FindOne(ctx, filter).Decode(&user)
}

func FindUserByUserName(userName string) (User, error) {
	var user User

	filter := bson.M{"user_name": userName}
	collection := common.GetMongoDb().Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return user, collection.FindOne(ctx, filter).Decode(&user)
}

func FindUserByEmail(email string) (User, error) {
	var user User

	filter := bson.M{"email": email}
	collection := common.GetMongoDb().Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return user, collection.FindOne(ctx, filter).Decode(&user)
}

func InsertUserPayload(payload UserPayload) (*mongo.InsertOneResult, error) {
	collection := common.GetMongoDb().Collection(collectionName)

	bsonBytes, err := bson.Marshal(payload)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, bsonBytes)
}
