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
