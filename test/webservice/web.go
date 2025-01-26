package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

type EventLog struct {
	EventType    string `json:"eventType" binding:"required"`
	EventDate    string `json:"eventDate" binding:"required"`
	EventTime    string `json:"eventTime"`
	EventContent string `json:"eventContent"`
	LogTimestamp string `json:"logTimestamp" binding:"required"`
}

type DBManager struct {
	db     *bolt.DB
	logger *log.Logger
}

func NewDBManager(DBLocation string) (*DBManager, error) {
	logFile, err := os.OpenFile("app_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := log.New(logFile, "DBManager: ", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := bolt.Open(DBLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Printf("Error opening database: %v", err)
		return nil, err
	}

	logger.Println("Database connection successfully opened")
	return &DBManager{db: db, logger: logger}, nil
}

func (m *DBManager) Close() error {
	err := m.db.Close()
	if err != nil {
		m.logger.Printf("Error closing database: %v", err)
	} else {
		m.logger.Println("Database connection successfully closed")
	}
	return err
}

// DBPut stores a key-value pair in the specified bucket and logs the operation.
func (m *DBManager) DBPut(ctx context.Context, BucketName, Key, Value string) error {
	m.logger.Printf("DBPut: Bucket=%s, Key=%s", BucketName, Key)

	done := make(chan error, 1)
	go func() {
		done <- m.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(BucketName))
			if err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}
			if err := b.Put([]byte(Key), []byte(Value)); err != nil {
				return fmt.Errorf("failed to put key-value pair: %w", err)
			}
			return nil
		})
	}()

	select {
	case <-ctx.Done():
		m.logger.Printf("DBPut cancelled or timed out: Bucket=%s, Key=%s", BucketName, Key)
		return ctx.Err()
	case err := <-done:
		if err != nil {
			m.logger.Printf("DBPut error: %v", err)
		} else {
			m.logger.Println("DBPut completed successfully")
		}
		return err
	}
}

// DBGet retrieves the value for a given key from the specified bucket and logs the operation.
func (m *DBManager) DBGet(ctx context.Context, BucketName, Key string) (string, error) {
	m.logger.Printf("DBGet: Bucket=%s, Key=%s", BucketName, Key)

	done := make(chan struct {
		value string
		err   error
	}, 1)

	go func() {
		var result string
		err := m.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BucketName))
			if b == nil {
				return fmt.Errorf("bucket %s not found", BucketName)
			}
			value := b.Get([]byte(Key))
			if value == nil {
				return fmt.Errorf("key %s not found in bucket %s", Key, BucketName)
			}
			result = string(value)
			return nil
		})
		done <- struct {
			value string
			err   error
		}{result, err}
	}()

	select {
	case <-ctx.Done():
		m.logger.Printf("DBGet cancelled or timed out: Bucket=%s, Key=%s", BucketName, Key)
		return "", ctx.Err()
	case res := <-done:
		if res.err != nil {
			m.logger.Printf("DBGet error: %v", res.err)
		} else {
			m.logger.Printf("DBGet completed successfully: Value=%s", res.value)
		}
		return res.value, res.err
	}
}

func (m *DBManager) LogHandler(c *gin.Context) {
	var event EventLog
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		m.logger.Printf("Invalid input: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second) // 2-second timeout
	defer cancel()

	bucketName := event.EventDate
	key := event.EventType + "_" + event.LogTimestamp
	value := event.EventTime + "_" + event.EventContent

	if err := m.DBPut(ctx, bucketName, key, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log event"})
		m.logger.Printf("Failed to log event: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event logged successfully"})
	m.logger.Println("Event logged successfully")
}
