package main

import (
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
			deletionError := os.Remove(*logPath)
			if deletionError != nil {
				fmt.Printf("I/O Error Deleting %q!\n%v\n", *logPath, deletionError)
				os.Exit(1)
			} else {
				fmt.Printf("%q deleted!\n", *logPath)
			}

		default:
			fmt.Println("TODO")
		}

	} else {
		fmt.Println("File (-f) cannot be empty!")
		os.Exit(1)
	}

}
