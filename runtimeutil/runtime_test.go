package runtimeutil

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetfunctionName(t *testing.T) {
	// This will print "name: main.foo"
	assert.Equal(t, GetFunctionName(TestGetfunctionName, '/', '.'), "TestGetfunctionName")

	// runtime/debug.FreeOSMemory
	assert.Equal(t, GetFunctionName(debug.FreeOSMemory), "runtime/debug.FreeOSMemory")
	// FreeOSMemory
	assert.Equal(t, GetFunctionName(debug.FreeOSMemory, '.'), "FreeOSMemory")
	// FreeOSMemory
	assert.Equal(t, GetFunctionName(debug.FreeOSMemory, '/', '.'), "FreeOSMemory")
}
