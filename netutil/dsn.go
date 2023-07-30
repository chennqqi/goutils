package netutil

import (
	"fmt"
	"strings"
)

type Dsn struct {
	Scheme string
	Source string
}

func ParseDsn(dsn string) (Dsn, error) {
	params := strings.Split(dsn, "://")
	if len(params) != 2 {
		return Dsn{}, fmt.Errorf("bad dsn string")
	}
	return Dsn{
		params[0], params[1],
	}, nil
}
