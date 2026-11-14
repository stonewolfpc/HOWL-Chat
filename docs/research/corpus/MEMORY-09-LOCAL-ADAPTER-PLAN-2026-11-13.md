# `MEMORY-09` — Local Adapter Feasibility and Implementation Plan

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-LAP-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PLAN — nothing executed. No server contacted, no model loaded, no prompt sent.**
**Authority:** owner direction, 2026-11-13: local-first, KoboldCpp first, adapter replaceable, no cloud spend.
**Context:** hosted HF route unavailable (P1 addendum A2 §9.4, enforced in A3). 62/62 offline checks green.

---

## 1. Why KoboldCpp first

| Option | Why it fits | Cost |
|---|---|---|
| **KoboldCpp** | One local HTTP server, no provider account, no credits, exposes both a native generate endpoint and an OpenAI-compatible chat endpoint, handles GGUF quantisation and CPU/GPU offload | Needs the binary running with a model loaded |
| **Direct `llama.cpp` (Python binding)** | Tighter control: explicit `n_ctx`, GPU layers, sampler values, token-level capture | Harness owns process lifecycle, memory sizing, and recovery — more to get wrong |
| **SillyTavern as front end** | Useful later as a *compatibility target* (the harness should behave like a client that works there) | Adds a layer between harness and model; not a first target |

**Decision:** KoboldCpp is the first execution target. `llama.cpp` direct stays the second, for when control over loading and token accounting matters more than setup simplicity. **The adapter stays replaceable**, so neither choice becomes permanent.

## 2. What already exists (inspected, not assumed)

| Component | State |
|---|---|
| `local_adapter.py` | Loopback-only enforcement (`localhost`, `127.0.0.1`, `::1`); explicit endpoint, no guessed port; `payload()` builds `{model, messages, stream:false, **params}`; `generate()` still raises `NotImplementedError` |
| Route policy | `local:<runtime>` selects it; hosted refused at three layers |
| Grant gate | Local execution is still *execution*: construction requires a P2-class grant |
| Provenance | Fields defined; local values (`model_revision`, `system_fingerprint`) will be `unavailable` unless the server supplies equivalents |

**Known gap:** the local adapter's **response handling is unverified**. Its request shape is tested; nothing has parsed a real (or realistically-shaped) local reply.

## 3. Phases and gates

| Phase | Question | Needs | Produces |
|---|---|---|---|
| **L0 · Inspect** | Does the adapter say what it does? | Nothing (done) | §2 above |
| **L1 · Offline shape verification** | Can we parse both documented reply shapes, including ugly ones? | Synthetic bodies only | Fixtures + parser + checks, **no server** |
| **L2 · Endpoint discovery** | Is a server listening, which one, which version, which model? | Explicit authorisation — loopback HTTP GET only, **no prompt** | A recorded identity block |
| **L3 · One generation** | Does a single prompt round-trip and capture cleanly? | Separate authorisation, call cap 1 | One record, examined by hand |
| **L4 · S1 behavioural run** | Does the memory concept hold? | S1 authorisation + assigned slots | The 18-call contrast |

**Nothing in L2–L4 is authorised today.** Each is a separate decision, and passing one authorises none of the others.

## 4. L1 — what offline verification must cover

Both surfaces are in scope because KoboldCpp exposes both and the choice between them is a design decision, not a detail:

| Surface | Shape to handle | Hazards |
|---|---|---|
| OpenAI-compatible `/v1/chat/completions` | `choices[0].message.content`, `usage`, `model`, `id`, `created` | Absent `usage`; `finish_reason` truncation; provider field missing |
| KoboldCpp native `/api/v1/generate` | `results[0].text`, token counts, stop reason | Field names differ from the chat surface; sampler semantics differ (`top_k`, `min_p`, `rep_pen`) |

**Required behaviours, each with a fixture:**
1. Content present → transcribed with `raw_body` preserved.
2. Truncated by token limit → recorded as `incomplete`, never as a model failure.
3. Empty `results`/`choices` → recorded, not crashed.
4. Unknown extra fields → ignored, still preserved in `raw_body`.
5. Error body (model not loaded, context exceeded) → structured error, not a silent empty string.
6. Both surfaces produce the **same** `GenerateResponse` shape, so downstream code cannot tell which was used except by the recorded route.

**Why this matters:** the hosted mistake was assuming a provider contract. Here the contract is knowable before any generation — the fixtures can be built from documented shapes, and corrected the moment a real body contradicts them at L3.

## 5. L2 — discovery, and why it is not "just a probe"

A discovery step should answer, without generating:

- Is something listening on the expected loopback port?
- Which server and **version** (identity drift is a reproducibility hazard, not a detail)?
- Which model file is loaded, and what **context size** did it actually get?
- Which surfaces are enabled (chat, native, both)?

**Recorded as identity, never inferred.** If the server cannot report its version or loaded model, those fields are `unavailable` — the run is not blocked, but its admissibility is qualified.

**Boundary:** a GET to a health/identity path is not a prompt. It is still network activity on loopback, which is why it needs its own authorisation rather than riding on L1.

## 6. L3 — the first generation

**One prompt. One call. Cap of 1.** The prompt is a fixed constant, not the S1 probe, so the first call tests the plumbing and nothing else.

Examined by hand against §4's expectations: did the body match the fixture shape? If not, **the fixture is corrected and L1 re-run before L3 is attempted again** — the parser follows the provider, never the reverse.

**Failure is a valid outcome.** A server that cannot answer is information; it is not a reason to loosen capture.

## 7. Runtime monitoring — a separate evidence stream

The owner's hardware (i7-10700K, 64 GB, GTX 1080 8 GB, per owner statement, unverified) is worth measuring. That measurement is **engineering evidence about the machine**, not behavioural evidence about memory.

| Observable | Source | Note |
|---|---|---|
| Model load time, RAM/VRAM after load | process + `nvidia-smi` | Distinguishes "won't fit" from "won't run" |
| VRAM headroom during generation | `nvidia-smi` sampling | Offload settings change behaviour, not just speed |
| CPU/GPU temperature, throttling | vendor/monitor tooling | Throttling affects latency, not content |
| Tokens/second, time-to-first-token | server-reported counts | Slow is acceptable; unexplained is not |
| Resource recovery after unload | process exit + memory | A leak here is a Howl risk later |

**Rules:**
- **Observational only.** Read-only sampling; no tuning, no overclocking, no driver changes during a series. The owner just resolved a driver fault — the monitoring must not become the next instability.
- **Never merged with behavioural results.** The no-bridge rule already in force applies: a performance number may not explain a behavioural outcome, or vice versa.
- **Recorded per run**, with the sampler and offload settings that produced it, because latency without configuration is not evidence.

## 8. Three separations, kept apart

| Stream | Question | Status |
|---|---|---|
| **Behavioural validation** | Does the memory concept hold? | Blocked until L1–L3 pass |
| **Runtime performance validation** | What does this hardware do, and stably? | Not started; observational only |
| **Howl integration** | Does the application deliver it? | Deferred — after both are independently understood |

A result in one stream never becomes evidence in another. This is the same discipline that kept the corpus stream and the engineering stream apart.

## 9. Reproducibility block (every local run)

Server name and **version**; model file name **and hash**; quantisation; `n_ctx` as loaded; GPU layer count / offload mode; sampler values actually sent; seed and whether the server honours it; the endpoint path used; surface (chat vs native); harness version; fixture id.

**A run missing any of these is inadmissible** — the same standard as hosted, with version and model hash replacing provider identity.

## 10. Risks, and the handling

| Risk | Why it bites | Handling |
|---|---|---|
| Silent context truncation | Replies look fine, history quietly dropped | Record `n_ctx` and prompt token count; flag near-limit runs |
| Sampler divergence | KoboldCpp's `top_k`/`min_p`/`rep_pen` are not HF's | Record values sent; compare only like-for-like conditions |
| Version drift between runs | Same model, different server behaviour | Version in the identity block; a version change is a **new condition** |
| CPU/GPU offload variance | Layer count changes numerics slightly | Fixed per series; change = new condition |
| Slow generation | Temptation to shorten prompts or cut reps | Prompts frozen by fixture; never traded away for speed |
| Monitoring as a new instability source | Owner just fixed a driver fault | Read-only sampling, no tuning |

## 11. Acceptance criteria per phase

| Phase | Passes when |
|---|---|
| **L1** | All six §4 behaviours have a fixture and a passing check; both surfaces map to one response type; suite stays green |
| **L2** | An identity block is recorded for the listening server, or the absence is documented as `unavailable` |
| **L3** | One call returns a body the L1 parser handles without correction, and one record is admissible |
| **L4** | The S1 18-call contrast runs under its own authorisation with complete provenance |

## 12. Not authorised by this plan

Server startup; any HTTP request to a local endpoint; model loading; any generation; S1 execution; Howl changes; corpus or standards changes; deleting or rewriting the hosted adapter (it stays as unreachable, unverified code); installing binaries or drivers.

## 13. Claim ceiling

A plan and a feasibility assessment. No server has been contacted, no model loaded, no prompt sent, no performance measured. It defines what would have to be true for the local route to be *validated* rather than merely *chosen* — and the difference between those two words is the whole reason the hosted route failed.
