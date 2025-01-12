package data

import (
	"fmt"
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
