package services

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/djamysh/TracerApp/models"
)

// route handler functions for the models.Activity model
func CreateActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the activity
	var activity models.Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = activity.CreateActivity()
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
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the request body to get the updated activity
	var activity models.Activity

	err = json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Overwrite in case of id value change attempt
	activity.ID = id

	err = activity.UpdateActivity(id)

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
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete the activity from the MongoDB collection
	err = models.DeleteActivity(id)

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

	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var activity *models.Activity
	activity, err = models.GetActivity(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the activity as a response
	json.NewEncoder(w).Encode(activity)
}

func GetActivityByNameHandler(w http.ResponseWriter, r *http.Request) {
	// Get the activity ID from the URL
	vars := mux.Vars(r)

	name := vars["name"]

	var activity *models.Activity
	activity, err := models.GetActivityByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the activity as a response
	json.NewEncoder(w).Encode(activity)
}

func GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {

	var activities []*models.Activity
	var err error
	if activities, err = models.GetActivities(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the activities as a response
	json.NewEncoder(w).Encode(activities)
}
