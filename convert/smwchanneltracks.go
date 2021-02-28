package convert

import (
	"fmt"
)

type SmwNote struct {
	key    string
	length uint8
	octave int
}

var noteDict = map[int]string{
	0:  "c",
	1:  "c+",
	2:  "d",
	3:  "d+",
	4:  "e",
	5:  "f",
	6:  "f+",
	7:  "g",
	8:  "g+",
	9:  "a",
	10: "a+",
	11: "b",
}

// SMW tempo conversion formula: BPM * (256/625)
func createSmwChannelTracks(noteTracks []noteTrack) [][]SmwNote {
	var smwTracks [][]SmwNote
	convertToSmwNoteLength := noteLengthConverter()
	for _, noteTrack := range noteTracks {
		var smwTrack []SmwNote
		lastStartTime := uint32(0)
		lastDuration := uint32(0)
		for j, note := range noteTrack.Notes {
			if j != 0 && note.StartTime == lastStartTime {
				continue // temporary way of dealing with chords
			}
			key, octave := noteValueToKey(note.Key)
			if restLength := getRestLength(note, lastStartTime, lastDuration, convertToSmwNoteLength); restLength != 0 {
				smwTrack = append(smwTrack, SmwNote{"r", restLength, octave})
			}
			length := convertToSmwNoteLength(note.Duration)
			smwNote := SmwNote{key, length, octave}
			if note.StartTime != lastStartTime {
				smwTrack = append(smwTrack, smwNote)
			}
			lastStartTime = note.StartTime
			lastDuration = note.Duration
		}
		smwTracks = append(smwTracks, smwTrack)
	}
	return smwTracks
}

func getRestLength(note midiNote, lastStartTime uint32, lastDuration uint32, converter func(duration uint32) uint8) uint8 {
	lastEndTime := lastStartTime + lastDuration
	restGapDuration := note.StartTime - lastEndTime
	restGapSmwLength := converter(restGapDuration)
	return restGapSmwLength
}

func noteValueToKey(noteValue uint8) (key string, octave int) {
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	if noteValue < 19 || noteValue > 88 {
		fmt.Printf("Error, invalid note value: %d\n", noteValue)
		return "", -1
	}

	key = noteDict[int(noteValue%12)]
	octave = int(noteValue/12) - 1
	return key, octave
}

func noteLengthConverter() func(duration uint32) uint8 {
	// TODO - Need to work out the math, but for the midi I'm working with,
	//  60 ticks is 1/32 note
	ticksPer32ndNote := 60.0
	constant := (1.0 / 32.0) / ticksPer32ndNote
	note32nd := constant * ticksPer32ndNote

	//note64th := note32nd / 2
	note16th := note32nd * 2
	note8th := note32nd * 4
	noteQuarter := note32nd * 8
	noteHalf := note32nd * 16
	noteWhole := note32nd * 32

	return func(duration uint32) uint8 {
		noteLength := float64(duration) * constant
		if noteLength <= note32nd {
			if round(noteLength, 0, note32nd) == 1 {
				return 0
			}
			return 32
		}
		if noteLength <= note16th {
			if round(noteLength, note32nd, note16th) == 1 {
				return 32
			}
			return 16
		}
		if noteLength <= note8th {
			if round(noteLength, note16th, note8th) == 1 {
				return 16
			}
			return 8
		}
		if noteLength <= noteQuarter {
			if round(noteLength, note8th, noteQuarter) == 1 {
				return 8
			}
			return 4
		}
		if noteLength <= noteHalf {
			if round(noteLength, noteQuarter, noteHalf) == 1 {
				return 4
			}
			return 2
		}
		if noteLength <= noteWhole {
			if round(noteLength, noteHalf, noteWhole) == 1 {
				return 2
			}
			return 1
		}

		return 1 // default to whole note if longer?
	}
}

// returns 1 if `smaller` is closer to n, or 2 if `bigger` is closer to n
func round(n, smaller, bigger float64) int {
	if smaller > bigger {
		fmt.Printf("Error: %f > %f", smaller, bigger)
		return 0
	}
	if n-smaller < bigger-n {
		return 1
	}
	return 2
}
