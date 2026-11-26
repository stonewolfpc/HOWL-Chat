# MEMORY-09 — Architectural contracts, revision 12

**Status: PROPOSAL REV 12 — refusal classes separated by whether data was read, and candidate group acquisition
defined as a total, ordered procedure. Not implemented, not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 11 (`40850419…`) **for §3.2's S-5 and S-8, §3.3, §4.1's "every group considered", §5's invariant 16 and §8 rows 17h′/17j only** |
| Stands unamended | REV 11 in every other respect; REV 1–10 preserved byte-identical |
| Authority | Owner's two clarifications, recorded verbatim (§1) |
| Ratification | **Not ratified.** On acceptance the operative snapshot text becomes REV 12 §2–§4 |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger |

---

## 0. What changed

| Item | REV 11 | REV 12 |
|---|---|---|
| Refusals | One rule, S-5: "refusals take no snapshot" | **Two classes.** Authorization refusals (S-5, unchanged in force) take no snapshot and read nothing. Capture refusals (S-9) occur after the gate, discard any partial state, and publish nothing derived from it |
| `inconsistent_snapshot` | Introduced without reconciling S-5 | Defined as a capture refusal: reachable only post-gate, no data-derived field, content-independent |
| Candidate groups | "every group considered" — subset left to the implementation | The candidate set is **the whole open group store**; relevance filtering happens after resolution, inside the capture |
| Member resolution | "covers every member of every group considered" | A **total** map over every candidate member reference, three-valued, with `unresolved` as an explicit value — never an absent key |
| Unresolved vs unexamined | Indistinguishable | Distinct: `unresolved` is a value; "not examined" is not representable, because an incomplete capture is refused |
| Invariants | 13–16 | 13–18: adds capture-procedure totality (17) and refusal-class separation (18) |

No rule, requirement, state value or tally changes. Nothing the engine may decide changes.

---

## 1. Owner instruction, recorded verbatim

> "**1. S-5 now contradicts the new failure behaviour.** … REV 11 retains: 'Refusals take no snapshot. A
> refused call stops at the gate before any data read, so it has no S.' But REV 11 §3.3 introduces another
> refusal: `refused / inconsistent_snapshot`. … Those are fundamentally different situations: Authorization
> refusal: the gate rejects the call before data access. No snapshot is taken. Capture refusal: the gate
> permits the call, capture is attempted, but no admissible snapshot can be established. The second case
> necessarily involves data reads or attempted snapshot acquisition. Saying all refusals take no snapshot is
> therefore no longer accurate. I recommend distinguishing the two explicitly. Preserve the original security
> guarantee for authorization refusals, while defining `inconsistent_snapshot` as a capture failure that
> returns no data-derived response fields. This does not weaken the refusal boundary."

> "**2. Define the candidate group set before resolving group members.** … The engine must first identify the
> candidate open groups, resolve their member references against the same coherent underlying state, and
> then select the groups relevant to the frozen identity set. Otherwise, the phrase 'every group considered'
> leaves room for different implementations to resolve different subsets. I would also explicitly represent
> unresolved member references in `resolution`, rather than allowing an absent mapping to ambiguously mean
> either unresolved or not examined. That preserves the distinction between a genuinely unresolved
> historical reference and incomplete snapshot acquisition."

> "I would request a narrowly scoped REV 12 addressing only: 1. Distinguish pre-read authorization refusals
> from post-attempt capture refusals. 2. Define candidate group acquisition and explicit unresolved member
> resolution. No changes to the memory semantics, ratified requirements, or authority hierarchy."

> "One other thing worth preserving: REV 11 correctly refuses to claim that the current single-threaded
> implementation has demonstrated concurrency safety. A contract can establish what future implementations
> must guarantee; it cannot turn an unexecuted concurrency scenario into verified behaviour."

**Scope:** refusal taxonomy and acquisition procedure. No semantics, no requirements, no hierarchy.

---

## 2. Two refusal classes

### 2.1 The contradiction

REV 11 S-5: "Refusals take no snapshot. A refused call stops at the gate before any data read, so it has no
`S`." REV 11 §3.3: "No admissible `S` could be obtained → `refused / inconsistent_snapshot`". The second
requires having attempted the reads the first forbids. One of them had to go, and it is S-5's universal
quantifier, not its guarantee.

### 2.2 The classes

| | Class A — authorization | Class B — capture |
|---|---|---|
| When | Before any data read | After the gate admits the call; during acquisition |
| Cause | The caller may not make this call, or the view does not exist | The caller may make it; the state could not be captured coherently |
| Data read | **None.** Not the record store, not the group store, not the registry | Reads attempted; whatever was read is discarded |
| `S` | Does not exist | No admissible `S` exists; partial state is dropped |
| Codes | `no_domain_authority`, `unknown_view` | `inconsistent_snapshot` |
| Views reachable | `propositions` (gate), any (`unknown_view`) | All three — each captures |
| Envelope | Fixed refusal shape; no data-derived field | **Identical** fixed shape; no data-derived field |

### 2.3 Rules

| Rule | Statement |
|---|---|
| S-5 (narrowed) | **Authorization refusals take no snapshot and read no data.** The gate, and view validation, complete before any store access; a Class A refusal has no `S` and its envelope contains no id, count, text or version. |
| S-9 (new) | **Capture refusals publish nothing derived.** A Class B refusal may have read state while attempting capture; it must discard it and return the same fixed refusal shape, with no field derived from any candidate snapshot, no partial collections, and no `version`. |
| S-10 (new) | **Ordering.** Class B is unreachable before the gate. An unauthorized caller therefore only ever receives a Class A code, and a Class B code can never disclose whether a domain holds data. |
| S-11 (new) | **Content independence.** The choice of refusal class and code must not depend on the content, size or emptiness of the stores. `inconsistent_snapshot` reports an acquisition outcome, not a fact about the domain. |

S-5's original guarantee is preserved exactly where it mattered — the security boundary is the gate, and a
Class A refusal still cannot touch a byte. What is abandoned is only the claim that *every* refusal is a
Class A refusal, which was true under REV 4 and false the moment REV 11 introduced a post-attempt failure.

### 2.4 Why this does not weaken the boundary

| Concern | Answer |
|---|---|
| Does a Class B code leak content? | No field is derived from data; S-11 makes the code independent of content. A caller learns "capture was contended or inconsistent", which is a property of attempts, not of the chat. |
| Can Class B be used to probe for existence? | Only post-gate, so the caller was already entitled to the data; and S-11 forbids the code varying with emptiness. |
| Does discarding partial state cost evidence? | The discarded state is inadmissible — it is a blend, so it is not evidence. Keeping it would be worse than losing it. |
| Does a Class A refusal still reveal the view list? | `unknown_view` answers only whether a view name exists, which REV 7 §5.1 already settled. |

---

## 3. Candidate groups and total member resolution

### 3.1 The acquisition procedure

Ordered, total, implemented once. Every step after the gate occurs inside the same capture attempt, under
the version bracket of S-7.

| Step | Action | On failure |
|---|---|---|
| A1 | Validate `view`; resolve `(domain, target)` | Class A: `unknown_view` |
| A2 | For `view == "propositions"`: `require(domain, actor, "promotion")` | Class A: `no_domain_authority` |
| A3 | **Candidate set:** `candidate_groups = every group in the open group store` — unfiltered | — |
| A4 | **Member references:** `member_refs = ⋃ members(g) for g ∈ candidate_groups`, as a set of record ids | — |
| A5 | **Total resolution:** `resolution = { m ↦ resolve(m) for every m ∈ member_refs }`, three-valued (§3.2) | — |
| A6 | **Accepted records:** `accepted` = eligible ∧ kind-accepted for this view | — |
| A7 | **Returned identities:** canon from `accepted`; authority from the registry, target-scoped | — |
| A8 | **Relevance:** `groups = [g ∈ candidate_groups : ∃ m ∈ members(g), pid(resolution[m]) ∈ returned_identities]` | — |
| A9 | **Establishment:** `establishes[g]` for every `g ∈ groups`, against `accepted` and `resolution` | — |
| A10 | **Version recheck** | Changed → discard everything, retry; attempts exhausted → Class B: `inconsistent_snapshot` |
| A11 | Assemble the response from `S` | — |

Three properties this fixes:

| Property | Why it matters |
|---|---|
| Candidate set is the **whole** store, not a pre-filtered subset | Relevance needs `resolution`, and `resolution` is part of the capture. Filtering first would require reading before the capture — the very thing S-7 forbids. Deterministic and implementation-independent: every implementation examines the same set. |
| Relevance is applied **after** resolution, inside one capture | A group becomes relevant because of what its members resolve to at the captured version, not because of a prior hint |
| `resolution` is total over `member_refs` | Every member of every candidate group has an entry; no consumer may treat absence as "not examined" |

Cost is a consequence, not a rule: acquisition is `O(open groups + members)` per call regardless of the
target's size. An implementation may add an index later, provided A3–A8 remain a single capture at one
version. The contract specifies the set and the order, not the data structure.

### 3.2 Resolution, three-valued

```
resolution[m] = resolved(proposition_id)      # the record exists; its identity is known
              | unresolved(record_absent)     # examined; the record store has no such record
```

| Value | Meaning | Effect |
|---|---|---|
| `resolved(pid)` | Examined, found, identity known | Usable by A8 and by `group_visible` / `group_contained` |
| `unresolved(record_absent)` | **Examined, not found** — a divergence between the group store and the record store | `group_visible` and `group_contained` both fail (REV 8 §3.2); authority view reports `establishes: false`, `cause: "member_unresolved"`; canon views surface `group_member_unresolved` in `integrity` only |

There is no third value, and in particular **no absence**. "Not examined" is unrepresentable because an
incomplete capture is a failed capture (A10/S-9), not a partially-filled one. This is the distinction the
owner asked for: a genuinely dangling historical reference and a truncated acquisition previously looked
identical in the data structure, and now they cannot.

A member whose record exists but whose proposition belongs to another `(domain, target)` is `resolved`, and
its group is simply not relevant to this response (A8). That is scope, not a fault, and it is not reported
as one.

---

## 4. Invariants 17–19

REV 9 §6 items 1–12 stand. REV 10–11 items 13–16 stand with 16 restated in terms of A8. Added:

17. **Capture procedure:** the candidate set is every open group; `resolution` is total over its member
    references; relevance is applied to `returned_identities` after resolution; establishment is computed
    for every relevant group. Every implementation examines the same groups and produces the same `S` for
    the same version.
18. **Resolution totality:** every member reference of every candidate group has exactly one value, and every
    value is either `resolved(pid)` or `unresolved(record_absent)`. No consumer may infer "not examined"
    from a missing key, because missing keys are not produced.
19. **Refusal separation:** Class A refusals read no data and contain no data-derived field; Class B refusals
    are reachable only after the gate, discard any partial capture, and contain no data-derived field. Both
    use one fixed envelope shape, and the class and code never vary with store content.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 5. Worked illustration

**A. Two refusals, side by side.**

| | Class A | Class B |
|---|---|---|
| Call | `propositions(world, target="king.status", actor="stranger")` | `world(world, "king.status")` during a contended capture |
| Reads before refusal | none | attempted |
| `S` | absent | discarded |
| Envelope | `{ "status": "refused", "code": "no_domain_authority" }` | `{ "status": "refused", "code": "inconsistent_snapshot" }` |
| Data-derived fields | none | none |
| What the caller learns | "you may not ask this" | "the state could not be captured" |

`unknown_view` behaves like the left column: validation precedes every read.

**B. Unresolved versus not examined.**

| Case | `resolution[m]` | Reported |
|---|---|---|
| Group names `truth-00099`, no such record exists at the captured version | `unresolved(record_absent)` | authority: `establishes: false`, `cause: "member_unresolved"`; canon: `integrity: ["group_member_unresolved"]` |
| The capture was truncated before reaching that group | — | no response at all: Class B refusal |
| Group names `truth-00014`, which was superseded | `resolved("prop-00003")` | authority: `establishes: false`, `cause: "member_not_accepted"`; canon: silent |

Three visibly different outcomes for three different situations, where REV 11 could have produced one.

---

## 6. Ratification status

| Item | State |
|---|---|
| REV 9 | Ratified within its declared sections, per the owner's direction; recorded in REV 10 §6, unchanged. |
| REV 10's snapshot contract | Not ratified; superseded by REV 11 §3–§5 and now by REV 12 §2–§4. |
| REV 11's snapshot contract | Accepted in principle, pending these two clarifications. On acceptance, §2 (refusal classes) and §3 (acquisition) of this record are the operative text, with REV 11 §3.2's S-5 and S-8 read as amended. |
| REV 8's unamended sections | Still flagged; unchanged by this revision. |
| Tallies | Unchanged: 57 rules, 118 requirements, same statuses. |
| Implementation, testing | Separate gates. Neither requested. |

---

## 7. Implementation list — deltas only

REV 11 §8 stands, with these rows replaced or added. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17h″ (REV 11 17h′) | `engine.py` | `_capture(...)` | Follows A1–A11 explicitly; A3 takes the whole open group store; A10 returns a classified failure |
| 17k | `engine.py` | `_resolve_members(refs)` | Total, three-valued map; `unresolved(record_absent)` for a missing record; never omits a key |
| 17l | `engine.py` | `_relevant_groups(candidates, returned_identities, resolution)` | Relevance applied after resolution; no registry or record-store access |
| 17m | `errors.py` | refusal classes | Two shapes sharing one envelope: `Refusal` (Class A) and `CaptureRefusal` (Class B); both carry only `code` |
| 17j′ (REV 11 17j) | `engine.py` | `_relevant_groups` | Superseded by 17l |
| 6⁷ (REV 11 6⁶) | `retrieval.py` | `world` | Propagates the refusing class unchanged; no field added to either |
| 18⁶ (REV 11 18⁵) | `retrieval.py` | `belief_facts` | Same |
| 9⁷ (REV 11 9⁶) | `engine.py` | `propositions(...)` | Same; gate remains A2, before every capture step |

REV 11's `_version`, `group_establishes`, `_accepted`, `_proposition_state`, `authoritative` and every other
row are unchanged. **D-A and D-B remain confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical), `eval/mem1/model.py`, the G1 harness, production Howl
Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–11.

---

## 8. Coherence

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14 | **C** |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** |
| Underspecification is not conflict | 8, 34 | K15 | **C** — `unresolved(record_absent)` is a divergence, reported as one, not as a dispute |
| No adjudication by the engine | 34 | K15, O-10 | **C** |
| History preserved, not re-reported as live | 20, 44 | D-02; K22 | **C** |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — Class B unreachable pre-gate (S-10) |
| Canon reads add no text surface | 18, 32, 47 | REV 4 §4.2; L2; H7 | **C** — canon relevance uses a frozen identity set derived from `accepted` |
| Fail closed | 8 | F4 | **C** — no admissible capture, no answer; no partial map, no inferred absence |
| Engine states what it did | 8 | K15 | **C** — refusal classes distinguish "may not ask" from "could not capture" |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded; R-B, R-C, R-D open.

---

## 9. What this record does and does not claim

**Still open, carried forward unchanged:** B3; belief-voiced dialogue; replay / F7 — S-7 and A10 make a
*response* reviewable against a version and a *failure* attributable to an attempt, not the *engine*
replayable; R-B, R-C, R-D.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- No concurrency safety. The current store is single-threaded and read-only on these paths, so A10's retry
  path is unreachable today. **That is vacuity, not evidence**, and this record repeats REV 11's refusal to
  treat an unexecuted scenario as verified behaviour. Class B is specified, not demonstrated.
- No claim about how `supersede` currently publishes its two writes; nothing was executed to determine it.
- No new rule, requirement, state value or capability. `inconsistent_snapshot` remains a proposed refusal
  code, tally-neutral and unreachable today; the refusal envelope gains a class distinction in the
  implementation, not a field in the response.
- No performance requirement. §3.1's cost is an observation; an index is permitted and changes no set.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9 (`698254ca…`), the
  map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–11 (`98447423…`, `74ac6f77…`, `de4f68b6…`,
  `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`, `5b3a9008…`, `7d703b77…`, `40850419…`).
- No transfer, destination interface, sharing, or cross-chat carryover.
- No change to the REV 10 §6 ratification record or its flagged reading.

**Sequence from here:** owner assessment of REV 12 → acceptance of the snapshot contract as REV 12 §2–§4
with REV 11 §3.2 as amended → separate implementation authorization → implementation → conformance review →
separate testing authorization. This record requests neither authorization.
