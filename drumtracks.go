package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DrumTrackGroup struct {
	trackNumber uint8
	noteGroups  [][]uint8
}

func specifyDrumTracks() ([]DrumTrackGroup, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var drumTrackGroups = make([]DrumTrackGroup, 0)
	var err error

	var drumTracks []uint8
	for true {
		fmt.Printf("Input drum tracks, space-separated (leave blank if none): ")
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			return drumTrackGroups, nil
		}

		if drumTracks, err = readLineOfInts(line); err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	for _, trackNum := range drumTracks {
		var newDtg DrumTrackGroup
		if newDtg, err = readDrumTrackGroup(scanner, trackNum); err != nil {
			return drumTrackGroups, err
		}
		if newDtg.noteGroups != nil {
			drumTrackGroups = append(drumTrackGroups, newDtg)
		}
	}

	return drumTrackGroups, nil
}

func readDrumTrackGroup(scanner *bufio.Scanner, trackNum uint8) (DrumTrackGroup, error) {
	var dtg = DrumTrackGroup{trackNumber: trackNum}

	for true {
		var notes []uint8
		var err error

		fmt.Printf("\tInput note group for track %d, space separated (q to finish): ", trackNum)
		scanner.Scan()
		line := scanner.Text()

		if line == "q" || line == "Q" {
			break
		}

		if notes, err = readLineOfInts(line); err != nil {
			fmt.Printf("\t\t%s\n", err)
		} else if numberAlreadyInAGroup(dtg.noteGroups, notes) {
			fmt.Printf("\t\tOne or more specified numbers already exists in a group for track %d\n", trackNum)
		} else if containsDuplicates(notes) {
			fmt.Printf("\t\tTrack group cannot have duplicates\n")
		} else {
			dtg.noteGroups = append(dtg.noteGroups, notes)
		}
	}

	return dtg, nil
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

func numberAlreadyInAGroup(existingGroups [][]uint8, newGroup []uint8) bool {
	allNums := make([]uint8, 0)
	for i := range existingGroups {
		for j := range existingGroups[i] {
			allNums = append(allNums, existingGroups[i][j])
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
