package main

import (
	"encoding/json"
	"fmt"
	"log_app/journal"
	"os"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	FiledSeparator = "\u001F"
	UnitSeparator  = "\u001D"
)

type dbConn struct {
	dbPath  string
	mode    os.FileMode
	options *bolt.Options
}

func newConn(dbPath string) dbConn {
	return dbConn{
		dbPath:  dbPath,
		mode:    0600,
		options: &bolt.Options{Timeout: 1 * time.Second},
	}
}

func openDB(conn dbConn) (*bolt.DB, error) {
	db, err := bolt.Open(conn.dbPath, conn.mode, conn.options)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	dpath := "/home/hj/apps/log_app/test/journal/journal_event.db"
	defaultConn := newConn(dpath)

	db, err := openDB(defaultConn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// getDBbucket(db)

	getDBBucketSize(db, "event-log")

	// getDBBucketKeys(db, "event-log")

	// viewDBBucketValueRaw(db, "event-log")

	getDBBucketValue(db, "event-log")
	fmt.Println("=========================================================")

	WriteBucketToJSON(db, "event-log", "event-log.json")

}

func WriteBucketToJSON(db *bolt.DB, bucketName string, outfile string) {
	done := make(chan struct {
		value []map[string]journal.EventLog
		err   error
	})

	go func() {
		var keyValuerResult []map[string]journal.EventLog
		var kVal journal.EventLog
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				err := json.Unmarshal(v, &kVal)
				if err != nil {
					fmt.Println("fail to parse database into struct")
					return err
				}
				keyValuerResult = append(keyValuerResult, map[string]journal.EventLog{string(k): kVal})
			}
			return nil
		})
		done <- struct {
			value []map[string]journal.EventLog
			err   error
		}{keyValuerResult, err}
	}()

	res := <-done
	if res.err != nil {
		fmt.Println("not return results from db")
		panic(res.err)
	} else {
		fmt.Printf("database bucket: %s, has values\n", bucketName)

		for _, v := range res.value {
			fmt.Println("values are:", v)
		}

		file, err := os.Create(outfile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(res.value); err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
	}
}

func getDBbucket(db *bolt.DB) {

	done := make(chan struct {
		value string
		err   error
	})

	go func() {
		var result string
		err := db.View(func(tx *bolt.Tx) error {
			tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				result += string(name) + "\n"
				return nil
			})
			return nil
		})

		done <- struct {
			value string
			err   error
		}{result, err}
	}()

	res := <-done

	if res.err != nil {
		panic(res.err)
	}

	fmt.Println(res.value)
}

func getDBBucketKeys(db *bolt.DB, bucketName string) {
	done := make(chan struct {
		value string
		err   error
	})

	go func() {
		var keyResult string
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				keyResult += string(k) + "\n"
			}
			return nil
		})
		done <- struct {
			value string
			err   error
		}{keyResult, err}
	}()

	res := <-done
	if res.err != nil {
		panic(res.err)
	} else {
		fmt.Printf("database bucket: %s, has keys:\n-------------------------------\n%v\n", bucketName, res.value)
	}
}

func viewDBBucketValueRaw(db *bolt.DB, bucketName string) {
	done := make(chan struct {
		value string
		err   error
	})

	go func() {
		var keyResult string
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				keyResult += string(k) + FiledSeparator + string(v) + UnitSeparator
			}
			return nil
		})
		done <- struct {
			value string
			err   error
		}{strings.Trim(keyResult, UnitSeparator), err}
	}()

	res := <-done
	if res.err != nil {
		panic(res.err)
	} else {
		fmt.Printf("database bucket: %s, has values\n", bucketName)
		eachEntry := strings.Split(res.value, UnitSeparator)
		for i, v := range eachEntry {
			keyPairValue := strings.Split(v, FiledSeparator)
			fmt.Printf("index: %d, value: %v\n", i, keyPairValue)
		}

	}
}

func getDBBucketValue(db *bolt.DB, bucketName string) {
	done := make(chan struct {
		value string
		err   error
	})

	go func() {
		var keyResult string
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				keyResult += string(k) + FiledSeparator + string(v) + UnitSeparator
			}
			return nil
		})
		done <- struct {
			value string
			err   error
		}{strings.Trim(keyResult, UnitSeparator), err}
	}()

	res := <-done
	if res.err != nil {
		panic(res.err)
	} else {
		fmt.Printf("database bucket: %s, has values\n", bucketName)

		eachEntry := strings.Split(res.value, UnitSeparator)

		for i, v := range eachEntry {
			fmt.Printf("index: %d, value:%v\n", i, v)
			fmt.Println("------------------------------")
		}

	}
}

func getDBBucketSize(db *bolt.DB, bucketName string) {
	done := make(chan struct {
		value int
		err   error
	})

	go func() {
		var bucketSize int
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				bucketSize++
			}
			return nil
		})
		done <- struct {
			value int
			err   error
		}{bucketSize, err}
	}()

	res := <-done
	if res.err != nil {
		panic(res.err)
	} else {
		fmt.Printf("database bucket: %s, has size: %v\n", bucketName, res.value)
	}
}
