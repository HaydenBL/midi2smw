package convert

func quantizeNotes(tracks []noteTrack, ticksPer64thNote uint32) []noteTrack {
	quantizeTimestamp := func(timestamp uint32) uint32 {
		timestampMod64thNote := timestamp % ticksPer64thNote
		timestampMinusMod := timestamp - timestampMod64thNote
		closest64thToTimestamp := closer(timestamp, timestampMinusMod, timestampMinusMod+ticksPer64thNote)
		return closest64thToTimestamp
	}

	for i, track := range tracks {
		for j, note := range track.Notes {
			tracks[i].Notes[j].StartTime = quantizeTimestamp(note.StartTime)
			tracks[i].Notes[j].Duration = quantizeTimestamp(note.Duration)
		}
	}
	return tracks
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
