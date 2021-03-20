package drumtrack

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Group struct {
	TrackNumber uint8
	NoteGroups  []NoteGroup
}

type NoteGroup struct {
	Notes  []uint8
	Sample int16
}

func SpecifyDrumTrackGroups() ([]Group, error) {
	sc := bufio.NewScanner(os.Stdin)
	var drumTrackGroups = make([]Group, 0)
	var err error

	var drumTracks []uint8
	for true {
		fmt.Printf("Input drum tracks, space-separated (leave blank if none): ")
		sc.Scan()
		line := sc.Text()
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
		var newDtg Group
		if newDtg, err = readDrumTrackGroup(sc, trackNum); err != nil {
			return drumTrackGroups, err
		}
		if newDtg.NoteGroups != nil {
			drumTrackGroups = append(drumTrackGroups, newDtg)
		}
	}

	return drumTrackGroups, nil
}

func readDrumTrackGroup(sc *bufio.Scanner, trackNum uint8) (Group, error) {
	var dtg = Group{TrackNumber: trackNum}

	for true {
		var notes []uint8
		var err error

		fmt.Printf("\tInput note group for track %d, space separated (q to finish): ", trackNum)
		sc.Scan()
		line := sc.Text()

		if strings.ToLower(line) == "q" {
			break
		}

		if notes, err = readLineOfInts(line); err != nil {
			fmt.Printf("\t\t%s\n", err)
		} else if numberAlreadyInAGroup(dtg.NoteGroups, notes) {
			fmt.Printf("\t\tOne or more specified numbers already exists in a group for track %d\n", trackNum)
		} else if containsDuplicates(notes) {
			fmt.Printf("\t\tTrack group cannot have duplicates\n")
		} else {
			dtg.NoteGroups = append(dtg.NoteGroups, NoteGroup{Notes: notes})
		}
	}

	for i, ng := range dtg.NoteGroups {
		if promptToSetSamples(sc, ng.Notes) {
			var err error
			if dtg.NoteGroups[i], err = setSample(sc, ng); err != nil {
				return Group{}, err
			}
		} else {
			dtg.NoteGroups[i].Sample = -1
		}
	}

	return dtg, nil
}

func promptToSetSamples(sc *bufio.Scanner, notes []uint8) bool {
	for true {
		fmt.Printf("\t\tSet samples for note group:")
		for _, note := range notes {
			fmt.Printf(" %d", note)
		}
		fmt.Printf("? (y/n): ")
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "y" {
			return true
		} else if strings.ToLower(line) == "n" {
			return false
		}
	}
	return false
}

func setSample(sc *bufio.Scanner, noteGroup NoteGroup) (NoteGroup, error) {
	sample, err := setSampleForTrackGroup(sc, noteGroup.Notes)
	if err != nil {
		return NoteGroup{}, err
	}
	noteGroup.Sample = int16(sample)
	return noteGroup, nil
}

func setSampleForTrackGroup(sc *bufio.Scanner, notes []uint8) (uint8, error) {
	for true {
		fmt.Printf("\t\tEnter sample number for note group:")
		for _, note := range notes {
			fmt.Printf(" %d", note)
		}
		fmt.Printf(": ")
		sc.Scan()
		line := sc.Text()

		if sample, err := strconv.ParseUint(line, 10, 8); err == nil {
			return uint8(sample), nil
		}
	}
	return 0, errors.New("error getting sample track for group")
}
