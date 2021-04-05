package smwtypes

import (
	"fmt"
	"strings"
)

type loopSection struct {
	loops int
	notes []SmwNote
}

func (ct ChannelTrack) StringCompressed() string {
	remainingNotes := ct.Notes
	sb := &strings.Builder{}
	for len(remainingNotes) > 1 {
		var section loopSection
		section, remainingNotes = getLoopSection(remainingNotes)
		ctx := &channelWriteContext{
			lastOctave:    ct.Notes[0].GetOctave(),
			lastSample:    ct.DefaultSample,
			defaultSample: ct.DefaultSample,
			sampleMap:     ct.SampleMap,
			sb:            &strings.Builder{},
		}
		sectionOutput := fmt.Sprintf("o%d %s", section.notes[0].GetOctave(), stringNotes(section.notes, ctx))
		if section.loops > 1 {
			sectionOutput = fmt.Sprintf("[%s]%d", sectionOutput, section.loops)
		}
		write(sb, sectionOutput)
	}
	return sb.String()
}

func getLoopSection(notes []SmwNote) (loopSection, []SmwNote) {
	longestLoopSection := loopSection{
		loops: 1,
		notes: []SmwNote{notes[0]},
	}
	current := []SmwNote{notes[0]}
	remainingNoteTrack := notes[1:] // second item to the end
	for len(current) <= len(notes)/2 {
		numLoops := getNumLoops(current, remainingNoteTrack)
		if numLoops > longestLoopSection.loops {
			longestLoopSection.notes = current
			longestLoopSection.loops = numLoops
		}
		current = append(current, remainingNoteTrack[0])
		remainingNoteTrack = remainingNoteTrack[1:]
	}
	return longestLoopSection, notes[len(longestLoopSection.notes)-1:]
}

func getNumLoops(notes, remainingTrack []SmwNote) int {
	numNotes := len(notes)
	count := 0
	for numNotes <= len(remainingTrack) {
		if !NoteSlicesEqual(notes, remainingTrack[:numNotes]) {
			break
		}
		remainingTrack = remainingTrack[numNotes:]
		count++
	}
	return count
}
