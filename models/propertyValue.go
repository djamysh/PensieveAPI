package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PropertyValue struct {
	Key   primitive.ObjectID
	Value interface{}
}

func (pv *PropertyValue) UnmarshalBSON(data []byte) error {
	m := make(map[string]interface{})
	err := bson.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	if key, ok := m["key"]; ok {
		if objectID, ok := key.(primitive.ObjectID); ok {
			pv.Key = objectID
		}
	}

	pv.Value = m["value"]
	return nil
}

func (pv PropertyValue) MarshalBSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["key"] = pv.Key
	m["value"] = pv.Value
	return bson.Marshal(m)
}
