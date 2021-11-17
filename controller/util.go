package controller

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Used for partial update on struct in a document
func Flatten(v interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)

	val, ok := asStruct(val)
	if !ok {
		return nil, errors.New("v must be a struct or a pointer to a struct")
	}

	m := make(map[string]interface{})
	if err := flattenFields(val, m, ""); err != nil {
		return nil, err
	}

	return m, nil
}

func flattenFields(v reflect.Value, m map[string]interface{}, p string) error {
	for i := 0; i < v.NumField(); i++ {
		tags, _ := bsoncodec.DefaultStructTagParser(v.Type().Field(i))

		if tags.Skip {
			continue
		}

		field := v.Field(i)
		if tags.OmitEmpty && field.IsZero() || !field.CanInterface() {
			continue
		}

		if _, ok := field.Interface().(bson.ValueMarshaler); !ok {
			if s, ok := asStruct(field); ok && hasExportedField(s) {
				fp := p
				if !tags.Inline {
					fp = p + tags.Name + "."
				}
				if err := flattenFields(s, m, fp); err != nil {
					return err
				}
				continue
			}
		}

		key := p + tags.Name
		if _, ok := m[key]; ok {
			return fmt.Errorf("duplicated key %s", key)
		}

		m[key] = field.Interface()
	}

	return nil
}

func asStruct(v reflect.Value) (reflect.Value, bool) {
	for {
		switch v.Kind() {
		case reflect.Struct:
			return v, true
		case reflect.Ptr:
			v = v.Elem()
		default:
			return reflect.Value{}, false
		}
	}
}

func hasExportedField(s reflect.Value) bool {
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).CanInterface() {
			return true
		}
	}
	return false
}

func removeSliceInt(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func removeSliceID(s []primitive.ObjectID, i int) []primitive.ObjectID {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
