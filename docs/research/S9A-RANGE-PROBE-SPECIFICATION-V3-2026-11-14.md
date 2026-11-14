# `S9A-RANGE-PROBE-SPECIFICATION` — v3

**Identity:** `HOWL-S9A-SPEC-3.1` · **Supersedes v3** (`38c9efbc…`), v2 (`c390b413…`), v1 (`5f9b589c…`)
**Status:** **SPECIFICATION ONLY — neither variant authorised, nothing executed.**
**Fixed in v3:** V-B's range was impossible (a 1 GiB offset into a 64 MiB file); V-B had no complete verdict mapping; the content-immutability limitation was stated for V-A only.
**Fixed in 3.1, after review:** V-A's range was impossible for every file its selection rule admits (1 GiB offset against a 512 MiB ceiling) — offset is now `0`; `b_eq_u2` is a populated boolean rather than `null`; the 256 MiB ceiling is defined as **total payload, writes included**.

> **Correction log.** v3 repeated v2's defect in the other variant. v2 put a 1 GiB offset on a 64 MiB scratch file; v3 fixed V-B and left V-A reading 1 GiB + 64 MiB from a file the rule caps at 512 MiB. Both are the same failure: a range specified by analogy instead of checked against the declared size. The range check is now an asserted precondition in both variants, not a reviewer's responsibility.
**Precedence:** replaces `STORAGE-DIAGNOSTIC-STAGING-PLAN` S-9a.

---

## 1. What the probe is for

| | |
|---|---|
| **Question** | Do two access paths return different bytes for the same range of the same file identity, under conditions fixed in advance? |
| **A divergence establishes** | that the paths disagreed on that range, with range and identity recorded |
| **A divergence does not establish** | which path is wrong, which component caused it, or anything about the S1 model files |
| **Agreement establishes** | **nothing about correctness** — a stale cache and a bad device agree whenever both serve the same wrong data |
| **Out of scope** | cache purging, cache-size changes, driver changes, stress, CHKDSK, surface scans, any `.gguf`/`.safetensors`/`.bin` payload, any write to a model destination |

### 1.1 Cache control is not available, and the output says so

`FlushFileBuffers` on a read handle fails (`GENERIC_WRITE` required); global cache control needs elevation; dismounting is destructive. **No invalidation is performed**, and every output carries `"no cache invalidation performed"` in its `limits`.

---

## 2. Variants and their ranges

| Variant | Target | Size | Range | Ground truth | Gate |
|---|---|---|---:|---|---|
| **V-A** | existing non-model file on `G:` | 64–512 MiB | offset **0**, length **64 MiB** | none | `G:` payload read |
| **V-B** | scratch file `G:\AI GGUF\_probe\<random>.bin` | **64 MiB** | offset **0**, length **64 MiB** | the written pattern | **gate 1** |

**Range errors, corrected.** v2's schema gave V-B the 1 GiB offset intended for V-A, which starts beyond the end of a 64 MiB scratch file. v3 corrected V-B and left V-A at 1 GiB — impossible from the other direction, since 1 GiB + 64 MiB = 1088 MiB exceeds the 512 MiB ceiling its own selection rule imposes. **Both variants now read offset 0, length 64 MiB**, and the selection rule's upper bound is unchanged because nothing about a 1 GiB offset was ever justified: the purpose is a fixed, aligned, whole-block range, not a position deep in the file.

Every run asserts `offset + length <= size` as a precondition (`range_check` in the output). The block size is 4 MiB and 64 MiB is 16 whole blocks, so **no tail exists** in either variant.

---

## 3. V-A — existing non-model file

### 3.1 Selection rule (evaluated and recorded before the run)

All of: on `G:`, outside `G:\AI GGUF\GGUF`; extension not `.gguf .safetensors .bin .ckpt .pt .pth .onnx`; size 64–512 MiB; size and mtime unchanged across two metadata reads ≥ 60 s apart; opened with `FILE_SHARE_READ | FILE_SHARE_WRITE`, never exclusively. If nothing satisfies the rule, **stop** — it is not relaxed mid-run.

### 3.2 Procedure

| Step | Action | I/O |
|---|---|---|
| A-0 | assert `offset + length <= size`; abort if not | metadata |
| A-1 | capture `identity_0` (`size`, `mtime_ns`, `volume_serial`, `file_id`, `in_use`) | metadata |
| A-2 | buffered read → `sha256_buffered` | 64 MiB |
| A-3 | capture `identity_1`; abort on change | metadata |
| A-4 | unbuffered read → `sha256_unbuffered` | 64 MiB |
| A-5 | capture `identity_2`; abort on change | metadata |
| A-6 | repeat unbuffered read → `sha256_unbuffered_2` | 64 MiB |
| A-7 | capture `identity_3`; abort on change | metadata |

Payload **192 MiB**. The repeat at A-6 separates a stable alternative content from a one-off bad read.

### 3.3 V-A verdicts

| Outcome | Verdict | Reading |
|---|---|---|
| buffered ≠ both unbuffered; unbuffered runs equal | `DIVERGED_STABLE` | the paths returned different bytes for the range; the unbuffered path is reproducible |
| the two unbuffered runs differ | `DIVERGED_UNSTABLE` | the unbuffered path is not reproducible on this range |
| all three equal | `AGREED` | the paths agreed on this range, this time — and nothing more |
| a precondition failed | `ABORTED` | not a negative result |

---

## 4. V-B — known-pattern scratch file (gate 1)

| Step | Action | Capture |
|---|---|---|
| B-1 | create the 64 MiB file with `NO_BUFFERING` from a recorded-seed PRNG stream | — |
| B-2 | `FlushFileBuffers` on the write handle | `flush: {requested: true, call_returned_success: true}` |
| B-3 | close write handle; capture `identity_0` | baseline |
| B-4 | unbuffered read of **offset 0, length 64 MiB** → `u1`; compare with expected `E` | `identity_1`; abort on change |
| B-5 | buffered read of the same range → `b` | `identity_2`; abort on change |
| B-6 | unbuffered read again → `u2` | `identity_3`; abort on change |

**Durability, stated precisely:** B-2 issues the documented flush and the call returned success. That does **not** establish that the device cache, controller or medium committed the data honestly. The oracle is the pattern, and it is only as trustworthy as the write path that produced it.

**Identity stability is not content immutability.** The four captures catch ordinary changes. They cannot exclude a write that preserves size and both timestamps — the ingestion tooling in this repository restores timestamps by design, so this is not a hypothetical.

### 4.1 V-B verdict mapping — complete over all combinations

Let `E` be the expected pattern digest, and `u1`, `b`, `u2` the three measured digests.

| # | `u1` | `b` | `u2` | Verdict | Reading |
|---:|---|---|---|---|---|
| 1 | = E | = E | = E | **`VERIFIED_AGREE`** | all three reads returned the pattern |
| 2 | = E | ≠ E | = E | **`BUFFERED_OFF_PATTERN`** | the buffered read returned bytes differing from the pattern; both unbuffered reads matched |
| 3 | = E | = E | ≠ E | **`UNBUFFERED_OFF_PATTERN`** | the later unbuffered read differed while the earlier one and the buffered one matched |
| 4 | = E | ≠ E | ≠ E | **`BOTH_LATER_OFF_PATTERN`** | both later reads differed from the pattern; record whether `b == u2` — if they agree with each other, the content differs from what was verified, and neither path is shown to be at fault |
| 5 | ≠ E | — | — | **`WRITE_UNVERIFIED`** | the oracle check failed. Stop. Nothing about read paths can be concluded from this run |
| 6 | — | — | — | **`ABORTED`** | a precondition (identity, error, timeout, byte ceiling) failed |

Rows 1–4 partition every case where the oracle held: each of `b` and `u2` is either equal or not equal to `E`, giving four combinations. Row 5 covers everything else. There is no unclassified outcome, and `WRITE_UNVERIFIED` is a **stop**, not a result.

**No row names a component.** Row 2 says the buffered path returned bytes that differ from the pattern; it does not say the system file cache is at fault, because the paths still share the memory manager, filter stack, class and port drivers, controller and device. "Cached stream" is reserved as a hypothesis.

---

## 5. Output schema

```json
{
  "spec": "HOWL-S9A-SPEC-3.1",
  "variant": "V-B",
  "path": "G:\\AI GGUF\\_probe\\<id>.bin",
  "size": 67108864,
  "range": {"offset": 0, "length": 67108864, "block": 4194304, "sector": 4096},
  "range_check": {"offset_plus_length_le_size": true},
  "identity": {"i0": {"...": "..."}, "i1": {"...": "..."}, "i2": {"...": "..."}, "i3": {"...": "..."}},
  "digests": {"expected": "...", "u1": "...", "buffered": "...", "u2": "..."},
  "relations": {"b_eq_u2": true},
  "flush": {"requested": true, "call_returned_success": true},
  "io": {"read_ops_delta": 0, "read_bytes_delta": 0, "elapsed_s": 0},
  "verdict": "VERIFIED_AGREE | BUFFERED_OFF_PATTERN | UNBUFFERED_OFF_PATTERN | BOTH_LATER_OFF_PATTERN | WRITE_UNVERIFIED | ABORTED",
  "limits": [
    "no cache invalidation performed",
    "flush success does not establish physical persistence",
    "identity stability does not imply content immutability",
    "divergence does not identify a component"
  ]
}

`spec` records the **exact revision** of this document, not a family version: evidence produced under a later correction must be attributable to the text that governed it, so the field is bumped whenever the body changes. This document emits `HOWL-S9A-SPEC-3.1`.

**`b_eq_u2` is always populated** in V-B whenever both `b` and `u2` exist, as `true` or `false` — never `null`. It is required evidence for `BOTH_LATER_OFF_PATTERN`: if the two later reads agree with each other and differ from the verified pattern, the content differs from what was verified, and that fact neither blames nor exonerates either path. V-A emits `null`, because it has no expected digest and the relation carries no meaning there.
```

---

## 6. Stop conditions

Stop, preserve output, do not retry if: any identity capture changes; any read errors, short-reads, or returns `ERROR_CRC` / `ERROR_DEVICE_HARDWARE_ERROR`; the file becomes `in_use`; an unbuffered call fails with `ERROR_INVALID_PARAMETER` (probe defect, not device); wall clock exceeds 10 minutes; total payload would exceed 256 MiB; or `u1 ≠ E` in V-B.

**The 256 MiB ceiling counts every payload byte in both directions.** V-B writes 64 MiB and reads 3 × 64 MiB, so a full run sits **exactly at** the ceiling: 256 MiB total, 64 written + 192 read. The check is `would_exceed`, not `reaches`, so a complete V-B run is permitted and any additional read is not. Implementations that count reads only would allow 320 MiB of actual I/O and are wrong. V-A is read-only at 192 MiB, and V-B's write count is fixed at 64 MiB, so no run can approach the ceiling without a second pass — which §6 does not authorise.

A stop is a result. Retrying the same range is a second experiment needing its own authorisation.

---

## 7. Review checklist before authorisation

| # | Item |
|---:|---|
| 1 | Exact command line in the work order |
| 2 | Reader's block/sector/buffer alignment reviewed line by line |
| 3 | Identity captured after every read, both variants |
| 4 | Range validated against file size (`offset + length ≤ size`) before the first read |
| 5 | "No cache invalidation performed" present in the JSON |
| 6 | File (V-A) or scratch path (V-B) fixed before the run |
| 7 | Gate named: `G:` payload read for V-A; gate 1 for V-B |
| 8 | No model path in any argument |
| 9 | Stop conditions implemented; `WRITE_UNVERIFIED` halts |
| 10 | Evidence directory write-once and named before the run |

---

## 8. Claim ceiling

**Establishes:** a specification whose range is possible for every file each variant admits, enforced by an asserted precondition; a complete six-verdict mapping for V-B; a populated `b_eq_u2`; a defined payload ceiling that includes writes; identity captures after every read; an explicit statement that identity stability does not imply content immutability; and limits printed in the output.

**Does not establish:** that the probe will be conclusive, that it identifies a layer, or that it bears on the S1 files. In the best case it separates "the buffered read returned bytes differing from a known pattern" from "the unbuffered path is unstable" — one bit, no more.
