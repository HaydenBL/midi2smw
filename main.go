package main

import (
	"flag"
	"fmt"
	"midi2smw/convert"
	"midi2smw/midi"
	"midi2smw/write"
	"os"
)

func main() {
	fileName, drumTracksFlag := parseFlags()
	begin(fileName, drumTracksFlag)
}

func parseFlags() (fileName string, drumTracksFlag bool) {
	drumTracksFlagPtr := flag.Bool("drumtracks", false, "Specify drum tracks and note groupings")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Error: no file name provided")
		os.Exit(1)
	}
	return flag.Args()[0], *drumTracksFlagPtr
}

func begin(fileName string, drumTracksFlag bool) {

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiFile, err := midi.Parse(fileName)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", fileName)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiFile, drumTracksFlag)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	write.AllTracks(tracks, midiFile.Bpm)

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}
