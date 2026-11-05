package gguf

// Projector pairing: choose the one projector that belongs to a model.
//
// The old loader returned the first file matching a list of generic names
// (mmproj-Q8_0.gguf, mmproj.gguf, then any mmproj-*.gguf in directory order).
// In a directory holding several audio models that silently attached the
// wrong projector, for example Qwen3-ASR-0.6B's projector to the 1.7B model.
// Directory order is not evidence of pairing, so this resolver never uses it.
//
// CONTRACT
//
//  1. Explicit identity wins: when ExpectedNames is supplied (the ASR registry
//     names its companion file), exactly one of those files is resolved.
//  2. Otherwise the companion is found by exact identity. The model filename
//     and the projector filename are reduced to an identity key by removing
//     separators, quantization/precision tokens, the "mmproj" token and
//     version tokens (v0_5). Both keys must be equal. This is what pairs
//     Qwen2.5-Omni-3B-Q4_K_M with Qwen2.5-Omni-3B-mmproj-Q8_0, and
//     Ultravox-Llama-3.2-1B with mmproj-Ultravox-v0_5-llama-3_2-1b-f16.
//  3. The projector's own metadata must declare an encoder, and that kind
//     must satisfy the model's requirement (vision-only architectures reject
//     an audio-only projector).
//  4. Zero matches is an error when a projector is required; more than one
//     match is always an error. An ambiguous or mismatched set is reported
//     with every candidate path so the user can rename or remove one.
//
// COMPATIBILITY NOTE
//
// clip.<modality>.projection_dim is compared with <arch>.embedding_length and
// recorded as DimensionMismatch. It is not a hard rejection: the verified
// Ultravox pair is 4096 vs 2048 yet is the correct pairing, because that
// projector stacks frames before projecting. Identity and encoder kind are the
// hard rules; the dimension is diagnostic evidence.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Pairing errors are distinct so callers can report why loading stopped.
var (
	ErrProjectorMissing   = errors.New("no companion projector found")
	ErrProjectorAmbiguous = errors.New("ambiguous companion projector")
	ErrProjectorMismatch  = errors.New("companion projector is incompatible")
)

// PairingRequest describes the model whose companion is being resolved.
type PairingRequest struct {
	ModelPath string
	// Arch is the model's GGUF architecture, used only for diagnostics.
	Arch string
	// EmbeddingLength is the model's hidden size, when known.
	EmbeddingLength int
	// ExpectedNames are explicit companion filenames. When set they are
	// authoritative and no identity scan is performed.
	ExpectedNames []string
	// Require makes a missing projector an error.
	Require bool
	// RequiredModality narrows acceptable encoder kinds. ProjectorUnknown with
	// Require=true accepts any declared kind.
	RequiredModality ProjectorKind
	// SearchDirs are additional directories to check, after the model's own
	// directory. Duplicates are removed.
	SearchDirs []string
}

// ProjectorPairing is the resolved companion.
type ProjectorPairing struct {
	Path           string
	Info           ProjectorInfo
	Identity       string
	ProjectionDim  int
	ModelEmbedding int
	// DimensionMismatch is diagnostic; see the compatibility note above.
	DimensionMismatch bool
	Notes             []string
}

// ResolveProjector finds the exact companion projector for a model.
func ResolveProjector(req PairingRequest) (ProjectorPairing, error) {
	if strings.TrimSpace(req.ModelPath) == "" {
		return ProjectorPairing{}, fmt.Errorf("%w: model path is empty", ErrProjectorMissing)
	}
	identity := ProjectorIdentity(req.ModelPath)
	if identity == "" {
		return ProjectorPairing{}, fmt.Errorf("%w: cannot derive an identity from %s", ErrProjectorMissing, filepath.Base(req.ModelPath))
	}

	dirs := pairingDirs(req.ModelPath, req.SearchDirs)

	if len(req.ExpectedNames) > 0 {
		return resolveExplicit(req, identity, dirs)
	}
	return resolveByIdentity(req, identity, dirs)
}

func resolveExplicit(req PairingRequest, identity string, dirs []string) (ProjectorPairing, error) {
	var found []string
	var attempted []string
	for _, name := range req.ExpectedNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		for _, dir := range dirs {
			candidate := filepath.Join(dir, name)
			attempted = append(attempted, candidate)
			if isRegularFile(candidate) {
				found = append(found, candidate)
			}
		}
	}
	found = uniquePaths(found)
	if len(found) == 0 {
		return ProjectorPairing{}, fmt.Errorf("%w: expected %s beside %s (checked %s)",
			ErrProjectorMissing, strings.Join(req.ExpectedNames, ", "),
			filepath.Base(req.ModelPath), strings.Join(attempted, ", "))
	}
	if len(found) > 1 {
		return ProjectorPairing{}, fmt.Errorf("%w: %s exists in more than one search directory: %s",
			ErrProjectorAmbiguous, strings.Join(req.ExpectedNames, ", "), strings.Join(found, ", "))
	}

	info := InspectProjector(found[0])
	if err := validateProjector(info, req); err != nil {
		return ProjectorPairing{}, err
	}
	return newPairing(info, identity, req, nil), nil
}

func resolveByIdentity(req PairingRequest, identity string, dirs []string) (ProjectorPairing, error) {
	var accepted []string
	var rejected []string

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".gguf") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			if samePath(path, req.ModelPath) {
				continue
			}
			if ProjectorIdentity(entry.Name()) != identity {
				continue
			}

			info := InspectProjector(path)
			if info.Err != "" {
				rejected = append(rejected, fmt.Sprintf("%s (unreadable: %s)", path, info.Err))
				continue
			}
			if info.Kind == ProjectorUnknown {
				rejected = append(rejected, fmt.Sprintf("%s (declares no vision or audio encoder)", path))
				continue
			}
			if !modalitySatisfied(info, req.RequiredModality) {
				rejected = append(rejected, fmt.Sprintf("%s (declares %s, model requires %s)", path, info.Kind, req.RequiredModality))
				continue
			}
			accepted = append(accepted, path)
		}
	}

	sort.Strings(accepted)
	sort.Strings(rejected)

	switch len(accepted) {
	case 1:
		return newPairing(InspectProjector(accepted[0]), identity, req, rejected), nil
	case 0:
		if !req.Require {
			// A text model may legitimately have no companion. Identity
			// candidates were still inspected, never assumed.
			return ProjectorPairing{Identity: identity, ModelEmbedding: req.EmbeddingLength, Notes: rejected}, nil
		}
		if len(rejected) > 0 {
			return ProjectorPairing{}, fmt.Errorf("%w: model %s (architecture %s) requires a %s projector, but no same-identity projector is usable: %s",
				ErrProjectorMismatch, filepath.Base(req.ModelPath), req.Arch, describeModality(req.RequiredModality), strings.Join(rejected, "; "))
		}
		return ProjectorPairing{}, fmt.Errorf("%w: model %s (architecture %s) requires a %s projector with identity %q, but none was found beside it",
			ErrProjectorMissing, filepath.Base(req.ModelPath), req.Arch, describeModality(req.RequiredModality), identity)
	default:
		return ProjectorPairing{}, fmt.Errorf("%w: model %s matches more than one projector: %s",
			ErrProjectorAmbiguous, filepath.Base(req.ModelPath), strings.Join(accepted, ", "))
	}
}

func newPairing(info ProjectorInfo, identity string, req PairingRequest, rejected []string) ProjectorPairing {
	pairing := ProjectorPairing{
		Path:           info.Path,
		Info:           info,
		Identity:       identity,
		ProjectionDim:  info.ProjectionDim,
		ModelEmbedding: req.EmbeddingLength,
		Notes:          rejected,
	}
	if req.EmbeddingLength > 0 && info.ProjectionDim > 0 && info.ProjectionDim != req.EmbeddingLength {
		pairing.DimensionMismatch = true
		pairing.Notes = append(pairing.Notes,
			fmt.Sprintf("projector projection_dim=%d differs from model embedding_length=%d (recorded, not rejected: compound projectors may stack frames)", info.ProjectionDim, req.EmbeddingLength))
	}
	return pairing
}

func validateProjector(info ProjectorInfo, req PairingRequest) error {
	if info.Err != "" {
		return fmt.Errorf("%w: %s is unreadable: %s", ErrProjectorMismatch, info.Path, info.Err)
	}
	if info.Kind == ProjectorUnknown {
		return fmt.Errorf("%w: %s declares no vision or audio encoder", ErrProjectorMismatch, info.Path)
	}
	if !modalitySatisfied(info, req.RequiredModality) {
		return fmt.Errorf("%w: %s declares %s, model requires %s", ErrProjectorMismatch, info.Path, info.Kind, describeModality(req.RequiredModality))
	}
	return nil
}

func modalitySatisfied(info ProjectorInfo, required ProjectorKind) bool {
	switch required {
	case ProjectorVision:
		return info.HasVision()
	case ProjectorAudio:
		return info.HasAudio()
	default:
		// No restriction: any declared encoder is usable.
		return info.Kind != ProjectorUnknown
	}
}

func describeModality(kind ProjectorKind) string {
	if kind == ProjectorUnknown {
		return "declared"
	}
	return string(kind)
}

// ProjectorIdentity reduces a model or projector filename to the identity both
// sides must share. It is intentionally conservative: it removes only tokens
// that are known not to identify the model (quantization, mmproj marker and
// version markers) and normalises separators.
func ProjectorIdentity(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return normalizeIdentity(base)
}

// ModelIdentity normalises a model name used as identity evidence. Unlike
// ProjectorIdentity it does not strip an arbitrary dotted suffix: metadata names
// such as "Llama 3.2 1B Instruct" are not filenames, and filepath.Ext would
// truncate them after "Llama 3". Only a known model-file extension is removed.
func ModelIdentity(name string) string {
	base := strings.TrimSpace(name)
	lower := strings.ToLower(base)
	for _, ext := range []string{".gguf", ".bin", ".ggml"} {
		if strings.HasSuffix(lower, ext) {
			base = base[:len(base)-len(ext)]
			break
		}
	}
	return normalizeIdentity(base)
}

func normalizeIdentity(base string) string {
	base = strings.ToLower(base)

	base = mmprojToken.ReplaceAllString(base, "$1")
	for {
		next := quantToken.ReplaceAllString(base, "$1")
		if next == base {
			break
		}
		base = next
	}
	for {
		next := versionToken.ReplaceAllString(base, "$1")
		if next == base {
			break
		}
		base = next
	}

	base = strings.NewReplacer("_", "-", ".", "-", " ", "-").Replace(base)
	base = repeatedSeparators.ReplaceAllString(base, "-")
	return strings.Trim(base, "-")
}

var (
	// q4_k_m, q8_0, iq4_xs, f16, bf16, mxfp4 ... are precision tokens.
	quantToken = regexp.MustCompile(`(?i)(^|[-_. ])(iq\d+([._][a-z0-9]+)*|q\d+([._][a-z0-9]+)*|f16|f32|bf16|fp16|fp32|mxfp4|mxfp8|nvfp4)([-_. ]|$)`)
	// v0_5, v1.0, v2 ... are release markers, not model identity.
	versionToken = regexp.MustCompile(`(?i)(^|[-_. ])v\d+([._]\d+)*([-_. ]|$)`)
	// mmproj is the file family marker shared by every projector.
	mmprojToken        = regexp.MustCompile(`(?i)(^|[-_. ])mmproj([-_. ]|$)`)
	repeatedSeparators = regexp.MustCompile(`-{2,}`)
)

func pairingDirs(modelPath string, extra []string) []string {
	dirs := []string{filepath.Dir(modelPath)}
	dirs = append(dirs, extra...)
	return uniquePaths(dirs)
}

func uniquePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		clean := filepath.Clean(path)
		key := strings.ToLower(clean)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, clean)
	}
	return out
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info != nil && info.Mode().IsRegular()
}

func samePath(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
