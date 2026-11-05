package gguf

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// writeGGUFFile writes a minimal GGUF v3 file (header plus metadata, no
// tensors) and returns its path. Tests use it so classification and pairing run
// through the real parser instead of a duplicated decision table.
func writeGGUFFile(t *testing.T, dir, name string, kvs map[string]any) string {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create dir: %v", err)
	}
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer file.Close()

	writeString := func(value string) {
		var size [8]byte
		binary.LittleEndian.PutUint64(size[:], uint64(len(value)))
		if _, err := file.Write(size[:]); err != nil {
			t.Fatalf("write string size: %v", err)
		}
		if _, err := file.WriteString(value); err != nil {
			t.Fatalf("write string: %v", err)
		}
	}

	if err := binary.Write(file, binary.LittleEndian, uint32(0x46554747)); err != nil { // "GGUF"
		t.Fatalf("write magic: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(3)); err != nil {
		t.Fatalf("write version: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(0)); err != nil { // tensors
		t.Fatalf("write tensor count: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(len(kvs))); err != nil {
		t.Fatalf("write kv count: %v", err)
	}

	for key, value := range kvs {
		writeString(key)
		switch typed := value.(type) {
		case bool:
			if err := binary.Write(file, binary.LittleEndian, uint32(7)); err != nil {
				t.Fatalf("write bool type: %v", err)
			}
			var flag byte
			if typed {
				flag = 1
			}
			if _, err := file.Write([]byte{flag}); err != nil {
				t.Fatalf("write bool: %v", err)
			}
		case uint32:
			if err := binary.Write(file, binary.LittleEndian, uint32(4)); err != nil {
				t.Fatalf("write uint32 type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, typed); err != nil {
				t.Fatalf("write uint32: %v", err)
			}
		case string:
			if err := binary.Write(file, binary.LittleEndian, uint32(8)); err != nil {
				t.Fatalf("write string type: %v", err)
			}
			writeString(typed)
		case []string:
			if err := binary.Write(file, binary.LittleEndian, uint32(9)); err != nil {
				t.Fatalf("write array type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, uint32(8)); err != nil {
				t.Fatalf("write array element type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, uint64(len(typed))); err != nil {
				t.Fatalf("write array count: %v", err)
			}
			for _, item := range typed {
				writeString(item)
			}
		default:
			t.Fatalf("unsupported test metadata type %T", value)
		}
	}

	return path
}
