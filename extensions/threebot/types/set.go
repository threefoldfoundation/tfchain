package types

import "strconv"

// sortedSetDifference allows you to take the difference of two (pre-)sorted sets,
// returning the indices of all elements that are in A, but not in B.
//
// cmp should return the result as: 0 if a==b, -1 if a < b, and +1 if a > b.
func sortedSetDifference(lenA, lenB int, cmp func(a int, b int) int) (indices []int) {
	var idxA, idxB, result int
	for idxA < lenA && idxB < lenB {
		switch result = cmp(idxA, idxB); result {
		case 0:
			idxA++
			idxB++
		case -1:
			indices = append(indices, idxA)
			idxA++
		case 1:
			idxB++
		default:
			panic("unexpected cmp result: " + strconv.Itoa(result))
		}
	}
	for ; idxA < lenA; idxA++ {
		indices = append(indices, idxA)
	}
	return
}

// sortedSetIntersection allows you to take the intersection of two (pre-)sorted sets,
// returning the indices of all elements that are in A AND in B.
//
// cmp should return the result as: 0 if a==b, -1 if a < b, and +1 if a > b.
func sortedSetIntersection(lenA, lenB int, cmp func(a int, b int) int) (indices []int) {
	var idxA, idxB, result int
	for idxA < lenA && idxB < lenB {
		switch result = cmp(idxA, idxB); result {
		case 0:
			indices = append(indices, idxA)
			idxA++
			idxB++
		case -1:
			idxA++
		case 1:
			idxB++
		default:
			panic("unexpected cmp result: " + strconv.Itoa(result))
		}
	}
	return
}
