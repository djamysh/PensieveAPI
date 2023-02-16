package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//ActivityID        int                  `bson:"activityID" json:"activityID"`
	Name              string               `bson:"name" json:"name" validate:"unique"`
	Description       string               `bson:"description" json:"description"`
	DefinedProperties []primitive.ObjectID `bson:"definedProperties" json:"definedProperties"`
}

func (activity *Activity) CreateActivity() error {

	// Insert the activity into the MongoDB collection
	activity.ID = primitive.NewObjectID()

	// Appending the default Timelings property
	// activity.DefinedProperties = append(activity.DefinedProperties, DefaultTimelingsPropertyID)

	_, err := ActivitiesCollection.InsertOne(context.TODO(), activity)
	return err

}

func (activity *Activity) UpdateActivity(id primitive.ObjectID) error {

	// Update the activity in the MongoDB collection
	activity.ID = id
	_, err := ActivitiesCollection.ReplaceOne(context.TODO(), bson.M{"_id": id}, activity)

	return err
}

func DeleteActivity(id primitive.ObjectID) error {

	// Delete the activity from the MongoDB collection
	_, err := ActivitiesCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err

}

func GetActivity(id primitive.ObjectID) (*Activity, error) {
	// Get the activity from the MongoDB collection
	var activity Activity
	err := ActivitiesCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&activity)
	return &activity, err
}

func GetActivityByName(name string) (*Activity, error) {
	// Get the activity from the MongoDB collection
	var activity Activity
	err := ActivitiesCollection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&activity)
	return &activity, err
}

func GetActivities() ([]*Activity, error) {
	// Get all the activities from the MongoDB collection
	cursor, err := ActivitiesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var activities []*Activity
	if err = cursor.All(context.TODO(), &activities); err != nil {
		return nil, err
	}
	return activities, nil
}
