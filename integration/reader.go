package integration

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
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

// Reader provide methods for consuming events on a given topic.
type Reader struct {
	reader *kafka.Reader
	debug  bool
}

func newReader(reader *kafka.Reader, debug bool) *Reader {
	return &Reader{
		reader,
		debug,
	}
}

// SubscribeTopic consumes the events of a topic and passes the event to the
// informed callback function. The method call blocks until an error occurs.
// The program may also specify a context to asynchronously cancel the blocking operation.
func (r *Reader) SubscribeTopic(ctx context.Context, callback func(event *Event, err error)) {
	for {
		m, err := r.reader.ReadMessage(ctx)
		if err != nil {
			callback(nil, err)
			break
		}

		if r.debug {
			r.logMessage(m)
		}

		callback(r.makeEvent(m), nil)
	}
}

// SubscribeOutboxEvent consumes the events of a topic and passes the events to the
// informed callback function. The method call blocks until an error occurs.
// The program may also specify a context to asynchronously cancel the blocking operation.
func (r *Reader) SubscribeOutboxEvent(ctx context.Context, callback func(event *OutboxEvent, err error)) {
	for {
		m, err := r.reader.ReadMessage(ctx)
		if err != nil {
			callback(nil, err)
			break
		}

		if r.debug {
			r.logMessage(m)
		}

		callback(r.makeOutboxEvent(m), nil)
	}
}

// Close closes the stream, preventing the program from reading any more
// events from it.
func (r *Reader) Close() {
	r.reader.Close()
}

func (r *Reader) makeEvent(m kafka.Message) *Event {
	headerMap := r.sliceToMap(m.Headers)
	return &Event{
		Topic:  m.Topic,
		Key:    m.Key,
		Value:  m.Value,
		Header: headerMap,
		Time:   m.Time,
	}
}

func (r *Reader) makeOutboxEvent(m kafka.Message) *OutboxEvent {
	headerMap := r.sliceToMap(m.Headers)

	id, ok := headerMap["id"]
	if !ok {
		id = "0"
	}

	eventType, ok := headerMap["eventType"]
	if !ok {
		eventType = "*"
	}

	event := &OutboxEvent{
		ID:            id,
		AggregateID:   string(m.Key),
		AggregateType: m.Topic,
		Type:          eventType,
		Payload:       m.Value,
	}
	return event
}

func (r *Reader) sliceToMap(sh []kafka.Header) map[string]string {
	headerMap := make(map[string]string)

	for _, v := range sh {
		headerMap[v.Key] = string(v.Value)
	}

	return headerMap
}

func (r *Reader) logMessage(m kafka.Message) {
	log.Printf("[>] ReadMessage - TIME: %s | TOPIC: %s | PARTITION: %d | OFFSET: %d | HEADER: %s | SIZE PAYLOAD: %d\n",
		m.Time, m.Topic, m.Partition, m.Offset, m.Headers, len(m.Value))
}
