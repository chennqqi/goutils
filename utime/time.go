package utime

import (
	"errors"
	"time"
)

var (
	DefaultFormat = time.RFC3339Nano
)

type Duration time.Duration
type Time time.Time

func (c Duration) String() string {
	return time.Duration(c).String()
}

func (c Duration) Duration() stdtime.Duration {
	return stdtime.Duration(c)
}

func (c Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(c).String(), nil
}

func (c *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	to, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*c = Duration(to)
	return err
}

func (c Duration) MarshalJSON() ([]byte, error) {
	s := c.String()
	//fmt.Println("Marshal:", s)
	return []byte(`"` + s + `"`), nil
}

func (c *Duration) UnmarshalJSON(raw []byte) error {
	if len(raw) < 3 {
		return errors.New("No data")
	}
	if raw[0] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	to, err := time.ParseDuration(string(raw))
	if err != nil {
		return nil
	}
	*c = Duration(to)
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte(`"null"`), nil
	}
	b := make([]byte, 0, len(DefaultFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, DefaultFormat)
	b = append(b, '"')
	return b, nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 1 {
		return errors.New("No data")
	}
	if data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	if string(data) == "null" {
		*t = Time(time.Unix(0, 0))
		return nil
	}

	now, err := time.ParseInLocation(DefaultFormat, string(data), time.Local)
	*t = Time(now)
	return err
}

func (t *Time) FromTime(tt time.Time) {
	*t = Time(tt)
}
