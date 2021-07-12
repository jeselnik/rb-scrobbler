package main

import (
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
