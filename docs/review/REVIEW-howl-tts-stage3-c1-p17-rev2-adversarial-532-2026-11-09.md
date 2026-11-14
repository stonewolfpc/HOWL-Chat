# Independent Adversarial Review — P17 Rev2 Candidate (532 lines)

**Review status:** COMPLETE — **NO R-5/R-6 BLOCKING RECURRENCE FOUND; NON-BLOCKING ROUTE-CLARITY FINDINGS REMAIN; NOT AN APPROVAL.**

**Reviewed artifact:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`

**Reviewed identity:** **532 lines**, SHA-256 `8c3b7df451aee1e19c6fa4b0902b7a75f7517034b63de8adc9acd35b819e551d`

**Prior blocking review:** `docs/review/REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-2026-11-09.md`, 343 lines, SHA-256 `78fa5b7ce45825c016e9afef7298212a63d588fd069f361488ee5c901177734f`. That review bound to the superseded 523-line artifact and its R-5/R-6 findings are re-audited here rather than accepted from its PASS/FAIL labels.

**Frozen inputs checked:** original P17, 166 lines, `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf`; decision register, 434 lines, `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`; contract rev3, 183 lines, `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e`.

**Basis:** Direct read of the exact 532-line artifact; line-by-line route audit of Case A, P4-A, S7, O3/O5, C-B, C-A, C-C, C-D, P4-D, and the no-orderly rules; a recheck that the frozen public state set remains thirteen values; no tests, builds, benchmarks, runtime experiments, source edits, contract edits, register edits, or target edits were performed. This review does not approve P17 or authorize C1.

---

## 1. Verdict

**No recurrence of the named R-5 or R-6 blocking bypasses was found.**

- **R-5: PASS.** The no-orderly-shutdown routes involving an established terminal now defer to Case A before any `unavailable` projection. The active-terminal and CP-1 guards remain applicable.
- **R-6: PASS for the named separate prior-public-failure case.** C-B now distinguishes the established internal-only classification from a separate prior public failure value and routes applicable, unresolved, incompatible, active, and CP-1 branches distinctly.

The candidate is **not approved** and is **not C1 authorization**. Several non-blocking route-clarity and open-scope findings remain; see §7. They are not the prior R-5/R-6 bypasses, but they should be resolved or explicitly retained as OPEN before the owner treats the routing text as stable.

The candidate's own header claim that the 532-line pass addresses R-5/R-6 is consistent with the text reviewed here; this review independently verifies the routing rather than accepting that claim.

---

## 2. R-5 — Case A precedence

### 2.1 Source rules checked

The relevant text is:

- S7, target line 186: an established terminal makes Case A control and makes the no-orderly `unavailable` fallback inapplicable.
- S8, target line 187: an excluded terminal is preserved only through Case A's CP-1 and active-terminal guards.
- O3, target line 203, and O5, target line 205: route A sends an established terminal to Case A, and a separate sentence makes the no-orderly fallback inapplicable when the terminal exists.
- O4, target line 204: an excluded `absent` routes an established terminal through Case A.
- P4-A terminal precedence, target line 308: whenever Gate 2 fails because an established terminal exists, Case A controls every projection branch below, and no `absent` or `unavailable` fallback overrides the active-terminal or CP-1 guards.
- Gate 2 failure, target line 311: Case A controls; both guards apply.
- P4-A no-orderly branch, target line 315: an established terminal makes Case A control and makes the `unavailable` fallback inapplicable.
- B-D, target line 352: an established terminal routes through Case A with both guards.
- Case A, target lines 321-326: no exit preserves the terminal; after exit, active and CP-1 are guarded, and a non-guarded terminal is preserved without `absent` overwriting it.
- **Case B, Case C, and Case D are defined by the absence of a terminal.** Their fallback text cannot silently overwrite a terminal because the terminal case does not enter those cases; if a terminal exists, Case A controls first.

### 2.2 Route audit

| Route that can compete with Case A | Line(s) | Terminal handling | Result |
|---|---:|---|---|
| S7 no-orderly fallback | 186 | Explicit Case A exception; terminal evidence immutable; guards govern projection | PASS |
| O3 route A and no-orderly text | 203 | Explicit Case A exception before `otherwise unavailable` | PASS |
| O5 route A and no-orderly text | 205 | Explicit Case A exception before `otherwise unavailable` | PASS |
| O4 `absent` excluded | 204 | Explicit Case A guard route | PASS |
| S8 `absent` excluded | 187 | Explicit Case A guard route | PASS |
| Gate 2 failure | 311 | Case A controls; active and CP-1 guards | PASS |
| P4-A global terminal precedence | 308 | Case A controls every branch below; no fallback override | PASS |
| P4-A no-orderly branch | 315 | Explicit terminal exception | PASS |
| P4-A Gate 3/Gate 4 no-orderly bullets | 312-313 | No local terminal sentence, but both are projection branches below line 308's global precedence | PASS; see N-7 |
| B-D | 352 | Case A controls; both guards; terminal set not defined | PASS |
| Case A no exit | 321 | Preserves the terminal; no contrary fallback | PASS |
| Case A active terminal after exit | 322 | Active value not preservable; public projection OPEN | PASS |
| Case A `crashed` terminal after exit | 323 | CP-1 OPEN; no preservation or exclusion decided | PASS |
| Case A non-guarded terminal | 324-325 | Preserves the terminal; orderly and no-orderly outcomes do not select `absent`/`unavailable` over it | PASS |
| Case A public projection | 326 | Subject to active and CP-1 guards | PASS |
| Case B/C/D no-terminal cases | 330-379 | Terminal is excluded by case definition; a terminal re-enters Case A | PASS |

### 2.3 R-5 conclusion

The prior counterexamples no longer hold in the current text:

- An active terminal is not replaced by `unavailable`; Case A holds the public projection OPEN.
- A non-failure, non-active terminal is not replaced by `unavailable`; Case A preserves the terminal.
- An incompatible terminal is not replaced by `unavailable`; Case A preserves the established terminal subject to the guards.
- A `crashed` terminal is not silently preserved or excluded; CP-1 remains OPEN.

**R-5 is closed at its named routes.** No new terminal-state definition was introduced, and the terminal-set definition remains OPEN. The generic `otherwise unavailable` phrases retained in O3/O5 are not bypasses because they occur after the explicit terminal exception.

---

## 3. R-6 — C-B and separate prior public values

### 3.1 C-B entry conditions

C-B, target lines 362-366, applies only when:

- the prior failure classification is established;
- its publication status is positively resolved as `excluded-internal-only`;
- no fatal observation exists;
- no terminal exists;
- the internal-only classification is not itself a public value.

The C-B text explicitly says that an internal-only classification does not establish that no separate prior public value exists (line 364), and that a separate public provider-state value is not governed by C-B's ordinary `absent`/`unavailable` fallback.

The final C-B condition is:

> Only when no separate prior public value remains applicable, unresolved, or CP-1 OPEN does C-B proceed to the ordinary `absent`/`unavailable` outcome.

That is the required gate: an applicable, unresolved, or CP-1 prior public value must interrupt C-B.

### 3.2 Separate prior-public-failure route matrix

| Separate prior public value V | C-B/C-A/C-C/P4-D route | Orderly shutdown | No orderly shutdown | Result |
|---|---|---|---:|---|
| None | C-B ordinary | `absent` if other gates hold | `unavailable` | PASS |
| Compatible and applicable, not active | C-A preserves V; `absent` excluded by Gate 5 | Preserve V | Preserve V | PASS |
| Compatibility or applicability unresolved | Explicit **OPEN** branch | **OPEN** | **OPEN** | PASS |
| Established incompatible | Historical evidence; ordinary P4-D; C1 retained | `absent` under P4-D | `unavailable` under P4-D | PASS in outcome; see N-7 |
| `crashed`, compatible and applicable | **CP-1 OPEN**; historical evidence; no preservation or exclusion decided | CP-1 OPEN | CP-1 OPEN | PASS |
| `crashed`, incompatible | Historical evidence; P4-D; CP-1 not engaged because the value is not otherwise preservable | `absent` under P4-D | `unavailable` under P4-D | PASS |
| Active operational value (including active `degraded`) | C-C; no active value survives | `absent` if remaining gates hold | `unavailable` | PASS |
| Compatible but not applicable | C-B ordinary; Gate 5 passes | `absent` | `unavailable` | PASS in outcome; see N-2 |
| Non-failure, non-active public provider-state value | No explicit C-A/C-C route in the text; final condition does not define whether it is applicable | Ambiguous | Ambiguous | **N-6 OPEN scope** |

### 3.3 C-B stress cases

1. **Internal-only classification plus separate applicable non-`crashed` value.**
   C-B routes the value through C-A and preserves it; `absent` is not eligible for that value. PASS.

2. **Internal-only classification plus separate `crashed` value.**
   C-B routes to **CP-1 OPEN** and retains the value as historical evidence. It does not select `absent`, `unavailable`, preservation, or exclusion. PASS.

3. **Internal-only classification plus separate value whose compatibility/applicability is unresolved.**
   C-B returns **OPEN**, including in the no-orderly branch; it does not fall through to `unavailable`. PASS.

4. **Internal-only classification plus separate incompatible value.**
   C-B reduces the value to history and uses ordinary P4-D. This is correct only if the internal-only classification is not itself an applicable public failure; the retained C1 classification creates a cross-route wording ambiguity. See N-7.

5. **Internal-only classification plus active operational value.**
   C-B routes to C-C; active values are not preserved, and the ordinary orderly/no-orderly result follows. PASS.

6. **Internal-only classification unresolved, separate applicable value.**
   The separate value is routed through C-A before the C-B ordinary branch; the internal-only classification's unresolved status does not erase the separate value. PASS.

7. **Internal-only classification resolved, no separate value.**
   C-B allows `absent` only with orderly shutdown and otherwise `unavailable`. PASS.

8. **Terminal or fatal observation also present.**
   C-B expressly defers to Case A or Case B. PASS.

### 3.4 R-6 conclusion

The named R-6 failure — C-B projecting `unavailable` over an applicable separate prior public failure or over CP-1 — is closed. The current text contains the required separate-value interruption and the CP-1 branch.

The remaining C-B issue is narrower: the text uses the broad term “separate prior public provider-state value” but defines explicit routes only for prior public **failure** values and active values. A non-failure, non-active public value has no explicit C-A/C-C route. That is an OPEN-scope finding, not a recurrence of the named R-6 counterexample. See N-6.

---

## 4. R-1 through R-4 regression audit

| Prior correction | Current evidence | Regression result |
|---|---|---|
| R-1: three-way establishment; unresolved blocks `absent` | P4-A line 290; gate 3 line 296; B-B line 343; Case C introduction line 356; C-D line 372 | PASS |
| R-2: published / excluded-internal-only / unresolved are distinct | P4-A gate 4 line 297; C-B entry conditions line 362; no “not published” conflation | PASS |
| R-3: CP-1 remains OPEN and is not indirectly resolved | S7 line 186; S8 line 187; O3 line 203; O4 line 204; O5 line 205; P4-A lines 298, 308, 314-315; B-C line 350; C-A line 358; C-B line 364; Case A line 323; public-crash section line 381 | PASS |
| R-4: active operational values cannot survive confirmed exit | P4-A gate 5 line 298; S8 line 187; O4 line 204; B-C line 350; C-A line 358; C-C line 368; Case A line 322 | PASS |

**No R-1 through R-4 regression was found.** The three-way establishment and publication states remain present and distinct. `S5` remains verbatim RETAINED, and `S5-scope` remains a separate OPEN tracking row, not a gate.

The previously noted C-D cross-reference to B-C remains (N-1). It is a route-label defect, not a regression in the substantive Gate 5/CP-1 outcome.

---

## 5. CP-1 and terminal-set audit

**CP-1 remains OPEN.** It is named as a separate tracking case in the unresolved-cases table (line 483), in S14 (line 193), and in the public-crash summary (line 381). The candidate does not select a new public `crashed` route.

**Terminal-set definition remains OPEN.** This is stated in §3.4, the Case A heading, Case A line 319, and B-D line 352. The active guard constrains preservation; it does not define which states may be terminal.

The Case A active guard does not assume that every terminal is safe: active and indeterminate-active values are not preserved; `crashed` is diverted to CP-1. The remaining non-active, non-`crashed` terminal preservation is the proposal's existing Case A disposition, not a new terminal-set definition.

---

## 6. Controls and document integrity

- **S5:** exact retained row present at line 182.
- **S5-scope:** exactly one table row; remains OPEN and not a gate.
- **E1–E9:** nine rows present and unresolved.
- **U1–U8:** eight rows present, all OPEN; CP-1 is separately tracked as OPEN.
- **Case structure:** Case A, Case B, Case C, and Case D present; C-D remains separate from Case C.
- **Tables:** eight tables; all rows have consistent column counts.
- **Trailing whitespace:** none detected.
- **Public state set:** the source `TTSProviderState` constants remain the thirteen values `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, and `backoff`. The target adds no public state.
- **Internal-only labels:** `outcome_unknown`, `terminated_by_action`, and `quarantined` remain internal and are not public provider states.
- **Authority:** original P17, decision register, and contract rev3 hashes remain `01d887c1…`, `3c64797b…`, and `ed0da1a7…` respectively.
- **No tests, builds, benchmarks, or runtime experiments were run.** Source-anchor labels were not re-derived as implementation evidence.

---

## 7. Non-blocking findings

### N-1 — C-D still cross-references B-C

C-D, target line 372, says an independently established prior public value remains governed by “B-C and Gate 5.” C-D is under Case C and has no fatal observation; B-C is the Case B rule titled “Fatal observation with a previously published public failure value.” The substantive Gate 5/CP-1 outcome is correct, but the label is cross-case. Prefer “C-A and Gate 5” or “Gate 5” directly.

**Disposition:** non-blocking route clarity; it should be corrected before the routing text is treated as final.

### N-2 — Compatible-but-not-applicable values lack an explicit case sentence

A prior public failure value can be compatible yet not applicable. Gate 5 passes for that value, so ordinary `absent`/`unavailable` follows. C-B's final condition treats the value as non-blocking, but C-A and B-C do not state the branch explicitly. This is consistent with the global predicate but is not self-evident from the case text.

**Disposition:** non-blocking clarification.

### N-3 — C-C's “operates at P4-A gate 5” wording conflates filtering with gate failure

C-C, target line 368, says the active-state exclusion “operates at P4-A gate 5.” Gate 5 describes an active state as not applicable, rather than necessarily failing the gate. The intended result — discard the active projection, then use the remaining orderly/no-orderly rules — is clear, but the wording can be read as saying Gate 5 fails while the remaining gates are still satisfied.

**Disposition:** non-blocking wording clarification.

### N-4 — “Published with the public value identified” and “independently justified” need an explicit relationship

Gate 4 passes when publication status is `published` with the public value identified. Gate 5 requires an independently justified prior public failure value. The text should state that publication evidence sufficient for Gate 4 also supports independent justification, or that a published-but-unjustified value remains `unresolved`. Otherwise a reader can treat a value as passing Gate 4 and failing Gate 5 for lack of the same evidence.

**Disposition:** non-blocking clarification; the surrounding text generally treats publication as evidence, but the distinction is not explicit.

### N-5 — Multi-classification selection remains underspecified

The gates apply to “every prior failure classification,” but the case branches are largely singular. C-B now handles one established internal-only classification plus one separate public value; it does not state a general selection rule when multiple classifications or multiple prior public values exist. The named R-6 combination is covered, but a general multi-classification rule remains OPEN.

**Disposition:** non-blocking; retain OPEN rather than inventing a selection rule.

### N-6 — C-B's “provider-state value” language is broader than its failure-value routes

C-B, target line 364, says a separate prior public provider-state value is not governed by C-B's ordinary fallback. It then defines routes for a non-active prior public **failure** and an active operational state. A non-active, non-failure public provider-state value has no explicit C-A/C-C route and no explicit OPEN statement. The final condition may treat it as not applicable, which could allow the ordinary `absent`/`unavailable` result without deciding what happened to that value.

This is outside the named R-6 counterexample, which concerns a separate prior public **failure** value. If the requirement is read literally to include every public provider-state value, this combination remains an OPEN governance gap.

**Disposition:** non-blocking for the named R-6 failure; explicitly OPEN if the broader reading is intended.

### N-7 — Local fallback text still depends on the global P4-A precedence rule

P4-A's Gate 3 and Gate 4 no-orderly bullets (lines 312-313) still say “OD-4's `unavailable` fallback applies” without repeating the terminal exception. The global terminal-precedence sentence at line 308 is explicit and does control every branch below, so this is not an R-5 bypass. Local repetition of the terminal condition would make the branch safe against excerpting.

C-B's incompatible-value branch also routes to P4-D while retaining the internal-only classification (line 364); the outcome is ordinary `absent`/`unavailable`, but P4-D's “no applicable … failure classification remains” wording is not obviously satisfied by the retained C1 classification.

**Disposition:** non-blocking route clarity.

---

## 8. Disposition

**R-5:** PASS at all audited competing no-orderly routes. No `unavailable` fallback overrides an established terminal.

**R-6:** PASS for the named separate prior-public-failure/CP-1 cases. No `unavailable` or `absent` fallback silently discards an applicable, unresolved, incompatible, active, or CP-1 separate prior public failure.

**Remaining findings:** N-1 through N-7 are non-blocking route-clarity or OPEN-scope findings. They do not reproduce R-5 or R-6.

**P17 remains UNAPPROVED.** This review is not an approval and does not authorize a correction pass. **C1 implementation remains UNAUTHORIZED.** Contract rev3 remains frozen.

A later correction pass, if authorized, should address the non-blocking findings rather than reopen the closed R-5/R-6 routes. No source, test, build, benchmark, runtime, contract, register, authority, or target edit was made by this review.
