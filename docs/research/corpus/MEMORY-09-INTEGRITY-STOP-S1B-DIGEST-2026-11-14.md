# MEMORY-09 — Integrity Stop: `S1-B` file digest is not stable

**Status:** **STOPPED.** `S1-B` requalification **refused** by the digest gate. `S1-C` not attempted. Server stopped.
**Date:** 2026-11-14 · **Trigger:** the qualification gate, exactly as designed.
**Evidence:** `eval/runs/DIGEST-MISMATCH-S1-B-2026-11-14T140558Z/evidence.json`

---

## 1. What happened

`S1-A` requalified cleanly at the S1-v1.1 profile. Then `S1-B` was loaded and the recorder refused it:

```
REFUSED: S1-B file does not match the frozen baseline
  expected b5b1913e…ede0 (5,335,285,728 B)
  actual   f46b1049…9dcd2 (5,335,285,728 B)
```

**Same byte count. Different content.** Then the file produced *different digests on consecutive reads*.

## 2. The digest timeline

| Time (local) | Condition | SHA-256 |
|---|---|---|
| 13:45 | baseline creation, server idle | `b5b1913e…ede0` |
| 13:50 | `S1-B` qualification v1 | `b5b1913e…ede0` (matched) |
| 14:03 | requalification, **server running** | `f46b1049…9dcd2` |
| 14:03 | second read, server running | `e6dca6eb…51e23` |
| 14:04 | server **stopped**, read | `e6dca6eb…51e23` |
| 14:04 | server **stopped**, read again | `e6dca6eb…51e23` |
| 14:05 | server **stopped**, read | `eac8960a…2e3d0` |

**Four distinct digests from one file**, with size constant at 5,335,285,728 B and **mtime unchanged since 2026-09-04 20:53:48**.

## 3. The control

`LFM2.5-8B-A1B-heretic.Q4_K_M.gguf`, same directory, same drive, hashed during the same window:

```
f820ef3f…c5c7e  ==  baseline f820ef3f…c5c7e   ✅
```

So this is **not** a broken hashing path, a wrong file, or a transient tooling error. One file on that volume is unstable; another on the same volume is not.

## 4. What this rules out, and what it leaves

**Ruled out:** harness bug (control passes), wrong path (same file, same size, same mtime), a deliberate edit that updated mtime (mtime untouched), my own writes (nothing in this session writes to that directory).

**Left, and outside the harness's authority to resolve:**

1. **A background writer** — a sync client, downloader, or application rewriting or re-caching the file. SillyTavern or a mod manager are plausible candidates on this machine; so is a cloud-sync folder if `G:` is indexed.
2. **Storage instability** — a failing sector, cable/controller fault, or external-drive power management returning inconsistent data on repeated reads. A file whose bytes differ between reads is the signature of this, not of an edit.

**A file that returns different bytes on consecutive reads cannot be used as experimental evidence.** Which bytes were loaded at any moment is unknowable, so a digest in a record would be a claim about a file that no longer describes what ran.

## 5. What was not done

- **No substitution.** The reserve candidate was not used to fill the slot.
- **No baseline edit.** `s1-readiness-baseline*.json` untouched; the baseline digest is still the truth the gate compares against.
- **No retry loop.** One refusal, stopped immediately.
- **No `S1-C` attempt.** Sequential order means B's failure ends the batch, even though C is independent.
- **No quarantine copy.** Copying an unstable file would produce a second artifact of unknown identity; that needs a decision, not a reflex.

## 6. Status of the three slots

| Slot | At S1-v1.1 profile | Evidence |
|---|---|---|
| **S1-A** | **QUALIFIED** | 20/33 layers, 20,480 ctx, 5,379 MiB used (65.7%), recovery Δ −55 MiB, `generations: 0` |
| **S1-B** | **BLOCKED — file identity unstable** | refused before any record was written |
| **S1-C** | **NOT ATTEMPTED** | batch stopped at B |

`S1-A`'s estimate (~5,295 MiB) was within 84 MiB of the measured 5,379 MiB — the KV-and-layer model was sound, which makes the B failure more clearly a storage/identity problem and not a configuration one.

## 7. What the owner needs to decide

1. **Identify writers** — check whether anything (SillyTavern, a launcher, a sync client, a download tool) has `G:\AI GGUF` open or indexed. This is a five-minute check with the right tool and cannot be done from inside the harness.
2. **Check the volume** — a read-stability test or filesystem scan of `G:` before trusting any file on it.
3. **Re-obtain the file** from its original source, verify the digest **three times**, and only then requalify.
4. **Consider a hashed copy** on a volume already proven stable, with the original left untouched, so future qualifications cannot be invalidated by whatever is rewriting this one.

Until one of those happens, `S1-B` cannot be qualified and **S1 cannot run**: two of three slots would be unverified, and the population rule forbids substituting the third.

## 8. The gate worked

A qualification gate that never refuses anything is decoration. This one refused a file whose digest had been recorded as correct an hour earlier — before a single prompt was sent, and before that file's bytes could contaminate an experimental record. The refusal is the most useful thing the harness has produced today.

---

## 9. Correction (owner review, 2026-11-14) — what four digests do and do not prove

§4 above ordered the causes wrongly and is corrected here. **Four different hashes establish that the observed contents were inconsistent. They do not establish failing storage hardware.**

| Explanation | Compatible with constant size and unchanged mtime? |
|---|---|
| A process rewriting the file while preserving timestamps | **yes** |
| Storage path — disk, cable, controller, cache | **yes** |
| System memory instability corrupting reads | **yes** |
| Filesystem-level problem (NTFS metadata/cache) | **yes** |

The unchanged size and mtime **do not discriminate** between these. Ranking them, as §4 did, was inference presented too close to fact.

### 9.1 The control's limited meaning

`LFM2.5-8B-A1B` hashing consistently shows that **one other file** on `G:` read consistently during that window. It does **not** certify the volume: a localised bad region, or a writer targeting one path, produces exactly this pattern.

### 9.2 Preserve-first protocol (followed)

| Rule | Status |
|---|---|
| Do not hash the subject again (each read ≈ 5.3 GB and does not narrow the cause) | **followed** — no read after 14:05 |
| No CHKDSK /R, surface scan, benchmark or broad drive test | **followed** — none run |
| No in-place repair of the original | **followed** |
| Check for an active writer, read-only and narrowly scoped | **done** — see 9.3 |
| Preserve all four digests, timestamps, size, logs and sequence | **done** — `timeline.json` (7 observations, 4 distinct digests) and `diagnostics.json` in the evidence directory |

### 9.3 What the narrow check found

- **No non-shell process** references the model store in its command line. *Caveat:* command lines do not reveal open handles; a writer that never names the path would not appear. **Sysinternals `handle.exe` was not available**, so true handle enumeration was **not performed** — that gap is stated, not papered over.
- **OneDrive** and **OneDrive.Sync.Service** are running.
- Volume `G:` — label **"Second backup"**, NTFS, Fixed, `HealthStatus: Healthy`.
- Physical disk 2 — `WDC WD10SPCX-75KHST0`, SATA, `Healthy`, `Online`.

**"Healthy" is firmware-reported metadata, not a read-integrity test.** It means no condition was flagged; it does not mean the volume is proven stable.

### 9.4 Open questions this check could not answer

1. Does the label "Second backup" mean backup software writes here periodically, timestamps preserved?
2. Is `G:` inside any OneDrive-synced or indexed path?
3. Would a handle enumeration show a holder of this specific file?

Until one of these is answered, **no further model loads** and **no manifest change**. `S1-A` remains qualified at the v1.1 profile; `S1-B` blocked; `S1-C` unqualified at v1.1.

---

## 10. What was hashed — source-level confirmation

A reasonable hypothesis was raised: that these digests might be of the model's **generated output**, in which case differing values would be ordinary sampling nondeterminism rather than an integrity problem. Checked against the code and the commands actually run:

| Question | Answer | Evidence |
|---|---|---|
| What did the shell hash? | `sha256sum "$f"` where `$f` is the **`.gguf` path** | the only path bound to `$f` in every command that produced a digest |
| What does the recorder hash? | `file_digest(path)` → `path.open("rb")`, streamed in 1 MiB blocks | `eval/harness/qualify.py` |
| Was any response text ever fed to it? | **No** | `file_digest` accepts a path only; no call site passes a string of model output |
| Where do response hashes live? | `response_hash`, computed over canonical JSON of the **body**, in the L3 record | `eval/harness/l3.py` — a different field, never `model.sha256` |
| Is the tool itself unstable? | No — hashing a text file twice gives identical digests instantly | demo run, both `8b880b25…` |

Three independent signals confirm these were file reads:

1. The first two observations **matched** the manifest digest computed from the same path at baseline. Generated text could not have matched a file's digest.
2. Every observation carries the file's byte count (5,335,285,728 B), which `file_digest` returns from the bytes it actually read.
3. Each digest required reading ~5.3 GB and took tens of seconds; hashing a short string is instantaneous.

**So the nondeterminism explanation is excluded**, and the anomaly stands: the bytes returned from one path on `G:` were not constant across reads. Which of the four contents, if any, is the original remains unknown — and the loaded model in the 14:03 cycle cannot be tied to any recorded digest.
