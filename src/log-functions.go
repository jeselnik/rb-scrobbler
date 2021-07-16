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
	AUDIOSCROBBLER_HEADER = "#AUDIOSCROBBLER/"
	SEPARATOR             = "\t"
	LISTENED              = "L"
	ARTIST_INDEX          = 0
	ALBUM_INDEX           = 1
	TITLE_INDEX           = 2
	RATING_INDEX          = 5
	TIMESTAMP_INDEX       = 6
)

/* Take a path to a file and return a representation of that file in a
string slice where every line is a value */
func importLog(path *string) ([]string, error) {
	logFile, err := os.Open(*path)
	if err != nil {
		return []string{}, err
	}

	/* It's important to not log.Fatal() or os.Exit() within this function
	as the following statement will not execute and "let go" of the given
	file */
	defer logFile.Close()

	logInBytes, err := ioutil.ReadAll(logFile)
	if err != nil {
		return []string{}, err
	}

	logAsLines := strings.Split(string(logInBytes), "\n")
	/* Ensure that the file is actually an audioscrobbler log */
	if !strings.Contains(logAsLines[0], AUDIOSCROBBLER_HEADER) {
		return []string{}, errors.New("invalid .scrobbler.log")
	} else {
		return logAsLines, nil
	}
}

/* Take a string, split it, convert time if needed and return a track */
func logLineToTrack(line, offset string) (Track, error) {
	splitLine := strings.Split(line, SEPARATOR)

	/* Check the "RATING" index instead of looking for "\tL\t" in a line,
	just in case a track or album is named "L". If anything like this exists
	and was skipped the old method would false positive it as listened
	and then it'd be submitted */

	if splitLine[RATING_INDEX] == LISTENED {
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

		return track, nil

	} else {
		return Track{}, errors.New("Track was skipped")
	}
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
