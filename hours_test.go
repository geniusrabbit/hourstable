package hourstable

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test_TestHour(t *testing.T) {
	var tests = []struct {
		hours      Hours
		weekDay    time.Weekday
		hour       byte
		allActive  bool
		notActive  bool
		testResult bool
	}{
		{
			hours:      MustHoursByString("*"),
			weekDay:    time.Monday,
			hour:       1,
			allActive:  true,
			notActive:  false,
			testResult: true,
		},
		{
			hours:      MustHoursByString("1000000"),
			weekDay:    time.Sunday,
			hour:       1,
			allActive:  false,
			notActive:  false,
			testResult: false,
		},
		{
			hours:      MustHoursByString("1001100"),
			weekDay:    time.Sunday,
			hour:       4,
			allActive:  false,
			notActive:  false,
			testResult: true,
		},
		{
			hours:      MustHoursByString("10011001111110011001111.1001100"),
			weekDay:    time.Monday,
			hour:       4,
			allActive:  false,
			notActive:  false,
			testResult: true,
		},
		{
			hours: MustHoursByString(DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString +
				DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString),
			weekDay:    time.Monday,
			hour:       4,
			allActive:  false,
			notActive:  true,
			testResult: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if test.hours.TestHour(test.weekDay, test.hour) != test.testResult {
				t.Errorf("test hour fail: %d, %d => %t", test.weekDay, test.hour, test.testResult)
			}

			if test.hours.IsAllActive() != test.allActive {
				t.Errorf("IsAllActive should be %v", test.allActive)
			}

			if test.hours.IsNoActive() != test.notActive {
				t.Errorf("IsNoActive should be %v", test.notActive)
			}

			test.hours.SetHour(test.weekDay, test.hour, !test.testResult)
			if test.hours.TestHour(test.weekDay, test.hour) == test.testResult {
				t.Errorf("test2 hour fail: %d, %d => %t", test.weekDay, test.hour, !test.testResult)
			}
		})
	}
}

func Test_HourJSONMarshal(t *testing.T) {
	type item struct {
		Hours Hours `json:"hours"`
	}

	var tests = []struct {
		hoursString string
		result      string
	}{
		{
			hoursString: "*",
			result:      `{"hours":"*"}`,
		},
		{
			hoursString: ActiveWeekHoursString,
			result:      `{"hours":"*"}`,
		},
		{
			hoursString: ActiveDayHoursString + DisabledDayHoursString + ActiveDayHoursString,
			result: `{"hours":"` + ActiveDayHoursString + DisabledDayHoursString + ActiveDayHoursString +
				DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString + `"}`,
		},
		{
			hoursString: "000000000000111111011111" + DisabledDayHoursString + ActiveDayHoursString +
				DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString + ActiveDayHoursString,
			result: `{"hours":"` + "000000000000111111011111" + DisabledDayHoursString + ActiveDayHoursString +
				DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString + ActiveDayHoursString + `"}`,
		},
	}

	for _, test := range tests {
		it := &item{Hours: MustHoursByString(test.hoursString)}
		data, err := json.Marshal(it)

		if err != nil {
			t.Errorf("Hour JSON encode error: %s", err.Error())
		}

		if string(data) != test.result {
			t.Errorf("Invalid data encodeing [%s] should be [%s]", string(data), test.result)
		}
	}
}

func Test_HourJSONUnmarshal(t *testing.T) {
	type item struct {
		Hours Hours `json:"hours"`
	}

	var tests = []struct {
		hoursJSON string
		result    Hours
	}{
		{
			hoursJSON: `{"hours":"*"}`,
			result:    MustHoursByString("*"),
		},
		{
			hoursJSON: `{"hours":"` + ActiveWeekHoursString + `"}`,
			result:    MustHoursByString("*"),
		},
		{
			hoursJSON: `{"hours":"` + ActiveDayHoursString + `"}`,
			result:    MustHoursByString(ActiveDayHoursString),
		},
		{
			hoursJSON: `{"hours":"` + ActiveDayHoursString + DisabledDayHoursString + ActiveDayHoursString + `"}`,
			result:    MustHoursByString(ActiveDayHoursString + DisabledDayHoursString + ActiveDayHoursString),
		},
	}

	for _, test := range tests {
		var (
			it  item
			err = json.Unmarshal([]byte(test.hoursJSON), &it)
		)

		if err != nil {
			t.Errorf("Hour JSON decode error: %s", err.Error())
		}

		if !it.Hours.Equal(test.result) {
			t.Errorf("Invalid data encodeing [%s] should be [%s]", it.Hours.String(), test.result.String())
		}
	}
}

func TestHours_Value(t *testing.T) {
	tests := []struct {
		name     string
		hours    Hours
		expected string
	}{
		{
			name:     "nil hours (all active)",
			hours:    nil,
			expected: "*",
		},
		{
			name:     "empty hours",
			hours:    make(Hours, 24),
			expected: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name: "business hours",
			hours: func() Hours {
				h := make(Hours, 24)
				h.SetHour(time.Monday, 9, true)
				h.SetHour(time.Monday, 10, true)
				return h
			}(),
			expected: "000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.hours.Value()
			if err != nil {
				t.Errorf("Value() error = %v", err)
				return
			}
			if value != tt.expected {
				t.Errorf("Value() = %v, expected %v", value, tt.expected)
			}
		})
	}
}

func TestHours_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected Hours
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "string input - all active",
			input:    "*",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "string input - empty",
			input:    "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: make(Hours, 24),
			wantErr:  false,
		},
		{
			name:     "[]byte input",
			input:    []byte("*"),
			expected: nil,
			wantErr:  false,
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
		{
			name:    "invalid string format",
			input:   strings.Repeat("1", 200), // Too long, should trigger error
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h Hours
			err := h.Scan(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Scan() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Scan() error = %v", err)
				return
			}

			if !h.Equal(tt.expected) {
				t.Errorf("Scan() result = %v, expected %v", h, tt.expected)
			}
		})
	}
}

func TestHours_Merge(t *testing.T) {
	tests := []struct {
		name     string
		h1       Hours
		h2       Hours
		expected Hours
	}{
		{
			name:     "merge with nil (empty base)",
			h1:       nil,
			h2:       MustHoursByString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			expected: nil, // nil base doesn't change
		},
		{
			name: "merge nil into existing",
			h1:   make(Hours, 24),
			h2:   nil,
			expected: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] = 0xff
				}
				return h
			}(),
		},
		{
			name: "merge two valid hours",
			h1: func() Hours {
				h := make(Hours, 24)
				h.SetHour(time.Monday, 9, true)
				return h
			}(),
			h2: func() Hours {
				h := make(Hours, 24)
				h.SetHour(time.Tuesday, 10, true)
				return h
			}(),
			expected: func() Hours {
				h := make(Hours, 24)
				h.SetHour(time.Monday, 9, true)
				h.SetHour(time.Tuesday, 10, true)
				return h
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h1.Merge(tt.h2)
			if !tt.h1.Equal(tt.expected) {
				t.Errorf("Merge() result = %v, expected %v", tt.h1, tt.expected)
			}
		})
	}
}

func TestHours_TestTime(t *testing.T) {
	// Create business hours: Monday-Friday 9AM-5PM
	businessHours := make(Hours, 24)
	for day := time.Monday; day <= time.Friday; day++ {
		for hour := 9; hour < 17; hour++ {
			businessHours.SetHour(day, byte(hour), true)
		}
	}

	tests := []struct {
		name     string
		hours    Hours
		testTime time.Time
		expected bool
	}{
		{
			name:     "nil hours (all active)",
			hours:    nil,
			testTime: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), // Monday 10:30 AM
			expected: true,
		},
		{
			name:     "business hours - active time",
			hours:    businessHours,
			testTime: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), // Monday 10:30 AM
			expected: true,
		},
		{
			name:     "business hours - inactive time (weekend)",
			hours:    businessHours,
			testTime: time.Date(2024, 1, 13, 10, 30, 0, 0, time.UTC), // Saturday 10:30 AM
			expected: false,
		},
		{
			name:     "business hours - inactive time (early morning)",
			hours:    businessHours,
			testTime: time.Date(2024, 1, 15, 6, 30, 0, 0, time.UTC), // Monday 6:30 AM
			expected: false,
		},
		{
			name:     "business hours - inactive time (late evening)",
			hours:    businessHours,
			testTime: time.Date(2024, 1, 15, 20, 30, 0, 0, time.UTC), // Monday 8:30 PM
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hours.TestTime(tt.testTime)
			if result != tt.expected {
				t.Errorf("TestTime() = %v, expected %v for time %v", result, tt.expected, tt.testTime)
			}
		})
	}
}

func TestHoursToBinary(t *testing.T) {
	tests := []struct {
		name      string
		hours     string
		dayOfWeek time.Weekday
		expected  func() Hours
	}{
		{
			name:      "all active wildcard",
			hours:     "*",
			dayOfWeek: time.Monday,
			expected: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] |= byte(0x01) << byte(time.Monday)
				}
				return h
			},
		},
		{
			name:      "empty string (inactive)",
			hours:     "",
			dayOfWeek: time.Monday,
			expected: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] &= ^(byte(0x01) << byte(time.Monday))
				}
				return h
			},
		},
		{
			name:      "specific hours pattern",
			hours:     "110000000000000000000000",
			dayOfWeek: time.Sunday,
			expected: func() Hours {
				h := make(Hours, 24)
				h[0] |= byte(0x01) << byte(time.Sunday)
				h[1] |= byte(0x01) << byte(time.Sunday)
				return h
			},
		},
		{
			name:      "partial pattern with different day",
			hours:     "001100",
			dayOfWeek: time.Friday,
			expected: func() Hours {
				h := make(Hours, 24)
				h[2] |= byte(0x01) << byte(time.Friday)
				h[3] |= byte(0x01) << byte(time.Friday)
				return h
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := make(Hours, 24)
			hoursToBinary(h, tt.hours, tt.dayOfWeek)

			expected := tt.expected()
			if !h.Equal(expected) {
				t.Errorf("hoursToBinary() result = %v, expected %v", h, expected)
			}
		})
	}
}

func TestBinaryToHours(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() Hours
		dayOfWeek time.Weekday
		expected  string
	}{
		{
			name: "all inactive",
			setup: func() Hours {
				return make(Hours, 24)
			},
			dayOfWeek: time.Monday,
			expected:  "000000000000000000000000",
		},
		{
			name: "all active for specific day",
			setup: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] |= byte(0x01) << byte(time.Tuesday)
				}
				return h
			},
			dayOfWeek: time.Tuesday,
			expected:  "111111111111111111111111",
		},
		{
			name: "mixed pattern",
			setup: func() Hours {
				h := make(Hours, 24)
				h[0] |= byte(0x01) << byte(time.Wednesday)
				h[5] |= byte(0x01) << byte(time.Wednesday)
				h[23] |= byte(0x01) << byte(time.Wednesday)
				return h
			},
			dayOfWeek: time.Wednesday,
			expected:  "100001000000000000000001",
		},
		{
			name: "different day should return all zeros",
			setup: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] |= byte(0x01) << byte(time.Monday)
				}
				return h
			},
			dayOfWeek: time.Sunday, // Different day
			expected:  "000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.setup()
			result := binaryToHours(h, tt.dayOfWeek)
			if result != tt.expected {
				t.Errorf("binaryToHours() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestBinaryToHoursShort(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() Hours
		dayOfWeek time.Weekday
		expected  string
	}{
		{
			name: "all active should return wildcard",
			setup: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] |= byte(0x01) << byte(time.Monday)
				}
				return h
			},
			dayOfWeek: time.Monday,
			expected:  "*",
		},
		{
			name: "all inactive should return empty string",
			setup: func() Hours {
				return make(Hours, 24)
			},
			dayOfWeek: time.Tuesday,
			expected:  "",
		},
		{
			name: "mixed pattern should return full string",
			setup: func() Hours {
				h := make(Hours, 24)
				h[9] |= byte(0x01) << byte(time.Wednesday)
				h[10] |= byte(0x01) << byte(time.Wednesday)
				return h
			},
			dayOfWeek: time.Wednesday,
			expected:  "000000000110000000000000",
		},
		{
			name: "partial active at end",
			setup: func() Hours {
				h := make(Hours, 24)
				h[22] |= byte(0x01) << byte(time.Friday)
				h[23] |= byte(0x01) << byte(time.Friday)
				return h
			},
			dayOfWeek: time.Friday,
			expected:  "000000000000000000000011",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.setup()
			result := binaryToHoursShort(h, tt.dayOfWeek)
			if result != tt.expected {
				t.Errorf("binaryToHoursShort() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestHours_Clone(t *testing.T) {
	tests := []struct {
		name  string
		hours Hours
	}{
		{
			name:  "nil hours",
			hours: nil,
		},
		{
			name:  "empty hours",
			hours: make(Hours, 24),
		},
		{
			name: "business hours",
			hours: func() Hours {
				h := make(Hours, 24)
				h.SetHour(time.Monday, 9, true)
				h.SetHour(time.Friday, 17, true)
				return h
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := tt.hours.Clone()

			// Test that cloned is equal to original
			if !cloned.Equal(tt.hours) {
				t.Errorf("Clone() result not equal to original")
			}

			// Test that it's a separate copy (if not nil)
			if tt.hours != nil && cloned != nil {
				// Modify original
				if len(tt.hours) > 0 {
					original := tt.hours[0]
					tt.hours[0] = ^original

					// Cloned should not be affected
					if cloned[0] != original {
						t.Errorf("Clone() not independent copy - modification affected clone")
					}
				}
			}
		})
	}
}

func TestMustHoursByString_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustHoursByString should panic on invalid input")
		}
	}()

	// This should panic due to invalid length
	invalidString := ""
	for i := 0; i < 200; i++ {
		invalidString += "1"
	}
	MustHoursByString(invalidString)
}

func TestHoursByString_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "valid asterisk",
			input:   "*",
			wantErr: false,
		},
		{
			name:    "valid full week string",
			input:   ActiveWeekHoursString,
			wantErr: false,
		},
		{
			name: "too long string",
			input: func() string {
				s := ""
				for i := 0; i < 200; i++ {
					s += "1"
				}
				return s
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HoursByString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HoursByString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHours_Equal_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		h1   Hours
		h2   Hours
		want bool
	}{
		{
			name: "both nil",
			h1:   nil,
			h2:   nil,
			want: true,
		},
		{
			name: "nil vs empty",
			h1:   nil,
			h2:   make(Hours, 24),
			want: false,
		},
		{
			name: "different lengths - h1 longer",
			h1:   make(Hours, 25),
			h2:   make(Hours, 24),
			want: true, // should be equal if extra bytes are 0
		},
		{
			name: "different lengths - h2 longer",
			h1:   make(Hours, 24),
			h2:   make(Hours, 25),
			want: true, // should be equal if extra bytes are 0
		},
		{
			name: "different lengths with non-zero extra",
			h1: func() Hours {
				h := make(Hours, 25)
				h[24] = 1 // non-zero in extra byte
				return h
			}(),
			h2:   make(Hours, 24),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h1.Equal(tt.h2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHours_SetHour_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		initial  Hours
		weekDay  time.Weekday
		hour     byte
		active   bool
		expected Hours
	}{
		{
			name:    "set on empty hours",
			initial: make(Hours, 24),
			weekDay: time.Monday,
			hour:    9,
			active:  true,
			expected: func() Hours {
				h := make(Hours, 24)
				h[9] |= byte(0x01) << byte(time.Monday)
				return h
			}(),
		},
		{
			name: "toggle same hour twice",
			initial: func() Hours {
				h := make(Hours, 24)
				h.SetHour(time.Monday, 9, true)
				return h
			}(),
			weekDay: time.Monday,
			hour:    9,
			active:  true, // Same as current state - should not change
			expected: func() Hours {
				h := make(Hours, 24)
				h[9] |= byte(0x01) << byte(time.Monday)
				return h
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy to avoid modifying the test case
			var result Hours
			if tt.initial != nil {
				result = make(Hours, len(tt.initial))
				copy(result, tt.initial)
			}

			result.SetHour(tt.weekDay, tt.hour, tt.active)

			if !result.Equal(tt.expected) {
				t.Errorf("SetHour() result = %v, expected %v", result.String(), tt.expected.String())
			}
		})
	}
}

func TestHours_IsAllActive_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		hours Hours
		want  bool
	}{
		{
			name:  "nil hours",
			hours: nil,
			want:  true,
		},
		{
			name:  "empty hours",
			hours: make(Hours, 0),
			want:  true,
		},
		{
			name:  "short hours array",
			hours: make(Hours, 10),
			want:  false,
		},
		{
			name: "partial active hours",
			hours: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 20; i++ {
					h[i] = 0x7f // all days active
				}
				// Leave last 4 hours inactive
				return h
			}(),
			want: false,
		},
		{
			name: "all active hours",
			hours: func() Hours {
				h := make(Hours, 24)
				for i := 0; i < 24; i++ {
					h[i] = 0x7f // all days active
				}
				return h
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hours.IsAllActive(); got != tt.want {
				t.Errorf("IsAllActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_Hours(b *testing.B) {
	var (
		hours = []Hours{
			MustHoursByString(ActiveWeekHoursString),
			MustHoursByString("1001100"),
			MustHoursByString("1000000"),
			MustHoursByString("10011001111110011001111"),
		}
		now = time.Now()
	)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var i = 0
		for pb.Next() {
			h := hours[i%len(hours)]
			_ = h.TestHour(time.Weekday(i%7), byte(i%24))
			_ = h.TestTime(now)
			i++
		}
	})
}
