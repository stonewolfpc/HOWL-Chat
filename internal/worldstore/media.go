package worldstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxMediaBytes = 8 << 20

var (
	// ErrMediaMalformed reports a manifest whose JSON or required structure is
	// invalid. It is distinct from a structurally valid but unsupported schema.
	ErrMediaMalformed = errors.New("media manifest is malformed")
	// ErrMediaUnsupported reports a structurally valid manifest whose format or
	// role is not supported by this migration slice.
	ErrMediaUnsupported = errors.New("media manifest is unsupported")
	// ErrMediaSchema is retained for callers that treated all schema failures as
	// one class. New code should distinguish malformed from unsupported.
	ErrMediaSchema = ErrMediaMalformed
	// ErrMediaInline reports binary/base64 payloads in a media manifest.
	ErrMediaInline = errors.New("media manifest contains inline media bytes")
)

var mediaRoles = map[string]bool{
	"thumbnail":  true,
	"gallery":    true,
	"attachment": true,
}

// NormalizeMediaManifest validates a reference-only media manifest and returns
// a deterministic JSON representation. Unknown fields are preserved. The
// canonical shape is:
//
//	{"format":"howl-media-refs-v1","items":[{"mediaId":"...","role":"..."}]}
//
// The legacy reference-only shape {"images":[...]} is accepted for LF-1
// migration and preserved as supplied; no image bytes are copied.
func NormalizeMediaManifest(data []byte) ([]byte, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var doc map[string]interface{}
	if err := dec.Decode(&doc); err != nil || doc == nil {
		return nil, fmt.Errorf("%w: decode object", ErrMediaMalformed)
	}
	if inlineMediaIn(doc) {
		return nil, ErrMediaInline
	}
	items, hasItems := doc["items"].([]interface{})
	if hasItems {
		format, hasFormat := doc["format"].(string)
		if !hasFormat || strings.TrimSpace(format) == "" {
			return nil, fmt.Errorf("%w: canonical shape requires format howl-media-refs-v1", ErrMediaMalformed)
		}
		if format != "howl-media-refs-v1" {
			return nil, fmt.Errorf("%w: format must be howl-media-refs-v1", ErrMediaUnsupported)
		}
		for i, raw := range items {
			item, ok := raw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%w: item %d is not an object", ErrMediaMalformed, i)
			}
			mediaID, _ := item["mediaId"].(string)
			if strings.TrimSpace(mediaID) == "" {
				return nil, fmt.Errorf("%w: item %d has no mediaId", ErrMediaMalformed, i)
			}
			role, _ := item["role"].(string)
			if !mediaRoles[role] {
				return nil, fmt.Errorf("%w: item %d has unsupported role %q", ErrMediaUnsupported, i, role)
			}
			if provider, ok := item["provider"]; ok && provider != nil {
				if _, isString := provider.(string); !isString {
					return nil, fmt.Errorf("%w: item %d has a non-string provider", ErrMediaMalformed, i)
				}
			}
		}
	} else {
		images, ok := doc["images"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: expected format howl-media-refs-v1 with items, or legacy images", ErrMediaMalformed)
		}
		for i, raw := range images {
			item, ok := raw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%w: legacy image %d is not an object", ErrMediaMalformed, i)
			}
			id, _ := item["id"].(string)
			mediaID, _ := item["mediaId"].(string)
			if strings.TrimSpace(id) == "" && strings.TrimSpace(mediaID) == "" {
				return nil, fmt.Errorf("%w: legacy image %d has no id/mediaId", ErrMediaMalformed, i)
			}
		}
	}
	normalized, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("%w: encode: %v", ErrMediaMalformed, err)
	}
	return normalized, nil
}

// LoadMedia reads and normalizes the optional media manifest through the
// canonical item lookup. A missing file is absent, not an error.
func (s *Store) LoadMedia(id string) ([]byte, bool, error) {
	if s == nil || s.reg == nil {
		return nil, false, fmt.Errorf("content library registry is not initialised")
	}
	id = strings.TrimSpace(id)
	if !validID.MatchString(id) {
		return nil, false, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	search, err := s.Lookup(id)
	if err != nil {
		return nil, false, err
	}
	if !search.Found || search.Item == nil {
		if !search.Complete {
			return nil, false, fmt.Errorf("%w: an item root is unavailable", ErrUnavailable)
		}
		return nil, false, ErrNotFound
	}
	if search.AuthorityIndeterminate {
		return nil, false, fmt.Errorf("%w: item authority is indeterminate", ErrUnavailable)
	}
	path := filepath.Join(search.Item.Root, id, "media.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("%w: read media: %v", ErrUnavailable, err)
	}
	if len(data) > maxMediaBytes {
		return nil, false, fmt.Errorf("%w: media.json exceeds %d bytes", ErrInvalidDocument, maxMediaBytes)
	}
	normalized, err := NormalizeMediaManifest(data)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrInvalidDocument, err)
	}
	return normalized, true, nil
}

func inlineMediaIn(value interface{}) bool {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, item := range typed {
			switch key {
			case "data", "base64", "blob", "dataUrl", "data_url", "imageData", "image_data":
				if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
					return true
				}
			case "image", "images", "thumbnail", "avatar", "media", "legacyRef":
				if inlineMediaIn(item) {
					return true
				}
			default:
				if inlineMediaIn(item) {
					return true
				}
			}
		}
	case []interface{}:
		for _, item := range typed {
			if inlineMediaIn(item) {
				return true
			}
		}
	case string:
		trimmed := strings.TrimSpace(typed)
		return strings.HasPrefix(strings.ToLower(trimmed), "data:image/") ||
			(strings.Contains(strings.ToLower(trimmed), ";base64,") && len(trimmed) > 64)
	}
	return false
}
