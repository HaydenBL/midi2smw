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

func createSmwChannelTracksForAllTracks(noteTracks []noteTrack, ticksPer64thNote uint32) [][]SmwNote {
	var smwTracks [][]SmwNote
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)
	for _, noteTrack := range noteTracks {
		smwTrack := createSmwChannelTrack(noteTrack.Notes, noteLengthConverter)
		smwTracks = append(smwTracks, smwTrack)
	}
	return smwTracks
}

func createSmwChannelTrack(notes []midiNote, noteLengthConverter func(uint32) []uint8) []SmwNote {
	var smwTrack []SmwNote
	var tick, lastNoteEndTime uint32
	var activeNote *midiNote
	for tick = 0; !trackDone(notes, tick); tick++ {
		if activeNote != nil {
			if tick != activeNote.StartTime+activeNote.Duration {
				continue
			} else {
				activeNote = nil
			}
		}
		activeNote = getNoteWithStartTime(notes, tick)
		if activeNote == nil {
			continue
		} else {
			// insert rest
			restLength := tick - lastNoteEndTime
			if restLength != 0 {
				lengths := noteLengthConverter(restLength)
				restSmwNote := SmwNote{key: "r", lengthValues: lengths, octave: 0}
				smwTrack = append(smwTrack, restSmwNote)
			}
		}
		key, octave := noteValueToSmwKey(*activeNote)
		lengths := noteLengthConverter(activeNote.Duration)
		smwNote := SmwNote{key, lengths, octave}
		smwTrack = append(smwTrack, smwNote)
		lastNoteEndTime = activeNote.StartTime + activeNote.Duration
	}
	return smwTrack
}

func trackDone(notes []midiNote, tick uint32) bool {
	lastNote := notes[len(notes)-1]
	return tick > lastNote.StartTime+lastNote.Duration
}

func getNoteWithStartTime(notes []midiNote, tick uint32) *midiNote {
	var potentialNotes = make([]midiNote, 0)
	for _, note := range notes {
		if note.StartTime == tick {
			potentialNotes = append(potentialNotes, note)
		}
	}
	if len(potentialNotes) > 0 {
		highestNote := midiNote{}
		for _, note := range potentialNotes {
			if note.Key > highestNote.Key {
				highestNote = note
			}
		}
		return &highestNote
	} else {
		return nil
	}
}

func noteValueToSmwKey(note midiNote) (key string, octave int) {
	noteValue := note.Key
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	if noteValue < 19 || noteValue > 88 {
		fmt.Printf("Error, cannot convert note value %d to SMW note (out of range)\n", noteValue)
		return "", -999
	}
	key = noteDict[int(noteValue%12)]
	octave = int(noteValue/12) - 1
	return key, octave
}

func getNoteLengthConverter(ticksPer64thNote uint32) func(duration uint32) (lengths []uint8) {
	ticksPer32ndNote := ticksPer64thNote * 2
	ticksPer16thNote := ticksPer32ndNote * 2
	ticksPer8thNote := ticksPer16thNote * 2
	ticksPerQuarterNote := ticksPer8thNote * 2
	ticksPerHalfNote := ticksPerQuarterNote * 2
	ticksPerWholeNote := ticksPerHalfNote * 2

	noteLengthToSmwLength := func(duration uint32) (uint8, uint32) {
		if duration > ticksPerWholeNote {
			return 1, duration - ticksPerWholeNote
		}
		num64thNotes := duration / ticksPer64thNote
		half := num64thNotes
		acc := 0
		for half != 1 {
			half = half / 2
			acc++
		}
		switch acc {
		case 0:
			return 64, duration - ticksPer64thNote
		case 1:
			return 32, duration - ticksPer32ndNote
		case 2:
			return 16, duration - ticksPer16thNote
		case 3:
			return 8, duration - ticksPer8thNote
		case 4:
			return 4, duration - ticksPerQuarterNote
		case 5:
			return 2, duration - ticksPerHalfNote
		case 6:
			return 1, duration - ticksPerWholeNote
		}
		fmt.Println("You shouldn't be here!")
		return 0, 0
	}

	return func(duration uint32) []uint8 {
		lengths := make([]uint8, 0)
		if duration == 0 {
			return lengths
		}
		length, remainder := noteLengthToSmwLength(duration)
		lengths = append(lengths, length)
		for remainder != 0 {
			length, remainder = noteLengthToSmwLength(remainder)
			lengths = append(lengths, length)
		}
		return lengths
	}
}
