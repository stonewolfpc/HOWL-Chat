# Independent Adversarial Review — P17 Rev2 Candidate (536 lines)

**Review status:** COMPLETE — **BLOCKING N-6 SCOPE CONTRADICTION; R-5/R-6 PASS; NOT AN APPROVAL.**

**Reviewed artifact:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`

**Reviewed identity:** **536 lines**, SHA-256 `d3ca43a4feefa0c9f1108f3499b4b7796dfb0908c0acd1a0bfa6d5570dad6b40`

**Prior review:** `docs/review/REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-532-2026-11-09.md`, 250 lines, SHA-256 `1910534f4b174845a41c899a670958458eb2362f9bf9cd2980e28212236d9aef`. That review applies only to the 532-line predecessor and is not an approval of this revision.

**Frozen inputs checked:** original P17, 166 lines, `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf`; decision register, 434 lines, `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`; contract rev3, 183 lines, `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e`.

**Basis:** Direct read of the exact 536-line artifact; comparison against the 250-line review and the accepted OD boundaries; line-by-line audit of compatibility/applicability, C-B, C-A, C-C, C-D, P4-D, Case A, and the N-5/N-6 OPEN statements. No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or candidate edits were performed. This review does not approve P17 or authorize C1.

---

## 1. Verdict

**BLOCKING.**

R-5 and R-6 remain closed. N-1, N-2, N-3, N-4, and N-7 are corrected without an unintended new projection rule. N-5 remains correctly OPEN.

**N-6 is not consistently OPEN.** The new §13 statement says the treatment of a non-failure, non-active prior public provider-state value in C-B remains OPEN, but C-B’s ordinary fallback still permits that value to fall through to `absent` or `unavailable`. The specification therefore simultaneously says “OPEN” and supplies a default. That is a substantive scope contradiction, not a wording polish issue.

The candidate is **not approved** and is **not C1 authorization**.

---

## 2. Blocking finding

### B-1 — C-B still silently defaults the N-6 non-failure public-value combination

**Exact text:**

- C-B, target line 366: “A separate prior public provider-state value is not governed by C-B's `absent` or `unavailable` fallback.”
- C-B, target line 366: “Route a non-active prior public failure through C-A and an active operational state through C-C.”
- C-B, target line 368: “Only when no separate prior public value remains applicable, unresolved, or **CP-1 OPEN** does C-B proceed to the ordinary `absent`/`unavailable` outcome.”
- §13, target line 528: “Two scope limitations remain explicitly **OPEN**: selection among multiple prior failure classifications, and treatment of a separate non-failure, non-active prior public provider-state value in C-B.”

**Counterexample:**

1. C1 is established with publication `excluded-internal-only`.
2. There is no fatal observation and no terminal; this is Case C.
3. A separate prior public provider-state value V exists.
4. V is non-active and is not a prior public **failure** value.
5. V is not `crashed`, not applicable as a failure, and not CP-1 OPEN.

C-B line 366 says V is not governed by C-B’s fallback, but the only defined routes are for a non-active prior public **failure** (C-A) and an active operational state (C-C). V fits neither. Line 368 then allows C-B to proceed to `absent`/`unavailable` because V is not “applicable, unresolved, or CP-1 OPEN” in the failure-value sense. The result is a default for a case §13 says is OPEN.

The problem is not that C-B should invent a route for V. The problem is that it must **stop** and remain OPEN for V rather than allow the ordinary fallback to decide the combination. The current text does not state that stop condition.

**Required correction:**

Add an explicit C-B stop condition for any separate prior public provider-state value that is neither:

- a prior public failure value routed through C-A/C-C/Gate 5; nor
- an active operational value routed through C-C.

For that uncovered combination, C-B must remain **OPEN** and must not proceed to `absent` or `unavailable`. No preservation or projection rule should be invented.

This is the only blocking finding in the 536-line candidate.

---

## 3. Compatibility / applicability review

### 3.1 Four required distinctions

| Value status | Route required | Current text | Result |
|---|---|---|---|
| Compatible and applicable | Preserve the prior public failure where Gate 5 permits | P4-A line 300; C-A line 360; B-C line 352 | PASS |
| Compatibility or applicability unresolved | **OPEN**; no `absent`/`unavailable` substitution | P4-A line 300; C-A line 360; B-C line 352; C-B line 366 | PASS |
| Established incompatible | Historical evidence; ordinary OD-4 handling | C-A line 362; B-C line 352; C-B line 366 | PASS |
| `crashed` and otherwise compatible/applicable | **CP-1 OPEN**; not preservation or exclusion | P4-A line 300; C-A line 360; B-C line 352; C-B line 366 | PASS |
| Compatible but not applicable | Not preserved; Gate 5 passes; ordinary OD-4 handling | P4-A line 302; C-A line 360; B-C line 352 | PASS |

The new sentences correctly distinguish compatibility from applicability:

- P4-A line 302: “Compatibility does not establish applicability; a compatible but non-applicable value does not fail gate 5.”
- C-A line 360: “Compatibility alone does not establish applicability; a compatible but non-applicable value is not preserved and does not fail gate 5.”
- B-C line 352: “Compatibility does not establish applicability; a compatible but non-applicable value is not preserved and follows ordinary OD-4 exit handling.”

### 3.2 No unintended new projection rule found

The additions do not create a new public state or a new terminal definition. They restate Gate 5’s existing distinction:

- compatibility is necessary but not sufficient for preservation;
- a compatible value may still be non-applicable;
- non-applicability permits ordinary OD-4 handling only where the global predicate is otherwise satisfied.

The new C-A sentence is consistent with P4-A. The B-C sentence uses ordinary OD-4 handling rather than a cross-case P4-D rule. No unintended projection rule was found.

### 3.3 Publication versus justification

The new P4-A sentence at line 302 correctly distinguishes:

- proof that a value was published, which establishes publication status; and
- independent justification for using that value as current projection evidence.

It does not itself create a preservation rule. Gate 4 and Gate 5 remain separate. PASS for N-4.

---

## 4. R-5 and R-6 regression audit

### 4.1 R-5 — Case A precedence

The current text retains the terminal-precedence rule and the local terminal exceptions:

| Competing route | Line(s) | Result |
|---|---:|---|
| S7 no-orderly | 188 | PASS — Case A controls when a terminal exists |
| O3 no-orderly | 205 | PASS — explicit terminal exception before `otherwise unavailable` |
| O5 no-orderly | 207 | PASS — explicit terminal exception before `otherwise unavailable` |
| O4 `absent` excluded | 206 | PASS — Case A guards |
| S8 `absent` excluded | 189 | PASS — Case A guards |
| Gate 2 failure | 313 | PASS — active and CP-1 guards |
| P4-A terminal precedence | 310 | PASS — controls every branch below |
| P4-A Gate 3 no-orderly | 314 | PASS — now repeats the terminal exception |
| P4-A Gate 4 no-orderly | 315 | PASS — now repeats the terminal exception |
| P4-A no-orderly | 317 | PASS — explicit terminal exception |
| B-D | 354 | PASS — Case A guards |
| Case A | 323-328 | PASS — active and CP-1 guards; no new terminal set |

No `unavailable` fallback overrides an established terminal. **R-5 remains closed.**

### 4.2 R-6 — Separate prior public failure values

The named R-6 counterexamples remain closed:

| Separate prior public failure V | Current route | Result |
|---|---|---|
| Compatible and applicable | C-A preserves; `absent` excluded | PASS |
| Compatibility/applicability unresolved | **OPEN** | PASS |
| Established incompatible | Historical; ordinary OD-4 | PASS |
| `crashed`, otherwise compatible/applicable | **CP-1 OPEN** | PASS |
| Active operational state | C-C | PASS |
| No separate value | C-B ordinary outcome | PASS |

C-B no longer silently projects `unavailable` over an applicable, unresolved, incompatible, active, or CP-1 separate **failure** value. **R-6 remains closed for its named scope.**

The remaining problem is the broader N-6 combination of a non-failure public value, not the R-6 failure-value counterexamples.

---

## 5. N-1 through N-7 status

| Finding | Current evidence | Status |
|---|---|---|
| N-1 — C-D/B-C cross-reference | C-D line 374 now routes through C-A and Gate 5 | CORRECTED |
| N-2 — compatibility vs applicability | P4-A line 302; C-A line 360; B-C line 352 | CORRECTED |
| N-3 — Gate 5 filter wording | C-C line 370 says the active exclusion is an eligibility filter, not a separate gate | CORRECTED |
| N-4 — publication vs justification | P4-A line 302 distinguishes publication status from independent justification | CORRECTED |
| N-5 — multiple classifications | §13 line 528 retains a multi-classification selection rule as OPEN; no rule invented | OPEN, correctly |
| N-6 — non-failure public value | §13 line 528 says OPEN, but C-B lines 366-368 still provide an ordinary fallback | **BLOCKING B-1** |
| N-7 — local fallback wording | P4-A lines 314-315 repeat the terminal exception | CORRECTED |

N-5 is respected: the proposal does not adopt a selection rule. N-6 is not respected consistently because C-B’s fallback still decides the uncovered combination.

---

## 6. R-1 through R-4 regression audit

| Prior correction | Current evidence | Result |
|---|---|---|
| R-1: three-way establishment; unresolved blocks `absent` | P4-A lines 292-298; B-B line 342; Case C line 358; C-D line 374 | PASS |
| R-2: published / excluded-internal-only / unresolved are distinct | P4-A lines 299-302; C-B line 364 | PASS |
| R-3: CP-1 remains OPEN and is not indirectly resolved | S7 line 188; S8 line 189; O3-O5 lines 205-207; P4-A lines 300, 310, 314-317; B-C line 352; C-A line 360; C-B line 366; Case A line 325; public-crash line 383 | PASS |
| R-4: active operational values cannot survive confirmed exit | P4-A line 300; S8 line 189; O4 line 206; B-C line 352; C-A line 360; C-C line 370; Case A line 324 | PASS |

No R-1 through R-4 regression was found.

---

## 7. Controls and document integrity

- **S5:** exact retained wording present at line 184.
- **S5-scope:** exactly one row; remains OPEN and not a gate.
- **E1–E9:** nine rows present and unresolved.
- **U1–U8:** eight rows present and OPEN; CP-1 is separately tracked as OPEN.
- **CP-1:** remains OPEN. No new public `crashed` route is selected.
- **Terminal-set definition:** remains OPEN.
- **Public states:** the source `TTSProviderState` set remains the thirteen values `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, and `backoff`.
- **Tables:** eight tables; all rows have consistent column counts.
- **Trailing whitespace:** none detected.
- **Authority:** original P17, decision register, and contract rev3 hashes are unchanged.
- **No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or target edits were performed.**

---

## 8. Disposition

**BLOCKING: B-1 / N-6.**

The 536-line candidate is a successful N-1–N-4 and N-7 correction, and it preserves R-1–R-6. It is not ready for approval because its own OPEN scope statement for a non-failure, non-active prior public provider-state value is contradicted by C-B’s ordinary `absent`/`unavailable` fallback.

The next authorized correction should be narrow: add the C-B stop/OPEN condition for that uncovered public-value combination, without inventing a preservation rule, a public state, or a multiple-classification selection rule.

**P17 remains UNAPPROVED. Contract rev3 remains frozen. C1 implementation remains UNAUTHORIZED.**

This is a review artifact, not an approval and not a correction pass. Any further correction must receive a fresh independent review bound to its exact final line count and SHA-256.
