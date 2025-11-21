package utils_test

import (
	"context"
	"errors"
	"testing"

	utils "github.com/heaveless/dbz-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestWithFallback_PrimarySuccess(t *testing.T) {
	ctx := context.Background()

	primary := func(ctx context.Context) (int, error) {
		return 10, nil
	}
	secondary := func(ctx context.Context) (int, error) {
		return 20, nil
	}
	shouldFallback := func(err error) bool {
		return true
	}

	v, err := utils.WithFallback(ctx, primary, secondary, shouldFallback)
	assert.NoError(t, err)
	assert.Equal(t, 10, v)
}

func TestWithFallback_PrimaryFail_ShouldFallback(t *testing.T) {
	ctx := context.Background()

	primary := func(ctx context.Context) (int, error) {
		return 0, errors.New("fail")
	}
	secondary := func(ctx context.Context) (int, error) {
		return 99, nil
	}
	shouldFallback := func(err error) bool {
		return true
	}

	v, err := utils.WithFallback(ctx, primary, secondary, shouldFallback)
	assert.NoError(t, err)
	assert.Equal(t, 99, v)
}

func TestWithFallback_PrimaryFail_ShouldNotFallback(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("fail")
	primary := func(ctx context.Context) (int, error) {
		return 0, expectedErr
	}
	secondary := func(ctx context.Context) (int, error) {
		return 99, nil
	}
	shouldFallback := func(err error) bool {
		return false
	}

	v, err := utils.WithFallback(ctx, primary, secondary, shouldFallback)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, v)
}
