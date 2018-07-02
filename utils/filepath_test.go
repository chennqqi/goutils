package utils

import (
	"fmt"
	"testing"
)

func Test_CleanFileName(t *testing.T) {
	tests := []string{
		`../../../../../tmp`,
		`..\..\..\..\abc`,
		`./../zzz`,
	}
	for i := 0; i < len(tests); i++ {
		n := CleanFileName("data", tests[i])
		fmt.Println(tests[i], "[clean]->", n)
		t.Log(n)
	}
}
