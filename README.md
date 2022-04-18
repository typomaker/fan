# FAN
Fan is [golang](https://go.dev) library for combining channels.

## Installation

```bash
go get github.com/typomaker/fan
```

## Usage
Use to implement fanin and fanout concepts.

### Fanout
```golang
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
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)