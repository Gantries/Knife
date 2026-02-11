// Package serde provides serialization and deserialization utilities.
//
// It offers generic functions for JSON and YAML serialization/deserialization
// with type safety through Go generics.
package serde

import (
	"encoding/json"

	"sigs.k8s.io/yaml"
)

func Serialize[T any](v T) (buf []byte, err error) {
	if buf, err = json.Marshal(v); err == nil {
		return
	}
	return nil, err
}

func Deserialize[T any](b []byte) (t *T, err error) {
	var inst T
	if err = json.Unmarshal(b, &inst); err == nil {
		return &inst, nil
	}
	return nil, err
}

func DeserializeArray[T any](b []byte) (a []T, err error) {
	var r []T
	if err = json.Unmarshal(b, &r); err == nil {
		return
	}
	return nil, err
}

func SerializeYAML[T any](v T) (buf []byte, err error) {
	if buf, err = yaml.Marshal(v); err == nil {
		return
	}
	return nil, err
}

func DeserializeYAML[T any](b []byte) (t *T, err error) {
	var r T
	if err = yaml.Unmarshal(b, &r); err == nil {
		return &r, nil
	}
	return nil, err
}

func DeserializeYAMLArray[T any](b []byte) (a []T, err error) {
	if err = yaml.Unmarshal(b, &a); err == nil {
		return
	}
	return nil, err
}

// IsJSON checks whether the provided string can be unmarshal into a JSON object.
// It verifies that the string is not just a JSON-formatted string, but an actual JSON object.
// It returns true if the string is a valid JSON object, and false otherwise.
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// IsJSONArray checks if the provided string represents a valid JSON array.
// It attempts to unmarshal the string into a slice of empty interfaces.
// If the unmarshal is successful, the function returns true, indicating that
// the string is a well-formed JSON array. Otherwise, it returns false.
func IsJSONArray(s string) bool {
	var arr []interface{}
	return json.Unmarshal([]byte(s), &arr) == nil
}
