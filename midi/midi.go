package midi

import (
	"fmt"
	"os"
)

// Much of this parsing code was written with the help of OLC's midi parsing implementation in C++
// https://github.com/OneLoneCoder/olcPixelGameEngine/blob/master/Videos/OneLoneCoder_PGE_MIDI.cpp

var (
	globalTempo uint32 = 0
	globalBPM   uint32 = 0
)

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
	Events     []Event
}

const (
	NoteOff EventType = 0
	NoteOn  EventType = 1
	Other   EventType = 3
)

func Parse(fileName string) ([]Track, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return nil, err
	}
	defer file.Close()

	numTracks, err := parseHeader(file)
	if err != nil {
		return []Track{}, err
	}
	var midiTracks []Track
	for track := 0; track < int(numTracks); track++ {
		track, err := parseTrack(file)
		if err != nil {
			return []Track{}, err
		}
		midiTracks = append(midiTracks, track)
	}

	fmt.Printf("\nFound %d tracks\n", len(midiTracks))

	return midiTracks, nil
}
