package convert

func quantizeNotesOnAllTracks(tracks []NoteTrack, ticksPer64thNote uint32) []NoteTrack {
	quantizer := getQuantizer(ticksPer64thNote)
	for i, track := range tracks {
		tracks[i].Notes = quantizeNotes(track.Notes, quantizer)
	}
	return tracks
}

func quantizeNotes(notes []MidiNote, quantizer func(uint32) uint32) []MidiNote {
	quantizedNotes := make([]MidiNote, 0)
	for _, note := range notes {
		note.StartTime = quantizer(note.StartTime)
		note.Duration = quantizer(note.Duration)
		if note.Duration > 0 {
			quantizedNotes = append(quantizedNotes, note)
		}
	}
	return quantizedNotes
}

func getQuantizer(ticksPer64thNote uint32) func(uint32) uint32 {
	return func(timestamp uint32) uint32 {
		timestampMod64thNote := timestamp % ticksPer64thNote
		timestampMinusMod := timestamp - timestampMod64thNote
		closest64thToTimestamp := closer(timestamp, timestampMinusMod, timestampMinusMod+ticksPer64thNote)
		return closest64thToTimestamp
	}
}

func closer(n, subtracted, added uint32) uint32 {
	if n-subtracted < added-n {
		return subtracted
	} else {
		return added
	}
}
