package main

import (
	"log"
	web "log_app/webservice"

	"github.com/gin-gonic/gin"
)

func main() {
	dbManager, err := web.NewDBManager("app.db")
	if err != nil {
		log.Fatalf("Error initializing DBManager: %v", err)
	}
	defer dbManager.Close()

	r := gin.Default()
	r.POST("/log_today", dbManager.LogHandler)

	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
