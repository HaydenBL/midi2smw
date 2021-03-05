package convert

import (
	"reflect"
	"testing"
)

func TestCreateSmwChannelTrack(t *testing.T) {
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
			Duration:  60, // endtime 180
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

	smwNotes := createSmwChannelTrack(notes, noteLengthConverter)
	if !reflect.DeepEqual(smwNotes, expected) {
		t.Fatalf("Not equal?")
	}

}

func TestNoteValueToSmwKey(t *testing.T) {
	var k uint8
	// too low for SMW note
	for k = 0; k < 19; k++ {
		key, octave := noteValueToSmwKey(midiNote{Key: k})
		if key != "" || octave != -999 {
			t.Fatalf("Should error when key value too low")
		}
	}
	// valid SMW note range
	for k = 19; k < 89; k++ {
		key, _ := noteValueToSmwKey(midiNote{Key: k})
		if key == "" {
			t.Fatalf("Shouldn't error when key value is within range")
		}
	}
	// too high for SMW note
	for k = 89; k < 128; k++ {
		key, octave := noteValueToSmwKey(midiNote{Key: k})
		if key != "" || octave != -999 {
			t.Fatalf("Should error when key value too high")
		}
	}
}
