package main

import (
	"fmt"

	"howl-chat/internal/backend/types"
)

// spawnSettingsFor returns the settings for one spawn with the chat template that
// belongs to THIS model.
//
// The template must not be written into the saved settings: doing that made the
// first family's template stick forever, because the "is a template already set?"
// check stayed true for every later model — a Qwen template sent to a Gemma model
// is worse than no template at all. Here a template the user typed is respected,
// while a template this app injected for a previous family is replaced.
func (a *App) spawnSettingsFor(profile *types.ModelProfile) map[string]interface{} {
	live := a.settings.Model()
	out := make(map[string]interface{}, len(live)+1)
	for k, v := range live {
		out[k] = v
	}
	if profile == nil || profile.Template == "" {
		fmt.Printf("[LoadModel] chat template: none available for family %q\n", a.familyName())
		return out
	}
	current, _ := out["custom_jinja_template"].(string)
	if current == "" || current == a.engine.AutoTemplate() {
		out["custom_jinja_template"] = profile.Template
		a.engine.SetAutoTemplate(profile.Template)
		fmt.Printf("[LoadModel] chat template: model-native for family %q (%d bytes)\n",
			profile.Family, len(profile.Template))
	} else {
		fmt.Printf("[LoadModel] chat template: user override (%d bytes), family %q\n",
			len(current), profile.Family)
	}
	return out
}
