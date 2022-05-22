package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/gommon/log"
	"go.bug.st/serial"
)

// Global State
var ClockModeEnabled = false

type SystemState struct {
	ClockMode bool `json:"clock_mode"` //true = clock mode enabled, false = clock mode disabled
}

type PortMap struct {
	Port     serial.Port
	PortPath string
	PanelID  string
	State    *PanelState
}

type BCDData struct {
	PanelID string `json:"panel_id"`
	Value   string `json:"value"`
}

type AlphaData struct {
	PanelID string `json:"panel_id"`
	Alpha   string `json:"alpha"`
}

type DirectData struct {
	PanelID string `json:"panel_id"`
	A       string `json:"a"`
	B       string `json:"b"`
	C       string `json:"c"`
	D       string `json:"d"`
	E       string `json:"e"`
	F       string `json:"f"`
	G       string `json:"g"`
}

type PanelState struct {
	Mode       string      `json:"mode"`
	BCDData    *BCDData    `json:"bcd_data"`
	AlphaData  *AlphaData  `json:"alpha_data"`
	DirectData *DirectData `json:"direct_data"`
}

func (p *PanelState) Repr() string {
	result := fmt.Sprintf("Mode: %s", p.Mode)
	if p.Mode == "bcd" {
		result += fmt.Sprintln("PanelID:", p.BCDData.PanelID)
		result += fmt.Sprintln("Value:", p.BCDData.Value)
	} else if p.Mode == "alpha" {
		result += fmt.Sprintln("PanelID:", p.AlphaData.PanelID)
		result += fmt.Sprintln("Alpha:", p.AlphaData.Alpha)
	} else if p.Mode == "direct" {
		result += fmt.Sprintln("PanelID:", p.DirectData.PanelID)
		result += fmt.Sprintln("A:", p.DirectData.A)
		result += fmt.Sprintln("B:", p.DirectData.B)
		result += fmt.Sprintln("C:", p.DirectData.C)
		result += fmt.Sprintln("D:", p.DirectData.D)
		result += fmt.Sprintln("E:", p.DirectData.E)
		result += fmt.Sprintln("F:", p.DirectData.F)
		result += fmt.Sprintln("G:", p.DirectData.G)
	}
	return result
}

func LoadTZData(filename string) (*[]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var slice []string
	err = json.Unmarshal(data, &slice)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &slice, nil

}
