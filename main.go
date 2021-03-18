package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/drumtrack"
	"midi2smw/midi"
	"midi2smw/write"
)

func main() {
	begin()
}

func begin() {
	filename := "dean_town.mid"

	drumTrackGroups := getDrumTracksGroups()

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiFile, err := midi.Parse(filename)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", filename)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiFile, drumTrackGroups)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	write.AllTracks(tracks, midiFile.Bpm)

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}

func getDrumTracksGroups() []drumtrack.Group {
	var err error
	var bah []drumtrack.Group
	if bah, err = drumtrack.SpecifyDrumTrackGroups(); err != nil {
		fmt.Println(err)
	}
	return bah
}
