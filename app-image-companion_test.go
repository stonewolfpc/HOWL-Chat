package main

import "testing"

// The picker and the validator must agree, or a model the UI reports as fully
// equipped is refused at load time with a message about the file it just chose.
func TestCompanionFamilyUsable(t *testing.T) {
	cases := []struct {
		target, file string
		want         bool
		why          string
	}{
		{"anima", "anima", true, "same architecture"},
		{"anima", "", true, "unclassified file cannot be ruled out"},
		{"anima", "unknown", true, "scanner gave up, not a mismatch"},
		{"anima", "sd1.5", false, "an sd1.5 VAE decodes noise under Anima"},
		{"anima", "sdxl", false, "same reason"},
		{"z_image", "flux", true, "both decode with ae.safetensors"},
		{"flux", "z_image", true, "sharing is symmetric"},
		{"qwen_image", "flux", false, "qwen_image has its own VAE"},
		{"z_image", "qwen_image", false, "and does not borrow it"},
		{"Flux", "FLUX", true, "family names are case-insensitive"},
	}
	for _, c := range cases {
		if got := companionFamilyUsable(c.target, c.file); got != c.want {
			t.Errorf("companionFamilyUsable(%q, %q) = %v, want %v (%s)", c.target, c.file, got, c.want, c.why)
		}
	}
}
