package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DBName = "PensieveAPI"
var ActivitiesCollectionName = "activities"
var ActivitiesCollection *mongo.Collection
var EventsCollectionName = "events"
var EventsCollection *mongo.Collection
var PropertiesCollectionName = "properties"
var PropertiesCollection *mongo.Collection

func Connect2DB() {
	// Connect to MongoDB
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Estabilishing the connection
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		// Case of serious problem
		log.Fatal(err)
	}

	log.Println("Database connection has been estabilished.")

}

func init() {
	Connect2DB()
	ActivitiesCollection = Client.Database(DBName).Collection(ActivitiesCollectionName)
	EventsCollection = Client.Database(DBName).Collection(EventsCollectionName)
	PropertiesCollection = Client.Database(DBName).Collection(PropertiesCollectionName)

}
