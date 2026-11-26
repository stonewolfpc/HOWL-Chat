package audio

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"howl-chat/internal/audio/types"
	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/paths"
)

// ModelSource records where a resolved model file was found.
type ModelSource string

const (
	ModelSourceEnv    ModelSource = "env"
	ModelSourceNested ModelSource = "nested"
	ModelSourceFlat   ModelSource = "flat"
)

// Resolver errors are deliberately distinct so callers can report whether the
// engine is unsupported, the policy is absent, chose nothing, or chose
// something invalid.
var (
	ErrEngineUnsupported      = errors.New("unsupported ASR engine")
	ErrModelNotFound          = errors.New("asr model not found")
	ErrPolicyNotConfigured    = errors.New("asr model policy not configured")
	ErrNoModelSelected        = errors.New("asr model policy selected no model")
	ErrPolicyInvalidSelection = errors.New("asr model policy returned an invalid selection")
)

// ModelFile is a located ASR model file. Type is empty only for an untyped
// HOWL_ASR_MODEL candidate; a policy must assign an approved type before
// selecting that candidate.
type ModelFile struct {
	Path   string
	Type   types.ASRModelType
	Source ModelSource
}

// SelectionPolicy chooses from the candidates the resolver discovered. It may
// filter or reorder them, but may not introduce paths or change canonical
// identities.
type SelectionPolicy interface {
	Candidates(installed []ModelFile) ([]ModelFile, error)
}

// ResolverConfig carries the resolver's injectable dependencies. Root must be
// an absolute path and Getenv must be non-nil. Stat is the sole filesystem
// primitive used by the resolver; nil means os.Stat.
type ResolverConfig struct {
	Root       string
	Getenv     func(string) string
	Stat       func(path string) (os.FileInfo, error)
	SearchDirs []string
}

// Resolver resolves an ASR model without acquiring, downloading, or running it.
type Resolver struct {
	cfg ResolverConfig
}

type knownASRModel struct {
	Type types.ASRModelType
	File string
}

// approvedASRModels is derived from the single registry declaration. The
// resolver cannot disagree with the registry, the planner, or the settings UI
// about which engines may be loaded locally (DEBT-T3): an entry is approved
// only when it is local, admitted, and names the file it loads.
var approvedASRModels = buildApprovedASRModels()

func buildApprovedASRModels() []knownASRModel {
	modelTypes := make([]types.ASRModelType, 0, len(types.ASRModelRegistry))
	for modelType := range types.ASRModelRegistry {
		modelTypes = append(modelTypes, modelType)
	}
	sort.Slice(modelTypes, func(i, j int) bool { return modelTypes[i] < modelTypes[j] })

	out := make([]knownASRModel, 0, len(modelTypes))
	for _, modelType := range modelTypes {
		profile := types.ASRModelRegistry[modelType]
		if !profile.Local || !profile.ApprovedLocal || profile.LocalFile == "" {
			continue
		}
		out = append(out, knownASRModel{Type: modelType, File: profile.LocalFile})
	}
	return out
}

// NewResolver creates a resolver from explicit dependencies. Production root
// and environment discovery belongs to the caller, not this constructor.
func NewResolver(cfg ResolverConfig) (*Resolver, error) {
	if strings.TrimSpace(cfg.Root) == "" {
		return nil, fmt.Errorf("invalid resolver configuration: root is required")
	}
	if !filepath.IsAbs(cfg.Root) {
		return nil, fmt.Errorf("invalid resolver configuration: root must be absolute: %q", cfg.Root)
	}
	if cfg.Getenv == nil {
		return nil, fmt.Errorf("invalid resolver configuration: Getenv is required")
	}
	cfg.Root = filepath.Clean(cfg.Root)
	return &Resolver{cfg: cfg}, nil
}

// Resolve locates a model for an explicit engine or applies the user's policy
// for "auto". It performs no acquisition and never invents model identity.
func (r *Resolver) Resolve(engine string, policy SelectionPolicy) (ModelFile, error) {
	engine = strings.TrimSpace(engine)
	if engine == "" {
		return ModelFile{}, fmt.Errorf("%w: engine is blank", ErrEngineUnsupported)
	}
	if engine == "auto" {
		return r.resolveAuto(policy)
	}
	return r.resolveExplicit(engine)
}

func (r *Resolver) resolveExplicit(engine string) (ModelFile, error) {
	modelType := types.ASRModelType(engine)
	fileName, ok := approvedASRModelFile(modelType)
	if !ok {
		return ModelFile{}, fmt.Errorf("%w: %q", ErrEngineUnsupported, engine)
	}

	if override, ok := r.override(); ok {
		if err := r.requireRegular(override); err != nil {
			return ModelFile{}, err
		}
		return ModelFile{Path: override, Type: modelType, Source: ModelSourceEnv}, nil
	}

	locations := r.modelLocations(fileName)
	for _, path := range locations {
		if r.isRegular(path) {
			source := ModelSourceFlat
			if strings.Contains(filepath.ToSlash(path), "/models/whisper/") {
				source = ModelSourceNested
			}
			return ModelFile{Path: path, Type: modelType, Source: source}, nil
		}
	}
	if len(locations) == 0 {
		return ModelFile{}, fmt.Errorf("%w: no search locations for %s", ErrModelNotFound, fileName)
	}
	return ModelFile{}, fmt.Errorf("%w: tried %s", ErrModelNotFound, strings.Join(locations, ", "))
}

func (r *Resolver) resolveAuto(policy SelectionPolicy) (ModelFile, error) {
	// Policy presence is the first check so a nil policy performs no model
	// existence lookups at all.
	if policy == nil {
		return ModelFile{}, ErrPolicyNotConfigured
	}

	override, hasOverride := r.override()
	if hasOverride {
		// A configured override that is broken fails closed before discovery.
		if err := r.requireRegular(override); err != nil {
			return ModelFile{}, err
		}
	}

	candidates := r.canonicalCandidates()
	if hasOverride && !candidatePathExists(candidates, override) {
		candidates = append(candidates, ModelFile{Path: override, Source: ModelSourceEnv})
	}
	if len(candidates) == 0 {
		return ModelFile{}, fmt.Errorf("%w: no ASR model candidates under %s", ErrModelNotFound, r.root())
	}

	selected, err := policy.Candidates(candidates)
	if err != nil {
		return ModelFile{}, err
	}
	if len(selected) == 0 {
		return ModelFile{}, ErrNoModelSelected
	}
	return validatePolicySelection(candidates, selected)
}

func (r *Resolver) canonicalCandidates() []ModelFile {
	var out []ModelFile
	seen := make(map[string]struct{})

	add := func(path string, modelType types.ASRModelType, source ModelSource) {
		key := filepath.Clean(path)
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		out = append(out, ModelFile{Path: key, Type: modelType, Source: source})
	}

	for _, model := range approvedASRModels {
		path := r.nestedPath(model.File)
		if r.isRegular(path) {
			add(path, model.Type, ModelSourceNested)
		}
	}
	for _, model := range approvedASRModels {
		path := r.flatPath(model.File)
		if r.isRegular(path) {
			add(path, model.Type, ModelSourceFlat)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		left, right := sourceRank(out[i].Source), sourceRank(out[j].Source)
		if left != right {
			return left < right
		}
		return out[i].Path < out[j].Path
	})
	return out
}

func validatePolicySelection(candidates, selected []ModelFile) (ModelFile, error) {
	byPath := make(map[string]ModelFile, len(candidates))
	for _, candidate := range candidates {
		byPath[filepath.Clean(candidate.Path)] = candidate
	}

	seen := make(map[string]types.ASRModelType, len(selected))
	var result ModelFile
	for _, choice := range selected {
		path := filepath.Clean(choice.Path)
		candidate, ok := byPath[path]
		if !ok {
			return ModelFile{}, fmt.Errorf("%w: path %q is not a discovered candidate", ErrPolicyInvalidSelection, choice.Path)
		}
		if previous, duplicate := seen[path]; duplicate && previous != choice.Type {
			return ModelFile{}, fmt.Errorf("%w: path %q was selected with conflicting identities", ErrPolicyInvalidSelection, choice.Path)
		}
		seen[path] = choice.Type

		if candidate.Type != "" {
			if choice.Type != candidate.Type {
				return ModelFile{}, fmt.Errorf("%w: canonical candidate %q cannot be retyped from %q to %q", ErrPolicyInvalidSelection, choice.Path, candidate.Type, choice.Type)
			}
			if result.Path == "" {
				result = candidate
			}
			continue
		}

		if !isApprovedASRModelType(choice.Type) {
			return ModelFile{}, fmt.Errorf("%w: untyped override %q was selected without an approved type", ErrPolicyInvalidSelection, choice.Path)
		}
		if result.Path == "" {
			result = ModelFile{Path: path, Type: choice.Type, Source: candidate.Source}
		}
	}

	if result.Path == "" {
		return ModelFile{}, ErrNoModelSelected
	}
	return result, nil
}

func (r *Resolver) root() string {
	return filepath.Clean(r.cfg.Root)
}

func (r *Resolver) override() (string, bool) {
	value := strings.TrimSpace(r.cfg.Getenv("HOWL_ASR_MODEL"))
	if value == "" {
		return "", false
	}
	if filepath.IsAbs(value) {
		return filepath.Clean(value), true
	}
	return filepath.Join(r.root(), value), true
}

// modelLocations lists the directory-agnostic places a model file may live:
// the historical models/whisper and models/ paths, then every configured ASR
// inventory directory (directly and up to a bounded depth).
func (r *Resolver) modelLocations(fileName string) []string {
	locations := []string{r.nestedPath(fileName), r.flatPath(fileName)}
	for _, dir := range r.searchDirs() {
		locations = append(locations, filepath.Join(dir, fileName))
		locations = append(locations, r.walkForFile(dir, fileName, 3)...)
	}
	return dedupePaths(locations)
}

func (r *Resolver) searchDirs() []string {
	dirs := append([]string{}, r.cfg.SearchDirs...)
	if r.cfg.Getenv != nil {
		dirs = append(dirs, filepath.SplitList(strings.TrimSpace(r.cfg.Getenv("HOWL_ASR_SEARCH_DIRS")))...)
	}
	out := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(r.root(), dir)
		}
		out = append(out, filepath.Clean(dir))
	}
	return dedupePaths(out)
}

func (r *Resolver) walkForFile(root, fileName string, maxDepth int) []string {
	var out []string
	root = filepath.Clean(root)
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			rel, relErr := filepath.Rel(root, path)
			if relErr == nil && rel != "." && strings.Count(rel, string(os.PathSeparator)) >= maxDepth {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(entry.Name(), fileName) {
			out = append(out, path)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

func dedupePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, path := range paths {
		clean := filepath.Clean(path)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		out = append(out, clean)
	}
	return out
}

// ResolveProjector finds the audio projector paired with a resolved model. The
// projector is looked for beside the model first, then in the same configured
// inventory; an explicit HOWL_ASR_MMPROJ override wins.
func (r *Resolver) ResolveProjector(model ModelFile, profile *types.ASRModelProfile) (ModelFile, error) {
	if profile == nil || profile.ProjectorKind != types.ProjectorKindAudio {
		return ModelFile{}, nil
	}
	if override, ok := r.projectorOverride(); ok {
		if err := r.requireRegular(override); err != nil {
			return ModelFile{}, err
		}
		return ModelFile{Path: override, Type: model.Type, Source: ModelSourceEnv}, nil
	}

	// A companion beside its own model is the pairing. Identical filenames in two
	// configured directories are mirrors of one inventory, not two candidates:
	// the shared resolver treats that as ambiguity and refuses to load a model
	// whose projector is right next to it. Only a projector that is not beside
	// the model goes through the shared search, where ambiguity stays an error.
	if name := strings.TrimSpace(profile.ProjectorFile); name != "" {
		beside := filepath.Join(filepath.Dir(model.Path), name)
		if r.isRegular(beside) {
			return ModelFile{Path: beside, Type: model.Type, Source: ModelSourceFlat}, nil
		}
	}

	// The registry's ProjectorFile is an explicit identity: when present the
	// shared resolver looks for exactly that file and never substitutes a
	// different projector from the directory. Without it, the shared resolver
	// falls back to exact identity matching against the model file name.
	expected := make([]string, 0, 1)
	if name := strings.TrimSpace(profile.ProjectorFile); name != "" {
		expected = append(expected, name)
	}
	search := append([]string{
		filepath.Join(r.root(), "models"),
		filepath.Join(r.root(), "models", "whisper"),
	}, r.searchDirs()...)

	pairing, err := gguf.ResolveProjector(gguf.PairingRequest{
		ModelPath:        model.Path,
		ExpectedNames:    expected,
		Require:          true,
		RequiredModality: gguf.ProjectorAudio,
		SearchDirs:       search,
	})
	if err != nil {
		return ModelFile{}, fmt.Errorf("resolve audio projector: %w", err)
	}
	if !r.isRegular(pairing.Path) {
		return ModelFile{}, fmt.Errorf("%w: %s", ErrModelNotFound, pairing.Path)
	}
	return ModelFile{Path: pairing.Path, Type: model.Type, Source: ModelSourceFlat}, nil
}

func (r *Resolver) projectorOverride() (string, bool) {
	if r.cfg.Getenv == nil {
		return "", false
	}
	value := strings.TrimSpace(r.cfg.Getenv("HOWL_ASR_MMPROJ"))
	if value == "" {
		return "", false
	}
	if filepath.IsAbs(value) {
		return filepath.Clean(value), true
	}
	return filepath.Join(r.root(), value), true
}

func (r *Resolver) nestedPath(fileName string) string {
	return filepath.Join(r.root(), "models", "whisper", fileName)
}

func (r *Resolver) flatPath(fileName string) string {
	return filepath.Join(r.root(), "models", fileName)
}

func (r *Resolver) isRegular(path string) bool {
	info, err := r.stat(path)
	return err == nil && info != nil && info.Mode().IsRegular()
}

func (r *Resolver) requireRegular(path string) error {
	info, err := r.stat(path)
	if err != nil {
		return fmt.Errorf("%w: %s: %v", ErrModelNotFound, path, err)
	}
	if info == nil || !info.Mode().IsRegular() {
		return fmt.Errorf("%w: %s", ErrModelNotFound, path)
	}
	return nil
}

// stat is the only path through which the resolver touches the filesystem.
func (r *Resolver) stat(path string) (os.FileInfo, error) {
	if r.cfg.Stat != nil {
		return r.cfg.Stat(path)
	}
	return os.Stat(path)
}

func sourceRank(source ModelSource) int {
	switch source {
	case ModelSourceNested:
		return 0
	case ModelSourceFlat:
		return 1
	default:
		return 2
	}
}

var (
	installedResolverOnce sync.Once
	installedResolver     *Resolver
)

// resolveInstalledASRModel is the production discovery the pipeline and the
// capability contract share. It returns "" when the model is not installed;
// callers then report the expected path instead of inventing a load path.
func resolveInstalledASRModel(modelType types.ASRModelType) string {
	resolver := productionResolver()
	if resolver == nil {
		return ""
	}
	file, err := resolver.Resolve(string(modelType), nil)
	if err != nil {
		return ""
	}
	return file.Path
}

// resolveInstalledASRProjector returns the matching audio projector for a
// model whose main file was resolved. Empty means "not installed", never "not
// required".
func resolveInstalledASRProjector(modelType types.ASRModelType) string {
	profile, ok := types.GetASRProfile(modelType)
	if !ok || profile.ProjectorKind != types.ProjectorKindAudio {
		return ""
	}
	resolver := productionResolver()
	if resolver == nil {
		return ""
	}
	main, err := resolver.Resolve(string(modelType), nil)
	if err != nil {
		return ""
	}
	projector, err := resolver.ResolveProjector(main, profile)
	if err != nil {
		return ""
	}
	return projector.Path
}

func productionResolver() *Resolver {
	installedResolverOnce.Do(func() {
		resolver, err := NewResolver(ResolverConfig{
			Root:       paths.ProjectRoot(),
			Getenv:     os.Getenv,
			SearchDirs: configuredASRSearchDirs(),
		})
		if err != nil {
			return
		}
		installedResolver = resolver
	})
	return installedResolver
}

// configuredASRSearchDirs reads the directory-agnostic ASR inventory. Paths
// live in config, never in code, so moving models between drives needs no
// rebuild.
func configuredASRSearchDirs() []string {
	root := paths.ProjectRoot()
	var dirs []string
	if data, err := os.ReadFile(filepath.Join(root, "config", "asr-model-dirs.json")); err == nil {
		var list []string
		if json.Unmarshal(data, &list) == nil {
			dirs = append(dirs, list...)
		}
	}
	if data, err := os.ReadFile(filepath.Join(root, "config", "model-scan-dirs.json")); err == nil {
		var byType map[string]string
		if json.Unmarshal(data, &byType) == nil {
			for dir := range byType {
				dirs = append(dirs, dir)
			}
		}
	}
	return dirs
}

func approvedASRModelFile(modelType types.ASRModelType) (string, bool) {
	for _, model := range approvedASRModels {
		if model.Type == modelType {
			return model.File, true
		}
	}
	return "", false
}

func isApprovedASRModelType(modelType types.ASRModelType) bool {
	_, ok := approvedASRModelFile(modelType)
	return ok
}

func candidatePathExists(candidates []ModelFile, path string) bool {
	clean := filepath.Clean(path)
	for _, candidate := range candidates {
		if filepath.Clean(candidate.Path) == clean {
			return true
		}
	}
	return false
}
