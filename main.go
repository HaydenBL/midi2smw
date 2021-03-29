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
	fileName, splitTracksFlag, samplesFlag := parseFlags()
	begin(fileName, splitTracksFlag, samplesFlag)
}

func parseFlags() (fileName string, splitTracksFlag, samplesFlag bool) {
	splitTracksFlagPtr := flag.Bool("split", false, "Specify tracks to split with note groupings")
	samplesFlagPtr := flag.Bool("samples", false, "Specify samples for notes")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Error: no file name provided")
		os.Exit(1)
	}
	return flag.Args()[0], *splitTracksFlagPtr, *samplesFlagPtr
}

func begin(fileName string, splitTracksFlag, samplesFlag bool) {

	fmt.Printf("========== BEGIN PARSING ==========\n\n")

	midiFile, err := midi.Parse(fileName)
	if err != nil {
		fmt.Printf("Error parsing midi file: %s\n", fileName)
		return
	}

	fmt.Printf("\n\n\n========== BEGIN CONVERTING ==========\n\n")

	tracks := convert.Convert(midiFile, splitTracksFlag, samplesFlag)

	fmt.Printf("\n\n\n========== BEGIN WRITING ==========\n\n")

	outputFile, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("Error creating file")
	}
	defer outputFile.Close()

	trackPrinter := trackoutput.NewPrinter(tracks, midiFile.Bpm)
	if err := trackPrinter.Print(outputFile); err != nil {
		log.Fatalf("Error writing to file %s: %s", outputFile.Name(), err)
	}

	fmt.Printf("\n\n\n========== COMPLETE ==========\n")
}
