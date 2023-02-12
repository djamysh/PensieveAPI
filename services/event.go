package services

import (
	"encoding/json"
	"net/http"

	"github.com/djamysh/TracerApp/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the event
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = event.CreateEvent()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the event was created successfully
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	var events []*models.Event
	var err error
	if events, err = models.GetEvents(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the properties as a response
	json.NewEncoder(w).Encode(events)
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// How to check whether YCM working
	//var event *models.Event
	var event interface{} // for testing purposes
	event, err = models.MockGetEvent(id)
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
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Delete the activity from the MongoDB collection
	err = models.DeleteEvent(id)
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
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Parse the request body to get the updated property
	var event models.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = event.UpdateEvent(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the property was updated successfully
	json.NewEncoder(w).Encode(event)
}
