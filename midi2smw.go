package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Type uint8

const (
	NoteOff Type = 0
	NoteOn  Type = 1
	Other   Type = 3
)

type MidiEvent struct {
	event     Type
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

// midi events
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

func main() {
	parseFile("dean_town.mid")
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

func parseFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	defer file.Close()

	numTracks := parseHeader(file)

	for track := 0; track < int(numTracks); track++ {
		parseTrack(file)
	}

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

func parseTrack(file *os.File) {
	fmt.Println("----- TRACK FOUND -----")

	var val32 uint32 = 0
	var val16 uint16 = 0
	var eof bool = false

	// Read track header
	// First 4 bytes, file ID (always MTrk)
	binary.Read(file, binary.BigEndian, &val32)
	// Next 4 bytes are track length
	if err := binary.Read(file, binary.BigEndian, &val32); err != nil {
		eof = err == io.EOF
	}
	trackLength := val32

	var endOfTrack = false
	var previousStatus uint8 = 0

	for !endOfTrack && !eof {
		var statusTimeDelta uint32 = 0
		var status uint8 = 0

		statusTimeDelta = readValue(file)
		binary.Read(file, binary.BigEndian, &status)

		// Sometimes midi files optimize data by putting consecutive midi events with the same
		// status byte next to each other, not repeating the status bytes.
		// If we encounter a byte without the status flag set, it means we've run into this case
		// and we have to seek back one byte because it was actually an event!
		if status < 0x80 {
			status = previousStatus
			file.Seek(-1, 1) // seek back 1 byte from current position
		}

		if (status & 0xF0) == VoiceNoteOff {
			var channel, noteID, noteVelocity uint8
			previousStatus = status
			channel = status & 0x0F
			binary.Read(file, binary.BigEndian, noteID)
			if err := binary.Read(file, binary.BigEndian, noteVelocity); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == VoiceNoteOn {
			var channel, noteID, noteVelocity uint8
			previousStatus = status
			channel = status & 0x0F
			binary.Read(file, binary.BigEndian, noteID)
			if err := binary.Read(file, binary.BigEndian, noteVelocity); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == VoiceAftertouch {
			var channel, noteID, noteVelocity uint8
			previousStatus = status
			channel = status & 0x0F
			binary.Read(file, binary.BigEndian, noteID)
			if err := binary.Read(file, binary.BigEndian, noteVelocity); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == VoiceControlChange {
			var channel, noteID, noteVelocity uint8
			previousStatus = status
			channel = status & 0x0F
			binary.Read(file, binary.BigEndian, noteID)
			if err := binary.Read(file, binary.BigEndian, noteVelocity); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == VoiceProgramChange {
			var channel, programID uint8
			previousStatus = status
			channel = status & 0x0F
			if err := binary.Read(file, binary.BigEndian, programID); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == VoiceChannelPressure {
			var channel, channelPressure uint8
			previousStatus = status
			channel = status & 0x0F
			if err := binary.Read(file, binary.BigEndian, channelPressure); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == VoicePitchBend {
			var channel, LS7B, MS7B uint8
			previousStatus = status
			channel = status & 0x0F
			binary.Read(file, binary.BigEndian, LS7B)
			if err := binary.Read(file, binary.BigEndian, MS7B); err != nil {
				eof = err == io.EOF
			}

		} else if (status & 0xF0) == SystemExclusive {

		} else {
			fmt.Printf("Unrecognized status byte: %b\n", status)
		}

	}
}
