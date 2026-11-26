package sdcpp

import (
	"regexp"
	"strconv"
)

// Placement is where sd-server decided to keep the weights of the model it just
// loaded. On a card that cannot hold the model, this is the difference between
// "loaded" and "loaded and will take four minutes an image", and until now it
// only ever reached the console: the app would report a model as loaded while
// every byte of it sat in system RAM, which is exactly the state a user with
// 8 GB of VRAM needs to see.
type Placement struct {
	Known       bool               `json:"known"`
	VRAMMB      float64            `json:"vramMB"`
	RAMMB       float64            `json:"ramMB"`
	Components  map[string]string  `json:"components"` // component -> RAM | VRAM | N/A
	ComponentMB map[string]float64 `json:"componentMB,omitempty"`
	Backends    string             `json:"backends"` // the auto-fit plan verbatim
	// VAETiled records that sd.cpp's own decode fallback ran. Howl pre-enables
	// tiling when its estimate predicts pressure; this is the observed event when
	// the estimate was wrong or tiling was not pre-enabled.
	VAETiled      bool   `json:"vaeTiled"`
	VAETilingNote string `json:"vaeTilingNote,omitempty"`
	// Process carries Howl's own working-set observations. It is not part of the
	// engine's per-component claim.
	Process map[string]Footprint `json:"process,omitempty"`
}

// Offloaded reports that the model is running with no weights resident in VRAM.
func (p Placement) Offloaded() bool {
	return p.Known && p.RAMMB > 0 && p.VRAMMB == 0
}

var (
	// [I] total params memory size = 6083.61MB (VRAM 0.00MB, RAM 6083.61MB): text_encoders 1137.00MB(RAM), diffusion_model 4786.61MB(RAM), vae 160.00MB(RAM), ...
	rePlacementTotal = regexp.MustCompile(`total params memory size = ([0-9.]+)MB \(VRAM ([0-9.]+)MB, RAM ([0-9.]+)MB\):?(.*)`)
	rePlacementPart  = regexp.MustCompile(`([a-z_]+)\s+([0-9.]+)MB\((RAM|VRAM|N/A)\)`)
	// [I] auto-fit: --backend "diffusion=CUDA0,te=CUDA0,vae=CUDA0" --params-backend "diffusion=cpu,te=cpu,vae=cpu"
	rePlacementPlan = regexp.MustCompile(`auto-fit: --backend "([^"]*)" --params-backend "([^"]*)"`)
	reVAEDecodeOOM  = regexp.MustCompile(`VAE decode (compute failed|ran out of memory)`)
	reVAETileChange = regexp.MustCompile(`Reducing VAE decode tiles from ([0-9]+x[0-9]+) to ([0-9]+x[0-9]+)`)
)

// parsePlacement reads one line of sd-server output. It returns false for every
// line that is not a placement summary, so the caller can feed it the whole log.
func parsePlacement(line string) (Placement, bool) {
	if m := rePlacementPlan.FindStringSubmatch(line); m != nil {
		return Placement{Known: true, Backends: m[1] + " | params: " + m[2]}, true
	}
	m := rePlacementTotal.FindStringSubmatch(line)
	if m == nil {
		return Placement{}, false
	}
	// The format has the total first, then the per-component split; keep the
	// per-component figure when it is there because it is the honest one.
	parts := map[string]string{}
	sizes := map[string]float64{}
	for _, pm := range rePlacementPart.FindAllStringSubmatch(m[4], -1) {
		parts[pm[1]] = pm[3]
		if mb, err := strconv.ParseFloat(pm[2], 64); err == nil {
			sizes[pm[1]] = mb
		}
	}
	ram, _ := strconv.ParseFloat(m[3], 64)
	vram, _ := strconv.ParseFloat(m[2], 64)
	return Placement{Known: true, VRAMMB: vram, RAMMB: ram, Components: parts, ComponentMB: sizes}, true
}

// merge keeps what is already known and overrides it with anything the new line
// states, so an auto-fit plan line and a totals line can arrive in any order.
func (p Placement) merge(next Placement) Placement {
	if !next.Known {
		return p
	}
	if next.Backends != "" {
		p.Backends = next.Backends
	}
	if next.RAMMB > 0 || next.VRAMMB > 0 || len(next.Components) > 0 {
		p.RAMMB, p.VRAMMB, p.Components = next.RAMMB, next.VRAMMB, next.Components
		p.ComponentMB = next.ComponentMB
	}
	if next.VAETiled {
		p.VAETiled = true
		if next.VAETilingNote != "" {
			p.VAETilingNote = next.VAETilingNote
		}
	}
	p.Known = true
	return p
}

// parsePlacementEvent reads sd.cpp's own decode-retry lines. They are not a
// placement summary, but they are an execution fact the UI must not hide: the
// engine already tiled because the decode did not fit.
func parsePlacementEvent(line string) (Placement, bool) {
	if reVAEDecodeOOM.MatchString(line) {
		return Placement{Known: true, VAETiled: true, VAETilingNote: line}, true
	}
	if reVAETileChange.MatchString(line) {
		return Placement{Known: true, VAETiled: true, VAETilingNote: line}, true
	}
	return Placement{}, false
}

// recordPlacement feeds one line of server output into the placement state.
func (e *Engine) recordPlacement(line string) {
	next, ok := parsePlacement(line)
	if !ok {
		next, ok = parsePlacementEvent(line)
	}
	if !ok {
		return
	}
	e.placementMu.Lock()
	e.placement = e.placement.merge(next)
	e.placementMu.Unlock()
}

// GetPlacement returns where the loaded model's weights actually live.
func (e *Engine) GetPlacement() Placement {
	e.placementMu.Lock()
	defer e.placementMu.Unlock()
	p := e.placement
	if p.Components == nil {
		p.Components = map[string]string{}
	}
	if p.ComponentMB == nil {
		p.ComponentMB = map[string]float64{}
	}
	p.Process = e.GetFootprints()
	return p
}

// resetPlacement drops the previous model's numbers; a stale "6083 MB in RAM"
// next to a freshly loaded 2 GB SDXL checkpoint would be a lie.
func (e *Engine) resetPlacement() {
	e.placementMu.Lock()
	e.placement = Placement{}
	e.placementMu.Unlock()
}

// placementSummary is the one sentence the UI shows under a loaded model.
func placementSummary(p Placement) string {
	switch {
	case !p.Known:
		return ""
	case p.Offloaded():
		return "Loaded with weights in system RAM (" + formatMB(p.RAMMB) + " offloaded, 0 MB in VRAM). Expect minutes per image on this setup."
	case p.RAMMB > 0:
		return "Loaded — " + formatMB(p.VRAMMB) + " of weights in VRAM, " + formatMB(p.RAMMB) + " in system RAM."
	default:
		return "Loaded — all " + formatMB(p.VRAMMB) + " of weights in VRAM."
	}
}

func formatMB(mb float64) string {
	if mb >= 1024 {
		return strconv.FormatFloat(mb/1024, 'f', 1, 64) + " GB"
	}
	return strconv.FormatFloat(mb, 'f', 0, 64) + " MB"
}
