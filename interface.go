// Package progress offers a simple, standardized interface for monitoring progress in your applications.
//
// This package exposes the progress.Reader interface, designed to be similar in spirit to io.Reader but
// for progress reporting. This interface makes tracking long-running tasks and their progress a breeze.
//
// The interface defines two key methods:
//
// - DoneChan() returns a channel that signals when the task is done, and a boolean indicating if the task is completed.
// - Count() provides the current count of finished tasks and the total number of tasks.
//
// In addition to the interface, the package provides several useful tools:
//
// - progress.NewLongRunningJob: A concurrency-safe utility for instrumenting long-running functions with a progress.Reader.
// - progress.Extend(Reader): A simple wrapper that adds useful methods to a Reader, such as estimating remaining time.
// - progress.Logger(Reader): Starts a goroutine to log updates from a Reader.
// - progress.WaitGroup: A drop-in replacement for sync.WaitGroup that implements the Reader interface.
//
// By using this package, you can make your long-running tasks and functions more informative and user-friendly,
// improving the overall experience for your library consumers.
//
// For those interested in contributing, there are several areas where the functionality of this package could be extended,
// such as satisfying progress.Reader with just a count, returning a progress.Reader from a slice or channel, creating
// a pretty terminal progress bar, or experimenting with passing return types through.
package progress

// Reader is an interface that all progress trackers should satisfy.
// This allows for other libraries to work with progress trackers generally,
// and for library developers to create compatible progress trackers.
type Reader interface {
	// DoneChan returns a channel that will be closed upon completion.
	// The boolean return specifies whether the channel is already closed.
	DoneChan() (chan struct{}, bool)

	// Count returns the amount of items processed and the total amount.
	// The amount of items is guaranteed to be less than or equal to the total.
	// The total amount is guaranteed to be at least 1.
	// If they are equal, it is guaranteed that InProgress will return false.
	// Both values may increase at any time and are guaranteed never to decrease.
	// Count must be safe to be called concurrently.
	Count() (uint64, uint64)
}
