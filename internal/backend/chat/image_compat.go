package chat

import "encoding/json"

func parseImageValue(raw interface{}) *Image {
	if raw == nil {
		return nil
	}

	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil
		}
		return &Image{Path: v}
	case map[string]interface{}:
		imageJSON, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		var image Image
		if err := json.Unmarshal(imageJSON, &image); err != nil {
			return nil
		}
		if image.Path == "" && image.Thumbnail == "" {
			return nil
		}
		return &image
	default:
		return nil
	}
}

func decodeWorldContext(data []byte) (*WorldContext, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	image := parseImageValue(raw["image"])
	delete(raw, "image")

	clean, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var world WorldContext
	if err := json.Unmarshal(clean, &world); err != nil {
		return nil, err
	}
	world.Image = image
	return &world, nil
}

func decodeScenarioContext(data []byte) (*ScenarioContext, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	image := parseImageValue(raw["image"])
	delete(raw, "image")

	clean, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var scenario ScenarioContext
	if err := json.Unmarshal(clean, &scenario); err != nil {
		return nil, err
	}
	scenario.Image = image
	return &scenario, nil
}

func decodeCharacterContext(data []byte) (*CharacterContext, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	image := parseImageValue(raw["image"])
	delete(raw, "image")

	clean, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var character CharacterContext
	if err := json.Unmarshal(clean, &character); err != nil {
		return nil, err
	}
	character.Image = image
	return &character, nil
}
