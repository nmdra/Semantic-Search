package embed

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Embedder = (*MockEmbedder)(nil)

type MockEmbedder struct {
	Response []float32
	Err      error
}

func (m *MockEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	return m.Response, m.Err
}

func TestMockEmbedder(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &MockEmbedder{
			Response: make([]float32, 768),
		}

		ctx := context.Background()
		input := "test input"

		result, err := mock.Embed(ctx, input)

		assert.NoError(t, err, "expected no error on successful embed")
		assert.NotNil(t, result, "expected a non-nil result")
		assert.Equal(t, 768, len(result), "expected result to have 768 dimensions")
	})

	t.Run("Error", func(t *testing.T) {
		expectedErr := errors.New("mock error")
		mock := &MockEmbedder{
			Err: expectedErr,
		}

		ctx := context.Background()
		result, err := mock.Embed(ctx, "fail case")

		assert.Error(t, err, "expected an error on embed failure")
		assert.ErrorIs(t, err, expectedErr, "error should match mock error")
		assert.Nil(t, result, "expected result to be nil on error")
	})

	t.Run("Empty input", func(t *testing.T) {
		mock := &MockEmbedder{
			Response: make([]float32, 768),
		}

		ctx := context.Background()
		result, err := mock.Embed(ctx, "")

		assert.NoError(t, err, "expected no error on empty input")
		assert.NotNil(t, result, "expected result even with empty input")
		assert.Equal(t, 768, len(result), "expected result to have 768 dimensions")
	})
}

func BenchmarkMockEmbedder(b *testing.B) {
	mock := &MockEmbedder{
		Response: make([]float32, 768),
	}

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, _ = mock.Embed(ctx, "benchmark input")
	}
}
