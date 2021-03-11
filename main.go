package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/midi"
	"midi2smw/write"
)

func main() {
	beginParsing()
}

func beginParsing() {
	filename := "silver_street.mid"
	var ticksPer64thNote uint32 = 10 // hardcoding for the track I'm working with, figure this out later

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiTracks, err := midi.Parse(filename)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", filename)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiTracks, ticksPer64thNote)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	write.AllTracks(tracks)

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}
