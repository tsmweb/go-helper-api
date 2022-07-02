/*
Package event implements routines to produce event log for a topic in Apache Kafka.

Send an event log to a kafka topic:

	err := event.Init(context.Background(), []string{"localhost:9094"}, "CLIENT_ID", "TOPIC_NAME")
	if err != nil {
	// ...
	defer event.Close()

	e := event.New("localhost", "user", "New User", "New user added to the database")

	if err = event.Send(e); err != nil {
	// ...
*/

package event

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tsmweb/go-helper-api/kafka"
	"log"
	"sync"
	"time"
)

// Event structure represents an event log.
type Event struct {
	Host      string `json:"host"`
	User      string `json:"user,omitempty"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Timestamp string `json:"timestamp"`
}

// New creates an Event instance.
func New(host, user, title, detail string) *Event {
	return &Event{
		Host:      host,
		User:      user,
		Title:     title,
		Detail:    detail,
		Timestamp: time.Now().Format("2006-02-01 15:04:05"),
	}
}

func (e Event) toJSON() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return b
}

var (
	chEvent chan *Event
	wg      sync.WaitGroup
	running bool
	mu      sync.RWMutex // guard running

	ErrClosed  = errors.New("closed eventlog")
	ErrRunning = errors.New("is already running")
)

// Init creates a new producer for Apache Kafka and initializes the routines for sending events.
func Init(ctx context.Context, brokerUrls []string, clientID string, topic string) error {
	mu.RLock()
	if running {
		mu.RUnlock()
		return ErrRunning
	}
	mu.RUnlock()

	chEvent = make(chan *Event)
	producer := kafka.New(brokerUrls, clientID).NewProducer(topic)
	running = true
	wg.Add(1)

	go func() {
		defer wg.Done()

		for event := range chEvent {
			if err := producer.Publish(ctx, []byte(event.Host), event.toJSON()); err != nil {
				log.Printf("eventlog.Send() \nError: %v\n", err.Error())
			}
		}

		producer.Close()
	}()

	return nil
}

// Close closes the event communication channel and the Apache Kafka producer.
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if running {
		close(chEvent)
		running = false
		wg.Wait()
	}
}

// Send sends the event to the Apache Kafka topic.
func Send(event *Event) error {
	mu.RLock()
	defer mu.RUnlock()

	if !running {
		return ErrClosed
	}

	chEvent <- event
	return nil
}
