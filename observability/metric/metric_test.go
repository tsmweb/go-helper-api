package metric

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	kafka "github.com/tsmweb/go-helper-api/observability/internal/mock"
)

func TestMetric_New(t *testing.T) {
	m, err := newMetric("localhost")
	if err != nil {
		t.Error(err)
	}

	b := m.toJSON()
	if b == nil {
		t.Error("Error toJSON()")
	} else {
		t.Log(string(b))
	}
}

func TestMetric_Start(t *testing.T) {
	producer := new(kafka.Producer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			t.Log(string((args[2].([][]byte))[0]))
		})
	producer.On("Close")

	err := Start("localhost", 2, producer)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(6 * time.Second)
	Stop()
}
