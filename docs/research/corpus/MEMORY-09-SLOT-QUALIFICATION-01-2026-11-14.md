# MEMORY-09 — Slot Qualification `WO-SLOT-QUAL-01` (S1-B, S1-C)

**Status:** both slots **qualified**. Two zero-generation load cycles completed, evidence filed, frozen baseline unmodified.
**Date:** 2026-11-14 · **Authority:** `WO-SLOT-QUAL-01` (owner, 2026-11-14).
**Scope:** compatibility of two model files with KoboldCpp 1.120 on this machine. **Nothing about MEMORY-09 behaviour.** No prompt was sent.

---

## 1. Result

| Slot | Model | Load | Identity | Generations | Recovery | Verdict |
|---|---|---|---|---|---|---|
| **S1-B** | Gemma-4-E4B Q4_K_M | 16 s | matched | **0** | Δ 0 MiB | **PASS** |
| **S1-C** | LFM2.5-8B-A1B Q4_K_M | 32 s | matched | **0** | Δ 0 MiB | **PASS** |

With `S1-A` already demonstrated at L3, **all three population slots now have a demonstrated local load path**. That is a compatibility statement, not an experimental one.

---

## 2. Per-slot evidence

### 2.1 `S1-B` — Gemma-4-E4B

| Field | Value |
|---|---|
| File | `G:\AI GGUF\GGUF\Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf` |
| Bytes / SHA-256 | 5,335,285,728 · `b5b1913e44aa64d2c12b6159bb782774dd98362528211ac830860595d293ede0` — **matched baseline before load** |
| Server | KoboldCpp 1.120, `127.0.0.1:5001` |
| Reported identity | `koboldcpp/Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M1` |
| Config | context 8192, `gpulayers=36`, `q8_0` KV, 8 threads |
| Offload (server log) | **36/43 layers to GPU** |
| GPU buffer | 2,424.87 MiB |
| CPU-mapped buffer | **3,173.18 MiB** |
| `n_layer` | 42 |
| VRAM | 1,245 → 4,423 MiB (Δ **+3,178**) |
| GPU temp | 28 → 30 °C |
| Recovery | port `000` closed, VRAM back to **1,245 MiB** (Δ 0) |
| Generations | **0** — three identity GETs only |

**Observation with consequences:** more than half this model's weights stayed **CPU-mapped** (3,173 MiB) despite 36 of 43 layers nominally on GPU. At S1 this will likely mean slow generation compared with `S1-A`. It is a property of the file and the 8 GB ceiling, not a configuration error — but it must be recorded before latency is compared across slots.

### 2.2 `S1-C` — LFM2.5-8B-A1B

| Field | Value |
|---|---|
| File | `G:\AI GGUF\GGUF\LFM2.5-8B-A1B-heretic.Q4_K_M.gguf` |
| Bytes / SHA-256 | 5,155,565,440 · `f820ef3f78704ffcaaeed061f29a28b398c1d348d69aca43a142b63e4d1c5c7e` — **matched baseline before load** |
| Server | KoboldCpp 1.120, `127.0.0.1:5001` |
| Reported identity | `koboldcpp/LFM2.5-8B-A1B-heretic.Q4_K_M` |
| Config | context 8192, `gpulayers=20`, `q8_0` KV, 8 threads |
| Offload (server log) | **20/25 layers to GPU** |
| GPU buffer | 4,215.99 MiB |
| CPU-mapped buffer | 897.98 MiB |
| `n_layer` | 24 |
| VRAM | 1,245 → 5,896 MiB (Δ **+4,651**) |
| GPU temp | 29 → 31 °C |
| Headroom after load | **2,164 MiB free** |
| Recovery | port `000` closed, VRAM back to **1,245 MiB** (Δ 0) |
| Generations | **0** — three identity GETs only |

`S1-C` loads more fully onto the GPU than `S1-B`, and leaves less headroom. Both are within the 400 MiB abort floor by a wide margin at this context.

---

## 3. Two details worth keeping

**Identity string sanitisation.** KoboldCpp reports `…Q4_K_M1` for the file `…Q4_K_M(1).gguf` — it strips punctuation from the name. A raw string comparison would have reported a **false mismatch on a correct load**, and the tempting fix (renaming the file) would have broken the baseline hash. The match rule is now normalised (casefold, alphanumerics only) and recorded in each record as `match_rule`. The frozen baseline was not touched.

**Layer counts are per-model and frozen at qualification.** `S1-A` 30/32 (L3), `S1-B` 36/43, `S1-C` 20/25. These differ because the models differ. Per §10.2 they carry into S1 unchanged; a mid-batch change would make slots incomparable while the manifest still claimed constant conditions.

---

## 4. Verification and integrity

| Check | Result |
|---|---|
| File digest before load (both) | **matched baseline** |
| Server-reported identity (both) | **matched** |
| `generations` in both records | **0** |
| Requests issued | the three identity GETs only |
| Recovery gate (both) | **pass**, Δ 0 MiB, port closed |
| Frozen baseline modified | **no** — `eval/s1-readiness-baseline.json` still `ad42f28d…` |
| Offline suite | **89/89** |
| Registry | PASS, 0 anomalies |

Evidence:
- `eval/runs/SLOT-QUAL-S1-B-2026-11-14T135031Z/` (`record.json`, `record.sha256`, `recovery.json`)
- `eval/runs/SLOT-QUAL-S1-C-2026-11-14T135147Z/` (same)
- Server logs: `eval/runs/kobold-S1B-qual.log`, `eval/runs/kobold-S1C-qual.log`

Records are write-once; recovery is a **separate** file so the record is never rewritten after the fact.

---

## 5. What this does and does not establish

**Establishes:** both files load under KoboldCpp 1.120 on this machine, report identities matching the frozen baseline, and release VRAM completely on unload. The population is viable as three sequential slots.

**Does not establish:** anything about MEMORY-09. No prompt was sent, no reply exists, no behaviour was observed. A model that loads correctly can still fail the experiment, and a model that failed to load would not have been replaced by the reserve candidate.

**Does not establish:** that six consecutive calls per model are safe, that latency is comparable across slots (§2.1), or that `max_tokens=256` behaves the same on all three. Those are S1 observations.

---

## 6. Next

Qualification does **not** authorise the 18-call experiment. The next step is a separate owner decision on S1 execution, made with these results in hand. `NIP-01` remains a separate design track and is not affected by this document.
