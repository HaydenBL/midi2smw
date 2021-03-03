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

func createSmwChannelTracks(noteTracks []noteTrack, ticksPer64thNote uint32) [][]SmwNote {
	var smwTracks [][]SmwNote
	convertToSmwNoteLength := noteLengthConverter(ticksPer64thNote)
	for _, noteTrack := range noteTracks {
		var smwTrack []SmwNote
		for _, note := range noteTrack.Notes {
			key, octave := noteValueToKey(note)
			lengths := convertToSmwNoteLength(note.Duration)
			smwNote := SmwNote{key, lengths, octave}
			smwTrack = append(smwTrack, smwNote)
		}
		smwTracks = append(smwTracks, smwTrack)
	}
	return smwTracks
}

func noteValueToKey(note midiNote) (key string, octave int) {
	noteValue := note.Key
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	if note.isRest {
		return "r", 0
	}
	if noteValue < 19 || noteValue > 88 {
		fmt.Printf("Error, invalid note value: %d\n", noteValue)
		return "", -1
	}
	key = noteDict[int(noteValue%12)]
	octave = int(noteValue/12) - 1
	return key, octave
}

func noteLengthConverter(ticksPer64thNote uint32) func(duration uint32) (lengths []uint8) {
	ticksPer32ndNote := ticksPer64thNote * 2
	ticksPer16thNote := ticksPer32ndNote * 2
	ticksPerQuarterNote := ticksPer16thNote * 2
	ticksPerHalfNote := ticksPerQuarterNote * 2
	ticksPerWholeNote := ticksPerHalfNote * 2

	noteLengthToSmwLength := func(duration uint32) (uint8, uint32) {
		if duration <= ticksPer64thNote {
			return 64, ticksPer64thNote - duration
		}
		if duration <= ticksPer32ndNote {
			return 32, ticksPer32ndNote - duration
		}
		if duration <= ticksPer16thNote {
			return 16, ticksPer16thNote - duration
		}
		if duration <= ticksPerQuarterNote {
			return 4, ticksPerQuarterNote - duration
		}
		if duration <= ticksPerHalfNote {
			return 2, ticksPerHalfNote - duration
		}
		if duration <= ticksPerWholeNote {
			return 1, ticksPerWholeNote - duration
		}
		return 1, duration - ticksPerWholeNote
	}

	return func(duration uint32) []uint8 {
		lengths := make([]uint8, 0)
		length, remainder := noteLengthToSmwLength(duration)
		lengths = append(lengths, length)
		for remainder != 0 {
			length, remainder = noteLengthToSmwLength(remainder)
			lengths = append(lengths, length)
		}
		return lengths
	}
}
