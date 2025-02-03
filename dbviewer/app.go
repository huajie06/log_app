package dbviewer

import (
	"fmt"
	"log"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

// the package is to read a bbolt DB and show the values on a html

// page 1: list all the bucket, or maybe first 20, with an option to search for a bucket
// page 2: show the key-value pair of the bucket, also maybe first 20, with an option to search for a key
// TODO: add edit functions

type DBManager struct {
	db     *bolt.DB
	logger *log.Logger
}

func NewDBManager(DBLocation string) (*DBManager, error) {

	logFile, err := os.OpenFile("db-viewer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile, "DBManager: ", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := bolt.Open(DBLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Printf("Error opening database: %v", err)
	}
	logger.Println("Database succesfully opened")

	return &DBManager{db: db, logger: logger}, nil
}

func FetchBucketName(m *DBManager) ([]string, error) {
	m.logger.Println("fetching from database...")

	done := make(chan struct {
		result []string
		err    error
	})
	go func() {
		var result []string
		err := m.db.View(func(tx *bolt.Tx) error {
			tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				result = append(result, string(name))
				return nil
			})
			return nil
		})
		done <- struct {
			result []string
			err    error
		}{result: result, err: err}
	}()

	res := <-done

	return res.result, res.err

}

func FetchBucketKeyVal(m *DBManager, bucketName string) (map[string]string, error) {

	done := make(chan struct {
		result map[string]string
		err    error
	})

	go func() {
		result := make(map[string]string)
		err := m.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %s not found", bucketName)
			}

			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				result[string(k)] = string(v)
			}
			return nil
		})
		done <- struct {
			result map[string]string
			err    error
		}{result: result, err: err}
	}()

	res := <-done

	return res.result, res.err

}
