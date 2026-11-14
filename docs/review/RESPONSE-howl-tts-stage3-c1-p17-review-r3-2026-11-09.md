# C1 — Targeted Manifest Correction (CR-01–CR-08)

**Date:** 2026-11-09
**Status:** AUTHOR CORRECTION. **P17 NOT APPROVED. RC-03/RC-04 REMAIN PARTIAL.** Not rev4, not P18.
**Responds to:** `REVIEW-howl-tts-stage3-c1-p17-manifest-2026-11-09.md`, preserved verbatim.
**Supersedes for manifest purposes:** r2 response `57cc58f5…`; that artifact remains frozen.
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

No source, contract, test, or build change. The r2 manifest is not edited in place; this response replaces its counting method and records why.

## 0. Fourth-stage corrections (FR-01–FR-04)

The fourth review is preserved at `REVIEW-howl-tts-stage3-c1-p17-r3-2026-11-09.md`. Its four scoped corrections are applied here.

| FR | Correction |
|---|---|
| FR-01 | O14 attached to S59, S61, and S66, so Synthesize cancellation is associated where it is handled |
| FR-02 | M08 corrected from `client.go:294` to `client.go:295`; the mechanism inventory is unchanged |
| FR-03 | 61 relabeled a provisional condition-pair count, not a verified branch count |
| FR-04 | §6 now reads only ID cells; prose mentions cannot satisfy coverage |

These change associations and labels, not source facts. Emission and mechanism counts are unchanged.

## 1. CR mapping and corrections

The review lists five numbered corrections and one additional distinction, then requires CR-01–CR-08. The mapping used here is stated so no identifier is invented silently.

| CR | Basis | Correction |
|---|---|---|
| CR-01 | correction 1 | propagate O12/O13/O14 to `Synthesize` sites `:756`, `:771`, `:783` |
| CR-02 | correction 2 | represent all five ignored `Terminate` calls and the ignored `Cancel` |
| CR-03 | correction 2 | add the `started = false` operation at `client.go:671` |
| CR-04 | correction 3 | move E7 from `:341` to the silent return at `:342` |
| CR-05 | correction 4 | add the missing O2, O3, O6 operations at `:149`, `:190`, `:303` |
| CR-06 | correction 5 | narrow the inclusion rule explicitly; do not adjust it silently |
| CR-07 | additional distinction | separate emission anchors from mechanism anchors; add a provisional condition-pair count |
| CR-08 | additional distinction | add an independent coverage verification procedure and scope statement |

### CR-01 — O12/O13/O14 associations

- `client.go:756` now carries **O12** (send failure), **O18** (provider crash), **R13**.
- `client.go:771` now carries **O13** (synthesis timeout), **O19**.
- `client.go:783` now carries **O14** (synthesis cancellation), **O19**.
- `client.go:421` remains the `callOn` send site under **O12**.
- The r2 table's O18/O19-only appearance is withdrawn as an association error.

### CR-02 — E8 completeness

All six ignored results are now represented:

| Line | Call | Ignored result |
|---|---|---|
| `client.go:248` | `Terminate` | handshake failure cleanup |
| `client.go:654` | `Terminate` | shutdown failure branch |
| `client.go:661` | `Terminate` | shutdown decode error |
| `client.go:664` | `Terminate` | shutdown not accepted |
| `client.go:861` | `Terminate` | cleanupRejected |
| `client.go:871` | `Cancel` | cancellation watcher |

`:668` is not counted here because its result is returned.

### CR-03 — E9 state clearing

`client.go:668–674` performs four distinct operations: terminate result (`:668`), `closing = true` (`:670`), `started = false` (`:671`), `transport = nil` (`:672`). The previously missing `:671` is now a separate emission anchor. All three clear operations are tagged **E9**.

### CR-04 — E7 line

The silent return is `client.go:342`; `:341` is the mutex unlock. The r2 row is corrected and the `:341` site is withdrawn.

### CR-05 — O2/O3/O6 operations

| Path | Added site | Operation |
|---|---|---|
| O2 | `supervisor.go:149` | `setUnavailable` publication after factory failure |
| O3 | `supervisor.go:190` | `lastFailure` assignment for unexpected exit |
| O6 | `client.go:303` | `failAll` dispatch after stream reclassification |

O17's adverse publications (`:581`, `:586`, `:590`) were already present. Its routine `loading` and success `ready` writes are excluded by the narrowed rule in §2, not omitted silently.

### CR-06 — Inclusion rule

**Narrowed rule, applied to the 68 rows below:**

Include a source line only when it performs one of:

1. failure creation or reclassification;
2. a terminal fact or adverse state publication (`absent`, `exited`, `crashed`, `unavailable`, `incompatible`, `degraded`, `stopping`, `timeout`, `cancelled`);
3. an evidence store, overwrite, dispatch, or drop;
4. a lifecycle action issue or an unrecorded action result.

Exclude declarations, routine progress or success states (`starting`, `loading`, `busy`, `ready`, `unloaded` when no failure is attached), pure validation returns, and forwarding that does not change state.

This rule is narrower than r2's wording. It is not being changed to preserve a count: applying it removes three mechanism declarations and adds four operations, for 68.

### CR-07 — Three counts

- **Emission anchors:** 68 (§3).
- **Mechanism anchors:** 17 (§4).
- **Condition-pair count (provisional):** 61 (§5), distinct `(function, condition)` pairs. This is not a verified branch count; stable branch identifiers or source locations are required for that.

`S20` (`:246`) and `S25` (`:291`) are no longer emission rows. `:292`, the `Kill` alias return, is also a mechanism, not an independent action.

### CR-08 — Coverage verification

§6 provides commands that read only the ID cells of the emission and mechanism-finding tables and check every O1–O20, E1–E9, and T1–T10. Prose mentions do not satisfy the check. The declared scope is exactly that union. Coverage beyond it is not claimed.

## 2. Declared scope

The manifest covers the union of P17 O1–O20, E1–E9, and T1–T10 as interpreted by the §1 rule. It does not enumerate arbitrary code paths, validation failures outside that union, or non-Go files. A passing coverage check proves association, not completeness beyond the declared scope.

## 3. Emission anchors (68)

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
| S14 | supervisor.go:149 | Start | factory error | state unavailable | start | O2 |
| S15 | supervisor.go:178 | observeExit | stopping and zero exit | fact absent | exit | O4 |
| S16 | supervisor.go:180 | observeExit | stopping and nonzero exit | fact exited | exit | O5, E3 |
| S17 | supervisor.go:182 | observeExit | stopping, either exit | state absent | exit | O4, O5, E3 |
| S18 | supervisor.go:185 | observeExit | unexpected exit | state crashed | exit | O3 |
| S19 | supervisor.go:186 | observeExit | unexpected exit | failure | exit | O3 |
| S20 | supervisor.go:190 | observeExit | unexpected exit | store lastFailure | exit | O3 |
| S21 | supervisor.go:238 | SetState | any publication | overwrite lastFailure | publication | E2 |
| S22 | supervisor.go:275 | Terminate | ctx done | action kill | termination | O20, T10 |
| S23 | supervisor.go:276 | Terminate | ctx done | result cause | termination | R06 |
| S24 | supervisor.go:278 | Terminate | grace expiry | action kill | termination | O20, T10 |
| S25 | supervisor.go:285 | Terminate | second grace expiry | failure timeout | termination | O20, R06 |
| S26 | client.go:248 | Start | handshake failure | action terminate, ignored | handshake | O15, T1, E8 |
| S27 | client.go:251 | Start | incompatible or minor unsupported | state incompatible | handshake | O15 |
| S28 | client.go:253 | Start | other handshake failure | state unavailable | handshake | O15, O16 |
| S29 | client.go:301 | readLoop | internal_error | reclassify helper_crashed | transport | O6, E6 |
| S30 | client.go:303 | readLoop | receive error | dispatch failAll | transport | O6 |
| S31 | client.go:305 | readLoop | invalid frame, too large, violation | action kill | protocol | O7, T2 |
| S32 | client.go:311 | readLoop | decode response error | dispatch failure | protocol | O8 |
| S33 | client.go:312 | readLoop | decode response error | action kill | protocol | O8, T3 |
| S34 | client.go:323 | dispatch | empty response ID | violation | protocol | O9 |
| S35 | client.go:331 | dispatch | duplicate terminal | violation | protocol | O10 |
| S36 | client.go:342 | dispatch | retired request ID | silent drop | protocol | E7 |
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
| S52 | client.go:668 | Shutdown | acknowledgement accepted | action, result returned | shutdown | T8, E9 |
| S53 | client.go:670 | Shutdown | after terminate | state closing | shutdown | E9 |
| S54 | client.go:671 | Shutdown | after terminate | state started false | shutdown | E9 |
| S55 | client.go:672 | Shutdown | after terminate | state transport nil | shutdown | E9 |
| S56 | client.go:756 | Synthesize | send error | state crashed | transport | O12, O18, R13 |
| S57 | client.go:771 | Synthesize | timer expiry | state timeout | request | O13, O19 |
| S58 | client.go:783 | Synthesize | cancel grace expiry | state cancelled | request | O14, O19 |
| S59 | client.go:793 | Synthesize | pending failure and cancel requested | state ready | request | O14, O19 |
| S60 | client.go:796 | Synthesize | pending failure otherwise | state degraded | request | O19, R02 |
| S61 | client.go:803 | Synthesize | wire error and cancel requested | state ready | request | O14, O19 |
| S62 | client.go:806 | Synthesize | wire error otherwise | state degraded | request | O19 |
| S63 | client.go:814 | Synthesize | result decode error | state degraded | output | O19 |
| S64 | client.go:820 | Synthesize | result validation error | state degraded | output | O19 |
| S65 | client.go:828 | Synthesize | WAV validation error | state degraded | output | O19 |
| S66 | client.go:839 | Synthesize | cancelled or timed out pre-acceptance | state ready | request | O14, O19 |
| S67 | client.go:861 | cleanupRejected | process alive and output may exist | action, ignored | cleanup | E8, T9 |
| S68 | client.go:871 | watchCancellation | ctx done | action, ignored | cleanup | E8 |

## 4. Mechanism anchors (17)

These explain how operations work and are not counted as emissions.

| MID | Site | Mechanism |
|---|---|---|
| M01 | discovery.go:36 | Discover |
| M02 | supervisor.go:197 | setUnavailable |
| M03 | supervisor.go:233 | SetState |
| M04 | supervisor.go:246 | Terminate |
| M05 | supervisor.go:291 | Kill |
| M06 | supervisor.go:292 | Kill alias to Terminate |
| M07 | client.go:191 | Start |
| M08 | client.go:295 | readLoop |
| M09 | client.go:321 | dispatch |
| M10 | client.go:349 | waitForReader |
| M11 | client.go:364 | protocolViolation |
| M12 | client.go:369 | failAll |
| M13 | client.go:398 | callOn |
| M14 | client.go:645 | Shutdown |
| M15 | client.go:859 | cleanupRejected |
| M16 | client.go:866 | watchCancellation |
| M17 | protocol.go:502 | validateHello |

Mechanism-level findings:

| Mechanism finding | Description | Applies to |
|---|---|---|
| E4 | no intent parameter; Kill aliases Terminate | M04, M05, M06 |

## 5. Three counts and regeneration

| Count | Value | Definition |
|---|---|---|
| Emission anchors | 68 | unique source lines performing an evidence operation |
| Mechanism anchors | 17 | declarations or shared implementations |
| Condition-pair count (provisional) | 61 | unique function and condition pairs; not a verified branch count |

```bash
D=docs/review/RESPONSE-howl-tts-stage3-c1-p17-review-r3-2026-11-09.md
grep -c '^| S[0-9]' "$D"                       # 68
grep -c '^| M[0-9]' "$D"                       # 17
awk -F'|' '/^\| S[0-9]/ {for(i=1;i<=NF;i++) gsub(/^ +| +$/,"",$i); print $4"|"$5}' "$D" \
  | sort -u | wc -l                            # condition pairs (provisional)
awk -F'|' '/^\| S[0-9]/ {gsub(/^ +| +$/,"",$3); print $3}' "$D" | sort | uniq -d   # no duplicates
```

The emission count is derived from the table, not asserted. B is derived from distinct function-and-condition pairs, so it cannot rise merely because one condition spans more lines.

## 6. Independent coverage verification

The check reads only ID cells, never prose:

```bash
D=docs/review/RESPONSE-howl-tts-stage3-c1-p17-review-r3-2026-11-09.md
awk -F'|' '/^\| S[0-9]/ {gsub(/^ +| +$/,"",$8); print $8}' "$D" \
  | tr ',' '\n' | sed 's/^ *//; s/ *$//' | grep -v '^$' | sort -u > /tmp/p17-emission-ids
awk -F'|' '/^\| E[0-9] \|/ {gsub(/^ +| +$/,"",$2); print $2}' "$D" \
  | sort -u > /tmp/p17-mechanism-ids
for id in $(seq 1 20); do grep -qx "O$id" /tmp/p17-emission-ids || grep -qx "O$id" /tmp/p17-mechanism-ids || echo "MISSING O$id"; done
for id in $(seq 1 9);  do grep -qx "E$id" /tmp/p17-emission-ids || grep -qx "E$id" /tmp/p17-mechanism-ids || echo "MISSING E$id"; done
for id in $(seq 1 10); do grep -qx "T$id" /tmp/p17-emission-ids || grep -qx "T$id" /tmp/p17-mechanism-ids || echo "MISSING T$id"; done
```

Expected: no output. A prose mention cannot satisfy this check. It verifies cell-level association within the declared scope, not that the scope is exhaustive.

## 7. Preserved dispositions

R01–R13, E1–E9, and U1–U8 keep their prior identifiers and dispositions. R01/R03/R05/R06 remain resolved at citation level; RC-03/RC-04 remain partial until an independent reviewer confirms the manifest; U1–U8 remain open. No finding is closed by this response.

## 8. Status

Manifest corrected under an explicit narrowed rule. Emission anchors and mechanism anchors are counted; the condition-pair count is provisional. **P17 NOT APPROVED. RC-03/RC-04 PARTIAL. Contract amendment and C1 implementation UNAUTHORIZED.** The third count is provisional and is not a verified branch count.
