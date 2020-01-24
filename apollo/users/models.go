package users

import "time"

const table = "users"

type User struct {
	Id        string
	Email     string
	Password  string
	CreatedAt time.Time
}

// LoginPayload Represents the paylod for POST /login
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupPayload Represents the paylod for POST /signup
type SignupPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//import (
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//const collectionName = "users"
//
//type User struct {
//	Id        primitive.ObjectID `json:"id" bson:"_id"`
//	UserName  string             `json:"user_name" bson:"user_name"`
//	Email     string             `json:"email" bson:"email"`
//	Password  string             `json:"password" bson:"password"`
//	FirstName string             `json:"first_name" bson:"first_name"`
//	LastName  string             `json:"last_name" bson:"last_name"`
//	Address   string             `json:"address" bson:"address"`
//}

//func FindUserByLoginPayload(payload LoginPayload) (User, error) {
//	if len(payload.UserName) > 0 {
//		return FindUserByUserName(payload.UserName)
//	} else if len(payload.Email) > 0 {
//		return FindUserByEmail(payload.Email)
//	}
//
//	return User{}, errors.New("user_name and email are empty")
//}
//
//func FindUserById(hexId string) (User, error) {
//	var (
//		err  error
//		user User
//	)
//
//	id, err := primitive.ObjectIDFromHex(hexId)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	filter := bson.M{"_id": id}
//	collection := database.GetDB(env.MongoDb).Collection(collectionName)
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//
//	return user, collection.FindOne(ctx, filter).Decode(&user)
//}

//func FindUserByUserName(userName string) (User, error) {
//	var user User
//
//	filter := bson.M{"user_name": userName}
//	collection := database.GetDB(env.MongoDb).Collection(collectionName)
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//
//	return user, collection.FindOne(ctx, filter).Decode(&user)
//}
//
//func FindUserByEmail(email string) (User, error) {
//	var user User
//
//	filter := bson.M{"email": email}
//	collection := database.GetDB(env.MongoDb).Collection(collectionName)
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//
//	return user, collection.FindOne(ctx, filter).Decode(&user)
//}
//
//func InsertUserPayload(payload UserPayload) (*mongo.InsertOneResult, error) {
//	collection := database.GetDB(env.MongoDb).Collection(collectionName)
//
//	bsonBytes, err := bson.Marshal(payload)
//	if err != nil {
//		return nil, err
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(env.MongoRetrySeconds)*time.Second)
//	defer cancel()
//
//	return collection.InsertOne(ctx, bsonBytes)
//}
