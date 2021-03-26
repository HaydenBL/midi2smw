package midi

import (
	"fmt"
	"os"
)

// Much of this parsing code was written with the help of OLC's midi parsing implementation in C++
// https://github.com/OneLoneCoder/olcPixelGameEngine/blob/master/Videos/OneLoneCoder_PGE_MIDI.cpp

type EventType uint8

type Event struct {
	Event     EventType
	Key       uint8
	Velocity  uint8
	DeltaTick uint32
}

type Track struct {
	Name       string
	Instrument string
	Length     uint32
	Bpm        uint32
	Events     []Event
}

type File struct {
	MidiTracks       []Track
	Bpm              uint32
	TicksPer64thNote uint32
}

const (
	NoteOff EventType = 0
	NoteOn  EventType = 1
	Other   EventType = 2
)

func Parse(fileName string) (File, error) {
	var mf File

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return File{}, err
	}
	defer file.Close()

	var numTracks uint16
	numTracks, mf.TicksPer64thNote, err = parseHeader(file)
	if err != nil {
		return File{}, err
	}

	for i := 0; i < int(numTracks); i++ {
		track, err := parseTrack(file)
		if err != nil {
			return File{}, err
		}
		if mf.Bpm == 0 {
			mf.Bpm = track.Bpm
		}
		mf.MidiTracks = append(mf.MidiTracks, track)
	}

	fmt.Printf("\nFound %d tracks\n", len(mf.MidiTracks))
	fmt.Printf("%dBPM\n", mf.Bpm)
	fmt.Printf("Ticks per 64th note: %d\n", mf.TicksPer64thNote)

	return mf, nil
}
