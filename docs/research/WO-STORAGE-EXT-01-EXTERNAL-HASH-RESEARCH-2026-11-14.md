# WO-STORAGE-EXT-01 — External Reference Hash Research

**Identity:** `HOWL-EXT-HASH-01` · **Date:** 2026-11-14 · **Mode:** external research only
**Storage access:** none. No local volume, model file, hash, or handle was touched. All local values below are quoted from already-filed run artefacts.
**Outcome in one line:** an authoritative content hash **does** exist for the exact S1-B Gemma artifact, and it equals **one** of the seven observed digests — the composite read from the 19:16:03 run, whose body was unbuffered and whose final 4,064 bytes were buffered. It does not identify the responsible layer.
**Status:** **FROZEN** 2026-11-14. Any later change requires a new identity (`HOWL-EXT-HASH-02`); this text is not edited after freeze.

**Three conclusions, kept separate — do not merge them:**

1. **S1-B** has a repository-verified reference, and one local composite read reproduces it.
2. **S1-A and S1-C** have **no** authoritative publisher reference for the exact converted GGUFs. Third-party references are not equivalent and are not substituted.
3. **Root cause: unresolved.** The reference match anchors the discrepancy; it does not locate the faulty or changing layer.

---

## 1. Method

| Step | What was done |
|---|---|
| 1 | Located candidate publisher repositories by model name, quantisation and naming pattern |
| 2 | Read each repository's public API record (`/api/models/<repo>?blobs=true`), which returns per-file `size` and the **LFS object hash** — git-lfs records `sha256` as the content hash of the stored blob |
| 3 | Matched on **filename + exact byte size**, not filename similarity |
| 4 | Checked for separately published manifests (`SHA256SUMS`, checksum sections in model cards) |
| 5 | Compared results against local values **quoted from filed artefacts only** |

**Why the LFS hash is the strongest available source.** A Git LFS pointer is content-addressed: `sha256` in the pointer is computed over the blob that the server stores and serves. It is machine-generated at upload and cannot drift from the file without producing a new pointer. That is materially stronger than a hand-written checksum table in a README. Its limit is equally clear: it attests to **the blob on the server**, not to what the publisher intended locally, and not to what any local disk holds.

**Retrieval date:** 2026-11-14. Repository head shas are recorded below so a later change is detectable.

**Terminology: `repository-verified`, deliberately not `publisher-verified`.** The pointer authenticates **the bytes it represents** — nothing more. It does not independently establish how the model was originally built, which tooling produced it, whether the publisher's local copy was ever the same bytes, or that the artifact is free of defects. The class label in §3 reflects exactly that scope, and no finding in this report is labelled `publisher-verified`. The file is served as a Xet-backed object; the `sha256` reported is the git-lfs content hash the repository publishes for the pointer, which is the value the Hub displays.

---

## 2. Findings

### 2.1 Primary target — `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M`

| Field | Value |
|---|---|
| Repository | `HauhauCS/Gemma-4-E4B-Uncensored-HauhauCS-Aggressive` |
| Repo head sha | `45b6a334b4bcd1d7f37179df58b3b1d66a184e5d` |
| Repo created / last modified | 2026-04-02 / 2026-04-06 |
| File (exact) | `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` |
| Size | **5,335,285,728 B** — **exactly** the local size |
| LFS sha256 | **`d0027dd3a9128d9323e9f282c8bf010a8526c46477584535991dc1a869b56e96`** |
| Upload commit | `4fb9c3b` (2026-04-02), blob `09fa6f72809d5f304b66e5a4a988560c2975a161` |
| Base model | `google/gemma-4-e4b-it`; abliterated, imatrix quants |
| Classification | **Repository-verified** — publisher-hosted repo, machine-generated content hash, exact filename and exact byte size |
| Served as | Xet-backed object; `sha256` is the git-lfs content hash published for the pointer |
| Independent corroboration | **owner, 2026-11-14:** the public repository page was inspected manually and displays the same SHA-256 and the same 5,335,285,728 B for the named GGUF. Manual inspection, not a second tool |

The repository publishes a full quant ladder (Q2_K_P … Q8_K_P, IQ3_M, IQ4_XS, mmproj), each with its own LFS sha256. Only the Q4_K_M line is quoted here; the others differ, which rules out a family-level checksum being mistaken for a file-level one.

**Local comparison, from `eval/runs/S1B-G-BUFFERED-VS-UNBUFFERED-20261114T191603Z.txt`:**

| Observation | Digest | Equals published hash? |
|---|---|---|
| Buffered read, 19:16:04–19:16:23 | `341fb28ccaa52b57a50b7704476449703f3814338094ff1ab85b7c0d42f67c35` | **no** |
| Unbuffered read, 19:16:03–19:17:12 | `d0027dd3a9128d9323e9f282c8bf010a8526c46477584535991dc1a869b56e96` | **yes — exact** |
| Recorded baseline (`v1.1`, 13:45) | `b5b1913e44aa64d2c12b6159bb782774dd98362528211ac830860595d293ede0` | **no** |
| Unit C, 2026-11-05 | `7096c0dd243da9762247386d476801714971e53e6894b4001710287a522fea1b` | **no** |
| `f46b1049…`, `e6dca6eb…`, `eac8960a…` | (eight-hex prefixes as filed) | **no** |

The run log's own accounting for the matching read: `size 5335285728`, `aligned_len 5335281664`, `tail_len 4064`, `unbuffered_bytes 5335281664`, `buffered_bytes 4064`, `read_calls 1273`, `short_sector_aligned_reads 0`, `identity_stable_after_read: true`. 5,335,281,664 + 4,064 = 5,335,285,728: the composite read covered the whole file with no gap and no short read.

**Local filename** is `…-Q4_K_M(1).gguf`. The `(1)` is the Windows duplicate-name suffix, not part of the artifact name; the association rests on exact name-without-suffix plus **exact byte size** plus the digest match itself.

### 2.2 `Qwen3.5-9B` Q4_K_M

| Field | Value |
|---|---|
| Publisher repo | `Qwen/Qwen3.5-9B` — **safetensors only, no GGUF** |
| Implication | a Q4_K_M GGUF is necessarily a **third-party conversion**; no publisher-originated GGUF hash can exist |
| Converter repo | `unsloth/Qwen3.5-9B-GGUF`, head `3885219b6810b007914f3a7950a8d1b469d598a5`, last modified 2026-03-02 |
| File | `Qwen3.5-9B-Q4_K_M.gguf` — **5,680,522,464 B** |
| sha256 | `03b74727a860a56338e042c4420bb3f04b2fec5734175f4cb9fa853daf52b7e8` |
| Classification | **Third-party reference** — reputable converter, but not the model publisher |

**Collision evidence, and why it matters here:** a separate redistributor publishes a file named `Qwen3.5-9B.Q4_K_M.gguf` (note the dot before `Q4_K_M`) with sha256 `9fa52e37c8293ee02cf277634f8d3a392dc3239b8ccef2315b4c8c9ef223e7fb`, sourced from `mradermacher/Qwen3.5-9B-GGUF`, size ≈5.6 GB. Two plausible-looking Qwen 9B Q4_K_M files, **different names, different hashes, different byte counts**. A filename match alone would have picked the wrong reference.

### 2.3 `LFM2.5-8B-A1B-heretic` Q4_K_M

| Field | Value |
|---|---|
| Repository | `saidutta69/LFM2.5-8B-A1B-heretic`, head `67138c75b2642cef3ce9f9c11e8d4fd77d2ce43d` |
| Created / last modified | 2026-10-04 / 2026-11-01 |
| Base model | `LiquidAI/LFM2.5-8B-A1B`, decensored with **Heretic v1.4.0** abliteration |
| File | `LFM2.5-8B-A1B-heretic-Q4_K_M.gguf` — **5,155,564,448 B** |
| sha256 | `abf60be987673cf8b53c95a3fd0a97b6e6fddd2c81ad75ee2874a0a4299fddf4` |
| Published manifest | `reproduce/SHA256SUMS` exists but lists **only the four safetensors shards** — no GGUF line |
| Classification | **Third-party reference** — the "heretic" derivative is community-produced; Liquid AI publishes `LiquidAI/LFM2.5-8B-A1B-GGUF` but no heretic variant |

**Collision evidence:** at least three other repositories carry the name `LFM2.5-8B-A1B-heretic` (`coder3101`, `Dingdust`, and others) with **safetensors only** and their own ablation settings (Heretic v1.2.0/ARA). A fourth family, `zaakirio/LFM2.5-8B-A1B-Uncensored-GGUF`, is a *different artifact* under a similar name. Name matching without size and hash would be meaningless.

---

## 3. Exact hashes, consolidated

| Artifact | Size (B) | SHA-256 | Provenance class |
|---|---:|---|---|
| `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` | 5,335,285,728 | `d0027dd3a9128d9323e9f282c8bf010a8526c46477584535991dc1a869b56e96` | **Repository-verified** (publisher-hosted) |
| `Qwen3.5-9B-Q4_K_M.gguf` | 5,680,522,464 | `03b74727a860a56338e042c4420bb3f04b2fec5734175f4cb9fa853daf52b7e8` | Third-party reference |
| `LFM2.5-8B-A1B-heretic-Q4_K_M.gguf` | 5,155,564,448 | `abf60be987673cf8b53c95a3fd0a97b6e6fddd2c81ad75ee2874a0a4299fddf4` | Third-party reference |

**Publisher-verified: none of the three.** Not one of these artifacts is published by the model's own publisher as a GGUF — the Gemma case is publisher-hosted by the artifact's *creator* (HauhauCS is both), the other two are derivative conversions.

---

## 4. What the match does and does not mean

### 4.1 It is a reference match at artifact level

At 19:16:03 on 2026-11-14, a read of the file returned bytes whose SHA-256 equals the hash the hosting repository publishes for the exact artifact at the exact byte size. A 256-bit coincidence is not a candidate explanation. So **at least one read path returned the canonical content of the published artifact**, and the file's dimensions on `G:` are the canonical ones.

### 4.2 The evidence is asymmetric, and the asymmetry decides how far it goes

| Result | Strength |
|---|---|
| **Match** to the published hash | **strong positive.** It establishes that this read reproduced the published bytes |
| **Mismatch** | **weak negative.** It cannot distinguish (a) a bad read from (b) genuinely different bytes behind that path — a locally modified, repacked, or restored file under the same name produces exactly the same mismatch signature |

So the six non-matching observations are **not** shown to be bad reads. They are shown to differ from the published content — and for a mismatch the cause remains open.

### 4.3 Consequences, stated at their real strength

| Observation | New status |
|---|---|
| Unbuffered `d0027dd3…` (19:16) | **matches the published artifact.** The strongest single fact in the investigation so far |
| Buffered `341fb28c…` (19:16, same identity, same minute) | differs from the published content; **cause undetermined** |
| Baseline `b5b1913e…` (13:45, 13:50) | differs from the published content. The qualification record's "matched baseline" therefore means **self-consistency only**, not agreement with the publisher |
| Unit C `7096c0dd…` (Nov 5) | differs from the published content; cause undetermined |

### 4.4 The prohibition, restated because it is easy to cross

**A matching external hash does not identify the responsible read path.** Specifically it does **not** establish:

- that the unbuffered path is *correct in general* — it establishes that one unbuffered read, on one file, at one time, returned the published bytes;
- that the buffered path is *defective* — a buffered read returning non-canonical bytes is equally consistent with the page cache holding older bytes, with a filter modifying data, with a genuinely different stored state, or with a read error;
- that `G:` holds the canonical bytes on disk — a cache holding correct bytes and a disk holding correct bytes are indistinguishable from this evidence alone;
- anything about the other six digests' provenance.

The matching read is also **not a pure unbuffered read**: its last 4,064 bytes were fetched buffered (the file is not a whole multiple of the 4,096-byte sector). The composite digest matched, which means both the sector-aligned body and that buffered tail matched. It does not license the sentence "the unbuffered path returned the published file" without that qualifier.

---

## 5. Ambiguity that remains

1. **No local digest exists for Qwen or LFM.** All seven observations concern S1-B (Gemma). The published hashes for the other two are reference values only — nothing local can be compared against them, and their local filenames and sizes were not read (correctly: no storage access).
2. **The other six Gemma digests have no external counterpart.** They are not published anywhere; they are local observations of a file that is published with one hash.
3. **The two dated records still lack a shared file ID** — the November 5 record carries no file identity, so identity across 5 and 14 November remains unestablished (see `S3-METADATA-TEST-RECONCILIATION` §2).
4. **Repacking cannot be excluded for the mismatching reads.** A structurally valid GGUF with different payload bytes hashes differently and would be indistinguishable from a bad read at this evidence level.
5. **The LFS hash attests to the server blob.** If the repository were ever re-uploaded, the pointer would change; the head sha and retrieval date are recorded so that can be checked.

---

## 6. Conclusion

**Can an independent reference resolve the observed digest disagreement?**

**The strongest defensible statement, in full:**

> The 19:16:03 composite read reproduced the exact byte sequence identified by the publisher-hosted Git LFS pointer for the exact Gemma filename and byte size. The other six recorded digests do not match that published blob.

### 6.1 S1-B — reference established

The canonical reference exists and is now anchored to one observation. That value is no longer merely another local reading: it is a read that reproduced the content the hosting repository publishes for the artifact. The investigation has one anchor instead of seven mutually contradicting values.

**The matching read is a composite, not a clean unbuffered validation.** Its body (5,335,281,664 B, 1,273 sector-aligned calls, zero short reads) was unbuffered; its final 4,064 B — the file is not a whole multiple of the 4,096-byte sector — were read buffered. The composite digest matched, so body **and** tail matched. This is evidence for the composite read. It is **not** a validation of the unbuffered path in isolation, and no experiment has yet separated the two paths under control.

### 6.2 S1-A and S1-C — no authoritative reference

Both are community conversions. `Qwen/Qwen3.5-9B` ships safetensors only, so no publisher GGUF hash can exist; the heretic derivative has no Liquid AI counterpart at all. Their hashes in §3 are **third-party references** and must never be promoted, averaged with, or treated as equivalent to the S1-B reference.

### 6.3 Root cause — unresolved

The six mismatches are **real mismatches against the reference**. They are not thereby proof of six faulty reads: the local file could have contained different bytes at those observations, and a mismatch cannot distinguish that from a read fault.

**The buffered path is not convicted.** The matching composite read included a buffered tail, and no controlled experiment isolates the two paths. The root cause remains open across at least four candidates that this evidence cannot separate: stale or inconsistent host memory, storage-path behaviour, a filter or agent in the stack, and actual changes to file contents.

The honest summary: the disagreement is now **anchored**, not **explained**.

---

## 7. Claim ceiling

**Establishes:** that a machine-generated content hash exists for the exact S1-B Gemma artifact at its exact byte size; that one local observation — the 19:16:03 unbuffered read — equals it; that the other six differ from it; that no publisher-verified GGUF hash exists for Qwen3.5-9B or LFM2.5-8B-A1B-heretic; and that filename-similar artifacts with different hashes are abundant for both.

**Does not establish:** which layer produced the disagreeing bytes; whether `G:` holds canonical bytes on disk; whether any of the six mismatches was a read fault or a genuine content difference; whether the file changed on disk at any point; or that the unbuffered access mode is reliable in general.

**Not done, per the work order:** no volume access, no local hash, no model read, no scratch file, no volume handle, no USN query, no S-9a execution, no patch application, no source/manifest/registry change, no digest promoted to authoritative status in any registry.
