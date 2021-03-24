package convert

import (
	"bufio"
	"fmt"
	"midi2smw/utils"
	"os"
	"strings"
)

func SpecifySamples(noteTracks []NoteTrack) []NoteTrack {
	sc := bufio.NewScanner(os.Stdin)

	for true {
		index := promptToSetSamples(sc, noteTracks)
		if index == -1 {
			break
		}

		sampleMap := getSampleMap(sc)
		noteTracks[index].SampleMap = sampleMap
	}
	return noteTracks
}

func promptToSetSamples(sc *bufio.Scanner, noteTracks []NoteTrack) int {
	for true {
		fmt.Println("-- Specify index of track to add samples (q to quit)")
		for i, track := range noteTracks {
			fmt.Printf("\t%d -\tName:\t\t\t%s\n", i, track.Name)
			fmt.Printf("\t\tNotes:\t\t\t%d\n\n", len(track.Notes))
		}
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "q" {
			return -1
		}

		index, err := utils.ReadInt(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		return int(index)
	}
	return -1
}

func getSampleMap(sc *bufio.Scanner) map[uint8]uint8 {
	m := make(map[uint8]uint8)
	for true {
		fmt.Printf("\t\tEnter note values, space separated (q to quit): ")
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "q" {
			return m
		}
		notes, err := utils.ReadLineOfUInt8s(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("\t\tEnter sample number: ")
		sc.Scan()
		line = sc.Text()
		sample, err := utils.ReadInt(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, note := range notes {
			m[note] = sample
		}
	}
	return m
}
