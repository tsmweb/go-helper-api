package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// Producer injects mock kafka.Producer dependency.
type Producer struct {
	mock.Mock
}

// Publish represents the simulated method for the Publish feature in the kafka.Producer layer.
func (m *Producer) Publish(ctx context.Context, key []byte, values ...[]byte) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

// Close represents the simulated method for the Close feature in the kafka.Producer layer.
func (m *Producer) Close() {}
