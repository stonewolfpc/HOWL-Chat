# MEMORY-09 — Architectural contracts, revision 6

**Status: PROPOSAL REV 6 — world results preserve proposition identity; query kind is a third explicit
selection. Not implemented, not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 5 (`a7e7c013…`) **for §3.3, §3.4, §4 item 3, §6 turn 8 and §9 rows 17b/9″ only**; and REV 3 §7.1 item 6's wording |
| Stands unamended | REV 5 in every other respect; REV 1–4 preserved byte-identical |
| Authority | Owner's two-boundary instruction, recorded verbatim (§1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger. Corrects a proposal and one unapplied implementation-list row |
| Ratification | **Not ratified.** Implementation authorization remains a separate gate |

---

## 0. What changed

| Item | REV 5 | REV 6 |
|---|---|---|
| `world()` across identities | One aggregate; two distinct settled identities became `unresolved_relation` | **Per-proposition results**; top-level `multiple` when several identities apply, with no relation asserted |
| Relation between identities | Implied by aggregation | Reported only when an **open conflict group** establishes it |
| Record eligibility | Shared predicate, kind unmentioned | Unchanged: access only, kind-agnostic |
| Query kind | Implicit in each retrieval function | **Third explicit dimension**; `authoritative()` requires `view`, and refuses an unknown one |
| Belief ascriptions in world answers | Excluded by convention, enforced nowhere | **Excluded structurally** by query kind; a separate `belief_facts` view is their only canon surface |
| Partition invariant | Over one eligible set | **Per view**, plus an exact-complement check between `world` and `belief_facts` |

No rule, requirement or tally changes. The 57/118 codes are untouched.

---

## 1. Owner instruction, recorded verbatim

> "`world()` calls two distinct propositions an `unresolved_relation`, it risks reintroducing the very
> conflation U-1 was intended to prevent. There is an important difference between: 1. The engine cannot
> determine which of two incompatible accounts governs the same proposition. 2. The query identifies an
> applicability containing multiple independently declared propositions. The second is not necessarily an
> unresolved relation. It may simply be an underspecified query."

> "My recommendation: `world()` should preserve proposition identity in its result. When multiple
> identities apply to the requested target, return the separate proposition results rather than
> collapsing them into one `unresolved_relation`. The caller can then see what is established for each
> proposition without the engine inventing a relationship between them."

> "C-2 requires `world()` to exclude `belief_ascription` propositions, while a separate belief-fact view
> can retrieve them. The contract should make this distinction explicit: Eligibility … Identity
> selection … Query-kind selection … Otherwise, sharing the eligibility predicate could accidentally make
> belief-ascription records eligible for ordinary world-truth answers."

**Scope:** these two boundaries only. No architecture change, no ratification, no implementation.

---

## 2. Why the aggregation was wrong

REV 5 computed `world()` from `_eligible_truth(domain, target)` and then applied a single rule — two or
more active records means `unresolved_relation`. That rule is correct **within one proposition** and
false **across propositions**, and the surface could not tell the difference because it had discarded
identity before deciding.

It also contradicted text already in the proposal set. REV 2 §3.2 states the rule plainly:

> "Different propositions with the same applicability are answered separately. This is the owner's U-1
> requirement — 'different propositions about the same subject must not automatically conflict' — made
> structural rather than conventional."

REV 5's world surface was therefore not a gap in the architecture; it was a proposal that failed to
implement its own §3.2. REV 3 §7.1 item 6 ("`unresolved_relation` for ≥2 live records; never `settled`")
is the row that permitted it, and REV 6 narrows that row to **≥2 live records on one proposition**.

| Question | Correct answer | Wrong answer |
|---|---|---|
| Two live records, one identity, no assertion | `unresolved_relation` — the engine cannot tell which governs | — |
| Two identities, one record each | Both `settled`; the query is underspecified as to identity | `unresolved_relation` |
| Two identities, open group spanning them | Both `conflict`; the relation **was** asserted | `unresolved_relation` |

---

## 3. Three dimensions, three functions

The owner named eligibility, identity and query kind. They are now separate, applied in this order, and
none re-implements another.

```
eligible   = _eligible_truth(domain, target)          # access: may this caller see the record at all?
accepted   = [r for r in eligible if _accepts(view, kind_of(r))]   # query kind
grouped    = group accepted by proposition_id          # identity
```

### 3.1 Eligibility — access only, kind-agnostic

```
_eligible_truth(domain, target) ->
    [r for r in self._records.values()
     if r.truth_apt and r.owner_character is None
     and r.domain == domain and r.target == target and r.status == ACTIVE]
```

Unchanged from REV 5. It answers **visibility**, and deliberately says nothing about proposition kind:
putting the kind test here would make a false ascription record *inaccessible* rather than *inapplicable*,
which is the wrong failure and would hide the record from the view that should answer it.

### 3.2 Query kind — which question is being asked

```
_accepts(view, kind):
    "world"        -> kind == "claim_about_world"
    "belief_facts" -> kind == "belief_ascription"
    "propositions" -> kind in {"claim_about_world", "belief_ascription"}   # authoring surface
    otherwise      -> False            # fail-closed
```

| Property | Consequence |
|---|---|
| `authoritative(domain, target, *, view)` — `view` is **required**, keyword-only | No default view; omission is a refusal, not a silent choice of world truth |
| Unknown view | Refused / empty, never widened (F4) |
| `kind` is closed at two values (REV 2 §3.1) | The two canon views are exact complements; a third kind would break the check in §4, which is the intended alarm |
| The kind test lives **inside** `authoritative()`, not only in `retrieval.world` | A new call site cannot forget it |

### 3.3 Identity — which proposition a record belongs to

```
_proposition_records(prop_id, domain, target, view) ->
    [r for r in _eligible_truth(domain, target)
     if r.proposition_id == prop_id and _accepts(view, kind_of(r))]
```

Unchanged in shape from REV 5, with the view parameter. Access is inherited through `_eligible_truth`;
identity is exact; no fallback, no widening.

### 3.4 State, per identity, view-independent

`_proposition_state(prop_id)` is unchanged from REV 5 §3.4: `conflict` if one of its records is in an
open group, else `unresolved_relation` at two or more active records, `settled` at one, `none` at zero.
State is a property of an identity, not of a query, so it takes no view.

---

## 4. `world()` — identity preserved

**Return shape**

```
engine.world(world, "king.status")
→ { "status": "multiple",
    "reason": "multiple_propositions",
    "relation": "not_established",
    "view": "world",
    "open_groups": [],
    "propositions": [
      { "proposition_id": "prop-00003", "kind": "claim_about_world", "state": "settled",
        "records": [ { "record_id": "truth-00012", "text": "the king is alive",
                       "created_turn": 4, "promotion_id": "promo-00004" } ] },
      { "proposition_id": "prop-00007", "kind": "claim_about_world", "state": "settled",
        "records": [ { "record_id": "truth-00019", "text": "the king is dead",
                       "created_turn": 6, "promotion_id": "promo-00009" } ] } ] }
```

| Field | Rule |
|---|---|
| `propositions` | One entry per accepted identity on the applicability, ordered by `proposition_id` (F7) |
| `records` inside an entry | `_proposition_records(...)`; that identity only |
| `state` | `_proposition_state(entry.proposition_id)` |
| Top-level `records` | Present **only** when exactly one entry is returned; omitted entirely for `multiple`, so a caller cannot read a partial set as the whole answer |
| `relation` | `"not_established"` unless an open group spans returned entries; then `"asserted"` and `"open_groups"` names it |
| `label` | **Not returned.** REV 4 §4.4 confines labels to the domain-authority audience; `world()` is not that audience. Identity is preserved by `proposition_id` plus record text |

**Status**

| Situation | `status` | Note |
|---|---|---|
| No accepted identity | `none` | Abstention, not error |
| One identity | That identity's state: `settled`, `unresolved_relation`, `conflict`, or `none` | Same as today for the single-identity case |
| Several identities, no open group spanning them | `multiple` | **The query is underspecified as to identity**, not in unresolved conflict |
| Several identities, open group spanning them | `conflict` | The relationship was asserted; `open_groups` names it |

`multiple` is deliberately not a synonym for `unresolved_relation`. It says: this applicability holds
several independently declared propositions, each answered separately; the engine has not established a
relation between them and does not imply one. A caller that needs one answer must name a proposition, or
ask a human to establish a relation.

**What this costs.** A caller that used to receive a flat `records` list now receives entries and must
read them. That is the intended cost: the flat list was the conflation.

**What it does not change.** A single identity with two live records still reports `unresolved_relation`;
an open group still reports `conflict`; nothing is ever selected as settled on the caller's behalf.

---

## 5. `belief_facts()` — the other canon view

C-2 (REV 2 §4.1–4.2) already requires that a `belief_ascription` proposition answers **only** "what does
<character> believe" and is never returned by a world-truth query. REV 5 left that as prose; REV 6 makes
it a query kind.

```
engine.belief_facts(world, "king.status")
→ { "status": "settled", "view": "belief_facts", "relation": "not_established",
    "propositions": [
      { "proposition_id": "prop-00011", "kind": "belief_ascription",
        "about_character": "sable", "state": "settled",
        "records": [ { "record_id": "truth-00031", "text": "Sable believes the king was poisoned", … } ] } ] }
```

| Rule | Reason |
|---|---|
| Only `belief_ascription` identities appear | `_accepts("belief_facts", kind)`; a world proposition is invisible here too |
| `about_character` is returned | Without it the answer cannot say whose belief is authoritative; it is a property of the canon record |
| This is **not** Sable's belief record | It returns authoritative records *about* Sable (O-5, rule 56). Sable's private belief stays unlinked, unenumerated and owner-only (REV 4 §4.3) |
| Same audience as `world()` | Both are canon surfaces, not character perspectives; neither confers knowledge (rules 18, 32) |
| No actor gate | Consistent with the existing `world()`; this revision introduces no new gate and removes none. Enumeration remains the only authority-gated surface (REV 4 §3) |

**The leak this closes.** With eligibility shared, an ascription record was *visible* to `world()` and
excluded only because no code filtered it. Now the exclusion is a named predicate applied inside
`authoritative()`, and the two views are provably disjoint (§6).

---

## 6. Invariants

Per view, for an authorised caller:

1. Every returned record satisfies `_accepts(view, kind_of(record))`.
2. Every returned record's `proposition_id` equals its entry's id.
3. No record appears in two entries — single-valued `proposition_id`.
4. The union of entries' records equals the view-accepted eligible set.
5. State is computed from the entry's own records and open groups only.
6. **Complement:** `world` records ∪ `belief_facts` records = `_eligible_truth(domain, target)`, and the
   intersection is empty.

A violation is an **integrity failure reported as such**, not filtered away. Item 6 is the check that
would have caught the C-2 leak: an ascription appearing in a world answer fails disjointness, and a world
record disappearing from both views fails completeness.

These invariants constrain the proposed mechanism. They add nothing to the ratified map and are not
executed here.

---

## 7. Conformance of the two corrections

A coherence check, not conformance evidence.

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Distinct identities answered separately | 6, 41 | U-1; REV 2 §3.2; J13, J14 | **C** — `multiple`, no relation invented |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** — `conflict` requires an open group |
| Underspecification is not conflict | 8, 34 | K15 (no assessment) | **C** — stated as a query fact, not a verdict |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** — query-kind filter inside `authoritative()` |
| Belief canon ≠ belief record | 15, 57 | O-5; F2; D7 | **C** — ascription records are domain-owned; belief stays private |
| Per-identity state, no cross-talk | 6 | J12, J14 | **C** — state from own records |
| Deterministic order | 34 | F7 | **C** — by proposition id, then record id |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded and audience-constrained; R-B,
R-C, R-D open.

---

## 8. Implementation list — deltas only

REV 5 §9 stands, with these rows replaced. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 6′ (REV 3 §7.1) | `retrieval.py` | `world` | Returns the per-proposition envelope; `multiple`/`conflict` as §4; never flattens entries into one record list |
| 17b′ | `engine.py` | `_accepts(view, kind)` | The query-kind predicate; unknown view false |
| 17c | `engine.py` | `authoritative(domain, target, *, view)` | Required view; applies `_accepts` **inside**; groups by identity; top-level status per §4 |
| 9‴ | `engine.py` | `propositions(...)` | Uses `view="propositions"` (both kinds); otherwise REV 5 §5.1 |
| 18 | `retrieval.py` | new `belief_facts` | `authoritative(..., view="belief_facts")`; returns `about_character`; no record outside ascription identities |
| — | `__init__.py` | `__all__` | Export `belief_facts` alongside the REV 3 §7.2 list |

REV 3 §7.1 item 6 is narrowed to "≥2 live records **on one proposition** → `unresolved_relation`"; the
cross-identity aggregate is withdrawn. `supersede`/D-A, `retract`/D-B and every other row are unchanged.

**Untouched:** `eval/mem1/model.py` and the G1 harness; production Howl; r32; plan v3; Amendment 1; v9;
the map; the ledger; REV 1–5.

---

## 9. What this record does and does not claim

**Open:** U-6 / B3 — untouched. F7 replay — no facility, determinism unproven. R-B, R-C, R-D as before.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported or compiled.
- No behavioural verification. Current-behaviour statements remain source inspection.
- No conformance proof. §7 checks coherence only.
- No new rule or requirement. `multiple` is a result status of a proposed query, not a new architecture
  state, and `relation: "not_established"` adds no field to any ratified record.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9
  (`698254ca…`), the map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–5 (`98447423…`, `74ac6f77…`,
  `de4f68b6…`, `1da08f2c…`, `a7e7c013…`).
- No transfer, sharing, destination interface, or cross-chat carryover.

**Sequence from here:** owner assessment of the complete proposal for ratification → separate
implementation authorization → implementation → conformance review → separate testing authorization.
