package sdcpp

import "strings"

// AutoencoderGroup names the swapped-in autoencoder a family decodes latents
// with. Two files share a VAE only when they name the same group.
//
// This exists because "same family or unknown" is too strict and "any VAE" is
// too loose: Z-Image reuses the Flux autoencoder (its vae/config.json is a copy
// of flux-dev's), and Anima reuses the Qwen-Image one, while putting a Flux AE
// under Anima decodes every latent as noise. The distinction is one fact about
// the models, so it lives in one function instead of in each caller's guess.
func AutoencoderGroup(family string) string {
	switch strings.ToLower(strings.TrimSpace(family)) {
	case "flux", "z_image":
		return "flux_ae"
	case "anima", "qwen_image":
		return "qwen_image_ae"
	}
	// flux2 has its own flux2-vae, and sd1.5/sd2.1/sdxl/sd3 each have their own,
	// so those families only ever accept their own family.
	return ""
}

// SharesAutoencoder reports whether two families can use the same VAE file.
func SharesAutoencoder(a, b string) bool {
	ga, gb := AutoencoderGroup(a), AutoencoderGroup(b)
	return ga != "" && ga == gb
}
