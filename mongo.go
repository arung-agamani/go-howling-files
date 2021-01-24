package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

// UserAuthObject Data type representing UserAuth schema from mongoDb
type UserAuthObject struct {
	// ID          primitive.ObjectID `bson:"_id"`
	UserID      string    `json:"userId"`
	AccessToken string    `json:"accessToken"`
	LastAccess  time.Time `json:"lastAccess"`
}

func getMongoClient() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionURI := os.Getenv("MONGO_URL")
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(connectionURI)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
		log.Println("MongoDB Connected")
	})
	return clientInstance, clientInstanceError
}

func getUserAuth(token string) (UserAuthObject, error) {
	result := UserAuthObject{}
	client, err := getMongoClient()
	if err != nil {
		return result, err
	}
	filter := bson.D{
		primitive.E{
			Key:   "accessToken",
			Value: token,
		},
	}
	collection := client.Database("blogDatabase").Collection("userauths")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
