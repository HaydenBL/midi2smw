package convert

func getIndexForDrumTrackGroup(note MidiNote, noteGroups []NoteGroup) int {
	for i, noteGroup := range noteGroups {
		if contains(note.Key, noteGroup.Notes) {
			return i
		}
	}
	return len(noteGroups)
}

func contains(num uint8, arr []uint8) bool {
	for _, n := range arr {
		if num == n {
			return true
		}
	}
	return false
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
