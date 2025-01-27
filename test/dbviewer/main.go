package main

import (
	"fmt"
	"math/rand"
	"time"

	bolt "go.etcd.io/bbolt"

	dv "log_app/dbviewer"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	keyLength   = 10 // Length of the random key
	valueLength = 15 // Length of the random value
)

// randString generates a random string of a given length
func randString(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rng.Intn(len(letterBytes))]
	}
	return string(b)
}

// generateRandomPairs generates n random key-value pairs
func generateRandomPairs(n int) map[string]string {
	pairs := make(map[string]string)
	for i := 0; i < n; i++ {
		key := randString(keyLength)
		value := randString(valueLength)
		pairs[key] = value
	}
	return pairs
}

func main() {
	// /home/hj/apps/log_app/test/journal/journal_event.db

	// dbPath := "test-dbviewer.db"

	// createDummyDB(dbPath)

	dbPath := "/home/hj/apps/log_app/test/journal/journal_event.db"

	dv.ViewerWeb(dbPath)
}

func createDummyDB(dbPath string) {
	// dbPath := "test-dbviewer.db"

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	buketTest := []string{"abc", "bcd", "xxx", "yyy"}

	pairs := generateRandomPairs(20)

	for _, v := range buketTest {
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return fmt.Errorf("fail to create bucket: %s", v)
			}

			for key, value := range pairs {
				err := b.Put([]byte(key), []byte(value))
				if err != nil {
					fmt.Println("put key-value pair: %w", err)
					continue
				}
			}
			return nil
		})
	}
}
