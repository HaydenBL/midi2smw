package midi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
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

func parseHeader(file *os.File) (numTrackChunks uint16, ticksPer64thNote uint32, err error) {
	var val32 uint32 = 0
	var val16 uint16 = 0

	// First 4 bytes, file ID (always MThd)
	if err := binary.Read(file, binary.BigEndian, &val32); err != nil {
		return 0, 0, err
	}
	// Next 4 bytes, length of header
	if err := binary.Read(file, binary.BigEndian, &val32); err != nil {
		return 0, 0, err
	}
	// Next 2 bytes, format details
	if err := binary.Read(file, binary.BigEndian, &val16); err != nil {
		return 0, 0, err
	}
	// Next 2 bytes, number of tracks
	if err := binary.Read(file, binary.BigEndian, &numTrackChunks); err != nil {
		return 0, 0, err
	}
	var division uint16
	// Next 2 bytes, time division
	if err := binary.Read(file, binary.BigEndian, &division); err != nil {
		return 0, 0, err
	}

	// If bit 15 is zero, bits 0-14 is ticks per quarter note
	if division&0x8000 == 0 {
		ticksPerQuarterNote := division & 0x7FFF
		ticksPer64thNote = uint32(ticksPerQuarterNote / 16)
	} else {
		return 0, 0, errors.New("unsupported time format")
	}

	return numTrackChunks, ticksPer64thNote, nil
}

func parseTrack(file *os.File) (Track, error) {
	fmt.Println("----- TRACK FOUND")

	var track = &Track{Name: "Unnamed Track"}
	var sc = &scanContext{}
	var MTrk uint32 = 0
	var endOfTrack = false

	// Read track header
	// First 4 bytes, file ID (always MTrk)
	binary.Read(file, binary.BigEndian, &MTrk)
	// Next 4 bytes are track length
	if err := binary.Read(file, binary.BigEndian, &track.Length); err != nil {
		if err == io.EOF {
			endOfTrack = true
		} else {
			return *track, err
		}
	}

	for !endOfTrack {
		var event *Event = nil
		var err error

		sc.statusTimeDelta = readValue(file)
		if err = readStatus(file, sc); err != nil {
			if err == io.EOF {
				endOfTrack = true
			} else {
				return *track, err
			}
		}
		if event, err = handleStatus(file, track, sc); err != nil {
			if err == io.EOF {
				endOfTrack = true
			} else {
				return *track, err
			}
		}
		if event != nil {
			track.Events = append(track.Events, *event)
		}

	}

	return *track, nil
}
