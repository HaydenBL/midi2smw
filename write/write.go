package write

import (
	"fmt"
	"midi2smw/convert"
)

func AllTracks(tracks []convert.SmwTrack) {
	for i, track := range tracks {
		fmt.Printf("---- Printing track %d\n\n", i)
		Track(track)
		fmt.Println()
	}
}

func Track(track convert.SmwTrack) {
	for i, channel := range track.ChannelTracks {
		fmt.Printf("-- Printing channel %d\n", i)
		Channel(channel)
		fmt.Println()
	}
}

func Channel(channel []convert.SmwNote) {
	lastOctave := channel[0].Octave
	fmt.Printf("Start octave: %d\n", lastOctave)
	for _, smwNote := range channel {
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
}
