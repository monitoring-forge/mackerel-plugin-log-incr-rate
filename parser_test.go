package main

import "testing"

// test simpleCounter.Parse() and simpleCounter.Finish()
func TestSimpleCounter(t *testing.T) {
	lc := &simpleCounter{}
	if err := lc.Parse([]byte("test")); err != nil {
		t.Errorf("Parse() error = %v", err)
	}
	lc.Finish(10)
	if got := lc.GetTotal(); got != 1 {
		t.Errorf("GetTotal() = %v, want %v", got, 1)
	}
	if got := lc.GetDuration(); got != 10 {
		t.Errorf("GetDuration() = %v, want %v", got, 10)
	}
}
