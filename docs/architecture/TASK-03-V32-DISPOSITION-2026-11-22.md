# TASK-03 — v3.2 disposition: retained as reviewed draft, not adopted

| Field | Value |
|---|---|
| Purpose | Record the owner's disposition decision on reconciliation contract v3.2 and name the two decisions that follow it |
| Status | **Owner disposition. Nothing adopted. No implementation authorized** |
| Basis | v3.2 as revised (row 584); the four-bucket audit (row 585); rows 583–585 |
| Edits | None. This record amends no clause, adds no state, row, reason, test or bucket |
| Recovery token | `TASK-03-V32-DISPOSITION-2026-11-22` |

Labels: **[Observed]** read in this workspace · **[Source-derived]** stated by a cited row or clause · **[Calculated]** arithmetic on recorded figures · **[Inferred]** reasoning from those, not observed · **[Unestablished]** no evidence in this session.

---

## §1 Boundary of this record

**DP-01 — This is a disposition, not a contract amendment.** It records a status decision about v3.2. It does not edit v3.2, and it changes no clause, decision row, state, reason or acceptance test. **[Observed]** v3.2's bytes are unchanged by this record.

**DP-02 — Authority.** The disposition is the owner's, stated on 2026-11-22. This record is the register of that decision, not its source and not a substitute for the owner's words.

**DP-03 — No adoption, no implementation, and no implication of either.** Retention is not adoption. A pass-with-limitations audit is not adoption. This record is not adoption. Nothing here authorizes a gate, lock, index producer, repair, migration, deletion, marker change, storage-format change, memory redesign or retry policy.

**DP-04 — Expansion stops here.** The decision this record dispositions is that the contract is not to be expanded further. No clause, row, state, reason, transition, test or bucket is added by it.

---

## §2 The decision

**DP-05 — v3.2 is retained as a reviewed draft.** Retained means: current, reviewable, citable only as a draft, and not withdrawn. The text retained is the revision registered as row 584 — 104 clauses, 21 decision rows (1, 2, 2b, 3–20), 62 acceptance tests, hash `302519a3300bdd451ccbcf9cd440305bb74a2a99bc3de8bd4f6129269ad604f7`. **[Observed]**

**DP-06 — Row 585 is the audit record of that draft.** The four-bucket audit is the evidentiary basis of this disposition. It is not superseded by this record and is not folded into the contract. **[Source-derived]** row 585's own scope statement.

**DP-07 — Not adopted. Not rejected.** Both outcomes remain available. Each requires an explicit owner act. Silence, this disposition, an implemented gate, or a passing test is not adoption — the rule `C-01` states for FC-1 applies unchanged here. **[Source-derived]** row 583's audit and row 585's `RQ-24`.

**DP-08 — Nothing historical moves.** v3.1 and rows 547–582 remain byte-identical; row 583 remains the record of v3.2's first revision and is superseded for identity only by row 584. **[Source-derived]** rows 584's scope statement.

---

## §3 The five preserved conclusions

**DP-09 — The reachable-state rules have no identified open-question dependency.** No remaining question changes the rule of a currently reachable row; the twenty-one rows decide from the same inputs after the audit as before it. **[Source-derived]** row 585 `RQ-04`, `RQ-23(d)`. **[Unestablished]** by execution: this is a claim about the register, not a demonstration that any row behaves as written.

**DP-10 — Coverage is deliberately incomplete without an independent index producer.** Membership is what row 1 and row 20 test against, and no code path writes a membership record: `GetChatList` enumerates `chats` with `os.ReadDir` and `continue`s past directories whose `meta.json` fails to load, and `AppConfig` carries only `LastChatID`. **[Observed]** in `internal/backend/chat/persistence.go`. **[Source-derived]** `V3-13a`, `V3-13d`, `V3-89`, `V3-92`. **[Inferred]** consequence: `Missing` and `Orphaned` are operationally unavailable, because `index_state = present` cannot be satisfied by any build of the current product. That is a coverage limit recorded in the contract, not a defect the scanner may compensate for.

**DP-11 — The mutation gate is not demonstrated implementable; three blockers remain.** The blocking set is OQ-7a (path centralisation, today a three-convention heuristic with no constructor), the undesigned enforcement self-test (a startup test covers the paths someone remembered), and V2-OQ-2 for the marker half. **[Source-derived]** row 585, bucket 2. "Not demonstrated" is a statement about the evidence; it is not proof of impossibility, and it is not a finding that the gate is wrong. **[Unestablished]** whether the three can be closed, and at what cost.

**DP-12 — The separate writer defects remain separately owned.** The final-save representation window, the retargeting defect, the unatomic creation path, retry identity, repair and non-history ownership are defects of writer ordering, identity and lifecycle, not reconciliation questions. **[Source-derived]** rows 570, 574, 580, 585 bucket 4. Absorbing any of them into the verdict would repeat the mistake row 558 found in the addendum stack.

**DP-13 — The textual adoption conditions are not evidence of runtime correctness.** `V3-88`–`V3-92` and `RQ-23` audit text and structure: no duplicate clause identifier, no undefined reference, no state added, every reason already in the vocabulary. None of that observes behaviour. **[Source-derived]** rows 584–585.

**DP-14 — What verification actually was performed, stated as it was.** `go build ./...` exited 0 and the product tree carries 0 tracked or staged Go changes. **[Observed]** No test of the contract, of any proposed enforcement, or of any row's behaviour was executed for v3.2, for row 585 or for this record. A compile is evidence about the compiler; it is not evidence that the procedure decides correctly, that a gate can be built, or that a verdict matches the bytes on disk.

---

## §4 What retention changes, and what it does not

| Item | Effect of this disposition |
|---|---|
| v3.2 text | Unchanged. Any later edit changes the bytes, hence the hash, hence requires a new registry row and supersedes row 584 for identity **[Calculated]** |
| Row 583 | Still the record of the first revision; its hash stays stale by design |
| `V3-73a` | Still reads "subject to owner adoption", because adoption was not granted **[Observed]** at `:433` |
| Open questions | 23 unchanged; none closed, none reordered **[Source-derived]** row 585 |
| Buckets | Still decision domains, not a schedule; the numbers carry no priority **[Source-derived]** `RQ-25` |
| Schedule or backlog | None created. Retention sets no deadline and implies no sequence |
| Implementation | Still unauthorized. Retention grants no permission by being "reviewed" |

---

## §5 The next foundational decision is two decisions

**DP-15 — Two separate decisions are now open, and neither is taken here.**

- **D-2 — Resolve the gate's three implementation blockers.** The bucket 2 blocking set of DP-11.
- **D-3 — Formally accept the contract's current coverage boundary.** No independent index producer; `Missing` and `Orphaned` unavailable; memories outside the verdict; `Healthy` an availability claim within the assessed scope, with its five non-claims.

**DP-16 — They are separable, and this record orders neither.** D-2 without D-3 resolves implementability while leaving the operating boundary unaccepted. D-3 without D-2 accepts a boundary whose enforcement is undemonstrated. Both are coherent positions; the choice between them is the owner's next decision, not a consequence of this one.

**DP-17 — Neither requires reopening the historical investigations.** Rows 577–582 are closed inputs. D-2 consumes evidence already recorded in rows 567–569, 571, 572 and 574. D-3 consumes the coverage statements already in `V3-13a`, `V3-13d`, `V3-89` and row 582. **[Source-derived]** No new investigation is a precondition of either.

**DP-18 — D-3 is a coverage decision, not an adoption decision.** Accepting the boundary would settle whether the current limit is an acceptable operating envelope. It would not adopt v3.2, and adopting v3.2 would not by itself accept the boundary. DP-07's requirement of an explicit act stands for both.

**DP-19 — Whichever is taken, the fail-closed behaviour is not traded away.** `index_state = absent` yields no `Missing` and no `Orphaned`; a scan error never becomes `Missing`; `Unavailable` never becomes `Healthy`; a refusal is never a content verdict; an unparseable `meta.json` is `identity_lost`. **[Source-derived]** `V3-13d`, rows 1, 4 and 20.

---

## §6 Standing after this disposition

| Item | State |
|---|---|
| v3.2 | Retained reviewed draft, row 584 text, 104 clauses / 21 rows / 62 tests |
| Audit | Row 585, complete, not superseded |
| Adoption | **Not granted** |
| Implementation | **Not authorized** |
| D-2 gate blockers | Open, owner decision |
| D-3 coverage boundary | Open, owner decision |
| Open questions | 23 unchanged, none closed here |
| Tests executed | None |

---

## §7 Checks this record implies

Not executed. These check the register against the contract, not code against behaviour.

- **DP-T1** v3.2's SHA-256 still equals the value recorded in row 584. A difference means the draft was edited after this disposition and this record's reference is stale.
- **DP-T2** v3.1 and rows 547–582 still match their recorded hashes; this record must not be the cause of any change to them.
- **DP-T3** The three blockers named in DP-11 are still exactly bucket 2's blocking set. A blocker appearing or disappearing without a new audit row means the register and this disposition disagree.
- **DP-T4** `V3-73a` still ties closure to owner adoption. A change to that clause without an adoption record is a defect in the register, not an adoption.

---

**Scope of this record.** It records a disposition decision and nothing else. It amends no contract, registers no open question as resolved, and grants no implementation authorization. D-2 and D-3 are named as open and are not taken, recommended or ordered here.
