/*
Package integration is a simple wrapper for the kafka-go segmentio library, providing tools
for consuming and producing events and working with the event outbox pattern.

Send a outbox event to a kafka topic:

	profile := profile{
		ID:       "123",
		Name:     "Paul",
		Lastname: "Mark",
	}

	// newClientCreateEvent return *profileCreateEvent which implements the ExportedEvent interface.
	profileCreateEvent := newClientCreateEvent(profile)
	event := integration.BuildOutboxEvent(profileCreateEvent)

	kafka := integration.NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "profile")
	kafka.Debug(true)
	writer := kafka.NewWriter("CLIENT")

	defer writer.Close()

	err := writer.SendOutboxEvent(context.Background(), event)
	if err != nil {
	// ...

Consume outbox event from a kafka topic:

	kafka := integration.NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "profile")
	kafka.Debug(true)
	reader := kafka.NewReader("ClientSubscribeTest", "CLIENT")

	defer reader.Close()

	ctx := context.Background()

	fnCallback := func(event *OutboxEvent, err error) {
		// ...
	}

	reader.SubscribeOutboxEvent(ctx, fnCallback)
	// ...
*/
package integration

import (
	"github.com/segmentio/kafka-go"
	"time"
)

// Kafka is a simple wrapper for the kafka-go segmentio library, providing tools
// for consuming and producing events and working with the event outbox pattern.
type Kafka struct {
	kafkaBrokerUrls []string
	clientId        string
	debug           bool
}

// NewKafka creates a new Kafka instance taking as an parameter an array of
// kafka brokers and the client ID.
func NewKafka(kafkaBrokerUrls []string, clientId string) *Kafka {
	return &Kafka{
		kafkaBrokerUrls: kafkaBrokerUrls,
		clientId:        clientId,
		debug:           false,
	}
}

// Debug enables logging of incoming events.
func (k *Kafka) Debug(debug bool) {
	k.debug = debug
}

func (k *Kafka) dialer() *kafka.Dialer {
	return &kafka.Dialer{
		ClientID:  k.clientId,
		DualStack: true,
		Timeout:   10 * time.Second,
	}
}

// NewWriter creates a new Writer to produce events on a topic.
func (k *Kafka) NewWriter(topic string) *Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(k.kafkaBrokerUrls...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: time.Millisecond,
		Compression:  kafka.Snappy,
	}
	return newWriter(w)
}

// NewReader creates a new Reader to consume events from a topic.
func (k *Kafka) NewReader(groupID, topic string) *Reader {
	config := kafka.ReaderConfig{
		Brokers:         k.kafkaBrokerUrls,
		GroupID:         groupID,
		Topic:           topic,
		Dialer:          k.dialer(),
		MinBytes:        10e3,        // 10KB
		MaxBytes:        10e6,        // 10MB
		MaxWait:         time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
		//CommitInterval: time.Second, // flushes commits to KafkaWrap every second
	}
	r := kafka.NewReader(config)
	return newReader(r, k.debug)
}