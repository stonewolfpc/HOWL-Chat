package durable

import "time"

// LockFile exposes the store's exclusive file lock for callers that need the same
// primitive over a file of their own. Conversation persistence uses it to serialise
// one conversation without adopting the document store.
//
// A zero timeout attempts once and reports ErrBusy rather than waiting, which is
// what a refusal-on-contention caller wants. The returned function releases the
// lock; it is safe to call once.
func LockFile(path string, timeout time.Duration) (func(), error) {
	lock, err := acquireLock(path, timeout)
	if err != nil {
		return nil, err
	}
	return func() { lock.release() }, nil
}
