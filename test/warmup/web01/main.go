package main

import (
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

	router.LoadHTMLFiles("/home/hj/apps/log_app/journal/src/event-form.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "event-form.html", gin.H{"APIRouteForEventForm": "/api/event"})
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
