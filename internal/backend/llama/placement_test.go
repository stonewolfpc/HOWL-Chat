package llama

import (
	"strings"
	"testing"
)

func TestParsePlacementReadsVerbatimCPUOnlyLoad(t *testing.T) {
	// Captured from runtime/llama/llama-server.exe with --log-verbosity 4 on
	// this machine. The timestamps and spacing are part of the input because
	// the parser sees the engine's real log, not a tidy fixture.
	const output = `0.00.117.300 I srv  llama_server: initializing ...
0.00.642.084 I load_tensors: offloading 0 repeating layers to GPU
0.00.642.089 I load_tensors: offloaded 0/33 layers to GPU
0.00.642.092 I load_tensors:   CPU_Mapped model buffer size =  3024.67 MiB
0.00.647.457 I llama_context:        CPU  output buffer size =     1.96 MiB
0.00.651.484 I llama_kv_cache:        CPU KV buffer size =    32.00 MiB
0.00.659.302 I sched_reserve:      CUDA0 compute buffer size =    32.56 MiB
0.00.659.308 I sched_reserve:  CUDA_Host compute buffer size =     4.03 MiB
0.01.437.351 I srv    load_model: initializing, n_slots = 4, n_ctx_slot = 256, kv_unified = 'true'
`
	got := ParsePlacement(output)
	if !got.Known || !got.Complete() {
		t.Fatalf("placement not recognized: %+v", got)
	}
	if got.GPULayers != 0 || got.TotalLayers != 33 {
		t.Fatalf("layers = %d/%d, want 0/33", got.GPULayers, got.TotalLayers)
	}
	if got.GPUMB != 0 {
		t.Fatalf("CPU-only load reported GPU model buffers: %v", got.GPUMB)
	}
	if got.RAMMB != 3024.67 {
		t.Fatalf("RAM model buffer = %v, want 3024.67", got.RAMMB)
	}
	if got.Context != 256 {
		t.Fatalf("context = %d, want the engine's 256", got.Context)
	}
	if _, ok := got.Buffers["CPU"]; ok {
		t.Fatalf("CPU output/KV buffers must not be counted as model weights: %v", got.Buffers)
	}
	if _, ok := got.Buffers["CUDA_Host"]; ok {
		t.Fatalf("host compute buffers must not be counted as model weights: %v", got.Buffers)
	}
	if got.TotalGPUMB() != 32.56 {
		t.Fatalf("reported GPU total = %v, want the 32.56 MiB compute buffer", got.TotalGPUMB())
	}
	if diff := got.TotalRAMMB() - (3024.67 + 32 + 4.03); diff > 0.001 || diff < -0.001 {
		t.Fatalf("reported RAM total = %v, want weights + CPU KV + host compute", got.TotalRAMMB())
	}
}

func TestParsePlacementReadsSplitGPUAndCPU(t *testing.T) {
	const output = `0.00.642.089 I load_tensors: offloaded 24/33 layers to GPU
0.00.642.092 I load_tensors:      CUDA0 model buffer size =  5210.25 MiB
0.00.642.094 I load_tensors: CPU_Mapped model buffer size =   680.50 MiB
0.01.437.351 I srv    load_model: initializing, n_slots = 1, n_ctx_slot = 8192, kv_unified = 'true'
`
	got := ParsePlacement(output)
	if got.GPULayers != 24 || got.TotalLayers != 33 {
		t.Fatalf("layers = %d/%d, want 24/33", got.GPULayers, got.TotalLayers)
	}
	if got.GPUMB != 5210.25 || got.RAMMB != 680.50 {
		t.Fatalf("buffers = GPU %v RAM %v", got.GPUMB, got.RAMMB)
	}
	if got.Context != 8192 {
		t.Fatalf("context = %d, want 8192", got.Context)
	}
	if summary := PlacementSummary(got); !strings.Contains(summary, "CUDA0 5210 MiB") || !strings.Contains(summary, "CPU_Mapped 680 MiB") {
		t.Fatalf("summary lost a buffer: %q", summary)
	}
}

func TestParsePlacementReadsWorkingBuffersAndBatch(t *testing.T) {
	// Captured lines from the desktop run's verbose llama-server log.
	const output = `0.00.437.888 I llama_context: n_batch               = 512
0.00.642.089 I load_tensors: offloaded 8/29 layers to GPU
0.00.642.092 I load_tensors:      CUDA0 model buffer size =   117.00 MiB
0.00.642.094 I load_tensors: CPU_Mapped model buffer size =   186.00 MiB
0.00.456.385 I llama_kv_cache:        CPU KV buffer size =    68.25 MiB
0.00.457.908 I llama_kv_cache:      CUDA0 KV buffer size =    22.75 MiB
0.00.467.723 I sched_reserve:      CUDA0 compute buffer size =    27.25 MiB
0.00.467.730 I sched_reserve:  CUDA_Host compute buffer size =    11.27 MiB
0.00.438.410 I llama_context:  CUDA_Host  output buffer size =     1.53 MiB
`
	got := ParsePlacement(output)
	if got.Batch != 512 {
		t.Fatalf("batch = %d, want 512", got.Batch)
	}
	if got.KVGPUMB != 22.75 || got.KVRAMMB != 68.25 {
		t.Fatalf("KV split = GPU %v RAM %v", got.KVGPUMB, got.KVRAMMB)
	}
	if got.ComputeGPUMB != 27.25 || got.ComputeRAMMB != 11.27 {
		t.Fatalf("compute split = GPU %v RAM %v", got.ComputeGPUMB, got.ComputeRAMMB)
	}
	if got.TotalGPUMB() != 117+22.75+27.25 {
		t.Fatalf("total GPU = %v", got.TotalGPUMB())
	}
	if got.TotalRAMMB() != 186+68.25+11.27 {
		t.Fatalf("total RAM = %v", got.TotalRAMMB())
	}
}

func TestParsePlacementLastAttemptWins(t *testing.T) {
	// A failed first attempt and a KV retry both print placement. The health
	// check is answered by the second allocation, so that one must be reported.
	const output = `0.00.642.089 I load_tensors: offloaded 0/33 layers to GPU
0.00.642.092 I load_tensors: CPU_Mapped model buffer size =  3024.67 MiB
0.02.000.000 I load_tensors: offloaded 12/33 layers to GPU
0.02.000.092 I load_tensors:      CUDA0 model buffer size =  1500.00 MiB
0.02.000.094 I load_tensors: CPU_Mapped model buffer size =  1600.00 MiB
0.03.000.000 I srv    load_model: initializing, n_slots = 1, n_ctx_slot = 4096, kv_unified = 'true'
`
	got := ParsePlacement(output)
	if got.GPULayers != 12 || got.GPUMB != 1500 || got.RAMMB != 1600 || got.Context != 4096 {
		t.Fatalf("last attempt not reported: %+v", got)
	}
}

func TestParsePlacementDoesNotInventAPlacementFromUnrelatedOutput(t *testing.T) {
	got := ParsePlacement("model loaded\nlistening on http://127.0.0.1:8080\n")
	if got.Known || got.Complete() {
		t.Fatalf("unrelated output reported a placement: %+v", got)
	}
}
