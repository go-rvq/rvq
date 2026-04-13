package context_utils

import "context"

func OrContextValue[T any](ctx *context.Context, key any, or func() T) T {
	if *ctx == nil {
		*ctx = context.Background()
	}
	v, ok := (*ctx).Value(key).(T)
	if ok {
		return v
	}
	v = or()
	*ctx = context.WithValue(*ctx, key, v)
	return v
}

func OrContextValueE[T any](ctx *context.Context, key any, or func() (T, error)) (T, error) {
	if *ctx == nil {
		*ctx = context.Background()
	}
	v, ok := (*ctx).Value(key).(T)
	if ok {
		return v, nil
	}
	v, err := or()
	if err != nil {
		return v, err
	}
	*ctx = context.WithValue(*ctx, key, v)
	return v, nil
}
