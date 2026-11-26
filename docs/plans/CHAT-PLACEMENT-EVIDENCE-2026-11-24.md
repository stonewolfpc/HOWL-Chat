# Chat placement evidence — 2026-11-24

**Status:** Implemented and exercised in both supported hosts. This is the first
slice of the resource-orchestration pillar, not the finished orchestrator.

**Related:** `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`,
`docs/research/INFERENCE-MEMORY-AND-OFFLOAD-SURVEY-2026-11-24.md` (defect 8),
`docs/plans/IMAGE-HARDWARE-OPTIMIZATION-PLAN-2026-11-24.md`.

## 1. Decision

Howl already calculated a chat plan (context, batch, GPU layers, admission
estimate) and told llama-server to use it. It could not say what llama-server
actually placed. The user saw a load result that was only Howl's arithmetic.

This slice separates the three claims the pillar requires:

| Claim | Owner | Example |
|---|---|---|
| Theoretical requirement | `resource.Authority` estimate | 627 MB VRAM / 693 MB RAM |
| Planned placement | chat admission + spawn settings | 8 GPU layers, context 3096, batch 512 |
| Actual engine placement | llama-server's own load log, parsed after health | 8/15 layers, context 3328, batch 512, 96 MiB VRAM / 96 MiB RAM weights, 245 MiB / 239 MiB including KV and compute |

The actual report never feeds back into admission. A stale observation is not a
second capacity owner; the next load still reads the live resource ledger.

## 2. What changed

- `app-llama.go` requests `--log-verbosity 4` when the probed binary advertises
  `--log-verbosity`. Without level 4 llama-server does not print the layer split
  or buffer sizes. This changes diagnostics, not execution.
- `internal/backend/llama/placement.go` parses the engine's own statements:
  `offloaded N/M layers to GPU`, `DEVICE model buffer size`, `DEVICE KV buffer
  size`, `DEVICE compute buffer size`, `n_ctx_slot`, and `n_batch`.
- `internal/runtime/engine.go` holds the last load's `ChatPlacement` report.
  It is cleared when the old server is stopped, a load fails, or the model is
  unloaded, so a notification cannot quote a dead process.
- `app-chat-placement.go` builds the plan/actual report and exposes
  `GetChatPlacement` for diagnostics.
- `app-model.go` parses the bounded loader capture at the health check and emits
  `model:placement` before `model:loaded`.
- `frontend/dist/js/loading-progress.js` listens for `model:placement` and shows
  the full report as an information notification.
- `internal/backend/loader/diagnostics.go` adds a load-phase guard. Level 4
  logging includes prompt and sampled text after the model is loaded (observed
  in the direct probe); the guard forwards only through `model loaded`, and the
  load capture is reset after the health check.

The parser is deliberately a record of reported facts. Host compute and output
buffers are not counted as model weights; they are reported only in the engine
total. `CUDA_Host` compute is host memory, not VRAM.

## 3. Verification

Both hosts used the isolated root
`C:/Users/StoneWolfPC/howl-chat-placement-iso-20261124`.

### Browser dev host (`-tags uiharness`, HTTP 8799)

Clicked **Load** in the model picker after setting **GPU Offload Layers** to 8
through Model Settings.

- Model: `LFM2.5-230M-QAD-Q4_0.gguf`
- Plan: 8 GPU layers, context 3096, batch 512; admission estimate 627 MB VRAM /
  693 MB RAM.
- Actual: **8/15** GPU layers, context **3328**, batch **512**, weights
  **96 MiB VRAM / 96 MiB RAM**, engine total **245 MiB VRAM / 239 MiB RAM**.
- The visible notification said: “Model loaded as planned … llama-server rounded
  the requested context 3096 up to its token block … reported engine total …”.
- `GetChatPlacement` returned the same numbers; `matches_plan` was true.

### Packaged desktop host (`-tags "desktop,production"`, WebView2)

Drove the same model picker and **Load** button with native mouse input.

- Model: `granite-4.0-350m-Q4_K_M.gguf`
- Plan: 8 GPU layers, context 3096, batch 512; admission estimate 613 MB VRAM /
  959 MB RAM.
- Actual: **8/29** GPU layers, context **3328**, batch **512**, weights
  **117 MiB VRAM / 186 MiB RAM**, engine total **167 MiB VRAM / 266 MiB RAM**.
- The information notification was visible in the packaged window with the same
  plan/actual text. Screenshot: `build/verification/desktop-placement.png`.

### Post-load privacy check

After the load, a direct completion was sent to the running llama-server with a
unique marker (`PLACEMENT_GUARD_SECRET_123` in the browser host,
`PLACEMENT_GUARD_DESKTOP_456` in the desktop host). Both returned HTTP 200, the
placement line remained in the host log once, and the marker was absent from
both harness and desktop logs. This exercises the load-phase guard with real
engine output; it is not a chat-workflow test.

The actual spawned command line was observed with `--log-verbosity 4`, e.g.:

```text
llama-server.exe --host 127.0.0.1 --port 8080
  -m ...\granite-4.0-350m-Q4_K_M.gguf -c 3096 --threads 8 -b 512 -ngl 8
  --fit off --log-verbosity 4 ...
```

Checks run after the final code change:

- `go test ./...` — pass.
- `go build -tags uiharness -o build/bin/howl-ui-harness.exe .` — pass.
- `go build -tags "desktop,production" -o build/bin/howl-chat-s1.exe .` — pass.

## 4. Ownership

| State | Owner |
|---|---|
| Fit and admission estimate | `resource.Authority` (`internal/resource/admission.go`) |
| Spawn plan (context, batch, GPU layers) | `LoadModel`/`spawnSettings` |
| Last plan/actual report | `runtime.Engine` |
| Actual layer/buffer facts | llama-server log; `llama.ParsePlacement` is stateless |
| Notification / explanation | `loading-progress.js` |

No competing admission writer was added. The estimate and the actual report are
different authorities over different facts by design.

## 5. Limitations and rejected claims

- **Load-time evidence only.** There is no continuous VRAM monitoring after
  load. The report describes allocation at the health check, not later growth
  from prompt cache, compute graphs or other processes.
- **Not total card usage.** “Engine total” is the sum of reported model, KV and
  compute buffers. It excludes process RSS, driver overhead, CUDA context and
  other applications. It must not be presented as free/used VRAM.
- **Batch evidence is partial.** `n_batch` is parsed; `n_ubatch` and the actual
  scheduler split are not. The planned batch is reported separately.
- **Context rounding is allowed.** llama.cpp rounds a requested context up to
  its token block. The report treats an upward difference of at most 256 tokens
  as the same plan and names both numbers. A lower actual context remains a
  mismatch.
- **Verbose logs cost space during load.** A direct 32-token probe produced
  19,676 bytes at level 4 versus 1,610 bytes at the default level. The guard
  drops post-load output, so prompt text and per-token logs do not accumulate;
  a 256 KiB capture can still lose placement lines on a very noisy load, in
  which case the report says “unverified” rather than guessing.
- **Prompt privacy depends on the guard.** llama-server logs prompt text at
  level 4 (`<|start_of_role|>user...`). The guard is covered by unit tests for
  split markers and inference content; it forwards no post-load line. The
  capture is also cleared at health as a second boundary.
- **No release verification.** Unload still reports the process stopped; it does
  not prove the driver reclaimed the memory.
- **No iGPU, storage residency, TTS placement or coexistence arbitration in this
  slice.** Those remain design work in the orchestration design; this slice only
  makes chat's existing placement visible and auditable.
- **Admission estimates are not benchmarks.** The observed gap between the
  estimate and engine buffers (e.g. 627 MB estimated vs 245 MB reported total) is
  evidence that the estimate includes reserve and conservatism. It is not yet a
  calibrated predictor and is not tuned by this slice.

## 6. Next slices

1. Surface the same plan/actual contract for image jobs in one shared report
   shape instead of separate image-only fields.
2. Add release evidence (process exit + measured VRAM return) before claiming a
   model or component was reclaimed.
3. Measure which placed layers actually run on GPU during decode, not only how
   many were loaded there. A direct probe showed level-4 logs contain prompt
   text, so preserving privacy while measuring decode needs a filter or a
   separate safe metric.
4. Reconcile the over-conservative admission estimate with observed engine
   buffers before using estimates for tighter placement decisions.
5. Extend the authority to iGPU/CPU/storage candidates only after each backend's
   actual execution path and capability are measured.
