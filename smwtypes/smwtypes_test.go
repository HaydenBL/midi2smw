package smwtypes

import "testing"

func TestGetKeyFromKeyValue(t *testing.T) {
	var k uint8
	// too low for SMW note
	for k = 0; k < 19; k++ {
		key := getKeyFromKeyValue(k)
		if key != "r" {
			t.Fatalf("Should error when key value too low")
		}
	}
	// valid SMW note range
	for k = 19; k < 89; k++ {
		key := getKeyFromKeyValue(k)
		if key == "r" {
			t.Fatalf("Shouldn't error when key value is within range")
		}
	}
	// too high for SMW note
	for k = 89; k < 128; k++ {
		key := getKeyFromKeyValue(k)
		if key != "r" {
			t.Fatalf("Should error when key value too high")
		}
	}
}
