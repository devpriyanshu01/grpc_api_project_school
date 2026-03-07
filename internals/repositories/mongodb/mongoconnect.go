package mongodb

import (
	"context"
	"log"

	"github.com/devpriyanshu01/grpc_api_project_school/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient(ctx context.Context) (*mongo.Client, error) {
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("username:password@mongodb://localhost:27017"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		utils.ErrorHandler(err, "Failed to connect to mongodb database.")
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		utils.ErrorHandler(err, "Error while pinging mongodb.")
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}
