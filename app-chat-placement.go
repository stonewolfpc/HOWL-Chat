package main

// app-chat-placement.go — the chat load's plan/actual report.
//
// H0c decides a plan before the spawn. This file records what llama-server
// actually did with it after the health check. The report is for the user and
// diagnostics; it never feeds back into admission, because a stale observation
// would become a competing owner of the next decision.

import (
	"fmt"
	"strconv"

	llamabackend "howl-chat/internal/backend/llama"
	"howl-chat/internal/resource"
	"howl-chat/internal/runtime"
)

// chatPlacementReport joins the requested plan with the engine's own placement
// statements. The two are kept separate in the message because "planned" and
// "actual" are different claims.
func chatPlacementReport(modelPath string, settings map[string]interface{}, decision resource.AdmissionDecision, actual llamabackend.Placement) *runtime.ChatPlacement {
	plannedLayers := settingInt(settings, "gpu_layers", 0)
	plannedContext := settingInt(settings, "context_size", 4096)
	plannedBatch := settingInt(settings, "batch_size", 512)

	source := "requested"
	if decision.Adjustment != nil {
		if decision.Manual {
			source = "manual fallback"
		} else {
			source = "fitted"
		}
	}

	report := &runtime.ChatPlacement{
		Model:            modelPath,
		PlanSource:       source,
		PlannedGPULayers: plannedLayers,
		PlannedContext:   plannedContext,
		PlannedBatch:     plannedBatch,
		PlannedVRAMMB:    decision.Estimate.RequiredVRAMMB,
		PlannedRAMMB:     decision.Estimate.RequiredRAMMB,
		Actual:           actual,
	}
	// llama.cpp rounds a requested context up to its own token block. That is a
	// completed plan, not a different one; report the rounding explicitly so the
	// user still sees both numbers. A lower context or a different layer count
	// is a real departure.
	roundedContext := actual.Known && actual.Context > plannedContext && actual.Context-plannedContext <= 256
	contextMatches := actual.Known && (actual.Context == plannedContext || roundedContext)
	layerMatches := actual.Known && actual.GPULayers == plannedLayers
	batchMatches := actual.Batch == 0 || actual.Batch == plannedBatch
	report.MatchesPlan = layerMatches && contextMatches && batchMatches && actual.Complete()

	plan := fmt.Sprintf("planned %d GPU layer(s), context %d, batch %d; admission estimated %d MB VRAM / %d MB RAM",
		plannedLayers, plannedContext, plannedBatch, report.PlannedVRAMMB, report.PlannedRAMMB)

	switch {
	case !actual.Known:
		report.Message = fmt.Sprintf("Model loaded. %s. llama-server did not report its placement, so the actual footprint is unverified.", plan)
	case !actual.Complete():
		report.Message = fmt.Sprintf("Model loaded. %s. llama-server reported only part of its placement (layers %d/%d, %s), so the actual footprint is unverified.",
			plan, actual.GPULayers, actual.TotalLayers, placementActual(actual))
	case report.MatchesPlan && roundedContext:
		report.Message = fmt.Sprintf("Model loaded as planned: %s. Engine placed %d/%d GPU layers at %s (llama-server rounded the requested context %d up to its token block); %s.",
			plan, actual.GPULayers, actual.TotalLayers, actualContext(actual), plannedContext, placementActual(actual))
	case report.MatchesPlan:
		report.Message = fmt.Sprintf("Model loaded as planned: %s. Engine placed %d/%d GPU layers at %s; %s.",
			plan, actual.GPULayers, actual.TotalLayers, actualContext(actual), placementActual(actual))
	default:
		report.Message = fmt.Sprintf("Model loaded with a different placement than planned: %s. Engine placed %d/%d GPU layers at %s; %s.",
			plan, actual.GPULayers, actual.TotalLayers, actualContext(actual), placementActual(actual))
	}
	return report
}

func settingInt(settings map[string]interface{}, key string, def int) int {
	value, err := strconv.Atoi(intSetting(settings, key, def))
	if err != nil {
		return def
	}
	return value
}

func actualContext(p llamabackend.Placement) string {
	if p.Context <= 0 {
		return "an unreported context"
	}
	if p.Batch > 0 {
		return fmt.Sprintf("context %d, batch %d", p.Context, p.Batch)
	}
	return fmt.Sprintf("context %d (batch not reported)", p.Context)
}

func placementActual(p llamabackend.Placement) string {
	if p.GPUMB <= 0 && p.RAMMB <= 0 {
		return "no model buffer sizes reported"
	}
	result := fmt.Sprintf("model buffers %.0f MiB VRAM / %.0f MiB RAM", p.GPUMB, p.RAMMB)
	if p.KVGPUMB+p.KVRAMMB+p.ComputeGPUMB+p.ComputeRAMMB > 0 {
		result += fmt.Sprintf("; reported engine total %.0f MiB VRAM / %.0f MiB RAM including KV and compute",
			p.TotalGPUMB(), p.TotalRAMMB())
	}
	return result
}

// GetChatPlacement returns the last chat load's plan and observed placement.
// It is empty before a load reaches a healthy server.
func (a *App) GetChatPlacement() map[string]interface{} {
	report := a.engine.ChatPlacement()
	if report == nil {
		return map[string]interface{}{
			"known":   false,
			"message": "No chat placement has been reported yet.",
		}
	}
	return map[string]interface{}{
		"known":              true,
		"model":              report.Model,
		"plan_source":        report.PlanSource,
		"planned_gpu_layers": report.PlannedGPULayers,
		"planned_context":    report.PlannedContext,
		"planned_batch":      report.PlannedBatch,
		"planned_vram_mb":    report.PlannedVRAMMB,
		"planned_ram_mb":     report.PlannedRAMMB,
		"matches_plan":       report.MatchesPlan,
		"actual": map[string]interface{}{
			"known":          report.Actual.Known,
			"gpu_layers":     report.Actual.GPULayers,
			"total_layers":   report.Actual.TotalLayers,
			"gpu_mb":         report.Actual.GPUMB,
			"ram_mb":         report.Actual.RAMMB,
			"kv_gpu_mb":      report.Actual.KVGPUMB,
			"kv_ram_mb":      report.Actual.KVRAMMB,
			"compute_gpu_mb": report.Actual.ComputeGPUMB,
			"compute_ram_mb": report.Actual.ComputeRAMMB,
			"total_gpu_mb":   report.Actual.TotalGPUMB(),
			"total_ram_mb":   report.Actual.TotalRAMMB(),
			"context":        report.Actual.Context,
			"batch":          report.Actual.Batch,
			"buffers":        report.Actual.Buffers,
		},
		"message": report.Message,
	}
}
