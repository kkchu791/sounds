package utils

import (
	"reflect"
	"testing"
)

// go test ./internal/domain/utils/ -run TestAssignReplicasToBrokers -v
func TestAssignReplicasToBrokers(t *testing.T) {
	result := AssignReplicasToBrokers(10, 3, 5, -1, -1)

	expected := map[int][]int{
		0: {0, 1, 2},
		1: {1, 2, 3},
		2: {2, 3, 4},
		3: {3, 4, 0},
		4: {4, 0, 1},
		5: {0, 2, 3},
		6: {1, 3, 4},
		7: {2, 4, 0},
		8: {3, 0, 1},
		9: {4, 1, 2},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}
