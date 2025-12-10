package kafka

// ============================================
// FILE: order-service/main.go
// ============================================

// Layer 3: USER APPLICATION

import (
	"fmt"
	"time"

	"yourname/kafkaclient"
)

func main() {
	fmt.Println("=== KAFKA (PULL MODEL) ===\n")

	// Create consumer
	consumer, _ := kafkaclient.NewConsumer(
		"localhost:9092",
		"order-processors",
		"consumer-1",
	)
	defer consumer.Close()

	// Subscribe to partition 0
	consumer.Subscribe("orders", 0)

	// Start consuming in background
	go func() {
		for {
			// PULL from broker
			msg, err := consumer.Poll(1000)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if msg != nil {
				fmt.Printf("\n[📦 Order Service] Received: %s\n", msg.Value)
				processOrder(msg)

				// Commit progress
				consumer.Commit()
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	// Create producer
	producer, _ := kafkaclient.NewProducer("localhost:9092")
	defer producer.Close()

	// Produce messages
	fmt.Println("\n--- Publishing Orders ---\n")

	producer.Produce("orders", "order_1001", "iPhone 15 - $999")
	time.Sleep(1 * time.Second)

	producer.Produce("orders", "order_1002", "MacBook Pro - $2499")
	time.Sleep(1 * time.Second)

	producer.Produce("orders", "order_1003", "AirPods - $249")
	time.Sleep(2 * time.Second)

	fmt.Println("\n--- Done ---")
}

func processOrder(msg *kafkaclient.Message) {
	fmt.Printf("  💳 Processing payment\n")
	fmt.Printf("  📧 Sending email\n")
	fmt.Printf("  ✅ Done!\n")
}
