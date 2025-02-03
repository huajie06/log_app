package data

import (
	"fmt"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

func DBinit(path string) {

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

}

func DBRetrieve(dbPath, bucket, key string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		fmt.Println("db open error")
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	return err

}

func DBLoopBucket(dbPath string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("db open error")
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println(string(name))
			return nil
		})
		return nil
	})
	return err
}

func DBPrintAll_v0(dbPath string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("db open error")
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("==============================")
			fmt.Println(string(name))
			buck := tx.Bucket([]byte(string(name)))
			c := buck.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key=%s, value=%s\n", k, v)
			}
			fmt.Println("==============================")

			return nil
		})
		return nil
	})
	return err
}

func DBPrintAll(dbPath string) error {
	// Open the database with timeout
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// View transaction
	return db.View(func(tx *bolt.Tx) error {
		// Track if database is empty
		bucketCount := 0

		// Iterate over each bucket
		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucketCount++
			fmt.Printf("\n=== Bucket: %s ===\n", string(name))

			// Use the bucket passed to ForEach instead of getting it again
			c := b.Cursor()
			keyCount := 0

			for k, v := c.First(); k != nil; k, v = c.Next() {
				keyCount++
				// Handle nil values (nested buckets)
				if v == nil {
					fmt.Printf("Key: %s (nested bucket)\n", k)
					// Optionally handle nested buckets
					printNestedBucket(b.Bucket(k), 1)
				} else {
					fmt.Printf("Key: %s, Value: %s\n", k, v)
				}
			}

			fmt.Printf("Total keys in bucket: %d\n", keyCount)
			fmt.Println(strings.Repeat("=", 40))
			return nil
		})

		if bucketCount == 0 {
			fmt.Println("Database is empty (no buckets found)")
		} else {
			fmt.Printf("\nTotal buckets: %d\n", bucketCount)
		}

		return err
	})
}

// Helper function to print nested buckets with indentation
func printNestedBucket(b *bolt.Bucket, depth int) {
	if b == nil {
		return
	}

	indent := strings.Repeat("  ", depth)
	c := b.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v == nil {
			fmt.Printf("%s└── %s (nested bucket)\n", indent, k)
			printNestedBucket(b.Bucket(k), depth+1)
		} else {
			fmt.Printf("%s└── %s: %s\n", indent, k, v)
		}
	}
}

func DBStore(dbPath, bucket, key, value string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		// return fmt.Errorf("create bucket: %s", err)
		fmt.Println("db open error")
		return err
	}
	defer db.Close()

	// eventType, eventDate, eventTime, eventContent, logTimestamp := CreateTestData()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			// return fmt.Errorf("create bucket: %s", err)
			fmt.Println("bucket open/create error")
			return err
		}

		err = b.Put([]byte(key), []byte(value))
		if err != nil {
			fmt.Println("db save error")
			return err
		}
		return nil
	})
	return err
}

// create some test data

func CreateTestData() (r1, r2, r3, r4, r5 string) {

	eventType := "read"
	eventDate := "2024_1231"
	eventTime := "afternoon"
	eventContent := "some long reading task completed"
	logTimestamp := "2025_1231_000000"

	return eventType, eventDate, eventTime, eventContent, logTimestamp
}
