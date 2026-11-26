package llama

import "testing"

func TestParseHelpFlagsReadsAdvertisedLongFlags(t *testing.T) {
	help := "-ctk,  --cache-type-k TYPE              KV cache data type for K\n" +
		"-fa,   --flash-attn [on|off|auto]       set Flash Attention use\n" +
		"-fit,  --fit [on|off]                   whether to adjust unset arguments\n" +
		"-kvo,  --kv-offload, -nkvo, --no-kv-offload\n" +
		"                                        keep the KV cache in RAM\n" +
		"-lm,   --load-mode MODE                 model loading mode\n" +
		"-mg,   --main-gpu INDEX                 the GPU to use\n"

	flags := parseHelpFlags(help)
	for _, want := range []string{"--cache-type-k", "--flash-attn", "--fit", "--kv-offload", "--no-kv-offload", "--load-mode", "--main-gpu"} {
		if !flags[want] {
			t.Fatalf("help parser missed %s: %v", want, flags)
		}
	}
	if flags["--not-a-flag"] {
		t.Fatal("help parser invented a flag")
	}
}

func TestUnprobedCapabilitySupportsNothingOptional(t *testing.T) {
	caps := CapabilitySet{Probed: false, Flags: map[string]bool{"--flash-attn": true}}
	if caps.Supports("--flash-attn") {
		t.Fatal("an unprobed binary must not be treated as supporting a flag")
	}
	caps.Probed = true
	if !caps.Supports("--flash-attn") {
		t.Fatal("a probed binary must expose its advertised flags")
	}
}

func TestParseVersion(t *testing.T) {
	version, build, commit := parseVersion("version: 0.8.0-dev (build 12015, commit a34de8618)\nbuilt with Clang")
	if version != "0.8.0-dev" || build != "12015" || commit != "a34de8618" {
		t.Fatalf("version parse = %q/%q/%q", version, build, commit)
	}
}
