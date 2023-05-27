# `progress.Reader`


**Progress** is something that deserves a standard, simple interface. Almost every program and library has long-running tasks. This library exposes the `progress.Reader` interface, which is minimal but flexible enough to build all kinds of progress bars or monitors.

    progress.Reader interface
	    DoneChan() (chan struct{}, bool)
	    Count() (uint64, uint64)

To encourage adoption, this package also bootstraps a few helpful utilities.

[![Run unit tests](https://github.com/libfor/progress/actions/workflows/test_on_push.yaml/badge.svg)](https://github.com/libfor/progress/actions/workflows/test_on_push.yaml) 

### `progress.Extend(Reader)`

This simple wrapper provides builds some nice functionality on top of a Reader, such as: 

    PerSecond() float64
    Remaining() time.Duration
    Percentage() float64
    InProgress() bool

### `progress.Logger`

Easily spin up a goroutine that watches a Reader and prints helpful log lines, like such:

    searching:  49% (20.69/s, 1.256483582s remaining)
    searching:  56% (18.01/s, 1.221804358s remaining)
    searching:  80% (20.37/s, 490.981158ms remaining)
    searching: 100% (21.98/s, 0s remaining)

### `progress.WaitGroup`

Drop in replacement for sync.WaitGroup that satisfies the Reader interface.

### `progress.NewLongRunningJob`

Easy to use concurrency-safe type for instrumenting your long-running functions with a Reader implementation. It's painless to support `progress` and might make the consumers of your function very happy.

### `todo:`

Here's some simple ideas for those looking to contribute.

- [x] Satisfy `progress.Reader` with just a count.
- [ ] Return a `progress.Reader` from a slice or channel.
- [ ] Create a pretty terminal prgoress bar.
- [ ] Experiment with passing return types through.