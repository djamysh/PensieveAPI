package models

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/djamysh/TracerApp/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID             primitive.ObjectID                 `bson:"_id,omitempty" json:"id"`
	ActivityID     primitive.ObjectID                 `bson:"activityID" json:"activityID"`
	PropertyValues map[primitive.ObjectID]interface{} `bson:"propertyValues" json:"propertyValues"`
	//Timestamp      int64                              `bson:"timestamp" json:"timestamp"`
}

var TypeMap = map[string]reflect.Type{
	"string":        reflect.TypeOf(""),
	"integer":       reflect.TypeOf(int64(0)),
	"float":         reflect.TypeOf(float64(0)),
	"timelings":     reflect.TypeOf(map[string]int64{}),
	"integer array": reflect.TypeOf([]int64{}),
	"string array":  reflect.TypeOf([]string{}),
	"float array":   reflect.TypeOf([]float64{}),
}

var TypeErr = &utils.EventError{Message: "Invalid data type"}
var TimestampErr = &utils.EventError{Message: "Invalid UNIX timestamp"}

// TODO: Improve this function,
func ControlEvent(event *Event) error {
	// default types : integer, float, string, timelings or array of those

	activity, err := GetActivity(event.ActivityID)

	if err != nil {
		return err
	}

	for _, propertyID := range activity.DefinedProperties {
		property, err := GetProperty(propertyID)
		if err != nil {
			return err
		}

		propertyValue := event.PropertyValues[propertyID]
		valueType := reflect.TypeOf(propertyValue)
		if TypeMap[property.ValueDataType].Name() == valueType.Name() {

			if valueType == TypeMap["timelings"] {

				// Checling wheter the given timeling is valid or not
				for _, timeling := range propertyValue.(map[string]int64) {
					if time.Unix(timeling, 0).IsZero() {
						return TimestampErr // given timestamp is not valid
					}
				}
			}

			continue

		} else {
			msg := fmt.Sprintf("PropertyID : %s, Given Property Value Type : %s Expected Property Value Type : %s ", propertyID, reflect.TypeOf(propertyValue), TypeMap[property.ValueDataType].Name())
			return errors.New(msg)

			// return TypeErr
		}
	}
	return nil
}

func (event *Event) CreateEvent() error {
	// Insert the event into the MongoDB collection
	event.ID = primitive.NewObjectID()
	//event.PropertyValues[DefaultTimelingsPropertyID] = map[string]int64{"CreatedAt": int64(time.Now().Unix())}

	/*
		createdAt := map[string]int64{"CreatedAt": int64(time.Now().Unix())}
		primitiveCreatedAt := primitive.M{}
		for key, value := range createdAt {
			primitiveCreatedAt[key] = value
		}
		event.PropertyValues[DefaultTimelingsPropertyID] = primitiveCreatedAt
	*/
	err := ControlEvent(event)
	if err != nil {
		return err
	}

	// defaultTimelings := map[string]int64{"CreatedAt": int64(time.Now().Unix())}
	// event.PropertyValues[DefaultTimelingsPropertyID] = utils.MarshallTimelings(defaultTimelings)

	_, err = EventsCollection.InsertOne(context.TODO(), event)

	return err
}

func (event *Event) UpdateEvent(id primitive.ObjectID) error {

	// Update the event in the MongoDB collection
	event.ID = id

	err := ControlEvent(event)
	if err != nil {
		return err

	}

	_, err = EventsCollection.ReplaceOne(context.TODO(), bson.M{"_id": id}, event)

	return err
}

func DeleteEvent(id primitive.ObjectID) error {

	// Delete the event from the MongoDB collection
	_, err := EventsCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err

}

func GetEvent(id primitive.ObjectID) (*Event, error) {
	// Get the event from the MongoDB collection
	var event Event

	err := EventsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&event)
	return &event, err
}

func MockGetEvent(id primitive.ObjectID) (*interface{}, error) {
	// Get the event from the MongoDB collection
	//var event Event
	var event interface{}

	val := EventsCollection.FindOne(context.TODO(), bson.M{"_id": id})
	err := val.Decode(&event)
	//err := EventsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&event)
	return &event, err
}

func GetEvents() ([]*Event, error) {
	// Get all the activities from the MongoDB collection
	cursor, err := EventsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []*Event
	if err = cursor.All(context.TODO(), &events); err != nil {
		return nil, err
	}
	return events, nil
}
