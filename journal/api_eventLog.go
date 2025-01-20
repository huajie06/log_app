package journal

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// from GenAI tool
const (
	// Record separator (RS) for separating entire records
	EntrySeparator = "\u001E"
	// Unit separator (US) for separating fields within a record
	FiledSeparator = "\u001F"
	// Group separator (GS) for separating items in arrays
	UnitSeparator = "\u001D"
)

type EventLog struct {
	EventType    string `json:"eventType" binding:"required"`
	EventDate    string `json:"eventDate" binding:"required"`
	EventTime    string `json:"eventTime"`
	EventContent string `json:"eventContent"`
	LogTimestamp string `json:"logTimestamp"`
}

func SerializeLog(e EventLog) (string, string) {
	dbKey := strings.Join([]string{e.LogTimestamp, e.EventType}, FiledSeparator)
	dbValue := strings.Join([]string{e.EventDate, e.EventTime, e.EventContent}, FiledSeparator)

	return dbKey, dbValue
}

func (m *DBManager) EventLogHandler(c *gin.Context) {
	var event EventLog

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second) // 2-second timeout
	defer cancel()

	bucketName := "event-log"
	dbKey, dbValue := SerializeLog(event)

	if err := m.DBput(ctx, bucketName, dbKey, dbValue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log event"})
		m.logger.Printf("Failed to log event: %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event logged successfully"})
	m.logger.Println("Event logged successfully")
}
