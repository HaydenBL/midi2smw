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
	noteValues  []uint8
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
		var input DrumTrackGroup
		if input, err = readDrumTrackGroup(scanner, trackNum); err != nil {
			return drumTrackGroups, err
		}
		drumTrackGroups = append(drumTrackGroups, input)
	}

	return drumTrackGroups, nil
}

func readDrumTrackGroup(scanner *bufio.Scanner, trackNum uint8) (DrumTrackGroup, error) {
	var dtg = DrumTrackGroup{trackNumber: trackNum}

	for true {
		var notes []uint8
		var err error

		fmt.Printf("\tInput note group for track %d, space separated: ", trackNum)
		scanner.Scan()
		line := scanner.Text()

		if notes, err = readLineOfInts(line); err != nil {
			fmt.Println(err)
		} else {
			dtg.noteValues = notes
			break
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
	nums = removeDuplicates(nums)
	return nums, nil
}

func removeDuplicates(nums []uint8) []uint8 {
	seen := make(map[uint8]bool, len(nums))
	j := uint8(0)
	for _, v := range nums {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = true
		nums[j] = v
		j++
	}
	return nums[:j]
}
