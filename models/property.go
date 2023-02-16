package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ValidDataTypes = [7]string{"string", "integer", "float", "timelings", "string array", "integer array", "float array"}

// var DefaultTimelingsName = "Default Timelings"
// var DefaultTimelingsPropertyID primitive.ObjectID

type Property struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name" validate:"unique"`
	Description   string             `bson:"description" json:"description"`
	ValueDataType string             `bson:"valueDataType" json:"valueDataType"`
}

/*
// Create a null Property so that at the delete moment deleted property reference will be set to null property
func DefaultTimelingsProperty() (primitive.ObjectID, error) {
	// Check wheter the default timelings Property is created or not, if not  creates.
	var defaultTimelingsProperty *Property

	if err := Client.Database(DBName).Collection("properties").FindOne(context.TODO(), bson.M{"name": DefaultTimelingsName}).Decode(&defaultTimelingsProperty); err != nil {

		// definition of default timelings property
		defaultTimelingsProperty = &Property{ID: primitive.NewObjectID(), Name: "Default Timelings", Description: "Default timelings property", ValueDataType: "timelings"}

		// Creating a record for the default timelings property
		if err = defaultTimelingsProperty.CreateProperty(); err != nil {

			return primitive.NilObjectID, err
		}

		return defaultTimelingsProperty.ID, nil
	}

	return defaultTimelingsProperty.ID, nil
}
*/

func (P *Property) IsValidType() bool {
	// Check wheter the given data type is valid
	for _, t := range ValidDataTypes {
		if P.ValueDataType == t && t != "" {
			// Note :
			// Since the last character in a string is often followed by a zero-value byte,
			// the code point of the last character can appear to be an empty string.
			// Remember that a string array is also a string, and a string is an array of characters.
			// Therefore although it is an array of strings, the last element is the null terminator.
			return true
		}
	}
	return false
}

func (property *Property) CreateProperty() error {

	_, err := PropertiesCollection.InsertOne(context.TODO(), property)

	return err
}

func (property *Property) UpdateProperty(id primitive.ObjectID) error {

	_, err := PropertiesCollection.ReplaceOne(context.TODO(), bson.M{"_id": id}, property)
	return err

}

func DeleteProperty(id primitive.ObjectID) error {

	// Delete the property from the MongoDB collection
	_, err := PropertiesCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func GetProperty(id primitive.ObjectID) (*Property, error) {
	// Get the property from the MongoDB collection
	var property *Property
	err := PropertiesCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&property)
	return property, err
}
func GetPropertyByName(name string) (*Property, error) {
	var property *Property
	err := PropertiesCollection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&property)
	return property, err
}

func GetProperties() ([]*Property, error) {
	// Get all the properties from the MongoDB collection
	cursor, err := PropertiesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var properties []*Property
	if err = cursor.All(context.TODO(), &properties); err != nil {
		return nil, err
	}
	return properties, nil
}
