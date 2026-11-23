package chat

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// This test is the enforcement half of the mutation boundary. A startup probe
// would only cover the paths someone remembered; this fails the build whenever a
// filesystem write appears anywhere in the package except the four helpers in
// writers.go, which is a property of the code rather than of a scenario.

// filesystemWriteFunctions are the os calls that mutate the filesystem. Reads are
// deliberately absent: only writes are boundary-controlled.
var filesystemWriteFunctions = map[string]bool{
	"WriteFile": true,
	"OpenFile":  true,
	"Create":    true,
	"Truncate":  true,
	"Rename":    true,
	"MkdirAll":  true,
	"Mkdir":     true,
	"Remove":    true,
	"RemoveAll": true,
	"Chmod":     true,
	"Chtimes":   true,
}

// boundaryWriterFiles maps each file to the functions that may call them. The
// recycle files are platform-split on purpose: there is no runtime fallback from
// Recycle Bin deletion to permanent deletion, only a different build target.
var boundaryWriterFiles = map[string]map[string]bool{
	"writers.go": {
		"ensureDir":       true,
		"writeFileAtomic": true,
		"appendLine":      true,
		"claimDir":        true,
	},
	"recycle_windows.go": {"recycleTree": true},
	"recycle_other.go":   {"recycleTree": true},
}

type writeViolation struct {
	File     string
	Function string
	Call     string
}

// collectWriteViolations parses one source file and reports every filesystem write
// made outside the boundary helpers. src may be a filename hint for error messages
// or nil to read from disk; a []byte or string can be supplied for a synthetic file.
func collectWriteViolations(file string, src interface{}) ([]writeViolation, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, file, src, 0)
	if err != nil {
		return nil, err
	}

	var violations []writeViolation
	for _, decl := range parsed.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		function := fn.Name.Name
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			selector, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := selector.X.(*ast.Ident)
			if !ok || ident.Name != "os" || !filesystemWriteFunctions[selector.Sel.Name] {
				return true
			}
			if boundaryWriterFiles[filepath.Base(file)][function] {
				return true
			}
			violations = append(violations, writeViolation{
				File:     filepath.Base(file),
				Function: function,
				Call:     "os." + selector.Sel.Name,
			})
			return true
		})
	}
	return violations, nil
}

func TestWritesOnlyThroughBoundaryHelpers(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob package files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no Go files found; the analysis test is not running in the package directory")
	}

	var all []writeViolation
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		found, err := collectWriteViolations(file, source)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		all = append(all, found...)
	}

	if len(all) > 0 {
		sort.Slice(all, func(i, j int) bool { return all[i].Function < all[j].Function })
		var lines []string
		for _, violation := range all {
			lines = append(lines, violation.File+" "+violation.Function+" calls "+violation.Call)
		}
		t.Fatalf("filesystem writes outside the boundary helpers (use writers.go instead):\n%s", strings.Join(lines, "\n"))
	}
}

// TestAnalysisDetectsABypass proves the checker is capable of failing. A test that
// has never rejected anything is not evidence of enforcement (PA-T1).
func TestAnalysisDetectsABypass(t *testing.T) {
	const synthetic = `package chat

import "os"

func bypass(path string) error {
	return os.RemoveAll(path)
}
`
	violations, err := collectWriteViolations("synthetic.go", synthetic)
	if err != nil {
		t.Fatalf("parse synthetic source: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected exactly one violation from the synthetic bypass, got %d: %+v", len(violations), violations)
	}
	if violations[0].Function != "bypass" || violations[0].Call != "os.RemoveAll" {
		t.Fatalf("unexpected violation: %+v", violations[0])
	}
}
