package logFile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const AUDIOSCROBBLER_HEADER = "#AUDIOSCROBBLER/"

var ErrInvalidLog = errors.New("invalid .scrobbler.log")

/* Take a path to a file and return a representation of that file in a
string slice where every line is a value */
func ImportLog(path *string) ([]string, error) {
	var logAsLines []string

	logFile, err := os.Open(*path)
	if err != nil {
		return logAsLines, err
	}

	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		logAsLines = append(logAsLines, scanner.Text())
	}

	if !strings.Contains(logAsLines[0], AUDIOSCROBBLER_HEADER) {
		return []string{}, ErrInvalidLog
	} else {
		return logAsLines, nil
	}
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
