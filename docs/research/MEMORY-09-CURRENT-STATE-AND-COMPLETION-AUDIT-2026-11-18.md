# MEMORY-09 — Current State and Completion Audit

**Status: AUDIT RECORD — evidence-based status only.**

This document creates no decision, authorization, gate, finding, acceptance, or registration. It
does not modify any existing file, does not implement, does not execute, and does not resolve any
contradiction it reports. Where two records disagree, both are preserved and the disagreement is
flagged rather than settled.

| Field | Value |
|---|---|
| Date | 2026-11-18 |
| Audit type | Read-only repository inspection and document comparison |
| Authority created | **None** |
| Files created | **1** — this document |
| Files modified | **0** |
| Tests / validators / probes / model calls executed | **0** |
| Servers started | **0** |
| Git actions | **None** — no staging, commit, push, or history change |
| Registry | **Not registered** (by instruction) |
| Scope | `docs/`, `docs/research/`, `docs/review/`, `docs/architecture/`, `eval/mem1/`, `eval/harness/`, `eval/s3/`, `eval/runs/`, root `STAGE_STATUS.md`, Git working-tree state |

---

## 0. Evidence key

Every substantive claim carries one of these tags.

| Tag | Meaning |
|---|---|
| **OBSERVED** | Verified directly in the current repository during this audit |
| **SOURCE-DERIVED** | Explicitly stated by an authoritative document, cited by path and section |
| **CALCULATED** | Derived arithmetically from verified counts; the arithmetic is shown |
| **INFERRED** | A reasoned interpretation; the supporting evidence is given |
| **UNESTABLISHED** | Evidence is absent or insufficient |
| **SUPERSEDED** | An older status is no longer current; the superseding evidence is named |

---

## 1. Executive status: where MEMORY-09 actually stands

MEMORY-09 is **design-complete, partially implemented, and unexecuted**. There is no point in its
history at which a model was used to test the memory system.

| Question | Answer | Tag |
|---|---|---|
| Is the architecture decided? | **Yes.** 57 rules ratified (O-6b), 118 conformance requirements defined, r32 design contract frozen and intact | SOURCE-DERIVED |
| Is it specified? | **Yes, with one blocked decision.** D3 is specification-complete but corpus-blocked; D9's text is unrecoverable | SOURCE-DERIVED |
| Is it implemented? | **Yes, in part.** A 6-file, 1,951-line engine package is committed at `eval/mem1/m09/` | OBSERVED |
| Is the implementation reviewed? | **Yes, for earlier states only.** R1 covers 886 lines; the Gate 2 report covers 1,632; the supplement covers 1,803. **No document covers the committed 1,951-line state** | OBSERVED |
| Has anything been executed? | **No conformance probe, no suite, no subject run.** `py_compile` is the only reported execution anywhere in the implementation trail | SOURCE-DERIVED |
| Has an experiment run? | **S1: 5 behavioural calls** (a slot/threshold study, not a memory test). **S2: never run. S3: closed** (reranker not shipped; Gate B 6/7) | SOURCE-DERIVED |
| Can local-model conversations test the memory system today? | **No.** Nine independent prerequisites are unmet; two are structural rather than procedural | OBSERVED / INFERRED |
| Is there duplicated work? | **No exact duplication found. Three substantial overlaps, one of which the recent work never examined** | See §8–§9 |

**The single most important finding of this audit:** the implemented engine (`eval/mem1/m09/`) and
the experimental conversation harness (`eval/mem1/conversation.py`) are **not connected**. Nothing
outside `m09/` references `m09` at all, and `conversation.py` contains zero references to it
[OBSERVED: `grep -rn "m09" eval/ --include=*.py`, excluding `eval/mem1/m09/` → no matches;
`grep -c "m09" eval/mem1/conversation.py` → 0]. The experiment measures the harness's own memory
model in `eval/mem1/model.py`; the ratified v9 semantics live in a separate package with no caller.
"Point local models at the memory system" is therefore not a configuration step — it is a missing
integration, and no authorization currently permits building it.

---

## 2. Terminology crosswalk

The project runs **four** independent numbering schemes over three different subjects. They are
never interchangeable and must not be read as stages of one sequence.

### 2.1 The four schemes

| Scheme | Labels | Subject | Governing source |
|---|---|---|---|
| **A. Experimental stages** | S1, S2, S3 | The *experiment programme*: slot qualification, injection design, embedder/reranker evaluation | `MEMORY-09-CHECKPOINT-2026-11-15.md`; `MEMORY-09-S2-MEMORY-INJECTION-DESIGN-2026-11-15.md` §12 |
| **B. MEMORY-03 gates** | G0 … G5 | Research gates: protocol design, corpus, manifest/harness, runtime execution, human review, publication | `MEMORY-03` §9, quoted verbatim in `MEMORY-09-S2-MEMORY-INJECTION-DESIGN-2026-11-15.md` §12 |
| **C. November implementation gates** | Gate 1, Gate 2 | Conformance probes (Gate 1) and `eval/mem1` implementation (Gate 2) | `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` §"Gate 1" / §"Gate 2" |
| **D. Engine-internal labels** | "G1", "G4" inside `eval/mem1` | The **harness's own** acceptance generations and its execution gate | `eval/mem1/write-guard-v2/g1-receipt-writeguard-v2.json` (`"gate": "G1"`); `conversation.py:177` |

### 2.2 The collisions, stated plainly

| Collision | Consequence |
|---|---|
| **"G1"** = *corpus and labels* (B) vs **"Gate 1"** = *conformance probes* (C) vs **"G1"** = *the harness acceptance generation* (D) | Three meanings. "G1 passed" means three different things |
| **"G2"** = *model/harness manifest* (B) vs **"Gate 2"** = *implementation authorization* (C) | These are unrelated |
| **"G4"** = *human review* (B) vs **"G4 authorization"** = *network execution permission* (`conversation.py:177`) | The source refuses a network responder with "requires G4 authorization", but MEMORY-03 assigns runtime execution to **G3** and human review to **G4** |
| **"G5"** = *publication* (B) vs "G5 execution" / "G5 rubric" in the D1–D10 record | The D1–D10 record's G5 is a scored run with a rubric, not a publication gate |

### 2.3 What the owner's "Stage 1 and Stage 2 implementation are complete" resolves to

Two defensible readings; both are given because the repository does not contain the phrase.

| Reading | Corresponds to | State |
|---|---|---|
| **Implementation sequence** = Gate 1 then Gate 2 | Gate 1: conformance probes permitted (O-7 + Amendment 1) · Gate 2: implementation authorized and delivered | **Complete for Gate 2's stated scope** — with the coverage caveat in §4.4 |
| **Implementation passes** = first pass, then correction | R1 (886 lines) → Gate 2 report (1,632) → Supplement/correction (1,803) → committed tree (1,951) | **Complete as passes, but the final pass is unreported** |

What it does **not** mean, and cannot: that S1 finished (it stopped at 5 of 18), that S2 ran (it did
not), or that anything was tested.

**Recommended convention (not applied):** write `Gate 1`/`Gate 2` for scheme C, `G0–G5` for scheme
B, `S1–S3` for scheme A, and `harness-G1` / `exec-gate` for scheme D. No existing document was
edited to do this.

---

## 3. Architecture and design completion

| Artifact | Identity | State | Tag |
|---|---|---|---|
| Architecture v9 — 57 rules | `docs/research/MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | **RATIFIED as design baseline** (O-6b), tier 3 under O-6a. Header still lists Gate 2 as "unanswered" — **stale**, flagged in R1 and left unfixed | SOURCE-DERIVED + SUPERSEDED |
| Conformance map — 118 requirements | `docs/research/MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | Requirements ratified; "0 implemented" line now **stale** (see §4) | SOURCE-DERIVED + SUPERSEDED |
| Stage contract r32 | 91,311 B, `e88ef2d1…` | **Frozen, intact, unmodified** | SOURCE-DERIVED |
| Owner decisions (O-series) | `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | All confirmed except **O-2** (parked) | SOURCE-DERIVED |
| D1–D10 | `MEMORY-09-G5-SEQUENCING-D1-D10-DECISION-RECORD-2026-11-18.md` | Closed; D9 text unavailable | SOURCE-DERIVED |
| Preregistration + 3 amendments | 4 frozen documents, hashes verified inside the v2 receipt | Frozen | OBSERVED (via receipt) |
| Experiment plan v3 | 21,616 B, `6c4078bd…` | Frozen; **approval status not established in this audit** | SOURCE-DERIVED / UNESTABLISHED |

Requirement counts by family, as recorded on the map:

`A 5 · B 11 · C 9 · D 10 · E 5 · F 9 · G 3 · H 10 · I 3 · J 16 · K 27 · L 10`
**CALCULATED:** 5+11+9+10+5+9+3+10+3+16+27+10 = **118** ✓

**Design is not the bottleneck.** The bottleneck is that design decisions have not been executed
against anything, and the corpus that would let them be executed does not exist.

---

## 4. Implementation completion and source evidence

### 4.1 What exists — OBSERVED

`eval/mem1/m09/` — 6 modules, **1,951 lines**, tracked, clean against HEAD.

| File | Lines | sha256 (verified this audit) |
|---|---:|---|
| `__init__.py` | 38 | `2b72a75f870b1ea370806691ec8ca17388b70d5eca8d0a0837b41ed704444cad` |
| `authority.py` | 236 | `f344582350af88f123aadf7540752ec7079e7041d5875ef0b4479d597e7824fc` |
| `engine.py` | 1,209 | `7505e9acca8e32293d13b73fdb1e5f09e972a8bf6ebce439b8b98cef34d45f22` |
| `errors.py` | 115 | `5f33e710359b690862878720b6a1e16602bf0d81953e2dbb0625d413edac9c78` |
| `records.py` | 252 | `ad2d4d11501c357b6a96bcf053be47d730c7ac99a945ae03105f684db7dcc2a5` |
| `retrieval.py` | 101 | `cd2d380bd5d9d1f731ed385dd8a19405307196eaf126fd907ed88ac660a83286` |
| **Total** | **1,951** | — |

Committed as `830cba4 feat(memory-09): add the complete m09 authority engine package`
(`git show --stat`: 6 files, 1,951 insertions). Working tree matches HEAD for this path
[OBSERVED: `git diff --stat HEAD -- eval/mem1/m09` → empty].

### 4.2 The line-count trail — four states, three reports

| State | Lines | Evidence |
|---|---:|---|
| R1 review | **886** | R1 §1 inventory: 31 + 74 + 154 + 156 + 391 + 80 |
| Gate 2 report | **1,632** | Report §2: 35 + 156 + 1,045 + 103 + 192 + 101 |
| Correction supplement | **1,803** | Supplement §2: 35 + 236 + 1,118 + 103 + 210 + 101 |
| **Committed tree** | **1,951** | OBSERVED `wc -l` |

**CALCULATED delta from the last covering report:** `1,951 − 1,803 = 148` lines, distributed as
`engine +91`, `records +42`, `errors +12`, `__init__ +3`, `authority ±0`, `retrieval ±0`.

**No document in the repository accounts for those 148 lines.** The supplement says it resolves the
five withheld items; R2 records further corrections (including deletion of a dead branch, which
would *reduce* lines). Neither reaches 1,951. This is the audit's clearest evidence gap: the
artifact that exists is **newer than every report about it**.

### 4.3 What the implementation actually contains — SOURCE-DERIVED

From R1/R2 and the Gate 2 supplement, and consistent with the module inventory:

- Four record classes with structural validation; containment, ownership, domain and truth-aptness
  kept as separate fields.
- Claims, accounts, beliefs, promotion, disclosure, lifecycle, declarations (`engine.py`).
- A proposition registry, conflict groups, assertion of incompatibility, capture/assembly, gated
  `authoritative`, version bracket (supplement §2).
- Authority as `(kind, id)` with `Grant.operations` and **`require(domain, actor, operation)` with
  no default** — so a promotion grant cannot silently authorise authorship (supplement §2.1).
- `author_ambient` implemented under the owner's correction.
- Fail-closed retrieval contexts; conflict presentation (`retrieval.py`).
- `EXPOSES_TRANSFER = False`; no transfer, share or carryover symbol exists.

### 4.4 Coverage against the 118 requirements — partial, and caveated

| Capacity | State | Tag |
|---|---|---|
| Rules mapped to source (57) | Mapped in R1 §2 for the 886-line state; two **Partial** (rule 13 conflict detection, rule 39 ambient authorship — since implemented) | SOURCE-DERIVED |
| Requirements mapped to evidence (118) | Gate 2 report §7 maps them; performed against the 1,632-line state | SOURCE-DERIVED |
| Requirements with **executed** evidence | **0** | SOURCE-DERIVED |
| Requirements with **runtime** evidence | **0** | SOURCE-DERIVED |

The 118-requirement mapping has not been re-performed for the committed tree. Whether the 148
unreported lines changed any requirement's status is **UNESTABLISHED**.

### 4.5 The reported open items — verified as *reported*, not as *current*

| Item | Report says | Audit finding | Tag |
|---|---|---|---|
| Six modules | 6 | **Confirmed** | OBSERVED |
| 886 lines | 886 | **Superseded** — 1,951 | OBSERVED |
| `author_ambient` / `Grant.operations` blocked | Blocked in Gate 2 report §5 | **Resolved by the supplement**; `authority.py` now 236 lines with `OPERATIONS` and `require(..., operation)`. Report §5 is void by its own supplement | SOURCE-DERIVED |
| Ten assumptions, two flagged | Report §8 | Still present as a record; **not re-verified** against the committed tree | SOURCE-DERIVED / UNESTABLISHED |
| Nine residual risks | Report §9 | Includes concurrency **not demonstrated**, `reason` vs `code` (since changed to `code` per supplement), and `proposition_unregistered` as the only invented diagnostic | SOURCE-DERIVED |
| Concurrency undemonstrated | §9 item 8 | **Still true.** `_capture`'s version test is sound only in a single-threaded store; Class B unreachable | SOURCE-DERIVED |

### 4.6 Authorization: was the implementation authorized?

**Yes.** The exact owner record is
`docs/research/MEMORY-09-OWNER-DECISIONS-2026-11-17.md`, section **"Gate 2 — implementation of
`eval/mem1` (AUTHORISED — implementation only; testing not authorised)"**:

> **Decision (owner, 2026-11-17): Gate 2 is authorised — implementation only.** The sequence is
> **implementation first, architecture review second, testing not at all** until the owner
> explicitly approves the next stage.

Scope and exclusions, verbatim from the same section: `eval/mem1` only · no production Howl
integration · no changes to r32, plan v3, Amendment 1 or ratified decision records · no new
architecture, rules, fields, requirements or probes · no transfer/share/carryover interface · do not
resolve open questions through implementation assumptions · **do not run conformance probes or other
tests under this authorisation**.

**Discrepancy resolved, not hidden.** The claim that "the ledger records Gate 2 as unauthorized" is
**SUPERSEDED**: the same document preserves the earlier "not authorised" paragraph immediately above
the decision, explicitly marked *"the position **before** the Gate 2 decision; preserved unamended."*
Four other documents still carry the stale "Gate 2 unanswered" line (§14.1).

**Two limits that survive:**

1. Gate 2 authorized **writing code**, not executing it. Testing is separately unauthorised.
2. The 148 unreported lines were committed under `830cba4`. Whether that commit falls inside Gate 2's
   scope cannot be settled from the records; the correction supplement shows the owner twice
   narrowed and then widened the instruction, so scope was moving. Recorded as **UNESTABLISHED**,
   not as a violation.

---

## 5. Conformance and validation evidence

### 5.1 What has actually been validated — OBSERVED

| Evidence | Result | What it proves | What it does not |
|---|---|---|---|
| `write-guard-v2` receipt — `status: PASS` | **142/142 checks**, 18 components, 33 negative checks, 5 limitations, 11 dependencies | Static, zero-model, in-process self-check of the G1 harness | Nothing about memory behaviour |
| Protected set in that receipt | **12 files** (`stage1: 8`, `accepted_v2: 4`), guard `protected_files: 12`, checks `WG01–WG24` | Identity protection of 12 artifacts under 24 checks | Nothing about the engine |
| `g1-receipt.json` | **96/96 checks** PASS | Earlier G1 acceptance | Superseded by v2 for the check count |
| `G4 preflight receipt` | 13 probe receipts; T01–T07+ PASS | Preflight of probe accounting | **Not a G4 run.** Receipt itself records `g4_authorized: false`, `g4_executed: false` |
| Frozen identities | 8 artifacts + 4 preregistration docs + 2 preserved docs, all `match: true` | Integrity of frozen inputs at that moment | Not re-verified in this audit |

### 5.2 The 142 / 38 / 24 / 12 relationships — evidence for Finding 6 and Finding 5

- **142** = the `write-guard-v2` **check** count [OBSERVED in the receipt].
- **24** = checks `WG01–WG24`, i.e. **24 checks over 12 artifacts** [OBSERVED; receipt evidence field
  reads *"all twelve reference artifacts identity-protected; G1 zero-model check set including
  WG01-WG24"*].
- **12** = protected **artifacts**: 8 stage-1 + 4 accepted-v2 [OBSERVED].
- **38** = validators named by D6 for the v3 generation [SOURCE-DERIVED: canonical §4 D6].

**INFERRED, offered as evidence and not as a ruling:** the apparent "12 vs 24" conflict in Finding 5
may be a **unit mismatch** — 12 artifacts against 24 checks over them — rather than two different
protected sets. If so, Finding 5's premise ("if both numbers count protected artifacts") does not
hold. This is **not** a resolution: D6's v3 set has never been enumerated and compared, and only the
owner may close the finding. Similarly, 38 validators and 142 checks are different quantities, so
Finding 6's "subset, replacement or addition" question remains genuinely open.

### 5.3 What "0 implemented / 0 executed" now means

The conformance map's status line — *"0 implemented · 0 registered in any framework · 0 executed"* —
was accurate when written. After Gate 2 it is **SUPERSEDED in one respect and still true in two**:
implementation exists in `eval/mem1/m09/`; **0 requirements have been executed, and none are
registered in any test framework** [OBSERVED: no test runner, no suite, no probe invocation].

---

## 6. Gate authorization and completion history

### 6.1 MEMORY-03 gates (scheme B) — as recorded

| Gate | Scope | Recorded status | Evidence |
|---|---|---|---|
| **G0** | Protocol design and threshold review | Design complete; freeze needs owner approval. **G0 receipt exists** and states a pass authorizes nothing beyond scope | `MEMORY-09-G0-PREFLIGHT-2026-11-15.md`; hash `5c245283…` verified in the v2 receipt |
| **G1** | Corpus, labels, fixtures | **Not complete.** No corpus with the 16 frozen conversation specifications exists | Canonical §4 D3 ("Corpus-blocked"); §12 below |
| **G2** | Model/hardware manifest and harness build | Harness built and self-checked (96/96 → 142/142). Manifest binding for an actual run not established | §5.1; `G12` still open |
| **G3** | **Runtime execution of any stage** | **Never entered.** Not authorized, not refused, not requested | S2 design §12; `conversation.py:177` refuses network responders |
| **G4** | Human review and analysis | **Never entered.** A *preflight* receipt exists; that is not a G4 run | `eval/mem1/g4-preflight-receipt.json` (`g4_executed: false`) |
| **G5** | Publication | **Closed.** No publication authorization exists for any result, and U4-ii(a) requires it to predate the first envelope | Canonical §4 D8-U4-ii(a); §8 |

### 6.2 November implementation gates (scheme C) — chronological

| Step | Date | Authority | Scope | Result |
|---|---|---|---|---|
| **Gate 1** | 2026-11-17 | Owner, **O-7 + Amendment 1** (`6afd2a0d…`) | *Prospectively permits only* conformance probes for the 118 requirements | Permitted; **probes were never written or run** |
| **Gate 2 — first instruction** | 2026-11-17 | Owner | Implementation, `authority.py` **frozen** | Delivered; §5 blocked |
| **Gate 2 — report** | 2026-11-17 | — | Report of the above | Report §5 "blocked item" |
| **Gate 2 — correction** | 2026-11-17 | Owner: the prohibition *"was too strong"* | `authority.py` may change; `author_ambient` required | Supplement; **report §5 void** |
| **R2 review** | 2026-11-17 | Gate 2 | Focused classification of five owner findings | A/B/C classifications; two **C** (architecture change required) |
| **Commit `830cba4`** | undated in record | — | The full 1,951-line package | Committed; **no covering report** |

### 6.3 What no gate authorized

| Action | Status |
|---|---|
| Running any conformance probe | **Not authorized** (Gate 2 prohibition, verbatim) |
| Running the self-test suite | **Not authorized** |
| Compiling / linting as verification | **Not authorized** (report §10.3 notes `py_compile` was deliberately not run) |
| Network responder / model call | **Refused in code** (`ExecutionNotAuthorised`) |
| Corpus preparation | **Not authorized** |
| Publication | **Not authorized** |

---

## 7. Experimental history: S1, S2, S3

### 7.1 S1 — model-slot qualification and behavioural calls

| Property | Value | Tag |
|---|---|---|
| Slots qualified | **A** (5,379 MiB, Δ−55) · **B** (4,041 MiB, Δ+12) · **C** (6,461 MiB, Δ−24, **identity caveat**: `f820ef3f…` has no external reference; nearest published file differs by 992 B) | SOURCE-DERIVED |
| Calls filed | **5 of 18** planned; hard ceiling 30; **25 unspent** | SOURCE-DERIVED |
| Content-admissible | **2** — both condition A, both `stop` (64 and 19 completion tokens) | SOURCE-DERIVED |
| R0 (256 tokens) | 4 observations, 2 admissible | SOURCE-DERIVED |
| R1 (1,024 tokens) | 1 observation, **0 admissible** | SOURCE-DERIVED |
| Condition A | 3 calls: r1 `stop` 64, r2 `stop` 19, r3 `length` 256 | SOURCE-DERIVED |
| Condition B | 2 calls, both `length` (256, then 1,024) — **no answer ever produced** | SOURCE-DERIVED |
| Amendment 01 | Activated; overlay `IN FORCE` at `8d05620a…`; call 05 is the first R1 call | SOURCE-DERIVED |
| Scoring protocol | `MEMORY-09-S1-SCORING-v1`, frozen `a7bd7208…`; six axes; `length` never scored | SOURCE-DERIVED |
| Validation | `validate` 89/89; corpus PASS; batch audit PASS on 5 calls | SOURCE-DERIVED |
| Pause point | After call 05. Three options tabled (keep 1,024; change surface; accept `T=length` as the finding). **None taken** | SOURCE-DERIVED |

**What S1 actually established:** that at these budgets this model deliberates past termination. It
did **not** test memory — there is no memory layer in the request path (call 05 request body, 342 B,
`max_tokens` + 2 messages only). The sampler question was closed by **reading the executable's
embedded adapters**, not by spending a generation: effective chain `rep_pen 1.0 → top_k 100 → top_p
0.95 → temp 0.8 → dist` for all three slots, derived not observed.

### 7.2 S2 — memory-injection experimental design

| Property | Value | Tag |
|---|---|---|
| State | **Design only.** v1 preserved unmodified (`eval/s2-schedule-draft.json`, `03897f44…`); Rev 2 added | SOURCE-DERIVED |
| Trials | **48** (4 characters × 2 blocks × 6) | SOURCE-DERIVED |
| Validation | **PASS, 0 violations, 48 trials** | SOURCE-DERIVED |
| Foreign-arm coverage | **8 of 12** ordered pairs; balanced (each character donates 2, receives 2); the 4 uncovered are the offset-3 diagonal and are structurally unreachable with two shifts | SOURCE-DERIVED |
| Message-count fork | Parameterised: 4-input (5 total, 43,087 B) vs 5-input (6 total, 49,423 B); both validate; 4-input recommended, owner decides | SOURCE-DERIVED |
| Termination policy | **Proposed, not adopted** — fix `max_tokens` before the block; `length` is terminal, never retried; non-terminating replies are not evidence about memory | SOURCE-DERIVED |
| Ran | **Nothing.** No server, no request, no trial, no arm observed | SOURCE-DERIVED |
| Open owner decisions | **5** — five-message reading · 48 vs smaller first block · `max_tokens` · which gate authorizes execution · S2's own scoring axes | SOURCE-DERIVED |

Two self-corrections worth preserving: the 48 is a **provisional cost estimate, not a decided
sample** (preregistration requires sample size from measured pilot variance), and the six-trial run
was **misnamed a pilot** — it is an **instrument check**.

### 7.3 S3 — closed

| Property | Value | Tag |
|---|---|---|
| Gate A | **14/14 PASS** | SOURCE-DERIVED |
| Gate B | **6/7 PARTIAL** — **P5 ("Who trained the horse?") is an open failure** | SOURCE-DERIVED |
| Reranker | Evaluated, **insufficient, not shipped** | SOURCE-DERIVED |
| Resource envelope | **No configuration meets VRAM ≤384 MiB *and* p50 ≤40 ms** | SOURCE-DERIVED |
| Status | **Closed.** Apparatus frozen; explicit prohibition on re-fetching or re-running | SOURCE-DERIVED |
| Reopened by this audit | **No** | — |

Any plan that treats the baseline as clean must carry P5 as an open failure. The experiment plan
states this explicitly.

---

## 8. Recovery work history and duplication/overlap audit

### 8.1 The three distinct "recovery" subjects

| Subject | What it is | Documents |
|---|---|---|
| **R3 file recovery** | Recovering a lost *revision* (r3) of the stage contract from disk or Git | `MEMORY-09-R3-RECOVERY-CLOSURE-2026-11-17.md` |
| **F3-C storage recovery** | Crash/power-loss recovery, projection repair, recovery gate, incarnation recovery, new-root rules for **production storage** | `ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md`, frozen master, adjudication package |
| **Operational recovery** (Nov 18) | What to do when **chain authority** in the memory engine is contested — owner as external resolver, no in-system adjudicator | The six 2026-11-18 documents (§10) |

These are three different questions sharing one word.

### 8.2 Classification of each recent activity

| Activity | Class | Basis |
|---|---|---|
| R3 recovery closure (11-17) | **1 — entirely new work** | Bounded search of 9,419 files and 4,628 Git objects; conclusive negative |
| D1–D10 decision record (11-18) | **2 — restatement/consolidation** | Canonical transcription of decisions supplied in conversation; explicitly "not proof these decisions were previously integrated into repository files" (§9) |
| Operational-recovery analysis (11-18) | **5 — related work at a different layer** (+ **3 — refinement** of the canonical §5.3.2 analysis) | Builds on Finding 2; adds triggers, preservation, continuity, freshness |
| Owner-decision brief (11-18) | **2 — restatement** of positions already in canonical §7 item 2, plus a new neutral-options framing | §1–§2 restate; options A1/A2/B–G are new presentation |
| Delegated decision statement (11-18) | **2 + new authorization record** | Records choices already available in the brief; the *authorization* is new |
| Procedure specification (11-18) | **1 — new work**, with **5 — parallel** to F3-C at another layer | No prior engine-level recovery procedure exists |
| Muse reviews (11-18, two) | **1 — new work** | Independent reads; no prior review of this specification |
| Two wording corrections | **3 — refinement** | S4 cross-reference; "successful stop" phrasing |

### 8.3 Did the November 16–18 work repeat anything?

**No exact duplication was found.** Specifically:

- It did not re-run R3 recovery (different object).
- It did not re-decide the (b)/(c) deferrals — it restated them and deferred again.
- It did not touch S1/S2/S3.
- It did not restate the F3-C protocol — it never cited it (see §9).

**One real gap in the trail, not a duplication:** the Nov 18 work reconciled E/F against
**D8-U4-iii / §7 item 11** — the MEMORY-09-side rules — but not against F3-C. The canonical record's
D8-U4-iii already permits reporting repair *"only mechanically: locating an expected artifact that
was previously omitted, or establishing that two paths hold identical bytes"*, and the §5.4 caveat
already notes there is **no testable discriminator** for that line. The November 18 analysis arrives
at the same place from the other direction (E: "no testable discriminator"). That is convergence on
an existing owner question, and it is **correctly labelled as unresolved**.

---

## 9. F3-C versus MEMORY-09 recovery relationship

### 9.1 What F3-C already decided — SOURCE-DERIVED

| Provision | Location | Content |
|---|---|---|
| Recovery in frozen scope | `docs/FROZEN-ARTIFACTS.md` row 2 | Frozen master covers "operations; **recovery**; tombstones" |
| OD-02a / OD-02b split | `ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md` §2 | Replacement mechanics (02a) and **durable recovery infrastructure** (02b) are separate instruments; **either alone authorizes nothing**; the joint transaction requires both |
| Recovery components | same, lines ~57–73 | "a read gate; **recovery coordination**; **a new root**" |
| Missing-marker crash recovery | same, §3 | A commit is authoritative over incomplete targets; recovery must reinstall package-derived files and verify the manifest **before ordinary access** |
| No guessed commits | same, line ~193 | "a record without a valid package is a **hard evidence block, not a guessed commit**" |
| No root by implication | `ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md` | "**No new root may be introduced by implication**" |
| Deterministic abandonment | `ADDENDUM-2` §6 | Serialized ordering under the shared lock decides the winner; validated at lease issuance **and** final commit |
| Status of both instruments | `ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md` | **OD-02a: Awaiting owner · OD-02b: Awaiting owner** |

### 9.2 Are they the same question?

**No.** The layers differ:

| | F3-C | MEMORY-09 |
|---|---|---|
| Object | Storage files, projection packages, chat roots, locks | Records, accounts, promotions, evidence envelopes, receipts |
| Failure | Crash mid-write; partial projection; lost marker | Contested chain authority; missing envelope; competing heads |
| Authority | Single-writer lock + recovery gate + incarnation record | Chain-based acceptance; owner as external resolver |
| Escrow | OD-02a/OD-02b owner decisions | §7 item 11 owner question; E and F unresolved |

**Merging them would be wrong,** and this audit does not do it.

### 9.3 Does any F3-C decision already answer E / F / S5 / S7?

| MEMORY-09 question | F3-C analogue | Verdict |
|---|---|---|
| **E** — repairing an omitted artifact when no discriminator exists | "Hard evidence block, not a guessed commit"; reinstall-and-verify before access | **Does not answer.** F3-C's case has a *manifest* to verify against; E's problem is that no discriminator exists. Analogy only |
| **F** — what anchors a new root | "No new root may be introduced by implication"; new roots are a declared recovery component under OD-02b | **Constrains by analogy, does not decide.** It supplies a principle (no implicit root) that the engine question already independently expresses |
| **S5** — stop rather than answer F | n/a | **Does not answer** |
| **S7** — branch completeness | Multi-chat recovery ordering under a shared lock | **Does not answer** |
| **E/F generally** | Separation of recovery authority from replacement authority | **Structural parallel:** F3-C's insistence that recovery is a *separately granted instrument* mirrors D4's "execution authorization is not memory authority". This is the most transferable principle, and no document records anyone noticing it |

**Reconciliation status: UNRECONCILED.** No document states whether F3-C's recovery principles
apply, constrain, or are irrelevant to the memory engine. The relationship is identified here for the
first time and is left open by instruction. Both architectures currently have recovery as an
**open owner question** — F3-C via OD-02a/OD-02b, MEMORY-09 via §7 item 11 — so neither can lend
authority to the other.

---

## 10. The November 18 recovery work — verified state

Six documents, all untracked, none registered, all created 2026-11-18.

| Document | Lines / bytes | sha256 (prefix) | Class | Tag |
|---|---|---|---|---|
| `…OPERATIONAL-RECOVERY-ANALYSIS-2026-11-18.md` | 248 / 13,780 | `e7c82ea1100be217…` | Reviewed analysis; **not owner-approved** | OBSERVED |
| `…OPERATIONAL-RECOVERY-OWNER-DECISION-BRIEF-2026-11-18.md` | 270 / 15,069 | `d010a99870fa99f5…` | Neutral options A1/A2/B–G | OBSERVED |
| `…OPERATIONAL-RECOVERY-DECISION-STATEMENT-2026-11-18.md` | 127 / 6,056 | `f4935d822dcf0727…` | **Delegated decision taken by ChatGPT under owner delegation**; §4 is the only granted authority | OBSERVED |
| `…OPERATIONAL-RECOVERY-PROCEDURE-SPECIFICATION-2026-11-18.md` | 240 / 12,580 | `37e0b8a5373a6adf…` | **DeepSeek-authored specification**; descriptive only | OBSERVED |
| `…OPERATIONAL-RECOVERY-REVIEW-PASS-2-MUSE-2026-11-18.md` | 94 / 7,091 | `6ad9ebc3b2290bf0…` | **Muse** sign-off of Pass 1 corrections; not authored by the assistant; provenance unconfirmed by owner | OBSERVED |
| `…OPERATIONAL-RECOVERY-SPECIFICATION-REVIEW-MUSE-2026-11-18.md` | 283 / 14,248 | `9541289d50717e16…` | **Muse** review of the specification; transcription with attribution block | OBSERVED |

### 10.1 What each is, and is not

| Content | Classification |
|---|---|
| Finding 2 remains `false_closure`, partially resolved, open on operational recovery | **Owner position** (canonical §7 item 2) |
| Recovery limitation acknowledged unresolved (option A1) | **Delegated decision** (ChatGPT, under owner delegation) |
| Written procedure commissioned (option B) | **Delegated decision** |
| Specification authoring and review authorized; implementation, execution, rehearsal, registration **not** | **Owner authorization** (decision statement §4) |
| The two wording corrections applied (S4 cross-ref; "successful stop") | **Reviewer recommendation, owner-accepted, DeepSeek applied** |
| Traversal / stopping / freshness behaviour | **Author specification** — never executed |
| "Review complete, descriptive only" | **Reviewer gate verdict** — not acceptance |

### 10.2 The boundary that survives

Nothing in these six documents demonstrates operational recovery. No procedure was implemented,
tested, rehearsed, or executed. Finding 2's third dimension is as open as it was before the work
began. The specification's own §14 records: Implemented **No**, Executed **No**, Corpus **No**,
Tested **No**, Safety **No**.

---

## 11. D1–D10 status

| ID | Subject | Status | Spec-complete | Implemented | Validated | Experiment evidence | Limitation |
|---|---|---|---|---|---|---|---|
| **D1** | No equivalence conclusion | Closed | Yes | n/a (prohibition) | n/a | None | `g5_verdict` unavailable; no `indistinguishable` category may exist |
| **D2** | Amendment 2 §1 governs; cells never collapsed | Closed | Yes | n/a | n/a | None | `net_only` / `selection_only` never merged |
| **D3** | Evidence architecture and counts | **Specification-complete; corpus-blocked** | Yes | Partial (harness enforces N=16 / 13 / 3) | Static only | None | **No corpus**; no `plan_sha256`; G4 plan-row discrepancy preserved |
| **D4** | Authorization provenance and identity | Closed | Yes | **No** | No | None | **Cannot be instantiated**: no corpus, no G07 tokenizer, no G12 binding, **no grant mechanism** |
| **D5** | Immutable evidence; opaque `run_id` | Closed | Yes | Partial | No | None | `docs/authorizations/` **does not exist** (OBSERVED) |
| **D6** | `write-guard-v3`, all-or-nothing | Closed | Yes | **No** — v3 absent | No | None | v3 targets `MEM1_VERSION 0.2.0` / schema 2; tree has **`0.1.0-g1`** and **no `receipt_schema_version`** (OBSERVED) |
| **D7** | Disposition precedence, cell-scoped statements | Closed (a, b) | Yes | No | No | None | **Finding 7**: vocabulary not formally mapped |
| **D8** | Evidence integrity, invalidity, terminal record | Closed (U1–U5) | Yes | No | No | None | **E/F open**; two silent-preference defects recorded and **not repaired** |
| **D9** | — | Closed | **No — text unrecoverable** | No | No | None | Must not be reconstructed from neighbours |
| **D10** | Report identity, determinism, G-06 | Closed | Yes | Partial | No | None | R2 finding: belief-vs-truth distinction exists only in the promoting actor's prose (**class C**) |

**D3's dependency**, in one line: no corpus → no frozen probe identities → no `plan_sha256` → nothing
under D4 can be issued → no run.

---

## 12. Exact remaining prerequisites for local-model chat testing

The goal — *use local-model conversations to test the memory system* — requires the chain below. Each
row states what exists, what is missing, whether it blocks, and who must authorize it.

| # | Prerequisite | Exists now | Missing | Blocks next step? | Authority required | Tag |
|---|---|---|---|---|---|---|
| **P1** | 16 frozen conversation specifications | **No.** No corpus artifact set found anywhere in the tree | All 16 | **Yes — hard block** | Owner (corpus preparation explicitly outside every existing authorization) | OBSERVED / SOURCE-DERIVED |
| **P2** | 13 probe identities per conversation | Count enforced in code (`PROBES_PER_CONVERSATION = 13`, `probe.py:70`); `PROBE_QUOTAS` defined | The **frozen identities themselves** — they are corpus-dependent | Yes, follows P1 | Owner | OBSERVED |
| **P3** | `plan_sha256` | **No** | Derived from P1 | Yes | — (derived, not authorized) | SOURCE-DERIVED |
| **P4** | 144 session envelopes, 1,872 probe receipts | **0 of each** | All | Yes | — | CALCULATED: 16 × 3 reps × 3 scored arms (A0/A1/D) = **144**; 144 × 13 = **1,872** ✓ |
| **P5** | **G07** tokenizer binding | **No** | Runtime tokenizer bound to a recorded identity | Yes — named G4 blocker | Owner | SOURCE-DERIVED |
| **P6** | **G12** model revision + sampling binding | **No** for a run. S1 *derived* the sampling chain for three slots — related, **not** a G12 binding | Recorded, frozen revision + sampling for the run's model | Yes — named G4 blocker | Owner | SOURCE-DERIVED |
| **P7** | Grant mechanism | **No.** `ExecutionNotAuthorised` is raised **unconditionally** for any network responder | A mechanism by which an authorization could actually be honoured — **a code change** | **Yes — hard block, structural** | Owner (Gate 2 prohibits new scope; a new authorization is needed) | OBSERVED (`conversation.py:177`) |
| **P8** | Live responder / chat bridge | **No.** Only `ScriptedResponder` (`conversation.py:62`, `network = False`) and `P5Responder` (`p5.py:151`) exist | Any network or local-model responder class | **Yes — hard block, structural** | Owner | OBSERVED |
| **P9** | Engine ↔ conversation integration | **No.** Nothing outside `m09/` references `m09`; `conversation.py` has 0 references to it | A bridge, if the engine is what is being tested | **Yes — and undefined by any record** | Owner (architecture question, not implementation) | OBSERVED |
| **P10** | Execution authorization | **No** | The gate the code calls "G4 authorization" (MEMORY-03 would call it G3) | Yes | Owner, separate instrument | SOURCE-DERIVED |
| **P11** | Publication authorization | **No** | A pre-execution publication scope (U4-ii(a)) — **separate from execution** | Not for the run; **yes for reporting any result or refusal** | Owner, separate instrument | SOURCE-DERIVED |
| **P12** | Environment evidence (5 classes) | Partly: harness records provenance | authorized / expected / configured / effective / observed as five distinct claims per session | No — but required evidence | — | SOURCE-DERIVED |

### 12.1 Two blockers are structural, not paperwork

P7, P8 and P9 cannot be discharged by signing a document. **There is no code path by which an
authorized run could occur**, because:

1. the only responders are scripted;
2. the guard refuses network responders before any grant is consulted;
3. the implemented engine has no caller inside the experimental harness.

**INFERRED:** without a deliberate owner decision on *which* memory implementation the experiment is
meant to test — the harness model in `eval/mem1/model.py` (6,152 lines across the harness modules) or
the ratified engine in `eval/mem1/m09/` — the corpus work in P1 could be built against the wrong
subject. This is the most consequential ambiguity in the path.

### 12.2 What does *not* block local-model testing

- Design decisions D1–D10: closed, and none prevents an authorized run.
- Findings 3–8: procedural/reporting questions, not run preconditions.
- S1's truncation result: a property of a different instrument.
- S3's closure: fenced off, and required to stay closed.

---

## 13. Repository and documentary state

| Item | State | Tag |
|---|---|---|
| HEAD | `830cba4 feat(memory-09): add the complete m09 authority engine package` | OBSERVED |
| `m09` vs HEAD | **identical** — no uncommitted changes in the engine | OBSERVED |
| Pre-existing modified files | **6, unchanged, untouched:** `config/logs/startup.log`, `config/model-registry.json`, `docs/DOCUMENT-REGISTRY.md`, `docs/research/OWNER-ENVIRONMENT-HISTORY-2026-11-14.md`, `eval/harness/__pycache__/validate.cpython-312.pyc`, `eval/harness/validate.py` | OBSERVED |
| Untracked files | **168**, of which **89** are `MEMORY-09*` documents under `docs/research/` | OBSERVED |
| Documentation volume | `docs/`: **579** `.md`; `docs/research/`: 288; `docs/review/`: 224; `docs/architecture/`: 17; `MEMORY-09*`: 155; `docs/research/corpus/`: 170 files | CALCULATED (file counts) |
| Registry | Header claims **444 rows**; **506** lines open with a numbered row cell; **579** `.md` files exist under `docs/` | OBSERVED |
| Registry freshness | **Stale.** 0 of the 2026-11-18 documents are registered; `grep -c "2026-11-18"` → 0 | OBSERVED |
| `docs/authorizations/` | **Absent**, although D5 names it as the committed home for authorization documents | OBSERVED |
| `eval/` volume | 411 files (185 under `eval/runs`, 87 under `eval/s3`, 54 under `eval/harness`) | CALCULATED |

**The registry is not a usable index of current work.** It predates the entire November 16–18 corpus
and should not be used to decide what exists.

---

## 14. Contradictions, stale records, and unresolved evidence gaps

### 14.1 Stale "Gate 2 unauthorised" statements — four documents

| Document | Location | Status |
|---|---|---|
| `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | header, "Open inputs" | **Stale.** Flagged in R1 §"documentation deviation", left unfixed because editing a ratified artifact was prohibited |
| `MEMORY-09-SESSION-CORRECTIONS-2026-11-17.md` | §7 "Standing evidence boundaries" | **Stale** |
| `MEMORY-09-EXPERIMENT-PLAN-AMENDMENT-1-CONFORMANCE-2026-11-17.md` | gate table | **Stale** |
| `MEMORY-09-FUTURE-INTEGRATION-TESTING-NOTE-2026-11-17.md` | gate table | **Stale** |

Governing evidence: `MEMORY-09-OWNER-DECISIONS-2026-11-17.md`, §"Gate 2 — implementation of
`eval/mem1`". **Preserved, not corrected.**

### 14.2 The unreported implementation delta — the most serious gap

**1,951 committed lines, no covering report.** Three reports describe 886 / 1,632 / 1,803. There is no
record of the +148 lines, of what they implement, or of which requirements they touch. Any claim that
"the 118 requirements are mapped to the implementation" rests on a mapping of a tree that no longer
exists in that form.

### 14.3 The checkpoint's own header contradicts its body

`MEMORY-09-CHECKPOINT-2026-11-15.md` header: *"behavioural calls PAUSED before call 03"* and *"4 of
18 planned"*. Body: §12 records call 03, §13 call 04, §15 call 05, and §15.4 states **5 of 18**. The
header was written mid-series and never reconciled. The body governs; the header is **SUPERSEDED by
its own later sections**.

### 14.4 The "142-check baseline" is identified but its relationship is not

**OBSERVED:** 142 = `write-guard-v2` self-check count. **SOURCE-DERIVED:** D6 states 38 validators for
v3. Whether 38 replaces, subsets or adds to 142 remains **open** (Finding 6). The two are different
quantities and must not be compared directly.

### 14.5 Findings 5's premise may be a unit mismatch

See §5.2 — 12 artifacts vs 24 checks (WG01–WG24). Offered as evidence; **not** a closure.

### 14.6 D6's version targets are unmet

**OBSERVED:** `MEM1_VERSION = "0.1.0-g1"`; no `receipt_schema_version` symbol exists in
`eval/mem1/*.py`. D6 sets v3 targets of `0.2.0` / schema `2`. D6 is a closed decision about an
instrument that has not been built. `write-guard-v3` does not exist as a directory.

### 14.7 What could not be established

| Question | Why |
|---|---|
| Whether commit `830cba4` falls inside Gate 2's scope | No authorization document names a commit; scope moved twice |
| Experiment plan v3's approval status | Not established from the documents inspected |
| Whether the 10 assumptions still describe the committed tree | Requires a fresh mapping, which is out of scope here |
| Whether any F3-C provision was considered and rejected for MEMORY-09 | No record either way |
| Full contents of `docs/review/` (224 files) and `docs/research/` (288 files) | Not read file-by-file; see §16.2 |

---

## 15. Final evidence-based readiness assessment

### 15.1 The five states, separated

| State | Verdict | Basis |
|---|---|---|
| **Design-ready** | **YES** | v9 ratified, 118 requirements defined, r32 frozen, D1–D10 closed |
| **Implementation-complete** | **PARTIAL — and unreported at the committed state** | 1,951 lines committed; latest covering report is 148 lines behind; no report covers `830cba4` |
| **Conformance-verified** | **NO** | 0 requirements executed, 0 probes run, 0 tests run |
| **Runtime-authorized** | **NO** | Execution gate never entered; network responders refused in code |
| **Experiment-ready** | **NO** | Corpus absent, G07/G12 open, no grant mechanism, no live responder, no engine↔harness bridge |

**One-line classification: DESIGN-COMPLETE · PARTIALLY IMPLEMENTED · UNEXECUTED · NOT
EXPERIMENT-READY.**

### 15.2 The eight questions, answered directly

1. **What is genuinely finished?** The architecture and its ratification; the 118-requirement map;
  the frozen design contract; D1–D10; the preregistration and its three amendments; S3 (closed,
  including its negative result); the harness and its identity guard (`write-guard-v2`, 142/142).
2. **What is implemented but not behaviourally demonstrated?** The whole `m09` engine — 1,951 lines
  covering authority, records, promotion, disclosure, beliefs, conflict, retrieval. Its correctness is
  a source claim. R2 already names one **class C** gap (belief vs truth) that implementation cannot
  close and one **class B** gap (false conflicts from target equality).
3. **What has already been validated?** Static, zero-model checks only: 142/142 harness self-checks,
  96/96 earlier generation, S2 schedule validation (0 violations), frozen-identity verification,
  `validate` 89/89, corpus PASS, S1 batch audit PASS. No behavioural validation of the engine exists.
4. **What experiments have actually run?** S1: five generations at two token budgets, of which two are
  content-admissible. S3: Gate A 14/14, Gate B 6/7 with P5 open — **a negative result that must be
  preserved**. S2: nothing. G5: nothing. G3/G4: never entered.
5. **What work has been repeated or substantially overlaps?** No exact duplication. Three overlaps:
  (a) the November 18 recovery work restates the (b)/(c) deferrals rather than re-deciding them —
  correct, and correctly bounded; (b) it converges on the D8-U4-iii repair boundary already recorded —
  convergence, not duplication; (c) it runs parallel to F3-C recovery without ever citing it —
  **an unreconciled relationship, not a repeat** (§9).
6. **What is the exact remaining path?** P1–P12 in §12. The hard sequence is: corpus → frozen probe
  identities → `plan_sha256` → G07 + G12 bindings → grant mechanism → live responder → engine↔harness
  bridge → execution authorization → publication authorization. **Five of those do not exist in any
  form**, and three of them require code that no current authorization permits writing.
7. **Next smallest action that repeats nothing and violates nothing?** **Answer one question, in
  writing: is the experiment meant to test `eval/mem1/model.py` (what the harness actually drives) or
  `eval/mem1/m09/` (the ratified engine)?** Everything downstream — corpus shape, probe identities,
  the bridge, the grant mechanism — depends on that answer, and no document currently states it.
  Building the corpus before answering it risks a corpus built against the wrong subject, which is the
  one remaining way to genuinely duplicate work.
8. **What must NOT be reopened?** S3 (closed, with P5 preserved); the S1 scoring protocol and its
  amendment history; the 57 rules and 118 requirements; r32; D1–D10 as decisions; the (b)/(c)
  deferrals; `write-guard-v2`'s standing 142/142 identity guarantee; the F3-C frozen contracts; and
  Finding 2's classification (`false_closure`, partially resolved — **not** reclassified, **not**
  closed).

### 15.3 The one thing this audit most wants the owner to see

The project is not short of design, review, or documentation discipline. It is short of **one
decision and one bridge**: which memory implementation the experiment tests, and the code that would
let a conversation reach it. Until those exist, no corpus, no gate, and no authorization can produce a
local-model test of the memory system — and building any of them first is how the next round of work
would get repeated.

---

## 16. Complete source index for inspected primary records

### 16.1 Read in full or in the cited sections during this audit

| Document | Sections inspected |
|---|---|
| `docs/research/MEMORY-09-G5-SEQUENCING-D1-D10-DECISION-RECORD-2026-11-18.md` | 162–331 (D1–D6), 376–520 (D8), 700–859 (§5.3–§5.5), 840–987 (§6–§10) |
| `docs/research/MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | 520–619 (Gate 1, Gate 2, D-4) |
| `docs/research/MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | 1–100 (header, A–F), 100–211 (G–L, review answers) |
| `docs/research/MEMORY-09-CHECKPOINT-2026-11-15.md` | 1–120, 121–310, 400–563 |
| `docs/research/MEMORY-09-S2-MEMORY-INJECTION-DESIGN-2026-11-15.md` | 95–184 (§8–§14) |
| `docs/research/MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` | 1–130 |
| `docs/research/MEMORY-09-IMPLEMENTATION-REVIEW-R1-2026-11-17.md` | 1–85 |
| `docs/research/MEMORY-09-IMPLEMENTATION-REVIEW-R2-2026-11-17.md` | 1–90 |
| `docs/research/MEMORY-09-GATE2-IMPLEMENTATION-REPORT-2026-11-17.md` | 1–22, 43–100, 330–389 |
| `docs/research/MEMORY-09-GATE2-REPORT-SUPPLEMENT-2026-11-17.md` | 1–60 |
| `docs/research/MEMORY-09-R3-RECOVERY-CLOSURE-2026-11-17.md` | 1–110 |
| `docs/research/MEMORY-09-OPEN-DECISIONS-2026-11-15-REV-2.md` | full |
| `docs/research/MEMORY-09-REVIEW-LEDGER-2026-11-15-REV-2.md` | full |
| `docs/research/MEMORY-09-SESSION-CORRECTIONS-2026-11-17.md` | full |
| `docs/research/MEMORY-09-J3-SCOPE-DETERMINATION-2026-11-17.md` | 1–120 |
| `docs/research/MEMORY-09-OPERATIONAL-RECOVERY-ANALYSIS-2026-11-18.md` | 1–95 |
| The five other 2026-11-18 recovery documents | Content and hashes verified; prose read in earlier passes |
| `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md` | recovery sections (§2, §3, abandonment ordering) |
| `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md` | OD-02 boundary, frozen §4.1 ruling |
| `docs/review/ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md` | OD-02a/OD-02b status |
| `docs/FROZEN-ARTIFACTS.md` | 1–60 |
| `docs/DOCUMENT-REGISTRY.md` | 1–90 (header, rules, summary) |
| `docs/RECOVERY-GUIDE.md` | 1–70 (confirmed unrelated: spare-laptop restore) |
| `eval/mem1/__init__.py`, `m09/*.py` (all 6), `conversation.py` (62–84, 165–189) | direct inspection |
| `eval/mem1/write-guard-v2/g1-receipt-writeguard-v2.json` | parsed in full |
| `eval/mem1/g1-receipt.json`, `eval/mem1/g4-preflight-receipt.json` | header blocks |

### 16.2 Exclusions — what this audit did **not** inspect

| Excluded | Reason |
|---|---|
| 224 files under `docs/review/` beyond those listed | Sampled for relevance by keyword and title; not read individually |
| ~200 of 288 files under `docs/research/` | Same |
| `docs/architecture/` (17 files) except via the registry | Not read directly |
| 170 files under `docs/research/corpus/` | Corpus-governance records, not required for this status question |
| 185 files under `eval/runs/` | S1 records verified through the checkpoint and audit summaries, not opened individually |
| `eval/s3/` (87 files) | S3 status taken from the plan and checkpoint |
| Full Git history (50 commits) | Only HEAD and working-tree state inspected |
| Production Howl application code | Explicitly outside MEMORY-09's ratified scope |

**No claim of exhaustive audit is made.** The inventory searched is stated above; documents outside it
may contain material this record does not reflect.

---

### Verification of this document's own creation

| Check | Result |
|---|---|
| Files created | 1 — this document |
| Existing files modified | 0 |
| Tests, validators, probes, model calls executed | 0 |
| Servers or responders started | 0 |
| Git staging, commit, push | None |
| Registry registration | Not performed |
| Frozen artifacts touched | None |

**Boundary statement.** This document reports state. It does not change it. Finding 2 remains
`false_closure`, partially resolved and open on operational recovery. No gate is opened, no decision
is made, and no result is claimed.
