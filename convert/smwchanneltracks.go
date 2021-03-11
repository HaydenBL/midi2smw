package convert

import (
	"fmt"
)

type SmwNote struct {
	Key          string
	LengthValues []uint8
	Octave       int
}
type SmwTrack struct {
	ChannelTracks [][]SmwNote // if a midi track has chords/overlapping notes, we'll throw them into multiple channels
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

func createSmwChannelTracksForAllTracks(noteTracks []noteTrack, ticksPer64thNote uint32) []SmwTrack {
	var smwTracks []SmwTrack
	longestTrackLength := getLongestTrackLength(noteTracks)
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)
	for _, noteTrack := range noteTracks {
		smwTrack := createSmwChannelTrack(noteTrack.Notes, longestTrackLength, noteLengthConverter)
		smwTracks = append(smwTracks, smwTrack)
	}
	return smwTracks
}

func getLongestTrackLength(noteTracks []noteTrack) (longestTrackLength uint32) {
	for _, track := range noteTracks {
		trackLength := getTrackLength(track)
		if trackLength > longestTrackLength {
			longestTrackLength = trackLength
		}
	}
	return longestTrackLength
}

func getTrackLength(track noteTrack) uint32 {
	lastNote := track.Notes[len(track.Notes)-1]
	return lastNote.StartTime + lastNote.Duration
}

func createSmwChannelTrack(notes []midiNote, length uint32, noteLengthConverter func(uint32) []uint8) SmwTrack {
	var smwTrack SmwTrack
	// scan through track and create SMW channels until until no more notes
	for len(notes) > 0 {
		var smwChannel []SmwNote
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
			notes, activeNote = extractHighestNoteAtStartTime(notes, tick)
			if activeNote == nil {
				continue
			} else {
				// insert rest
				restLength := tick - lastNoteEndTime
				if restLength != 0 {
					lengths := noteLengthConverter(restLength)
					restSmwNote := SmwNote{Key: "r", LengthValues: lengths, Octave: 0}
					smwChannel = append(smwChannel, restSmwNote)
				}
			}
			key, octave := noteValueToSmwKey(*activeNote)
			lengths := noteLengthConverter(activeNote.Duration)
			smwNote := SmwNote{key, lengths, octave}
			smwChannel = append(smwChannel, smwNote)
			lastNoteEndTime = activeNote.StartTime + activeNote.Duration
		}
		if lastNoteEndTime < length {
			// pad out ending with rest so we don't prematurely loop when a track ends
			restLength := length - lastNoteEndTime
			lengths := noteLengthConverter(restLength)
			restSmwNote := SmwNote{Key: "r", LengthValues: lengths, Octave: 0}
			smwChannel = append(smwChannel, restSmwNote)
		}
		smwTrack.ChannelTracks = append(smwTrack.ChannelTracks, smwChannel)
	}
	return smwTrack
}

func trackDone(notes []midiNote, tick uint32) bool {
	if len(notes) == 0 {
		return true
	}
	lastNote := notes[len(notes)-1]
	return tick > lastNote.StartTime+lastNote.Duration
}

func extractHighestNoteAtStartTime(notes []midiNote, tick uint32) ([]midiNote, *midiNote) {
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
		notes = removeNote(notes, tick, highestNote.Key)
		return notes, &highestNote
	} else {
		return notes, nil
	}
}

func removeNote(notes []midiNote, tick uint32, key uint8) []midiNote {
	for i, note := range notes {
		if note.StartTime == tick && note.Key == key {
			copy(notes[i:], notes[i+1:])
			notes[len(notes)-1] = midiNote{}
			notes = notes[:len(notes)-1]
			return notes
		}
	}
	fmt.Printf("Could not remove note with start time %d, key %d", tick, key)
	return notes
}

func noteValueToSmwKey(note midiNote) (key string, octave int) {
	noteValue := note.Key
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	if noteValue < 19 || noteValue > 88 {
		fmt.Printf("Error, cannot convert note value %d to SMW note (out of range). Inserting rest\n", noteValue)
		return "r", 0
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
