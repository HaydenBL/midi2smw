package convert

import (
	"reflect"
	"testing"
)

func TestCreateSmwChannelTrack_singleTrack(t *testing.T) {
	var notes []MidiNote
	var ticksPer64thNote uint32 = 30
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)

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

	expected := []SmwNote{
		Note{
			KeyValue:     24,
			LengthValues: []uint8{64},
		},
		Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
		Rest{
			LengthValues: []uint8{32},
		},
		Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
	}

	smwNotes := createSmwChannelTrack(NoteTrack{Notes: notes}, 210, noteLengthConverter)
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
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)

	notes = []MidiNote{
		{
			Key:       24,
			StartTime: 0,
			Duration:  30, // endtime 30
		},
	}

	expected := []SmwNote{
		Note{
			KeyValue:     24,
			LengthValues: []uint8{64},
		},
		Rest{
			LengthValues: []uint8{64},
		},
	}

	smwNotes := createSmwChannelTrack(NoteTrack{Notes: notes}, 60, noteLengthConverter)
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
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)

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

	expectedTrack1 := []SmwNote{
		Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
		Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
	}

	expectedTrack2 := []SmwNote{
		Rest{
			LengthValues: []uint8{64},
		},
		Note{
			KeyValue:     24,
			LengthValues: []uint8{32},
		},
		Rest{
			LengthValues: []uint8{64},
		},
	}

	smwNotes := createSmwChannelTrack(NoteTrack{Notes: notes}, 120, noteLengthConverter)
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

func TestGetKeyFromKeyValue(t *testing.T) {
	var k uint8
	// too low for SMW note
	for k = 0; k < 19; k++ {
		key := getKeyFromKeyValue(k)
		if key != "r" {
			t.Fatalf("Should error when key value too low")
		}
	}
	// valid SMW note range
	for k = 19; k < 89; k++ {
		key := getKeyFromKeyValue(k)
		if key == "r" {
			t.Fatalf("Shouldn't error when key value is within range")
		}
	}
	// too high for SMW note
	for k = 89; k < 128; k++ {
		key := getKeyFromKeyValue(k)
		if key != "r" {
			t.Fatalf("Should error when key value too high")
		}
	}
}
