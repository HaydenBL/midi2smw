package smwtypes

import (
	"log"
	"testing"
)

func TestChannelTrackString(t *testing.T) {
	channel := ChannelTrack{
		Notes: []SmwNote{
			Note{
				KeyValue:     24,
				LengthValues: []uint8{8, 16},
			},
			Note{
				KeyValue:     38,
				LengthValues: []uint8{16},
			},
			Rest{
				LengthValues: []uint8{2, 4},
			},
			Note{
				KeyValue:     24,
				LengthValues: []uint8{16, 32},
			},
		},
		DefaultSample: 0,
		SampleMap: map[uint8]uint8{
			38: 1,
		},
	}

	expected := "c8^16>@1d16r2^4<@0c16^32"
	actual := channel.String()

	if actual != expected {
		t.Fatalf("Error in track output:\nExpected:\t%s\nActual:\t\t%s", expected, actual)
	}
}

func TestStringCompressed(t *testing.T) {
	channel := ChannelTrack{
		Notes: []SmwNote{
			Note{
				KeyValue:     44,
				LengthValues: []uint8{8, 16},
			},
			Note{
				KeyValue:     58,
				LengthValues: []uint8{16},
			},
			Rest{
				LengthValues: []uint8{2, 4},
			},
			Note{
				KeyValue:     38,
				LengthValues: []uint8{16},
			},
			// loop break
			Note{
				KeyValue:     44,
				LengthValues: []uint8{8, 16},
			},
			Note{
				KeyValue:     58,
				LengthValues: []uint8{16},
			},
			Rest{
				LengthValues: []uint8{2, 4},
			},
			Note{
				KeyValue:     38,
				LengthValues: []uint8{16},
			},
			// end of loops
			Note{
				KeyValue:     40,
				LengthValues: []uint8{8, 16},
			},
		},
		DefaultSample: 0,
		SampleMap: map[uint8]uint8{
			38: 1,
		},
	}

	expected := "[@0 g+8^16>a+16r2^4<@1d16]2 @0 e8^16 "
	actual := channel.StringCompressed()

	if actual != expected {
		t.Fatalf("Error in compressed track output:\nExpected:\t%s\nActual:\t\t%s", expected, actual)
	}
}

func TestGetNumLoops(t *testing.T) {

	remainingTrack := []SmwNote{
		Note{
			KeyValue:     24,
			LengthValues: []uint8{8, 16},
		},
		Note{
			KeyValue:     38,
			LengthValues: []uint8{16},
		},
		Rest{
			LengthValues: []uint8{2, 4},
		},
		// loop break
		Note{
			KeyValue:     24,
			LengthValues: []uint8{8, 16},
		},
		Note{
			KeyValue:     38,
			LengthValues: []uint8{16},
		},
		Rest{
			LengthValues: []uint8{2, 4},
		},
		// end of loops
		Note{
			KeyValue:     20,
			LengthValues: []uint8{8, 16},
		},
	}

	numLoops := getNumLoops(3, remainingTrack)

	if numLoops != 2 {
		log.Fatalln("getNumLoops returned unexpected value")
	}
}

func TestGetLoopSection(t *testing.T) {
	track := []SmwNote{
		Note{
			KeyValue:     24,
			LengthValues: []uint8{8, 16},
		},
		Note{
			KeyValue:     38,
			LengthValues: []uint8{16},
		},
		Rest{
			LengthValues: []uint8{2, 4},
		},
		// loop break
		Note{
			KeyValue:     24,
			LengthValues: []uint8{8, 16},
		},
		Note{
			KeyValue:     38,
			LengthValues: []uint8{16},
		},
		Rest{
			LengthValues: []uint8{2, 4},
		},
		// end of loops
		Note{
			KeyValue:     20,
			LengthValues: []uint8{8, 16},
		},
	}

	expectedLoopSection := loopSection{
		loops: 2,
		notes: []SmwNote{
			Note{
				KeyValue:     24,
				LengthValues: []uint8{8, 16},
			},
			Note{
				KeyValue:     38,
				LengthValues: []uint8{16},
			},
			Rest{
				LengthValues: []uint8{2, 4},
			},
		},
	}

	expectedRemainingTrack := []SmwNote{
		Note{
			KeyValue:     20,
			LengthValues: []uint8{8, 16},
		},
	}

	actualLoopSection, actualRemainingTrack := getLoopSection(track)

	if expectedLoopSection.loops != actualLoopSection.loops ||
		!NoteSlicesEqual(expectedLoopSection.notes, actualLoopSection.notes) ||
		!NoteSlicesEqual(expectedRemainingTrack, actualRemainingTrack) {
		log.Fatalln("Unexpected return from getLoopSection")
	}
}
