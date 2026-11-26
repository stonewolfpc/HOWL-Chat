package llama

// placement.go — what llama-server actually allocated, read from its own log.
//
// The planned context, batch and GPU layer count are Howl's arithmetic. The
// engine's allocation is a separate fact: it can clamp, ignore, or spill a
// request. llama-server states that fact only at load verbosity 4, which is why
// buildLlamaServerArgs requests that level and the loader keeps the output.
// Parsing the statement is the difference between "our plan says it fits" and
// "the engine says what it placed".

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	placementOffloadPattern = regexp.MustCompile(`offloaded\s+(\d+)/(\d+)\s+layers to GPU`)
	placementBufferPattern  = regexp.MustCompile(`\b([A-Za-z0-9_]+)\s+(model|KV|compute|output) buffer size\s*=\s*([0-9]+(?:\.[0-9]+)?)\s*MiB`)
	placementContextPattern = regexp.MustCompile(`n_ctx_slot\s*=\s*(\d+)`)
	placementBatchPattern   = regexp.MustCompile(`\bn_batch\s*=\s*(\d+)`)
)

// Placement is the allocation llama-server reported while loading. It is
// deliberately a record of observations, not a second estimate.
type Placement struct {
	Known        bool               `json:"known"`
	GPULayers    int                `json:"gpu_layers"`
	TotalLayers  int                `json:"total_layers"`
	GPUMB        float64            `json:"gpu_mb"`
	RAMMB        float64            `json:"ram_mb"`
	Context      int                `json:"context"`
	Batch        int                `json:"batch"`
	KVGPUMB      float64            `json:"kv_gpu_mb"`
	KVRAMMB      float64            `json:"kv_ram_mb"`
	ComputeGPUMB float64            `json:"compute_gpu_mb"`
	ComputeRAMMB float64            `json:"compute_ram_mb"`
	Buffers      map[string]float64 `json:"buffers,omitempty"`
}

// Complete reports whether the engine stated the layer split and at least one
// model buffer. A context line alone proves the server started, not where the
// weights were placed.
func (p Placement) Complete() bool {
	return p.Known && p.TotalLayers > 0 && (p.GPUMB > 0 || p.RAMMB > 0)
}

// TotalGPUMB is the reported model + KV + compute allocation on GPU devices.
// It is an engine-buffer total, not process RSS and not total card usage.
func (p Placement) TotalGPUMB() float64 { return p.GPUMB + p.KVGPUMB + p.ComputeGPUMB }

// TotalRAMMB is the reported model + KV + compute allocation in host memory.
func (p Placement) TotalRAMMB() float64 { return p.RAMMB + p.KVRAMMB + p.ComputeRAMMB }

// ParsePlacement reads the last load's placement statements. A process can
// print several attempts (KV retry or a second load); the last statement wins,
// which is the allocation that served the health check.
func ParsePlacement(output string) Placement {
	p := Placement{Buffers: map[string]float64{}}
	for _, line := range strings.Split(output, "\n") {
		if m := placementOffloadPattern.FindStringSubmatch(line); m != nil {
			// A new load attempt starts a new allocation. Reset the previous
			// attempt's buffers before recording this one; otherwise a retry
			// would be reported as the sum of both attempts.
			p.GPULayers, _ = strconv.Atoi(m[1])
			p.TotalLayers, _ = strconv.Atoi(m[2])
			// Reset allocations only. n_batch/n_ctx can be printed before this
			// line; clearing them here would discard a fact already observed in
			// the same attempt. A later attempt prints its own values.
			p.GPUMB, p.RAMMB = 0, 0
			p.KVGPUMB, p.KVRAMMB, p.ComputeGPUMB, p.ComputeRAMMB = 0, 0, 0, 0
			p.Buffers = map[string]float64{}
			p.Known = true
		}
		if m := placementContextPattern.FindStringSubmatch(line); m != nil {
			p.Context, _ = strconv.Atoi(m[1])
		}
		if m := placementBatchPattern.FindStringSubmatch(line); m != nil {
			p.Batch, _ = strconv.Atoi(m[1])
		}
		if m := placementBufferPattern.FindStringSubmatch(line); m != nil {
			name, kind := m[1], m[2]
			mib, err := strconv.ParseFloat(m[3], 64)
			if err != nil || mib <= 0 {
				continue
			}
			gpu := isGPUModelBuffer(name)
			cpu := strings.HasPrefix(strings.ToUpper(name), "CPU") || strings.Contains(strings.ToLower(name), "host")
			switch strings.ToLower(kind) {
			case "model":
				// CPU_Host and CUDA_Host are compute areas, not weight buffers.
				if strings.Contains(strings.ToLower(name), "host") {
					continue
				}
				p.Buffers[name] = mib
				if gpu {
					p.GPUMB += mib
				} else if cpu {
					p.RAMMB += mib
				}
				p.Known = true
			case "kv":
				if gpu {
					p.KVGPUMB += mib
				} else if cpu {
					p.KVRAMMB += mib
				}
			case "compute":
				if gpu {
					p.ComputeGPUMB += mib
				} else if cpu {
					p.ComputeRAMMB += mib
				}
			}
		}
	}
	return p
}

func isGPUModelBuffer(name string) bool {
	upper := strings.ToUpper(name)
	for _, prefix := range []string{"CUDA", "METAL", "VULKAN", "SYCL", "ROCM", "HIP"} {
		if strings.HasPrefix(upper, prefix) && !strings.Contains(upper, "HOST") {
			return true
		}
	}
	return false
}

// PlacementSummary renders the observed buffer split in a stable order for
// notifications and diagnostics.
func PlacementSummary(p Placement) string {
	if len(p.Buffers) == 0 {
		return ""
	}
	names := make([]string, 0, len(p.Buffers))
	for name := range p.Buffers {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, name+" "+strconv.FormatFloat(p.Buffers[name], 'f', 0, 64)+" MiB")
	}
	return strings.Join(parts, ", ")
}
