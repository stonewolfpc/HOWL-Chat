# C1 / P17 — Final Read-Only Consistency Review

**Date:** 2026-11-09

**Status:** Read-only review. **PASS — no blocking findings.** This review does not approve P17, revise P17, amend contract rev3, update the decision register, or authorize implementation. U1–U8 remain OPEN.

**Reviewed artifact:** `docs/review/PLAN-howl-tts-stage3-c1-p17-revision-2026-11-09.md`

- SHA-256: `56c181db3e494cd030ac6784ef7c6b803b7aac7d99a175d570d3ffd1e7a21fb7`
- Lines: 374

**Baseline checked against:**

- Decision register: `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`, 434 lines; OD-1, OD-3, and OD-9 inspected directly.
- P17: `01d887c1…`, 166 lines, unapproved.
- Contract rev3: `ed0da1a7…`, 183 lines, frozen and not accepted.
- Reconciliation matrix: `3a06f8e3…`, 281 lines; C-01 through C-10 are the controlling conflicts.

## 1. Scope and method

This review checked the revised plan against the accepted owner decisions and the reconciliation history. It verified the local file, hash, line count, preservation inventory, unresolved-case treatment, the four corrected conflicts, the OD-9 consistency note, and the ten-conflict coverage.

The review is document-level. It does not re-audit the source claims behind E1–E9, execute U1–U8, select technical mechanisms, or verify any implementation.

## 2. Required checks

### 2.1 E1–E9 preservation — PASS

The plan now preserves all nine evidence-loss findings:

- E1–E6 at lines 224–229.
- E7 at line 230: late terminal evidence can be dropped by dispatch (`client.go:339–342`).
- E8 at line 231: cancellation and cleanup paths ignore errors from `Cancel` and `Terminate`.
- E9 at line 232: `Shutdown` calls `Terminate`, then clears lifecycle and transport state without establishing successful process reaping.

The preservation requirement is carried into the revision order at line 310 (`retain E1–E9`) and into verification criterion 13 at line 339 (`E1–E9 remain present and unresolved where applicable`). No residual `E1–E6`-only limitation was found.

**Result:** E1–E9 are consistently protected; the earlier omission is corrected.

### 2.2 U1–U8 openness — PASS

All eight cases remain listed at lines 238–245. The plan states that no U case may be closed by the revision plan (line 247), that U1–U8 remain OPEN in the status and OD-9 boundaries (lines 5, 46, and 179), and that they remain OPEN or have separately justified dispositions at criterion 14 (line 340). The conclusion repeats that they remain open (line 366).

No U case is assigned a determined outcome, and no precedence rule is used to close one. The plan’s references to separately justified dispositions preserve OD-9’s escape path without asserting that any disposition exists.

**Result:** U1–U8 remain open everywhere; no case is silently resolved.

### 2.3 C-01 and C-03 public-state boundary — PASS

**C-01** (lines 71–74) now requires only that P17 stop projecting fatal protocol/transport observations to `crashed` without an observed exit, retain the evidence, avoid inferring process exit, and avoid an unsupported terminal or public projection. It states that the public-state treatment when exit is not established remains explicitly unresolved (line 71) and that the prior state’s continued validity is itself an open technical question (line 73). It no longer prescribes “preserve the last independently justified state.”

**C-03** (line 93) permits OD-4’s accepted public fallback for a confirmed exit without orderly shutdown. For the distinct no-established-exit case, it does not invent a fallback: it prohibits assigning a public terminal and refers the treatment back to C-01 as unresolved. The earlier “preserve the last justified state” formulation is gone.

**Result:** Neither C-01 nor C-03 prescribes an unsupported public-state fallback. The only public fallback retained is OD-4’s accepted rule for a confirmed exit without orderly shutdown.

### 2.4 C-06 and OD-1 — PASS

C-06 (lines 126–129) now states that diagnostic and projection roles are separate but not mutually exclusive:

- evidence may be retained diagnostically;
- evidence may support a provider-state projection only under an explicit, justified publication rule;
- a diagnostic label or diagnostic-only record does not select a provider state by itself;
- the underlying evidence remains retained.

This matches OD-1’s accepted “projection-supporting as well as diagnostic” authority. It preserves the S4 contradiction as a correction target without turning “diagnostic” into a prohibition on projection. The verification language at line 129 tests the publication rule rather than requiring every diagnostic observation to be non-projecting.

**Result:** C-06 accurately reflects OD-1 and removes the over-rigid reading.

### 2.5 C-08 and OD-3 — PASS

C-08 (lines 148–151) requires one immutable first terminal per process instance, one linearization point, retention of later evidence, and protection of a replacement generation from old-instance evidence. It explicitly limits immutability to the terminal decision and states that it does not decide whether later evidence may affect diagnostics or non-terminal projections (line 148). The technical-decision field repeats that those effects remain unresolved (line 150).

This does not overstate OD-3. OD-3 accepts the immutable terminal and retained later evidence, while expressly excluding “whether and how later evidence may affect diagnostics or non-terminal projections.” C-08 leaves that exclusion intact and verification criterion 10 (line 336) expressly avoids claiming all later evidence is inert.

**Result:** C-08 preserves terminal immutability without deciding OD-3’s excluded questions.

### 2.6 G-01 and the OD-9 register boundary — PASS

G-01 (lines 175–181) identifies the historical OD-9 condition tied to OD-1 being deferred and notes that OD-1 was later accepted. It requires a read-only consistency review, prohibits rewriting the register as part of the P17 revision, and requires separate authorization for any correction to OD-9’s historical wording.

The prohibition is reinforced at line 261 (Phase 0: flag, do not rewrite), line 216 (rewriting the register as part of the revision is out of scope), and criterion 19 (line 345: no register edit is bundled into the P17 revision). The register’s current OD-9 entry already records that OD-1’s acceptance removed the OD-1 condition but not the U1–U8 condition.

**Result:** G-01 flags the historical issue without authorizing a register edit.

### 2.7 Ten-conflict coverage and no new policy — PASS

The plan contains exactly ten conflict sections, C-01 through C-10, at lines 65–174. Their subjects match the matrix’s direct-conflict register:

- C-01: fatal protocol/transport to `crashed` without observed exit.
- C-02: nonzero exit while stopping to `absent`.
- C-03: public `outcome_unknown`.
- C-04: handshake missing-capability split.
- C-05: public `quarantined`.
- C-06: diagnostics selecting projections.
- C-07: missing generation binding.
- C-08: missing terminal latch/late-evidence rule.
- C-09: unsupported “known action cause.”
- C-10: undefined `seq` ordering.

G-01 is a separate governance consistency check, not an eleventh conflict and not a new P17 policy. The plan does not add a provider state, rename the thirteen-value set, make `quarantined` or `terminated_by_action` public, authorize a contract amendment, or approve P17.

**Result:** All ten conflicts remain covered; no new policy or public state was introduced.

## 3. Residual observations

No blocking findings.

Two non-blocking observations:

1. The plan’s line 370 statement that “No artifact other than this plan was modified” describes the plan revision itself. This assessment adds a review record but does not modify the plan.
2. The review is document-level. It verifies consistency with the cited decisions and matrix, not the underlying source audit or the truth of E7–E9’s source locations.

## 4. Verdict

**PASS.** The revised plan is consistent with the accepted decisions on the seven requested points. It is suitable as the basis for a separately authorized P17 revision proposal.

This verdict does **not** approve P17, close U1–U8, authorize a P17 edit, amend contract rev3, or authorize source changes, tests, builds, benchmarks, runtime experiments, or C1 implementation.

## 5. Verification and limitations

The reviewed plan remains at SHA-256 `56c181db3e494cd030ac6784ef7c6b803b7aac7d99a175d570d3ffd1e7a21fb7`, 374 lines. The decision register remains `3c64797b…`, 434 lines; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; the reconciliation matrix remains `3a06f8e3…`. No P17, contract, register, source, test, configuration, or generated file was changed.

**Limitations:** this review did not re-run the source audit, execute U1–U8, observe runtime interleavings, decide any technical mechanism, or verify implementation behavior.

**No implementation was performed or authorized.**
