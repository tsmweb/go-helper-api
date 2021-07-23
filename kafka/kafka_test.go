package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestProducer_Publish(t *testing.T) {
	kafka := New([]string{"localhost:9094"}, "TEST")
	producer := kafka.NewProducer("USERS")
	defer producer.Close()

	for i := 0; i < 10; i++ {
		u := user{
			ID:   uuid.New().String(),
			Name: fmt.Sprintf("Test-%d", i),
			Age:  i + 10,
		}

		if err := producer.Publish(context.Background(), []byte(u.ID), u.toJSON()); err != nil {
			t.Fatalf("producer.Publish() - Error: %v", err)
		}
	}

	t.Log("published messages")
}

func TestConsumer_Subscribe(t *testing.T) {
	kafka := New([]string{"localhost:9094"}, "TEST")
	kafka.Debug(true)
	consumer := kafka.NewConsumer("UserSubscribeTest", "USERS")
	defer consumer.Close()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	callbackFn := func(event *Event, err error) {
		if err != nil {
			if strings.Contains(err.Error(), "deadline exceeded") {
				t.Logf("Subscribe() %s", err)
			} else {
				t.Errorf("Subscribe() - Error: %s", err)
			}
		} else {
			t.Log("--------------------------------------------------------")
			t.Logf("[>] TOPIC: %s", event.Topic)
			t.Logf("[>] KEY: %s", string(event.Key))
			t.Logf("[>] Value: %s", string(event.Value))
			t.Logf("[>] Time: %v", event.Time)
		}
	}

	consumer.Subscribe(ctx, callbackFn)
}

type user struct {
	ID   string
	Name string
	Age  int
}

func (u *user) toJSON() []byte {
	us, err := json.Marshal(u)
	if err != nil {
		return nil
	}
	return us
}
