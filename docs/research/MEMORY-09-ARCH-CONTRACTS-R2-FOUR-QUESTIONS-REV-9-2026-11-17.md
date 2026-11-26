# MEMORY-09 — Architectural contracts, revision 9

**Status: PROPOSAL REV 9 — `open_groups` and `establishing_groups` are separated, and `relation` is defined
against the establishing set in every view. Not implemented, not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 8 (`b38f1062…`) **for §4's `relation` definition and biconditionals, §4.2's envelope table, §5's view table, §6 items 11, §7 examples A and B, §9 rows 6‴, 18″, 9‴″ only** |
| Stands unamended | REV 8 in every other respect; REV 1–7 preserved byte-identical |
| Authority | Owner's naming-and-invariant instruction, recorded verbatim (§1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger |
| Ratification | **Not ratified.** Implementation authorization remains a separate gate |

---

## 0. What changed

| Item | REV 8 | REV 9 |
|---|---|---|
| The group collection | One key, `open_groups`, holding two different sets depending on audience | Two keys: `open_groups` (all groups this view may show) and `establishing_groups` (the subset where `group_establishes` holds) |
| `relation` | Defined against "open groups" — the ambiguous collection | Defined against `establishing_groups` in every view |
| The biconditional | `relation == "asserted"` ⟺ `open_groups` non-empty ⟺ some entry is `conflict` — false for the authority view | `relation == "asserted"` ⟺ `establishing_groups` non-empty ⟺ some entry is `conflict`, in every view |
| Canon views | `open_groups` = establishing only | `open_groups` = `establishing_groups` (identical lists, stated as invariant 12) |
| Authority view | `open_groups` = establishing plus historical, distinguished only by a per-entry flag | `open_groups` = all relevant, `establishing_groups` = the governing subset; non-emptiness of `open_groups` implies nothing |
| Key set | View-dependent meaning | Uniform across views — one key, one meaning |

No rule, requirement or tally changes. Nothing is added to the ratified set.

---

## 1. Owner instruction, recorded verbatim

> "**The `open_groups` biconditional does not hold for the authority view.** … For example: An identity has
> one active record and one superseded record. An open group references both. The authority view lists the
> group with `establishes: false` and `cause: "member_not_accepted"`. No proposition is in conflict because
> the group does not establish over the accepted records. The authority view therefore has a non-empty
> `open_groups`, but `relation` must be `not_established`."

> "This is not a problem with preserving historical assertions. It is a problem with using the same
> collection name for two different sets."

> "**Recommended correction.** Keep the behavior. Clarify the distinction between: `open_groups`: all
> relevant open groups visible to the authority audience; `establishing_groups`: the subset that actually
> establishes over the accepted records. For canon views, `open_groups` can continue to contain only
> establishing groups. Define `relation` against `establishing_groups`, or equivalently against the set of
> groups for which `group_establishes(...)` is true. Then state the invariant precisely: `relation ==
> "asserted"` iff at least one group establishes in the current view. In canon views, `open_groups` contains
> exactly those groups. In the authority view, `open_groups` may also contain non-establishing historical
> groups, so its non-emptiness does not imply conflict."

> "I recommend a narrowly scoped REV 9 correcting only the `open_groups` / `establishing_groups`
> distinction and the scope of the biconditional. Preserve REV 1–8 and all ratified sources."

**Scope:** the naming and the invariant only. No architecture change, no ratification, no implementation.

---

## 2. The defect, stated exactly

REV 8 §6 item 11 read:

> "**Relation:** `relation == "asserted"` ⟺ `open_groups` non-empty ⟺ some entry is `conflict`."

and REV 8 §4.2 gave `open_groups` two definitions:

| View | REV 8 §4.2 |
|---|---|
| Canon | "Establishing groups only — every member is in `accepted`" |
| Authority | "Every open group touching a returned identity" |

Three statements, not two of which can all be true. The owner's counter-example is minimal and complete:

| Step | State |
|---|---|
| 1 | `prop-00003` holds `truth-00012` (active) and `truth-00014` (superseded at turn 9) |
| 2 | `group-00004` names both; asserted at turn 6, both then live |
| 3 | `propositions(world, target="king.status", actor="operator")` — authority view, so `truth-00014` is listed as a record of the identity even though it is retired |
| 4 | `group_visible` fails: `truth-00014 ∉ accepted_record_ids` |
| 5 | `group_establishes` false → `prop-00003` state is `settled`, not `conflict` |
| 6 | Authority `open_groups` lists `group-00004` with `establishes: false`, `cause: "member_not_accepted"` |
| 7 | So `open_groups` is non-empty while `relation` is `"not_established"` — the biconditional's first leg fails |

The fault is not behavioural. Steps 1–5 are exactly what REV 8 intends, and step 6 is exactly what the
authority audience needs. The fault is that REV 8 used one name for two sets and then wrote a
cross-view equation about the name. REV 7's version of the same mistake was one predicate for two
questions (§3.1); REV 8's is one collection for two audiences. The correction is the same shape as before:
separate the two things, then state the relation against the one that actually governs.

It is worth naming why this was easy to miss: in every canon view the two sets coincide, so the equation
held in the two views that are cheapest to reason about and failed only in the one that is gated and
therefore least exercised. An invariant that holds in the easy cases is not an invariant.

---

## 3. Two collections, one meaning each

### 3.1 Definitions

```
relevant_groups(view, domain, target) =
    [g for g in open_group_store
       if any(proposition_of(m) in returned_identities(view) for m in g.members)]

establishing(g, view, accepted) = group_establishes(g, view, accepted)      # REV 8 §3.2, unchanged

open_groups(view)         = relevant_groups(view, ...)      # canon: filtered to establishing
establishing_groups(view) = [g for g in open_groups(view) if establishing(g, view, accepted)]
```

| Collection | Means | Audience |
|---|---|---|
| `open_groups` | Every group this view is permitted to show | Both; the permitted set differs |
| `establishing_groups` | The groups that govern `state` in this view | Both; the set is computed identically |

`open_group_store` means every group not marked closed. No operation in v9 or R2 marks a group closed — the
contract defines no engine-initiated closure, because closing an assertion is adjudication (rule 34). If a
closure operation is later authorised, `open_groups` means the not-closed set and nothing else in this
record changes.

### 3.2 Canon views: the two collections coincide

```
open_groups(world)         == establishing_groups(world)
open_groups(belief_facts)  == establishing_groups(belief_facts)
```

Both are literal, identical lists. This is not redundancy for its own sake: it is what lets every invariant
be written once, for every view, with no audience condition attached. A field whose meaning is silently
narrowed by view is the defect this revision exists to correct, and reintroducing it under a second name
would be the same error with a better label.

Canon disclosure is unchanged: a group with a non-accepted member affects no state and appears in no canon
collection; it surfaces only as an `integrity` code, with no id, text or count (REV 7 §4, D-02 behaviour 8).

### 3.3 Authority view: the two collections differ, and that is the point

```
establishing_groups(propositions) ⊆ open_groups(propositions)
open_groups(propositions) non-empty  ⇏  any state is "conflict"
establishing_groups(propositions) non-empty  ⟺  some entry state is "conflict"
```

The authority `open_groups` contains every group naming at least one record whose proposition is among the
returned identities, which includes the cases the canon views must stay silent about:

| Group | In authority `open_groups` | In `establishing_groups` | `cause` |
|---|---|---|---|
| All members current and kind-accepted | yes | yes | — |
| One member superseded or revoked | yes | no | `member_not_accepted` |
| Member in another domain or target | yes | no | `member_not_accepted` |
| Members of differing kinds | yes | no | `spans_query_kinds` |
| Member id unresolved | yes | no | `member_unresolved` |
| Every member retired | yes — the identity is still returned | no | `member_not_accepted` |

The last row is why the authority collection is defined against *returned identities* rather than returned
records: an identity whose records are all superseded is still enumerated (`state: none`, `records: []`,
REV 4 §7), and the assertions made about it are precisely what an adjudicating audience needs to see. A
group is therefore "relevant" to the authority view if it names any record of any returned identity,
retired or not.

Faults appear twice to this audience, deliberately: a `group_id` in `integrity`, and a full entry in
`open_groups` with `establishes: false` and `cause`. The code is the machine-readable signal; the entry is
the diagnosis. Canon views get only the code.

### 3.4 Naming the state source

`_proposition_state` (REV 8 §3.4) already filters: `live = [g for g in groups if group_establishes(...)]`.
REV 9 changes no logic there; it changes only what the envelope publishes and which set `relation` is
defined on. The state function and `establishing_groups` must be derived from the same predicate — one
function, called once per group, memoised per call — so that a conflict entry and a governing group entry
cannot disagree.

---

## 4. `relation`, defined against the establishing set

**Decision.** `relation` states whether an incompatibility relation has been explicitly asserted among the
records in this answer. REV 8's meaning stands; only its operand changes.

```
"asserted"        iff establishing_groups(view) is non-empty
"not_established" otherwise
```

| Biconditional | Holds in | Why |
|---|---|---|
| `relation == "asserted"` ⟺ `establishing_groups` non-empty | **every view** | This is the definition |
| `relation == "asserted"` ⟺ some entry has `state == "conflict"` | **every view** | Establishing members are in `accepted`, so some entry is in conflict; a conflicting entry has an establishing group (REV 8 §3.4) |
| `relation == "asserted"` ⟹ `open_groups` non-empty | every view | `establishing_groups ⊆ open_groups` |
| `open_groups` non-empty ⟹ `relation == "asserted"` | **canon only** | In the authority view this implication is false by design (§3.3) |
| `relation == "asserted"` ⟺ `open_groups` non-empty | **canon only** | Same; stated as a canon-specific consequence, not as the general equation |

REV 8's error was writing the last line as though it were the general case. The general case is the first
line, and it now holds where it is written.

### 4.1 Envelope

```
engine.world(world, "king.status")
→ { "status": "conflict", "relation": "asserted", "view": "world",
    "open_groups":         [ { "group_id": "group-00004", "members": ["truth-00012","truth-00018"],
                               "propositions": ["prop-00003","prop-00003"] } ],
    "establishing_groups": [ { "group_id": "group-00004", "members": ["truth-00012","truth-00018"],
                               "propositions": ["prop-00003","prop-00003"] } ],
    "integrity": [], "propositions": [ … ] }

engine.propositions(world, target="king.status", actor="operator")
→ { "status": "ok", "viewer": "authority", "knowledge_conferred": false,
    "open_groups":         [ { "group_id": "group-00004", "establishes": false,
                               "cause": "member_not_accepted",
                               "members": ["truth-00012","truth-00014"],
                               "propositions": ["prop-00003","prop-00003"] },
                             { "group_id": "group-00007", "establishes": true, … } ],
    "establishing_groups": [ { "group_id": "group-00007", … } ],
    "integrity": [], "propositions": [ … ] }
```

The second example is the owner's counter-example made concrete: `relation` is `"not_established"`, one
entry is `conflict` only if `group-00007` establishes on its own records, and `open_groups` is non-empty
without implying anything about conflict.

| Field | Canon | Authority |
|---|---|---|
| `open_groups` | = `establishing_groups` | everything relevant; the only place `establishes`/`cause` appear |
| `establishing_groups` | identical to `open_groups` | governing subset |
| `relation` | from `establishing_groups` | from `establishing_groups` |
| `integrity` | codes only | codes with `group_id` |

`establishes` is omitted from canon entries because it is always `true` there; in the authority view it is
present on every `open_groups` entry, and `cause` appears only when it is `false`. `establishing_groups`
entries never carry either field — membership in that collection is the assertion.

---

## 5. Effect on the three views

| View | Identity set | `open_groups` | `establishing_groups` | `relation` from |
|---|---|---|---|---|
| `world` | accepted records only | establishing only | identical | establishing |
| `belief_facts` | accepted records only | establishing only | identical | establishing |
| `propositions` (gated) | registry, target-scoped | all relevant | governing subset | establishing |

The third row's identity source is the capability REV 4 protects: listing an identity whose records are all
superseded is the adjudication service, and the gate at the callable boundary is what makes it legitimate.
REV 8 §5's reading of invariant 9 stands unchanged.

---

## 6. Invariants

Per view, for a caller permitted that view:

1. Every returned record satisfies `_accepts(view, kind_of(record))`.
2. Every returned record's `proposition_id` equals its entry's id.
3. No record appears in two entries — single-valued `proposition_id`.
4. The union of entries' records equals the view-accepted eligible set.
5. **Participation:** an entry's `state` is a function of its own accepted records and of open groups that
   `group_establishes(...)` in this view. No other input affects it.
6. **No absorbable fault or history:** a non-establishing group changes no `state` and no `relation`.
7. **Complement:** `world` records ∪ `belief_facts` records = `_eligible_truth(domain, target)`; the
   intersection is empty.
8. **Gate:** every `propositions`-view call returning a non-refusal has completed
   `require(domain, actor, "promotion")` before its first data read; every refusal contains no
   registry-derived field.
9. **Scope:** every record id, group member, text and count in a canon-view response belongs to a record in
   `accepted`. No member of a listed group lies outside it.
10. **Kind:** every record of an identity carries that identity's kind.
11. **Relation:** `relation == "asserted"` ⟺ `establishing_groups` non-empty ⟺ some entry is `conflict`.
    In every view, without audience conditions.
12. **Collection identity:** in canon views `open_groups` and `establishing_groups` are equal lists; in the
    authority view `establishing_groups ⊆ open_groups`, and non-emptiness of `open_groups` implies nothing
    about `relation` or about any `state`.

Item 11 no longer mentions `open_groups`, which is what made it false in REV 8. Item 12 is new and is the
statement the owner asked for: it fixes each collection's meaning and confines the one-way implication to
the audience where it holds.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 7. Worked examples

**A. The retired member, all three views.**

| | `world` | `belief_facts` | `propositions` (operator) |
|---|---|---|---|
| `prop-00003` state | `settled` | `none` | `settled` |
| `relation` | `not_established` | `not_established` | `not_established` |
| `open_groups` | `[]` | `[]` | `[group-00004: establishes false, cause member_not_accepted]` |
| `establishing_groups` | `[]` | `[]` | `[]` |
| `integrity` | `[]` | `[]` | `[]` |

No view reports `conflict`; only the adjudicating audience learns that an assertion exists and why it is
not governing. REV 8 already produced this behaviour — REV 9 only makes the envelope say so without a
contradiction.

**B. Governing and historical at once.**

| Call | `relation` | `establishing_groups` | `open_groups` | Conflicts |
|---|---|---|---|---|
| One identity, `g1` across two live records | `asserted` | `[g1]` | `[g1]` | 1 |
| Two identities, `g2` across live records, `g3` retired | `asserted` | `[g2]` | `[g2, g3]` (authority) | 2 |
| Two identities, only `g3` retired | `not_established` | `[]` | `[g3]` (authority) | 0 |
| Two identities, no groups | `not_established` | `[]` | `[]` | 0 |

Row three is the case that broke REV 8's equation; row four is the case that made the equation look true.

**C. Refusal boundary, unchanged from REV 7 §5.**

| Call | Result |
|---|---|
| `authoritative(..., view="propositions")`, no grant | `refused / no_domain_authority`; no ids, no counts |
| `authoritative(..., view="world")`, no grant | Served; no gate |
| `authoritative(..., view="bogus")` | `refused / unknown_view` |

A refusal contains no collection at all, so the two-key envelope adds no surface to it.

---

## 8. Coherence of the correction

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14 | **C** |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** — establishment unchanged; only its publication is split |
| Underspecification is not conflict | 8, 34 | K15 | **C** |
| No adjudication by the engine | 34 | K15, O-10 | **C** — a retired assertion is reported, never resolved or closed |
| History preserved, not re-reported as live | 20, 44 | D-02; K22 | **C** |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — gate unchanged |
| Canon reads add no text surface | 18, 32, 47 | REV 4 §4.2; L2; H7 | **C** — canon `open_groups` still establishing-only |
| Fail closed | 8 | F4 | **C** |
| Engine states what it did | 8 | K15 | **C** — `relation` now names a set the response contains |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded; R-B, R-C, R-D open.

---

## 9. Implementation list — deltas only

REV 8 §9 stands, with these rows replaced. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 6⁗ (REV 8 6‴) | `retrieval.py` | `world` | Emits `open_groups` (= establishing) and `establishing_groups` (identical); `relation` from the latter |
| 18‴ (REV 8 18″) | `retrieval.py` | `belief_facts` | Same as `world` |
| 9⁗ (REV 8 9‴″) | `engine.py` | `propositions(...)` | `open_groups` = all relevant with `establishes`/`cause`; `establishing_groups` = governing subset; `relation` from the latter |
| 17g | `engine.py` | `_relevant_groups(view, …)` | Derives both collections from one pass; `establishing` computed once per group per call |
| 17d″ | `engine.py` | `group_establishes(...)` | Unchanged from REV 8; called through one accessor so state and `establishing_groups` cannot diverge |

REV 8's `_accepted`, `_proposition_state`, `authoritative`, `assert_incompatible` refusal, `supersede`/D-A
and `retract`/D-B rows are unchanged. **D-A and D-B remain confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical), `eval/mem1/model.py`, the G1 harness, production
Howl Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–8.

---

## 10. What this record does and does not claim

**Still open, carried forward unchanged:** B3; belief-voiced dialogue; replay / F7 (no facility,
determinism unproven — ordering by `proposition_id` then `record_id` remains the only determinism claim);
R-B, R-C, R-D. None is a reason to expand this revision.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- No behavioural verification. Statements about current behaviour remain source inspection from R1/R2.
- No conformance proof. §8 checks coherence only.
- No new rule, requirement or state value. `establishing_groups` is a second projection of a set the
  response already contains, not a new capability; `open_groups` keeps its REV 8 authority meaning and its
  canon narrowing, now stated as such.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9
  (`698254ca…`), the map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–8 (`98447423…`, `74ac6f77…`,
  `de4f68b6…`, `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`).
- No transfer, destination interface, sharing, or cross-chat carryover.

**Sequence from here:** owner assessment of the complete proposal (REV 1–9) for ratification → separate
implementation authorization → implementation → conformance review → separate testing authorization. This
record does not request ratification and does not request implementation.
