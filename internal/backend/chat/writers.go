package chat

import (
	"fmt"
	"os"
)

// This file holds the only filesystem write calls in the conversation package.
// Everything else goes through these four helpers, and
// boundary_analysis_test.go fails if a write appears anywhere else. That
// property, not a startup probe, is what makes the write boundary structural.

// ensureDir creates a directory (and its parents) when it is missing.
func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// requireExistingDir refuses a directory that is not there. Mutators that are
// not the creation route use this instead of ensureDir, so a write that merely
// mentions an identifier can never bring a conversation back into existence.
func requireExistingDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

// writeFileAtomic replaces a document by writing a sibling temporary file and
// renaming it over the target, so a reader never observes a half-written
// document and a crash cannot truncate the live one.
func writeFileAtomic(filename string, data []byte, perm os.FileMode) error {
	tmp := filename + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	if err := os.Rename(tmp, filename); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// claimDir creates a directory only if it does not already exist. Creation uses
// this instead of ensureDir so two conversations can never silently share one
// directory: an existing directory is a collision to report, never a target to
// adopt.
func claimDir(path string) error {
	return os.Mkdir(path, 0o755)
}

// appendLine appends one complete line to a log file.
//
// If the existing file does not end in a newline, the tail is terminated in this
// same call before the new line is written. An unterminated record is therefore
// sealed and keeps its own line; the append is never concatenated onto it, which
// is the rule TASK-03 row 579 requires.
func appendLine(path string, line []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() > 0 {
		last := make([]byte, 1)
		if _, err := f.ReadAt(last, info.Size()-1); err != nil {
			return err
		}
		if last[0] != '\n' {
			if _, err := f.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	if _, err := f.Write(line); err != nil {
		return err
	}
	return nil
}
