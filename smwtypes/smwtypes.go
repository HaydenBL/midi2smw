package smwtypes

var noteDict = map[uint8]string{
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

func NotesEqual(note1, note2 SmwNote) bool {
	if len(note1.GetLengthValues()) != len(note2.GetLengthValues()) {
		return false
	}
	for i := range note1.GetLengthValues() {
		if note1.GetLengthValues()[i] != note2.GetLengthValues()[i] {
			return false
		}
	}
	return note1.GetKeyValue() == note2.GetKeyValue() &&
		note1.GetKey() == note2.GetKey() &&
		note1.GetOctave() == note2.GetOctave()
}

func NoteSlicesEqual(noteSlice1, noteSlice2 []SmwNote) bool {
	if len(noteSlice1) != len(noteSlice2) {
		return false
	}
	for i := range noteSlice1 {
		if !NotesEqual(noteSlice1[i], noteSlice2[i]) {
			return false
		}
	}
	return true
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
