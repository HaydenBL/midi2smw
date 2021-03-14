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
		file.Seek(-1, 1) // seek back 1 byte from current position
	}
	return nil
}

func handleStatus(file *os.File, track *Track, sc *scanContext) (*Event, error) {

	if (sc.status & 0xF0) == voiceNoteOff {
		//var channel uint8
		var noteID, noteVelocity uint8
		sc.previousStatus = sc.status
		//channel = status & 0x0F

		binary.Read(file, binary.BigEndian, &noteID)
		if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
			return nil, err
		}

		return &Event{NoteOff, noteID, noteVelocity, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceNoteOn {
		//var channel uint8
		var noteID, noteVelocity uint8
		sc.previousStatus = sc.status
		//channel = sc.status & 0x0F
		binary.Read(file, binary.BigEndian, &noteID)
		if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
			return nil, err
		}

		if noteVelocity == 0 {
			return &Event{NoteOff, noteID, noteVelocity, sc.statusTimeDelta}, nil
		} else {
			return &Event{NoteOn, noteID, noteVelocity, sc.statusTimeDelta}, nil
		}

	} else if (sc.status & 0xF0) == voiceAftertouch {
		//var channel uint8
		var noteID, noteVelocity uint8
		sc.previousStatus = sc.status
		//channel = sc.status & 0x0F
		binary.Read(file, binary.BigEndian, &noteID)
		if err := binary.Read(file, binary.BigEndian, &noteVelocity); err != nil {
			return nil, err
		}

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceControlChange {
		//var channel uint8
		var controlID, controlValue uint8
		sc.previousStatus = sc.status
		//channel = sc.status & 0x0F
		binary.Read(file, binary.BigEndian, &controlID)
		if err := binary.Read(file, binary.BigEndian, &controlValue); err != nil {
			return nil, err
		}

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceProgramChange {
		//var channel uint8
		var programID uint8
		sc.previousStatus = sc.status
		//channel = sc.status & 0x0F
		if err := binary.Read(file, binary.BigEndian, &programID); err != nil {
			return nil, err
		}

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voiceChannelPressure {
		//var channel uint8
		var channelPressure uint8
		sc.previousStatus = sc.status
		//channel = sc.status & 0x0F
		if err := binary.Read(file, binary.BigEndian, &channelPressure); err != nil {
			return nil, err
		}

		return &Event{Other, 0, 0, sc.statusTimeDelta}, nil

	} else if (sc.status & 0xF0) == voicePitchBend {
		//var channel uint8
		var LS7B, MS7B uint8
		sc.previousStatus = sc.status
		//channel = sc.status & 0x0F
		binary.Read(file, binary.BigEndian, &LS7B)
		if err := binary.Read(file, binary.BigEndian, &MS7B); err != nil {
			return nil, err
		}

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
