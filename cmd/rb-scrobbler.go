package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/Jeselnik/rb-scrobbler/internal/logFile"
	"github.com/Jeselnik/rb-scrobbler/internal/track"
	"github.com/shkh/lastfm-go/lastfm"
)

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	offset := flag.Float64("o", 0,
		"Time difference from UTC (format +10 or -10.5)")
	nonInteractive := flag.String("n", "",
		`Non Interactive Mode:
Automatically ("keep", "delete" or "delete-on-success") at end of program`)
	auth := flag.Bool("auth", false, "First Time Authentication")
	colours := flag.Bool("nc", true,
		"No Terminal Colours (Default behaviour on Windows)")
	flag.Parse()

	/* Windows CMD doesn't support esc code colours, default
	it to false */
	if runtime.GOOS == "windows" {
		*colours = false
	}

	api := lastfm.New(API_KEY, API_SECRET)

	/* First time Authentication */
	if *auth {
		/* https://www.last.fm/api/desktopauth */

		/* Create folder to store session key */
		err := os.Mkdir(getConfigDir(), os.ModePerm)
		if err != nil && !errors.Is(err, os.ErrExist) {
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
		err = os.WriteFile(getKeyFilePath(), []byte(sessionKey), os.ModePerm)
		if errors.Is(err, os.ErrExist) {
			err = os.Remove(getKeyFilePath())
			if err != nil {
				log.Fatal(err)
			}
			err = os.WriteFile(getKeyFilePath(), []byte(sessionKey),
				os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
		}

		var userName string
		getInfo, err := api.User.GetInfo(lastfm.P{})
		if err != nil {
			userName = "User"
		} else {
			userName = getInfo.Name
		}

		fmt.Printf("%s authenticated\n", userName)
	}

	/* When given a file, start executing here */
	if *logPath != "" {

		tracks, err := logFile.ImportLog(logPath, offset)
		if err != nil {
			log.Fatal(err)
		}

		/* Login here, after tracks have been parsed and are ready to send */
		sessionKey, err := getSavedKey()
		if err != nil {
			log.Fatal(err)
		}
		api.SetSession(sessionKey)

		success, fail := track.Scrobble(api, tracks, colours)
		fmt.Printf("\nFinished: %d tracks scrobbled, %d failed, %d total\n",
			success, fail, len(tracks))

		/* Handling of file (manual/non interactive delete/keep) */
		os.Exit(logFile.HandleFile(nonInteractive, logPath, fail))

	} else if !(*auth) {
		fmt.Println("File (-f) cannot be empty!")
		os.Exit(1)
	}

}
