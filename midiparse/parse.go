package midiparse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type eventType uint8

type midiEvent struct {
	event     eventType
	key       uint8
	velocity  uint8
	deltaTick uint32
}

type midiNote struct {
	key       uint8
	velocity  uint8
	startTime uint32
	duration  uint32
}

type midiTrack struct {
	name       string
	instrument string
	events     []midiEvent
	notes      []midiNote
	maxNote    uint8
	minNote    uint8
}

const (
	noteOff eventType = 0
	noteOn  eventType = 1
	other   eventType = 3
)

const (
	voiceNoteOff         uint8 = 0x80
	voiceNoteOn          uint8 = 0x90
	voiceAftertouch      uint8 = 0xA0
	voiceControlChange   uint8 = 0xB0
	voiceProgramChange   uint8 = 0xC0
	voiceChannelPressure uint8 = 0xD0
	voicePitchBend       uint8 = 0xE0
	systemExclusive      uint8 = 0xF0
)

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

func parseTrack(file *os.File) midiTrack {
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

	track := midiTrack{}

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

		if (status & 0xF0) == voiceNoteOff {
			//var channel uint8
			var noteID, noteVelocity uint8
			previousStatus = status
			//channel = status & 0x0F

			binary.Read(file, binary.BigEndian, &noteID)
			if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, midiEvent{noteOff, noteID, noteVelocity, statusTimeDelta})

		} else if (status & 0xF0) == voiceNoteOn {
			//var channel uint8
			var noteID, noteVelocity uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &noteID)
			if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
				eof = err == io.EOF
			}

			if noteVelocity == 0 {
				track.events = append(track.events, midiEvent{noteOff, noteID, noteVelocity, statusTimeDelta})
			} else {
				track.events = append(track.events, midiEvent{noteOn, noteID, noteVelocity, statusTimeDelta})
			}

		} else if (status & 0xF0) == voiceAftertouch {
			//var channel uint8
			var noteID, noteVelocity uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &noteID)
			if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, midiEvent{other, 0, 0, 0})

		} else if (status & 0xF0) == voiceControlChange {
			//var channel uint8
			var controlID, controlValue uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &controlID)
			if err := binary.Read(file, binary.BigEndian, &controlValue); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, midiEvent{other, 0, 0, 0})

		} else if (status & 0xF0) == voiceProgramChange {
			//var channel uint8
			var programID uint8
			previousStatus = status
			//channel = status & 0x0F
			if err := binary.Read(file, binary.BigEndian, &programID); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, midiEvent{other, 0, 0, 0})

		} else if (status & 0xF0) == voiceChannelPressure {
			//var channel uint8
			var channelPressure uint8
			previousStatus = status
			//channel = status & 0x0F
			if err := binary.Read(file, binary.BigEndian, &channelPressure); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, midiEvent{other, 0, 0, 0})

		} else if (status & 0xF0) == voicePitchBend {
			//var channel uint8
			var LS7B, MS7B uint8
			previousStatus = status
			//channel = status & 0x0F
			binary.Read(file, binary.BigEndian, &LS7B)
			if err := binary.Read(file, binary.BigEndian, &MS7B); err != nil {
				eof = err == io.EOF
			}

			track.events = append(track.events, midiEvent{other, 0, 0, 0})

		} else if (status & 0xF0) == systemExclusive {
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
