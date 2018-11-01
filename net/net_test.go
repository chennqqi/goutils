package net

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRoute(t *testing.T) {
	r, e := ListGateway()
	assert.Nil(t, e)
	for _, v := range r {
		t.Log(v)
	}
}

func TestGetGatewayByNic(t *testing.T) {
	ifs, err := net.Interfaces()
	assert.Nil(t, err)
	for i, ifi := range ifs {
		ip, err := GetGatewayByNic(ifi.Name)
		if err != ErrNotFound {
			assert.Nil(t, err)
			t.Log(i, ip.String())
		}
	}
}

func TestGetDefaultGateway(t *testing.T) {
	ip, err := GetDefaultGateway()
	assert.Nil(t, err)
	assert.NotEqual(t, ip.String(), "")
	assert.NotEqual(t, ip.String(), "0.0.0.0")
	t.Log("default gateway:", ip.String())
}
