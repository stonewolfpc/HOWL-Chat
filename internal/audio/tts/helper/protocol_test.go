package helper

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDuplicateKeysAreFatal(t *testing.T) {
	line := []byte(`{"protocol":{"major":1,"minor":0},"id":"req-0123456789abcdef-0000000000000001","ok":true,"result":{},"result":{}}`)
	_, err := decodeResponseFrame(line)
	f, ok := AsFailure(err)
	if !ok || f.Code != "invalid_frame" {
		t.Fatalf("error = %v, want invalid_frame", err)
	}
}

func TestUnknownResponseFieldIsRejected(t *testing.T) {
	line := []byte(`{"protocol":{"major":1,"minor":0},"id":"req-0123456789abcdef-0000000000000001","ok":true,"result":{},"extra":true}`)
	_, err := decodeResponseFrame(line)
	f, ok := AsFailure(err)
	if !ok || f.Code != "invalid_request" {
		t.Fatalf("error = %v, want invalid_request", err)
	}
}

func TestOversizedFrameIsRejected(t *testing.T) {
	line := bytes.Repeat([]byte("x"), MaxFrameBytes+1)
	if err := prepareFrame(line); err == nil {
		t.Fatal("oversized frame was accepted")
	} else if f, ok := AsFailure(err); !ok || f.Code != "frame_too_large" {
		t.Fatalf("error = %v, want frame_too_large", err)
	}
}

func TestRequestFrameCarriesExactVersionAndNoExtensions(t *testing.T) {
	frame, err := encodeRequestFrame("req-0123456789abcdef-0000000000000001", "health", healthParams{})
	if err != nil {
		t.Fatal(err)
	}
	var env map[string]any
	if err := json.Unmarshal(bytes.TrimSuffix(frame, []byte("\n")), &env); err != nil {
		t.Fatal(err)
	}
	if env["op"] != "health" {
		t.Fatalf("op = %v", env["op"])
	}
	if _, ok := env["extensions"]; ok {
		t.Fatal("health request unexpectedly carried extensions")
	}
}

func TestRetiredResponseIsDiscarded(t *testing.T) {
	ok := true
	c := &Client{
		pending: make(map[string]*pendingRequest),
		retired: map[string]struct{}{"req-retired": {}},
	}
	if !c.dispatch(responseEnvelope{Protocol: ProtocolVersion, ID: "req-retired", OK: &ok, Result: json.RawMessage(`{}`)}) {
		t.Fatal("retired response terminated the session")
	}
}

func TestDuplicateTerminalIsProtocolViolation(t *testing.T) {
	ok := true
	ch := make(chan pendingOutcome, 1)
	c := &Client{
		pending: map[string]*pendingRequest{
			"req-active": {op: "health", ch: ch, terminal: true},
		},
		retired: make(map[string]struct{}),
	}
	if c.dispatch(responseEnvelope{Protocol: ProtocolVersion, ID: "req-active", OK: &ok, Result: json.RawMessage(`{}`)}) {
		t.Fatal("duplicate terminal was accepted")
	}
	if c.readFailure == nil || c.readFailure.Code != "protocol_violation" {
		t.Fatalf("failure = %+v, want protocol_violation", c.readFailure)
	}
}

func TestValidateHelloRejectsInvalidCapabilitySets(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*helloResult)
	}{
		{
			name: "unregistered service capability",
			mutate: func(h *helloResult) {
				h.Capabilities = append(h.Capabilities, capability{ID: "protocol.hello", Version: 1})
			},
		},
		{
			name: "unregistered provider capability",
			mutate: func(h *helloResult) {
				h.Providers[0].Capabilities = append(h.Providers[0].Capabilities, capability{ID: "synthesis.file", Version: 1})
			},
		},
		{
			name: "wrong capability version",
			mutate: func(h *helloResult) {
				h.Capabilities[0].Version = 2
			},
		},
		{
			name: "missing mandatory service capability",
			mutate: func(h *helloResult) {
				filtered := h.Capabilities[:0]
				for _, c := range h.Capabilities {
					if c.ID != CapabilityUnload {
						filtered = append(filtered, c)
					}
				}
				h.Capabilities = filtered
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := defaultHello()
			tt.mutate(&h)
			err := validateHello(h)
			f, ok := AsFailure(err)
			if !ok || f.Code != "protocol_violation" {
				t.Fatalf("error = %v, want protocol_violation", err)
			}
		})
	}
}

func TestValidateHelloNamesMissingCoreCapabilities(t *testing.T) {
	h := defaultHello()
	h.Capabilities = []capability{{ID: CapabilityCatalog, Version: 1}}
	err := validateHello(h)
	f, ok := AsFailure(err)
	if !ok || f.Code != "protocol_violation" {
		t.Fatalf("error = %v, want protocol_violation", err)
	}
	if f.Detail["missing"] == "" {
		t.Fatalf("missing capability detail = %+v", f.Detail)
	}
}

func TestSynthesisValidationRejectsEscapingOutput(t *testing.T) {
	params := synthesisParams{
		ModelIdentity: "id", Text: "hello", VoiceID: "voice", Speed: 1,
		Output: outputSpec{Directory: `C:\allowed`, Filename: "out.wav"},
	}
	result := synthesisResult{
		OutputPath: `C:\elsewhere\out.wav`, Format: "wav", Codec: "pcm_s16le",
		SampleRateHz: 24000, Channels: 1, BitDepth: 16, DurationSeconds: 1,
		Bytes: 100, SHA256: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		ModelIdentity: "id",
	}
	err := validateSynthesisResult(result, params)
	if f, ok := AsFailure(err); !ok || f.Code != "output_failed" {
		t.Fatalf("error = %v, want output_failed", err)
	}
}
