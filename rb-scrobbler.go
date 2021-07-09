package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	AUDIOSCROBBLER_HEADER = "#AUDIOSCROBBLER/"
	LISTENED              = "\tL\t"
	SEPARATOR             = "\t"
	ARTIST_INDEX          = 0
	ALBUM_INDEX           = 1
	TITLE_INDEX           = 2
	TIMESTAMP_INDEX       = 6
)

func importLog(path *string) ([]string, error) {
	logFile, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logInBytes, err := ioutil.ReadAll(logFile)
	if err != nil {
		log.Fatal(err)
	}

	logAsLines := strings.Split(string(logInBytes), "\n")
	if !strings.Contains(logAsLines[0], AUDIOSCROBBLER_HEADER) {
		return logAsLines, errors.New("Not a valid .scrobbler.log!")
	} else {
		return logAsLines, nil
	}
}

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	flag.Parse()

	scrobblerLog, err := importLog(logPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(scrobblerLog[0])
}
