package fan

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Fan[T any] struct {
	mtx    sync.RWMutex
	ctx    context.Context
	newin  chan *in[T]
	newout chan *out[T]
}

// New instantiate a new fan
// to avoid goroutine leaks, use a cancelable context to stop the background fan process
func New[T any](ctx context.Context) *Fan[T] {
	var f = &Fan[T]{}
	f.ctx = ctx
	f.newin = make(chan *in[T], 128)
	f.newout = make(chan *out[T], 128)
	go f.loop()
	return f
}

// In register the input channel in a fan
// use cancelable context for remove a input channel from a fan
// returns error if context of a fan is closed
func (f *Fan[T]) In(ctx context.Context, ch ...(<-chan T)) (err error) {
	if f.ctx.Err() != nil {
		return fmt.Errorf("fan: %w", f.ctx.Err())
	}
	for _, c := range ch {
		f.newin <- newIn(ctx, c)
	}
	return
}

// Out register the output channel in a fan
// use cancelable context for remove a output channel from a fan
// returns error if context of a fan is closed
func (f *Fan[T]) Out(ctx context.Context, ch ...(chan<- T)) (err error) {
	if f.ctx.Err() != nil {
		return fmt.Errorf("fan: %w", f.ctx.Err())
	}
	for _, c := range ch {
		f.newout <- newOut(ctx, c)
	}
	return
}
func (f *Fan[T]) loop() {
	var ins []*in[T]
	var outs []*out[T]
	var msgs []T
	for {
		if msgs == nil || cap(msgs) > 256 {
			msgs = make([]T, 0, 256)
		} else if len(msgs) != 0 {
			msgs = msgs[:0]
		}

		select {
		case <-f.ctx.Done():
			return
		case c := <-f.newout:
			outs = append(outs, c)
			for len(f.newout) > 0 {
				outs = append(outs, <-f.newout)
			}
		case p := <-f.newin:
			ins = append(ins, p)
			for len(f.newin) > 0 {
				ins = append(ins, <-f.newin)
			}
		case <-time.After(time.Millisecond):
		}

		for i := 0; i < len(ins); i++ {
			select {
			case <-ins[i].ctx.Done():
				ins = append(ins[:i], ins[i+1:]...)
			default:
			}
		}
		for i := 0; i < len(outs); i++ {
			select {
			case <-outs[i].ctx.Done():
				outs = append(outs[:i], outs[i+1:]...)
			default:
			}
		}

		for _, in := range ins {
			msgs = append(msgs, in.pull()...)
		}
		for _, out := range outs {
			out.push(msgs...)
		}
	}
}
