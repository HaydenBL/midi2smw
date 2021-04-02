package convert

import (
	"midi2smw/smwtypes"
	"reflect"
	"testing"
)

func TestCreateSmwChannelTrack_singleTrack(t *testing.T) {
	var notes []MidiNote
	var ticksPer64thNote uint32 = 30
	noteGenerator := smwtypes.GetNoteGenerator(ticksPer64thNote)

	notes = []MidiNote{
		{
			Key:       24,
			StartTime: 0,
			Duration:  30, // endtime 30
		},
		{
			Key:       24,
			StartTime: 30,
			Duration:  60, // endtime 90
		},
		{
			Key:       24,
			StartTime: 150,
			Duration:  60, // endtime 210
		},
	}

	expected := []smwtypes.SmwNote{
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{64},
		},
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
		smwtypes.Rest{
			LengthValues: []uint8{32},
		},
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
	}

	smwNotes := createSmwChannelTrack(NoteTrack{Notes: notes}, 210, noteGenerator)
	if len(smwNotes.ChannelTracks) != 1 {
		t.Fatalf("Expected 1 track, got %d", len(smwNotes.ChannelTracks))
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[0].Notes, expected) {
		t.Fatalf("Not equal?")
	}

}

func TestCreateSmwChannelTrack_padsEndingProperly(t *testing.T) {
	var notes []MidiNote
	var ticksPer64thNote uint32 = 30
	noteGenerator := smwtypes.GetNoteGenerator(ticksPer64thNote)

	notes = []MidiNote{
		{
			Key:       24,
			StartTime: 0,
			Duration:  30, // endtime 30
		},
	}

	expected := []smwtypes.SmwNote{
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{64},
		},
		smwtypes.Rest{
			LengthValues: []uint8{64},
		},
	}

	smwNotes := createSmwChannelTrack(NoteTrack{Notes: notes}, 60, noteGenerator)
	if len(smwNotes.ChannelTracks) != 1 {
		t.Fatalf("Expected 1 track, got %d", len(smwNotes.ChannelTracks))
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[0].Notes, expected) {
		t.Fatalf("Not equal?")
	}

}

func TestCreateSmwChannelTrack_multiTrack(t *testing.T) {
	var notes []MidiNote
	var ticksPer64thNote uint32 = 30
	noteGenerator := smwtypes.GetNoteGenerator(ticksPer64thNote)

	notes = []MidiNote{
		{
			Key:       24,
			StartTime: 0,
			Duration:  60, // endtime 60
		},
		{
			Key:       24,
			StartTime: 30,
			Duration:  60, // endtime 90
		},
		{
			Key:       24,
			StartTime: 60,
			Duration:  60, // endtime 120
		},
	}

	expectedTrack1 := []smwtypes.SmwNote{
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
	}

	expectedTrack2 := []smwtypes.SmwNote{
		smwtypes.Rest{
			LengthValues: []uint8{64},
		},
		smwtypes.Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
		smwtypes.Rest{
			LengthValues: []uint8{64},
		},
	}

	smwNotes := createSmwChannelTrack(NoteTrack{Notes: notes}, 120, noteGenerator)
	if len(smwNotes.ChannelTracks) != 2 {
		t.Fatalf("Expected 2 tracks, got %d", len(smwNotes.ChannelTracks))
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[0].Notes, expectedTrack1) {
		t.Fatalf("First track not equal!")
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[1].Notes, expectedTrack2) {
		t.Fatalf("Second track not equal!")
	}
}
