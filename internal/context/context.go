package context

import (
	"context"
)

type keyType int

const (
	txKey keyType = iota
)

func WithTxID(ctx context.Context, txID int64) context.Context {
	return context.WithValue(ctx, txKey, txID)
}

func TxIDFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(txKey).(int64)

	return v
}
