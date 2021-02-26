package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/midiparse"
)

func main() {
	filename := "dean_town.mid"

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiTracks, err := midiparse.Parse(filename)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", filename)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	convert.Convert(midiTracks)

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")

}
