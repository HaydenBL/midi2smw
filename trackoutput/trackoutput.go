package trackoutput

import (
	"fmt"
	"io"
	"math"
	"midi2smw/convert"
	"text/template"
)

type Printer struct {
	bpm    uint8
	tracks []convert.SmwTrack
}

func NewPrinter(tracks []convert.SmwTrack, bpm uint32) Printer {
	return Printer{
		bpm:    bpmToSmwTempo(bpm),
		tracks: tracks,
	}
}

func (p *Printer) Print(writer io.Writer) error {
	config := p.getOutputConfig()
	t := template.Must(template.ParseFiles("trackoutput/output.tmpl"))
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

func write(writer io.Writer, format string, a ...interface{}) {
	formattedString := fmt.Sprintf(format, a...)
	_, err := writer.Write([]byte(formattedString))
	if err != nil {
		fmt.Printf("Error writing string: %s", formattedString)
	}
}

func writeAllTracks(writer io.Writer, tracks []convert.SmwTrack) {
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

func writeTrack(writer io.Writer, track convert.SmwTrack) {
	for i, channel := range track.ChannelTracks {
		fmt.Printf("-- Channel %d\n", i)
		writeChannel(writer, channel)
		fmt.Println()
	}
}

func writeChannel(writer io.Writer, channel convert.ChannelTrack) {
	notes := channel.Notes
	lastOctave := notes[0].Octave
	lastSample := channel.DefaultSample

	for _, smwNote := range notes {
		if smwNote.Key == "r" {
			for i, note := range smwNote.LengthValues {
				if i == 0 {
					write(writer, "r%d", note)
				} else {
					write(writer, "^%d", note)
				}
			}
		} else {
			if smwNote.Octave > lastOctave {
				for i := 0; i < smwNote.Octave-lastOctave; i++ {
					write(writer, ">")
				}
			} else if smwNote.Octave < lastOctave {
				for i := 0; i < lastOctave-smwNote.Octave; i++ {
					write(writer, "<")
				}
			}
			for i, note := range smwNote.LengthValues {
				if i == 0 {
					// Check if we need to swap the sample
					sample, ok := channel.SampleMap[smwNote.KeyValue]
					if !ok {
						sample = channel.DefaultSample
					}
					if sample != lastSample {
						lastSample = sample
						write(writer, "@%d", sample)
					}

					write(writer, "%s%d", smwNote.Key, note)
				} else {
					write(writer, "^%d", note)
				}
			}
			lastOctave = smwNote.Octave
		}
	}
	write(writer, "\n")
}
