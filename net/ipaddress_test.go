package net

import (
	stdnet "net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrivateAddr(t *testing.T) {
	var tests = []struct {
		ip     string
		expect bool
	}{
		{"127.0.0.5", true},
		{"127.0.0.1", true},
		{"100.2.5.1", false},
		{"10.10.10.10", true},
	}
	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		ip := stdnet.ParseIP(test.ip)
		assert.Equal(t, IsPrivateIP(ip), test.expect)
	}
}
