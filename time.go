package kis

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime is a custom time type that allows for JSON unmarshalling.
type CustomTime struct {
	time.Time
}

// dateLayout is the layout for the date format. UTC datetime in ISO 8601 format (YYYY-MM-DDThh:mm:ssZ).
const dateLayout = "2006-01-02T15:04:05"

// UnmarshalJSON unmarshals a JSON string into a CustomTime type.
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || len(s) == 0 {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(dateLayout, s)
	if err != nil {
		panic(fmt.Errorf("error parsing time: %w", err))
	}
	return
}

// MarshalJSON marshals a CustomTime type into a JSON string.
func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(ct.Time.Format(dateLayout)), nil
}
