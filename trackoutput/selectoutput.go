package trackoutput

import (
	"bufio"
	"errors"
	"fmt"
	"midi2smw/convert"
	"os"
	"strconv"
	"strings"
)

type outputConfig struct {
	bpm          uint8
	trackOutputs []trackOutput
}

type trackOutput struct {
	name          string
	defaultSample uint8
	startOctave   int
	noteOutput    string
}

func getOutputConfig(p Printer) outputConfig {
	tracks := p.tracks
	config := outputConfig{bpm: p.bpm}
	sc := bufio.NewScanner(os.Stdin)

	for i := 0; i < 8; i++ {
		writeAllTracks(os.Stdout, tracks)
		fmt.Printf("Enter track to add to output (q to quit): ")
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "q" {
			return config
		}

		index, err := readInt(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if index < 0 || int(index) >= len(tracks) {
			fmt.Println("Index out of range")
			continue
		}
		config.trackOutputs = append(config.trackOutputs, getTrackOutput(sc, tracks[index]))
	}

	return config
}

func getTrackOutput(sc *bufio.Scanner, track convert.SmwTrack) trackOutput {
	to := trackOutput{name: track.Name}
	for true {
		fmt.Printf("Track %s:\n", track.Name)
		writeTrack(os.Stdout, track)
		fmt.Printf("Enter channel to add to output: ")
		sc.Scan()
		line := sc.Text()
		index, err := readInt(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if index < 0 || int(index) >= len(track.ChannelTracks) {
			fmt.Println("Index out of range")
			continue
		}

		channel := track.ChannelTracks[index]
		sb := strings.Builder{}
		writeChannel(&sb, channel)
		to.startOctave = channel.Notes[0].Octave
		to.defaultSample = channel.DefaultSample
		to.noteOutput = sb.String()
		return to
	}
	fmt.Println("Error getting track output")
	return to
}

func readInt(str string) (uint8, error) {
	var num64 uint64
	var err error

	if num64, err = strconv.ParseUint(str, 10, 8); err != nil {
		return 0, errors.New("error parsing line")
	}
	if num64 > 255 {
		return 0, errors.New("number too large (max 255)")
	}
	return uint8(num64), nil
}
