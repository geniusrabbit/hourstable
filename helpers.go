package hourstable

import (
	"bytes"
	"strings"
	"time"
)

// Default constants...
const (
	ActiveDayHoursString   = "111111111111111111111111"
	DisabledDayHoursString = "000000000000000000000000"
	AllActiveHoursString   = "*" // default value to save the space
	ActiveWeekHoursString  = "" +
		ActiveDayHoursString + // Sunday
		ActiveDayHoursString + // Monday
		ActiveDayHoursString + // Tuesday
		ActiveDayHoursString + // Wednesday
		ActiveDayHoursString + // Thursday
		ActiveDayHoursString + // Friday
		ActiveDayHoursString // Saturday
)

func hoursToBinary(timetable Hours, hours string, dayOfWeek time.Weekday) {
	// All hours is on
	if hours == AllActiveHoursString {
		for i := 0; i < 24; i++ {
			timetable[i] |= byte(0x01) << byte(dayOfWeek)
		}
		return
	}
	// Erace all data if empty values
	if hours == "" {
		for i := 0; i < 24; i++ {
			timetable[i] &= ^(byte(0x01) << byte(dayOfWeek))
		}
		return
	}
	for i, c := range hours {
		if c == '1' {
			timetable[i] |= byte(0x01) << byte(dayOfWeek)
		} else {
			timetable[i] &= ^(byte(0x01) << byte(dayOfWeek))
		}
	}
}

func binaryToHours(timetable Hours, dayOfWeek time.Weekday) string {
	var buff bytes.Buffer
	for _, hour := range timetable {
		if hour&(byte(0x01)<<byte(dayOfWeek)) != 0 {
			buff.WriteByte('1')
		} else {
			buff.WriteByte('0')
		}
	}
	return buff.String()
}

func binaryToHoursShort(timetable Hours, dayOfWeek time.Weekday) string {
	var (
		shortAll  = true
		shortNone = true
	)
	for _, hour := range timetable {
		if hour&(byte(0x01)<<byte(dayOfWeek)) == 0 {
			shortAll = false
		} else {
			shortNone = false
		}
		if !shortAll && !shortNone {
			break
		}
	}
	if shortAll {
		return AllActiveHoursString
	}
	if shortNone {
		return ""
	}
	return binaryToHours(timetable, dayOfWeek)
}

// ActiveHoursRangeString returns preformatted string with marked active houts according to range
func ActiveHoursRangeString(from, to byte) string {
	if from <= 0 && to >= 23 {
		return ActiveDayHoursString
	}
	if from > 23 || from >= to || (from <= 0 && to <= 0) {
		return DisabledDayHoursString
	}

	var buff bytes.Buffer
	if from > 0 {
		buff.WriteString(strings.Repeat("0", int(from)))
	}
	if to > 0 {
		buff.WriteString(strings.Repeat("1", int(to-from)))
	}
	if to < 23 {
		buff.WriteString(strings.Repeat("0", int(24-to)))
	}

	return buff.String()
}
