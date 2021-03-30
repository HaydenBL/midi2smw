package main

import (
	"flag"
	"fmt"
	"log"
	"midi2smw/convert"
	"midi2smw/midi"
	"midi2smw/trackoutput"
	"os"
)

func main() {
	flags := parseFlags()
	begin(flags)
}

type flagData struct {
	inputFileName   string
	outputFileName  string
	splitTracksFlag bool
	samplesFlag     bool
}

func parseFlags() flagData {
	splitTracksFlagPtr := flag.Bool("split", false, "Specify tracks to split with note groupings")
	samplesFlagPtr := flag.Bool("samples", false, "Specify samples for notes")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Error: no file name provided")
		os.Exit(1)
	}
	var outputFileName string
	if len(flag.Args()) > 1 {
		outputFileName = flag.Args()[1]
	}
	return flagData{
		inputFileName:   flag.Args()[0],
		outputFileName:  outputFileName,
		splitTracksFlag: *splitTracksFlagPtr,
		samplesFlag:     *samplesFlagPtr,
	}
}

func begin(flags flagData) {

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiFile, err := midi.Parse(flags.inputFileName)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", flags.inputFileName)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiFile, flags.splitTracksFlag, flags.samplesFlag)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	outputFileName := "output.txt"
	if flags.outputFileName != "" {
		outputFileName = flags.outputFileName
	}
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Error creating file")
	}
	defer outputFile.Close()

	trackPrinter := trackoutput.NewPrinter(tracks, midiFile.Bpm)
	if err := trackPrinter.Print(outputFile); err != nil {
		log.Fatalf("Error writing to file %s: %s", outputFile.Name(), err)
	}
	fmt.Printf("\nOutput written to %s\n", outputFileName)
}
