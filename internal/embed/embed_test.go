package embed

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure MockEmbedder satisfies the Embedder interface
var _ Embedder = (*MockEmbedder)(nil)

type MockEmbedder struct {
	Response []float32
	Err      error
}

func (m *MockEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	return m.Response, m.Err
}

func TestMockEmbedder_Success(t *testing.T) {
	mock := &MockEmbedder{
		Response: make([]float32, 768),
	}

	ctx := context.Background()
	input := "test input"
	result, err := mock.Embed(ctx, input)

	assert.NoError(t, err, "expected no error on successful embed")
	assert.NotNil(t, result, "expected a non-nil result")
	assert.Equal(t, 768, len(result), "expected result to have 768 dimensions")
}

func TestMockEmbedder_Error(t *testing.T) {
	mock := &MockEmbedder{
		Err: errors.New("mock error"),
	}

	ctx := context.Background()
	result, err := mock.Embed(ctx, "fail case")

	assert.Error(t, err, "expected an error on embed failure")
	assert.Nil(t, result, "expected result to be nil on error")
}
