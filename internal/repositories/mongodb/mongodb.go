package mongodb

import (
	"context"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client = Connect()

func createUserIndexes(client *mongo.Client) error {
	userIndexModel := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"phone_number": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	if _, err := client.Database(config.GetMongoDBConfig().Database).Collection("users").Indexes().CreateMany(context.Background(),
		userIndexModel); err != nil {
		return err
	}

	return nil
}

func Connect() *mongo.Client {
	if config.GetMongoDBConfig().URI == "" {
		log.Fatal("MongoDB Host is not set")
	}

	clientOptions := options.Client().ApplyURI(config.GetMongoDBConfig().URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	ctx, cancel := utils.ContextWithTimeout(20)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	if err := createUserIndexes(client); err != nil {
		log.Fatalf("MongoDB create user indexes error: %v", err)
	}

	log.Println("Connected to MongoDB")
	return client
}

// GetCollection returns a collection
func GetCollection(collectionName string) *mongo.Collection {
	return client.Database(config.GetMongoDBConfig().Database).Collection(collectionName)
}
