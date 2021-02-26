package convert

import (
	"fmt"
	"midi2smw/midiparse"
)

func Convert(midiTracks []midiparse.MidiTrack) {
	fmt.Println("Converting midi tracks...")
	filterEmptyTracks(midiTracks)
}

func filterEmptyTracks(tracks []midiparse.MidiTrack) []midiparse.MidiTrack {
	var nonEmptyTracks []midiparse.MidiTrack
	for _, track := range tracks {
		if len(track.Events) != 0 {
			nonEmptyTracks = append(nonEmptyTracks, track)
		}
	}

	fmt.Printf("Removed %d empty tracks\n", len(tracks)-len(nonEmptyTracks))
	return nonEmptyTracks
}
