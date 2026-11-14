# Independent Adversarial Review — P17 Rev2 Candidate (538 lines)

**Review status:** COMPLETE — **BLOCKING B-2; B-1 STOP IS PRESENT BUT NOT INTEGRATED; R-5/R-6 PASS; NOT AN APPROVAL.**

**Reviewed artifact:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`

**Reviewed identity:** **538 lines**, SHA-256 `0cf14a2a4bdb1a73d34fbdf32264aa7bef8752837a42b90c7998842b21ca3dca`

**Prior review:** `docs/review/REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-536-2026-11-09.md`, 201 lines, SHA-256 `fd9e8ee98e550ea2419b9bcb0079116d2509e1783c61b5a25eed726c3a32fc78`. That review applies only to the superseded 536-line candidate and is not an approval of this revision.

**Frozen inputs checked:** original P17, 166 lines, `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf`; decision register, 434 lines, `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`; contract rev3, 183 lines, `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e`.

**Basis:** Direct read of the exact 538-line artifact, including the B-1 addition and the complete C-B/C-A/C-C/C-D/P4-D route text. The review compared the new sentence against the final C-B fallback rather than accepting the correction report’s route table or its “no route now reaches” claim. No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or candidate edits were performed. This review does not approve P17 or authorize C1.

---

## 1. Verdict

**BLOCKING.**

R-5 and R-6 remain closed at their named routes. The B-1 sentence was added at target line 368 and correctly states that a separate prior public value that is neither a prior public failure nor an active operational state does not proceed to the ordinary `absent`/`unavailable` outcome.

However, B-1 is not integrated into C-B’s final gate. At target line 370, C-B still says it proceeds whenever no separate prior public value “remains applicable, unresolved, or CP-1 OPEN.” A non-failure, non-active value is precisely a value that can be present yet not applicable in the failure-value sense. For that exact value, line 368 says **OPEN**, while line 370 still supplies the ordinary fallback. The candidate therefore has two competing routing rules for the N-6 combination, with no stated precedence rule.

The candidate is **not approved** and is **not C1 authorization**.

---

## 2. Blocking finding

### B-2 — The B-1 stop is not applied by the final C-B gate

**Exact text:**

- C-B, target line 368: “If a separate prior public provider-state value is present and is neither a prior public failure value nor an active operational state, C-B does not proceed to the ordinary `absent`/`unavailable` outcome; the combination remains OPEN.”
- C-B, target line 370: “Only when no separate prior public value remains applicable, unresolved, or **CP-1 OPEN** does C-B proceed to the ordinary `absent`/`unavailable` outcome: with orderly shutdown and the other gates satisfied, `absent` is eligible; without orderly shutdown, `unavailable` applies.”
- §13, target line 530: “Two scope limitations remain explicitly **OPEN**: selection among multiple prior failure classifications, and treatment of a separate non-failure, non-active prior public provider-state value in C-B.”

**Counterexample:**

1. C-B is established: the prior failure classification is established and its publication status is positively resolved as `excluded-internal-only`.
2. No fatal observation and no terminal exist, so Case C/C-B is the relevant case.
3. A separate prior public provider-state value **V = `unloaded`** is established by publication evidence.
4. `unloaded` is neither a prior public failure value nor an active operational state. It is a member of the frozen thirteen-value provider-state set.
5. Orderly shutdown is independently established and the other C-B gates are satisfied.

Line 368 says C-B does not proceed and the combination remains **OPEN**. But V is not a prior public failure value, so it does not “remain applicable” in the failure-value sense, is not unresolved, and is not CP-1 OPEN. On the literal final-gate condition at line 370, nothing in that paragraph excludes V, and the paragraph then states that **`absent` is eligible** once its condition is met. Without orderly shutdown, the same wording makes **`unavailable`** eligible.

A reader could try to preserve line 368 by treating line 370’s “Only when” as a necessary rather than sufficient condition. The artifact supplies no such precedence rule, however, and the correction report expressly claims that no C-B route reaches the ordinary outcome. A normative routing document cannot leave the answer dependent on which paragraph the reader treats as controlling. The same conflict exists with any other non-failure, non-active public value, not just `unloaded`. The final paragraph does not repeat the line-368 stop, does not say “and no separate non-failure, non-active value is present,” and does not establish that the stop controls the ordinary fallback.

**Why this is blocking rather than a wording polish issue:**

- The artifact’s own §13 says the treatment of this value in C-B remains **OPEN**.
- Line 368 now says the same value remains **OPEN**.
- Line 370 nevertheless defines the exact ordinary outcome that the OPEN statement excludes.
- The correction report’s claim that “No C-B route now reaches ordinary `absent`/`unavailable`” is not true on the face of the final C-B gate.
- A reader or implementer following the final gate can project `absent` or `unavailable` over V, while a reader following the stop condition keeps the outcome OPEN. The specification has no tie-breaker.

**Required correction:**

Integrate the stop into the final C-B gate. The narrow correction is to qualify line 370 so that it proceeds only when there is **no separate prior public value present that is neither a prior public failure value nor an active operational state**, and no separate prior public failure value remains applicable, unresolved, or CP-1 OPEN. The final fallback must not be reachable while the uncovered non-failure, non-active value is present.

This should not invent a preservation rule, a public projection, or a multiple-classification selection rule. The only required outcome for the uncovered combination is **OPEN**.

This is the only blocking finding in the 538-line candidate. R-5/R-6 and the prior N-1–N-4/N-7 corrections are otherwise preserved.

---

## 3. B-1 review

### 3.1 The added sentence is correct in isolation

Target line 368 correctly:

- separates the established internal-only classification from a separate prior public provider-state value;
- applies C-B only when no fatal observation and no terminal exist;
- routes a non-active prior public **failure** through C-A;
- routes an active operational state through C-C;
- identifies the non-failure, non-active public-value combination as an **OPEN** stop;
- preserves the CP-1 branch for a prior public `crashed` value;
- reduces an incompatible prior public failure to historical evidence and routes ordinary confirmed-exit handling through P4-D.

The sentence does not create a new public state, a new preservation rule, or a terminal definition.

### 3.2 The stop is not the final C-B rule

The B-1 correction adds a sentence before the old final gate but does not change that gate. The final gate remains broader than the stop and continues to supply the ordinary fallback for a value that is not a failure value. The correction is therefore incomplete, which is B-2.

### 3.3 No unintended new projection rule found in the added sentence

The new stop itself only withholds the ordinary outcome; it does not select `absent`, `unavailable`, preservation, or a replacement state. The problem is the retained final fallback, not the new sentence’s intent.

---

## 4. R-5 regression audit

| Competing route | Line(s) | Result |
|---|---:|---|
| S7 no-orderly | 190 | PASS — Case A controls when an established terminal exists |
| O3 no-orderly | 207 | PASS — explicit terminal exception before the fallback |
| O5 no-orderly | 209 | PASS — explicit terminal exception before the fallback |
| O4 `absent` excluded | 208 | PASS — Case A guards |
| S8 `absent` excluded | 191 | PASS — Case A guards |
| Gate 2 failure | 315 | PASS — active and CP-1 guards |
| P4-A terminal precedence | 312 | PASS — controls every branch below |
| P4-A Gate 3 no-orderly | 316 | PASS — local terminal exception present |
| P4-A Gate 4 no-orderly | 317 | PASS — local terminal exception present |
| P4-A no-orderly | 319 | PASS — explicit terminal exception |
| B-D | 356 | PASS — Case A guards |
| Case A | 325–330 | PASS — active and CP-1 guards; no terminal-set definition |

No `unavailable` fallback overrides an established terminal. **R-5 remains closed.**

---

## 5. R-6 regression audit

The named separate prior-public-**failure** routes remain correct:

| Separate prior public failure V | Current route | Result |
|---|---|---|
| Compatible and applicable | C-A preserves; `absent` excluded | PASS |
| Compatibility/applicability unresolved | **OPEN** | PASS |
| Established incompatible | Historical; ordinary P4-D | PASS |
| `crashed`, otherwise compatible/applicable | **CP-1 OPEN** | PASS |
| Active operational state | C-C | PASS |
| No separate value | C-B ordinary outcome | PASS |

C-B no longer silently projects `unavailable` over an applicable, unresolved, incompatible, active, or CP-1 separate **failure** value. **R-6 remains closed for its named scope.**

B-2 is outside the named R-6 failure-value scope: it concerns a non-failure, non-active separate public value. It is nevertheless a blocking open-scope consistency defect because §13 and line 368 make that combination OPEN.

---

## 6. N-1 through N-7 status

| Finding | Current evidence | Status |
|---|---|---|
| N-1 — C-D/B-C cross-reference | C-D line 376 now routes through C-A and Gate 5 | CORRECTED |
| N-2 — compatibility vs applicability | P4-A line 304; C-A line 362; B-C line 354 | CORRECTED |
| N-3 — Gate 5 filter wording | C-C line 372 says the active exclusion is an eligibility filter, not a separate gate | CORRECTED |
| N-4 — publication vs justification | P4-A line 304 distinguishes publication status from independent justification | CORRECTED |
| N-5 — multiple classifications | §13 line 530 retains a multi-classification selection rule as OPEN; no rule invented | OPEN, correctly |
| N-6 — non-failure public value | Line 368 adds the stop, but line 370 still supplies the ordinary fallback | **BLOCKING B-2** |
| N-7 — local fallback wording | P4-A lines 316–317 repeat the terminal exception | CORRECTED |

### 6.1 Non-blocking residual: C-B/P4-D retained internal-only classification

C-B line 368 routes an incompatible separate public failure to P4-D while saying “the internal-only classification remains retained.” P4-D line 380 begins: “No applicable prior terminal, fatal observation, or failure classification remains.”

An `excluded-internal-only` classification is not a public failure value, so this can be read as consistent: the classification is retained internally but is not applicable as a public failure. The wording is nevertheless ambiguous because P4-D says “failure classification” without the “public/applicable” qualifier. If the author wants the route to be self-evident, P4-D should say “no applicable prior public failure classification” or C-B should explain why the retained internal-only classification does not defeat P4-D’s precondition.

**Disposition:** non-blocking route-clarity finding; it does not create a current public projection.

### 6.2 Scope observation: B-B’s broad “separate prior public value” phrase

B-B line 350 says that if a separate prior public value exists, “B-C/Gate 5 governs it (or C-C if it is an active operational state).” B-C is expressly limited to a previously published public **failure** value. A non-failure, non-active public value is not expressly routed by B-C. This is not the B-1/C-B counterexample and is outside the B-1 scope, but it is the same kind of broad-value wording. If the owner intends the separate-public-value rule to extend beyond C-B, this needs a separate scope decision; this review does not treat it as a new blocking finding for the C-B-only B-1 correction.

---

## 7. R-1 through R-4 preservation

| Prior correction | Current evidence | Result |
|---|---|---|
| R-1: three-way establishment; unresolved blocks `absent` | P4-A lines 294–300; B-B line 344; Case C line 360; C-D line 376 | PASS |
| R-2: published / excluded-internal-only / unresolved are distinct | P4-A line 301; C-B line 366 | PASS |
| R-3: CP-1 remains OPEN and is not indirectly resolved | S7 line 190; S8 line 191; O3–O5 lines 207–209; P4-A lines 302, 312, 316–319; B-C line 354; C-A line 362; C-B line 368; Case A line 327; public-crash line 385 | PASS |
| R-4: active operational values cannot survive confirmed exit | P4-A line 302; S8 line 191; O4 line 208; B-C line 354; C-A line 362; C-C line 372; Case A line 326 | PASS |

No R-1 through R-4 regression was found. `S5` remains verbatim at line 186; `S5-scope` is a single OPEN tracking row at line 188 and is not a gate.

---

## 8. Controls and document integrity

- **S5:** exact retained wording present at line 186.
- **S5-scope:** one row; remains OPEN and not a gate.
- **E1–E9:** nine rows present and unresolved.
- **U1–U8:** eight rows present and OPEN; CP-1 is separately tracked as OPEN.
- **CP-1:** remains OPEN at lines 385 and 487. No new public `crashed` route is selected.
- **Terminal-set definition:** remains OPEN in §3.4, the Case A heading, Case A line 323, and B-D line 356.
- **Public states:** the target retains the frozen thirteen-value set; no public state is added, removed, or renamed.
- **Tables:** eight tables; all rows have consistent column counts.
- **Trailing whitespace:** none detected.
- **Authority:** original P17, decision register, and contract rev3 hashes are unchanged.
- **No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or target edits were performed.**

---

## 9. Disposition

**BLOCKING: B-2 / N-6.**

The B-1 sentence is a correct start but is not the final C-B rule. The final C-B gate still allows the ordinary `absent`/`unavailable` outcome for the exact non-failure, non-active public-value combination that line 368 and §13 declare OPEN. The next authorized correction should be narrow: qualify the final C-B gate so the line-368 stop controls it. No preservation rule, public state, or multi-classification selection rule should be invented.

R-5, R-6, R-1–R-4, N-1–N-4, and the local N-7 fix otherwise pass. N-5 remains correctly OPEN; the C-B/P4-D retained-classification wording is a non-blocking clarity observation.

**P17 remains UNAPPROVED. Contract rev3 remains frozen. C1 implementation remains UNAUTHORIZED.**

This is a review artifact, not an approval and not a correction pass. Any further correction must receive a fresh independent review bound to its exact final line count and SHA-256.
