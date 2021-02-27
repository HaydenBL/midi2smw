package convert

import "fmt"

var noteDict = map[int]string{
	0:  "c",
	1:  "c+",
	2:  "d",
	3:  "d+",
	4:  "e",
	5:  "f",
	6:  "f+",
	7:  "g",
	8:  "g+",
	9:  "a",
	10: "a+",
	11: "b",
}

// SMW tempo conversion formula: BPM * (256/625)
func createSmwChannelTracks(noteTracks []noteTrack) {

	for _, noteTrack := range noteTracks {
		for _, note := range noteTrack.Notes {
			// TODO
			noteValueToKey(note.Key)
		}
	}
}

func noteValueToKey(noteValue uint8) string {
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	if noteValue < 19 || noteValue > 88 {
		fmt.Printf("Error, invalid note value: %d\n", noteValue)
		return ""
	}

	note := int(noteValue % 12)
	octave := int(noteValue/12) - 1
	return fmt.Sprintf("%s%d", noteDict[note], octave)
}
