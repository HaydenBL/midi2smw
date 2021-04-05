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
	sb := &strings.Builder{}
	remainingNotes := ct.Notes
	previousSectionLastOctave := remainingNotes[0].GetOctave()
	previousSectionLastSample := ct.DefaultSample

	for len(remainingNotes) > 0 {
		var section loopSection
		section, remainingNotes = getLoopSection(remainingNotes)

		sectionFirstNote := section.notes[0]
		sectionLastNote := section.notes[len(section.notes)-1]

		sectionFirstSample := getSample(sectionFirstNote.GetKeyValue(), ct.DefaultSample, ct.SampleMap)
		sectionLastSample := getSample(sectionLastNote.GetKeyValue(), ct.DefaultSample, ct.SampleMap)

		sectionOutput := ct.stringNotes(section.notes)
		if sectionFirstSample != previousSectionLastSample || sectionFirstSample != sectionLastSample {
			sectionOutput = fmt.Sprintf("@%d %s", sectionFirstSample, sectionOutput)
		}
		if sectionFirstNote.GetOctave() != previousSectionLastOctave || sectionFirstNote.GetOctave() != sectionLastNote.GetOctave() {
			sectionOutput = fmt.Sprintf("o%d %s", section.notes[0].GetOctave(), sectionOutput)
		}
		if section.loops > 1 {
			sectionOutput = fmt.Sprintf("[%s]%d", sectionOutput, section.loops)
		}
		write(sb, fmt.Sprintf("%s ", sectionOutput))

		previousSectionLastOctave = sectionLastNote.GetOctave()
		previousSectionLastSample = sectionLastSample
	}
	return sb.String()
}

func getLoopSection(notes []SmwNote) (loopSection, []SmwNote) {
	longestLoopSection := loopSection{
		loops: 1,
		notes: []SmwNote{notes[0]},
	}
	currentSliceLength := 0
	for currentSliceLength <= len(notes)/2 {
		numLoops := getNumLoops(currentSliceLength, notes)
		if numLoops > longestLoopSection.loops {
			newLongestLoopNotes := notes[:currentSliceLength]
			longestLoopSection.notes = newLongestLoopNotes
			longestLoopSection.loops = numLoops
		}
		currentSliceLength++
	}
	notesToRemove := len(longestLoopSection.notes) * longestLoopSection.loops
	return longestLoopSection, notes[notesToRemove:]
}

func getNumLoops(sliceLength int, track []SmwNote) int {
	if sliceLength < 1 {
		return 0
	}
	sliceToCompare := track[:sliceLength]
	remainingTrack := track
	count := 0
	for sliceLength <= len(remainingTrack) {
		if !NoteSlicesEqual(sliceToCompare, remainingTrack[:sliceLength]) {
			break
		}
		remainingTrack = remainingTrack[sliceLength:]
		count++
	}
	return count
}
