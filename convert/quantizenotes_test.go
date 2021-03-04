package convert

import "testing"

func TestQuantizeNotes(t *testing.T) {
	var ticksPer64thNote uint32 = 30
	quantizer := getQuantizer(ticksPer64thNote)
	var notes []midiNote

	notes = []midiNote{
		{
			StartTime: 0,
			Duration:  14,
		},
		{
			StartTime: 0,
			Duration:  15,
		},
		{
			StartTime: 0,
			Duration:  44,
		},
		{
			StartTime: 0,
			Duration:  45,
		},
	}

	notes = quantizeNotes(notes, quantizer)

	if len(notes) > 3 {
		t.Fatalf("Wah")
	}
	if notes[0].StartTime != 0 || notes[0].Duration != 30 {
		t.Fatalf("Boop")
	}
	if notes[1].StartTime != 0 || notes[1].Duration != 30 {
		t.Fatalf("Floop")
	}
	if notes[2].StartTime != 0 || notes[2].Duration != 60 {
		t.Fatalf("Floop")
	}
}
