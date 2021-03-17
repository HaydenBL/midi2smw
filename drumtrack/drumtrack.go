package drumtrack

import (
	"bufio"
	"fmt"
	"os"
)

type Group struct {
	TrackNumber uint8
	NoteGroups  [][]uint8
}

func SpecifyDrumTracks() ([]Group, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var drumTrackGroups = make([]Group, 0)
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
		var newDtg Group
		if newDtg, err = readDrumTrackGroup(scanner, trackNum); err != nil {
			return drumTrackGroups, err
		}
		if newDtg.NoteGroups != nil {
			drumTrackGroups = append(drumTrackGroups, newDtg)
		}
	}

	return drumTrackGroups, nil
}

func readDrumTrackGroup(scanner *bufio.Scanner, trackNum uint8) (Group, error) {
	var dtg = Group{TrackNumber: trackNum}

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
		} else if numberAlreadyInAGroup(dtg.NoteGroups, notes) {
			fmt.Printf("\t\tOne or more specified numbers already exists in a group for track %d\n", trackNum)
		} else if containsDuplicates(notes) {
			fmt.Printf("\t\tTrack group cannot have duplicates\n")
		} else {
			dtg.NoteGroups = append(dtg.NoteGroups, notes)
		}
	}

	return dtg, nil
}
