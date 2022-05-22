package main

import (
	"testing"
)

func TestGetLetter(t *testing.T) {
	ans, err := GetLetter("a")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if ans != 0x7D {
		t.Errorf("Expected 0x7D, got %d", ans)
	}
	ans, err = GetLetter("A")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if ans != 0x77 {
		t.Errorf("Expected 0x77, got %d", ans)
	}
}

func TestFailToGetLetter(t *testing.T) {
	_, err := GetLetter("z")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

}

func TestConvertToBitString(t *testing.T) {
	// A
	letterA, _ := GetLetter("A")

	ans := ConvertIntToBitstring(letterA)
	if ans != "1110111" {
		t.Errorf("Expected 1110111, got %s", ans)
	}
	// a
	letterA, _ = GetLetter("a")
	ans = ConvertIntToBitstring(letterA)
	if ans != "1111101" {
		t.Errorf("Expected 1111101, got %s", ans)
	}
}

func TestConvertBitstringToDataStruct(t *testing.T) {
	// A
	letterA, _ := GetLetter("A")
	ans := ConvertIntToBitstring(letterA)
	// if ans != "1110111" {
	bs := ConvertBitstringToDataStruct(ans)
	if bs.A != "1" &&
		bs.B != "1" &&
		bs.C != "1" &&
		bs.D != "0" &&
		bs.E != "1" &&
		bs.F != "1" &&
		bs.G != "1" {
		t.Errorf("Expected A, got %s", ans)
	}
}

func TestGetValidLetters(t *testing.T) {
	ans := GetValidLetters()

	if len(ans) != 25 {
		t.Errorf("Expected 25, got %d", len(ans))
	}
	if ans[0] != "A" {
		t.Errorf("Expected A, got %s", ans[0])
	}
	if ans[24] != "y" {
		t.Errorf("Expected y, got %s", ans[23])
	}

}
