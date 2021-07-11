package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func importLog(path *string) ([]string, error) {
	logFile, err := os.Open(*path)
	if err != nil {
		var emptySlice []string
		return emptySlice, err
	}
	defer logFile.Close()

	logInBytes, err := ioutil.ReadAll(logFile)
	if err != nil {
		var emptySlice []string
		return emptySlice, err
	}

	logAsLines := strings.Split(string(logInBytes), "\n")
	if !strings.Contains(logAsLines[0], AUDIOSCROBBLER_HEADER) {
		return logAsLines, errors.New("Not a valid .scrobbler.log!")
	} else {
		return logAsLines, nil
	}
}

func logLineToTrack(line, offset string) Track {
	splitLine := strings.Split(line, SEPARATOR)
	var timestamp string

	if offset != "0h" {
		timestamp = convertTimeStamp(splitLine[TIMESTAMP_INDEX], offset)
	} else {
		timestamp = splitLine[TIMESTAMP_INDEX]
	}

	track := Track{
		artist:    splitLine[ARTIST_INDEX],
		album:     splitLine[ALBUM_INDEX],
		title:     splitLine[TITLE_INDEX],
		timestamp: timestamp,
	}

	return track
}

func convertTimeStamp(timestamp, offset string) string {
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	/* 0 = zero milliseconds */
	trackTime := time.Unix(timestampInt, 0)
	newOffset, err := time.ParseDuration(offset)
	if err != nil {
		log.Fatal(err)
	}

	convertedTime := trackTime.Add(-newOffset)
	return strconv.FormatInt(convertedTime.Unix(), 10)
}

func deleteLogFile(path *string) {
	deletionError := os.Remove(*path)
	if deletionError != nil {
		fmt.Printf("Error Deleting %q!\n%v\n", *path, deletionError)
		os.Exit(1)
	} else {
		fmt.Printf("%q Deleted!", *path)
	}
}
