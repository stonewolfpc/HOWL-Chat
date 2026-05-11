package melotts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/audio/tts"
	"howl-chat/internal/audio/types"
)

// Client talks to an HTTP MeloTTS deployment.
//
// Official MeloTTS ships a Gradio WebUI first party; deployments often expose a small HTTP layer
// (docs and community setups use POST /synthesize with JSON plus OpenAI-compatible /v1/audio/speech).
// This client attempts /synthesize first, then falls back to /v1/audio/speech when the server responds 404 or 405.
//
// Configure base URL via TTSConfig.ModelPath or HOWL_MELOTTS_URL (for example http://127.0.0.1:8888 ).
// Optional HOWL_MELOTTS_API_KEY is sent as Authorization: Bearer for gated endpoints.
type Client struct {
	*tts.BaseSynthesizer
	httpClient *http.Client
	baseURL    *url.URL
	apiKey     string
}

// New allocates a synthesizer wired for MeloTTS HTTP backends.
func New() *Client {
	return &Client{
		BaseSynthesizer: tts.NewBaseSynthesizer(types.TTSMeloTTS),
		httpClient:      &http.Client{Timeout: 8 * time.Minute},
		apiKey:          strings.TrimSpace(os.Getenv("HOWL_MELOTTS_API_KEY")),
	}
}

// Initialize validates the remote endpoint and caches directories.
func (c *Client) Initialize(ctx context.Context, cfg types.TTSConfig) error {
	raw := strings.TrimSpace(cfg.ModelPath)
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("HOWL_MELOTTS_URL"))
	}
	if raw == "" {
		raw = "http://127.0.0.1:8888"
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return types.NewAudioError(
			types.ErrCodeModelNotFound,
			fmt.Sprintf("MeloTTS base URL invalid (%q): set HOWL_MELOTTS_URL or TTSConfig.ModelPath to http(s)://host:port", raw),
			err,
		)
	}
	c.baseURL = u

	outDir := strings.TrimSpace(cfg.CacheDir)
	if outDir == "" {
		outDir = filepath.Join(os.TempDir(), "howl-melotts-out")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to create MeloTTS output directory", err)
	}

	c.SetInitialized(cfg)
	return nil
}

// Synthesize calls the configured Melo HTTP server and stores returned audio as WAV bytes on disk.
func (c *Client) Synthesize(ctx context.Context, req types.SynthesisRequest, callback types.ProgressCallback) (*types.SynthesisResult, error) {
	if err := c.CheckInitialized(); err != nil {
		return nil, err
	}
	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}
	if callback != nil {
		callback(10, "contacting melotts backend")
	}

	outDir := c.outputDir()

	wavBody, ctype, synthErr := c.tryLegacySynth(ctx, req)
	if synthErr != nil {
		return nil, synthErr
	}
	if len(wavBody) == 0 {
		wavBody, ctype, synthErr = c.tryOpenAICompatSynth(ctx, req)
		if synthErr != nil {
			return nil, synthErr
		}
	}
	if len(wavBody) == 0 {
		return nil, types.NewAudioError(types.ErrCodeSynthesisFailed,
			"melotts server accepted the request but returned no audio (/synthesize and /v1/audio/speech); verify Melo endpoint and voice ids", nil)
	}

	dest := filepath.Join(outDir, fmt.Sprintf("melo_%d.wav", time.Now().UnixNano()))
	if err := os.WriteFile(dest, wavBody, 0o644); err != nil {
		return nil, types.NewAudioError(types.ErrCodeSynthesisFailed, "persist melotts audio", err)
	}

	meta := types.AudioMetadata{Format: types.FormatWAV}
	if ctype != "" {
		mt, _, _ := mime.ParseMediaType(ctype)
		if mt == "audio/wav" || mt == "audio/x-wav" || strings.Contains(ctype, "wav") {
			meta.Format = types.FormatWAV
		}
	}
	meta.SampleRate = 44100
	meta.Channels = 1
	meta.BitDepth = 16

	if md, err := types.GetAudioMetadata(dest); err == nil {
		meta = *md
	}

	if callback != nil {
		callback(100, "complete")
	}

	return &types.SynthesisResult{
		AudioPath: dest,
		Metadata:  meta,
		Duration:  meta.Duration,
	}, nil
}

// SynthesizeStream buffers server output locally then emits chunks identical to Piper implementation.
func (c *Client) SynthesizeStream(ctx context.Context, req types.SynthesisRequest, chunks chan<- types.AudioChunk) error {
	defer close(chunks)

	res, err := c.Synthesize(ctx, req, nil)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(res.AudioPath)
	if err != nil {
		return types.NewAudioError(types.ErrCodeSynthesisFailed, "read synthesized melotts audio", err)
	}
	const sz = 32 * 1024
	for i, off := 0, 0; off < len(data); i++ {
		end := off + sz
		if end > len(data) {
			end = len(data)
		}
		select {
		case <-ctx.Done():
			return types.NewAudioError(types.ErrCodeCancelled, "melotts stream cancelled", ctx.Err())
		case chunks <- types.AudioChunk{Data: data[off:end], Metadata: res.Metadata, Index: i, IsLast: end == len(data)}:
		}
		off = end
	}
	return nil
}

// GetVoices maps registry entries for UX surfaces.
func (c *Client) GetVoices() []types.Voice {
	prof, ok := types.GetTTSProfile(types.TTSMeloTTS)
	if !ok {
		return nil
	}
	voices := make([]types.Voice, 0, len(prof.SupportedVoices))
	for _, id := range prof.SupportedVoices {
		langHint := normalizeMeloLanguage(id)
		voices = append(voices, types.Voice{ID: id, Name: id, Language: strings.ToLower(langHint)})
	}
	return voices
}

// SupportsSSML reports Melo inference path (HTTP) does not accept SSML in this adapter.
func (c *Client) SupportsSSML() bool { return false }

// Release is a noop for HTTP client state.
func (c *Client) Release() error { return nil }

func (c *Client) outputDir() string {
	if d := strings.TrimSpace(c.GetConfig().CacheDir); d != "" {
		return d
	}
	return filepath.Join(os.TempDir(), "howl-melotts-out")
}

func (c *Client) tryLegacySynth(ctx context.Context, req types.SynthesisRequest) ([]byte, string, error) {
	u := cloneURL(*c.baseURL)
	u.Path = strings.TrimSuffix(u.Path, "/") + "/synthesize"
	body := map[string]any{
		"text":        req.Text,
		"language":    normalizeMeloLanguage(req.Language),
		"speaker_id":  normalizeSpeakerID(req.Voice),
		"speed":       req.Speed,
		"length_scale": invertSpeed(req.Speed),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, "", types.NewAudioError(types.ErrCodeSynthesisFailed, "marshal legacy melotts request", err)
	}
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, "", err
	}
	rq.Header.Set("Content-Type", "application/json")
	c.attachAuth(rq)

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, "", types.NewAudioError(types.ErrCodeSynthesisFailed, "melotts /synthesize request failed", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusMethodNotAllowed, http.StatusNotFound:
		return nil, "", nil
	default:
		b, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		bodyStr := ""
		if readErr == nil {
			bodyStr = strings.TrimSpace(string(b))
		}
		return nil, "", types.NewAudioError(types.ErrCodeSynthesisFailed,
			fmt.Sprintf("melotts /synthesize HTTP %s: %s", resp.Status, bodyStr), nil)
	}

	b, readErr := io.ReadAll(io.LimitReader(resp.Body, 256<<20))
	if readErr != nil {
		return nil, "", readErr
	}
	ct := resp.Header.Get("Content-Type")

	if wav, handled, ferr := coerceToWAV(b, ct); handled {
		if ferr != nil {
			return nil, ct, ferr
		}
		return wav, ct, nil
	}
	return nil, ct, types.NewAudioError(types.ErrCodeSynthesisFailed,
		"unexpected melotts /synthesize body (need audio/wav or JSON with audio fields)", nil)
}

func (c *Client) tryOpenAICompatSynth(ctx context.Context, req types.SynthesisRequest) ([]byte, string, error) {
	u := cloneURL(*c.baseURL)
	u.Path = strings.TrimSuffix(u.Path, "/") + "/v1/audio/speech"
	model := strings.TrimSpace(os.Getenv("HOWL_MELOTTS_OPENAI_MODEL"))
	if model == "" {
		model = "tts-1"
	}
	payload := map[string]any{
		"model":           model,
		"input":           req.Text,
		"voice":           normalizeSpeakerID(req.Voice),
		"response_format": "wav",
		"speed":           clampSpeechSpeed(req.Speed),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, "", types.NewAudioError(types.ErrCodeSynthesisFailed, "marshal openai-compat melotts request", err)
	}
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	rq.Header.Set("Content-Type", "application/json")
	c.attachAuth(rq)

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, "", types.NewAudioError(types.ErrCodeSynthesisFailed, "melotts /v1/audio/speech request failed", err)
	}
	defer resp.Body.Close()

	b, readErr := io.ReadAll(io.LimitReader(resp.Body, 256<<20))
	if readErr != nil {
		return nil, "", readErr
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", types.NewAudioError(types.ErrCodeSynthesisFailed,
			fmt.Sprintf("melotts OpenAI-compat HTTP %s: %s", resp.Status, strings.TrimSpace(string(b))), nil)
	}
	ct := resp.Header.Get("Content-Type")

	if wav, handled, ferr := coerceToWAV(b, ct); handled {
		return wav, ct, ferr
	}
	return nil, ct, types.NewAudioError(types.ErrCodeSynthesisFailed,
		"missing audio bytes from OpenAI-compatible melotts response", nil)
}

func (c *Client) attachAuth(rq *http.Request) {
	if c.apiKey == "" {
		return
	}
	rq.Header.Set("Authorization", "Bearer "+c.apiKey)
}

func cloneURL(u url.URL) *url.URL {
	u2 := u
	return &u2
}

func normalizeSpeakerID(raw string) string {
	s := strings.TrimSpace(raw)
	if s != "" {
		return s
	}
	prof, ok := types.GetTTSProfile(types.TTSMeloTTS)
	if ok && len(prof.SupportedVoices) > 0 {
		return prof.SupportedVoices[0]
	}
	return "EN"
}

func normalizeMeloLanguage(lang string) string {
	s := strings.TrimSpace(strings.ToUpper(lang))
	if s == "" {
		return "EN"
	}
	if len(s) >= 2 {
		switch s[:2] {
		case "ZH", "JA", "KO", "ES", "FR":
			switch s[:2] {
			case "JA":
				return "JP"
			case "KO":
				return "KR"
			default:
				return s[:2]
			}
		}
	}
	return "EN"
}

func invertSpeed(speed float64) float64 {
	if speed <= 0 {
		return 1
	}
	return 1 / clampSpeechSpeed(speed)
}

func clampSpeechSpeed(v float64) float64 {
	if v <= 0 {
		return 1
	}
	if v < 0.25 {
		return 0.25
	}
	if v > 4 {
		return 4
	}
	return v
}

func coerceToWAV(body []byte, ct string) (wav []byte, handled bool, err error) {
	mt, _, _ := mime.ParseMediaType(ct)
	if mt == "" {
		mt = ct
	}
	if strings.Contains(mt, "wav") || (len(body) > 12 && bytes.Equal(body[0:4], []byte("RIFF")) && bytes.Equal(body[8:12], []byte("WAVE"))) {
		return body, true, nil
	}
	if ct == "" || strings.Contains(ct, "json") {
		var wrap map[string]json.RawMessage
		if err := json.Unmarshal(body, &wrap); err != nil {
			return nil, false, nil
		}
		for _, key := range []string{"audio_wav_base64", "audio_base64", "audio", "data", "wav_base64"} {
			if raw, ok := wrap[key]; ok && len(raw) >= 2 {
				var encoded string
				if err := json.Unmarshal(raw, &encoded); err == nil && encoded != "" {
					dec, derr := base64.StdEncoding.DecodeString(encoded)
					if derr == nil && len(dec) > 0 {
						return dec, true, nil
					}
				}
				var nested map[string]string
				if err := json.Unmarshal(raw, &nested); err == nil {
					for _, sub := range nested {
						if sub == "" {
							continue
						}
						if dec, derr := base64.StdEncoding.DecodeString(sub); derr == nil && len(dec) > 0 {
							return dec, true, nil
						}
					}
				}
			}
		}
		return nil, true, types.NewAudioError(types.ErrCodeSynthesisFailed,
			"MeloTTS JSON response missing recognised audio_* base64 payload", nil)
	}
	return nil, false, nil
}
