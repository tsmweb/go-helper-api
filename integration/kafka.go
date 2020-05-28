/*
Package integration provides a custom Kafka wrapper to work with the outbox pattern.

Send a event to a kafka topic:

	profile := profile{
		ID:       "123",
		Name:     "Paul",
		Lastname: "Mark",
	}

	// newClientCreateEvent return *profileCreateEvent which implements the ExportedEvent interface.
	profileCreateEvent := newClientCreateEvent(profile)
	event := integration.BuildEventMessage(profileCreateEvent)

	kafka := integration.NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "profile")
	kafka.Debug(true)
	writer := kafka.NewWriter("CLIENT")

	defer writer.Close()

	err := kafka.SendEvent(context.Background(), writer, event)
	if err != nil {
	// ...

Consume event from a kafka topic:

	kafka := integration.NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "profile")
	kafka.Debug(true)
	reader := kafka.NewReader("ClientSubscribeTest", "CLIENT")

	defer reader.Close()

	ctx := context.Background()

	fnCallback := func(event EventMessage, err error) {
		// ...
	}

	kafka.SubscribeEvent(ctx, reader, fnCallback)
	// ...
 */
package integration

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"log"
	"time"
)

// Kafka is a custom Kafka wrapper to work with the outbox pattern.
type Kafka struct {
	kafkaBrokerUrls []string
	clientId        string
	debug           bool
}

// NewKafka creates a new Kafka instance taking as an parameter an array of
// kafka brokers and the profile ID.
func NewKafka(kafkaBrokerUrls []string, clientId string) *Kafka {
	return &Kafka{
		kafkaBrokerUrls: kafkaBrokerUrls,
		clientId:        clientId,
		debug:           false,
	}
}

// Debug enables logging of incoming messages.
func (k *Kafka) Debug(debug bool) {
	k.debug = debug
}

func (k *Kafka) dialer() *kafka.Dialer {
	return &kafka.Dialer{
		ClientID: k.clientId,
		Timeout:  10 * time.Second,
	}
}

// NewWriter creates a new kafka.Writer to write messages on a topic.
func (k *Kafka) NewWriter(topic string) *kafka.Writer {
	config := kafka.WriterConfig{
		Brokers:          k.kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.Hash{},
		Dialer:           k.dialer(),
		RequiredAcks:     -1,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	w := kafka.NewWriter(config)
	return w
}

// SendEvent produces and sends an event message for a kafka topic.
// The context passed as first argument may also be used to asynchronously
// cancel the operation.
func (k *Kafka) SendEvent(ctx context.Context, writer *kafka.Writer, event EventMessage) error {
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

	err := writer.WriteMessages(ctx, message)
	return err
}

// NewReader creates a new kafka.Reader to read messages for a topic.
func (k *Kafka) NewReader(groupID, topic string) *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:         k.kafkaBrokerUrls,
		GroupID:         groupID,
		Topic:           topic,
		Dialer:          k.dialer(),
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	return reader
}

// SubscribeEvent consumes the events of a topic and passes the events to the
// informed callback function. The method call blocks until an error occurs.
// The program may also specify a context to asynchronously cancel the blocking operation.
func (k *Kafka) SubscribeEvent(ctx context.Context, reader *kafka.Reader, callback func(event EventMessage, err error)) {
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			callback(EventMessage{}, err)
			break
		}

		if k.debug {
			k.logMessage(m)
		}

		callback(k.makeEventMessage(m), nil)
	}
}

func (k *Kafka) makeEventMessage(m kafka.Message) EventMessage {
	headerMap := k.sliceToMap(m.Headers)

	id, ok := headerMap["id"]
	if !ok {
		id = "0"
	}

	eventType, ok := headerMap["eventType"]
	if !ok {
		eventType = "*"
	}

	event := EventMessage{
		ID:            id,
		AggregateID:   string(m.Key),
		AggregateType: m.Topic,
		Type:          eventType,
		Payload:       m.Value,
	}

	return event
}

func (k *Kafka) sliceToMap(sh []kafka.Header) map[string]string {
	headerMap := make(map[string]string)

	for _, v := range sh {
		headerMap[v.Key] = string(v.Value)
	}

	return headerMap
}

func (k *Kafka) logMessage(m kafka.Message) {
	log.Printf("[>] ReadMessage - TIME: %s | TOPIC: %s | PARTITION: %d | OFFSET: %d | HEADER: %s | SIZE PAYLOAD: %d\n",
		m.Time, m.Topic, m.Partition, m.Offset, m.Headers, len(m.Value))
}
