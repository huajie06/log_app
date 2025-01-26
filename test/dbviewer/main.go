package main

import (
	"fmt"
	dv "log_app/dbviewer"
)

func main() {
	// /home/hj/apps/log_app/test/journal/journal_event.db

	dbPath := "/home/hj/apps/log_app/test/journal/journal_event.db"

	m, err := dv.NewDBManager(dbPath)
	if err != nil {
		fmt.Println("error initiate conn")
	}

	result, err := dv.FetchBucketName(m)
	if err != nil {
		fmt.Println("error fetch db conn")
	}
	for i, v := range result {
		fmt.Printf("index: %d, buckename: %v\n", i, v)
	}

	bucketName := "event-log"
	bres, err := dv.FetchBucketKeyVal(m, bucketName)
	if err != nil {
		fmt.Println("error fetch bucket: ", err)
	}

	for k, v := range bres {
		fmt.Printf("key is :%v, value is :%v\n", k, v)
	}

}
