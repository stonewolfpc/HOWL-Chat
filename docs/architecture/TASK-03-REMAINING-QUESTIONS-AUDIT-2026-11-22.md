# TASK-03 — bounded audit of the remaining open questions

| Field | Value |
|---|---|
| Purpose | Split every remaining TASK-03 question into the four decision domains the owner named, and say exactly what each can change about v3.2 |
| Status | **Audit and proposal. It re-splits the existing register; it opens no investigation of its own. Nothing adopted. No implementation authorized.** |
| Baseline audited | v3.2 (row 583, revised under row 584) — 104 clauses, 21 decision rows, 62 tests |
| Predecessor | row 576 (`OQR-01`–`OQR-32`), which consolidated the distributed register; this record re-splits that register against v3.2 and adds one entry for the index producer |
| Recovery token | `TASK-03-REMAINING-QUESTIONS-AUDIT-2026-11-22` |

Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## §1 — Scope, and the finding that constrains all four buckets

**RQ-01 — What this is.** A single table per decision domain, covering every still-open question from row 576 plus the one opened by v3.2's index revision. Each entry says what the question can change, what it cannot, and whether the contract already absorbs its absence. It is the four-bucket answer to the owner's request, not a new programme.

**RQ-02 — What it is not.** It reopens no resolved question (OQ-1, OQ-2, OQ-3b-b, row 550's reporting effect, `V3-44`'s wording, memories), investigates nothing new, decides nothing about repair, migration, locks or caller surfaces, and authorizes nothing. Where it disagrees with row 576 it says so and says why.

**RQ-03 — The assignment rule, stated so it can be checked.** A question is assigned to the bucket where its **decision** lands. A dependency on another bucket is named inside the entry and is not counted twice. Three register entries — row 550, memories, `OQ-7c-ii` — have a **closed contract half** and an **open defect half**; both halves are named, the open half is counted in its bucket and the closed half under Closed.

**RQ-04 — Finding: no remaining question changes the rule of a currently reachable row.** This is the substantive result. Since v3.2, what remains is three items that change *which inputs exist* (OQ-6, OQ-7, OQ-8), eight about *whether the gate can be built*, four about *what the caller sees*, eight *defects and separate work*, and two *measurements*. Nothing left can make a reachable row decide differently. That is why adoption is decidable now rather than after them — and equally why adopting the contract does not make any of them go away.

---

## §2 — Bucket 1: questions that can change the v3.2 decision table

| Question | What it changes | What it cannot change | Already fail-closed? |
|---|---|---|---|
| **OQ-6** settlement record | Adds a settled-generation condition to rows 16–17, and makes `V3-44` case 3 a describable input | Rows 18–19 still report `attribution_unestablished`; no reachable row moves; no backfill is permitted (row 562 `S-16`) | Yes — rows 16–17 are unreachable without the gate, so their condition is dormant |
| **OQ-7** enforcement record | The only legitimate source of `gate_enforced = true`; makes rows 16–17 reachable | Rows 18–19's verdict; nothing about rows 1–15 or 20 | Yes — `gate_enforced` is false on any doubt (`V3-36`) |
| **OQ-8** index producer (new) | Whether rows 1 and 20 can be evaluated at all | Their conditions, which are now fixed (`V3-13d`, `V3-89`) | Yes — `index_state ≠ present` reports neither `Missing` nor `Orphaned` |

**RQ-05 — OQ-6 is a condition on two rows, not a change of decision.** Row 562's recommendation adds `settled_generation ≥ G-1` to rows 16–17. Both rows are unreachable until an enforcement record exists, so the amendment is invisible to every conversation the current product can produce. Its real value is that it would let a *future* generation ≥2 conversation with a completed clear be described at all.

**RQ-06 — OQ-7 is a capability, and it is the only one that makes two rows live.** `V3-36` fixes that `gate_enforced` cannot be declared, configured or inferred. There is no second route by which rows 16–17 become reachable, and no route by which they become reachable *silently*: a false `gate_enforced` sends such content to rows 18–19, which is the fail-closed answer.

**RQ-07 — OQ-8 is the entry this record adds, and it is narrower than it looks.** It does not ask what the index should contain. It records that no code path writes a membership record (Observed: `GetChatList` at `internal/backend/chat/persistence.go:485` enumerates `<rootDir>/chats` and skips directories whose `meta.json` fails to load; `AppConfig` holds only `LastChatID`), so the product's "index" is a filtered directory listing. The contract now defines the input and refuses to invent the artefact (`V3-89`). Until it exists, `Missing` and `Orphaned` are unreportable — a capability gap, stated in the contract, not a defect in the table.

**RQ-08 — Consequence for adoption, stated plainly.** Adopting v3.2 adopts a procedure that says two of its six states cannot be produced by today's product. That is the same shape as rows 16–17 and is preferable to the alternative: an index assumed to exist, against which every directory would be either "indexed" or "orphaned" by default. The default would have been a claim, and the default here is silence.

---

## §3 — Bucket 2: whether the mutation gate can actually be implemented

| Question | State | Blocking? |
|---|---|---|
| **OQ-7a** path centralisation | Verdict: the containment assertion is a **heuristic today** — three path conventions and no directory constructor (row 565), so a build assertion over literal write sites would miss an assembled path | **Yes** |
| **OQ-7b** lock scope | Answered: per conversation, in a sibling namespace outside `chats/`, with a deterministic total order for the two-acquisition operations; placement forced by creation, deletion and wipe (row 567) | No — unimplemented, not undecided |
| **OQ-7c-i** root authority | Answered enough: the history root is authoritative; no lock may key on `ProjectRoot()`; the claim is per user per machine, and roaming profiles are a local serialiser only (row 569) | No |
| **V2-OQ-2** `cleared` in the marker vs a separate record | Undecided schema question that precedes any marker field; a second document adds a crash domain and a two-document consistency problem | **Yes for OQ-6**, not for the gate |
| **OQ-5** manifest versioning | Classified outside by row 576; **reclassified here as a conditional prerequisite** — see RQ-10 | **Conditional** |
| Enforcement self-test | Undesigned. `V3-36` permits only a build-bound positive result; row 563 shows a startup self-test covers exactly the paths someone remembered | **Yes** |
| **OQ-3b-c** marker placement and scaling | Unanswered placement question that any marker implementation must settle first | No |
| **OQ-3a-c** a future import inheriting the representation gap | Not a current writer; becomes a `V3-35` obligation when that writer exists | No |

**RQ-09 — The gate's blocking set is small and mostly structural.** Path centralisation and the self-test are the two items that cannot be decided by writing more contract: one needs one directory constructor and one write boundary, the other needs a capability shape (unexported writer, single `Mutate`) rather than a checklist. Lock scope and root authority are already decided well enough to build against. Settlement, versioning and marker placement become blocking only for the *marker* half.

**RQ-10 — OQ-5's reclassification, with its exact condition.** Row 576 put manifest versioning outside. That holds only while the marker schema does not change. If OQ-6 is pursued — any new field — then an older build reading a newer marker has no defined behaviour, and the unknown-schema refusal path becomes a gate requirement rather than a storage nicety. So: **outside if OQ-6 is not pursued; prerequisite if it is.** Recording the condition rather than the class is the honest form, because the class depends on a decision not yet made.

**RQ-11 — The ordering that already exists is technical, not administrative.** OQ-6 must not precede OQ-7 (row 562 `S-20`, row 563 `O7-24`): a settlement field without an enforced boundary records what a writer asserted while appearing to prove more. No other pair has a fixed order.

**RQ-12 — The smallest buildable slice, stated as scope and not as a proposal.** A refusal-only boundary with per-conversation locks, the marker CAS, the tail rule and the containment assertion is implementable without OQ-5, OQ-6 or any marker field. It cannot set `gate_enforced` without the self-test, so rows 16–17 stay unreachable and generation ≥2 content stays Pending. Nothing here authorizes building it.

---

## §4 — Bucket 3: caller behaviour and unsaved-work reporting

| Question | What it decides | Feeds a row? |
|---|---|---|
| **OQ-7d-iii** caller-visible outcomes | Row 575's seven-field result over stored / captured / displayed, and the disposition of work that was captured but not written | **No** |
| **Operation identity** and the retry policy | Whether a refused write may be retried under a durable key, or refusal is final. Owner policy; row 574 showed identity is not derivable from existing data | **No** |
| **OQ-3c-c** `inconsistent` actionability | What the user can do after the verdict | No |
| **OQ-3a-b** `SaveChatAs` repointing | Which conversation is selected after duplication | No |

**RQ-13 — This bucket is where refusal becomes visible, and `V3-27` already makes silence non-conforming.** A caller that discards a typed refusal is not a policy choice; it is a defect against a rule already in the contract. That is why shutdown reporting (`runStage` reporting only to stdout, row 575) appears here rather than in bucket 4: the probe exists, the verdict string does not.

**RQ-14 — A boolean cannot carry this bucket.** Row 575's sixth position is decisive: `appendChatMessage` can capture and display a message and return `nil` when no conversation is selected. "Saved: false" and "not attempted" and "captured but unwritten" are three different facts, and only one of them is a failure the user can act on.

**RQ-15 — Retry policy is not a reconciliation input.** `V3-12` forbids attribution from anything but the marker and the transition, so no retry record can feed a row. Refusal-only is sufficient for the gate; operation identity is needed only for ambiguous retry. It is a policy choice between two sound designs, which is why it sits here and not in bucket 1.

**RQ-16 — The retargeting rule belongs here even though its defect is bucket 4.** `O7-157` — every mutation uses the identity bound to the content it mutates; a deferred operation that finds its conversation no longer bound refuses rather than retargeting — is a product rule about caller binding. No lock enforces it, because the second writer acquires the wrong conversation's lock legitimately.

---

## §5 — Bucket 4: separate foundational defects and separate work

| Item | Why it is not reconciliation's | Status |
|---|---|---|
| **Row 550** creation collision | Reporting effect closed (row 580 / `V3-72a`); the defect is a product requirement — collision safety is separate from crash-safe creation | Owned elsewhere, unrepaired |
| **OQ-3c-a** loader precedence | Display only; `V3-05` and `V3-06` already separate display from verdict | Owned, unrepaired |
| **OQ-3c-b** `RewriteHistory` resurrection | Produces inputs rows 12/13 already classify; the harm is reached by *reading* (`V3-43`) | Owned, unrepaired |
| **Row 570** final-save window | Every clean shutdown passes through log ⊆ snapshot; writer ordering, not classification | Owned, unrepaired |
| **OQ-4** repair path | Changes files between scans; `V3-47` and `V3-55` forbid the scanner from doing it | Not started |
| Torn-line repair | Same class as OQ-4; row 7 already reports. The gate's terminate-or-refuse (`V3-08a`) is not repair | Not started |
| **OQ-7d** non-history ownership | Twelve files, seven application-owned; directory adjacency establishes nothing; user memories are the only irreplaceable non-history data | Boundary stated, ownership open |
| **`OQ-7c-ii`** retargeting defect | Closed as contract wording (`V3-13a`); the shipped defect remains — an autosave captured in X can overwrite Y | Owned, unrepaired |
| Memories architecture | Not a defect of this contract: row 582 puts it outside the verdict entirely. Separate foundational work, with today's implementation explicitly not the target | Separate |

**RQ-17 — The test for this bucket is whether absorbing the item would change a verdict.** Every entry fails that test. A repair path, a loader preference, a writer rewrite and a shutdown ordering all change files or display; reconciliation reads bytes and reports. Absorbing any of them would repeat the mistake row 558 found in the addendum stack.

**RQ-18 — Creation is the one item that escalated out of this bucket.** Its reporting effect closed, but the underlying non-atomicity is now `V3-72a`'s product-correctness requirement — no creation may silently adopt, merge with, overwrite or redirect writes into an existing conversation. It is stated in the contract and owned outside it, which is the only arrangement that does not let the scanner compensate.

---

## §6 — The residual: measurements that decide nothing

**RQ-19 — Two questions are measurements and are not a bucket.** OQ-3a-a (collision frequency) and OQ-3b-a (per-scan cost) can change a risk sentence and cannot change a state, a capability or a caller outcome. `V3-09` already fixes the empty thresholds by observation, and row 551 explicitly refuses to make per-scan cost a budget. They are retained as evidence and assigned nowhere, because a bucket implies a decision and these have none.

**RQ-20 — Accounting, so nothing is unclassified.** Thirty-two register entries plus OQ-8: **3** in bucket 1, **8** in bucket 2, **4** in bucket 3, **8** in bucket 4, **2** residual, **10** closed. That is 35 assignments for 33 entries, because row 550 and memories each carry a closed half and an open half. No entry is unassigned. The ten closed: OQ-1 (row 577), OQ-2 (row 578), OQ-3 (row 548), OQ-3b-b (row 579), OQ-7c-iii (row 572), OQ-7d-i (row 573), OQ-7d-ii (row 574), row 550's reporting effect (row 580), `V3-44`'s wording (row 581), the memories disclosure (row 582).

---

## §7 — What adoption would and would not mean

**RQ-21 — Adoption is adoption of the procedure.** It would make v3.2 the single authoritative ordered decision procedure and would close the six baseline questions permanently for this contract. It would close none of the 23 open entries, and it would not create a backlog with an implied order.

**RQ-22 — Adoption is not authorization, and this record does not weaken that boundary.** Gate, CAS, locks, root observation, physical identity, tail rule, index producer and memory redesign are all unimplemented and none is authorized by the contract or by this audit.

**RQ-23 — The four conditions for adoption, and their state.** (a) the index producer gap is stated in the contract, not assumed away — holds (`V3-89`, `V3-92`); (b) the coverage limit is stated, so a clean result cannot be read as proof that no unindexed data exists — holds (`V3-13a`, `V3-90`); (c) an input that previously fell between "candidate" and "row 4" is classified — holds (row 4, `V3-91`); (d) no remaining question changes a reachable row's rule — holds (`RQ-04`).

**RQ-24 — Nothing was found that should block adoption.** That is a finding about the contract, not a recommendation to adopt it, and not a substitute for the owner's decision.

**RQ-25 — The four buckets are decision domains, not a schedule.** This record orders nothing beyond the technical constraint in RQ-11. A reader who treats the bucket numbers as priority has added a decision that does not exist here.

---

## §8 — Checks this audit implies

Not executed; these are checks on the register and on any future consolidation, not tests of code.

- **RQ-T1** Every entry in row 576's register appears exactly once as an open assignment or under Closed; an entry that is missing or double-counted fails.
- **RQ-T2** No bucket-1 entry is described as changing the rule of a reachable row; a bucket-1 item that claims otherwise fails.
- **RQ-T3** A scan with `index_state = absent` reports zero `Missing` and zero `Orphaned`, on every root class (`VT3-59`, `VT3-61`).
- **RQ-T4** A clean sweep reports its candidate count and the coverage non-claim, never a conversation total (`VT3-55`, `VT3-62`).

---

## §9 — Standing

| Item | State after this audit |
|---|---|
| Baseline | v3.2, revised once (rows 583–584) — 104 clauses, 21 decision rows, 62 tests |
| Bucket 1 — changes the table's inputs | 3: OQ-6, OQ-7, OQ-8 |
| Bucket 2 — gate implementability | 8, of which 3 blocking (OQ-7a, self-test, V2-OQ-2 for the marker half) and OQ-5 conditional |
| Bucket 3 — caller behaviour | 4 |
| Bucket 4 — separate defects and work | 8 |
| Residual — measurements | 2 |
| Closed | 10 |
| Unclassified | **none** |
| Adoption | not granted; not recommended here |
| Implementation | **not authorized** |
