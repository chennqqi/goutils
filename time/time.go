package time

import (
	"errors"
	stdtime "time"
)

type Duration stdtime.Duration

func (c Duration) MarshalYAML() (interface{}, error) {
	return stdtime.Duration(c).String(), nil
}

func (c *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	to, err := stdtime.ParseDuration(s)
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

func (c Duration) String() string {
	return stdtime.Duration(c).String()
}

func (c *Duration) UnmarshalJSON(raw []byte) error {
	if len(raw) < 3 {
		return errors.New("No data")
	}
	if raw[0] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	to, err := stdtime.ParseDuration(string(raw))
	if err != nil {
		return nil
	}
	*c = Duration(to)
	return nil
}
