package fakehelper

import "os"

// publishNoReplace gives the final output a directory entry without ever
// replacing a file that may have appeared after the pre-write existence check.
// A hard link is atomic and fails if the destination exists; the temporary
// name is then unlinked. os.Rename is deliberately not used because Windows
// implements it with MOVEFILE_REPLACE_EXISTING.
func publishNoReplace(part, final string) error {
	if err := os.Link(part, final); err != nil {
		return err
	}
	if err := os.Remove(part); err != nil {
		_ = os.Remove(final)
		return err
	}
	return nil
}
