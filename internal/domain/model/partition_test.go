package model

import "testing"

func TestPartition_Len(t *testing.T) {
	// setup
	expected := 0
	// call the thing
	p := NewPartition()
	result := p.Len()
	// check the result

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
