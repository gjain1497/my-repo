// package main

// import (
// 	"errors"
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type Message struct {
// 	Topic  string
// 	Data   interface{}
// 	Offset int64 // ← NEW: Position in topic's message log
// }

// type Topic struct {
// 	name        string
// 	subscribers []*Subscriber
// 	messages    []Message // ← NEW: Store ALL messages (never delete!)
// 	nextOffset  int64     // ← NEW: Counter for offsets
// 	mu          sync.RWMutex
// }

// type PubSub struct {
// 	topics map[string]*Topic
// 	mu     sync.RWMutex
// }

// func NewPubSub() *PubSub {
// 	return &PubSub{
// 		topics: make(map[string]*Topic),
// 	}
// }

// func (ps *PubSub) CreateTopic(topicName string) error {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	_, ok := ps.topics[topicName]
// 	if ok {
// 		return errors.New("Topic already exists")
// 	}
// 	topic := &Topic{
// 		name:        topicName,
// 		subscribers: []*Subscriber{},
// 		messages:    []Message{}, // ← NEW: Empty message log
// 		nextOffset:  0,           // ← NEW: Start from 0
// 	}
// 	ps.topics[topicName] = topic
// 	return nil
// }

// func (ps *PubSub) Subscribe(topicName string, subscriber *Subscriber) error {
// 	ps.mu.RLock()
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock()

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	topic.mu.Lock()
// 	topic.subscribers = append(topic.subscribers, subscriber)
// 	topic.mu.Unlock()

// 	// Start subscriber
// 	subscriber.start()
// 	return nil
// }

// func (ps *PubSub) Unsubscribe(topicName string, subscriberId string) error {
// 	ps.mu.RLock()
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock()

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	topic.mu.Lock()
// 	var updatedSubscribers []*Subscriber
// 	var removedSubscriber *Subscriber
// 	for _, subscriber := range topic.subscribers {
// 		if subscriber.id != subscriberId {
// 			updatedSubscribers = append(updatedSubscribers, subscriber)
// 		} else {
// 			removedSubscriber = subscriber
// 		}
// 	}
// 	topic.subscribers = updatedSubscribers
// 	topic.mu.Unlock()

// 	if removedSubscriber != nil {
// 		removedSubscriber.Close()
// 	} else {
// 		return errors.New("Subscriber not found")
// 	}

// 	return nil
// }

// func (ps *PubSub) Publish(topicName string, data interface{}) error {
// 	ps.mu.RLock()
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock()

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	topic.mu.Lock()

// 	// NEW: Create message with offset
// 	msg := Message{
// 		Topic:  topicName,
// 		Data:   data,
// 		Offset: topic.nextOffset,
// 	}

// 	// NEW: Store message in log (NEVER delete!)
// 	topic.messages = append(topic.messages, msg)

// 	// NEW: Increment offset counter
// 	topic.nextOffset++

// 	// Get subscribers
// 	subscribers := topic.subscribers

// 	topic.mu.Unlock()

// 	// Publish to all subscribers
// 	for _, subscriber := range subscribers {
// 		// Non-blocking send (in case subscriber is slow)
// 		select {
// 		case subscriber.channel <- msg:
// 			// Sent successfully
// 		default:
// 			// Channel full, skip (or handle differently)
// 			fmt.Printf("[Warning] Subscriber %s channel full, dropping message\n", subscriber.id)
// 		}
// 	}

// 	return nil
// }

// // ============================================
// // NEW: RESET OFFSET - Replay messages!
// // ============================================

// func (ps *PubSub) ResetOffset(topicName, subscriberId string, newOffset int64) error {
// 	ps.mu.RLock()
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock()

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	topic.mu.RLock()
// 	defer topic.mu.RUnlock()

// 	// Find subscriber
// 	var targetSubscriber *Subscriber
// 	for _, subscriber := range topic.subscribers {
// 		if subscriber.id == subscriberId {
// 			targetSubscriber = subscriber
// 			break
// 		}
// 	}

// 	if targetSubscriber == nil {
// 		return errors.New("Subscriber not found")
// 	}

// 	// Replay messages from newOffset onwards
// 	fmt.Printf("[PubSub] 🔄 Replaying messages for '%s' from offset %d\n", subscriberId, newOffset)

// 	for i := newOffset; i < int64(len(topic.messages)); i++ {
// 		msg := topic.messages[i]

// 		// Send to subscriber (non-blocking)
// 		select {
// 		case targetSubscriber.channel <- msg:
// 			fmt.Printf("[PubSub] Replayed offset %d to '%s'\n", msg.Offset, subscriberId)
// 		default:
// 			fmt.Printf("[Warning] Subscriber %s channel full during replay\n", subscriberId)
// 		}
// 	}

// 	return nil
// }

// // ============================================
// // SUBSCRIBER
// // ============================================

// type Subscriber struct {
// 	id      string
// 	channel chan Message //each subscriber has its own queue, in actual impltmentation it is getting
// 	//each request throguh a tcp connection, broker is pushing into it (push based mechansim)
// 	handler func(Message)
// }

// func NewSubscriber(id string, handler func(Message)) *Subscriber {
// 	return &Subscriber{
// 		id:      id,
// 		channel: make(chan Message, 10),
// 		handler: handler,
// 	}
// }

// func (s *Subscriber) start() {
// 	go func() {
// 		for msg := range s.channel {
// 			s.handler(msg)
// 		}
// 	}()
// }

// func (s *Subscriber) Close() {
// 	close(s.channel)
// }

// // ============================================
// // MAIN - DEMONSTRATING ALL REQUIREMENTS
// // ============================================

// func main() {
// 	fmt.Println("=== MESSAGE QUEUE WITH ALL REQUIREMENTS ===\n")

// 	// Create pubsub system
// 	pubsub := NewPubSub()

// 	// ============================================
// 	// REQUIREMENT 1: Support multiple topics
// 	// ============================================
// 	pubsub.CreateTopic("gaming")
// 	pubsub.CreateTopic("sports")
// 	pubsub.CreateTopic("news")

// 	fmt.Println("✅ Multiple topics created: gaming, sports, news\n")

// 	// ============================================
// 	// REQUIREMENT 3: Subscribers subscribe to topic
// 	// ============================================
// 	sub1 := NewSubscriber("notification-service", func(msg Message) {
// 		fmt.Printf("[Notification] Received (offset=%d): %v\n", msg.Offset, msg.Data)
// 	})

// 	sub2 := NewSubscriber("analytics-service", func(msg Message) {
// 		fmt.Printf("[Analytics] Received (offset=%d): %v\n", msg.Offset, msg.Data)
// 	})

// 	sub3 := NewSubscriber("email-service", func(msg Message) {
// 		fmt.Printf("[Email] Received (offset=%d): %v\n", msg.Offset, msg.Data)
// 	})

// 	// Subscribe all to gaming topic
// 	pubsub.Subscribe("gaming", sub1)
// 	pubsub.Subscribe("gaming", sub2)
// 	pubsub.Subscribe("gaming", sub3)

// 	fmt.Println("✅ 3 subscribers subscribed to 'gaming'\n")

// 	// ============================================
// 	// REQUIREMENT 2: Publisher publishes to topic
// 	// REQUIREMENT 4: All subscribers receive message
// 	// REQUIREMENT 6: Subscribers run in parallel
// 	// ============================================

// 	fmt.Println("--- Publishing Messages ---\n")

// 	pubsub.Publish("gaming", "New video: Elden Ring Gameplay")
// 	time.Sleep(100 * time.Millisecond)

// 	pubsub.Publish("gaming", "New video: God of War Review")
// 	time.Sleep(100 * time.Millisecond)

// 	pubsub.Publish("gaming", "New video: Cyberpunk 2077")
// 	time.Sleep(100 * time.Millisecond)

// 	fmt.Println("\n✅ All 3 subscribers received all 3 messages (fan-out)\n")

// 	// ============================================
// 	// REQUIREMENT 5: Reset offset (replay messages)
// 	// ============================================

// 	fmt.Println("--- Testing Offset Reset (Replay) ---\n")

// 	// Reset email-service to offset 0 (replay ALL messages)
// 	fmt.Println("🔄 Resetting email-service to offset 0...\n")

// 	err := pubsub.ResetOffset("gaming", "email-service", 0)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}

// 	time.Sleep(200 * time.Millisecond)

// 	fmt.Println("\n✅ Email-service replayed all 3 messages!\n")

// 	// Reset analytics-service to offset 1 (replay from second message)
// 	fmt.Println("🔄 Resetting analytics-service to offset 1...\n")

// 	err = pubsub.ResetOffset("gaming", "analytics-service", 1)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}

// 	time.Sleep(200 * time.Millisecond)

// 	fmt.Println("\n✅ Analytics-service replayed messages from offset 1 onwards!\n")

// 	// ============================================
// 	// Additional Demo: Publish more after reset
// 	// ============================================

// 	fmt.Println("--- Publishing New Message ---\n")

// 	pubsub.Publish("gaming", "New video: Starfield Exploration")
// 	time.Sleep(100 * time.Millisecond)

// 	fmt.Println("\n✅ All subscribers received new message\n")

// 	// ============================================
// 	// Cleanup
// 	// ============================================

// 	fmt.Println("--- Unsubscribing ---\n")
// 	pubsub.Unsubscribe("gaming", "notification-service")

// 	fmt.Println("\n--- Done ---")
// }
