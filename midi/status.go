package midi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type scanContext struct {
	status          uint8
	previousStatus  uint8
	statusTimeDelta uint32
}

func readStatus(file *os.File, sc *scanContext) error {
	if err := binary.Read(file, binary.BigEndian, &sc.status); err != nil {
		return err
	}

	// Sometimes midi files optimize data by putting consecutive midi events with the same
	// status byte next to each other, not repeating the status bytes.
	// If we encounter a byte without the status flag set, it means we've run into this case
	// and we have to seek back one byte because it was actually an event!
	if sc.status < 0x80 {
		sc.status = sc.previousStatus
		_, err := file.Seek(-1, 1) // seek back 1 byte from current position
		if err != nil {
			fmt.Printf("Error seeking backwards in file")
			return err
		}
	}
	return nil
}

func handleStatus(file *os.File, track *Track, sc *scanContext) (*Event, error) {

	if (sc.status & 0xF0) == voiceNoteOff {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		noteID := readByte(file)
		noteVelocity := readByte(file)

		return &Event{NoteOff, noteID, noteVelocity, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceNoteOn {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		noteID := readByte(file)
		noteVelocity := readByte(file)

		if noteVelocity == 0 {
			return &Event{NoteOff, noteID, noteVelocity, sc.statusTimeDelta}, nil
		} else {
			return &Event{NoteOn, noteID, noteVelocity, sc.statusTimeDelta}, nil
		}

	} else if (sc.status & 0xF0) == voiceAftertouch {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		noteID := readByte(file)
		noteVelocity := readByte(file)
		_, _ = noteID, noteVelocity // not doing anything with these for now

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceControlChange {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		controlID := readByte(file)
		controlValue := readByte(file)
		_, _ = controlID, controlValue // not doing anything with these for now

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceProgramChange {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		programID := readByte(file)
		_ = programID // not doing anything with this for now

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceChannelPressure {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		channelPressure := readByte(file)
		_ = channelPressure // not doing anything with this for now

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voicePitchBend {
		sc.previousStatus = sc.status
		channel := sc.status & 0x0F
		_ = channel // not doing anything with channel for now

		LS7B := readByte(file)
		MS7B := readByte(file)
		_, _ = LS7B, MS7B // not doing anything with these for now

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == systemExclusive {
		sc.previousStatus = 0
		if sc.status == 0xF0 {
			fmt.Printf("System exclusive message begin: %s\n", readString(file, readValue(file)))
		}

		if sc.status == 0xF7 {
			fmt.Printf("System exclusive message end: %s\n", readString(file, readValue(file)))
		}

		if sc.status == 0xFF {
			bpm, endOfTrack := handleMetaType(file, track)
			if track.Bpm == 0 {
				track.Bpm = bpm
			}
			if endOfTrack {
				return nil, io.EOF
			}
			return &Event{Other, 0, 0, sc.statusTimeDelta}, nil
		}

	} else {
		fmt.Printf("Unrecognized status byte: %d\n", sc.status)
	}
	return nil, nil
}
