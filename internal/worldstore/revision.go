package worldstore

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// revisionPrefix marks the revision algorithm. A future change to canonicalization
// must change the prefix so old revisions cannot be compared as if equivalent.
const revisionPrefix = "rv1-"

// Revision computes the frozen content revision: SHA-256 over canonical JSON of
// the document, truncated to 16 hex characters. Canonical JSON sorts object keys,
// preserves array order, emits UTF-8 without normalization, distinguishes null
// from omission, and normalizes number representation.
//
// Excluded from the hash: envelope timestamp fields (updatedAt/updated_at) and
// media payload fields inside image entries. Changing a timestamp or an embedded
// image blob therefore does not change content identity.
func Revision(document map[string]interface{}) string {
	sum := sha256.Sum256([]byte(canonicalJSON(document)))
	return revisionPrefix + hex.EncodeToString(sum[:])[:16]
}

func canonicalJSON(value interface{}) string {
	var b strings.Builder
	writeCanonical(&b, value, "")
	return b.String()
}

func writeCanonical(b *strings.Builder, value interface{}, parent string) {
	switch typed := value.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if typed {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case string:
		encoded, _ := json.Marshal(typed)
		b.Write(encoded)
	case json.Number:
		b.WriteString(canonicalNumber(typed))
	case int:
		b.WriteString(strconv.Itoa(typed))
	case int64:
		b.WriteString(strconv.FormatInt(typed, 10))
	case float64:
		b.WriteString(strconv.FormatFloat(typed, 'g', -1, 64))
	case map[string]interface{}:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			if isExcludedKey(parent, key) {
				continue
			}
			keys = append(keys, key)
		}
		sort.Strings(keys)
		b.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			encodedKey, _ := json.Marshal(key)
			b.Write(encodedKey)
			b.WriteByte(':')
			writeCanonical(b, typed[key], key)
		}
		b.WriteByte('}')
	case []interface{}:
		b.WriteByte('[')
		for i, item := range typed {
			if i > 0 {
				b.WriteByte(',')
			}
			writeCanonical(b, item, parent)
		}
		b.WriteByte(']')
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			b.WriteString("null")
			return
		}
		b.Write(encoded)
	}
}

// isExcludedKey keeps timestamps and media payloads out of content identity.
// `parent` is the containing key, so "updated" under metadata is a timestamp
// while a user field named "updated" at the top level is content.
func isExcludedKey(parent, key string) bool {
	if parent == "" && (key == "updatedAt" || key == "updated_at") {
		return true
	}
	if parent == "metadata" && key == "updated" {
		return true
	}
	if parent == "image" || parent == "images" {
		switch key {
		case "data", "medium", "original", "path", "thumbnail":
			return true
		}
	}
	return false
}

// canonicalNumber gives every numeric literal one deterministic representation.
// Integers keep their exact literal, so precision beyond float64 is preserved;
// fractional and exponent forms are normalized through float64.
func canonicalNumber(number json.Number) string {
	raw := string(number)
	if !strings.ContainsAny(raw, ".eE") {
		return raw
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return raw
	}
	return strconv.FormatFloat(value, 'g', -1, 64)
}
