// Package fakehelper is a provider-neutral Protocol v1 test helper. It is not
// imported by Howl's audio manager or production TTS paths.
package fakehelper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	protocolMajor = 1
	protocolMinor = 0

	maxFrameBytes = 4 * 1024 * 1024
	maxTextBytes  = 64 * 1024
	maxIDBytes    = 128
	maxFilename   = 128
	maxDirectory  = 4096
	maxDetail     = 16

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
	requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)
	filenamePattern  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,123}\.wav$`)
)

type protocolVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

func (v protocolVersion) valid() bool {
	return v.Major >= 0 && v.Major <= 65535 && v.Minor >= 0 && v.Minor <= 65535
}

func (v protocolVersion) current() bool {
	return v.Major == protocolMajor && v.Minor == protocolMinor
}

var currentProtocol = protocolVersion{Major: protocolMajor, Minor: protocolMinor}

type requestEnvelope struct {
	Protocol   protocolVersion `json:"protocol"`
	ID         string          `json:"id"`
	Op         string          `json:"op"`
	Params     json.RawMessage `json:"params"`
	Extensions json.RawMessage `json:"extensions,omitempty"`
}

type responseEnvelope struct {
	Protocol protocolVersion `json:"protocol"`
	ID       string          `json:"id"`
	OK       bool            `json:"ok"`
	Result   json.RawMessage `json:"result,omitempty"`
	Error    *wireError      `json:"error,omitempty"`
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
	Client        string          `json:"client"`
	ClientVersion string          `json:"client_version"`
	ProtocolMin   protocolVersion `json:"protocol_min"`
	ProtocolMax   protocolVersion `json:"protocol_max"`
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

func errorFrame(id string, e *wireError) ([]byte, error) {
	return json.Marshal(responseEnvelope{
		Protocol: currentProtocol,
		ID:       id,
		OK:       false,
		Error:    e,
	})
}

func successFrame(id string, result any) ([]byte, error) {
	raw, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return json.Marshal(responseEnvelope{
		Protocol: currentProtocol,
		ID:       id,
		OK:       true,
		Result:   raw,
	})
}

func newWireError(code, message string) *wireError {
	return &wireError{Code: code, Message: message}
}

func decodeParams(raw json.RawMessage, dst any) *wireError {
	if !isObject(raw) {
		return newWireError("invalid_request", "params must be an object").withField("params")
	}
	return decodeStrict(raw, dst)
}

func (e *wireError) withField(field string) *wireError {
	if e != nil {
		e.Field = field
	}
	return e
}

func isObject(raw []byte) bool {
	raw = bytes.TrimSpace(raw)
	return len(raw) >= 2 && raw[0] == '{' && raw[len(raw)-1] == '}'
}

func decodeStrict(data []byte, dst any) *wireError {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return newWireError("invalid_request", "strict decode failed").withField("frame")
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return newWireError("invalid_request", "frame has trailing data")
	}
	return nil
}

// decodeRequestFrame validates one line and returns the decoded envelope plus a
// non-fatal wire error. fatal means the session must send one frame-level error
// and exit 4; a non-fatal error is returned to the request ID.
func decodeRequestFrame(line []byte) (requestEnvelope, *wireError, bool) {
	var env requestEnvelope
	if e, fatal := prepareFrame(line); e != nil {
		return env, e, fatal
	}
	if err := validateDuplicateKeys(line); err != nil {
		return env, newWireError("invalid_frame", "frame has duplicate or invalid JSON keys"), true
	}
	if err := decodeStrict(line, &env); err != nil {
		// Preserve a usable ID for the non-fatal unknown-core-field case. A
		// malformed envelope without one is fatal.
		id := rawID(line)
		if id != "" && requestIDPattern.MatchString(id) {
			env.ID = id
			return env, newWireError("invalid_request", "strict envelope decode failed").withField("frame"), false
		}
		return env, newWireError("invalid_frame", "frame does not have a valid request envelope"), true
	}
	if !env.Protocol.valid() || !env.Protocol.current() {
		return env, newWireError("protocol_incompatible", "request protocol does not match v1.0"), true
	}
	if env.ID == "" || !requestIDPattern.MatchString(env.ID) || len(env.ID) > maxIDBytes {
		return env, newWireError("invalid_request", "request ID is missing or invalid"), env.ID == ""
	}
	if !isObject(env.Params) {
		return env, newWireError("invalid_request", "params must be an object").withField("params"), false
	}
	if len(env.Extensions) > 0 {
		trimmed := bytes.TrimSpace(env.Extensions)
		if !isObject(trimmed) {
			return env, newWireError("invalid_request", "extensions must be an object").withField("extensions"), false
		}
		if string(trimmed) != "{}" {
			return env, newWireError("unsupported_extension", "extensions are not advertised"), false
		}
	}
	return env, nil, false
}

func prepareFrame(line []byte) (*wireError, bool) {
	if len(line) > maxFrameBytes {
		return newWireError("frame_too_large", "frame exceeds 4 MiB"), true
	}
	if len(line) == 0 || bytes.HasPrefix(line, []byte{0xef, 0xbb, 0xbf}) {
		return newWireError("invalid_frame", "frame is empty or has a BOM"), true
	}
	if bytes.IndexByte(line, 0) >= 0 {
		return newWireError("invalid_frame", "frame contains NUL"), true
	}
	if !utf8.Valid(line) {
		return newWireError("invalid_frame", "frame is not valid UTF-8"), true
	}
	return nil, false
}

func rawID(line []byte) string {
	var raw map[string]json.RawMessage
	if json.Unmarshal(line, &raw) != nil {
		return ""
	}
	var id string
	if json.Unmarshal(raw["id"], &id) != nil {
		return ""
	}
	return id
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

func readFrame(r *bufio.Reader) ([]byte, *wireError) {
	var frame []byte
	for {
		part, err := r.ReadSlice('\n')
		if len(frame)+len(part) > maxFrameBytes+1 {
			return nil, newWireError("frame_too_large", "frame exceeds 4 MiB")
		}
		frame = append(frame, part...)
		if err == bufio.ErrBufferFull {
			continue
		}
		if err != nil {
			if err == io.EOF && len(frame) == 0 {
				return nil, nil
			}
			return nil, newWireError("invalid_frame", "frame is incomplete")
		}
		break
	}
	if len(frame) == 0 || frame[len(frame)-1] != '\n' {
		return nil, newWireError("invalid_frame", "frame must end with LF")
	}
	frame = frame[:len(frame)-1]
	if len(frame) > maxFrameBytes {
		return nil, newWireError("frame_too_large", "frame exceeds 4 MiB")
	}
	return frame, nil
}

func validateText(text string) *wireError {
	if text == "" || !utf8.ValidString(text) || strings.IndexByte(text, 0) >= 0 || len(text) > maxTextBytes {
		return newWireError("invalid_request", "text is empty, invalid UTF-8, or exceeds 64 KiB").withField("params.text")
	}
	return nil
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
