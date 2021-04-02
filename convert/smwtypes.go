package convert

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
	Key          string
	KeyValue     uint8
	LengthValues []uint8
	Octave       uint8
}

func (n Note) GetKey() string {
	return n.Key
}

func (n Note) GetKeyValue() uint8 {
	return n.KeyValue
}

func (n Note) GetOctave() uint8 {
	return n.Octave
}

func (n Note) GetLengthValues() []uint8 {
	return n.LengthValues
}

// An rest note to be played
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
