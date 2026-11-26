# Sampler Capture — Route A Result

**Identity:** `MEMORY-09-SAMPLER-A-1` · **Date:** 2026-11-15 · **Authorised:** owner ("proceed with Route A only, with zero behavioral generations authorized")
**Question:** can the server's effective sampler chain be observed **without generating**?
**Answer:** **No — not in this build, at `--debugmode 1`.**

## 1. Safeguard: observability only

The owner required proof that enabling logging does not change the experiment. The launch `Namespace(...)` was parsed from both logs (171 arguments each) and compared field by field.

| Field | Qualification run | Probe run |
|---|---|---|
| `debugmode` | `0` | **`1`** |
| `quiet` | `True` | **`False`** |
| everything else (169 fields) | — | **identical** |

`gpulayers=20`, `contextsize=20480`, `quantkv='q8_0'`, `threads=8`, `batchsize=512`, `usemmap=True`, `usecuda=['0']`, `gendefaults=''` — all unchanged. **No model, sampler or runtime setting differs.** The only two differences control verbosity.

## 2. What was requested

Four non-generating requests, nothing that produces tokens:

| Request | Status | Response |
|---|---|---|
| `GET /api/extra/version` | 200 | KoboldCpp 1.120, `llm: true` |
| `GET /api/v1/model` | 200 | `koboldcpp/Gemma-4-E4B-…Q4_K_M1` |
| `POST /api/extra/generate/check` | 200 | `{"results": [{"text": ""}]}` — **no completion text** |
| `POST /api/extra/tokencount` | 200 | accepted |

## 3. What appeared

**Nothing.** No `sampler chain`, `sampling:`, `samplers order`, `top_k =`, `top_p =` or `temp =` line exists in the log — not at startup, not after any of the four requests. Zero error lines.

Searches run against: `eval/runs/kobold-S1B-v12-samplerprobe.log` (12 KB, 171-key Namespace line included).

## 4. What this establishes

1. **At `--debugmode 1`, this build does not print the sampler chain** during load or on a dry-run check. Verbosity is not enough.
2. **`gendefaults=''` is confirmed again** — no CLI sampler overrides have ever been in play, in any run.
3. The effective values therefore remain **unobserved**; `temperature=0.8` and `top_p=0.95` are still only the *request* parameters.

## 5. What it does not establish — the trap in the obvious next step

It does **not** establish that one generation would reveal the chain. It is equally consistent that this build never logs it at any verbosity. So Route B (one generation) is **not guaranteed to succeed**: it may spend a generation and return the same nothing.

The honest framing before any spend:

| If the chain appears during generation | If it does not |
|---|---|
| samplers captured at a cost of 1 generation, and the S1 request parameters can be checked against the server's actual chain | samplers are unobservable this way at all, and the only remaining evidence is behavioural (output distribution) or the llama.cpp defaults as compiled in |

Anyone authorising Route B should authorise it knowing it may buy nothing.

## 6. Machine state

Load with `debugmode 1`: **4,110 MiB** used (4,041 MiB on the plain qualification run — **69 MiB more**, observed but unexplained; logging buffers are a plausible and untested cause). 30 °C.

After stop: **no KoboldCpp process, no listener on 5001, 1,794 MiB used**, 6,266 free. Clean unload; this run added no qualification record, because a plain load at `debugmode 1` is not a new configuration.

## 7. Disposition

Route A is **exhausted**. No generation was spent: budget remains **1 of 30**. The next decision belongs to the owner: spend one generation on a possible capture, or record the request parameters as an explicit, labelled assumption.
