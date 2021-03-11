package midi

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	metaSequence          uint8 = 0x00
	metaText              uint8 = 0x01
	metaCopyright         uint8 = 0x02
	metaTrackName         uint8 = 0x03
	metaInstrumentName    uint8 = 0x04
	metaLyrics            uint8 = 0x05
	metaMarker            uint8 = 0x06
	metaCuePoint          uint8 = 0x07
	metaProgramName       uint8 = 0x08
	metaDevicePort        uint8 = 0x09
	metaChannelPrefix     uint8 = 0x20
	metaMIDIPort          uint8 = 0x21
	metaEndOfTrack        uint8 = 0x2F
	metaSetTempo          uint8 = 0x51
	metaSMPTEOffset       uint8 = 0x54
	metaTimeSignature     uint8 = 0x58
	metaKeySignature      uint8 = 0x59
	metaSequencerSpecific uint8 = 0x7F
)

func handleMetaType(file *os.File, track Track) (endOfTrack bool) {
	var metaType, length uint8

	binary.Read(file, binary.BigEndian, &metaType)
	length = uint8(readValue(file))

	switch metaType {

	case metaSequence:
		var val1, val2 uint8
		binary.Read(file, binary.BigEndian, &val1)
		binary.Read(file, binary.BigEndian, &val2)
		fmt.Printf("Sequence number: %d %d\n", val1, val2)

	case metaText:
		fmt.Printf("Meta text: %s\n", readString(file, uint32(length)))

	case metaCopyright:
		fmt.Printf("Copyright: %s\n", readString(file, uint32(length)))

	case metaTrackName:
		track.Name = readString(file, uint32(length))
		fmt.Printf("Track name: %s\n", track.Name)

	case metaInstrumentName:
		track.Instrument = readString(file, uint32(length))
		fmt.Printf("Instrument name: %s\n", track.Instrument)

	case metaLyrics:
		fmt.Printf("Lyrics: %s\n", readString(file, uint32(length)))

	case metaMarker:
		fmt.Printf("Marker: %s\n", readString(file, uint32(length)))

	case metaCuePoint:
		fmt.Printf("Cue: %s\n", readString(file, uint32(length)))

	case metaProgramName:
		fmt.Printf("Meta program name: %s\n", readString(file, uint32(length)))

	case metaDevicePort:
		fmt.Printf("Meta device port: %s\n", readString(file, uint32(length)))

	case metaChannelPrefix:
		var prefix uint8
		binary.Read(file, binary.BigEndian, &prefix)
		fmt.Printf("Prefix: %d\n", prefix)

	case metaMIDIPort:
		var port uint8
		binary.Read(file, binary.BigEndian, port)
		fmt.Printf("MIDI port: %d\n", port)

	case metaEndOfTrack:
		endOfTrack = true

	case metaSetTempo:
		// Tempo is in microseconds per quarter note
		globalTempo = 0
		var b uint8
		binary.Read(file, binary.BigEndian, &b)
		globalTempo |= uint32(b) << 16
		binary.Read(file, binary.BigEndian, &b)
		globalTempo |= uint32(b) << 8
		binary.Read(file, binary.BigEndian, &b)
		globalTempo |= uint32(b) << 0
		globalBPM = 60000000 / globalTempo

		fmt.Printf("Tempo: %d (%d bpm)\n", globalTempo, globalBPM)

	case metaSMPTEOffset:
		var h, m, s, fr, ff uint8
		binary.Read(file, binary.BigEndian, &h)
		binary.Read(file, binary.BigEndian, &m)
		binary.Read(file, binary.BigEndian, &s)
		binary.Read(file, binary.BigEndian, &fr)
		binary.Read(file, binary.BigEndian, &ff)
		fmt.Printf("SMPTE: H: %d M: %d S: %d FR: %d FF: %d\n", h, m, s, fr, ff)

	case metaTimeSignature:
		var val1, val2 uint8

		binary.Read(file, binary.BigEndian, &val1)
		binary.Read(file, binary.BigEndian, &val2)
		fmt.Printf("Time signature: %d/%d\n", val1, 2<<val2)

		binary.Read(file, binary.BigEndian, &val1)
		fmt.Printf("Clocks per tick: %d\n", val1)

		// A MIDI "Beat" is 24 ticks, so specify how many 32nd notes constitute a beat
		binary.Read(file, binary.BigEndian, &val1)
		fmt.Printf("32per24Clocks: %d\n", val1)

	case metaKeySignature:
		var keySignature, minorKey uint8
		binary.Read(file, binary.BigEndian, &keySignature)
		binary.Read(file, binary.BigEndian, &minorKey)

		fmt.Printf("Key signature: %d\n", keySignature)
		fmt.Printf("Minor key: %d\n", minorKey)

	case metaSequencerSpecific:
		fmt.Printf("Sequencer specific: %s", readString(file, uint32(length)))

	default:
		fmt.Printf("Unrecognized MetaEvent: %c\n", metaType)

	}

	return
}
