package midiparse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var (
	globalTempo uint32 = 0
	globalBPM   uint32 = 0
)

type EventType uint8

const (
	NoteOff EventType = 0
	NoteOn  EventType = 1
	Other   EventType = 3
)

const (
	MetaSequence          uint8 = 0x00
	MetaText              uint8 = 0x01
	MetaCopyright         uint8 = 0x02
	MetaTrackName         uint8 = 0x03
	MetaInstrumentName    uint8 = 0x04
	MetaLyrics            uint8 = 0x05
	MetaMarker            uint8 = 0x06
	MetaCuePoint          uint8 = 0x07
	MetaChannelPrefix     uint8 = 0x20
	MetaEndOfTrack        uint8 = 0x2F
	MetaSetTempo          uint8 = 0x51
	MetaSMPTEOffset       uint8 = 0x54
	MetaTimeSignature     uint8 = 0x58
	MetaKeySignature      uint8 = 0x59
	MetaSequencerSpecific uint8 = 0x7F
)

const (
	VoiceNoteOff         uint8 = 0x80
	VoiceNoteOn          uint8 = 0x90
	VoiceAftertouch      uint8 = 0xA0
	VoiceControlChange   uint8 = 0xB0
	VoiceProgramChange   uint8 = 0xC0
	VoiceChannelPressure uint8 = 0xD0
	VoicePitchBend       uint8 = 0xE0
	SystemExclusive      uint8 = 0xF0
)

type MidiEvent struct {
	event     EventType
	key       uint8
	velocity  uint8
	deltaTick uint32
}

type MidiNote struct {
	key       uint8
	velocity  uint8
	startTime uint32
	duration  uint32
}

type MidiTrack struct {
	name       string
	instrument string
	events     []MidiEvent
	notes      []MidiNote
	maxNote    uint8
	minNote    uint8
}

func Parse(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	defer file.Close()

	numTracks := parseHeader(file)

	var midiTracks []MidiTrack
	for track := 0; track < int(numTracks); track++ {
		track := parseTrack(file)
		midiTracks = append(midiTracks, track)
	}
}

func readString(file *os.File, length uint32) string {
	b := make([]byte, length)
	n, _ := file.Read(b)
	str := string(b[:n])
	return str
}

// Values are chained together by using the most significant bit as a flag, indicating
// whether or not another byte should be read. The lower 7 bits contain the actual data
// and we'll just shift them all into a 32 bit integer while the flag is set
func readValue(file *os.File) uint32 {
	var finalValue uint32 = 0
	var aByte uint8 = 0

	binary.Read(file, binary.BigEndian, &aByte)
	finalValue = uint32(aByte)

	// If MSB is set, we need to read more bytes in
	if (finalValue & 0x80) != 0 {
		finalValue &= 0x7F                             // Extract bottom 7 bits of read byte
		for ok := true; ok; ok = (aByte & 0x80) != 0 { // Loop while MSB is 1
			// Read next byte
			binary.Read(file, binary.BigEndian, &aByte)

			// Shift 7 bits in, apply value from last byte read into their position
			finalValue = (finalValue << 7) | (uint32(aByte) & 0x7F)
		}
	}

	return finalValue
}

func parseHeader(file *os.File) (numTrackChunks uint16) {
	var val32 uint32 = 0
	var val16 uint16 = 0

	// First 4 bytes, file ID (always MThd)
	binary.Read(file, binary.BigEndian, &val32)
	// Next 4 bytes, length of header
	binary.Read(file, binary.BigEndian, &val32)
	// Next 2 bytes, format details
	binary.Read(file, binary.BigEndian, &val16)
	// Next 2 bytes, number of tracks
	binary.Read(file, binary.BigEndian, &val16)
	numTrackChunks = val16
	// Next 2 bytes, division
	binary.Read(file, binary.BigEndian, &val16)

	return numTrackChunks
}

func parseTrack(file *os.File) MidiTrack {
	fmt.Println("----- TRACK FOUND -----")

	var val32 uint32 = 0
	var eof = false

	// Read track header
	// First 4 bytes, file ID (always MTrk)
	binary.Read(file, binary.BigEndian, &val32)
	// Next 4 bytes are track length
	if err := binary.Read(file, binary.BigEndian, &val32); err != nil {
		eof = err == io.EOF
	}

	var endOfTrack = false
	var previousStatus uint8 = 0

	track := MidiTrack{}

	for !endOfTrack && !eof {
		var statusTimeDelta uint32 = 0
		var status uint8 = 0

		statusTimeDelta = readValue(file)
		if err := binary.Read(file, binary.BigEndian, &status); err != nil {
			eof = err == io.EOF
		}

		// Sometimes midi files optimize data by putting consecutive midi events with the same
		// status byte next to each other, not repeating the status bytes.
		// If we encounter a byte without the status flag set, it means we've run into this case
		// and we have to seek back one byte because it was actually an event!
		if status < 0x80 {
			status = previousStatus
			file.Seek(-1, 1) // seek back 1 byte from current position
		}

		if (status & 0xF0) == VoiceNoteOff {
			//var channel uint8
			var noteID, noteVelocity uint8
			previousStatus = status
			//channel = status & 0x0F

			binary.Read(file, binary.BigEndian, &noteID)
			if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, MidiEvent{NoteOff, noteID, noteVelocity, statusTimeDelta})

		} else if (status & 0xF0) == VoiceNoteOn {
			//var channel uint8
			var noteID, noteVelocity uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &noteID)
			if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
				eof = err == io.EOF
			}

			if noteVelocity == 0 {
				track.events = append(track.events, MidiEvent{NoteOff, noteID, noteVelocity, statusTimeDelta})
			} else {
				track.events = append(track.events, MidiEvent{NoteOn, noteID, noteVelocity, statusTimeDelta})
			}

		} else if (status & 0xF0) == VoiceAftertouch {
			//var channel uint8
			var noteID, noteVelocity uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &noteID)
			if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, MidiEvent{Other, 0, 0, 0})

		} else if (status & 0xF0) == VoiceControlChange {
			//var channel uint8
			var controlID, controlValue uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &controlID)
			if err := binary.Read(file, binary.BigEndian, &controlValue); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, MidiEvent{Other, 0, 0, 0})

		} else if (status & 0xF0) == VoiceProgramChange {
			//var channel uint8
			var programID uint8
			previousStatus = status
			//channel = status & 0x0F
			if err := binary.Read(file, binary.BigEndian, &programID); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, MidiEvent{Other, 0, 0, 0})

		} else if (status & 0xF0) == VoiceChannelPressure {
			//var channel uint8
			var channelPressure uint8
			previousStatus = status
			//channel = status & 0x0F
			if err := binary.Read(file, binary.BigEndian, &channelPressure); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, MidiEvent{Other, 0, 0, 0})

		} else if (status & 0xF0) == VoicePitchBend {
			//var channel uint8
			var LS7B, MS7B uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &LS7B)
			if err := binary.Read(file, binary.BigEndian, &MS7B); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, MidiEvent{Other, 0, 0, 0})

		} else if (status & 0xF0) == SystemExclusive {
			previousStatus = 0
			if status == 0xF0 {
				fmt.Printf("System exclusive message begin: %s\n", readString(file, readValue(file)))
			}

			if status == 0xF7 {
				fmt.Printf("System exclusive message end: %s\n", readString(file, readValue(file)))
			}

			if status == 0xFF {
				endOfTrack = handleMetaType(file, track)
			}

		} else {
			fmt.Printf("Unrecognized status byte: %d\n", status)
		}

	}

	return track
}

func handleMetaType(file *os.File, track MidiTrack) (endOfTrack bool) {
	var metaType, length uint8

	binary.Read(file, binary.BigEndian, &metaType)
	length = uint8(readValue(file))

	switch metaType {

	case MetaSequence:
		var val1, val2 uint8
		binary.Read(file, binary.BigEndian, &val1)
		binary.Read(file, binary.BigEndian, &val2)
		fmt.Printf("Sequence number: %d %d\n", val1, val2)

	case MetaText:
		fmt.Printf("Meta text: %s\n", readString(file, uint32(length)))

	case MetaCopyright:
		fmt.Printf("Copyright: %s\n", readString(file, uint32(length)))

	case MetaTrackName:
		track.name = readString(file, uint32(length))
		fmt.Printf("Track name: %s\n", track.name)

	case MetaInstrumentName:
		track.instrument = readString(file, uint32(length))
		fmt.Printf("Instrument name: %s\n", track.instrument)

	case MetaLyrics:
		fmt.Printf("Lyrics: %s\n", readString(file, uint32(length)))

	case MetaMarker:
		fmt.Printf("Marker: %s\n", readString(file, uint32(length)))

	case MetaCuePoint:
		fmt.Printf("Cue: %s\n", readString(file, uint32(length)))

	case MetaChannelPrefix:
		var prefix uint8
		binary.Read(file, binary.BigEndian, &prefix)
		fmt.Printf("Prefix: %d\n", prefix)

	case MetaEndOfTrack:
		endOfTrack = true

	case MetaSetTempo:
		// Tempo is in microseconds per quarter note
		if globalTempo == 0 {
			var b uint8
			binary.Read(file, binary.BigEndian, &b)
			globalTempo |= uint32(b) << 16
			binary.Read(file, binary.BigEndian, &b)
			globalTempo |= uint32(b) << 8
			binary.Read(file, binary.BigEndian, &b)
			globalTempo |= uint32(b) << 0
			globalBPM = 60000000 / globalTempo

			fmt.Printf("Tempo: %d (%d bpm)\n", globalTempo, globalBPM)
		}

	case MetaSMPTEOffset:
		var h, m, s, fr, ff uint8
		binary.Read(file, binary.BigEndian, &h)
		binary.Read(file, binary.BigEndian, &m)
		binary.Read(file, binary.BigEndian, &s)
		binary.Read(file, binary.BigEndian, &fr)
		binary.Read(file, binary.BigEndian, &ff)
		fmt.Printf("SMPTE: H: %d M: %d S: %d FR: %d FF: %d\n", h, m, s, fr, ff)

	case MetaTimeSignature:
		var val1, val2 uint8

		binary.Read(file, binary.BigEndian, &val1)
		binary.Read(file, binary.BigEndian, &val2)
		fmt.Printf("Time signature: %d/%d\n", val1, 2<<val2)

		binary.Read(file, binary.BigEndian, &val1)
		fmt.Printf("Clocks per tick: %d\n", val1)

		// A MIDI "Beat" is 24 ticks, so specify how many 32nd notes constitute a beat
		binary.Read(file, binary.BigEndian, &val1)
		fmt.Printf("32per24Clocks: %d\n", val1)

	case MetaKeySignature:
		var keySignature, minorKey uint8
		binary.Read(file, binary.BigEndian, &keySignature)
		binary.Read(file, binary.BigEndian, &minorKey)

		fmt.Printf("Key signature: %d\n", keySignature)
		fmt.Printf("Minor key: %d\n", minorKey)

	case MetaSequencerSpecific:
		fmt.Printf("Sequencer specific: %s", readString(file, uint32(length)))

	default:
		fmt.Printf("Unrecognized MetaEvent: %c\n", metaType)

	}

	return
}
