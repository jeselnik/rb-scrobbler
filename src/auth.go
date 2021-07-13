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

func getSavedKey() string {
	tokenPath := getKeyFilePath()
	tokenFile, err := os.Open(tokenPath)
	if err != nil {
		log.Fatal(err)
	}

	tokenInBytes, err := ioutil.ReadAll(tokenFile)
	if err != nil {
		log.Fatal(err)
	}

	return string(tokenInBytes)
}
