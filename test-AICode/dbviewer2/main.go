package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
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

// DBManager manages the BoltDB connection
type DBManager struct {
	db *bolt.DB
}

// NewDBManager creates a new DBManager instance
func NewDBManager(dbPath string) (*DBManager, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &DBManager{db: db}, nil
}

// FetchBucketName fetches all bucket names from the database
func FetchBucketName(m *DBManager) ([]string, error) {
	var buckets []string
	err := m.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bucket names: %w", err)
	}
	return buckets, nil
}

// FetchBucketKeyVal fetches all key-value pairs from a bucket
func FetchBucketKeyVal(m *DBManager, bucketName string) (map[string]string, error) {
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
	if err != nil {
		return nil, fmt.Errorf("failed to fetch key-value pairs: %w", err)
	}
	return result, nil
}

// createDummyDB creates a dummy database with random data
func createDummyDB(dbPath string) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bucketTest := []string{"aaa", "bbb", "ccc", "ddd"}
	pairs := generateRandomPairs(10)

	for _, v := range bucketTest {
		err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return fmt.Errorf("failed to create bucket: %s", v)
			}

			for key, value := range pairs {
				err := b.Put([]byte(key), []byte(value))
				if err != nil {
					return fmt.Errorf("failed to put key-value pair: %w", err)
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

const ProjectRoot = "/home/hj/apps/log_app/"

// ViewerWeb starts the web server to view the database
func ViewerWeb(dbPath string) {
	r := gin.Default()

	htmlFilePath := filepath.Join(ProjectRoot, "dbviewer/*.html")
	r.LoadHTMLGlob(htmlFilePath)

	m, err := NewDBManager(dbPath)
	if err != nil {
		fmt.Println("Error initiating connection:", err)
		return
	}

	result, err := FetchBucketName(m)
	if err != nil {
		fmt.Println("Error fetching bucket names:", err)
		return
	}

	// Route to list all buckets
	r.GET("/db", func(c *gin.Context) {
		c.HTML(http.StatusOK, "buckets.html", gin.H{
			"dbfile":     dbPath,
			"bucketlist": result,
		})
	})

	// Route to show all key-value pairs for a specific bucket
	r.GET("/db/:bucketname", func(c *gin.Context) {
		bucketName := c.Param("bucketname")

		bres, err := FetchBucketKeyVal(m, bucketName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching bucket: %v", err)
			return
		}

		c.HTML(http.StatusOK, "keyvals.html", gin.H{
			"bucketname": bucketName,
			"data":       bres,
		})
	})

	r.Run()
}

func main() {
	dbPath := "test-dbviewer.db"
	createDummyDB(dbPath)
	ViewerWeb(dbPath)
}
