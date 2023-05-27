## ⏱️ `type progress.Reader interface`

[![Go Reference](https://pkg.go.dev/badge/github.com/libfor/progress.svg)](https://pkg.go.dev/github.com/libfor/progress) [![Run unit tests](https://github.com/libfor/progress/actions/workflows/test_on_push.yaml/badge.svg)](https://github.com/libfor/progress/actions/workflows/test_on_push.yaml) 

**Progress monitoring** is a crucial part of many applications, and deserves a simple, standardized interface. The `progress.Reader` interface aims to do just that.

    progress.Reader interface
	    DoneChan() (chan struct{}, bool)
	    Count() (finished uint64, total uint64)


To drive adoption and provide practical utility, we've also included some handy tools right out of the box! 📦

### `progress.Extend(Reader)`

This simple wrapper adds helpful functions to a Reader:

    PerSecond() float64
    Remaining() time.Duration
    Percentage() float64
    InProgress() bool

### `progress.Logger(Reader)`

Instantly start a goroutine that keeps an eye on a Reader and logs informative updates! 👀 Here's what it looks like:

    searching:  49% (20.69/s, 1.256483582s remaining)
    searching:  56% (18.01/s, 1.221804358s remaining)
    searching:  80% (20.37/s, 490.981158ms remaining)
    searching: 100% (21.98/s, 0s remaining)

### `progress.NewLongRunningJob`

This concurrency-safe utility makes instrumenting your long-running functions with a Reader implementation absolutely painless. Support progress and make your function consumers extremely happy! 😃

### `progress.WaitGroup`

Drop in replacement for sync.WaitGroup that satisfies the Reader interface.


## 📝 Todo List

Here's some simple ideas for those looking to contribute.

- [x] Satisfy `progress.Reader` with just a count.
- [ ] Return a `progress.Reader` from a slice or channel.
- [ ] Create a pretty terminal progress bar.
- [ ] Experiment with passing return types through.