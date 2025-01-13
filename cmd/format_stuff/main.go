package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello world")
	// Get the current time
	currentTime := time.Now()

	// Print the current time in default format
	fmt.Println("Current date and time:", currentTime)

	// Print the current time in a custom format
	fmt.Println("Formatted date and time:", currentTime.Format("2006-01-02 15:04:05"))

	fmt.Println("-------------------------")
	standardTime := "2024-01-12 15:04:05"
	fmt.Println("StandardTime date and time:", currentTime)
	t1, err := time.Parse("2006-01-02 15:04:05", standardTime)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Parsed standard time:", t1)
}
