package convert

import (
	"fmt"
)

func splitAllTracks(tracks []NoteTrack, trackSplitMap TrackSplitMap) []NoteTrack {
	newNoteTracks := make([]NoteTrack, 0)
	for i, track := range tracks {
		noteGroups, ok := trackSplitMap[i]
		if !ok {
			newNoteTracks = append(newNoteTracks, track)
			continue
		}
		splitTracks := splitIntoTracks(track, noteGroups)
		splitTracks = filterEmptyNoteTracks(splitTracks)
		for _, newTrack := range splitTracks {
			newNoteTracks = append(newNoteTracks, newTrack)
		}
	}
	return newNoteTracks
}

func splitIntoTracks(track NoteTrack, noteGroups []NoteGroup) []NoteTrack {
	if len(noteGroups) == 0 {
		return []NoteTrack{track}
	}

	newNoteTracks := make([]NoteTrack, len(noteGroups)+1)
	for _, note := range track.Notes {
		index := getIndexForNoteGroup(note, noteGroups)
		newNoteTracks[index].Notes = append(newNoteTracks[index].Notes, note)
	}
	newNoteTracks = setSplitTrackNames(track.Name, newNoteTracks)
	return newNoteTracks
}

func setSplitTrackNames(oldName string, tracks []NoteTrack) []NoteTrack {
	for i := range tracks {
		tracks[i].Name = fmt.Sprintf("%s - Split %d", oldName, i+1)
	}
	return tracks
}
