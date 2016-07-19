package main

import (
	"fmt"
)

func WildcardCmp(txt string, wild string) bool {
	wildLen := len(wild)
	txtLen := len(txt)
	wildIdx := 0
	txtIdx := 0

	for txtIdx < txtLen && (wild[wildIdx] != byte('*')) {
		if (wild[wildIdx] != txt[txtIdx]) && (wild[wildIdx] != '?') {
			return false
		}
		wildIdx++
		txtIdx++
	}

	var mpIdx int
	var cpIdx int

	for txtIdx < txtLen {
		if wild[wildIdx] == byte('*') {
			wildIdx++
			if wildIdx == wildLen {
				return true
			}
			mpIdx = wildIdx
			cpIdx = txtIdx + 1
		} else if wild[wildIdx] == txt[txtIdx] || wild[wildIdx] == '?' {
			wildIdx++
			txtIdx++
		} else {
			wildIdx = mpIdx
			txtIdx = cpIdx
		}
	}

	for wild[wildIdx] == '*' {
		wildIdx++
	}
	return wildIdx == wildLen
}

func main() {
	fmt.Println(WildcardCmp("aaaaa", "a*"))
	fmt.Println(WildcardCmp("aaaaa", "a?"))
	fmt.Println(WildcardCmp("aaaaa", "a????"))
}
