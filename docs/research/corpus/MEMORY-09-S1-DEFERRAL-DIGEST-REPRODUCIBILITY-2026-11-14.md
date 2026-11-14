# MEMORY-09 — S1 Deferred: Digests Not Reproducible On This Machine

**Status:** S1 behavioural execution **deferred**. No generation was made under this record.
**Date:** 2026-11-14 · **Decision by:** owner · **Manifest in force:** S1-v1.1 (v1.2 withdrawn, never used)

---

## 1. Why

The S1 identity gate compares each model file's SHA-256 against the frozen manifest before any qualification is recorded. On 2026-11-14 the same files returned **different digests at different times**, with byte counts and modification times unchanged throughout.

| File | Earlier reading | Later reading | Frozen baseline |
|---|---|---|---|
| S1-A on `G:` | `0e124fab…` | **`2ca636d9…`** | `2ca636d9…` |
| S1-A on `H:` | `de3b8d28…` | `cd4b4430…` | — |
| S1-B on `G:` | `290f7ead…` | `341fb28c…` (twice) | `b5b1913e…` |
| S1-B on `H:` | `6ff1b33f…` | `01c357cd…` | — |
| S1-C on `G:` | `f820ef3f…` | `f820ef3f…` | `f820ef3f…` |
| S1-C on `H:` | `7fca804a…` | `7fca804a…` | — |

Byte counts were constant for every file (`5,627,044,224` / `5,335,285,728` / `5,155,565,440`), and S1-B's modification time stayed at `2026-09-04 20:53:48`.

## 2. What that means

**`G:` S1-A returned its historical baseline digest again** after previously reading differently. Content that genuinely changed does not drift back to its own past hash, and it does not do so without moving the modification time. The simplest explanation consistent with all six observations is that **the reads vary while the stored bytes do not** — corrupted pages held in the file cache, which is why consecutive reads within one session agree and differ from the next session.

This **reverses** the interpretation recorded earlier in the day. The originals on `G:` are the more likely intact copies; the `H:` copies are more likely faithful reproductions of a corrupted read. On that basis:

- `eval/s1-readiness-baseline-v1.2.json` was **withdrawn** (deleted) before any qualification ran under it.
- No re-baseline onto the `H:` copies was frozen.
- No model was copied, deleted or modified during any of this.

## 3. Functional state, unchanged by the above

All three candidates on both drives load and generate: 6/6 loads, 6/6 generations, zero load errors, identical GGUF structure per pair, VRAM 49–70 % and 35–42 °C. The machine reads these files well enough to run them. It cannot read them the same way twice.

## 4. Decision

**S1 behavioural execution is deferred until digests are reproducible on this machine.** The experiment's provenance rests on being able to state which bytes were executed; a manifest that changes every hour cannot support that statement.

## 5. Resumption criteria

All three must hold before S1 is reconsidered:

1. The driver-stack fault (iCUE / Afterburner, established separately) is resolved or the replacement is validated.
2. **Three consecutive cold reads of each candidate return identical digests**, and those digests match the manifest in force.
3. A preflight returns PASS with hashing enabled (`preflight-s1` without `--no-hash`).

Until then the preflight is expected to report FAIL, and that FAIL is correct behaviour rather than a defect to be worked around.

## 6. Boundaries observed

No generation. No S1 call. No deletion, no copy, no conversion, no repair. No model file written. The v1.1 manifest, the v1 manifest and every qualification record remain on disk unchanged.

## 7. Harness change made alongside this record

`eval/harness/cli.py::_baseline_models` read `eval/s1-readiness-baseline.json` directly, so qualification kept gating on **v1** hashes after v1.1 superseded them, while the preflight used the newest manifest. It now uses the preflight's selection rule. Offline validation after the change: **89/89 passed**.
