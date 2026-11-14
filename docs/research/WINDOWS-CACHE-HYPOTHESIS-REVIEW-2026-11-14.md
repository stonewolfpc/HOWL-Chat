# `WINDOWS-CACHE-HYPOTHESIS-REVIEW`

**Identity:** `HOWL-CACHE-REVIEW-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-02` Track 2
**Status:** review and correction. **No read performed, no volume opened, no cache purged, no test run.**
**Corrects:** `WINDOWS-11-STORAGE-EVIDENCE-MATRIX` §3.2 and its rank-2 test. The correction is stated in §7 and §10 rather than applied silently to the earlier document, which stands as filed.

---

## 0. The claim under review

The previous matrix concluded:

> "**M-01/M-02 (cache-mediated reads) is the only mechanism that fits every observation with no additional assumption.**"

That sentence is withdrawn. It was an inference dressed as a result. This review tests it against documentation and finds three defects in it:

1. **"Cache" is a location, not a cause.** The system cache holds bytes in RAM. For cached bytes to differ from what a read returned earlier, *something must have changed RAM* — a memory fault, a driver writing into those pages, or a refill from a device that returned different bytes. Naming the cache names where a difference becomes visible, not what produced it.
2. **The uniqueness was not demonstrated**, only asserted. At least four hypotheses fit the observations; none is identified by them.
3. **The proposed discriminating test does not discriminate.** A buffered-versus-unbuffered comparison can disagree for entirely ordinary reasons, and may not disagree at all even when the cache is wrong. Details in §5 and §6.

---

## 1. What buffered `ReadFile` actually guarantees

| Statement | Source | Class |
|---|---|---|
| By default, "read operations read file data from an area in system memory known as the system file cache, rather than from the physical disk" | File Caching, Microsoft Learn | **MS-DOC** |
| The cache is a **write-back** cache, managed per file object, under the cache manager | File Caching | **MS-DOC** |
| Data in the cache is written to disk "at intervals determined by the operating system" | File Caching | **MS-DOC** |
| `ReadFile` copies bytes from that cache into the caller's buffer | File Caching + `ReadFile` semantics | **MS-DOC** |

**What a successful buffered read therefore guarantees:**

- that the requested byte range was copied into the caller's buffer;
- that every *cache-enabled* reader of that file sees the same bytes at that moment (cached I/O is coherent with itself);
- nothing else.

**What it does not guarantee:**

| Non-guarantee | Consequence for this investigation |
|---|---|
| That the bytes equal what is on the medium | a digest is a statement about the cache's contents, not the platter's |
| That the sector was readable at all | a fully cached read can succeed while the device is unreadable |
| Durability or ordering | irrelevant to reads, decisive for the copy path |
| Any checksum | a wrong-but-consistent page is returned without complaint |

**The corollary that drives everything below:** a buffered read returning `X` twice is evidence that *the cache held X twice*. A buffered read returning `X` then `Y` is evidence that *the cache contents changed*. It is not, by itself, evidence about the disk, because the change could have occurred in RAM.

---

## 2. What `FILE_FLAG_NO_BUFFERING` bypasses — and what it does not

| Bypasses | Does **not** bypass |
|---|---|
| System caching of the file's **data** for reads and writes ("all read and write operations directly access the physical disk") | **The device's own cache.** `CreateFile` documentation: "This flag does not affect hard disk caching or memory mapped files" — **MS-DOC** |
| Readahead into the cache for that handle | **File metadata caching.** File Caching: "the file metadata may still be cached" — **MS-DOC** |
| — | **Memory-mapped views.** Explicit in the same `CreateFile` sentence — **MS-DOC** |
| — | **The filter/minifilter stack.** Minifilters attach to the file system stack and see I/O regardless of the caller's cache flags — **MS-DOC** (architecture), **INFERRED** for this specific claim |
| — | **The volume, class and port drivers**, and the physical link |

So `NO_BUFFERING` narrows *one* layer — the system cache — and leaves four others untouched. Any test built on it tests the system cache layer and nothing else.

---

## 3. Alignment requirements for unbuffered reads

Documented constraints, all of which apply to a probe:

| Constraint | Requirement |
|---|---|
| File offset | exact multiple of the volume's sector size |
| Transfer length | exact multiple of the volume's sector size |
| Buffer address | aligned to the physical sector size (relevant for 4 KiB-native devices) |
| Violation | the call fails (`ERROR_INVALID_PARAMETER`) rather than being silently adjusted |

Practical consequence: a probe must read whole sectors and cannot compare a partial trailing region the way a hash does. A hash of a 5.3 GB file ends at a non-sector boundary; a probe must round, and the rounded read is not the same range as the hash. That is a design constraint on any future test, not a defect in the code.

---

## 4. Can a buffered and unbuffered read disagree with no external writer?

**Yes — and this is the part the previous matrix got wrong.** Two ordinary, documented paths produce disagreement:

| Path | Mechanism | Class |
|---|---|---|
| **Dirty cache** | a buffered write left pages dirty; a buffered read returns the new bytes from RAM; an unbuffered read fetches the old bytes from the device | **MS-DOC** (write-back semantics) |
| **Unbuffered write then cached read** | an unbuffered write by any process bypasses the cache; cached pages remain stale until invalidated | **MS-DOC** for the bypass; **technical commentary** for the coherence consequence |

And in the opposite direction — **the test may not discriminate at all**:

| Claim | Source | Class |
|---|---|---|
| For a file already in cache, a non-buffered request may still be satisfied from the cache, because the cache must stay coherent regardless of access mode | NTDEV / OSR developer discussion (kernel developers) | **FIELD** (expert commentary, not documentation) |
| Microsoft's own text says unbuffered operations "directly access the physical disk" | File Caching | **MS-DOC** |

These two statements are in tension. The resolution matters: if the commentary is right for resident pages, an unbuffered probe *cannot* reveal a corrupt cache page, and the test would return "agree" while the cache is wrong. If the documentation is right, the test is informative. **The behaviour has not been measured on this machine and cannot be asserted either way.**

**Verdict on the test:** it is worth running eventually because it is cheap and bounded, but it is **not a discriminator**, and its result must be interpreted with both readings in view. It has been demoted from rank 2 to a conditional, secondary step in the staging plan.

---

## 5. Conflation audit — six different things called "cache"

| # | Thing | Where it lives | What a divergence there would mean | Distinguishable by a digest? |
|---|---|---|---|---|
| 1 | **System file cache** (page cache) | host RAM | host memory held different bytes | no — a digest is the only signal |
| 2 | **Mapped view** of the model (loader) | the *same* physical pages as #1 for that file | nothing new — it is #1 seen through a mapping | no |
| 3 | **Device cache** (drive's own RAM) | the drive | the drive returned different bytes for the same LBA | no |
| 4 | **Controller / HBA / link buffers** | controller, cable | transient corruption on the wire | no |
| 5 | **Host RAM outside the cache** (the hasher's own buffer) | host RAM | a fault corrupting any read into memory | no |
| 6 | **The filesystem's metadata cache** | host RAM | stale size/time, not stale content | no (and size/mtime did not move) |

Three consequences the record must respect:

- **#1 and #2 are the same memory.** A mapped and a buffered reader of one file are expected to agree, always. "Mapped pages" is not a separate explanation from "system cache"; treating them as two mechanisms inflates the hypothesis count without adding information.
- **#1, #2 and #5 are all "host memory was wrong".** They differ in *which* memory and *which* agent corrupted it, not in kind. A digest cannot separate them.
- **#3 and #4 are "the device path returned different bytes".** Also one family, separated only by which component failed.

So the mechanism space is really **three families**, not six: host memory, device path, and an agent that changed the bytes deliberately (writer/minifilter). That is the honest granularity.

---

## 6. Do the observations uniquely identify any of them?

| Observation | Host memory | Device path | Writer / minifilter |
|---|---|---|---|
| 4 distinct digests, `S1-B`, constant size and mtime | fits | fits | fits |
| Two consecutive idle reads identical, a third different | fits (pages stable then changed) | fits (region readable then not) | fits (wrote between reads) |
| **`S1-A` returned its baseline again later** | fits (bad pages evicted, refilled good) | fits (region recovered / cache refilled) | **poor** — requires a writer that restored the original bytes exactly |
| Control file on the same volume stable | fits (different pages) | fits (different region) | fits (path-specific) |
| GGUF structure identical across loads | fits (localised) | fits (localised) | fits (in-place patch) |
| `cublas64_12.dll` copy differed once, identical on re-copy | fits (one bad read) | fits (one bad read) | **poor** — no plausible reason to rewrite an unrelated DLL |
| `H:` readings disagree across time (if X-2) | fits | **poor** — two different devices would both have to fail | poor |
| `H:` readings disagree across time (if X-1) | fits | fits | fits |

**Every column fits most rows.** The only rows with discriminating power point *away* from a writer: the baseline returning exactly, and the DLL. Neither points to a specific mechanism.

---

## 7. Competing hypotheses, with assumption counts

Assumptions are counted as *independent things that must be true beyond the observations themselves*. Counting is not ranking: a hypothesis with three assumptions can be truer than one with two. It is included because the previous matrix reasoned as though explanatory coverage were evidence.

| ID | Hypothesis | Core claim | Additional assumptions | What it uniquely predicts | What would refute it |
|---|---|---:|---|---|
| **H-1** | **Writer / minifilter changed the bytes** | something with write access altered payload in place | 3: (a) a process with write access exists; (b) it acts on these paths; (c) it preserves size and mtime | open write-capable handle or a USN data-overwrite record at a time when no read was in progress | handle enumeration on the exact path showing no writer **and** USN showing no writes over the window |
| **H-2** | **Host memory fault** (cache pages, hasher buffer, or both) | RAM held or produced wrong bytes | 3: (a) a fault exists; (b) it strikes these pages; (c) it is intermittent enough to spare the control | WHEA/memory-error history; the same corruption appearing in unrelated files over time; a failing memory test | clean extended memory test plus no recurrence anywhere else |
| **H-3** | **Device-path read instability** (`G:`) | the drive/controller/link returns different bytes for one LBA range | 3: (a) a localised bad region; (b) silent (no SMART, no event); (c) intermittent | SMART reallocated/pending sectors or link CRC counters; a read-verify error localised to the same offset | clean SMART, clean counters, and identical bytes on repeated controlled reads |
| **H-4** | **Agent inside the storage stack** (AV real-time scan rewriting, backup filter) | a registered filter altered or substituted data | 3: (a) a filter matches these extensions; (b) it preserves timestamps; (c) it acts in idle windows | the filter present in `fltmc`; a detection or restore log for these paths | no such filter, and no product log referencing them |
| **H-5** | **Tools / hashing** | the measurement was wrong | 0 standing: refuted — a text file hashes identically twice, and reads matched the manifest | — | already refuted |
| **H-6** | **Whole-file replacement** | the file was swapped | refuted: size, mtime and GGUF structure all unchanged | — | already refuted |
| **H-7** | **Wrong-on-arrival** | the download was already corrupt | 1: the upload itself was bad | a stable digest that never matches any reference | a publisher hash matching one of the observed digests |

**H-7 is worth separating out because it is orthogonal.** It does not explain *variation* at all — a stable wrong file cannot produce four digests — but it does undercut the assumption that "baseline" means "correct". Both can be true: the file arrived wrong *and* the reads varied. The record currently cannot separate them, because no publisher hash exists.

---

## 8. What the observations actually constrain

| Constraint | Strength |
|---|---|
| No write was observed, and no write-capable handle was ever enumerated | moderate — absence of observation, not observation of absence |
| Bytes changed while nothing in the workflow ran (observations 5→7) | **strong** — excludes the harness and the loader for that change |
| A baseline digest was reproduced later | **strong** — argues against a persistent rewrite |
| An unrelated DLL showed the same metadata-identical, content-different signature | **moderate** — one occurrence, but same shape, different file type |
| Two volumes may both have varied (`H:` readings) | **currently unresolvable** — depends on X-1/X-2/X-3 |
| Size and mtime never moved | **strong** against ordinary user-mode rewrite |

The strongest single constraint is the last-but-three: **the change happened while nothing was running.** That is compatible with a writer acting in that window, with a memory fault, and with a device returning different bytes — and it rules out the two actors we can name.

---

## 9. What would discriminate (summary; detail in the staging plan)

| Question | Discriminates | Cost |
|---|---|---|
| Were the `H:` copies re-made between readings? (X-1/X-2/X-3) | one volume vs two → device-path hypothesis gains or loses its hardest requirement | metadata only |
| Is a write-capable handle open on the exact path *now*? | H-1/H-4, **only if positive** | one enumeration |
| Does a USN record exist for the path over the window? | H-1/H-4, historically | one bounded journal query (gated) |
| Do SMART and link counters show errors? | H-3 | one attribute read |
| Is there WHEA/memory-error history? | H-2 | log read |
| Do buffered and unbuffered reads of a *non-model* range agree? | weakly: cache layer only, ambiguous in both directions | one bounded read |
| Do three cold reads of a candidate agree? | nothing about cause; it is a resumption gate, not a diagnostic |

---

## 10. Corrected confidence statement

Replacing the sentence withdrawn in §0:

> **No mechanism is identified.** Cache-mediated reads are *compatible* with every observation, but compatibility is not identification: the cache is where a differing byte becomes visible, and any of three families — host memory, device path, or an agent that changed the bytes — would produce the same visible result. At least four hypotheses (H-1 to H-4) fit the record, none is uniquely predicted by it, and the strongest available discriminator is not a cache property at all but the unresolved question of whether the `H:` copies were re-made. Until that is settled, assigning confidence to any single mechanism would overstate the evidence.

**Directional reading, stated as such:** the two observations with any discriminating power (the baseline returning exactly, and an unrelated DLL showing the same signature) both point away from a deliberate writer and toward a read-path phenomenon. That is a lean, not a finding — and it is the opposite of proof, because a writer that restores bytes and a filter that substitutes on read produce identical appearances.

---

## 11. Claim ceiling

**Establishes:** what a buffered read does and does not guarantee; exactly one layer that `NO_BUFFERING` bypasses and four it leaves alone; that buffered and unbuffered reads can disagree with no writer present, and may not disagree even when the cache is wrong; that "cache", "mapped pages", "device cache" and "RAM" were being used as if they were one thing; and that no observation uniquely identifies a mechanism.

**Does not establish:** which family caused the divergence, whether the disk or the host is at fault, or that the `H:` anomaly is real rather than an undocumented re-copy.
