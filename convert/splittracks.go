package convert

import (
	"fmt"
)

func splitAllTracks(tracks []NoteTrack, midiTracksWithNoteGroups []MidiTrackWithNoteGroups) []NoteTrack {
	if len(tracks) != len(midiTracksWithNoteGroups) {
		fmt.Printf("Error splitting tracks - lengths not equal")
		return tracks
	}

	newNoteTracks := make([]NoteTrack, 0)
	for i := range tracks {
		splitTracks := splitIntoTracks(tracks[i], midiTracksWithNoteGroups[i].NoteGroups)
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
