package convert

import (
	"fmt"
	"midi2smw/smwtypes"
)

func createSmwChannelTracksForAllTracks(noteTracks []NoteTrack, ticksPer64thNote uint32) []smwtypes.SmwTrack {
	var smwTracks []smwtypes.SmwTrack
	longestTrackLength := getLongestTrackLength(noteTracks)
	noteGenerator := smwtypes.NewNoteGenerator(ticksPer64thNote)
	for _, noteTrack := range noteTracks {
		smwTrack := createSmwChannelTrack(noteTrack, longestTrackLength, noteGenerator)
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

func createSmwChannelTrack(noteTrack NoteTrack, length uint32, noteGenerator smwtypes.NoteGenerator) smwtypes.SmwTrack {
	var smwTrack smwtypes.SmwTrack
	smwTrack.Name = noteTrack.Name
	notes := noteTrack.Notes
	// scan through track and create SMW channels until until no more notes
	for len(notes) > 0 {
		var smwNoteChannel []smwtypes.SmwNote
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
					restSmwNote := noteGenerator.NewRest(restLength)
					smwNoteChannel = append(smwNoteChannel, restSmwNote)
				}
			}
			smwNote := noteGenerator.NewNote(activeNote.Key, activeNote.Duration)
			smwNoteChannel = append(smwNoteChannel, smwNote)
			lastNoteEndTime = activeNote.StartTime + activeNote.Duration
		}
		if lastNoteEndTime < length {
			// pad out ending with rest so we don't prematurely loop when a track ends
			restLength := length - lastNoteEndTime
			restSmwNote := noteGenerator.NewRest(restLength)
			smwNoteChannel = append(smwNoteChannel, restSmwNote)
		}

		newTrack := smwtypes.ChannelTrack{
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
