package convert

import (
	"fmt"
	"midi2smw/convert/drumtrack"
)

func splitAllTracks(tracks []noteTrack, midiTracksWithNoteGroups []drumtrack.MidiTrackWithNoteGroups) []noteTrack {
	if len(tracks) != len(midiTracksWithNoteGroups) {
		fmt.Printf("Error splitting drum tracks")
		return tracks
	}

	newNoteTracks := make([]noteTrack, 0)
	for i := range tracks {
		splitTracks := splitIntoTracks(tracks[i], midiTracksWithNoteGroups[i].NoteGroups)
		splitTracks = filterEmptyNoteTracks(splitTracks)
		for _, newTrack := range splitTracks {
			newNoteTracks = append(newNoteTracks, newTrack)
		}
	}
	return newNoteTracks
}

func splitIntoTracks(track noteTrack, noteGroups []drumtrack.NoteGroup) []noteTrack {
	if len(noteGroups) == 0 {
		return []noteTrack{track}
	}

	newNoteTracks := make([]noteTrack, len(noteGroups)+1)
	for _, note := range track.Notes {
		index := getIndexForDrumTrackGroup(note, noteGroups)
		newNoteTracks[index].Notes = append(newNoteTracks[index].Notes, note)
	}
	newNoteTracks = setSplitTrackNames(track.Name, newNoteTracks)
	return newNoteTracks
}

func setSplitTrackNames(oldName string, tracks []noteTrack) []noteTrack {
	for i := range tracks {
		tracks[i].Name = fmt.Sprintf("%s - Split %d", oldName, i+1)
	}
	return tracks
}
