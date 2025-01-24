package main

import (
	"log"
	"log_app/journal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	htmlDir, err := journal.HtmlDir()
	if err != nil {
		panic(err)
	}
	router.LoadHTMLFiles(htmlDir)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "event-form.html", gin.H{"APIRouteForEventForm": "/api/eventlog"})
	})

	dbManager, err := journal.NewDBManager("journal_event.db", "journal_app.log")
	if err != nil {
		log.Fatalf("Error initializing DBManager: %v", err)
	}
	defer dbManager.Close()

	router.POST("/api/eventlog/", dbManager.EventLogHandler)

	// router.Run("localhost:8001")
	router.Run()
}
