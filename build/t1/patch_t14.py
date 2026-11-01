import io

def rd(p): return io.open(p, encoding="utf-8", newline="").read()
def wr(p, s): io.open(p, "w", encoding="utf-8", newline="").write(s)

# ---------------------------------------------------------------- 1. App field
p = "app.go"; s = rd(p)
if "autoTemplate" not in s:
    anchor = "\tlastImageModel     string"
    assert s.count(anchor) == 1, "app anchor %d" % s.count(anchor)
    s = s.replace(anchor, anchor + "\n\t// autoTemplate is the chat template this app injected into the last\n\t// model load. It is tracked so a stale auto template can be replaced when\n\t// a different family loads, while a user override is left alone.\n\tautoTemplate string")
    wr(p, s); print("app.go: autoTemplate field added")

# ---------------------------------------------------------------- 2. template resolution
p = "app-model.go"; s = rd(p)
old_block = '''			if _, exists := a.modelSettings["custom_jinja_template"]; !exists {
				a.modelSettings["custom_jinja_template"] = profile.Template
				fmt.Println("[LoadModel] using model's native chat template")
			}'''
if old_block in s:
    s = s.replace(old_block, '''			// The template is NOT written into the saved settings here. Doing
			// that made the first family's template stick forever: the second
			// family loaded with if _, exists… still true and was prompted in
			// the wrong format. It is injected per load in spawnSettingsFor().''')
    print("app-model.go: sticky template write removed")

if "spawnSettingsFor" not in s:
    helper = '''
// spawnSettingsFor returns the settings for one spawn with the chat template that
// belongs to THIS model. A template the user typed is respected; a template this
// app injected for a previous family is replaced, because a Qwen template sent to
// a Gemma model is worse than no template at all.
func (a *App) spawnSettingsFor(profile *types.ModelProfile) map[string]interface{} {
	out := make(map[string]interface{}, len(a.modelSettings)+1)
	for k, v := range a.modelSettings {
		out[k] = v
	}
	if profile == nil || profile.Template == "" {
		fmt.Printf("[LoadModel] chat template: none available for family %q\\n", a.familyName())
		return out
	}
	current, _ := out["custom_jinja_template"].(string)
	if current == "" || current == a.autoTemplate {
		out["custom_jinja_template"] = profile.Template
		a.autoTemplate = profile.Template
		fmt.Printf("[LoadModel] chat template: model-native for family %q (%d bytes)\\n",
			profile.Family, len(profile.Template))
	} else {
		fmt.Printf("[LoadModel] chat template: user override (%d bytes), family %q\\n",
			len(current), profile.Family)
	}
	return out
}
'''
    anchor = "// UnloadModel kills the llama-server process, freeing all model memory."
    assert s.count(anchor) == 1, "unload anchor %d" % s.count(anchor)
    s = s.replace(anchor, helper.strip() + "\n\n" + anchor)
    print("app-model.go: spawnSettingsFor added")

if "spawnSettings := a.spawnSettingsFor(" not in s:
    anchor = '\tfmt.Println("[LoadModel] spawning llama-server")'
    assert s.count(anchor) == 1, "spawn anchor %d" % s.count(anchor)
    s = s.replace(anchor, "\tspawnSettings := a.spawnSettingsFor(profile)\n" + anchor)
    n = s.count("spawnLlamaServerWithModel(modelPath, a.modelSettings,")
    assert n == 2, "spawn call count %d" % n
    s = s.replace("spawnLlamaServerWithModel(modelPath, a.modelSettings,",
                  "spawnLlamaServerWithModel(modelPath, spawnSettings,")
    print("app-model.go: spawn calls use per-load settings")
wr(p, s)

# ---------------------------------------------------------------- 3. family base + log
p = "app-prompt.go"; s = rd(p)
if "familySamplingDefaults" not in s:
    anchor = "func (a *App) buildInferenceOptions() *llama.InferenceOptions {\n\topts := llama.NewInferenceOptions()"
    assert s.count(anchor) == 1, "inference anchor %d" % s.count(anchor)
    s = s.replace(anchor, anchor + '''

	// Per-family defaults first: the publishers' own recommendations differ too
	// much for one set to serve them all. Anything the user has set below still
	// wins, because the sampler overrides are applied after this.
	familySource := familySamplingDefaults(opts, a.familyName())''')

    tail = "\topts.StopStrings = stops"
    assert s.count(tail) == 1, "stops tail %d" % s.count(tail)
    s = s.replace(tail, '''\topts.StopStrings = stops
\tfmt.Printf("[inference] family=%q source=%q temp=%.2f top_p=%.2f top_k=%d min_p=%.2f repeat=%.2f n_predict=%d\\n",
\t\ta.familyName(), familySource, opts.Temperature, opts.TopP, opts.TopK, opts.MinP, opts.RepeatPenalty, opts.NPredict)''')
    print("app-prompt.go: family defaults applied + effective values logged")
wr(p, s)

# ---------------------------------------------------------------- 4. UI shows effective values
p = "frontend/dist/js/sampler-settings.js"; s = rd(p)
if "refreshEffectiveSampler" not in s:
    anchor = "// Initialize on DOM load"
    assert s.count(anchor) == 1, "init anchor %d" % s.count(anchor)
    helper = '''// After a model loads, show the values the engine is actually using. The
// family defaults live in the backend, so without this the sliders would keep
// displaying the last saved number while inference ran on something else.
async function refreshEffectiveSampler() {
    const app = window.go?.main?.App;
    if (!app?.GetEffectiveSampler) return;
    try {
        const eff = await app.GetEffectiveSampler();
        if (!eff) return;
        const shown = {};
        ["temperature", "top_p", "top_k", "min_p", "repeat_penalty", "repeat_last_n",
         "presence_penalty", "frequency_penalty"].forEach(k => {
            if (typeof eff[k] === "number") shown[k] = eff[k];
        });
        setSamplerSettings(shown);
        console.log('[sampler] showing effective values for family:', eff.family || 'unknown');
    } catch (err) {
        console.warn('[sampler] could not read effective values:', err);
    }
}

'''
    s = s.replace(anchor, helper + anchor)
    s = s.replace("document.addEventListener('DOMContentLoaded', function() {\n    initSamplerSettings();\n});",
                  "document.addEventListener('DOMContentLoaded', function() {\n    initSamplerSettings();\n    refreshEffectiveSampler();\n});")
    wr(p, s); print("sampler-settings.js: effective values shown")
