package smwtypes

import (
	"fmt"
	"io"
	"strings"
)

func write(writer io.Writer, format string, a ...interface{}) {
	formattedString := fmt.Sprintf(format, a...)
	_, err := writer.Write([]byte(formattedString))
	if err != nil {
		fmt.Printf("Error writing string: %s", formattedString)
	}
}

type channelWriteContext struct {
	lastOctave    uint8
	lastSample    uint8
	defaultSample uint8
	sampleMap     map[uint8]uint8
	sb            *strings.Builder
}

func (ct ChannelTrack) String() string {
	ctx := channelWriteContext{
		lastOctave:    ct.Notes[0].GetOctave(),
		lastSample:    ct.DefaultSample,
		defaultSample: ct.DefaultSample,
		sampleMap:     ct.SampleMap,
		sb:            &strings.Builder{},
	}

	for _, smwNote := range ct.Notes {
		if rest, ok := smwNote.(Rest); ok {
			writeRest(rest, &ctx)
		} else {
			shiftOctaveIfNecessary(smwNote.GetOctave(), &ctx)
			writeNote(smwNote, &ctx)
			ctx.lastOctave = smwNote.GetOctave()
		}
	}
	return ctx.sb.String()
}

func writeNote(note SmwNote, ctx *channelWriteContext) {
	for i, length := range note.GetLengthValues() {
		if i == 0 {
			// Check if we need to swap the sample
			sample, ok := ctx.sampleMap[note.GetKeyValue()]
			if !ok {
				sample = ctx.defaultSample
			}
			if sample != ctx.lastSample {
				ctx.lastSample = sample
				write(ctx.sb, "@%d", sample)
			}
			write(ctx.sb, "%s%d", note.GetKey(), length)
		} else {
			write(ctx.sb, "^%d", length)
		}
	}
}

func writeRest(rest Rest, ctx *channelWriteContext) {
	for i, length := range rest.LengthValues {
		if i == 0 {
			write(ctx.sb, "r%d", length)
		} else {
			write(ctx.sb, "^%d", length)
		}
	}
}

func shiftOctaveIfNecessary(octave uint8, ctx *channelWriteContext) {
	if octave > ctx.lastOctave {
		for i := uint8(0); i < octave-ctx.lastOctave; i++ {
			write(ctx.sb, ">")
		}
	} else if octave < ctx.lastOctave {
		for i := uint8(0); i < ctx.lastOctave-octave; i++ {
			write(ctx.sb, "<")
		}
	}
}
