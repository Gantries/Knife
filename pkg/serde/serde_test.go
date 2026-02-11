package serde

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerde(t *testing.T) {
	v := struct {
		Name string `json:"name,omitempty" default:"foo"`
		Age  int    `json:"age,omitempty" default:"1"`
	}{
		Name: "Go",
		Age:  10,
	}
	s, e := Serialize(v)
	assert.Nil(t, e)
	assert.True(t, string(s) == `{"name":"Go","age":10}`)
	a, e := Deserialize[map[string]any]([]byte(`{"name":"Go","age":10}`))
	assert.Nil(t, e)
	assert.True(t, (*a)["name"] == "Go")
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid JSON object", `{"key": "value"}`, true},
		{"Invalid JSON", `{"key": "value"`, false},
		{"Empty string", "", false},
		{"JSON array", `["key"]`, false},
		{"Pure number", "123", false},
		{"String", `"string"`, false},
		{"Boolean", `true`, false},
		{"Null", `null`, true},
		{"Nested JSON object", `{"key": {"nested": "value"}}`, true},
		{"JSON object with array", `{"key": ["value"]}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSON(tt.input); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSONArray(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid JSON array", `["key"]`, true},
		{"Invalid JSON", `["key"`, false},
		{"Empty string", "", false},
		{"JSON object", `{"key": "value"}`, false},
		{"Pure number", "123", false},
		{"String", `"string"`, false},
		{"Boolean", `true`, false},
		{"Null", `null`, true},
		{"Nested JSON array", `[{"nested": "value"}]`, true},
		{"JSON array with object", `[{"key": "value"}]`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONArray(tt.input); got != tt.want {
				t.Errorf("IsJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
