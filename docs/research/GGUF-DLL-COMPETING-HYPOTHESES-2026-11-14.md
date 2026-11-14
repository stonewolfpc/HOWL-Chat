# `GGUF-DLL-COMPETING-HYPOTHESES`

**Identity:** `HOWL-HYP-2` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-03` Track 3
**Status:** reassessment. **No read, no volume access, no test run.**
**Supersedes the hypothesis ranking** in `WINDOWS-CACHE-HYPOTHESIS-REVIEW` §7 and the leaning in its §10, both of which were written before the layer-divergence result was located in the record.

---

## 1. The observation set

| ID | Observation | Provenance |
|---|---|---|
| **O-1** | `G:` S1-B: baseline `b5b1913e…` at 13:45, matched at 13:50 | integrity record; qualification v1 |
| **O-2** | `G:` S1-B: `f46b1049…` then `e6dca6eb…` at 14:03, server running | integrity record |
| **O-3** | `G:` S1-B: `e6dca6eb…` twice at 14:04, server stopped | integrity record |
| **O-4** | `G:` S1-B: `eac8960a…` at 14:05, idle | integrity record — the decisive idle transition |
| **O-5** | `G:` S1-A: baseline `2ca636d9…` at 13:45 → `0e124fab…` at migration → `2ca636d9…` again later | baseline; `h-migration-candidates`; deferral table |
| **O-6** | `G:` S1-C: `f820ef3f…` at baseline, at migration, and at the later independent check — **never varied** | `independent-hash-check.log` |
| **O-7** | Six `H:` model paths and `G:` paths all accepted an exclusive open at 13:03:38 — 0 blocked | `handle-enumeration-2026-11-14.txt` |
| **O-8** | `H:` copies at migration: A `de3b8d28…` (1 run), B `6ff1b33f…` (3 identical runs), C `7fca804a…` (1 run); all recorded `H_stable: true` | `H-MODEL-MIGRATION-CANDIDATES` |
| **O-9** | `H:` S1-C differs from `G:` S1-C; first difference at byte **58,444,208** | `independent-hash-check.log`, `cmp` |
| **O-10** | Same file, same minute, not in use, identical file id and volume serial: buffered `341fb28c…` vs unbuffered `d0027dd3…` — **different** | `S1B-G-BUFFERED-VS-UNBUFFERED-…191603Z.txt` |
| **O-11** | `cublas64_12.dll` copied `H:`→`D:` arrived with same size and mtime but a different digest from its source; both sides stable over 3 reads; a re-copy matched the source | `H-MIGRATION-VERIFY*` |
| **O-12** | Filters attached to `G:`: `WdFilter` (Defender), `bfs`, `FileInfo`, `Wof`, `gameflt`, `UCPD` — all Microsoft; no third-party AV filter | `filter-inventory-2026-11-14.txt` |
| **O-13** | `G:` geometry: logical sector 512, physical 4096; `H:` logical 512, physical 512; different volume serials | `geometry-G/H.txt` |

---

## 2. What O-10 changes

Every other observation is a difference **between two moments**, and every hypothesis can accommodate a state change. O-10 is a difference **at one moment, inside the stack**: two read mechanisms returned different byte streams from one file object while nothing held it open and no metadata moved.

That licenses one narrow, strong statement:

> At 19:16 on 2026-11-14, at least one layer of the read path returned bytes that another layer did not have. The file did not have a single well-defined content for readers at that instant.

It does **not** say which layer was wrong. Both streams were internally reproducible — `341fb28c…` had been seen twice before, and the unbuffered run's own byte accounting is complete. **Two stable but different streams** is the important shape: it is the signature of two persistent sources of truth, not of random noise.

That shape narrows the mechanism to three readings, all of which produce exactly it:

| Reading | What holds `B` (buffered stream) | What holds `U` (unbuffered stream) | Assumptions |
|---|---|---|---|
| **D-1 — stale cache after an unbuffered write** | an old, coherent cache population | the new device content | (a) an agent wrote the file with unbuffered I/O; (b) nothing invalidated the cache afterwards |
| **D-2 — device returns different bytes than the cache holds** | what was cached from an earlier, good read | wrong bytes from the device now | (a) a device-path fault; (b) the cache was populated before it began |
| **D-3 — a persistent wrong page in cache** | memory content that is wrong but stable | the device's true bytes | (a) a corrupting agent with access to those pages; (b) persistence across minutes |

**D-1's attraction:** it explains the *stability of both streams*, needs only two assumptions, and requires no fault in any component — only an agent that wrote unbuffered. Its cost: it requires a writer, and no writer was ever observed.
**D-2's attraction:** it makes the host innocent and matches the S1 hypothesis of a failing read path. Its cost: a device that returns *consistently the same wrong bytes* for a minute is a strange failure mode; most read faults are noisy.
**D-3's attraction:** it makes the baseline meaningful again (disk = `b5b1913e`, cached = wrong). Its cost: stable corruption of specific pages for minutes is not typical of random memory faults either.

None of the three is excluded. **None is favoured by evidence.** D-1 is the most economical; "most economical" is not "supported", and it is not adopted here.

One further fact from O-10 matters: **neither stream equals the baseline.** So at 19:16 the disk did not obviously hold `b5b1913e` either. Any reading must explain three distinct contents, not two — which is what makes D-2 and D-3 (which leave the disk at the baseline) less tidy than they first appear, and D-1 (disk changed once, cache stale) tidier.

---

## 3. The hypothesis space, as the work order defines it

Assumptions are counted as independent things that must be true beyond the observations. Explanatory coverage is **not** treated as evidence, and technical familiarity is not a ranking criterion.

| ID | Hypothesis | Explains | Does **not** explain | Additional assumptions | Evidence for | Evidence against | Distinguishable now? |
|---|---|---|---|---|---|---|---|
| **H-1** | **External writer** (backup, launcher, sync, user action) | O-2…O-5, O-10 via D-1, O-8 as copies of changed sources | O-6 (S1-C never changed); O-5's return to baseline needs a restoration | 3: (a) write access; (b) acts on these paths; (c) preserves size/mtime | volume label "Second backup"; timestamp-preserving writers are common; D-1 needs a writer | O-7 (no handle at 13:03 — but 10 h before O-10); no scheduled task references `G:`; O-6 stable | **no** — needs a handle at the right time, or USN history |
| **H-2** | **Filesystem filter / storage-management agent** | O-10 (a filter could serve or substitute different data per path); O-2…O-5 | O-6; requires a rule matching model files | 3: (a) a filter acts on content; (b) mtime-preserving; (c) silent | O-12: five Microsoft filters attached to `G:` | `WdFilter` does not rewrite content; no third-party AV filter; O-6 unaffected | **no** — would need per-filter instrumentation |
| **H-3** | **Host memory instability** | O-2…O-5, O-10 via D-3, O-5's return to baseline | why the streams are each *stable*; why `H:` (different machine path, same RAM) shows no such variation in its filed runs | 3: (a) a fault; (b) striking these pages; (c) intermittent | non-ECC consumer machine; a cached-read divergence is a memory-content divergence | two stable streams is not the shape of random corruption; no WHEA history on file | **no** — needs WHEA history or a memory test |
| **H-4** | **Device / controller read instability on `G:`** | O-2…O-5; O-10 via D-2 | O-6 (same device, stable); the *stability* of the wrong stream | 2–3: (a) a localised fault; (b) silent; (c) returning repeatable wrong bytes | O-9 shows `H:` content differs from `G:` after a copy; a failing region would produce that | SMART `OK`; control file `LFM2.5` stable throughout; repeatable wrong bytes is atypical | **no** — needs SMART deltas / read-verify on a **non-model** range |
| **H-5** | **Windows storage-stack behaviour** (coherence, dirty/stale cache, lazy write) | the *form* of O-10 — it is exactly what a stale cache looks like | **cannot generate bytes by itself.** It is a mechanism of divergence, not a source of content | 1 (that the stack behaves as documented) + one of H-1/H-3/H-4 to supply wrong bytes | File Caching / File Buffering documentation; O-10's shape | none — but it is not a cause on its own | **not a standalone hypothesis** — must be paired |
| **H-6** | **Source/destination identity confusion** | O-8/O-9 (a copy read a different object) | O-10 (identity was recorded: same file id, same serial, within one minute) | 2: (a) two objects with equal size/mtime; (b) the copy resolved the wrong one | `(1)` in the S1-B filename is a re-download marker; duplicates with near-identical names exist | O-10 records a single file id across both reads; `cmp` found a concrete offset difference | **largely excluded for O-10**, unresolved for O-8/O-9 |
| **H-7** | **Historical provenance error** (already wrong on arrival) | why no local digest ever matches a publisher hash; nothing else | the *variation* — a stable wrong file cannot produce six digests | 1: the published file was bad, or the download was | baseline post-dates mtime by 10 weeks | no publisher hash exists to test it | **no** — needs a public hash, not storage |
| **H-8** | **Combinations** (e.g. H-1 + H-5; H-4 + H-5) | the full shape, including O-10 with two stable streams | — | the union, minus shared assumptions | as above | — | **no** |

### 3.1 What is *not* a hypothesis

- **"The cache"** is not a cause; it is a place where a divergence becomes visible. H-5 exists only to say that the stack can hold two coherent views — it must be paired with something that put wrong bytes into one of them.
- **"Windows 11"** is a *scope*, not a mechanism. No acknowledged Windows 11 defect matching this signature exists (`WINDOWS-11-STORAGE-EVIDENCE-MATRIX` §2.1). Invoking the OS generically is not an explanation.
- **"The file was corrupt"** restates the observation. Which layer, when, and how are the questions.

---

## 4. Pairwise distinguishability

| Pair | Distinguishable by | Currently distinguished? |
|---|---|---|
| H-1 vs H-3 | a handle or USN record at the relevant minutes, plus WHEA history | **no** |
| H-1 vs H-4 | whether the device returns the *same* wrong bytes repeatedly (H-4) or the file changed once (H-1) | **no** — both produce repeatable wrong bytes |
| H-3 vs H-4 | reading a **non-model** file on `G:` both ways; and checking another volume for the same behaviour | **no** |
| H-1 / H-3 / H-4 vs H-6 | recorded file identity | **yes for O-10** (identity recorded); no for O-8 |
| H-7 vs everything else | a publisher hash for one candidate | **no** — not obtained |
| D-1 vs D-2 vs D-3 | one buffered read after the cache is dropped, compared with the unbuffered stream | **no** — not performed |

---

## 5. What the evidence actually supports

| Statement | Support |
|---|---|
| At 19:16 the read stack held two different, individually stable, byte streams for one file object | **confirmed by observation** |
| `G:` S1-B took at least six distinct values over six hours with constant size and mtime | **confirmed by observation** |
| `G:` S1-C never varied | **confirmed by observation** |
| A copy on this machine can produce a stable, metadata-identical, content-different destination | **confirmed by observation** (O-11) |
| Some agent or fault changed bytes on `G:` | **not established** |
| The host is at fault | **not established** |
| The device is at fault | **not established** |
| Any hypothesis above is more likely than another | **not established** — assumption counts differ by at most one, and no discriminating test has run¹ |

¹ One discriminating test *has* run — O-10 — and it distinguished the **family**: whatever is happening involves the read stack, not only file history. It did not distinguish **members**.

---

## 6. If one more question could be asked

| Rank | Question | Distinguishes | Cost |
|---:|---|---|---|
| 1 | Read one fixed range of a **non-model** file on `G:` buffered, drop the cache for that file only, read it again, compare with an unbuffered read | D-1 vs D-2 vs D-3 — the first real separation available | one bounded read, three passes, no model |
| 2 | Do the same on `H:` and on `C:` | H-4 (device) vs H-3 (host) | same, two more volumes |
| 3 | WHEA history + memory configuration | H-3 | log read only |
| 4 | SMART deltas and link counters | H-4 | attribute read |
| 5 | Handle + USN for the 19:16 window (historical) | H-1, historically | gated |

None is authorised by this work order.

---

## 7. Claim ceiling

**Establishes:** a complete observation set with provenance; that the divergence is inside the read stack rather than in file history; three readings of that divergence with their assumption counts; eight hypotheses with what each explains, misses, needs, and is contradicted by; and that no hypothesis is currently distinguishable from its competitors.

**Does not establish:** the cause. It also does not establish that the *most economical* reading (D-1, a stale cache after an unbuffered write) is true — it is recorded as the tidiest, and tidiness is not evidence.
