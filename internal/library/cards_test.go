package library

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"
)

func TestScanCharactersFindsHowlAndSillyTavernCards(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "HowlCard"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "HowlCard", "character.json"), []byte(`{"name":"Howl Card"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tavern.json"), []byte(`{"spec":"chara_card_v2","data":{"name":"Tavern Card"}}`), 0644); err != nil {
		t.Fatal(err)
	}
	png := characterPNG(t, `{"name":"PNG Card"}`)
	if err := os.WriteFile(filepath.Join(root, "card.png"), png, 0644); err != nil {
		t.Fatal(err)
	}

	reg := Registration{ID: "lib", Category: Characters, Path: root, Role: RoleRead, Origin: OriginExternal, Format: FormatSillyTavern}
	items, err := ScanCharacters(reg)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("got %d cards, want 3: %+v", len(items), items)
	}
	seen := map[string]bool{}
	for _, item := range items {
		seen[item.Name] = true
		if item.Editable {
			t.Fatalf("read-only card %q is editable", item.Name)
		}
		if !item.External {
			t.Fatalf("external card %q was reported as Howl-owned", item.Name)
		}
	}
	for _, name := range []string{"Howl Card", "Tavern Card", "PNG Card"} {
		if !seen[name] {
			t.Fatalf("missing %q from %+v", name, items)
		}
	}
}

func TestScanCharactersWriteRoleMarksCardsEditable(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "card.json"), []byte(`{"name":"Editable"}`), 0644); err != nil {
		t.Fatal(err)
	}
	items, err := ScanCharacters(Registration{ID: "lib", Category: Characters, Path: root, Role: RoleRead | RoleWrite, Origin: OriginExternal, Format: FormatSillyTavern})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || !items[0].Editable {
		t.Fatalf("write library did not mark card editable: %+v", items)
	}
}

func characterPNG(t *testing.T, cardJSON string) []byte {
	t.Helper()
	var out bytes.Buffer
	out.Write([]byte("\x89PNG\r\n\x1a\n"))
	writeChunk(&out, "IHDR", []byte{0, 0, 0, 1, 0, 0, 0, 1, 8, 6, 0, 0, 0})
	payload := append([]byte("chara\x00"), []byte(base64.StdEncoding.EncodeToString([]byte(cardJSON)))...)
	writeChunk(&out, "tEXt", payload)
	writeChunk(&out, "IEND", nil)
	return out.Bytes()
}

func writeChunk(out *bytes.Buffer, kind string, body []byte) {
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(body)))
	out.Write(length)
	out.WriteString(kind)
	out.Write(body)
	crc := crc32.NewIEEE()
	crc.Write([]byte(kind))
	crc.Write(body)
	sum := make([]byte, 4)
	binary.BigEndian.PutUint32(sum, crc.Sum32())
	out.Write(sum)
}
