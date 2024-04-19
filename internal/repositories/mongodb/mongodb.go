package mongodb

import (
	"context"
	"fmt"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client = Connect()

func Connect() *mongo.Client {
	mongoURI := config.GetMongoDBConfig().URI
	if mongoURI == "" {
		log.Fatal("MongoDB Host is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	if err := pingMongoDB(client); err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	if err := createUserIndexes(client); err != nil {
		log.Fatalf("MongoDB create user indexes error: %v", err)
	}

	log.Println("Connected to MongoDB")
	return client
}
func pingMongoDB(client *mongo.Client) error {
	ctx, cancel := helpers.ContextWithTimeout(20)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping error: %v", err)
	}
	return nil
}

func createUserIndexes(client *mongo.Client) error {
	collection := client.Database(config.GetMongoDBConfig().Database).Collection(config.GetMongoDBConfig().Collections.Users)
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	return err
}

// GetCollection returns a collection
func GetCollection(collectionName string) *mongo.Collection {
	return client.Database(config.GetMongoDBConfig().Database).Collection(collectionName)
}
