package main

import (
	"testing"
)

func TestIsValidUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		expectedValue bool
	}{
		{
			name:          "Valid user (root)",
			username:      "root",
			expectedValue: true,
		},
		{
			name:          "Valid user (nobody)",
			username:      "nobody",
			expectedValue: true,
		},
		{
			name:          "Invalid user - starts with number",
			username:      "1kofabhhfsbf6krb",
			expectedValue: false,
		},
		{
			name:          "Invalid user - contains uppercase letters",
			username:      "XD5hObIMZF2zKS7W",
			expectedValue: false,
		},
		{
			name:          "Invalid user - too long",
			username:      "idfkjcacexia1dyji5iwcfweoliamzpn1",
			expectedValue: false,
		},
		{
			name:          "Invalid user - contains invalid characters",
			username:      "moctcg!@",
			expectedValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isValidUser(test.username)
			if result != test.expectedValue {
				t.Errorf("Expected %v, got %v", test.expectedValue, result)
			}
		})
	}
}
