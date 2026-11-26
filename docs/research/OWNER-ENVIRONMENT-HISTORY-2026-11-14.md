# Owner Environment History — Reconciliation

**Identity:** `HOWL-ENV-HIST-1.1` · **Date:** 2026-11-14 · **Mode:** documentation only
**1.1 adds:** owner answers to the timing, memory-configuration and bugcheck questions, and a mechanism-admissibility table that uses them. §2.3's bound is now resolved rather than open.
**Evidence class:** **owner statement.** Everything in §1 is reported by the owner from memory, logs and an analysis performed in another session. None of it was independently verified here: no event log, minidump, SMART data, DIMM report or driver inventory was read in this work.
**Storage access:** none.

---

## 1. What the owner reports

| # | Report | Approximate timing |
|---|---|---|
| E-1 | **OS reinstalled.** Memory controller suspected weak. | September 2026 |
| E-2 | **Two of four 32 GB DIMMs removed** (4×32 → 2×32, i.e. 128 → 64 GB). BSODs stopped for about a month, then resumed. | September 2026 |
| E-3 | **`E:` had paging-file problems.** Powering the drive down made the BSODs "nearly entirely" stop. | during the fault period |
| E-4 | **`E:` physically removed from the machine.** | 2026-11-13 ("yesterday") |
| E-5 | **Minidump review (separate session)** found one consistency: a **driver-stack collision** — iCUE, an ASUS Wi-Fi driver, and Afterburner. ASUS driver updated; Afterburner removed; system now runs much more smoothly. | before 2026-11-14 |
| E-6 | The **"44k" driver error** was traced to a **dying USB drive**, now disconnected. | before 2026-11-14 |
| E-7 | **ASUS Wi-Fi driver updated and Afterburner removed *before* `E:` was removed** — therefore before 2026-11-13, and before the 2026-11-14 19:16:03 run | 2026-11-13 or just before |
| E-8 | **Bugcheck `0x1E`** (`KMODE_EXCEPTION_NOT_HANDLED`), with **2 GB dumps**; all dumps reviewed in the separate session; the **driver stack was the only constant** across them | during the fault period |
| E-9 | **Remaining 64 GB in slots A2/B2** per the board manual. Every stick and every slot tested individually; **no memory or board fault reported**. Memory controller remains the best guess: removing two sticks stopped the BSODs for a month | 2026-11 |

`E:` is the drive the three GGUFs were **copied from** (`build/move-to-g.sh` L92, `E:\AI GGUF\GGUF` → `G:\AI GGUF\GGUF`). It is no longer in the machine.

---

## 2. What this changes

### 2.1 S3 §4's open question is partly answered

`S3-METADATA-TEST-RECONCILIATION` §4 recorded "Were the three GGUFs on `G:` before the Windows reinstall? — **unknown**, no install date is recorded anywhere". The owner now reports a **September 2026** reinstall. That narrows the window from "unknown" to "a month, not a date", and it is an **owner statement**, not a record. The file mtime is 2026-09-04 20:53:48, which may fall before or after the reinstall; nothing here resolves which.

### 2.2 The source drive is gone

Any hope of verifying the ingestion by inspecting `E:`, or by re-copying and comparing, is closed. The only surviving ingestion-time artefact is `G:\AI GGUF\e-manifest.tsv`, written by the migration script itself, which would give expected size and mtime per file at copy time. **It is now the single most valuable unread document in this investigation**, and reading one small text file is a far smaller ask than any probe.

### 2.3 A new mechanism class enters the record: page-file corruption

If the Windows page file lived on a failing `E:`, then memory pages evicted to that drive and read back could return **wrong bytes to any process** — including a hashing tool, a loader, or the diagnostic reader itself. This is not a read-path fault at all; it is corruption at rest in memory, laundered through storage.

It has exactly the signature the investigation has been struggling to explain: identical size, identical mtime, identical file identity, no write recorded, and a digest that differs between observations. It also explains why errors appeared across unrelated tools rather than in one path.

**The bound is now resolved, and it excludes the observation that matters most.** `E:` was removed on **2026-11-13**, before every 2026-11-14 observation. A page file on a drive that is no longer in the machine cannot serve pages. So page-file corruption is **inadmissible for the 19:16:03 divergence and for the 13:45 baseline**, and admissible only for readings taken before 13 November — which includes the 5 November digest.

### 2.4 Two more candidates, neither established

| Candidate | Status |
|---|---|
| **Driver-stack collision** (iCUE / ASUS Wi-Fi / Afterburner) | a credible explanation for the **BSODs**; not evidence that any I/O returned wrong data. Kernel-path corruption is possible, unshown |
| **RAM / memory controller** | still a candidate. Removing 64 GB changed the symptom; changing a symptom is not a diagnosis |

---

## 3. What this does **not** change

- The repository-verified hash for S1-B, and the fact that the 19:16:03 composite read matched it (`HOWL-EXT-HASH-01`, frozen).
- The fact that the other six observations do not match it.
- The absence of any authoritative reference for S1-A and S1-C.
- The root cause: **unidentified**.
- Every closed gate.

---

## 4. The constraint the match places on "the copy was bad"

This is the part of the new information that meets the hash finding, and it cuts in one direction.

The 19:16:03 read was **99.99992 % unbuffered**: 5,335,281,664 B in 1,273 sector-aligned calls with zero short reads, then 4,064 B buffered. The composite hashed to the published value, so **the cache-bypassing portion returned canonical bytes**, and the buffered tail did too.

An unbuffered read does not consult the system file cache — it still passes through the filter stack and the class/port drivers, so "device-level" would be too strong. The supported claim is narrower: at 19:16:03 a cache-bypassing path returned bytes that **hashed to the published artifact**. That makes "the copy from failing `E:` wrote corrupt contents into this file" **disfavoured** — it would require the storage to have been corrected afterwards, and no write is recorded (mtime unchanged since 2026-09-04).

**Stated at its real strength:** disfavoured, not excluded. `/COPY:DAT` can restore a timestamp after a write, so a later repairing write would not necessarily show in metadata; and an unbuffered read is *bypassed* by the cache, not *guaranteed* to bypass every layer in the stack.

The consequence, for S1-B only: the six mismatching digests now point more strongly at the **observing side** — cache contents, memory, page-file corruption, a filter, or the tooling — than at the bytes stored on the drive. That is a narrowing of where to look. It is not an identification of what to blame, and it says nothing about S1-A or S1-C, which have no reference to compare against.

---

## 5. Questions — status

| # | Question | Status |
|---|---|---|
| 1 | `E:` removed before or after 19:16:03 on 11-14? | **Answered: before.** Removed 2026-11-13 |
| 2 | Page file exclusively on `E:`, or system-managed? | **Unanswered** — now matters only for pre-13 November readings |
| 3 | Which DIMM slots hold the remaining 64 GB? | **Answered: A2/B2**, per the board manual |
| 4 | ASUS update / Afterburner removal before the 19:16:03 run? | **Answered: before** — just before `E:` came out |
| 5 | Any storage-stack or memory-management bugcheck codes, or WHEA events? | **Partly answered:** `0x1E` with 2 GB dumps; driver stack the only constant. **No storage-stack or WHEA code was reported** — absence of a report is not absence of an event |

### 5.1 Mechanism admissibility, by observation window

| Mechanism | Readings before 2026-11-13 (incl. 5 Nov) | Readings on 2026-11-14 (13:45, 19:16) |
|---|---|---|
| Page-file corruption on failing `E:` | **admissible** | **excluded** — drive removed 11-13 |
| Driver-stack collision (iCUE / ASUS / Afterburner) | **admissible** | **weakened** — both changes installed before the run |
| Stale or poisoned cache pages | admissible | **leading candidate** (§5.2) |
| Memory / memory-controller error in the read path | admissible | **admissible, untested under load** |
| A filter or agent in the storage stack | admissible | admissible |
| Genuinely different bytes stored on `G:` | admissible | **disfavoured** — the cache-bypassing read returned canonical bytes at 19:16:03 (§4) |

### 5.2 The narrowing for 19:16:03

Two of the four most attractive mechanisms are now removed or weakened **for the one observation that has a reference**, by the owner's own timeline: the drive that might have corrupted pages was gone, and the driver stack that might have misbehaved had already been corrected.

What remains is dominated by one state that survives both changes: **pages already resident in the system cache**. A cache population survives a drive's removal and survives a driver update. If those pages were filled from wrong reads *before* 13 November, they would keep serving that content afterwards, and a buffered read of them would never touch the device.

The observed shape fits that: the buffered digest `341fb28c…` was recorded **at least three times across hours** and never once matched the cache-bypassing read's digest, while that single read produced the canonical content. A resident wrong population explains a stable wrong stream without requiring any fault to be active at 19:16.

**This is the leading candidate, not a finding.** It is not established by anything in this record, and it makes a specific prediction that no current evidence tests: a buffered read of a file whose pages were **never** resident should match its unbuffered read. That prediction is exactly what the V-B scratch-file probe would test — a file written and read for the first time, with a known pattern — and it needs gate 1, because it writes to `G:`.

---

## 6. Operational consequence: the machine is no longer the observation environment

The observations were made on a machine with four DIMMs, a failing `E:` with an active page file, Afterburner installed, and an older ASUS Wi-Fi driver. None of those conditions holds now.

Any future probe therefore **replicates a question, not an environment**. A clean result would not show the original mechanism is gone; a dirty result would not show the same mechanism is present. This must be attached to every future staging plan, and it weakens the case for spending gates on the original file rather than on a controlled scratch file (S-9a V-B), which tests the mechanism under *current* conditions and needs no assumption about the old ones.

### 6.1 The practical priority, given no replacement budget

The owner has stated plainly that replacing hardware is not an option. That changes what this investigation is *for*. It is no longer a route to an RMA; it is a route to **knowing which copies can be trusted**, so that a backup is worth making and the files are not quietly re-verified against a bad reading.

What that implies, stated only as far as the evidence goes:

| File | What is known | What it means for backing up |
|---|---|---|
| `S1-B` Gemma | one cache-bypassing read matched the repository-verified reference | consistent with an intact copy, **not proof of one** — stability was never demonstrated. Verify future copies with a cache-bypassing read, and treat a single matching read as one observation |
| `S1-A`, `S1-C` | no reference exists; no authoritative hash obtainable | **integrity unknown and unknowable by hash.** A second independent copy is the only mitigation available |
| Any file already read by the old stack | a cached stream may be wrong while the disk is right | **re-verify from a fresh process after a restart, or with an unbuffered read**, before trusting a previously computed digest |

The third row is a precaution, not a proven method: it rests on one matched observation and one candidate mechanism. It is cheap, reversible and consistent with the evidence, which is the standard for acting under uncertainty rather than the standard for a finding.

---

## 7. Claim ceiling

**Establishes:** that owner-reported environment history materially widens the mechanism list (page-file corruption, driver-stack collision, reduced memory population, a removed source drive); that the source drive can no longer be inspected; that the surviving ingestion manifest is now the best documentary artefact; and that the match against the repository-verified reference constrains — without refuting — the "bad copy from `E:`" explanation for S1-B.

**Does not establish:** that any of these mechanisms caused any observed digest; that the memory subsystem is damaged; that the driver collision produced wrong data rather than crashes; that the OS reinstall predates or postdates the 2026-09-04 mtime; or that the current machine would or would not reproduce the disagreement. Every item in this document is owner-reported, unverified here, and no gate moved.

*Anchored, not explained — and now with a plausible reason to think the anchor and the disagreement have different parents.*
