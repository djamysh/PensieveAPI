package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var DBName = "TracerApp"
var validDataTypes = [5]string{"string", "percentage", "integer", "float", "timelings"}

var defaultTimelingsName = "Default Timelings"
var defaultTimelingsPropertyID primitive.ObjectID

// expected formats
// string -> string
// percentage -> x element of [0,100]
// integer -> x element of Z
// float -> x element of R
// timelings -> map of timestamps with tags

// each activity must have the property timelings. As a default if the activity has
// a start and end time specific then it must have two key-pair in the timelings map,
// which are {'start':timestamp of start, 'end': timestamp of end}, if the activity data
// does not rely on start and end, then it must have 'instant' timestamp, which corresponds
// to the moment which data submitted. For example recording the cigarette consumption
// can be instantenous timeling while recording the thermodynamics study duration being
// an interval of start and end. This is for the default timelings property. Each activity
// will have either 'IntervalTimelings' property or 'InstantTimeling' property.
// After a brainstorm I think defining as Instant or Interval is restricting the flexability.
// *I think every Activity must have Default timeling property but not restricted as
// instant or interval.
// ** After a bit of coding I decided that Interval Timelings is a bit of optional, user additionally
// can also define it, I will automatically define the default timeling while creating the event in the collection.

// each activity must have the default Note property. I am not sure, it may be redundant.

// Get (Property&Activity) by name functions must be implemented

// I am thinking of removing percentage datatype, because if you consider
// all of the other datatypes they are well defined datatypes(timelings as map[string]int64)
// however percentage requires additional constraint of x element of [0,100]. It is is similar
// to that start-end timelings is some constraint of timelings. I will remove the percentage
// datatype, in future I am planning to add some built-in constraints for percentage, start-end
// intervals or similar usefull built-in features.

type Activity struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//ActivityID        int                  `bson:"activityID" json:"activityID"`
	Name              string               `bson:"name" json:"name" validate:"unique"`
	Description       string               `bson:"description" json:"description"`
	DefinedProperties []primitive.ObjectID `bson:"definedProperties" json:"definedProperties"`
}

type Property struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name" validate:"unique"`
	Description   string             `bson:"description" json:"description"`
	ValueDataType string             `bson:"valueDataType" json:"valueDataType"`
}

type Event struct {
	ID             primitive.ObjectID                 `bson:"_id,omitempty" json:"id"`
	ActivityID     primitive.ObjectID                 `bson:"activityID" json:"activityID"`
	PropertyValues map[primitive.ObjectID]interface{} `bson:"propertyValues" json:"propertyValues"`
	//Timestamp      int64                              `bson:"timestamp" json:"timestamp"`
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
	// Check if default timelings properties are created
	defaultTimelingsPropertyID = DefaultTimelingsProperty()

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

func DefaultTimelingsProperty() primitive.ObjectID {
	// Check wheter the default timelings Property is created or not, if not  creates.
	var defaultTimelingsProperty *Property

	if err := client.Database(DBName).Collection("properties").FindOne(context.TODO(), bson.M{"name": defaultTimelingsName}).Decode(&defaultTimelingsProperty); err != nil {
		// definition of default timelings property
		defaultTimelingsProperty = &Property{ID: primitive.NewObjectID(), Name: "Default Timelings", Description: "Default timelings property", ValueDataType: "timelings"}
		// inserting to the properties collection
		client.Database(DBName).Collection("properties").InsertOne(context.TODO(), defaultTimelingsProperty)
		return defaultTimelingsProperty.ID
	}

	return defaultTimelingsProperty.ID
}

func (P *Property) isValidType() bool {
	// Check wheter the given data type is valid
	for _, t := range validDataTypes {
		if P.ValueDataType == t {
			return true
		}
	}
	return false
}

func cleanInput(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
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

	// Appending the default Timelings property
	activity.DefinedProperties = append(activity.DefinedProperties, defaultTimelingsPropertyID)

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

	property.ValueDataType = cleanInput(property.ValueDataType)
	if !property.isValidType() {

		// When the given input is invalid
		http.Error(w, "Invalid data type.", http.StatusBadRequest)
		return

	}

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

	property.ValueDataType = cleanInput(property.ValueDataType)
	if property.isValidType() {

		// When the given input is invalid
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return

	}
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
	event.PropertyValues[defaultTimelingsPropertyID] = int64(time.Now().Unix())
	//event.Timestamp = int64(time.Now().Unix())
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
