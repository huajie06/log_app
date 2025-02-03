package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

var cache sync.Map

type Job struct {
	BucketName string
	ResultChan chan map[string]string
	ErrorChan  chan error
}

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

	bucketTest := []string{"abc", "bcd", "xxx", "yyy"}
	pairs := generateRandomPairs(20)

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

func worker(m *DBManager, jobs <-chan Job) {
	for job := range jobs {
		result, err := FetchBucketKeyValWithCache(m, job.BucketName)
		job.ResultChan <- result
		job.ErrorChan <- err
	}
}

func startWorkerPool(m *DBManager, numWorkers int) chan Job {
	jobs := make(chan Job, 100) // Buffer to hold pending jobs
	for i := 0; i < numWorkers; i++ {
		go worker(m, jobs)
	}
	return jobs
}

func FetchBucketKeyValWithCache(m *DBManager, bucketName string) (map[string]string, error) {
	// Check the cache first
	if cachedResult, ok := cache.Load(bucketName); ok {
		return cachedResult.(map[string]string), nil
	}

	// If not in cache, fetch from the database
	result, err := FetchBucketKeyVal(m, bucketName)
	if err != nil {
		return nil, err
	}

	// Store the result in the cache
	cache.Store(bucketName, result)
	return result, nil
}

const ProjectRoot = "/home/hj/apps/log_app/"

func main() {
	dbPath := "test-dbviewer.db"
	createDummyDB(dbPath)

	m, err := NewDBManager(dbPath)
	if err != nil {
		fmt.Println("Error initiating connection:", err)
		return
	}

	jobs := startWorkerPool(m, 10) // Start a pool of 10 workers

	r := gin.Default()
	htmlFilePath := filepath.Join(ProjectRoot, "dbviewer/*.html")
	r.LoadHTMLGlob(htmlFilePath)

	r.GET("/db/:bucketname", func(c *gin.Context) {
		bucketName := c.Param("bucketname")

		job := Job{
			BucketName: bucketName,
			ResultChan: make(chan map[string]string, 1),
			ErrorChan:  make(chan error, 1),
		}
		jobs <- job

		select {
		case result := <-job.ResultChan:
			c.HTML(http.StatusOK, "keyvals.html", gin.H{
				"bucketname": bucketName,
				"data":       result,
			})
		case err := <-job.ErrorChan:
			c.String(http.StatusInternalServerError, "Error fetching bucket: %v", err)
		}
	})

	r.Run()
}
