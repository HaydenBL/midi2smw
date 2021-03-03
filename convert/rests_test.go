package convert

import "testing"

func TestInsertRestsIntoTrack(t *testing.T) {
	var notes []midiNote

	notes = []midiNote{
		{
			StartTime: 0,
			Duration:  10,
		},
		{
			StartTime: 20,
			Duration:  10,
		},
	}

	notes = insertRestsIntoTrack(notes)

	if len(notes) != 3 {
		t.Fatalf("Oh no!")
	}
	if !notes[1].isRest {
		t.Fatalf("Shoulda been a rest!")
	}
	if notes[1].StartTime != 10 || notes[1].Duration != 10 {
		t.Fatalf("Whaaaat")
	}
	if notes[2].StartTime != 20 || notes[2].Duration != 10 {
		t.Fatalf("That ain't right")
	}

	notes = []midiNote{
		{
			StartTime: 20,
			Duration:  10,
		},
	}

	notes = insertRestsIntoTrack(notes)

	if len(notes) != 2 {
		t.Fatalf("Flip")
	}
	if !notes[0].isRest {
		t.Fatalf("That's not a rest!")
	}
	if notes[0].StartTime != 0 || notes[0].Duration != 20 {
		t.Fatalf("Bah")
	}
	if notes[1].StartTime != 20 || notes[1].Duration != 10 {
		t.Fatalf("Whawhawhaaaa")
	}
}
