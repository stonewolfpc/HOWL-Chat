package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// civitaiClient is a long-lived HTTP client for Civitai API calls.
var civitaiClient = &http.Client{
	Timeout: 15 * time.Second,
}

// LookupByHash queries Civitai's public API to identify a model by its SHA256.
// It returns the model name (e.g. "xinsir-canny") when found, or an empty
// string when the model is not on Civitai or the request fails.
// This is a public, no-auth endpoint suitable for offline model identification.
func LookupByHash(hash string) (string, error) {
	url := "https://civitai.com/api/v1/model-versions/by-hash/" + hash
	resp, err := civitaiClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("civitai request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", nil // not on Civitai — not an error
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("civitai returned HTTP %d", resp.StatusCode)
	}

	var result struct {
		Name      string `json:"name"`
		ModelName string `json:"model"` // nested model object — actually a struct
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("civitai decode: %w", err)
	}

	// The API returns the version name in "name" (e.g. "xinsir-canny").
	if result.Name != "" {
		return result.Name, nil
	}
	return "", nil
}

// HashFile computes the SHA256 of a file. The result is a 64-char lowercase hex string.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// IdentifyFromFile tries to find a human-readable name for a model file that has
// no usable metadata title. It checks the manifest cache first, then queries
// Civitai by SHA256 on a cache miss.
//
// Returns the name if found, or an empty string to fall back to the existing
// UUID/family naming. On success the manifest is updated with the cached name.
func IdentifyFromFile(path string, m *Manifest) (string, error) {
	// Check manifest cache first — never re-query if we already know the name.
	if m != nil && m.CivitaiName != "" {
		return m.CivitaiName, nil
	}

	// Compute SHA256 — this reads the full file. Only reach here when the
	// name is genuinely unknown, so the cost is justified.
	hash, err := HashFile(path)
	if err != nil {
		return "", fmt.Errorf("hash file: %w", err)
	}

	name, err := LookupByHash(hash)
	if err != nil {
		return "", err // network error — caller should fall back
	}
	if name == "" {
		return "", nil // not on Civitai — not an error
	}

	// Cache the result in the manifest so rescan is instant.
	if m != nil {
		m.CivitaiName = name
		m.Hash = hash // also store the hash for reference/debug
	}
	return name, nil
}
