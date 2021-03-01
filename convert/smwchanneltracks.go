package convert

import (
	"fmt"
)

type SmwNote struct {
	key          string
	lengthValues []uint8
	octave       int
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
			if restLengths := getRestLength(note.StartTime, lastStartTime, lastDuration, convertToSmwNoteLength); len(restLengths) > 0 && restLengths[0] != 0 {
				smwTrack = append(smwTrack, SmwNote{"r", restLengths, octave})
			}
			lengths := convertToSmwNoteLength(note.Duration)
			smwNote := SmwNote{key, lengths, octave}
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

func getRestLength(currentStartTime, lastStartTime uint32, lastDuration uint32, converter func(duration uint32) []uint8) []uint8 {
	lastEndTime := lastStartTime + lastDuration
	restGapDuration := currentStartTime - lastEndTime
	restGapSmwLengths := converter(restGapDuration)
	return restGapSmwLengths
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

func noteLengthConverter() func(duration uint32) (lengths []uint8) {
	// TODO - Need to work out the math, but for the midi I'm working with,
	//  60 ticks is 1/32 note
	ticksPer32ndNote := 60.0
	constant := (1.0 / 32.0) / ticksPer32ndNote
	note32nd := constant * ticksPer32ndNote

	note64th := note32nd / 2
	note16th := note32nd * 2
	note8th := note32nd * 4
	noteQuarter := note32nd * 8
	noteHalf := note32nd * 16
	noteWhole := note32nd * 32

	noteLengthToSmwLength := func(noteLength float64) (uint8, float64) {
		if noteLength <= note64th {
			if round(noteLength, 0, note64th) == 1 {
				return 0, 0
			}
			return 64, note64th
		}
		if noteLength <= note32nd {
			if round(noteLength, note64th, note32nd) == 1 {
				return 64, note64th
			}
			return 32, note32nd
		}
		if noteLength <= note16th {
			if round(noteLength, note32nd, note16th) == 1 {
				return 32, note32nd
			}
			return 16, note16th
		}
		if noteLength <= note8th {
			if round(noteLength, note16th, note8th) == 1 {
				return 16, note16th
			}
			return 8, note8th
		}
		if noteLength <= noteQuarter {
			if round(noteLength, note8th, noteQuarter) == 1 {
				return 8, note8th
			}
			return 4, noteQuarter
		}
		if noteLength <= noteHalf {
			if round(noteLength, noteQuarter, noteHalf) == 1 {
				return 4, noteQuarter
			}
			return 2, noteHalf
		}
		if noteLength <= noteWhole {
			if round(noteLength, noteHalf, noteWhole) == 1 {
				return 2, noteHalf
			}
			return 1, noteWhole
		}

		return 1, noteWhole // default to whole note if longer?
	}

	return func(duration uint32) (lengths []uint8) {
		noteLength := float64(duration) * constant
		length, lengthTicks := noteLengthToSmwLength(noteLength)
		lengths = append(lengths, length)
		remainder := noteLength - lengthTicks
		length, lengthTicks = noteLengthToSmwLength(remainder)
		for ; length != 0; length, lengthTicks = noteLengthToSmwLength(remainder) {
			lengths = append(lengths, length)
			remainder = remainder - lengthTicks
		}
		return lengths
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
