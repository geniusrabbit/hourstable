package hourstable

import (
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestHoursYAMLMarshal(t *testing.T) {
	tests := []struct {
		name     string
		hours    Hours
		expected string
	}{
		{
			name:     "all_active",
			hours:    nil, // nil represents all active
			expected: "'*'\n",
		},
		{
			name:     "empty_hours",
			hours:    make(Hours, 24),
			expected: "\"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"\n",
		},
		{
			name: "business_hours",
			hours: func() Hours {
				h := make(Hours, 24)
				// Monday-Friday 9AM-5PM
				for day := time.Monday; day <= time.Friday; day++ {
					for hour := 9; hour < 17; hour++ {
						h.SetHour(day, byte(hour), true)
					}
				}
				return h
			}(),
			expected: "\"000000000000000000000000000000000111111110000000000000000111111110000000000000000111111110000000000000000111111110000000000000000111111110000000000000000000000000000000\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(tt.hours)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("Marshal result mismatch\nExpected: %q\nGot:      %q", tt.expected, string(data))
			}
		})
	}
}

func TestHoursYAMLUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected Hours
	}{
		{
			name:     "all_active",
			yaml:     "'*'",
			expected: nil,
		},
		{
			name:     "all_inactive",
			yaml:     "\"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"",
			expected: make(Hours, 24),
		},
		{
			name: "business_hours",
			yaml: "\"000000000000000000000000000000000111111110000000000000000111111110000000000000000111111110000000000000000111111110000000000000000111111110000000000000000000000000000000\"",
			expected: func() Hours {
				h := make(Hours, 24)
				// Monday-Friday 9AM-5PM
				for day := time.Monday; day <= time.Friday; day++ {
					for hour := 9; hour < 17; hour++ {
						h.SetHour(day, byte(hour), true)
					}
				}
				return h
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hours Hours
			err := yaml.Unmarshal([]byte(tt.yaml), &hours)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if !hours.Equal(tt.expected) {
				t.Errorf("Unmarshal result mismatch\nExpected: %v\nGot:      %v", tt.expected, hours)
			}
		})
	}
}

func TestHoursObjectYAMLMarshal(t *testing.T) {
	tests := []struct {
		name     string
		hours    HoursObject
		contains []string // Check if output contains these substrings
	}{
		{
			name:  "all_active",
			hours: nil, // nil represents all active
			contains: []string{
				"mon: '*'",
				"tue: '*'",
				"wed: '*'",
				"thu: '*'",
				"fri: '*'",
				"sat: '*'",
				"sun: '*'",
			},
		},
		{
			name: "business_hours",
			hours: func() HoursObject {
				h := make(Hours, 24)
				// Monday-Friday 9AM-5PM
				for day := time.Monday; day <= time.Friday; day++ {
					for hour := 9; hour < 17; hour++ {
						h.SetHour(day, byte(hour), true)
					}
				}
				return HoursObject(h)
			}(),
			contains: []string{
				"mon: \"000000000111111110000000\"",
				"tue: \"000000000111111110000000\"",
				"wed: \"000000000111111110000000\"",
				"thu: \"000000000111111110000000\"",
				"fri: \"000000000111111110000000\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(tt.hours)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			output := string(data)
			for _, substr := range tt.contains {
				if !containsString(output, substr) {
					t.Errorf("Output should contain %q, but got: %s", substr, output)
				}
			}
		})
	}
}

func TestHoursObjectYAMLUnmarshal(t *testing.T) {
	yamlData := `
mon: "000000000111111110000000"
tue: "000000000111111110000000"  
wed: "000000000111111110000000"
thu: "000000000111111110000000"
fri: "000000000111111110000000"
sat: ""
sun: ""
`

	var hours HoursObject
	err := yaml.Unmarshal([]byte(yamlData), &hours)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Check that Monday 10AM is active (should be true)
	if !Hours(hours).TestHour(time.Monday, 10) {
		t.Error("Expected Monday 10AM to be active")
	}

	// Check that Saturday 10AM is inactive (should be false)
	if Hours(hours).TestHour(time.Saturday, 10) {
		t.Error("Expected Saturday 10AM to be inactive")
	}

	// Check that Monday 8AM is inactive (should be false - before 9AM)
	if Hours(hours).TestHour(time.Monday, 8) {
		t.Error("Expected Monday 8AM to be inactive")
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findIndex(s, substr) >= 0
}

func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
