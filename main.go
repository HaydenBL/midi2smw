package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/midi"
	"midi2smw/write"
)

func main() {
	//begin()
	bah, _ := specifyDrumTracks()
	fmt.Println(bah[0].trackNumber)
}

func begin() {
	filename := "dean_town.mid"

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiFile, err := midi.Parse(filename)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", filename)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiFile)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	write.AllTracks(tracks, midiFile.Bpm)

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}
