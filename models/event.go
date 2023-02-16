package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model for test purposes
type Event struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActivityID     primitive.ObjectID `bson:"activityID" json:"activityID"`
	PropertyValues []PropertyValue    `bson:"propertyValues" json:"propertyValues"`
}

func CreateEvent(activityID primitive.ObjectID, propertyValues []PropertyValue) (*Event, error) {

	event := &Event{
		ID:             primitive.NewObjectID(),
		ActivityID:     activityID,
		PropertyValues: propertyValues,
	}

	// Insert the event into the MongoDB collection
	insertResult, err := EventsCollection.InsertOne(context.TODO(), event)
	if err != nil {
		return nil, err
	}

	event.ID = insertResult.InsertedID.(primitive.ObjectID)
	return event, nil
}

func GetEvent(id primitive.ObjectID) (*Event, error) {
	// Get the event from the MongoDB collection
	var event Event
	err := EventsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEvents retrieves a list of events from the database
func GetEvents() ([]Event, error) {
	var events []Event
	cur, err := EventsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &events); err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateEvent updates a specific event in the database
func UpdateEvent(id primitive.ObjectID, update bson.M) (*Event, error) {
	var event Event
	if err := EventsCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, update).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

// DeleteEvent deletes a specific event from the database
func DeleteEvent(id primitive.ObjectID) error {
	if _, err := EventsCollection.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}
