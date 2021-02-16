package integration

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestWriter_SendMessage(t *testing.T) {
	kafka := NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "Profile")
	kafka.Debug(true)
	writer := kafka.NewWriter("profile")
	defer writer.Close()

	for i := 0; i < 10; i++ {
		p := profile{
			ID:       uuid.New().String(),
			Name:     "Paul",
			Lastname: "Mark",
		}

		err := writer.SendEvent(context.Background(), []byte(p.ID), p.ToJSON())
		if err != nil {
			t.Errorf("SendEvent error: %s", err)
		}
	}

	t.Log("sent messages")
}

func TestReader_SubscribeTopic(t *testing.T) {
	kafka := NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "Profile")
	kafka.Debug(true)
	reader := kafka.NewReader("ProfileSubscribeTest", "profile")
	defer reader.Close()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	fnCallback := func(event *Event, err error) {
		if err != nil {
			if strings.Contains(err.Error(), "deadline exceeded") {
				t.Logf("SubscribeTopic %s", err)
			}  else {
				t.Errorf("SubscribeTopic error: %s", err)
			}
		} else {
			t.Log("--------------------------------------------------------")
			t.Logf("[>] KEY: %s", string(event.Key))
			t.Logf("[>] Value: %s", string(event.Value))
			t.Logf("[>] Time: %v", event.Time)
		}
	}

	reader.SubscribeTopic(ctx, fnCallback)
}

func TestWriter_SendEvent(t *testing.T) {
	kafka := NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "Profile")
	kafka.Debug(true)
	writer := kafka.NewWriter("profile_event")
	defer writer.Close()

	for i := 0; i < 10; i++ {
		p := profile{
			ID:       uuid.New().String(),
			Name:     "Paul",
			Lastname: "Mark",
		}

		pce := newProfileCreateEvent(p)
		event := BuildOutboxEvent(pce)

		err := writer.SendOutboxEvent(context.Background(), event)
		if err != nil {
			t.Errorf("SendOutboxEvent event %v, error: %s", event, err)
		}
	}

	t.Log("sent messages")
}

func TestReader_SubscribeEvent(t *testing.T) {
	kafka := NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "Profile")
	kafka.Debug(true)
	reader := kafka.NewReader("ProfileSubscribeTest", "profile_event")
	defer reader.Close()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//ctx := context.Background()

	fnCallback := func(event *OutboxEvent, err error) {
		if err != nil {
			if strings.Contains(err.Error(), "deadline exceeded") {
				t.Logf("SubscribeOutboxEvent %s", err)
			}  else {
				t.Errorf("SubscribeOutboxEvent error: %s", err)
			}
		} else {
			t.Log("--------------------------------------------------------")
			t.Logf("[>] ID: %s", event.ID)
			t.Logf("[>] AggregateID: %s", event.AggregateID)
			t.Logf("[>] AggregateType: %s", event.AggregateType)
			t.Logf("[>] Type: %s", event.Type)
			t.Logf("[>] Payload: %s", event.Payload)
		}
	}

	reader.SubscribeOutboxEvent(ctx, fnCallback)
}

type profile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (p *profile) ToJSON() []byte {
	pj, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return pj
}

type profileCreateEvent struct {
	prof profile
}

func newProfileCreateEvent(prof profile) *profileCreateEvent {
	return &profileCreateEvent{prof}
}

func (p *profileCreateEvent) AggregateID() string {
	return p.prof.ID
}

func (p *profileCreateEvent) AggregateType() string {
	return "profile"
}

func (p *profileCreateEvent) Type() string {
	return "ProfileCreated"
}

func (p *profileCreateEvent) Payload() []byte {
	b, err := json.Marshal(p.prof)
	if err != nil {
		return nil
	}
	return b
}