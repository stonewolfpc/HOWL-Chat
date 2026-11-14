# F3-C Master Architecture and Implementation Contract

**Date:** 2026-11-11  
**Status:** **FROZEN — ARCHITECTURE ONLY.** This document is the candidate implementation authority for the F3-C persistence contract. It does not accept F3-C, authorize P1a, or claim implementation evidence.  
**Supersedes:** the isolated Master-1 through Master-3c amendment sequence as the consolidated implementation authority. The frozen F3-C decision records and accepted native evidence remain supporting authority.  
**Authorization:** No source changes, tests, P1a implementation, or F3-D work are authorized by this document.

---

## 1. Scope and invariants

F3-C governs durable chat identity, revisions, commit receipts, recovery, migration, wipe, and the cast-confirmation persistence boundary.

1. Chat identity is `(rootID, chatID)`, never a display name.
2. Chat IDs are not reused within a root incarnation.
3. Wipe creates a new root incarnation.
4. Every durable revision has exactly one commit node and one commit record.
5. Every commit node has exactly one valid receipt, or a specifically permitted reconstructable missing receipt.
6. Commit inventory and receipt inventory must both be complete and mutually consistent before verified status.
7. Historical proof is independent of current target integrity.
8. A missing receipt may be reconstructed only from exact stored body bytes in its commit record.
9. A malformed receipt or commit record is a hard block and is never overwritten.
10. Authoritative App-state mutations are serialized per chat under `stateApplyMu`.
11. Every root-dependent operation participates in the root barrier.
12. Publication occurs only after the durable commit succeeds.
13. A durable append is a complete commit of its candidate state, not a promise completed later by an NDJSON-only write.
14. Historical and current state are never conflated.
15. Every authoritative mutation caller must be inventoried; the API list is not the mutation surface.

The architecture freeze is not F3-C acceptance. Acceptance requires implementation and the evidence gates in §12.

---

## 2. Inherited frozen decisions

### 2.1 Six F3-C decisions

| ID | Decision | Contract effect |
|---|---|---|
| F3-C-Q1 | Folder per item | World/scenario storage remains `<root>/<id>/world.json` or `scenario.json`; F3-C chat persistence does not alter it. |
| F3-C-Q2 | Stable identity | Identity is opaque; no title/name authorization; legacy IDs preserved or explicitly mapped. |
| F3-C-Q3 | Explicit journaled import | Import is create-only, idempotent by identity plus content hash, journaled, reversible; browser data untouched. |
| F3-C-Q4 | Pinned snapshots | Chat snapshots persist until explicit refresh; no load/edit refresh. F3-C does not introduce implicit snapshot mutation. |
| F3-C-Q5 | Preserve unavailable | Missing references remain unavailable, excluded from prompts, and repairable; never substituted. |
| F3-C-Q6 | Stable lore ownership | World/scenario lore resolves by stable item ID; rename does not change ownership; same-name recreation inherits nothing. |

### 2.2 LF-1

LF-1 authorizes one explicit, journaled, reversible legacy name-folder migration. Legacy name folders are never read at runtime, never a fallback, never written to, and never deleted. Only the explicit migration package may read them.

F3-C chat persistence does not add a runtime name-folder reader. It does not write canonical item folders. `world_lore.json` and `scenario_lore.json` remain item-owned and outside the chat target set.

### 2.3 Acceptance gates

A1–A13 remain the F3-C acceptance matrix. A9, A10, and A12 are accepted on native evidence; A13 is a passing regression gate. A11 remains an availability contract, not a persistence exception.

Important evidence limit: A9/A10/A12 were executed before this consolidated persistence contract and do **not** prove that the new target set, receipts, revisions, or migration code exist. They prove the accepted import, single-writer, and native prompt boundaries only.

---

## 3. Identity, state, and revision domains

### 3.1 Root identity

`<root>/f3c-root-identity.json`:

```json
{
  "rootId": "uuid",
  "incarnation": 1,
  "createdAt": "RFC3339"
}
```

`rootId` is immutable within an incarnation. Wipe installs a new root ID and increments incarnation.

### 3.2 Chat metadata

`ChatMeta` adds:

```go
RootID                string `json:"rootId"`
Revision              uint64 `json:"revision"`
LastCommitOperationID string `json:"lastCommitOperationId"`
```

Rules:

- `Revision == 0` means legacy/uninitialized.
- New chats start at revision 1.
- One successful durable logical mutation advances revision by exactly 1.
- `LastCommitOperationID` identifies the exact current commit.
- Creation, Save As, legacy initialization, full save, append, replace, clear, compound save, memory migration, history migration, and cast confirmation each produce one revision.
- Deletion produces a tombstone, not a new revision.
- State-only mutation produces no revision.

Chat `Revision` is distinct from the `worldstore` `rv1-...` document revision. There is no shared counter or implied equality.

### 3.3 Authoritative state

```text
chatMessages
world
scenario
characters
lorebooks
```

Derived state:

```text
triggeredLore    // cache only; no session, epoch, or revision
```

Session state:

```text
stateSessionID
stateEpoch
stateDirty
statePoisoned
```

### 3.4 Revision and target integrity

A chat revision is valid only when its commit node, receipt, current target manifest, and complete inventories agree. A world/scenario document revision is a separate `worldstore` property and is never inferred from the chat revision.

---

## 4. Durable layout and path identity

### 4.1 Layout

```text
<root>/
  f3c-root-identity.json
  f3c-wipe-intent.json
  config.json
  chats/<chatID>/
    meta.json
    characters.json
    history.json
    history.ndjson
    lorebooks.json
    memories_snapshot.json
    scenario.json              # optional
    world.json                 # optional
    f3c-ordinary-journal.json
  f3c-creation-intents/<F3C-PATH-CREATION-1 hash>.json
  f3c-commits/<F3C-PATH-COMMIT-1 hash>/
    r00000000000000000001/<operationID>.json
  f3c-receipts/<F3C-PATH-RECEIPT-1 hash>/<operationID>.json
  f3c-tombstones/<F3C-PATH-TOMBSTONE-1 hash>/<operationID>.json
  f3c-memory-migrations/<F3C-PATH-MEMORY-1 hash>.json
  f3c-tmp/<uuid>.json
  f3c-cast-journal.json
```

`f3c-tmp` is outside every visible inventory. Temporary records, receipts, targets, and tombstones are never placed inside visible inventory directories.

### 4.2 Path hashes

```text
pathHash(domain, parts...) =
  lowercaseHex(
    SHA-256(domain || 0x00 || join(parts, 0x00))
  )
```

Domains:

```text
F3C-PATH-COMMIT-1       -> rootID, chatID
F3C-PATH-RECEIPT-1      -> rootID, chatID
F3C-PATH-TOMBSTONE-1    -> rootID, chatID
F3C-PATH-CREATION-1     -> rootID, operationID
F3C-PATH-MEMORY-1       -> rootID, chatID
```

Rules:

- exact UTF-8 parts, no normalization or case folding;
- one `0x00` separator, no trailing separator;
- exactly 64 lowercase hexadecimal characters;
- operation filenames use canonical lowercase UUIDv4;
- creation, commit, receipt, tombstone, and memory-migration domains never share a hash.

---

## 5. Canonicalization, target bytes, and hashing

### 5.1 `F3C-CANON-1`

RFC 8785 JSON Canonicalization Scheme with:

- schema validation before canonicalization;
- duplicate-key rejection at every depth;
- unknown-field rejection;
- invalid UTF-8 and lone surrogates rejected;
- no Unicode normalization;
- sorted object keys; ordered arrays;
- no NaN, Infinity, negative zero, fractions, or exponents;
- unsigned integers `0..2^53-1`;
- RFC 3339 timestamps;
- lowercase hexadecimal hashes;
- target manifests sorted by exact path.

### 5.2 Request fingerprint

```text
requestFingerprint =
  SHA-256("F3C-REQUEST-1" || 0x00 || canonical request bytes)
```

The request binds operation type, root/chat identity, operation ID, base revision, operation-specific input, Save As source binding, and cast token binding. It is stored in receipts, commit records, creation intents, tombstones, and operation-map entries.

### 5.3 Hash domains

```text
receiptBodyHash =
  SHA-256("F3C-RECEIPT-BODY-1" || 0x00 || canonical body bytes)

commitRecordHash =
  SHA-256("F3C-COMMIT-RECORD-1" || 0x00 ||
          canonical record excluding commitRecordHash)

receiptEnvelopeHash =
  SHA-256("F3C-RECEIPT-ENVELOPE-1" || 0x00 ||
          canonical envelope excluding receiptEnvelopeHash)

treeHash =
  SHA-256("F3C-TREE-HASH-1" || 0x00 || canonical tree manifest)

tombstoneHash =
  SHA-256("F3C-TOMBSTONE-1" || 0x00 ||
          canonical tombstone excluding tombstoneHash)

memoryMigrationJournalHash =
  SHA-256("F3C-MEMORY-MIGRATION-1" || 0x00 ||
          canonical pre-commit migration journal excluding
          commitRecordHash, receiptEnvelopeHash, and journalHash)

migrationFinalJournalHash =
  SHA-256("F3C-MEMORY-MIGRATION-1" || 0x00 ||
          canonical final migration journal excluding journalHash)
```

### 5.4 Acyclic construction

1. Canonical receipt body without hash fields.
2. Exact canonical UTF-8 body bytes.
3. `receiptBodyHash`.
4. Commit record with exact `receiptBodyBytes` and body hash.
5. `commitRecordHash`.
6. Receipt envelope with body bytes, body hash, commit-record hash.
7. `receiptEnvelopeHash`.

No hash input contains its own output. `receiptBodyBytes` uses RFC 4648 standard Base64 with required padding, no URL-safe alphabet, no whitespace, and byte-identical re-encoding.

### 5.5 Complete target set

| Relative path | Presence | Owner |
|---|---|---|
| `meta.json` | required | Chat persistence |
| `characters.json` | required | Chat persistence |
| `history.json` | required | Chat persistence |
| `history.ndjson` | required | Chat persistence |
| `lorebooks.json` | required | Chat-scoped snapshot; not canonical item lore |
| `memories_snapshot.json` | required | Chat-owned generated lore |
| `scenario.json` | optional; absent if none | Chat snapshot |
| `world.json` | optional; absent if none | Chat snapshot |

`memories.json` is **not** a canonical target. It is a legacy/separate memory artifact and is never a substitute for `memories_snapshot.json`.

### 5.6 Exact target serialization — `F3C-TARGET-1`

JSON targets use:

- Go `encoding/json.MarshalIndent(value, "", "  ")`;
- default HTML escaping;
- UTF-8; LF; no trailing newline; no BOM;
- frozen Go struct declaration order;
- deterministic Go map-key order;
- required arrays normalized to `[]`, never `null`.

`meta.json` order:

```text
id, name, createdAt, updatedAt, worldName, scenarioName,
characterNames, lorebookNames, worldId, scenarioId,
activeCharacters, removedCharacters, addedCharacters,
activeLorebooks, removedLorebooks, addedLorebooks,
rootId, revision, lastCommitOperationId
```

`history.json` is ordered `[]ChatHistory`:

```text
id, sender, text, role, timestamp, hasImage, characterId
```

`characters.json` is ordered `[]CharacterContext`: embedded `Character` fields (`id`, `name`, `description`, `personality`, `backstory`, `goals`, `relationships`, `references`, `age`, `role`, `background`, `priority`, `version`, `voice_config`), then `image` with `path`, `thumbnail`.

`lorebooks.json` is the ordered chat-scoped array.

`memories_snapshot.json`:

```json
{
  "world_memories": [],
  "scenario_memories": [],
  "character_memories": {}
}
```

Field order is `world_memories`, `scenario_memories`, `character_memories`. Character map keys use deterministic sorted order. Empty required fields are `[]` and `{}`.

`world.json` / `scenario.json` use the corresponding frozen context structs when present.

`history.ndjson`:

- compact Go `json.Marshal` per `ChatHistory`;
- exact one LF after each entry;
- no blank lines;
- one trailing LF after the last entry;
- empty history is a zero-byte file;
- no indentation, BOM, or CRLF.

`history.json` and `history.ndjson` must decode to the same ordered entries.

### 5.7 Target hashes and manifest

```text
tgt1-<64 lowercase hexadecimal characters>
```

The hash is SHA-256 over exact persisted bytes, not canonicalized JSON.

Manifest:

```json
{
  "format": "f3c-target-manifest-1",
  "targets": [
    {"path":"characters.json","presence":"present","bytes":123,"hash":"tgt1-..."},
    {"path":"scenario.json","presence":"absent"}
  ]
}
```

Rules:

- exactly the eight target paths are represented;
- sorted by path in ascending bytewise UTF-8 order;
- `present` requires `bytes` and `hash`;
- `absent` forbids the file;
- the manifest is canonicalized under `F3C-CANON-1` inside receipts/records;
- target bytes themselves are never canonicalized before hashing.

Verification outcomes:

| Condition | Outcome |
|---|---|
| Required target missing | `ErrTargetManifestIncomplete` |
| Optional absent but file exists | `ErrTargetConflict` |
| Optional present but file missing | `ErrTargetManifestIncomplete` |
| Readable, hash differs | `ErrTargetConflict` |
| Unreadable | `ErrTargetUnknown` |
| Readable but malformed | `ErrTargetMalformed` |
| History files disagree | `ErrTargetConflict` |
| Meta identity/revision mismatch | `ErrTargetConflict` |

Historical manifests are never compared with current files.

### 5.8 Compatibility decisions

These are contract decisions, not claims that the current writer already conforms:

1. `memories_snapshot.json` is canonical, but no production writer currently exists in this tree. Implementation must add it.
2. The legacy snapshot is read from `paths.StateRoot()/Chats/<id>`; it migrates to `<root>/chats/<id>`.
3. `lorebooks.json` is required because current chat persistence writes and reads it; it remains a chat-scoped copy.
4. Current `SaveChat` does not remove stale optional world/scenario files. F3-C must remove them or fail.
5. Current Go serialization turns nil slices into `null`; F3-C requires `[]`.
6. `memories.json` is preserved, reported, and migrated only by explicit ownership decision; it is never silently deleted or treated as the snapshot.

---

## 6. Commit records, receipts, chain, and inventories

### 6.1 Receipt body

```text
version, rootID, chatID, operationID, operationType,
requestFingerprint, baseRevision, committedRevision,
candidateFingerprint, targetManifest,
previousCommitOperationID, previousCommitRecordHash,
previousReceiptEnvelopeHash, commitSequence, createdAt,
result payload,
memoryMigrationJournalHash       // optional; memory migration only
```

The body excludes all hash fields.

### 6.2 Commit record

```text
version, rootID, chatID, operationID, operationType,
requestFingerprint, baseRevision, committedRevision,
candidateFingerprint, targetManifest,
receiptBodyBytes, receiptBodyHash,
previousCommitOperationID, previousCommitRecordHash,
previousReceiptEnvelopeHash, commitSequence, createdAt,
memoryMigrationJournalHash       // optional
commitRecordHash
```

### 6.3 Receipt envelope

```text
receiptBodyBytes, receiptBodyHash, commitRecordHash, receiptEnvelopeHash
```

### 6.4 Current commit proof

Current commit is `meta.Revision + meta.LastCommitOperationID`.

Verified current state requires:

- exactly one commit record at that revision and operation;
- exactly one valid receipt, or one successfully reconstructed permitted missing receipt;
- valid receipt body, record, and envelope hashes;
- receipt body bytes equal record body bytes;
- `committedRevision == meta.Revision`;
- root/chat identity match;
- current targets match the manifest;
- complete commit inventory;
- complete receipt inventory;
- one-to-one commit/receipt mapping.

Zero, multiple, malformed, or mismatched evidence fails.

### 6.5 Historical proof

Historical verification checks hashes, operation identity, predecessor links, revision links, unique operations, no cycles, and receipt/record body agreement. It does not compare historical manifests with current files.

### 6.6 Revision chain

- Revision 1: `baseRevision = 0`, no predecessor; creation, Save As, or legacy initialization.
- Revision N > 1: predecessor revision N-1, `committedRevision = N`.
- Later saves do not invalidate the creation manifest.
- Missing historical receipts may be reconstructed from valid record body bytes; malformed receipts never are.

### 6.7 Commit inventory

Path:

```text
f3c-commits/<F3C-PATH-COMMIT-1 hash>/r<20-digit-revision>/<operationID>.json
```

Strict scan under `tree.R` and `chatWriteMu`:

- exactly revisions 1..current;
- one node per revision;
- one successor below current;
- no record above current;
- no duplicate operation, revision, or unaccounted record;
- no missing/empty revision directory;
- no malformed, unknown-version, hidden, temporary, symlink, junction, reparse, hard-linked, or non-regular entry;
- no unexpected filename or identity mismatch;
- no broken predecessor or cycle.

### 6.8 Receipt inventory

Path:

```text
f3c-receipts/<F3C-PATH-RECEIPT-1 hash>/<operationID>.json
```

Strict scan:

- directory and ancestors checked for symlinks/reparse points;
- regular `.json` entries with exact `<operationID>.json` names;
- embedded root/chat/operation must match directory and filename;
- corresponding commit node exists;
- receipt bytes equal record `receiptBodyBytes`;
- body hash, commit-record hash, and envelope hash verify;
- operation type, revisions, candidate fingerprint, and manifest agree;
- no receipt without a commit node;
- no duplicate or conflicting receipt;
- exactly one receipt per node or a permitted reconstructable missing receipt;
- unexpected, hidden, temporary, symlink, reparse, hard-linked, nested, or uninspectable entries block.

Error distinction:

```text
ErrReceiptInventoryIncomplete   // cannot inspect fully
ErrReceiptInventoryInvalid      // inspected and invalid
ErrReceiptEvidenceGap           // missing and unreconstructable
ErrReceiptConflict              // conflicting evidence
ErrReceiptMalformed             // malformed receipt
```

Missing receipts are reconstructed create-only from record body bytes. Malformed receipts are hard blocks. Reconstruction collisions fail; they never overwrite.

### 6.9 Temporary installation and retention

Commit records, receipts, tombstones, and targets are written under `f3c-tmp`, flushed, and installed create-only by same-volume atomic move. Final collisions fail. Crash leftovers are cleaned without being seen by inventory scanners. Required evidence remains while a chain depends on it. Compaction requires a separately specified checkpoint; a cache is not proof.

---

## 7. Locks, operation map, and publication order

### 7.1 Locks

| Lock | Scope | Protects |
|---|---|---|
| `castApplyMu` | root | cast apply/recovery; wipe |
| `chatSwitchMu` | root | active-session transitions |
| `chatTreeMu` | root tree | root barrier vs wipe |
| `stateApplyMu` | per chat | authoritative App-state mutation |
| `chatWriteMu` | per chat | targets, records, receipts, verified reads |
| `rootJournalMu` | root | root journals |
| `chatStateMu` | App | map, active identity, session, epoch, publication |

### 7.2 Order

```text
castApplyMu
  -> chatSwitchMu
    -> chatTreeMu.RLock/Lock
      -> stateApplyMu(chat)       // stable order
      -> chatWriteMu(chat)        // stable order
        -> rootJournalMu

chatStateMu is a leaf lock.
```

Rules:

1. `chatStateMu` is brief for registration/snapshot or publication and released before acquiring any other lock.
2. No lock is acquired while holding `chatStateMu`.
3. No disk I/O under `chatStateMu`.
4. No code acquires `chatWriteMu` then waits for `stateApplyMu`.
5. `stateApplyMu` may precede `chatWriteMu`.
6. Multi-chat locks use stable `(rootID, chatID)` order.
7. No `chatSwitchMu` while holding tree/state/write locks.
8. No `castApplyMu` while holding `chatSwitchMu`.
9. `rootJournalMu` only under the relevant `chatWriteMu`.
10. Verification helpers never acquire `stateApplyMu` under `chatWriteMu`.

### 7.3 Standard durable mutation

1. Register operation under `chatStateMu`; release.
2. Acquire `tree.R`; revalidate root.
3. Acquire `stateApplyMu`.
4. Revalidate session, epoch, poison, base.
5. Acquire `chatWriteMu`; acquire `rootJournalMu` if needed.
6. Verify/write targets, record, receipt.
7. Release journal/write locks.
8. While retaining `stateApplyMu`, acquire `chatStateMu` briefly.
9. Publish exact candidate; advance epoch if content changed; clear dirty only if exact snapshot.
10. Release state and tree.

Publication occurs after durable locks. A verifier may observe a committed receipt before publication; that is the defined `CommittedButNotPublished` window.

### 7.4 Autosave

At `chatWriteMu` ownership, revalidate captured root, chat, `stateSessionID`, epoch, base revision, poison, and discard generation.

- Session changed → `ErrSessionInvalidated`; no durable write.
- Explicitly discarded → `ErrAutosaveDiscarded`; no durable write.
- Root/chat/base changed → `ErrStaleChatState`.
- Poisoned → abort.
- Epoch advanced but session unchanged → may commit snapshot, but must not publish stale state or clear dirty.

An autosave waiting after a switch cannot write the discarded outgoing state.

### 7.5 Switch

1. Snapshot active identity under `chatStateMu`; release.
2. `chatSwitchMu`.
3. `tree.R`.
4. Revalidate root/chat.
5. `stateApplyMu(outgoing)` and `stateApplyMu(incoming)` in stable order.
6. `chatWriteMu(outgoing)` and incoming verification lock in stable order.
7. Save outgoing or explicit discard.
8. Failed save aborts; outgoing remains active.
9. Verify incoming.
10. Release write locks.
11. Publish new identity/session under `chatStateMu`.
12. Apply incoming state.

### 7.6 Operation map

Key `(rootID, chatID, operationID)`; states in-flight/completed.

- registration under `chatStateMu`, then release;
- duplicate waiters never acquire `stateApplyMu`;
- in-flight entries never evicted;
- completed entries evicted only at zero waiters;
- cancellation before execution ownership removes the entry;
- cancellation after execution ownership does not abort the owner;
- historical results return `historical: true`;
- wipe invalidation closes old-root entries once after the root commit.

---

## 8. Operations

### 8.1 Creation

Sequence:

1. Persist creation intent.
2. Persist `prepared`.
3. Persist and sync `writing`.
4. Create chat directory create-only.
5. Write targets.
6. Read back and verify.
7. Create commit record.
8. Create receipt.
9. Mark journal `committed`.
10. Publish configuration/activation.

No directory is created while `prepared`. `writing` authorizes all creation artifacts.

Recovery:

- `prepared`, no artifacts → may durably advance to `writing`; no candidate adoption.
- `prepared`, any directory/target/record/receipt → conflict.
- `writing`, no directory → create and continue.
- `writing`, partial targets → complete only with full candidate evidence.
- all targets candidate → complete record/receipt and commit.
- invalid/malformed receipt → block; no replacement chat ID.

Creation identity revalidation before advancing or continuing:

```text
creation intent
journal
allocated chatID
rootID
operationID
operationType
requestFingerprint
Save As source binding
baseRevision = 0
```

A phase alone does not authorize another operation's artifacts. Mismatch → `ErrCreationIdentityConflict`.

### 8.2 Save As

Creation intent binds source root/chat, source revision, source commit operation, source record hash, source receipt envelope hash, source manifest, source fingerprint, and copied target hashes.

- source must be durable and verified;
- active dirty source must flush first or Save As aborts;
- source write lock held while copying;
- new chat starts at revision 1;
- no activation;
- no replacement chat ID;
- source mismatch during recovery abandons the intent with `ErrSourceChanged`; no substitution.

### 8.3 Full/compound saves

`SaveChatData`, world/scenario add/clear:

- tree.R, stateApply, chatWrite;
- candidate built without publication;
- revalidate base at write ownership;
- write/readback/record/receipt;
- publish exact candidate and advance epoch;
- pre-commit failure leaves prior state;
- commit-but-publish failure poisons and reloads.

### 8.4 Append/replace/clear

- copy-on-write full candidate;
- tree.R, stateApply, chatWrite;
- both history targets and complete manifest;
- readback, record, receipt;
- publish exact candidate; one revision; epoch increments; clean.
- failure before commit leaves prior state; after commit before publication poisons.

### 8.5 Streaming and dirty state

- token events are authoritative state mutations; epoch advances; dirty true; no revision.
- later complete append/autosave commits full snapshot.
- dirty means current state differs from last exactly published durable snapshot.
- clean autosave is a no-op; dirty autosave commits full snapshot.
- stale autosave never publishes or clears dirty.
- no NDJSON-only persistence.

### 8.6 Delete

Tombstone path:

```text
f3c-tombstones/<F3C-PATH-TOMBSTONE-1 hash>/<operationID>.json
```

Fields:

```text
format, version, rootID, chatID, operationID,
requestFingerprint, baseRevision,
lastCommitOperationID, lastCommitRecordHash,
lastReceiptEnvelopeHash, deletedAt, tombstoneHash
```

Rules:

- create-only; active deletion uses switch;
- failed outgoing save aborts;
- verify last commit evidence;
- tombstone before directory removal;
- same request idempotent; different request conflict;
- already-tombstoned chat returns `ErrAlreadyDeleted`;
- malformed/hash-mismatched tombstone hard block;
- missing chat without tombstone is `ErrChatStateUnknown`;
- tombstone present + directory present completes deletion;
- no revision and no ID reuse.

Tombstone inventory: exactly one tombstone per chat; no hidden/temp/symlink/reparse/hard-link/unaccounted entries; uninspectable/invalid/duplicate/evidence gap block.

### 8.7 Cast

Preview is non-durable. Confirmation requires token and operation ID; token binds chat, base, candidate, scope, expiry. Committed expired token resolves to receipt; never reapplies; different candidate conflicts. Uses the same commit/record/receipt/revision protocol. Cast confirmation remains absent from current source.

### 8.8 Memory migration

Canonical `memories_snapshot.json`. Legacy `memories.json` is preserved and separate.

Migration journal:

```text
<root>/f3c-memory-migrations/<F3C-PATH-MEMORY-1 hash>.json
```

Fields:

```text
format, version, rootID, chatID, operationID,
requestFingerprint, sourcePath, sourceArtifactKind,
sourceBytesHash, sourceEntryCount, destinationPath,
destinationBytesHash, destinationEntryCount, phase,
journalPreCommitHash, commitRecordHash,
receiptEnvelopeHash, createdAt, journalHash
```

`journalPreCommitHash` and `journalHash` are defined in §5.3. `commitRecordHash` and `receiptEnvelopeHash` appear only in the committed final journal. Sequence:

1. Read/validate legacy source.
2. Record source path, hash, and entry count.
3. Create destination create-only.
4. Read back and verify.
5. Include in target manifest, record, and receipt.
6. After chat commit succeeds, mark migration `committed`.

Outcomes:

| Destination | Source | Outcome |
|---|---|---|
| Absent | Valid | Migrate exact bytes; source preserved. |
| Present | Absent | Preserve; `already_canonical`. |
| Both identical | Both | Adopt; no rewrite; source preserved. |
| Both different | Both | Preserve both; conflict; no merge/overwrite. |
| Unreadable/malformed source | — | Preserve; block; no destination invented. |
| Destination created, chat commit fails | — | `commit-pending`; not complete; source preserved. |

Receipt body includes `memoryMigrationJournalHash`, which is the pre-commit journal hash. The final migration journal includes `journalPreCommitHash`, `commitRecordHash`, `receiptEnvelopeHash`, `phase = committed`, and `journalHash`; `journalPreCommitHash` must equal the receipt field. A verified destination with missing/corrupt migration journal is **not** migration-complete: it is `memory_migration_evidence_gap`. Receipt and migration evidence must agree.

`memories.json` is never marked migrated by discovery and never becomes a substitute.

### 8.9 History migration

When NDJSON absent and JSON exists:

1. Parse JSON array strictly.
2. Require non-empty unique message IDs.
3. Preserve order and IDs exactly.
4. Generate exact NDJSON.
5. Write create-only.
6. Read back both.
7. Verify same count, order, IDs, and all fields.
8. Include both in manifest/record/receipt.
9. Commit.

Rules:

- no index-derived IDs;
- no renumbering;
- duplicate/empty/malformed ID blocks;
- malformed NDJSON blocks; no silent fallback;
- existing NDJSON and JSON must agree; otherwise `ErrHistoryMismatch`;
- neither exists → empty `[]` JSON and zero-byte NDJSON.

### 8.10 Wipe

Wipe intent at `<root>/f3c-wipe-intent.json`, hash domain `F3C-WIPE-INTENT-1`.

Phases: `prepared`, `staged`, `root-committing`, `root-committed`, `cleanup`, `committed`.

Independent state:

```text
appStageState:      none | started | renamed | verified
galleryStageState:  absent | started | renamed | verified
preCommitGate:      pending | passed
```

Required staged components: `chats/`, `f3c-commits/`, `f3c-receipts/`, `f3c-tombstones/`, `f3c-creation-intents/`, `f3c-memory-migrations/`, `f3c-cast-journal.json` if present, and gallery root.

Tree hash `F3C-TREE-MANIFEST-1`:

```json
{"path":"relative/path","type":"file","size":123,"sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}
{"path":"relative/dir","type":"dir"}
```

- relative paths, `/` separators, exact UTF-8, no normalization;
- sorted bytewise;
- `file` with exact size and SHA-256;
- `dir`, including empty;
- no links/reparse points/hard links/unsupported objects; fail closed;
- manifest canonicalized; tree hash domain `F3C-TREE-HASH-1`.

Pre-commit gate requires every existing component staged, staged hashes/counts matching, no source remains, no unexpected component, complete inventory, and gallery proof. Failure restores/aborts; new identity never installed.

Gallery is a separate boundary; same-volume rename; cross-volume abort unless separately specified. A successful chat stage cannot hide incomplete gallery staging.

Commit point remains the new root identity file. `preCommitGate=passed` and `root-committing` are preconditions, not commit points.

Crash recovery reconciles actual source/staged paths; duplicate/lost/ambiguous states block with `ErrWipeBoundaryAmbiguous`. Before `root-committing`, restore old root. After new identity, finish cleanup. Gallery cleanup failure is `wipe_cleanup_incomplete`, not rollback.

---

## 9. Journal recovery

### 9.1 General

Before recovery writes:

1. Validate root identity.
2. Validate ordinary and cast journals.
3. Both nonterminal → `ErrDualJournalAmbiguous`.
4. Identity/version/schema mismatch → block.
5. Verify operation ID and request fingerprint.
6. Read before/candidate/readback evidence.
7. Read record/receipt evidence.
8. Choose recovery.

### 9.2 Ordinary phases

| Phase | Evidence | Recovery |
|---|---|---|
| `prepared` | all-before, no record | Abort; remove temps/journal. |
| `prepared` | all-candidate | Conflict; no adoption, no record. |
| `prepared` | mixed | Ambiguous. |
| `prepared` | any record/receipt | Conflict. |
| `writing` | all-before | Remove temps; failed. |
| `writing` | all-candidate, base current, no record | Create record/receipt; commit. |
| `writing` | all-candidate, record valid | Ensure/reconstruct receipt; commit. |
| `writing` | mixed | Ambiguous. |
| `written` | all-candidate, base current | Create record/receipt if absent; commit. |
| `written` | base stale, no record | Block. |
| `written` | mixed | Ambiguous. |
| `verified` | all-candidate, record valid | Ensure/reconstruct receipt; commit. |
| `verified` | all-candidate, no record, base current | Create record/receipt; commit. |
| `verified` | mixed/stale | Block/ambiguous. |
| `committed` | record/receipt valid | Publish if needed; clean. |
| `committed` | receipt missing | Reconstruct. |
| `committed` | record missing/malformed | Block. |
| `committed` | receipt malformed | Hard block. |
| `failed` | all-before | Clean. |
| `failed` | mixed | Ambiguous. |
| `ambiguous` | any | No serve/write; explicit recovery. |

### 9.3 Cast phases

Same evidence rules. `prepared` never adopts. `writing`/later may complete if token/base valid. Committed cast with later revisions is recognized, not reapplied.

### 9.4 Creation

See §8.1. `prepared` with any artifact conflicts. Identity must agree before advancement.

### 9.5 Wipe and gallery

See §8.10. No ambiguous guess; commit point only root identity.

### 9.6 Memory/history migration

Memory migration never completes without both a valid chat receipt and a valid final migration journal. History migration requires both target hashes and equal ordered entries/IDs.

### 9.7 Failure taxonomy

```text
ErrReceiptMalformed
ErrReceiptInventoryIncomplete
ErrReceiptInventoryInvalid
ErrReceiptConflict
ErrReceiptEvidenceGap
ErrCommitInventoryIncomplete
ErrCommitInventoryInvalid
ErrCommitFork
ErrCommitChainBroken
ErrCurrentStateMismatch
ErrRootIncarnationChanged
ErrChatStateUnknown
ErrCreationIntentMissing
ErrCreationIdentityConflict
ErrSourceChanged
ErrPreparedJournalConflict
ErrIdempotencyConflict
ErrStaleChatState
ErrSessionInvalidated
ErrAutosaveDiscarded
ErrAlreadyDeleted
ErrTargetConflict
ErrTargetUnknown
ErrTargetMalformed
ErrTargetManifestIncomplete
ErrHistoryIdentityMissing
ErrHistoryIdentityDuplicate
ErrHistoryMismatch
ErrHistoryReadbackMismatch
ErrHistoryMalformed
ErrMemoryMigrationConflict
ErrMemorySourceUnreadable
ErrMemorySourceMalformed
ErrMemoryMigrationEvidenceGap
ErrTombstoneInventoryIncomplete
ErrTombstoneInventoryInvalid
ErrTombstoneMalformed
ErrTombstoneConflict
ErrTombstoneEvidenceGap
ErrWipeStagingIncomplete
ErrWipeGalleryStagingIncomplete
ErrWipeBoundaryAmbiguous
ErrWipeIntentAmbiguous
ErrWipeRestoreFailed
ErrDualJournalAmbiguous
CommittedButNotPublished
stateDirty
statePoisoned
```

Neither incomplete inventory nor invalid inventory is reported as verified.

---

## 10. Required APIs and mutation surface

Mutating APIs carry operation ID, chat ID, base revision, state session ID, and base state epoch. Server derives root and operation type. Callers never supply operation type.

Illustrative public surface:

```text
NewChat(name, operationID)
SaveChatAs(sourceChatID, name, operationID)
LoadChat(chatID, operationID)
DeleteChat(chatID, baseRevision, operationID)
SaveChat(fullState, context)
SaveChatData(data, context)
AppendChatMessage(role, content, sender, context)
ReplaceChatHistory(entries, context)
ClearChatHistory(context)
AddCharacterToChat(character, context)
RemoveCharacterFromChat(characterID, context)
ResetSessionCharacters(context)
SetLorebooks(entries, context)
UpdateLorebookContext(activeIDs, context)
AddWorldToChat(world, context)
AddScenarioToChat(scenario, context)
ClearWorldFromSession(context)
ClearScenarioFromSession(context)
PreviewChatScenarioCast(chatID, baseRevision)
ConfirmChatScenarioCast(chatID, token, baseRevision, operationID)
WipeAllData(operationID)
```

Operation types include creation, Save As, full/autosave/final save, append, history replace/clear, context changes, delete, cast preview/confirm, legacy revision init, history migration, memory migration, and wipe.

The API list is not the mutation surface. Every actual caller that changes `chatMessages`, `world`, `scenario`, `characters`, or `lorebooks` must use the same `tree.R` + `stateApplyMu` serialization and epoch rules. Direct writes outside that path are defects. `LoadChatData` is read-only; hydration occurs only through verified load/recovery.

---

## 11. Acceptance and evidence program

No tests are authorized now. Eventual acceptance must include:

The controlling implementation-verification method is Gate v2, `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`. Gate v2 does not replace this evidence list, authorize implementation, or grant acceptance. Each item below requires a claim → evidence → adversarial counterexample → result record under Gate v2.

1. Canonicalization, request fingerprints, Base64, hash chain.
2. Receipt inventory and reconstruction.
3. Commit fork/orphan/missing inventory.
4. Target exact bytes, `tgt1-` manifest, memory snapshot, history identity.
5. Append/dirty/autosave; session invalidation; stale snapshots.
6. Locks/publication; autosave vs switch; wipe vs save.
7. Save As source binding; operation ID conflicts.
8. Tombstones and deletion recovery.
9. Creation journal boundaries and identity revalidation.
10. Memory migration conflict/commit/receipt agreement.
11. History migration order/ID preservation.
12. Wipe tree hashes, gallery boundary, staging failure, root identity commit.
13. Fault injection at every journal boundary.
14. Filesystem recovery (symlinks, reparse points, hard links, partial staging).
15. Race testing.
16. Native Wails/WebView2 paths.
17. Independent Mimo Flash Fast challenge before final approval.

A passing unit-test suite does not prove these guarantees. Gate v2 requires every row to be resolved; a test name without an assertion, input, observed result, and negative case is not evidence.

---

## 12. Known source gaps at freeze

Current source does not implement:

- root identity/incarnation;
- chat revisions/commit IDs;
- immutable commit records/receipts/inventories;
- canonical hashing;
- creation intents;
- deletion tombstones;
- operation maps/waiters;
- `stateApplyMu`/`chatSwitchMu`/root barrier;
- receipt reconstruction;
- cast confirmation;
- a writer for `memories_snapshot.json`.

Known defects include:

- `LoadChat` continues after failed outgoing save;
- `SaveChatData` applies state before persistence;
- `AppendHistoryEntry` writes only NDJSON;
- `NewChat` uses timestamp IDs and `MkdirAll`;
- `SaveChatAs` is two-step;
- `DeleteChat` has no tombstone;
- `WipeAllData` lacks root identity/waiter invalidation;
- `buildChatHistory` derives message IDs from array length;
- optional world/scenario files are not removed when cleared;
- nil slices can serialize as `null`.

These are implementation gaps, not contradictions in the contract.

---

## Appendix A — Decision history

| Revision | Contribution |
|---|---|
| Master-1 | Consolidated identity, revision chain, receipts, locks, recovery, and cast boundary. |
| Master-2 | Added strict receipt inventory, request fingerprints, Save As source binding, tombstones, phase recovery, and publication ordering. |
| Master-3 | Fixed prepared-journal adoption, restored exact target hashing, added autosave session revalidation, domain-separated paths, wipe commit point, and tombstone integrity. |
| Master-3a | Corrected `memories_snapshot.json`, creation ordering, wipe staging completeness, and target evidence. |
| Master-3b | Added memory conflict policy, tree hashing, gallery boundary, target ownership matrix, and creation identity. |
| Master-3c | Added memory commit boundary, gallery crash reconciliation, and history identity preservation. |

## Appendix B — Superseded wording register

| Superseded wording | Replaced by |
|---|---|
| `memories.json` as canonical chat target | `memories_snapshot.json`; `memories.json` is legacy/separate. |
| `prepared + all-candidate` may adopt | `prepared` never adopts; conflict. |
| Creation may create directory while `prepared` | `writing` durable before directory creation. |
| Autosave revalidates only root/base | Also session, epoch, poison, discard generation. |
| `sha256(rootID||chatID)` paths | Domain-separated path hashes. |
| Wipe without staging gate | Tree-hash staging gate before `root-committing`. |
| Gallery covered by app root | Separate gallery boundary and crash state. |
| History migration may renumber IDs | Preserve IDs; reject missing/duplicate identities. |
| Receipt alone proves memory migration | Receipt and migration journal must agree. |
| Target manifest without target bytes | Exact `F3C-TARGET-1` and `tgt1-` contract. |

Superseded wording is retained in the review sequence and summarized here. This consolidated contract is authoritative where the two differ.

## Appendix C — Consistency audit

- Q1–Q6: no conflict.
- LF-1: no runtime name-folder read/write; chat migration only touches chat-owned memory/history artifacts.
- A1–A13: no conflict. A9/A10/A12 accepted; A13 passing; A11 remains a consumer/availability contract.
- A9/A10/A12 evidence predates this persistence contract and does not prove its new targets or receipts.
- Memory migration: receipt and journal must agree; missing/corrupt journal is an evidence gap, not committed.
- Wipe: staging states are preconditions; root identity remains the only commit point.
- Cast: still unimplemented; no binding is claimed.

## Appendix D — Freeze boundary

This freeze fixes the architecture so implementation can be measured against a stable contract. It does not accept F3-C, authorize P1a, or authorize source changes or tests. F3-C acceptance requires implementation, fault-injection, race, filesystem recovery, native Wails/WebView2 evidence, and the independent Mimo challenge.
