package worldstore

import "strings"

// RefKind is the category a reference points at. References are typed so a
// lorebook identity can never be resolved as a character or a world.
type RefKind string

const (
	RefWorld     RefKind = "world"
	RefScenario  RefKind = "scenario"
	RefCharacter RefKind = "character"
	RefLorebook  RefKind = "lorebook"
)

// Reference is the frozen reference shape. Resolution uses Kind and ID only;
// Title is a display cache and LegacyName marks a name-only record awaiting an
// explicit mapping. A display name never authorizes or resolves a reference.
type Reference struct {
	Kind       RefKind `json:"kind"`
	ID         string  `json:"id"`
	Title      string  `json:"title,omitempty"`
	LegacyName string  `json:"legacyName,omitempty"`
}

// Catalog lists the known identities per kind for resolution.
type Catalog map[RefKind]map[string]bool

// Resolution reports whether a reference resolved and why not when it did not.
type Resolution struct {
	Reference
	Resolved bool   `json:"resolved"`
	Reason   string `json:"reason,omitempty"`
}

// NormalizeReferences converts a stored reference list into the typed shape.
// Entries with an identity keep it; name-only entries become LegacyName and stay
// unresolved. Duplicates collapse on identity, or on legacy name when no identity
// exists, preserving the first title seen.
func NormalizeReferences(kind RefKind, raw interface{}) []Reference {
	var items []interface{}
	switch typed := raw.(type) {
	case []interface{}:
		items = typed
	case []map[string]interface{}:
		items = make([]interface{}, len(typed))
		for i, item := range typed {
			items[i] = item
		}
	case nil:
		return nil
	default:
		return nil
	}

	out := make([]Reference, 0, len(items))
	seen := make(map[string]int, len(items))
	for _, item := range items {
		ref, ok := normalizeReference(kind, item)
		if !ok {
			continue
		}
		key := referenceKey(ref)
		if index, exists := seen[key]; exists {
			if out[index].Title == "" {
				out[index].Title = ref.Title
			}
			continue
		}
		seen[key] = len(out)
		out = append(out, ref)
	}
	return out
}

func normalizeReference(kind RefKind, item interface{}) (Reference, bool) {
	switch typed := item.(type) {
	case string:
		name := strings.TrimSpace(typed)
		if name == "" {
			return Reference{}, false
		}
		return Reference{Kind: kind, LegacyName: name}, true
	case map[string]interface{}:
		explicitKind := strings.TrimSpace(stringField(typed, "kind"))
		if explicitKind != "" && explicitKind != string(kind) {
			return Reference{}, false
		}
		id := strings.TrimSpace(stringField(typed, "id"))
		title := strings.TrimSpace(stringField(typed, "title"))
		name := strings.TrimSpace(stringField(typed, "name"))
		legacy := strings.TrimSpace(stringField(typed, "legacyName"))
		if title == "" {
			title = name
		}
		if id != "" {
			return Reference{Kind: kind, ID: id, Title: title, LegacyName: legacy}, true
		}
		if legacy == "" {
			legacy = title
		}
		if legacy == "" {
			return Reference{}, false
		}
		return Reference{Kind: kind, LegacyName: legacy, Title: title}, true
	default:
		return Reference{}, false
	}
}

func referenceKey(ref Reference) string {
	if ref.ID != "" {
		return string(ref.Kind) + "\x00id:" + ref.ID
	}
	return string(ref.Kind) + "\x00legacy:" + strings.ToLower(ref.LegacyName)
}

// ResolveReferences resolves each reference by Kind and ID only. A legacy
// name-only reference is never matched to a catalog item, even when a title
// happens to be identical; that would silently substitute one item for another.
func ResolveReferences(refs []Reference, catalog Catalog) []Resolution {
	out := make([]Resolution, 0, len(refs))
	for _, ref := range refs {
		resolution := Resolution{Reference: ref}
		switch {
		case ref.ID == "":
			resolution.Reason = "legacy_name_unresolved"
		case catalog[ref.Kind][ref.ID]:
			resolution.Resolved = true
		default:
			resolution.Reason = "not_found"
		}
		out = append(out, resolution)
	}
	return out
}

// UnresolvedReferences returns the references that did not resolve, in input
// order. Callers display them as unavailable and exclude them from prompts.
func UnresolvedReferences(refs []Reference, catalog Catalog) []Reference {
	var out []Reference
	for _, resolution := range ResolveReferences(refs, catalog) {
		if !resolution.Resolved {
			out = append(out, resolution.Reference)
		}
	}
	return out
}
