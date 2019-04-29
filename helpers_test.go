package hourstable

import "testing"

func Test_ActiveHoursRangeString(t *testing.T) {
	var tests = []struct {
		from, to byte
		result   string
	}{
		{
			from:   0,
			to:     10,
			result: "111111111100000000000000",
		},
		{
			from:   10,
			to:     24,
			result: "000000000011111111111111",
		},
		{
			from:   10,
			to:     20,
			result: "000000000011111111110000",
		},
	}

	for _, test := range tests {
		if rangeHours := ActiveHoursRangeString(test.from, test.to); rangeHours != test.result {
			t.Errorf("invalid range [%s] must be [%s]", rangeHours, test.result)
		}
	}
}
