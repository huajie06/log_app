package main

import (
	"encoding/json"
	"fmt"
)

type EventLog struct {
	EventType    string `json:"eventType" binding:"required"`
	EventDate    string `json:"eventDate" binding:"required"`
	EventTime    string `json:"eventTime"`
	EventContent string `json:"eventContent"`
	LogTimestamp string `json:"logTimestamp"`
}

func main() {

	eventlog := EventLog{
		EventType:    "run",
		EventDate:    "2022-01-01",
		EventTime:    "Afternoon",
		EventContent: "some randome stuff",
		LogTimestamp: "2024-01-01 12:00:00",
	}
	// fmt.Println(eventlog)

	m, err := json.Marshal(&eventlog)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)

	fmt.Println(string(m))
}
