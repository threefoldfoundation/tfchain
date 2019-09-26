package types

import (
	"reflect"
	"strings"
	"testing"
)

func TestSortedSetDifference(t *testing.T) {
	arrA := []string{"a", "b", "d", "e", "g"}
	arrB := []string{"a", "c", "d", "f", "h", "z"}
	indices := sortedSetDifference(len(arrA), len(arrB), func(idxA, idxB int) int {
		return strings.Compare(arrA[idxA], arrB[idxB])
	})
	expectedIndices := []int{1, 3, 4}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Error("unexpected result for 'a \\ b':", expectedIndices, "!=", indices)
	}
	indices = sortedSetDifference(len(arrB), len(arrA), func(idxA, idxB int) int {
		return strings.Compare(arrB[idxA], arrA[idxB])
	})
	expectedIndices = []int{1, 3, 4, 5}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Error("unexpected result for 'b \\ a':", expectedIndices, "!=", indices)
	}
}

func TestSortedSetDifference_Empty(t *testing.T) {
	arrA := []string{}
	arrB := []string{}
	indices := sortedSetDifference(len(arrA), len(arrB), func(idxA, idxB int) int {
		return strings.Compare(arrA[idxA], arrB[idxB])
	})
	if len(indices) != 0 {
		t.Fatal("unexpected sorted set difference indices:", indices)
	}
}

func TestSortedSetDifference_EmptyResult(t *testing.T) {
	arrA := []string{"a", "c"}
	arrB := []string{"a", "b", "c", "d", "e"}
	indices := sortedSetDifference(len(arrA), len(arrB), func(idxA, idxB int) int {
		return strings.Compare(arrA[idxA], arrB[idxB])
	})
	if len(indices) != 0 {
		t.Fatal("unexpected sorted set difference indices:", indices)
	}
}

func TestSortedSetIntersection(t *testing.T) {
	arrA := []string{"a", "b", "d", "e", "g"}
	arrB := []string{"a", "c", "d", "f", "h", "z"}
	indices := sortedSetIntersection(len(arrA), len(arrB), func(idxA, idxB int) int {
		return strings.Compare(arrA[idxA], arrB[idxB])
	})
	expectedIndices := []int{0, 2}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Error("unexpected result for 'intersection(a, b)':", expectedIndices, "!=", indices)
	}
	indices = sortedSetIntersection(len(arrB), len(arrA), func(idxA, idxB int) int {
		return strings.Compare(arrB[idxA], arrA[idxB])
	})
	expectedIndices = []int{0, 2}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Error("unexpected result for 'intersection(b, a)':", expectedIndices, "!=", indices)
	}
}

func TestSortedSetIntersection_Empty(t *testing.T) {
	arrA := []string{}
	arrB := []string{}
	indices := sortedSetIntersection(len(arrA), len(arrB), func(idxA, idxB int) int {
		return strings.Compare(arrA[idxA], arrB[idxB])
	})
	if len(indices) != 0 {
		t.Fatal("unexpected sorted set intersection indices:", indices)
	}
	indices = sortedSetIntersection(len(arrB), len(arrA), func(idxA, idxB int) int {
		return strings.Compare(arrB[idxA], arrA[idxB])
	})
	if len(indices) != 0 {
		t.Fatal("unexpected sorted set intersection indices:", indices)
	}
}

func TestSortedSetIntersection_EmptResulty(t *testing.T) {
	arrA := []string{"a", "c", "e"}
	arrB := []string{"b", "d"}
	indices := sortedSetIntersection(len(arrA), len(arrB), func(idxA, idxB int) int {
		return strings.Compare(arrA[idxA], arrB[idxB])
	})
	if len(indices) != 0 {
		t.Fatal("unexpected sorted set intersection indices:", indices)
	}
	indices = sortedSetIntersection(len(arrB), len(arrA), func(idxA, idxB int) int {
		return strings.Compare(arrB[idxA], arrA[idxB])
	})
	if len(indices) != 0 {
		t.Fatal("unexpected sorted set intersection indices:", indices)
	}
}
