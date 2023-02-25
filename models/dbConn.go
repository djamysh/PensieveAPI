package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func CreateUniqueFieldInCollection(collection *mongo.Collection, field string, order int) {
	// order-> -1:desending 1:ascending order

	// Create a unique index on the given field of the collection
	_, err := collection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bson.M{field: order},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	Connect2DB()
	ActivitiesCollection = Client.Database(DBName).Collection(ActivitiesCollectionName)
	EventsCollection = Client.Database(DBName).Collection(EventsCollectionName)
	PropertiesCollection = Client.Database(DBName).Collection(PropertiesCollectionName)

	// Create a unique index on the 'name' field of the PropertiesCollection collection
	CreateUniqueFieldInCollection(PropertiesCollection, "name", 1)
	// Create a unique index on the 'name' field of the ActivitiesCollection collection
	CreateUniqueFieldInCollection(ActivitiesCollection, "name", 1)

}
