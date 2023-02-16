package services

import (
	"encoding/json"
	"net/http"

	"github.com/djamysh/TracerApp/models"
	"github.com/djamysh/TracerApp/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the property
	var property models.Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert the property into the MongoDB collection
	property.ID = primitive.NewObjectID()

	property.ValueDataType = utils.CleanInput(property.ValueDataType)

	if !property.IsValidType() {

		// When the given input is invalid
		http.Error(w, "Invalid data type.", http.StatusBadRequest)
		return

	}

	err = property.CreateProperty()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating that the property was created successfully
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(property)
}

func UpdatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	// ?? Is it possible to change ID with the request
	// ?> No, because there is a redefinition of property.ID in the
	// following lines. Even if there is a ID value in the request
	// it will be overwritten with the given parameter ID.

	// Get the property ID from the URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	// Parse the request body to get the updated property
	var property models.Property
	err = json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update the property in the MongoDB collection
	property.ID = id

	property.ValueDataType = utils.CleanInput(property.ValueDataType)
	if !property.IsValidType() {
		// If not a valid property value type

		// When the given input is invalid
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	err = property.UpdateProperty(id)
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
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.DeleteProperty(id)
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
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	property, err := models.GetProperty(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the property as a response
	json.NewEncoder(w).Encode(property)
}

func GetPropertyByNameHandler(w http.ResponseWriter, r *http.Request) {
	// Get the property ID from the URL
	vars := mux.Vars(r)
	name := vars["name"]

	property, err := models.GetPropertyByName(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the property as a response
	json.NewEncoder(w).Encode(property)
}

func GetPropertiesHandler(w http.ResponseWriter, r *http.Request) {

	var properties []*models.Property
	var err error
	if properties, err = models.GetProperties(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the properties as a response
	json.NewEncoder(w).Encode(properties)
}
