package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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

func logLineToTrack(line string) Track {
	splitLine := strings.Split(line, SEPARATOR)
	timestamp, err := strconv.ParseUint(splitLine[TIMESTAMP_INDEX], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	track := Track{
		artist:    splitLine[ARTIST_INDEX],
		album:     splitLine[ALBUM_INDEX],
		title:     splitLine[TITLE_INDEX],
		timestamp: timestamp,
	}

	return track
}
