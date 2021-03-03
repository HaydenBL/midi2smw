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

func removeOverlapping(tracks []noteTrack) []noteTrack {
	removeIndex := func(slice []midiNote, s int) []midiNote {
		return append(slice[:s], slice[s+1:]...)
	}

	for i, track := range tracks {
		overlapping := make([]int, 0)
		for j := 0; j < len(track.Notes); j++ {
			for k := j + 1; k < len(track.Notes); k++ {
				if overlap(track.Notes[j], track.Notes[k]) {
					overlapping = append(overlapping, k)
				}
			}
		}
		for j := len(overlapping) - 1; j >= 0; j-- {
			overlappingIndex := overlapping[j]
			tracks[i].Notes = removeIndex(tracks[i].Notes, overlappingIndex)
		}
	}

	return tracks
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
