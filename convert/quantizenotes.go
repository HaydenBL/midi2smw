package convert

func quantizeNotesOnAllTracks(tracks []noteTrack, ticksPer64thNote uint32) []noteTrack {
	quantizer := getQuantizer(ticksPer64thNote)
	for i, track := range tracks {
		tracks[i].Notes = quantizeNotes(track.Notes, quantizer)
	}
	return tracks
}

func quantizeNotes(notes []midiNote, quantizer func(uint32) uint32) []midiNote {
	quantizedNotes := make([]midiNote, 0)
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

func removeOverlappingOnAllTracks(tracks []noteTrack) []noteTrack {
	modifiedTracks := make([]noteTrack, 0)
	for _, track := range tracks {
		modifiedTracks = append(modifiedTracks, removeOverlapping(track))
	}
	return modifiedTracks
}

func removeOverlapping(track noteTrack) noteTrack {
	removeIndex := func(slice []midiNote, s int) []midiNote {
		return append(slice[:s], slice[s+1:]...)
	}
	rmv := func(notes []midiNote) (newNotes []midiNote, more bool) {
		for i := range notes {
			if i == 0 {
				continue
			}
			lastNote := notes[i-1]
			currentNote := notes[i]
			if overlap(lastNote, currentNote) {
				notes = removeIndex(notes, i)
				return notes, true
			}
		}
		return notes, false
	}
	var more bool
	track.Notes, more = rmv(track.Notes)
	for more {
		track.Notes, more = rmv(track.Notes)
	}

	return track
}

func closer(n, subtracted, added uint32) uint32 {
	if n-subtracted < added-n {
		return subtracted
	} else {
		return added
	}
}

func overlap(note1 midiNote, note2 midiNote) bool {
	note1EndTime := note1.StartTime + note1.Duration
	note2EndTime := note2.StartTime + note2.Duration
	//x1 <= y2 && y1 <= x2
	return note1.StartTime < note2EndTime && note2.StartTime < note1EndTime
}
