package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/midiparse"
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

	fmt.Printf("========== BEGIN PRINTING ==========\n\n")

	for i, track := range tracks {
		fmt.Printf("--- PRINTING TRACK #%d ---\n", i)
		testPrint(track)
	}

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}

// temporary, just to test this thing
func testPrint(smwTrack []convert.SmwNote) {
	if len(smwTrack) == 0 {
		return
	}
	lastOctave := smwTrack[0].Octave
	fmt.Printf("Start octave: %d\n", lastOctave)
	for _, smwNote := range smwTrack {
		if smwNote.Key == "r" {
			for i, note := range smwNote.LengthValues {
				if i == 0 {
					fmt.Printf("r%d", note)
				} else {
					fmt.Printf("^%d", note)
				}
			}
		} else {
			if smwNote.Octave > lastOctave {
				for i := 0; i < smwNote.Octave-lastOctave; i++ {
					fmt.Printf(">")
				}
			} else if smwNote.Octave < lastOctave {
				for i := 0; i < lastOctave-smwNote.Octave; i++ {
					fmt.Printf("<")
				}
			}
			for i, note := range smwNote.LengthValues {
				if i == 0 {
					fmt.Printf("%s%d", smwNote.Key, note)
				} else {
					fmt.Printf("^%d", note)
				}
			}
			lastOctave = smwNote.Octave
		}
	}
	fmt.Println()
	fmt.Println()
}
