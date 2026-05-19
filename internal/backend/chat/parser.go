package chat

import (
	"encoding/json"
	"strings"
)

// ParseImageFromValue converts a raw interface{} image value (string path or
// map with path/thumbnail keys) into a typed *Image.  Returns nil for blank
// or unrecognised values.
func ParseImageFromValue(imageVal interface{}) *Image {
	if imageVal == nil {
		return nil
	}
	switch v := imageVal.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return nil
		}
		return &Image{Path: v}
	case map[string]interface{}:
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		var img Image
		if err := json.Unmarshal(b, &img); err != nil {
			return nil
		}
		if img.Path == "" && img.Thumbnail == "" {
			return nil
		}
		return &img
	default:
		return nil
	}
}

// ParseScenarioContextFromMap converts a frontend scenario map into a typed
// ScenarioContext.  The frontend stores the display name as "title"; this
// function maps it to the canonical "name" field when "name" is absent.
func ParseScenarioContextFromMap(m map[string]interface{}) (ScenarioContext, error) {
	clean := make(map[string]interface{}, len(m))
	for k, v := range m {
		clean[k] = v
	}

	if name, has := clean["name"]; !has || name == "" || name == nil {
		if title, has := clean["title"]; has && title != "" {
			clean["name"] = title
		}
	}

	image := ParseImageFromValue(clean["image"])
	delete(clean, "image")

	b, err := json.Marshal(clean)
	if err != nil {
		return ScenarioContext{}, err
	}
	var sc ScenarioContext
	if err := json.Unmarshal(b, &sc); err != nil {
		return ScenarioContext{}, err
	}
	sc.Image = image
	return sc, nil
}

// ParseCharacterContextFromMap converts a frontend character map into a typed
// CharacterContext.  The image field is handled separately to avoid storing
// raw base64 blobs inside the struct.
func ParseCharacterContextFromMap(m map[string]interface{}) (CharacterContext, error) {
	clean := make(map[string]interface{}, len(m))
	for k, v := range m {
		clean[k] = v
	}

	// Hoist frontend metadata{age,role,background,priority} sub-object into top-level
	// fields so they map onto the flat Character struct fields.
	if meta, ok := clean["metadata"].(map[string]interface{}); ok {
		for _, key := range []string{"age", "role", "background", "priority", "version"} {
			if v, exists := meta[key]; exists {
				if _, already := clean[key]; !already {
					clean[key] = v
				}
			}
		}
	}
	delete(clean, "metadata")

	image := ParseImageFromValue(clean["image"])
	delete(clean, "image")

	b, err := json.Marshal(clean)
	if err != nil {
		return CharacterContext{}, err
	}
	var cc CharacterContext
	if err := json.Unmarshal(b, &cc); err != nil {
		return CharacterContext{}, err
	}
	cc.Image = image
	return cc, nil
}

// ParseWorldContextFromMap converts a frontend world map into a typed WorldContext.
// The frontend stores the display name as "title"; this function maps it to the
// canonical "name" field when "name" is absent, matching ParseScenarioContextFromMap.
func ParseWorldContextFromMap(m map[string]interface{}) (WorldContext, error) {
	clean := make(map[string]interface{}, len(m))
	for k, v := range m {
		clean[k] = v
	}

	// Hoist "title" into "name" when "name" is missing or empty (frontend uses title)
	if name, has := clean["name"]; !has || name == "" || name == nil {
		if title, has := clean["title"]; has && title != "" {
			clean["name"] = title
		}
	}

	image := ParseImageFromValue(clean["image"])
	delete(clean, "image")

	b, err := json.Marshal(clean)
	if err != nil {
		return WorldContext{}, err
	}
	var wc WorldContext
	if err := json.Unmarshal(b, &wc); err != nil {
		return WorldContext{}, err
	}
	wc.Image = image
	return wc, nil
}
