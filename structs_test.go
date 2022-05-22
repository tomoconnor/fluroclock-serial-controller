package main

import "testing"

func TestLoadTZData(t *testing.T) {
	td, err := LoadTZData("tz.json")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(*td) != 538 {
		t.Errorf("Expected 538 timezones, got %d", len(*td))
	}

	if (*td)[0] != "Africa/Abidjan" {
		t.Errorf("Expected Africa/Abidjan, got %s", (*td)[0])
	}

	if (*td)[537] != "Zulu" {
		t.Errorf("Expected Zulu, got %s", (*td)[537])
	}

}
