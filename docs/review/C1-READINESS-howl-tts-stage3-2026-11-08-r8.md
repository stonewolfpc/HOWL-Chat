# C1 Readiness — Revision 8

**Date:** 2026-11-08
**Subject:** Howl TTS Stage 3 — Capability C1 (helper lifecycle supervisor contract)
**Status:** **NOT ACCEPTANCE-READY — C1 UNAUTHORIZED — NO IMPLEMENTATION**

## 1. Purpose and scope

This document states, without rounding up, how far C1 has progressed after design revision 8. It
separates four things that are routinely conflated: design content, author reconciliation,
independent review, and acceptance. Only the first is complete.

Scope of this revision: documentation only. No source, test, build, runtime, benchmark, or
deployment action was performed.

## 2. Artifact identity

| Artifact | SHA-256 | Lines / bytes |
|---|---|---|
| **C1 design revision 8** | `a8c40fac7d86f3cd34ad34e7f656668684830e76dfb605ae6817fa7119638e06` | 391 / 30,637 |
| C1 design revision 7 (normative for unchanged sections) | `0867acb7216a0e8ee2802b6f08275d997e4138f089ae51acac373c10eea4af54` | 1069 / 161,934 |
| Independent r7 review (the findings r8 answers) | `c025ba0a32c776996fd90a59cca0635da5532a43200a5f6b8dad39ff91730f9d` | 200 / 13,835 |
| Governing helper contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2586 / 182,904 |
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 |
| Owner decisions (D10 complete; six pending) | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 |

The r8 design is a delta document: changed sections are restated in full, unchanged sections are
incorporated by explicit reference to the r7 hash above (r8 §12). The normative corpus is closed and
unambiguous; it is not a single file.

## 3. Readiness verdict

| Dimension | State | Class |
|---|---|---|
| Design content for R-24 … R-31 (r7-closed) | retained, re-verified for absence of regression | `[Observed]` |
| Design content for R-32 … R-37 | closed in r8 §4–§8 | `[Source-derived]` |
| Findings register completeness | 37 identifiers, each with a disposition | `[Calculated]` |
| Acceptance criteria | 33 criteria, each with one verification method (§8.1 map) | `[Calculated]` |
| Criterion execution | **0 of 33 executed**; all tests are specifications | `[Observed]` |
| Independent review of r8 | **not performed** | `[Observed]` |
| Owner acceptance | **not granted** | `[Observed]` |
| Implementation | **none**; C1 surfaces do not exist in source | `[Observed]` |
| C1-AC-15 resource determinations | **open** — 13 rows need written determinations, 8 I/O categories need reviewed N/A entries | `[Unestablished]` |
| `MaxPending=32` | **open** — no accepted authorised source; R8-AC-4 disposition not recorded | `[Unestablished]` |
| HW-1.6 independent review | **open** | `[Unestablished]` |

**Verdict:** the design is ready for a *third* independent review. It is **not** ready for acceptance
and **not** ready for implementation.

## 4. Evidence ladder

| Rung | Meaning | r6 | r7 | r8 |
|---|---|---|---|---|
| 1. Design content | normative text exists and is internally consistent | reached | reached | reached |
| 2. Author reconciliation | author has checked findings against governing text | partial | reached | reached |
| 3. Independent review | a review not written by the design author passes the revision | — | **FAIL — narrow** | **No** |
| 4. Acceptance | owner accepts the capability contract | No | No | **No** |
| 5. Implementation | C1 code exists, reviewed, and tested | No | No | **No** |

Reaching rung 3 requires a reviewer to accept r8; reaching rung 4 requires the owner; reaching rung 5
requires explicit authorization. Neither r3 nor r4 is claimed here.

## 5. What r8 changes, and why it is reviewable

1. **Update authority is enforced, not asserted** (r8 §4). The caller-constructible decision is gone.
   A reset requires a ticket whose 32-byte nonce exists in the issuing authority's registry, consumed
   once, against a binding the supervisor builds from its own state. Approval is the absence of an
   error. The design names its own limits (`unsafe`/`reflect`, and the un-audited C6/C7 implementation)
   instead of claiming more.
2. **Selection ownership has a carrier** (r8 §5). The "direct Catalog only" prose is replaced by a
   run token that only the in-flight direct Catalog operation can mint, closed on return.
3. **R-34 needs no owner decision** (r8 §6). The complete A2-D5 clause conditions case (a) on
   configuration validity, not liveness; the liveness clause belongs to case (b) only.
4. **Budget exhaustion is exercised** (r8 §6.3, C1-AC-33). The definition is required by A2-D5 (b)
   and now has a criterion and a test rather than being removed to tidy a count.
5. **Baseline labels are corrected** (r8 §7) without rewriting what r6 and r7 concluded.

## 6. Explicitly unresolved — not closed by r8

- **C1-AC-15** — resource and I/O determinations are still missing. Closing R-33 … R-37 does not
  close this.
- **`MaxPending=32`** — still without an accepted authorised source. Unchanged by this revision.
- **HW-1.6**, and the independent review and owner acceptance rungs above.
- Every criterion's executed evidence.

## 7. What would move this forward

For the reviewer: re-test r8 §4 and §5 by trying to forge a reset without an issued ticket and to
emit a selection from the adapter path. Those two attempts are the whole of R-33 and R-35.

For the owner: C1-AC-15 determinations, the `MaxPending` decision, and HW-1.6. None of these are
authorised by this document and none should be inferred from it.

## 8. What this document does not claim

It does not claim that any test passes, that the design is correct, that the contract is satisfied,
that C1-AC-15 or `MaxPending=32` is resolved, or that implementation may begin. A design claim is not a
demonstrated property, and no property here has been demonstrated by execution.

## 9. Authorization status

**C1 remains unauthorized.** No implementation, dependency, configuration, contract, or plan change
was made. No source or test file was modified. The pre-existing dirty working tree was left as found.
