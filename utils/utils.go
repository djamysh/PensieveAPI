package utils

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CleanInput(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}

func MarshallTimelings(timelings map[string]int64) bson.M {
	primitiveTimeling := primitive.M{}
	for key, value := range timelings {
		primitiveTimeling[key] = value
	}
	return primitiveTimeling
}

/*
string -> interface{} type of map marshalling function,but currently not in use
func MarshallMap(mapVar interface{}) (primitive.M, error) {
        if reflect.TypeOf(mapVar).Kind() != reflect.Map {
                return nil, errors.New("input is not a map")
        }

        m, ok := mapVar.(map[string]interface{})
        if !ok {
                return nil, errors.New("failed to convert input to map[string]interface{}")
        }

        primitiveMarshall := primitive.M{}
        for key, value := range m {
                switch reflect.TypeOf(value).Kind() {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                        primitiveMarshall[key] = value
                case reflect.Float32, reflect.Float64:
                        primitiveMarshall[key] = value
                case reflect.Map:
                        nestedMap, err := MarshallMap(value)
                        if err != nil {
                                return nil, err
                        }
                        primitiveMarshall[key] = nestedMap
                default:
                        return nil, errors.New("unsupported type in map")
                }
        }

        return primitiveMarshall, nil
}
*/
