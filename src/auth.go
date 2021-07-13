/* This set of functions deals with getting and saving the
session key in a cross OS manner */

package main

import (
	"io/ioutil"
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
func getSavedKey() string {
	keyPath := getKeyFilePath()
	keyFile, err := os.Open(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	keyInBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	return string(keyInBytes)
}
