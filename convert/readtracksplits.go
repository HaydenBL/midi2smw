package convert

import (
	"bufio"
	"fmt"
	"midi2smw/midi"
	"os"
	"strings"
)

type TrackSplitMap map[int][]NoteGroup // Maps a track index to an array of note grouping

type NoteGroup struct {
	Notes []uint8
}

func SpecifyTrackSplits(midiTracks []midi.Track) TrackSplitMap {
	sc := bufio.NewScanner(os.Stdin)

	trackSplits := make(TrackSplitMap)

	for true {
		index := promptToSplitTracks(sc, midiTracks)
		if index == -1 {
			break
		}

		noteGroups := readTrackGroups(sc)
		if len(noteGroups) > 0 {
			trackSplits[index] = noteGroups
		}
	}

	return trackSplits

}

func promptToSplitTracks(sc *bufio.Scanner, midiTracks []midi.Track) int {
	for true {
		fmt.Println("-- Specify index of track to split (q to quit)")
		for i, track := range midiTracks {
			fmt.Printf("\t%d -\tName:\t\t\t%s\n", i, track.Name)
			fmt.Printf("\t\tInstrument:\t\t%s\n", track.Instrument)
			fmt.Printf("\t\tEvents:\t\t\t%d\n\n", len(track.Events))
		}
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "q" {
			return -1
		}

		if index, err := readInt(line); err != nil {
			fmt.Println(err)
		} else if int(index) > len(midiTracks)-1 {
			fmt.Println("Index out of range")
		} else {
			return int(index)
		}
	}
	return -1
}

func readTrackGroups(sc *bufio.Scanner) []NoteGroup {
	noteGroups := make([]NoteGroup, 0)

	for true {
		var notes []uint8
		var err error

		fmt.Printf("\tInput note group, space separated (q to finish): ")
		sc.Scan()
		line := sc.Text()

		if strings.ToLower(line) == "q" {
			break
		}

		if notes, err = readLineOfUInt8s(line); err != nil {
			fmt.Printf("\t\t%s\n", err)
		} else if numberAlreadyInAGroup(noteGroups, notes) {
			fmt.Printf("\t\tOne or more specified numbers already exists in a group for track\n")
		} else if containsDuplicates(notes) {
			fmt.Printf("\t\tTrack group cannot have duplicates\n")
		} else {
			noteGroups = append(noteGroups, NoteGroup{Notes: notes})
		}
	}

	return noteGroups
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
