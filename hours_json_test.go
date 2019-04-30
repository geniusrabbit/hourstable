//
// @project GeniusRabbit 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

package hourstable

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_JSONEncodeDecode(t *testing.T) {
	type item struct {
		Hours HoursJSON `json:"hours"`
	}
	var tests = []struct {
		timetable string
		result    string
	}{
		{
			timetable: `{"hours":{}}`,
			result:    `{"hours":{}}`,
		},
		{
			timetable: `{"hours":{"mon":"111111111111111111111111"}}`,
			result:    `{"hours":{"mon":"*"}}`,
		},
		{
			timetable: `{"hours":{"mon":"111111111111111111111111","sun":"000000000000000000000000"}}`,
			result:    `{"hours":{"mon":"*"}}`,
		},
		{
			timetable: `{"hours":{"mon":"111111111111111111111111","sun":"000001110000000000000000"}}`,
			result:    `{"hours":{"mon":"*","sun":"000001110000000000000000"}}`,
		},
		{
			timetable: `{"hours":{"mon":"111111111111111111111111","sun":"000001110000000000000000","fri":"11"}}`,
			result:    `{"hours":{"mon":"*","fri":"110000000000000000000000","sun":"000001110000000000000000"}}`,
		},
	}

	for _, test := range tests {
		var it item
		if err := json.Unmarshal([]byte(test.timetable), &it); err != nil {
			t.Errorf("invalid timetable unmarshal: %s", err.Error())
		}
		if data, _ := json.Marshal(it); string(data) != test.result {
			t.Errorf("invalid timetable marshal [%s] must be [%s]", string(data), test.result)
		}
	}
}

func Test_TestHoursJSON(t *testing.T) {
	var tests = []struct {
		hours      HoursJSON
		weekDay    time.Weekday
		hour       byte
		allActive  bool
		notActive  bool
		testResult bool
	}{
		{
			hours:      HoursJSON(MustHoursByString("*")),
			weekDay:    time.Monday,
			hour:       1,
			allActive:  true,
			notActive:  false,
			testResult: true,
		},
		{
			hours:      HoursJSON(MustHoursByString("1000000")),
			weekDay:    time.Sunday,
			hour:       1,
			allActive:  false,
			notActive:  false,
			testResult: false,
		},
		{
			hours:      HoursJSON(MustHoursByString("1001100")),
			weekDay:    time.Sunday,
			hour:       4,
			allActive:  false,
			notActive:  false,
			testResult: true,
		},
		{
			hours:      HoursJSON(MustHoursByString("10011001111110011001111.1001100")),
			weekDay:    time.Monday,
			hour:       4,
			allActive:  false,
			notActive:  false,
			testResult: true,
		},
		{
			hours: HoursJSON(MustHoursByString(DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString +
				DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString + DisabledDayHoursString)),
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
