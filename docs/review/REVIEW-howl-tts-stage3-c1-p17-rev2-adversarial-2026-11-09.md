# Independent Adversarial Review — P17 Rev2 Candidate (523 lines)

**Review status:** COMPLETE — **BLOCKING FINDINGS; NOT AN APPROVAL.**

**Reviewed artifact:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`

**Reviewed identity:** **523 lines**, SHA-256 `9cd45c468ffbaeca1512e1cfeaeba312397b7300dc7569e5b0595eae77ac3e36`

**Comparison artifact:** 513-line candidate, SHA-256 `49e9ecc2c812a361d03e739fbb571bfacc61f93fe45cffd30a8cb0bddb1055d6`

**Frozen inputs checked:** original P17, 166 lines, `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf`; decision register, 434 lines, `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`; contract rev3, 183 lines, `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e`

**Basis:** Direct read of the exact 523-line artifact; direct source-anchor audit; comparison against the accepted OD-1 through OD-9 boundaries and frozen contract. No tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or target edits were performed. This review does not approve P17 or authorize C1.

---

## 1. Verdict

**BLOCKING.**

The named R-1 and R-2 corrections are substantively effective. The named R-3 and R-4 guards are present at their explicit Case A entry points. However, the correction is not globally coherent:

- **R-5:** S7 and the O3/O5 no-orderly rows still permit an established terminal to be projected `unavailable`, bypassing Case A's terminal guards. This is a remaining active-terminal and terminal-immutability bypass.
- **R-6:** C-B can route an established internal-only classification to `unavailable` even when a separate prior public failure value exists. That can discard an applicable prior public failure and can bypass CP-1 for a prior public `crashed` value. C-B's `absent` branch also does not explicitly route a separate prior public value through Gate 5/CP-1.

The target therefore does not yet constitute a coherent correction of R-3 and R-4 across every route, and CP-1 is not globally preserved.

**Disposition:** P17 remains **UNAPPROVED**. Contract rev3 remains **FROZEN**. C1 remains **UNAUTHORIZED**. No correction pass is performed or authorized by this review.

---

## 2. Blocking findings

### R-5 — BLOCKING: the no-orderly fallback in S7/O3/O5 bypasses Case A's terminal guards

**Exact text:**

- S7, target line 183: “A confirmed process exit without established orderly shutdown projects `unavailable` unless the proposed P4-A preserves an independently justified prior **public** failure value.”
- O3, target line 200: after routing an established terminal to Case A, the same row says “Preserve an established, applicable prior public failure; ... otherwise `unavailable` applies.”
- O5, target line 202: the same “otherwise `unavailable` applies” fallback.
- Case A, target lines 316–321: after exit, an established terminal is preserved unless the active-terminal guard or CP-1 guard applies; an active terminal is **OPEN**, a `crashed` terminal is **CP-1 OPEN**, and a non-active/non-`crashed` terminal is preserved.

**Counterexample 1 — active terminal:**

1. A prior immutable terminal decision exists for the generation, with public value `ready`.
2. The process exit is confirmed.
3. Orderly shutdown is not established.
4. No separate prior public failure value is established.

Case A's active-terminal guard (line 317) says the active value is not preservable, the immutable terminal remains, and the public projection is **OPEN**. The O3/O5 route-A text also sends the terminal through Case A. But S7 line 183 and the O3/O5 “otherwise `unavailable`” fallback project `unavailable`. Both cannot be correct.

**Counterexample 2 — non-failure terminal:**

1. A prior immutable terminal decision exists with public value `absent` (the terminal set is explicitly OPEN, so this is not excluded by the target).
2. The process exit is confirmed.
3. Orderly shutdown is not established.
4. There is no prior public failure value.

Case A line 320 preserves the terminal; S7 line 183 projects `unavailable`. This is a direct conflict even without an active state.

**Counterexample 3 — incompatible terminal:**

1. A prior immutable terminal decision exists with public value `incompatible`.
2. The process exit is confirmed.
3. Orderly shutdown is not established.
4. The terminal value is established as post-exit-incompatible.

Case A preserves the immutable terminal; S7's only exception is an independently justified prior **public failure value that P4-A preserves**. An incompatible value is not preserved by P4-A gate 5, so S7 produces `unavailable` while Case A preserves the terminal.

**Why this is blocking:** The target says every route into Case A applies the active-terminal and CP-1 guards, including Gate 2 and B-D. S7 is a competing publication rule that does not defer to Case A. O3/O5 both route to Case A and then repeat a no-orderly fallback that does not include a terminal exception. The result is an established terminal with two possible public outcomes, and an active terminal can be replaced by `unavailable` rather than held OPEN.

**Required correction:** Make the no-orderly fallback conditional on the absence of an established immutable terminal, or state explicitly that Case A's terminal projection controls whenever Gate 2 fails. S7 and the O3/O5 rows must say:

> If an established immutable terminal exists, Case A controls and the no-orderly `unavailable` fallback does not apply.

The correction must not define the terminal set; it only needs to give Case A precedence over the fallback.

---

### R-6 — BLOCKING: C-B can bypass a separate prior public failure and CP-1

**Exact text:**

- C-B, target line 357: “If its establishment, threshold, publication status, and public value are resolved—including an established internal-only classification with publication excluded and no public value—it does not defeat P4-A gate 5: with orderly shutdown and the other gates satisfied, `absent` is eligible; without orderly shutdown, `unavailable` applies.”
- The same line says: “If establishment, threshold, publication, or public value is unresolved, ... the no-orderly-shutdown outcome remains `unavailable` because no public failure value is established.”
- Case C introduction, target line 351: “The three subcategories below are exhaustive for this case. A prior classification that remains unresolved as to establishment fails Gate 3 and remains **OPEN**; absent an independently established prior public value, it must not be routed through C-A, C-B, or C-C to `absent`.”
- P4-A gate 5, target line 295, and the no-orderly projection, target line 310, require an applicable prior public failure to be preserved; a prior public `crashed` value is CP-1 OPEN.

**Counterexample 1 — separate prior public value:**

1. Classification C1 is established, and its publication status is positively `excluded-internal-only`. C1 has no public value.
2. A separate prior public failure value V is independently established, post-exit-compatible, and applicable. V is not active and is not `crashed`.
3. There is no fatal observation and no terminal; this is Case C.
4. The process exit is confirmed; orderly shutdown is **not** established.

C-B line 357 says `unavailable` applies. P4-A gate 5 and the no-orderly projection say V must be preserved; OD-4 says `unavailable` is not permission to discard a more specific, independently justified prior failure. The outcome is contradictory.

**Counterexample 2 — CP-1:**

1. Same as above, but V is a previously published public `crashed`.
2. The exit is confirmed; orderly shutdown is not established.

C-B says `unavailable`. P4-A gate 5 says a prior public `crashed` value that would otherwise be compatible and applicable is **CP-1 OPEN**; the no-orderly projection and S7 repeat CP-1 OPEN. C-B therefore indirectly decides a CP-1 route against preservation by projecting `unavailable`, contrary to S14's statement that CP-1 is not decided by implication.

**Counterexample 3 — orderly shutdown:**

1. C1 is established and `excluded-internal-only`.
2. Separate prior public V is applicable, or V is `crashed` and otherwise compatible/applicable.
3. Exit is confirmed; orderly shutdown is established.

C-B says `absent` is eligible if “the other gates” are satisfied. Gate 5 is not satisfied for V, so `absent` is not eligible. But C-B does not state the resulting route: preserve V, or CP-1 OPEN for `crashed`. The global P4-A projection supplies the answer, but C-B's own text is incomplete and its no-orderly counterpart is wrong.

**Why this is blocking:** The target explicitly applies Gate 3/Gate 4 to “every prior failure classification” and Gate 5 to prior public values. A record can contain an internal-only classification and a separate public failure value. C-B does not defer to C-A/B-C/Gate 5 when that happens. The C-B `unavailable` output can discard an applicable prior public failure and can bypass CP-1. This also contradicts the target's claim that the three Case C subcategories are exhaustive and that a prior public value is not routed to `absent` by silence. If “no public value” is intended to make C-B inapplicable whenever a separate public value exists, then Case C has no route for that combination instead; the gap remains either way.

**Required correction:** Define Case C selection precedence, or make C-B explicitly defer:

- any prior public failure value is routed through C-A/B-C and Gate 5;
- a prior public `crashed` value triggers CP-1 OPEN;
- C-B's `unavailable` fallback applies only when no applicable or unresolved prior public value exists.

The same audit should check the B-B `excluded-internal-only` branch (target line 341), which has the same single-classification assumption. Its no-orderly branch already has the correct prior-public-value caveat, but its orderly branch does not explicitly state the preserve/CP-1 outcome when a separate public value exists.

---

## 3. Gate 3/4 partition

The partition is closed for a single prior failure classification. Gate 3 covers establishment; Gate 4 covers publication; Gate 5 covers compatibility and applicability. The previous R-1 and R-2 gaps do not recur for a single classification.

| Establishment | Publication status | Gate 3 | Gate 4 | Result / route | Disposition |
|---|---|---|---|---|---|
| No classification advanced | n/a | pass | n/a | Gate 5 and ordinary Case B-A/P4-D | PASS |
| Unresolved | n/a | **fail** | not reached | orderly shutdown: **OPEN**; no orderly: `unavailable` | PASS |
| Rejected | n/a | pass | n/a | no prior failure; P4-D, unless a separate prior public value remains | PASS |
| Established | `excluded-internal-only` | pass | pass | C-B; `absent` if orderly and no separate public value; no orderly `unavailable` | **R-6 when a separate public value exists** |
| Established | `published`, value identified | pass | pass | Gate 5; preserve compatible/applicable, OPEN if unresolved, historical if incompatible; `crashed` → CP-1 | PASS for the single-classification route |
| Established | `published`, value unidentified | pass | **fail** | orderly: **OPEN**; no orderly: `unavailable` | PASS |
| Established | `unresolved` (including missing evidence) | pass | **fail** | orderly: **OPEN**; no orderly: `unavailable` | PASS |
| Established | published active operational state | pass | pass | active filter at Gate 5; C-C; no active state survives; `absent` if orderly | PASS |
| Established | published `crashed`, compatible/applicable | pass | pass | **CP-1 OPEN** | PASS |
| Established | published `crashed`, incompatible | pass | pass | historical; ordinary `absent`/`unavailable` | PASS |
| Established | published compatible, applicability unresolved | pass | pass | Gate 5 fails; **OPEN** | PASS |

**Conclusion:** Gate 3/4 is not the remaining problem. The remaining problem is case routing when a separate prior public value exists or when a terminal reaches S7/O3/O5.

---

## 4. Publication outcome audit

The target's three publication outcomes are present and distinct.

| Publication outcome | Gate 4 | Required downstream handling | Audit |
|---|---|---|---|
| `published` with public value identified | pass | compatibility and applicability under Gate 5; preserve/OPEN/historical; `crashed` → CP-1 | PASS for a single classification |
| `excluded-internal-only` | pass | no public value from that classification; C-B; `absent` if orderly; no orderly `unavailable` | **R-6 if a separate public value exists** |
| `unresolved` | **fail** | orderly **OPEN**; no orderly `unavailable`; missing evidence is not exclusion | PASS |

Additional required distinctions:

- Published, value active: C-C; no active state survives; `absent` if orderly. PASS.
- Published, value incompatible: historical; ordinary `absent`/`unavailable`. PASS.
- Published, value compatible and applicable: preserve. PASS.
- Published, value compatible but **not** applicable: P4-A Gate 5 passes, so ordinary `absent`/`unavailable` follows. C-A/B-C do not explicitly state this branch; see N-2.
- Published, value compatibility or applicability unresolved: **OPEN**. PASS.
- Published `crashed`, compatible/applicable: **CP-1 OPEN**. PASS, except the C-B route in R-6.
- Published `crashed`, incompatible: historical; no CP-1 preservation. PASS.
- Merely asserted or internal-only without positive exclusion: not established/not excluded; Gate 3 or Gate 4 blocks `absent`. PASS.
- Missing publication evidence: `unresolved`, not “not published.” PASS.

**Conclusion:** The three outcomes themselves are correctly separated. The C-B route fails to honor a separate prior public value, not the publication-status taxonomy.

---

## 5. Case A entry-route audit

| Route | Target line | Terminal guards applied? | Audit |
|---|---|---|---|
| Gate 2 fails | 306 | yes: active guard and CP-1 guard | PASS |
| O3 route A | 200 | yes in route-A clause; no-orderly fallback still says “otherwise `unavailable`” | **R-5 ambiguity** |
| O5 route A | 202 | yes in route-A clause; no-orderly fallback still says “otherwise `unavailable`” | **R-5 ambiguity** |
| O4, `absent` excluded | 201 | yes: Case A CP-1/active guards | PASS |
| S8, `absent` excluded | 184 | yes: Case A CP-1/active guards | PASS |
| B-D | 347 | yes: active guard and CP-1 guard | PASS |
| Public `crashed` section | 372 | yes: Case A routes terminal `crashed` to CP-1 | PASS |
| S7 no-orderly fallback | 183 | **no terminal exception**; projects `unavailable` unless a prior public failure is preserved | **R-5 BLOCKING** |
| Case A no exit | 316 | not applicable; preserves terminal | PASS |
| Case A confirmed exit, active terminal | 317 | active guard; retains terminal; public **OPEN** | PASS |
| Case A confirmed exit, `crashed` terminal | 318 | CP-1 guard; public **OPEN** | PASS |
| Case A confirmed exit, no guard, orderly | 319 | preserves terminal; no `absent` | PASS |
| Case A confirmed exit, no guard, no orderly | 320 | preserves terminal; `crashed` only internal | PASS |

**Conclusion:** Every explicit Case A entry now contains the guards. The defect is that S7 remains a competing no-orderly rule and O3/O5 retain a generic fallback that does not defer to Case A. This is not a missing guard at an explicit entry; it is a bypass around the entry.

---

## 6. Rejected and unresolved classification routing

| Classification state | Route | Audit |
|---|---|---|
| Unresolved establishment, orderly shutdown | B-B line 338; Gate 3; **OPEN** | PASS |
| Unresolved establishment, no orderly shutdown | `unavailable` | PASS |
| Rejected, no prior public value | B-B line 339; C-D line 363; P4-D | PASS |
| Rejected, separate prior public value | B-B line 339 refers to B-C/Gate 5; C-D line 363 refers to B-C/Gate 5 | PASS in outcome, but C-D's B-C reference is a cross-case label; see N-1 |
| Established, publication unresolved | B-B line 340; **OPEN** with orderly | PASS |
| Established, publication excluded-internal-only | B-B line 341; C-B line 357 | **R-6 when a separate public value exists** |
| Established, publication published | B-C line 345 | PASS |
| Rejected classification with a prior public value | C-D line 363 | B-C is a Case B rule requiring a fatal observation; C-D has no fatal observation. The content still routes through Gate 5, but the cross-reference is wrong. N-1 |
| Prior classification rejected under accepted semantics | C-D line 363 | PASS |

**Conclusion:** Unresolved and rejected single-classification routing is coherent. The remaining defects are C-B's separate-public-value handling and the C-D cross-reference.

---

## 7. CP-1 route audit

| Route | CP-1 handling | Audit |
|---|---|---|
| Gate 2 failure → Case A | active guard then CP-1 guard; public **OPEN** | PASS |
| O3 route A | Case A CP-1 guard | PASS |
| O5 route A | Case A CP-1 guard | PASS |
| O4/S8 `absent` excluded | Case A CP-1 guard | PASS |
| B-D | Case A CP-1 guard | PASS |
| P4-A Gate 5 | compatible/applicable `crashed` → CP-1 **OPEN** | PASS |
| P4-A no-orderly | prior public `crashed` → CP-1 **OPEN** | PASS |
| B-C | compatible/applicable `crashed` → CP-1 **OPEN** | PASS |
| C-A | compatible/applicable `crashed` → CP-1 **OPEN** | PASS |
| B-B no-orderly | prior public `crashed` → CP-1 **OPEN** | PASS |
| C-D | prior public value via B-C/Gate 5, including CP-1 | PASS in outcome; N-1 label issue |
| C-B, no orderly, separate prior public `crashed` | says `unavailable` | **R-6 BLOCKING** |
| C-B, orderly, separate prior public `crashed` | says `absent` if “other gates”; does not state CP-1 | **R-6 gap** |
| Terminal `crashed` | Case A CP-1 guard | PASS |
| Public `crashed` summary | CP-1 remains OPEN | PASS |

**Conclusion:** CP-1 is genuinely open for every route that names it. It is not globally open because C-B's no-orderly fallback resolves a separate prior public `crashed` route to `unavailable`, and its orderly branch does not state CP-1 handling. This is exactly the kind of indirect settlement the review was asked to detect.

---

## 8. Terminal-set audit

The terminal-set definition remains genuinely OPEN:

- §3.4 target line 135 says the complete terminal set remains OPEN.
- Case A target line 314 says the section does not define the terminal set.
- Case A target line 317's active guard lists `ready`, active `degraded`, and any state whose meaning asserts current operation, but it does not enumerate the terminal set and expressly says it “does not resolve the terminal set.”
- B-D target line 347 says the guards do not decide CP-1 or define the terminal set.
- S8/O4 route terminal values through Case A without defining the terminal set.

The active-terminal guard does constrain **preservation** of active states; it does not decide which states may be terminal decisions. The terminal set is still an open question. PASS.

---

## 9. Mandatory reviewer challenges

| # | Challenge | Disposition |
|---|---|---|
| 1 | Verify every retained line and code claim against source; do not accept a classification from the draft. | **PARTIAL — all cited anchors checked; no line-by-line re-derivation of every retained sentence; no runtime verification.** |
| 2 | Challenge each fatal assignment, especially O6, O12, O15, O16, O18. | **PASS for cited anchors.** O6 is `internal_error → helper_crashed`, fatal, no kill for that code; O12 is returned only, no termination; O15 is handshake terminate then incompatible/unavailable; O16 is the nine-exit `validateHello` path; O18 publishes `crashed` while the process may be alive. No material discrepancy. |
| 3 | Review U1–U8 as open; no tests/builds/runtime experiments. | **PASS.** U1–U8 remain OPEN; CP-1 remains separate; no experiments were run. |
| 4 | Terminal immutability must not erase late evidence; old-generation evidence must not modify a replacement; every Case A entry applies active/CP-1 guards. | **PARTIAL — R-5.** Explicit Case A entries have the guards; S7/O3/O5's no-orderly fallback does not defer to them. |
| 5 | No public projection uses `outcome_unknown`, `quarantined`, or `terminated_by_action`; thirteen states unchanged. | **PASS.** The internal-only labels remain internal; the public enum is unchanged. |
| 6 | E1–E9 remain visible and unresolved; no implementation claim. | **PASS.** E1–E9 are retained; the target claims no implementation evidence. |
| 7 | Return finding-level dispositions; do not approve P17. | **PASS.** This report is finding-level and does not approve P17. |
| 8 | Record exact artifact path, revision, lines, SHA-256. | **PASS.** The identity is at the top of this report. |

The target's own reviewer challenge 4 (“every route into Case A applies the active-terminal and CP-1 guards, including Gate 2 and B-D”) is not fully satisfied because S7 is a competing route outside Gate 2/B-D.

---

## 10. Source-anchor audit

The cited source anchors were checked against the local source at the reviewed revision. No material discrepancy was found.

| Evidence | Source anchor | Audit |
|---|---|---|
| O3 unexpected exit | `supervisor.go:185–190` | `observeExit` publishes `crashed` for a non-`stopping` exit. PASS |
| O4 zero exit while stopping | `supervisor.go:178–180` | `observeExit` publishes `absent` for a nil error in the `stopping` branch. PASS |
| O5 nonzero exit while stopping | `supervisor.go:179–181` | `observeExit` publishes `absent` after setting exited. PASS |
| O6 internal_error → helper_crashed | `client.go:299–303` | `readLoop` rewrites `internal_error` and calls `failAll`; no kill for that code. PASS |
| O12 returned helper_crashed | `client.go:421` | `asFailure` path returns; no termination. PASS |
| O15 handshake | `client.go:248–253` | Terminates, waits, then publishes incompatible or unavailable. PASS |
| O16 validateHello | `protocol.go:502–545` | Nine validation exits. PASS |
| O18 live write failure | `client.go:756` | Publishes `crashed` on a send failure while the process may be alive. PASS |
| O19 ready/degraded | `client.go:771,783,793,796` | timeout, cancelled, ready, degraded paths. PASS |
| E1–E9 | cited source lines | Retained as unresolved evidence-loss findings; no claim that they are fixed. PASS |

No tests, builds, benchmarks, or runtime experiments were run, so this audit establishes document-to-anchor consistency only.

---

## 11. Non-blocking findings

### N-1 — C-D cross-references B-C despite no fatal observation

C-D (target line 363) says a prior public value is governed by “B-C,” but C-D is under Case C, which has no fatal observation, while B-C is the Case B rule for a fatal observation with a published public failure value. The substantive Gate 5/CP-1 treatment is correct; the cross-reference should be to C-A/C-C or to Gate 5 directly. **Non-blocking, but it should be corrected for route clarity.**

### N-2 — compatible-but-not-applicable prior public values have no explicit case branch

C-A and B-C explicitly handle compatible/applicable, unresolved, and incompatible values. A value established as compatible but **not applicable** is not “unresolved,” and Gate 5 passes, so the global P4-A projection yields `absent`/`unavailable`. The case text does not state this branch. The global predicate still prevents an incorrect `absent` when compatibility remains unresolved. **Non-blocking clarification.**

### N-3 — C-C's “operates at P4-A gate 5” wording should distinguish filtering from gate failure

C-C (target line 359) says the active-state exclusion “operates at P4-A gate 5.” Gate 5's actual text treats an active state as **not applicable**, so Gate 5 can hold while the active value is discarded. If “operates at gate 5” is intended to mean Gate 5 **fails**, C-C's subsequent “remaining gates are satisfied” instruction is inconsistent. The intended result (no active state; `absent` if orderly) is clear. **Non-blocking wording clarification.**

### N-4 — “published with public value identified” and “independently justified” should be explicitly related

Gate 4 passes when publication status is `published` with the public value identified. Gate 5 requires an “independently justified” prior public failure value. The target should state that “published with the public value identified” includes the evidence of publication necessary to make the value independently justified, or treat a published-but-unjustified value as `unresolved`. Otherwise a reader can pass Gate 4 and then treat the value as not applicable for Gate 5. **Non-blocking because the surrounding text generally treats publication evidence as justification; clarify to prevent a future gap.**

### N-5 — the multi-classification selection rule should be stated

The target applies Gate 3/Gate 4 to every prior failure classification and Gate 5 to prior public values, but the Case B/C subcategories are written in the singular. R-6 is the blocking consequence for C-B. Once C-B is corrected, add an explicit selection rule for records with multiple classifications and public values. **Non-blocking after R-6 is fixed.**

---

## 12. Status of the prior findings

| Prior finding | Status in the 523-line candidate |
|---|---|
| R-1: advanced classification + resolved threshold + unresolved establishment can reach `absent` | **CLOSED.** Three-way establishment; unresolved blocks `absent`; B-B/C routes it to OPEN. |
| R-2: “not published” conflated with “internal-only” | **CLOSED.** Published / excluded-internal-only / unresolved are distinct; missing evidence is unresolved. |
| R-3: Case A can preserve a terminal carrying a prior public `crashed` without CP-1 | **CLOSED for explicit Case A entries.** Case A has a CP-1 guard. Remaining CP-1 gap is C-B (R-6). |
| R-4: Case A can preserve an active terminal | **CLOSED for explicit Case A entries.** Case A has an active-terminal guard. Remaining active-terminal path is S7/O3/O5 (R-5). |

**Net:** R-1 and R-2 are closed. R-3 and R-4 are closed at the named Case A entries but not across every route. The correction pass therefore does not yet close the underlying bypass class.

---

## 13. Disposition and next step

**BLOCKING FINDINGS:** R-5 and R-6.

**Non-blocking:** N-1 through N-5.

**Current governance:**

- P17: **UNAPPROVED**.
- Contract rev3: **FROZEN**.
- C1 implementation: **UNAUTHORIZED**.
- CP-1: **OPEN**, except that C-B's current no-orderly text indirectly resolves one route; R-6 must be corrected.
- Terminal-set definition: **OPEN**; the Case A guards do not define it.
- U1–U8: **OPEN**.
- E1–E9: **visible and unresolved**.
- No target, contract, register, source, test, build, or runtime change was made by this review.

**Required next step:** a separately authorized correction pass must close R-5 and R-6, then rebind this artifact with a new line count and SHA-256 before another review. This report does not authorize that pass and does not approve P17 or C1.
