# Sampler Verification — c4ai-command-r7b, 2048 context

Model: `c4ai-command-r7b-12-2024-abliterated-Q4_K_M.gguf`
Server: shipped `runtime/llama/llama-server.exe`, `-c 2048 -ngl 0 --port 8899`
Method: requests sent with the exact JSON shape `buildRequest()` produces; identical prompt each time.

## Runtime results (server-confirmed)

| Case | Hash | Meaning |
|------|------|---------|
| temperature 0.0, seed 42 | `0fb90d21b0` | baseline |
| same, repeated | `0fb90d21b0` | **deterministic** — temperature honored |
| temperature 1.3, seed 42 | `325e69deb6` | different text |
| same, repeated | `325e69deb6` | reproducible per seed — seed honored |
| temperature 1.3, top_k 1 | `0fb90d21b0` | collapses to the temp-0 result — **top_k honored** |
| presence_penalty 1.5 | `4b1031c2aa` | changes output — penalty honored |
| repeat_penalty 2.0 vs 1.0 | `5abbb53549` / `916f028060` | differ — repeat penalty honored |
| mirostat 2 (tau 5, eta 0.1) | `17c37c978e` | accepted, output changes |
| dry_multiplier 0.8 | `23ad9d0a07` | accepted, output changes |
| dynatemp_range 1.0 | `3d64b37c96` | accepted, output changes |
| min_p 0.5 | `685979ece2` | accepted |
| typical_p 0.9 | `3fa911c3de` | accepted |
| `top_a: 2.0` | `a6fa4d554a` | **ignored** |
| `epsilon_cutoff: 0.9` | `a6fa4d554a` | **ignored** |
| `smoothing_factor: 0.5` | `a6fa4d554a` | **ignored** |
| plain baseline temp 0.8 seed 11 | `a6fa4d554a` | proves the three above are byte-identical to no-op |

`/props` reported `n_ctx: 2048` for the loaded model, matching the `-c` flag.

## What this build of llama.cpp supports

From `llama-server.exe --help`: temperature, top-k, top-p, min-p, top-n-sigma, xtc-probability/threshold,
typical-p, repeat-last-n, repeat-penalty, presence-penalty, frequency-penalty,
dry-multiplier/base/allowed-length/penalty-last-n, dynatemp-range/exp, mirostat/lr/ent, samplers, seed.

**Not supported at all:** top-a, epsilon-cutoff, eta-cutoff, smoothing-factor, smoothing-curve,
encoder-repeat-penalty, no-repeat-ngram.

## UI changes made because of that

- Removed the duplicate **Temperature** card that sat in Advanced Controls (`temperature-advanced`): it was a
  second control for the same value, wired to nothing, so moving it changed the display but never the request.
- Removed the seven cards for unsupported samplers, plus their entries in `samplerSettings`, the preset keys,
  the conflict map, the id→backend mapping, the Go request fields, and the six tutorial entries that taught them.

## Remaining card inventory (18, all backed by a real server parameter)

Temperature, Top-P, Top-K, Min-P, Repeat Penalty, Repeat Last-N, Presence Penalty, Frequency Penalty,
Typical-P, Mirostat Mode, Mirostat Tau, Mirostat Eta, Dynamic Temperature Range, Dynamic Temperature Exponent,
DRY Multiplier, DRY Allowed Length, DRY Base, Seed.

## UI-level checks (served dist, real clicks)

- 18 cards, **no duplicate labels**.
- Scanned flyout text for mojibake: none (`Ã`, `â`, `Â` all absent after the encoding repair of 76 sequences).
- Each slider dispatched a real `input` event: value mirrored into its number box, and `getSamplerSettings()`
  returned a key per control (24 keys including `_enabled` flags).
- Flyout stays inside the viewport at 1280×900.

## Not covered

The Wails bridge itself (native window) — these tests drove the same code and the same HTTP contract, but not the
Electron-WebView2 bindings inside the packaged window.
