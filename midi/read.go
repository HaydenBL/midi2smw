package midi

import (
	"encoding/binary"
	"fmt"
	"os"
)

func readByte(file *os.File) uint8 {
	var i uint8
	err := binary.Read(file, binary.BigEndian, &i)
	if err != nil {
		fmt.Printf("Error reading: %d\n", err)
	}
	return i
}

func readString(file *os.File, length uint32) string {
	b := make([]byte, length)
	n, _ := file.Read(b)
	str := string(b[:n])
	return str
}

// Values are chained together by using the most significant bit as a flag, indicating
// whether or not another byte should be read. The lower 7 bits contain the actual data
// and we'll just shift them all into a 32 bit integer while the flag is set
func readValue(file *os.File) uint32 {
	var finalValue uint32 = 0

	aByte := readByte(file)
	finalValue = uint32(aByte)

	// If MSB is set, we need to read more bytes in
	if (finalValue & 0x80) != 0 {
		finalValue &= 0x7F                             // Extract bottom 7 bits of read byte
		for ok := true; ok; ok = (aByte & 0x80) != 0 { // Loop while MSB is 1
			// Read next byte
			aByte = readByte(file)

			// Shift 7 bits in, apply value from last byte read into their position
			finalValue = (finalValue << 7) | (uint32(aByte) & 0x7F)
		}
	}

	return finalValue
}
