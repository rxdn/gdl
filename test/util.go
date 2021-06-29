package main

import "testing"

func MustMatch(t *testing.T, name string, got interface{}, want interface{}) {
	if want != got {
		t.Errorf("%s mismatch - expected: %v, got %v", name, want, got)
	}
}
