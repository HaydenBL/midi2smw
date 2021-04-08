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
	return ct.stringNotes(ct.Notes)
}

func (ct ChannelTrack) stringNotes(notes []SmwNote) string {
	if len(notes) == 0 {
		return ""
	}

	ctx := &channelWriteContext{
		lastOctave:    notes[0].GetOctave(),
		lastSample:    ct.DefaultSample,
		defaultSample: ct.DefaultSample,
		sampleMap:     ct.SampleMap,
		sb:            &strings.Builder{},
	}

	for _, smwNote := range notes {
		if rest, ok := smwNote.(Rest); ok {
			writeRest(rest, ctx)
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
	if len(note.GetLengthValues()) < 1 {
		return
	}

	for i, length := range note.GetLengthValues() {
		if i == 0 {
			write(ctx.sb, "%s%d", note.GetKey(), length)
		} else {
			write(ctx.sb, "^%d", length)
		}
	}
}

func writeRest(rest Rest, ctx *channelWriteContext) {
	lv := rest.GetLengthValues()
	if len(lv) >= 2 && lv[0] == lv[1] {
		rest = writeRestSuperLoop(rest, ctx)
	}
	writeNote(rest, ctx)
}

func writeRestSuperLoop(rest Rest, ctx *channelWriteContext) Rest {
	lv := rest.GetLengthValues()
	numLoops := 1
	for numLoops < len(lv) {
		if lv[numLoops] != lv[numLoops-1] {
			break
		}
		numLoops++
	}
	write(ctx.sb, " [[%s%d]]%d", rest.GetKey(), lv[0], numLoops)

	remainingLengths := lv[numLoops:]
	remainingRest := Rest{remainingLengths}
	return remainingRest
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
	sample := getSample(keyValue, ctx.defaultSample, ctx.sampleMap)
	if sample != ctx.lastSample {
		ctx.lastSample = sample
		write(ctx.sb, "@%d", sample)
	}
}

func getSample(keyValue uint8, defaultSample uint8, sampleMap map[uint8]uint8) uint8 {
	sample, ok := sampleMap[keyValue]
	if !ok {
		sample = defaultSample
	}
	return sample
}
