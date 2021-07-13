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

/* Take a path to a file and return a representation of that file in a
string slice where every line is a value */
func importLog(path *string) ([]string, error) {
	logFile, err := os.Open(*path)
	if err != nil {
		var emptySlice []string
		return emptySlice, err
	}

	/* It's important to not log.Fatal() or os.Exit() within this function
	as the following statement will not execute and "let go" of the given
	file */
	defer logFile.Close()

	logInBytes, err := ioutil.ReadAll(logFile)
	if err != nil {
		var emptySlice []string
		return emptySlice, err
	}

	logAsLines := strings.Split(string(logInBytes), "\n")
	/* Ensure that the file is actually an audioscrobbler log */
	if !strings.Contains(logAsLines[0], AUDIOSCROBBLER_HEADER) {
		return logAsLines, errors.New("invalid .scrobbler.log")
	} else {
		return logAsLines, nil
	}
}

/* Take a string, split it, convert time if needed and return a track */
func logLineToTrack(line, offset string) Track {
	splitLine := strings.Split(line, SEPARATOR)
	var timestamp string

	/* Time conversion - the API wants it in UTC timezone */
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

/* Convert back/to UTC from localtime */
func convertTimeStamp(timestamp, offset string) string {
	/* Log stores it in unix epoch format. Convert to an int
	so it can be manipulated with the time package */
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	/* Convert the epoch into a Time type for conversion.
	The 0 stands for 0 milliseconds since epoch which isn't
	needed */
	trackTime := time.Unix(timestampInt, 0)

	/* Take the offset flag and convert it into the duration
	which will be added/subtracted */
	newOffset, err := time.ParseDuration(offset)
	if err != nil {
		log.Fatal(err)
	}

	/* The duration is negative so that entries behind UTC are
	brought forward while entries ahead are brought back */
	convertedTime := trackTime.Add(-newOffset)
	return strconv.FormatInt(convertedTime.Unix(), 10)
}

func deleteLogFile(path *string) {
	deletionError := os.Remove(*path)
	if deletionError != nil {
		fmt.Printf("Error Deleting %q!\n%v\n", *path, deletionError)
		os.Exit(1)
	} else {
		fmt.Printf("%q Deleted!\n", *path)
	}
}
