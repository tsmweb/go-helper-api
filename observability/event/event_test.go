package event

import (
	"context"
	"errors"
	"testing"
)

func Test_Init(t *testing.T) {
	err := Init(context.Background(), []string{"localhost:9094"}, "CLIENT_ID", "TOPIC_NAME")
	if err != nil {
		t.Error(err)
	}

	err = Init(context.Background(), []string{"localhost:9094"}, "CLIENT_ID", "TOPIC_NAME")
	if err != nil && errors.Is(err, ErrRunning) {
		t.Log(err.Error())
	}
}

func Test_Send(t *testing.T) {
	err := Init(context.Background(), []string{"localhost:9094"}, "CLIENT_ID", "TOPIC_NAME")
	if err != nil {
		t.Error(err)
	}

	event := New("localhost", "Test", "Object Not Found",
		"Could not find the requested object.")

	if err = Send(event); err != nil {
		t.Error(err)
	}

	Close()

	err = Send(event)
	if err != nil && errors.Is(err, ErrClosed) {
		t.Log(err)
	}
}
