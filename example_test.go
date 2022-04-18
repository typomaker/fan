package fan_test

import (
	"context"
	"fmt"
	"time"

	"github.com/typomaker/fan"
)

func ExampleNew() {
	// use cancelable context to avoid goroutine leaks
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// instantiate a new fan of `string` type
	var f = fan.New[string](ctx)

	// define the output channel
	var out = make(chan string)
	// register output in a fan
	f.Out(ctx, out)

	// used in a example for syncs
	time.Sleep(time.Millisecond)

	// define the input channel
	var in = make(chan string)
	// use cancelable context to remove a input channel from a fan
	var ctxin, cancelin = context.WithCancel(context.Background())
	defer cancelin()
	// register a input channel in a fan
	f.In(ctxin, in)
	// send test message
	in <- "through a fan"

	fmt.Println(<-out)
	//Output: through a fan
}
