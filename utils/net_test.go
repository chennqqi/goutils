package utils

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnectIP(t *testing.T) {
	type test struct {
		Proto  string
		Addr   string
		Expect string
	}
	var tests [4]test
	ip, err := GetHostIP()
	assert.Nil(t, err)
	assert.NotEqual(t, "", ip)
	t.Log("HOST IP:", ip)

	tip := fmt.Sprintf("%v:0", ip)
	conn0, err := net.Listen("tcp", tip)
	assert.Nil(t, err)
	defer conn0.Close()
	tests[0].Proto = "tcp"
	tests[0].Addr = conn0.Addr().String()
	tests[0].Expect = ip

	tip1 := fmt.Sprintf("%v:8081", ip)
	conn1, err := net.ListenPacket("udp", tip1)
	fmt.Println("tip1", tip1, err)
	assert.Nil(t, err)
	defer conn1.Close()
	tests[1].Proto = "udp"
	tests[1].Addr = conn1.LocalAddr().String()
	tests[1].Expect = ip

	conn2, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:0"))
	assert.Nil(t, err)
	defer conn2.Close()
	tests[2].Proto = "udp"
	tests[2].Addr = conn2.Addr().String()
	tests[2].Expect = "127.0.0.1"

	conn3, err := net.ListenPacket("udp", fmt.Sprintf("127.0.0.1:8081"))
	assert.Nil(t, err)
	defer conn3.Close()
	tests[3].Proto = "udp"
	tests[3].Addr = conn3.LocalAddr().String()
	tests[3].Expect = "127.0.0.1"

	for i := 0; i < len(tests); i++ {
		te := &tests[i]
		result, e := GetLocalConnectIP(te.Proto, te.Addr)
		assert.Nil(t, e)
		assert.Equal(t, te.Expect, result)
	}
}
