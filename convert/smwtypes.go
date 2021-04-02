package convert

import "fmt"

type SmwTrack struct {
	Name          string
	ChannelTracks []ChannelTrack // if a midi track has chords/overlapping notes, we'll throw them into multiple channels
}

type ChannelTrack struct {
	Notes         []SmwNote
	DefaultSample uint8
	SampleMap     map[uint8]uint8
}

type SmwNote interface {
	GetKey() string
	GetOctave() uint8
	GetKeyValue() uint8
	GetLengthValues() []uint8
}

// An actual note to be played (i.e. not a rest)
// Implements the SmwNote interface
type Note struct {
	KeyValue     uint8
	LengthValues []uint8
}

func (n Note) GetKey() string {
	return getKeyFromKeyValue(n.KeyValue)
}

func (n Note) GetKeyValue() uint8 {
	return n.KeyValue
}

func (n Note) GetOctave() uint8 {
	return getOctaveFromKeyValue(n.KeyValue)
}

func (n Note) GetLengthValues() []uint8 {
	return n.LengthValues
}

// A rest note to be played
// Implements the SmwNote interface
type Rest struct {
	LengthValues []uint8
}

func (r Rest) GetKey() string {
	return "r"
}

func (r Rest) GetKeyValue() uint8 {
	return 0
}

func (r Rest) GetOctave() uint8 {
	return 0
}

func (r Rest) GetLengthValues() []uint8 {
	return r.LengthValues
}

func getKeyFromKeyValue(noteValue uint8) string {
	if !noteValueWithinSmwRange(noteValue) {
		fmt.Printf("ERROR: note value %d not in SMW range. Using rest\n", noteValue)
		return "r"
	}
	return noteDict[noteValue%12]
}

func getOctaveFromKeyValue(noteValue uint8) uint8 {
	if !noteValueWithinSmwRange(noteValue) {
		fmt.Printf("ERROR: note value %d note in SMW range. Using 0\n", noteValue)
		return 0
	}
	return noteValue/12 - 1
}

func noteValueWithinSmwRange(noteValue uint8) bool {
	// Lowest SMW note is g0 == 19
	// Highest SMW note is e6 == 88
	return noteValue >= 19 && noteValue <= 88
}
