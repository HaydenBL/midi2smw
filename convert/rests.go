package convert

func insertRests(tracks []noteTrack) []noteTrack {
	for i, track := range tracks {
		tracks[i] = insertRestsIntoTrack(track)
	}
	return tracks
}

func insertRestsIntoTrack(track noteTrack) noteTrack {
	newNoteTrack := make([]midiNote, 0)
	for i := 0; i < len(track.Notes); i++ {
		newNoteTrack = append(newNoteTrack, track.Notes[i])
		var lastNote midiNote
		if i == 0 {
			lastNote = midiNote{
				Key:       0,
				Velocity:  0,
				isRest:    false,
				StartTime: 0,
				Duration:  0,
			}
		} else {
			lastNote = track.Notes[i-1]
		}
		currentNote := track.Notes[i]
		restLength := getRestLength(currentNote, lastNote)
		if restLength > 0 {
			newNoteTrack = append(newNoteTrack, midiNote{
				Key:       0,
				Velocity:  0,
				StartTime: lastNote.StartTime + lastNote.Duration,
				Duration:  restLength,
				isRest:    true,
			})
		}
	}
	track.Notes = newNoteTrack
	return track
}

func getRestLength(currentNote midiNote, lastNote midiNote) uint32 {
	lastEndTime := lastNote.StartTime + lastNote.Duration
	restGapDuration := currentNote.StartTime - lastEndTime
	return restGapDuration
}
