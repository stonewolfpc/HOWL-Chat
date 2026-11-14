# AUDIT — F3-C Import Boundaries

**Date:** 2026-11-10
**Type:** read-only audit; no implementation, UI wiring, or F3-D changes
**Auditor:** working-tree inspection, focused source tracing, and existing test execution
**Scope:** reference compatibility, recovery ownership, Wails binding boundaries, and deliberate design refinements against the frozen F3-C contracts

## 1. Verdict

> **Post-audit native status:** A9 and A10 later passed in an isolated native run, and A12 later
> passed in a separate capture build. See `VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md` and
> `VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`. The findings, gate table, and “not closed”
> statements below remain the audit-time record.

**F3-C import was not closed at the time of this audit.** The code compiled and the existing import/store tests passed, but this audit found unresolved contract and boundary failures:

- imported reference arrays are not normalized to the frozen typed shape;
- a registered canonical root is not rejected as a legacy source root;
- media is ignored when deciding whether an existing item can be adopted;
- an empty `ownerId` is not always stamped with the mapped item identity;
- failed/quarantined staging can outlive an operation while rollback still reports success;
- the generated Wails binding surface does not expose the new methods;
- execution results mislabel all quarantined items as conflicts.

The passing unit tests do not test these boundaries. At audit time native A9, A10, and A12
evidence remained outstanding; all three later passed. A13 remains a regression gate rather than
import acceptance.

## 2. Gate status

| Gate | Status | Basis |
|---|---|---|
| F3-C import boundary | **NOT CLOSED** | Findings F-01 through F-09 below |
| Reference compatibility | **FAIL** | No normalizer call in the import conversion; frozen §4/§6 conflict |
| Recovery ownership | **PARTIAL / OPEN** | Commit-time ownership checks are strong, but failed/quarantined staging cleanup and abort/commit serialization are incomplete |
| Wails binding boundary | **OPEN** | Go methods exist, but generated `frontend/wailsjs` bindings contain no `ContentImport` symbol and no App-level test invokes them |
| A9 import acceptance | **OPEN** | No native import run; findings above are not narrowed by native evidence |
| A10 single-writer runtime | **OPEN** | Source audit is separate; no native chat load/save/generated-lore runtime proof was run |
| A13 regression | **PASS only as regression** | Existing suites pass; this does not close A9 or A10 |
| A12 native binary evidence | **AUDIT-TIME OPEN; LATER PASS** | Not performed in this audit; later recorded in `VERIFICATION-F3C-NATIVE-A12-2026-11-10.md` |

A9, A10, and A13 are intentionally separate. A green A13 does not erase the failures above.

## 3. Positive findings verified in the working tree

These are real strengths, but they do not offset the findings that follow.

- `contentimport` does not call `library.Registry.WriteFile` or `WriteNewFile` for item content; canonical writes are routed through `worldstore` import methods.
- `CommitImport` takes the process-wide `writeMu`, re-reads the destination registration, requires save and write roles, requires the same root, requires `Lookup` completeness, requires final absence, and verifies staged hashes before rename (`internal/worldstore/import.go:329`).
- Staging is operation-scoped and beside the final folder; the first staging segment is not a valid item identity (`internal/worldstore/import.go:20-21`, `63-64`).
- `AbortImport` verifies the owner marker and refuses to remove when a final item exists (`internal/worldstore/import.go:423`).
- `RemoveImported` requires the journaled revision and every journaled file hash, rejects unexpected files, and is reached only from explicit rollback (`internal/worldstore/import.go:447`).
- Crash tests cover after journal, after stage, incomplete staging, after prepared, after rename, after committed, and terminal marker no-op; source hashes are compared around the crash cases (`internal/contentimport/service_test.go:324`, `internal/contentimport/edge_test.go:100`, `internal/contentimport/marker_test.go:10`).
- The focused and regression commands listed in §7 pass.

## 4. Findings

### F-01 — HIGH — Reference normalization is not implemented

**Status: FAIL against frozen contract.**

The frozen reference contract requires imported legacy references to be converted to the typed reference shape, with name-only references represented as unresolved `legacyName` references, and canonical writes to normalize duplicate `howl.lorebooks`/`lorebooks` lists:

- `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md:120-133` (typed shape, name-only import rule, normalized canonical list).
- `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md:154` (preflight must normalize references).
- `docs/review/DESIGN-F3C-IMPORT-2026-11-10.md:155-156` (known reference fields use `worldstore.NormalizeReferences`; transformations reported in the plan).

Observed implementation:

- `internal/contentimport/prepare.go:266` `prepareItem` deep-copies the source body and sets the destination `id`; it never calls `worldstore.NormalizeReferences`.
- There is no call to `NormalizeReferences`, `ResolveReferences`, or `UnresolvedReferences` anywhere in `internal/contentimport` (verified by repository search). The normalizer exists only in `internal/worldstore/refs.go:36`.
- Source arrays such as `scenarios`, `characters`, `lorebooks`, and `howl.lorebooks` are preserved verbatim. Name-only entries are not converted to `{kind, id: "", legacyName: ...}`.
- No plan item reports reference transformations, duplicate collapse, or unresolved name-only references.
- No import test asserts imported reference shape; `internal/worldstore/refs_test.go` tests the normalizer in isolation, not the import conversion.

Impact: imported canonical documents can contain legacy/name-only arrays that do not match the frozen canonical reference shape. Consumers may continue to see raw title-only entries, and there is no import-time record of which references were transformed or left unresolved. This is a contract failure even though the raw data is preserved losslessly.

The design document contains a same-day §15.F refinement saying raw references are intentionally preserved. That refinement does not satisfy the frozen change-control rule at `FREEZE:6` (“any change to a frozen item requires a new dated decision record”). It is an implementation note, not a separate decision record, and it narrowly addresses raw preservation while leaving the required typed/legacy-name transformation unimplemented. This is exactly the kind of quiet weakening the audit must stop.

**Required resolution before A9 closure:** either implement the typed reference conversion and report transformations, or produce an explicit dated decision record that supersedes the frozen reference rule and states the consumer/behavioral consequences. Do not mark F-01 PASS because raw fields survived.

### F-02 — HIGH — Legacy roots are not checked against the registry

**Status: FAIL against LF-1 source selection.**

The request type documents legacy roots as “explicit, ephemeral source selections and are never registered” (`internal/contentimport/types.go:34-38`), and LF-1 requires an explicitly selected, unregistered legacy input root.

The only source validation in `loadLegacySources` is:

- path normalization and regular-directory checks;
- rejection when the source overlaps the destination root (`internal/contentimport/prepare.go:81-106`);
- directory enumeration and folder reads.

There is no check against `s.reg` or any registration registry. A caller can supply a registered canonical `Worlds/` or `Scenarios/` root, or a root containing one, as `LegacyRoots`; the code will read it as legacy input. The destination root is protected from overlap, but other registered roots are not.

Impact: this violates the LF-1 source boundary and can cause a canonical root to be treated as migration input. It does not write the source, but it can import canonical items from a root that the contract says must not be selected as legacy migration input. It can also create duplicate/ambiguous identity situations across registered roots. No existing test covers this.

**Required resolution:** reject any legacy root that is registered in the content registry (or is inside/contains a registered root), with an explicit plan reason; add a test using a registered sibling root.

### F-03 — HIGH — Media is ignored when adopting an existing item

**Status: FAIL for idempotent content handling.**

`destinationDisposition` compares the destination document and lore, then returns `ActionAdopt` without comparing `media.json`:

- `internal/contentimport/service.go:236` compares `search.Item.Document` and `prepared.document`;
- it iterates lore files and compares semantic hashes;
- it never calls a media loader or compares `prepared.media` with destination media;
- `PlanItem.MediaHash` exists (`internal/contentimport/types.go:97`) and is included in the plan hash, but is not used in the adoption decision.

Failure cases:

- Source has a reference-only `media.json`, destination has no `media.json`: the item is adopted, and the source media references are silently absent from the canonical item.
- Source has no media, destination has media: the item is adopted, and the destination media is retained even though the source content differs.
- Source and destination have different media references: the item is adopted without reporting a conflict.

The frozen import contract says idempotence is by identity plus content hash and that `media.json` is part of item content (`FREEZE:150-158`). A document/lore-only adoption is not identity-plus-content idempotence. `StageImport` and `RemoveImported` track media hashes, so the omission is specifically in the adoption path.

**Required resolution:** compare reference-only media semantically and return `media_conflict`/`identity_conflict` when it differs; adopt only when document, lore, and media all match. Add source-with-media and destination-media mismatch tests.

### F-04 — HIGH — Empty `ownerId` is not stamped

**Status: FAIL against design §6.**

`prepareLoreFile` stamps only when the `ownerId` key is absent:

```go
if _, ok := entry["ownerId"]; !ok {
    if owner, ok := entry["owner"].(string); ok && strings.TrimSpace(owner) == "" {
        entry["ownerId"] = destinationID
    } else if _, hasOwner := entry["owner"]; !hasOwner {
        entry["ownerId"] = destinationID
    }
}
```

`internal/contentimport/prepare.go:327-331`.

An entry such as `{"id":"lb_1","ownerId":""}` has the key present, so the condition is false and the empty owner remains empty. The same applies to whitespace-only `ownerId`. The design requires an empty owner to be stamped with the mapped canonical item ID (`docs/review/DESIGN-F3C-IMPORT-2026-11-10.md:171`). The current test uses an absent owner (`internal/contentimport/service_test.go:277`), not an empty one, so it misses this case.

Impact: imported lore can remain unowned and be excluded from prompt eligibility even though the item mapping established ownership. That is a functional ownership loss.

**Required resolution:** treat absent, empty, and whitespace-only `ownerId`/`owner_id`/`owner` consistently; stamp the mapped ID and preserve unknown fields. Add tests for `ownerId: ""`, whitespace, and `owner_id` variants.

### F-05 — MEDIUM/HIGH — Failed or quarantined staging can outlive rollback

**Status: PARTIAL for recovery ownership.**

The commit path is guarded, but the failure path is not:

- On commit error, `Execute` sets `StateFailed` and ignores the `AbortImport` error (`internal/contentimport/service.go:395-400`).
- On a non-`saved` commit result, it sets `StateQuarantined` and again ignores the abort error (`internal/contentimport/service.go:403-407`).
- `Resume` skips `StateQuarantined` and `StateFailed` (`internal/contentimport/service.go:458`).
- `Rollback` phase one inspects adopted, ambiguous, prepared, planned, and committed items, but not failed or quarantined staging (`internal/contentimport/service.go:540-620`).
- If no collision/ambiguity is detected, phase two leaves those items alone and can set `OperationRolledBack` and remove the marker.
- `AbortImport` is not serialized by `writeMu`; a concurrent rollback/abort and resume/commit can race on the same staged item. `CommitImport` is locked, but `AbortImport` only checks the final path and removes staging (`internal/worldstore/import.go:423-440`).

Impact: an operation can report successful rollback while leaving operation-scoped staging, and ignored abort errors are not surfaced. The canonical item may be safe, but staging ownership and cleanup evidence are incomplete, which is contrary to the audit requirement for journaled recovery and explicit cleanup. The design says incomplete owned staging should be removed only when its marker matches; the implementation neither cleans nor reports a matched-but-failed cleanup as an item outcome.

**Required resolution:** treat abort failure as an item state (not an ignored error), record leftover staging explicitly, make rollback report it as `rollback_partial`/`ambiguous` instead of `rolled_back`, and serialize `AbortImport` with commit/remove or document why the race is impossible. Add a test where `AbortImport` fails and rollback still has staging.

### F-06 — MEDIUM — Wails generated bindings are absent; binding boundary is unproven

**Status: OPEN.**

The Go methods are present and correctly routed:

- `PreviewContentImport` — `app-imports.go:12`
- `ExecuteContentImport` — `app-imports.go:22`
- `ResumeContentImport` — `app-imports.go:32`
- `RollbackContentImport` — `app-imports.go:42`
- `ListContentImportJournal` — `app-imports.go:51`
- Services are initialized in `app.go:254-262`.

However:

- A repository search for `ContentImport` under `frontend/wailsjs` returned **zero** matches. The generated `App.js`/`App.d.ts` surface does not expose the new methods.
- There is no App-level test invoking these bindings. Existing binding tests cover lorebook import, not content import.
- No native UI calls these methods; the import page is explicitly out of scope.

Impact: until the Wails generation step runs, the documented renderer binding boundary does not exist in the checked-in frontend artifacts. The Go layer may build, but the native/renderer contract has not been exercised. This is OPEN, not a claim that the Go methods are wrong.

**Required resolution:** regenerate the Wails bindings as part of the authorized build step (or explicitly document that generation is deferred and that no renderer may consume them), and add an App-level binding test or a native smoke call. This does not authorize UI wiring in this audit.

### F-07 — MEDIUM — Result reporting counts every quarantine as a conflict

`resultFromOperation` sets:

```go
Conflicted: countState(operation.Items, StateQuarantined),
```

at `internal/contentimport/service.go:728`.

`StateQuarantined` includes unsupported files, inline media, unmapped owners, destination unavailable, and other non-conflict reasons. The preflight builder correctly increments `Conflicted` only for `identity_conflict` and `lore_conflict` (`internal/contentimport/service.go:213-221`), but execute/resume results do not preserve that distinction.

Impact: a UI consuming `Result.Plan.Conflicted` will report ordinary quarantine as conflicts, weakening the explicit outcome reporting. It does not affect canonical writes.

**Required resolution:** count only conflict reasons in `resultFromOperation`, and preserve `Unsupported`/`Quarantined` counts. Add a result-reporting test with an unsupported file or inline media.

### F-08 — LOW — Ownership marker path/name differs from the design record

The design’s staging section names `.howl-import-owner.json` as the ownership marker and describes it with the staging item (`docs/review/DESIGN-F3C-IMPORT-2026-11-10.md:474-483`). The implementation uses:

- `importOwnerFile = "owner.json"` (`internal/worldstore/import.go:21`);
- `OwnerPath = .../<itemID>.owner.json`, outside the item folder (`internal/worldstore/import.go:63-64`);
- owner file is created with mode `0600` and removed during commit/abort (`internal/worldstore/import.go:247-291`, `413`, `440`).

The outside-folder placement is safer for publication isolation, but the checked-in design record is stale. This is not a frozen-contract failure by itself; it is a documentation/implementation mismatch that must be reconciled before claiming the design describes the implementation.

**Required resolution:** update the design record to the actual marker name/location, or change the implementation to the recorded design. Do not leave both descriptions as authoritative.

### F-09 — LOW — `media.json` schema is not validated

`readLegacyFolder` rejects inline byte-like keys and accepts any other valid JSON object as media:

- `internal/contentimport/prepare.go:197-215` decodes and calls `inlineMediaIn`;
- `worldstore.StageImport` only checks `json.Valid` for media (`internal/worldstore/import.go:243`).

The frozen media contract requires the `howl-media-refs-v1` shape with `items` and roles (`FREEZE:51-61`). Arbitrary reference-only JSON can be copied without schema validation, and `inlineMediaIn` only recognizes specific keys (`data`, `base64`, `blob` under selected media keys). A string data URL in an `images` array, or bytes under an unrecognized key such as `imageData`, may not be detected.

Impact: the import can preserve a malformed media envelope or miss a different inline-byte representation. This is a data-contract gap, not a canonical-write ownership failure.

**Required resolution:** validate the media envelope/roles and make inline-media detection exhaustive or explicitly document its supported shapes. Add malformed-schema and alternate inline-byte tests.

## 5. Reference compatibility matrix

| Frozen expectation | Observed | Result |
|---|---|---|
| `scenarios`, `characters`, `howl.lorebooks`, `lorebooks` normalized to typed `{kind,id,title,legacyName}` | Raw source arrays copied verbatim by `prepareItem` | **FAIL** |
| Name-only references become `{kind, id: "", legacyName}` and are reported unresolved | No conversion; no plan report | **FAIL** |
| Duplicate `howl.lorebooks`/`lorebooks` written as one normalized list | Both remain as supplied | **FAIL** |
| Unknown per-reference fields preserved | Raw arrays preserve them, but not in the typed shape | **PARTIAL** |
| ID-only resolution at consumption | `worldstore.ResolveReferences` is ID-only, but import does not feed it normalized refs | **PARTIAL** |
| Import test exercises reference conversion | No import reference test; only isolated normalizer tests | **FAIL (coverage)** |

## 6. Recovery ownership matrix

| State/evidence | Expected | Observed | Result |
|---|---|---|---|
| `planned`, no staging/final | fail without write | handled | PASS |
| `planned`, complete owned staging | verify and commit | handled | PASS |
| `planned`, incomplete/unknown staging | ambiguous; no canonical write | handled; staging left in place | PASS for safety / OPEN for cleanup |
| `prepared`, final absent | revalidate and commit | handled | PASS |
| `prepared`, final present | verify and record committed | handled | PASS |
| `committed`, hashes intact | no-op | handled | PASS |
| `committed`, hashes changed | ambiguous, no deletion | handled | PASS |
| `failed` after commit failure | report abort outcome; no orphan staging | abort error ignored; item failed | **FAIL** |
| `quarantined` after non-saved commit | report abort outcome; no orphan staging | abort error ignored; rollback ignores it | **FAIL** |
| adopted item rollback | refuse | handled | PASS |
| changed lore rollback | refuse, retain item | handled | PASS |
| concurrent abort and commit | serialized or provably safe | abort is not under `writeMu` | **OPEN** |

## 7. Wails binding matrix

| Method | Go `App` method | Generated `frontend/wailsjs` | App-level test | Native use |
|---|---|---|---|---|
| `PreviewContentImport` | present `app-imports.go:12` | absent | absent | none |
| `ExecuteContentImport` | present `app-imports.go:22` | absent | absent | none |
| `ResumeContentImport` | present `app-imports.go:32` | absent | absent | none |
| `RollbackContentImport` | present `app-imports.go:42` | absent | absent | none |
| `ListContentImportJournal` | present `app-imports.go:51` | absent | absent | none |

The service initialization is present at `app.go:254-262`. The missing generated artifacts and absent App-level tests are the binding-boundary gap; they do not by themselves prove the Go code is wrong.

## 8. Test evidence

Executed during this audit with `GOFLAGS=-buildvcs=false`:

| Command | Result |
|---|---|
| `go test ./internal/contentimport ./internal/worldstore -count=1` | both `ok` |
| `go test ./internal/backend/lorebook ./internal/backend/chat ./internal/backend/memory . -count=1` | all `ok`; A13 exit 0 |
| `go test ./internal/worldstore ./internal/backend/chat ./internal/contentimport . -count=1` | all `ok` |
| `go test -race ./internal/worldstore ./internal/contentimport -count=1` | both `ok` |
| `go vet ./...` | exit 0 |
| `node tests/frontend/character-identity-ownership.test.js` | fail 0; prior reported 4/4 |
| `node tests/frontend/memory-identity-ownership.test.js` | fail 0; prior reported 2/2 |

No test is added or changed by this audit. Passing tests cover the implemented happy paths and several crash boundaries, but none of the following: typed reference conversion, registered legacy-root rejection, media adoption comparison, empty `ownerId`, failed-staging cleanup, generated Wails symbols, or native import.

## 9. Required next evidence before closure

1. Fix or formally supersede F-01 with a dated decision record; do not treat §15.F alone as change control.
2. Reject registered/overlapping legacy roots and test it.
3. Compare media during adoption and test references-only media conflicts.
4. Stamp empty/whitespace owners consistently and test all owner-key variants.
5. Make abort failures and leftover staging visible; make rollback outcomes match actual staging state; serialize abort with commit/remove.
6. Regenerate or explicitly defer the Wails bindings and prove one App-level call path.
7. Run isolated native A9 import evidence.
8. Run native A10 evidence proving chat load/save and generated-lore synchronization do not change canonical item trees.
9. Keep A13 as regression evidence only; do not use it to close A9 or A10.

No gate in this audit is closed by compilation alone.

## 10. Correction pass disposition — 2026-11-10

The original findings above remain the audit history. This section records the focused correction
pass authorized after the audit; it does not erase or relabel the findings.

| Finding | Disposition | Correction and focused evidence | Remaining limitation |
|---|---|---|---|
| F-01 | **CORRECTED — frozen contract implemented, not amended** | `internal/contentimport/references.go` normalizes known fields losslessly to typed references, preserves unknown document/per-reference fields, collapses duplicates, and reports legacy-name, missing-ID, malformed, duplicate, unknown-kind, and kind-mismatch cases. Runtime `worldstore.NormalizeReferences` skips explicit kind mismatch. Tests: `TestImportNormalizesReferencesLosslesslyAndResolvesByID`, `TestNormalizeReferencesRejectsExplicitKindMismatch`. | Unknown-kind entries remain preserved but unresolved; no name/title fallback exists. |
| F-02 | **CORRECTED** | `loadLegacySources` canonicalizes the source and checks every registered root, including parent/child overlap. Destination overlap keeps its original reason; other registered/equivalent roots return `source_registered`. Test: `TestLegacyRootRejectsRegisteredEquivalentAndNestedPaths`. | A source that is simultaneously the destination reports `source_destination_overlap` first; both reject with no writes. |
| F-03 | **CORRECTED** | `destinationDisposition` calls `Store.LoadMedia` and compares normalized manifests. Missing, present-but-different, and unreadable destination media are `media_conflict`; only identical document, lore, and media adopt. Test: `TestMediaAdoptionRequiresIdenticalManifest`. | Media bytes are not migrated; F3-F remains out of scope. |
| F-04 | **CORRECTED** | `prepareLoreFile` stamps absent, empty, and whitespace `ownerId`; non-empty foreign owners remain `unmapped_owner` and are not rewritten. Test: `TestEmptyAndWhitespaceOwnerIDsAreStamped` plus existing unmapped-owner coverage. | `owner`/`owner_id` are validated as foreign evidence but only `ownerId` is written as the canonical field. |
| F-05 | **CORRECTED** | `AbortImport` takes `writeMu`, verifies the marker and exact staged file set, and refuses unexpected files. Abort errors become `cleanup_failed`/`ambiguous`; resume and rollback inspect failed, quarantined, prepared, cleanup, and committed states; rollback never reports `rolled_back` with owned staging/marker evidence. Tests: `TestCleanupFailureRemainsVisibleAndRollbackLeavesStaging`, `TestAbortImportRefusesUnexpectedStagingContents`, `TestAbortAndCommitSerializeToOneOutcome`. | Commit cleanup after a successful rename remains best-effort; surviving evidence is visible and blocks `rolled_back`. |
| F-06 | **ARTIFACT DEFECT CORRECTED; NATIVE EXPOSURE LATER PASSED** | `build\wails.cmd generate module` exited 0 and regenerated `frontend/wailsjs`; all five methods are present in `App.js`/`App.d.ts` with `contentimport` request/result types. `build\wails.cmd build` exited 0 and produced `build\bin\howl-chat.exe`. | At audit time no native UI import run had occurred; A9 later passed in an isolated native run. Existence of bindings alone was not native runtime acceptance. |
| F-07 | **CORRECTED** | `resultFromOperation` counts only `identity_conflict`, `lore_conflict`, and `media_conflict` in `Conflicted`; unsupported and other quarantines remain `Quarantined`/`Unsupported`. Result also reports committed writes. Tests: `TestResultConflictsExcludeNonConflictQuarantine`, `TestResultHasNoStringlyConflictForInlineMedia`. | Per-item reasons remain the authoritative detail; aggregate counts are summary only. |
| F-08 | **CORRECTED DOCUMENTATION** | Design §15.C now records the actual `owner.json` at `<itemID>.owner.json` beside the staging folder and the exact-contents cleanup rule. No marker rename or journal invalidation occurred. | Existing journals retain the same shape. |
| F-09 | **CORRECTED** | `worldstore.NormalizeMediaManifest` validates canonical `howl-media-refs-v1`/`items` and legacy `images`, preserves unknown fields, detects inline strings recursively, and is used by both legacy reading and staging. Tests: `TestMediaManifestValidationAndUnknownExtensions`, existing legacy-media import test. | Supported shapes only; arbitrary unsupported envelopes quarantine rather than migrate. |

### Evidence run for this disposition

Executed on 2026-11-10 from the project root with `GOFLAGS=-buildvcs=false` unless noted:

| Command | Result |
|---|---|
| `go test ./internal/contentimport ./internal/worldstore -count=1` | both `ok` |
| `go test ./internal/backend/lorebook ./internal/backend/chat ./internal/backend/memory . -count=1` | all `ok`; A13 exit 0 |
| `go test -race ./internal/worldstore ./internal/contentimport -count=1` | both `ok`; race detector exit 0 |
| `go vet ./...` | exit 0 |
| `node tests/frontend/character-identity-ownership.test.js` | 4/4 passed, exit 0 |
| `node tests/frontend/memory-identity-ownership.test.js` | 2/2 passed, exit 0 |
| `build\wails.cmd generate module` | exit 0; `ContentImport` symbols present |
| `build\wails.cmd build` | exit 0; native binary produced |

**Verification limitation:** the baseline was recorded before implementation edits, but the focused
regression tests were authored after the first edit. This pass therefore does not claim an
independently captured red-before/green-after run for each finding. The pre-correction source
evidence is in §3–§6, and the post-correction tests are listed above.

No gate is closed by compilation alone. At this point **A9 was still open** (no native import
acceptance) and **A10 was still open** (no native chat load/save and generated-lore synchronization
evidence). Both later passed; see §12. **A13 remains a regression gate and passed.** F3-C import is
not declared closed by this audit.

## 11. Final targeted-verification corrections — 2026-11-10

This disposition addresses only the defects found after §10. It preserves the findings and evidence
above; where a row below conflicts with an earlier row, this section is the later disposition.

| Finding | Disposition | Evidence and residual limitation |
|---|---|---|
| C3/C7 marker identity and isolation | **CORRECTED** | Completion markers are operation-scoped and record category, operation, plan, destination library/root, sorted item IDs, terminal status, and time. Creation is exclusive; replacement requires full ownership. `writeMarker`, `verifyMarker`, and `removeMarker` enforce this. Tests: `TestMarkerWriteRefusesForeignOwner`, `TestMarkerWriteRefusesForeignCategory`, `TestCompletedMarkerResumeIsNoOp`. |
| C8 journal/marker disagreement | **CORRECTED** | `Execute`, `Resume`, and `Rollback` read and verify the marker before recovery; the read-only `Operations` listing also surfaces `integrity_failure`. Missing, malformed, stale, contradictory, and foreign markers do not delete canonical content. Tests: `TestMissingMarkerBlocksCompletedResume`, `TestMalformedMarkerBlocksCompletedResume`, `TestContradictoryMarkerBlocksCompletedResume`, `TestStaleMarkerBlocksResume`. |
| C10 cleanup aggregation | **CORRECTED** | `countQuarantined` includes quarantined, failed, ambiguous, cleanup-failed, and cleanup-complete states. `cleanup_complete` is limited to cleaned failed/quarantined dispositions; a prepared item removed during rollback is `rolled_back` and is not counted. Mixed cleanup reports two quarantined, zero conflicts, status `ambiguous`. Tests: `TestCleanupFailureAggregatesAsQuarantine`, `TestPreparedRollbackIsNotAggregatedAsQuarantine`. |
| In-progress action/count mismatch | **CORRECTED** | `planned`/`prepared` now return `pending`, not `quarantine`; aggregate quarantine remains zero. `TestInProgressStatesArePendingNotQuarantined`. |
| C11 marker fault coverage | **CORRECTED FOR MARKERS; PARTIAL FOR POST-RENAME CLEANUP** | Isolated tests cover missing, malformed, contradictory, stale, foreign-collision, and rollback-marker removal. A post-rename cleanup error remains source-inspected rather than injected. |
| M5/M7/M10 media classification | **CORRECTED** | Malformed and unsupported manifests are separate (`malformed_media_manifest`/`unsupported_media_manifest`); missing format is malformed; legacy and canonical shapes preserve unknown fields; legacy adoption is tested. Tests: `TestNormalizeMediaManifestClassifiesMalformedAndUnsupported`, `TestNormalizeMediaManifestPreservesUnknownFieldsInBothShapes`, `TestLegacyMediaEnvelopeAdoptsWithUnknownFields`. |
| R9/R12 reference evidence | **CORRECTED** | Unknown-kind, kind-mismatch, missing-ID, malformed, and legacy-name-only outcomes are classified and preserved; unknown per-reference fields survive lossless update/reload; browser records compare unchanged. `TestImportClassifiesAllReferenceEdgeCases`, `TestBrowserPreviewWritesNothingAndImportPreservesUnknowns`. |
| Runtime reference consumers | **STILL AN EXPLICIT BOUNDARY** | No production consumer of the resolver exists in this F3-C slice. This is documented as an integration boundary; no name/title fallback, editor wiring, or prompt wiring was added. |

### Correction evidence run

| Command | Result |
|---|---|
| `go test ./internal/contentimport ./internal/worldstore -count=1` | both `ok` |
| `go test ./internal/backend/lorebook ./internal/backend/chat ./internal/backend/memory . -count=1` | all `ok`; A13 exit 0 |
| `go test -race ./internal/worldstore ./internal/contentimport -count=1` | both `ok`; exit 0 |
| `go vet ./...` | exit 0 |
| `gofmt -l internal/contentimport internal/worldstore` | no output |
| Frontend identity tests | 4/4 and 2/2, exit 0 |
| `build\wails.cmd generate module` | exit 0 |
| `build\wails.cmd build` | exit 0; binary produced |

The correction pass did not close A9 or A10. No native import, chat, or synchronization acceptance
was performed during that pass. A13 remains green and separate. A later native run closed A9 and A10;
see §12 and `VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md`.

## 12. Native acceptance addendum — 2026-11-10

| Gate | Status | Evidence |
|---|---|---|
| A9 native import | **PASS** | Isolated native bindings: preview `ready=2`/`quarantined=3`/`writes=2`; execute `completed_partial`; resume no-op; rollback `rolled_back`; idempotent adopt; conflict `identity_conflict` with zero writes; LF-1 source byte-identical before/after. |
| A10 native single writer | **PASS** | Canonical world/scenario items created; `LoadChat`, `SaveChat`, `SaveChatData`, canonical load, and generated-lore synchronization completed; canonical trees byte-identical before/after; chat files only under `Chats/A10Chat`. |
| A12 native binary/prompt boundary | **PASS** | Separate capture build; see `VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`. |

Drivers and evidence are under `build/tmp/a9-acceptance/`, `build/tmp/a10-acceptance/`, and
`build/tmp/a12-acceptance/`. The production binary was rebuilt without test tags and contained none
of `__howl_native_test_seam`, `HOWL_NATIVE_TEST_SEAM`, `X-Native-Test-Seam-Token`,
`NativeTestSaveChatData`, `native_test_seam`, `promptcapture`, `prompt-capture`, or
`GetPromptCapture`.
