# `STORAGE-DIAGNOSTIC-STAGING-PLAN`

**Identity:** `HOWL-STAGE-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-02` Track 3
**Status:** **PLAN ONLY — nothing below has been executed. No gate is opened by this document.**
**Design rule:** each step must be independently authorisable, independently stoppable, and must not depend on an earlier step having produced a positive result.

---

## 1. Ordering principle

Steps are ordered by three criteria, in this priority:

1. **Least storage contact** — metadata before content, system metadata before volume metadata, volume metadata before payload.
2. **Discriminating power per byte** — a step that splits the hypothesis space is worth more than one that confirms a favourite.
3. **Non-destructibility** — nothing that consumes the evidence, dwells on the failing region, or purges state comes before the cheap questions are answered.

Explicitly **not** initial steps, and why:

| Rejected early step | Why not |
|---|---|
| Cache purge | Destroys the only instrument that has produced any signal, before it is understood |
| Repeated reads of a candidate model | ~5.3 GB each, changes no hypothesis without a prior discriminator, and consumes the artefact's state |
| `chkdsk /r`, surface scan, read-verify | Hours on a 5400 rpm HDD, can force remaps, and destroys the ability to observe the original symptom |
| Parallel hashing | Maximises I/O and heat on a device whose stability is the question |
| Copy-and-compare the candidates | Produces a third artefact of unknown provenance (`H:` is already that artefact) |

---

## 2. The sequence

Legend: **Reads content?** = does the step read file payload. **Writes storage?** = any write to any volume. **Touches target?** = `G:`, `H:`, or a model file.

### S-0 — Record reconciliation (currently the only step requiring nothing)

| Field | Value |
|---|---|
| Information collected | whether the two `H:` readings in the record can be two readings of one unchanged file: filed hashes, the migration record's stability claims, any copy re-run, any `H:` re-copy entry |
| Reads content? | **no** |
| Writes storage? | **no** (reads filed documents only) |
| I/O volume | kilobytes |
| Elevation | no |
| Positive establishes | that the X-1/X-2/X-3 ambiguity is a **filing** question, resolvable without touching a disk — and what the record actually claims |
| Negative does not establish | anything about storage; a record that is merely silent is not a record that a re-copy did not happen |
| Touches `G:`/`H:`/model? | **no** |

Status: **done** — see `GGUF-FILE-PROVENANCE-MATRIX` §4. It resolved the ambiguity *to a question* (X-1 vs X-2 vs X-3) without resolving which answer is true.

### S-1 — OS build, cumulative update, boot and shutdown history

| Field | Value |
|---|---|
| Information collected | Windows edition/build/CU; last boot time; whether an unexpected shutdown occurred between 13:45 and 14:10 on 2026-11-14 (Kernel-Boot / EventLog 6008/41) |
| Reads content? | no |
| Writes storage? | no (event log is written by the OS regardless; this step writes nothing) |
| I/O volume | kilobytes |
| Elevation | build: no; **event log: yes** (and it is a **closed gate**) |
| Positive establishes | which Windows storage regressions are even applicable; whether a state transition (dirty shutdown, resume) occurred inside the observation window, which H-2-adjacent ordering effects require |
| Negative does not establish | that no state transition occurred at a finer grain; that the build is unaffected by an unlisted defect |
| Touches `G:`/`H:`/model? | **no** |

### S-2 — Storage stack inventory: filter drivers, controller, drivers

| Field | Value |
|---|---|
| Information collected | `fltmc filters` (minifilter list and altitudes); storage controller and driver versions; volume's filesystem parameters; write-cache policy |
| Reads content? | no |
| Writes storage? | no |
| I/O volume | kilobytes |
| Elevation | **yes** for `fltmc`; driver enumeration is usually unprivileged |
| Positive establishes | whether an H-1/H-4 agent is *registered* in the stack, and whether any product could plausibly see model I/O |
| Negative does not establish | that no filter acts (a driver can act without appearing as a minifilter; an AV can use a kernel callback); that a listed filter is responsible |
| Touches `G:`/`H:`/model? | **no** (stack metadata is volume-independent) |

### S-3 — `H:` metadata only (the X-1/X-2 discriminator)

| Field | Value |
|---|---|
| Information collected | creation time, last-write time, size and file ID for the three `H:` copies and for `D:\H-MIGRATION-2026-11-14` staging; directory entry timestamps |
| Reads content? | **no** |
| Writes storage? | no |
| I/O volume | kilobytes (metadata) |
| Elevation | no |
| Positive establishes | if creation time post-dates the first `H:` reading, a **re-copy occurred** (X-1) and every `H:` anomaly belongs to the `G:` read path. If creation time predates both readings and the digests differ, two volumes produced varying reads (X-2) — the single most consequential fact available |
| Negative does not establish | that no in-place write occurred: a metadata-preserving writer leaves creation time intact. robocopy sets creation time at copy unless `/DCOPY:T`; the migration's flags are not recorded, so an unchanged creation time is **suggestive, not decisive** |
| Touches `G:`/`H:`/model? | **`H:` metadata only** — no payload, no model content |

Ranked first among storage-touching steps because its result halves the hypothesis space.

### S-4 — Handle enumeration on the exact path

| Field | Value |
|---|---|
| Information collected | for each process: handle presence on `G:\AI GGUF\GGUF\` and on the three exact file paths; where available, the granted access mask |
| Reads content? | no |
| Writes storage? | no |
| I/O volume | negligible |
| Elevation | **yes** — enumerating another process's handles requires `SeDebugPrivilege`; `handle.exe` is not installed, so a tool decision is part of authorising this step |
| Positive establishes | that a process **has the file open at that instant**. With a granted-access mask of write/data-write, that it holds a **write-capable** handle — which is still not proof it wrote, only that it could |
| Negative does not establish | that no writer exists: absence now says nothing about a scheduled writer, an intermittent one, a historical write, or a kernel-mode accessor that never takes a user-mode handle |
| Touches `G:`/`H:`/model? | **`G:`** — a handle-adjacent query; the file is not opened by this step's own code |

**Framing constraint carried from the work order:** a handle is evidence of access, not of a write. The step must report the handle and, where available, the access mask — never the word "written" — unless a write is otherwise established.

### S-5 — Writer capability: ACLs and ownership

| Field | Value |
|---|---|
| Information collected | `icacls` on `G:\AI GGUF` and the three files: which principals have write/modify, whether inheritance is disabled, owner |
| Reads content? | no |
| Writes storage? | no |
| I/O volume | negligible |
| Elevation | no (read of security descriptors) |
| Positive establishes | the population of accounts that *could* write; a writable-by-everyone directory raises H-1's plausibility |
| Negative does not establish | that any of them did; capability is not action |
| Touches `G:`/`H:`/model? | **`G:`** metadata |

### S-6 — Scheduled and background writers inventory

| Field | Value |
|---|---|
| Information collected | scheduled tasks and services whose action references `G:`, `AI GGUF`, backup or sync binaries; installed backup/sync products; startup entries |
| Reads content? | no |
| Writes storage? | no |
| I/O volume | kilobytes |
| Elevation | yes for the full task store (system tasks); user tasks readable unprivileged |
| Positive establishes | a **named** scheduled writer — the cheapest resolution available for H-1 |
| Negative does not establish | that none exists: a product's background thread need not have a task, and command lines do not contain open paths |
| Touches `G:`/`H:`/model? | **no** |

### S-7 — Device health counters

| Field | Value |
|---|---|
| Information collected | SMART attributes (reallocated, pending, uncorrectable, power-on hours, CRC/link error count) for disks 0 and 2; `Get-StorageReliabilityCounter`; controller error counters |
| Reads content? | no |
| Writes storage? | no |
| I/O volume | kilobytes |
| Elevation | yes for some counters |
| Positive establishes | support for H-3 if attributes are degrading, or a link-error count consistent with intermittent transfer faults |
| Negative does not establish | that the device is healthy: SMART is firmware-reported and frequently blind to localised read faults. `HealthStatus: Healthy` was already recorded and proves nothing |
| Touches `G:`/`H:`/model? | **no** |

### S-8 — Memory and machine-error history

| Field | Value |
|---|---|
| Information collected | WHEA-Logger events, prior memory-diagnostic results, RAM configuration and whether ECC is present; page-file/memory pressure at the observation time if logged |
| Reads content? | no |
| Writes storage? | no |
| I/O volume | kilobytes |
| Elevation | yes (system log) — **event log gate** |
| Positive establishes | support for H-2 (a correctable-error history would materially raise it) |
| Negative does not establish | absence of a memory fault: uncorrected, infrequent faults in non-ECC memory are typically invisible |
| Touches `G:`/`H:`/model? | **no** |

### S-9a — Bounded buffered/unbuffered comparison on a **non-model** range

| Field | Value |
|---|---|
| Information collected | SHA-256 of one fixed sector-aligned range (e.g. 16 MiB) of a non-model file on `G:`, read once buffered and once with `FILE_FLAG_NO_BUFFERING`, plus the same range from a second file |
| Reads content? | **yes — payload, but not a model** |
| Writes storage? | no (read handles only; no cache purge) |
| I/O volume | ~64 MiB total |
| Elevation | no |
| Positive establishes | a disagreement between access modes on a non-model file would show the cache layer and the device disagreeing — but see below |
| Negative does not establish | that the cache and device agree, because a resident page may satisfy an unbuffered request (the documented and the practitioner accounts differ — `WINDOWS-CACHE-HYPOTHESIS-REVIEW` §4). Agreement is therefore uninformative; only disagreement is informative |
| Touches `G:`/`H:`/model? | **`G:` payload, non-model file** |

### S-9b — Three cold reads of one candidate (resumption gate, not a diagnostic)

| Field | Value |
|---|---|
| Information collected | three digests of `S1-C` (the only candidate with a stable `G:` history), separated by process exit |
| Reads content? | **yes — a model file** |
| Writes storage? | no |
| I/O volume | ~15.5 GB |
| Elevation | no |
| Positive establishes | nothing about cause. It answers resumption criterion 2 for one slot only |
| Negative does not establish | that the volume is unreliable in general, or that any other file is affected |
| Touches `G:`/`H:`/model? | **`G:` model payload** |

This is deliberately last: it is the most expensive step, the only one that reads a model, and the least informative per byte.

---

## 3. Decision logic — what combinations mean

| Result | Reading | What it does **not** mean |
|---|---|---|
| S-3 shows `H:` creation time after the first reading | X-1: a re-copy happened; every anomaly is on the `G:` read path | that `G:` is faulty — only that the copies cannot be independent witnesses |
| S-3 shows unchanged creation time and differing digests | X-2: two devices produced varying reads → host-level cause gains weight | that the host is at fault; a metadata-preserving writer still fits |
| S-4 finds a write-capable handle | a concrete H-1/H-4 candidate, named and inspectable | that it wrote anything, or when |
| S-4 finds nothing | no user-mode handle at that instant | nothing about the past, or about a scheduled writer |
| S-6 names a backup task touching `G:` | the strongest writer candidate; prioritise S-3 history against its schedule | that its writes explain the digests without further checks |
| S-7 shows degrading SMART | H-3 supported, and no further qualification should be attempted on that device | that SMART caused anything by itself |
| S-8 shows WHEA correctable errors | H-2 materially supported | that the specific pages were affected |
| S-9a disagrees | cache layer and device disagree on a non-model file → read-path phenomenon confirmed, still not which layer | which layer |
| S-9a agrees | **uninformative** (see §2, S-9a negative) | that the cache is sound |

**No single row is conclusive, and several rows being empty is not exoneration.** The plan's purpose is to make each question answerable once, cheaply, rather than to converge on a narrative.

---

## 4. Authorisation map

| Step | Gate | Currently open? |
|---|---|---|
| S-0 | none (documents) | yes — completed |
| S-1 (build only) | system query | no |
| S-1 (event log) | **event-log enumeration** | **no** |
| S-2 | system query + elevation | no |
| S-3 | `H:` metadata access | no |
| S-4 | `G:` + handle tool + elevation | no |
| S-5 | `G:` metadata | no |
| S-6 | system query + elevation | no |
| S-7 | device counters + elevation | no |
| S-8 | **event log** | **no** |
| S-9a | **`G:` payload read** | no |
| S-9b | **model read** + S1-B journal gate context | no |

The three diagnostic gates named in the standing rule — `G:` disposable-file validation, event-log enumeration, and the bounded S1-B journal query — remain **closed**, and this plan adds no authorisation of its own.

---

## 5. Claim ceiling

**Establishes:** a staged, minimum-access sequence in which each step names exactly what it collects, whether it reads content, whether it writes, how much I/O it costs, whether it needs elevation, what a positive result establishes and what a negative result does not.

**Does not establish:** that any step will be conclusive; that the sequence is optimal in the presence of a result not yet seen; or that the investigation converges. Step S-3 is ordered first because it is cheap and splits the space — not because it is expected to produce an answer.
