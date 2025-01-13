package main

import (
	"context"
	"fmt"
	"log"
	"log_app/data"
	"time"
)

func main() {
	// test1()

	// fmt.Println("======================================")
	// fmt.Println("=========== test2 below================")
	// fmt.Println("======================================")
	// test2()

	// fmt.Println("======================================")
	// data.DBLoopBucket("test_file.db")

	fmt.Println("========= print all data=================")
	test3()

}

func test3() {
	// data.DBLoopBucket("/home/hj/apps/log_app/cmd/web/app.db")
	data.DBPrintAll("/home/hj/apps/log_app/cmd/web/app.db")

}

func test2() {
	dbManager, err := data.NewDBManager("test_file.db")
	if err != nil {
		log.Fatalf("Error initializing DBManager: %v", err)
	}
	defer dbManager.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Insert key-value pair
	if err := dbManager.DBPut(ctx, "ExampleBucket", "Key1", "Value1"); err != nil {
		log.Printf("Error putting key-value pair: %v", err)
	} else {
		fmt.Println("Key-Value pair successfully inserted")
	}

	// Retrieve the value for a key
	value, err := dbManager.DBGet(ctx, "ExampleBucket", "Key1")
	if err != nil {
		log.Printf("Error getting key-value pair: %v", err)
	} else {
		fmt.Printf("Retrieved value: %s\n", value)
	}
}

func test1() {
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
