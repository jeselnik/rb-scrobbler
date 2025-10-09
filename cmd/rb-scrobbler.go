package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jeselnik/rb-scrobbler/internal/logFile"
	"github.com/jeselnik/rb-scrobbler/internal/track"
	"github.com/twoscott/gobble-fm/session"
)

const SECONDS_IN_HOUR = 3600

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

	fm := session.NewClient(API_KEY, API_SECRET)

	/* First time Authentication */
	if *auth {
		/* https://www.last.fm/api/desktopauth */

		/* Create folder to store session key */
		configDir, err := getConfigDir()
		if err != nil {
			log.Fatal(err)
		}

		err = os.Mkdir(configDir, 0700)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}

		/* "Step 2" */
		token, err := fm.Auth.Token()
		if err != nil {
			log.Fatal(err)
		}

		/* Step 3 */
		authURL := fm.AuthTokenURL(token)
		fmt.Printf("Go to %q, allow access and press ENTER\n", authURL)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

		/* "Step 4" */
		loginErr := fm.TokenLogin(token)
		if loginErr != nil {
			log.Fatal(loginErr)
		}

		sessionKey := fm.SessionKey
		/* Save session key in $CONFIG/rb-scrobbler */
		keyPath, err := getKeyFilePath()
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(keyPath, []byte(sessionKey), 0600)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("User successfully authenticated.")
	}

	if *logPath == "" {
		fmt.Println("File (-f) cannot be empty!")
		os.Exit(1)
	}

	/* When given a file, start executing here */
	offsetSec := *offset * SECONDS_IN_HOUR

	tracks, err := logFile.ImportLog(logPath, offsetSec, colours)
	if err != nil {
		log.Fatal(err)
	}

	/* Login here, after tracks have been parsed and are ready to send */
	sessionKey, err := getSavedKey()
	if err != nil {
		log.Fatal(err)
	}
	fm.SetSessionKey(sessionKey)

	success, fail := track.Scrobble(fm, tracks, colours)
	fmt.Printf("\nFinished: %d tracks scrobbled, %d failed, %d total\n",
		success, fail, len(tracks))

	/* Handling of file (manual/non interactive delete/keep) */
	os.Exit(logFile.HandleFile(nonInteractive, logPath, fail))

}
