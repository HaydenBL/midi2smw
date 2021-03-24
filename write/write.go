package write

import (
	"fmt"
	"math"
	"midi2smw/convert"
)

func AllTracks(tracks []convert.SmwTrack, bpm uint32) {
	for i, track := range tracks {
		fmt.Printf("---- Printing track %d", i)
		if track.Name != "" {
			fmt.Printf(" (%s)", track.Name)
		}
		fmt.Printf("\n\n")
		Track(track)
		fmt.Println()
	}
	fmt.Printf("SMW tempo: %d\n", bpmToSmwTempo(bpm))
}

func bpmToSmwTempo(bpm uint32) uint8 {
	const multiplier = float64(256) / 625
	tempo := math.Round(float64(bpm) * multiplier)
	return uint8(tempo)
}

func Track(track convert.SmwTrack) {
	for i, channel := range track.ChannelTracks {
		fmt.Printf("-- Printing channel %d\n", i)
		Channel(channel)
		fmt.Println()
	}
}

func Channel(channel convert.ChannelTrack) {
	notes := channel.Notes
	lastOctave := notes[0].Octave
	lastSample := channel.DefaultSample
	fmt.Printf("Start octave: %d\n", lastOctave)
	fmt.Printf("Sample: %d\n", channel.DefaultSample)

	for _, smwNote := range notes {
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
					// Check if we need to swap the sample
					sample, ok := channel.SampleMap[smwNote.KeyValue]
					if !ok {
						sample = channel.DefaultSample
					}
					if sample != lastSample {
						lastSample = sample
						fmt.Printf("@%d", sample)
					}

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
