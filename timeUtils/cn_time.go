package timeUtils

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

var (
	CstFormat = "2006-01-02 15:04:05"
)

type CstTime time.Time

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t CstTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(t).Format(CstFormat))
	return []byte(str), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *CstTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if len(data) == 0 || string(data) == "null" || string(data) == "" || string(data) == `""` {
		return nil
	}
	data = bytes.Trim(data, "\"")
	// Fractional seconds are handled implicitly by Parse.
	cst, err := time.ParseInLocation(`"`+CstFormat+`"`, string(data), time.Local)
	if err != nil {
		if strings.Contains(err.Error(), "cannot parse") {
			cst, err = time.ParseInLocation(CstFormat, string(data), time.Local)
		}
	}
	*t = CstTime(cst)
	return err
}
