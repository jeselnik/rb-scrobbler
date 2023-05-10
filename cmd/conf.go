package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DEF_PATH    = ""
	DEF_OFFSET  = "0h"
	DEF_NONINT  = ""
	DEF_COLOURS = true
)

type config struct {
	Path, Offset, NonInteractive string
	Colours                      bool
}

func getConfigFileDir() string {
	return filepath.Join(getConfigDir(), "rb-scrobbler.conf")
}

func createDefaultConfig() error {
	c := config{DEF_PATH, DEF_OFFSET, DEF_NONINT, DEF_COLOURS}
	inJson, err := json.Marshal(c)
	if err != nil {
		return err
	}

	writeErr := os.WriteFile(getConfigFileDir(), inJson, os.ModePerm)
	return writeErr
}

func openConfigFile() (config, error) {
	var c config

	file, err := os.Open(getConfigFileDir())

	if errors.Is(err, os.ErrNotExist) {
		createDefaultConfig()
		file, err = os.Open(getConfigFileDir())
	}

	if err != nil {
		return c, err
	}
	defer file.Close()

	fileToBytes, convError := ioutil.ReadAll(file)
	if convError != nil {
		return c, err
	}

	unmarErr := json.Unmarshal(fileToBytes, &c)
	return c, unmarErr
}
