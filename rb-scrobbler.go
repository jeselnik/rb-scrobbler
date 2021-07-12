package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Track struct {
	artist    string
	album     string
	title     string
	timestamp string
}

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	offset := flag.String("o", "0h", "Offset from UTC")
	nonInteractive := flag.String("n", "", "Non Interactive Mode (Delete or Keep log with no user input.)")
	auth := flag.Bool("auth", false, "First time authentication")
	flag.Parse()

	if *auth {
		fmt.Println("Auth stream TODO")
	}

	if *logPath != "" {
		scrobblerLog, err := importLog(logPath)
		if err != nil {
			log.Fatal(err)
		}

		var tracks []Track
		for i := FIRST_TRACK_LINE_INDEX; i < len(scrobblerLog)-1; i++ {
			if strings.Contains(scrobblerLog[i], LISTENED) {
				tracks = append(tracks, logLineToTrack(scrobblerLog[i], *offset))
			}
		}

		switch *nonInteractive {
		case "keep":
			fmt.Printf("%q kept\n", *logPath)

		case "delete":
			deleteLogFile(logPath)

		default:
			reader := bufio.NewReader(os.Stdin)
			var input string
			fmt.Printf("Delete %q? [y/n] ", *logPath)
			input, err := reader.ReadString('\n')
			fmt.Print("\n")
			if err != nil {
				fmt.Printf("Error reading input! File %q not deleted.\n%v\n", *logPath, err)
			} else if strings.ContainsAny(input, "y") || strings.ContainsAny(input, "Y") {
				deleteLogFile(logPath)
			} else {
				fmt.Printf("%q kept\n", *logPath)
			}
		}

	} else {
		fmt.Println("File (-f) cannot be empty!")
		os.Exit(1)
	}

}
