package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Touch(fileName string) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(fileName, currentTime, currentTime)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func InClockModeFromRestart() bool {
	_, err := os.Stat("/etc/fclock-clock-enable")
	return err == nil
}

func ResetDisplay() {
	a, _ := json.Marshal(BCDData{
		PanelID: "1",
		Value:   "0",
	})
	b, _ := json.Marshal(BCDData{
		PanelID: "2",
		Value:   "0",
	})
	c, _ := json.Marshal(BCDData{
		PanelID: "3",
		Value:   "0",
	})
	d, _ := json.Marshal(BCDData{
		PanelID: "4",
		Value:   "0",
	})

	aresp, err := http.Post("http://localhost:9000/panel/numeric", "application/json", bytes.NewBuffer(a))
	if err != nil {
		log.Fatal(err)
	}
	defer aresp.Body.Close()
	if aresp.StatusCode != 200 {
		log.Fatal("error sending panel 1")
	}
	bresp, err := http.Post("http://localhost:9000/panel/numeric", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer bresp.Body.Close()
	if bresp.StatusCode != 200 {
		log.Fatal("error sending panel 2")
	}
	cresp, err := http.Post("http://localhost:9000/panel/numeric", "application/json", bytes.NewBuffer(c))
	if err != nil {
		log.Fatal(err)
	}
	defer cresp.Body.Close()
	if cresp.StatusCode != 200 {
		log.Fatal("error sending panel 3")
	}
	dresp, err := http.Post("http://localhost:9000/panel/numeric", "application/json", bytes.NewBuffer(d))
	if err != nil {
		log.Fatal(err)
	}
	defer dresp.Body.Close()
	if dresp.StatusCode != 200 {
		log.Fatal("error sending panel 4")
	}
}
