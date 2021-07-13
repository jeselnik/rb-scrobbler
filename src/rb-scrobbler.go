package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shkh/lastfm-go/lastfm"
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
		token, _ := api.GetToken()

		authURL := api.GetAuthTokenUrl(token)

		/* "Step 3" */
		fmt.Printf("Go to %q, allow access and press ENTER\n", authURL)
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')

		/* "Step 4" */
		api.LoginWithToken(token)
		sessionKey := api.GetSessionKey()

		/* Save session key in config dir/rb-scrobbler */
		os.WriteFile(getKeyFilePath(), []byte(sessionKey), os.ModePerm)
	}

	/* When given a file, start executing here */

	if *logPath != "" {
		scrobblerLog, err := importLog(logPath)
		if err != nil {
			log.Fatal(err)
		}

		var tracks []Track
		/* length -1 since you go out of bounds otherwise. Only iterate from where trakcs
		actually show up */
		for i := FIRST_TRACK_LINE_INDEX; i < len(scrobblerLog)-1; i++ {
			if strings.Contains(scrobblerLog[i], LISTENED) {
				tracks = append(tracks, logLineToTrack(scrobblerLog[i], *offset))
			}
		}

		/* Login here, after tracks have been parsed and are ready to send */
		api.SetSession(getSavedKey())

		for _, track := range tracks {
			p := lastfm.P{"artist": track.artist, "album": track.album, "track": track.title, "timestamp": track.timestamp}

			_, err := api.Track.Scrobble(p)
			if err != nil {
				fmt.Printf("[FAIL] %s - %s\n", track.artist, track.title)
			} else {
				fmt.Printf("[OK] %s - %s\n", track.artist, track.title)
			}
		}

		/* Handling of file (manual/non interactive delete/keep) */

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