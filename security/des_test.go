package security

import (
	"testing"
)

func Test3des(t *testing.T) {
	testdata := "abcdefghijklmnopqrstuvwxyz01234567890"
	testkey := "testkey....."
	dat, err := TripleEcbDesEncrypt([]byte(testdata), []byte(testkey))
	if err != nil {
		t.Error(err)
		return
	}
	raw, err := TripleEcbDesDecrypt(dat, []byte(testkey))
	if err != nil {
		t.Error(err)
		return
	}
	if string(raw) == testdata {
		t.Log("OK")
	} else {
		t.Log("ERROR")
	}
}
