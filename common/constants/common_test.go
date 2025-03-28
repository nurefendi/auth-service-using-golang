package constants

import "testing"

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{
			name:     "PRODUCTION",
			actual:   PRODUCTION,
			expected: "production",
		},
		{
			name:     "LOCAL",
			actual:   LOCAL,
			expected: "local",
		},
		{
			name:     "DATE_FORMAT_DEFAULT",
			actual:   DATE_FORMAT_DEFAULT,
			expected: "2006-01-02 15:04:05",
		},
		{
			name:     "DATE_FORMAT_YYYY_MM_DD",
			actual:   DATE_FORMAT_YYYY_MM_DD,
			expected: "2006-01-02",
		},
		{
			name:     "CHANNEL_ID",
			actual:   CHANNEL_ID,
			expected: "X-CHANNEL-ID",
		},
		{
			name:     "BEARER",
			actual:   BEARER,
			expected: "Bearer ",
		},
		{
			name:     "CHANNEL_SYSTEM",
			actual:   CHANNEL_SYSTEM,
			expected: "INTERNAL_SYSTEM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, tt.actual)
			}
		})
	}
}