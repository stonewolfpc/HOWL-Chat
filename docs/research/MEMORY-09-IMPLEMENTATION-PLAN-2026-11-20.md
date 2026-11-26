# MEMORY-09 — Implementation Plan

## Bounded L-2/P-1 implementation, with a class-source stop before auto-capture

**Status:** implementation **plan**, not execution. Implementation is authorized within IAR-04/IAR-14, but **no code,
schema, migration, test or data change has been made**. The plan identifies a **blocking class-source decision**:
`MemoryRecord` has no `knowledge_class`, and no classifier rule supplies one. The plan does not invent it.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Authorization | IAR approved (IAA-01) |
| Storage | JSON files under `Chats/<chatID>/`; **no SQL migrations** |
| New file | `Chats/<chatID>/memory_lifecycle.json` (version 1) |
| Primary blocker | `knowledge_class` has no source in the current code; stop before auto-capture |
| Tests now | **0**; baseline run is recorded |
| Scope | L-2, P-1 integration, three guardrails, S-3 rebuild/integrity, status-slot exclusion |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Pre-change baseline

**IPL-01 (Observed — focused test baseline).** `go test ./internal/backend/memory .` passed before any change:

```
ok  howl-chat/internal/backend/memory  (cached)
ok  howl-chat                          9.248s
```

**IPL-02 (Observed — repository baseline).** `git status --porcelain` shows the same **6** pre-existing modified files
(`config/logs/startup.log`, `config/model-registry.json`, `docs/DOCUMENT-REGISTRY.md`,
`docs/research/OWNER-ENVIRONMENT-HISTORY-2026-11-14.md`, `eval/harness/__pycache__/validate.cpython-312.pyc`,
`eval/harness/validate.py`). They are not touched.

**IPL-03 (Observed — map/canon baseline).** The conformance map is **118**; r33/v9/ledger/map hashes are unchanged
(`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`).

**IPL-04 (Observed — no SQL).** `glob **/*.sql` finds **0** files. Persistence is JSON:
`Chats/<chatID>/memories.json`, `candidates.json`, `auto_capture_config.json`, `continuity_config.json`.

**IPL-05 (Observed — the current record).** `internal/backend/types/memory.go` `MemoryRecord` has `Scope`, `Status`,
`OwnerID`, `OwnerKind`, but **no `KnowledgeClass`**. `Status` is documented as `active/archived`; `NewMemoryRecord`
defaults it to `"active"`; `AddMemory` and `ApproveCandidate` also default to `"active"`.

**IPL-06 (Observed — no class source).** `internal/backend/memory/classify.go` assigns scope, owner, title, summary,
content, tags and `pending` status. It does **not** assign a knowledge class. A repository search finds no
`knowledge_class` / `KnowledgeClass` in code. The class baseline has **no source** today.

---

## 2. Scope and phases

**IPL-07 (Calculated — the plan applies adopted decisions).** This is not a new design round. L-2, P-1 and S-3 are
inputs. The plan implements them within the IAR scope.

**IPL-08 (Calculated — phases).**

| Phase | Work | Gate |
|---|---|---|
| **0** | Decide the **knowledge-class source** for new records | **blocking owner decision**; no auto-capture implementation before it |
| **1** | Lifecycle types, store, formation entry, status slot | scope/class values supplied; status slot explicit |
| **2** | P-1 scope/class changes, independent versions, projection | no cross-dimension ordering |
| **3** | S-3 rebuild/integrity, status exclusion, inspection | focused tests pass |
| **4** | Legacy handling/test baseline | existing records reported, never silently repaired |

**IPL-09 (Calculated — why Phase 0 blocks).** L-2 requires a `class_baseline`. The current code cannot supply one:
`MemoryRecord` has no field, the DTO has no field, and the classifier has no rule. Adding the field is mechanical;
choosing its value for auto-captured events is a **semantic decision** (event-type mapping, default class, or author
input) that this plan does not make.

**IPL-10 (Calculated — the minimum unblocked work).** The lifecycle framework, scope baseline, status slot, P-1
machinery and tests can be built with explicit class values in tests. Auto-capture and any caller that cannot supply a
class **stop** at validation. Full L-2/P-1 acceptance for class waits for Phase 0.

---

## 3. File-by-file plan

**IPL-11 (Observed — exact files).**

| File | Change | Minimum |
|---|---|---|
| `internal/backend/types/memory.go` | add `KnowledgeClass string` to `MemoryRecord` and `MemoryCandidate`; no status-vocabulary change | yes |
| `internal/backend/memory/lifecycle.go` | **new**: lifecycle types, JSON store, formation append, P-1 changes, projection, check, inspection errors | yes |
| `internal/backend/memory/lifecycle_test.go` | **new**: formation/P-1/status/rebuild tests | yes |
| `internal/backend/memory/service.go` | `AddMemory` appends a formation entry; atomic-ish save order; lifecycle diagnostics | yes |
| `internal/backend/memory/candidates.go` | `ApproveCandidate` appends a formation entry; class required | yes |
| `app-memory.go` | DTO `KnowledgeClass`; lifecycle status inspection field; conversion | yes |
| `internal/backend/memory/classify.go` | **no change in this plan** — class source is Phase 0 | stop |
| `internal/backend/memory/continuity.go`, `dedup.go`, `cluster.go`, `converter.go` | **no change** — legacy paths; they do not create lifecycle entries | out of scope |
| `internal/backend/chat/persistence.go` | **no change** — second writer to `memories.json`; identified risk | out of scope |
| `frontend/wailsjs/go/models.ts` | regenerate if the DTO changes; do not hand-edit generated bindings | conditional |

**IPL-12 (Calculated — why these files).** `service.go` and `candidates.go` are the two record-creation paths.
`app-memory.go` is the frontend boundary. `types/memory.go` carries the class field. `lifecycle.go` is the new
authority and check. No other file needs a change for the minimum.

---

## 4. Storage and migration

**IPL-13 (Calculated — new lifecycle file).** A new `Chats/<chatID>/memory_lifecycle.json` holds lifecycle entries.
`memories.json` remains the record cache. No SQL migration is needed.

**IPL-14 (Calculated — proposed shape, not applied).**

```json
{
  "version": 1,
  "formations": [
    {
      "memory_id": "mem_...",
      "entry_id": "life_...",
      "scope_baseline": "world",
      "class_baseline": "experiential",
      "status_slot": { "state": "unresolved" },
      "provenance_ref": ""
    }
  ],
  "scope_changes": [
    { "memory_id": "mem_...", "entry_id": "life_...", "version": 1, "value": "scenario" }
  ],
  "class_changes": [
    { "memory_id": "mem_...", "entry_id": "life_...", "version": 1, "value": "belief" }
  ]
}
```

The arrays are per-dimension; there is **no global order array** and no shared cross-stream sequence.

**IPL-15 (Calculated — status slot).** The slot is a tagged object, not a string:

| Stored form | Meaning |
|---|---|
| `{"state":"unresolved"}` | explicit unresolved; no value |
| `{"state":"resolved","value":"..."}` | a future N-b value; current code never writes it |
| missing `status_slot` | **absent** — validation error, not unresolved |
| unknown `state` / missing `value` | **invalid** — validation error |
| malformed JSON / unreadable file | **corrupted** — load error |

The four conditions are represented distinctly; they must not collapse into `null`, `""`, or `active`.

**IPL-16 (Calculated — migration).** New records write a formation entry. Existing `memories.json` records are **not
auto-backfilled**: they have no `knowledge_class`, and their `Status` is the legacy `active/archived` field, not the
lifecycle status. Backfilling would invent both a class and a status; the plan **stops** and returns for a migration
decision.

**IPL-17 (Calculated — legacy reporting).** An existing record without a formation entry is reported as
**incomplete/unmigrated** by the lifecycle check, not silently treated as valid and not auto-repaired. This is the
S-3 trust rule applied to legacy data.

---

## 5. Write order and integrity

**IPL-18 (Calculated — formation first, record second).** `AddMemory`/`ApproveCandidate` validate the class and status
slot, append the formation entry (temp-file + rename), then save the record. If the record write fails, the formation
entry remains and the record is absent; inspection reports an orphan formation. A retry can write the record. The
reverse order would leave a record without a baseline; the chosen order keeps the stream authoritative.

**IPL-19 (Calculated — no silent repair).** A missing formation, missing record, gap, duplicate, mismatch, absent slot,
invalid slot or corrupt file is surfaced. None is repaired by inference.

**IPL-20 (Calculated — concurrency).** A service-level mutex serializes lifecycle writes for the active chat; same-
dimension versions are assigned under it. Cross-dimension versions remain independent. This is the minimum mechanism;
the invariant is what matters, not the lock choice.

**IPL-21 (Calculated — atomic write).** `memory_lifecycle.json` uses `write temp + rename`, matching
`PersistenceManager.saveJSON`. `memories.json` keeps its current writer; the plan does not change it.

---

## 6. P-1 integration

**IPL-22 (Calculated — change API).** `AppendScopeChange(memoryID, value)` and `AppendClassChange(memoryID, value)`
read the current maximum for that dimension, assign **max+1**, and append. The first change is **1**. A duplicate
`(memory, dimension, version)` is an integrity error.

**IPL-23 (Calculated — projection).**

- scope: `formation.scope_baseline`, then the highest `scope_change` version; if none, baseline.
- class: `formation.class_baseline`, then the highest `class_change` version; if none, baseline.
- status: the `status_slot`; no projection value is invented.

**IPL-24 (Calculated — no cross-dimension ordering).** Scope and class use separate arrays/sequences. No entry orders
a scope change against a class change. No global lifecycle order is created; N-c remains open.

**IPL-25 (Calculated — no acquisition dependency).** The lifecycle file carries its own values. `provenance_ref` is
optional, typed, and never read for values; no acquisition event ID is used as the lifecycle identity.

---

## 7. S-3 rebuild, integrity and inspection

**IPL-26 (Calculated — rebuild).** `Rebuild(memoryID)` returns the scope/class projection and the status-slot state
from `memory_lifecycle.json` alone. It does not read `memories.json`, the acquisition stream, or a stored join.

**IPL-27 (Calculated — check).** `Check(memoryID)` compares the rebuilt scope/class to the record cache and reports:

| Condition | Result |
|---|---|
| missing formation | incomplete/unmigrated |
| internal gap | incomplete/untrusted; integrity error if completeness known |
| out-of-order | pending; hold higher version |
| duplicate version | integrity error/conflict |
| record-cache mismatch | integrity error |
| missing latest | **undetectable by sequence**; limitation preserved |
| absent/invalid/corrupt slot | distinct validation/load errors |

**IPL-28 (Calculated — inspection).** A Go diagnostics API reports the slot as a state
(`unresolved`/`resolved`/`absent`/`invalid`/`corrupted`), never as a valid status value. It also reports the legacy
`Status` separately, clearly labeled, so operators cannot confuse it with the lifecycle slot.

**IPL-29 (Calculated — status exclusion).** No lifecycle projection, fold or check reads `MemoryRecord.Status` as the
lifecycle status. Tests enforce that. The legacy `Status` remains a UI/compat field; it is not the authority.

---

## 8. Tests

**IPL-30 (Calculated — new tests).**

| Test | What it proves |
|---|---|
| `TestFormationEntryCarriesScopeAndClass` | L-2 baseline values are in the entry |
| `TestScopeAndClassReconstructIndependently` | each dimension uses only its own baseline/changes |
| `TestNoCrossDimensionOrdering` | a scope change does not move class, and vice versa |
| `TestFormationEntryIsNotAChange` | guardrail 1; the fold starts from the formation entry |
| `TestStatusSlotUnresolvedRoundTrip` | `{"state":"unresolved"}` survives JSON without becoming null/active |
| `TestStatusSlotStatesAreDistinct` | unresolved vs absent vs invalid vs corrupted |
| `TestStatusSlotIsNotAStatusValue` | `provisional`/`active` are rejected as slot states |
| `TestRebuildFromLifecycleOnly` | S-3 rebuild reads only the lifecycle file |
| `TestMismatchAndDuplicateAreIntegrityErrors` | §6.4 behavior |
| `TestMissingLatestIsUndetectable` | limitation preserved; no false completeness |
| `TestKnowledgeClassMissingIsRejected` | Phase 0 stop; no invented class |
| `TestLegacyRecordWithoutFormationIsReported` | no silent backfill/repair |
| `TestProvenanceRefIsNotUsedForValues` | guardrail 3 |
| `TestNoAcquisitionDependency` | lifecycle rebuild without acquisition input |

**IPL-31 (Observed — DTO tests).** `app_memory_lifecycle_test.go` (root) can cover `KnowledgeClass` and the
lifecycle-status inspection field in `typesMemoryToRecord`. No frontend/UI test is in scope.

**IPL-32 (Calculated — focused verification).** After implementation:

```
go test ./internal/backend/memory .
go test .
```

A full `go test ./...` is optional/CI; the focused packages cover the changed paths. No browser, UI or performance
test is in scope.

---

## 9. Stop conditions

**IPL-33 (Observed — stop and return to the owner).**

| Trigger | Why |
|---|---|
| **class source is unresolved** | Phase 0; no invented event-type → class mapping |
| **backfill existing records** | requires class and legacy-status decisions |
| **treat legacy `Status` as lifecycle status** | AC-11 / status-slot exclusion |
| **use acquisition ordering or values** | OD-a separate structures |
| **create a global lifecycle order** | N-c; not authorized |
| **change continuity/dedup semantics** | outside this scope; separate decision |
| **change `SaveChat`'s `memories.json` ownership** | dual-writer risk; separate decision |
| **introduce SQL/DB schema** | not needed; stop if the JSON path is insufficient |
| **resolve N-b/N-c/authority** | excluded by IAR-10/IAA-07 |

---

## 10. Risks and out of scope

**IPL-34 (Observed — known risks).**

| Risk | Handling |
|---|---|
| **dual writer** — `chat.PersistenceManager.SaveChat` writes `memories.json` from `ChatData.Memories` | the lifecycle file is separate; a missing record is reported, not repaired; changing the dual writer is a separate decision |
| **legacy records** — no class, legacy status | reported as unmigrated; no backfill |
| **continuity/dedup** — merge/move records without lifecycle entries | out of scope; tests must not claim they are lifecycle-complete |
| **JSON atomicity** — two files cannot be one transaction | formation-first order; orphan/missing states are visible |
| **missing latest** | documented; Decision B |

**IPL-35 (Observed — out of scope).** Decision B, N-b, N-c, formation authority, cross-stream ordering, Howl/UX,
requirement capture, canonical capture, deployment, UI work.

---

## 11. Traceability

**IPL-36 (Observed).**

| What | Where |
|---|---|
| Authorization | IAA-01; IAR-04/14 |
| L-2 | L2D-01…L2D-19 |
| P-1 | P1D-04…P1D-11 |
| S-3 | S3A-27 |
| Guardrails | LCR-27; IAR-07 |
| Status safeguard | IAA-03; AC-11 |
| Current record model | `internal/backend/types/memory.go` |
| Current storage | `internal/backend/memory/service.go` |
| Class gap | `classify.go`; no `KnowledgeClass` in code |
| Baseline tests | `go test ./internal/backend/memory .` passed |

---

## 12. Record status

**IPL-37 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — map **118** |
| Implementation started | **0** — plan only |
| Tests run for the plan | baseline `go test ./internal/backend/memory .` passed |
| New implementation tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **implementation plan** — no canonical capture |

**End of plan.** The minimum implementation is a lifecycle-side formation entry plus P-1 changes in a new
`memory_lifecycle.json`, with the status slot as a tagged unresolved object and S-3 rebuild/check over the lifecycle
file alone (IPL-13…IPL-29). **The class source is a blocking Phase 0 decision**: `MemoryRecord` has no
`knowledge_class` and the classifier supplies none; the plan does not invent one (IPL-05/06/09/33). Scope and status
can be implemented now; class acceptance and auto-capture stop until the owner decides the source. Legacy records are
reported, not backfilled or repaired (IPL-16/17). **No code, schema, migration or test has been written by this plan**
(IPL-37).
