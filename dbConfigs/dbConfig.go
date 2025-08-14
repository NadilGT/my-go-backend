package dbConfigs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var UserCollection *mongo.Collection

func ConnectMongoDB(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Connect(ctx); err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB")
	MongoClient = client

	db := MongoClient.Database("Evolza")
	fmt.Println(db.Name())

	UserCollection = db.Collection("employees")

	return client
}
