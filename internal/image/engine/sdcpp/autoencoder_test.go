package sdcpp

import "testing"

func TestAutoencoderGroups(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"z_image", "flux", true}, // Z-Image's VAE is the Flux one
		{"flux", "z_image", true},
		{"anima", "qwen_image", true}, // Anima ships the Qwen-Image VAE
		{"qwen_image", "anima", true},
		{"anima", "flux", false}, // the mismatch that shipped broken models
		{"anima", "sd1.5", false},
		{"flux", "flux2", false}, // flux2 has its own VAE
		{"sdxl", "sd1.5", false},
		{"", "flux", false},
		{"unknown", "unknown", false},
	}
	for _, c := range cases {
		if got := SharesAutoencoder(c.a, c.b); got != c.want {
			t.Errorf("SharesAutoencoder(%q,%q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
	if got := AutoencoderGroup(" Z_Image "); got != "flux_ae" {
		t.Errorf("group should ignore case and spacing, got %q", got)
	}
}
