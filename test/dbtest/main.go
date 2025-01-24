package main

import (
	"encoding/csv"
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
	dpath := "/home/hj/apps/log_app/cmd/journal/journal_event.db"
	defaultConn := newConn(dpath)

	db, err := openDB(defaultConn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	getDBbucket(db)

	getDBBucketSize(db, "event-log")

	getDBBucketKeys(db, "event-log")

	// viewDBBucketValueRaw(db, "event-log")

	getDBBucketValue(db, "event-log")

	// WriteBucketToCSV(db, "event-log", "event-log.csv")

	// WriteBucketToJSON(db, "event-log", "event-log.json")
}

func WriteBucketToJSON(db *bolt.DB, bucketName string, outfile string) {
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
				keyResult += string(k) + string(v)
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
		fmt.Println("not return results from db")
		panic(res.err)
	} else {
		fmt.Printf("database bucket: %s, has values\n", bucketName)

		procResult := returnDBbucket(res.value)
		fmt.Println("--------------------")
		// for _, v := range procResult {
		// 	fmt.Printf("%+v\n", v)
		// }

		file, err := os.Create(outfile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close() // Ensure the file is closed when the function exits

		// Convert to JSON and write to the file
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ") // Pretty print with indentation

		if err := encoder.Encode(procResult); err != nil {
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
		}{keyResult, err}
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
				keyResult += string(k) + string(v)
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
		fmt.Printf("database bucket: %s, has values\n", bucketName)
		fmt.Println(res.value)
		// procResult := returnDBbucket(res.value)
		// fmt.Println("--------------------")
		// for _, v := range procResult {
		// 	fmt.Printf("%+v\n", v)
		// }

	}
}

func printDBbucket(data string) {
	eachEntry := strings.Split(data, UnitSeparator)
	for _, v := range eachEntry {
		keyPairValue := strings.Split(v, FiledSeparator)

		r1 := journal.EventLog{LogTimestamp: keyPairValue[0],
			EventType:    keyPairValue[1],
			EventDate:    keyPairValue[2],
			EventTime:    keyPairValue[3],
			EventContent: keyPairValue[4]}
		fmt.Printf("%+v\n", r1)
	}
}

func returnDBbucket(data string) []journal.EventLog {
	eachEntry := strings.Split(data, UnitSeparator)

	var result []journal.EventLog

	for _, v := range eachEntry {
		keyPairValue := strings.Split(v, FiledSeparator)

		r1 := journal.EventLog{LogTimestamp: keyPairValue[0],
			EventType:    keyPairValue[1],
			EventDate:    keyPairValue[2],
			EventTime:    keyPairValue[3],
			EventContent: keyPairValue[4]}

		result = append(result, r1)
	}
	return result
}

func WriteBucketToCSV(db *bolt.DB, bucketName string, csvfile string) error {
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
		fmt.Println("fail to conver to csv")
		return res.err
	} else {

		procResult := returnDBbucket(res.value)
		file, err := os.Create(csvfile)
		if err != nil {
			fmt.Println("csv creation fail")
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"LogTimestamp", "EventType", "EventDate", "EventTime", "EventContent"}
		writer.Write(headers)

		for _, v := range procResult {
			row := []string{
				v.LogTimestamp,
				v.EventType,
				v.EventDate,
				v.EventTime,
				v.EventContent,
			}
			writer.Write(row)
		}
		return nil
	}

}
