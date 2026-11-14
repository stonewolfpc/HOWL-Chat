# F3-C Import — Final Targeted Verification Before Native Acceptance

**Date:** 2026-11-10  
**Scope:** Read-only verification of three boundaries: reference round-trip, post-rename cleanup/recovery, and media normalization/adoption consistency.  
**Change policy:** No code, tests, generated bindings, or UI wiring were changed for this report. No command was executed to produce it. Test names below document existing coverage; they are not a claim that the tests were rerun in this pass.

## Evidence discipline

- `PASS` means the property is established by static inspection of the cited implementation and, where noted, by an existing test.
- `FAIL` means the cited implementation contradicts the requested property, or the requested safety check is absent.
- `NOT VERIFIED` means the property is not established by the available source and tests. It is not a claim that the property is false.
- Existing tests were written after implementation edits. They do not provide independently captured red-before/green-after evidence. This limitation remains.
- No native UI, native execution, or end-to-end acceptance behavior is claimed here.

## Gate status

| Gate | Status | Note |
|---|---|---|
| A9 — native import acceptance | **PASS** | Later isolated native run; see §9 and `VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md`. |
| A10 — native chat load/save and generated-lore synchronization | **PASS** | Later isolated native run; canonical item trees unchanged; see §9. |
| A13 — existing regression gate | **PASSING (prior correction report)** | Not re-executed here. The correction report records the regression gate as passing. |

---

## 1. Reference round-trip

### 1.1 Trace

The import path is:

1. Browser records enter through `loadBrowserSources` at `internal/contentimport/prepare.go:44`.
2. Legacy folders enter through `loadLegacySources` at `internal/contentimport/prepare.go:82` and `readLegacyFolder` at `internal/contentimport/prepare.go:153`.
3. `prepareItem` deep-copies the body, assigns the destination identity, and normalizes supported reference fields at `internal/contentimport/prepare.go:280` using `normalizeReferences` at `internal/contentimport/references.go:24`.
4. `StageImport` writes the document, lore, and media through `internal/worldstore/import.go:166`; `CommitImport` renames staging to the canonical folder at `internal/worldstore/import.go:331` and `:403`.
5. `Store.Load` reads the canonical document through `internal/worldstore/store.go:383`.
6. A lossless update is performed by `Store.Update` at `internal/worldstore/store.go:202`, which calls `mergeDocument` at `internal/worldstore/store.go:469`.
7. The updated document is reloaded through `Store.Load`.
8. The runtime resolution helpers are `NormalizeReferences` at `internal/worldstore/refs.go:40`, `ResolveReferences` at `internal/worldstore/refs.go:121`, and `UnresolvedReferences` at `internal/worldstore/refs.go:138`.

### 1.2 Findings

| # | Claim | Status | Evidence and residual limitation |
|---|---|---|---|
| R1 | Browser and legacy source extraction preserve unknown top-level fields. | **PASS** | `loadBrowserSources` and `readLegacyFolder` retain the decoded document map; `prepareItem` deep-copies it at `internal/contentimport/prepare.go:280`. `TestBrowserPreviewWritesNothingAndImportPreservesUnknowns` asserts `legacyNotes`, nested `howl.customFlag`, and character `custom` survive to `Store.Load`; `TestLegacyImportPreservesLoreAndLeavesSource` asserts `legacyDocField.keep` survives. |
| R2 | Supported references are converted to the typed `{kind,id}` shape during import. | **PASS** | `normalizeReferenceEntry` writes `kind` and `id` at `internal/contentimport/references.go:117-183`; `normalizeReferenceField` assigns the typed list back to the document at `:53-64`. `TestImportNormalizesReferencesLosslesslyAndResolvesByID` asserts typed `scenario`, `character`, and `lorebook` entries after `Store.Load`. |
| R3 | Unknown per-reference fields survive the source-to-store round-trip. | **PASS** | `deepCopyReferenceMap` and `mergeReferenceFields` retain fields not owned by the normalizer; `normalizeReferenceEntry` only sets `kind`, `id`, `title`, and `legacyName` where required. `TestImportNormalizesReferencesLosslesslyAndResolvesByID` asserts `custom`, `legacyUnknown`, and merged `duplicateUnknown` after `Store.Load`. |
| R4 | Unknown fields survive staging and commit. | **PASS** | The prepared map is serialized as the staged document at `internal/worldstore/import.go:166-328` and committed without field rewriting at `:331-419`; `Store.Load` decodes the stored map. The existing round-trip test asserts the fields after load. |
| R5 | Unknown fields survive a lossless document update and reload. | **PASS (source)** | `mergeDocument` starts with `deepCopyMap(existing)` and only applies patch keys at `internal/worldstore/store.go:469-507`; `Store.Update` re-encodes that map at `:202-289`. **Test gap:** `TestImportNormalizesReferencesLosslesslyAndResolvesByID` updates only `title` and then asserts reference list lengths, not the custom per-reference fields. No existing test asserts every unknown field after `Update`. |
| R6 | Supported references resolve strictly by `{kind,id}`. | **PASS** | `ResolveReferences` checks `ref.ID == ""`, then `catalog[ref.Kind][ref.ID]`; it never consults `Title` or `LegacyName` at `internal/worldstore/refs.go:121-136`. `TestResolveReferencesIsIDOnly` covers same-title unknown identity and legacy-name-only non-resolution; `TestResolveReferencesKeepsKindsSeparate` covers cross-kind separation. |
| R7 | No title, name, or `legacyName` fallback exists in runtime resolution. | **PASS** | The resolver uses only `Kind` and `ID`; `Title` is a display cache and `LegacyName` is an unresolved marker (`internal/worldstore/refs.go:12-18`, `:118-136`). The import normalizer converts names to `legacyName` rather than assigning an ID (`internal/contentimport/references.go:117-183`). |
| R8 | Kind mismatches, missing IDs, unknown kinds, malformed entries, and legacy-name-only entries receive explicit import classifications. | **PASS (source)** | `normalizeReferenceEntry` returns `unknown_kind`, `kind_mismatch`, `missing_id`, `malformed`, or `legacy_name` reasons at `internal/contentimport/references.go:117-186`; `normalizeReferenceList` preserves those entries and reports them at `:66-115`; `buildPlan` exposes the issues through `ReferenceIssues` at `internal/contentimport/service.go:208-210` and `internal/contentimport/types.go:78-95`. |
| R9 | Every reference edge case is exercised end-to-end through import, store, update, save, and reload. | **NOT VERIFIED** | `TestImportNormalizesReferencesLosslesslyAndResolvesByID` covers typed entries, legacy-name-only entries, duplicates, unknown fields, `Store.Update`, and `Store.Load`. It does not cover imported `unknown_kind`, `kind_mismatch`, `missing_id`, malformed non-array fields, or a post-update custom-field assertion. `TestNormalizeReferencesRejectsExplicitKindMismatch` exercises the separate `worldstore` normalizer, not the contentimport classifier or the stored document. |
| R10 | An actual runtime consumer resolves imported references through `ResolveReferences` or `UnresolvedReferences`. | **NOT VERIFIED** | Repository search finds no production call site outside `internal/worldstore/refs.go` itself; the only callers are tests (`TestImportNormalizesReferencesLosslesslyAndResolvesByID`, `TestResolveReferencesIsIDOnly`, `TestResolveReferencesKeepsKindsSeparate`, `TestUnresolvedReferencesPreservesOrder`). The resolver is implemented and unit-tested, but there is no production consumer wiring to verify. Whether another runtime path resolves raw document fields independently is outside this import boundary and is not established here. |
| R11 | A legacy source folder remains byte-identical. | **PASS** | `TestLegacyImportPreservesLoreAndLeavesSource` hashes the source tree before and after execute and compares with `reflect.DeepEqual`; `TestCrashBoundariesResumeWithoutSourceMutation` repeats that check across all injected crash boundaries. |
| R12 | A browser source record remains byte-identical. | **NOT VERIFIED** | The browser record is in-memory, so there is no source-file byte image to compare. `TestBrowserPreviewWritesNothingAndImportPreservesUnknowns` marshals the request record after preview and checks validity and absence of `"imported"`, but it never compares the result with the declared `before` value; `before` is unused at `internal/contentimport/service_test.go:75-113`. Source inspection indicates `prepareItem` copies rather than mutates the request map, but the test does not prove full byte equality. |

### 1.3 Reference boundary verdict

**NOT VERIFIED.**

The typed normalization, lossless field retention, staging/commit/load path, and ID-only resolver semantics are implemented and supported by existing tests. The complete requested round-trip is not verified because:

- no production consumer calls the runtime resolver;
- the full edge-case matrix is not exercised end-to-end;
- post-update unknown-field retention is not asserted;
- browser-record byte identity is not compared against its `before` image.

These are evidence gaps, not proof that the stored round-trip is wrong.

---

## 2. Post-rename cleanup and recovery

### 2.1 Trace

- Successful commit: `CommitImport` validates ownership and staged hashes, renames staging to the final path, then performs best-effort cleanup at `internal/worldstore/import.go:331-419`.
- Cleanup after rename: `os.Remove(ownerPath)`, `os.Remove(staging directory)`, and `os.Remove(import staging root)` are explicitly ignored at `internal/worldstore/import.go:413-415`.
- Exact staging ownership: `importOwner` records operation ID, item ID, category, plan hash, document hash, lore hashes, media hash, and revision; `verifyImportOwner` compares them at `internal/worldstore/import.go:551-575`.
- Abort: `AbortImport` validates paths and owner, then refuses if the final item exists at `internal/worldstore/import.go:425-445`.
- Rollback deletion: `RemoveImported` verifies revision and every journaled file hash before removing the canonical folder at `internal/worldstore/import.go:447-515`.
- Recovery: `Resume` at `internal/contentimport/service.go:463-595`; rollback at `:597-736`; final state aggregation at `:753-777`; result reporting at `:779-836`.
- Completion marker: `writeMarker` at `internal/contentimport/journal.go:75-95`; `removeMarker` at `:97-115`. `loadOperation` reads only the journal at `:32-40`. There is no `loadMarker`.

### 2.2 Findings

| # | Claim | Status | Evidence and residual limitation |
|---|---|---|---|
| C1 | A successful canonical commit remains valid if staging cleanup fails. | **PASS (source)** | The rename occurs at `internal/worldstore/import.go:403`; only after it succeeds are cleanup errors discarded at `:413-415`; the function returns `saved` and the canonical revision at `:417-419`. `Store.Load` reads the final folder independently of staging. **Test gap:** no test injects a cleanup failure after a successful rename. `TestCleanupFailureRemainsVisibleAndRollbackLeavesStaging` injects unexpected staging content before rollback, not a post-rename cleanup failure. |
| C2 | Surviving staging-owner evidence is tied to the exact operation and item. | **PASS** | The staging path includes operation ID and item ID at `internal/worldstore/import.go:255-257`; the owner object records both plus plan and hashes at `:270-285`; `verifyImportOwner` checks all of those fields at `:551-575`. `LoadStaging` supplies the same operation/item/plan identity to `Resume` at `internal/contentimport/service.go:498-527`. `TestAuthorityChangeAfterStagingStopsResume` covers changed destination authority; `TestCrashBoundariesResumeWithoutSourceMutation` covers recovery at each staging/commit boundary. |
| C3 | The completion marker is tied to the exact operation and item. | **FAIL** | `marker` contains `OperationID`, `PlanHash`, `Status`, and `CompletedAt` but no item identity at `internal/contentimport/types.go:219-226`. `writeMarker` writes one category-level marker at `internal/contentimport/journal.go:75-95`. It is operation-scoped, not item-scoped. |
| C4 | `Resume` distinguishes committed content from abandoned or foreign staging. | **PASS (source)** | Committed items are re-verified by `verifyJournalItem` and become `ambiguous` if hashes no longer match at `internal/contentimport/service.go:476-486`. Prepared/planned items are loaded only with the exact operation ID, item ID, library ID, and plan hash through `LoadStaging` at `:518-540` and `:548-566`. A foreign owner or hash mismatch cannot be adopted. **Test gap:** no dedicated test places foreign staging at the same operation/item path; the crash and authority tests cover adjacent failure modes. |
| C5 | Cleanup cannot remove a canonical item merely because staging evidence remains. | **PASS** | `AbortImport` refuses when the final path exists at `internal/worldstore/import.go:438-441`. `RemoveImported` requires revision and every journaled file hash to match before deleting the canonical folder at `:447-515`. `Rollback` refuses committed items when staging or staging-owner evidence remains at `internal/contentimport/service.go:653-661`. `TestRollbackRemovesCleanImportAndRefusesChangedLore` covers the changed-content refusal; `TestAdoptedItemBlocksRollback` covers adopted-content refusal. |
| C6 | A later operation cannot claim or delete the previous operation's staging artifacts. | **PASS** | Staging is under `ImportStagingDir/<operationID>/<itemID>` at `internal/worldstore/import.go:255-257`; `AbortImport` validates that exact root, operation, item, and owner at `:425-445` and `:551-575`. A different operation cannot satisfy the owner check. `TestConcurrentImportsClaimIdentityOnce` covers competing imports for one identity. |
| C7 | A later terminal operation cannot overwrite the previous operation's completion marker. | **FAIL** | `writeMarker` writes the single category-level path from `markerPath()` without reading or comparing an existing marker at `internal/contentimport/journal.go:75-95`. `removeMarker` checks `OperationID` before removal at `:97-115`, but `writeMarker` has no equivalent check. `TestCompletedMarkerResumeIsNoOp` starts from an agreeing marker; no test covers a later operation replacing an earlier marker. |
| C8 | Journal/marker disagreement stops safely. | **FAIL** | `Resume` loads only the journal with `loadOperation` at `internal/contentimport/service.go:467` and `internal/contentimport/journal.go:32-40`; it never reads the marker. The only marker read is `removeMarker`, and that is used during successful rollback at `internal/contentimport/service.go:713`. A disagreement is therefore not detected: `Resume` proceeds from the journal, and `writeMarker` may overwrite the disagreeing marker. `TestCompletedMarkerResumeIsNoOp` covers only an agreeing completed marker. |
| C9 | A `cleanup_failed` item leaves the operation in a safe, explicit state. | **PASS** | `finalize` treats `StateCleanupFailed` as ambiguous and partial, setting `OperationAmbiguous` at `internal/contentimport/service.go:753-777`. `Resume` and `Rollback` also set `StateCleanupFailed` with a `cleanup_failed` reason when owned staging cannot be removed at `:497-511` and `:682-695`. |
| C10 | Partial cleanup is accurately represented in aggregate operation result reporting. | **FAIL** | `resultFromOperation` counts `Quarantined` only as `countState(StateQuarantined)` at `internal/contentimport/service.go:815`; `StateCleanupFailed` (and `StateFailed`) are not counted there. Yet `actionForState` maps their default case to `ActionQuarantine` at `:827-835`. A result can therefore report an item with quarantine action while `Plan.Quarantined` remains zero. No test asserts result counts for a `cleanup_failed` operation. |
| C11 | Every path that can change `committed`, `failed`, `cleanup_failed`, `ambiguous`, or `rolled_back` is covered by a test. | **NOT VERIFIED** | Execute transitions are at `internal/contentimport/service.go:324-461`; Resume transitions at `:463-595`; Rollback at `:597-736`. Existing tests cover crash boundaries (`TestCrashBoundariesResumeWithoutSourceMutation`), changed authority (`TestAuthorityChangeAfterStagingStopsResume`), changed-content rollback (`TestRollbackRemovesCleanImportAndRefusesChangedLore`), adopted rollback refusal (`TestAdoptedItemBlocksRollback`), cleanup refusal (`TestCleanupFailureRemainsVisibleAndRollbackLeavesStaging`), and an agreeing terminal marker (`TestCompletedMarkerResumeIsNoOp`). There is no test for post-rename cleanup failure, marker disagreement, later-operation marker replacement, or cleanup-failed result aggregation. |

### 2.3 State-transition audit

| State | Paths that set or change it | Safety observation |
|---|---|---|
| `committed` | Execute after `after_rename`/`after_committed` at `internal/contentimport/service.go:437-452`; Resume prepared item after verified recovery at `:535-540`; Resume committed item after verification at `:476-486` | Committed content is journal-verified; a surviving staging owner can block rollback. |
| `failed` | Execute when staging cannot be created or commit fails at `internal/contentimport/service.go:356-370` and `:405-423`; Resume planned item with no complete staging at `:566-579` | Final existence is checked during Resume and Rollback before cleanup. |
| `cleanup_failed` | Resume when owned staging cannot be verified/removed at `internal/contentimport/service.go:497-511`; Rollback at `:681-695` | Operation becomes `ambiguous`; aggregate `Plan.Quarantined` does not include it. |
| `ambiguous` | Resume when committed hashes fail, a final exists for a non-committed item, or planned staging is incomplete at `internal/contentimport/service.go:476-490`, `:518-530`, `:544-565` | Rollback refuses the operation at `:615-620`. |
| `rolled_back` | Rollback after all committed items are removed at `internal/contentimport/service.go:708-714` | Marker removal checks the operation ID; no canonical deletion occurs without hash/revision verification. |
| `cleanup_done` | Resume or Rollback after owned staging is successfully removed at `internal/contentimport/service.go:512-515` and `:699-705` | Treated as partial in `finalize`; not quarantine-actioned by the default at `:827-835` only because success paths set it explicitly. |

### 2.4 Recovery boundary verdict

**FAIL.**

The core staging ownership and canonical-deletion guards are sound by inspection, and a successful rename is not invalidated by ignored cleanup errors. The boundary fails the requested recovery guarantees because:

1. journal/marker disagreement is never detected;
2. `writeMarker` can overwrite a previous operation's category-level marker;
3. `cleanup_failed` is not represented in the aggregate quarantine count;
4. there is no test for post-rename cleanup failure, marker disagreement, or later-operation marker replacement.

---

## 3. Media normalization and adoption consistency

### 3.1 Trace

- Canonical shape: `{"format":"howl-media-refs-v1","items":[{"mediaId":"...","role":"..."}]}`.
- Legacy reference-only shape: `{"images":[...]}`.
- `NormalizeMediaManifest` validates either shape, rejects inline media, and returns deterministic JSON while preserving unknown fields at `internal/worldstore/media.go:37-96`.
- Legacy source preparation calls the normalizer and records the normalized hash at `internal/contentimport/prepare.go:253-265`.
- Staging calls the same normalizer again and computes `mediaHash` from the normalized bytes at `internal/worldstore/import.go:243-246` and `:278-280`.
- `LoadMedia` normalizes the stored file on read at `internal/worldstore/media.go:101-138`.
- Adoption comparison loads the canonical media and compares normalized bytes with `bytes.Equal` at `internal/contentimport/service.go:239-283`.

### 3.2 Findings

| # | Claim | Status | Evidence and residual limitation |
|---|---|---|---|
| M1 | The same normalization rule is applied during preparation and staging. | **PASS** | Both call `worldstore.NormalizeMediaManifest`: preparation at `internal/contentimport/prepare.go:253`; staging at `internal/worldstore/import.go:243`. |
| M2 | The same normalization rule is applied during load. | **PASS** | `Store.LoadMedia` calls `NormalizeMediaManifest` before returning bytes at `internal/worldstore/media.go:133`. |
| M3 | Adoption compares media through the same normalized representation. | **PASS** | `destinationDisposition` calls `LoadMedia` and compares with `bytes.Equal` at `internal/contentimport/service.go:256-270`; `LoadMedia` has already normalized the canonical bytes. |
| M4 | Media hashing uses the normalized bytes at every boundary that stores a hash. | **PASS** | Preparation hashes normalized media into `sourceHashes` at `internal/contentimport/prepare.go:265`; staging hashes normalized media into `mediaHash` at `internal/worldstore/import.go:278-280`; the owner stores the same hash and `verifyImportOwner` checks it at `:551-575`. |
| M5 | Unknown top-level and per-item media fields survive normalization, staging, commit, and load. | **PASS (source)** | `NormalizeMediaManifest` decodes into a map and re-encodes the whole map, so unknown keys are retained at `internal/worldstore/media.go:37-96`. `TestMediaAdoptionRequiresIdenticalManifest` asserts `unknownTop` survives; its canonical item also carries `custom`, but the test does not assert `custom` after load. `TestMediaManifestValidationAndUnknownExtensions` asserts an extended media item commits but does not inspect its fields after load. **Test gap:** per-item and legacy-envelope unknown fields are not asserted end-to-end. |
| M6 | Canonical and legacy envelopes cannot compare equal when their meaningful references differ. | **PASS** | Adoption uses exact normalized byte equality at `internal/contentimport/service.go:269`. Different `mediaId`, role, provider, or shape produces different bytes. A canonical and legacy envelope also cannot compare equal to one another even when semantically equivalent, which is stricter than the requirement and therefore safe for adoption. |
| M7 | Missing, malformed, unsupported, and inline-binary media remain distinguishable. | **FAIL** | Missing media is represented by absence (`LoadMedia` returns `mediaExists=false` at `internal/worldstore/media.go:119-129`) and by `media_conflict` during adoption at `internal/contentimport/service.go:263-268`. Inline media has the distinct `inline_media_unmigrated` reason at `internal/contentimport/prepare.go:256-260`. However, malformed and unsupported reference-only manifests are both mapped to `unsupported_media_manifest` at `internal/contentimport/prepare.go:260-263`; the code does not preserve a separate malformed outcome. `TestMediaManifestValidationAndUnknownExtensions` confirms a malformed format reports `unsupported_media_manifest`, the same reason used for an unsupported schema. |
| M8 | Inline-binary media is rejected rather than migrated. | **PASS** | `inlineMediaIn` detects inline payload keys and `data:image`/base64 strings at `internal/worldstore/media.go:140-171`; `NormalizeMediaManifest` returns `ErrMediaInline` before writing at `:47-49`. `TestInlineMediaIsQuarantinedNotDropped` and `TestResultHasNoStringlyConflictForInlineMedia` assert quarantine and no write. |
| M9 | Adoption cannot ignore a media difference. | **PASS** | Missing-on-one-side and byte differences both return `media_conflict` at `internal/contentimport/service.go:263-270`. `TestMediaAdoptionRequiresIdenticalManifest` covers identical adoption, a different `mediaId`, and a missing canonical manifest. |
| M10 | Both supported envelope shapes traverse preparation, staging, commit, load, and adoption. | **NOT VERIFIED** | Preparation, staging, and load all accept both shapes through `NormalizeMediaManifest`. `TestLegacyImportPreservesLoreAndLeavesSource` imports a legacy `{"images":[...]}` source and asserts the file exists, but does not adopt it. `TestMediaAdoptionRequiresIdenticalManifest` uses the canonical `items` shape. No existing test performs adoption with a legacy-envelope source and a legacy-envelope canonical destination. |
| M11 | No media bytes are migrated. | **PASS** | Only normalized JSON is written: preparation stores normalized bytes at `internal/contentimport/prepare.go:264`; staging writes `item.Media` as `media.json` at `internal/worldstore/import.go:321-324`; inline media is rejected before staging. The code has no binary copy path. |

### 3.3 Media boundary verdict

**FAIL.**

Normalization, staging, load, hashing, exact adoption comparison, unknown-field retention, and inline rejection are consistent by source inspection. The boundary fails the requested distinction guarantee because malformed and unsupported reference-only manifests collapse to the same `unsupported_media_manifest` reason. End-to-end legacy-envelope adoption and per-item unknown-field retention are also not verified.

---

## 4. Consolidated gaps

| Gap | Boundary | Status | Impact |
|---|---|---|---|
| No production consumer calls `NormalizeReferences`, `ResolveReferences`, or `UnresolvedReferences`. | Reference | **NOT VERIFIED** | Runtime consumer resolution cannot be confirmed from this repository. |
| No end-to-end fixture for `unknown_kind`, `kind_mismatch`, `missing_id`, malformed non-array fields, or post-update custom fields. | Reference | **NOT VERIFIED** | The classifier exists, but the full requested edge matrix is not evidenced. |
| Browser-record byte equality is not compared with the captured `before` value. | Reference | **NOT VERIFIED** | The test checks validity and absence of `"imported"`, not full equality. |
| Completion marker has no read/compare path; `writeMarker` can overwrite. | Recovery | **FAIL** | Journal/marker disagreement does not stop; a later operation can replace the marker. |
| `cleanup_failed` is absent from `Plan.Quarantined` while action is quarantine. | Recovery | **FAIL** | Aggregate result counts are inconsistent with per-item action. |
| No fault injection for post-rename cleanup failure or marker disagreement. | Recovery | **NOT VERIFIED** | The safety properties are source-inspected only. |
| Malformed and unsupported media manifests share `unsupported_media_manifest`. | Media | **FAIL** | The requested four-way distinction is not met. |
| Legacy-envelope adoption and per-item unknown media fields are not asserted. | Media | **NOT VERIFIED** | Source supports them; tests do not demonstrate them end-to-end. |

## 5. What this pass establishes

- The reference normalizer/classifier and the ID-only resolver are implemented with lossless field retention and no title/name fallback.
- Staging ownership is operation/item/plan/hash scoped, and canonical deletion is revision/hash guarded.
- A successful rename is not rolled back by cleanup errors; the canonical item remains readable.
- Media normalization is consistently applied at preparation, staging, load, and adoption, and adoption uses exact normalized bytes.
- Inline media is rejected rather than copied.

## 6. What remains unverified or failing

- Actual runtime consumer resolution for imported references.
- Complete reference edge-case and post-update evidence.
- Journal/marker disagreement handling.
- Later-operation completion-marker isolation.
- Accurate aggregate reporting for `cleanup_failed`.
- Distinct malformed versus unsupported media classification.
- End-to-end legacy media adoption and per-item unknown-field retention.

## 7. Stopping point

This report stopped the targeted verification. No code, tests, generated bindings, UI, or F3-D work
was changed by that pass. A9 and A10 were later exercised natively and passed; see §9. A13 remains
passing per the prior correction report.

---

## 8. Correction addendum — 2026-11-10

This appendix records the authorized correction pass that followed the read-only verification. It
supersedes only the failing or not-verified rows named below. Sections 1–7 remain the unchanged
pre-correction record.

### 8.1 Gate status after correction

| Gate | Status | Evidence |
|---|---|---|
| Import/store focused | **PASS** | `go test ./internal/contentimport ./internal/worldstore -count=1` — both `ok` |
| A13 regression | **PASS** | `go test ./internal/backend/lorebook ./internal/backend/chat ./internal/backend/memory . -count=1` — all `ok` |
| Race detector | **PASS** | `go test -race ./internal/worldstore ./internal/contentimport -count=1` — both `ok` |
| Vet / formatting | **PASS** | `go vet ./...` exit 0; `gofmt -l internal/contentimport internal/worldstore` empty |
| Frontend identity | **PASS** | Character suite 4/4; memory suite 2/2 |
| Wails bindings/build | **PASS** | `build\wails.cmd generate module` exit 0; `build\wails.cmd build` exit 0, binary produced |
| A9 native import acceptance | **OPEN** | Not run in this correction pass. |
| A10 native chat/lore synchronization | **OPEN** | Not run in this correction pass. |

### 8.2 Corrected findings

| Finding | Before | After | Source and regression test |
|---|---|---|---|
| Completion marker identity and isolation | One category-level marker; `writeMarker` could replace it. | Marker path is operation-scoped; marker carries category, operation, plan, destination library/root, sorted item IDs, terminal status, and completion time. Creation is exclusive; replacement requires full marker ownership. | `journal.go` (`markerPath`, `writeMarker`, `sameMarker`, `writeNewFileExclusive`); `TestMarkerWriteRefusesForeignOwner`, `TestMarkerWriteRefusesForeignCategory` |
| Missing marker on a completed operation | `Resume` followed the journal without reading a marker. | `verifyMarker` requires a matching marker for `completed`/`completed_partial`; missing, malformed, stale, contradictory, or foreign markers return `OperationIntegrityFailure` without changing canonical content. | `journal.go` (`verifyMarker`); `service.go` (`Execute`, `Resume`, `Rollback`); `TestMissingMarkerBlocksCompletedResume`, `TestMalformedMarkerBlocksCompletedResume`, `TestContradictoryMarkerBlocksCompletedResume`, `TestStaleMarkerBlocksResume` |
| Marker ownership on rollback | Rollback could leave a completion marker after `rollback_partial`. | `finishRollbackConflict` records `rollback_partial` and removes the marker; `removeMarker` now also validates marker version and category before deletion. | `service.go` (`finishRollbackConflict`); `journal.go` (`removeMarker`); `TestRollbackPartialRemovesCompletionMarker` |
| `cleanup_failed` aggregation | Item action was `quarantine` but `Plan.Quarantined` counted only `quarantined`. | `countQuarantined` counts `quarantined`, `failed`, `ambiguous`, `cleanup_failed`, and `cleanup_complete` when the latter records a cleaned failed/quarantined disposition. A prepared item removed during rollback is `rolled_back`, not `cleanup_complete`, so it does not enter the aggregate. Mixed cleanup failure reports two quarantined items and zero conflicts, with operation status `ambiguous`. | `service.go` (`countQuarantined`, `resultFromOperation`, `Rollback`); `TestCleanupFailureAggregatesAsQuarantine`, `TestPreparedRollbackIsNotAggregatedAsQuarantine` |
| In-progress action/count agreement | `planned`/`prepared` actioned as `quarantine` while excluded from the aggregate. | `planned`/`prepared` map to explicit `pending`; aggregate quarantine remains zero. | `types.go` (`ActionPending`); `service.go` (`actionForState`); `TestInProgressStatesArePendingNotQuarantined` |
| Malformed versus unsupported media | Both collapsed to `unsupported_media_manifest`. | `ErrMediaMalformed` and `ErrMediaUnsupported` are distinct; missing/blank canonical `format` is malformed, a known format with an unsupported shape is unsupported. Preparation emits `malformed_media_manifest` or `unsupported_media_manifest` accordingly. | `worldstore/media.go`; `contentimport/prepare.go`; `TestNormalizeMediaManifestClassifiesMalformedAndUnsupported`, `TestImportClassifiesAllReferenceEdgeCases` |
| Legacy media envelope and unknown fields | Adoption of the legacy shape was not asserted end-to-end. | Legacy `images` sources import, reload, and adopt; top-level and per-image unknown fields survive; different references remain `media_conflict`. Canonical per-item unknown fields also survive. | `TestLegacyMediaEnvelopeAdoptsWithUnknownFields`, `TestNormalizeMediaManifestPreservesUnknownFieldsInBothShapes`, `TestMediaAdoptionRequiresIdenticalManifest` |
| Reference edge matrix | Only the classifier was present; malformed, missing-ID, kind-mismatch, unknown-kind, and post-update field retention lacked one consolidated fixture. | `TestImportClassifiesAllReferenceEdgeCases` covers all five outcomes, verifies typed `{kind,id}` preservation, verifies unknown per-reference fields after a lossless `Store.Update` and reload, and checks the browser source map is unchanged. | `contentimport/references.go`; `TestImportClassifiesAllReferenceEdgeCases`; `TestBrowserPreviewWritesNothingAndImportPreservesUnknowns` |

### 8.3 Marker state transitions

- `executing`, `ambiguous`, `rolled_back`, and `rollback_partial` require **no** completion marker. A surviving marker in one of these states is stale evidence and produces `integrity_failure`.
- `completed` and `completed_partial` require a marker whose category, operation ID, plan hash, destination identity, item set, and status all match the journal. Missing or malformed evidence is never inferred into completion.
- `Execute` and `Resume` call `verifyMarker` before recovery and call `syncMarker` after finalization. `Rollback` calls `verifyMarker` before deletion and removes the marker after a successful or partial rollback.
- The read-only journal listing validates each loaded operation's marker and reports `integrity_failure` in the returned copy without mutating the journal.
- `writeMarker` creates exclusively when absent. If a marker exists, it may only be replaced when `sameMarker` proves full ownership; otherwise it returns `ErrMarkerConflict` and leaves the foreign marker untouched.

### 8.4 Reference boundary disposition

The import/store slice now has evidence for every requested reference classification and lossless
field retention. There is still **no production consumer** of `NormalizeReferences`,
`ResolveReferences`, or `UnresolvedReferences` in this F3-C slice. That absence is an explicit
current integration boundary, not a claim of end-to-end runtime resolution. No name/title fallback
or broader reference contract was added.

### 8.5 Remaining limitations

- A9 and A10 later passed in an isolated native run (§9); that run did not exercise the import UI or model inference.
- No independently captured red-before/green-after run; tests were authored with the correction.
- A post-rename cleanup error is still source-inspected rather than fault-injected; successful rename remains canonical even when cleanup is best-effort.
- No fsync/power-loss guarantee and no cross-process writer serialization.
- Marker/root ownership remains a desktop single-process contract.

## 9. Native acceptance addendum — 2026-11-10

A9 and A10 were exercised against the real Wails bindings in an isolated native WebView2 process
using a marked root outside the project. A9 observed `ready=2`, `quarantined=3`, and `writes=2`;
execute completed partially, explicit resume was a no-op, rollback removed only journaled creations,
a separate idempotency fixture adopted without writes, conflict preview reported `identity_conflict` with zero writes, and the
LF-1 source tree was byte-identical before and after. A10 created canonical world/scenario items,
loaded chat data by stable ID, called `LoadChat`, `SaveChat`, `SaveChatData`, and generated-lore
synchronization, and left both canonical item trees byte-identical. The generated-lore path was
reached through a build-tagged native seam; the production binary contained no seam symbols.

Full evidence: `VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md` and
`VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`. A12 passed separately.
