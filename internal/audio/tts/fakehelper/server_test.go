package fakehelper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testRequestFrame(t *testing.T, id, op string, params any) string {
	t.Helper()
	raw, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	frame, err := json.Marshal(requestEnvelope{
		Protocol: currentProtocol, ID: id, Op: op, Params: raw,
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(frame) + "\n"
}

func TestRunHelloHealthAndShutdown(t *testing.T) {
	var out, diag bytes.Buffer
	input := strings.NewReader(
		testRequestFrame(t, "req-1", "hello", helloParams{
			Client: "howl-chat", ClientVersion: "test",
			ProtocolMin: currentProtocol, ProtocolMax: currentProtocol,
		}) +
			testRequestFrame(t, "req-2", "health", healthParams{}) +
			testRequestFrame(t, "req-3", "shutdown", struct{}{}),
	)
	if code := Run(input, &out, &diag, modeNormal, ""); code != 0 {
		t.Fatalf("Run exit = %d, want 0; diag=%q", code, diag.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("response lines = %d, want 3: %q", len(lines), out.String())
	}
	for i, line := range lines {
		var env responseEnvelope
		if err := json.Unmarshal([]byte(line), &env); err != nil {
			t.Fatalf("line %d is not JSON: %v", i, err)
		}
		if !env.OK || env.Error != nil || len(env.Result) == 0 {
			t.Fatalf("line %d is not a success terminal: %+v", i, env)
		}
	}
	if diag.Len() != 0 {
		t.Fatalf("stderr diagnostics = %q, want empty", diag.String())
	}
}

func TestRunRejectsDuplicateKeysFatal(t *testing.T) {
	var out, diag bytes.Buffer
	line := `{"protocol":{"major":1,"minor":0},"id":"req-1","op":"hello","params":{},"params":{}}`
	code := Run(strings.NewReader(line+"\n"), &out, &diag, modeNormal, "")
	if code != 4 {
		t.Fatalf("Run exit = %d, want 4", code)
	}
	var env responseEnvelope
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &env); err != nil {
		t.Fatal(err)
	}
	if env.Error == nil || env.Error.Code != "invalid_frame" || env.ID != "" {
		t.Fatalf("duplicate-key response = %+v", env)
	}
}

func TestRunRejectsUnknownFieldAndUnknownOperation(t *testing.T) {
	var out, diag bytes.Buffer
	input := strings.NewReader(
		testRequestFrame(t, "req-1", "hello", helloParams{
			Client: "howl-chat", ClientVersion: "test",
			ProtocolMin: currentProtocol, ProtocolMax: currentProtocol,
		}) +
			`{"protocol":{"major":1,"minor":0},"id":"req-2","op":"health","params":{"unexpected":true}}` + "\n" +
			`{"protocol":{"major":1,"minor":0},"id":"req-3","op":"not_an_operation","params":{}}` + "\n" +
			testRequestFrame(t, "req-4", "shutdown", struct{}{}),
	)
	if code := Run(input, &out, &diag, modeNormal, ""); code != 0 {
		t.Fatalf("Run exit = %d, want 0; diag=%q", code, diag.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("responses = %d, want 4: %q", len(lines), out.String())
	}
	var unknownField, unknownOp responseEnvelope
	if err := json.Unmarshal([]byte(lines[1]), &unknownField); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(lines[2]), &unknownOp); err != nil {
		t.Fatal(err)
	}
	if unknownField.Error == nil || unknownField.Error.Code != "invalid_request" {
		t.Fatalf("unknown field response = %+v", unknownField)
	}
	if unknownOp.Error == nil || unknownOp.Error.Code != "unsupported_operation" {
		t.Fatalf("unknown operation response = %+v", unknownOp)
	}
}

func TestRunRejectsInvalidProtocolVersion(t *testing.T) {
	var out, diag bytes.Buffer
	line := `{"protocol":{"major":2,"minor":0},"id":"req-1","op":"hello","params":{}}`
	code := Run(strings.NewReader(line+"\n"), &out, &diag, modeNormal, "")
	if code != 4 {
		t.Fatalf("Run exit = %d, want 4", code)
	}
	var env responseEnvelope
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &env); err != nil {
		t.Fatal(err)
	}
	if env.ID != "req-1" || env.Error == nil || env.Error.Code != "protocol_incompatible" {
		t.Fatalf("invalid-version response = %+v", env)
	}
}

func TestRunMalformedFrameExitsFourWithFrameError(t *testing.T) {
	var out, diag bytes.Buffer
	code := Run(strings.NewReader("not a frame\n"), &out, &diag, modeNormal, "")
	if code != 4 {
		t.Fatalf("Run exit = %d, want 4", code)
	}
	var env responseEnvelope
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &env); err != nil {
		t.Fatalf("error frame: %v", err)
	}
	if env.ID != "" || env.OK || env.Error == nil || env.Error.Code != "invalid_frame" {
		t.Fatalf("error frame = %+v", env)
	}
}

func TestServerRejectsSecondNonControlOperation(t *testing.T) {
	s := &server{
		workerCh: make(chan workerRequest, 3),
		stopCh:   make(chan struct{}),
	}
	defer s.stop()
	if !s.enqueue(requestEnvelope{ID: "req-1", Op: "synthesize"}) {
		t.Fatal("first non-control operation was rejected")
	}
	if s.enqueue(requestEnvelope{ID: "req-2", Op: "synthesize"}) {
		t.Fatal("second non-control operation was accepted while busy")
	}
	s.finishWork()
	if !s.enqueue(requestEnvelope{ID: "req-3", Op: "catalog"}) {
		t.Fatal("helper did not accept work after the active operation finished")
	}
}

func TestSyntheticWAVIsCanonicalAndDeterministic(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.wav")
	second := filepath.Join(dir, "second.wav")
	hash1, size1, duration1, err := writeSyntheticWAV(first)
	if err != nil {
		t.Fatal(err)
	}
	hash2, size2, duration2, err := writeSyntheticWAV(second)
	if err != nil {
		t.Fatal(err)
	}
	if hash1 != hash2 || size1 != size2 || duration1 != duration2 {
		t.Fatalf("synthetic output is not deterministic: %s/%d/%v vs %s/%d/%v", hash1, size1, duration1, hash2, size2, duration2)
	}
	if size1 != 44+syntheticSamples*2 || duration1 != float64(syntheticSamples)/syntheticRate {
		t.Fatalf("synthetic metadata = size %d duration %v", size1, duration1)
	}
	if err := validateSyntheticWAV(first, hash1, size1); err != nil {
		t.Fatal(err)
	}
}

func TestRunAdvertisesFrozenCapabilities(t *testing.T) {
	var out, diag bytes.Buffer
	input := strings.NewReader(testRequestFrame(t, "req-1", "hello", helloParams{
		Client: "howl-chat", ClientVersion: "test",
		ProtocolMin: currentProtocol, ProtocolMax: currentProtocol,
	}))
	if code := Run(input, &out, &diag, modeNormal, ""); code != 0 {
		t.Fatalf("Run exit = %d, want 0; diag=%q", code, diag.String())
	}
	var env responseEnvelope
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &env); err != nil {
		t.Fatal(err)
	}
	var result helloResult
	if err := json.Unmarshal(env.Result, &result); err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{
		CapabilityCatalog: true, CapabilityLoad: true, CapabilityUnload: true,
		CapabilitySynthesize: true, CapabilityCancel: true,
	}
	if len(result.Capabilities) != len(want) {
		t.Fatalf("service capabilities = %+v, want five registry capabilities", result.Capabilities)
	}
	for _, cap := range result.Capabilities {
		if !want[cap.ID] || cap.Version != 1 {
			t.Fatalf("unexpected service capability = %+v", cap)
		}
	}
	if len(result.Providers) != 1 || !hasCapability(result.Providers[0].Capabilities, CapabilitySynthesize) {
		t.Fatalf("provider capabilities = %+v, want synthesize", result.Providers)
	}
}

func TestHandleSynthesizeRejectsUnadvertisedOptionalCapability(t *testing.T) {
	var out bytes.Buffer
	s := &server{
		out:           bufio.NewWriter(&out),
		mode:          modePartialCapabilities,
		loaded:        true,
		modelIdentity: "fake:model:1",
	}
	params, err := json.Marshal(synthesisParams{
		ModelIdentity: "fake:model:1", Text: "hello", VoiceID: fakeVoice, Speed: 1.25,
		Output: outputSpec{Directory: t.TempDir(), Filename: "out.wav"},
	})
	if err != nil {
		t.Fatal(err)
	}
	s.handleSynthesize(requestEnvelope{
		Protocol: currentProtocol, ID: "req-1", Op: "synthesize", Params: params,
	}, make(chan struct{}))
	if err := s.out.Flush(); err != nil {
		t.Fatal(err)
	}
	var env responseEnvelope
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &env); err != nil {
		t.Fatal(err)
	}
	if env.Error == nil || env.Error.Code != "unsupported_capability" || env.Error.Detail["capability"] != CapabilitySpeed {
		t.Fatalf("optional-capability response = %+v", env)
	}
}

func TestPublishNoReplaceRefusesExistingFinal(t *testing.T) {
	dir := t.TempDir()
	part := filepath.Join(dir, "out.wav.part-test")
	final := filepath.Join(dir, "out.wav")
	if err := os.WriteFile(part, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(final, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := publishNoReplace(part, final); err == nil {
		t.Fatal("publication replaced an existing final file")
	}
	raw, err := os.ReadFile(final)
	if err != nil || string(raw) != "existing" {
		t.Fatalf("existing final changed: %q, %v", raw, err)
	}
	if err := os.Remove(final); err != nil {
		t.Fatal(err)
	}
	if err := publishNoReplace(part, final); err != nil {
		t.Fatal(err)
	}
	raw, err = os.ReadFile(final)
	if err != nil || string(raw) != "new" {
		t.Fatalf("published final = %q, %v", raw, err)
	}
	if _, err := os.Stat(part); !os.IsNotExist(err) {
		t.Fatalf("part file remains after publish: %v", err)
	}
}

func TestValidateOutputPathRejectsTraversalAndCollision(t *testing.T) {
	dir := t.TempDir()
	if e := validateOutputPath(dir, "../escape.wav"); e == nil || e.Code != "invalid_request" {
		t.Fatalf("traversal error = %+v", e)
	}
	if e := validateOutputPath(dir, "ok.wav"); e != nil {
		t.Fatalf("valid output rejected: %+v", e)
	}
	if err := os.WriteFile(filepath.Join(dir, "ok.wav"), []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if e := validateOutputPath(dir, "ok.wav"); e == nil || e.Code != "output_failed" {
		t.Fatalf("collision error = %+v", e)
	}
}
