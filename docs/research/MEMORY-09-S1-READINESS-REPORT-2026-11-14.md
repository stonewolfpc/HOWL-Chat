# MEMORY-09 — S1 Readiness Report (resume checkpoint)

**Identity:** `MEMORY-09-S1-READINESS-RPT-1` · **Date:** 2026-11-14 · **Mode:** assessment only
**Server:** not started. **Generations:** none. **S1:** not authorised, not started.
**Filed outside `corpus/`** so that no corpus item, standard, tag or score is touched.

---

## 1. Latest verified checkpoint

| Artefact | Value |
|---|---|
| Readiness record | `MEMORY-09-LOCAL-S1-SEQUENTIAL-READINESS-2026-11-14.md` (S1-v1.1 design, §10 frozen baseline) |
| Machine manifest, **guarded by the suite** | `eval/s1-readiness-baseline.json` (id `S1-readiness-baseline-1`, 2,903 B) — `ad42f28d9402e34624e0b3bf4bd7fed169099a188174b9c12accf5b99f7b05fc` |
| Machine manifest, **governs execution** | `eval/s1-readiness-baseline-v1.1.json` (id `S1-readiness-baseline-1.1`, 3,679 B) — `14a249bce1673cf7186b963566fb71f2b4314e75e1abb31253a1f3a136129981` |
| Plan-of-record | S1 plan row 490, `d15ebacea5fe3408fc40117ffbf781f82042fcd9a01293abafab270f7fdac9e9` |
| Offline suite | **88/88** (`python -m eval.harness.validate` reports 89/89 including the manifest self-check) |
| Corpus check | **PASS** |
| Last execution | L3-01, `eval/runs/L3-2026-11-14T080303Z/` — ceiling 1, calls made 1, server stopped |

## 2. Completed tests and results

| Test | Result | Evidence |
|---|---|---|
| **L3-01** controlled generation | HTTP 200, 2,875 ms, `finish_reason=length`, reply consumed 48 tokens inside a `Thinking Process` block | L3 record; 87/87 evidence verification, 10/10 verifier |
| **Slot qualification `S1-A`** | **QUALIFIED** — 20/33 layers, 20,480 ctx, `q8_0` KV, 5,379 MiB used (65.7 %), recovery Δ −55 MiB, `generations: 0` | `eval/runs/SLOT-QUAL-S1-A-*` |
| **Slot qualification `S1-B`** | was **BLOCKED** by the digest gate — **now QUALIFIED** at S1-v1.2, 2026-11-15T00:42:49Z | `eval/runs/SLOT-QUAL-S1-B-2026-11-15T004249Z/` |
| **Slot qualification `S1-C`** | **NOT ATTEMPTED at 20,480** — and "not attempted" was wrong as originally written: `S1-C` did load and was recorded at 20/25 layers, 8,192 ctx | `eval/runs/SLOT-QUAL-S1-C-2026-11-14T135147Z/` |
| **S1 behavioural series** | **not started** | — |

Budget spent: **1** generation (L3-01). Slot qualification is **outside** the ceiling by design.

## 3. The three identities — and what the digests are now known to be

**Read this before trusting any hash column below.** An external check run today resolved part of the digest question (report: `docs/research/WO-STORAGE-EXT-01-EXTERNAL-HASH-RESEARCH-2026-11-14.md`).

| Slot | File | Bytes | Manifest hash | External status |
|---|---|---|---|---|
| **S1-A** | `Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` | 5,627,044,224 | `2ca636d9…a4f5` | **matches the repository-hosted LFS pointer exactly** — repository-verified reference (`2ca636d9e81d3d23…a4f5`, size 5,627,044,224) |
| **S1-B** | `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf` | 5,335,285,728 | `b5b1913e…ede0` | **DOES NOT MATCH.** The published hash is `d0027dd3…6e96` — and the 19:16:03 **unbuffered-body composite read** returned exactly that. The manifest value is a reading that contradicts the artifact |
| **S1-C** | `LFM2.5-8B-A1B-heretic.Q4_K_M.gguf` | 5,155,565,440 | `f820ef3f…c5c7e` | **No authoritative reference.** The closest published file (`saidutta69/LFM2.5-8B-A1B-heretic-Q4_K_M.gguf`) is 5,155,564,448 B — **992 bytes different**, different filename punctuation. Not the same artifact as far as evidence goes |

Architectures, for the record and unchanged: `qwen35` 32 layers / 262k · `gemma4` 42 layers / 131k · `lfm2moe` 24 layers / 128k.

**S1-A is the only slot whose identity is now externally corroborated.** S1-B's manifest hash is demonstrably not the file's content. S1-C's identity rests on local reads alone.

## 4. Required sampler and configuration

| Control | Frozen value | Status |
|---|---|---|
| `temperature` / `top_p` / `max_tokens` | 0.8 / 0.95 / 256 | frozen, unchanged (owner decision **D2**: keep 256, truncation is data) |
| Surface | `POST /v1/chat/completions` | frozen; native `/api/v1/generate` excluded |
| Seed | none | local adapter reports no usable seed |
| Full sampler capture | required before any S1 call | **still open** — the preset carried no sampler keys and the server defaults were never captured |
| Per-slot runtime profile | recorded in the manifest | S1-A: **20 layers / 20,480 ctx / `q8_0`** |

**The 8,192-ctx proposal is superseded, and this was checked rather than assumed.** `S1-readiness-baseline-1.1` states the runtime profile explicitly: **20 GPU layers, 20,480 context, 90 % VRAM target, 400 MiB hard floor**, for models up to 12B, "retaining the v1 manifest untouched as history". S1-A's qualified profile (20/33 layers, 20,480 ctx, 5,379 MiB = 65.7 %) is exactly that profile. The 8,192 figure in the readiness record's §6 item 4 is a superseded proposal, not an open conflict.

**Residual gap, however:** `validate.py` line 1472 guards `eval/s1-readiness-baseline.json` — the **v1** file — while v1.1 is what governs execution. The suite therefore does not currently protect the manifest that will be used. Reported, not fixed: changing a check is a source change and is unauthorised here.

## 5. Remaining prerequisites

| # | Prerequisite | Blocks |
|---|---|---|
| 1 | **Owner decision on S1-B identity** — see §6 | everything |
| 2 | Sampler capture attempt per model | any S1 call |
| 3 | `gemma4` and `lfm2moe` compatibility unverified on this KoboldCpp build | B and C qualification |
| 4 | S1-C identity (no reference; 992-byte size difference from the only same-named published file) | C admissibility as a third family |
| 5 | VRAM headroom confirmation at load — 8 GiB card, ~1.7 GiB already in use by desktop | B/C load attempt |
| 6 | Suite guards the v1 manifest, not the v1.1 manifest that governs execution | long-term drift protection |

## 6. Proposed next authorised action

**Not S1. Not a generation.** One decision, then at most two load tests.

The decision: what S1-B's expected hash is now that the artifact-level evidence says `d0027dd3…` is the published content and `b5b1913e…` is a reading that contradicts it.

| Option | Effect | Cost |
|---|---|---|
| **B-1 Re-baseline S1-B to the published hash, re-qualify** | honest correction: the gate compares against the artifact's actual content | a baseline version bump (`S1-v1.2`), one re-read to confirm, then qualification |
| **B-2 Keep `b5b1913e`, treat any mismatch as failure** | the slot can never qualify — a gate that refuses a correct read forever | S1 stays blocked indefinitely |
| **B-3 Drop S1-B, run two families** | population 2, calls 12, ceiling re-derived | a **different experiment**; must be declared as `S1-v2`, not a quiet shrink |

Recommendation (not a decision): **B-1**, because the file's bytes are now anchored to an external reference and keeping a contradicted hash is not conservatism — it is a gate calibrated to an error.

After that decision, the only authorised step I would propose is **slot qualification for B and C, `generations: 0`**, one model at a time, with the §4.1 steps 1–3 and 6–7, filed as `SLOT-QUAL-*`.

## 7. Uncertainty and discrepancies

1. **S1-B's manifest hash is wrong or the stored file is.** Leaning towards the hash being wrong, because an unbuffered-body read matched the repository-verified reference — but **the reads have never been shown to be stable**, and one matching read does not make an input trustworthy. Both the integrity stop and this report are consistent with that; neither proves the earlier readings' cause.

**Terminology correction:** earlier drafts of this report and of the environment record said "device-level read". That is withdrawn. `FILE_FLAG_NO_BUFFERING` bypasses the **system file cache**; it does not bypass the filter stack (`WdFilter`, `FileInfo`, `Wof`, `bfs`, `gameflt`, `UCPD` are all attached to `G:`) or the class/port drivers. The accurate description is **unbuffered-body composite read** — and the tool's output is a composite, because the final 4,064 bytes are read buffered.

**Update 2026-11-15 — the probe was run.** `MEMORY-09-S1B-STABILITY-RESULT-2026-11-15.md`: 10 reads (5 buffered, 5 composite), one identity, one digest — `d0027dd3…6e96` on both paths, 0/10 for `b5b1913e`. Stability is now measured rather than assumed. The manifest decision below remains the owner's.

**Update 2026-11-15 — `S1-B` qualified under S1-v1.2.** `SLOT-QUAL-S1-B-2026-11-15T004249Z`. File verified 5,335,285,728 B / `d0027dd3…6e96` against the active baseline; KoboldCpp 1.120 reported `koboldcpp/Gemma-4-E4B-…Q4_K_M1`; **20/43** layers offloaded, CUDA0 buffer 1,570.63 MiB, CPU-mapped 4,027.42 MiB, `n_ctx 20736`, 720 tensors; VRAM 1,750 → 4,041 MiB (49.3 %), 29 °C; **0 generations**, three identity GETs; unload recovery Δ 12 MiB with the port closed. `S1-B` status: **QUALIFIED**. `S1-C` remains unqualified at the 20,480 profile.
2. **S1-C may not be the artifact named in the manifest.** 992-byte size difference from the only same-named published file.
3. **The suite validates the v1 manifest**, while `v1.1` is the execution authority. Not a S1 blocker; a real gap in drift protection.
4. **Sampler defaults remain uncaptured**, so "exact sampling values" is still incomplete for every slot.
5. **All three slots have loaded — this corrects an earlier claim in this report.** Evidence: `S1-C` at `20/25` layers, 8,192 ctx (`SLOT-QUAL-S1-C-2026-11-14T135147Z`); `S1-B` at `36/43`, 8,192 (`SLOT-QUAL-S1-B-2026-11-14T135031Z`) **and again at `20/43`, 20,480 at 14:02Z with no record filed** (`eval/runs/kobold-S1B-v11-qual.log`); `S1-A` at `20/33`, 20,480 (`SLOT-QUAL-S1-A-2026-11-14T140219Z`). What remains true: only A and B have loaded at the 20,480 profile, and `S1-C` has no reference digest.
6. **The integrity stop remains in force.** Nothing here lifts it; only an owner decision does.

## 8. GPU state at the time of this report (operational, not experimental)

`nvidia-smi`: GTX 1080, **8,192 MiB total, 1,730 MiB used (21.1 %), 6,330 MiB free, 11 % utilisation**. **No compute processes present** — no llama, kobold, python or server process. The resident memory is desktop compositing and normal applications (explorer, ShellHost, StartMenu, SearchHost, EdgeWebView, iCUE, two Spotify instances, ChatGPT, WidgetBoard). **Nothing was killed and nothing needed to be**: no weights are stashed, no test process is loose.

At the owner's ~90 % ceiling that is ~7.37 GiB usable, of which ~6.33 GiB is free now. S1-A qualified at 5.38 GiB, inside that.

On "max the layers for e4b": the governing profile is **20 layers / 20,480 ctx for models up to 12B**, which is what S1-A was qualified at. `gemma4` has 42 layers, so 20 offloaded leaves 22 on CPU — it will load, but slower, and the number must be the recorded qualification value rather than a per-run choice. 8 GiB total with ~1.7 GiB already held by the desktop is the binding constraint here, not the model size alone.

---

## Claim ceiling

**Establishes:** the checkpoint, the completed tests and their results, the three identities with external status, the frozen controls and the open ones, the blocking decision, and the current GPU state.

**Does not establish:** that S1-B or S1-C will load; which digests were faulty; that S1-A's qualification profile is comparable to an 8,192-ctx profile; or that S1 may proceed. No server was started, no prompt sent, no file hashed, no corpus item changed.
