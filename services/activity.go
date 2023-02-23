package services

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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

// Function for finding the difference between the old and new definedProperties
func diffDefinedProperties(oldDefinedProperties, newDefinedProperties []primitive.ObjectID) map[primitive.ObjectID]bool {
	// Create a map of the elements in the old slice along with their frequency of occurrence.
	freq := make(map[primitive.ObjectID]int)
	// ChangeMap : false -> deleted 	true -> added
	changeMap := make(map[primitive.ObjectID]bool)

	for _, elem := range oldDefinedProperties {
		freq[elem]++
	}

	// Iterate over the new slice and find added elements, and decrement the frequency of existing elements.
	for _, elem := range newDefinedProperties {
		if freq[elem] > 0 {
			freq[elem]--
		} else {
			changeMap[elem] = true
			// added = append(added, elem)
		}
	}

	// Iterate over the map of the old slice and find removed elements.
	for elem, count := range freq {
		for i := 0; i < count; i++ {
			changeMap[elem] = false
			//removed = append(removed, elem)
		}
	}

	return changeMap
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

	//TODO: What happens if the user defines the same property twice
	update := make(map[string]interface{})

	updateRelationsFlag := false
	if activity.Name != "" {
		update["name"] = activity.Name

	}
	if activity.Description != "" {
		update["description"] = activity.Description
	}
	if activity.DefinedProperties != nil {
		update["definedProperties"] = activity.DefinedProperties
		updateRelationsFlag = true
	}
	update_bsonM := bson.M(update)

	//err = activity.UpdateActivity(id)
	oldValue, err := models.UpdateActivity(id, update_bsonM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if updateRelationsFlag {

		previousActivity, err := models.GetActivity(id)
		if err != nil {
			// TODO-errHandling
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Finds the changed properties
		changedProperties := diffDefinedProperties(previousActivity.DefinedProperties, activity.DefinedProperties)

		// Gets the events that are related to updated activity
		relatedEvents, err := models.GetEventsByFilter(bson.M{"activityID": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for _, relatedEvent := range relatedEvents {

			// For easier implementation converting back to map[objectID]interface{}
			propertyValues := PropertyValueBackConvertion(relatedEvent.PropertyValues)
			for propertyID, state := range changedProperties {
				if state {
					// added property
					// Get new property
					property, err := models.GetProperty(propertyID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					// Setting the null/default value
					propertyValues[propertyID] = TypeNullMap[property.ValueDataType]

				} else {
					// removed property
					// deleting the property value from propertyValues map
					delete(propertyValues, propertyID)
				}

			}
			// Converting back to DB submitable format
			propertyValuesSlice := PropertyValueConvertion(propertyValues)
			// Updating the new propertyValues
			_, err := models.UpdateEvent(relatedEvent.ID, bson.M{"propertyValues": propertyValuesSlice})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}

	// Send a response indicating that the activity was updated successfully
	json.NewEncoder(w).Encode(oldValue)
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
