# CORRECTION BASELINE — F3-C Import Audit Findings

**Date:** 2026-11-10
**Status:** historical pre-correction baseline; corrections applied 2026-11-10 — see audit §10
**Input:** `docs/review/AUDIT-F3C-IMPORT-BOUNDARIES-2026-11-10.md`

This file records the confirmed findings and the regression tests that must fail before the
correction pass begins. It is an audit record, not an implementation claim.

## Confirmed findings

| ID | Confirmed? | Source evidence | Baseline behavior |
|---|---|---|---|
| F-01 reference normalization | **Yes** | `internal/contentimport/prepare.go:266`; no `NormalizeReferences` call in `internal/contentimport`; `internal/worldstore/refs.go:36` normalizer exists but is not used by import | Raw arrays are written; name-only entries are not typed as `legacyName`; duplicate `lorebooks` lists remain |
| F-02 registered legacy roots | **Yes** | `internal/contentimport/prepare.go:81-138` checks only destination overlap; no registry query | A registered canonical root can be read as `LegacyRoots` |
| F-03 media adoption | **Yes** | `internal/contentimport/service.go:236-260` compares document and lore, not media; `PlanItem.MediaHash` exists | Same document/lore with different media is classified `adopt` |
| F-04 empty owner IDs | **Yes** | `internal/contentimport/prepare.go:327-331` stamps only when `ownerId` is absent | `{"ownerId":""}` remains empty |
| F-05 staging cleanup/rollback | **Yes** | `internal/contentimport/service.go:398,406` ignore abort errors; `Rollback` ignores failed/quarantined staging; `AbortImport` lacks `writeMu` | Cleanup failure can be invisible; rollback can report `rolled_back` with staging left |
| F-06 Wails bindings | **Yes, exposure nuance** | `build/wails.cmd` documents `generate module`; `build/tools/wails.exe` exists; `frontend/wailsjs` has no `ContentImport`; `main.go` binds the whole `App` | The runtime boundary may work, but checked-in generated bindings are stale and no App-level test proves the new methods |
| F-07 result classification | **Yes** | `internal/contentimport/service.go:727-728` sets `Conflicted` to every `quarantined` item | Unsupported/inline/unavailable quarantines are reported as conflicts |
| F-08 marker record | **Yes** | `internal/worldstore/import.go:21,63-64` uses `owner.json` at `<itemID>.owner.json`; design §15.C names `.howl-import-owner.json` with the item | Design record is stale; implementation is safer (marker outside final folder) |
| F-09 media schema | **Yes** | `internal/contentimport/prepare.go:197-251` accepts any non-inline JSON; `internal/worldstore/import.go:243` only checks `json.Valid` | Malformed/unknown media manifests are copied; inline detection is heuristic |

## Findings proving not defects

None of F-01 through F-09 is withdrawn. F-06 is narrowed: the Go methods and Wails `Bind: []interface{}{app}` are present, so the missing generated symbols are a checked-in artifact defect and an unverified renderer boundary, not proof that the Go methods are unreachable at runtime. Regeneration is still required to make the checked-in bindings API-complete.

## Baseline regression tests to add

These tests are intended to fail on the current tree and pass after correction:

1. `TestImportNormalizesReferencesLosslesslyAndResolvesByID`
2. `TestLosslessUpdatePreservesNormalizedReferences`
3. `TestLegacyRootRejectsRegisteredEquivalentAndNestedPaths`
4. `TestMediaAdoptionRequiresIdenticalManifest`
5. `TestEmptyAndWhitespaceOwnerIDsAreStamped`
6. `TestCleanupFailureIsVisibleAndRollbackRefusesOrRetries`
7. `TestResultConflictsExcludeNonConflictQuarantine`
8. `TestMediaManifestValidationAndUnknownExtensions`
9. `TestMediaMissingVersusPresentIsConflict`

## Corrections that will be attempted

- Implement lossless typed reference normalization in the import boundary, preserving unknown per-reference fields and reporting legacy-name, missing-ID, malformed, duplicate, and unknown-kind classifications.
- Canonicalize and compare legacy roots against every registered root, including equivalent, parent, and child paths.
- Add a store-owned media manifest normalizer/loader; compare media during adoption and validate media before staging.
- Stamp absent or empty owners while preserving foreign owners.
- Lock `AbortImport` with `writeMu`; surface abort failures; inspect failed/quarantined/prepared states during rollback; require no owned artifacts before `rolled_back`.
- Regenerate Wails bindings with the established wrapper if it succeeds; do not hand-edit generated output.
- Correct result aggregation.
- Update design/report records for the actual marker location and the corrected reference/media behavior.

## Baseline limitations

- This file does not claim a correction yet.
- Native A9, A10, and A12 evidence remains outstanding.
- No UI wiring or F3-D work is authorized by this pass.
