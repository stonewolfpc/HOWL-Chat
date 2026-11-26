package sdcpp

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/types"
)

// ggufHeader is the fixed prefix of a GGUF file: magic, version, tensor count,
// then the number of key/value metadata pairs. version 3 is little-endian for
// this prefix, which is all that is needed to tell a usable checkpoint from a
// stripped one.
type ggufHeader struct {
	Version     uint32
	TensorCount uint64
	KVCount     uint64
}

func readGGUFHeader(path string) (h ggufHeader, isGGUF bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return h, false, err
	}
	defer f.Close()

	var magic [4]byte
	if _, err := io.ReadFull(f, magic[:]); err != nil {
		return h, false, err
	}
	if string(magic[:]) != "GGUF" {
		return h, false, nil
	}
	if err := binary.Read(f, binary.LittleEndian, &h.Version); err != nil {
		return h, true, err
	}
	if err := binary.Read(f, binary.LittleEndian, &h.TensorCount); err != nil {
		return h, true, err
	}
	if err := binary.Read(f, binary.LittleEndian, &h.KVCount); err != nil {
		return h, true, err
	}
	return h, true, nil
}

// validateCheckpoint rejects a checkpoint sd-server cannot identify, before it
// is spawned. A GGUF with zero key/value pairs has no general.architecture, so
// sd.cpp answers "get sd version from file failed" and the app waited the whole
// readiness timeout before telling the user to update the binary — advice that
// cannot work, because the file is the problem.
func validateCheckpoint(target *types.ModelInfo) error {
	if target == nil {
		return nil
	}
	path := target.Files["model"]
	if path == "" || !strings.EqualFold(filepath.Ext(path), ".gguf") {
		return nil
	}
	h, isGGUF, err := readGGUFHeader(path)
	if err != nil {
		return fmt.Errorf("cannot read '%s': %w", filepath.Base(path), err)
	}
	if !isGGUF {
		return nil
	}
	if h.KVCount == 0 {
		// Measured counterexample: a classic full checkpoint converted to GGUF
		// with no key/value pairs (stable-diffusion-v2-1-Q8_0) loads and samples
		// in sd-server, because sd.cpp identifies that layout from its tensor
		// names rather than from metadata. Refusing here claimed such a file
		// "can never load", which was false and blocked a working SD2.1.
		//
		// The engine still fails fast on a file it truly cannot identify:
		// waitReadyOrExit watches the child, and its "get sd version from file
		// failed" branch explains the real cause instead of blaming a version.
		fmt.Printf("[ImageEngine] %s has 0 GGUF metadata pairs with %d tensors; letting sd-server identify it from tensor layout\n",
			filepath.Base(path), h.TensorCount)
	}
	return nil
}

// loadAssumedDiskMBps is a deliberately pessimistic read rate. Loading also
// builds graphs, so the budget is generous rather than a promise.
const loadAssumedDiskMBps = 40

const (
	loadBaseTimeout = 90 * time.Second
	loadMaxTimeout  = 10 * time.Minute
)

// loadTimeoutFor sizes the readiness wait to the bytes that must be read. A
// fixed two minutes is fine for a 5 GB checkpoint on an SSD and far too short
// for a 17 GB encoder on a spinning disk, where the load then "failed" even
// though the machine had the RAM for it.
func loadTimeoutFor(files map[string]string) time.Duration {
	var bytes int64
	for _, p := range files {
		if p == "" {
			continue
		}
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			bytes += fi.Size()
		}
	}
	seconds := int(bytes / (loadAssumedDiskMBps << 20))
	d := loadBaseTimeout + time.Duration(seconds)*time.Second
	if d > loadMaxTimeout {
		d = loadMaxTimeout
	}
	return d
}

// stderrTail appends what sd-server actually said to a load failure. Without it
// a timeout is reported as a bare "not ready", which hides an out-of-memory or a
// missing DLL behind a number of seconds.
func stderrTail(buf *strings.Builder, max int) string {
	if buf == nil {
		return ""
	}
	s := strings.TrimSpace(buf.String())
	if s == "" {
		return ""
	}
	if len(s) > max {
		s = s[len(s)-max:]
	}
	return "\nLast output from sd-server:\n" + s
}
