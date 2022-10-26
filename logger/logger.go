package main

import (
	"log"
	"os"
)

func newLog(name string, logger *log.Logger) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	logger.SetOutput(file)
}
