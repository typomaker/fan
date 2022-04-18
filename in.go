package fan

import "context"

type in[T any] struct {
	ctx context.Context
	ch  (<-chan T)
}

func newIn[T any](ctx context.Context, ch <-chan T) *in[T] {
	var p = &in[T]{
		ctx: ctx,
		ch:  ch,
	}
	return p
}
func (it *in[T]) pull() (msgs []T) {
	select {
	case msg := <-it.ch:
		msgs = append(msgs, msg)
		for len(it.ch) > 0 {
			msgs = append(msgs, <-it.ch)
		}
	default:
	}
	return
}
