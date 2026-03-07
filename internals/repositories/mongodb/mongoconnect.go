package mongodb

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient(ctx context.Context) (*mongo.Client, error) {
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("username:password@mongodb://localhost:27017"))	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))	
	if err != nil {
		log.Println("Error connecting to MongoDB.", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}