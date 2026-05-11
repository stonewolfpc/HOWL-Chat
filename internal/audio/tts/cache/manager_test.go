package cache

import (
	"testing"

	"howl-chat/internal/audio/types"
)

func TestManagerStoreAndGet(t *testing.T) {
	mgr := NewManager(1)
	key := "voice_hello"
	value := &types.SynthesisResult{AudioPath: "a.wav", Duration: 1.2}

	if err := mgr.Store(key, value); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	got, ok := mgr.Get(key)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got.AudioPath != "a.wav" {
		t.Fatalf("unexpected path: %s", got.AudioPath)
	}
	stats := mgr.GetStats()
	if stats.HitCount != 1 {
		t.Fatalf("expected hit count 1, got %d", stats.HitCount)
	}
}

func TestManagerEvictsWhenOverLimit(t *testing.T) {
	mgr := NewManager(1) // 1MB
	large := &types.SynthesisResult{
		AudioPath: "big.wav",
		Duration:  20,
		Metadata: types.AudioMetadata{
			BitRate: 512,
		},
	}

	if err := mgr.Store("a", large); err != nil {
		t.Fatalf("Store a failed: %v", err)
	}
	if err := mgr.Store("b", large); err != nil {
		t.Fatalf("Store b failed: %v", err)
	}

	stats := mgr.GetStats()
	if stats.SizeBytes > stats.MaxSizeBytes {
		t.Fatalf("cache exceeded max size: %d > %d", stats.SizeBytes, stats.MaxSizeBytes)
	}
}
