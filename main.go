package main

import (
	"fmt"
	"midi2smw/convert"
	"midi2smw/midiparse"
)

func main() {
	filename := "dean_town.mid"

	midiTracks, err := midiparse.Parse(filename)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", filename)
		return
	}

	convert.Convert(midiTracks)

}
