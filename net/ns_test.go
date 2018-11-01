package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNs(t *testing.T) {
	ns, err := GetLocalNS()
	assert.Nil(t, err)
	assert.NotEqual(t, len(ns.NSRecord), 0)
}
