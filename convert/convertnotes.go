package convert

import (
	"fmt"
	"midi2smw/midi"
)

type MidiNote struct {
	Key       uint8
	Velocity  uint8
	StartTime uint32
	Duration  uint32
}

type NoteTrack struct {
	Name          string
	Notes         []MidiNote
	MaxNote       uint8
	MinNote       uint8
	DefaultSample uint8
	SampleMap     map[uint8]uint8
}

func convertToNotes(midiTracks []midi.Track, splitTracks bool) []NoteTrack {
	var noteTracks = make([]NoteTrack, len(midiTracks))

	for trackIndex, track := range midiTracks {
		var wallTime uint32
		var notesBeingProcessed []MidiNote

		for _, event := range track.Events {
			wallTime += event.DeltaTick
			if event.Event == midi.NoteOn {
				notesBeingProcessed = append(notesBeingProcessed, MidiNote{event.Key, event.Velocity, wallTime, 0})
			}
			if event.Event == midi.NoteOff {
				i, note := findNoteIndex(notesBeingProcessed, event.Key)
				if i != -1 {
					note.Duration = wallTime - note.StartTime
					noteTracks[trackIndex].Name = track.Name
					noteTracks[trackIndex].Notes = append(noteTracks[trackIndex].Notes, note)
					noteTracks[trackIndex].MinNote = minUint8(noteTracks[trackIndex].MinNote, note.Key)
					noteTracks[trackIndex].MaxNote = maxUint8(noteTracks[trackIndex].MaxNote, note.Key)
					notesBeingProcessed = deleteAtIndex(notesBeingProcessed, i)
				}
			}
		}
	}

	noteTracks = filterEmptyNoteTracks(noteTracks)
	if splitTracks {
		trackSplitMap := SpecifyTrackSplits(midiTracks)
		noteTracks = splitAllTracks(noteTracks, trackSplitMap)
	}

	return noteTracks
}

func filterEmptyNoteTracks(tracks []NoteTrack) []NoteTrack {
	var nonEmptyTracks []NoteTrack
	for _, track := range tracks {
		if len(track.Notes) != 0 {
			nonEmptyTracks = append(nonEmptyTracks, track)
		}
	}

	if len(nonEmptyTracks) < len(tracks) {
		fmt.Printf("Removed %d tracks with no note data\n", len(tracks)-len(nonEmptyTracks))
	}
	return nonEmptyTracks
}
