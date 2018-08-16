package re

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpGroup(t *testing.T) {
	exp := regexp.MustCompile(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`)
	group := RegexpGroup(exp, "1234.5678.9")
	assert.Equal(t, "1234", group["first"])
	assert.Equal(t, "9", group["second"])
	assert.Equal(t, 2, len(group))
}
