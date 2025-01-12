package main

import (
	"fmt"
	"log_app/data"
)

func main() {
	path := "test_file.db"
	// data.DBinit(path)

	eventType := "read"
	eventDate := "2024_1231"
	eventTime := "afternoon"
	eventContent := "some long reading task completed"
	logTimestamp := "2025_1231_000000"

	key := eventType + "-" + logTimestamp
	value := eventTime + "-" + eventContent
	data.DBStore(path, eventDate, key, value)

	eventDate = "2024_1230"
	data.DBStore(path, eventDate, key, value)

	eventDate = "2025_1230"
	data.DBStore(path, eventDate, key, value)

	fmt.Println("------all bucket-------------")
	data.DBLoopBucket(path)

	fmt.Println("------ lookup something -------------")

	data.DBRetrieve(path, "2025_1230", key)

}
