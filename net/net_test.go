package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRoute(t *testing.T) {
	r, e := ListGateWay()
	assert.Nil(t, e)
	for _, v := range v {
		t.Log(v)
	}
}
