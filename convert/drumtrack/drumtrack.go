package drumtrack

import (
	"bufio"
	"fmt"
	"midi2smw/midi"
	"os"
	"strconv"
	"strings"
)

type MidiTrackWithNoteGroups struct {
	midi.Track
	NoteGroups []NoteGroup
}

type NoteGroup struct {
	Notes  []uint8
	Sample int16
}

func SpecifyDrumTrackGroups(midiTracks []midi.Track) []MidiTrackWithNoteGroups {
	sc := bufio.NewScanner(os.Stdin)

	tracksWithNoteGroups := make([]MidiTrackWithNoteGroups, len(midiTracks))
	for i, track := range midiTracks {
		tracksWithNoteGroups[i].Track = track
	}

	for true {
		var index int

		index = promptToSplitTracks(sc, midiTracks)
		if index == -1 {
			break
		}

		var noteGroups []NoteGroup
		noteGroups = readDrumTrackGroups(sc)
		if len(noteGroups) > 0 {
			tracksWithNoteGroups[index].NoteGroups = noteGroups
		}
	}

	return tracksWithNoteGroups

}

func promptToSplitTracks(sc *bufio.Scanner, midiTracks []midi.Track) int {
	var index uint8
	var err error
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

		if index, err = readInt(line); err != nil {
			fmt.Println(err)
		} else if int(index) > len(midiTracks)-1 {
			fmt.Println("Index out of range")
		} else {
			return int(index)
		}
	}
	return -1
}

func readDrumTrackGroups(sc *bufio.Scanner) []NoteGroup {
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

		if notes, err = readLineOfInts(line); err != nil {
			fmt.Printf("\t\t%s\n", err)
		} else if numberAlreadyInAGroup(noteGroups, notes) {
			fmt.Printf("\t\tOne or more specified numbers already exists in a group for track\n")
		} else if containsDuplicates(notes) {
			fmt.Printf("\t\tTrack group cannot have duplicates\n")
		} else {
			noteGroups = append(noteGroups, NoteGroup{Notes: notes, Sample: -1})
		}
	}

	for i, ng := range noteGroups {
		if promptToSetSamples(sc, ng.Notes) {
			noteGroups[i] = setSample(sc, ng)
		}
	}

	return noteGroups
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

func setSample(sc *bufio.Scanner, noteGroup NoteGroup) NoteGroup {
	sample := getSampleForNoteGroup(sc, noteGroup.Notes)
	noteGroup.Sample = int16(sample)
	return noteGroup
}

func getSampleForNoteGroup(sc *bufio.Scanner, notes []uint8) uint8 {
	for true {
		fmt.Printf("\t\tEnter sample number for note group:")
		for _, note := range notes {
			fmt.Printf(" %d", note)
		}
		fmt.Printf(": ")
		sc.Scan()
		line := sc.Text()

		if sample, err := strconv.ParseUint(line, 10, 8); err == nil {
			return uint8(sample)
		}
	}
	fmt.Println("Error getting sample for note group")
	return 0
}
