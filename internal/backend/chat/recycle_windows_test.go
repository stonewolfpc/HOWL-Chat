//go:build windows

package chat

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRecycleTreeMovesDisposableFixtureToRecycleBin exercises the real shell
// call on a disposable fixture under the test's temporary directory. It asserts
// only that the directory is gone and that the call reports success; whether the
// bytes landed in the Recycle Bin is Windows' decision.
func TestRecycleTreeMovesDisposableFixtureToRecycleBin(t *testing.T) {
	root := t.TempDir()
	victim := filepath.Join(root, "doomed-conversation")
	if err := os.MkdirAll(filepath.Join(victim, "nested"), 0o755); err != nil {
		t.Fatalf("create fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(victim, "history.ndjson"), []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("write fixture record: %v", err)
	}

	if !recycleBinAvailable(victim) {
		t.Log("no Recycle Bin reported for this volume; the shell still decides what happens")
	}

	if err := recycleTree(victim); err != nil {
		t.Fatalf("recycle rejected the disposable fixture: %v", err)
	}
	if _, err := os.Stat(victim); !os.IsNotExist(err) {
		t.Fatalf("fixture still exists after a successful recycle: %v", err)
	}
}

// TestRecycleTreeIsIdempotentForAMissingPath keeps "already gone" out of the
// failure class: deletion must never report a problem for an item that is not
// there, and must never create anything while checking.
func TestRecycleTreeIsIdempotentForAMissingPath(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "never-existed")
	if err := recycleTree(missing); err != nil {
		t.Fatalf("recycling a missing path reported a failure: %v", err)
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatalf("checking a missing path created it: %v", err)
	}
}
