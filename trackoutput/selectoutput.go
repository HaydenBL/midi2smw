package trackoutput

import (
	"bufio"
	"errors"
	"fmt"
	"midi2smw/smwtypes"
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
	StartOctave   uint8
	NoteOutput    string
}

func (p *Printer) getOutputConfig(specifyTracks bool) outputConfig {
	config := outputConfig{Bpm: p.bpm}

	var channelOutputs []channelOutput
	if specifyTracks {
		channelOutputs = manuallySpecifyChannelOutputs(p.tracks)
	} else {
		channelOutputs = getAllChannelOutputs(p.tracks)
	}
	config.ChannelOutputs = channelOutputs

	if len(channelOutputs) > 8 {
		fmt.Printf("More than 8 channels were inserted into the output.\n")
		fmt.Printf("You will have to manually remove tracks to insert your music.\n")
	}

	return config
}

func getAllChannelOutputs(tracks []smwtypes.SmwTrack) []channelOutput {
	channelOutputs := make([]channelOutput, 0)
	for _, track := range tracks {
		for _, channel := range track.ChannelTracks {
			co := smwChannelTrackToTrackOutput(channel, track.Name)
			channelOutputs = append(channelOutputs, co)
		}
	}
	return channelOutputs
}

func manuallySpecifyChannelOutputs(tracks []smwtypes.SmwTrack) []channelOutput {
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

func getChannelOutput(sc *bufio.Scanner, track smwtypes.SmwTrack) channelOutput {
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
		co := smwChannelTrackToTrackOutput(channel, track.Name)
		return co
	}
	fmt.Println("Error getting track output")
	return channelOutput{}
}

func smwChannelTrackToTrackOutput(channelTrack smwtypes.ChannelTrack, name string) channelOutput {
	return channelOutput{
		Name:          name,
		DefaultSample: channelTrack.DefaultSample,
		StartOctave:   channelTrack.Notes[0].GetOctave(),
		NoteOutput:    channelTrack.String(),
	}
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
