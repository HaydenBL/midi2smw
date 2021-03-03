package convert

import "testing"

func TestOverlap(t *testing.T) {
	var note1, note2 midiNote

	note1 = midiNote{
		StartTime: 0,
		Duration:  10,
	}
	note2 = midiNote{
		StartTime: 10,
		Duration:  10,
	}
	if res := overlap(note1, note2); res == true {
		t.Fatalf("Whoops")
	}

	note1 = midiNote{
		StartTime: 0,
		Duration:  10,
	}
	note2 = midiNote{
		StartTime: 9,
		Duration:  10,
	}
	if res := overlap(note1, note2); res != true {
		t.Fatalf("Whoops")
	}

	note1 = midiNote{
		StartTime: 0,
		Duration:  11,
	}
	note2 = midiNote{
		StartTime: 9,
		Duration:  10,
	}
	if res := overlap(note1, note2); res != true {
		t.Fatalf("Whoops")
	}
}

func TestRemoveOverlapping(t *testing.T) {
	track := noteTrack{
		Notes: []midiNote{
			{
				StartTime: 0,
				Duration:  10,
			},
			{
				StartTime: 5,
				Duration:  10,
			},
			{
				StartTime: 10,
				Duration:  10,
			},
		},
	}

	track = removeOverlapping(track)

	if len(track.Notes) != 2 {
		t.Fatalf("Expected 2, actual: %d", len(track.Notes))
	}
	if track.Notes[1].StartTime != 10 && track.Notes[1].Duration != 10 {
		t.Fatalf("Whoops")
	}
}
