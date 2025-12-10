package main

//The Channel is the BRIDGE, Not Part of Subscriber or Broker!

//Analogy: TCP Connection

// DISTRIBUTED SYSTEM (RabbitMQ/Kafka):

// Broker Process          TCP Connection          Subscriber Process
// ┌──────────────┐       ┌──────────┐            ┌────────────────┐
// │              │       │          │            │                │
// │  Broker      │──────►│   TCP    │───────────►│  Subscriber    │
// │  (writes)    │       │ (network)│            │  (reads)       │
// │              │       │          │            │                │
// └──────────────┘       └──────────┘            └────────────────┘

// TCP Connection is SHARED:
// - Broker writes to it
// - Subscriber reads from it
// - Neither "owns" it
// - Its the communication medium

// Our Implementation: Channel
// IN-PROCESS SYSTEM (Our Code):

// Broker              Channel                 Subscriber
// ┌──────────────┐   ┌──────────┐            ┌────────────────┐
// │              │   │          │            │                │
// │  Broker      │──►│ channel  │───────────►│  Subscriber    │
// │  (writes)    │   │ (memory) │            │  (reads)       │
// │              │   │          │            │                │
// └──────────────┘   └──────────┘            └────────────────┘

// Channel is SHARED:
// - Broker sends to it: channel <- msg
// - Subscriber reads from it: msg := <-channel
// - Neither "owns" it exclusively
// - It's the communication medium

// Who Creates the Channel?
// Current code:
// func NewSubscriber(id string, handler func(Message), broker *PubSub) *Subscriber {
//     s := &Subscriber{
//         id:      id,
//         channel: make(chan Message, 10), // ← Subscriber creates it
//         handler: handler,
//         broker:  broker,
//     }

//     go s.listen()
//     return s
// }

// func (s *Subscriber) SubscribeTo(topicName string) error {
//     // Subscriber gives channel to broker
//     return s.broker.Subscribe(topicName, s.id, s.channel) // ← Shares channel
// }
// ```

// **Subscriber creates the channel BUT then shares it with broker!**

// ---

// ## This is Like TCP Connection Setup
// ```
// TCP Connection Setup:

// 1. Subscriber: "I want to connect"
//    → Opens socket (subscriber.connect())

// 2. Connection established
//    → TCP connection created

// 3. Subscriber: "Here's my connection, send me messages"
//    → Gives connection handle to broker

// 4. Broker: "OK, I'll write to this connection"
//    → Stores connection

// 5. Data flows through SHARED connection
//    → Broker writes, subscriber reads
// ```

// **In our code:**
// ```
// Channel Setup:

// 1. Subscriber: "I want to receive messages"
//    → Creates channel

// 2. Channel exists in memory

// 3. Subscriber: "Here's my channel, send me messages"
//    → Gives channel to broker

// 4. Broker: "OK, I'll write to this channel"
//    → Stores channel

// 5. Data flows through SHARED channel
//    → Broker writes, subscriber reads
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ============================================
// SHARED DATA STRUCTURES
// Used by both broker and subscribers
// ============================================

type Message struct {
	Topic  string
	Data   interface{}
	Offset int64 // Position in topic's message log (like Kafka offset)
}

// ============================================
// BROKER SIDE STRUCTURES
// These represent the broker's view of the world
// ============================================

type Topic struct {
	name          string
	subscriptions []*Subscription // List of all subscriptions to this topic
	messages      []Message       // Message log - stores ALL messages for replay capability
	nextOffset    int64           // Counter for assigning offsets to new messages
	mu            sync.RWMutex
}

// Subscription - BROKER's representation of a subscriber
// This is the MINIMAL information broker needs about each subscriber
// Broker doesn't know about handler functions or processing logic
// Broker only needs to know WHERE to send messages (the channel)
type Subscription struct {
	subscriberID string
	channel      chan Message // ← Channel to send messages to
	// In actual distributed implementation (RabbitMQ/Kafka):
	// This channel would be replaced with a TCP connection
	// Broker would write messages over network to subscriber
}

// PubSub - THE MESSAGE BROKER
// This is the centralized message routing system
// Responsibilities:
//   - Manage topics
//   - Store messages
//   - Route messages to subscriber channels
//
// NOT responsible for:
//   - Starting subscriber goroutines
//   - Executing subscriber handlers
//   - Managing subscriber lifecycle
type PubSub struct {
	topics map[string]*Topic // All topics managed by this broker
	mu     sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		topics: make(map[string]*Topic),
	}
}

// ============================================
// BROKER METHODS
// These are the broker's API
// ============================================

// CreateTopic - BROKER method
// Creates a new topic for message publishing
func (ps *PubSub) CreateTopic(topicName string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.topics[topicName]
	if ok {
		return errors.New("Topic already exists")
	}

	topic := &Topic{
		name:          topicName,
		subscriptions: []*Subscription{},
		messages:      []Message{}, // Empty message log
		nextOffset:    0,           // Start from offset 0
	}
	ps.topics[topicName] = topic
	return nil
}

// Subscribe - BROKER method
// Registers a subscriber to receive messages from a topic
// Parameters:
//   - topicName: which topic to subscribe to
//   - subscriberID: unique identifier for the subscriber
//   - ch: the channel where broker will SEND messages
//
// IMPORTANT: Broker only stores the channel, nothing else!
// Broker doesn't know about:
//   - Subscriber's handler function
//   - Subscriber's goroutine
//   - How subscriber processes messages
//
// In distributed systems (RabbitMQ/Kafka):
//   - Instead of channel, this would be a TCP connection
//   - Broker would PUSH messages over network
func (ps *PubSub) Subscribe(topicName string, subscriberID string, ch chan Message) error {
	ps.mu.RLock()
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock()

	if !ok {
		return errors.New("Topic does not exist")
	}

	// Create broker's representation of subscriber
	// Just ID + channel, nothing more!
	subscription := &Subscription{
		subscriberID: subscriberID,
		channel:      ch, // ← Broker will SEND messages here
	}

	topic.mu.Lock()
	topic.subscriptions = append(topic.subscriptions, subscription)
	topic.mu.Unlock()

	fmt.Printf("[Broker] Subscriber '%s' registered to topic '%s'\n", subscriberID, topicName)

	return nil
}

// Unsubscribe - BROKER method
// Removes a subscriber from a topic
func (ps *PubSub) Unsubscribe(topicName string, subscriberID string) error {
	ps.mu.RLock()
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock()

	if !ok {
		return errors.New("Topic does not exist")
	}

	topic.mu.Lock()
	defer topic.mu.Unlock()

	var updatedSubscriptions []*Subscription
	var removedChannel chan Message

	for _, sub := range topic.subscriptions {
		if sub.subscriberID != subscriberID {
			updatedSubscriptions = append(updatedSubscriptions, sub)
		} else {
			removedChannel = sub.channel
		}
	}

	topic.subscriptions = updatedSubscriptions

	// Close the channel to signal subscriber to stop
	if removedChannel != nil {
		close(removedChannel)
		fmt.Printf("[Broker] Subscriber '%s' unsubscribed from topic '%s'\n", subscriberID, topicName)
		return nil
	}

	return errors.New("Subscriber not found")
}

// Publish - BROKER method
// Publishes a message to all subscribers of a topic
//
// BROKER's responsibilities:
//  1. Store message in log (for replay capability)
//  2. Send message to each subscriber's channel
//
// BROKER does NOT:
//   - Execute subscriber handlers
//   - Wait for subscribers to process
//   - Manage subscriber processing
//
// This is a PUSH-based mechanism:
//   - Broker actively SENDS messages to subscribers
//   - Subscribers passively RECEIVE from their channels
//   - Compare to PULL (Kafka): subscribers actively request messages
func (ps *PubSub) Publish(topicName string, data interface{}) error {
	ps.mu.RLock()
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock()

	if !ok {
		return errors.New("Topic does not exist")
	}

	topic.mu.Lock()

	// Create message with offset
	msg := Message{
		Topic:  topicName,
		Data:   data,
		Offset: topic.nextOffset,
	}

	// Store in message log (NEVER delete - allows replay)
	topic.messages = append(topic.messages, msg)
	topic.nextOffset++

	subscriptions := topic.subscriptions

	topic.mu.Unlock()

	fmt.Printf("[Broker] Publishing to topic '%s' at offset %d\n", topicName, msg.Offset)

	// PUSH to all subscriber channels
	// This is the PUSH mechanism!
	for _, sub := range subscriptions {
		select {
		case sub.channel <- msg: // ← BROKER PUSHES message
			// Sent successfully
		default:
			// Channel full - subscriber is slow or blocked
			fmt.Printf("[Broker Warning] Subscriber '%s' channel full, dropping message\n", sub.subscriberID)
		}
	}

	return nil
}

// ResetOffset - BROKER method
// Replays messages from a specific offset to a subscriber
// This allows subscribers to "rewind" and reprocess messages
// Useful for:
//   - Error recovery
//   - Reprocessing historical data
//   - New subscribers catching up
func (ps *PubSub) ResetOffset(topicName, subscriberID string, newOffset int64) error {
	ps.mu.RLock()
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock()

	if !ok {
		return errors.New("Topic does not exist")
	}

	topic.mu.RLock()
	defer topic.mu.RUnlock()

	// Find subscription
	var targetChannel chan Message
	for _, sub := range topic.subscriptions {
		if sub.subscriberID == subscriberID {
			targetChannel = sub.channel
			break
		}
	}

	if targetChannel == nil {
		return errors.New("Subscriber not found")
	}

	fmt.Printf("[Broker] Replaying messages for '%s' from offset %d\n", subscriberID, newOffset)

	// Replay messages from newOffset onwards
	for i := newOffset; i < int64(len(topic.messages)); i++ {
		msg := topic.messages[i]
		select {
		case targetChannel <- msg: // ← PUSH historical messages
			// Replayed successfully
		default:
			fmt.Printf("[Broker Warning] Channel full during replay\n")
		}
	}

	return nil
}

// ============================================
// CLIENT/SUBSCRIBER SIDE STRUCTURES
// These are completely separate from broker
// This is what application developers use
// ============================================

// Subscriber - CLIENT side structure
// This represents a message consumer in the application
//
// SUBSCRIBER's responsibilities:
//   - Listen on its own channel for messages
//   - Process messages with its handler function
//   - Manage its own goroutine lifecycle
//
// In this in-process implementation:
//   - channel is a Go channel (in-memory queue)
//   - Broker writes to channel, subscriber reads from channel
//
// In actual distributed implementation (RabbitMQ/Kafka):
//   - This channel would be replaced by TCP connection
//   - Broker would PUSH messages over network (TCP write)
//   - Subscriber would READ from network (TCP read)
//   - Example: broker.conn.Write(message) → subscriber.conn.Read()
type Subscriber struct {
	id      string
	channel chan Message // ← Each subscriber has its own queue
	// In actual distributed systems:
	//conn -> see layer_1 rabbit_mq.go inside brokers folder
	//for understanidn in very detail
	// - This would be a TCP connection
	// - Broker PUSHes by writing to connection
	// - Subscriber RECEIVEs by reading from connection
	// - Push-based: Broker initiates sending
	handler func(Message) // ← Subscriber's business logic
	broker  *PubSub       // ← Reference to broker (to call broker methods)
}

// NewSubscriber - CLIENT constructor
// Creates a new subscriber that is ALREADY listening
//
// Key point: Subscriber manages its OWN lifecycle
// - Starts its own goroutine immediately
// - Broker does NOT start this goroutine
// - Subscriber is ready to receive messages before subscribing
func NewSubscriber(id string, handler func(Message), broker *PubSub) *Subscriber {
	s := &Subscriber{
		id:      id,
		channel: make(chan Message, 10), // Buffered channel (queue)
		handler: handler,
		broker:  broker,
	}

	// SUBSCRIBER starts its OWN goroutine
	// This is CLIENT side responsibility, NOT broker's!
	go s.listen()

	fmt.Printf("[Subscriber %s] Created and listening\n", id)

	return s
}

// SubscribeTo - CLIENT method
// Subscriber registers itself with the broker
// Gives broker its channel so broker can PUSH messages
//
// =====================================
// IMPORTANT: Channel vs TCP Connection
// =====================================
//
// In this in-process implementation:
//   - We EXPLICITLY pass channel as parameter: s.broker.Subscribe(..., s.channel)
//   - Channel is passed directly in function call
//   - Broker receives channel and stores it
//   - Each subscriber has its OWN channel (separate queue)
//
// In distributed systems (TCP-based like RabbitMQ):
//   - Connection is IMPLICIT, NOT passed in request data
//   - Each subscriber has its OWN TCP connection (just like own channel!)
//   - Subscriber sends request OVER existing connection
//   - Broker already HAS the connection (from Accept())
//
// Detailed Flow (TCP-based):
//
//  1. Subscriber: net.Dial("tcp", "localhost:5672") → opens connection
//     NOTE: This happens inside Consumer.Subscribe() method
//     See: layer_2/brokers/rabbitmq.go for complete implementation
//     Just as we create different channel for each subscriber,
//     we create different TCP connection for each subscriber
//
//  2. Broker: listener.Accept() → receives connection as 'conn'
//     Each Accept() returns NEW connection for that subscriber
//
//  3. Broker: go handleClient(conn) → passes conn to handler
//     Separate goroutine per subscriber connection
//
//  4. Subscriber: encoder.Encode(SubscribeRequest{topic, id}) → sends THROUGH conn
//     Request data does NOT contain connection - it travels THROUGH connection
//
//  5. Broker (inside handleClient): decoder.Decode(&request) → receives FROM conn
//     Broker reads request that came through conn
//
//  6. Broker: knows request came from 'conn' (handleClient function parameter)
//     Connection is in scope, not in request data!
//
//  7. Broker: Subscribe(topic, id, conn) → uses conn from handleClient scope
//     Passes the connection it already has
//
// Connection Per Subscriber:
//
//	Subscriber 1 → Connection 1 → handleClient(conn1) goroutine
//	Subscriber 2 → Connection 2 → handleClient(conn2) goroutine
//	Subscriber 3 → Connection 3 → handleClient(conn3) goroutine
//	(Just like each subscriber has own channel in our implementation!)
//
// Key Difference:
//
//	Channel-based: broker.Subscribe(topic, id, channel) ← channel is EXPLICIT parameter
//	TCP-based:     Inside handleClient(conn), all requests come FROM that conn
//	               Connection is IMPLICIT - broker already has it from Accept()
//	               No need to pass conn in request data!
//
// Analogy:
//
//	Channel-based: "Here's my phone number (channel), call me on this"
//	TCP-based:     [Already on phone] "Send me messages on THIS line we're talking on"
//
// Reference: See layer_2/brokers/rabbitmq.go for complete TCP implementation
func (s *Subscriber) SubscribeTo(topicName string) error {
	// Call BROKER's Subscribe method
	// Pass OUR channel to broker (EXPLICIT in this implementation)
	// Broker will PUSH messages to this channel
	return s.broker.Subscribe(topicName, s.id, s.channel)
}

// UnsubscribeFrom - CLIENT method
// Subscriber unregisters from broker
func (s *Subscriber) UnsubscribeFrom(topicName string) error {
	return s.broker.Unsubscribe(topicName, s.id)
}

// listen - CLIENT private method
// This is the subscriber's message processing loop
// Runs in its own goroutine
//
// IMPORTANT: Broker NEVER calls this method!
// This is entirely subscriber's internal implementation
//
// Flow:
//  1. Wait for message on channel (BLOCKS if empty)
//  2. Message arrives (broker PUSHED it)
//  3. Execute handler (process message)
//  4. Go back to step 1
//
// In distributed systems:
//   - Instead of reading from channel: conn.Read()
//   - Blocks waiting for broker to PUSH over network
//   - Same concept, different transport mechanism
func (s *Subscriber) listen() {
	for msg := range s.channel {
		// Execute handler - this is subscriber's business logic
		// Broker doesn't know about this function
		// Broker doesn't execute this function
		s.handler(msg)
	}
	fmt.Printf("[Subscriber %s] Stopped listening (channel closed)\n", s.id)
}

// ============================================
// MAIN - USER APPLICATION
// This shows how a developer would use the system
// ============================================

func main() {
	fmt.Println("=== MESSAGE BROKER WITH CLEAR SEPARATION ===\n")

	// ============================================
	// STEP 1: Create BROKER (Server)
	// ============================================

	broker := NewPubSub()

	broker.CreateTopic("gaming")
	broker.CreateTopic("sports")
	broker.CreateTopic("news")

	fmt.Println("\n✅ Broker created with topics\n")

	// ============================================
	// STEP 2: Create SUBSCRIBERS (Clients)
	// Each subscriber is an independent client application
	// Each manages its own processing
	// ============================================

	// Subscriber 1 - Email notification service
	sub1 := NewSubscriber("notification-service", func(msg Message) {
		fmt.Printf("  [Notification] offset=%d: %v\n", msg.Offset, msg.Data)
		// Actual implementation would send email/SMS/push notification
	}, broker)

	// Subscriber 2 - Analytics service
	sub2 := NewSubscriber("analytics-service", func(msg Message) {
		fmt.Printf("  [Analytics] offset=%d: %v\n", msg.Offset, msg.Data)
		// Actual implementation would log to analytics database
	}, broker)

	// Subscriber 3 - Email service
	sub3 := NewSubscriber("email-service", func(msg Message) {
		fmt.Printf("  [Email] offset=%d: %v\n", msg.Offset, msg.Data)
		// Actual implementation would send confirmation email
	}, broker)

	fmt.Println("\n✅ Subscribers created (already listening on their channels)\n")

	// ============================================
	// STEP 3: SUBSCRIBERS register with BROKER
	// Each subscriber gives broker its channel
	// Broker will PUSH messages to these channels
	// ============================================

	sub1.SubscribeTo("gaming")
	sub2.SubscribeTo("gaming")
	sub3.SubscribeTo("gaming")

	fmt.Println("\n✅ All subscribers registered (broker knows where to PUSH)\n")

	// ============================================
	// STEP 4: PUBLISH messages
	// Broker PUSHes messages to all subscriber channels
	// Subscribers receive and process independently
	// ============================================

	fmt.Println("--- Publishing Messages (PUSH-based) ---\n")

	// Message 1
	broker.Publish("gaming", "New video: Elden Ring Gameplay")
	time.Sleep(100 * time.Millisecond)

	// Message 2
	broker.Publish("gaming", "New video: God of War Review")
	time.Sleep(100 * time.Millisecond)

	// Message 3
	broker.Publish("gaming", "New video: Cyberpunk 2077")
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n✅ All subscribers received and processed messages\n")

	// ============================================
	// STEP 5: REPLAY messages (Offset Reset)
	// Demonstrates replay capability
	// Email service will reprocess all messages
	// ============================================

	fmt.Println("--- Testing Replay (Reset to Offset 0) ---\n")

	broker.ResetOffset("gaming", "email-service", 0)
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n✅ Email service replayed all messages from beginning\n")

	// ============================================
	// STEP 6: UNSUBSCRIBE
	// Remove notification service from topic
	// ============================================

	fmt.Println("--- Unsubscribing Notification Service ---\n")

	sub1.UnsubscribeFrom("gaming")
	time.Sleep(100 * time.Millisecond)

	// ============================================
	// STEP 7: Publish after unsubscribe
	// Only 2 subscribers should receive this
	// ============================================

	fmt.Println("--- Publishing After Unsubscribe ---\n")

	broker.Publish("gaming", "New video: Starfield")
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n✅ Only analytics and email received (notification unsubscribed)\n")

	fmt.Println("--- Done ---")
}

/*
============================================
KEY CONCEPTS SUMMARY
============================================

1. BROKER (PubSub):
   - Manages topics and message routing
   - Stores messages for replay
   - PUSHes messages to subscriber channels
   - Does NOT execute subscriber handlers
   - Does NOT manage subscriber lifecycle

2. SUBSCRIPTION (Broker's view):
   - Just subscriber ID + channel
   - Broker doesn't know about handler or processing logic
   - Minimal information needed for routing

3. SUBSCRIBER (Client's view):
   - Has channel for receiving messages
   - Has handler for processing messages
   - Manages own goroutine
   - Registers itself with broker

4. PUSH vs PULL:
   - PUSH (this implementation): Broker sends to subscriber
   - PULL (Kafka-style): Subscriber requests from broker

5. In-Process vs Distributed:
   - In-Process (this): Channel-based, same process
   - Distributed (RabbitMQ/Kafka): Network-based, TCP connections

6. Responsibilities:
   - Broker: Route messages
   - Subscriber: Process messages
   - Clear separation!
============================================
*/
