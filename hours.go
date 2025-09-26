//
// @project GeniusRabbit 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package hourstable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// ErrTooMuchHoursForDecode tells that hours more then for a week
var ErrTooMuchHoursForDecode = errors.New("[hours] too much hours for decode, mpre then 24*7")

const daysBitMask = byte(0x7f)

// Hours type
type Hours []byte

// HoursByString returns hours value or error
func HoursByString(s string) (h Hours, err error) {
	if s == "" || s == "*" || s == AllActiveHoursString || s == ActiveWeekHoursString {
		return nil, nil
	}

	if len(s) > 7*24 {
		err = ErrTooMuchHoursForDecode
	}

	h = make([]byte, 24)
	for i, v := range s {
		if v == '1' {
			h[i%24] |= byte(0x01) << byte(i/24)
		}
	}
	return
}

// MustHoursByString returns hours value or panic
func MustHoursByString(s string) Hours {
	h, err := HoursByString(s)
	if err != nil {
		panic(err)
	}
	return h
}

// String implementation of fmt.Stringer
func (h Hours) String() string {
	if len(h) <= 0 {
		return AllActiveHoursString
	}

	var buff bytes.Buffer
	for dayOfWeek := time.Weekday(0); dayOfWeek < 7; dayOfWeek++ {
		buff.WriteString(binaryToHours(h, dayOfWeek))
	}

	return buff.String()
}

// Value implementation of valuer for database/sql
func (h Hours) Value() (driver.Value, error) {
	return h.String(), nil
}

// Scan - Implement the database/sql scanner interface
func (h *Hours) Scan(value any) (err error) {
	if value == nil {
		*h = nil
		return nil
	}

	var newHours Hours
	switch v := value.(type) {
	case []byte:
		if newHours, err = HoursByString(string(v)); err == nil {
			*h = newHours
		}
	case string:
		if newHours, err = HoursByString(v); err == nil {
			*h = newHours
		}
	default:
		err = fmt.Errorf("[hours] unsupported decode type %T", value)
	}
	return
}

// Merge from another hours
func (h Hours) Merge(h2 Hours) {
	if len(h) < 1 {
		return
	}
	if len(h2) < 1 {
		for i := 0; i < len(h); i++ {
			h[i] = 0xff
		}
	} else {
		for i := 0; i < len(h); i++ {
			h[i] |= h2[i]
		}
	} // end if
}

// IsAllActive then return the true
func (h Hours) IsAllActive() bool {
	if len(h) == 0 {
		return true
	}
	if len(h) < 24 {
		return false
	}
	for _, bt := range h {
		if bt&daysBitMask != daysBitMask {
			return false
		}
	}
	return true
}

// IsNoActive then return the true
func (h Hours) IsNoActive() bool {
	if len(h) < 1 {
		return false
	}
	for _, bt := range h {
		if bt&daysBitMask != 0 {
			return false
		}
	}
	return true
}

// Equal comarison of two hour tables
func (h Hours) Equal(h2 Hours) bool {
	if b1, b2 := h.IsAllActive(), h2.IsAllActive(); b1 || b2 {
		return b1 && b2
	}

	ln := len(h)
	if ln != len(h2) {
		if ln > len(h2) {
			for i := len(h2); i < ln; i++ {
				if h[i] != 0 {
					return false
				}
			}
			ln = len(h2)
		} else {
			for i := len(h); i < len(h2); i++ {
				if h2[i] != 0 {
					return false
				}
			}
		}
	}

	for i := 0; i < ln; i++ {
		if h[i]&daysBitMask != h2[i]&daysBitMask {
			return false
		}
	}

	return true
}

// TestHour hour
func (h Hours) TestHour(weekDay time.Weekday, hour byte) bool {
	return len(h) < 1 || (len(h) > int(hour) && h[hour]&(0x01<<byte(weekDay)) != 0)
}

// TestTime hour
func (h Hours) TestTime(t time.Time) bool {
	if len(h) < 1 {
		return true
	}
	return h.TestHour(t.Weekday(), byte(t.Hour()))
}

// SetHour as active or no
func (h *Hours) SetHour(weekDay time.Weekday, hour byte, active bool) {
	if h.TestHour(weekDay, hour) == active {
		return
	}

	if *h == nil {
		*h = make(Hours, 24)
	}

	if active {
		(*h)[hour] |= byte(0x01) << byte(weekDay)
	} else {
		(*h)[hour] &= ^(byte(0x01) << byte(weekDay))
	}
}

// MarshalJSON implements the functionality of json.Marshaler interface
func (h Hours) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON implements the functionality of json.Unmarshaler interface
func (h *Hours) UnmarshalJSON(data []byte) error {
	if string(data) == `"*"` || string(data) == `"`+ActiveWeekHoursString+`"` {
		*h = nil
		return nil
	}

	if bytes.HasPrefix(data, []byte{'"'}) {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		data = []byte(s)
	}

	newHours, err := HoursByString(string(data))
	if err != nil {
		return err
	}
	*h = newHours
	return nil
}

// MarshalYAML implements the functionality of yaml.Marshaler interface
func (h Hours) MarshalYAML() (any, error) {
	return h.String(), nil
}

// UnmarshalYAML implements the functionality of yaml.Unmarshaler interface
func (h *Hours) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	if s == "*" || s == ActiveWeekHoursString {
		*h = nil
		return nil
	}

	newHours, err := HoursByString(s)
	if err != nil {
		return err
	}
	*h = newHours
	return nil
}

// Clone returns a copy of Hours
func (h Hours) Clone() Hours {
	if h == nil {
		return nil
	}
	newHours := make(Hours, len(h))
	copy(newHours, h)
	return newHours
}

var (
	_ json.Marshaler   = (Hours)(nil)
	_ json.Unmarshaler = (*Hours)(nil)
	_ yaml.Marshaler   = (Hours)(nil)
	_ yaml.Unmarshaler = (*Hours)(nil)
	_ driver.Valuer    = (Hours)(nil)
	_ sql.Scanner      = (*Hours)(nil)
)
