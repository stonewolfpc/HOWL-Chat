# MEMORY-09 — Architectural contracts, revision 11

**Status: PROPOSAL REV 11 — a snapshot must be captured coherently, and the authority view's returned identity
set is part of it. Not implemented, not tested.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 10 (`7d703b77…`) **for §3.1's `S`, §3.2's rules, §3.3's first table row, §4's invariants 13, §5's illustration and §9 row 17h only** |
| Stands unamended | REV 10 in every other respect, including its §6 ratification record; REV 1–9 preserved byte-identical |
| Authority | Owner's two clarifications, recorded verbatim (§1) |
| Ratification | **Not ratified.** REV 10's snapshot contract was conditioned on these clarifications, so on acceptance the operative text is REV 11's §3, not REV 10's. REV 10 §6's record of the REV 9 ratification stands unchanged, including its flagged reading |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger |

---

## 0. What changed

| Item | REV 10 | REV 11 |
|---|---|---|
| Snapshot acquisition | "Taken once" — silent on whether the fields could have co-existed | S-7: every field of `S` must have co-existed at one committed version; atomic read, transactional read or version-check-and-retry, all acceptable |
| Snapshot use | S-1…S-6: no post-`S` mixing | Unchanged |
| `S` | `(domain, target, view, accepted, resolution, groups, establishes)` | `+ version` (admissibility only, never published) `+ returned_identities` (frozen, view-scoped) |
| Identity set | Implicit; the authority view's registry-derived set had no home | Captured in `S`; every consumer — enumeration, group relevance, scope — reads the frozen set |
| Registry reads | Not covered by S-2 | Covered: `returned_identities` is captured inside the same coherent read, so the gate precedes the capture and no post-`S` registry read exists |
| Failure to capture coherently | Undefined | Refuse, never publish a blend; a proposed `inconsistent_snapshot` code, unreachable in the current single-threaded store |

No rule, requirement, state value or tally changes.

---

## 1. Owner instruction, recorded verbatim

> "**1. Constructing `S` must itself be coherent.** REV 10 says `S` is taken once, then permits a concurrent
> write to make it stale. That is sensible. But taking one snapshot value is not necessarily the same as
> obtaining one atomic snapshot. Imagine `_snapshot()` reads: 1. The accepted records. 2. The proposition
> registry. 3. The group store. 4. The member-resolution data. A concurrent write between those reads could
> still produce an `S` that never existed as a coherent state. All subsequent response fields might derive
> perfectly from that `S`, but the snapshot itself would already combine different moments. S-1 through S-4
> prevent post-snapshot mixing. They don't yet explicitly prevent mixing during snapshot acquisition. The
> contract should require a coherent capture, using an atomic snapshot, transactionally consistent read, or
> equivalent version-check-and-retry mechanism. It need not prescribe which implementation."

> "**2. The authority identity set is missing from `S`.** … REV 9 deliberately allows the authority view to
> enumerate proposition identities whose records are all superseded. That means its identity set comes from
> the proposition registry, not merely from accepted records. But REV 10 defines `S = (domain, target, view,
> accepted, resolution, groups, establishes)`. There is no explicit registry snapshot or frozen
> returned-identity set. That matters because: Authority enumeration needs the registry to determine returned
> identities. Group relevance is defined against those returned identities. An identity with no accepted
> records may still have relevant historical groups. The response must not re-read the registry after `S` is
> taken. I recommend adding a frozen `returned_identities` value to `S`, populated according to the selected
> view and its authorization rules. Then group relevance, proposition enumeration, and state assembly all use
> that frozen identity set."

> "I would ask for a narrowly scoped REV 11 addressing only: 1. Coherent snapshot acquisition—not merely
> coherent use after acquisition. 2. Explicit inclusion of the authority view's returned identity set in `S`.
> No redesign, no new memory semantics, and no expansion of the ratified requirement set."

> "The underlying lesson is important: a snapshot is only a snapshot if all the data inside it could have
> existed together. A single variable containing several independently collected reads does not establish
> that."

**Scope:** capture coherence and the identity set. No architecture change, no new semantics.

---

## 2. Defect 1, stated exactly

REV 10 S-1 read: "`S` is taken once per call, before the first read that contributes to the response."
That answers *when* the reads happen relative to the response. It does not answer *whether the result is a
state*.

`supersede` is not one store write. It appends the replacement record and flips the old record's status.
Two mutations, in that order:

| Time | Record store, for one identity | A capture running concurrently |
|---|---|---|
| t0 | `truth-00012` ACTIVE | — |
| t1 | | Reads the record store for `accepted`; sees `truth-00012` ACTIVE |
| t2 | `truth-00018` appended, ACTIVE (the replacement) | — |
| t3 | | Continues the same pass; sees `truth-00018` ACTIVE |
| t4 | `truth-00012` flipped to SUPERSEDED | — |

The captured `accepted` contains **both** records of the identity as ACTIVE. No committed state ever had
that: before t2 the old record alone was live, after t4 the new one alone was. Every field derived from this
`S` is then internally perfect and jointly false — `state: unresolved_relation`, `relation: not_established`,
two live records where the caller's own read history shows one — and, worse, it is indistinguishable from a
genuine `unresolved_relation`, which is a real state the architecture must be able to report (rules 8, 34;
D-02 behaviour 8). The response invents a two-witness dispute.

This is not reachable only through `supersede`. The same class of tear exists wherever one logical change
touches more than one structure we read separately: the registry versus the record store (a minted
proposition with its first record, §4.2), the group store versus member resolution, and any future write
that the contract permits to be non-atomic. The point of the clause is that the guarantee must not depend on
each mutating call happening to be single-structure and ordered luckily.

REV 10's S-4 already conceded that staleness is acceptable and mixing is not. It simply never stated what
*mixing during acquisition* is, or how to know it did not happen. That is the gap.

---

## 3. Coherent capture

### 3.1 `S`, revised

```
S = ( domain, target, view,          # call parameters, fixed
      version,                       # committed version at which every field below co-existed; internal only
      accepted,                      # accepted records, frozen sequence
      returned_identities,           # frozen identity set, view-scoped (§4)
      resolution,                    # record_id -> proposition_id, frozen
      groups,                        # open groups naming a returned identity, membership frozen
      establishes )                  # group_id -> bool, computed against this accepted and resolution
```

The additions are `version` and `returned_identities`. `version` is diagnostic: it exists so admissibility can
be stated and checked. It is **not** published in any envelope — adding a field to responses was not asked
for and is not needed by any caller.

### 3.2 Rules

| Rule | Statement |
|---|---|
| S-1 | **One take.** `S` is obtained once per call, before the first read that contributes to the response. |
| S-2 | **No re-read.** Every response field is a projection of `S`. No field is computed from a read made after `S` — including any read of the registry. |
| S-3 | **One map, two projections.** `establishes` is computed once per group against `S.accepted` and `S.resolution`; entry `state`s and `establishing_groups` are read from it, never recomputed. |
| S-4 | **Staleness permitted, mixing not.** A concurrent write may make `S` stale; the response must still be a function of exactly one `S`. No read-your-writes guarantee is offered within a call. |
| S-5 | **Refusals take no snapshot.** A refused call stops at the gate before any data read, so it has no `S` and cannot leak one. |
| S-6 | **Assembled means assembled from `S`.** "Response `R` is consistent" means: there exists an `S` such that `R = assemble(S)`. |
| **S-7** | **Coherent capture.** Every field of `S` must have co-existed at one committed state. `S` is admissible only if there is a single committed version `v` with `accepted = accepted@v`, `returned_identities = identities@v`, `resolution = resolution@v` and `groups = groups@v`. Atomic snapshot, transactionally consistent read, and version-check-and-retry all satisfy this; the contract prescribes none of them. |
| **S-8** | **Frozen identity set.** `returned_identities` is captured inside the same coherent read as the rest of `S`, populated by the rule for the selected view (§4.1). Enumeration, group relevance and scope checks read the frozen value; none may consult the registry afterwards. |

**S-7's soundness condition, stated because a version test is easy to get subtly wrong.** A
compare-then-recheck on a version counter proves admissibility only if every mutating call publishes its
changes and advances the version in **one step** — no reader may observe a partial write under an unchanged
version. If that cannot be guaranteed, the version test must be replaced by a read that cannot observe
partial state. The contract requires the property, not the mechanism, and it does not permit a version test
that is only *usually* sound.

### 3.3 Failure to capture

There are two admissible responses and no third:

| Outcome | Envelope |
|---|---|
| An admissible `S` was obtained | Served from it |
| No admissible `S` could be obtained | `refused / inconsistent_snapshot` — the same uniform refusal shape as REV 4 §7 / REV 7 §5.1, carrying no registry-derived field, no ids and no counts |

Serving a blended `S` is not an option, because a blend is a state that never existed and the engine's whole
job is to say what is so. `inconsistent_snapshot` is a **proposed diagnostic code in a refusal**, not a new
rule or requirement; it is unreachable in the current implementation, which is single-threaded and
read-only on these paths. It exists so that an implementation that cannot capture coherently has a defined
behaviour instead of an undefined one.

### 3.4 What capture coherence is not

| Settled by S-7 | Not settled |
|---|---|
| `S` describes a state that once existed | Whether it is the newest — liveness remains the caller's concern |
| No response can present a blended state | Whether two concurrent responses describe the same state — they may legitimately describe different ones |
| A write landing mid-capture cannot produce a phantom `unresolved_relation` | Replay. `S` is still not recorded or persisted; F7 stays **R** |

---

## 4. The frozen identity set

### 4.1 Population, per view

| View | `returned_identities` | Source | Gate |
|---|---|---|---|
| `world` | identities of `accepted` records | `accepted` | none |
| `belief_facts` | identities of `accepted` records | `accepted` | none |
| `propositions` | every proposition identity for `(domain, target)`, including those whose records are all superseded | the registry | `require(domain, actor, "promotion")` **before** the capture begins |

The third row is the identity set REV 9 §3.3 reasoned about and REV 10 failed to place inside `S`. It is
not derivable from `accepted`: an identity with no accepted records still enumerates (`state: none`,
`records: []`, REV 4 §7), and groups naming only its retired records are still relevant and inspectable.

| Consequence | Rule |
|---|---|
| Group relevance for the authority view is defined against `returned_identities`, not against `accepted` | `groups` is selected by "names at least one record whose proposition is in `returned_identities`" |
| That selection needs member resolution for candidate groups, so `resolution` must cover every member of every group considered | acquisition resolves members of every open group, inside the same coherent read |
| A group naming only retired records is relevant to the authority view | it appears in `open_groups` with `establishes: false`, `cause: member_not_accepted` (REV 9 §3.3) |
| Canon views cannot widen their identity surface | their `returned_identities` is derived from `accepted`, so invariant 9 still holds by construction |
| The gate still precedes every read | a refusal has no `S`, so it cannot have an identity set (S-5) |

### 4.2 Why the registry pair is a genuine tear, not a hypothetical

Minting a proposition and writing its first record are two writes. A capture that reads the registry and
then the record store can land between them:

| Time | State | Capture |
|---|---|---|
| t0 | identity `prop-00031` does not exist | — |
| t1 | `prop-00031` minted, no record yet | reads registry: `prop-00031` ∈ `returned_identities` |
| t2 | first record written, ACTIVE | reads records: no record of `prop-00031` |
| t3 | — | `S` reports `prop-00031` as `state: none, records: []` |

The caller is told an identity has no records while a live record exists for it. That state never existed at
t1 (no live record) or at t2 (record present) — only in the blend. S-7 removes it, and S-8 makes it
impossible to reach the registry again later to "fix it up", which would have been the other blend.

---

## 5. Invariants 13–16

REV 9 §6 items 1–12 stand unchanged. REV 10's items 13 and 14 are replaced and extended:

13. **Snapshot:** a non-refused response is a function of a single `S` — `accepted`, `returned_identities`,
    `resolution`, `groups`, `establishes`. Every `state`, `relation`, both collections and every enumerated
    identity derives from that `S`, and from no read made after it.
14. **One map, two projections:** for every group id `g` in the response, `g ∈ establishing_groups` ⟺
    `S.establishes[g]` ⟺ `g.establishes` (authority view) ⟺ `g` participates in every entry computed as
    `conflict`. The four cannot differ.
15. **Admissibility:** `S` corresponds to one committed version — its fields co-existed. A response is either
    served from an admissible `S` or refused; a blend is not a representable response.
16. **Identity scope:** `returned_identities` is frozen in `S` and equals the identity set of `accepted` in
    canon views, and the target-scoped registry identity set in the authority view. Enumeration, group
    relevance and scope checks read that value; no consumer re-derives it.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 6. Worked illustration

| | Admissible `S` from before the write | Admissible `S` from after | Blend — now refused |
|---|---|---|---|
| identity records | `truth-00012` | `truth-00018` | `truth-00012` **and** `truth-00018` |
| `state` | `settled` | `settled` | `unresolved_relation` |
| `relation` | per groups, from this `S` | per groups, from this `S` | `not_established` |
| `open_groups` (canon) | establishing groups of this state | same, other state | may name a record outside `accepted` |
| Existed as a state? | yes | yes | **no** |
| Response | served | served | `refused / inconsistent_snapshot` |

The middle two columns are both correct answers to different moments; the architecture does not choose
between them, and no caller is owed the newer one. The right-hand column is the only one that is a lie, and
it is now unrepresentable.

---

## 7. Ratification status

| Item | State |
|---|---|
| REV 9 | Ratified within its declared sections, per the owner's direction, recorded in REV 10 §6. Unchanged here. |
| REV 10's snapshot contract | **Not ratified.** The owner conditioned it on these two clarifications, so on acceptance the operative text is **REV 11 §3–§5**, and REV 10 §3 is superseded wholesale rather than amended in place. |
| REV 8's unamended sections | Still flagged, as REV 10 §6 said: whether they are part of the ratified chain or remain "accepted in principle" is the owner's call. This revision neither resolves nor presumes it. |
| Tallies | Unchanged: 57 rules, 118 requirements, same statuses. |
| Implementation, testing | Separate gates. Neither requested. |

---

## 8. Implementation list — deltas only

REV 10 §7 stands, with these rows replaced or added. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17h′ (REV 10 17h) | `engine.py` | `_capture(domain, target, view, actor)` | Obtains one admissible `S` per S-7/S-8, including `version`, `returned_identities` and complete member resolution; gate runs before it is called |
| 17i | `engine.py` | `_version` | Monotone counter advanced in the same step as each mutating call's publication; read only by `_capture` for admissibility, never published |
| 17j | `engine.py` | `_relevant_groups(S)` | Selects by `S.returned_identities`; no registry access |
| 6⁶ (REV 10 6⁵) | `retrieval.py` | `world` | Consumes `S` only |
| 18⁵ (REV 10 18⁴) | `retrieval.py` | `belief_facts` | Consumes `S` only |
| 9⁶ (REV 10 9⁵) | `engine.py` | `propositions(...)` | Consumes `S` only, including its identity set; enumeration reads `S.returned_identities` |

REV 10's `group_establishes`, `_accepted`, `_proposition_state`, `authoritative` and every other row are
unchanged. **D-A and D-B remain confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical), `eval/mem1/model.py`, the G1 harness, production Howl
Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–10.

---

## 9. Coherence

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14 | **C** |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** — established once, on an admissible state |
| Underspecification is not conflict | 8, 34 | K15 | **C** — and a phantom `unresolved_relation` is now unrepresentable |
| No adjudication by the engine | 34 | K15, O-10 | **C** |
| History preserved, not re-reported as live | 20, 44 | D-02; K22 | **C** — retired-record groups stay inspectable to the authority view via a frozen identity set |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — gate before capture; frozen identity set inside it |
| Canon reads add no text surface | 18, 32, 47 | REV 4 §4.2; L2; H7 | **C** — canon identities derive from `accepted` inside `S` |
| Fail closed | 8 | F4 | **C** — no admissible `S`, no answer |
| Engine states what it did | 8 | K15 | **C** — one response, one state that existed |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded; R-B, R-C, R-D open.

---

## 10. What this record does and does not claim

**Still open, carried forward unchanged:** B3; belief-voiced dialogue; replay / F7 — S-7 makes a *response*
reviewable against a version, not the *engine* replayable; R-B, R-C, R-D.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- No concurrency test. §2, §4.2 and §6 are constructions, not observations. The current store is
  single-threaded and read-only on these paths, so S-7 is satisfied vacuously today and that is not
  evidence it would be satisfied once writes can interleave.
- No claim that `supersede` currently writes in two steps as described — R1/R2 established the shape of the
  operation, not its atomicity, and nothing was executed to determine it. The clause is written to be
  independent of the answer.
- No new rule, requirement, state value or capability. `inconsistent_snapshot` is a proposed refusal code,
  tally-neutral and unreachable today; `version` is internal and unpublished.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9 (`698254ca…`), the
  map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–10 (`98447423…`, `74ac6f77…`, `de4f68b6…`,
  `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`, `5b3a9008…`, `7d703b77…`).
- No transfer, destination interface, sharing, or cross-chat carryover.
- No change to the REV 10 §6 ratification record, including its flagged reading of REV 8's unamended
  sections.

**Sequence from here:** owner assessment of REV 11 → ratification of the snapshot contract as REV 11 §3–§5
→ separate implementation authorization → implementation → conformance review → separate testing
authorization. This record requests neither authorization.
