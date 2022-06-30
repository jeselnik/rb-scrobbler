package logFile

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	AUDIOSCROBBLER_HEADER = "#AUDIOSCROBBLER/"
	NEWLINE               = "\n"
)

var ErrInvalidLog = errors.New("invalid .scrobbler.log")

/* Take a path to a file and return a representation of that file in a
string slice where every line is a value */
func ImportLog(path *string) ([]string, error) {
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

	logAsLines := strings.Split(string(logInBytes), NEWLINE)
	/* Ensure that the file is actually an audioscrobbler log */
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
