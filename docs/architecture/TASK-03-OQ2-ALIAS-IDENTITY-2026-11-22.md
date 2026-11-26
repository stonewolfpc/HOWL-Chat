# TASK-03 — OQ-2: alias paths and conversation identity

| Field | Value |
|---|---|
| Purpose | Decide whether two paths that reach one physical directory are one conversation location or two, and which identity may serve locating, deduplication, locking, reconciliation and ownership |
| Status | **Read-only investigation. Outcome: Clarify. V3.1 unedited. Nothing adopted, nothing authorized.** |
| Method | Source inspection of `paths.FileIdentity`, `resolvePhysical`, `CheckContained`, the conversation ID path and `SaveChat`/`LoadChat`/`GetChatList`/`SaveChatAs`; one disposable fixture probe with a junction, respellings, a copy and a rename (`build/tmp/oq2probe`, log `build/tmp/oq2probe.log`) |
| Recovery token | `TASK-03-OQ2-ALIAS-IDENTITY-2026-11-22` |

Second of row 576's six baseline questions, after OQ-1 (row 577). Resolved under the owner's frame: **Amend**, **Clarify**, **Retain**. This record answers identity only; it does not design a lock, a scanner or a migration.

---

## 1. The question, and why "one or two" is the wrong frame

**ALIAS-01 — Scope.** OQ-2 asks whether two paths that reach the same physical directory are one conversation identity or two. It does not decide OQ-3b-b's terminator rule, the row 550 collision's reporting, repair or retry.

**ALIAS-02 — The practical question is narrower than an architecture choice.** A scan enumerates entries; a lock names a protected directory; a verdict names a conversation the user recognises. Those three consumers can legitimately need different answers for the same bytes. Forcing one identifier to serve all three is the failure mode this record is written to prevent.

**ALIAS-03 — The answer, stated once: aliases add neither a location nor a conversation.** Two paths that resolve to one physical directory are **one candidate location** and **one logical conversation**. They must not produce two verdicts and must not produce two locks. A **copy** is the opposite case: a second physical location that may claim the same logical ID, and it is neither a second alias nor one healthy conversation — it is an identity conflict to report, never to merge or repair.

**ALIAS-04 — Therefore the system needs at least three identities, plus the one row 574 already names.** Locator, physical, logical and operation are different questions with different stability and different failure modes. The rest of this record defines them and says which consumer may use which.

---

## 2. Observed: what the filesystem actually does

**ALIAS-05 — Method.** A disposable fixture created `base/conv` with a `meta.json`, a directory junction `alias` → `base/conv` (`mklink /J`), and a byte-identical copy `base/copy`. The probe read `GetFileInformationByHandle` identity, `os.SameFile`, `GetFinalPathNameByHandle`, `filepath.EvalSymlinks`, `GetShortPathNameW` and `paths.CheckContained`. It wrote only inside `build/tmp/oq2fixture`; no product file was touched.

**ALIAS-06 — Identity results (Observed):**

| Path | Volume | File index | Same as target |
|---|---|---|---|
| `base\conv` | 1515597286 | 3096224744614400 | — |
| `alias` (junction) | 1515597286 | 3096224744614400 | **yes** |
| `base\copy` | 1515597286 | 3096224744614403 | **no** |
| `BASE\CONV` (case) | 1515597286 | 3096224744614400 | **yes** |
| `base\conv\.` | 1515597286 | 3096224744614400 | **yes** |
| `\\?\C:\…\base\conv` | 1515597286 | 3096224744614400 | **yes** |
| 8.3 short name | 1515597286 | 3096224744614400 | **yes** |

**ALIAS-07 — String equality fails on every one of those aliases.** `filepath.Clean(conv) == filepath.Clean(alias)` is false; `conv == strings.ToUpper(conv)` is false. Identity equality is true for all of them. A map keyed by path string therefore creates a second entry for each spelling, while a map keyed by `(volume, index)` collapses them correctly.

**ALIAS-08 — The standard library does not resolve a junction, and the existing helper knows it.** `filepath.EvalSymlinks(alias)` returned the alias path unchanged; `GetFinalPathNameByHandle(alias)` returned the target path. This confirms the source comment in `internal/paths/identity_windows.go` and means any implementation that uses `EvalSymlinks` as its resolver will treat the alias as a distinct location.

**ALIAS-09 — Physical identity survives rename; a junction does not.** After renaming `base\conv` → `base\conv-renamed`, the volume and file index were unchanged, so identity is stable across rename. The junction `alias` then failed with `ERROR_FILE_NOT_FOUND` and `IsNotExist = true`: a junction stores a **path**, not a file ID, so renaming the target breaks the alias. A dangling alias is not a conversation and must never become `Missing` for the target (see OQ-1, row 577).

**ALIAS-10 — A copy is a different physical directory with the same content.** Same volume, different file index, `SameFile = false`. That is the correct physical answer and the wrong logical answer if the copied `meta.json` still carries the source's ID.

**ALIAS-11 — PathService containment is physical, and the probe confirms it.** `paths.CheckContained(base, alias)` returned `nil` even though `alias` is lexically outside `base` and physically inside it. Containment through a junction already works by resolved physical path (row 519).

---

## 3. Does PathService already supply physical identity?

**ALIAS-12 — The primitive exists and is correct.** `internal/paths/identity_windows.go` defines an exported `FileIdentity{Volume, Index}` and an unexported `identityOf` that opens the directory with `FILE_FLAG_BACKUP_SEMANTICS` and reads `GetFileInformationByHandle`. `resolvePhysical` uses `GetFinalPathNameByHandle` because `EvalSymlinks` does not follow junctions.

**ALIAS-13 — But `identityOf` has no caller, and no exported function returns or compares a `FileIdentity`.** The type is to a `FileIdentity` what a declared struct is to a service: the capability is defined and unused. `CheckContained` resolves physical paths internally and returns only an error, so it does not hand the identity back.

**ALIAS-14 — Therefore PathService supplies containment, not identity.** The scanner cannot ask "are these two candidates the same directory?" and the lock cannot ask "what is this conversation's stable key?" through the current exported surface. The gap is an API gap, not a missing algorithm; the correct algorithm is already written and tested for containment.

**ALIAS-15 — `CheckContained` is deliberately a pre-flight check, not a TOCTOU guarantee.** Its own comment says a reparse point created after the call is not detected and that a write must use handle-relative operations preserving the validated identity. Identity for locking has the same shape: acquiring a lock named by a path is not the same as holding a lock on the directory that was verified.

---

## 4. Four identities, and what each is for

**ALIAS-16 — The table is the answer to "which identity serves which purpose":**

| Identity | Form | Stable across | Changes on | May be used for |
|---|---|---|---|---|
| **Locator** | a path string | nothing | case, spelling, `\\?\`, short name, junction, remap, rename | opening bytes; never identity, dedup or a lock key |
| **Physical** | `(volume serial, file index)` | rename, respelling, junction/symlink, mount alias, drive-letter remap | copy, restore to a new file, some network providers | one-candidate dedup, lock key, containment |
| **Logical** | `meta.ID`, else the directory name | rename, remap, copy (which is why copies conflict) | nothing by design; a manual copy wrongly inherits it | verdict attachment, ownership, user-visible identity, retry binding |
| **Operation** | row 574's durable key | — | — | retry only; never conversation identity |

**ALIAS-17 — The locator is not an identity, and the probe is the proof.** Every alias above reaches the same bytes and none is string-equal to another. Any equality test, map key, lock file name or verdict key derived from the path string is wrong for all of them.

**ALIAS-18 — Physical identity is a value with a validity guard.** It is usable only when `Volume != 0` and `Index != 0`. If two distinct resolved paths yield the same identity, the scan must **not** deduplicate them; it reports both and treats the identity as unreliable. Some SMB providers return zero or unstable file IDs; that is **Unestablished** here and the guard is fail-closed.

**ALIAS-19 — Logical identity is not the directory name, and the current code already conflates them.** `ChatMeta.ID` is the logical ID; `NewChat` sets it to the same timestamp string used as the directory name, `SaveChat` builds the directory from `Meta.ID`, `RenameChat` changes only `Name`, and `LoadChat` falls back to the directory name only when `meta.ID` is empty (`persistence.go:130-170, 719-729, 260-280`). No code enforces that the two agree.

**ALIAS-20 — And the list/load round trip is broken when they disagree (Observed).** `GetChatList` enumerates by `entry.Name()` but returns `ChatSummary.ID = meta.ID` (`persistence.go:502, 522`); `LoadChat(chatID)` builds the directory from its argument (`:260`). A directory renamed outside the app — or a manual copy whose `meta.json` still names the source — is listed with an ID that `LoadChat` cannot find. A junction alias inside `chats/` whose target carries a different `meta.ID` can make two listed rows resolve to one directory. Neither is hypothetical for a filesystem-visible store.

**ALIAS-21 — A copy is not an alias, and `SaveChatAs` already gets this right.** `SaveChatAs` calls `NewChat`, which generates a new timestamp ID, then copies the data but not the ID (`persistence.go:437-468`). A manual filesystem copy does not: it preserves `meta.ID`, producing two physical directories claiming one logical conversation. Reconciliation must report both and never merge, rename or delete either; the reporting shape belongs to the row 550 limitation, not to this record.

**ALIAS-22 — Message identity and operation identity stay separate.** Message IDs are positional (`O7-199`…`O7-200`); operation identity is row 574's durable key. Neither becomes a conversation key through this record.

---

## 5. Stability, including the cases that cannot be observed

**ALIAS-23 — Stability table for the events the owner named:**

| Event | Physical | Logical | Locator | Consequence |
|---|---|---|---|---|
| rename | same (Observed) | same (`meta.ID`) | old path dead | identity stable; string key breaks |
| case / short name / `\\?\` / trailing dot | same (Observed) | same | different string | dedup must be by identity |
| junction / symlink alias | same (Observed for junction) | same | different path | one location, one lock |
| `SaveChatAs` copy | new | **new** (by design) | new | correct: a copy is a new conversation |
| manual copy | new | **same** (inherited) | new | conflict to report, never merge |
| restore from backup | likely new | same | new | same conversation, new physical identity |
| drive-letter remap | same volume serial + index (Source-derived) | same | new letter | physical identity portable, locator not |
| root unavailable | unknown | known only from the index | unreachable | no dedup, no lock, no Missing (row 577) |

**ALIAS-24 — An unavailable root makes physical identity unknown, not empty.** No identity means no deduplication and no lock; it never means "these are different" and never means "this is absent". Row 577 already fixes the verdict side; this record fixes only the identity side.

**ALIAS-25 — Network providers are the open edge.** Whether a mapped drive and its UNC path return the same `(volume, index)`, whether an SMB server returns a stable file index, and whether a reconnect preserves it are all **Unestablished** — no share was available to test. The validity guard of ALIAS-18 is the contract's protection: if the provider cannot give a nonzero, collision-free identity, the scan does not dedup and the lock refuses rather than guessing.

---

## 6. The rules this implies

**ALIAS-26 — One candidate per physical directory.** A scan candidate is a physical conversation directory, and candidates are deduplicated by `FileIdentity` when it is computable. Path spelling, case, 8.3 names, junctions, symlinks and mount aliases never create a second candidate.

**ALIAS-27 — A verdict attaches to the logical identity, never to the locator.** The conversation the user sees is `meta.ID`; the directory name is how the bytes are reached. Deduplication may collapse two paths, but it may never collapse two different logical IDs into one verdict.

**ALIAS-28 — Two physical directories with one logical ID are reported, not merged.** Both are listed with their own physical evidence; neither is preferred, renamed, deleted or declared the original. The reporting shape is the row 550 limitation's job. The one rule this record states is negative: **dedup by physical identity, never by logical ID**, so the conflict cannot be silently erased.

**ALIAS-29 — A non-candidate gets no verdict.** A dangling junction, a directory without `meta.json`, or an alias entry is not a conversation; it must not produce `Missing` for a target that exists elsewhere (row 577, row 573's candidacy rule).

**ALIAS-30 — A lock key is physical identity when computable, and the lock refuses otherwise.** Two aliases must not obtain two locks for one directory, so a path-string lock is disqualified. A logical-ID lock is also insufficient: an alias and its target share the logical ID only if `meta.ID` agrees, and a manual copy would then share one lock across two directories. The key is `(volume, index)`; identity is verified before acquisition and again after it, and a mismatch or an uncomputable identity is a typed refusal, never a fallback to the path.

**ALIAS-31 — Ownership and retry binding use the logical identity.** A physical identity is not durable across restore or copy; ownership of a conversation must survive both. Row 575's conversation binding is therefore `meta.ID`, not a file ID and not a locator.

---

## 7. Outcome

**ALIAS-32 — Outcome: Clarify.** V3.1's rows are sound; what is missing is the meaning of "conversation directory" and of a scan candidate. Two paths to one directory do not contradict any row — they contradict the unstated assumption that one path equals one candidate. The correction is a term and a candidate rule, not a new state, a new reason or a changed decision row.

**ALIAS-33 — The bounded amendment:**

| Clause | Change |
|---|---|
| §1 terms | add **locator**, **physical identity** (with the nonzero/collision guard) and **logical identity** (`meta.ID`, else directory name) |
| scan boundary | one candidate per physical directory; dedup by physical identity when computable; never by path string and never by logical ID |
| verdict rule | attach to logical identity; a non-candidate gets no verdict |
| ownership/lock note | physical identity for the lock key and dedup; logical identity for ownership and retry binding |
| `V3-73` | OQ-2 closes, subject to owner adoption |

No row in `V3-39` changes. `V3-58`'s reason vocabulary is unchanged. The row 550 limitation still owns how a duplicated logical ID is reported.

**ALIAS-34 — What this record does not decide.** It does not design the lock file layout (row 567 owns placement), does not decide whether PathService should export identity (an implementation choice), does not add a duplicate-identity reason, does not define restore semantics, and does not choose a storage format. It says which identity each consumer may trust.

**ALIAS-35 — This is not an architecture preference.** Nothing here argues for path-based or handle-based storage. The probe shows the filesystem gives a usable physical identity on a local NTFS volume and leaves it unverified on a network provider; the contract takes the identity where it is trustworthy and refuses where it is not. If a later design replaces the directory store entirely, these distinctions still apply to whatever replaces it.

**ALIAS-36 — Nothing here authorizes the documentation to become a second product.** The deliverable is one term set and one candidate rule. It is not a refactor plan, not a scanner design and not a reason to keep the current layout.

---

## 8. Tests this would require

- **ALIAS-T1** A junction alias and its target produce one candidate and one verdict, keyed by the same `(volume, index)`.
- **ALIAS-T2** Case, trailing-dot, 8.3 and `\\?\` spellings of one directory produce one candidate; a path-string map produces five.
- **ALIAS-T3** `SaveChatAs` produces a distinct physical identity and a distinct logical ID; a manual copy produces a distinct physical identity and the same logical ID, reported as two locations, not merged.
- **ALIAS-T4** Renaming a conversation directory preserves its physical identity and its lock key; a path-keyed lock changes key and is rejected by the test.
- **ALIAS-T5** A dangling junction yields no conversation verdict and does not make the renamed target `Missing`.
- **ALIAS-T6** With identity uncomputable (root unavailable, or zero IDs injected), the scan does not dedup, does not lock, and does not report `Missing`.
- **ALIAS-T7** Two different paths yielding the same `(volume, index)` are not deduplicated; both are reported and the identity is marked unreliable.
- **ALIAS-T8** `GetChatList` and `LoadChat` agree on identity: a directory whose `meta.ID` differs from its name is not listed under an ID that cannot be loaded.
- **ALIAS-T9** The negative assertion: no dedup map, lock-file name or verdict key is derived from a path string or from `meta.ID` where physical identity is the question.

---

## 9. Nothing here is authorized

**No implementation follows.** This record authorizes no exported identity API, no scanner, no dedup pass, no lock key, no change to `paths.FileIdentity`, `CheckContained`, `LoadChat`, `GetChatList`, `SaveChatAs`, `RenameChat` or any directory layout; no migration, no rename, no repair of a duplicated logical ID, no new reason, no deletion enablement and no retry policy. It clarifies what identity means for three consumers and leaves every consumer unimplemented.
