package midiparse

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

type MidiEvent struct {
	Event     EventType
	Key       uint8
	Velocity  uint8
	DeltaTick uint32
}

type MidiTrack struct {
	Name       string
	Instrument string
	Events     []MidiEvent
}

const (
	NoteOff EventType = 0
	NoteOn  EventType = 1
	Other   EventType = 3
)

func Parse(fileName string) ([]MidiTrack, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return nil, err
	}
	defer file.Close()

	numTracks := parseHeader(file)

	var midiTracks []MidiTrack
	for track := 0; track < int(numTracks); track++ {
		track := parseTrack(file)
		midiTracks = append(midiTracks, track)
	}

	fmt.Printf("\nFound %d tracks\n", len(midiTracks))

	return midiTracks, nil
}
