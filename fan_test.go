package fan_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/typomaker/fan"
	"go.uber.org/goleak"
)

func TestFunCancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	var ctx, cancel = context.WithCancel(context.Background())
	cancel()
	var f = fan.New[int](ctx)
	var err error
	err = f.In(ctx, make(chan int))
	assert.NotNil(t, err)
	err = f.Out(ctx, make(chan int))
	assert.NotNil(t, err)
}

func TestInOut(t *testing.T) {
	defer goleak.VerifyNone(t)
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var f = fan.New[int](ctx)
	var in1 = make(chan int, 2)
	var in2 = make(chan int, 2)
	var out1 = make(chan int, 4)
	var out2 = make(chan int, 4)
	f.Out(ctx, out1, out2)
	time.Sleep(time.Millisecond)

	f.In(ctx, in1, in2)
	for i := 0; i < 2; i++ {
		in1 <- i
		in2 <- i
	}
	time.Sleep(time.Millisecond)

	assert.Len(t, out1, 4)
	assert.Len(t, out2, 4)
}
func TestInCancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var f = fan.New[int](ctx)
	var ctxin, cancelin = context.WithCancel(context.Background())
	var in = make(chan int, 2)
	f.In(ctxin, in)
	cancelin()
	time.Sleep(time.Millisecond)
	in <- 1
	in <- 2
	time.Sleep(time.Millisecond)
	assert.Len(t, in, 2)
}
func TestOutCancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var f = fan.New[int](ctx)

	var ctxout, cancelout = context.WithCancel(context.Background())
	var out = make(chan int, 1)
	f.Out(ctxout, out)
	cancelout()
	time.Sleep(time.Millisecond)

	var in = make(chan int, 1)
	in <- 1
	f.In(ctx, in)
	time.Sleep(time.Millisecond)

	assert.Len(t, out, 0)
}
