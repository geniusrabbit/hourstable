package hourstable

import (
	"encoding/json"
	"fmt"
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
