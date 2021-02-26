package midiparse

import (
	"fmt"
	"os"
)

var (
	globalTempo uint32 = 0
	globalBPM   uint32 = 0
)

func Parse(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	defer file.Close()

	numTracks := parseHeader(file)

	var midiTracks []midiTrack
	for track := 0; track < int(numTracks); track++ {
		track := parseTrack(file)
		midiTracks = append(midiTracks, track)
	}
}
