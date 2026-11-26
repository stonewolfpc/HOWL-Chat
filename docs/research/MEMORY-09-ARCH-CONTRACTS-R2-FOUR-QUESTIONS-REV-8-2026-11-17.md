# MEMORY-09 — Architectural contracts, revision 8

**Status: PROPOSAL REV 8 — group participation is established against the accepted record set, not the
identity set alone; `relation` is defined as assertion over the records in the answer. Not implemented,
not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 7 (`a281f3e5…`) **for §3.1's `group_contained`, §3.3's `_proposition_state` signature, §4's `open_groups`/`relation` rows, §6 items 5, 6 and 9, §7's third row and §9 rows 6″, 17d, 17e only** |
| Stands unamended | REV 7 in every other respect; REV 1–6 preserved byte-identical |
| Authority | Owner's two-boundary instruction, recorded verbatim (§1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger |
| Ratification | **Not ratified.** Implementation authorization remains a separate gate |

---

## 0. What changed

| Item | REV 7 | REV 8 |
|---|---|---|
| Group containment | Every member's **identity** among the answer's identities | Every member's **record id** in the answer's accepted record set, **and** the identity check kept as a second condition |
| A member that is superseded, revoked or out of scope | Could still establish a relation, and was listed although the answer never returned it | Cannot establish; silent in canon views, reported as history to the authority audience |
| A member whose record id does not resolve | Undefined — the identity check had nothing to compare | Fails closed; integrity code |
| `relation` | `"asserted"` when a group "spans returned entries" — undefined for a group inside one identity | `"asserted"` iff ≥1 group establishes over records in the answer; biconditional with `open_groups` and with `conflict` |
| `conflict` + `relation: not_established` | Possible for a single-identity conflict | **Impossible by construction** |
| Accepted set | Computed inline at each use | Computed once per call and passed to every consumer |

No rule, requirement or tally changes. The 57/118 codes are untouched.

---

## 1. Owner instruction, recorded verbatim

> "**1. Group containment needs to establish record visibility, not just identity presence.** …
> `group_contained`: every member's proposition identity is in `accepted_ids`. `group_establishes`: both
> conditions hold. That is not quite enough to establish that every group member is actually part of the
> answer. The eligibility predicate also checks `truth_apt`, ownership, domain, target, and active status.
> A group member could have an accepted kind and an identity represented in the answer, yet not itself be
> an eligible record returned by that answer."

> "For example, an active eligible record and an inactive record share an identity. If an open group
> references both, identity containment alone could allow the inactive record to influence the active
> answer."

> "Recommendation: Define containment against the exact accepted record-ID set, not only the identity set.
> A group establishes a relation in a view only when every member is among the records accepted for that
> view and applicability. Keep the identity-set check as an additional condition. This makes the statement
> in §4 that all listed group members are already visible an actual invariant rather than an assumption."

> "**2. Clarify `relation` for conflicts within one proposition.** … What happens when an establishing
> group relates two records belonging to the same proposition? The proposition is in conflict, but the
> group does not span multiple proposition entries. … `status: conflict` / `relation: not_established`. The
> two fields would appear to contradict one another."

> "Recommendation: Clarify whether `relation` means: an incompatibility relation has been explicitly
> asserted among records in the answer; or a relation has been asserted between distinct proposition
> identities. Either interpretation can work. The field name and current wording need to agree."

> "I recommend a narrowly scoped REV 8 addressing only those two points. Preserve REV 1–7 and all ratified
> sources. No code changes, tests, new requirements, or changes to the 57/118 tallies."

**Scope:** these two boundaries only. No architecture change, no ratification, no implementation.

---

## 2. Defect 1, stated exactly

REV 7 §3.1 read:

> "`group_establishes(g, view, ids) = group_accepted(g, view) and group_contained(g, ids)`"
>
> "`group_contained(g, ids)` = every member's proposition_id is in ids"

and REV 7 §4 then claimed:

> "`open_groups` | Only groups with `group_establishes(...)` true. Every id inside is already visible in
> this view, so listing adds no disclosure"

The second sentence was **false** under the first. Identity containment proves that a member's proposition
is *represented* in the answer; it says nothing about whether the member record itself survives
`_eligible_truth` — `truth_apt`, `owner_character is None`, `domain`, `target`, and `status == ACTIVE`.

**Trace, entirely within the proposed contract.**

| Step | State |
|---|---|
| 1 | `prop-00003` ("the king is alive") holds `truth-00012` (active) and `truth-00014` (superseded at turn 9) |
| 2 | At turn 6 the operator asserted `group-00004` between `truth-00012` and `truth-00014` — two live records on one identity, both `claim_about_world`, a legal same-kind assertion even under REV 7 §3.4 |
| 3 | `world(world, "king.status")` → `accepted` contains `truth-00012` only; `accepted_ids = {prop-00003}` |
| 4 | `group_accepted` — true (both members are `claim_about_world`) |
| 5 | `group_contained` — true (`truth-00014`'s identity is `prop-00003`, which is in `accepted_ids`) |
| 6 | `group_establishes` — **true**, so `prop-00003` reports `conflict` |
| 7 | `open_groups` lists `group-00004`, whose member `truth-00014` the answer never returned |

Three failures in one outcome, and none of them is a coding slip — all three are what the text says:

- **State.** A record the view cannot see changed the view's verdict, which is precisely the class of leak
  §3.1 of REV 7 existed to close; identity containment merely moved the boundary from "kind" to "record".
- **Disclosure.** Step 7 emits a record id — and, in the authority view, a text — outside the answer's
  accepted set, contradicting REV 7 invariant 9 and REV 4's L2 rule that enumeration adds no text surface.
- **Unresolvable answer.** REV 6 §4 requires that a caller be able to see what is established and what is
  disputed. A `conflict` whose second member is not in the answer cannot be examined, and the caller has no
  way to learn that the missing evidence is historical rather than withheld.

The root cause is the same one REV 6 and REV 7 each corrected one level down: the test was written against
a *set that is convenient*, not against *the set the answer is actually made of*.

---

## 3. Participation against the accepted set

### 3.1 The accepted set is the unit of participation

```
accepted = [r for r in _eligible_truth(domain, target) if _accepts(view, kind_of(r))]
accepted_record_ids = {r.record_id    for r in accepted}
accepted_prop_ids   = {r.proposition_id for r in accepted}
```

| Set | Proves | Does not prove |
|---|---|---|
| `accepted_record_ids` | This record is in the answer: truth-apt, domain-owned, same domain, same target, active, kind accepted | — |
| `accepted_prop_ids` | A proposition is represented in the answer | That any particular record of it is present |
| Group membership (record ids) | What the assertion names | Whether the named records are current or even resolvable |

`accepted` is computed **once** per call and passed to every consumer. A second derivation is a second
opportunity for the two to diverge, and divergence is exactly what §3.3 exists to catch.

### 3.2 Three conditions, not two

```
group_accepted(g, view)        = all(_accepts(view, kind_of(m))   for m in g.members)
group_visible(g, accepted_ids) = all(m in accepted_ids            for m in g.members)   # NEW, record-level
group_contained(g, prop_ids)   = all(proposition_of(m) in prop_ids for m in g.members)
group_establishes(g, view, accepted) =
        group_accepted(g, view) and group_visible(g, ids(accepted)) and group_contained(g, props(accepted))
```

| Condition | Catches | Failure mode if omitted |
|---|---|---|
| `group_accepted` | A relation answering a different question | REV 6's defect — a belief relation in a world answer |
| `group_visible` | A relation resting on a record this answer did not return | This revision's defect — invisible evidence deciding a visible verdict, and its id disclosed |
| `group_contained` | A relation reaching an applicability this answer did not cover | A cross-target relation; also a registry/record-store divergence |

`group_contained` is now logically implied by `group_visible`, because a member's proposition is read from
the record and that record is in `accepted`. It is kept deliberately, as the owner asked, for two reasons:
it states the applicability claim independently of how membership is resolved, and it fails closed if a
future implementation supplies the two sets from different sources. A redundant check that costs one set
lookup and catches a divergence is worth keeping; a redundant check that hides which condition failed is
not — so the envelope reports the failing condition by name (§5).

`proposition_of(m)` is a lookup in the record store. A member whose id does not resolve yields no
proposition, so **both** `group_visible` and `group_contained` fail. That is the fail-closed answer, and it
is the only defensible one: a group naming a record that does not exist cannot be evaluated at all.

### 3.3 Non-establishing causes, and which are faults

REV 7 §3.1 collapsed every non-establishing group into one row ("straddling"). With visibility added, that
row must split, because three different things now prevent establishment and only two of them are faults.

| Cause | What it means | Fault? | Canon view | Authority view |
|---|---|---|---|---|
| Kind straddle | One relation claims two questions | **Integrity fault** | `integrity: ["group_spans_query_kinds"]`, no ids | `open_groups` entry, `establishes: false`, `cause: "spans_query_kinds"` |
| Member not current or not in scope | Superseded, revoked, or another applicability | **Not a fault** — this is history (rules 20, 44) | Silent | `open_groups` entry, `establishes: false`, `cause: "member_not_accepted"` |
| Member id unresolved | The registry and record store have diverged | **Integrity fault** | `integrity: ["group_member_unresolved"]`, no ids | same, `cause: "member_unresolved"` |

The middle row is the important one and it is new work, not a restatement. A conflict asserted at turn 6
between two then-live records stops being reportable in a view once one member is superseded — not because
the assertion was withdrawn, but because that view can no longer see half of it. Reporting `conflict` on
the strength of a retired record would assert a dispute about evidence the caller cannot inspect. The
assertion is preserved, is still listed to the audience entitled to the whole graph, and returns to
governing state if a human reactivates the record or restates the assertion.

This also makes the three-way relation of REV 7 §3.1 exactly right in shape and merely under-specified in
its first row: wholly in view establishes, wholly out is silent, and *partly* in is now split by cause
instead of assumed to be one fault.

### 3.4 State, revised again

```
_proposition_state(prop_id, *, view, accepted) ->
    records = [r for r in accepted if r.proposition_id == prop_id]
    groups  = [g for g in self._open_groups if prop_id in propositions_of(g)]
    live    = [g for g in groups if group_establishes(g, view, accepted)]

    if live:                    "conflict"
    elif len(records) >= 2:     "unresolved_relation"
    elif len(records) == 1:     "settled"
    else:                       "none"
```

| Property | Consequence |
|---|---|
| `accepted` replaces `accepted_ids` | Record-level and identity-level conditions are both derivable from one value; no call site can pass the weaker one |
| `records` is read from `accepted` | The count and the containment test cannot disagree about what is in the answer |
| Conflict before counts | An established relation is still never downgraded |
| Non-establishing groups tested nowhere | They cannot force `conflict`, and cannot force `unresolved_relation` either |

---

## 4. `relation`, defined

**Decision.** `relation` states whether an incompatibility relation has been **explicitly asserted among
the records in this answer**. It is not a statement about identity topology.

```
"asserted"        iff at least one open group establishes over records in this answer
"not_established" otherwise
```

| Consequence | Effect |
|---|---|
| Single identity, two live records, establishing group | `status: "conflict"`, `relation: "asserted"` — the two fields now agree |
| Several identities, establishing group across them | `status: "conflict"`, `relation: "asserted"` |
| Several identities, no establishing group | `status: "multiple"`, `relation: "not_established"` |
| One identity, two live records, no group | `status: "unresolved_relation"`, `relation: "not_established"` |
| No accepted identity | `status: "none"`, `relation: "not_established"` |

| Biconditional | Why it holds |
|---|---|
| `relation == "asserted"` ⟺ `open_groups` non-empty | `open_groups` contains exactly the establishing groups |
| `relation == "asserted"` ⟺ some entry has `state == "conflict"` | Every establishing group's members are in `accepted`, so at least one entry is in conflict; and a conflicting entry has an establishing group by §3.4 |
| `relation == "not_established"` ⟹ no entry is `conflict` | The contrapositive of the second |

These are checkable, not aspirational: `conflict` with `relation: "not_established"` is now a representable
impossibility rather than an edge case left to the implementation.

### 4.1 Why this meaning and not the identity-spanning one

| Consideration | "Among records in the answer" (chosen) | "Between distinct identities" (rejected) |
|---|---|---|
| Reads correctly for a single-identity conflict | `conflict` + `asserted` — consistent | `conflict` + `not_established` — reads as "nothing was asserted", the opposite of the record |
| What the field names | The assertion actually made | A structural property of the registry |
| What the answer can decide | Whether a group establishes here — already computed | Whether that group's member identities differ — the same computation, used as a proxy |
| Cost | None; `open_groups` already exposes membership | Needs a third value or a second field to say "asserted, within one identity" |

The last row is decisive on its own: identity scope is **already derivable** from `open_groups`, whose
members are visible by §3.2. A caller that needs to know whether a dispute crosses identities counts
distinct `proposition_id`s in the group it can already see. Encoding that in `relation` would make the
field carry two facts and would force the confusing single-identity result above. No new top-level field is
introduced; the envelope is the one REV 6 defined with the two keys REV 7 added.

### 4.2 Envelope, and where the detail lives

```
engine.world(world, "king.status")
→ { "status": "conflict", "relation": "asserted", "view": "world",
    "open_groups": [ { "group_id": "group-00004",
                       "members": ["truth-00012", "truth-00018"],
                       "propositions": ["prop-00003", "prop-00003"] } ],
    "integrity": [], "propositions": [ … ] }

engine.propositions(world, target="king.status", actor="operator")
→ { "status": "ok", "viewer": "authority", "knowledge_conferred": false,
    "open_groups": [ { "group_id": "group-00004", "establishes": true,  … },
                     { "group_id": "group-00007", "establishes": false,
                       "cause": "member_not_accepted", "members": [ … ] } ],
    "integrity": [], "propositions": [ … ] }
```

| Field | Canon views | Authority view |
|---|---|---|
| `open_groups` | Establishing groups only — every member is in `accepted`, so membership is already visible | Every open group touching a returned identity |
| `establishes` / `cause` | Absent (all are `true` by construction) | Present; `cause` only when `false` |
| `integrity` | Codes only, no ids, no texts, no counts | Same codes, `group_id` attached |

The two `propositions` entries in the first example carry the same id because the group is internal to one
identity — which is the case §4 exists to make legible, not to flatten. `propositions` is the member list,
so a caller distinguishes a single-identity dispute from a cross-identity one without a dedicated field.

---

## 5. Effect on the three views

| View | Identity set source | `open_groups` | `integrity` |
|---|---|---|---|
| `world` | `accepted` records only | Establishing groups | Codes only |
| `belief_facts` | `accepted` records only | Establishing groups | Codes only |
| `propositions` (gated) | **The registry**, target-scoped — this is the capability REV 4 protects | Every group touching a returned identity, with `establishes`/`cause` | Codes plus `group_id` |

The third row is not an inconsistency, and REV 7 invariant 9 must be read with it in view: canon views
derive identities from accepted records so that they cannot become enumeration surfaces by accident; the
authority view derives them from the registry because listing an identity whose records are all superseded
(`state: none`, `records: []`, REV 4 §7) is the adjudication service REV 4 authorises. The gate is what
makes that difference legitimate rather than a leak, which is why §5 of REV 7 put it at the callable
boundary.

---

## 6. Invariants

Per view, for a caller permitted that view:

1. Every returned record satisfies `_accepts(view, kind_of(record))`.
2. Every returned record's `proposition_id` equals its entry's id.
3. No record appears in two entries — single-valued `proposition_id`.
4. The union of entries' records equals the view-accepted eligible set.
5. **Participation:** an entry's `state` is a function of its own accepted records and of open groups that
   `group_establishes(...)` in this view. No other input affects it.
6. **No absorbable fault or history:** a non-establishing group changes no `state` and no `relation`, and
   appears in a canon view only as an `integrity` code.
7. **Complement:** `world` records ∪ `belief_facts` records = `_eligible_truth(domain, target)`; the
   intersection is empty.
8. **Gate:** every `propositions`-view call returning a non-refusal has completed
   `require(domain, actor, "promotion")` before its first data read; every refusal contains no
   registry-derived field.
9. **Scope, corrected:** every record id, group member, text and count in a canon-view response belongs to
   a record in `accepted`. No member of a listed group lies outside it.
10. **Kind:** every record of an identity carries that identity's kind.
11. **Relation:** `relation == "asserted"` ⟺ `open_groups` non-empty ⟺ some entry is `conflict`.

Item 9 is the sentence REV 7 §4 asserted without support; it is now implied by `group_visible`, which is
why it is stated as an invariant rather than left as a claim in prose. Item 11 is new. Item 6's wording
widens from "fault" to "fault or history" because §3.3 distinguishes them.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 7. Worked examples

**A. The inactive member — REV 7 versus REV 8.** State as §2 steps 1–2.

| | REV 7 | REV 8 |
|---|---|---|
| `prop-00003` state | `conflict` | `settled` |
| `relation` | `asserted` | `not_established` |
| `open_groups` | `group-00004`, member `truth-00014` — never returned | `[]` |
| `integrity` | `[]` | `[]` (history is not a fault) |
| `propositions` (operator) | not distinguished | lists `group-00004`, `establishes: false`, `cause: "member_not_accepted"` |

The operator still sees the assertion and can see why it is not governing. The world answer no longer
asserts a dispute on evidence it declined to produce.

**B. Single-identity conflict, cross-identity comparison.**

| State | `status` | `relation` | `open_groups` members |
|---|---|---|---|
| One identity, two live records, group `g1` between them | `conflict` | `asserted` | `["truth-00012", "truth-00014"]`, one distinct `proposition_id` |
| Two identities, one live record each, group `g2` across them | `conflict` | `asserted` | `["truth-00012", "truth-00019"]`, two distinct ids |
| Same as the second, group never asserted | `multiple` | `not_established` | `[]` |

Both disputes now read the same way at top level, which is correct: the difference between them is a
property of the group, and the group is inspectable.

**C. The refusal boundary, unchanged from REV 7 §5.**

| Call | Result |
|---|---|
| `authoritative(..., view="propositions")`, no grant | `refused / no_domain_authority`; no ids, no counts |
| `authoritative(..., view="world")`, no grant | Served; no gate, unchanged |
| `authoritative(..., view="bogus")` | `refused / unknown_view` |

---

## 8. Coherence of the two corrections

A check that the corrections do not contradict the ratified set, not conformance evidence.

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14 | **C** — kind, visibility and identity all required |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** — unchanged from REV 7 |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** — establishment now means "on records this answer contains" |
| Underspecification is not conflict | 8, 34 | K15 | **C** — `multiple` unchanged |
| No adjudication by the engine | 34 | K15, O-10 | **C** — a retired assertion is reported, not resolved |
| History is preserved, not re-reported as live | 20, 44 | D-02; K22 | **C** — superseded member: assertion kept, state not governed |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — gate unchanged |
| Canon reads add no text surface | 18, 32, 47 | REV 4 §4.2; L2; H7 | **C** — group members now inside `accepted` |
| Fail closed | 8 | F4 | **C** — unresolved member fails both containment tests |
| Engine states what it did | 8 | K15 | **C** — `relation` names the assertion made |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded and audience-constrained; R-B,
R-C, R-D open.

---

## 9. Implementation list — deltas only

REV 7 §9 stands, with these rows replaced. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17f | `engine.py` | `_accepted(domain, target, view)` | Computes the accepted record list once; every consumer takes it as an argument |
| 17d′ (REV 7 17d) | `engine.py` | `group_establishes(group, view, accepted)` | `group_accepted` + `group_visible` (record ids) + `group_contained` (identities); unresolved member fails both |
| 17e′ (REV 7 17e) | `engine.py` | `_proposition_state(prop_id, *, view, accepted)` | Records read from `accepted`; establishing groups only; conflict before counts |
| 6‴ (REV 7 6″) | `retrieval.py` | `world` | `relation` per §4; `open_groups` establishing-only, all members inside `accepted`; `integrity` codes only |
| 18″ (REV 7 18′) | `retrieval.py` | `belief_facts` | Same rules as `world` |
| 9‴″ (REV 7 9‴′) | `engine.py` | `propositions(...)` | Registry-sourced identity set, gated; `open_groups` includes non-establishing groups with `establishes`/`cause` |
| 17c″ (REV 7 17c′) | `engine.py` | `authoritative(...)` | Unchanged except that it computes `accepted` once and passes it; gate order per REV 7 §5.1 |

REV 7's `assert_incompatible` cross-kind refusal, `supersede`/D-A, `retract`/D-B and every other row are
unchanged. **D-A and D-B remain confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical), `eval/mem1/model.py`, the G1 harness, production
Howl Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–7.

---

## 10. What this record does and does not claim

**Still open, carried forward unchanged:** B3 (event versus definition); belief-voiced dialogue; replay /
F7 (no facility, determinism unproven — the only determinism claim remains ordering by `proposition_id`
then `record_id`); R-B, R-C, R-D. None is a reason to expand this revision.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- No behavioural verification. Statements about current behaviour remain source inspection from R1/R2.
- No conformance proof. §8 checks coherence only.
- No new rule or requirement. `cause` is a diagnostic field of a proposed envelope; `relation`'s meaning is
  a definition of an existing field, not a new state; neither adds a code to the map.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9
  (`698254ca…`), the map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–7 (`98447423…`, `74ac6f77…`,
  `de4f68b6…`, `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`).
- No transfer, destination interface, sharing, or cross-chat carryover.

**Sequence from here:** owner assessment of the complete proposal (REV 1–8) for ratification → separate
implementation authorization → implementation → conformance review → separate testing authorization. This
record does not request ratification and does not request implementation.
