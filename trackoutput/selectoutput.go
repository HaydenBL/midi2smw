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
	Bpm            uint8
	ChannelOutputs []channelOutput
}

type channelOutput struct {
	Name          string
	DefaultSample uint8
	StartOctave   int
	NoteOutput    string
}

func (p *Printer) getOutputConfig() outputConfig {
	tracks := p.tracks
	config := outputConfig{Bpm: p.bpm}
	config.ChannelOutputs = manuallySpecifyChannelOutputs(tracks)
	return config
}

func manuallySpecifyChannelOutputs(tracks []convert.SmwTrack) []channelOutput {
	sc := bufio.NewScanner(os.Stdin)
	outputs := make([]channelOutput, 0)

	for true {
		writeAllTracks(os.Stdout, tracks)
		fmt.Printf("Enter track to add to output (q to quit): ")
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "q" {
			return outputs
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
		channel := getChannelOutput(sc, tracks[index])
		outputs = append(outputs, channel)
	}

	return outputs
}

func getChannelOutput(sc *bufio.Scanner, track convert.SmwTrack) channelOutput {
	to := channelOutput{Name: track.Name}
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
		to.StartOctave = channel.Notes[0].Octave
		to.DefaultSample = channel.DefaultSample
		to.NoteOutput = sb.String()
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
