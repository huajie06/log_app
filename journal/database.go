package journal

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

func NewDBManager(DBLocation, LogfileLocation string) (*DBManager, error) {
	logFile, err := os.OpenFile(LogfileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile, "DBManager: ", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := bolt.Open(DBLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Printf("Error opening database: %v", err)
	}
	logger.Println("Database succesfully created")

	return &DBManager{db: db, logger: logger}, nil
}

func (dbm *DBManager) Close() error {
	err := dbm.db.Close()
	if err != nil {
		dbm.logger.Printf("Error close database: %v", err)
	} else {
		dbm.logger.Printf("DB cloded")
	}
	return err
}

func (m *DBManager) DBput(ctx context.Context, BucketName, Key, Value string) error {
	m.logger.Printf("DBput: bucket=%s, key=%s", BucketName, Key)

	done := make(chan error)

	go func() {
		done <- m.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(BucketName))
			if err != nil {
				return fmt.Errorf("fail to create bucket: %s", BucketName)
			}
			if err := b.Put([]byte(Key), []byte(Value)); err != nil {
				return fmt.Errorf("fail to put key-value pair: %w", err)
			}
			return nil
		})
	}()

	select {
	case <-ctx.Done():
		m.logger.Printf("DBPut canncelled or timeout: Bucket=%s, Key=%s", BucketName, Key)
		return ctx.Err()
	case err := <-done:
		if err != nil {
			m.logger.Printf("DBput error: %v", err)
		} else {
			m.logger.Println("DBput successful")
		}
		return err
	}

}
func (m *DBManager) DBget(ctx context.Context, BucketName, Key string) (string, error) {
	m.logger.Printf("DBget: BucketName=%s, Key=%s", BucketName, Key)

	done := make(chan struct {
		value string
		err   error
	})

	go func() {
		var result string
		err := m.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BucketName))
			if b == nil {
				return fmt.Errorf("bucket: %s not found", BucketName)
			}
			value := b.Get([]byte(Key))
			if value == nil {
				return fmt.Errorf("key: %s not found", Key)
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
		m.logger.Printf("DBget cancelled or timeout, BucketName=%s, Key=%s", BucketName, Key)
		return "", ctx.Err()
	case res := <-done:
		if res.err != nil {
			m.logger.Println("DBget error: ", res.err)
		} else {
			m.logger.Printf("DBget success on  BucketName=%s, Key=%s", BucketName, Key)
		}
		return res.value, res.err
	}
}
