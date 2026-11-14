# `S9A-RANGE-PROBE-SPECIFICATION` — v2

**Identity:** `HOWL-S9A-SPEC-2` · **Supersedes:** `HOWL-S9A-SPEC-1` (`5f9b589c…`)
**Status:** **SPECIFICATION ONLY — neither variant is authorised, nothing was executed.**
**What changed from v1:** V-B's interpretation no longer names the cache as the faulty component; V-B gains identity and metadata captures after every read; V-A gains a final identity capture after its last read and states that metadata stability does not prove content stability; the flush claim is corrected to "flush requested successfully".
**Precedence:** replaces `STORAGE-DIAGNOSTIC-STAGING-PLAN` S-9a.

---

## 1. What the probe is for

| | |
|---|---|
| **Question** | Do two access paths return different bytes for the same range of the same file identity, under conditions fixed in advance? |
| **A divergence establishes** | that the two paths disagreed on that range, with range and identity recorded |
| **A divergence does not establish** | which path is wrong, which component caused it, or anything about the S1 model files |
| **Agreement establishes** | **nothing about correctness** — a stale cache and a bad device agree whenever both serve the same wrong data |
| **Out of scope** | cache purging, cache-size changes, driver changes, stress, CHKDSK, surface scans, any `.gguf`/`.safetensors`/`.bin` payload, any write to a model destination |

### 1.1 Cache control is not available, and the output must say so

| Route | Reality |
|---|---|
| `FlushFileBuffers` on a read handle | **fails** — the handle must hold `GENERIC_WRITE` |
| `SetSystemFileCacheSize` / `NtSetSystemInformation` | requires elevation and changes system state; out of scope |
| Dismount / remount the volume | destroys the investigation; forbidden |

So **no invalidation is performed**. The probe measures two access paths in whatever cache state exists, and every output carries `"limits": ["no cache invalidation performed", ...]`. A reader of the JSON must not infer that a purge happened.

---

## 2. Variants

| Variant | Target | Ground truth | Gate |
|---|---|---|---|
| **V-A** | an existing non-model file on `G:` | **none** — content unknown by construction | `G:` payload read |
| **V-B** | a scratch file written with a known pattern | **yes** — the pattern is the oracle | **gate 1** (`G:` disposable file) |

V-B is the more informative experiment and should be preferred when gate 1 opens. Neither reads a model.

---

## 3. V-A — existing non-model file

### 3.1 Selection rule (evaluated before the run, recorded)

All of: on `G:`, outside `G:\AI GGUF\GGUF`; extension not `.gguf .safetensors .bin .ckpt .pt .pth .onnx`; size 64–512 MiB; size and mtime unchanged across two metadata reads ≥ 60 s apart; opened with `FILE_SHARE_READ | FILE_SHARE_WRITE` (never exclusively). If nothing satisfies the rule, **stop** — the rule is not relaxed during the run.

**Limitation, stated because the rule relies on it:** two metadata reads 60 s apart do **not** prove content stability. A writer can change content while leaving size and mtime as they were — `S3-METADATA-TEST-RECONCILIATION` shows the ingestion tooling itself can restore timestamps. Selection stability is a screening convenience, not a guarantee.

### 3.2 Range

Offset a fixed 4096-aligned value recorded before the run (e.g. `1 GiB`); length `64 MiB`; block `4 MiB`; buffer 4096-aligned. `G:` is logical 512 / physical 4096, so 4096 satisfies both.

### 3.3 Procedure

| Step | Action | I/O |
|---|---|---|
| A-1 | capture `identity_0` (`size`, `mtime_ns`, `volume_serial`, `file_id`, `in_use`) | metadata |
| A-2 | buffered read → `sha256_buffered` | 64 MiB |
| A-3 | capture `identity_1`; abort on any change | metadata |
| A-4 | unbuffered read → `sha256_unbuffered` | 64 MiB |
| A-5 | capture `identity_2`; abort on any change | metadata |
| A-6 | repeat unbuffered read → `sha256_unbuffered_2` | 64 MiB |
| A-7 | **capture `identity_3`; abort on any change** | metadata |

Four identity captures, one after every read. Total payload **192 MiB**. The repeat at A-6 exists because a single differing pair cannot distinguish a stable alternative content from a one-off bad read.

### 3.4 Tail handling

The range is 16 whole 4 MiB blocks, so no tail exists. If the range is ever shortened, hash body and tail **separately** and report both — never one combined value. This is the T-15 confound, removed.

---

## 4. V-B — known-pattern scratch file (gate 1)

| Step | Action | Identity / metadata |
|---|---|---|
| B-1 | create `G:\AI GGUF\_probe\<random>.bin`, 64 MiB, written with `NO_BUFFERING` from a recorded-seed PRNG stream | — |
| B-2 | `FlushFileBuffers` on the write handle | **"flush requested successfully"** — not proof that every downstream layer committed |
| B-3 | close write handle; capture `identity_0` | baseline |
| B-4 | unbuffered read → compare with expected pattern digest | capture `identity_1`; abort on change |
| B-5 | buffered read → `sha256_buffered` | capture `identity_2`; abort on change |
| B-6 | unbuffered read → `sha256_unbuffered` | capture `identity_3`; abort on change |

Cleanup is a write and must be recorded; the file is left in place until the evidence is filed.

**On durability:** step B-2 issues the documented flush for the write handle and the call returned success. It does **not** establish that the device cache, controller, or medium physically committed the data honestly. The reading is "flush requested successfully"; the oracle is the expected pattern, and it is only as trustworthy as the write path that produced it.

---

## 5. Interpretation — byte-level facts, not components

### 5.1 V-A

| Outcome | Reading |
|---|---|
| buffered ≠ unbuffered, both unbuffered runs equal | the two paths returned different bytes for the range; the unbuffered path is reproducible on this range |
| the two unbuffered runs differ | the unbuffered path is not reproducible on this range |
| all three equal | the paths agreed on this range, this time — **and nothing more** |

### 5.2 V-B

| Outcome | Reading |
|---|---|
| B-4 matches; B-5 differs | **the buffered path returned bytes that differ from the known pattern** |
| B-4 matches; B-6 differs | **the later unbuffered read returned bytes that differ from the known pattern** |
| both later reads differ | both paths returned incorrect bytes, **or** the file changed after the verification read |
| all match | the tested paths agreed with the expected pattern during this experiment |

**"Cached stream" is reserved as a hypothesis, never a finding.** A divergence in B-5 does not locate the fault in the Windows system file cache: the paths still share the memory manager, the filter stack, the class and port drivers, the controller and the device. The pattern gives ground truth about *returned bytes*, not a location for the defect.

### 5.3 The word "wrong"

Used above only in the sense "differs from the pattern that was written". It is not a claim about which component is defective.

---

## 6. Output schema

```json
{
  "spec": "HOWL-S9A-SPEC-2",
  "variant": "V-B",
  "started_utc": "...", "ended_utc": "...",
  "path": "...",
  "range": {"offset": 1073741824, "length": 67108864, "block": 4194304, "sector": 4096},
  "identity": {
    "i0": {"size": 0, "mtime_ns": 0, "volume_serial": "0x...", "file_id": "...", "in_use": false},
    "i1": {"...": "..."}, "i2": {"...": "..."}, "i3": {"...": "..."}
  },
  "digests": {"expected": "...", "buffered": "...", "unbuffered": "...", "unbuffered_2": "..."},
  "flush": {"requested": true, "call_returned_success": true},
  "io": {"read_ops_delta": 0, "read_bytes_delta": 0, "elapsed_s": 0},
  "verdict": "DIVERGED_STABLE | DIVERGED_UNSTABLE | AGREED | ABORTED",
  "limits": [
    "no cache invalidation performed",
    "flush success does not establish physical persistence",
    "divergence does not identify a component"
  ]
}
```

`DIVERGED_STABLE` = buffered differs from both unbuffered runs, which agree. `DIVERGED_UNSTABLE` = the two unbuffered runs disagree. `AGREED` = all equal.

---

## 7. Stop conditions

Stop, preserve output, do not retry if: any identity capture changes; any read errors, short-reads, or returns `ERROR_CRC`/`ERROR_DEVICE_HARDWARE_ERROR`; the file becomes `in_use`; the unbuffered call fails with `ERROR_INVALID_PARAMETER` (a probe defect, not a device one); wall clock exceeds 10 minutes; total payload exceeds 256 MiB.

A stop is a result. Retrying against the same range is a second experiment requiring its own authorisation.

---

## 8. Review checklist before authorisation

| # | Item |
|---:|---|
| 1 | Exact command line in the work order, not improvised |
| 2 | Reader's block/sector/buffer alignment reviewed line by line |
| 3 | Identity captured after **every** read, in both variants |
| 4 | Range block-exact, or body and tail hashed separately |
| 5 | "No cache invalidation performed" present in the JSON, not only in prose |
| 6 | File selected by §3.1 rule before the run |
| 7 | V-B preferred; gate 1 named in the authorisation |
| 8 | No model path in any argument |
| 9 | Stop conditions implemented, not merely documented |
| 10 | Evidence directory write-once and named before the run |

---

## 9. Claim ceiling

**Establishes:** a bounded, reproducible specification that fixes file, range, alignment, accounting and stop rules in advance; captures identity after every read in both variants; states the limits of cache control and of flushing; and interprets results as byte-level facts rather than component attributions.

**Does not establish:** that the probe will be conclusive, that it identifies a layer, or that it bears on the S1 files. In the best case it separates "the buffered path returned bytes that differ from a known pattern" from "the unbuffered path is unstable" — one bit more than the investigation has today, and no more.
