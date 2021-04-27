package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoWrapper(dbString string) *mongo.Database {

	clientOptions := options.Client().ApplyURI(dbString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		// Stop bidding if not connected
		panic(err)
	} else {
		fmt.Println("Connected to MongoDB: " + dbString)
	}

	db := client.Database("fin")
	return db
}
