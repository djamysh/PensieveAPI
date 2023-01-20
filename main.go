package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var DBName = "TracerApp"

type Activity struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//ActivityID        int                  `bson:"activityID" json:"activityID"`
	Name              string               `bson:"name" json:"name"`
	Description       string               `bson:"description" json:"description"`
	DefinedProperties []primitive.ObjectID `bson:"definedProperties" json:"definedProperties"`
}

type Property struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	ValueDataType string             `bson:"valueDataType" json:"valueDataType"`
}

type Event struct {
	ID             primitive.ObjectID                 `bson:"_id,omitempty" json:"id"`
	ActivityID     primitive.ObjectID                 `bson:"activityID" json:"activityID"`
	PropertyValues map[primitive.ObjectID]interface{} `bson:"propertyValues" json:"propertyValues"`
	Timestamp      int64                              `bson:"timestamp" json:"timestamp"`
}

func main() {
	// Connect to MongoDB
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Define the router
	r := mux.NewRouter()

	// Add the routes
	r.HandleFunc("/activities", CreateActivityHandler).Methods("POST")
	r.HandleFunc("/activities/{id}", UpdateActivityHandler).Methods("PUT")
	r.HandleFunc("/activities/{id}", DeleteActivityHandler).Methods("DELETE")
	r.HandleFunc("/activities", GetActivitiesHandler).Methods("GET")
	r.HandleFunc("/activities/{id}", GetActivityHandler).Methods("GET")
	r.HandleFunc("/properties", CreatePropertyHandler).Methods("POST")
	r.HandleFunc("/properties/{id}", UpdatePropertyHandler).Methods("PUT")
	r.HandleFunc("/properties/{id}", DeletePropertyHandler).Methods("DELETE")
	r.HandleFunc("/properties", GetPropertiesHandler).Methods("GET")
	r.HandleFunc("/properties/{id}", GetPropertyHandler).Methods("GET")
	r.HandleFunc("/events", CreateEventHandler).Methods("POST")
	r.HandleFunc("/events/{id}", GetEventHandler).Methods("GET")
	r.HandleFunc("/events", GetEventsHandler).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", r))
}

func CreateActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the activity
	var activity Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert the activity into the MongoDB collection
	activity.ID = primitive.NewObjectID()
	_, err = client.Database(DBName).Collection("activities").InsertOne(context.TODO(), activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the activity was created successfully
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}

func UpdateActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Get the activity ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Parse the request body to get the updated activity
	var activity Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update the activity in the MongoDB collection
	activity.ID = id
	_, err = client.Database(DBName).Collection("activities").ReplaceOne(context.TODO(), bson.M{"_id": id}, activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the activity was updated successfully
	json.NewEncoder(w).Encode(activity)
}

func DeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Get the activity ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Delete the activity from the MongoDB collection
	_, err := client.Database(DBName).Collection("activities").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the activity was deleted successfully
	w.WriteHeader(http.StatusNoContent)
}

func GetActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Get the activity ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Get the activity from the MongoDB collection
	var activity Activity
	err := client.Database(DBName).Collection("activities").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the activity as a response
	json.NewEncoder(w).Encode(activity)
}

func GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	// Get all the activities from the MongoDB collection
	cursor, err := client.Database(DBName).Collection("activities").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var activities []Activity
	if err = cursor.All(context.TODO(), &activities); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the activities as a response
	json.NewEncoder(w).Encode(activities)
}

func CreatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the property
	var property Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert the property into the MongoDB collection
	property.ID = primitive.NewObjectID()
	_, err = client.Database(DBName).Collection("properties").InsertOne(context.TODO(), property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the property was created successfully
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(property)
}

func UpdatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Parse the request body to get the updated property
	var property Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update the property in the MongoDB collection
	property.ID = id
	_, err = client.Database(DBName).Collection("properties").ReplaceOne(context.TODO(), bson.M{"_id": id}, property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the property was updated successfully
	json.NewEncoder(w).Encode(property)
}

func DeletePropertyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Delete the property from the MongoDB collection
	_, err := client.Database(DBName).Collection("properties").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the property was deleted successfully
	w.WriteHeader(http.StatusNoContent)
}

func GetPropertyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Get the property from the MongoDB collection
	var property Property
	err := client.Database(DBName).Collection("properties").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the property as a response
	json.NewEncoder(w).Encode(property)
}

func GetPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	// Get all the properties from the MongoDB collection
	cursor, err := client.Database(DBName).Collection("properties").Find(context.TODO(), bson.M{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var properties []Property
	if err = cursor.All(context.TODO(), &properties); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the properties as a response
	json.NewEncoder(w).Encode(properties)
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the event
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert the event into the MongoDB collection
	event.ID = primitive.NewObjectID()
	event.Timestamp = int64(time.Now().Unix())
	_, err = client.Database(DBName).Collection("events").InsertOne(context.TODO(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the event was created successfully
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all the properties from the MongoDB collection
	cursor, err := client.Database(DBName).Collection("events").Find(context.TODO(), bson.M{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var events []Event
	if err = cursor.All(context.TODO(), &events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the properties as a response
	json.NewEncoder(w).Encode(events)
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Get the property from the MongoDB collection
	var event Event
	err := client.Database(DBName).Collection("events").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the property as a response
	json.NewEncoder(w).Encode(event)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// Get the activity ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Delete the activity from the MongoDB collection
	_, err := client.Database(DBName).Collection("events").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the activity was deleted successfully
	w.WriteHeader(http.StatusNoContent)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	// Parse the request body to get the updated property
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update the property in the MongoDB collection
	event.ID = id
	_, err = client.Database(DBName).Collection("properties").ReplaceOne(context.TODO(), bson.M{"_id": id}, event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the property was updated successfully
	json.NewEncoder(w).Encode(event)
}
