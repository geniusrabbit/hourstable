//
// @project GeniusRabbit 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

package hourstable

import (
	"encoding/json"
	"testing"
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
