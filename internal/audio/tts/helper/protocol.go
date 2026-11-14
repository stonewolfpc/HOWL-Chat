package helper

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	ProtocolMajor = 1
	ProtocolMinor = 0

	MaxFrameBytes     = 4 * 1024 * 1024
	MaxTextBytes      = 64 * 1024
	MaxIDBytes        = 128
	MaxFilenameBytes  = 128
	MaxDirectoryBytes = 4096
	MaxErrorDetail    = 16
)

const (
	CapabilityCatalog    = "howl.tts.catalog"
	CapabilityLoad       = "howl.tts.load"
	CapabilityUnload     = "howl.tts.unload"
	CapabilitySynthesize = "howl.tts.synthesize"
	CapabilityCancel     = "howl.tts.cancel"
	CapabilitySpeed      = "howl.tts.speed"
	CapabilityLanguage   = "howl.tts.language"
	CapabilitySSML       = "howl.tts.ssml"
	CapabilityStreaming  = "howl.tts.streaming"
)

var (
	requestIDPattern   = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)
	filenamePattern    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,123}\.wav$`)
	providerPattern    = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)
	enginePattern      = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
	namespacePattern   = regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9_-]*)+$`)
	capabilityVersions = map[string]int{
		CapabilityCatalog: 1, CapabilityLoad: 1, CapabilityUnload: 1,
		CapabilitySynthesize: 1, CapabilityCancel: 1, CapabilitySpeed: 1,
		CapabilityLanguage: 1, CapabilitySSML: 1, CapabilityStreaming: 1,
	}
	mandatoryServiceCapabilities = []string{
		CapabilityCatalog, CapabilityLoad, CapabilityUnload, CapabilitySynthesize, CapabilityCancel,
	}
)

// Version is the exact protocol version object used by every frame.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

func (v Version) valid() bool {
	return v.Major >= 0 && v.Major <= 65535 && v.Minor >= 0 && v.Minor <= 65535
}

func (v Version) String() string { return fmt.Sprintf("%d.%d", v.Major, v.Minor) }

var ProtocolVersion = Version{Major: ProtocolMajor, Minor: ProtocolMinor}

type requestEnvelope struct {
	Protocol   Version         `json:"protocol"`
	ID         string          `json:"id"`
	Op         string          `json:"op"`
	Params     json.RawMessage `json:"params"`
	Extensions json.RawMessage `json:"extensions,omitempty"`
}

type responseEnvelope struct {
	Protocol   Version         `json:"protocol"`
	ID         string          `json:"id"`
	OK         *bool           `json:"ok"`
	Result     json.RawMessage `json:"result,omitempty"`
	Error      *wireError      `json:"error,omitempty"`
	Extensions json.RawMessage `json:"extensions,omitempty"`
}

type wireError struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Retryable bool              `json:"retryable"`
	Field     string            `json:"field,omitempty"`
	Detail    map[string]string `json:"detail,omitempty"`
}

type capability struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type extension struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type providerDescriptor struct {
	ProviderID      string            `json:"provider_id"`
	ProviderVersion string            `json:"provider_version"`
	Engine          string            `json:"engine"`
	DisplayName     string            `json:"display_name"`
	Runtime         map[string]string `json:"runtime"`
	Capabilities    []capability      `json:"capabilities"`
}

type helloParams struct {
	Client        string  `json:"client"`
	ClientVersion string  `json:"client_version"`
	ProtocolMin   Version `json:"protocol_min"`
	ProtocolMax   Version `json:"protocol_max"`
}

type helloResult struct {
	HelperVersion string               `json:"helper_version"`
	Providers     []providerDescriptor `json:"providers"`
	Capabilities  []capability         `json:"capabilities"`
	Extensions    []extension          `json:"extensions"`
	Limits        limits               `json:"limits"`
}

type limits struct {
	MaxFrameBytes int `json:"max_frame_bytes"`
	MaxTextBytes  int `json:"max_text_bytes"`
}

type catalogParams struct {
	ProviderID         string `json:"provider_id"`
	Engine             string `json:"engine"`
	IncludeUnavailable bool   `json:"include_unavailable"`
}

type speedRange struct {
	Supported bool    `json:"supported"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
}

type voiceInfo struct {
	VoiceID     string `json:"voice_id"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
	Detail      string `json:"detail"`
}

type catalogResult struct {
	ProviderID         string      `json:"provider_id"`
	Engine             string      `json:"engine"`
	CatalogVersion     int         `json:"catalog_version"`
	Available          bool        `json:"available"`
	Detail             string      `json:"detail,omitempty"`
	DefaultVoice       string      `json:"default_voice,omitempty"`
	Languages          []string    `json:"languages"`
	Speed              speedRange  `json:"speed"`
	SSMLSupported      bool        `json:"ssml_supported"`
	StreamingSupported bool        `json:"streaming_supported"`
	Voices             []voiceInfo `json:"voices"`
}

type loadAssets struct {
	Root   string `json:"root"`
	Model  string `json:"model,omitempty"`
	Voices string `json:"voices,omitempty"`
	Config string `json:"config,omitempty"`
	Data   string `json:"data,omitempty"`
}

type loadParams struct {
	ProviderID    string     `json:"provider_id"`
	Engine        string     `json:"engine"`
	ModelIdentity string     `json:"model_identity"`
	Assets        loadAssets `json:"assets"`
}

type loadResult struct {
	ProviderID      string            `json:"provider_id"`
	Engine          string            `json:"engine"`
	ProviderVersion string            `json:"provider_version"`
	ModelIdentity   string            `json:"model_identity"`
	Loaded          bool              `json:"loaded"`
	AlreadyLoaded   bool              `json:"already_loaded"`
	Runtime         map[string]string `json:"runtime"`
}

type outputSpec struct {
	Directory string `json:"directory"`
	Filename  string `json:"filename"`
}

type synthesisParams struct {
	ModelIdentity string     `json:"model_identity"`
	Text          string     `json:"text"`
	VoiceID       string     `json:"voice_id"`
	Language      string     `json:"language,omitempty"`
	Speed         float64    `json:"speed"`
	Output        outputSpec `json:"output"`
}

type synthesisResult struct {
	OutputPath      string  `json:"output_path"`
	Format          string  `json:"format"`
	Codec           string  `json:"codec"`
	SampleRateHz    int     `json:"sample_rate_hz"`
	Channels        int     `json:"channels"`
	BitDepth        int     `json:"bit_depth"`
	DurationSeconds float64 `json:"duration_seconds"`
	Bytes           int64   `json:"bytes"`
	SHA256          string  `json:"sha256"`
	ModelIdentity   string  `json:"model_identity"`
}

type healthParams struct{}

type healthProcess struct {
	Alive    bool  `json:"alive"`
	PID      int   `json:"pid"`
	UptimeMS int64 `json:"uptime_ms"`
}

type healthResult struct {
	Process         healthProcess `json:"process"`
	State           string        `json:"state"`
	Ready           bool          `json:"ready"`
	AcceptingWork   bool          `json:"accepting_work"`
	ActiveRequestID *string       `json:"active_request_id"`
	QueueDepth      int           `json:"queue_depth"`
	LastError       *wireError    `json:"last_error"`
}

type cancelParams struct {
	RequestID string `json:"request_id"`
}

type cancelResult struct {
	RequestID       string `json:"request_id"`
	CancelRequested bool   `json:"cancel_requested"`
}

type unloadResult struct {
	Unloaded bool `json:"unloaded"`
}

type shutdownResult struct {
	Accepted bool `json:"accepted"`
}

func newSessionID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func newRequestID(session string, counter uint64) string {
	return fmt.Sprintf("req-%s-%016d", session, counter)
}

func encodeRequestFrame(id, op string, params any) ([]byte, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if !isJSONObject(raw) {
		return nil, fmt.Errorf("params must be a JSON object")
	}
	frame, err := json.Marshal(requestEnvelope{
		Protocol: ProtocolVersion,
		ID:       id,
		Op:       op,
		Params:   raw,
	})
	if err != nil {
		return nil, err
	}
	if len(frame) > MaxFrameBytes {
		return nil, newFailure("frame_too_large", "request frame exceeds 4 MiB")
	}
	return append(frame, '\n'), nil
}

func decodeResponseFrame(line []byte) (responseEnvelope, error) {
	var env responseEnvelope
	if err := prepareFrame(line); err != nil {
		return env, err
	}
	if err := decodeStrict(line, &env); err != nil {
		return env, err
	}
	if !env.Protocol.valid() || env.Protocol != ProtocolVersion {
		return env, newFailure("protocol_violation", "response protocol does not match negotiated v1.0")
	}
	if env.ID != "" && !requestIDPattern.MatchString(env.ID) {
		return env, newFailure("protocol_violation", "response ID is invalid")
	}
	if env.OK == nil {
		return env, newFailure("protocol_violation", "response has no terminal flag")
	}
	if *env.OK {
		if env.Error != nil || !isJSONObject(env.Result) {
			return env, newFailure("protocol_violation", "success response shape is invalid")
		}
	} else {
		if env.Error == nil || len(env.Result) != 0 {
			return env, newFailure("protocol_violation", "error response shape is invalid")
		}
		if err := validateWireError(env.Error); err != nil {
			return env, err
		}
	}
	return env, nil
}

func decodeRequestFrame(line []byte) (requestEnvelope, error) {
	var env requestEnvelope
	if err := prepareFrame(line); err != nil {
		return env, err
	}
	if err := decodeStrict(line, &env); err != nil {
		return env, err
	}
	if !env.Protocol.valid() || env.Protocol != ProtocolVersion {
		return env, newFailure("protocol_violation", "request protocol does not match v1.0")
	}
	if env.ID == "" || !requestIDPattern.MatchString(env.ID) {
		return env, newFailure("invalid_request", "request ID is missing or invalid")
	}
	if !isJSONObject(env.Params) {
		return env, newFailure("invalid_request", "params must be an object").withField("params")
	}
	if len(env.Extensions) != 0 && string(env.Extensions) != "{}" {
		return env, newFailure("unsupported_extension", "extensions are not accepted before negotiation").withField("extensions")
	}
	return env, nil
}

func prepareFrame(line []byte) error {
	if len(line) > MaxFrameBytes {
		return newFailure("frame_too_large", "frame exceeds 4 MiB")
	}
	if len(line) == 0 || bytes.HasPrefix(line, []byte{0xef, 0xbb, 0xbf}) {
		return newFailure("invalid_frame", "frame is empty or has a BOM")
	}
	if bytes.IndexByte(line, 0) >= 0 {
		return newFailure("invalid_frame", "frame contains NUL")
	}
	if !utf8.Valid(line) {
		return newFailure("invalid_frame", "frame is not valid UTF-8")
	}
	if err := validateDuplicateKeys(line); err != nil {
		return newFailure("invalid_frame", "frame has duplicate or invalid JSON keys")
	}
	return nil
}

func decodeStrict(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return newFailure("invalid_request", "strict decode failed").withField("frame")
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return newFailure("invalid_request", "frame has trailing data")
	}
	return nil
}

func decodeParams(raw json.RawMessage, dst any) error {
	if !isJSONObject(raw) {
		return newFailure("invalid_request", "params must be an object").withField("params")
	}
	return decodeStrict(raw, dst)
}

func decodeResult(raw json.RawMessage, dst any) error {
	if !isJSONObject(raw) {
		return newFailure("protocol_violation", "result must be an object")
	}
	return decodeStrict(raw, dst)
}

func isJSONObject(raw []byte) bool {
	raw = bytes.TrimSpace(raw)
	return len(raw) >= 2 && raw[0] == '{' && raw[len(raw)-1] == '}'
}

func validateDuplicateKeys(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := scanJSONValue(dec); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	return nil
}

func scanJSONValue(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := make(map[string]struct{})
		for dec.More() {
			keyTok, err := dec.Token()
			if err != nil {
				return err
			}
			key, ok := keyTok.(string)
			if !ok {
				return fmt.Errorf("object key is not a string")
			}
			if _, duplicate := seen[key]; duplicate {
				return fmt.Errorf("duplicate object key %q", key)
			}
			seen[key] = struct{}{}
			if err := scanJSONValue(dec); err != nil {
				return err
			}
		}
		end, err := dec.Token()
		if err != nil || end != json.Delim('}') {
			return fmt.Errorf("unterminated object")
		}
	case '[':
		for dec.More() {
			if err := scanJSONValue(dec); err != nil {
				return err
			}
		}
		end, err := dec.Token()
		if err != nil || end != json.Delim(']') {
			return fmt.Errorf("unterminated array")
		}
	default:
		return fmt.Errorf("unexpected delimiter %q", delim)
	}
	return nil
}

func validateWireError(e *wireError) error {
	if e == nil {
		return newFailure("protocol_violation", "missing error object")
	}
	if !knownErrorCode(e.Code) {
		return newFailure("protocol_violation", "unknown error code")
	}
	if e.Message == "" || len(e.Message) > 512 || !utf8.ValidString(e.Message) {
		return newFailure("protocol_violation", "error message is invalid")
	}
	if len(e.Field) > 128 {
		return newFailure("protocol_violation", "error field is too long")
	}
	if len(e.Detail) > MaxErrorDetail {
		return newFailure("protocol_violation", "error detail has too many entries")
	}
	for k, v := range e.Detail {
		if len(k) == 0 || len(k) > 64 || len(v) > 256 {
			return newFailure("protocol_violation", "error detail is invalid")
		}
	}
	return nil
}

func knownErrorCode(code string) bool {
	switch code {
	case "invalid_frame", "frame_too_large", "protocol_violation", "invalid_request",
		"unsupported_operation", "protocol_incompatible", "protocol_minor_unsupported",
		"duplicate_request_id", "request_not_found", "request_not_cancellable", "busy",
		"helper_missing", "helper_start_failed", "helper_crashed", "dependency_missing",
		"provider_unknown", "provider_not_loaded", "provider_already_loaded", "model_missing",
		"model_identity_mismatch", "voice_missing", "voice_unsupported", "language_unsupported",
		"speed_unsupported", "unsupported_capability", "unsupported_extension",
		"synthesis_failed", "output_failed", "cancelled", "timeout", "shutdown", "internal_error":
		return true
	default:
		return false
	}
}

func validateHello(h helloResult) error {
	if h.HelperVersion == "" || len(h.HelperVersion) > 64 {
		return newFailure("protocol_violation", "helper version is invalid")
	}
	if h.Limits.MaxFrameBytes != MaxFrameBytes || h.Limits.MaxTextBytes != MaxTextBytes {
		return newFailure("protocol_violation", "helper limits do not match v1")
	}
	seenProviders := make(map[string]struct{})
	for _, p := range h.Providers {
		if !providerPattern.MatchString(p.ProviderID) || !enginePattern.MatchString(p.Engine) ||
			p.ProviderVersion == "" || len(p.ProviderVersion) > 64 || p.DisplayName == "" ||
			len(p.DisplayName) > 128 || p.Runtime == nil {
			return newFailure("protocol_violation", "provider descriptor is invalid")
		}
		if _, ok := seenProviders[p.ProviderID]; ok {
			return newFailure("protocol_violation", "duplicate provider ID")
		}
		seenProviders[p.ProviderID] = struct{}{}
		if err := validateCapabilities(p.Capabilities); err != nil {
			return err
		}
	}
	if err := validateCapabilities(h.Capabilities); err != nil {
		return err
	}
	if missing := missingCapabilities(h.Capabilities, mandatoryServiceCapabilities); len(missing) > 0 {
		f := newFailure("protocol_violation", "mandatory service capabilities are missing")
		f.Detail = map[string]string{"missing": strings.Join(missing, ",")}
		return f
	}
	seenExt := make(map[string]struct{})
	for _, e := range h.Extensions {
		if !namespacePattern.MatchString(e.Name) || len(e.Name) > 64 || e.Version < 1 {
			return newFailure("protocol_violation", "extension descriptor is invalid")
		}
		if _, ok := seenExt[e.Name]; ok {
			return newFailure("protocol_violation", "duplicate extension")
		}
		seenExt[e.Name] = struct{}{}
	}
	return nil
}

func validateCapabilities(caps []capability) error {
	seen := make(map[string]struct{}, len(caps))
	for _, c := range caps {
		if !namespacePattern.MatchString(c.ID) || len(c.ID) > 64 || c.Version < 1 {
			return newFailure("protocol_violation", "capability descriptor is invalid")
		}
		want, known := capabilityVersions[c.ID]
		if !known {
			return newFailure("protocol_violation", "capability is not in the v1 registry")
		}
		if c.Version != want {
			return newFailure("protocol_violation", "capability version is not v1")
		}
		if _, ok := seen[c.ID]; ok {
			return newFailure("protocol_violation", "duplicate capability")
		}
		seen[c.ID] = struct{}{}
	}
	return nil
}

func missingCapabilities(caps []capability, required []string) []string {
	have := make(map[string]struct{}, len(caps))
	for _, c := range caps {
		have[c.ID] = struct{}{}
	}
	var missing []string
	for _, id := range required {
		if _, ok := have[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing
}

func hasCapability(caps []capability, id string) bool {
	for _, c := range caps {
		if c.ID == id {
			return true
		}
	}
	return false
}

func validateCatalog(c catalogResult, requested catalogParams) error {
	if c.ProviderID != requested.ProviderID || c.Engine != requested.Engine {
		return newFailure("protocol_violation", "catalog identity mismatch")
	}
	if c.CatalogVersion != 1 {
		return newFailure("protocol_violation", "unsupported catalog version")
	}
	if !finite(c.Speed.Min) || !finite(c.Speed.Max) || c.Speed.Min > c.Speed.Max {
		return newFailure("protocol_violation", "catalog speed range is invalid")
	}
	if c.DefaultVoice != "" && !catalogHasVoice(c.Voices, c.DefaultVoice) {
		return newFailure("protocol_violation", "catalog default voice is missing")
	}
	seen := make(map[string]struct{}, len(c.Voices))
	for _, v := range c.Voices {
		if v.VoiceID == "" || len(v.VoiceID) > 128 {
			return newFailure("protocol_violation", "catalog voice ID is invalid")
		}
		if _, ok := seen[v.VoiceID]; ok {
			return newFailure("protocol_violation", "duplicate catalog voice")
		}
		seen[v.VoiceID] = struct{}{}
	}
	return nil
}

func catalogHasVoice(voices []voiceInfo, id string) bool {
	for _, v := range voices {
		if v.VoiceID == id {
			return true
		}
	}
	return false
}

func validateLoad(l loadResult, requested loadParams) error {
	if l.ProviderID != requested.ProviderID || l.Engine != requested.Engine ||
		l.ModelIdentity != requested.ModelIdentity || !l.Loaded {
		return newFailure("protocol_violation", "load result identity mismatch")
	}
	if l.Runtime == nil {
		return newFailure("protocol_violation", "load runtime is missing")
	}
	return nil
}

func validateSynthesisResult(s synthesisResult, requested synthesisParams) error {
	if s.Format != "wav" || s.Codec != "pcm_s16le" {
		return newFailure("output_failed", "synthesis format is not v1 WAV")
	}
	if s.Channels != 1 || s.BitDepth != 16 {
		return newFailure("output_failed", "synthesis must be mono 16-bit")
	}
	if s.SampleRateHz < 8000 || s.SampleRateHz > 48000 {
		return newFailure("output_failed", "synthesis sample rate is out of range")
	}
	if !finite(s.DurationSeconds) || s.DurationSeconds <= 0 || s.Bytes <= 0 {
		return newFailure("output_failed", "synthesis audio metadata is invalid")
	}
	if len(s.SHA256) != 64 || strings.ToLower(s.SHA256) != s.SHA256 {
		return newFailure("output_failed", "synthesis hash is invalid")
	}
	if s.ModelIdentity != requested.ModelIdentity {
		return newFailure("model_identity_mismatch", "synthesis identity does not match")
	}
	if !filepath.IsAbs(s.OutputPath) || filepath.Base(s.OutputPath) != requested.Output.Filename ||
		filepath.Clean(filepath.Dir(s.OutputPath)) != filepath.Clean(requested.Output.Directory) {
		return newFailure("output_failed", "synthesis output path is outside the request directory")
	}
	return nil
}

func validateHealth(h healthResult) error {
	switch h.State {
	case "unloaded", "loading", "ready", "synthesizing", "stopping":
	default:
		return newFailure("protocol_violation", "health state is invalid")
	}
	if h.QueueDepth != 0 {
		return newFailure("protocol_violation", "health queue depth must be zero")
	}
	if h.ActiveRequestID != nil {
		if h.State != "synthesizing" || !requestIDPattern.MatchString(*h.ActiveRequestID) {
			return newFailure("protocol_violation", "health active request is invalid")
		}
	}
	if h.LastError != nil {
		return validateWireError(h.LastError)
	}
	return nil
}

func validateText(text string) error {
	if text == "" || !utf8.ValidString(text) || strings.IndexByte(text, 0) >= 0 || len(text) > MaxTextBytes {
		return newFailure("invalid_request", "text is empty, invalid UTF-8, or exceeds 64 KiB").withField("params.text")
	}
	return nil
}

func validateOutput(dir, filename string) error {
	if !filepath.IsAbs(dir) || len(dir) > MaxDirectoryBytes {
		return newFailure("invalid_request", "output directory must be absolute and bounded").withField("params.output.directory")
	}
	if len(filename) > MaxFilenameBytes || !filenamePattern.MatchString(filename) ||
		filepath.Base(filename) != filename || strings.Contains(filename, "..") {
		return newFailure("invalid_request", "output filename is invalid").withField("params.output.filename")
	}
	if reparse, err := hasReparseComponent(filepath.Clean(dir)); err != nil || reparse {
		return newFailure("output_failed", "output directory contains a reparse point or is unavailable").withField("params.output.directory")
	}
	return nil
}

func validateAssets(a loadAssets) error {
	if !filepath.IsAbs(a.Root) || filepath.Clean(a.Root) != a.Root {
		return newFailure("invalid_request", "assets.root must be an absolute cleaned directory").withField("params.assets.root")
	}
	for name, value := range map[string]string{
		"model": a.Model, "voices": a.Voices, "config": a.Config, "data": a.Data,
	} {
		if value == "" {
			continue
		}
		if filepath.IsAbs(value) || filepath.Clean(value) != value || strings.HasPrefix(value, ".."+string(filepath.Separator)) || value == ".." {
			return newFailure("invalid_request", "asset path escapes root").withField("params.assets." + name)
		}
	}
	return nil
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func readFrame(r *bufio.Reader) ([]byte, error) {
	var frame []byte
	for {
		part, err := r.ReadSlice('\n')
		if len(frame)+len(part) > MaxFrameBytes+1 {
			return nil, newFailure("frame_too_large", "frame exceeds 4 MiB")
		}
		frame = append(frame, part...)
		if err == bufio.ErrBufferFull {
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}
	if len(frame) == 0 || frame[len(frame)-1] != '\n' {
		return nil, io.ErrUnexpectedEOF
	}
	frame = frame[:len(frame)-1]
	if len(frame) > MaxFrameBytes {
		return nil, newFailure("frame_too_large", "frame exceeds 4 MiB")
	}
	return frame, nil
}
