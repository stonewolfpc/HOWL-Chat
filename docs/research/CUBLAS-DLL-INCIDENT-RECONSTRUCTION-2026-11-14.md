# `CUBLAS-DLL-INCIDENT-RECONSTRUCTION`

**Identity:** `HOWL-CUBLAS-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-03` Track 2
**Status:** reconstruction only. **No file read, no digest computed, no volume accessed.**
**Artefacts:** `eval/runs/H-MIGRATION-VERIFY-2026-11-14.json`, `H-MIGRATION-VERIFY-CORRECTION.json`, `H-MIGRATION-RECORD-2026-11-14.json`, `H-RESTORE-VERIFY-2026-11-14.json`, `h-migration-verify.log`, `h-restore-verify.log`.

**Correction carried into this document.** The prior provenance matrix described this incident as a copy *from `G:`* and as a second instance of "reads of `G:`". **Both statements are wrong.** The source was `H:\AI apps`, before that volume was reformatted; `G:` was not involved at all. See `STORAGE-INVESTIGATION-CLAIM-CORRECTIONS` C-02.

---

## 1. Recorded sequence

| # | Step | Recorded detail |
|---:|---|---|
| 1 | **Source state** | `H:` was exFAT, `HealthStatus: Warning`, `OperationalStatus: Full Repair Needed`, disk SMART status **OK** |
| 2 | **Copy** | `H:\AI apps` → `D:\H-MIGRATION-2026-11-14\AI apps`; robocopy reported 25,744 files, **0 failed**; whole tree 78,259 files / 13,385,936,341 B |
| 3 | **Verification** | 25,744 files hash-checked, **1 mismatch**, 0 size mismatches, 0 read errors, 0 missing |
| 4 | **The file** | `llama.cpp-prism-b9596-cuda12\runtime\cublas64_12.dll`, 100,033,536 B on both sides, **same mtime on both sides** |
| 5 | **The two digests** | source `e40202fe4223c1cd2d2dce7beec59e1ed61c7801bd827309183be9b50e358f4c` · destination `714858cdb63540c6ac434c5c9bc062c15a2d6b17281f332699082319377a8dcd` |
| 6 | **Stability at detection** | **both sides stable over three reads** |
| 7 | **Re-copy** | destination re-copied from the same source path |
| 8 | **After re-copy** | `after_recopy_matches_source: true`; destination stable over **two further reads**; mismatch count 0 |
| 9 | **Then** | `H:` was quick-formatted to NTFS; the staging tree was restored `D:`→`H:`; restore verification hashed 78,259 files with **0 mismatches** |
| 10 | **Conclusion filed** | the migration record calls it "a silent copy error: same size, same mtime, different bytes, no robocopy failure reported … a single copy pass cannot certify a tree on this machine" |

**Not recorded:** the source's digest before the copy began (only after), any file id or volume serial for either side, whether the source was read by anything else at the time, and the exact clock times of steps 2–8.

---

## 2. What "re-copy fixed it" establishes

| Statement | Status |
|---|---|
| The destination file, as first written, did **not** have the same content as the source **as measured after the copy** | **established** — a second copy from the same path produced a file matching the source digest |
| The first destination file was stable across three reads — it was not a transient, still-settling write | **established** |
| The source was stable across three reads at the same time | **established** |
| One copy pass serialised bytes that another copy pass from the same path did not | **established** — this is the incident's whole content |
| The source file was correct on disk | **not established** — no external reference exists for this DLL. The only "truth" used was a digest of the same path |
| The destination was wrong *and the source right* | **not established** — the re-copy picked the source as the reference by definition |
| A cache, memory fault, or device fault caused it | **not established** — no mechanism is identified by the record |
| `G:` was involved | **false** — the source was `H:` |

### 2.1 The distinction the work order asks for

| Candidate reading | Supported? |
|---|---|
| **Successful replacement of an incorrect destination** | **yes**, in the operational sense: the tree is now internally consistent and restore-verified |
| **Successful refresh of stale or inconsistent cached data** | **not supported as a finding** — plausible mechanism, indistinguishable from the others in this record (see §3) |
| **A different source file being copied** | **no** — same path, same size, same mtime, and the re-copy matched the source digest |
| **A path or identity change** | **no evidence either way** — no file id was captured, so identity was never *tested*, only assumed from path + size + mtime |
| **Unexplained result that cannot yet distinguish these** | **this is the correct classification** |

The last row is the honest one: the incident is a **confirmed non-reproducible copy discrepancy with an unresolved mechanism**, not a confirmed corruption with a known cause.

---

## 3. Candidate mechanisms, with the assumptions each needs

| # | Mechanism | Assumptions required | Fits what is recorded? | Confidence |
|---:|---|---|---|---|
| M-1 | **Bad read from the source volume** — `H:` was exFAT reporting `Full Repair Needed`; one read returned wrong bytes | (a) the filesystem state produced at least one bad read; (b) it was localised to one file in 78,259 | **yes**, and it is the only *documented* anomaly in the chain | **plausible, unproven** |
| M-2 | **Bad write to the destination** — `D:` received wrong bytes despite robocopy reporting success | (a) a write-path fault on `D:`; (b) silent | yes | plausible, unproven |
| M-3 | **Cache staleness in the copy path** — source pages stale relative to the medium at copy time, refreshed before verification | (a) something dirtied or bypassed the cache; (b) timing such that verification read fresh | yes | plausible, unproven |
| M-4 | **Filter or AV interference** — `WdFilter`, `bfs`, `FileInfo`, `Wof`, `gameflt` are all attached to `H:` (filed filter inventory) | (a) a filter altered or substituted content; (b) preserved mtime; (c) left no detection | possible; no product other than Microsoft's own filters is present, and none of them is known to substitute file content | **weak** |
| M-5 | **Host memory fault** corrupting the copy buffer | (a) a fault; (b) striking one buffer among 78,259 files | yes | plausible, unproven |
| M-6 | **Source changed between copy and verification** | (a) a writer on `H:`; (b) not detected by three stable reads afterwards | possible but contrived | weak |
| M-7 | **The re-copy, not the first copy, was the anomalous one** — i.e. the destination now holds a *second* wrong stream | (a) two independent bad reads; (b) both stable | yes — and it is why "re-copy fixed it" cannot be treated as proof | **not excluded** |

**M-7 deserves emphasis.** The record's own conclusion treats the *first* copy as wrong. Nothing in it establishes that, because the source digest was measured around the same period and is itself the reference. If the source read that produced `e40202fe…` was the corrupted one, the "corrected" destination now holds a faithfully copied bad read. Only a publisher hash could settle it.

---

## 4. What this incident does and does not contribute to the S1 question

| Contribution | Status |
|---|---|
| Demonstrates that a copy on this machine produced a file with identical size, identical mtime and different content | **established** — the same *signature* as the S1 files |
| Demonstrates that the signature can arise **outside** any model path, on a different volume, in a different file type | **established** |
| Demonstrates that one copy pass cannot certify a tree here | **established**, and stated by the record itself |
| Implicates `G:` | **no** — different volume, not involved |
| Implicates Windows 11's storage stack generically | **no** — one incident, no mechanism |
| Shares a mechanism with S1-B | **unknown** — same signature, no shared evidence. A shared signature is not a shared cause |

The S1-B case has something this one lacks: a verified **layer divergence** (buffered vs unbuffered in one minute). The DLL case has something S1-B lacks: a **second copy** that produced different bytes from the same path. Neither record connects them, and no document should present the pair as one mechanism until one does.

---

## 5. Claim ceiling

**Establishes:** the incident's exact sequence, its correct source volume, the four things "re-copy fixed it" does establish, the four it does not, and seven candidate mechanisms with the assumptions each needs — including the inverted reading in which the *re-copy* carries the bad bytes.

**Does not establish:** that the DLL was corrupt, which side of the copy was wrong, or that any mechanism is shared with the S1 model files. It is an unexplained copy discrepancy on a volume that was reporting `Full Repair Needed`, and it is recorded as such.
