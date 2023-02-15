package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model for test purposes
type TestEvent struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActivityID     primitive.ObjectID `bson:"activityID" json:"activityID"`
	PropertyValues []PropertyValue    `bson:"propertyValues" json:"propertyValues"`
}

func MockCreateEvent(activityID primitive.ObjectID, propertyValues map[primitive.ObjectID]interface{}) (*TestEvent, error) {
	// Create a new event
	propertyValuesSlice := make([]PropertyValue, 0, len(propertyValues))
	for key, value := range propertyValues {
		propertyValuesSlice = append(propertyValuesSlice, PropertyValue{
			Key:   key,
			Value: value,
		})
	}

	event := &TestEvent{
		ID:             primitive.NewObjectID(),
		ActivityID:     activityID,
		PropertyValues: propertyValuesSlice,
	}

	// Insert the event into the MongoDB collection
	insertResult, err := EventsCollection.InsertOne(context.TODO(), event)
	if err != nil {
		return nil, err
	}

	event.ID = insertResult.InsertedID.(primitive.ObjectID)
	return event, nil
}

func MockGetEvent(id primitive.ObjectID) (*TestEvent, error) {
	// Get the event from the MongoDB collection
	var event TestEvent
	err := EventsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEvents retrieves a list of events from the database
func MockGetEvents() ([]TestEvent, error) {
	var events []TestEvent
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
func MockUpdateEvent(id primitive.ObjectID, update bson.M) (*TestEvent, error) {
	var event TestEvent
	if err := EventsCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, update).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

// DeleteEvent deletes a specific event from the database
func MockDeleteEvent(id primitive.ObjectID) error {
	if _, err := EventsCollection.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}
