package h

import (
	"context"
	"errors"
)

func GetUsedIDFromContext(ctx context.Context, key string) (int, error) {
	userId, ok := ctx.Value(key).(int)
	if !ok {
		return 0, errors.New("Invalid userId")
	}
	return userId, nil
}
