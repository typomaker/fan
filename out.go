package fan

import (
	"context"
	"time"
)

type out[T any] struct {
	ctx context.Context
	ch  (chan<- T)
}

func newOut[T any](ctx context.Context, ch chan<- T) *out[T] {
	var p = &out[T]{
		ctx: ctx,
		ch:  ch,
	}
	return p
}
func (it *out[T]) push(msgs ...T) {
	for _, msg := range msgs {
		select {
		case it.ch <- msg:
		case <-time.After(time.Microsecond):
		}
	}
}
