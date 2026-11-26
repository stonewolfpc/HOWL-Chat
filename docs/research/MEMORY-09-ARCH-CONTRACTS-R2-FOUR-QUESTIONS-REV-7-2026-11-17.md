# MEMORY-09 — Architectural contracts, revision 7

**Status: PROPOSAL REV 7 — conflict participation is relative to the accepted view; the `propositions`
view carries its authority gate at the callable boundary. Not implemented, not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 6 (`41421dce…`) **for §3.4, §4's `open_groups`/`relation`/`integrity` rows, §5's "no actor gate" row, §6 items 5–6 and §8 rows 6′, 9‴, 17b′, 17c only** |
| Stands unamended | REV 6 in every other respect; REV 1–5 preserved byte-identical |
| Authority | Owner's two-boundary instruction, recorded verbatim (§1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger. Corrects a proposal and two unapplied implementation-list rows |
| Ratification | **Not ratified.** Implementation authorization remains a separate gate |

---

## 0. What changed

| Item | REV 6 | REV 7 |
|---|---|---|
| Conflict participation | `_proposition_state(prop_id)` — a group touching the identity marked it `conflict` regardless of kind | **Relative to the accepted view**: only a group whose every member is accepted in that view, and whose members' identities are all in the answer, can affect state |
| Cross-kind groups at creation | Not addressed; could be created | **Refused, fail closed** — `assert_incompatible` will not build one |
| Mixed group found in the registry | Would have marked a world identity `conflict` | **Cannot affect state**; reported as an integrity code with no foreign ids or texts in a canon view |
| `authoritative(..., view="propositions")` | View accepted both kinds; gate lived in the `propositions()` wrapper only | **Gate is an invariant of the view**: refusal precedes any registry read, `actor` is required for that view |
| Canon views | No actor gate | **Unchanged** — no gate added, none removed |
| Response envelope | `open_groups` | `open_groups` (establishing groups only) + `integrity` (codes; details only in the authority view) |

No rule, requirement or tally changes. The 57/118 codes are untouched.

---

## 1. Owner instruction, recorded verbatim

> "**1. Open conflict groups must not cross-contaminate query kinds.** REV 6 makes query kind explicit, but
> §3.4 says `_proposition_state(prop_id)` is view-independent and marks a proposition as conflicting if one
> of its records belongs to an open group. That creates a possible leak … That would allow a belief-fact
> relationship to affect an ordinary world-truth answer, despite the new query-kind boundary."

> "Recommendation: Define conflict participation relative to the accepted view. A proposition's state in a
> view should only be affected by an open group whose relevant members are in that view. If cross-kind
> conflict groups are prohibited instead, state that explicitly and make their creation fail closed. This
> does not change the conflict model. It makes the new query-kind boundary hold consistently."

> "**2. The `propositions` view must not bypass its authority gate.** … If `authoritative(...,
> view=\"propositions\")` is callable without the actor/capability check, the new view could become a second
> route around REV 4's enumeration restriction."

> "Recommendation: Make the authority gate an invariant of the `propositions` view itself, not merely a
> convention imposed by the `propositions()` wrapper. Unauthorized calls must refuse before reading the
> registry, consistent with REV 4."

> "I recommend one narrowly scoped REV 7 addressing only: 1. Conflict-group participation across query-kind
> boundaries. 2. Enforcement of the `propositions` authority gate at the actual callable boundary.
> Preserve REV 1–6 and all ratified sources. No implementation, tests, or changes to the 57/118
> requirements."

**Scope:** these two boundaries only. No architecture change, no ratification, no implementation. The
owner's remaining open items — B3, belief-voiced dialogue, replay — are explicitly **not** reasons to
reopen these contracts and are carried forward unchanged (§10).

---

## 2. Defect 1, stated exactly

REV 6 §3.4 read:

> "`_proposition_state(prop_id)` is unchanged from REV 5 §3.4: `conflict` if one of its records is in an
> open group, else `unresolved_relation` at two or more active records, `settled` at one, `none` at zero.
> State is a property of an identity, not of a query, so it takes no view."

The last sentence was the error, and it was an error of the same family as REV 5's. REV 5 discarded
identity before deciding and reported a relation nobody had asserted; REV 6 restored identity and
filtered records by kind, but then let a relation cross the very boundary the filter had just drawn.

**Trace, entirely within the proposed contract.**

| Step | State |
|---|---|
| 1 | `prop-00003` — one world record `truth-00012`, "the king is alive" |
| 2 | `prop-00011` — one ascription record `truth-00031`, "Sable believes the king was poisoned" |
| 3 | An open group `group-00002` names `truth-00012` and `truth-00031` as incompatible |
| 4 | `world(world, "king.status")` → `_accepts("world", …)` drops `truth-00031` correctly; `prop-00003` survives with one record |
| 5 | `_proposition_state("prop-00003")` sees `group-00002`, which contains a record of `prop-00003`, and returns **`conflict`** |
| 6 | The world answer reports `conflict` because of a relation whose other member the query was constructed to exclude |

Nothing in step 5 is a coding slip. It is what REV 6 §3.4 literally says. And it defeats the invariant
REV 6 §6 item 5 states — "state is computed from the entry's own records and open groups only" — because a
group is not an entry's own unless its members are in the same view.

Two facts made this reachable rather than hypothetical:

- REV 3 §5.4 permits an assertion to "name records from different propositions that share an applicability"
  and says nothing about kind, so a mixed group is a legal authoring act under the text as it stood;
- REV 6 §5 states that `propositions` accepts both kinds, so a mixed group is visible to the authoring
  surface and invisible to both canon surfaces — the worst possible combination, since it would be
  undiagnosable from `world()` and `belief_facts()` while still altering their verdicts.

---

## 3. Conflict participation, made view-relative

### 3.1 Three group relations to a view, not two

A group is not simply "open or closed" relative to a question. It is **wholly in**, **wholly out**, or
**straddling**. Only the last is a fault, and REV 6 collapsed all three.

```
group_kind(g)      -> the set of kinds of g's member records
group_accepted(g, view)   = all(_accepts(view, kind_of(member)) for member in g.members)
group_contained(g, ids)   = every member's proposition_id is in ids
group_establishes(g, view, ids) = group_accepted(g, view) and group_contained(g, ids)
```

| Relation | Meaning | Effect on state | Reported |
|---|---|---|---|
| Wholly in view — `group_accepted` and contained | The relation is about the question being asked | May set `conflict` | `open_groups`, full detail (all members already visible) |
| Wholly out of view — no member accepted | Not this question's business — e.g. a world group in `belief_facts` | **None**, silently | **Not mentioned.** Listing it would be a disclosure, not an answer |
| Straddling — some accepted, some not | An invariant violation: one relation claiming both kinds are one question | **None** | `integrity` code, **no ids, no texts, no counts** in a canon view; full detail in the authority view only |

The middle row is not a loophole. A group whose members are all `claim_about_world` is invisible to
`belief_facts` for exactly the reason a world record is: it answers a different question. Suppressing it
there is correctness, not silence.

The third row is a fault because §3.2's rule — kind is single-valued per proposition — makes it
contradictory: two records asserting incompatibility while answering different questions. The engine may
not resolve that, so it must not let either answer absorb it.

### 3.2 Kind is single-valued per identity

New, and load-bearing for §3.1:

```
every record r with r.proposition_id == pid has kind_of(r) == kind_of_proposition(pid)
```

A proposition's kind is a property of the declared identity (REV 2 §3.1, REV 4 §4.1), so every record
pointing at it inherits that kind and no record may carry a different one. Without this, `group_accepted`
would be ill-defined — the same proposition would be partly in and partly out of a view — and `_accepts`
could not be applied at identity granularity at all. A violation is an integrity failure, reported as such.

This is also what lets §5 state the authority gate in terms of the *view* rather than the individual
record: if kind were per record, "this view" would not be a well-defined scope for anything.

### 3.3 State, revised

```
_proposition_state(prop_id, *, domain, target, view, accepted_ids) ->
    records = _proposition_records(prop_id, domain, target, view)      # accepted kinds only
    groups  = [g for g in self._open_groups if prop_id in propositions_of(g)]
    live    = [g for g in groups if group_establishes(g, view, accepted_ids)]

    if live:                    "conflict"           # a relation asserted *in this view*
    elif len(records) >= 2:     "unresolved_relation" # one identity, several accounts, nothing asserted
    elif len(records) == 1:     "settled"
    else:                       "none"
```

The signature change is the whole correction: `view` and `accepted_ids` are now inputs, so state cannot be
computed without them, and a call site cannot forget the boundary the way REV 6's could.

| Property | Consequence |
|---|---|
| `accepted_ids` is the answer's own identity set | A group reaching an identity this answer did not return cannot affect this answer |
| Conflict is tested **before** record counts | An established relation is still reported as `conflict`, not downgraded to `unresolved_relation` |
| A non-establishing group is tested **nowhere** | It cannot force `conflict`, and it cannot force `unresolved_relation` either — it is not a weaker relation, it is not this view's relation |
| `none` at zero records | Unchanged: absent by construction, not filtered |

### 3.4 Cross-kind groups are refused at creation

The participation rule protects evaluation. It cannot, by itself, prevent the registry from holding a
relation that no view can evaluate correctly on both sides. So authoring is narrowed too — **fail closed,
before mutation**:

```
assert_incompatible(record_a, record_b, applicability, statement, *, actor)
    kinds = {kind_of(record_a), kind_of(record_b)}
    if len(kinds) != 1:  refuse("cross_kind_group")     # no group created, no turn advanced, nothing mutated
```

| Question | Answer |
|---|---|
| Is this a new restriction on expressiveness? | It removes one expression: "world claim X is incompatible with ascription Y". That expression has no meaning the architecture can honour |
| Why not? | An incompatibility claim asserts that two records cannot both hold of one applicability. An ascription does not hold of the world at all — C-2 (REV 2 §4.1–4.2) confines it to "what does `<character>` believe", and rules 7 and 56 keep the two apart. Asserting rivalry between them would make a canon view adjudicate a belief, which is the leak this revision exists to close |
| What is the correct expression instead? | Belief diverging from canon is **expected**, not conflict: canon says the king is dead, Sable believes he lives, both stand (rule 15 / O-5 / D7). A genuine canon rivalry is asserted between two same-kind records |
| What if a mixed group already exists? | Only reachable by import or a bug, since creation refuses. §3.1 row 3 handles it: no state effect, opaque integrity code, detail to the authority audience |

`assert_incompatible` is a REV 3 §7.1 item that remains **[PA]** and unimplemented, so this narrows a
proposal, not behaviour. No existing record can be affected: the operation does not exist in
`eval/mem1/m09/`.

### 3.5 What the conflict model still is

| Statement | Status |
|---|---|
| Conflict requires an **established open group** | Unchanged |
| Conflict is never inferred from disagreement, multiplicity or inaction | Unchanged (rule 41, O-10, J12) |
| Two identities with no asserted relation are `multiple`, not conflict | Unchanged (REV 6 §4) |
| One identity with ≥2 live records and no assertion is `unresolved_relation` | Unchanged |
| A group affects only the identities it names, in views that contain it whole | **Tightened** — this revision |
| The engine never resolves a conflict on the caller's behalf | Unchanged |

No state is added, no state is renamed, and no new relation type appears. `conflict` means what it meant:
a human or authorised actor asserted that these records cannot both stand, and it is still open.

---

## 4. Effect on the three views

**Envelope, revised in two places.** `open_groups` now carries only establishing groups; `integrity` is
new and normally empty.

```
engine.world(world, "king.status")
→ { "status": "settled", "relation": "not_established", "view": "world",
    "open_groups": [], "integrity": [],
    "propositions": [ … ] }

engine.propositions(world, target="king.status", actor="operator")
→ { "status": "ok", "viewer": "authority", "knowledge_conferred": false,
    "integrity": [ { "code": "group_spans_query_kinds", "group_id": "group-00002" } ],
    "propositions": [ … ] }
```

| Field | Rule |
|---|---|
| `open_groups` | Only groups with `group_establishes(...)` true. Every id inside is already visible in this view, so listing adds no disclosure |
| `relation` | `"asserted"` only when an establishing group spans returned entries; otherwise `"not_established"` |
| `integrity` | Codes only in `world` / `belief_facts`: `"group_spans_query_kinds"`, `"proposition_kind_mismatch"`. **No group id, no foreign record id, no text, no count** — the count of ascription members is itself information under D-02 behaviour 8 |
| `integrity` in `propositions` | Same codes with `group_id` attached. The authority audience already reads both kinds, so detail is not a disclosure there (REV 4 §4.4 reasoning) |
| `state` per entry | `_proposition_state(entry.id, …, view=view, accepted_ids=…)` |
| `status` | Unchanged from REV 6 §4: `none` / `settled` / `unresolved_relation` / `conflict` / `multiple` |

**Why `integrity` rather than silence.** REV 6 §6 already says a violation is "an integrity failure
reported as such, not filtered away". Silently dropping a straddling group would make a corrupt registry
look correct, and neither `world()` nor `belief_facts()` would ever reveal it. A code with no payload keeps
both commitments: the fault is visible, and nothing outside the view is disclosed.

**Why not `unresolved_relation`.** That state already means something specific — one identity, several
accounts, nothing asserted. Reusing it for a boundary fault would tell the caller to look for a dispute
that does not exist on this applicability, which is the REV 5 conflation in a new costume.

---

## 5. The `propositions` gate, as an invariant of the view

REV 6 §3.2 put the kind test **inside** `authoritative()` on the argument that "a new call site cannot
forget it". That argument applies with equal force to the authority test, and REV 6 left the authority test
in the wrapper. Applying one half of a principle is how a second route appears.

### 5.1 Callable boundary and order

```
authoritative(self, domain, target, *, view, actor=None):

  1  view == "propositions"       → requires actor; no default
  2  view not in KNOWN_VIEWS      → refuse("unknown_view")            # fail closed, F4
  3  view == "propositions":
  3a     actor is None            → refuse("no_domain_authority")
  3b     authority.require(domain, actor, "promotion")                # R-A item 4
  3c     AuthorityError           → refuse("no_domain_authority")     # uniform, detail withheld
  4  eligible = _eligible_truth(domain, target)                       # first data read
  5  accepted = [r for r in eligible if _accepts(view, kind_of(r))]
  6  identity set, state, envelope as §3–§4
```

| Property | Consequence |
|---|---|
| Step 3 precedes step 4 | No record, no identity, no count is touched before authorisation |
| Step 3c is uniform | "Domain not declared" and "no live `promotion` grant" and "`actor=None`" are the same refusal: `{"status": "refused", "reason": "no_domain_authority"}` and nothing else. Uniformity is imposed at this boundary; `authority.py` is not rewritten to achieve it |
| `actor` is required for one view, ignored by the others | `world` and `belief_facts` gain no gate and lose none. Passing an `actor` to them does not widen eligibility by one row (REV 4 §3, "why authority is not access") |
| The gate is on the **view**, not on the wrapper | Any future call site — a test harness, a second surface, a replay reader — reaches the registry only through step 3. The wrapper's own check becomes a duplicate, not the defence |
| `require(domain, actor, "promotion")` | R-A's per-operation form (REV 3 §3 item 4). The implemented `AuthorityRegistry.require(domain, actor)` lacks the operation argument; it is unmodified by this record and this is not a claim that it conforms |

### 5.2 The wrapper keeps its check, and why that is not redundant

`propositions(domain, target, *, actor)` still authorises before building its envelope, because the
envelope itself — `viewer: "authority"`, `knowledge_conferred: false`, the L1 metadata, the refusal
envelope — is part of the surface REV 4 §3 defines. The second check is idempotent, read-only, and occurs
before any data read; a `require` call is not a mutation.

What changes is the **claim**: REV 4 §3 said enumeration is available to an actor holding a live
`promotion` capability "and to no one else". REV 7 makes that true of the view rather than of one
function's body.

### 5.3 The other second route, closed by the same reasoning

Enumeration has two ways to leak: a caller reaches the registry (gated above), or a canon surface starts
listing identities the registry knows. REV 6's pipeline already prevents the second —
`grouped = group accepted by proposition_id` derives the identity set **from accepted records**, never from
the registry's identity list — so `world()` cannot reveal an identity that has no world record, and cannot
reveal one that has only ascription records. REV 7 states this as invariant 9 (§6) rather than leaving it
an artefact of the code shape, because it is the property that keeps `world()` from being an enumeration
surface by accident.

---

## 6. Invariants

Per view, for a caller permitted that view:

1. Every returned record satisfies `_accepts(view, kind_of(record))`.
2. Every returned record's `proposition_id` equals its entry's id.
3. No record appears in two entries — single-valued `proposition_id`.
4. The union of entries' records equals the view-accepted eligible set.
5. **Participation:** an entry's `state` is a function of its own view-accepted records and of open groups
   that `group_establishes(...)` in this view. No other input affects it.
6. **No absorbable fault:** a group that straddles the view changes no `state` and no `relation`, and
   appears only as an `integrity` code.
7. **Complement:** `world` records ∪ `belief_facts` records = `_eligible_truth(domain, target)`, and the
   intersection is empty.
8. **Gate:** every `propositions`-view call that returns a non-refusal has completed
   `require(domain, actor, "promotion")` in the same call before its first data read; every refusal
   contains no registry-derived field.
9. **Scope:** every id and every text in a canon-view response belongs to a record accepted in that view.
   The identity set of `world` / `belief_facts` is derived from accepted records, never from the registry's
   identity list.
10. **Kind:** every record of an identity carries that identity's kind (§3.2).

Items 5, 6, 8, 9, 10 are this revision's additions; 1–4 and 7 are REV 6's, unchanged. A violation is an
**integrity failure reported as such**, not filtered away.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 7. Worked example — the straddling group, in three views

State: as §2 steps 1–3. Assume a mixed group exists, which §3.4 now prevents from being created.

| View | Identity set | State of `prop-00003` | `relation` | `open_groups` | `integrity` |
|---|---|---|---|---|---|
| `world` | `prop-00003` (ascription has no world record) | `settled` | `not_established` | `[]` | `["group_spans_query_kinds"]` |
| `belief_facts` | `prop-00011` | `settled` | `not_established` | `[]` | `["group_spans_query_kinds"]` |
| `propositions` (operator) | `prop-00003`, `prop-00011` | `conflict` for both | `asserted` | `group-00002`, both members, kind `claim_about_world` + `belief_ascription` | `group-00002` |

Read across the row: the world answer is not altered by a belief-fact relation; the belief answer is not
altered by a world relation; the authority view sees both and reports the assertion it was actually given.
Nothing is resolved, and nothing is hidden from the only audience entitled to the whole picture.

Second example — the normal case, unchanged:

| Turn | Event | `world` | `propositions` |
|---|---|---|---|
| 4 | `prop-00003` promoted: "the king is alive" | `settled` | one identity, `settled` |
| 6 | `prop-00007` minted and promoted: "the king is dead" | `multiple`, relation `not_established`, two entries | two identities, each `settled` |
| 8 | operator asserts incompatibility between `truth-00012` and `truth-00019` | `conflict`, `open_groups: [group-00001]` | both `conflict` |
| 9 | model without a grant calls `authoritative(..., view="propositions")` | — | `refused / no_domain_authority`, no ids, no counts |

---

## 8. Coherence of the two corrections

A check that the corrections do not contradict the ratified set, not conformance evidence.

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14; REV 2 §3.2 | **C** — participation is view-relative |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** — kind filter plus a state that cannot see the other kind |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** — still requires an establishing open group |
| Underspecification is not conflict | 8, 34 | K15 | **C** — `multiple` unchanged; faults use `integrity`, not state |
| No automatic adjudication | 34 | K15, O-10 | **C** — engine reports, never resolves |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — gate is an invariant of the view |
| Refusal leaks nothing | 21 | REV 4 §7; D-02 behaviour 8 | **C** — refusal precedes reads; uniform code |
| Belief canon ≠ belief record | 15, 57 | O-5; F2; D7 | **C** — untouched; §3.4 makes the distinctions harder to blur |
| Fail closed on unknown input | 8 | F4 | **C** — unknown view; unknown kind; mixed group |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded and audience-constrained; R-B,
R-C, R-D open.

---

## 9. Implementation list — deltas only

REV 6 §8 stands, with these rows replaced. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17b″ (REV 6 17b′) | `engine.py` | `_accepts(view, kind)` | Unchanged in body; `KNOWN_VIEWS` named once so §5.1 step 2 and the dispatch cannot drift |
| 17c′ (REV 6 17c) | `engine.py` | `authoritative(domain, target, *, view, actor=None)` | View required; **gate at step 3 for `propositions`**, before any read; uniform `no_domain_authority`; per-view dispatch; envelope per §4 |
| 6″ (REV 6 6′) | `retrieval.py` | `world` | Unchanged shape; `open_groups` filtered to establishing groups; `integrity` codes only; no foreign id, text or count |
| 9‴′ (REV 6 9‴) | `engine.py` | `propositions(...)` | Keeps its own check before building the envelope; delegates retrieval to `authoritative(..., view="propositions", actor=actor)`; reports `integrity` with `group_id` |
| 18′ (REV 6 18) | `retrieval.py` | `belief_facts` | Unchanged shape; same `integrity` rule as `world` |
| 17d | `engine.py` | `group_establishes(group, view, accepted_ids)`, `group_accepted`, `group_contained` | New; the only place §3.1's three-way relation is decided |
| 17e | `engine.py` | `_proposition_state(prop_id, *, domain, target, view, accepted_ids)` | View-relative; conflict before counts; non-establishing groups ignored |
| 12′ (REV 3 row 12) | `engine.py` | `assert_incompatible` | Refuses `cross_kind_group` before mutation; still spans propositions of one kind (REV 3 §5.4 unaffected) |
| 19 | `records.py` / `engine.py` | kind of a proposition | One accessor used by `_accepts`, `group_accepted` and the §3.2 integrity check; no second derivation |

REV 6's `supersede`/D-A, `retract`/D-B and every row not listed are unchanged. **D-A and D-B remain
confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical), `eval/mem1/model.py`, the G1 harness, production
Howl Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–6.

---

## 10. What this record does and does not claim

**Still open, carried forward unchanged — not reasons to reopen these contracts:**

- **B3** — event versus definition. Unresolved; the owner's instruction not to expand the contract quietly
  stands.
- **Belief-voiced dialogue** — how a character's own belief reaches its own voice is a retrieval and
  prompt-assembly question, outside these interfaces.
- **Replay / F7** — no facility; determinism remains unproven, and the only determinism claim here is
  ordering by `proposition_id` then `record_id`.
- **R-B, R-C, R-D** — as recorded in REV 3.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- No behavioural verification. Statements about current behaviour remain source inspection of
  `eval/mem1/m09/` as reviewed in R1/R2.
- No conformance proof. §8 checks coherence only.
- No new rule or requirement. `integrity` is a diagnostic field of a proposed envelope, not an
  architecture state; `cross_kind_group` is a refusal code, not a new rule; neither adds a code to the map.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9
  (`698254ca…`), the map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–6 (`98447423…`, `74ac6f77…`,
  `de4f68b6…`, `1da08f2c…`, `a7e7c013…`, `41421dce…`).
- No transfer, destination interface, sharing, or cross-chat carryover.

**Sequence from here:** owner assessment of the complete proposal (REV 1–7) for ratification →
separate implementation authorization → implementation → conformance review → separate testing
authorization. This record does not request ratification and does not request implementation.
