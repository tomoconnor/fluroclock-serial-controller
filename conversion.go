package main

import (
	"fmt"
	"sort"
	"strings"
)

func getLetterLut() map[string]int {
	return map[string]int{
		"A": 0x77,
		"a": 0x7D,
		"b": 0x1F,
		"C": 0x4E,
		"c": 0x0D,
		"d": 0x3D,
		"E": 0x4F,
		"F": 0x47,
		"G": 0x5E,
		"H": 0x37,
		"h": 0x17,
		"I": 0x06,
		"J": 0x3C,
		"L": 0x0E,
		"n": 0x15,
		"O": 0x7E,
		"o": 0x1D,
		"P": 0x67,
		"q": 0x73,
		"r": 0x05,
		"S": 0x5B,
		"t": 0x0F,
		"U": 0x3E,
		"u": 0x1C,
		"y": 0x3B,
	}
}
func GetValidLetters() []string {
	lut := getLetterLut()
	var letters []string
	for k := range lut {
		letters = append(letters, k)
	}
	sort.Strings(letters)
	return letters
}

func GetLetter(letter string) (int, error) {
	lut := getLetterLut()
	if lut[letter] != 0 {
		return lut[letter], nil
	} else if lut[strings.ToUpper(letter)] != 0 {
		return lut[strings.ToUpper(letter)], nil
	} else if lut[strings.ToLower(letter)] != 0 {
		return lut[strings.ToLower(letter)], nil
	}
	return 0, fmt.Errorf("invalid letter: %s", letter)
}

func ConvertIntToBitstring(i int) string {
	return fmt.Sprintf("%07b", i)
}

func ConvertBitstringToDataStruct(bitstring string) DirectData {
	var d DirectData
	d.A = bitstring[0:1]
	d.B = bitstring[1:2]
	d.C = bitstring[2:3]
	d.D = bitstring[3:4]
	d.E = bitstring[4:5]
	d.F = bitstring[5:6]
	d.G = bitstring[6:7]
	return d
}
