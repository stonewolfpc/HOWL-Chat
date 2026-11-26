package main

import (
	llamabackend "howl-chat/internal/backend/llama"
)

// LlamaCapabilityInfo is the renderer-facing view of what the shipped
// llama-server actually advertises. The settings UI disables an option with
// this reason instead of offering a control whose flag the binary does not
// understand.
type LlamaCapabilityInfo struct {
	Binary   string          `json:"binary"`
	Version  string          `json:"version"`
	Build    string          `json:"build"`
	Commit   string          `json:"commit"`
	Probed   bool            `json:"probed"`
	Error    string          `json:"error,omitempty"`
	Supports map[string]bool `json:"supports"`
}

// GetLlamaCapabilities probes the llama-server binary in use and reports the
// controls it can honour. Sentinels correspond to the chat-side switches, not
// to a second feature list: each one is answered by a real --help flag.
func (a *App) GetLlamaCapabilities() LlamaCapabilityInfo {
	binary := findLlamaServerForRuntime()
	caps := llamabackend.ProbeCLI(binary)
	return LlamaCapabilityInfo{
		Binary:  binary,
		Version: caps.Version,
		Build:   caps.Build,
		Commit:  caps.Commit,
		Probed:  caps.Probed,
		Error:   caps.Error,
		Supports: map[string]bool{
			"flash_attn":    caps.Supports("--flash-attn"),
			"cache_type_k":  caps.Supports("--cache-type-k"),
			"cache_type_v":  caps.Supports("--cache-type-v"),
			"kv_offload":    caps.Supports("--kv-offload"),
			"no_kv_offload": caps.Supports("--no-kv-offload"),
			"load_mode":     caps.Supports("--load-mode"),
			"fit":           caps.Supports("--fit"),
			"main_gpu":      caps.Supports("--main-gpu"),
		},
	}
}
