package enums

import "testing"

func TestHttpMethodConstants(t *testing.T) {
	tests := []struct {
		name     string
		actual   HttpMethod
		expected string
	}{
		{"GET", GET, "GET"},
		{"PUT", PUT, "PUT"},
		{"DELETE", DELETE, "DELETE"},
		{"POST", POST, "POST"},
		{"PATCH", PATCH, "PATCH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.actual) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.actual)
			}
		})
	}
}

func TestNameMethod(t *testing.T) {
	tests := []struct {
		method   HttpMethod
		expected string
	}{
		{GET, "GET"},
		{PUT, "PUT"},
		{DELETE, "DELETE"},
		{POST, "POST"},
		{PATCH, "PATCH"},
		{HttpMethod("UNKNOWN"), "GET"}, // Test default case
	}

	for _, tt := range tests {
		t.Run(string(tt.method), func(t *testing.T) {
			result := tt.method.Name()
			if result != tt.expected {
				t.Errorf("For method %s, expected %s, got %s", tt.method, tt.expected, result)
			}
		})
	}
}

func TestGetValueFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected HttpMethod
	}{
		{"GET", GET},
		{"PUT", PUT},
		{"DELETE", DELETE},
		{"POST", POST},
		{"PATCH", PATCH},
		{"UNKNOWN", GET}, // Test default case
		{"", GET},        // Test empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetValue(tt.input)
			if result != tt.expected {
				t.Errorf("For input %s, expected %s, got %s", tt.input, tt.expected, result)
			}
		})
	}
}