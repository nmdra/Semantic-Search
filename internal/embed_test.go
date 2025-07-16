package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEmbedder struct {
	Response []float32
	Err      error
}

func (m *MockEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	return m.Response, m.Err
}

func TestEmbed_Success(t *testing.T) {
	mock := &MockEmbedder{
		Response: make([]float32, 768),
	}

	result, err := mock.Embed(context.Background(), "hello world")
	assert.NoError(t, err)
	assert.Equal(t, 768, len(result))
}

func TestEmbed_Failure(t *testing.T) {
	mock := &MockEmbedder{
		Err: errors.New("embedding failed"),
	}

	result, err := mock.Embed(context.Background(), "hello world")
	assert.Error(t, err)
	assert.Nil(t, result)
}
