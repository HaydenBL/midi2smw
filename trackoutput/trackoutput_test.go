package trackoutput

import (
	"midi2smw/convert"
	"strings"
	"testing"
)

func TestWriteChannel(t *testing.T) {
	channel := convert.ChannelTrack{
		Notes: []convert.SmwNote{
			convert.Note{
				KeyValue:     24,
				LengthValues: []uint8{8, 16},
			},
			convert.Note{
				KeyValue:     38,
				LengthValues: []uint8{16},
			},
			convert.Rest{
				LengthValues: []uint8{2, 4},
			},
			convert.Note{
				KeyValue:     24,
				LengthValues: []uint8{16, 32},
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
