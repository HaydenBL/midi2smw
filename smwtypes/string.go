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
	ctx := &channelWriteContext{
		lastOctave:    ct.Notes[0].GetOctave(),
		lastSample:    ct.DefaultSample,
		defaultSample: ct.DefaultSample,
		sampleMap:     ct.SampleMap,
		sb:            &strings.Builder{},
	}

	for _, smwNote := range ct.Notes {
		if rest, ok := smwNote.(Rest); ok {
			writeNote(rest, ctx)
			continue
		}
		shiftOctaveIfNecessary(smwNote.GetOctave(), ctx)
		switchSampleIfNecessary(smwNote.GetKeyValue(), ctx)
		writeNote(smwNote, ctx)
		ctx.lastOctave = smwNote.GetOctave()
	}
	return ctx.sb.String()
}

func writeNote(note SmwNote, ctx *channelWriteContext) {
	for i, length := range note.GetLengthValues() {
		if i == 0 {
			write(ctx.sb, "%s%d", note.GetKey(), length)
		} else {
			write(ctx.sb, "^%d", length)
		}
	}
}

func shiftOctaveIfNecessary(octave uint8, ctx *channelWriteContext) {
	if octave == ctx.lastOctave {
		return
	}
	var shiftToken string
	if octave > ctx.lastOctave {
		shiftToken = ">"
	} else {
		shiftToken = "<"
	}
	diff := int(octave) - int(ctx.lastOctave)
	if diff < 0 {
		diff = -diff
	}
	for i := 0; i < diff; i++ {
		write(ctx.sb, shiftToken)
	}
}

func switchSampleIfNecessary(keyValue uint8, ctx *channelWriteContext) {
	sample, ok := ctx.sampleMap[keyValue]
	if !ok {
		sample = ctx.defaultSample
	}
	if sample != ctx.lastSample {
		ctx.lastSample = sample
		write(ctx.sb, "@%d", sample)
	}
}
