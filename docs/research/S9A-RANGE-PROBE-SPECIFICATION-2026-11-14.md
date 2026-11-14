# `S9A-RANGE-PROBE-SPECIFICATION`

**Identity:** `HOWL-S9A-SPEC-1` · **Date:** 2026-11-14 · **Status:** **SPECIFICATION ONLY — not authorised, not executed.**
**Scope:** a bounded, read-only comparison of two access paths on one fixed byte range of one **non-model** file, plus an optional ground-truth variant that requires gate 1.
**Precedence:** this document replaces `STORAGE-DIAGNOSTIC-STAGING-PLAN` S-9a. The staging plan's version described an experiment it could not perform (per-file cache invalidation); this one does not.

---

## 1. What the probe is for, and what it is not for

| | |
|---|---|
| **Question** | Do two access paths return different bytes for the same range of the same file identity, under conditions the spec fixes in advance? |
| **What a divergence establishes** | that the two paths disagreed on that range at those two moments, with the range and identity recorded |
| **What a divergence does not establish** | which path is wrong, which layer caused it, or anything about the S1 model files |
| **What agreement establishes** | **nothing about correctness** — a stale cache and a bad device agree whenever both serve the same wrong data |
| **Explicitly out of scope** | cache purging, cache-size changes, driver changes, stress, CHKDSK, surface scans, any `.gguf`/`.safetensors`/`.bin` payload, any write to a model destination |

### 1.1 The cache-control problem, stated before anything else

The staging plan said "drop the file's cache". **That is not available**, and the spec must not pretend otherwise:

| Route | Reality |
|---|---|
| `FlushFileBuffers` on a read handle | **fails** — the handle must have `GENERIC_WRITE` |
| `SetSystemFileCacheSize` / `NtSetSystemInformation` global cache control | requires elevation and changes system state; out of scope |
| Server-side cache invalidation | not applicable — this is a local fixed volume |
| Unmounting / dismounting the volume | destructive to the investigation; forbidden |

So this spec performs **no invalidation**. It measures two access paths and their *existing* cache state, and it says so in every output. The consequence is not fatal: T-15 already produced a divergence without any invalidation, and a repeat with tighter accounting still narrows what the divergence is.

---

## 2. Variants

| Variant | Target | Ground truth | Gate |
|---|---|---|---|
| **V-A** | an existing non-model file on `G:` | **none** — content is unknown by construction | `G:` payload read |
| **V-B** | a scratch file written with a known pattern | **yes** — the pattern is the oracle | **gate 1** (`G:` disposable file) |

**V-B is the scientifically stronger test** and should be preferred when gate 1 opens: it can say not just "the paths disagreed" but "one path returned the bytes that were written". V-A is weaker but needs a smaller authorisation. Neither reads a model.

---

## 3. V-A — existing non-model file

### 3.1 File selection rule (deterministic, decided before the run)

Candidate must satisfy **all** of:

1. located on `G:`, outside `G:\AI GGUF\GGUF`;
2. extension not in `.gguf .safetensors .bin .ckpt .pt .pth .onnx`;
3. size ≥ 64 MiB and ≤ 512 MiB (large enough to defeat trivial caching, small enough to bound I/O);
4. size and mtime unchanged across two metadata reads ≥ 60 s apart;
5. not openable exclusively would be fine — but the probe **must not** take an exclusive handle if a legitimate process may hold it; use `FILE_SHARE_READ | FILE_SHARE_WRITE`.

Record the selected path, size, mtime (ns), volume serial, file id, and the selection rule's evaluation. If no file satisfies the rule, **stop** — do not relax the rule during the run.

### 3.2 Range

| Parameter | Value | Reason |
|---|---|---|
| Offset | a fixed 4096-aligned offset chosen before the run and recorded (e.g. `1 GiB`) | alignment for unbuffered I/O; reproducibility |
| Length | `64 MiB` | bounded |
| Block size | `4 MiB` | matches the existing `unbuffered_digest.py` reader |
| Sector alignment | 4096 | `G:` is logical 512 / physical 4096; 4096 satisfies both |
| Buffer alignment | 4096 (page-aligned allocation) | required for `NO_BUFFERING` |

### 3.3 Procedure

| Step | Action | I/O |
|---|---|---|
| A-1 | capture identity + metadata (`size`, `mtime_ns`, `volume_serial`, `file_id`, `in_use`) | metadata |
| A-2 | buffered read of the range → `sha256_buffered` | 64 MiB |
| A-3 | capture identity + metadata again; **abort if any field changed** | metadata |
| A-4 | unbuffered read of the range (`FILE_FLAG_NO_BUFFERING`, 4 MiB blocks) → `sha256_unbuffered` | 64 MiB |
| A-5 | capture identity + metadata again; abort if changed | metadata |
| A-6 | **repeat the unbuffered read once** → `sha256_unbuffered_2` | 64 MiB |
| A-7 | report | — |

Total payload: **192 MiB**, on a file that is not a model. `A-6` exists because a single differing pair cannot distinguish a stable alternative content from a one-off bad read — the exact gap T-15 left.

### 3.4 Tail handling — the T-15 confound removed

T-15's unbuffered digest covered the body unbuffered and 4,064 trailing bytes buffered, so a difference confined to the tail would have been misattributed. Here the range is a whole number of 4 MiB blocks (64 MiB = 16 blocks), so **no tail exists**. If the range must be shortened for any reason, hash body and tail **separately** and report both digests; never a single combined value.

---

## 4. V-B — known-pattern scratch file (gate 1)

| Step | Action | Notes |
|---|---|---|
| B-1 | create `G:\AI GGUF\_probe\<random>.bin`, size 64 MiB, written with `NO_BUFFERING` from a deterministic PRNG stream (seed recorded) | the write itself is the only write |
| B-2 | `FlushFileBuffers` on the write handle | makes the written content durable |
| B-3 | close, capture identity | — |
| B-4 | unbuffered read of the whole file → compare with the expected digest | **oracle check**: if this differs, the file was not written as intended; stop and report |
| B-5 | buffered read of the whole file → `sha256_buffered` | — |
| B-6 | unbuffered read again → `sha256_unbuffered` | — |
| B-7 | compare all three against the expected pattern digest | — |

Interpretation becomes decidable: if B-4 matched and B-5 differs, the **cached** stream is wrong. If B-4 matched and B-6 differs, the **unbuffered** stream is unstable. If B-4 fails, the *write* path is implicated and nothing about reads can be concluded from this run.

Cleanup is itself a write; the spec requires it to be recorded, and the file left in place until the evidence is filed.

---

## 5. Output schema

```json
{
  "spec": "HOWL-S9A-SPEC-1",
  "variant": "V-A",
  "started_utc": "...", "ended_utc": "...",
  "path": "...", "range": {"offset": 1073741824, "length": 67108864, "block": 4194304, "sector": 4096},
  "identity_before": {"size": 0, "mtime_ns": 0, "volume_serial": "0x...", "file_id": "...", "in_use": false},
  "identity_after": { "...": "..." },
  "digests": {"buffered": "...", "unbuffered": "...", "unbuffered_2": "...", "expected": null},
  "io": {"read_ops_delta": 0, "read_bytes_delta": 0, "elapsed_s": 0},
  "verdict": "DIVERGED_STABLE | DIVERGED_UNSTABLE | AGREED | ABORTED",
  "limits": ["no cache invalidation performed", "no ground truth for this file (V-A)"]
}
```

`DIVERGED_STABLE` = buffered differs from both unbuffered runs, which agree with each other. `DIVERGED_UNSTABLE` = the two unbuffered runs disagree. `AGREED` = all three equal.

---

## 6. Interpretation matrix

| Verdict | Reading | Not established |
|---|---|---|
| `DIVERGED_STABLE` (V-A) | the two paths hold different, individually reproducible content for one range | which is right; the layer; anything about models |
| `DIVERGED_STABLE` (V-B, B-4 passed) | **the cache holds content the device does not** | how the cache became wrong |
| `DIVERGED_UNSTABLE` | the unbuffered path itself is not reproducible on this range | whether that is device, controller, driver or memory |
| `AGREED` | the paths agreed on this range, this time | **nothing.** Agreement is not correctness, and a stale cache agrees with a bad device |
| `ABORTED` | a precondition failed | nothing — an aborted probe is not a negative result |

---

## 7. Stop conditions

Stop immediately, preserve output, and do not retry if:

1. any identity field changes between captures;
2. any read returns an error, short read, or `ERROR_CRC` / `ERROR_DEVICE_HARDWARE_ERROR`;
3. the file becomes `in_use`;
4. the unbuffered call fails with `ERROR_INVALID_PARAMETER` (alignment defect in the probe, not in the device);
5. wall-clock exceeds **10 minutes**;
6. total payload read exceeds **256 MiB**.

A stop is a result. Retrying a stopped probe against the same range is a second experiment and needs its own authorisation.

---

## 8. Review checklist — required before authorisation

| # | Item | Why |
|---:|---|---|
| 1 | The exact command line is in the work order, not improvised at the keyboard | reproducibility |
| 2 | The reader's block, sector and buffer-alignment handling is reviewed line by line | an alignment defect produces `ERROR_INVALID_PARAMETER`, or worse, a silently wrong read |
| 3 | The reader reports the identity before, between and after each pass | the T-15 lesson |
| 4 | Body and tail are hashed separately, or the range is block-exact | the T-15 confound |
| 5 | The spec states "no cache invalidation" in the output, not just in the document | a reader of the JSON must not infer a purge happened |
| 6 | The selected file satisfies §3.1 by rule, evaluated before the run | selection after seeing results is a bias |
| 7 | V-B, if used, is the preferred variant and gate 1 is named | ground truth is the whole point |
| 8 | No model path appears in any argument | boundary |
| 9 | Stop conditions are implemented, not merely documented | a stop condition that cannot fire is decoration |
| 10 | The evidence directory is write-once and named before the run | preserves the first result |

---

## 9. Claim ceiling

**Establishes:** a bounded, reproducible specification that fixes the file, the range, the alignment, the accounting and the stop rules before any execution; removes T-15's tail confound; offers a ground-truth variant; and states the limits of cache control rather than implying a purge.

**Does not establish:** that the probe will be conclusive, that it identifies a layer, or that any of it bears on the S1 files. In the best case it separates "the cache holds different bytes" from "the unbuffered path is unstable" — which is one bit more than the investigation has today, and no more.
