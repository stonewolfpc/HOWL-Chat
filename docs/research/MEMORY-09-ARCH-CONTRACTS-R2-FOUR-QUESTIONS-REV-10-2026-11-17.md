# MEMORY-09 — Architectural contracts, revision 10

**Status: PROPOSAL REV 10 — one response is assembled from one snapshot of accepted records, member
resolution and establishment results. Accompanies the owner's ratification of REV 9 within its declared
sections. Not implemented, not tested.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 9 (`5b3a9008…`) **for §3.4's final paragraph only**; adds invariants 13 and 14 and one implementation row |
| Stands unamended | REV 9 in every other respect; REV 1–8 preserved byte-identical |
| Authority | Owner's qualification and ratification direction, recorded verbatim (§1) |
| Ratification | **Recorded in §6.** REV 9 ratified within its declared sections, per the owner's direction; one reading is flagged there for confirmation |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger |

---

## 0. What changed

| Item | REV 9 | REV 10 |
|---|---|---|
| Establishment predicate | Called once per group, memoised — an implementation note in §3.4 | One snapshot value `S`; every field of the response is a function of `S` and nothing else |
| The accepted set | Computed once and passed (REV 8 §3.1) | Frozen into `S` before assembly, together with member resolution |
| Protection against mid-assembly change | None — correctness followed only if no write intervened | Stated as invariant 13, with the concurrency and staleness rules |
| `establishes` and `establishing_groups` | Must agree because both filter the same predicate | Must agree because both project one map inside `S` (invariant 14) |
| Observable behaviour in a quiet system | — | **Unchanged.** The clause constrains assembly, not verdicts |

No rule, requirement, state value or tally changes. Nothing new the engine may decide.

---

## 1. Owner instruction, recorded verbatim

> "The important distinction is now explicit: `open_groups` answers: Which open assertions may this audience
> inspect? `establishing_groups` answers: Which assertions actually govern state in this view? `relation`
> answers: Has incompatibility been explicitly asserted among the accepted records in this answer? That
> separation fixes the counter-example without suppressing historical assertions or weakening the authority
> view."

> "**One qualification before ratification.** I would have [the implementer] clarify one small but important
> implementation detail in §3.4: 'The state calculation and both published collections must use the same
> per-call snapshot of accepted records and group-establishment results.' The report says the predicate is
> called once per group and memoised. Good. But the contract should explicitly prevent the accepted-record
> set or group store from changing between the calculation of entry states and the construction of the two
> collections. Otherwise, the mathematical invariant could be correct while the response is assembled from
> inconsistent snapshots. This is a precision clarification, not a reason to reopen the architecture."

> "**My recommendation.** REV 9 accepted in principle. I recommend: 1. Resolve that snapshot-consistency
> wording. 2. Ratify REV 9 as the replacement for REV 8 within its explicitly declared sections. 3. Preserve
> every other REV 8 section and all REV 1–7 history. 4. Keep implementation and testing as separate
> authorization gates."

> "The remaining open items—B3, belief-voiced dialogue, and replay/F7—should stay open. They are not reasons
> to expand this correction."

**Scope:** snapshot consistency, and the ratification record the owner directed. No architecture change.

---

## 2. The gap, stated exactly

REV 9 §3.4 ended:

> "The state function and `establishing_groups` must be derived from the same predicate — one function,
> called once per group, memoised per call — so that a conflict entry and a governing group entry cannot
> disagree."

That constrains the **predicate**. It says nothing about the **inputs**. A response is assembled in several
passes over state the engine reads live:

| Pass | Reads |
|---|---|
| 1 | `_eligible_truth(domain, target)` → the accepted list, and the per-entry record counts |
| 2 | the group store → membership, and each member's `proposition_id` |
| 3 | the same accepted list again → `group_visible` / `group_contained` |
| 4 | the same group store again → the entries written into `open_groups` and `establishing_groups` |

Memoising a pure function of `(group, accepted)` does not help if `accepted` is a live view re-derived in
pass 3, or if pass 4 re-reads the group store. The predicate then has two arguments' worth of drift, and
the fields can contradict one another while each line of code is exactly as specified.

**Trace, entirely inside the proposed contract.** A `supersede` for `truth-00014` commits between passes 1
and 4 — a second caller, or an operator, in the window while a response is being built:

| Step | Result |
|---|---|
| 1 | Accepted includes `truth-00012` and `truth-00014`; `group-00004` establishes; entry `prop-00003` computed as `conflict` |
| 2 | The group store is re-read after the write; membership is unchanged, so the group is still listed |
| 3 | The accepted list is re-derived; `truth-00014` is now retired, so `group_establishes` is **false** |
| 4 | `establishing_groups` is `[]`, `relation` is `"not_established"`, and the entry still says `conflict` |

Invariant 11 (`relation == "asserted"` ⟺ some entry `conflict`) is violated — not because any single
definition is wrong, but because the two halves of the equation were answered from different worlds. The
response is then both misleading and, in the canon views, a disclosure fault: `open_groups` would name a
record that pass-3's accepted set excluded.

This is the fourth instance of one pattern, and worth recording as such because it is the reason the fix is
structural rather than another predicate:

| Revision | One thing, two meanings | Corrected by |
|---|---|---|
| REV 7 | One predicate for two questions (kind) | Record-level participation in REV 8 |
| REV 8 | One collection for two audiences | `open_groups` / `establishing_groups` in REV 9 |
| REV 9 | One derivation, two reads | One snapshot in REV 10 |

Each time, the contract was internally consistent and the composition was not.

---

## 3. The snapshot

### 3.1 Definition

A non-refused response is assembled from a value taken at the start of the call, before any state read used
to build it:

```
S = ( domain, target, view,             # call parameters, fixed for the response
      accepted,                         # accepted records: eligible ∧ kind-accepted, as a frozen sequence
      resolution,                       # record_id -> proposition_id, for every member named by any group
      groups,                           # open groups naming a returned identity, membership frozen
      establishes )                     # group_id -> bool, computed against this `accepted` and `resolution`
```

`S` is a value, not a view. After it is taken, no field of the response may be derived from a later read of
the record store, the group store or the registry, from a clock, or from a second pass over
`_eligible_truth`. Everything the response contains — entry `state`s, record counts, `relation`,
`open_groups`, `establishing_groups`, `integrity` — is a projection of `S`.

| Element | Why it must be in `S` | What breaks without it |
|---|---|---|
| `accepted` | It is the operand of `group_visible`, `group_contained` and the record counts | An entry's count and its group's visibility could disagree (the trace above) |
| `resolution` | `group_contained` reads a member's proposition from the record store | A member could resolve during one test and not the other; `unresolved` could be reported for a record that resolves, or vice versa |
| `groups` + membership | The two published collections are projections of the same membership | `open_groups` could list a group that `establishing_groups` never considered, or omit one it did |
| `establishes` | `state` and `establishing_groups` both consume it | The two could differ while the predicate is correctly memoised — memoised against *what* is the question |

### 3.2 Rules

| Rule | Statement |
|---|---|
| S-1 | **One take.** `S` is taken once per call, before the first read that contributes to the response. |
| S-2 | **No re-read.** Every response field is a projection of `S`. No field is computed from a post-`S` read. |
| S-3 | **One map, two projections.** `establishes` is computed once per group against `S.accepted` and `S.resolution`; entry `state`s and `establishing_groups` are read from it, never recomputed. |
| S-4 | **Staleness permitted, mixing not.** A concurrent write may make `S` stale; the response must still be a function of exactly one `S`. No read-your-writes guarantee is offered within a call, and none is required: mutating calls (`assert_incompatible`, `supersede`, `retract`) are separate calls whose effects appear in later snapshots. |
| S-5 | **Refusals take no snapshot.** A refused `propositions` or `authoritative` call stops at the gate before any data read (REV 7 §5.1), so it has no `S` and cannot leak one. |
| S-6 | **Assembled means assembled from `S`.** "Response `R` is consistent" is defined as: there exists an `S` such that `R = assemble(S)`. This is a testable property of a single response, without concurrency. |

### 3.3 What this does and does not settle

| Settled | Not settled |
|---|---|
| A response cannot mix two accepted sets, two resolutions, or two establishment results | Whether any particular `S` is the latest — that is the caller's concurrency control |
| `relation`, both collections and every `state` come from one `S`, so invariant 11 holds by construction | Whether repeated identical calls return identical bytes — ordering determinism is only the `proposition_id` / `record_id` ordering of REV 6/7 |
| A response is reproducible from its `S`, which makes it inspectable in review | **Replay.** `S` is not recorded, not persisted and not offered. F7 stays **R**: no replay facility exists and determinism remains unproven |

The last row matters enough to be blunt: this clause makes a *response* reconstructible from a value the
contract defines. It does not make the *engine* replayable, and it must not be read as closing F7.

**Current implementation.** `eval/mem1/m09/` is single-threaded and read-only on these paths, with no
`await` and no callback between the passes, so `S` is achievable by copying references into locals at entry
and never re-reading. The rules above state what must remain true if that ever stops being so; they are
satisfied trivially today and are not claimed to be tested.

---

## 4. Invariants 13 and 14

Added to REV 9 §6's list, which otherwise stands unchanged:

13. **Snapshot:** a non-refused response is a function of a single `S`. Every `state`, `relation`,
    `open_groups` entry and `establishing_groups` entry is derived from the same `accepted`, `resolution`
    and `establishes`. No field is derived from a read made after `S`.
14. **One map, two projections:** for every group id `g` in the response,
    `g ∈ establishing_groups` ⟺ `S.establishes[g]` ⟺ `g.establishes` (authority view) ⟺ `g` participates in
    every entry `state` computed as `conflict`. The four statements cannot differ.

Item 14 is the owner's memoisation requirement promoted from an implementation note to an invariant, with
the operand named. Item 13 is what makes items 11 and 12 unviolatable by assembly order rather than true by
inspection.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 5. Worked illustration

A quiet system and the trace of §2, side by side, on the same identity and group:

| | No write during assembly | `supersede(truth-00014)` between passes |
|---|---|---|
| `prop-00003` state | `conflict` | `settled` — computed from one `S`, which saw the write before it was taken |
| `relation` | `asserted` | `not_established` |
| `open_groups` (world) | `[group-00004]` | `[]` |
| `establishing_groups` (world) | `[group-00004]` | `[]` |
| `open_groups` (authority) | `[group-00004: establishes true]` | `[group-00004: establishes false, cause member_not_accepted]` |
| Consistent? | yes | yes — both are responses to *some* `S`; the right-hand column is a response to the newer one |

Without S-1 to S-4, the same write instead yields a **third** column that no `S` can produce: entry
`conflict`, `relation` `not_established`, `establishing_groups` `[]`, and `open_groups` naming
`truth-00014` although the accepted set excluded it. That column is the defect, and it is invisible to any
test that does not interleave a write — which is why it is a contract clause and not a test to be written
later.

---

## 6. Ratification, as directed

The owner directed: resolve the wording, ratify REV 9 as the replacement for REV 8 within its declared
sections, preserve every other REV 8 section and all of REV 1–7, keep implementation and testing separate.
That direction is executed as follows, and only as follows.

**Recorded ratification — REV 9, within its declared sections, effective 2026-11-17.**

| REV 9 section, now ratified | Replaces |
|---|---|
| §3.1–§3.3 — `open_groups` / `establishing_groups` separated; canon equality; authority relevance defined against returned identities | REV 8 §4.2's single-collection table |
| §4 — `relation` defined against `establishing_groups`; the four biconditionals with the canon-only qualifier | REV 8 §4's definition and biconditionals |
| §5 — three-view table | REV 8 §5 |
| §6 items 11 and 12 | REV 8 §6 item 11 |
| §7 — examples A and B revised | REV 8 §7 |
| §9 — rows 6⁗, 18‴, 9⁗, 17g, 17d″ | REV 8 §9's corresponding rows |

**Ratified with it:** REV 10 §3 (`S` and rules S-1…S-6), §4 (invariants 13 and 14), §9 row 17h — the
snapshot wording this ratification was conditioned on.

**Carried unchanged and unamended:** every other section of REV 8 and of REV 1–7, byte-identical.

**Still not ratified, and not implied by this record:** the ratified requirement set is unchanged at 57
rules and 118 requirements with their existing tallies; no section of v9, the map, r32, plan v3 or
Amendment 1 is amended; the ledger is untouched (no registration was directed).

**One reading flagged for confirmation.** The owner's disposition line reads "REV 9 accepted in principle",
and the direction then says "ratify REV 9 as the replacement for REV 8 within its explicitly declared
sections" while "preserv[ing] every other REV 8 section". I have recorded that as: REV 9's declared sections
are ratified contract text, and REV 8's remaining sections stand as the operative text they were when REV 9
replaced part of them — i.e. the ratified contract is REV 1–10 read as a chain, REV 9 superseding REV 8
where declared, REV 10 superseding REV 9 §3.4's last paragraph. If the owner intended to ratify only REV 9's
sections and leave REV 8's remainder at "accepted in principle", then this section overstates the record and
should be corrected before any implementation authorization; nothing outside this document has been changed
on that reading.

---

## 7. Implementation list — deltas only

REV 9 §9 stands, with these rows replaced or added. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17h | `engine.py` | `_snapshot(domain, target, view)` | Takes `accepted`, `resolution`, `groups`, `establishes` once; returns an immutable value; the only entry point for state reads on a response path |
| 17d‴ (REV 9 17d″) | `engine.py` | `group_establishes(...)` | Pure over `S`; results memoised **in `S`**, not per call site |
| 6⁵ (REV 9 6⁗) | `retrieval.py` | `world` | Consumes `S` only; no post-`S` read |
| 18⁴ (REV 9 18‴) | `retrieval.py` | `belief_facts` | Same |
| 9⁵ (REV 9 9⁗) | `engine.py` | `propositions(...)` | Consumes `S` only; gate still before `S` is taken (S-5) |

REV 9's `_relevant_groups`, `_accepted`, `_proposition_state`, `authoritative` and every other row are
unchanged. **D-A and D-B remain confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical), `eval/mem1/model.py`, the G1 harness, production Howl
Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–9.

---

## 8. Coherence

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14 | **C** |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** — now established once, not once per pass |
| Underspecification is not conflict | 8, 34 | K15 | **C** |
| No adjudication by the engine | 34 | K15, O-10 | **C** |
| History preserved, not re-reported as live | 20, 44 | D-02; K22 | **C** |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — gate still before any read, now before `S` |
| Canon reads add no text surface | 18, 32, 47 | REV 4 §4.2; L2; H7 | **C** — one `S` cannot widen the accepted set mid-response |
| Fail closed | 8 | F4 | **C** |
| Engine states what it did | 8 | K15 | **C** — the response now reports one state of the world, not a blend |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded; R-B, R-C, R-D open.

---

## 9. What this record does and does not claim

**Still open, carried forward unchanged:** B3; belief-voiced dialogue; replay / F7 — this revision defines a
snapshot, not a log, and explicitly does not close it; R-B, R-C, R-D.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- No behavioural verification, and no claim that the current single-threaded read paths are snapshotted —
  only that the rules are trivially satisfiable there.
- No concurrency test. §5 is a construction, not an observation.
- No new rule, requirement, state value or capability. S-1…S-6 constrain assembly; they change no verdict in
  a quiet system.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9 (`698254ca…`), the
  map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–9 (`98447423…`, `74ac6f77…`, `de4f68b6…`,
  `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`, `5b3a9008…`).
- No transfer, destination interface, sharing, or cross-chat carryover.

**Sequence from here:** implementation authorization (separate, not requested) → implementation →
conformance review → separate testing authorization. This record requests neither.
