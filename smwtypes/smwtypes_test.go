package smwtypes

import "testing"

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
