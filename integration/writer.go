package integration

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

// Writer provide methods for producing events for a given topic.
type Writer struct {
	writer *kafka.Writer
}

func newWriter(writer *kafka.Writer) *Writer {
	return &Writer{writer}
}

// SendEvent produces and sends an event for a kafka topic.
// The context passed as first argument may also be used to asynchronously
// cancel the operation.
func (w *Writer) SendEvent(ctx context.Context, key, value []byte) error {
	message := kafka.Message{
		Key: key,
		Value: value,
	}
	return w.writer.WriteMessages(ctx, message)
}

// SendOutboxEvent produces and sends an event message for a kafka topic.
// The context passed as first argument may also be used to asynchronously
// cancel the operation.
func (w *Writer) SendOutboxEvent(ctx context.Context, event *OutboxEvent) error {
	headers := []kafka.Header{
		{Key: "id", Value: []byte(event.ID)},
		{Key: "eventType", Value: []byte(event.Type)},
	}

	message := kafka.Message{
		Key:     []byte(event.AggregateID),
		Value:   event.Payload,
		Headers: headers,
		Time:    time.Time{},
	}
	return w.writer.WriteMessages(ctx, message)
}

// Close flushes pending writes, and waits for all writes to complete before
// returning. Calling Close also prevents new writes from being submitted to
// the Writer, further calls to WriteMessages and the like will fail with
// io.ErrClosedPipe.
func (w *Writer) Close() {
	w.writer.Close()
}
