package convert

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

func readLineOfUInt8s(str string) ([]uint8, error) {
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

func getIndexForNoteGroup(note MidiNote, noteGroups []NoteGroup) int {
	for i, noteGroup := range noteGroups {
		if contains(note.Key, noteGroup.Notes) {
			return i
		}
	}
	return len(noteGroups)
}

func contains(num uint8, arr []uint8) bool {
	for _, n := range arr {
		if num == n {
			return true
		}
	}
	return false
}

func findNoteIndex(notes []MidiNote, key uint8) (int, MidiNote) {
	for i, note := range notes {
		if note.Key == key {
			return i, note
		}
	}
	return -1, MidiNote{}
}

func deleteAtIndex(notes []MidiNote, i int) []MidiNote {
	copy(notes[i:], notes[i+1:])
	notes[len(notes)-1] = MidiNote{}
	notes = notes[:len(notes)-1]
	return notes
}

func minUint8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func maxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
