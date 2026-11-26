package sdcpp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ComponentBytes is the byte split of a model file's tensor data by role. It is
// read from the safetensors header only: no tensor payload is read and no model
// is loaded. sd.cpp places these modules independently, so charging the whole
// file to one module overstates the VRAM requirement for a full checkpoint.
type ComponentBytes struct {
	Diffusion   int64
	Conditioner int64
	VAE         int64
	Other       int64
}

// Total is the tensor payload measured by the header.
func (c ComponentBytes) Total() int64 {
	return c.Diffusion + c.Conditioner + c.VAE + c.Other
}

// Known reports whether at least one tensor was classified into a model role.
func (c ComponentBytes) Known() bool {
	return c.Diffusion > 0 || c.Conditioner > 0 || c.VAE > 0
}

const (
	maxSafetensorsHeader = 64 << 20
	safetensorsAlign     = 8
)

// ReadComponentBytes reads the component split for a model file. Safetensors is
// parsed from its header. GGUF and every other format return false so callers
// keep the conservative whole-file estimate; a GGUF file must not be guessed
// from its name.
func ReadComponentBytes(path string) (ComponentBytes, bool) {
	if !strings.EqualFold(filepath.Ext(path), ".safetensors") {
		return ComponentBytes{}, false
	}
	parts, err := readSafetensorsComponents(path)
	if err != nil || !parts.Known() {
		return ComponentBytes{}, false
	}
	return parts, true
}

type safetensorsTensor struct {
	DataOffsets []int64 `json:"data_offsets"`
}

func readSafetensorsComponents(path string) (ComponentBytes, error) {
	f, err := os.Open(path)
	if err != nil {
		return ComponentBytes{}, err
	}
	defer f.Close()

	var headerLen uint64
	if err := binary.Read(f, binary.LittleEndian, &headerLen); err != nil {
		return ComponentBytes{}, err
	}
	if headerLen == 0 || headerLen > maxSafetensorsHeader {
		return ComponentBytes{}, fmt.Errorf("safetensors header length %d is outside the supported range", headerLen)
	}
	header := make([]byte, headerLen)
	if _, err := io.ReadFull(f, header); err != nil {
		return ComponentBytes{}, err
	}

	var tensors map[string]safetensorsTensor
	if err := json.Unmarshal(header, &tensors); err != nil {
		return ComponentBytes{}, err
	}

	var out ComponentBytes
	for name, tensor := range tensors {
		if name == "__metadata__" || len(tensor.DataOffsets) != 2 {
			continue
		}
		size := tensor.DataOffsets[1] - tensor.DataOffsets[0]
		if size <= 0 {
			continue
		}
		switch componentRole(name) {
		case "diffusion":
			out.Diffusion += size
		case "conditioner":
			out.Conditioner += size
		case "vae":
			out.VAE += size
		default:
			out.Other += size
		}
	}
	return out, nil
}

// componentRole maps sd.cpp tensor prefixes to the module sd-server places.
// Prefixes are checked before "other" so an unknown family falls back to the
// whole checkpoint instead of silently dropping bytes.
func componentRole(name string) string {
	switch {
	case strings.HasPrefix(name, "model.diffusion_model."),
		strings.HasPrefix(name, "diffusion_model."):
		return "diffusion"
	case strings.HasPrefix(name, "conditioner."),
		strings.HasPrefix(name, "cond_stage_model."),
		strings.HasPrefix(name, "text_encoders."),
		strings.HasPrefix(name, "te."):
		return "conditioner"
	case strings.HasPrefix(name, "first_stage_model."),
		strings.HasPrefix(name, "vae."),
		strings.HasPrefix(name, "autoencoder."):
		return "vae"
	default:
		return ""
	}
}
