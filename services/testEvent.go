package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/djamysh/TracerApp/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateEventRequest struct {
	ActivityID     string                 `json:"activityID"`
	PropertyValues map[string]interface{} `json:"propertyValues"`
}

var TypeNullMap = map[string]interface{}{
	"string":        "",
	"integer":       int64(0),
	"float":         float64(0),
	"timelings":     map[string]int64{},
	"integer array": []int64{},
	"string array":  []string{},
	"float array":   []float64{},
}

// TODO: Implement a cache for inner calls to increase performance, use redis.

// TODO: Improve this function,
// TODO: Find a better name for this function, it does not give accurate info
// because ControlEvent function also processes the given data into to a
// useable format, converts the PropertyValues into the []PropertyValue slice.
func MockControlEvent(event *CreateEventRequest, previousEvent *models.TestEvent) (*models.TestEvent, error) {
	// default types : integer, float, string, timelings or array of those

	// Convert the activityID string to an ObjectID
	activityID, err := primitive.ObjectIDFromHex(event.ActivityID)
	if err != nil {
		// Handle error
	}

	// Get the corresponding activity entry
	activity, err := models.GetActivity(activityID)
	if err != nil {
		return nil, err
	}

	undefinedProperties := make(map[primitive.ObjectID]*models.Property)

	// Checking data type consistency with given property values' data types
	// TODO: make neater way of error response
	for _, propertyID := range activity.DefinedProperties {
		// Get the property entry
		property, err := models.GetProperty(propertyID)
		if err != nil {
			return nil, err
		}

		// Get the corresponding value

		propertyValue, isPresent := event.PropertyValues[propertyID.Hex()]

		// If the propertyValue is not given
		if !isPresent {
			undefinedProperties[propertyID] = property

		} else {

			// Determine the data type
			valueType := reflect.TypeOf(propertyValue)

			// If the given value's data type is valid
			if models.TypeMap[property.ValueDataType].Name() == valueType.Name() {

				// if the property is timelings
				if valueType == models.TypeMap["timelings"] {

					// Checking wheter the given timeling is valid or not
					for _, timeling := range propertyValue.(map[string]int64) {
						if time.Unix(timeling, 0).IsZero() {
							return nil, models.TimestampErr // given timestamp is not valid
						}
					}
				}

				continue

			} else {
				msg := fmt.Sprintf("PropertyID : %s, Given Property Value Type : %s Expected Property Value Type : %s ", propertyID, reflect.TypeOf(propertyValue), models.TypeMap[property.ValueDataType].Name())
				return nil, errors.New(msg)

				// return TypeErr
			}

		}

	}

	// If the given property values are valid according to given Properties.

	// Convert the propertyValues map keys from strings to ObjectIDs
	propertyValues := make(map[primitive.ObjectID]interface{})
	for key, value := range event.PropertyValues {
		id, err := primitive.ObjectIDFromHex(key)
		if err != nil {
			return nil, err
		}
		propertyValues[id] = value
	}

	// Convert map[objectID]string to []PropertyValues
	propertyValuesSlice := make([]models.PropertyValue, 0, len(propertyValues))
	for key, value := range propertyValues {
		propertyValuesSlice = append(propertyValuesSlice, models.PropertyValue{
			Key:   key,
			Value: value,
		})
	}

	//TODO: Consider checking validity of the previousEvent parameter, in case not equal to nil.

	// Append the undefined properties into the propertyValueSlice with accurate values
	if previousEvent == nil {
		// Create Event call
		for propertyID, property := range undefinedProperties {

			// Setting the null element
			element := models.PropertyValue{
				Key:   propertyID,
				Value: TypeNullMap[property.ValueDataType],
			}

			// Appending to the propertyValuesSlice
			propertyValuesSlice = append(propertyValuesSlice, element)

		}
	} else {
		// Update Event call
		for _, pair := range previousEvent.PropertyValues {

			// Checking wheter the propertyID is defined in undefinedProps map
			_, isExist := undefinedProperties[pair.Key]

			// if it is append the pair to propertyValuesSlice
			if isExist {
				propertyValuesSlice = append(propertyValuesSlice, pair)
			}
		}

	}

	/*
		// Append the undefined properties into the propertyValueSlice
		for propertyID, property := range undefinedProperties {
			// If this propertyValue is previously defined or
			// in other words this is a update process
			// Fetch the previous value of the property and
			// set to undefined property value, else
			// set null value accordingly.

			if previousEvent == nil {
				// Create Event call

				// Setting the null value
				// [!] I don't know maybe setting directly nil is a better option.
				element := models.PropertyValue{
					Key:   propertyID,
					Value: TypeNullMap[property.ValueDataType],
				}
				propertyValuesSlice = append(propertyValuesSlice, element)

			} else {
				// Update Event call
				// Finding the property and setting the value
				for _, pair := range previousEvent.PropertyValues {
					if pair.Key == propertyID {
						propertyValuesSlice = append(propertyValuesSlice, pair)
					}

				}

			}
		}
	*/

	// Pass the processed data into the new model.
	var checkedEvent models.TestEvent
	checkedEvent.ActivityID = activityID
	checkedEvent.PropertyValues = propertyValuesSlice

	return &checkedEvent, nil
}

func MockCreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request CreateEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	event, err := MockControlEvent(&request, nil)
	if err != nil {

		// TODO: make better way of error handling
		http.Error(w, err.Error(), http.StatusBadRequest)
		//Handle error
	}

	// Call the CreateEvent function
	event, err = models.MockCreateEvent(event.ActivityID, event.PropertyValues)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusBadRequest)
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

// UpdateEventHandler updates a specific event based on
// the passed ID and returns the old event as a response
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

	var event, previousEvent *models.TestEvent

	previousEvent, err = models.MockGetEvent(id)
	if err != nil {
		// If there is a problem with obtaining the previous
		// Event value from the given EventID
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	event, err = MockControlEvent(&updateEvent, previousEvent)
	if err != nil {
		// If there is a problem with controlling the event
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//TODO: Consider AppendUpdate for array data types,

	update := bson.M{"$set": bson.M{"activityID": event.ActivityID, "propertyValues": event.PropertyValues}}

	oldEvent, err := models.MockUpdateEvent(id, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(oldEvent)
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
