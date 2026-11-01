package registry

import "testing"

func TestIsUnderDir(t *testing.T) {
	cases := []struct {
		path, dir string
		want      bool
	}{
		{`E:\AI GGUF\GGUF\a.gguf`, `E:\AI GGUF\GGUF`, true},
		{`E:\AI GGUF\GGUF`, `E:\AI GGUF\GGUF`, true},
		// The prefix trap: a sibling directory that merely starts with the name.
		{`E:\AI GGUF\GGUF2\a.gguf`, `E:\AI GGUF\GGUF`, false},
		{`E:\AI GGUF\GGUFETA`, `E:\AI GGUF\GGUF`, false},
		// Windows paths compare case-insensitively.
		{`e:\ai gguf\gguf\a.gguf`, `E:\AI GGUF\GGUF`, true},
		{`C:\somewhere\a.gguf`, `E:\AI GGUF\GGUF`, false},
		{``, `E:\AI GGUF\GGUF`, false},
		{`E:\AI GGUF\GGUF\a.gguf`, ``, false},
	}
	for _, c := range cases {
		if got := IsUnderDir(c.path, c.dir); got != c.want {
			t.Errorf("IsUnderDir(%q, %q) = %v, want %v", c.path, c.dir, got, c.want)
		}
	}
}

func TestKindForPathPrefersDeepestDirectory(t *testing.T) {
	r := &Registry{
		dirs: map[string]string{
			`E:\models`:    KindLLM,
			`E:\models\sd`: KindImage,
		},
		entries: map[string]Entry{},
	}
	cases := map[string]string{
		`E:\models\sd\flux1-dev.gguf`: KindImage,
		`E:\models\chat\qwen3.gguf`:   KindLLM,
		`E:\models\sd2\x.gguf`:        KindLLM, // sibling of sd, still inside models
		`E:\models2\x.gguf`:           "",      // prefix trap, not inside models
		`D:\unregistered\x.gguf`:      "",
	}
	for path, want := range cases {
		if got := r.KindForPath(path); got != want {
			t.Errorf("KindForPath(%q) = %q, want %q", path, got, want)
		}
	}
}
