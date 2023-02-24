package services

import (
	"encoding/json"
	"net/http"

	"github.com/djamysh/PensieveAPI/models"
	"github.com/djamysh/PensieveAPI/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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

func GetPropertysRelatedActivities(propertyID primitive.ObjectID) ([]models.Activity, error) {
	filter := bson.M{
		"definedProperties": bson.M{
			"$elemMatch": bson.M{
				"$eq": propertyID,
			},
		},
	}
	// Get the related activities
	relatedActivities, err := models.GetActivitiesByFilter(filter)
	return relatedActivities, err

}

func DeletePropertysRelations(propertyID primitive.ObjectID) error {
	// Get related activites
	relatedActivities, err := GetPropertysRelatedActivities(propertyID)
	if err != nil {
		return err

	}

	for _, relatedActivity := range relatedActivities {
		// Gets the events that are related to related activity
		relatedEvents, err := models.GetEventsByFilter(bson.M{"activityID": relatedActivity.ID})
		if err != nil {
			return err
		}

		for _, relatedEvent := range relatedEvents {
			deletePropertyValueIndex := -1

			// finding the index of the pair that has the updated property in it's Key value
			for idx, pair := range relatedEvent.PropertyValues {
				if pair.Key == propertyID {
					deletePropertyValueIndex = idx
					break
				}
			}

			if deletePropertyValueIndex == -1 {
				// Something is definitely wrong. Because filter must bring the related Events
				// that contain property in their propertyValues
				panic("Error with obtaining the related events properly.")
			}

			// removing the corresponding property value pair
			relatedEvent.PropertyValues = append(relatedEvent.PropertyValues[:deletePropertyValueIndex], relatedEvent.PropertyValues[deletePropertyValueIndex+1:]...)

			// Update the property values of the corresponding related event
			_, err := models.UpdateEvent(relatedEvent.ID, bson.M{"propertyValues": relatedEvent.PropertyValues})
			if err != nil {
				return err
			}

		}

		deletePropertyIDIndex := -1

		// finding the index of the propertyID that has the deleted property.
		for idx, value := range relatedActivity.DefinedProperties {
			if value == propertyID {
				deletePropertyIDIndex = idx
				break
			}
		}

		if deletePropertyIDIndex == -1 {
			// Something is definitely wrong. Because filter must bring the related Activities
			// that contain property in their propertyValues
			panic("Error with obtaining the related activity properly.")
		}

		// deleting the deleted propertyID
		relatedActivity.DefinedProperties = append(relatedActivity.DefinedProperties[:deletePropertyIDIndex], relatedActivity.DefinedProperties[deletePropertyIDIndex+1:]...)

		// update activity
		_, err = models.UpdateActivity(relatedActivity.ID, bson.M{"definedProperties": relatedActivity.DefinedProperties})
		if err != nil {
			return err
		}
	}
	return nil

}

func UpdatePropertysRelations(propertyID primitive.ObjectID, newValueDataType string) error {
	// find the related activites
	// find the related events that are related to related activitiesa
	// delete the property value pairs from the related events' propertyValues that has the deleted property as their Key value
	// delete the property from the related activities defined properties

	// Get the related activities
	relatedActivities, err := GetPropertysRelatedActivities(propertyID)
	if err != nil {
		return err
	}

	for _, relatedActivity := range relatedActivities {
		// Gets the events that are related to related activity
		relatedEvents, err := models.GetEventsByFilter(bson.M{"activityID": relatedActivity.ID})
		if err != nil {
			return err
		}

		for _, relatedEvent := range relatedEvents {
			updatePropertyValueIndex := -1

			// finding the index of the pair that has the updated property in it's Key value
			for idx, pair := range relatedEvent.PropertyValues {
				if pair.Key == propertyID {
					updatePropertyValueIndex = idx
					break
				}
			}

			if updatePropertyValueIndex == -1 {
				// Something is definitely wrong. Because filter must bring the related Events
				// that contain property in their propertyValues
				panic("Error with obtaining the related events properly.")
			}

			// Set the changed property value to null value of the new value data type
			relatedEvent.PropertyValues[updatePropertyValueIndex].Value = TypeNullMap[newValueDataType]

			// Update the property values of the corresponding related event
			_, err := models.UpdateEvent(relatedEvent.ID, bson.M{"propertyValues": relatedEvent.PropertyValues})
			if err != nil {
				return err
			}

		}
	}
	return nil
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

	update := make(map[string]interface{})

	updateRelationsFlag := false

	if property.Name != "" {
		update["name"] = property.Name
	}
	if property.Description != "" {
		update["description"] = property.Description
	}
	if property.ValueDataType != "" {
		// Checking given Value data type
		property.ValueDataType = utils.CleanInput(property.ValueDataType)
		if !property.IsValidType() {
			// If not a valid property value type
			// When the given input is invalid
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		update["valueDataType"] = property.ValueDataType
		updateRelationsFlag = true
	}

	bsonUpdate := bson.M(update)
	oldValue, err := models.UpdateProperty(id, bsonUpdate)

	if oldValue.ValueDataType == update["valueDataType"] {
		// If the new given value data type is the same with previous
		// don't need to update relations flag, it is just
		// overwritting the same value to valueDataType it may
		// lead to incorrect nullification of related event's
		// propertyValues, this is important.

		updateRelationsFlag = false
	}

	//err = property.UpdateProperty(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if updateRelationsFlag {
		if err := UpdatePropertysRelations(id, update["valueDataType"].(string)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	// Send a response indicating that the property was updated successfully
	json.NewEncoder(w).Encode(oldValue)
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

	err = DeletePropertysRelations(id)
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
