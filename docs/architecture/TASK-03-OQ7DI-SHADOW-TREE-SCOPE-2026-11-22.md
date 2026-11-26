# TASK-03 — OQ-7d-i: the scope of the installation-tree store

| Field | Value |
|---|---|
| Purpose | What `StateRoot()/Chats` actually is — legacy store, active data with its own lifecycle, test-only artifact, or a mixture — and what that makes the reconciliation boundary |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Source inspection of every reader and writer of that namespace, the declared root model, and the wipe and delete paths |
| Recovery token | `TASK-03-OQ7DI-SHADOW-TREE-SCOPE-2026-11-22` |

Continues OQ-7's numbering (`O7-172`…). Resolves OQ-7d-i as opened by `O7-79` (row 566), and corrects two of that row's claims (§1).

---

## 1. Corrections to row 566

**O7-172 — The path base is `StateRoot()`, not `ProjectRoot()`.** In production the two are equal (`isolation.go:185-190` returns `ProjectRoot()` when not isolated), so the observed layout is the same. Under isolation they are not: the namespace becomes `<isolation>/data/Chats`, while the conversation store becomes `<isolation>/data/persisted-chats`. So the two trees remain disjoint in both modes, but for a different reason than row 566 stated, and an isolated run cannot see production memories at all.

**O7-173 — Withdrawn: `chat.json` has a reader.** Row 566 listed "export path (`contentItemFile`, `library.Chats`)". `ExportContentItem` returns an error for `library.Chats` **before** reaching `contentItemFile` (`app-libraries.go:230-231`), and no other code reads `chat.json`. It is **write-only**: created by the test seam and read by nothing.

---

## 2. What is in the namespace

**O7-174 — `StateRoot()/Chats` is not a store; it is a shared directory namespace with four populations and four different reachability classes:**

| Population | Files | Written live? | Read live? | Class |
|---|---|---|---|---|
| user-authored memories | `memories.json`, `memories_backup.json` (`memory/service.go:451-472`) | **yes** — `AddMemory`, `EditMemory`, `DeleteMemory`, `ImportCharacterMemories` | **yes** — `GetMemories`, `GetMemoryDetail` via `App` bindings | **active product data** |
| settings | `auto_capture_config.json` (`:265-300`) | **yes** — `App.SetAutoCaptureConfig` | **yes** — `App.GetAutoCaptureConfig` | active product data |
| settings, writer missing | `continuity_config.json` (`continuity.go:155-172`) | **no caller found** | **yes** — `LoadContinuityConfig` | dormant writer, live reader |
| test-seam artifacts | `history.json`, `chat.json`, `memories_snapshot.json` (`chat/service.go:200, 518, 536`) | **seam only** — `NativeTestSaveChatData` | `memories_snapshot` only | **test-only** |
| cache | `candidates.json` (`candidates.go:59-79`) | reachable via `ProcessMessageForEvents` | yes | production caller **Unestablished** |

**O7-175 — Verdict: a mixture, and the mixture is the finding.** All four of the owner's candidate answers are true of *part* of the namespace, and none is true of it as a whole. That is exactly why it cannot be treated as one store with one lifecycle.

---

## 3. Which tree does the product recognize?

**O7-176 — Conversations are recognised only in the history root.** `LoadChat`, `GetChatList`, `GetActiveChatMeta` and `DeleteChat` all operate on `<rootDir>/chats` (`persistence.go:257, 485, 550, 474`). **No conversation is ever opened from `StateRoot()/Chats`** — not by any binding, not by any loader, not by the frontend. The namespace has no `meta.json` anywhere in its writers, so nothing in it is even shaped like a conversation.

**O7-177 — The declared root model does not locate the live conversation store either.** `paths.Roots()` classifies `RootApplicationData` as `%AppData%\Howl` (`roots.go:157-160`), while the persistence manager actually writes `%AppData%\Roaming\HOWLChat` (`persistence.go:55-72`). The same model labels `RootInstallation` (= `ProjectRoot()`) `Writable: false`, class `installation`, described as "read-only resources" — yet the memory package writes beneath it in production (O7-174). `roots.go` states this plainly: it is additive and "nothing here changes where the running application reads or writes a byte" until ADR-001 §12 is authorized.

**O7-178 — Therefore the root model is a target, not an authority.** Reconciliation must not derive the conversation boundary from `Roots()`, and must not read its app-data root as the place conversations live. The only correct derivation remains `pm.RootDir()/chats`, as V3.1 has it.

**O7-179 — And the installation root is declared non-writable while being written.** Recorded as evidence, not as a defect to repair here: live memory writes go under a root the PathService labels read-only. Any future enforcement of that label would break the memories store, which is why the two cannot be resolved separately.

---

## 4. What happens when both trees hold data

**O7-180 — Nothing.** There is no merge, no precedence rule, no cross-tree comparison, and no detection anywhere. `GetChatList` enumerates one directory. The consequence is asymmetric and it matters for the contract:

| Situation | Scanner | Application |
|---|---|---|
| both trees hold `<id>/memories.json` | sees one (`<rootDir>/chats`) | reads and writes the **other** |
| shadow tree holds `history.json` (seam ran) | sees nothing there | never reads it |
| shadow tree holds an orphan for a deleted conversation | sees nothing | may still read its memories |

**O7-181 — So a `Healthy` conversation verdict says nothing about that conversation's memories.** The scanner's verdict covers the history root's directory; the memories the user actually edits live in a namespace the scanner never inspects, under a root that moves with launch location (`OQ-7c-i` `O7-117`, `O7-122`). A conversation can scan `Healthy` while its shadow memories were orphaned by launching the app from a different directory.

**O7-182 — This is the second memories store, not a second history store, and the contract must say so.** `O7-65` already corrected the schema divergence (`[]interface{}` from `ChatData.Memories` versus `[]types.MemoryRecord`). OQ-7d-i adds the lifecycle: the two stores have **different lifetimes** (§5), so neither can be derived from the other and neither may be reported as the other.

---

## 5. Lifecycle: the namespace outlives the conversation

**O7-183 — Wipe removes one tree.** `WipeAllData` stages `<RootDir>/chats` and the gallery (`app-wipe.go:45-52`). Every file in `StateRoot()/Chats` survives — memories included. So **"wipe all data" does not remove user-authored memories.** Whether that is the intended product behaviour is not this record's call; that it is the current behaviour is established.

**O7-184 — And the memory scope is not reset by the wipe.** `WipeAllData` calls `a.memoryService.SetChatID("")` (`app-wipe.go:78`), and `SetChatID` **returns early on an empty value** (`memory/service.go:33-38`). The service therefore keeps the last chat ID, and the next memory write recreates `Chats/<that id>/`. A wiped conversation's memory directory can be recreated after the wipe.

**O7-185 — Delete leaves the namespace untouched.** `DeleteChat` removes `<rootDir>/chats/<id>` only (`persistence.go:471-483`). A deleted conversation's shadow directory persists indefinitely, keyed by an ID that no longer names any conversation.

**O7-186 — So the namespace accumulates orphans by construction:** one class created by every deletion, one by every wipe, and one by every launch from a different directory. Orphans are not conversations, and their presence is not evidence that a conversation exists.

---

## 6. What the reconciliation boundary therefore is

**O7-187 — The boundary stays exactly where V3.1 put it: `<rootDir>/chats`, history root only.** OQ-7d-i changes nothing about the scan. It closes the *question*, not the boundary.

**O7-188 — Three declarations the contract must carry, so the boundary is not mistaken for a wider guarantee:**

| Declaration | Content |
|---|---|
| scope | a conversation verdict is a statement about `<rootDir>/chats/<id>` — its history and its app-data files — and about nothing else |
| exclusion | `StateRoot()/Chats` is outside the boundary, whatever it contains, including memories |
| prohibition | a shadow directory must never be promoted into the scan: it carries no `meta.json`, so it is not a conversation and cannot be one |

**O7-189 — Candidacy is decided by `meta.json`, which is what makes the exclusion safe rather than arbitrary.** Every shadow artifact is nameless in conversation terms: no identity file, only content and settings. Requiring an identity document for candidacy is the structural guard that keeps the namespace out of the scan without a path rule that could drift.

**O7-190 — Orphans must not be reported at all.** A shadow directory for a deleted or wiped conversation is not `Missing`, not `Unavailable`, not `Pending` — it is not a conversation, so it has no verdict. Reporting it would be inventing a conversation from an artifact, which is the same error class as promoting a scan error to `Missing`.

**O7-191 — The boundary's silence must be visible.** Because the exclusion is real (O7-181), an operator reading a `Healthy` verdict must be able to see that memories were not assessed. A verdict that does not say what it did not cover invites exactly the inference `O7-78` forbade.

**O7-192 — Nothing here is authorized.** No path consolidation, no memories-store unification, no relocation of the shadow namespace, no wipe or delete change, no `SetChatID` change, no `Roots()` change, no boundary extension, no migration, no backfill, no deletion enablement. This record establishes what the second namespace is and keeps it outside the contract.

---

## 7. Consequences for the rest of the branch

**O7-193 — OQ-7d-ii (retry identity) must be designed against one conversation store plus one excluded namespace, not two conversation stores.** That is the value of resolving scope first: the retry question has a single content domain to be idempotent over. The shadow namespace's writes are settings and memories, and their retry semantics are not governed by the history contract.

**O7-194 — The writer population is confirmed as spanning two roots and four packages** (row 566's standing conclusion), but only one of those roots contains conversations. Enforcement (OQ-7) therefore has a single-root target and does not need a cross-root lock.

**O7-195 — The launch dependence remains the sharpest live risk in this namespace**, and it is now concrete rather than theoretical: user-authored memories sit under a root that resolves from the executable's location or the working directory (`O7-117`), while conversations do not. Recorded here; owned by OQ-7c-i and unauthorized.

---

## 8. Tests this would require

- **O7-T55** A conversation directory and a shadow directory with the same ID coexist; `GetChatList` returns exactly one, and a scan reports exactly one.
- **O7-T56** A shadow directory containing `history.json` but no `meta.json` is never a scan candidate.
- **O7-T57** A conversation with shadow memories scans `Healthy`, and the verdict's scope statement names memories as unassessed.
- **O7-T58** After `WipeAllData`, the shadow namespace is unchanged and `SetChatID("")` is shown not to reset the memory scope.
- **O7-T59** After `DeleteChat`, the shadow directory persists and produces no verdict.
- **O7-T60** `Roots()`'s application-data path differs from `pm.RootDir()`, and the scanner uses the latter — asserted structurally.
- **O7-T61** An isolated run's shadow namespace is inside the isolation root and disjoint from the production one.
- **O7-T62** `chat.json` is written by the seam and read by no production path — asserted by reachability.
- **O7-T63** The same binary launched from two directories resolves two shadow namespaces and one conversation root, documenting `O7-195`.
