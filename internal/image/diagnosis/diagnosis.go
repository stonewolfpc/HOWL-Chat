// Package diagnosis turns an engine or admission failure into a category a user
// can act on.
//
// The engine reports everything through one channel: an HTTP error body, a
// process exit, or an admission string. Those are not the same problem, and
// treating them as one made a model the engine cannot identify look like a VRAM
// shortfall. Each category here answers "what should the user try next?".
package diagnosis

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Category is the kind of failure, independent of where it was observed.
type Category string

const (
	// UnsupportedArchitecture means the engine does not know this architecture
	// at all. More VRAM will not help, and neither will another copy of the same
	// weights.
	UnsupportedArchitecture Category = "unsupported_architecture"
	// UnsupportedWeightFormat means the engine could not identify the model from
	// this file. The architecture may well be supported through another release
	// of the same family; the quantized/GGUF variant is what failed. This is a
	// different user action from an unknown architecture, so it is a different
	// category.
	UnsupportedWeightFormat Category = "unsupported_weight_format"
	// InvalidModelFile means the file is missing, truncated, or not a model this
	// engine can open.
	InvalidModelFile Category = "invalid_model_file"
	// InsufficientResources means the workload legitimately needs more memory
	// than the post-reserve budget allows.
	InsufficientResources Category = "insufficient_resources"
	// PlacementFailure means the requested placement (backend, parameter backend,
	// flags) was rejected or could not be applied.
	PlacementFailure Category = "placement_failure"
	// EngineStartupFailure means the engine process never became usable.
	EngineStartupFailure Category = "engine_startup_failure"
	// GenerationFailure means the engine started and accepted the request but the
	// job itself failed.
	GenerationFailure Category = "generation_failure"
)

// Stage says where the failure was observed, which disambiguates identical
// engine text appearing at different points.
type Stage string

const (
	StageStartup    Stage = "startup"
	StagePlacement  Stage = "placement"
	StageGeneration Stage = "generation"
)

// Diagnosis is a classified failure with a user-facing explanation.
type Diagnosis struct {
	Category   Category
	Stage      Stage
	Summary    string
	Suggestion string
	Retryable  bool
	Raw        string
}

// Describe is the single line a job or log carries. It always names the
// category so a later reader cannot mistake one problem for another.
func (d Diagnosis) Describe() string {
	line := fmt.Sprintf("[%s] %s", d.Category, d.Summary)
	if d.Suggestion != "" {
		line += " " + d.Suggestion
	}
	if raw := strings.TrimSpace(d.Raw); raw != "" {
		line += " (engine said: " + raw + ")"
	}
	return line
}

// Error is the wrapped form, so callers can keep using errors.As.
type Error struct {
	Diagnosis Diagnosis
}

func (e *Error) Error() string { return e.Diagnosis.Describe() }

func (e *Error) Unwrap() error { return nil }

// As returns the diagnosis already carried by err, if any.
func As(err error) (Diagnosis, bool) {
	var diagnosed *Error
	if errors.As(err, &diagnosed) {
		return diagnosed.Diagnosis, true
	}
	return Diagnosis{}, false
}

// Wrap classifies err and returns it with the diagnosis attached. An already
// classified error is returned unchanged, so a failure is never re-labeled on
// its way up.
func Wrap(stage Stage, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := As(err); ok {
		return err
	}
	return &Error{Diagnosis: Diagnose(stage, err)}
}

// Diagnose classifies err without wrapping it.
func Diagnose(stage Stage, err error) Diagnosis {
	raw := strings.TrimSpace(err.Error())
	lower := strings.ToLower(raw)
	d := Diagnosis{Stage: stage, Raw: raw}

	switch {
	case containsAny(lower, "cannot identify model family", "could not identify model family",
		"unable to identify model family", "cannot determine model family", "cannot identify model type"):
		// Identification failed. Whether that is the architecture or the file
		// depends on what the engine was handed: a format it does not read says
		// nothing about the architecture.
		switch modelFileExt(raw) {
		case ".gguf", ".pt", ".pth", ".bin":
			d.Category = UnsupportedWeightFormat
			d.Summary = "The engine cannot identify the model from this file; this weight format or quantized variant is not supported, even though the architecture may be."
			d.Suggestion = "Use the original Safetensors release for this model family, or update the engine; memory settings will not change this."
		default:
			d.Category = UnsupportedArchitecture
			d.Summary = "The engine does not know this architecture, so it cannot run it."
			d.Suggestion = "Choose a model file this engine supports, or update the engine; memory settings will not change this."
		}
	case containsAny(lower, "unsupported model", "unknown model architecture", "model type not supported",
		"architecture not supported", "unknown architecture"):
		d.Category = UnsupportedArchitecture
		d.Summary = "The engine does not know this architecture, so it cannot run it."
		d.Suggestion = "Choose a model file this engine supports, or update the engine; memory settings will not change this."
	case containsAny(lower, "no such file", "file not found", "cannot open", "failed to load model",
		"invalid model", "not a valid", "unsupported format", "unknown format", "unexpected end of file",
		"corrupt", "truncated"):
		d.Category = InvalidModelFile
		d.Summary = "The model file could not be opened or parsed."
		d.Suggestion = "Check that the file exists and finished downloading, then re-scan the model folder."
	case containsAny(lower, "out of memory", "failed to allocate", "not enough memory", "oom",
		"exceeds the budget", "needs ", "mb vram", "insufficient"):
		d.Category = InsufficientResources
		d.Summary = "The workload does not fit the memory budget that admission computed."
		d.Suggestion = "Use a smaller or more quantized model, lower the resolution, or unload another workload."
		d.Retryable = false
	case containsAny(lower, "unknown argument", "invalid argument", "invalid value", "invalid backend",
		"backend", "params-backend", "max-vram", "bad flag"):
		d.Category = PlacementFailure
		d.Summary = "The engine rejected the requested placement."
		d.Suggestion = "Reset Backend and Params Backend to Auto so placement can be planned, then try again."
		d.Retryable = true
	case containsAny(lower, "sd-server exited", "engine stopped", "failed to start", "could not start",
		"connection refused", "dial tcp", "process exited", "engine not ready", "not listening"):
		d.Category = EngineStartupFailure
		d.Summary = "The engine process did not become usable."
		d.Suggestion = "Check the engine binary and model folder, then retry; the previous engine was restarted."
		d.Retryable = true
	default:
		d.Category = GenerationFailure
		d.Summary = "The engine started but the generation failed."
		d.Suggestion = "Retry the job; if it repeats, check the engine log for the failing stage."
		d.Retryable = true
	}

	if stage == StageStartup && d.Category == GenerationFailure {
		// A startup failure with no recognisable text is still a startup failure.
		d.Category = EngineStartupFailure
		d.Summary = "The engine process did not become usable."
		d.Suggestion = "Check the engine binary and model folder, then retry."
		d.Retryable = true
	}
	return d
}

// modelFileExt reads the file the engine named. The path is the only reliable
// signal for format-versus-architecture: the engine reports both as "cannot
// identify model family".
func modelFileExt(raw string) string {
	lower := strings.ToLower(raw)
	const marker = "from '"
	if i := strings.Index(lower, marker); i >= 0 {
		rest := raw[i+len(marker):]
		if j := strings.Index(rest, "'"); j > 0 {
			return strings.ToLower(filepath.Ext(rest[:j]))
		}
	}
	for _, ext := range []string{".gguf", ".safetensors", ".sft", ".ckpt", ".pt", ".pth", ".bin"} {
		if strings.Contains(lower, ext) {
			return ext
		}
	}
	return ""
}

func containsAny(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}
