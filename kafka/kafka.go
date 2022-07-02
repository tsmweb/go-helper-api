/*
Package kafka is a simple wrapper for the kafka-go segmentio library, providing tools
for consuming and producing events.

Send an event to a kafka topic:

	p := profile{
		ID:       "123",
		Name:     "Paul",
		Lastname: "Mark",
	}
	pj, _ := json.Marshal(profile)

	k := kafka.New([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "profile")
	k.Debug(true)

	producer := k.NewProducer("CLIENT")
	defer producer.Close()

	err := producer.Publish(context.Background(), []byte(profile.ID), pj)
	if err != nil {
	// ...

Consume event from a kafka topic:

	k := kafka.New([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "profile")
	k.Debug(true)
	consumer := k.NewConsumer("ClientSubscribeTest", "CLIENT")
	defer consumer.Close()

	ctx := context.Background()

	callbackFn := func(event *kafka.Event, err error) {
		// ...
	}

	consumer.Subscribe(ctx, callbackFn)
	// ...
*/

package kafka

import (
	"context"
	"time"
)

// Event structure that represents a Kafka event.
type Event struct {
	Topic  string
	Key    []byte
	Value  []byte
	Header map[string]string
	Time   time.Time
}

// Kafka it is an abstraction to send and consume events from an
// event kafka service.
type Kafka interface {

	// NewProducer creates a new Producer to produce events on a topic.
	NewProducer(topic string) Producer

	// NewConsumer creates a new Consumer to consume events from a topic.
	NewConsumer(groupID, topic string) Consumer

	// Debug enables logging of incoming events.
	Debug(debug bool)
}

// Producer provide methods for producing events for a given topic.
type Producer interface {
	// Publish produces and sends an event for a kafka topic.
	// The context passed as first argument may also be used to asynchronously
	// cancel the operation.
	Publish(ctx context.Context, key []byte, values ...[]byte) error

	// Close flushes pending writes, and waits for all writes to complete before
	// returning.
	Close()
}

// Consumer provide methods for consuming events on a given topic.
type Consumer interface {
	// Subscribe consumes the events of a topic and passes the event to the
	// informed callback function. The method call blocks until an error occurs.
	// The program may also specify a context to asynchronously cancel the blocking operation.
	Subscribe(ctx context.Context, callbackFn func(event *Event, err error))

	// Close closes the stream, preventing the program from reading any more
	// events from it.
	Close()
}
