# `MEMORY-09` — Model Adapter Interface (Python-First)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-MADP-1` · **Version:** 1.1 · **Date:** 2026-11-13 · **Status:** **INTERFACE SPEC — no adapter implemented, no call made.**
**Revision 1.1:** transport fixed to one documented client; provenance expanded; failure handling made fail-closed (owner review, 2026-11-13).
**Revision 1.2:** local route admitted as a loopback-only second transport; spend policy applied before construction (row 484).
**Authority:** owner direction (Python-first follow-on to `WO-ARCH-06`).

---

## 1. Why an adapter, not a direct call

The harness must **not depend on Howl's Go runtime.** Coupling the independent test environment to an application still under development would make every result a statement about that application's state, not about the memory concept.

The adapter is the seam: **one interface, several backends**, so the same tests run against hosted inference now and `llama.cpp` later.

## 2. Interface

```
generate(request) -> Response
```

**Request:** `messages` (exact text), `model_id` (hub id + revision), `provider` (explicit), `params` (temperature, top_p, top_k, min_p, penalties, max_tokens, stop), `seed`, `stream=false`.

**Response (the record, not the text):**

| Field | Why |
|---|---|
| `text` | Convenience only |
| `raw_body` | **The primary evidence** — unmodified provider response |
| `request_echo` | Exactly what was sent, params as serialised |
| `provider_requested` / `provider_used` | The pinned value, and what actually served |
| `model_requested` / `model_reported` | Including revision where the API returns it |
| `system_fingerprint` | Where supplied — the only backend-identity hint available |
| `usage` | Token counts if provided |
| `latency_ms`, `timestamp` | Informational |
| `request_id` | Correlation |
| `client_version` | Library version that made the call |
| `error` | Structured on failure |

**Rule: `text` is a convenience; `raw_body` is the evidence.** A record without `raw_body` is not admissible.

## 3. Transport — one recipe, chosen explicitly

Two related interfaces exist and **must not be mixed**:

| Interface | Pinning mechanism | Status |
|---|---|---|
| `huggingface_hub.InferenceClient` | separate `provider=` argument (default unset — routes automatically; documented as `auto`) | **CHOSEN** |
| OpenAI-compatible router (`/v1/chat/completions`) | provider suffix on the model id (`model:groq`) | Recorded as an alternative; **not used** |

**Decision: `huggingface_hub.InferenceClient` is the sole transport for the first adapter**, because it was the Python-first choice and exposes provider selection as a first-class argument rather than string surgery on the model id.

**Rule T-1:** `provider="auto"` is **prohibited** in evaluation calls. Every call names its provider. **An unset provider is prohibited too:** the installed client (1.29.0) defaults the argument to `None`, not to the string `"auto"`, so a harness that accepted "not literally auto" would silently route — the adapter rejects unset, empty and `auto` alike.

**Verified 2026-11-13** by signature inspection of `huggingface_hub` 1.29.0: `InferenceClient(provider=...)` exists, and `chat_completion` takes `model`, `messages`, `stream` plus sampling parameters. No call was made.
**Rule T-2:** the model id is passed unmodified — no `:provider`, `:fastest`, `:cheapest`, `:preferred` suffixes, which belong to the other interface. Mixing the two recipes is the specific error this section exists to prevent.

## 4. Provider pinning is necessary, not sufficient

Pinning stops the harness choosing a different provider. It does **not** establish that the provider serves the same backend revision, quantisation or inference stack across dates. Provider metadata is recorded, and **its absence is recorded as `unavailable`, never inferred**.

**Rule P-1:** a provider change is a **new condition**, not a rerun.
**Rule P-2:** a provider silently changing its backend is a **possible confound**, recorded as such — not assumed away.

## 5. Failure handling — fail closed

| Class | Handling |
|---|---|
| Transport/auth | Fail fast; record; no retry storm |
| Rate limit (429) | Bounded backoff, same request id, max 3 |
| Provider unavailable (503) | **Record failure for the pinned provider. No fallback.** |
| Empty/truncated content | Record `incomplete` — never a model failure |
| Content filter | Record verbatim; the filter's decision is part of the observation |

**Rule F-1:** a failed pinned call **must never** trigger automatic fallback to another provider. The condition is preserved; trying another provider requires an explicit new run with its own record.

## 6. Determinism

**Rule D-1:** seeds are sent when the interface allows, and marked `honoured: unknown` unless the provider documents support **and** a test demonstrates it. Hosted seeding is not assumed. T-01 measures spread rather than presupposing none.

## 7. Capability declaration

```
capabilities = {seed, logprobs, stop_sequences, streaming, max_context}
```

A test needing a capability the adapter lacks is **`NOT RUN`**, never approximated.

## 8. Provenance block (expanded)

Per call: adapter name/version, **client library version**, hub model id **and revision where available**, provider requested and used, params sent, returned model identifier, `system_fingerprint` if supplied, prompt hash, response hash, state hashes before/after, fixture id/version, harness version, test id, timestamp, provider headers where present. **Any field the API does not supply is stored as `unavailable`.** A run missing any field without that marker is not admissible.

## 9. Local adapter (later stage, same interface)

A loopback HTTP adapter for **KoboldCpp, SillyTavern, `llama.cpp` and other OpenAI-compatible local servers**, implementing the identical contract and adding `n_ctx`, quantisation, GPU layers and peak RAM/VRAM at S3. A test written for hosted inference runs locally unchanged; only adapter and identity differ. The per-family defaults (row 472) re-enter here as *configuration*, reproduced rather than imported from the Go implementation.

**Loopback only:** `localhost`, `127.0.0.1`, `::1`. A non-loopback host is refused — a remote host wearing a local name is a hosted provider and must face the spend policy. No default endpoint or port is guessed; a local route without an explicit endpoint is an error.

**Route selection:** the provider string decides (`local:<runtime>` → local adapter), and `build_adapter` runs the spend policy *before* constructing anything (row 484). A provider change is still a new condition, and nothing falls back across routes.

## 10. What the adapter cannot do

Score, judge, interpret; change parameters between calls; substitute one model for another; report success for a response it did not receive.

## 11. Claim ceiling

An interface specification. No adapter exists, no token spent, no model called. It establishes what a compliant adapter must record — not that any provider is consistent across dates, which only repeated measurement can show.
