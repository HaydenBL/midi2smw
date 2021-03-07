package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/midiparse"
	"midi2smw/write"
)

func main() {
	beginParsing()
}

func beginParsing() {
	filename := "dean_town.mid"

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiTracks, err := midiparse.Parse(filename)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", filename)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiTracks)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	write.AllTracks(tracks)

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}
