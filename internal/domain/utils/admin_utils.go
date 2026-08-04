package utils

import "math/rand/v2"

func AssignReplicasToBrokers(pc, rf, bn, fixedStartIdx, startPId int) map[int][]int {
	ret := make(map[int][]int)
	var startIdx int
	if fixedStartIdx >= 0 {
		startIdx = fixedStartIdx
	} else {
		startIdx = rand.IntN(bn)
	}

	currPId := max(0, startPId)

	var nRShift int
	if fixedStartIdx >= 0 {
		nRShift = fixedStartIdx
	} else {
		nRShift = rand.IntN(bn)
	}

	for range pc {
		if (currPId > 0) && (currPId%bn == 0) {
			nRShift += 1
		}
		firstRIndex := (currPId + startIdx) % bn
		rBuf := []int{firstRIndex}

		for j := range rf - 1 {
			rBuf = append(rBuf, replicaIndex(firstRIndex, nRShift, j, bn))
		}

		ret[currPId] = rBuf
		currPId += 1
	}

	return ret

}

func replicaIndex(firstRIndex, secondRShift, rIndex, bn int) int {
	shift := 1 + (secondRShift+rIndex)%(bn-1)
	return (firstRIndex + shift) % bn
}
