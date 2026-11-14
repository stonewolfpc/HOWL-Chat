# MEMORY-09 — S1 Model Population and Sequential Execution Readiness (WO-LOCAL-S1-02)

**Status:** **READINESS ASSESSMENT — no server started, no model loaded, no call made, S1 not authorised.**
**Date:** 2026-11-14 · **Order:** `WO-LOCAL-S1-02` (planning, inspection, documentation only).
**Supersedes:** nothing. Amends `MEMORY-09-LOCAL-S1-EXECUTION-PLAN-2026-11-14.md` (row 490) by revision note — see §7.

---

## 1. The execution assumption, corrected

The previous plan left an ambiguity: three slots read as three simultaneously-resident models. That is wrong and is now explicit.

**Three distinct model identities, one resident at a time.** VRAM is 8 GiB; a 9B Q4_K_M model plus KV cache already occupies ~6.6 GiB (L3 measurement). Two resident inference models are not merely unnecessary — they are impossible on this hardware.

Sequential execution is an **operational** adaptation. Population stays at **three distinct identities**, calls stay at **18**, ceiling stays at **30**.

---

## 2. Population inspection (read-only)

Method: project registry (`config/model-registry.json`) plus bounded **GGUF header parsing** within `G:\AI GGUF\GGUF` — magic, version, `general.architecture`, `{arch}.block_count`, `{arch}.context_length`, `general.file_type`. No model was loaded, no server contacted, no other drive scanned. Vision/audio/image artefacts, `mmproj` projectors and TTS vocoders were excluded by name.

**27 chat-model candidates** were found. Family distribution is heavily weighted to the `llama` architecture (many 12B finetunes of overlapping lineage), with singles in `qwen35`, `gemma4`, `chatglm`, `cohere2`, `lfm2moe`, `qwen2vl`, `qwen3vl` and `gpt-oss`.

Two entries deserve a note: `Qwen3.8-27B…IQ2_M` (9.61 GB, 65 layers) cannot fit 8 GiB at any sensible offload, and `GPT-OSS-20B…MXFP4` (11.28 GB) is likewise out. Neither is proposed; a slot that cannot load is not a slot.

---

## 3. Ranked candidate assessment

### 3.1 Rank 1 — `S1-A` · Qwen 3.5 9B (verified identity)

| Field | Value |
|---|---|
| File | `G:\AI GGUF\GGUF\Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` |
| Bytes / SHA-256 | 5,627,044,224 · `2ca636d9e81d3d23ca9b60c234fe185d30ec082eeba69ce770fdb0c76559a4f5` |
| Arch / layers / train ctx | `qwen35` · 32 · 262,144 |
| Quantisation | **Q4_K_M** — GGUF `general.file_type = 15` |
| Memory estimate | ~5.24 GiB weights; ~6.6 GiB total at 20k ctx `q8_0` **measured in L3**, not estimated |
| KoboldCpp compatibility | **Demonstrated** — loaded and answered one call under 1.120 |
| Special behaviour | **Thinking mode observed**: L3 reply opened a `Thinking Process` block and spent the full 48-token budget there. Documented here as an observation of this file, not of the family |

This is the only candidate whose load, identity and generation path are empirically established. It anchors the population.

### 3.2 Rank 2 — `S1-B` · Gemma 4 E4B

| Field | Value |
|---|---|
| File | `G:\AI GGUF\GGUF\Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf` |
| Bytes / SHA-256 | 5,335,285,728 · `b5b1913e44aa64d2c12b6159bb782774dd98362528211ac830860595d293ede0` |
| Arch / layers / train ctx | `gemma4` · 42 · 131,072 |
| Quantisation | **Q4_K_M** — `general.file_type = 15` |
| Memory estimate | ~4.97 GiB weights; **estimate** ~5.5–6.2 GiB resident at 8k ctx `q8_0` |
| KoboldCpp compatibility | **Uncertain** — different architecture family (`gemma4`); no local load attempted. KoboldCpp ships Gemma-family support, but this build/file combination is unverified |
| Special behaviour | **Unknown.** No documented thinking mode found; must not be assumed either way |

Distinct tokenizer and architecture from Qwen — a genuinely different identity, not a re-quantisation. The 42-layer count means layer offload must be chosen per model; that choice is recorded, not hidden.

### 3.3 Rank 3 — `S1-C` · LFM 2.5 8B-A1B (MoE)

| Field | Value |
|---|---|
| File | `G:\AI GGUF\GGUF\LFM2.5-8B-A1B-heretic.Q4_K_M.gguf` |
| Bytes / SHA-256 | 5,155,565,440 · `f820ef3f78704ffcaaeed061f29a28b398c1d348d69aca43a142b63e4d1c5c7e` |
| Arch / layers / train ctx | `lfm2moe` · 24 · 128,000 |
| Quantisation | **Q4_K_M** — `general.file_type = 15` |
| Memory estimate | ~4.80 GiB weights; **estimate** ~5.2–5.8 GiB resident at 8k ctx |
| KoboldCpp compatibility | **Uncertain** — MoE architecture (`lfm2moe`); no local load attempted |
| Special behaviour | **Unknown.** Mixture-of-experts may give a different latency profile; that is a property of the model, not a condition |

Chosen deliberately for **architectural** diversity (MoE, 24 layers) rather than smallness alone. It is the lowest-VRAM candidate of the three, which also makes it the safest first *new* load to shake out compatibility.

### 3.4 Held in reserve (not chosen)

| Model | Bytes | Why not now |
|---|---|---|
| `LongWriter-glm4-9b-abliterated.IQ4_XS` | 4.94 GB | `chatglm` — strong first alternate if `gemma4` or `lfm2moe` fails to load |
| `c4ai-command-r7b-12-2024-abliterated-Q4_K_M` | 4.71 GB | `cohere2`, 8,192 train ctx only |
| `Qwen3VL-8B-Uncensored-HauhauCS-Balanced-Q4_K_M` | 4.68 GB | `qwen3vl`, vision-capable; redundant with `qwen35` in lineage |
| 12B `llama` finetunes (e.g. `Mistral-Helcyon-Mercury-12B`, `Wayfarer-2-12B`) | 6.96 GB | Would require partial CPU offload (~20/40 layers); slower and a different runtime profile. Viable as a **second population**, not as a substitute |

**Shortfall:** none. Three defensible identities exist with verified files and hashes.

---

## 4. Sequential execution design

### 4.1 Per-model cycle (no overlap, ever)

```
for model in (S1-A, S1-B, S1-C):          # fixed order, see §4.4
  1. BASELINE   binding check: nothing listening on 5001; VRAM/RAM/temp recorded
  2. LOAD       one server process, loopback only, per-model ctx + gpu layers
  3. VERIFY     GET /api/v1/model must equal the expected identity
                sha256 of the model file re-checked against the manifest entry
                sampler set captured or explicitly `unavailable`
  4. EXECUTE    6 calls: 2 conditions × 3 reps, in plan order
                one write-once record per call; ledger counts before each send
  5. EVIDENCE   verify each record (hash, parsed==raw), then the batch manifest
  6. UNLOAD     terminate the server process; confirm port dead
  7. RECOVERY   VRAM free must return to within 250 MiB of the step-1 baseline
                temperature recorded; cool-down wait before the next model
```

**Hard rules.** One server process at a time. No model substitution — if a model fails to load, that slot is **reported as failed**, not filled by the reserve candidate without a new decision. No concurrent loads. No retry of a truncated reply. No change to prompts, fixture, state serialization or the frozen sampler set.

### 4.2 Evidence separation by identity

```
eval/runs/S1-<stamp>/
  manifest.json                       # order, identities, hashes, ceiling, monitoring series
  S1-A/  call-01-A-A-r1/ … call-06-A-B-r3/
  S1-B/  call-07-B-A-r1/ … call-12-B-B-r3/
  S1-C/  call-13-C-A-r1/ … call-18-C-B-r3/
```

The directory carries the slot, and the record repeats `model_requested`, the file sha256 and the server-reported identity. A record whose `returned_model` disagrees with its slot directory is inadmissible — the check fails rather than being reconciled by hand.

### 4.3 Per-model configuration must be recorded, not assumed

Context size and GPU layer count differ by model (32 / 42 / 24 layers). Those values are **runtime configuration**, not experimental conditions, but they differ per slot and so must be written into the manifest beside each identity, with the reason ("42 layers, 34 offloaded; 8 GiB VRAM"). A reader must be able to see that a difference between slots came from architecture, not from a changed test.

### 4.4 Order, time and recovery concerns

| Concern | Why it is real | Mitigation (proposed) |
|---|---|---|
| **Order effects** | Fixed order A→B→C means any thermal or driver drift correlates with slot | Record a per-model baseline; after the run, compare baselines. Randomising order is *available* but changes the procedure — decision for the owner |
| **Thermal drift across a batch** | 6 consecutive calls per model heat the card | Cool-down to a recorded threshold (proposal: GPU ≤ 45 °C or 3 min, whichever is longer) between models |
| **VRAM fragmentation / leak** | A server that does not fully release can poison the next model | Recovery check is a **gate**: if VRAM does not return within 250 MiB, stop the batch and report |
| **Time-of-day / machine state** | Background desktop load (browser, updates) varies | Record VRAM baseline before each model; do not treat a busy baseline as an error |
| **Partial batch** | Model 2 fails to load after model 1 ran | Batch is filed as **partial**, 6/18 calls, with the failure recorded. No silent replacement |

---

## 5. Hardware safety review of the abort thresholds

The previous plan applied one temperature number to "the GPU". Reviewed against this machine:

| Component | Telemetry actually available | Threshold (proposed) |
|---|---|---|
| **GTX 1080** | `nvidia-smi` temperature, VRAM, utilisation | ~~warn 78 °C, abort 83 °C~~ — **SUPERSEDED by D3 (§8): warn 74 °C, abort 78 °C.** The original value was wrong in principle: 83 °C *is* the card's thermal limit, so aborting there coincides with throttling instead of preceding it |
| **VRAM** | free/total via `nvidia-smi` | **abort if free < 400 MiB** after load, or if a call OOMs |
| **System RAM** | 64 GB DDR4; free via OS counters | abort if available < 8 GB before a load |
| **CPU (i7-10700K)** | **no reliable per-core temperature without extra tooling** | **recorded as unavailable**; monitor utilisation and OS responsiveness instead. **No CPU temperature threshold is claimed** |
| **Driver / display** | none programmatic | any driver reset, display dropout, or TDR event → stop and report |

Two deliberate changes from the previous plan: the threshold is now attached to a **named component** with a rationale, and the **CPU** is explicitly excluded from temperature-based abort because the telemetry does not exist here. A threshold nobody can measure is worse than an acknowledged gap.

No clock, power limit, fan curve, driver or BIOS setting is touched at any point.

---

## 6. Unresolved decisions and risks

| # | Item | Type |
|---|---|---|
| 1 | Confirm the three identities (§3) and the fixed order A→B→C | decision |
| 2 | Accept `gemma4` and `lfm2moe` as **compatibility-unverified** until each loads, or restrict population to verified families | decision / risk |
| 3 | `max_tokens=256` against Qwen's observed thinking mode — keep (truncation as data) or version an amendment | decision |
| 4 | Per-model context (proposal: 8,192 for all three, ≥ prompt size by orders of magnitude) | decision |
| 5 | Per-model GPU layers — load-time choice, must be recorded; failure to fit is a reported failure, not a tuning loop | risk |
| 6 | Cool-down rule (§4.4) and whether order is randomised | decision |
| 7 | CPU temperature unavailable → no CPU thermal abort exists | acknowledged gap |
| 8 | Two of three identities have never been run through this harness | risk, mitigable only by a separately authorised load test |

**Risk statement.** L3 established one model on this build. It did not establish that `gemma4` or `lfm2moe` load at all, nor that a 42-layer model behaves like a 32-layer one under the same VRAM ceiling. The first new load is the real compatibility test, and a failure there is a **result about the population**, not an obstacle to route around.

---

## 7. Revision record for the S1 plan (row 490)

| Change | From | To |
|---|---|---|
| Slot concurrency | ambiguous ("three slots") | **explicitly sequential: one model resident at a time** |
| Population | "three distinct local models (recommended)" | **three named candidates with verified hashes** (§3) |
| Sampler capture | pre-call requirement | unchanged, plus per-model capture recorded in the manifest |
| Budget policy | keep 256 or version an amendment | unchanged — still an open owner decision |
| Ceiling | 18 planned + 12 retries, hard 30 | unchanged |
| Abort thresholds | single temperature | **per-component, with CPU excluded as unmeasurable** (§5) |
| New | — | per-model load/verify/unload/recovery cycle and evidence layout (§4) |

Nothing in the scientific design changed. Conditions, reps, prompts, state serialization and frozen params are untouched.

---

## 8. Owner decisions recorded (2026-11-14)

**D1 — Population accepted.** `S1-A` Qwen3.5-9B, `S1-B` Gemma-4-E4B, `S1-C` LFM2.5-8B-A1B are the planned population: three architecture families, verified file identities.

- `S1-B` and `S1-C` remain **conditional** until each individual load and identity check succeeds.
- **No backfill.** If either fails to load, the batch stops and the shortfall is reported. The reserve candidate (`LongWriter-glm4-9b`) is not substituted, because a substituted model is a different experiment with the same name.

**D2 — `max_tokens = 256` stays frozen.** Truncation is data. Changing the budget after observing Qwen's thinking block would be a post-observation adjustment, so no amendment is made here; a larger budget belongs in a separately versioned experiment (`S1-v2`), never as a mid-run tweak and never as a re-run of truncated cells.

**D3 — GPU abort threshold corrected; §5 is superseded on this point.**

| Component | Was (§5) | **Now** |
|---|---|---|
| GTX 1080 | warn 78 °C, abort 83 °C | **warn 74 °C, abort 78 °C** |
| VRAM free | abort < 400 MiB | unchanged |
| RAM available | abort < 8 GB | unchanged |
| CPU temperature | unavailable | unchanged — no threshold claimed |

The earlier value was wrong in principle, not just in number: 83 °C *is* the 1080's default thermal limit, so an abort there would coincide with the card's own throttle rather than precede it. **78 °C sits 5 °C below the target**, which is margin rather than description.

**D4 — Advertised context is not usable context.** GGUF `{arch}.context_length` (262k / 131k / 128k for A/B/C) states what the architecture advertises, not what this machine can run. The per-model **configured** context (proposal: 8,192) is a runtime configuration recorded in the manifest; the S1 prompt is a few dozen tokens. No claim about long-context behaviour is made anywhere in S1.

**D5 — Track separation affirmed.** This remains readiness work. It authorises neither S1 execution nor the separate `NIP-01` natural-interaction pilot. The two answer different questions — S1 is the controlled assisted/unassisted contrast under frozen conditions; NIP-01 is natural interaction — and results from one must never be presented as evidence for the other.

## 9. What this document does not do

It does not start KoboldCpp, load a model, send a prompt, run a budget probe, execute S1, or alter corpus items, standards or Howl Chat. **Execution requires separate owner authorisation**, and two of the three identities remain unproven on this hardware.

---

## 10. Frozen readiness baseline (2026-11-14)

### 10.1 Slot qualification is not observation

The next compatibility checks are **slot qualification**, not preliminary S1 work.

| | Slot qualification | S1 observation |
|---|---|---|
| Question | Does this candidate load, identify correctly, and release resources under the local path? | Does the assisted condition differ from the unassisted condition? |
| Steps used | §4.1 steps **1–3 and 6–7** (baseline, load, verify, unload, recovery) | §4.1 steps **1–7** |
| Generations | **none** — no prompt is sent | 6 per model |
| Output | one qualification record per slot | 6 evidence records per model |
| Counts against the 18/30 ceiling | **no** | yes |
| May be cited about MEMORY-09 | **never** | as `[OBS]`, single-source |

**Load verification can establish that a candidate works with the local inference path. It cannot establish anything about MEMORY-09 behaviour**, and no statement about the architecture may rest on a qualification record. A generation performed "just to check it works" would be a scaffolded S1 call wearing a qualification label — not permitted.

Qualification records are filed under `eval/runs/SLOT-QUAL-<SLOT>-<stamp>/` and carry: file sha256, server version, `/api/v1/model` identity, configured context, GPU layers, load time, VRAM before/after, recovery delta, and an explicit `generations: 0` field.

### 10.2 No opportunistic tuning between models

Per-model context and GPU-layer count are fixed **at qualification** and carried into S1 **unchanged**. If a model needs different settings mid-batch — because it OOMs, truncates badly, or runs hot — the batch **stops** and the plan is re-versioned. Re-tuning between models would change the runtime conditions that the manifest claims were held constant, and would make two slots incomparable while appearing comparable.

### 10.3 Baseline values, pinned

Machine-readable manifest: **`eval/s1-readiness-baseline.json`** (`S1-readiness-baseline-1`). Full digests, exact byte counts — abbreviations are for prose only and are never the identity.

| Artefact | Bytes | SHA-256 (full) |
|---|---|---|
| S1 plan (row 490) | 12,074 | `d15ebacea5fe3408fc40117ffbf781f82042fcd9a01293abafab270f7fdac9e9` |
| Fixture `F-ST-01-S1` v1 | — (computed) | `3dd41c8ea03c4101295e1bded40824cc3c41de99c154a285ef7fe5af0966412d` |
| Memory state | — (computed) | `aa98d95faf1d18d910c4ba124c0e8d3bf688d3a6d5feee4ca5ff05d30b9d5700` |
| Prompt A | — (computed) | `076de19e9087bb4a3b3d39464516c3417495aa82e27c7efccb971283cd075af7` |
| Prompt B | — (computed) | `97213346a4ac7ec9a0f3f363ee8f9c15fb73cc0667bdfcc80c03bdb8e385b998` |
| `S1-A` Qwen3.5-9B | 5,627,044,224 | `2ca636d9e81d3d23ca9b60c234fe185d30ec082eeba69ce770fdb0c76559a4f5` |
| `S1-B` Gemma-4-E4B | 5,335,285,728 | `b5b1913e44aa64d2c12b6159bb782774dd98362528211ac830860595d293ede0` |
| `S1-C` LFM2.5-8B-A1B | 5,155,565,440 | `f820ef3f78704ffcaaeed061f29a28b398c1d348d69aca43a142b63e4d1c5c7e` |
| L3 preset (reference only) | 3,591 | `9ed540a5da4954f102f92fb70de26b87607eef8068446bbcc786f9f1f3a8db2b` |
| Frozen params | — | `temperature=0.8`, `top_p=0.95`, `max_tokens=256` |
| Offline suite at freeze | — | **88 checks, 88 pass** (87 at freeze; check 88 verifies this manifest and was added before any execution) |

Fixture, state and prompt digests are **computed** from the code rather than transcribed, so a mismatch means the code changed. Model digests are **read from the files** at freeze time. Check 88 re-derives the computed values and the plan hash against the manifest on every suite run, so drift is caught without loading a model.

### 10.4 Change control

Any change to §10.3 — a prompt, a hash, a parameter, a candidate file, a threshold — is a **new version** (`S1-v2`, plan v1.2), never an edit to this baseline. The intended execution path is fixed: **qualify B and C under 10.1, then run only the authorised sequence, in the documented order.**
