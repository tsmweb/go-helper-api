package event

import (
	"errors"
	"github.com/stretchr/testify/mock"
	kafka "github.com/tsmweb/go-helper-api/observability/internal/mock"
	"testing"
)

func TestEvent_Init(t *testing.T) {
	producer := new(kafka.Producer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			t.Log(string((args[2].([][]byte))[0]))
		})
	producer.On("Close")

	err := Init(producer)
	if err != nil {
		t.Error(err)
	}

	err = Init(producer)
	if err != nil && errors.Is(err, ErrRunning) {
		t.Log(err.Error())
	}
}

func TestEvent_Send(t *testing.T) {
	producer := new(kafka.Producer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			t.Log(string((args[2].([][]byte))[0]))
		})
	producer.On("Close")

	err := Init(producer)
	if err != nil {
		t.Error(err)
	}

	event := New("localhost", "Test", "Object Not Found", Warning,
		"Could not find the requested object.")

	t.Log(string(event.toJSON()))

	if err = Send(event); err != nil {
		t.Error(err)
	}

	Close()

	err = Send(event)
	if err != nil && errors.Is(err, ErrClosed) {
		t.Log(err)
	}
}
