package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Order represents a customer order
type Order struct {
	ID          string  `json:"id"`
	CustomerID  string  `json:"customer_id"`
	Amount      float64 `json:"amount"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

// Global channels for order processing
var (
	newOrders       = make(chan Order, 100) // Buffer for incoming orders
	processedOrders = make(chan Order, 100) // Buffer for processed orders
	errorChannel    = make(chan error, 10)  // Channel for error handling
)

func main() {
	// Start background workers
	go orderProcessor()
	go errorHandler()
	go orderLogger()

	// Set up HTTP endpoints
	http.HandleFunc("/order", handleNewOrder)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleNewOrder handles incoming HTTP requests for new orders
func handleNewOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send order to processing queue
	newOrders <- order

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Order queued for processing",
		"id":     order.ID,
	})
}

// orderProcessor simulates processing multiple orders concurrently
func orderProcessor() {
	for order := range newOrders {
		// Start a new goroutine for each order
		go func(ord Order) {
			// Simulate processing time
			time.Sleep(time.Second * 2)

			// Simulate potential errors
			if ord.Amount <= 0 {
				errorChannel <- fmt.Errorf("invalid amount for order %s", ord.ID)
				return
			}

			// Process the order
			ord.ProcessedAt = time.Now().Format(time.RFC3339)

			// Send to processed orders channel
			processedOrders <- ord

		}(order)
	}
}

// errorHandler handles processing errors
func errorHandler() {
	for err := range errorChannel {
		log.Printf("Error processing order: %v", err)
		// In a real application, you might:
		// - Send to error monitoring service
		// - Trigger alerts
		// - Attempt retry logic
	}
}

// orderLogger logs processed orders and simulates saving to database
func orderLogger() {
	// Use WaitGroup for batch processing
	var wg sync.WaitGroup
	batch := make([]Order, 0, 10)

	for order := range processedOrders {
		batch = append(batch, order)

		// Process in batches of 10
		if len(batch) >= 10 {
			ordersBatch := make([]Order, len(batch))
			copy(ordersBatch, batch)

			wg.Add(1)
			go func(orders []Order) {
				defer wg.Done()
				saveBatchToDB(orders)
			}(ordersBatch)

			batch = batch[:0] // Clear the batch
		}

		log.Printf("Processed order: %s for customer: %s, amount: %.2f",
			order.ID, order.CustomerID, order.Amount)
	}
}

// saveBatchToDB simulates saving orders to a database
func saveBatchToDB(orders []Order) {
	// Simulate database operation
	time.Sleep(time.Millisecond * 500)
	log.Printf("Saved batch of %d orders to database", len(orders))
}
