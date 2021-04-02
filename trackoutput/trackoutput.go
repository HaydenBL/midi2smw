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

func (p *Printer) Print(writer io.Writer, specifyTracks bool) error {
	config := p.getOutputConfig(specifyTracks)
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

func write(writer io.Writer, format string, a ...interface{}) {
	formattedString := fmt.Sprintf(format, a...)
	_, err := writer.Write([]byte(formattedString))
	if err != nil {
		fmt.Printf("Error writing string: %s", formattedString)
	}
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
		writeChannel(writer, channel)
		fmt.Println()
	}
}

func writeChannel(writer io.Writer, channel smwtypes.ChannelTrack) {
	notes := channel.Notes
	lastOctave := notes[0].GetOctave()
	lastSample := channel.DefaultSample

	for _, smwNote := range notes {
		if rest, ok := smwNote.(smwtypes.Rest); ok {
			for i, note := range rest.LengthValues {
				if i == 0 {
					write(writer, "r%d", note)
				} else {
					write(writer, "^%d", note)
				}
			}
		} else {
			if smwNote.GetOctave() > lastOctave {
				for i := uint8(0); i < smwNote.GetOctave()-lastOctave; i++ {
					write(writer, ">")
				}
			} else if smwNote.GetOctave() < lastOctave {
				for i := uint8(0); i < lastOctave-smwNote.GetOctave(); i++ {
					write(writer, "<")
				}
			}
			for i, note := range smwNote.GetLengthValues() {
				if i == 0 {
					// Check if we need to swap the sample
					sample, ok := channel.SampleMap[smwNote.GetKeyValue()]
					if !ok {
						sample = channel.DefaultSample
					}
					if sample != lastSample {
						lastSample = sample
						write(writer, "@%d", sample)
					}

					write(writer, "%s%d", smwNote.GetKey(), note)
				} else {
					write(writer, "^%d", note)
				}
			}
			lastOctave = smwNote.GetOctave()
		}
	}
}
