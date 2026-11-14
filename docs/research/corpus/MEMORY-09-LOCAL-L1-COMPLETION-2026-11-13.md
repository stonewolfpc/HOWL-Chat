# `MEMORY-09` — Local Adapter L1: Serialisation and Transcription

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-LAL1-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **L1 COMPLETE — offline only. No server contacted, no model loaded, no prompt sent.**
**Authority:** owner direction, 2026-11-13 (local adapter, KoboldCpp first, no live test until the adapter is inspected).
**Plan:** row 486 §3–§4.

---

## 1. Interface inspection — result, including the blocker

| Check | Result |
|---|---|
| `koboldcpp` / `koboldcpp.exe` on PATH | **not found** |
| `llama-server` / `llama-cli` / `main.exe` on PATH | **not found** |
| KoboldCpp or `llama.cpp` tree inside the project | **absent** |
| Model inventory | `config/model-registry.json` lists GGUF files with names, paths, sizes and per-model recommendations (e.g. `Rocinante-X-12B-v1-absolute-heresy.Q4_K_M`, 7.1 GB) |

> **Revision note (2026-11-13, after the owner supplied a location).** The search above covered PATH and the project tree only. The binaries are on another drive: `H:\AI apps\` holds `koboldcpp.exe` (607 MB), two `llama.cpp` CUDA 12.4 builds (b10076, b9596) and SillyTavern 1.18.0. The "no binary found" conclusion is therefore **a scope limit, not an absence**, and is superseded by `MEMORY-09-LOCAL-INTERFACE-DISCOVERY-2026-11-13.md` (row 488). The text below is retained as the record of what that search could and could not establish; KoboldCpp's **version** remains unavailable from the filesystem.

**Conclusion, stated precisely:** models are registered, but **no inference server binary was found** *within the searched scope*. The installed KoboldCpp/`llama.cpp` *versions* therefore cannot be inspected, and a version determines the exact API surface. Selection of KoboldCpp as first target is **provisional** — decided on documented capability, not on a verified local build.

**This blocks L2, not L1.** L1 was designed to proceed without a server, and did.

## 2. What L1 built

| File | Bytes | Role |
|---|---:|---|
| `eval/harness/local_surfaces.py` | 6,543 | Both request serialisers, surface detection, both transcribers, error handling |
| `eval/harness/local_fixtures.py` | 2,926 | Nine synthetic bodies + one unrecognised body |
| `eval/harness/adapter.py` | 12,520 | `finish_reason` field and `incomplete` property added to the response record |

Checks **63–74** (twelve new). Suite: **74/74, exit 0**, socket guard active, 0 scratch.

## 3. Design decisions recorded

| Decision | Reason |
|---|---|
| **Both surfaces supported** (chat-compatible and native) | KoboldCpp exposes both; choosing one is a design decision, not a formatting detail |
| **Native prompt is flattened**, template named `system\n\nuser` | The native surface takes one prompt string. Flattening changes the prompt text, so **a cross-surface comparison is a different condition, not a re-run** |
| **Parameter translation is recorded** (`max_tokens → max_length`, `presence_penalty → rep_pen`, …) | A silent rename is how sampler values drift between harness and server |
| **Unmappable parameter refused, not dropped** | Dropping a sampler setting changes the condition silently; check 65 fails if a parameter vanishes |
| **Truncation is `incomplete`, not an error** | `finish_reason == "length"` is a fact about the request, not a model failure |
| **Empty results and error bodies become records** | A crash loses evidence; a record labelled `empty results` keeps it |

## 4. Defect found by the fixtures, and fixed

The first run failed check 72: an error body (`{"detail": …}`) carries **no** `choices` and **no** `results`, so surface detection raised **before** the error path ran. That is precisely what a local server's 500 looks like — the most common failure would have been the one that crashed the harness.

**Fix:** `transcribe_any` now checks for `error`/`detail` **first**, labels a shapeless body `surface="unknown"`, and returns a record. A regression case was added for `{"error": "boom"}`.

## 5. What L1 establishes, and what it does not

**Establishes:** both request bodies serialise deterministically from one `GenerateRequest`; parameter translation is explicit and lossless or refused; both documented reply shapes map to a single record type with identical fields; truncation, empty results, unknown fields and error bodies are all handled without raising.

**Does not establish:** that a real KoboldCpp or `llama.cpp` server returns these shapes. Every fixture is **documentation-derived, not captured**. The parser is written to follow the provider — if L3 contradicts a fixture, the fixture is corrected and L1 re-run, not the body reinterpreted.

## 6. Next gate

**L2 (endpoint discovery) needs two things that do not exist yet:** an authorisation for loopback HTTP, and a path to a running server. Without a binary path from the owner, L2 cannot start — and no amount of further offline work changes that.

## 7. Claim ceiling

Twelve checks of serialisation and transcription logic against synthetic bodies. No server, no model, no prompt, no measurement. Nothing here says local inference works; it says the harness will not mistake a truncated reply, an empty result or a 500 for a normal answer.
