package stringutil

import (
	"testing"
)

func testwildcard(parttern, txt string, t *testing.T) {
	if !WildcardCmp(txt, parttern) {
		t.Logf("%v %v %v", txt, "Match", parttern)
	} else {
		t.Logf("%v %v %v", txt, "Match", parttern)
	}
}

func testwildcardAssert(parttern, txt string, t *testing.T, assert bool) {
	r := WildcardCmp(txt, parttern)
	var s string
	if r {
		s = "MATCH"
	} else {
		s = "NOT MATCH"
	}
	if r != assert {
		t.Fatalf("AssertFails %s %s %s", txt, s, parttern)
	} else {
		t.Log("AssertOK", txt, s, parttern)
	}
}

func Test_wildcard(t *testing.T) {
	testwildcardAssert("g*ks", "geeks", t, true)           // Yes
	testwildcardAssert("ge?ks*", "geeksforgeeks", t, true) // Yes
	testwildcardAssert("abc*bcd", "abcdhghgbcd", t, true)  // Yes
	testwildcardAssert("*c*d", "abcd", t, true)            // Yes
	testwildcardAssert("*?c*d", "abcd", t, true)           // Yes

	testwildcardAssert("g*ks", "geeks", t, true)           // Yes
	testwildcardAssert("ge?ks*", "geeksforgeeks", t, true) // Yes
	testwildcardAssert("abc*bcd", "abcdhghgbcd", t, true)  // Yes
	testwildcardAssert("*c*d", "abcd", t, true)            // Yes
	testwildcardAssert("*?c*d", "abcd", t, true)           // Yes

	testwildcardAssert("g*k", "gee", t, false)      // No because 'k' is not in second
	testwildcardAssert("abc*c?d", "abcd", t, false) // No because second must have 2 instances of 'c'
	testwildcardAssert("g*k", "gee", t, false)      // No because 'k' is not in second
	testwildcardAssert("*pqrs", "pqrst", t, false)  // No because 't' is not in first
	testwildcardAssert("abc*c?d", "abcd", t, false) // No because second must have 2 instances of 'c'
	testwildcardAssert("*pqrs", "pqrst", t, false)  // No because 't' is not in first
}
