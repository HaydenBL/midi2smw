package convert

import "midi2smw/drumtrack"

func getIndexForDrumTrackGroup(note midiNote, dtg drumtrack.Group) int {
	for i, noteGroup := range dtg.NoteGroups {
		if contains(note.Key, noteGroup) {
			return i
		}
	}
	return len(dtg.NoteGroups)
}

func contains(num uint8, arr []uint8) bool {
	for _, n := range arr {
		if num == n {
			return true
		}
	}
	return false
}

func findNoteIndex(notes []midiNote, key uint8) (int, midiNote) {
	for i, note := range notes {
		if note.Key == key {
			return i, note
		}
	}
	return -1, midiNote{}
}

func deleteAtIndex(notes []midiNote, i int) []midiNote {
	copy(notes[i:], notes[i+1:])
	notes[len(notes)-1] = midiNote{}
	notes = notes[:len(notes)-1]
	return notes
}

func minUint8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func maxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
