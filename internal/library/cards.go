package library

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Item is a discovered character card. It is a view over the original file;
// discovery never rewrites, moves, or imports it.
type Item struct {
	ID          string   `json:"id"`
	LibraryID   string   `json:"library_id"`
	Category    Category `json:"category"`
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Format      Format   `json:"format"`
	Source      Origin   `json:"source"`
	External    bool     `json:"external"`
	Editable    bool     `json:"editable"`
	Available   bool     `json:"available"`
	Description string   `json:"description,omitempty"`
	Personality string   `json:"personality,omitempty"`
	Scenario    string   `json:"scenario,omitempty"`
	FirstMes    string   `json:"first_mes,omitempty"`
	MesExample  string   `json:"mes_example,omitempty"`
}

// ScanCharacters discovers Howl directory cards and SillyTavern JSON/PNG cards
// in one registered root. It scans one level, matching how both formats are
// normally laid out, and never opens more than a bounded amount of card data.
func ScanCharacters(reg Registration) ([]Item, error) {
	if reg.Category != Characters {
		return nil, fmt.Errorf("library %q is not a character library", reg.ID)
	}
	entries, err := os.ReadDir(reg.Path)
	if err != nil {
		return nil, fmt.Errorf("read character library %q: %w", reg.Path, err)
	}
	out := make([]Item, 0)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		path := filepath.Join(reg.Path, entry.Name())
		if entry.IsDir() {
			card := filepath.Join(path, "character.json")
			if item, ok := readCardFile(reg, card, FormatHowl); ok {
				out = append(out, item)
			}
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		switch ext {
		case ".json":
			if item, ok := readCardFile(reg, path, FormatSillyTavern); ok {
				out = append(out, item)
			}
		case ".png":
			if item, ok := readPNGCard(reg, path); ok {
				out = append(out, item)
			}
		}
	}
	return out, nil
}

const maxCardBytes = 16 << 20

func readCardFile(reg Registration, path string, format Format) (Item, bool) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() || info.Size() > maxCardBytes {
		return Item{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Item{}, false
	}
	card, name, ok := characterCard(data)
	if !ok {
		return Item{}, false
	}
	return itemFor(reg, path, format, card, name), true
}

func readPNGCard(reg Registration, path string) (Item, bool) {
	data, err := os.ReadFile(path)
	if err != nil || len(data) > maxCardBytes {
		return Item{}, false
	}
	payload, ok := pngTextChunk(data, "chara")
	if !ok {
		return Item{}, false
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(payload)))
	if err != nil {
		return Item{}, false
	}
	card, name, ok := characterCard(decoded)
	if !ok {
		return Item{}, false
	}
	return itemFor(reg, path, FormatSillyTavern, card, name), true
}

func characterCard(data []byte) (map[string]interface{}, string, bool) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil || raw == nil {
		return nil, "", false
	}
	card := raw
	if nested, ok := raw["data"].(map[string]interface{}); ok {
		card = nested
		if name := firstString(nested, "name", "char_name"); name != "" {
			return nested, name, true
		}
	}
	if name := firstString(raw, "name", "char_name", "title"); name != "" {
		return card, name, true
	}
	return nil, "", false
}

func firstString(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := m[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func itemFor(reg Registration, path string, format Format, card map[string]interface{}, name string) Item {
	id := sha256.Sum256([]byte(string(reg.Category) + "\x00" + strings.ToLower(filepath.Clean(path))))
	return Item{
		ID:          hex.EncodeToString(id[:12]),
		LibraryID:   reg.ID,
		Category:    reg.Category,
		Name:        name,
		Path:        path,
		Format:      format,
		Source:      reg.Origin,
		External:    reg.Origin != OriginHowl,
		Editable:    reg.Role.CanWrite() && format != FormatUnknown,
		Available:   true,
		Description: firstString(card, "description", "personality"),
		Personality: firstString(card, "personality"),
		Scenario:    firstString(card, "scenario"),
		FirstMes:    firstString(card, "first_mes", "first_message"),
		MesExample:  firstString(card, "mes_example", "example_dialogue"),
	}
}

// pngTextChunk extracts a tEXt value. Character cards store their JSON in the
// "chara" keyword, base64-encoded. Only the chunk table is read; image bytes are
// never decoded or copied.
func pngTextChunk(data []byte, keyword string) ([]byte, bool) {
	const signature = "\x89PNG\r\n\x1a\n"
	if len(data) < len(signature)+12 || !bytes.Equal(data[:8], []byte(signature)) {
		return nil, false
	}
	for offset := 8; offset+12 <= len(data); {
		length := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		kind := string(data[offset+4 : offset+8])
		if length < 0 || offset+12+length > len(data) {
			return nil, false
		}
		body := data[offset+8 : offset+8+length]
		if kind == "tEXt" {
			split := bytes.IndexByte(body, 0)
			if split >= 0 && string(body[:split]) == keyword {
				return body[split+1:], true
			}
		}
		offset += 12 + length
	}
	return nil, false
}
