package audio

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/audio/types"
)

type testPolicy func([]ModelFile) ([]ModelFile, error)

func (f testPolicy) Candidates(installed []ModelFile) ([]ModelFile, error) {
	return f(installed)
}

func envMap(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func writeModelFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create model directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("model"), 0o644); err != nil {
		t.Fatalf("write model file: %v", err)
	}
}

func mustResolver(t *testing.T, cfg ResolverConfig) *Resolver {
	t.Helper()
	resolver, err := NewResolver(cfg)
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}
	return resolver
}

func TestNewResolverRequiresExplicitConfiguration(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name string
		cfg  ResolverConfig
		want string
	}{
		{
			name: "empty root",
			cfg:  ResolverConfig{Getenv: envMap(nil)},
			want: "root is required",
		},
		{
			name: "relative root",
			cfg:  ResolverConfig{Root: "relative-root", Getenv: envMap(nil)},
			want: "root must be absolute",
		},
		{
			name: "nil Getenv",
			cfg:  ResolverConfig{Root: root},
			want: "Getenv is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolver, err := NewResolver(test.cfg)
			if err == nil {
				t.Fatal("NewResolver accepted invalid configuration")
			}
			if resolver != nil {
				t.Fatalf("NewResolver returned %+v with an error", resolver)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want message containing %q", err, test.want)
			}
		})
	}
}

func TestResolverApprovedExplicitEngineSet(t *testing.T) {
	root := t.TempDir()
	approved := []struct {
		engine string
		file   string
	}{
		{engine: "whisper-tiny", file: "ggml-tiny.bin"},
		{engine: "whisper-base", file: "ggml-base.bin"},
		{engine: "whisper-small", file: "ggml-small.bin"},
	}
	for _, test := range approved {
		writeModelFile(t, filepath.Join(root, "models", "whisper", test.file))
	}

	resolver := mustResolver(t, ResolverConfig{Root: root, Getenv: envMap(nil)})
	for _, test := range approved {
		t.Run(test.engine, func(t *testing.T) {
			file, err := resolver.Resolve(test.engine, nil)
			if err != nil {
				t.Fatalf("Resolve: %v", err)
			}
			wantPath := filepath.Join(root, "models", "whisper", test.file)
			if file.Path != wantPath {
				t.Fatalf("path = %q, want %q", file.Path, wantPath)
			}
			if file.Type != types.ASRModelType(test.engine) {
				t.Fatalf("type = %q, want %q", file.Type, test.engine)
			}
			if file.Source != ModelSourceNested {
				t.Fatalf("source = %q, want %q", file.Source, ModelSourceNested)
			}
		})
	}
}

func TestResolverRejectsUnsupportedAndBlankEngines(t *testing.T) {
	resolver := mustResolver(t, ResolverConfig{Root: t.TempDir(), Getenv: envMap(nil)})
	tests := []struct {
		name   string
		engine string
	}{
		{name: "medium", engine: "whisper-medium"},
		{name: "large", engine: "whisper-large-v3"},
		{name: "qwen", engine: "qwen3-asr-0.6b"},
		{name: "empty", engine: ""},
		{name: "whitespace", engine: "  \t\n "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := resolver.Resolve(test.engine, nil)
			if !errors.Is(err, ErrEngineUnsupported) {
				t.Fatalf("error = %v, want ErrEngineUnsupported", err)
			}
			if errors.Is(err, ErrModelNotFound) {
				t.Fatalf("error = %v, must not be ErrModelNotFound", err)
			}
			if file != (ModelFile{}) {
				t.Fatalf("file = %+v, want zero value", file)
			}
		})
	}
}

func TestResolverAutoDiscoveryExcludesUnapprovedModels(t *testing.T) {
	root := t.TempDir()
	writeModelFile(t, filepath.Join(root, "models", "whisper", "ggml-base.bin"))
	writeModelFile(t, filepath.Join(root, "models", "whisper", "ggml-medium.bin"))
	writeModelFile(t, filepath.Join(root, "models", "whisper", "ggml-large-v3.bin"))

	resolver := mustResolver(t, ResolverConfig{Root: root, Getenv: envMap(nil)})
	var got []ModelFile
	policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
		got = append([]ModelFile(nil), candidates...)
		for _, candidate := range candidates {
			if candidate.Type == types.ASRWhisperBase {
				return []ModelFile{candidate}, nil
			}
		}
		return nil, nil
	})

	file, err := resolver.Resolve("auto", policy)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(got) != 1 || got[0].Type != types.ASRWhisperBase {
		t.Fatalf("candidates = %+v, want only whisper-base", got)
	}
	if file.Type != types.ASRWhisperBase {
		t.Fatalf("resolved type = %q, want %q", file.Type, types.ASRWhisperBase)
	}
}

func TestResolverExplicitMissingModelReportsBothPaths(t *testing.T) {
	root := t.TempDir()
	resolver := mustResolver(t, ResolverConfig{Root: root, Getenv: envMap(nil)})

	file, err := resolver.Resolve("whisper-base", nil)
	if !errors.Is(err, ErrModelNotFound) {
		t.Fatalf("error = %v, want ErrModelNotFound", err)
	}
	if file != (ModelFile{}) {
		t.Fatalf("file = %+v, want zero value", file)
	}

	nested := filepath.Join(root, "models", "whisper", "ggml-base.bin")
	flat := filepath.Join(root, "models", "ggml-base.bin")
	for _, path := range []string{nested, flat} {
		if !filepath.IsAbs(path) {
			t.Fatalf("test path %q is not absolute", path)
		}
		if !strings.Contains(err.Error(), path) {
			t.Fatalf("error = %v, want attempted path %q", err, path)
		}
	}
}

func TestResolverExplicitEnginePrecedenceAndIdentity(t *testing.T) {
	t.Run("override wins and keeps explicit type", func(t *testing.T) {
		root := t.TempDir()
		nestedBase := filepath.Join(root, "models", "whisper", "ggml-base.bin")
		flatBase := filepath.Join(root, "models", "ggml-base.bin")
		override := filepath.Join(root, "override", "tiny.bin")
		writeModelFile(t, nestedBase)
		writeModelFile(t, flatBase)
		writeModelFile(t, override)

		resolver := mustResolver(t, ResolverConfig{
			Root:   root,
			Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": override}),
		})
		file, err := resolver.Resolve("whisper-base", nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if file.Path != override {
			t.Fatalf("path = %q, want override %q", file.Path, override)
		}
		if file.Type != types.ASRWhisperBase {
			t.Fatalf("type = %q, want explicit engine %q", file.Type, types.ASRWhisperBase)
		}
		if file.Source != ModelSourceEnv {
			t.Fatalf("source = %q, want %q", file.Source, ModelSourceEnv)
		}
	})

	t.Run("nested precedes flat", func(t *testing.T) {
		root := t.TempDir()
		nestedBase := filepath.Join(root, "models", "whisper", "ggml-base.bin")
		flatBase := filepath.Join(root, "models", "ggml-base.bin")
		writeModelFile(t, nestedBase)
		writeModelFile(t, flatBase)

		resolver := mustResolver(t, ResolverConfig{Root: root, Getenv: envMap(nil)})
		file, err := resolver.Resolve("whisper-base", nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if file.Path != nestedBase || file.Source != ModelSourceNested {
			t.Fatalf("resolved %q (%s), want nested %q", file.Path, file.Source, nestedBase)
		}
	})

	t.Run("flat fallback", func(t *testing.T) {
		root := t.TempDir()
		flatBase := filepath.Join(root, "models", "ggml-base.bin")
		writeModelFile(t, flatBase)

		resolver := mustResolver(t, ResolverConfig{Root: root, Getenv: envMap(nil)})
		file, err := resolver.Resolve("whisper-base", nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if file.Path != flatBase || file.Source != ModelSourceFlat {
			t.Fatalf("resolved %q (%s), want flat %q", file.Path, file.Source, flatBase)
		}
	})

	t.Run("relative override is rooted", func(t *testing.T) {
		root := t.TempDir()
		relative := filepath.Join("override", "small.bin")
		absolute := filepath.Join(root, relative)
		writeModelFile(t, absolute)

		resolver := mustResolver(t, ResolverConfig{
			Root:   root,
			Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": relative}),
		})
		file, err := resolver.Resolve("whisper-small", nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if file.Path != absolute {
			t.Fatalf("path = %q, want %q", file.Path, absolute)
		}
		if !filepath.IsAbs(file.Path) {
			t.Fatalf("path = %q, want absolute", file.Path)
		}
	})

	t.Run("missing override fails closed", func(t *testing.T) {
		root := t.TempDir()
		writeModelFile(t, filepath.Join(root, "models", "whisper", "ggml-base.bin"))
		missing := filepath.Join(root, "missing", "model.bin")

		resolver := mustResolver(t, ResolverConfig{
			Root:   root,
			Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": missing}),
		})
		file, err := resolver.Resolve("whisper-base", nil)
		if !errors.Is(err, ErrModelNotFound) {
			t.Fatalf("error = %v, want ErrModelNotFound", err)
		}
		if file != (ModelFile{}) {
			t.Fatalf("file = %+v, want zero value", file)
		}
		if !strings.Contains(err.Error(), missing) {
			t.Fatalf("error = %v, want absolute path %q", err, missing)
		}
	})

	t.Run("unsupported explicit engine", func(t *testing.T) {
		resolver := mustResolver(t, ResolverConfig{Root: t.TempDir(), Getenv: envMap(nil)})
		if _, err := resolver.Resolve("qwen3-asr-0.6b", nil); !errors.Is(err, ErrEngineUnsupported) {
			t.Fatalf("error = %v, want ErrEngineUnsupported", err)
		}
	})
}

func TestResolverAutoPolicyStates(t *testing.T) {
	root := t.TempDir()
	nestedBase := filepath.Join(root, "models", "whisper", "ggml-base.bin")
	override := filepath.Join(root, "override", "custom.bin")
	writeModelFile(t, nestedBase)
	writeModelFile(t, override)
	resolver := mustResolver(t, ResolverConfig{
		Root:   root,
		Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": override}),
	})

	t.Run("untyped override may be assigned a known type", func(t *testing.T) {
		policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
			for _, candidate := range candidates {
				if candidate.Source == ModelSourceEnv {
					candidate.Type = types.ASRWhisperSmall
					return []ModelFile{candidate}, nil
				}
			}
			return nil, nil
		})
		file, err := resolver.Resolve("auto", policy)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if file.Path != override || file.Type != types.ASRWhisperSmall || file.Source != ModelSourceEnv {
			t.Fatalf("resolved %+v, want env override with declared type", file)
		}
	})

	t.Run("untyped override may not be assigned an unapproved type", func(t *testing.T) {
		policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
			for _, candidate := range candidates {
				if candidate.Source == ModelSourceEnv {
					candidate.Type = types.ASRModelType("whisper-medium")
					return []ModelFile{candidate}, nil
				}
			}
			return nil, nil
		})
		if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrPolicyInvalidSelection) {
			t.Fatalf("error = %v, want ErrPolicyInvalidSelection", err)
		}
	})

	t.Run("canonical candidate keeps its identity", func(t *testing.T) {
		policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
			for _, candidate := range candidates {
				if candidate.Type == types.ASRWhisperBase {
					return []ModelFile{candidate}, nil
				}
			}
			return nil, nil
		})
		file, err := resolver.Resolve("auto", policy)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if file.Path != nestedBase || file.Type != types.ASRWhisperBase {
			t.Fatalf("resolved %+v, want canonical base", file)
		}
	})

	t.Run("canonical candidate may not be retyped", func(t *testing.T) {
		policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
			for _, candidate := range candidates {
				if candidate.Type == types.ASRWhisperBase {
					candidate.Type = types.ASRWhisperSmall
					return []ModelFile{candidate}, nil
				}
			}
			return nil, nil
		})
		if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrPolicyInvalidSelection) {
			t.Fatalf("error = %v, want ErrPolicyInvalidSelection", err)
		}
	})

	t.Run("foreign path is rejected", func(t *testing.T) {
		policy := testPolicy(func([]ModelFile) ([]ModelFile, error) {
			return []ModelFile{{Path: filepath.Join(root, "foreign.bin"), Type: types.ASRWhisperBase}}, nil
		})
		if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrPolicyInvalidSelection) {
			t.Fatalf("error = %v, want ErrPolicyInvalidSelection", err)
		}
	})

	t.Run("untyped override requires a known type", func(t *testing.T) {
		policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
			for _, candidate := range candidates {
				if candidate.Source == ModelSourceEnv {
					return []ModelFile{candidate}, nil
				}
			}
			return nil, nil
		})
		if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrPolicyInvalidSelection) {
			t.Fatalf("error = %v, want ErrPolicyInvalidSelection", err)
		}
	})

	t.Run("conflicting duplicate is rejected", func(t *testing.T) {
		policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
			for _, candidate := range candidates {
				if candidate.Source == ModelSourceEnv {
					first := candidate
					first.Type = types.ASRWhisperSmall
					second := candidate
					second.Type = types.ASRWhisperBase
					return []ModelFile{first, second}, nil
				}
			}
			return nil, nil
		})
		if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrPolicyInvalidSelection) {
			t.Fatalf("error = %v, want ErrPolicyInvalidSelection", err)
		}
	})

	t.Run("empty policy result", func(t *testing.T) {
		policy := testPolicy(func([]ModelFile) ([]ModelFile, error) { return nil, nil })
		if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrNoModelSelected) {
			t.Fatalf("error = %v, want ErrNoModelSelected", err)
		}
	})

	t.Run("nil policy", func(t *testing.T) {
		if _, err := resolver.Resolve("auto", nil); !errors.Is(err, ErrPolicyNotConfigured) {
			t.Fatalf("error = %v, want ErrPolicyNotConfigured", err)
		}
	})
}

func TestResolverAutoMissingOverrideFailsBeforePolicy(t *testing.T) {
	root := t.TempDir()
	writeModelFile(t, filepath.Join(root, "models", "whisper", "ggml-base.bin"))
	missing := filepath.Join(root, "missing", "model.bin")

	resolver := mustResolver(t, ResolverConfig{
		Root:   root,
		Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": missing}),
	})
	policyCalls := 0
	policy := testPolicy(func([]ModelFile) ([]ModelFile, error) {
		policyCalls++
		return nil, nil
	})
	_, err := resolver.Resolve("auto", policy)
	if !errors.Is(err, ErrModelNotFound) {
		t.Fatalf("error = %v, want ErrModelNotFound", err)
	}
	if policyCalls != 0 {
		t.Fatalf("policy called %d times, want 0", policyCalls)
	}
}

func TestResolverAutoEmptyCandidateSet(t *testing.T) {
	resolver := mustResolver(t, ResolverConfig{Root: t.TempDir(), Getenv: envMap(nil)})
	policy := testPolicy(func([]ModelFile) ([]ModelFile, error) {
		t.Fatal("policy must not run without candidates")
		return nil, nil
	})
	if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrModelNotFound) {
		t.Fatalf("error = %v, want ErrModelNotFound", err)
	}
}

func TestResolverAutoCandidateOrderAndDedup(t *testing.T) {
	root := t.TempDir()
	nestedBase := filepath.Join(root, "models", "whisper", "ggml-base.bin")
	nestedTiny := filepath.Join(root, "models", "whisper", "ggml-tiny.bin")
	flatSmall := filepath.Join(root, "models", "ggml-small.bin")
	override := filepath.Join(root, "override", "custom.bin")
	for _, path := range []string{nestedBase, nestedTiny, flatSmall, override} {
		writeModelFile(t, path)
	}

	var got []ModelFile
	policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
		got = append([]ModelFile(nil), candidates...)
		return nil, nil
	})
	resolver := mustResolver(t, ResolverConfig{
		Root:   root,
		Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": override}),
	})
	if _, err := resolver.Resolve("auto", policy); !errors.Is(err, ErrNoModelSelected) {
		t.Fatalf("error = %v, want ErrNoModelSelected", err)
	}
	want := []string{nestedBase, nestedTiny, flatSmall, override}
	if len(got) != len(want) {
		t.Fatalf("candidate count = %d, want %d: %+v", len(got), len(want), got)
	}
	for i, path := range want {
		if got[i].Path != path {
			t.Fatalf("candidate %d = %q, want %q", i, got[i].Path, path)
		}
	}

	dedup := mustResolver(t, ResolverConfig{
		Root:   root,
		Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": nestedTiny}),
	})
	got = nil
	if _, err := dedup.Resolve("auto", policy); !errors.Is(err, ErrNoModelSelected) {
		t.Fatalf("dedup error = %v, want ErrNoModelSelected", err)
	}
	if len(got) != 3 {
		t.Fatalf("dedup candidate count = %d, want 3: %+v", len(got), got)
	}
}

func TestResolverAutoNilPolicyDoesNotStat(t *testing.T) {
	calls := 0
	resolver := mustResolver(t, ResolverConfig{
		Root:   t.TempDir(),
		Getenv: envMap(nil),
		Stat: func(string) (os.FileInfo, error) {
			calls++
			t.Errorf("Stat called during nil-policy resolution")
			return nil, os.ErrNotExist
		},
	})
	file, err := resolver.Resolve("auto", nil)
	if !errors.Is(err, ErrPolicyNotConfigured) {
		t.Fatalf("error = %v, want ErrPolicyNotConfigured", err)
	}
	if file != (ModelFile{}) {
		t.Fatalf("file = %+v, want zero value", file)
	}
	if calls != 0 {
		t.Fatalf("Stat calls = %d, want 0", calls)
	}
}

func TestResolverAutoPolicyUsesInjectedStat(t *testing.T) {
	root := t.TempDir()
	nestedBase := filepath.Join(root, "models", "whisper", "ggml-base.bin")
	writeModelFile(t, nestedBase)

	calls := 0
	resolver := mustResolver(t, ResolverConfig{
		Root:   root,
		Getenv: envMap(nil),
		Stat: func(path string) (os.FileInfo, error) {
			calls++
			return os.Stat(path)
		},
	})
	policy := testPolicy(func(candidates []ModelFile) ([]ModelFile, error) {
		for _, candidate := range candidates {
			if candidate.Type == types.ASRWhisperBase {
				return []ModelFile{candidate}, nil
			}
		}
		return nil, nil
	})
	file, err := resolver.Resolve("auto", policy)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if file.Path != nestedBase {
		t.Fatalf("path = %q, want %q", file.Path, nestedBase)
	}
	if calls == 0 {
		t.Fatal("Stat calls = 0, want positive control to observe the injected seam")
	}
}

func TestResolverDanglingSymlinkFailsClosed(t *testing.T) {
	root := t.TempDir()
	link := filepath.Join(root, "dangling.bin")
	if err := os.Symlink(filepath.Join(root, "missing-target.bin"), link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	resolver := mustResolver(t, ResolverConfig{
		Root:   root,
		Getenv: envMap(map[string]string{"HOWL_ASR_MODEL": link}),
	})
	file, err := resolver.Resolve("whisper-base", nil)
	if !errors.Is(err, ErrModelNotFound) {
		t.Fatalf("error = %v, want ErrModelNotFound", err)
	}
	if file != (ModelFile{}) {
		t.Fatalf("file = %+v, want zero value", file)
	}
	if !strings.Contains(err.Error(), link) {
		t.Fatalf("error = %v, want symlink path %q", err, link)
	}
}
