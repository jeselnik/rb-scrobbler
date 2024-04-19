/* This set of functions deals with getting and saving the
session key in a cross OS manner */

package main

import (
	"io"
	"os"
	"path/filepath"
)

func getConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "rb-scrobbler"), nil
}

func getKeyFilePath() (string, error) {
	configDir, err := getConfigDir()
	return filepath.Join(configDir, "rb-scrobbler.key"), err
}

/* Open saved key from disk */
func getSavedKey() (string, error) {
	keyPath, err := getKeyFilePath()
	if err != nil {
		return "", err
	}

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
