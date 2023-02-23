package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ValidDataTypes = [5]string{"string", "number", "timelings", "string array", "number array"}

type Property struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name" validate:"unique"`
	Description   string             `bson:"description" json:"description"`
	ValueDataType string             `bson:"valueDataType" json:"valueDataType"`
}

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

func UpdateProperty(id primitive.ObjectID, update bson.M) (*Property, error) {
	var property Property

	if err := PropertiesCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, bson.M{"$set": update}).Decode(&property); err != nil {
		return nil, err
	}
	return &property, nil
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
