package host

import (
	"context"
	"sync"
)

// Streams owns the one active generation token.
//
// A turn registers its cancel function and receives a token. Shutdown cancels
// whatever is active; a finishing generation may only clear its own token, so
// an older turn cannot erase a newer one's cancellation.
type Streams struct {
	mu     sync.Mutex
	cancel context.CancelFunc
	gen    uint64
}

// NewStreams returns an idle stream authority.
func NewStreams() *Streams { return &Streams{} }

// Begin records the cancel function for one generation and returns its token.
func (s *Streams) Begin(cancel context.CancelFunc) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gen++
	s.cancel = cancel
	return s.gen
}

// End clears the token only if it is still the current generation.
func (s *Streams) End(token uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.gen == token {
		s.cancel = nil
	}
}

// Cancel stops the active generation, if any. A finished stream is a no-op.
func (s *Streams) Cancel() {
	s.mu.Lock()
	cancel := s.cancel
	s.gen++
	s.cancel = nil
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}
