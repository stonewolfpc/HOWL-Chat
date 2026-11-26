package main

import (
	"strings"
	"testing"

	llamabackend "howl-chat/internal/backend/llama"
)

// llamaCapsForTests is the flag set this machine's shipped binary advertises,
// captured from its own --help. Tests use an explicit set so a different build
// does not change what the argument tests mean.
func llamaCapsForTests() llamabackend.CapabilitySet {
	flags := map[string]bool{}
	for _, flag := range []string{
		"--fit", "--flash-attn", "--cache-type-k", "--cache-type-v",
		"--kv-offload", "--no-kv-offload", "--load-mode", "--main-gpu",
		"--vocab-only", "--mmproj", "--log-verbosity",
	} {
		flags[flag] = true
	}
	return llamabackend.CapabilitySet{Probed: true, Flags: flags}
}

func joinedLlamaArgs(settings map[string]interface{}) string {
	return strings.Join(buildLlamaServerArgs("model.gguf", "", settings, false, llamaCapsForTests()), " ")
}

func TestLlamaArgsRequestPlacementEvidence(t *testing.T) {
	got := joinedLlamaArgs(map[string]interface{}{})
	if !strings.Contains(got, "--log-verbosity 4") {
		t.Fatalf("placement evidence requires verbosity 4 in the spawn: %q", got)
	}
	// An unprobed/older binary must not receive a flag its help did not list.
	legacy := strings.Join(buildLlamaServerArgs("model.gguf", "", map[string]interface{}{}, false,
		llamabackend.CapabilitySet{Probed: true, Flags: map[string]bool{"--fit": true}}), " ")
	if strings.Contains(legacy, "--log-verbosity") {
		t.Fatalf("unsupported verbosity flag reached an older binary: %q", legacy)
	}
}

func TestLlamaArgsCarryKVOffloadAndFlashAttentionControls(t *testing.T) {
	settings := map[string]interface{}{
		"context_size":         4096,
		"batch_size":           256,
		"gpu_layers":           12,
		"flash_attention_mode": "off",
		"kv_cache_type":        "q8_0",
		"offload_kqv":          false,
		"use_mmap":             true,
		"use_mlock":            false,
	}
	got := joinedLlamaArgs(settings)

	for _, want := range []string{
		"--flash-attn off",
		"--cache-type-k q8_0",
		"--cache-type-v q8_0",
		"--no-kv-offload",
		"--load-mode mmap",
		"--fit off",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("spawned args %q do not contain %q", got, want)
		}
	}
	if strings.Contains(got, "--kv-offload --no-kv-offload") {
		t.Fatalf("offload directions must be mutually exclusive: %q", got)
	}
}

func TestLlamaArgsFlashAttentionIsTriState(t *testing.T) {
	cases := map[string]string{
		"on":   "--flash-attn on",
		"off":  "--flash-attn off",
		"auto": "--flash-attn auto",
	}
	for mode, want := range cases {
		got := joinedLlamaArgs(map[string]interface{}{"flash_attention_mode": mode})
		if !strings.Contains(got, want) {
			t.Fatalf("mode %q produced %q, want %q", mode, got, want)
		}
	}
}

func TestLlamaArgsLegacyFlashBooleanKeepsOldMeaning(t *testing.T) {
	if got := joinedLlamaArgs(map[string]interface{}{"flash_attention": true}); !strings.Contains(got, "--flash-attn on") {
		t.Fatalf("legacy true should mean on: %q", got)
	}
	// The old false value sent nothing, which the engine treated as auto. It is
	// migrated to an explicit auto choice rather than silently becoming off.
	if got := joinedLlamaArgs(map[string]interface{}{"flash_attention": false}); !strings.Contains(got, "--flash-attn auto") {
		t.Fatalf("legacy false should mean auto: %q", got)
	}
}

func TestLlamaArgsKVLocationSendsOneDirection(t *testing.T) {
	onGPU := joinedLlamaArgs(map[string]interface{}{"offload_kqv": true})
	if !strings.Contains(onGPU, "--kv-offload") || strings.Contains(onGPU, "--no-kv-offload") {
		t.Fatalf("offload_kqv=true must send only --kv-offload: %q", onGPU)
	}
	inRAM := joinedLlamaArgs(map[string]interface{}{"offload_kqv": false})
	if !strings.Contains(inRAM, "--no-kv-offload") || strings.Contains(inRAM, "--kv-offload ") {
		t.Fatalf("offload_kqv=false must send only --no-kv-offload: %q", inRAM)
	}
}

func TestResolveKVTypePassesQuantizedThroughOnAndAuto(t *testing.T) {
	for _, cache := range []string{"q4_0", "q8_0"} {
		for _, mode := range []string{"on", "auto"} {
			kvType, warning := resolveKVType(map[string]interface{}{
				"kv_cache_type":        cache,
				"flash_attention_mode": mode,
			})
			if kvType != cache || warning != "" {
				t.Fatalf("%s with Flash Attention %q must be sent unchanged so the engine can try it, got %q / %q", cache, mode, kvType, warning)
			}
		}
		for _, mode := range []string{"off", ""} {
			kvType, warning := resolveKVType(map[string]interface{}{
				"kv_cache_type":        cache,
				"flash_attention_mode": mode,
			})
			if kvType != "f16" || !strings.Contains(warning, "requires Flash Attention") {
				t.Fatalf("%s with Flash Attention %q must downgrade to f16 with a reason, got %q / %q", cache, mode, kvType, warning)
			}
		}
	}

	kvType, warning := resolveKVType(map[string]interface{}{
		"kv_cache_type":        "f16",
		"flash_attention_mode": "auto",
	})
	if kvType != "f16" || warning != "" {
		t.Fatalf("f16 with Flash Attention auto must be sent unchanged, got %q / %q", kvType, warning)
	}
}

func TestQuantizedKVRejectionRetriesOnceWithF16(t *testing.T) {
	// Captured verbatim from runtime/llama/llama-server.exe on this machine:
	// --flash-attn off --cache-type-k q4_0 --cache-type-v q4_0
	const injected = "0.00.572.074 E llama_init_from_model: quantized V cache requires flash_attn to be enabled"
	settings := map[string]interface{}{
		"kv_cache_type":        "q4_0",
		"flash_attention_mode": "auto",
	}
	before := strings.Join(buildLlamaServerArgs("model.gguf", "", settings, false, llamaCapsForTests()), " ")
	t.Logf("injected engine output: %s", injected)
	t.Logf("first spawn args: %s", before)
	if !strings.Contains(before, "--cache-type-v q4_0") {
		t.Fatalf("the requested q4_0 cache must reach the first spawn: %q", before)
	}

	retries := 0
	requested, warning, ok := prepareQuantizedKVRetry(injected, settings, &retries)
	if !ok || requested != "q4_0" {
		t.Fatalf("injected engine rejection must trigger one retry, got ok=%t requested=%q", ok, requested)
	}
	if !strings.Contains(warning, "q4_0") || !strings.Contains(warning, "f16") || !strings.Contains(warning, "Flash Attention was not enabled") {
		t.Fatalf("warning must name the request, the replacement, and the reason: %q", warning)
	}

	after := strings.Join(buildLlamaServerArgs("model.gguf", "", settings, false, llamaCapsForTests()), " ")
	t.Logf("warning: %s", warning)
	t.Logf("retry spawn args: %s", after)
	if !strings.Contains(after, "--cache-type-k f16") || !strings.Contains(after, "--cache-type-v f16") || strings.Contains(after, "q4_0") {
		t.Fatalf("retry spawn must carry f16 K and V and no quantized cache: %q", after)
	}
	wantAfter := strings.ReplaceAll(before, "--cache-type-k q4_0 --cache-type-v q4_0", "--cache-type-k f16 --cache-type-v f16")
	if after != wantAfter {
		t.Fatalf("retry must change only the KV cache types; before=%q after=%q want=%q", before, after, wantAfter)
	}

	_, _, again := prepareQuantizedKVRetry(injected, settings, &retries)
	t.Logf("second identical request retried: %t", again)
	if again {
		t.Fatal("the same request must never be retried a second time")
	}
	_, _, retried := prepareQuantizedKVRetry(injected, map[string]interface{}{"kv_cache_type": "f16"}, new(int))
	t.Logf("f16 request retried: %t", retried)
	if retried {
		t.Fatal("an f16 request must not be retried for a quantized-cache rejection")
	}
}

func TestLlamaArgsTensorSplitSliderReachesEngine(t *testing.T) {
	got := joinedLlamaArgs(map[string]interface{}{"tensor_split": float64(30)})
	if !strings.Contains(got, "--tensor-split 30,70") {
		t.Fatalf("numeric tensor split must reach the engine as two proportions: %q", got)
	}
	got = joinedLlamaArgs(map[string]interface{}{"tensor_split": "3,1"})
	if !strings.Contains(got, "--tensor-split 3,1") {
		t.Fatalf("explicit tensor split string must pass through: %q", got)
	}
}

func TestLlamaArgsRejectUnsupportedOrRemovedKVTypes(t *testing.T) {
	for _, value := range []string{"paged", "flash-attention", "not-a-type"} {
		got := joinedLlamaArgs(map[string]interface{}{"kv_cache_type": value})
		if strings.Contains(got, "--cache-type-k") || strings.Contains(got, "--cache-type-v") {
			t.Fatalf("removed/unknown KV type %q must not reach the engine: %q", value, got)
		}
	}
}

func TestLlamaArgsOmitOptionalFlagsTheBinaryDoesNotAdvertise(t *testing.T) {
	caps := llamabackend.CapabilitySet{Probed: true, Flags: map[string]bool{}}
	settings := map[string]interface{}{
		"flash_attention_mode": "off",
		"kv_cache_type":        "q8_0",
		"offload_kqv":          false,
		"use_mmap":             true,
		"main_gpu":             1,
		"vocab_only":           true,
	}
	got := strings.Join(buildLlamaServerArgs("model.gguf", "", settings, false, caps), " ")
	for _, unwanted := range []string{
		"--flash-attn", "--cache-type-k", "--cache-type-v",
		"--kv-offload", "--no-kv-offload", "--load-mode",
		"--main-gpu", "--vocab-only", "--fit off",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("unadvertised flag %q reached the command line: %q", unwanted, got)
		}
	}
}
