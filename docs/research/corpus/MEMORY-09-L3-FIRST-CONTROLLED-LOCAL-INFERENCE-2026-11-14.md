# MEMORY-09 — L3: First Controlled Local Inference (WO-L3-01)

**Status:** L3 requirement **met** (one real local generation through the harness, raw evidence preserved, provenance reproducible). MEMORY-09 is **not** evaluated, S1 has **not** started, Howl Chat is untouched, and no behavioural claim is made anywhere in this document.

**Order:** `WO-L3-01`, owner-authorized 2026-11-13.
**Run:** 2026-11-14T08:03:03Z, `eval/runs/L3-2026-11-14T080303Z/`.
**Ceiling:** one call, and one call was made.

---

## 1. Observed facts

### 1.1 Server identity (endpoints, not the executable)

| Field | Value | Source |
|---|---|---|
| Server | KoboldCpp **1.120** | `GET /api/extra/version` |
| Endpoint | `http://127.0.0.1:5001` | loopback only |
| Model resident (pre-load, model-free instance) | `inactive`, `"llm": false` | `/api/v1/model`, `/api/extra/version` |
| Model resident (post-load) | `koboldcpp/Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M` | `/api/v1/model` |
| Server-reported context | 20480 | requested at launch; not re-read post-load |

### 1.2 Preset, and the path it did not resolve

`H:\AI apps\qwen 9b.kcpps` — **3,591 bytes**, sha256 `9ed540a5da4954f102f92fb70de26b87607eef8068446bbcc786f9f1f3a8db2b`.

| Preset key | Value |
|---|---|
| `model_param` | `C:/Users/StoneWolfPC/Downloads/Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` — **absent; `Downloads` holds no `.gguf` at all** |
| `contextsize` | 20480 |
| `gpulayers` | 30 |
| `threads` | 8 |
| `batchsize` | 512 |
| `quantkv` | `q8_0` |
| `usemmap` | true |
| `usecuda` | `["normal", "0"]` |
| `host` | `""` — **binds all interfaces; overridden to `127.0.0.1`** |
| `launch` | true — **overridden; no browser was opened** |
| sampler settings (`temp`, `top_p`, `top_k`, `min_p`, `rep_pen`, `gendefaults`) | **absent from the preset** → server defaults, recorded as unavailable |

**Resolution:** the same filename exists at `G:\AI GGUF\GGUF\Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` and is named by the project's own model registry (`config/model-registry.json`). The owner directed the search after the `Downloads` path failed. **Only the directory was substituted; the filename, size and hash were verified independently of the preset.**

| Model file | Value |
|---|---|
| Path | `G:\AI GGUF\GGUF\Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` |
| Bytes | 5,627,044,224 |
| SHA-256 | `2ca636d9e81d3d23ca9b60c234fe185d30ec082eeba69ce770fdb0c76559a4f5` |
| Quantisation | **Q4_K_M appears in the filename only** — GGUF metadata was not parsed (see §4) |

### 1.3 Launch

```
koboldcpp.exe --model "G:\AI GGUF\GGUF\Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf" \
  --contextsize 20480 --gpulayers 30 --threads 8 --batchsize 512 --quantkv q8_0 \
  --usemmap --usecuda 0 --host 127.0.0.1 --port 5001 --skiplauncher --quiet
```

The model-free L2 instance (PID 24824 child, 11032 bootloader) was stopped first, because KoboldCpp loads a model at launch and offers no load-by-API path without admin. **Port 5001 answered `000` between the two processes** — no overlap.

Load completed in **≈28 s** to a responsive server. Log recorded `Load Text Model OK: True`, `Chat completion heuristic: ChatML (Generic)`, and `Active Modules: TextGeneration`.

### 1.4 The one request

| Field | Value |
|---|---|
| Test id | `L3-01` |
| System | "You are a controlled inference check. Answer in one short sentence." |
| User | "Name the colour of a clear midday sky." |
| Params | `temperature=0.0`, `max_tokens=48` |
| Seed | none sent |
| Streaming | false |
| Material | synthetic only — no personal, sensitive or corpus-derived text |

### 1.5 The one response

| Field | Value |
|---|---|
| HTTP status | 200 |
| Surface | `chat_completions` |
| Latency (client-measured) | **2,875 ms** |
| Usage (server-supplied) | `prompt_tokens=37`, `completion_tokens=48`, `total_tokens=85` |
| `finish_reason` | **`length`** → `incomplete = True` |
| Text length | 173 characters, beginning ` thinking\nThinking Process:\n\n1. **Analyze the Request:**…` |
| Body keys | `choices`, `created`, `id`, `model`, `object`, `usage` |

The response was **truncated at the token limit**, not completed. That is a recorded outcome, not a failure to be retried: the ceiling was one call and a repeat would be a new authorisation.

---

## 2. Calculated measurements

- **Generation speed:** 48 completion tokens in 2,875 ms ≈ **16.7 tok/s**. This includes prompt evaluation and network/loopback overhead, so it is an upper-bound-ish composite, not a pure decode rate.
- **Prompt evaluation:** 37 prompt tokens; not timed separately.
- **VRAM delta across the load:** 6,573 − 1,160 = **5,413 MiB** attributable to model + KV cache at 20,480 context with `q8_0` KV.
- **VRAM headroom after load:** **1,487 MiB** free of 8,192 MiB.

## 3. Inferences (stated as such, not as facts)

- The model is a **reasoning-style** model: it opened a `Thinking Process` block inside `message.content`, which consumed the entire 48-token budget before naming a colour. The truncation is the **expected consequence** of `max_tokens=48` against that behaviour, not evidence of a defect.
- With `temperature=0.0` the effect is still preamble-first; sampling temperature is not the cause.
- `Q4_K_M` is very likely the quantisation given the filename convention, but this was **not** read from GGUF metadata.
- 30 GPU layers plus 20,480-context `q8_0` KV fits in 8 GB **for this model on this build**; that is one observation, not a compatibility guarantee.

## 4. Unavailable or unverified

| Item | Why |
|---|---|
| Quantisation from metadata | GGUF header not parsed; filename is a label, not evidence |
| Sampler defaults actually applied | absent from preset; server defaults were not read back |
| Model revision / training provenance | not exposed by the endpoint |
| True `n_gpu_layers` and per-layer split | not read back from the server after load |
| CUDA runtime version used by the backend | log states `Auto Selected CUDA Backend (flag=0)`; no version string captured |
| Temperature/other hardware sensors | only GPU temperature is readable here |
| CPU utilisation | not sampled |
| Wall-clock load start | not timestamped precisely; "≈28 s" is derived from the readiness poll interval |

## 5. Requirement results

| Requirement | Result | Evidence |
|---|---|---|
| Confirm server at `127.0.0.1:5001` | **PASS** | identity probe, `listening: true` |
| Confirm version 1.120 | **PASS** | `/api/extra/version` |
| Inspect preset, resolve model unambiguously | **PASS with correction** | `Downloads` path absent; same filename + size + sha256 resolved on `G:` under owner direction |
| Record file identity before load | **PASS** | path, 5,627,044,224 B, sha256 `2ca636d9…` |
| Record preset contents and hash | **PASS** | 3,591 B, sha256 `9ed540a5…` |
| Record intended context / GPU layers / samplers | **PASS** | 20480 / 30 / unavailable (absent from preset) |
| Baseline hardware recorded before load | **PASS** | 1,160 MiB used, 6,900 free, 18 %, 27 °C |
| Load model, confirm resident and responsive | **PASS** | `/api/v1/model` → `koboldcpp/Qwen3.5-9B…`; `TextGeneration` active |
| No model substitution | **PASS** | granted sha256 matched the loaded file; no other preset used |
| Exactly one generation | **PASS** | ledger `attempted=1/1`, `per_slot={"L3-01": 1}` |
| No retries | **PASS** | single POST; failure paths in the adapter return records, they do not loop |
| Raw response preserved | **PASS** | `raw_body.json`, 668 B |
| Request + effective params preserved | **PASS** | `request_echo` in the record |
| Latency, tokens, speed recorded | **PASS** | 2,875 ms; 37/48/85; ≈16.7 tok/s |
| Completeness determined | **PASS** | `finish_reason=length`, `incomplete=true` |
| Record write-once | **PASS** | `record.json` + `record.sha256`; re-write raises |
| Record hash matches evidence | **PASS** | `l3-verify` → `hash vs raw: True` |
| Parsed matches raw | **PASS** | re-transcription identical |
| Model/runtime provenance present | **PASS** | runtime block: file, bytes, sha256, preset hash, server version, endpoint, identity block hash |
| Missing fields marked | **PASS** | state hashes and samplers `unavailable` |
| One-call limit respected | **PASS** | 1 of 1 |
| No hosted/network inference path | **PASS** | `provider_used = http://127.0.0.1:5001` |

**Verification 10/10 pass** (`python -m eval.harness.cli l3-verify --dir eval/runs/L3-2026-11-14T080303Z`). Offline suite **87/87**. Corpus registry **PASS**.

## 6. A defective check, corrected — and why it mattered

The first verification pass reported `no_hosted_route = False` on a call that plainly went to loopback. The predicate tested `"local" in provider_used`, but `provider_used` is the **endpoint URL** (`http://127.0.0.1:5001`), not the provider label. The second failure was `state_hash_before/after not recorded`, which is correct for L3: there is no memory state, so absence is the true value.

Both were defects in the **checker**, not the run. They were fixed by parsing the endpoint host against the loopback set and by treating state as **not applicable** with an explicit note. Because the call was already spent and retries are prohibited, the correction was validated by re-checking the **filed** record with a new `l3-verify` command that makes no calls — evidence was re-read, not re-created.

That the checker could be wrong in the direction of "looks like hosted" is the useful lesson: a verification that fails safe is worth more than one that passes quietly.

## 7. Final server state

Left **running**, per §F: `koboldcpp.exe` with the Qwen 9B model resident, `127.0.0.1:5001` only, VRAM 6,573 MiB used, GPU 3 %, 30 °C. No unload, no termination, no further inference.

## 8. What this does and does not establish

**Establishes:** the harness can perform one genuine local generation, capture the unedited body, transcribe it reproducibly, and file evidence whose hash, provenance and completeness survive re-verification.

**Does not establish:** that MEMORY-09 has any effect; that S1 passes; that Howl Chat integrates; that repeated runs are independent raters; that the model is reliable; that this configuration is stable under load. One truncated 48-token reply is an instrument check, not a result.
