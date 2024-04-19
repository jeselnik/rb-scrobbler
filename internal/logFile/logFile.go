package logFile

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/jeselnik/rb-scrobbler/internal/track"
)

const (
	AUDIOSCROBBLER_HEADER = "#AUDIOSCROBBLER/"
	REGEX_HEADERS         = "^#(AUDIOSCROBBLER/|TZ/|CLIENT/|ARTIST #ALBUM)"
	SEPARATOR             = '\t'
)

var ErrInvalidLog = errors.New("invalid .scrobbler.log")

func ImportLog(path *string, offset *float64) ([]track.Track, error) {

	var (
		logErr error = nil
		tracks []track.Track
	)

	f, err := os.Open(*path)
	if err != nil {
		return tracks, logErr
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	r := csv.NewReader(reader)
	r.Comma = SEPARATOR
	r.ReuseRecord = true

	headers, err := regexp.Compile(REGEX_HEADERS)
	if err != nil {
		return tracks, logErr
	}

	first := true
	for {
		line, err := r.Read()

		if err == io.EOF {
			break
		}

		if first && !strings.Contains(line[0], AUDIOSCROBBLER_HEADER) {
			logErr = ErrInvalidLog
			break
		}
		first = false

		if headers.MatchString(line[0]) {
			continue
		}

		if line[track.RATING_INDEX] != track.LISTENED {
			continue
		}

		trackObj, err := track.StringToTrack(line, *offset)
		if err != nil {
			continue
		}

		tracks = append(tracks, trackObj)

	}

	return tracks, logErr

}

func deleteLogFile(path *string) (exitCode int) {
	deletionError := os.Remove(*path)
	if deletionError != nil {
		fmt.Printf("Error Deleting %q!\n%v\n", *path, deletionError)
		exitCode = 1
	} else {
		fmt.Printf("%q Deleted!\n", *path)
	}
	return exitCode
}

func HandleFile(nonInteractive, logPath *string, fail uint) int {
	exitCode := 0

	switch *nonInteractive {
	case "keep":
		fmt.Printf("%q kept\n", *logPath)

	case "delete":
		exitCode = deleteLogFile(logPath)

	case "delete-on-success":
		if fail == 0 {
			exitCode = deleteLogFile(logPath)
		} else {
			fmt.Printf("Scrobble failures: %q not deleted.\n", *logPath)
			exitCode = 1
		}

	default:
		reader := bufio.NewReader(os.Stdin)
		var input string
		fmt.Printf("Delete %q? [y/n] ", *logPath)
		input, err := reader.ReadString('\n')
		fmt.Print("\n")
		if err != nil {
			fmt.Printf("Error reading input! File %q not deleted.\n%v\n",
				*logPath, err)
			exitCode = 1
		} else if strings.ContainsAny(input, "y") ||
			strings.ContainsAny(input, "Y") {
			deleteLogFile(logPath)
		} else {
			fmt.Printf("%q kept.\n", *logPath)
		}
	}

	return exitCode
}
