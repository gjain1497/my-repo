package kafka

// Layer 2: KAFKA CLIENT LIBRARY


// ============================================
// FILE: kafkaclient/consumer.go
// Package: kafkaclient
// ============================================

package kafkaclient

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Message struct {
	Offset int64
	Key    string
	Value  string
}

type Request struct {
	Type       string
	Topic      string
	Partition  int
	Offset     int64
	Key        string
	Value      string
	GroupID    string
	ConsumerID string
}

type Response struct {
	Type      string
	Messages  []Message
	Partition int
	Error     string
}

// ============================================
// CONSUMER
// ============================================

type Consumer struct {
	brokerAddr string
	groupID    string
	consumerID string
	topic      string
	partition  int    // Which partition assigned
	offset     int64  // Current read position
	conn       net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	mu         sync.Mutex
}

func NewConsumer(brokerAddr, groupID, consumerID string) (*Consumer, error) {
	return &Consumer{
		brokerAddr: brokerAddr,
		groupID:    groupID,
		consumerID: consumerID,
	}, nil
}

func (c *Consumer) Subscribe(topic string, partition int) error {
	c.topic = topic
	c.partition = partition

	// Connect to broker
	conn, err := net.Dial("tcp", c.brokerAddr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)

	// Get last committed offset for this consumer group
	req := Request{
		Type:      "GET_OFFSET",
		GroupID:   c.groupID,
		Partition: partition,
	}

	c.mu.Lock()
	c.encoder.Encode(req)
	var resp Response
	c.decoder.Decode(&resp)
	c.mu.Unlock()

	c.offset = resp.Messages[0].Offset

	fmt.Printf("[Consumer %s] Subscribed to %s-partition-%d, starting at offset %d\n",
		c.consumerID, topic, partition, c.offset)

	return nil
}

// ============================================
// POLL - Pull messages from broker
// ============================================

func (c *Consumer) Poll(timeoutMs int) (*Message, error) {
	c.mu.Lock()
	currentOffset := c.offset
	c.mu.Unlock()

	// Send FETCH request
	req := Request{
		Type:      "FETCH",
		Topic:     c.topic,
		Partition: c.partition,
		Offset:    currentOffset,
	}

	c.mu.Lock()
	err := c.encoder.Encode(req)
	if err != nil {
		c.mu.Unlock()
		return nil, err
	}

	// Wait for response
	var resp Response
	err = c.decoder.Decode(&resp)
	c.mu.Unlock()

	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	// Check if we got messages
	if len(resp.Messages) == 0 {
		return nil, nil // No new messages
	}

	// Take first message
	msg := resp.Messages[0]

	// Update offset
	c.mu.Lock()
	c.offset = msg.Offset + 1
	c.mu.Unlock()

	return &msg, nil
}

// ============================================
// COMMIT - Save progress
// ============================================

func (c *Consumer) Commit() error {
	c.mu.Lock()
	currentOffset := c.offset
	c.mu.Unlock()

	req := Request{
		Type:      "COMMIT",
		GroupID:   c.groupID,
		Partition: c.partition,
		Offset:    currentOffset,
	}

	c.mu.Lock()
	err := c.encoder.Encode(req)
	if err != nil {
		c.mu.Unlock()
		return err
	}

	var resp Response
	err = c.decoder.Decode(&resp)
	c.mu.Unlock()

	if err != nil {
		return err
	}

	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}

	return nil
}

func (c *Consumer) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
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
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, err
	}

	return &Producer{
		brokerAddr: brokerAddr,
		conn:       conn,
		encoder:    json.NewEncoder(conn),
		decoder:    json.NewDecoder(conn),
	}, nil
}

func (p *Producer) Produce(topic, key, value string) error {
	req := Request{
		Type:  "PRODUCE",
		Topic: topic,
		Key:   key,
		Value: value,
	}

	p.mu.Lock()
	err := p.encoder.Encode(req)
	if err != nil {
		p.mu.Unlock()
		return err
	}

	var resp Response
	err = p.decoder.Decode(&resp)
	p.mu.Unlock()

	if err != nil {
		return err
	}

	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}

	return nil
}

func (p *Producer) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}