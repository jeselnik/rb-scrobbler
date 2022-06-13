package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shkh/lastfm-go/lastfm"
)

const FIRST_TRACK_LINE_INDEX = 3

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	offset := flag.String("o", "0h",
		"Time difference from UTC (format +10h or -10.5h")
	nonInteractive := flag.String("n", "",
		"Non Interactive Mode: Automatically (\"keep\", \"delete\""+
			"or \"delete-on-success\") at end of program")
	auth := flag.Bool("auth", false, "First Time Authentication")
	colours := flag.Bool("nc", true, "No Terminal Colours")
	flag.Parse()

	api := lastfm.New(API_KEY, API_SECRET)

	/* First time Authentication */
	if *auth {
		/* https://www.last.fm/api/desktopauth */

		/* Create folder to store session key */
		err := os.Mkdir(getConfigDir(), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		/* "Step 2" */
		token, err := api.GetToken()
		if err != nil {
			log.Fatal(err)
		}

		/* Step 3 */
		authURL := api.GetAuthTokenUrl(token)
		fmt.Printf("Go to %q, allow access and press ENTER\n", authURL)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

		/* "Step 4" */
		loginErr := api.LoginWithToken(token)
		if loginErr != nil {
			log.Fatal(loginErr)
		}

		sessionKey := api.GetSessionKey()
		/* Save session key in $CONFIG/rb-scrobbler */
		os.WriteFile(getKeyFilePath(), []byte(sessionKey), os.ModePerm)

	}

	/* When given a file, start executing here */
	if *logPath != "" {
		scrobblerLog, err := importLog(logPath)
		if err != nil {
			log.Fatal(err)
		}

		var tracks Tracks
		/* length -1 since you go out of bounds otherwise.
		Only iterate from where tracks actually show up */
		for i := FIRST_TRACK_LINE_INDEX; i < len(scrobblerLog)-1; i++ {
			newTrack, listened := logLineToTrack(scrobblerLog[i], *offset)
			if listened {
				tracks = append(tracks, newTrack)
			}
		}

		/* Login here, after tracks have been parsed and are ready to send */
		sessionKey, err := getSavedKey()
		if err != nil {
			log.Fatal(err)
		}
		api.SetSession(sessionKey)

		success, fail := tracks.scrobble(api, colours)
		fmt.Printf("\nFinished: %d tracks scrobbled, %d failed, %d total\n",
			success, fail, len(tracks))

		/* Handling of file (manual/non interactive delete/keep) */
		os.Exit(logFileHandling(nonInteractive, logPath, fail))

	} else {
		fmt.Println("File (-f) cannot be empty!")
		os.Exit(1)
	}

}
