package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	flag.Parse()

	logFile, err := os.Open(*logPath)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
}
