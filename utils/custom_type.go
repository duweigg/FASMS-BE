package utils

import (
	"errors"
	"strings"
	"time"
)

// Custom Date type for parsing "YYYY-MM-DD"
type Date time.Time

const DateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	parsedTime, err := time.Parse(DateFormat, str)
	if err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}
	*d = Date(parsedTime)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	formatted := time.Time(d).Format(DateFormat)
	return []byte(`"` + formatted + `"`), nil
}

func (d Date) ToTime() time.Time {
	return time.Time(d)
}
