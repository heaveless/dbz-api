package utils

import "context"

func WithFallback[T any](
	ctx context.Context,
	primary func(ctx context.Context) (T, error),
	secondary func(ctx context.Context) (T, error),
	shouldFallback func(error) bool,
) (T, error) {
	var zero T

	v, err := primary(ctx)
	if err == nil {
		return v, nil
	}

	if !shouldFallback(err) {
		return zero, err
	}

	return secondary(ctx)
}
