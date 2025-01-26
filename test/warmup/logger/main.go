package main

import (
	"fmt"
	"log"
	"os"
)

// Create a global logger
var logger *log.Logger

func init() {
	// Open or create the log file
	file, err := os.OpenFile("app_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a logger that writes to the file with timestamp and file info
	logger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
}

// Function that simulates an error
func doSomething() error {
	return fmt.Errorf("something went wrong in doSomething")
}

// Another function that simulates an error
func doAnotherThing() error {
	return fmt.Errorf("another error occurred in doAnotherThing")
}

func main() {
	if err := doSomething(); err != nil {
		logger.Println(err)
	}

	if err := doAnotherThing(); err != nil {
		logger.Println(err)
	}

	fmt.Println("Errors have been logged to app_log.log")
}
