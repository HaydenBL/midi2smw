package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	parseFile("dean_town.mid")
}

func parseFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	defer file.Close()

	numTracks := readHeader(file)
	fmt.Println(numTracks)

	//ntc := readString(file, 4)
	//fmt.Println(ntc)

	for i := 0; i < 500; i++ {
		bah := readValue(file)
		fmt.Println(bah)
	}

}

func readString(file *os.File, length uint32) string {
	b := make([]byte, length)
	n, _ := file.Read(b)
	str := string(b[:n])
	return str
}

func readHeader(file *os.File) (numTrackChunks uint16) {
	var val32 uint32 = 0
	var val16 uint16 = 0

	// First 4 bytes, file ID (always MThd)
	binary.Read(file, binary.BigEndian, &val32)
	// Next 4 bytes, length of header
	binary.Read(file, binary.BigEndian, &val32)
	// Next 2 bytes, format details
	binary.Read(file, binary.BigEndian, &val16)
	// Next 2 bytes, number of tracks
	binary.Read(file, binary.BigEndian, &val16)
	numTrackChunks = val16
	// Next 2 bytes, division
	binary.Read(file, binary.BigEndian, &val16)

	return numTrackChunks
}

// Values are chained together by using the most significant bit as a flag, indicating
// whether or not another byte should be read. The lower 7 bits contain the actual data
// and we'll just shift them all into a 32 bit integer while the flag is set
func readValue(file *os.File) uint32 {
	var finalValue uint32 = 0
	var aByte uint8 = 0

	binary.Read(file, binary.BigEndian, &aByte)
	finalValue = uint32(aByte)

	// If MSB is set, we need to read more bytes in
	if (finalValue & 0x80) != 0 {
		finalValue &= 0x7F                             // Extract bottom 7 bits of read byte
		for ok := true; ok; ok = (aByte & 0x80) != 0 { // Loop while MSB is 1
			// Read next byte
			binary.Read(file, binary.BigEndian, &aByte)

			// Shift 7 bits in, apply value from last byte read into their position
			finalValue = (finalValue << 7) | (uint32(aByte) & 0x7F)
		}
	}

	return finalValue
}
