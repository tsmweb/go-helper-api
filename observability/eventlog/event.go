package eventlog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsmweb/go-helper-api/kafka"
	"log"
	"sync"
	"time"
)

type Event struct {
	Host      string `json:"host"`
	User      string `json:"user,omitempty"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Timestamp string `json:"timestamp"`
}

func NewEvent(host, user, title, detail string) *Event {
	return &Event{
		Host:      host,
		User:      user,
		Title:     title,
		Detail:    detail,
		Timestamp: time.Now().Format("2006-02-01 15:04:05"),
	}
}

func (e Event) Print() {
	fmt.Printf("[%v] TITLE [%s] | HOST [%s] | USER [%s] | DETAIL [%s]\n",
		e.Timestamp, e.Title, e.Host, e.User, e.Detail)
}

func (e Event) ToJSON() []byte {
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
			if err := producer.Publish(ctx, []byte(event.Host), event.ToJSON()); err != nil {
				log.Printf("eventlog.Send() \nError: %v\n", err.Error())
			}
		}

		producer.Close()
	}()

	return nil
}

func Close() {
	mu.Lock()
	defer mu.Unlock()

	if running {
		close(chEvent)
		running = false
		wg.Wait()
	}
}

func Send(event *Event) error {
	mu.RLock()
	defer mu.RUnlock()

	if !running {
		return ErrClosed
	}

	chEvent <- event
	return nil
}
