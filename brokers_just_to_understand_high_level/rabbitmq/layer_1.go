package main

// ============================================
// FILE: rabbitmq-broker/main.go
// Run: go run rabbitmq-broker/main.go
// Port: 5672
// ============================================

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// ============================================
// MESSAGE & REQUEST STRUCTURES
// ============================================

type Message struct {
	Queue string
	Data  string
}

type Request struct {
	Type       string // "SUBSCRIBE", "UNSUBSCRIBE", "PUBLISH"
	Queue      string
	ConsumerID string
	Data       string
}

type Response struct {
	Type  string
	Error string
}

// ============================================
// BROKER
// ============================================

type Broker struct {
	queues    map[string]*Queue
	consumers map[string]*ConsumerConnection // ← STORES ACTIVE CONNECTIONS!
	mu        sync.RWMutex
}

type Queue struct {
	name     string
	messages []Message
	mu       sync.RWMutex
}

// KEY STRUCTURE: Stores ACTUAL connection to consumer
// WHY? Because we need to PUSH messages to consumer later
type ConsumerConnection struct {
	consumerID string
	queue      string
	conn       net.Conn      // ← TCP connection to consumer (KEEP OPEN!)
	encoder    *json.Encoder // ← To PUSH messages over network
	mu         sync.Mutex
}

func NewBroker() *Broker {
	// STEP: Create broker with empty maps
	// WHY queues map? Store all queues by name
	// WHY consumers map? Store all active consumer connections (KEY for PUSH!)
	return &Broker{
		queues:    make(map[string]*Queue),
		consumers: make(map[string]*ConsumerConnection),
	}
}

// ============================================
// CREATE QUEUE
// ============================================

func (b *Broker) CreateQueue(name string) {
	// STEP: Lock broker to safely modify queues map
	// WHY Lock? Multiple goroutines might create queues simultaneously
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.queues[name]; !exists {
		// STEP: Create new queue with empty messages slice
		// WHY messages slice? Store messages that haven't been consumed yet
		b.queues[name] = &Queue{
			name:     name,
			messages: []Message{},
		}
		fmt.Printf("[Broker] Created queue: %s\n", name)
	}
}

// ============================================
// SUBSCRIBE - Consumer registers for PUSH
// ============================================

func (b *Broker) Subscribe(consumerID, queueName string, conn net.Conn, encoder *json.Encoder) error {
	// STEP 1: Lock broker to safely modify consumers map
	// WHY? Multiple consumers might subscribe simultaneously
	b.mu.Lock()
	defer b.mu.Unlock()

	// STEP 2: Store consumer's connection
	// WHY store conn? So we can PUSH messages to this consumer later!
	// WHY store encoder? To serialize messages and send over conn
	// THIS IS THE KEY DIFFERENCE FROM KAFKA!
	// Kafka: Doesn't store connections (consumers PULL)
	// RabbitMQ: MUST store connections (broker PUSHes)
	b.consumers[consumerID] = &ConsumerConnection{
		consumerID: consumerID,
		queue:      queueName,
		conn:       conn,    // ← Keep connection OPEN!
		encoder:    encoder, // ← Use this to PUSH!
	}

	fmt.Printf("[Broker] ✅ Consumer '%s' registered to queue '%s' (connection stored)\n",
		consumerID, queueName)

	return nil
}

// ============================================
// UNSUBSCRIBE
// ============================================

func (b *Broker) Unsubscribe(consumerID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// STEP 1: Find consumer
	if consumer, exists := b.consumers[consumerID]; exists {
		// STEP 2: Close connection
		// WHY? No longer need to push to this consumer
		consumer.conn.Close()

		// STEP 3: Remove from map
		// WHY? Free memory, consumer no longer active
		delete(b.consumers, consumerID)

		fmt.Printf("[Broker] Consumer '%s' unsubscribed\n", consumerID)
	}
}

// ============================================
// PUBLISH - Producer sends message
// ============================================

func (b *Broker) Publish(queueName, data string) error {
	// STEP 1: Get queue
	// WHY RLock? We're only reading queues map
	b.mu.RLock()
	queue, exists := b.queues[queueName]
	b.mu.RUnlock()

	if !exists {
		return fmt.Errorf("queue not found")
	}

	// STEP 2: Create message
	msg := Message{
		Queue: queueName,
		Data:  data,
	}

	// STEP 3: Store in queue
	// WHY store? In case we need to redeliver, or consumer is slow
	queue.mu.Lock()
	queue.messages = append(queue.messages, msg)
	queue.mu.Unlock()

	fmt.Printf("[Broker] 📥 Received message for queue '%s': %s\n", queueName, data)

	// STEP 4: IMMEDIATELY PUSH to all subscribed consumers
	// ⚡ THIS IS THE KEY PUSH MECHANISM!
	// WHY immediately? PUSH model = broker initiates delivery
	// Compare to Kafka: Messages just sit in partition, consumers PULL when ready
	b.pushToConsumers(queueName, msg)

	return nil
}

// ============================================
// PUSH TO CONSUMERS - The PUSH mechanism
// ============================================

func (b *Broker) pushToConsumers(queueName string, msg Message) {
	// STEP 1: Get all consumers
	// WHY RLock? We're reading consumers map
	b.mu.RLock()
	consumers := b.consumers
	b.mu.RUnlock()

	// STEP 2: Find consumers subscribed to this queue
	// WHY loop? Multiple consumers might be subscribed to same queue
	for _, consumer := range consumers {
		if consumer.queue == queueName {
			// STEP 3: PUSH to this consumer
			// WHY goroutine? Don't block if one consumer is slow
			// Each consumer gets message independently
			go b.pushMessage(consumer, msg)
		}
	}
}

func (b *Broker) pushMessage(consumer *ConsumerConnection, msg Message) {
	// STEP 1: Lock this consumer's encoder
	// WHY? Multiple messages might be pushed to same consumer simultaneously
	consumer.mu.Lock()
	defer consumer.mu.Unlock()

	// STEP 2: PUSH message over network
	// HOW? encoder.Encode() serializes msg to JSON and writes to TCP connection
	// IMPORTANT: Broker INITIATES this! Consumer doesn't ask for it!
	err := consumer.encoder.Encode(msg)
	if err != nil {
		fmt.Printf("[Broker] ❌ Failed to push to '%s': %v\n", consumer.consumerID, err)
		return
	}

	fmt.Printf("[Broker] 📤 Pushed message to '%s': %s\n", consumer.consumerID, msg.Data)
}

// ============================================
// NETWORK SERVER
// ============================================

func (b *Broker) Start(port string) {
	// STEP 1: Create TCP listener
	// WHY TCP? Reliable, ordered, connection-oriented
	// Binds to port, starts listening for connections
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[Broker] 🚀 RabbitMQ listening on %s\n", port)

	// STEP 2: Create default queue
	b.CreateQueue("orders")

	// STEP 3: Accept loop - wait for clients
	// WHY infinite loop? Server runs forever
	for {
		// BLOCKS here waiting for client to connect

		// Accept connection - broker gets the connection here!
		conn, err := listener.Accept()
		//  ^^^^
		//  Broker now has connection!
		if err != nil {
			continue
		}

		fmt.Printf("[Broker] 🔌 Client connected: %s\n", conn.RemoteAddr())

		// Handle this connection
		go b.handleClient(conn)
		//                ^^^^
		//                Pass connection to handler
	}
}

// Visual: How Connection is Known
// ```
//                     ONE TCP CONNECTION

// Subscriber                                          Broker
//     |                                                 |
//     | net.Dial("tcp", "localhost:5672")               |
//     |────────────────────────────────────────────────►| listener.Accept()
//     |                                                 | → conn = [connection object]
//     |                                                 |
//     |                                                 | go handleClient(conn)
//     |                                                 |    ▼
//     |                                        ┌────────────────────┐
//     |                                        │ handleClient(conn) │
//     |                                        │                    │
//     | encoder.Encode(SUB request)            │ decoder.Decode()   │
//     |───────────────────────────────────────►│ ← reads from conn │
//     |       (data flows through conn)        │                    │
//     |                                        │ Oh! This request   │
//     |                                        │ came from 'conn'!  │
//     |                                        │                    │
//     |                                        │ b.Subscribe(       │
//     |                                        │   topic,           │
//     |                                        │   id,              │
//     |                                        │   conn ← this one! │
//     |                                        │ )                  │
//     |                                        └────────────────────┘
//     |                                                 |
//     |                                                 | Now broker has:
//     |                                                 | Subscription{
//     |                                                 |   id: "sub1"
//     |                                                 |   conn: [this connection]
//     |                                                 | }

// ## Analogy

// **Channel-based (Explicit):**
// ```
// You: "Hey broker, here's my phone number (channel)"
// Broker: "Got it! I'll call this number when I have messages"
// ```

// **TCP-based (Implicit):**
// ```
// You: [Already on phone call with broker]
// You: "I want to subscribe to 'orders'"
// Broker: "Got it! I'll keep talking to you on THIS phone line"
func (b *Broker) handleClient(conn net.Conn) {
	//                         ^^^^^^^^^^^^^^
	//                         This conn belongs to ONE subscriber
	//                         All requests on this conn come from SAME subscriber
	// STEP 1: Ensure connection closes when function exits
	// WHY defer? Cleanup even if panic occurs
	defer conn.Close()

	// STEP 2: Create JSON encoder/decoder for this connection
	// WHY? To send/receive structured data (JSON) over TCP
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	// STEP 3: Request handling loop
	// WHY loop? Handle multiple requests on same connection
	for {
		var req Request

		// STEP 4: Wait for request from client
		// BLOCKS here until client sends something
		err := decoder.Decode(&req)
		if err != nil {
			fmt.Printf("[Broker] Client disconnected\n")
			return
		}

		var resp Response

		// STEP 5: Handle request based on type
		switch req.Type {
		case "SUBSCRIBE":
			// Consumer wants to subscribe
			// Pass: consumerID, queueName, connection, encoder
			// WHY pass conn & encoder? Need to store them for PUSHING later!
			// Pass the connection we already have!
			// We're inside handleClient(conn), so we have access to conn

			// Pass the connection we already have!
			// We're inside handleClient(conn), so we have access to conn
			err := b.Subscribe(req.ConsumerID, req.Queue, conn, encoder)
			if err != nil {
				resp.Error = err.Error()
			}

			// STEP 6: Send acknowledgment
			encoder.Encode(resp)

			// ⚠️ IMPORTANT: DON'T RETURN HERE!
			// WHY? Connection must stay OPEN for pushing messages later!
			// Loop continues, waits for next request (or nothing - stays open)

		case "UNSUBSCRIBE":
			// Consumer unsubscribing
			b.Unsubscribe(req.ConsumerID)
			return // NOW we return and close connection

		case "PUBLISH":
			// Producer publishing message
			err := b.Publish(req.Queue, req.Data)
			if err != nil {
				resp.Error = err.Error()
			}
			encoder.Encode(resp)
		}
	}
}

func main() {
	// STEP 1: Create broker
	broker := NewBroker()

	// STEP 2: Start server (blocks forever)
	broker.Start(":5672")
}
