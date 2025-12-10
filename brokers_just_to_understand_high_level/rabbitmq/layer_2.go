// Layer 2: RABBITMQ CLIENT LIBRARY

// ============================================
// FILE: rabbitmqclient/client.go
// Package: rabbitmqclient
// Users import: "yourname/rabbitmqclient"
// ============================================

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Message struct {
	Queue string
	Data  string
}

type Request struct {
	Type       string
	Queue      string
	ConsumerID string
	Data       string
}

type Response struct {
	Type  string
	Error string
}

// ============================================
// CONSUMER
// ============================================

type Consumer struct {
	brokerAddr string
	consumerID string
	queue      string
	conn       net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	msgChannel chan Message // ← Receives PUSHED messages from broker
	mu         sync.Mutex
}

func NewConsumer(brokerAddr, consumerID string) (*Consumer, error) {
	// STEP: Create consumer object
	// WHY msgChannel with buffer 100? To queue messages pushed by broker
	// If broker pushes faster than app processes, messages buffer here
	return &Consumer{
		brokerAddr: brokerAddr,
		consumerID: consumerID,
		msgChannel: make(chan Message, 100),
	}, nil
}

// ============================================
// SUBSCRIBE
// ============================================

func (c *Consumer) Subscribe(queue string) error {
	c.queue = queue

	// STEP 1: Open TCP connection to broker
	// WHY? Need network connection to communicate with broker
	// This initiates TCP 3-way handshake
	conn, err := net.Dial("tcp", c.brokerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	// STEP 2: Store connection and create encoder/decoder
	// WHY encoder? To send requests (JSON) to broker
	// WHY decoder? To receive responses/messages (JSON) from broker
	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)

	fmt.Printf("[Consumer %s] 🔌 Connected to broker\n", c.consumerID)

	// STEP 3: Send SUBSCRIBE request to broker
	// Tells broker: "I want messages from queue 'orders'"
	req := Request{
		Type:       "SUBSCRIBE",
		Queue:      queue,
		ConsumerID: c.consumerID,
	}

	err = c.encoder.Encode(req)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}

	// STEP 4: Wait for acknowledgment from broker
	// BLOCKS until broker responds
	// WHY wait? Need confirmation that subscription succeeded
	var resp Response
	err = c.decoder.Decode(&resp)
	if err != nil {
		return fmt.Errorf("failed to receive ack: %v", err)
	}

	if resp.Error != "" {
		return fmt.Errorf("subscribe error: %s", resp.Error)
	}

	fmt.Printf("[Consumer %s] ✅ Subscribed to '%s'\n", c.consumerID, queue)

	// STEP 5: Start goroutine to listen for PUSHED messages
	// WHY goroutine? Can't block Subscribe() waiting for messages
	// Subscribe() needs to return so main() can continue
	// Background goroutine waits for broker to push
	go c.listenForPushedMessages()

	return nil
}

// ============================================
// LISTEN FOR PUSHED MESSAGES
// This is THE KEY method for PUSH model
// ============================================

func (c *Consumer) listenForPushedMessages() {
	fmt.Printf("[Consumer %s] 👂 Waiting for broker to push messages...\n", c.consumerID)

	// STEP: Infinite loop waiting for pushes
	// WHY infinite? Should receive multiple messages over time
	// Only exits when connection closes
	for {
		var msg Message

		// STEP 1: BLOCK waiting for broker to push
		// This is PASSIVE - we don't ask, we just wait
		// Broker will push when message arrives
		// Compare to Kafka: Consumer would actively call Poll() here
		err := c.decoder.Decode(&msg)
		if err != nil {
			// Connection closed or error
			fmt.Printf("[Consumer %s] Connection closed\n", c.consumerID)
			close(c.msgChannel)
			return
		}

		fmt.Printf("[Consumer %s] 📨 Received PUSHED message: %s\n", c.consumerID, msg.Data)

		// STEP 2: Forward message to application via channel
		// WHY channel? Decouples network I/O from application logic
		// Application reads from channel at its own pace
		c.msgChannel <- msg
	}
}

// ============================================
// MESSAGES - Get channel to receive messages
// ============================================

func (c *Consumer) Messages() <-chan Message {
	// STEP: Return read-only channel
	// WHY read-only (<-chan)? Application should only read, not write
	// Prevents misuse
	return c.msgChannel
}

// ============================================
// UNSUBSCRIBE
// ============================================

func (c *Consumer) Unsubscribe() error {
	// STEP 1: Send UNSUBSCRIBE request
	req := Request{
		Type:       "UNSUBSCRIBE",
		ConsumerID: c.consumerID,
	}

	c.encoder.Encode(req)

	// STEP 2: Close connection
	// WHY? No longer need to receive messages
	// Broker will remove us from consumers map
	c.conn.Close()
	return nil
}

// ============================================
// PRODUCER
// ============================================

type Producer struct {
	brokerAddr string
	conn       net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	mu         sync.Mutex
}

func NewProducer(brokerAddr string) (*Producer, error) {
	// STEP 1: Connect to broker
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, err
	}

	// STEP 2: Create producer with encoder/decoder
	return &Producer{
		brokerAddr: brokerAddr,
		conn:       conn,
		encoder:    json.NewEncoder(conn),
		decoder:    json.NewDecoder(conn),
	}, nil
}

func (p *Producer) Publish(queue, data string) error {
	// STEP 1: Create PUBLISH request
	req := Request{
		Type:  "PUBLISH",
		Queue: queue,
		Data:  data,
	}

	// STEP 2: Send request to broker
	// WHY Lock? Multiple goroutines might publish simultaneously
	p.mu.Lock()
	err := p.encoder.Encode(req)
	if err != nil {
		p.mu.Unlock()
		return err
	}

	// STEP 3: Wait for acknowledgment
	var resp Response
	err = p.decoder.Decode(&resp)
	p.mu.Unlock()

	if err != nil {
		return err
	}

	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}

	fmt.Printf("[Producer] 📤 Published: %s\n", data)

	return nil
}

func (p *Producer) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
