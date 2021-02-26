package convert

import "midi2smw/midiparse"

type NoteTrack struct {
	Notes   []MidiNote
	MaxNote uint8
	MinNote uint8
}

func convertNotes(tracks []midiparse.MidiTrack) []NoteTrack {
	var noteTracks = make([]NoteTrack, len(tracks))

	for trackIndex, track := range tracks {
		var wallTime uint32
		var notesBeingProcessed []MidiNote

		for _, event := range track.Events {
			wallTime += event.DeltaTick
			if event.Event == midiparse.NoteOn {
				notesBeingProcessed = append(notesBeingProcessed, MidiNote{event.Key, event.Velocity, wallTime, 0})
			}
			if event.Event == midiparse.NoteOff {
				i, note := findNoteIndex(notesBeingProcessed, event.Key)
				if i != -1 {
					note.Duration = wallTime - note.StartTime
					noteTracks[trackIndex].Notes = append(noteTracks[trackIndex].Notes, note)
					noteTracks[trackIndex].MinNote = min(noteTracks[trackIndex].MinNote, note.Key)
					noteTracks[trackIndex].MaxNote = max(noteTracks[trackIndex].MaxNote, note.Key)
					notesBeingProcessed = deleteAtIndex(notesBeingProcessed, i)
				}
			}
		}
	}

	return noteTracks
}

func findNoteIndex(notes []MidiNote, key uint8) (int, MidiNote) {
	for i, note := range notes {
		if note.Key == key {
			return i, note
		}
	}
	return -1, MidiNote{}
}

func deleteAtIndex(notes []MidiNote, i int) []MidiNote {
	copy(notes[i:], notes[i+1:])
	notes[len(notes)-1] = MidiNote{}
	notes = notes[:len(notes)-1]
	return notes
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
