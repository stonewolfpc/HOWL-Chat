package contentimport

import (
	"fmt"
	"strings"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

func knownReferenceKind(kind string) bool {
	switch worldstore.RefKind(kind) {
	case worldstore.RefWorld, worldstore.RefScenario, worldstore.RefCharacter, worldstore.RefLorebook:
		return true
	default:
		return false
	}
}

// normalizeReferences converts supported reference fields to the frozen typed
// shape while retaining unknown per-reference fields. Name-only entries become
// legacyName-only references; malformed and unknown-kind entries are preserved
// and reported rather than discarded.
func (s *Service) normalizeReferences(document map[string]interface{}) []ReferenceIssue {
	var issues []ReferenceIssue
	add := func(list []ReferenceIssue) {
		issues = append(issues, list...)
	}

	if s.category == library.Worlds {
		add(normalizeReferenceField(document, "scenarios", worldstore.RefScenario))
		add(normalizeReferenceField(document, "characters", worldstore.RefCharacter))
		add(normalizeLorebookFields(document, true))
		return issues
	}

	add(normalizeReferenceField(document, "characters", worldstore.RefCharacter))
	if howl, ok := document["howl"].(map[string]interface{}); ok {
		if raw, exists := howl["world"]; exists {
			value, issue := normalizeReferenceEntry("howl.world", 0, raw, worldstore.RefWorld)
			if issue != nil {
				add([]ReferenceIssue{*issue})
			}
			if issue == nil || issue.Reason == "legacy_name" || issue.Reason == "missing_id" {
				howl["world"] = value
			}
		}
	}
	add(normalizeLorebookFields(document, false))
	return issues
}

func normalizeReferenceField(document map[string]interface{}, field string, kind worldstore.RefKind) []ReferenceIssue {
	raw, exists := document[field]
	if !exists {
		return nil
	}
	values, issues, ok := normalizeReferenceList(field, raw, kind)
	if !ok {
		return issues
	}
	document[field] = values
	return issues
}

func normalizeReferenceList(field string, raw interface{}, kind worldstore.RefKind) ([]interface{}, []ReferenceIssue, bool) {
	var items []interface{}
	switch typed := raw.(type) {
	case []interface{}:
		items = typed
	case []map[string]interface{}:
		items = make([]interface{}, len(typed))
		for i, item := range typed {
			items[i] = item
		}
	default:
		return nil, []ReferenceIssue{{
			Field: field, Index: -1, Reason: "malformed",
			Message: "The reference field is not an array; the original value was preserved.",
		}}, false
	}

	out := make([]interface{}, 0, len(items))
	seen := map[string]int{}
	var issues []ReferenceIssue
	for i, item := range items {
		value, issue := normalizeReferenceEntry(field, i, item, kind)
		if issue != nil {
			issues = append(issues, *issue)
			if issue.Reason == "unknown_kind" || issue.Reason == "kind_mismatch" || issue.Reason == "malformed" {
				out = append(out, value)
				continue
			}
		}
		key := referenceKey(value)
		if key != "" {
			if index, exists := seen[key]; exists {
				if existing, ok := out[index].(map[string]interface{}); ok {
					if incoming, ok := value.(map[string]interface{}); ok {
						mergeReferenceFields(existing, incoming)
					}
				}
				issues = append(issues, ReferenceIssue{
					Field: field, Index: i, Reason: "duplicate_reference",
					Kind: referenceKind(value), ID: referenceID(value), LegacyName: referenceLegacyName(value),
					Message: "A duplicate reference was collapsed; missing unknown fields were merged into the first entry.",
				})
				continue
			}
			seen[key] = len(out)
		}
		out = append(out, value)
	}
	return out, issues, true
}

func normalizeReferenceEntry(field string, index int, raw interface{}, expected worldstore.RefKind) (interface{}, *ReferenceIssue) {
	switch typed := raw.(type) {
	case string:
		name := strings.TrimSpace(typed)
		if name == "" {
			return raw, &ReferenceIssue{Field: field, Index: index, Reason: "malformed", Message: "An empty name-only reference was preserved."}
		}
		return map[string]interface{}{
			"kind":       string(expected),
			"id":         "",
			"legacyName": name,
		}, &ReferenceIssue{Field: field, Index: index, Reason: "legacy_name", LegacyName: name, Message: "Name-only reference is unresolved until explicitly repaired."}
	case map[string]interface{}:
		value := deepCopyReferenceMap(typed)
		explicitKind := strings.TrimSpace(stringField(value, "kind"))
		if explicitKind != "" {
			if !knownReferenceKind(explicitKind) {
				return value, &ReferenceIssue{Field: field, Index: index, Reason: "unknown_kind", Kind: explicitKind, Message: "The reference kind is not supported; the source value was preserved."}
			}
			if explicitKind != string(expected) {
				return value, &ReferenceIssue{Field: field, Index: index, Reason: "kind_mismatch", Kind: explicitKind, Message: "The reference kind does not match its field; the source value was preserved."}
			}
		}
		value["kind"] = string(expected)
		id := strings.TrimSpace(stringField(value, "id"))
		title := strings.TrimSpace(stringField(value, "title"))
		if title == "" {
			title = strings.TrimSpace(stringField(value, "name"))
		}
		legacy := strings.TrimSpace(stringField(value, "legacyName"))
		if legacy == "" {
			legacy = title
		}
		if id != "" {
			value["id"] = id
			if _, exists := value["title"]; !exists && title != "" {
				value["title"] = title
			}
			return value, nil
		}
		value["id"] = ""
		if legacy == "" {
			return value, &ReferenceIssue{Field: field, Index: index, Reason: "missing_id", Message: "The reference has neither an identity nor a legacy name; the source value was preserved."}
		}
		value["legacyName"] = legacy
		return value, &ReferenceIssue{Field: field, Index: index, Reason: "legacy_name", LegacyName: legacy, Message: "Name-only reference is unresolved until explicitly repaired."}
	default:
		return raw, &ReferenceIssue{Field: field, Index: index, Reason: "malformed", Message: "The reference is not an object or string; the source value was preserved."}
	}
}

func normalizeLorebookFields(document map[string]interface{}, world bool) []ReferenceIssue {
	howl, _ := document["howl"].(map[string]interface{})
	var howlRaw, topRaw interface{}
	var hasHowl, hasTop bool
	if howl != nil {
		howlRaw, hasHowl = howl["lorebooks"]
	}
	topRaw, hasTop = document["lorebooks"]
	if !hasHowl && !hasTop {
		return nil
	}

	var howlValues, topValues []interface{}
	var issues []ReferenceIssue
	if hasHowl {
		values, listIssues, ok := normalizeReferenceList("howl.lorebooks", howlRaw, worldstore.RefLorebook)
		if ok {
			howlValues = values
		}
		issues = append(issues, listIssues...)
	}
	if hasTop {
		values, listIssues, ok := normalizeReferenceList("lorebooks", topRaw, worldstore.RefLorebook)
		if ok {
			topValues = values
		}
		issues = append(issues, listIssues...)
	}

	merged, duplicateIssues := mergeReferenceLists("howl.lorebooks", howlValues, "lorebooks", topValues)
	issues = append(issues, duplicateIssues...)
	if world {
		if howl == nil {
			howl = map[string]interface{}{}
			document["howl"] = howl
		}
		if hasHowl || hasTop || len(merged) > 0 {
			howl["lorebooks"] = merged
		}
		delete(document, "lorebooks")
		return issues
	}
	if hasTop || hasHowl || len(merged) > 0 {
		document["lorebooks"] = merged
	}
	if howl != nil {
		delete(howl, "lorebooks")
	}
	return issues
}

func mergeReferenceLists(firstField string, first []interface{}, secondField string, second []interface{}) ([]interface{}, []ReferenceIssue) {
	out := make([]interface{}, 0, len(first)+len(second))
	seen := map[string]int{}
	var issues []ReferenceIssue
	add := func(field string, values []interface{}) {
		for i, value := range values {
			key := referenceKey(value)
			if key != "" {
				if index, exists := seen[key]; exists {
					if existing, ok := out[index].(map[string]interface{}); ok {
						if incoming, ok := value.(map[string]interface{}); ok {
							mergeReferenceFields(existing, incoming)
						}
					}
					issues = append(issues, ReferenceIssue{
						Field: field, Index: i, Reason: "duplicate_reference",
						Kind: referenceKind(value), ID: referenceID(value), LegacyName: referenceLegacyName(value),
						Message: "A duplicate reference was collapsed while canonicalizing the lorebook list.",
					})
					continue
				}
				seen[key] = len(out)
			}
			out = append(out, value)
		}
	}
	add(firstField, first)
	add(secondField, second)
	return out, issues
}

func referenceKey(value interface{}) string {
	if _, ok := value.(map[string]interface{}); !ok {
		return ""
	}
	kind := referenceKind(value)
	id := referenceID(value)
	if id != "" {
		return kind + "\x00id:" + id
	}
	legacy := referenceLegacyName(value)
	if legacy != "" {
		return kind + "\x00legacy:" + strings.ToLower(legacy)
	}
	return ""
}

func referenceKind(value interface{}) string {
	item, _ := value.(map[string]interface{})
	return stringField(item, "kind")
}

func referenceID(value interface{}) string {
	item, _ := value.(map[string]interface{})
	return strings.TrimSpace(stringField(item, "id"))
}

func referenceLegacyName(value interface{}) string {
	item, _ := value.(map[string]interface{})
	return strings.TrimSpace(stringField(item, "legacyName"))
}

func mergeReferenceFields(existing, incoming map[string]interface{}) {
	for key, value := range incoming {
		if _, exists := existing[key]; !exists {
			existing[key] = value
		}
	}
}

func deepCopyReferenceMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for key, value := range in {
		out[key] = deepCopyReferenceValue(value)
	}
	return out
}

func deepCopyReferenceValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return deepCopyReferenceMap(typed)
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i, item := range typed {
			out[i] = deepCopyReferenceValue(item)
		}
		return out
	default:
		return value
	}
}

func referenceIssuesMessage(issues []ReferenceIssue) string {
	return fmt.Sprintf("%d reference issue(s)", len(issues))
}
