package smwtypes

import "fmt"

type NoteGenerator interface {
	NewNote(keyValue uint8, noteLength uint32) SmwNote
	NewRest(restLength uint32) SmwNote
}

type noteGenerator struct {
	noteLengthConverter func(duration uint32) (lengths []uint8)
}

func NewNoteGenerator(ticksPer64thNote uint32) NoteGenerator {
	var ng noteGenerator
	ng.noteLengthConverter = getNoteLengthConverter(ticksPer64thNote)
	return ng
}

func (ng noteGenerator) NewNote(keyValue uint8, noteLength uint32) SmwNote {
	lengths := ng.noteLengthConverter(noteLength)
	var smwNote SmwNote
	if NoteValueWithinSmwRange(keyValue) {
		smwNote = Note{keyValue, lengths}
	} else {
		smwNote = Rest{lengths}
	}
	return smwNote
}

func (ng noteGenerator) NewRest(restLength uint32) SmwNote {
	lengths := ng.noteLengthConverter(restLength)
	return Rest{lengths}
}

func NoteValueWithinSmwRange(keyValue uint8) bool {
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	return keyValue >= 19 && keyValue <= 88
}

func getKeyFromKeyValue(keyValue uint8) string {
	if !NoteValueWithinSmwRange(keyValue) {
		fmt.Printf("ERROR: note value %d not in SMW range. Using rest\n", keyValue)
		return "r"
	}
	return noteDict[keyValue%12]
}

func getOctaveFromKeyValue(keyValue uint8) uint8 {
	if !NoteValueWithinSmwRange(keyValue) {
		fmt.Printf("ERROR: note value %d note in SMW range. Using 0\n", keyValue)
		return 0
	}
	return keyValue/12 - 1
}

func getNoteLengthConverter(ticksPer64thNote uint32) func(duration uint32) (lengths []uint8) {
	ticksPer32ndNote := ticksPer64thNote * 2
	ticksPer16thNote := ticksPer32ndNote * 2
	ticksPer8thNote := ticksPer16thNote * 2
	ticksPerQuarterNote := ticksPer8thNote * 2
	ticksPerHalfNote := ticksPerQuarterNote * 2
	ticksPerWholeNote := ticksPerHalfNote * 2

	noteLengthToSmwLength := func(duration uint32) (uint8, uint32) {
		if duration > ticksPerWholeNote {
			return 1, duration - ticksPerWholeNote
		}
		num64thNotes := duration / ticksPer64thNote
		half := num64thNotes
		acc := 0
		for half != 1 {
			half = half / 2
			acc++
		}
		switch acc {
		case 0:
			return 64, duration - ticksPer64thNote
		case 1:
			return 32, duration - ticksPer32ndNote
		case 2:
			return 16, duration - ticksPer16thNote
		case 3:
			return 8, duration - ticksPer8thNote
		case 4:
			return 4, duration - ticksPerQuarterNote
		case 5:
			return 2, duration - ticksPerHalfNote
		case 6:
			return 1, duration - ticksPerWholeNote
		}
		fmt.Println("You shouldn't be here!")
		return 0, 0
	}

	return func(duration uint32) []uint8 {
		lengths := make([]uint8, 0)
		if duration == 0 {
			return lengths
		}
		length, remainder := noteLengthToSmwLength(duration)
		lengths = append(lengths, length)
		for remainder != 0 {
			length, remainder = noteLengthToSmwLength(remainder)
			lengths = append(lengths, length)
		}
		return lengths
	}
}
