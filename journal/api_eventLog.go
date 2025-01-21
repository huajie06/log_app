package journal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func HtmlDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	res, err := findLogDirectory(dir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return filepath.Join(res, "./journal/src/event-form.html"), nil
}

func findLogDirectory(startPath string) (string, error) {
	maxLevels := 10
	targetFolder := "log_app"

	// Get absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	currentPath := absPath
	levelsUp := 0

	for levelsUp < maxLevels {
		// List directory contents
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return "", fmt.Errorf("failed to read directory %s: %w", currentPath, err)
		}

		// Check current directory for target folder
		for _, entry := range entries {
			if entry.IsDir() && strings.EqualFold(entry.Name(), targetFolder) {
				return filepath.Join(currentPath, entry.Name()), nil
			}
		}

		// Move up one directory level
		parentPath := filepath.Dir(currentPath)

		// Check if we've reached the root directory
		if parentPath == currentPath {
			return "", fmt.Errorf("reached root directory without finding %s", targetFolder)
		}

		currentPath = parentPath
		levelsUp++
	}

	return "", fmt.Errorf("directory %s not found within %d levels up from %s",
		targetFolder, maxLevels, startPath)
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
