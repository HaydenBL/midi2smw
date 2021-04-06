package trackoutput

import (
	"bufio"
	"errors"
	"fmt"
	"midi2smw/smwtypes"
	"os"
	"strconv"
	"strings"
	"sync"
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

func (p *Printer) getOutputConfig(specifyTracks, loop bool) outputConfig {
	config := outputConfig{Bpm: p.bpm}

	var channelOutputs []channelOutput
	if specifyTracks {
		channelOutputs = manuallySpecifyChannelOutputs(p.tracks, loop)
	} else {
		channelOutputs = getAllChannelOutputs(p.tracks, loop)
	}
	config.ChannelOutputs = channelOutputs

	if len(channelOutputs) > 8 {
		fmt.Printf("More than 8 channels were inserted into the output.\n")
		fmt.Printf("You will have to manually remove tracks to insert your music.\n")
	}

	return config
}

func getAllChannelOutputs(tracks []smwtypes.SmwTrack, loop bool) []channelOutput {
	namedChannels := make([]namedChannelTrack, 0)
	for _, track := range tracks {
		for _, channel := range track.ChannelTracks {
			nc := namedChannelTrack{channel, track.Name}
			namedChannels = append(namedChannels, nc)
		}
	}
	channelOutputs := namedChannelTracksToSmwOutputs(namedChannels, loop)
	return channelOutputs
}

func manuallySpecifyChannelOutputs(tracks []smwtypes.SmwTrack, loop bool) []channelOutput {
	sc := bufio.NewScanner(os.Stdin)
	namedChannels := make([]namedChannelTrack, 0)

	for true {
		writeAllTracks(os.Stdout, tracks)
		fmt.Printf("Enter track to add to output (q to quit): ")
		sc.Scan()
		line := sc.Text()
		if strings.ToLower(line) == "q" {
			break
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
		namedChannel := getNamedChannelTrackFromInput(sc, tracks[index])
		namedChannels = append(namedChannels, namedChannel)
	}
	outputs := namedChannelTracksToSmwOutputs(namedChannels, loop)
	return outputs
}

type namedChannelTrack struct {
	smwtypes.ChannelTrack
	Name string
}

func getNamedChannelTrackFromInput(sc *bufio.Scanner, track smwtypes.SmwTrack) namedChannelTrack {
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

		return namedChannelTrack{track.ChannelTracks[index], track.Name}
	}
	fmt.Println("Error getting track output")
	return namedChannelTrack{}
}

func namedChannelTracksToSmwOutputs(namedChannels []namedChannelTrack, loop bool) []channelOutput {
	wg := sync.WaitGroup{}
	outputs := make([]channelOutput, len(namedChannels))
	for i, channel := range namedChannels {
		wg.Add(1)
		go func(wg *sync.WaitGroup, nct namedChannelTrack, index int) {
			defer wg.Done()
			co := namedChannelTrackToChannelOutput(nct, loop)
			outputs[index] = co
		}(&wg, channel, i)
	}
	wg.Wait()
	return outputs
}

func namedChannelTrackToChannelOutput(nct namedChannelTrack, loop bool) channelOutput {
	var noteOutput string
	if loop {
		noteOutput = nct.StringCompressed()
	} else {
		noteOutput = nct.String()
	}

	return channelOutput{
		Name:          nct.Name,
		DefaultSample: nct.DefaultSample,
		StartOctave:   nct.Notes[0].GetOctave(),
		NoteOutput:    noteOutput,
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
