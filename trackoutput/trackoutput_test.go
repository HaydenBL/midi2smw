package trackoutput

import (
	"midi2smw/convert"
	"strings"
	"testing"
)

func TestWriteChannel(t *testing.T) {
	channel := convert.ChannelTrack{
		Notes: []convert.SmwNote{
			testNote{
				key:          "c",
				keyValue:     24,
				lengthValues: []uint8{8, 16},
				octave:       1,
			},
			testNote{
				key:          "d",
				keyValue:     38,
				lengthValues: []uint8{16},
				octave:       2,
			},
			testRest{
				lengthValues: []uint8{2, 4},
			},
			testNote{
				key:          "c",
				keyValue:     24,
				lengthValues: []uint8{16, 32},
				octave:       1,
			},
		},
		DefaultSample: 0,
		SampleMap: map[uint8]uint8{
			38: 1,
		},
	}

	expectedString := "c8^16>@1d16r2^4<@0c16^32"

	var sb strings.Builder
	writeChannel(&sb, channel)

	if sb.String() != expectedString {
		t.Fatalf("Error in track output:\nExpected:\t%s\nActual:\t\t%s", expectedString, sb.String())
	}
}

// Re-implementing structs to implement SmwNote for these tests
// Don't like this but whatever
type testNote struct {
	key          string
	keyValue     uint8
	lengthValues []uint8
	octave       uint8
}

func (n testNote) Key() string {
	return n.key
}

func (n testNote) KeyValue() uint8 {
	return n.keyValue
}

func (n testNote) Octave() uint8 {
	return n.octave
}

func (n testNote) LengthValues() []uint8 {
	return n.lengthValues
}

type testRest struct {
	lengthValues []uint8
}

func (r testRest) Key() string {
	return "r"
}

func (r testRest) KeyValue() uint8 {
	return 0
}

func (r testRest) Octave() uint8 {
	return 0
}

func (r testRest) LengthValues() []uint8 {
	return r.lengthValues
}
