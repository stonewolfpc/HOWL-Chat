# FC-4 — Content Ownership and Synchronization Contract

## Contract paper `FC-4-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract paper — review draft |
| Authority | D3 owner decision; `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Gate | M0B/M5 |
| Status | **PROPOSED — bypass inventory incomplete** |
| Implementation authorization | **None** |
| Tests/models/process runs | **0** |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Observed** · **Unestablished**.

---

## 1. Contract decision

**FC4-01 (Owner-approved — D3).** Canonical library stores own authored objects and revisions. Conversations hold immutable references plus reproducibility snapshots. Browser local storage becomes a compatibility/import view only after a bounded transition.

**FC4-02 (Proposed).** “Canonical” means the only write authority for the authored content of that object class. A copy, cache, snapshot, export or browser-local record does not become canonical by being present, newer, more complete or easier to read. This definition applies to authored content; it assigns no authority to an object class whose authority this paper leaves unestablished (FC4-04a).

**FC4-03 (Proposed).** Every object class has one canonical write authority, named readers, prohibited writers, a revision rule and a recovery/migration disposition. This requirement covers the authored object classes named in FC4-04; it does not extend to memory records, whose authority FC-4 leaves unestablished (FC4-04a).

---

## 2. Canonical object classes

**FC4-04 (Proposed — authority table).**

| Object class | Canonical owner | Conversation may hold | Browser storage |
|---|---|---|---|
| World | Canonical library | Immutable revision reference + snapshot | Compatibility/import only |
| Scenario | Canonical library | Immutable revision reference + snapshot | Compatibility/import only |
| Character | Canonical library | Immutable revision reference + snapshot | Compatibility/import only |
| Voice profile | Canonical library | Immutable profile revision + character binding | Compatibility/import only |
| User persona | Canonical library | Immutable revision reference + snapshot | Compatibility/import only |
| Lore object | Canonical library | Immutable revision reference + snapshot | Compatibility/import only |
| Conversation state | FC-1 ledger/state authority | Native | Not canonical |
| Memory record | Unestablished here — FC-7 boundary | Not defined by FC-4 | Not established here |

**FC4-04a (Unresolved — memory-record boundary).** FC-4 does not establish MEMORY-09 authority and does not define memory-record ownership, applicability or integration. FC-7 remains blocked by the delegated D8 decision. This paper neither grants nor denies authority to independently authorized memory sources; that question belongs outside FC-4.

**FC4-05 (Proposed).** The table names authority, not implementation. The actual store, file format and API remain design decisions.

**FC4-05a (Proposed — D7 storage boundary).** Voice-profile assets and character-voice assignments are authored/configured content under D3. The canonical library owns voice-profile revisions; character configuration binds a stable character ID to a voice-profile ID; a conversation references the binding but does not own it. A browser-local copy is not canonical.

**FC4-06 (Unresolved).** The exact boundary between a “character” and a “persona,” and whether either may have multiple active revisions in one conversation, requires owner confirmation before schema design.

---

## 3. Revision model

**FC4-07 (Proposed — immutable revision).**

- Every authored object has a stable object ID.
- Every saved revision has an immutable revision ID.
- The library may maintain a current-revision pointer; that pointer is mutable and is not a reproducibility reference.
- A conversation reference names object ID + revision ID, never only the display name.
- A snapshot records the exact revision IDs used.

**FC4-08 (Proposed).** Editing an object creates a new revision; it does not rewrite the revision already referenced by a conversation. A user may explicitly choose to update a conversation reference, but that is a new recorded transition.

**FC4-09 (Proposed — restoration).**

| Condition | Result |
|---|---|
| All referenced revisions exist | `valid` |
| Some revisions missing | `degraded`, named missing references |
| Object ID resolves but revision does not | `unresolved`; no latest fallback |
| Display name resolves but ID does not | Refuse; names are not identity |

**FC4-10 (Calculated).** This is how D1 and D3 remain compatible: the event ledger owns conversation transitions; the library owns authored revisions; the snapshot preserves exact revision references for restoration; no layer silently substitutes for another.

---

## 4. Browser storage and snapshots

**FC4-11 (Owner-approved — D3).** Browser local storage is not canonical. It may be used as a compatibility/import view only after a bounded transition plan.

**FC4-12 (Proposed).** A browser-local record may be imported only with an explicit identity mapping and provenance. If the mapping is ambiguous, the record is quarantined as unresolved; it is not merged by name or “most recent.”

**FC4-13 (Proposed).** A conversation snapshot is reproducibility evidence. It may contain display copies of authored content, but those copies are read-only and do not acquire library write authority. “Copy” in this paper means a redundant copy of authored content held in a snapshot, browser record or export; it does not refer to the independently owned memory copies defined in FC1-22f–FC1-22l, whose ownership and lifecycle FC-4 does not establish.

**FC4-14 (Proposed).** Deleting a library object does not silently delete conversation snapshots. Orphanhood alone does not authorize deletion. No automatic orphan cleanup and no cascading deletion of authored objects are authorized. This clause is scoped to authored objects and to reproducibility snapshots, and it does not govern memory records: FC-1 directs the cleanup of a parent's own original memory records as part of that parent's explicit deletion, and preserves independently owned copies through it (FC1-22c, FC1-22e, FC1-22i). Such a directed cleanup is bounded by the deletion the user performed and is not a general orphan sweep. FC-4 establishes no memory-record ownership, copy lifecycle or snapshot lifecycle. FC-1 addresses these lifecycles, while destination ownership remains unresolved as stated in FC1-22g; FC4-04a continues not to grant or deny memory-record storage authority. The interaction between orphan cleanup, privacy deletion and retained snapshots remains unresolved (`FC4-DELETION`) and must be decided before any automated orphan cleanup or privacy-driven deletion affecting retained snapshots. Legacy, unknown-origin and interrupted-lifecycle records remain unresolved in FC-1 (`FC1-D2-ORPHAN`, `FC1-CHAT-SCOPED`) and are not addressed here.

---

## 5. Synchronization and write authority

**FC4-15 (Proposed — write path).**

1. UI requests an authored-object change.
2. Canonical library validates identity and creates a revision.
3. Library returns object ID + revision ID.
4. Conversation records a new reference/transition if the user applies it.
5. Snapshot captures the exact revision if reproducibility is required.

**FC4-16 (Proposed — prohibited paths).**

- Browser local storage directly writing canonical library objects.
- Conversation save silently creating or replacing a library object.
- Display-name matching to resolve an object.
- “Latest revision” substitution during restoration.
- Import that merges by title without an ID map.
- A snapshot copy being edited as if it were the library.

**FC4-17 (Unresolved — bypass inventory).** The current frontend/backend bypass set is **Unestablished**. A file-by-file inventory is required before migration design. This paper does not claim that the table above matches the current implementation.

**FC4-18 (Unresolved — migration/quarantine).** The migration, quarantine, conflict and rollback policy for legacy/browser-local content is not decided. No cutover is authorized.

---

## 6. Errors, conflicts and recovery

**FC4-19a (Proposed — failure obligations).**

| Failure | Required behavior |
|---|---|
| Missing object ID | Refuse; do not resolve by name |
| Missing revision | `degraded`/`unresolved`; no latest fallback |
| Conflicting browser copy | Quarantine; preserve both; require mapping |
| Library write fails | No conversation reference claims success |
| Snapshot hash mismatch | `unresolved`; do not claim reproducibility |
| Partial import | Report partial/unknown; no silent merge |
| Concurrent library edits | New revision or explicit conflict; never silent overwrite |

**FC4-19 (Proposed).** Canonical identity is never inferred from a filename, folder, display name, browser key or snapshot copy. A failed identity resolution is visible and conservative.

---

## 7. Observability and privacy

**FC4-20 (Proposed).** Each canonical write records object ID, revision ID, previous revision, actor, timestamp, validation result and provenance.

**FC4-21 (Proposed).** A conversation reference records the exact revision used. An export or snapshot states whether it contains canonical references, display copies or both.

**FC4-22 (Proposed).** Authored content may be private. Browser compatibility copies, snapshots and import bundles must carry the same privacy classification as their source; copying does not declassify.

---

## 8. Acceptance design — proposed, not run

1. Create a character revision, reference it in a conversation, edit the character; the conversation still resolves the original revision.
2. Restore a conversation with all revisions present: exact content resolves.
3. Remove one revision: restoration reports `degraded` with the missing reference, not the latest revision.
4. A browser-local record with a matching title but no ID does not overwrite canonical content.
5. A browser-local import with an explicit ID map succeeds; an ambiguous import is quarantined.
6. Conversation save cannot create a canonical library object.
7. A snapshot copy is not editable as canonical content.
8. A deleted library object does not silently delete retained snapshots.
9. Concurrent edits produce distinct revisions or an explicit conflict.
10. No display name is used as a primary key.

---

## 9. Open questions

| ID | Question | Status |
|---|---|---|
| `FC4-BYPASS` | Which files/stores currently write each object class? | **Unestablished** |
| `FC4-MIGRATION` | What is the bounded browser-local transition? | **Unresolved** |
| `FC4-QUARANTINE` | How are ambiguous imports isolated? | **Unresolved** |
| `FC4-DELETION` | How do library deletion, snapshots and privacy deletion interact? | **Unresolved** |
| `FC4-PERSONA` | What is the character/persona boundary? | **Unresolved** |
| `FC4-REVISION-AUTHORITY` | Who may create/retire a canonical revision? | **Unresolved** |

---

## 10. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Owner decision | D3 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Conversation state | FC-1 (this set) |
| Prompt composition | FC-2 (this set) |
| Memory boundary | FC-7 (blocked by D8) |
| Evidence | FC-6 (this set) |
| Phase 1 composition | `HF1-PLAN-2026-11-21-01` |

**FC4-23 (Observed — no authorization).** This paper does not select a storage technology, migrate data, edit browser storage, change a canonical root, or authorize a cutover. The bypass inventory is not claimed complete.

**End of FC-4.** Canonical authority is singular for each authored object class covered by FC-4; references and snapshots preserve exact revisions; browser storage is a compatibility view, never a silent promotion path.
