package ginutils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var ErrNotExist = errors.New("Not Exist")

func GetQueryInt(c *gin.Context, q string) (int, error) {
	var v int
	queryValue, exist := c.GetQuery(q)
	if !exist {
		return 0, ErrNotExist
	}
	_, err := fmt.Sscanf(queryValue, "%d", &v)
	return v, err
}

func GetQueryInt64(c *gin.Context, q string) (int64, error) {
	var v int64
	queryValue, exist := c.GetQuery(q)
	if !exist {
		return 0, ErrNotExist
	}
	_, err := fmt.Sscanf(queryValue, "%d", &v)
	return v, err
}

func GetQueryfloat(c *gin.Context, q string) (float32, error) {
	var v float32

	queryValue, exist := c.GetQuery(q)
	if !exist {
		return v, ErrNotExist
	}
	_, err := fmt.Sscanf(queryValue, "%f", &v)
	return v, err
}

func GetQueryfloat64(c *gin.Context, q string) (float64, error) {
	var v float64
	queryValue, exist := c.GetQuery(q)
	if !exist {
		return v, ErrNotExist
	}
	_, err := fmt.Sscanf(queryValue, "%f", &v)
	return v, err
}
