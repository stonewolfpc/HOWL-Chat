# TASK-03 — open-question register audit and consolidation point

| Field | Value |
|---|---|
| Purpose | Consolidate every open TASK-03 question into one register and determine which can still change the accepted V3.1 baseline or a prerequisite of it |
| Status | **Audit and proposal. V3.1 unedited. No question closed by this record. No implementation authorized.** |
| Method | Every trailing "still open" table in rows 547–575 read and merged; each question classified by what it can change, not by how interesting it is |
| Recovery token | `TASK-03-OPEN-QUESTION-REGISTER-AUDIT-2026-11-22` |

Consumes `V3-71`…`V3-73` (row 560), the ownership boundary `V3-38`, and the open tables of rows 548–575. This record adds no investigation of its own; it classifies what already exists.

---

## 1. The register is distributed, and that is itself a finding

**OQR-01 — Observed: there is no open-question register.** The open questions exist only as trailing tables in ten documents: row 547 §9, row 549 §9, row 551 §9, row 552 §9, row 553 §8, row 560 `V3-71`/`V3-73`, row 563 §6, row 566, row 568, row 572 and row 574. No file lists them together, and no file states which are closed.

**OQR-02 — Calculated: this violates the pack's governing principle in its own housekeeping.** "No state may depend on a condition hidden in an amendment, footnote, transition narrative, or test expectation." The *state of the investigation* currently depends on ten footnotes. A reader cannot tell from V3.1 alone whether a question that would change a row is open, closed, or was never asked.

**OQR-03 — The classification used here is by consequence, with five values:**

| Class | Meaning | Blocks consolidation? |
|---|---|---|
| **baseline** | can change a §4 state decision or the terms on which a row is reached | **yes** |
| **prerequisite** | can change whether §2/§3 machinery exists, but no §4 row | no — V3.1 already declares the dependency |
| **evidence** | can change a measurement, a rationale or a reachability label, never a decision | no |
| **outside** | belongs to repair, migration, UI or another task | no |
| **closed** | answered by a named row; retained for traceability only | no |

A question that cannot be given a class is not open — it is unclassified, and unclassified is a defect of the register, exactly as an unreachable row is a defect of the table (`V3-39`, `V3-66`).

---

## 2. The consolidated register

**OQR-04 — Every question currently recorded, with its class and disposition.** "Disposition" is what consolidation must do about it, not what this record has done.

| # | Question | Source | Class | What it can change | Disposition |
|---|---|---|---|---|---|
| OQ-1 | UNC/network roots, where "absent" may mean "disconnected" | row 547 §9; `V3-73` | **baseline** | row 1 `Missing` versus row 2 `Unavailable` for a whole root class | decide an exclusion or a root-class condition |
| OQ-2 | Case sensitivity and junctions inside a conversation directory | row 547 §9; `V3-73` | **baseline** | whether two listed conversations are one physical directory (rows 1, 6, 12, 20) | define candidate identity as the resolved physical directory, or exclude |
| OQ-3 | Whether `meta.json` alone is content-present | row 548 | **closed** | — | answered by row 548: `meta.json` alone is not content |
| OQ-3a-a | Same-second creation collision frequency | row 549 | evidence | exposure only; the mechanism is confirmed | measurement, never a rule |
| OQ-3a-b | Whether `SaveChatAs` should repoint to the copy | row 549 | evidence | selection behaviour only | product decision, no verdict |
| OQ-3a-c | Whether a future import reuses `SaveChatAs` | row 549 | prerequisite | whether a future writer inherits the representation gap | bind to `V3-35` when that writer exists |
| OQ-3b-a | Exact per-format empty thresholds and per-scan cost | row 551 | evidence | constants and performance | `V3-09` already fixes 0 and 2 bytes; no decision |
| OQ-3b-b | Whether a missing final newline is Damaged, a warning, or nothing | row 551 | **baseline** | row 7 applies or does not for a well-defined byte string (`V3-08`) | **decide before consolidation** |
| OQ-3b-c | Where the marker document lives and how it scales | row 551 | prerequisite | §2 placement, not a row | implementation of §2 |
| OQ-3c-a | Loader preference for a torn log over a fuller snapshot | row 552; `V3-38` | evidence | what the user is shown, never the verdict (`V3-05`, `V3-06`) | separately owned; no baseline change |
| OQ-3c-b | `RewriteHistory` resurrection after a crash | row 552; `V3-38` | evidence | produces inputs rows 12/13 already classify | separately owned; no baseline change |
| OQ-3c-c | Whether `inconsistent` should be user-actionable | row 552 | outside | a post-verdict affordance | UI/product, not classification |
| OQ-4 | Whether Damaged conversations get a repair path, and what repair may do | `V3-73` | outside | files between scans; never a state | explicitly not in the contract (`V3-47`, `V3-55`) |
| OQ-5 | Manifest migration versioning across storage changes | `V3-73` | outside | row 20 index inputs only | later storage contract |
| V2-OQ-2 | `cleared` in the marker versus a separate retraction record | `V3-73` | prerequisite | §2 shape, rows 11 and 13 evidence | schema decision before marker implementation |
| Row 550 | Creation identity collision's effect on rows 6 and 12 | `V3-72`, `V3-73` | **baseline** | a legitimate creation can produce `inconsistent` / `unattributed` | limitation clause; never absorbed as repair |
| OQ-6 | Whether the marker records settlement of the last clear | `V3-71` | prerequisite | §2, and `V3-44` case 3 becoming an input | declared owner decision; rows already fail closed |
| OQ-7 | Whether an enforcement record exists | `V3-71` | prerequisite | §3, `V3-36`/`V3-37`, reachability of rows 16–19 | declared; rows already unreachable |
| OQ-7a | Whether paths can be centralised enough for the negative assertion to be proof | row 563 §6; row 565 | prerequisite | the strength of `O7-38`'s build assertion | enforcement prerequisite |
| OQ-7b | Lock scope: per conversation or per document | row 563 §6; row 567 | prerequisite | §3's single critical section | enforcement prerequisite |
| OQ-7c-i | Root derivation: what makes a lock root stable and authoritative | row 568; row 569 | evidence | lock identity across profiles; **not** the conversation scan (history root is stable) | answered enough by row 569; no row |
| OQ-7c-ii | Pointer consistency and execution-time retargeting | row 568; rows 571–572 | **baseline** | the truth of `V3-44`'s generation-1 Healthy justification | limitation clause (§4) |
| OQ-7c-iii | Whether anything requires true re-entrancy | row 568; row 572 | **closed** | — | answered by row 572: no product requirement |
| OQ-7d | Where non-history writers belong | row 563 §6; rows 566–567 | outside | nothing in the scan; boundary already excludes them | boundary statement only |
| OQ-7d-i | Scope of the installation-tree store | row 566; row 573 | **closed** | — | answered by row 573: mixture, outside the boundary |
| OQ-7d-ii | Whether retry identity is derivable | row 574 | **closed** (investigation) | the retry *policy* remains an owner choice | policy selection, not further investigation |
| OQ-7d-iii | Caller-visible outcomes and unsaved-work disposition | row 575 | prerequisite | §3's failure visibility (`V3-26`, `V3-27`) | proposal awaiting owner adoption |
| Torn-line repair | Whether a torn append is repaired, and by whom | row 574 `O7-210` | outside | files, not classification; row 7 already reports it | repair task, like OQ-4 |
| Operation identity | Whether a durable operation identity is adopted | row 574 §8; row 575 §8 | prerequisite | ambiguous retry only; **no §4 row** | separately specified capability |
| Row 570 | Final-save representation window | row 570 | outside | writer ordering; the scanner already reports the resulting rows | repair task |
| Memories | Whether `Healthy` may imply memories were assessed | row 573 | **baseline** (wording) | the meaning of the Healthy verdict, not a row | one mandatory sentence in V3.2 |
| Enforcement self-test | How `gate_enforced` is proved and bound to a build | `V3-36` | prerequisite | rows 16–19 remain unreachable until it exists | already specified as derived, default false |

**OQR-05 — Count, stated as arithmetic so it can be checked:** thirty-two entries — six baseline, ten prerequisite, six evidence, six outside, four closed. Of the six baseline entries, **five can change a §4 decision and one (memories) changes the meaning of the Healthy verdict without moving a row.**

---

## 3. The five questions that can still change the baseline

**OQR-06 — OQ-1 can flip row 1 to row 2 for an entire root class.** `V3-60` makes incomplete enumeration an error, but a disconnected UNC path that stats as absent is not an enumeration error — it is a successful observation of absence. `V3-39` row 1 then reports `Missing`, which is a claim about deletion the scanner cannot support. The baseline has no root-class condition. Either V3.2 adds one (network roots report `Unavailable / unreadable` without a stat) or it excludes network roots by statement. Silence is not an option, because silence selects row 1 by default.

**OQR-07 — OQ-2 can make two conversations one directory, or one directory two conversations.** The row 519 audit already established that identity must be physical, not lexical, because a junction can put a path inside the root while the bytes are outside it. The same fact decides this question: if the scan's candidate identity is the path string, case-variant entries and junction aliases are counted separately; if it is the resolved physical directory, they collapse to one candidate. `V3-39` does not say which. This is a **term** in §1, not a new row, but without it rows 1, 6, 12 and 20 have no defined input.

**OQR-08 — OQ-3b-b is the only open question that flips a row on a byte string.** `V3-08` requires a non-empty log to end with a newline; `V3-09` gives the 1-byte `"\n"` case. A complete final record without a terminator is therefore `malformed` today, and row 7 reports Damaged. The question asks whether it should instead be harmless because the loader parses it. Both answers are defensible; they are not compatible. **The fail-closed answer is to keep the strict rule** — a missing terminator is byte-indistinguishable from a write cut exactly at the record boundary, so treating it as harmless weakens the only structural torn-tail signal the format has. That is a decision V3.2 must state; leaving it open leaves row 7's boundary undefined.

**OQR-09 — Row 550's collision can produce a baseline state from a legitimate operation.** Two creations in one second share a directory; the later `SaveChat` overwrites `meta.json` and `history.json` while the earlier log survives (`RC5-1`, `RC5-3`). The scan then sees metadata describing one conversation and content belonging to another. `RC5-7` requires this be reported as an *anomaly*, but V3.1 has no anomaly state and no anomaly reason; the nearest outcomes are row 6 `inconsistent` and row 12 `unattributed`. Both are truthful and neither is the reported shape. V3.2 must either state that a collision is reported through those rows and is not a separate state — the honest, non-absorbing option — or add vocabulary. It may not remain an unstated effect on rows 6 and 12.

**OQR-10 — OQ-7c-ii undermines a sentence, not a row: `V3-44`'s generation-1 case.** That case says a generation-1 conversation is Healthy because "no preceding generation exists, so nothing can contaminate it". Row 572 shows a different contamination that needs no generation: an autosave captured while conversation X was selected can execute after selection moved to Y, load Y and overwrite it with X's in-memory content (`O7-157`'s defect, observed in shipped code). The scanner cannot see this — Y's files parse, the marker is `(none,0)` or `committed,1`, and row 10/14/15 reports Healthy. The verdict is correct about *observed content* and wrong if read as *this conversation's content*. V3.2 must narrow the clause: **Healthy asserts that content is present and attributable as §1 defines it; it never asserts authorship, writer identity or conversational binding**, and cross-conversation retargeting is outside what a scan can detect. That is a limitation on the meaning of a verdict, which is a baseline change even though no row moves.

**OQR-11 — The memories sentence is the sixth baseline change, and the smallest.** Row 573 establishes that `WipeAllData` removes only the history root, that user-authored memories live elsewhere and survive, and that the scan never assesses them. A `Healthy` verdict that does not say so invites the reader to conclude the whole conversation — memories included — is available. V3.2 needs one sentence: `Healthy` is a verdict about conversation history only; memories were not assessed. Row 573 already requires this; this audit records that it belongs in the baseline, not in the investigation paper.

---

## 4. Prerequisites, which do not block consolidation

**OQR-12 — A prerequisite can make a capability exist; it cannot change a classification.** OQ-6, V2-OQ-2, OQ-7, OQ-7a, OQ-7b, OQ-3b-c, OQ-3a-c, OQ-7d-iii, operation identity and the enforcement self-test all sit here. Each changes a schema, a lock, a boundary or a result type — never a §4 row.

**OQR-13 — The decisive property is that V3.1 already treats their absence as a state, not a defect.** Rows 16–19 are declared unreachable while `gate_enforced` is false (`V3-36`, `V3-37`, `V3-39`'s stated exception), and row 5 already stops a `Recovered` marker. So a missing prerequisite produces the fail-closed verdict the baseline already specifies. A missing *baseline decision* does not: row 1 will confidently say `Missing` on a disconnected share, and row 7 will confidently say `Damaged` on a terminator question nobody decided.

**OQR-14 — Therefore prerequisites may remain open at consolidation, provided each is named in the baseline.** What may not remain open is a question whose default answer silently favors one outcome. OQ-1 and OQ-3b-b are exactly that; OQ-7 is not, because its default is `false` and `false` is the fail-closed branch.

**OQR-15 — Operation identity, specifically, cannot extend the contract.** It is a record about whether an operation applied; no §4 row reads it, and `V3-12` forbids attribution from anything but the marker and the transition. Row 575 shows the separation: conversation binding is needed for refusal-only retry, operation identity only for ambiguous retry. Neither is a reconciliation input. This is the answer to the concern that the retry branch is opening a second contract in parallel: it is not, and this audit is the boundary that says so.

---

## 5. What changes only evidence or display

**OQR-16 — OQ-3a-a and OQ-3b-a are measurements, and neither is a rule.** Collision frequency is an exposure estimate (`RC5-A`); per-scan cost is a performance figure (`RC4-17` explicitly refuses to make it a budget). `V3-09` already fixes the empty thresholds at 0 and 2 bytes as observations. A measurement can change a sentence about risk, not a state.

**OQR-17 — OQ-3c-a and OQ-3c-b are writer defects the scanner observes, never imitates.** `V3-05` already says `LoadChat` is loadability and never validity; `V3-38` already assigns both defects elsewhere. A loader that shows a torn log's prefix can make the *display* disagree with the verdict — which `V3-06` permits, because display is not a state. A rewrite that resurrects content produces content rows 12/13 already classify. Neither can add a row or move one.

**OQR-18 — OQ-3c-c and OQ-3a-b are affordances.** Whether `inconsistent` is actionable, and whether `SaveChatAs` repoints, change what the user can do after a verdict and which conversation is selected. Neither changes what the verdict is.

**OQR-19 — OQ-7c-i is evidence about lock identity, not about conversation identity.** Row 569 establishes the history root as stable and `ProjectRoot()` as launch-dependent. The scan reads the history root, so the instability that matters to a lock does not move a verdict.

---

## 6. What is outside the contract

**OQR-20 — OQ-4, the torn-line repair question, row 570's window, and OQ-5 are outside, and the exclusion is structural.** `V3-47` forbids any automatic repair and `V3-55` forbids the scanner from rewriting, choosing or re-running a writer. A repair path changes files between scans; it cannot change the procedure that reads them. OQ-5 concerns the index used by row 20 and cannot create a seventh state.

**OQR-21 — OQ-7d and row 573's shadow tree are outside by boundary, not by omission.** `V3` scans `<rootDir>/chats`; candidacy is decided by `meta.json`; the installation-tree namespace gets no verdict. The boundary is already stated in row 573 and needs only to be carried into V3.2's scope sentence.

**OQR-22 — Deletion remains the one exempt mutation (`V3-34`) and creates no new question here.** Deletion during an unresolved transition is permitted and promotes nothing; deletion's partial-completion shape is row 8, already decided.

---

## 7. The discontinuation rule

**OQR-23 — A new investigation may be opened only if it can change a §4 row, a §2/§3 prerequisite, or is required by a legal or product obligation.** Anything else is recorded as a limitation, not investigated. This is the rule that stops the pack extending indefinitely: the last four records each opened a sub-branch, and each sub-branch was legitimate because it named a capability the baseline depended on. The register is now at the point where further branches would be about implementation detail.

**OQR-24 — Every open question must carry a class and a disposition.** A question with no class is not open. A question whose disposition is "investigate further" must name the row or prerequisite it could change; "it would be good to know" is a limitation.

**OQR-25 — A limitation is a first-class output, not a failure.** Several of the five baseline questions will close as limitations rather than as conditions: the honest statement that network roots are out of scope, that a collision reports through rows 6/12, that Healthy does not assert binding. A contract that states its blind spots is stronger than one that closes them by assumption.

---

## 8. The consolidation point, and the answer to the question asked

**OQR-26 — Answer: yes, six things can still change the accepted baseline, and no prerequisite can.** The six are OQ-1, OQ-2, OQ-3b-b, row 550's collision, the `V3-44` generation-1 wording, and the memories sentence. Everything else in the register is either closed, a prerequisite whose absence the baseline already handles fail-closed, evidence that cannot move a row, or outside the contract.

**OQR-27 — The closure set is bounded and can be written as a V3.2 amendment of six clauses, not a new investigation programme:**

| # | Amendment | Kind |
|---|---|---|
| 1 | Network roots: either a root-class condition reporting `Unavailable / unreadable`, or an explicit scope exclusion | decide |
| 2 | Candidate identity is the resolved physical directory, not the path string | term |
| 3 | A non-empty representation without its terminator stays `malformed` (row 7), with the torn-tail rationale | decide |
| 4 | A creation collision is reported through rows 6/12 and is not a separate state; never repaired | limitation |
| 5 | `Healthy` asserts present, attributable content; it does not assert authorship or writer binding; cross-conversation retargeting is undetectable | limitation |
| 6 | `Healthy` is a verdict about conversation history only; memories were not assessed | limitation |

**OQR-28 — V3.2 should also carry the register itself.** One table, in the baseline, with each question's class and disposition. That is the mechanical form of OQR-01–OQR-03: the state of the investigation becomes readable from the authoritative document instead of from ten trailing tables.

**OQR-29 — What V3.2 must not do.** It must not add a seventh state, an anomaly state, a per-record generation, a retry ledger, a repair path or a root-class taxonomy beyond item 1. Each of those would be an extension, and OQR-23 says an extension needs a row or a prerequisite it can change. None of them can name one.

**OQR-30 — The retry branch closes here.** OQ-7d-ii answered the investigation (row 574); OQ-7d-iii answered the caller consequences (row 575). The policy choice — refusal-only baseline versus durable operation identity — is an owner decision, not an open investigation. No further OQ-7d sub-branch is warranted, and none is opened by this audit.

**OQR-31 — Nothing in this audit changes V3.1.** The six amendments are proposals; the register is a proposal; the dispositions are recommendations. V3.1 remains the accepted baseline, unedited, and rows 547–575 remain the evidence.

**OQR-32 — And nothing is authorized.** No table edit, no term change, no state, no reason, no root-class condition, no physical-identity requirement, no limitation rendered into UI copy, no repair, no migration, no deletion enablement, and no implementation of any prerequisite. The audit's only product is the bounded closure set in OQR-27.

---

## 9. Checks this audit implies

These are checks on the register and the future consolidation, not tests of code, and none was executed.

- **OQR-T1** Every open question in rows 547–575 appears in the consolidated register with a class and a disposition; a question present in a source table but absent here fails.
- **OQR-T2** A disconnected-root fixture is reported `Unavailable`, never `Missing`, once OQR-27 item 1 is decided.
- **OQR-T3** A non-empty log lacking its final newline lands on the row OQR-27 item 3 decides, and no other input changes state as a result.
- **OQR-T4** A case-variant or junction alias of one physical conversation is counted once under OQR-27 item 2, not twice.
- **OQR-T5** A collision-shaped directory reports one of rows 6/12 and is never repaired, renamed or split (`RC5-8`, `RC5-9`).
- **OQR-T6** A `Healthy` verdict string states both limitations of OQR-27 items 5–6: no binding claim, memories not assessed.

---

## 10. Standing

| Item | Status after this audit |
|---|---|
| Open-question register | consolidated here; not yet in the baseline |
| Baseline-affecting questions | **6** — OQ-1, OQ-2, OQ-3b-b, row 550, `V3-44` wording, memories sentence |
| Prerequisite questions | 10, all fail-closed by the existing baseline |
| Evidence / outside / closed | 16 (6 / 6 / 4) |
| Retry branch | closed as investigation; policy unselected |
| V3.1 | unchanged, still the accepted baseline |
| Consolidation | bounded six-clause V3.2, plus the register table |
| Implementation | **not authorized** |
