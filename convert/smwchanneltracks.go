package convert

import (
	"fmt"
)

var noteDict = map[uint8]string{
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

func createSmwChannelTracksForAllTracks(noteTracks []NoteTrack, ticksPer64thNote uint32) []SmwTrack {
	var smwTracks []SmwTrack
	longestTrackLength := getLongestTrackLength(noteTracks)
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)
	for _, noteTrack := range noteTracks {
		smwTrack := createSmwChannelTrack(noteTrack, longestTrackLength, noteLengthConverter)
		smwTracks = append(smwTracks, smwTrack)
	}
	return smwTracks
}

func getLongestTrackLength(noteTracks []NoteTrack) (longestTrackLength uint32) {
	for _, track := range noteTracks {
		trackLength := getTrackLength(track)
		if trackLength > longestTrackLength {
			longestTrackLength = trackLength
		}
	}
	return longestTrackLength
}

func getTrackLength(track NoteTrack) uint32 {
	if len(track.Notes) < 1 {
		return 0
	}
	lastNote := track.Notes[len(track.Notes)-1]
	return lastNote.StartTime + lastNote.Duration
}

func createSmwChannelTrack(noteTrack NoteTrack, length uint32, noteLengthConverter func(uint32) []uint8) SmwTrack {
	var smwTrack SmwTrack
	smwTrack.Name = noteTrack.Name
	notes := noteTrack.Notes
	// scan through track and create SMW channels until until no more notes
	for len(notes) > 0 {
		var smwNoteChannel []SmwNote
		var tick, lastNoteEndTime uint32
		var activeNote *MidiNote

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
					restSmwNote := Rest{lengths}
					smwNoteChannel = append(smwNoteChannel, restSmwNote)
				}
			}
			key, octave := noteValueToSmwKey(*activeNote)
			lengths := noteLengthConverter(activeNote.Duration)
			smwNote := Note{key, activeNote.Key, lengths, octave}
			smwNoteChannel = append(smwNoteChannel, smwNote)
			lastNoteEndTime = activeNote.StartTime + activeNote.Duration
		}
		if lastNoteEndTime < length {
			// pad out ending with rest so we don't prematurely loop when a track ends
			restLength := length - lastNoteEndTime
			lengths := noteLengthConverter(restLength)
			restSmwNote := Rest{lengths}
			smwNoteChannel = append(smwNoteChannel, restSmwNote)
		}

		newTrack := ChannelTrack{
			Notes:         smwNoteChannel,
			DefaultSample: noteTrack.DefaultSample,
			SampleMap:     noteTrack.SampleMap,
		}
		smwTrack.ChannelTracks = append(smwTrack.ChannelTracks, newTrack)
	}
	return smwTrack
}

func trackDone(notes []MidiNote, tick uint32) bool {
	if len(notes) == 0 {
		return true
	}
	lastNote := notes[len(notes)-1]
	return tick > lastNote.StartTime+lastNote.Duration
}

func extractHighestNoteAtStartTime(notes []MidiNote, tick uint32) ([]MidiNote, *MidiNote) {
	var potentialNotes = make([]MidiNote, 0)
	for _, note := range notes {
		if note.StartTime == tick {
			potentialNotes = append(potentialNotes, note)
		}
	}
	if len(potentialNotes) > 0 {
		highestNote := MidiNote{}
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

func removeNote(notes []MidiNote, tick uint32, key uint8) []MidiNote {
	for i, note := range notes {
		if note.StartTime == tick && note.Key == key {
			copy(notes[i:], notes[i+1:])
			notes[len(notes)-1] = MidiNote{}
			notes = notes[:len(notes)-1]
			return notes
		}
	}
	fmt.Printf("Could not remove note with start time %d, key %d", tick, key)
	return notes
}

func noteValueToSmwKey(note MidiNote) (key string, octave uint8) {
	noteValue := note.Key
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	if noteValue < 19 || noteValue > 88 {
		fmt.Printf("Error, cannot convert note value %d to SMW note (out of range). Inserting rest\n", noteValue)
		return "r", 0
	}
	key = noteDict[noteValue%12]
	octave = noteValue/12 - 1
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
