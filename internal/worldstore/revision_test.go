package worldstore

import (
	"encoding/json"
	"testing"
)

func TestRevisionIgnoresKeyOrderAndNumberSpelling(t *testing.T) {
	a := map[string]interface{}{
		"title":  "Ashfall",
		"count":  1,
		"ratio":  1.5,
		"nested": map[string]interface{}{"b": true, "a": nil},
	}
	b := map[string]interface{}{
		"nested": map[string]interface{}{"a": nil, "b": true},
		"ratio":  json.Number("1.50"),
		"count":  json.Number("1.0"),
		"title":  "Ashfall",
	}
	if Revision(a) != Revision(b) {
		t.Fatalf("equivalent documents produced different revisions: %q vs %q", Revision(a), Revision(b))
	}
}

func TestRevisionIgnoresTimestampsAndMediaPayload(t *testing.T) {
	a := map[string]interface{}{
		"title":     "Ashfall",
		"updatedAt": "2026-01-01T00:00:00Z",
		"images": []interface{}{
			map[string]interface{}{"id": "img_1", "data": "AAA", "path": "C:/one.png"},
		},
		"metadata": map[string]interface{}{"updated": "2026-01-01T00:00:00Z", "seed": 7},
	}
	b := map[string]interface{}{
		"title":     "Ashfall",
		"updatedAt": "2026-12-31T23:59:59Z",
		"images": []interface{}{
			map[string]interface{}{"id": "img_1", "data": "BBB", "path": "D:/two.png"},
		},
		"metadata": map[string]interface{}{"updated": "2026-12-31T23:59:59Z", "seed": 7},
	}
	if Revision(a) != Revision(b) {
		t.Fatalf("timestamps or media payload changed the revision: %q vs %q", Revision(a), Revision(b))
	}
}

func TestRevisionChangesWithContent(t *testing.T) {
	base := map[string]interface{}{"title": "Ashfall", "refs": []interface{}{"a", "b"}}
	if Revision(base) == Revision(map[string]interface{}{"title": "Ashfall II", "refs": []interface{}{"a", "b"}}) {
		t.Fatal("title change did not change the revision")
	}
	if Revision(base) == Revision(map[string]interface{}{"title": "Ashfall", "refs": []interface{}{"b", "a"}}) {
		t.Fatal("array order change did not change the revision")
	}
	if Revision(map[string]interface{}{"a": nil}) == Revision(map[string]interface{}{}) {
		t.Fatal("null and omission produced the same revision")
	}
}

func TestRevisionKeepsLargeIntegerPrecision(t *testing.T) {
	a := map[string]interface{}{"n": json.Number("9007199254740993")}
	b := map[string]interface{}{"n": json.Number("9007199254740992")}
	if Revision(a) == Revision(b) {
		t.Fatal("integers beyond float64 precision collapsed to one revision")
	}
}

func TestRevisionDoesNotNormalizeUnicode(t *testing.T) {
	composed := map[string]interface{}{"title": "\u00fcber"}
	decomposed := map[string]interface{}{"title": "u\u0308ber"}
	if Revision(composed) == Revision(decomposed) {
		t.Fatal("distinct Unicode sequences produced the same revision")
	}
}
