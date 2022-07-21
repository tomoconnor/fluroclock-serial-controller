package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

// func SetPanelBCD()
func GetPortByID(id string, mapping []PortMap) (PortMap, error) {
	for _, m := range mapping {
		if m.PanelID == id {
			return m, nil
		}
	}
	return PortMap{}, fmt.Errorf("the PanelID %s not found", id)
}
func ClockModeUpdateTime(pm *[]PortMap) error {
	dt := time.Now()
	timestamp := dt.Format("15:04")
	log.Println("Time:", timestamp)
	h1 := timestamp[0:1]
	h2 := timestamp[1:2]
	m1 := timestamp[3:4]
	m2 := timestamp[4:5]

	// dataArray := []string{h1, h2, m1, m2}
	log.Println("Data:", h1, h2, m1, m2)
	porth1, err := GetPortByID("4", *pm)
	if err != nil {
		return err
	}
	porth2, err := GetPortByID("3", *pm)
	if err != nil {
		return err
	}
	portm1, err := GetPortByID("2", *pm)
	if err != nil {
		return err
	}
	portm2, err := GetPortByID("1", *pm)
	if err != nil {
		return err
	}

	seth1e := SendSerialString(fmt.Sprintf("+P4MIV%s-", h1), porth1.Port)
	if seth1e != nil {
		log.Printf("error updating time: %s\n", seth1e)
		return seth1e
	}
	f1 := BCDData{
		PanelID: "4",
		Value:   h1,
	}
	porth1.State.Mode = "clock"
	porth1.State.AlphaData = nil
	porth1.State.BCDData = &f1
	porth1.State.DirectData = nil

	seth2e := SendSerialString(fmt.Sprintf("+P3MIV%s-", h2), porth2.Port)
	if seth2e != nil {
		log.Printf("error updating time: %s\n", seth2e)
		return seth2e
	}
	f2 := BCDData{
		PanelID: "3",
		Value:   h2,
	}
	porth2.State.Mode = "clock"
	porth2.State.AlphaData = nil
	porth2.State.BCDData = &f2
	porth2.State.DirectData = nil

	setm1e := SendSerialString(fmt.Sprintf("+P2MIV%s-", m1), portm1.Port)
	if setm1e != nil {
		log.Printf("error updating time: %s\n", setm1e)
		return setm1e
	}

	f3 := BCDData{
		PanelID: "2",
		Value:   m1,
	}

	portm1.State.Mode = "clock"
	portm1.State.AlphaData = nil
	portm1.State.BCDData = &f3
	portm1.State.DirectData = nil

	setm2e := SendSerialString(fmt.Sprintf("+P1MIV%s-", m2), portm2.Port)
	if setm2e != nil {
		log.Printf("error updating time: %s\n", setm2e)
		return setm2e
	}
	f4 := BCDData{
		PanelID: "1",
		Value:   m2,
	}
	portm2.State.Mode = "clock"
	portm2.State.AlphaData = nil
	portm2.State.BCDData = &f4
	portm2.State.DirectData = nil

	// for i, m := range *pm {
	// 	value := dataArray[i]
	// 	// log.Printf("+P%sMIV%s-", m.PanelID, value)
	// 	sse := SendSerialString(fmt.Sprintf("+P%sMIV%s-", m.PanelID, value), m.Port)
	// 	if sse != nil {
	// 		log.Printf("error updating time: %s\n", sse)
	// 		return sse
	// 	}
	// 	fakeBCDData := BCDData{
	// 		PanelID: m.PanelID,
	// 		Value:   value,
	// 	}

	// 	m.State.Mode = "clock"
	// 	m.State.AlphaData = nil
	// 	m.State.BCDData = &fakeBCDData
	// 	m.State.DirectData = nil

	// }
	return nil
}

func GetPortMapping(panelID string, mapping []PortMap) (PortMap, error) {
	for _, m := range mapping {
		if m.PanelID == panelID {
			return m, nil
		}
	}
	return PortMap{}, fmt.Errorf("the PanelID %s not found", panelID)
}

// func (pm []PortMap) GetPortByID(PanelID string) *serial.Port {
// 	for _, m := range *pm {
// 		if m.PanelID == PanelID {
// 			return m.Port
// 		}
// 	}
// 	return nil
// }

// ports := make([]PortMap, 3)

func GetArduinoSerialPorts() []string {
	ValidPorts := make([]string, 0)

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found")
	} else {
		log.Printf("Found %d serial ports\n", len(ports))
		for _, port := range ports {
			log.Printf("Found port: %s\n", port)
			if strings.Contains(port, "ttyUSB") {
				ValidPorts = append(ValidPorts, port)
			}
		}
	}

	return ValidPorts

}

func SendSerialString(codes string, sport serial.Port) error {
	serialPort := sport
	codes_as_bytes := []byte(codes)
	log.Println("Sending:", codes)

	size, err := serialPort.Write(codes_as_bytes)

	if err != nil {
		log.Printf("Error writing to serial port: %s\n", err)
		return err
	}
	log.Printf("Wrote %d bytes to serial port\n", size)

	// log.Printf("%b\n", n)
	return nil
}

func getResponse(port serial.Port) string {
	serialPort := port
	buff := make([]byte, 100)
	response := make([]byte, 1)
	for {
		// Reads up to 100 bytes
		n, err := serialPort.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			log.Println("\nEOF")
			break
		}

		// log.Printf("%s", string(buff[:n]))
		response = append(response, buff[:n]...)
		// If we receive a newline stop reading
		if strings.Contains(string(buff[:n]), "-") {
			break
		}
	}
	return strings.TrimSpace(string(response))
}

func NewPortMapping(serialPort string) PortMap {
	m := serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	p, err := serial.Open(serialPort, &m)
	if err != nil {
		log.Fatal(err)
	}
	p.ResetOutputBuffer()
	p.ResetInputBuffer()
	p.Write([]byte("+?-"))

	time.Sleep(time.Second * 1)
	response := getResponse(p)
	log.Printf("Response: %s\n", response)
	// p.Close()

	pid := string(response[2])

	log.Printf("Panel ID: %s\n", pid)
	pm := PortMap{
		Port:     p,
		PanelID:  pid,
		PortPath: serialPort,
		State:    new(PanelState),
	}
	return pm

}

func main() {

	serialports := GetArduinoSerialPorts()
	nports := len(serialports)
	log.Println("Found USB serial ports x", nports)

	mapping := make([]PortMap, 0)

	for i, port := range serialports {
		log.Println("Port:", port)
		log.Println("idx:", i)
		mapping = append(mapping, NewPortMapping(port))
	}
	log.Println("Mapping:")
	for _, m := range mapping {
		log.Printf("%s -> %s\t State: %q\n", m.PanelID, m.PortPath, m.State.Repr())
	}

	// disable all displays.

	for _, m := range mapping {

		SendSerialString(fmt.Sprintf("+P%sMXA0B0C0D0E0F0G0-", m.PanelID), m.Port)
		time.Sleep(time.Second * 1)
	}

	StartServer(&mapping)

}
