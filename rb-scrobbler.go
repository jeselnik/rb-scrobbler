package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	AUDIOSCROBBLER_HEADER  = "#AUDIOSCROBBLER/"
	LISTENED               = "\tL\t"
	SEPARATOR              = "\t"
	FIRST_TRACK_LINE_INDEX = 3
	ARTIST_INDEX           = 0
	ALBUM_INDEX            = 1
	TITLE_INDEX            = 2
	TIMESTAMP_INDEX        = 6
)

type Track struct {
	artist    string
	album     string
	title     string
	timestamp uint64
}

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

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	flag.Parse()

	scrobblerLog, err := importLog(logPath)
	if err != nil {
		log.Fatal(err)
	}

	var tracks []Track
	for i := FIRST_TRACK_LINE_INDEX; i < len(scrobblerLog)-1; i++ {
		tracks = append(tracks, logLineToTrack(scrobblerLog[i]))
	}
}
