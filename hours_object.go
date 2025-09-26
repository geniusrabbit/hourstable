//
// @project GeniusRabbit 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package hourstable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

//easyjson:json
type timetableJSON struct {
	Monday    string `json:"mon,omitempty" yaml:"mon,omitempty"`
	Tuesday   string `json:"tue,omitempty" yaml:"tue,omitempty"`
	Wednesday string `json:"wed,omitempty" yaml:"wed,omitempty"`
	Thursday  string `json:"thu,omitempty" yaml:"thu,omitempty"`
	Friday    string `json:"fri,omitempty" yaml:"fri,omitempty"`
	Saturday  string `json:"sat,omitempty" yaml:"sat,omitempty"`
	Sunday    string `json:"sun,omitempty" yaml:"sun,omitempty"`
}

func (tt *timetableJSON) ToHours() Hours {
	hours := make(Hours, 24)
	tt.ToHoursObject(hours)
	return hours
}

func (tt *timetableJSON) ToHoursObject(hours Hours) error {
	hoursToBinary(hours, tt.Sunday, time.Sunday)
	hoursToBinary(hours, tt.Monday, time.Monday)
	hoursToBinary(hours, tt.Tuesday, time.Tuesday)
	hoursToBinary(hours, tt.Wednesday, time.Wednesday)
	hoursToBinary(hours, tt.Thursday, time.Thursday)
	hoursToBinary(hours, tt.Friday, time.Friday)
	hoursToBinary(hours, tt.Saturday, time.Saturday)
	return nil
}

func (tt *timetableJSON) FromHours(hours Hours) {
	tt.Sunday = binaryToHoursShort(hours, time.Sunday)
	tt.Monday = binaryToHoursShort(hours, time.Monday)
	tt.Tuesday = binaryToHoursShort(hours, time.Tuesday)
	tt.Wednesday = binaryToHoursShort(hours, time.Wednesday)
	tt.Thursday = binaryToHoursShort(hours, time.Thursday)
	tt.Friday = binaryToHoursShort(hours, time.Friday)
	tt.Saturday = binaryToHoursShort(hours, time.Saturday)
}

// HoursObject supports the JSON format of storing
type HoursObject Hours

// HoursByJSON decodes JSON format of timetable
func HoursByJSON(data []byte) (Hours, error) {
	var (
		timetable timetableJSON
		err       = json.Unmarshal(data, &timetable)
	)
	if err != nil {
		return nil, err
	}
	return timetable.ToHours(), nil
}

// String implementation of fmt.Stringer
func (h HoursObject) String() string {
	var timetable timetableJSON
	timetable.FromHours(Hours(h))
	data, _ := json.Marshal(&timetable)
	return string(data)
}

// Value implementation of valuer for database/sql
func (h HoursObject) Value() (driver.Value, error) {
	return h.MarshalJSON()
}

// Scan - Implement the database/sql scanner interface
func (h *HoursObject) Scan(value any) (err error) {
	if value == nil {
		*h = nil
		return nil
	}

	var newHours Hours
	switch v := value.(type) {
	case []byte:
		if newHours, err = HoursByJSON(v); err == nil {
			*h = HoursObject(newHours)
		}
	case string:
		if newHours, err = HoursByJSON([]byte(v)); err == nil {
			*h = HoursObject(newHours)
		}
	default:
		err = fmt.Errorf("[hours_json] unsupported decode type %T", value)
	}
	return
}

// Merge from another hours
func (h HoursObject) Merge(h2 Hours) {
	Hours(h).Merge(h2)
}

// IsAllActive then return the true
func (h HoursObject) IsAllActive() bool {
	return Hours(h).IsAllActive()
}

// IsNoActive then return the true
func (h HoursObject) IsNoActive() bool {
	return Hours(h).IsNoActive()
}

// Equal comarison of two hour tables
func (h HoursObject) Equal(h2 Hours) bool {
	return Hours(h).Equal(h2)
}

// TestHour hour
func (h HoursObject) TestHour(weekDay time.Weekday, hour byte) bool {
	return Hours(h).TestHour(weekDay, hour)
}

// TestTime hour
func (h HoursObject) TestTime(t time.Time) bool {
	return Hours(h).TestTime(t)
}

// SetHour as active or no
func (h *HoursObject) SetHour(weekDay time.Weekday, hour byte, active bool) {
	(*Hours)(h).SetHour(weekDay, hour, active)
}

// MarshalJSON implements the functionality of json.Marshaler interface
func (h HoursObject) MarshalJSON() ([]byte, error) {
	var timetable timetableJSON
	timetable.FromHours(Hours(h))
	return json.Marshal(&timetable)
}

// UnmarshalJSON implements the functionality of json.Unmarshaler interface
func (h *HoursObject) UnmarshalJSON(data []byte) error {
	newHours, err := HoursByJSON(data)
	if err != nil {
		return err
	}
	*h = HoursObject(newHours)
	return nil
}

// MarshalYAML implements the functionality of yaml.Marshaler interface
func (h HoursObject) MarshalYAML() (any, error) {
	var timetable timetableJSON
	timetable.FromHours(Hours(h))
	return &timetable, nil
}

// UnmarshalYAML implements the functionality of yaml.Unmarshaler interface
func (h *HoursObject) UnmarshalYAML(node *yaml.Node) error {
	var timetable timetableJSON
	if err := node.Decode(&timetable); err != nil {
		return err
	}
	*h = HoursObject(timetable.ToHours())
	return nil
}

// Clone returns a copy of HoursObject
func (h HoursObject) Clone() HoursObject {
	if h == nil {
		return nil
	}
	newHours := make(Hours, len(h))
	copy(newHours, h)
	return HoursObject(newHours)
}

var (
	_ json.Marshaler   = (HoursObject)(nil)
	_ json.Unmarshaler = (*HoursObject)(nil)
	_ yaml.Marshaler   = (HoursObject)(nil)
	_ yaml.Unmarshaler = (*HoursObject)(nil)
	_ driver.Valuer    = (HoursObject)(nil)
	_ sql.Scanner      = (*HoursObject)(nil)
)
