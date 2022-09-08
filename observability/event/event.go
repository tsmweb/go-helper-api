/*
Package event implements routines to produce event log for a topic in Apache Kafka.

Send an event log to a kafka topic:

	producer := kafka.New([]string{"localhost:9094"}, "CLIENT_ID").NewProducer("TOPIC_NAME")
	err := event.Init(producer)
	if err != nil {
	// ...
	defer event.Close()

	e := event.New("localhost", "user", "New User", event.Info, "New user added to the database")

	if err = event.Send(e); err != nil {
	// ...
*/

package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tsmweb/go-helper-api/kafka"
)

// EventType represents the event type ("info", "debug", "warning", ...).
type EventType int

const (
	// Info represents an information event.
	Info EventType = iota

	// Debug represents a debug event.
	Debug

	// Warning represents a warning event.
	Warning

	// Error represents an error event.
	Error
)

var eventTypeText = map[EventType]string{
	Info:    "info",
	Debug:   "debug",
	Warning: "warning",
	Error:   "error",
}

// String return the name of the eventType.
func (e EventType) String() string {
	return eventTypeText[e]
}

// Event structure represents an event log.
type Event struct {
	Host      string `json:"host"`
	User      string `json:"user,omitempty"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Detail    string `json:"detail"`
	Timestamp string `json:"timestamp"`
}

// New creates an Event instance.
func New(host string, user string, title string, eventType EventType, detail string) *Event {
	return &Event{
		Host:      host,
		User:      user,
		Title:     title,
		Type:      eventType.String(),
		Detail:    detail,
		Timestamp: time.Now().Format("2006-01-02T15:04:05-0700"), // yyyy-MM-dd'T'HH:mm:ssZ
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
func Init(producer kafka.Producer) error {
	mu.RLock()
	if running {
		mu.RUnlock()
		return ErrRunning
	}
	mu.RUnlock()

	ctx := context.Background()
	chEvent = make(chan *Event)
	running = true
	wg.Add(1)

	go func() {
		defer wg.Done()

		for event := range chEvent {
			if err := producer.Publish(ctx, []byte(event.Host), event.toJSON()); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
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
