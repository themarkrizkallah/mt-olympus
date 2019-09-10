package users

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"front_end_server/common"
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

func FindUserById(id primitive.ObjectID) (User, error) {
	var user User

	filter := bson.M{"_id": id}
	collection := common.GetMongoDb().Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

	return user, collection.FindOne(ctx, filter).Decode(&user)
}

func FindUserByUserName(userName string) (User, error) {
	var user User

	filter := bson.M{"user_name": userName}
	collection := common.GetMongoDb().Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

	return user, collection.FindOne(ctx, filter).Decode(&user)
}

func FindUserByEmail(email string) (User, error) {
	var user User

	filter := bson.M{"email": email}
	collection := common.GetMongoDb().Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

	return user, collection.FindOne(ctx, filter).Decode(&user)
}
