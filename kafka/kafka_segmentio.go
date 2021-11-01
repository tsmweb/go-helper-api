package kafka

import (
	"context"
	skafka "github.com/segmentio/kafka-go"
	"log"
	"time"
)

// kafka implementation for Kafka interface with segmentio/kafka-go library.
type kafka struct {
	kafkaBrokerUrls []string
	clientId        string
	debug           bool
}

// NewKafka creates a new instance of Kafka.
func New(kafkaBrokerUrls []string, clientId string) Kafka {
	return &kafka{
		kafkaBrokerUrls: kafkaBrokerUrls,
		clientId:        clientId,
		debug:           false,
	}
}

// Debug enables logging of incoming events.
func (k *kafka) Debug(debug bool) {
	k.debug = debug
}

func (k *kafka) dialer() *skafka.Dialer {
	return &skafka.Dialer{
		ClientID:  k.clientId,
		DualStack: true,
		Timeout:   10 * time.Second,
	}
}

// NewProducer creates a new Producer to produce events on a topic.
func (k *kafka) NewProducer(topic string) Producer {
	w := &skafka.Writer{
		Addr:         skafka.TCP(k.kafkaBrokerUrls...),
		Topic:        topic,
		RequiredAcks: skafka.RequireOne,
		BatchTimeout: time.Millisecond,
		Compression:  skafka.Snappy,
	}
	return newProducer(w)
}

// NewConsumer creates a new Consumer to consume events from a topic.
func (k *kafka) NewConsumer(groupID, topic string) Consumer {
	config := skafka.ReaderConfig{
		Brokers:         k.kafkaBrokerUrls,
		Topic:           topic,
		Dialer:          k.dialer(),
		MinBytes:        5,
		MaxBytes:        10e6,        // 10MB
		MaxWait:         time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
		//CommitInterval: time.Second, // flushes commits to KafkaWrap every second
	}
	if groupID != "-1" {
		config.GroupID = groupID
	}

	r := skafka.NewReader(config)
	return newConsumer(r, k.debug)
}

// producer provide methods for producing events for a given topic.
type producer struct {
	writer *skafka.Writer
}

func newProducer(w *skafka.Writer) Producer {
	return &producer{writer: w}
}

// Publish produces and sends an event for a kafka topic.
// The context passed as first argument may also be used to asynchronously
// cancel the operation.
func (p *producer) Publish(ctx context.Context, key []byte, values ...[]byte) error {
	var messages []skafka.Message

	for _, value := range values {
		message := skafka.Message{
			Key: key,
			Value: value,
		}
		messages = append(messages, message)
	}

	return p.writer.WriteMessages(ctx, messages...)
}

// Close flushes pending writes, and waits for all writes to complete before
// returning. Calling Close also prevents new writes from being submitted to
// the Writer, further calls to WriteMessages and the like will fail with
// io.ErrClosedPipe.
func (p *producer) Close() {
	p.writer.Close()
}

// Consumer provide methods for consuming events on a given topic.
type consumer struct {
	reader *skafka.Reader
	debug bool
}

func newConsumer(r *skafka.Reader, d bool) Consumer {
	return &consumer{
		reader: r,
		debug: d,
	}
}

// Subscribe consumes the events of a topic and passes the event to the
// informed callback function. The method call blocks until an error occurs.
// The program may also specify a context to asynchronously cancel the blocking operation.
func (c *consumer) Subscribe(ctx context.Context, callbackFn func(event *Event, err error)) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			callbackFn(nil, err)
			break
		}

		if c.debug {
			c.logMessage(m)
		}

		callbackFn(c.makeEvent(m), nil)
	}
}

// Close closes the stream, preventing the program from reading any more
// events from it.
func (c *consumer) Close() {
	c.reader.Close()
}

func (c *consumer) makeEvent(m skafka.Message) *Event {
	headerMap := c.sliceToMap(m.Headers)
	return &Event{
		Topic:  m.Topic,
		Key:    m.Key,
		Value:  m.Value,
		Header: headerMap,
		Time:   m.Time,
	}
}

func (c *consumer) sliceToMap(sh []skafka.Header) map[string]string {
	headerMap := make(map[string]string)

	for _, v := range sh {
		headerMap[v.Key] = string(v.Value)
	}

	return headerMap
}

func (c *consumer) logMessage(m skafka.Message) {
	log.Printf("[>] ReadMessage - TIME: %s | TOPIC: %s | PARTITION: %d | OFFSET: %d | HEADER: %s | SIZE PAYLOAD: %d\n",
		m.Time, m.Topic, m.Partition, m.Offset, m.Headers, len(m.Value))
}
