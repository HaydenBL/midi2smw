package trackoutput

import (
	"fmt"
	"io"
	"math"
	"midi2smw/smwtypes"
	"text/template"
)

type Printer struct {
	bpm    uint8
	tracks []smwtypes.SmwTrack
}

func NewPrinter(tracks []smwtypes.SmwTrack, bpm uint32) Printer {
	return Printer{
		bpm:    bpmToSmwTempo(bpm),
		tracks: tracks,
	}
}

func (p *Printer) Print(writer io.Writer, specifyTracks, loop bool) error {
	config := p.getOutputConfig(specifyTracks, loop)
	t := template.Must(template.New("output").Parse(outputTemplate))
	if err := t.Execute(writer, config); err != nil {
		return err
	}
	return nil
}

func bpmToSmwTempo(bpm uint32) uint8 {
	const multiplier = float64(256) / 625
	tempo := math.Round(float64(bpm) * multiplier)
	return uint8(tempo)
}

func writeAllTracks(writer io.Writer, tracks []smwtypes.SmwTrack) {
	for i, track := range tracks {
		fmt.Printf("---- Track %d", i)
		if track.Name != "" {
			fmt.Printf(" (%s)", track.Name)
		}
		fmt.Printf("\n\n")
		writeTrack(writer, track)
		fmt.Println()
	}
}

func writeTrack(writer io.Writer, track smwtypes.SmwTrack) {
	for i, channel := range track.ChannelTracks {
		fmt.Printf("-- Channel %d\n", i)
		if _, err := writer.Write([]byte(channel.String())); err != nil {
			fmt.Printf("Error writing track: %s", err)
		}
		fmt.Println()
	}
}
