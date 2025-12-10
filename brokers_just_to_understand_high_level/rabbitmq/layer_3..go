package main

// Layer 3: USER APPLICATION
// ============================================
// FILE: order-service/main.go
// This is what YOU (customer) write
// ============================================

import (
	"fmt"
	"time"

	"yourname/rabbitmqclient"
)

func main() {
	fmt.Println("=== RABBITMQ (PUSH MODEL) ===\n")

	// ============================================
	// CONSUMER SETUP
	// ============================================

	// STEP 1: Create consumer object
	// WHY? Need object to interact with broker
	// This just creates object in memory, no network call yet
	consumer, err := rabbitmqclient.NewConsumer(
		"localhost:5672", // Where broker is
		"email-service",  // Unique ID for this consumer
	)
	if err != nil {
		panic(err)
	}

	// STEP 2: Subscribe to queue
	// What happens:
	// - Opens TCP connection to broker
	// - Sends SUBSCRIBE request
	// - Broker stores our connection
	// - Broker can now PUSH to us
	// - Starts goroutine listening for pushes
	err = consumer.Subscribe("orders")
	if err != nil {
		panic(err)
	}

	// STEP 3: Start goroutine to process PUSHED messages
	// WHY goroutine? Don't block main() waiting for messages
	// This runs in background, processes messages as they arrive
	go func() {
		// STEP 4: Loop receiving messages
		// WHY for range? Channel pattern for receiving multiple values
		// BLOCKS when channel empty, waits for next message
		for msg := range consumer.Messages() {
			// Application receives message pushed by broker
			fmt.Printf("\n[📧 Email Service] Received: %s\n", msg.Data)

			// YOUR business logic
			processOrder(msg)
		}
	}()

	// STEP 5: Wait a bit for consumer to fully initialize
	time.Sleep(1 * time.Second)

	// ============================================
	// PRODUCER SETUP
	// ============================================

	// STEP 6: Create producer
	// Opens separate TCP connection to broker
	producer, err := rabbitmqclient.NewProducer("localhost:5672")
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// ============================================
	// PUBLISH MESSAGES
	// ============================================

	fmt.Println("\n--- Publishing Orders ---\n")

	// STEP 7: Publish first message
	// What happens:
	// 1. Producer sends PUBLISH request to broker
	// 2. Broker receives, stores in queue
	// 3. Broker IMMEDIATELY pushes to all subscribed consumers
	// 4. Consumer's listenForPushedMessages() receives it
	// 5. Forwards to msgChannel
	// 6. Application goroutine processes it
	producer.Publish("orders", "Order #1001: iPhone 15")
	time.Sleep(1 * time.Second)

	// STEP 8: Publish more messages
	producer.Publish("orders", "Order #1002: MacBook Pro")
	time.Sleep(1 * time.Second)

	producer.Publish("orders", "Order #1003: AirPods")
	time.Sleep(2 * time.Second)

	// STEP 9: Cleanup
	consumer.Unsubscribe()

	fmt.Println("\n--- Done ---")
}

// ============================================
// YOUR BUSINESS LOGIC
// ============================================

func processOrder(msg rabbitmqclient.Message) {
	// YOUR code to handle the message
	fmt.Printf("  💳 Processing payment\n")
	fmt.Printf("  📧 Sending email\n")
	fmt.Printf("  ✅ Done!\n")
}
