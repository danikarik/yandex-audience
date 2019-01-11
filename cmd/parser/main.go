package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	portablePath   = "./data/mcTrack-12-5-18-portable.csv"
	stationaryPath = "./data/mcTrack-12-5-18-stattionary.csv"
)

var (
	portableOut   = "./out/portable.txt"
	stationaryOut = "./out/stattionary.txt"
)

func main() {
	err := loadPortable()
	if err != nil {
		log.Fatalf("portable: %v", err)
	}
	err = loadStationary()
	if err != nil {
		log.Fatalf("stationary: %v", err)
	}
	log.Println("done.")
}

func loadPortable() error {
	buf, err := ioutil.ReadFile(portablePath)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(buf)
	reader := csv.NewReader(buffer)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	file, err := os.Create(portableOut)
	if err != nil {
		return err
	}
	defer file.Close()
	for i, line := range lines {
		if i == 0 {
			continue
		}
		mac := strings.ToUpper(strings.Replace(line[1], ":", "", -1))
		file.WriteString(mac)
		file.WriteString("\n")
	}
	return nil
}

func loadStationary() error {
	buf, err := ioutil.ReadFile(stationaryPath)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(buf)
	reader := csv.NewReader(buffer)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	file, err := os.Create(stationaryOut)
	if err != nil {
		return err
	}
	defer file.Close()
	for i, line := range lines {
		if i == 0 {
			continue
		}
		mac := strings.ToUpper(strings.Replace(line[2], ":", "", -1))
		file.WriteString(mac)
		file.WriteString("\n")
	}
	return nil
}
