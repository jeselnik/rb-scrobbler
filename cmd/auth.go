/* This set of functions deals with getting and saving the
session key in a cross OS manner */

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func getConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dir, "rb-scrobbler")
}

func getKeyFilePath() string {
	return filepath.Join(getConfigDir(), "rb-scrobbler.key")
}

/* Open saved key from disk */
func getSavedKey() (string, error) {
	keyPath := getKeyFilePath()
	keyFile, err := os.Open(keyPath)

	if err != nil {
		return "", err
	}
	defer keyFile.Close()

	keyInBytes, err := io.ReadAll(keyFile)
	if err != nil {
		return "", err
	}

	return string(keyInBytes), nil
}
