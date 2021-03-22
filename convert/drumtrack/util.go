package drumtrack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func readInt(str string) (uint8, error) {
	var num64 uint64
	var err error

	if num64, err = strconv.ParseUint(str, 10, 8); err != nil {
		return 0, errors.New("error parsing line")
	}
	if num64 > 255 {
		return 0, errors.New("number too large (max 255)")
	}
	return uint8(num64), nil
}

func readLineOfInts(str string) ([]uint8, error) {
	var num64 uint64
	var err error

	numStrings := strings.Split(str, " ")
	var nums = make([]uint8, 0)
	for _, numStr := range numStrings {
		if num64, err = strconv.ParseUint(numStr, 10, 8); err != nil {
			return []uint8{}, errors.New("error parsing line")
		}
		if num64 > 255 {
			return []uint8{}, errors.New(fmt.Sprintf("number %d too large (max 255)", num64))
		}
		nums = append(nums, uint8(num64))
	}
	return nums, nil
}

func numberAlreadyInAGroup(existingGroups []NoteGroup, newGroup []uint8) bool {
	allNums := make([]uint8, 0)
	for i := range existingGroups {
		for j := range existingGroups[i].Notes {
			allNums = append(allNums, existingGroups[i].Notes[j])
		}
	}
	for _, num := range newGroup {
		if numberExistsIn(num, allNums) {
			return true
		}
	}
	return false
}

func numberExistsIn(num uint8, arr []uint8) bool {
	for _, n := range arr {
		if n == num {
			return true
		}
	}
	return false
}

func containsDuplicates(arr []uint8) bool {
	seen := make(map[uint8]bool, len(arr))
	for _, v := range arr {
		if _, ok := seen[v]; ok {
			return true
		}
		seen[v] = true
	}
	return false
}
