package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRoute(t *testing.T) {
	r, e := ListGateWay()
	assert.Nil(t, e)
	for _, v := range r {
		t.Log(v)
	}
}
