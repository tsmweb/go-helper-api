package integration

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestKafka_SendEvent(t *testing.T) {
	kafka := NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "Profile")
	kafka.Debug(true)
	writer := kafka.NewWriter("PROFILE")

	defer writer.Close()

	for i := 0; i < 5; i++ {
		p := profile{
			ID:       uuid.New().String(),
			Name:     "Paul",
			Lastname: "Mark",
		}

		pce := newProfileCreateEvent(p)
		event := BuildEventMessage(pce)

		err := kafka.SendEvent(context.Background(), writer, event)
		if err != nil {
			t.Errorf("SendEvent event %v, error: %s", event, err)
		}
	}
}

func TestKafka_SubscribeEvent(t *testing.T) {
	kafka := NewKafka([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, "Profile")
	kafka.Debug(true)
	reader := kafka.NewReader("ProfileSubscribeTest", "PROFILE")

	defer reader.Close()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	fnCallback := func(event EventMessage, err error) {
		if err != nil {
			if strings.Contains(err.Error(), "deadline exceeded") {
				t.Logf("SubscribeEvent %s", err)
			}  else {
				t.Errorf("SubscribeEvent error: %s", err)
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

	kafka.SubscribeEvent(ctx, reader, fnCallback)
}

type profile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
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
	return "PROFILE"
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