package model

import (
	"log"
	"testing"
)

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

func TestPartition_Append(t *testing.T) {
	// setup
	p := NewPartition()
	m := NewMessage("Blues Chord", "indoor-sounds")

	p.Append(m)

	if len(p.messages) != 1 {
		t.Fatalf("got length %v, want %v", len(p.messages), 1)
	}

	// check the result
	expected := m
	result := p.messages[0] // this now won't panic cuz of the above if check

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestPartition_ReadOneMessage(t *testing.T) {
	// setup
	p := NewPartition()
	m := NewMessage("ambulance siren", "outdoor-sounds")
	p.Append(m)
	offset := 0
	result, err := p.Read(offset)
	expected := m
	if err != nil {
		log.Println(err)
	}

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestPartition_Read_OutOfBounds(t *testing.T) {
	// setup
	p := NewPartition()
	offset := 0
	result, err := p.Read(offset)

	if err == nil {
		t.Fatalf("expected an error when reading from empty partition, got nil")
	}

	if result != nil {
		t.Errorf("expected nil message, got %v", result)
	}

	expectedErr := "hey this offset is out of range"
	if err.Error() != expectedErr {
		t.Errorf("got error %q, want %q", err.Error(), expectedErr)
	}
}
