package convert

import (
	"reflect"
	"testing"
)

func TestCreateSmwChannelTrack_singleTrack(t *testing.T) {
	var notes []midiNote
	var ticksPer64thNote uint32 = 30
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)

	notes = []midiNote{
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
		{
			Key:          "c",
			LengthValues: []uint8{64},
			Octave:       1,
		},
		{
			Key:          "c",
			LengthValues: []uint8{32},
			Octave:       1,
		},
		{
			Key:          "r",
			LengthValues: []uint8{32},
			Octave:       0,
		},
		{
			Key:          "c",
			LengthValues: []uint8{32},
			Octave:       1,
		},
	}

	smwNotes := createSmwChannelTrack(notes, "", 210, noteLengthConverter)
	if len(smwNotes.ChannelTracks) != 1 {
		t.Fatalf("Expected 1 track, got %d", len(smwNotes.ChannelTracks))
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[0], expected) {
		t.Fatalf("Not equal?")
	}

}

func TestCreateSmwChannelTrack_padsEndingProperly(t *testing.T) {
	var notes []midiNote
	var ticksPer64thNote uint32 = 30
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)

	notes = []midiNote{
		{
			Key:       24,
			StartTime: 0,
			Duration:  30, // endtime 30
		},
	}

	expected := []SmwNote{
		{
			Key:          "c",
			LengthValues: []uint8{64},
			Octave:       1,
		},
		{
			Key:          "r",
			LengthValues: []uint8{64},
			Octave:       0,
		},
	}

	smwNotes := createSmwChannelTrack(notes, "", 60, noteLengthConverter)
	if len(smwNotes.ChannelTracks) != 1 {
		t.Fatalf("Expected 1 track, got %d", len(smwNotes.ChannelTracks))
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[0], expected) {
		t.Fatalf("Not equal?")
	}

}

func TestCreateSmwChannelTrack_multiTrack(t *testing.T) {
	var notes []midiNote
	var ticksPer64thNote uint32 = 30
	noteLengthConverter := getNoteLengthConverter(ticksPer64thNote)

	notes = []midiNote{
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
		{
			Key:          "c",
			LengthValues: []uint8{32},
			Octave:       1,
		},
		{
			Key:          "c",
			LengthValues: []uint8{32},
			Octave:       1,
		},
	}

	expectedTrack2 := []SmwNote{
		{
			Key:          "r",
			LengthValues: []uint8{64},
			Octave:       0,
		},
		{
			Key:          "c",
			LengthValues: []uint8{32},
			Octave:       1,
		},
		{
			Key:          "r",
			LengthValues: []uint8{64},
			Octave:       0,
		},
	}

	smwNotes := createSmwChannelTrack(notes, "", 120, noteLengthConverter)
	if len(smwNotes.ChannelTracks) != 2 {
		t.Fatalf("Expected 2 tracks, got %d", len(smwNotes.ChannelTracks))
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[0], expectedTrack1) {
		t.Fatalf("First track not equal!")
	}
	if !reflect.DeepEqual(smwNotes.ChannelTracks[1], expectedTrack2) {
		t.Fatalf("Second track not equal!")
	}
}

func TestNoteValueToSmwKey(t *testing.T) {
	var k uint8
	// too low for SMW note
	for k = 0; k < 19; k++ {
		key, octave := noteValueToSmwKey(midiNote{Key: k})
		if key != "r" || octave != 0 {
			t.Fatalf("Should error when key value too low")
		}
	}
	// valid SMW note range
	for k = 19; k < 89; k++ {
		key, _ := noteValueToSmwKey(midiNote{Key: k})
		if key == "r" {
			t.Fatalf("Shouldn't error when key value is within range")
		}
	}
	// too high for SMW note
	for k = 89; k < 128; k++ {
		key, octave := noteValueToSmwKey(midiNote{Key: k})
		if key != "r" || octave != 0 {
			t.Fatalf("Should error when key value too high")
		}
	}
}
