package convert

func insertRestsIntoAllTracks(tracks []noteTrack) []noteTrack {
	for i, track := range tracks {
		tracks[i].Notes = insertRestsIntoTrack(track.Notes)
	}
	return tracks
}

func insertRestsIntoTrack(notes []midiNote) []midiNote {
	newNoteTrack := make([]midiNote, 0)
	for i := 0; i < len(notes); i++ {
		var lastNote midiNote
		if i == 0 {
			if notes[i].StartTime > 0 {
				rest := getRestNote(0, notes[i].StartTime)
				newNoteTrack = append(newNoteTrack, rest)
			}
			newNoteTrack = append(newNoteTrack, notes[i])
			lastNote = midiNote{
				Key:       0,
				Velocity:  0,
				isRest:    false,
				StartTime: 0,
				Duration:  notes[i].StartTime,
			}
			continue
		} else {
			lastNote = notes[i-1]
		}
		currentNote := notes[i]
		restLength := getRestLength(currentNote, lastNote)
		if restLength > 0 {
			newNoteTrack = append(newNoteTrack, midiNote{
				Key:       0,
				Velocity:  0,
				StartTime: lastNote.StartTime + lastNote.Duration,
				Duration:  restLength,
				isRest:    true,
			})
			newNoteTrack = append(newNoteTrack, notes[i])
		} else {
			newNoteTrack = append(newNoteTrack, notes[i])
		}
	}
	notes = newNoteTrack
	return notes
}

func getRestNote(startTime, endTime uint32) midiNote {
	return midiNote{
		Key:       0,
		Velocity:  0,
		isRest:    true,
		StartTime: startTime,
		Duration:  endTime,
	}
}

func getRestLength(currentNote midiNote, lastNote midiNote) uint32 {
	lastEndTime := lastNote.StartTime + lastNote.Duration
	restGapDuration := currentNote.StartTime - lastEndTime
	return restGapDuration
}
