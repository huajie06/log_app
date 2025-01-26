package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

type DBManager struct {
	db     *bolt.DB
	logger *log.Logger
}

// NewDBManager initializes a new DBManager with a single BoltDB connection.
func NewDBManager(DBLocation string) (*DBManager, error) {
	// Create a log file
	logFile, err := os.OpenFile("app_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New(logFile, "DBManager: ", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := bolt.Open(DBLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Printf("Error opening database: %v", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	logger.Println("Database connection successfully opened")

	return &DBManager{db: db, logger: logger}, nil
}

// Close closes the BoltDB connection and logs the event.
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
