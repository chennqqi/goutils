package safeurl

import (
	"testing"
)

func TestA(t *testing.T) {
	v, e := QueryUnescape("aaa%A")
	if e != nil {
		t.Error("safeurl unescape failed")
	}
	t.Log(v, "OK")
}
