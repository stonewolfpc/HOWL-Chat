# `WINDOWS-11-STORAGE-EVIDENCE-MATRIX`

**Identity:** `HOWL-WIN11-EM-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-01` Track 2
**Status:** evidence plan only. **No inventory command run. No device queried. No volume touched. No event log read.**
**Method:** primary Microsoft documentation first, then filed incident evidence. One-model review.

**Rule applied throughout:** a documented Windows behaviour is not a diagnosis. A field report is not a defect. Similar symptoms are not a shared cause. Every row below carries the class of evidence it rests on, and the class is the whole point of the table.

**Evidence classes used:**

| Code | Meaning |
|---|---|
| **MS-DOC** | Confirmed Microsoft-documented behaviour |
| **MS-DEFECT** | Confirmed Microsoft-documented defect (acknowledged in a Microsoft source) |
| **FIELD** | Credible field report — reproducible by third parties, not acknowledged by Microsoft |
| **UNVERIFIED** | Unverified report — single source, unreproduced, or contradicted |
| **POSSIBLE** | Technical possibility only — follows from architecture, no incident evidence |
| **NOT-SUPPORTED** | Contradicted or unsupported by the evidence available here |

---

## 1. The environment: known and unknown

Nothing below was measured during this work order. Each entry names where the fact comes from, or states that it is unknown.

### 1.1 Known

| Fact | Value | Source |
|---|---|---|
| OS family | Windows 11 **Pro** | `MEMORY-09-MODEL-POPULATION-AND-HARDWARE-PLAN` L17 |
| Volume under investigation | `G:`, label **"Second backup"**, NTFS, Fixed, `HealthStatus: Healthy` | integrity record §9.3 |
| Physical device | disk 2, **WDC WD10SPCX-75KHST0**, SATA, `Healthy`, `Online` | integrity record §9.3 |
| OneDrive | running; sync root `C:\Users\StoneWolfPC\OneDrive` — **not** `G:` | execution audit §1.2 |
| Loader | KoboldCpp 1.120, `H:\AI apps\koboldcpp.exe`, every launch with `--usemmap` | execution audit §1.1 |
| Loader-side structure | GGUF V3, 48 KV pairs, 720 tensors, `n_layer=42`, `n_embd=2560` — identical across loads | execution audit §2 |
| File size / mtime | constant: 5,335,285,728 B; mtime `2026-09-04 20:53:48` | integrity record §2 |
| Scheduled tasks referencing `G:` | none found | execution audit §1.2 |
| Driver-stack fault | iCUE / Afterburner — recorded as established separately, **not re-verified here** | deferral record §5 |

### 1.2 Unknown — and deliberately not measured here

| Fact | Why it matters | Status |
|---|---|---|
| Windows **build and cumulative update** | every storage regression is build-specific; without it no update-level hypothesis can be tested | **unknown** |
| Storage controller and driver version | a controller driver bug is a live hypothesis | **unknown** |
| NVMe/SATA controller mode and chipset drivers | affects cache/flush semantics | **unknown** |
| Filesystem parameters (`fsutil` behaviour flags, cluster size) | affects cache granularity | **unknown** |
| **Filter drivers** in the stack (AV, backup, encryption, sync, indexing) | a minifilter can observe, delay, split or alter I/O | **unknown** |
| Write-cache policy on the device and in the stack | determines what "written" means at each layer | **unknown** |
| Drive firmware version | relevant to the Phison-class reports below | **unknown** |
| `H:` volume identity, device, filesystem | the `H:` copies are part of the evidence | **unknown** |
| RAM configuration and whether ECC | memory instability is one of four explanations that fit unchanged size and mtime | **unknown** |
| Whether a backup product owns "Second backup" and writes on a schedule | the earliest and cheapest writer hypothesis | **unknown** |

That list is the honest state of the environment. Any matrix row that depends on an unknown cannot be settled by argument; it needs one of the queries in §5.

---

## 2. Documented Windows behaviour the matrix rests on

| # | Behaviour | Class | Source | Consequence for this incident |
|---|---|---|---|---|
| B-01 | `FILE_FLAG_NO_BUFFERING` disables system caching for the file's data, with sector-alignment requirements on offset, length and buffer address | **MS-DOC** | File Buffering, Microsoft Learn | A buffered reader and an unbuffered reader of the same region can return different bytes **with no writer involved** |
| B-02 | With caching disabled, reads and writes go to the physical disk, **but file metadata may still be cached**; `FlushFileBuffers` is required to flush metadata | **MS-DOC** | File Caching, Microsoft Learn | "No buffering" is not "no state"; a metadata claim needs its own flush |
| B-03 | `FlushFileBuffers` writes all buffered information for a file to the device; the documented route to durability is `NO_BUFFERING | WRITE_THROUGH` | **MS-DOC** | A rename without a flush does not establish durability — this is the basis of writer defect **D-3** |
| B-04 | `WRITE_THROUGH` alone sends writes to disk immediately but does **not** flush the device's internal cache; combined with `NO_BUFFERING` it also requests a hardware flush | **MS-DOC** (secondary, Microsoft engineering blog) | *The Old New Thing*, 2021-07-29 | Explains how a write can be "complete" and still not be on the medium |
| B-05 | `WRITE_THROUGH` is described by the same source as "largely ineffective nowadays, because SATA drivers ignore the flush request" | **FIELD** (Microsoft blog aside, 2014) | *The Old New Thing*, 2014-03-06 | Weakens any plan that relies on `WRITE_THROUGH` for durability on this SATA device |
| B-06 | File system filter drivers (minifilters) attach to the stack and may **monitor, filter, modify, or prevent** I/O | **MS-DOC** | About File System Filter Drivers, Microsoft Learn | A minifilter is a legitimate place for a payload to be altered or delayed, and it is invisible to file-level tooling |
| B-07 | The change journal "logs only the fact of a change … and the reason"; repeated changes of the same kind may produce **one** record, and it "does not record enough information to allow reversing the change" | **MS-DOC** | Change Journal Records, Microsoft Learn | A single USN record cannot bound how many writes occurred; absence of a record is never proof of no write |
| B-08 | Buffered and memory-mapped access to one file share the same cached pages | **MS-DOC** (implied by B-01/B-02 and the cache architecture) | File Caching + File Buffering | A mapped loader and a hashing process can agree with each other and both differ from the medium |

### 2.1 Windows 11 storage regressions — what is actually documented

| Item | What is claimed | Class | Status |
|---|---|---|---|
| **KB5063878** (August 2025, 24H2, build 26100.4946) and the preceding KB5062660 preview | drives disappearing or becoming unreadable during **sustained large writes**; testers reported onset above ~60 % drive use after ~50 GB of continuous writing; Phison NAND controllers and DRAM-less models over-represented; some drives recovered on reboot, others did not | **FIELD** | Third-party reproduction by multiple testers; Microsoft opened an investigation |
| Microsoft's conclusion on KB5063878 | "After thorough investigation, Microsoft has found no connection between the August 2025 Windows security update and the types of hard drive failures reported"; no repro on up-to-date systems; no telemetry increase in disk failures or file corruption | **MS-DOC** (service-alert statement, reported by BleepingComputer) | **No confirmed defect.** The reports are not thereby explained, and the investigation does not exclude a narrower population |
| General "Windows 11 corrupts files" claims | no primary source located | **NOT-SUPPORTED** | Not usable |
| This incident as an instance of KB5063878 | would require: the same build, a Phison-class controller, sustained writes at the time of the observations, and drive-level symptoms | **NOT-SUPPORTED** | `G:` is a **SATA WDC** device reported `Healthy`; the observed behaviour is *stable size, unchanged mtime, varying content* — not a disappearing drive |

**The honest position:** a documented, acknowledged Windows 11 storage corruption defect matching this symptom was **not found**. That does not eliminate Windows. It means the Windows hypothesis has to be argued from mechanism (the rows below) rather than from precedent, and that the closest precedent on record was investigated and not confirmed.

---

## 3. Mechanism matrix

Each row: what it could explain, what it cannot, evidence for and against, what kind of access settling it needs, and whether it can be evaluated without touching `G:`, `H:` or a model file.

Access classes: **READ-ONLY QUERY** (metadata, no payload) · **BOUNDED READ** (fixed-size payload read) · **DISRUPTIVE** (may change state, consume the evidence, or stress hardware — requires separate authorisation).

| ID | Mechanism | Could explain | Cannot explain alone | Evidence for | Evidence against | Access needed | Without touching `G:`/`H:`/model? |
|---|---|---|---|---|---|---|---|
| **M-01** | **Cache-page corruption** — resident pages for the file hold wrong bytes (bad RAM, driver DMA into cached pages) while the medium is fine | repeated reads agreeing within a session; new session differing; `H:` copies stable but different (the copy serialised corrupted pages); size and mtime untouched | why one file and not its neighbour, unless the corruption is localised to those pages | observations 5 and 6 identical, 7 different, all idle; S1-A returning its baseline again later | the `LFM2.5-8B` control on the same volume read consistently — weakens volume-wide corruption, not page-local | BOUNDED READ (compare buffered vs unbuffered read of one region) | **partly** — needs a payload read, but not a model: any file's region shows the same cache behaviour |
| **M-02** | **Buffered-read illusion** — the digest gate hashes *cached* bytes, so "different digest" may describe cache state, not disk state | all four digests, without any file change | cannot explain a change that survives a reboot and a cache purge | `sha256sum` uses buffered `ReadFile`; the loader maps the same pages | none — this is a property of the measurement, not a claim about the disk | BOUNDED READ (unbuffered read of the same range) | **partly** as M-01 |
| **M-03** | **Minifilter interference** — AV/backup/encryption/sync driver observes, delays, splits or rewrites I/O | writes that bypass page cache entirely; reads served differently to different callers; a payload altered where no user-mode writer exists | would normally produce a visible file event, and often a change in mtime | B-06; several such products are installed (iCUE, Afterburner, OneDrive) | OneDrive's root is on `C:`; no filter was enumerated | READ-ONLY QUERY (`fltmc filters`, driver list) | **yes** |
| **M-04** | **Mapped-page / cache coherence** — a mapped reader and an independent reader disagree because one is served from a section object | digests differing while the loader is attached (observations 2–4) | observations 5–7 diverged with **nothing** attached | `--usemmap` on every launch; B-08 | divergence continued after the mapping died | POSSIBLE only | **yes** (argument only) |
| **M-05** | **Write-back ordering** — a writer's data reached the medium after metadata, or not at all, across a power/state transition | a file whose *content* differs while size/mtime agree | the machine was not power-cycled between observations 5 and 7 | B-03, B-04 | no unexpected reboot is recorded in the timeline | READ-ONLY QUERY (boot time, dirty-shutdown events) | **yes** |
| **M-06** | **Device-level read instability** — a failing region returns different bytes for the same LBA | varying content, constant size/mtime, one file affected | leaves the control file unaffected only by locality | WD SATA device, `Healthy` (firmware metadata, **not** a read test) | control file consistent; no SMART flag reported | DISRUPTIVE (surface/read-verify test) — **not authorised** | **no** |
| **M-07** | **Controller / cable / link errors** — SATA link retries returning stale or garbled sectors | localised corruption, silent to the filesystem | usually produces event-log entries and CRC counters | SATA link, 2.5" WD drive | none available | READ-ONLY QUERY (SMART attribute read, controller error counters) | **yes** for the counters |
| **M-08** | **Memory instability** — bad RAM corrupting read buffers | any of the observations; different content per read | would usually affect more than one file over hours | four distinct digests; page-cache hypothesis (M-01) needs a corruption source | the control file hashed consistently — an argument against *frequent* corruption, not against a bad page | READ-ONLY QUERY (WHEA/`mdsched` history); a memory test is DISRUPTIVE | **yes** for the history |
| **M-09** | **Backup product writing "Second backup"** | payload rewritten with timestamps preserved; exactly the observed signature | requires the product to write that path on a schedule, and to keep the file open long enough | volume label; backup-class software plausible on this machine | no scheduled task references `G:`; no non-shell process names the path | READ-ONLY QUERY (installed backup products; open-handle enumeration) — **handle enumeration is the decisive check** | **yes** |
| **M-10** | **Indexer / search / thumbnailer touching payload** | reads only, no content change | cannot alter bytes | — | indexers do not rewrite payload | READ-ONLY QUERY | **yes** |
| **M-11** | **NTFS metadata/cache fault** | an inconsistent view of the same file | would normally surface as structural change or remount events | none | GGUF structure identical across loads (§2) | DISRUPTIVE (`chkdsk` read-only is still a volume operation) | **no** |
| **M-12** | **Antivirus quarantine/remediation** rewriting a suspicious binary | content change with preserved timestamps | would normally log a detection | GGUF models are a plausible AV target | no detection record examined | READ-ONLY QUERY (AV history) | **yes** |
| **M-13** | **Power management parking the drive mid-read** | timeouts and retries, very rarely different bytes | cannot change bytes without a read error being surfaced | external/2.5" drive | device is fixed, not USB | READ-ONLY QUERY (power plan, selective-suspend policy) | **yes** |
| **M-14** | **The digest tool itself** | apparent non-determinism | already excluded: a text file hashes identically twice; two observations matched the manifest | integrity record §10 | — | none needed | **yes** |

### 3.1 What the existing evidence already constrains

| Observation | Constrains |
|---|---|
| Observations **5 and 6 identical**, **7 different**, all with the loader stopped | excludes "the loader's mapping caused it" for the 5→7 change; leaves cache state, an external writer, or hardware |
| **Size and mtime constant** across every observation | excludes whole-file replacement, and excludes most ordinary user-mode write paths (which update mtime) |
| **GGUF structure identical** across loads (48 KV, 720 tensors, same dims) | the damage is localised, not structural — consistent with a few corrupted regions |
| **Control file stable** on the same volume in the same window | argues against volume-wide failure; does **not** clear the volume (a localised bad region or a path-specific writer produces the same pattern) |
| **S1-A recovered its baseline digest later** | content did not change and change back; the *read* varied |
| `H:` copies stable but different from `G:` baselines | the copy captured whatever `G:` returned at that moment — evidence about the read path, not proof about the original |
| `Healthy` from `Get-PhysicalDisk` | **nothing about read integrity** — it is firmware-reported condition metadata |

### 3.2 What the evidence positively supports

**M-01/M-02 (cache-mediated reads) is the only mechanism that fits every observation with no additional assumption.** It explains: agreement within a session; disagreement across sessions; divergence while idle; the control file being unaffected (different pages); and `S1-A` returning its baseline again once the bad pages were gone.

It is a *fit*, not a finding. It predicts something testable and cheap: a buffered read and an unbuffered read of the same range should disagree if the cache is wrong and the disk is right.

---

## 4. Where Windows 11 remains a live contributor

Not eliminated, and not promoted:

1. **A build-specific storage-stack regression** cannot be excluded while the build is unknown. Testing it requires the build number, which is a read-only query.
2. **The documented KB5063878 class** is the closest precedent, and it was investigated without confirmation. Our device is SATA WD, not Phison NVMe, and the symptom is not drive loss. Precedent does not apply; the class is not thereby impossible.
3. **A minifilter** is the one Windows-adjacent component that can alter user-visible bytes without appearing as any file operation the user performed.
4. **Cache behaviour itself is Windows behaviour.** If the digests measured cache and not disk, then the "instability" is a property of the read stack — which is exactly where Windows lives.

---

## 5. Ranked tests, with authorisation class

| Rank | Test | Settles | Access | Authorised now? |
|---:|---|---|---|---|
| 1 | **Open-handle enumeration** on the exact file and directory | M-03, M-09, M-12 — the largest gap in the investigation | READ-ONLY QUERY (needs a handle tool; none installed) | **No** — gate |
| 2 | **Buffered vs unbuffered read of one fixed range** of a *non-model* file on `G:`, hashed both ways | M-01, M-02 without touching a model | BOUNDED READ | **No** — `G:` gate |
| 3 | **Build, CU and boot history** (`winver`, `Get-HotFix`, last boot, dirty shutdown) | M-05, and which regressions are even applicable | READ-ONLY QUERY | **No** — system query, needs a decision |
| 4 | **Filter driver list** (`fltmc filters`, storage-class driver versions) | M-03, M-07 | READ-ONLY QUERY | **No** |
| 5 | **SMART attributes + controller CRC counters** | M-06, M-07 | READ-ONLY QUERY | **No** |
| 6 | **Installed backup product inventory**; whether it owns "Second backup" | M-09 | READ-ONLY QUERY | **No** |
| 7 | **WHEA / memory error history** | M-08 | READ-ONLY QUERY | **No** |
| 8 | **Three cold reads of each candidate after a cache purge** (resumption criterion 2) | M-01 vs M-06 directly | BOUNDED READ ×3 per file | **No** — `G:` gate |
| 9 | Surface scan / read-verify / `chkdsk /r` | M-06, M-11 | DISRUPTIVE | **No** — and not recommended |

**Explicitly not recommended at any authorisation level in this work order:** stress tests, parallel hashing of large files, broad volume scans, registry or power-policy changes, driver rollbacks, firmware updates, and any "verify by copying" step that consumes the original.

---

## 6. Claim ceiling

**Establishes:** the documented Windows behaviours that permit a read to differ from the medium without any writer; the environment facts on record and the substantially larger set that is unknown; that no acknowledged Windows 11 defect matching this symptom was found, and that the closest precedent was investigated without confirmation; and a ranked, cheapest-first set of read-only queries.

**Does not establish:** that Windows 11 caused the divergence; that the cache is at fault; that the disk is healthy or unhealthy; that any mechanism is exonerated. Every row in §3 remains open, and the two leading hypotheses — an external writer and a cache-mediated read — are both consistent with all current evidence.
