# `GGUF-WRITER-SOURCE-AUDIT`

**Identity:** `HOWL-WRITER-AUDIT-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-01` Track 1
**Status:** source audit complete. **No writer executed. No model file opened. No target read.**
**Method:** static reading of this repository plus the filed incident evidence. One-model review.

---

## 1. The finding that reframes the track

**This repository contains no code that writes chat-model GGUF files.**

Every `.gguf` in the app's chat path is *discovered, registered, inspected and read*. None is created, copied or converted. So the files at `G:\AI GGUF\GGUF\` were not written by HOWL Chat's chat flow, and a writer defect in this repository cannot be the cause of the S1 digest instability.

That is a bounded negative claim, not an exoneration. It says what this source tree does. It does not say what other software on the machine did.

### 1.1 What the chat path actually does

| Stage | Location | Mechanism | Writes? |
|---|---|---|---|
| Choose a file | `app-model.go` L847–L862 `OpenModelFilePicker` | native dialog returns a path string | **no** |
| Models dir | `app-model.go` L839–L843 `GetModelsDir` | `MkdirAll(models/)` | directory only |
| Read metadata | `app-model.go` L864–L866 `LoadModelMetadata` | `gguf.NewReader().ReadProfile(path)` | **no** |
| Register / classify | `app-model-registry.go` L320–L327, `app-model-discovery.go` L459, L496, L532 | extension switch over a directory listing | **no** |
| Parse GGUF | `internal/backend/gguf/` (`reader.go`, `admission.go`, `projector.go`, `pairing.go`, `kind.go`) | read-only parsing | **no** — the package contains no `os.Create`, `os.OpenFile` or `io.Copy` outside `testgguf_test.go` |
| Hash for the gate | `eval/harness/qualify.py::file_digest` | `path.open("rb")`, 1 MiB stream | **no** |
| Launch the loader | `eval/harness/cli.py` L381 (`--executable` default `H:\AI apps\koboldcpp.exe`), launch args include `--usemmap` | external process reads the file | **no** — and see §1.3 |

### 1.2 Consequence for the incident

The authorisation split is therefore:

| Question | Whose code | In this repository |
|---|---|---|
| Was a model file ever written by HOWL Chat? | — | **No code path exists** |
| Was it copied by HOWL Chat? | — | **No copy function is used on model paths**; `copyFile` exists only in `internal/audio/asr/whisper/recognizer.go` L395 for audio inputs |
| Was it memory-mapped during the observations? | KoboldCpp 1.120, launched with `--usemmap` | launch argument is in the harness; the mapping is the loader's |
| Was it rewritten while the hash gate ran? | unknown — no handle enumeration exists | not testable from source |

### 1.3 The mmap detail, stated precisely

`--usemmap` makes the *loader* map the file. A read-only mapping cannot change file contents. What it does change is **what a concurrent read can be served from**: mapped pages and buffered reads share the system cache for that file. Two processes therefore need not see the same bytes as the disk holds, and neither needs to have written anything. This is a cache-visibility mechanism, not a writer, and it is treated as such in `WINDOWS-11-STORAGE-EVIDENCE-MATRIX` (M-04, M-05).

No source in this repository writes to a mapped model region. `os.Create`/`io.Copy` appear only for audio recordings, log files, JSON state, downloads and test fixtures.

---

## 2. Writers that do exist in this repository

Four real write paths, none of them used on chat models. They are audited because they are the code that *would* write a large binary here, and because their patterns are the ones a future GGUF import would inherit.

| ID | Path | Destination strategy | Flush | Length check | Hash check |
|---|---|---|---|---|---|
| **W-A** | `internal/image/deps/download.go` L29 `Download` | `<name>.part` → `os.Rename` | **none** (`Close` only) | `received != ContentLength` when length known | **none** |
| **W-B** | `internal/audio/asr/whisper/downloader_acquire.go` L308 `downloadHTTPFile` | `<dest>.partial` → `os.Rename` | **none** | **none** | **none** |
| **W-C** | `internal/audio/asr/whisper/downloader.go` L64–L77 zip extraction | **direct to final path**, `O_TRUNC` | **none** | **none** | **none** |
| **W-D** | `internal/image/gallery/gallery.go` L123–L148 | `CreateTemp` → `Sync` → `Rename` | **`Sync()`** | — | **none** |

**W-D is the only writer in the repository that fsyncs before rename.** The inconsistency is itself a finding: the correct pattern exists in the tree and was not applied to the two downloaders.

---

## 3. Failure-mode evaluation

Each row is evaluated against W-A…W-D. `n/a` = the mechanism does not apply to that writer.

| Failure mode | W-A deps | W-B whisper HTTP | W-C whisper zip | W-D gallery |
|---|---|---|---|---|
| **Partial write / short write** | loop writes `buf[:n]`; a `Write` error aborts and removes `.part`. `os.File.Write` returns an error rather than a silent short write, so this is handled | `io.Copy` — same guarantee, error removes `.partial` | `io.Copy` error returns with **partial file left at the final path** | `CreateTemp` + copy; error leaves a `.tmp-*` file |
| **Offset / seek behaviour** | sequential append via `Write`; no `Seek` | sequential via `io.Copy` | sequential | sequential |
| **Append vs overwrite** | fresh `.part` (`O_TRUNC` on create) — never appends to a live file | fresh `.partial` | **`O_TRUNC` over the live destination** | fresh temp |
| **Buffer reuse / lifetime** | one 1 MiB `buf`, reused per read; no retention past the call | `io.Copy` internal 32 KiB buffer | `io.Copy` internal | `io.Copy` internal |
| **Async completion / sync** | none — the function returns only after `Rename` | none | none | none |
| **Flush / close** | `Close()` before the length check and rename; **no `Sync`** | `Close()` then rename; **no `Sync`** | `Close()` per entry; **no `Sync`** | **`Sync()` then `Close()`** |
| **Temp promotion / rename** | same directory → same volume; `Rename` is the commit point | same | **no temp at all** | same directory |
| **Destination collision** | Windows `os.Rename` to an existing file fails → an existing model is not silently replaced | same | **`O_TRUNC` silently destroys an existing file of that name** | rename over exists fails |
| **Concurrent writers** | two concurrent `Download` calls for the same artifact share `<name>.part` → interleaved writes to one temp. No lock, no `O_EXCL` | same, shared `.partial` | not guarded | `CreateTemp` is unique per call |
| **Cancellation / error propagation** | `ctx` cancels the body read; error path removes `.part` | 45-minute `ctx` timeout; removes `.partial` | error returns, **partial destination remains** | error leaves temp |
| **Success reported before durability** | **yes** — returns after `Rename`; data may still be in the page cache. A power loss immediately after can leave a renamed file whose contents are not on the medium | **yes**, same | **yes** | **no** — `Sync` precedes rename |
| **Unexpected path resolution** | `filepath.Join(dir, art.FileName)`; a `..` in `FileName` would escape `dir`. Artifacts are code constants in `deps.go`, so not reachable from input today | `dest` is constructed by callers from a fixed name | zip slip is explicitly checked (`strings.HasPrefix(path, Clean(targetDir)+sep)`) — **correct** | `dir` + base name |
| **Final size / identity verified** | size only, and only when `ContentLength > 0` | **no** | **no** | **no** |

### 3.1 The defects this produces

These are defects in the writers that exist. **None is a defect in the chat path, because the chat path does not write.**

| ID | Defect | Severity | Evidence class |
|---|---|---|---|
| **D-1** | **W-A accepts a truncated body when the server sends no length.** `total := resp.ContentLength`; when it is `-1` (chunked) the guard `if total > 0 && received != total` never fires, so an interrupted-but-clean-EOF transfer is renamed into place as a complete artifact | **High for artifact integrity**, offline-testable | Observed in source |
| **D-2** | **W-C truncates the live destination and writes into it.** An extraction failure leaves a partial file at the final path, indistinguishable by name from a complete one | **High**, offline-testable | Observed in source |
| **D-3** | **W-A and W-B rename without `Sync`.** On power loss the rename can be durable while the data is not: NTFS journals metadata, not file contents. The repository already contains the correct pattern (W-D) | **Medium**, power-loss only | Observed in source; mechanism documented (`FlushFileBuffers`) |
| **D-4** | **No writer verifies a digest after writing.** Every artifact is trusted on length alone (W-A) or on nothing (W-B, W-C). A corrupted transfer or a corrupted read is accepted | **Medium**, offline-testable | Observed in source |
| **D-5** | **Concurrent same-artifact downloads share one temp path** (W-A, W-B). No lock, no `O_EXCL`, no per-attempt suffix | **Low**, requires two concurrent calls | Observed in source |

None of D-1…D-5 can be reached by the chat-model flow, because that flow does not write. They matter (a) as defects in their own right and (b) as the inherited pattern if a model importer is ever built.

---

## 4. Writer audit matrix

| ID | Location | Observed implementation | Potential failure mechanism | Permanent wrong bytes? | Inconsistent visibility during concurrent writes? | Evidence needed to settle it | Severity · confidence | Offline testable |
|---|---|---|---|---|---|---|---|---|
| **W-01** | `app-model.go` L847–L866; `app-model-registry.go` L320–L327; `app-model-discovery.go` L459–L532 | path returned by dialog, registered, parsed | none — no write call exists | **no** | **no** | none needed; absence is source-visible | Info · **High** | n/a |
| **W-02** | `internal/backend/gguf/**` | read-only GGUF parsing | none | **no** | **no** | none | Info · **High** | n/a |
| **W-03** | `internal/image/deps/download.go` L29–L114 | 1 MiB loop → `.part` → length check → rename | truncated body accepted when `ContentLength < 0` (D-1) | **yes**, a short artifact is committed | shared `.part` (D-5) | an HTTP fixture with no `Content-Length` and a short body; assert no file at `dest` | **High · High** | **yes** |
| **W-04** | same, L108–L113 | rename as commit point | no `Sync`; power loss can leave rename durable, data not (D-3) | only after power loss | no | not testable without fault injection; documented mechanism only | Medium · Medium | no |
| **W-05** | `internal/audio/asr/whisper/downloader_acquire.go` L308–L341 | `io.Copy` → `.partial` → rename | no length check at all; no `Sync` | **yes**, on a short transfer | shared `.partial` | fixture with truncated body and `Content-Length` present; assert removal | **High · High** | **yes** |
| **W-06** | `internal/audio/asr/whisper/downloader.go` L64–L77 | `O_TRUNC` direct to final path | partial file left at final path (D-2) | **yes** | destination destroyed under a concurrent reader | fixture whose zip entry errors mid-copy; assert final path state | **High · High** | **yes** |
| **W-07** | `internal/image/gallery/gallery.go` L123–L148 | temp + `Sync` + rename | rename over existing fails (safe); temp leak on error | **no** | unique temp | file-count assertion after a forced failure | Low · High | **yes** |
| **W-08** | `internal/library/registry.go` L385, L413; `internal/image/registry/registry.go` L437, L459; `internal/image/scanner/manifest.go` L118 | JSON state: temp → rename | state written without `Sync`; a crash can lose the last update while the file it describes is intact | no (state, not payload) | last-writer-wins | forced-crash test | Low · High | partly |
| **W-09** | `internal/contentimport/journal.go` L296–L320; `internal/worldstore/import.go` L331–L425 | staged write with an intent journal, `O_EXCL` staging | journal entry durable before payload → recoverable by design | **no** by design | `O_EXCL` prevents collision | recovery tests already exist in the package | Info · High | yes (already covered) |
| **W-10** | `internal/audio/asr/whisper/recognizer.go` L395 `copyFile` | copy for audio inputs only | not applied to model paths | n/a | — | none | Info · **High** | n/a |
| **W-11** | `eval/harness/qualify.py::file_digest`, `cli.py` L381 | read-only hashing; external loader launch | hashing shares the system cache with any concurrent reader (cache visibility, not writing) | **no** | **yes** — same file seen differently by two processes is possible without any write | an unbuffered read of the same region compared to a buffered read | Info · High | no (live) |
| **W-12** | no writer anywhere | — | **a writer outside this repository** remains the live hypothesis | possible | possible | handle enumeration on the exact path | **Highest investigative value** | no |

---

## 5. Can the app report completion before the bytes are durable?

| Flow | Completion signal | Durable at that moment? |
|---|---|---|
| Chat model load | HTTP 200 from the loader | n/a — nothing written |
| W-A dependency download | `Rename` returned | **no** — contents may be cached |
| W-B whisper binary | `Rename` returned | **no** |
| W-C whisper zip | per-entry `Close` | **no** |
| W-D gallery image | `Sync` returned, then `Rename` | **yes** |

So for two of the three download paths, "downloaded" means "renamed", not "on the medium". On a clean shutdown this is invisible; across a power loss it is the difference between a complete artifact and a same-name file full of zeros or stale clusters.

---

## 6. Ranked writer-path questions

1. **Which process held the file?** Not answerable from source. Handle enumeration on the exact path is the single decisive check, and it is unauthorised here.
2. **Did the volume's owner write to it on a schedule?** `G:` is labelled "Second backup". A backup product restoring or re-writing payload while preserving timestamps would explain every observation without any hardware fault.
3. **Was any in-repo writer ever pointed at `G:\AI GGUF`?** No call site does; `deps.Download` targets the image-deps directory and the whisper downloader targets the binaries directory. Verifiable by reading the two call sites (`app-image-deps.go` L210, `NewBinaryDownloader`).
4. **Can a hash gate be trusted when it uses buffered reads?** Yes for detecting *change over time*, no for asserting *the bytes on the medium*. A buffered read can be satisfied from cache. This is the mechanism behind `M-04`/`M-05` in the Windows matrix, and it means the four digests prove inconsistent *reads*, not necessarily a corrupt file.
5. **If a GGUF importer is ever added, which pattern should it copy?** W-D (temp + `Sync` + rename + digest), not W-A.

---

## 7. Boundaries observed

No writer executed. No download issued. No model file opened, hashed or listed. No volume queried. No file created outside this document. No source modified — the five defects are reported, not repaired; each repair is a separate, narrowly scoped change.

**Establishes:** the repository has no chat-model write path; the four writers it does have, their exact failure modes, and five concrete defects among them.

**Does not establish:** what wrote to `G:`, whether anything did, or that the storage stack is at fault. The writer hypothesis and the storage hypothesis both remain open, and this audit moves neither.
