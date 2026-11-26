# LFM2.5 CPU Probe — Checkpoint — 2026-11-19

**Identity:** `LFM25-CPU-PROBE-CHECKPOINT-2026-11-19`
**Status:** checkpoint written **before** any model was downloaded-complete, loaded, or run. No generation exists at the time of writing. Nothing committed, staged, or registered.
**Purpose:** (a) preserve the state of the tree and of MEMORY-09 at this point, (b) pre-register the probe design below so its prompt, metrics, and scoring rule cannot be adjusted after seeing output.

---

## 1. Tree state at this checkpoint

| Item | Value |
|---|---|
| HEAD | `830cba4` — `feat(memory-09): add the complete m09 authority engine package` |
| Tracked modifications | 6 (pre-existing; not touched here) — `config/logs/startup.log`, `config/model-registry.json`, `docs/DOCUMENT-REGISTRY.md`, `docs/research/OWNER-ENVIRONMENT-HISTORY-2026-11-14.md`, `eval/harness/__pycache__/validate.cpython-312.pyc`, `eval/harness/validate.py` |
| Staged | 0 |
| Untracked | 351 files at the start of this turn (research records, nothing new from this probe yet) |
| Machine | Intel 16-thread CPU, 68,603,957,248 B RAM (~63.9 GiB), 387 GB free on `C:` |
| Runtime | `runtime/llama/llama-server.exe` — llama.cpp `0.8.0-dev`, build 12015, commit `a34de8618`, Clang 20.1.8 win-x64 |

## 2. MEMORY-09 standing — unchanged by this checkpoint

This checkpoint **does not** move any MEMORY-09 gate, decision, finding, count, or artifact. Carried forward as-is:

- Gate 2 remains **implementation only** (authorised 2026-11-17); architecture review second; testing not at all. `eval/mem1/m09/` is committed at `830cba4` and clean against HEAD; package version `9.3.0-gate2-rev13`.
- S1 behavioural series remains **paused at 5 of 18** calls, **2 content-admissible** (r1, r2 of `S1-A × A`); `call-05` is `T=length`. No new call is authorised by anything here.
- Finding 2 remains `false_closure`, partially resolved, open on operational recovery. Totals unchanged.
- Still open with the owner: F3 vs J15 (correction cannot settle authority), permanence of conflict groups, A1 substitution clearance, and the unpersisted 118-row conformance mapping. K26's disposition (delegation via `Grant.granted_by`, prompt/turn pair recorded in the conversation layer, no new authorization subsystem) stands as previously argued; nothing here implements it.
- B2 disclosure-ingestion decision remains approved 2026-11-17 and recorded in its own stub, not in `OWNER-DECISIONS`.

**Explicitly not done here:** no edit to `eval/`, no change to any frozen artifact, hash, sidecar, receipt, registry, or gate document. The probe below is not a MEMORY-09 experiment execution and cannot supply evidence for or against any MEMORY-09 requirement.

## 3. The idea under test

Whether a **230M–1.2B CPU-only model** is good enough, fast enough, and cheap enough to do *memory-proposal triage* — deciding which candidate memories are worth keeping, which to reject, and which need the user's confirmation — on the same machine that runs the chat, with **no GPU and no VRAM budget**.

If a 230M or 350M model can do this reliably, memory handling could be shaped differently: triage could run inline, always-on, and per-turn, instead of being a separate/remote step. If only the 1.2B is usable, the trade-off is different. If none is usable, that is also data.

This is a **data-gathering probe**, not a hypothesis test with a frozen protocol. It is deliberately small: three models, one prompt each, one generation each.

## 4. Models (vendor-published, verified present on Hugging Face 2026-11-19)

| Model | Params | Release | GGUF repo | File used | Size |
|---|---|---:|---|---|---:|
| LFM2.5-230M | 230M | 2026-06-25 | `LiquidAI/LFM2.5-230M-GGUF` | `LFM2.5-230M-QAD-Q4_0.gguf` | 149.1 MB |
| LFM2.5-350M | 350M | 2026-03-31 | `LiquidAI/LFM2.5-350M-GGUF` | `LFM2.5-350M-QAD-Q4_0.gguf` | 219.3 MB |
| LFM2.5-1.2B-Instruct | 1.2B | 2026-01-05 | `LiquidAI/LFM2.5-1.2B-Instruct-GGUF` | `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf` | 695.8 MB |

All three use the **QAD (quantization-aware distillation) Q4_0** checkpoint, released 2026-08-19 and stated by the vendor to land at roughly 97% of the BF16 average. Same quantization family across all three so the comparison is not a comparison of quantizers.

## 5. Probe design — fixed now, before output

**Runtime:** `runtime/llama/llama-server.exe`, `-ngl 0` (CPU only), context 2,048, one slot, llama.cpp's own chat template from the GGUF, driven over `POST /v1/chat/completions`. Default thread count. Sampling is llama.cpp's **default** chain, sent explicitly so it is recorded: `temperature 0.8`, `top_p 0.95`, `top_k 40`, `repeat_penalty 1.0`. `max_tokens 400` as a stopping cap. One generation per model, no retries, no best-of.

**Prompt (one paragraph, verbatim, identical for all three):**

> You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory. For each item below, decide KEEP (durable, worth remembering), REJECT (transient or worthless), or HOLD (needs the user's confirmation), and give a reason of no more than twelve words. Watch for duplicates, contradictions, and anything marked as uncertain. Respond with exactly one line per item in the form `ID | DECISION | REASON`, no preamble and no extra commentary. Items: M1 — Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user). M2 — The user said "brb" and returned four minutes later. M3 — Sable is a fire mage who cannot swim. M4 — The user first said Sable has never met Corvin, then later said Corvin trained her. M5 — The user prefers short replies and dislikes long scene descriptions. M6 — A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon".

**Expected key (fixed before output):**

| ID | Expected | Why |
|---|---|---|
| M1 | **KEEP** | Durable identity fact, stated as canon |
| M2 | **REJECT** | Transient session chatter, no durable value |
| M3 | **KEEP** | Durable character trait |
| M4 | **HOLD** | Direct contradiction; must not silently pick a side |
| M5 | **KEEP** | Stable user preference, affects future replies |
| M6 | **HOLD** | Explicitly marked uncertain; needs confirmation |

**Scoring rule (fixed before output):** per item, **2** = expected decision *and* a reason that names the right ground (durability, transience, contradiction, uncertainty); **1** = expected decision with absent/vague reason; **0** = wrong decision. Maximum 12. Plus two non-item observations, recorded but not scored into the 12: **format compliance** (exactly six `ID | DECISION | REASON` lines, no preamble) and **invented content** (any fact, name, or detail not present in the prompt).

**Measured per model:** model file size; cold load time to `/health`; prompt tokens and prefill tok/s; generated tokens and decode tok/s; wall time; peak RSS of the server process tree sampled at 100 ms; system RAM before/after.

**Reporting rule:** the rubric, prompt, and parameters are the ones in this document. If a model's output is ambiguous, it is recorded as ambiguous, not rounded to the nearest bucket. No threshold or expectation is adjusted after seeing a reply, and no second generation replaces an inconvenient one.

**Addendum A1 — written after draw 1, before draws 2–3.** Draw 1 produced a complete six-line answer from the 1.2B but truncated/degenerate output from 230M (14 tokens: `M1 | KEEP | REJECT | HOLD |`) and 350M (7 tokens: `M6 | HOLD`). At temperature 0.8, one draw cannot separate a model's capability from a sampling draw, and that ambiguity would be the whole result. Therefore **two further draws per model** are added, same prompt, same sampling, same runtime, no changes to the rubric or key. Draw 1 is preserved unaltered under `raw/<model>/`; draws 2–3 are stored under `raw/<model>-d2/` and `raw/<model>-d3/`, and **all draws are reported**, scored individually. This raises n to 3; it does not replace, hide, or re-roll any inconvenient draw.

## 6. Boundaries

1. **No MEMORY-09 authority.** Results are descriptive observations about three small models. They are not evidence for any conformance requirement and cannot authorise testing, integration, or a gate.
2. **No production wiring.** No code in `eval/` or the Go app is touched; the harness lives under `experiments/lfm25-cpu/` and is standalone.
3. **No silent state change.** Downloads and outputs stay in `experiments/lfm25-cpu/`; model files are not added to the app's model directory, registry, or manifests.
4. **Nothing is committed.** As with every checkpoint in this series, committing remains the owner's call.

## 7. Artifacts this checkpoint anticipates

- `experiments/lfm25-cpu/run_probe.py` — harness
- `experiments/lfm25-cpu/results/results.json` — machine-readable results
- `experiments/lfm25-cpu/raw/<model>/` — server log and raw response per model
- `docs/research/LFM25-CPU-PROBE-RESULTS-2026-11-19.md` — the written comparison (created after the runs, following §5)
