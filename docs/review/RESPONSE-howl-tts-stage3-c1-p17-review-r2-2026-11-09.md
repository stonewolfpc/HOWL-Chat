# C1 — Targeted Author Response to Second-Stage Review (RC-01–RC-06)

**Date:** 2026-11-09
**Status:** AUTHOR RESPONSE. **P17 NOT APPROVED. RECONCILIATION NOT CONFIRMED AS COMPLETE.** Not rev4, not P18.
**Responds to:** `REVIEW-howl-tts-stage3-c1-p17-reconciliation-2026-11-09.md`, preserved verbatim.
**Subject:** reconciliation `b634cd20…`; P17 `01d887c1…`, unchanged.
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

No source, contract, test, or build change. Each correction below states what was claimed, what the source shows, and what is withdrawn or replaced.

## 1. Corrections RC-01–RC-06

### RC-01 — Restore `client.go:439/441` to `callOn`; withdraw the false reassignment

- **Second-stage review:** P17's 439/441 were correct `callOn` branches; the reconciliation wrongly reassigned them to `Synthesize`.
- **Verified:** `callOn` select is `client.go:429–442`; `timer.C` returns `timeout` at `:438–439`; `ctx.Done()` returns `cancelled` at `:440–441`.
- **Action:** the reconciliation's correction is **WITHDRAWN**. P17's O13/O14 citations were right. `Synthesize` has separate branches at `:766–772` and `:773–784`. R01 remains confirmed; only the mapping is restored.

### RC-02 — Correct the R03 transport-send citation

- **Second-stage review:** `:410–412` is encode failure; the send failure is `:420–421`.
- **Verified:** `encodeRequestFrame` error is handled at `client.go:409–414` and returns `invalid_request` at `:414`. `transport.Send` error is `client.go:420–421` and returns `helper_crashed` at `:421`.
- **Action:** R03's trace is corrected to `:420–421`. R03 itself stands: the send failure returns `helper_crashed` without terminating. The reconciliation's `:410–412` citation is withdrawn as a misattribution.

### RC-03 — Unique, branch-level site manifest

- **Review finding:** 45 was arithmetic over overlapping categories, not a reproducible count.
- **Action:** the 45 count is **WITHDRAWN**. The manifest in §3 lists one row per unique source line performing an evidence operation, with exact line, function, branch condition, type, phase, and source ID. A line named by multiple P17/E/T rows appears exactly once; its `IDs` cell lists every association. The 45 figure is retained only as a withdrawn historical claim.

### RC-04 — Recomputed count

- **Action:** the count is recomputed from the manifest, not asserted: **67 unique lines** (§3). The machine check in §4 reproduces both the count and the absence of duplicate lines.
- **Scope limit:** 67 is the size of the inventoried union, not proof that every lifecycle site exists. R08's completeness disposition is unchanged.

### RC-05 — Digest procedure and test-file inclusion

- **Verified:** the documented command below reproduces `d0df8b28…` exactly.
- **Scope:** recursive `internal/audio/tts/helper`, all 19 `*.go` files, including the 11 `*_test.go` files. It excludes non-Go files and the rest of the repository.
- **Action:** procedure documented in §5, including ordering, path form, and environment sensitivity.

### RC-06 — Correct the `callOn` range and branch references

- **Verified:** `callOn` spans `client.go:398–443`, not `:398–431`.
- **Branches:** begin/deliver `:402–407`; encode `client.go:409–414`; missing transport `:416–419`; send `client.go:420–421`; timer default `:424–426`; select `:429–442`; timeout `:438–439`; cancellation `:440–441`.
- **Action:** all citations corrected to this map. `waitForReader` remains `:349–362`; `protocolViolation` `:364–367`; `failAll` `:369–386`.

## 2. Preserved dispositions

| ID | Disposition after this response |
|---|---|
| R01 | confirmed; mapping restored, false correction withdrawn |
| R02 | confirmed |
| R03 | confirmed; citation corrected to `:420–421` |
| R04 | confirmed |
| R05 | confirmed |
| R06 | confirmed |
| R07 | confirmed conceptually; mapping replaced by §3 |
| R08 | open; completeness not established; 45 withdrawn, 67 is inventory size only |
| R09–R13 | unresolved model and compatibility questions |
| E1–E9 | individually retained; E7–E9 confirmed |
| U1–U8 | open |

No R, E, or U identifier is removed, merged, or closed by this response.

## 3. Unique site manifest

Rule: one row per unique source line that creates a failure, records a lifecycle fact, publishes or overwrites state, drops evidence, or issues a termination action. Forwarding `SetState` consumers are not separate rows; their creating line is listed once with the publication noted. `IDs` may name several rows; the site appears once.

| SID | Site | Function | Branch condition | Type | Phase | IDs |
|---|---|---|---|---|---|---|
| S01 | discovery.go:39 | Discover | helper path not absolute | failure | discovery | O1 |
| S02 | discovery.go:42 | Discover | helper path not a regular file | failure | discovery | O1 |
| S03 | discovery.go:61 | Discover | optional helper absent | failure | discovery | O1 |
| S04 | protocol.go:504 | validateHello | version empty or over 64 | failure | handshake | O16, R12 |
| S05 | protocol.go:507 | validateHello | limits mismatch | failure | handshake | O16, R12 |
| S06 | protocol.go:514 | validateHello | provider descriptor invalid | failure | handshake | O16, R12 |
| S07 | protocol.go:517 | validateHello | duplicate provider ID | failure | handshake | O16, R12 |
| S08 | protocol.go:521 | validateHello | provider capabilities invalid | delegated failure | handshake | O16, R12 |
| S09 | protocol.go:525 | validateHello | service capabilities invalid | delegated failure | handshake | O16, R12 |
| S10 | protocol.go:528 | validateHello | mandatory capabilities missing | failure plus detail | handshake | O16, R12 |
| S11 | protocol.go:535 | validateHello | extension descriptor invalid | failure | handshake | O16, R12 |
| S12 | protocol.go:538 | validateHello | duplicate extension | failure | handshake | O16, R12 |
| S13 | supervisor.go:147 | Start | factory error | failure | start | O2 |
| S14 | supervisor.go:178 | observeExit | stopping and zero exit | fact | exit | O4 |
| S15 | supervisor.go:180 | observeExit | stopping and nonzero exit | fact | exit | O5, E3 |
| S16 | supervisor.go:182 | observeExit | stopping, either exit | state absent | exit | O4, O5, E3 |
| S17 | supervisor.go:185 | observeExit | unexpected exit | state crashed | exit | O3 |
| S18 | supervisor.go:186 | observeExit | unexpected exit | failure | exit | O3 |
| S19 | supervisor.go:238 | SetState | any publication | overwrite lastFailure | publication | E2 |
| S20 | supervisor.go:246 | Terminate | entry | action | termination | E4, T1–T9 |
| S21 | supervisor.go:275 | Terminate | ctx done | action kill | termination | O20, T10 |
| S22 | supervisor.go:276 | Terminate | ctx done | result cause | termination | R06 |
| S23 | supervisor.go:278 | Terminate | grace expiry | action kill | termination | O20, T10 |
| S24 | supervisor.go:285 | Terminate | second grace expiry | failure timeout | termination | O20, R06 |
| S25 | supervisor.go:291 | Kill | entry | action alias | termination | E4 |
| S26 | supervisor.go:292 | Kill | default grace | action | termination | E4, T2–T4 |
| S27 | client.go:248 | Start | handshake failure | action terminate | handshake | O15, T1 |
| S28 | client.go:251 | Start | incompatible or minor unsupported | state incompatible | handshake | O15 |
| S29 | client.go:253 | Start | other handshake failure | state unavailable | handshake | O15, O16 |
| S30 | client.go:301 | readLoop | internal_error | reclassify as helper_crashed | transport | O6, E6 |
| S31 | client.go:305 | readLoop | invalid frame, too large, violation | action kill | protocol | O7, T2 |
| S32 | client.go:311 | readLoop | decode response error | dispatch failure | protocol | O8 |
| S33 | client.go:312 | readLoop | decode response error | action kill | protocol | O8, T3 |
| S34 | client.go:323 | dispatch | empty response ID | violation | protocol | O9 |
| S35 | client.go:331 | dispatch | duplicate terminal | violation | protocol | O10 |
| S36 | client.go:341 | dispatch | retired request ID | silent drop | protocol | E7 |
| S37 | client.go:345 | dispatch | unknown request ID | violation | protocol | O11 |
| S38 | client.go:359 | waitForReader | reader done | fact, no record | exit | E5, R04 |
| S39 | client.go:360 | waitForReader | timer expiry | timeout, no record | exit | E5, R04 |
| S40 | client.go:366 | protocolViolation | any violation | action kill | protocol | T4 |
| S41 | client.go:374 | failAll | readFailure nil | first-only gate | evidence | E1 |
| S42 | client.go:375 | failAll | first failure | store once | evidence | E1 |
| S43 | client.go:421 | callOn | send error | failure, no termination | transport | O12, R03 |
| S44 | client.go:439 | callOn | timer expiry | failure timeout | request | O13 |
| S45 | client.go:441 | callOn | ctx done | failure cancelled | request | O14 |
| S46 | client.go:581 | Load | call failure | state unavailable | load | O17 |
| S47 | client.go:586 | Load | decode failure | state unavailable | load | O17 |
| S48 | client.go:590 | Load | validation failure | state unavailable | load | O17 |
| S49 | client.go:654 | Shutdown | failure not crashed or timeout | action, ignored | shutdown | E8, T5 |
| S50 | client.go:661 | Shutdown | decode error | action, ignored | shutdown | E8, T6 |
| S51 | client.go:664 | Shutdown | not accepted | action, ignored | shutdown | E8, T7 |
| S52 | client.go:668 | Shutdown | acknowledgement accepted | action, result returned | shutdown | T8 |
| S53 | client.go:670 | Shutdown | after terminate | state clear | shutdown | E9 |
| S54 | client.go:672 | Shutdown | after terminate | transport clear | shutdown | E9 |
| S55 | client.go:756 | Synthesize | send error | state crashed | transport | O18, R13 |
| S56 | client.go:771 | Synthesize | timer expiry | state timeout | request | O19 |
| S57 | client.go:783 | Synthesize | cancel grace expiry | state cancelled | request | O19 |
| S58 | client.go:793 | Synthesize | pending failure and cancel requested | state ready | request | O19 |
| S59 | client.go:796 | Synthesize | pending failure otherwise | state degraded | request | O19, R02 |
| S60 | client.go:803 | Synthesize | wire error and cancel requested | state ready | request | O19 |
| S61 | client.go:806 | Synthesize | wire error otherwise | state degraded | request | O19 |
| S62 | client.go:814 | Synthesize | result decode error | state degraded | output | O19 |
| S63 | client.go:820 | Synthesize | result validation error | state degraded | output | O19 |
| S64 | client.go:828 | Synthesize | WAV validation error | state degraded | output | O19 |
| S65 | client.go:839 | Synthesize | cancelled or timed out pre-acceptance | state ready | request | O19 |
| S66 | client.go:861 | cleanupRejected | process alive and output may exist | action, ignored | cleanup | E8, T9 |
| S67 | client.go:871 | watchCancellation | ctx done | action, ignored | cleanup | E8 |

**Deduplications made explicit:** `client.go:756` was O12 and O18; `:771` was O13 and O19; `:783` was O14 and O19; `:248` was O15 and T1; `:305` was O7 and T2; `:312` was O8 and T3; `supervisor.go:180` was O5 and E3; `:301` was O6 and E6; `:654` was E8 and T5; `:275/:278` were O20 and T10. Each is one row above. The prior 45 counted some of these more than once and omitted branch lines.

## 4. Regenerating the count

```bash
D=docs/review/RESPONSE-howl-tts-stage3-c1-p17-review-r2-2026-11-09.md
grep -c '^| S[0-9]' "$D"                       # expect 67
awk -F'|' '/^\| S[0-9]/ {gsub(/^ +| +$/,"",$3); print $3}' "$D" \
  | sort | uniq -d                              # expect no output
awk -F'|' '/^\| S[0-9]/ {gsub(/ /,"",$2); print $2}' "$D" \
  | sort -c                                     # expect no output (sequential IDs)
```

## 5. Digest procedure (RC-05)

Run from the repository root; the relative path form is part of the hashed stream.

```bash
find internal/audio/tts/helper -type f -name '*.go' -print0 \
  | LC_ALL=C sort -z \
  | xargs -0 sha256sum \
  | sha256sum
```

Expected: `d0df8b2823389b305698fc06eda28e0e01f06259fdc69174cf7fb3851e271a51`.

- **Path scope:** recursive under `internal/audio/tts/helper`, regular files matching `*.go`, relative to the repository root.
- **Inclusion:** 19 files, including 11 `*_test.go` files. Non-Go files are excluded.
- **Ordering:** NUL-delimited byte order (`LC_ALL=C sort -z`), independent of locale.
- **Hashed stream:** the GNU `sha256sum` output, one line per file as `<hex>  <path>`; the digest therefore binds relative paths as well as contents.
- **Sensitivity:** running from another directory, using absolute paths, or renaming files changes the digest. A different path separator or coreutils implementation may also differ.
- **Reproduction:** rerunning the block above on the audited tree returns the expected value; a mismatch means path scope, ordering, or inclusion changed.

This documents `d0df8b28…`; it does not make it a Git identity. The helper remains untracked.

## 6. What this response does not do

- It does not approve P17 or confirm the reconciliation.
- It does not revise P17, the reconciliation, or any frozen artifact.
- It does not close R01–R13, E1–E9, or U1–U8.
- It does not claim the 67-line manifest is complete beyond the inventoried union.
- It does not authorize rev4, P18, implementation, tests, or builds.

## 7. Status

Corrections recorded. Count regenerated from the manifest. Digest procedure reproducible. **P17 NOT APPROVED. Reconciliation NOT CONFIRMED. Contract amendment and C1 implementation UNAUTHORIZED.**
