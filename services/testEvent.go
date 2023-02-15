package services

import (
	"encoding/json"
	"net/http"

	"github.com/djamysh/TracerApp/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateEventRequest struct {
	ActivityID     string                 `json:"activityID"`
	PropertyValues map[string]interface{} `json:"propertyValues"`
}

func MockCreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request CreateEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// Handle error
	}

	// Convert the activityID string to an ObjectID
	activityID, err := primitive.ObjectIDFromHex(request.ActivityID)
	if err != nil {
		// Handle error
	}

	// Convert the propertyValues map keys from strings to ObjectIDs
	propertyValues := make(map[primitive.ObjectID]interface{})
	for key, value := range request.PropertyValues {
		id, err := primitive.ObjectIDFromHex(key)
		if err != nil {
			// Handle error
		}
		propertyValues[id] = value
	}

	// Call the CreateEvent function
	event, err := models.MockCreateEvent(activityID, propertyValues)
	if err != nil {
		// Handle error
	}
	json.NewEncoder(w).Encode(event)
}

func MockGetEventHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var event *models.TestEvent
	event, err = models.MockGetEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the property as a response
	json.NewEncoder(w).Encode(event)
}

// GetEventsHandler retrieves a list of events and returns them as a response
func MockGetEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := models.MockGetEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// UpdateEventHandler updates a specific event based on the passed ID and returns the updated event as a response
func MockUpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Event ID", http.StatusBadRequest)
		return
	}

	var updateEvent CreateEventRequest
	if err = json.NewDecoder(r.Body).Decode(&updateEvent); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// [*] I am not sure this convertion wheter necessary or not.

	// Convert the propertyValues map keys from strings to ObjectIDs
	propertyValues := make(map[primitive.ObjectID]interface{})
	for key, value := range updateEvent.PropertyValues {
		id, err := primitive.ObjectIDFromHex(key)
		if err != nil {
			// Handle error
		}
		propertyValues[id] = value
	}

	// Convert map[primitiveObjectID]interface{} to []PropertyValues
	propertyValuesSlice := make([]models.PropertyValue, 0, len(propertyValues))
	for key, value := range propertyValues {
		propertyValuesSlice = append(propertyValuesSlice, models.PropertyValue{
			Key:   key,
			Value: value,
		})
	}

	// TODO
	// if the activity ID is not given
	// -> fetch the previous activity and apply the given property values
	// -> into the control event function
	// else
	// -> apply the given property values into the control event function

	// currently it is going to set the given values without check.

	update := bson.M{"$set": bson.M{"activityID": updateEvent.ActivityID, "propertyValues": propertyValuesSlice}}

	event, err := models.MockUpdateEvent(id, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(event)
}

func MockDeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the event to be deleted from the URL path
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteEvent function to delete the event from the database
	err = models.MockDeleteEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
