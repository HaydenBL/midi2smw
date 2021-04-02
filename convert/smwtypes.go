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
	Key() string
	Octave() uint8
	KeyValue() uint8
	LengthValues() []uint8
}

// An actual note to be played (i.e. not a rest)
// Implements the SmwNote interface
type Note struct {
	key          string
	keyValue     uint8
	lengthValues []uint8
	octave       uint8
}

func (n Note) Key() string {
	return n.key
}

func (n Note) KeyValue() uint8 {
	return n.keyValue
}

func (n Note) Octave() uint8 {
	return n.octave
}

func (n Note) LengthValues() []uint8 {
	return n.lengthValues
}

// An rest note to be played (i.e. not a rest)
// Implements the SmwNote interface
type Rest struct {
	lengthValues []uint8
}

func (r Rest) Key() string {
	return "r"
}

func (r Rest) KeyValue() uint8 {
	return 0
}

func (r Rest) Octave() uint8 {
	return 0
}

func (r Rest) LengthValues() []uint8 {
	return r.lengthValues
}
