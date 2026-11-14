package helper

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"howl-chat/internal/audio/types"
)

// Failure is the helper boundary's precise internal error. It is deliberately
// not collapsed into AudioError: callers need the protocol code and detail for
// diagnostics even when the public audio error is only a compatibility mapping.
type Failure struct {
	Code      string
	Message   string
	Retryable bool
	Field     string
	Detail    map[string]string
}

func (f *Failure) Error() string {
	if f == nil {
		return "<nil>"
	}
	msg := f.Message
	if msg == "" {
		msg = f.Code
	}
	if f.Field != "" {
		return fmt.Sprintf("%s (%s): %s", msg, f.Code, f.Field)
	}
	return fmt.Sprintf("%s (%s)", msg, f.Code)
}

func newFailure(code, message string) *Failure {
	return &Failure{Code: code, Message: message}
}

func (f *Failure) withField(field string) *Failure {
	if f == nil {
		return nil
	}
	f.Field = field
	return f
}

func (f *Failure) withRetryable(retryable bool) *Failure {
	if f == nil {
		return nil
	}
	f.Retryable = retryable
	return f
}

func (f *Failure) status() types.TTSFailure {
	if f == nil {
		return types.TTSFailure{}
	}
	return types.TTSFailure{
		Code:      f.Code,
		Message:   f.Message,
		Retryable: f.Retryable,
		Field:     f.Field,
		Detail:    cloneDetail(f.Detail),
	}
}

func cloneDetail(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// AudioError projects the internal failure onto the existing public error
// taxonomy. The original Failure remains in the chain via Cause.
func (f *Failure) AudioError() *types.AudioError {
	if f == nil {
		return &types.AudioError{
			Code:    types.ErrCodeSynthesisFailed,
			Message: "TTS helper failure classification is missing",
		}
	}
	return &types.AudioError{
		Code:        audioCodeFor(f.Code),
		Message:     f.Message,
		Cause:       f,
		Recoverable: f.Retryable,
	}
}

// AsFailure unwraps a helper failure from an error chain.
func AsFailure(err error) (*Failure, bool) {
	var f *Failure
	if errors.As(err, &f) {
		return f, true
	}
	var audioErr *types.AudioError
	if errors.As(err, &audioErr) && audioErr != nil {
		if f, ok := audioErr.Cause.(*Failure); ok {
			return f, true
		}
	}
	return nil, false
}

func audioCodeFor(code string) types.AudioErrorCode {
	switch code {
	case "helper_missing", "provider_unknown", "model_missing", "voice_missing":
		return types.ErrCodeModelNotFound
	case "protocol_incompatible", "protocol_minor_unsupported", "helper_start_failed",
		"dependency_missing", "provider_not_loaded", "provider_already_loaded",
		"model_identity_mismatch":
		return types.ErrCodeModelLoadFailed
	case "cancelled":
		return types.ErrCodeCancelled
	case "timeout":
		return types.ErrCodeTimeout
	default:
		return types.ErrCodeSynthesisFailed
	}
}

// DiagnosticString is bounded and never includes request text. It is suitable
// for status/log surfaces while preserving the precise code.
func (f *Failure) DiagnosticString() string {
	if f == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(f.Code)
	if f.Message != "" {
		b.WriteString(": ")
		b.WriteString(f.Message)
	}
	if f.Field != "" {
		b.WriteString(" [")
		b.WriteString(f.Field)
		b.WriteString("]")
	}
	if len(f.Detail) > 0 {
		keys := make([]string, 0, len(f.Detail))
		for k := range f.Detail {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString(" (")
		for i, k := range keys {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(f.Detail[k])
		}
		b.WriteString(")")
	}
	return b.String()
}
