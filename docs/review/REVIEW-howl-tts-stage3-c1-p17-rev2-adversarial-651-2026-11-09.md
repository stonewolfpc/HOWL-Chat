# Independent Adversarial Review — P17 Rev2 B-2-Corrected Candidate (651)

**Review status:** COMPLETE — **NO BLOCKING FINDING; B-2 PASSES; R-5/R-6 PASS; N-5/N-6 REMAIN OPEN; NOT AN APPROVAL.**

**Reviewed artifact:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`

**Reviewed identity:** **538 lines**, SHA-256 `651d38d3e2d97f981a606a39ea4a3117c2e748ff9f5941eb8bd9671b6e022f82`

**Prior blocking review:** `docs/review/REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-538-2026-11-09.md`, 200 lines, SHA-256 `3677834619d342441f58bf0e6780aaf3f4eba65edd1ddfa6c4de242538f7fc8e`. That review applied to the superseded `0cf14a2…` candidate and is not an approval of this artifact.

**Frozen inputs checked:** original P17, 166 lines, `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf`; decision register, 434 lines, `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`; contract rev3, 183 lines, `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e`.

**Basis:** Direct read of the exact 538-line artifact, including both ordinary C-B outcome branches, the complete C-A/C-C/C-D/P4-D route set, the B-1 sentence, the frozen controls, and the source-backed thirteen-value provider-state vocabulary. No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or candidate edits were performed. This review does not approve P17 or authorize C1.

---

## 1. Verdict

**No blocking finding.**

The B-2 correction is effective. C-B’s final gate at target line 370 now contains both required exclusions:

1. no separate prior public value remains applicable, unresolved, or **CP-1 OPEN**; and
2. no separate prior public value is present that is neither a prior public failure value nor an active operational state.

The same two exclusions are repeated in the no-orderly `unavailable` branch. A separate non-failure, non-active public value therefore cannot reach C-B’s ordinary `absent`/`unavailable` outcome in either branch.

The `unloaded` counterexample is covered by the second, class-wide clause. It is not hard-coded. The result remains **OPEN** under line 368 and §13.

R-5 and R-6 remain closed at their named routes. R-1 through R-4, C-A, C-C, CP-1, P4-D, B-1, and the prior N-1–N-4/N-7 corrections remain present. N-5 and N-6 remain **OPEN** as scope limitations, not as defects in the corrected gate.

The artifact is **not approved**, and P17 remains unapproved.

---

## 2. B-2 correction audit

### 2.1 The final gate now carries both exclusions

Target line 370 now says, in relevant part:

> “Only when no separate prior public value remains applicable, unresolved, or **CP-1 OPEN**, and no separate prior public value is present that is neither a prior public failure value nor an active operational state, does C-B proceed to the ordinary `absent`/`unavailable` outcome…”

This is the necessary correction. The first clause protects applicable, unresolved, and CP-1 failure values. The second clause protects the entire non-failure, non-active class.

The no-orderly branch repeats the same two conditions:

> “Without orderly shutdown, `unavailable` applies only if no separate prior public value remains applicable, unresolved, or **CP-1 OPEN**, and no separate prior public value is present that is neither a prior public failure value nor an active operational state; otherwise the separate prior public value governs under the preceding paragraph, including CP-1 or an **OPEN** branch.”

The repetition matters: the prior defect affected both ordinary outcomes, so a correction limited to the orderly `absent` branch would have left the no-orderly `unavailable` bypass in place.

### 2.2 The `unloaded` counterexample is closed

The concrete scenario remains:

1. C-B is established: the prior failure classification is established and its publication status is positively resolved as `excluded-internal-only`.
2. No fatal observation and no terminal exist.
3. A separate prior public value `V = unloaded` is established by publication evidence.
4. `unloaded` is a member of the frozen thirteen-value provider-state set and is neither a prior public failure value nor an active operational state.
5. Orderly shutdown is independently established and the other C-B gates are satisfied.

Under the revised line 370, condition 2 is false because `V` is present and is neither a failure value nor an active state. C-B does not proceed. Line 368 independently says the combination remains **OPEN**. The two rules now agree.

The same result follows without orderly shutdown: the second condition is false, so `unavailable` does not apply and the separate value governs under the preceding paragraph.

### 2.3 The correction is class-wide, not state-specific

The revised clause covers every prior public value that is neither a prior public failure value nor an active operational state. It does not name `unloaded`, does not enumerate the thirteen states, and does not create a new preservation or projection rule.

The clause is also conservative for a value whose classification is disputed: if the value might be non-failure and non-active, the condition blocks the ordinary C-B outcome rather than assuming a failure classification.

### 2.4 No new fallback rule was introduced

The correction only withholds C-B’s ordinary outcome. It does not select `absent`, `unavailable`, a preserved prior state, or a new state for the uncovered class. The specific outcome remains **OPEN**, consistent with §13.

### 2.5 Both final-gate branches are now symmetrical

| Branch | Required first exclusion | Required second exclusion | Result |
|---|---|---|---|
| Orderly shutdown → `absent` | Present at line 370 | Present at line 370 | PASS |
| No orderly shutdown → `unavailable` | Present at line 370 | Present at line 370 | PASS |
| Internal-only establishment/publication/public value unresolved | `absent` ineligible; orderly outcome **OPEN** | No bypass found | PASS |
| Separate non-failure, non-active value with no orderly shutdown | Second exclusion blocks `unavailable` | Separate value governs as **OPEN** | PASS |

No remaining route was found by which the uncovered class reaches the ordinary C-B outcome.

---

## 3. C-B route matrix

| Separate prior public value | Current route | Result |
|---|---|---|
| None | C-B ordinary outcome | PASS |
| Prior public failure, compatible and applicable | C-A preserves; C-B final gate blocked by “applicable” | PASS |
| Prior public failure, compatibility/applicability unresolved | C-A/Gate 5 **OPEN**; final gate blocked by “unresolved” | PASS |
| Prior public failure, established incompatible | C-A historical evidence; ordinary P4-D | PASS at the outcome level; see the separate C-B/P4-D wording observation |
| Prior public `crashed`, otherwise compatible/applicable | **CP-1 OPEN**; final gate blocked by “CP-1 OPEN” | PASS |
| Active operational state, including active `degraded` | C-C by line 368; C-C’s active-state and unresolved-interpretation rules remain applicable | PASS |
| Non-failure, non-active value, including `unloaded` | Line 368 **OPEN**; final-gate second exclusion blocks both ordinary outcomes | PASS |
| Terminal or fatal observation also present | Case A or Case B controls; C-B does not apply | PASS |
| Multiple separate values/classifications | N-5 remains **OPEN**; no selection rule is invented | OPEN, correctly |

The active row is not a C-B ordinary-outcome bypass: line 368 expressly routes active operational states through C-C, and C-C continues to exclude active preservation and to keep an unresolved post-exit `degraded` interpretation **OPEN**. The final gate’s second exclusion is limited to the non-failure, non-active class; it does not convert an active value into a C-B value.

---

## 4. C-A, C-C, CP-1, P4-D, and B-1

### 4.1 C-A

Target lines 362–364 retain the required distinctions:

- compatible and applicable prior public failure → preserve;
- compatibility or applicability unresolved → **OPEN**;
- `crashed`, otherwise compatible and applicable → **CP-1 OPEN**;
- established incompatible → historical evidence, then ordinary P4-D handling.

Compatibility does not establish applicability, and an active `degraded` state remains outside C-A. PASS.

### 4.2 C-C

Target line 372 retains the active-state exclusion and the separate post-exit interpretation rule for `degraded`. Active operational states are not preserved, and an asserted but unresolved post-exit interpretation remains **OPEN**. The active exclusion is described as a gate-5 eligibility filter, not as a new gate. PASS.

### 4.3 CP-1

CP-1 remains **OPEN** at line 385 and in §13. The revised final gate expressly names **CP-1 OPEN** in its first exclusion, so a separate compatible/applicable public `crashed` value cannot reach the ordinary C-B outcome. No new public `crashed` route is selected. PASS.

### 4.4 P4-D

P4-D remains the ordinary confirmed-exit route for no applicable prior controlling state. The incompatible-failure route still reduces that value to historical evidence before P4-D. The existing C-B/P4-D wording observation remains separate and non-blocking; the B-2 correction did not introduce a new P4-D contradiction.

### 4.5 B-1

The B-1 sentence remains at line 368:

> “If a separate prior public provider-state value is present and is neither a prior public failure value nor an active operational state, C-B does not proceed to the ordinary `absent`/`unavailable` outcome; the combination remains OPEN.”

It is now integrated with, rather than contradicted by, line 370. PASS.

---

## 5. R-1 through R-6

| Prior correction | Current evidence | Result |
|---|---|---|
| R-1: three-way establishment; unresolved blocks `absent` | P4-A lines 296–300; B-B line 344; Case C introduction line 360; C-D line 376 | PASS |
| R-2: published / excluded-internal-only / unresolved remain distinct | P4-A line 299; C-B line 366 | PASS |
| R-3: CP-1 remains OPEN and is not indirectly resolved | S7 line 190; S8 line 191; O3–O5 lines 207–209; P4-A lines 302, 312, 316–319; B-C line 354; C-A line 362; C-B lines 368–370; Case A line 327; public `crashed` line 385 | PASS |
| R-4: active operational values cannot survive confirmed exit | P4-A line 302; S8 line 191; O4 line 208; B-C line 354; C-A line 362; C-C line 372; Case A line 326 | PASS |
| R-5: terminal precedence over no-orderly fallbacks | S7 line 190; O3/O5 lines 207/209; O4/S8 lines 208/191; Gate 2 line 315; P4-A precedence line 312; Gates 3/4 lines 316–317; no-orderly line 319; B-D line 356; Case A lines 325–330 | PASS |
| R-6: separate prior public failure/CP-1 routes | C-A line 362; C-C line 372; CP-1 line 368/385; P4-D line 380; C-B final gate lines 368–370 | PASS for the named scope |

The revised C-B gate does not reopen R-5: no ordinary outcome is selected over an established terminal. It does not reopen R-6: applicable, unresolved, incompatible, active, and CP-1 failure values remain outside the ordinary C-B fallback.

---

## 6. N-1 through N-7 and §13

| Finding | Current evidence | Status |
|---|---|---|
| N-1 — C-D/B-C cross-reference | C-D line 376 routes through C-A and Gate 5 | CORRECTED |
| N-2 — compatibility vs applicability | P4-A line 304; C-A line 362; B-C line 354 | CORRECTED |
| N-3 — Gate 5 filter wording | C-C line 372 says active exclusion is an eligibility filter, not a separate gate | CORRECTED |
| N-4 — publication vs justification | P4-A line 304 distinguishes publication status from independent justification | CORRECTED |
| N-5 — multiple classifications | §13 line 530 retains selection as **OPEN**; no selection rule is invented | OPEN, correctly |
| N-6 — non-failure public value | Line 368 and revised line 370 both keep the combination **OPEN**; §13 line 530 retains the scope limitation | OPEN as a scope limitation; consistency defect corrected |
| N-7 — local fallback wording | P4-A lines 316–317 repeat the terminal exception | CORRECTED |

### 6.1 Separate B-B observation

B-B line 350 still says a separate prior public value is governed by “B-C/Gate 5” or C-C, while B-C is expressly limited to a previously published public **failure** value. A non-failure, non-active value is not expressly routed by B-C.

This remains a non-blocking scope observation. It is outside the B-2/C-B correction and does not turn the C-B fix into a failure. If the owner intends the broad B-B phrase to cover every public value, it needs a separate scope decision.

### 6.2 Separate C-B/P4-D observation

C-B’s incompatible-value route at line 368 retains the internal-only classification while routing ordinary confirmed-exit handling through P4-D. P4-D line 380 says no applicable prior terminal, fatal observation, or failure classification remains. The classification is retained internally and is not a public applicable failure, so the intended outcome is consistent, but the wording would be clearer with a public/applicable qualifier.

This remains non-blocking and separate from B-2.

### 6.3 Scope status

§13 line 530 still states that two scope limitations remain explicitly **OPEN**:

1. selection among multiple prior failure classifications; and
2. treatment of a separate non-failure, non-active prior public provider-state value in C-B.

The corrected line 370 makes the second limitation’s current result explicit without deciding a preservation or replacement rule. N-5 and N-6 therefore remain OPEN by design, not as hidden defaults.

---

## 7. Controls and integrity

- **S5:** exact retained wording present at line 186.
- **S5-scope:** one table row at line 188; remains **OPEN**, not a gate.
- **E1–E9:** nine rows present and unresolved.
- **U1–U8:** eight rows present and **OPEN**; CP-1 is separately tracked as **OPEN**.
- **CP-1:** remains **OPEN** at line 385 and §13; no new public `crashed` route.
- **Terminal-set definition:** remains **OPEN** in §3.4, the Case A heading, Case A line 323, and B-D line 356.
- **Public states:** source `TTSProviderState` at `internal/audio/types/provider.go:23–35` contains exactly the thirteen values: `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, `backoff`. The target adds no public state.
- **Internal-only labels:** `outcome_unknown`, `terminated_by_action`, and `quarantined` remain internal and are not public provider states.
- **Tables:** eight tables; all rows have consistent column counts.
- **Trailing whitespace:** none detected.
- **Authority:** original P17, decision register, and contract rev3 remain at the frozen hashes above.
- **Isolation limitation:** the target and the prior review are untracked. The workspace also contains unrelated modified files. Git therefore cannot independently establish that only the proposal changed during the B-2 pass. The reported single-file isolation remains a reported result, not a repository-cleanliness proof.
- **No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or target edits were performed.**

---

## 8. Disposition

**B-2:** PASS. The final C-B gate carries both required exclusions in both ordinary branches. The `unloaded` counterexample is closed by a class-wide rule, not a hard-coded state.

**R-1–R-6:** PASS at the audited routes. R-5 terminal precedence and R-6 separate-failure/CP-1 routing remain intact.

**C-A, C-C, CP-1, P4-D, B-1:** preserved. No new public projection or state was introduced.

**N-5/N-6:** remain **OPEN** scope limitations; §13 is consistent with lines 368–370.

**Residual non-blocking observations:** B-B’s broad separate-value wording and C-B’s retained internal-only classification under P4-D remain separate clarity items. The Git isolation limitation remains unresolved by design.

**P17 remains UNAPPROVED. Contract rev3 remains frozen. C1 implementation remains UNAUTHORIZED.**

This is a review artifact, not an approval, not a correction pass, and not implementation authorization. Any further change requires a fresh review bound to its exact path, line count, and SHA-256.
