package main

// import (
// 	"errors"
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type Message struct {
// 	Topic string
// 	Data  interface{}
// }

// type Topic struct {
// 	name        string
// 	subscribers []*Subscriber
// 	mu          sync.RWMutex
// }

// type PubSub struct {
// 	topics map[string]*Topic //[topic_name -> topic_object]
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
// 	}
// 	ps.topics[topicName] = topic
// 	return nil
// }

// // Benefits:
// // Don't hold PubSub lock while modifying topic
// // Better concurrency (other threads can access other topics)
// // Use RLock when only reading (multiple readers allowed
// func (ps *PubSub) Subscribe(topicName string, subscriber *Subscriber) error {
// 	//Lock Pubsub to read topics map
// 	ps.mu.RLock() //Use RLock (only reading)
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock() //Unlock early

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	//Lock Topic to modify subscribers
// 	topic.mu.Lock() //Lock the topic
// 	topic.subscribers = append(topic.subscribers, subscriber)
// 	topic.mu.Unlock() //Unlock the topic

// 	//Start subscriber
// 	subscriber.start() //Start listening
// 	return nil
// }

// func (ps *PubSub) Unsubscribe(topicName string, subscriberId string) error {
// 	// Remove subscriber from topic

// 	//Lock Pubsub to read topics map
// 	ps.mu.RLock() //Use RLock (only reading)
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock() //Unlock early

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	//Lock Topic to modify subscribers
// 	topic.mu.Lock() //Lock the topic
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
// 	topic.mu.Unlock() //Unlock the topic

// 	//close the removed subscriber's channel
// 	if removedSubscriber != nil {
// 		removedSubscriber.Close()
// 	} else {
// 		return errors.New("Susbcriber not found")
// 	}

// 	return nil
// }

// func (ps *PubSub) Publish(topicName string, data interface{}) error {
// 	//publish to all subscribers of this topic
// 	//getTopic from topic name
// 	ps.mu.RLock()
// 	topic, ok := ps.topics[topicName]
// 	ps.mu.RUnlock()

// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}

// 	topic.mu.RLock()
// 	//publish msg to all subscribers
// 	subscribers := topic.subscribers
// 	topic.mu.RUnlock()

// 	for _, subscriber := range subscribers {
// 		subscriber.channel <- Message{Topic: topicName, Data: data}
// 	}
// 	return nil
// }

// type Subscriber struct {
// 	id      string
// 	channel chan Message
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
// 	//start single goroutine that listens on the channel
// 	go func() {
// 		for msg := range s.channel {
// 			s.handler(msg)
// 		}
// 	}()
// }

// func (s *Subscriber) Close() {
// 	//stop the channel so nothing gets read from this
// 	close(s.channel)
// }

// func main() {
// 	//Create pubsub system
// 	pubsub := NewPubSub()

// 	//create topic
// 	pubsub.CreateTopic("gaming")

// 	//create subscribers
// 	sub1 := NewSubscriber("notification-service", func(msg Message) {
// 		fmt.Printf("[Notification] Received: %v\n", msg.Data)
// 	})

// 	sub2 := NewSubscriber("analytics-service", func(msg Message) {
// 		fmt.Printf("[Analytics] Received: %v\n", msg.Data)
// 	})

// 	sub3 := NewSubscriber("email-service", func(msg Message) {
// 		fmt.Printf("[Email] Received: %v\n", msg.Data)
// 	})

// 	// Subscribe them
// 	pubsub.Subscribe("gaming", sub1)
// 	pubsub.Subscribe("gaming", sub2)
// 	pubsub.Subscribe("gaming", sub3)

// 	// Publish messages
// 	fmt.Println("\n--- Publishing Message 1 ---")
// 	pubsub.Publish("gaming", "New video: Elden Ring Gameplay")
// 	time.Sleep(100 * time.Millisecond) // Let them process

// 	fmt.Println("\n--- Publishing Message 2 ---")
// 	pubsub.Publish("gaming", "New video: God of War Review")

// 	time.Sleep(100 * time.Millisecond)

// 	fmt.Println("\n--- Unsubscribing notification-service ---")
// 	pubsub.Unsubscribe("gaming", "notification-service")

// 	// Publish again
// 	fmt.Println("\n--- Publishing Message 3 ---")
// 	pubsub.Publish("gaming", "New video: Cyberpunk 2077")

// 	time.Sleep(100 * time.Millisecond)

// 	fmt.Println("\n--- Done ---")
// }
