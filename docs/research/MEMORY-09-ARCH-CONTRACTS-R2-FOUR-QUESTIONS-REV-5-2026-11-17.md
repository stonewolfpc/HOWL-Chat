# MEMORY-09 — Architectural contracts, revision 5

**Status: PROPOSAL REV 5 — eligibility separated from identity selection; audit disclosure stated.
Not implemented, not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 4 (`1da08f2c…`) **for §4.2, §5.1, §5.2, the §9 projection rows, and §10 rows 9′, 10′, 17 only** |
| Stands unamended | REV 4 in every other respect; REV 1 (`98447423…`), REV 2 (`74ac6f77…`), REV 3 (`de4f68b6…`) and REV 4 preserved byte-identical |
| Authority | Owner's consistency finding and clarification instruction, recorded verbatim (§1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed |
| Amends | **Nothing.** No edit to r32, v9, the 118-requirement map, the ledger, or any prior record |
| Ratification | **Not ratified.** One narrowly scoped correction to the proposal; implementation authorization remains separate |

---

## 0. What changed

| Item | REV 4 | REV 5 |
|---|---|---|
| `_truth_projection(domain, target)` | One predicate used as both access rule and record selection | **Two functions.** `_eligible_truth` is the access rule; `_proposition_records` adds the exact `proposition_id` |
| Records per proposition entry | Every active truth record on the applicability | **Only records whose `proposition_id` is that identity** |
| `state` | Computed from the applicability's records | **Computed per identity**, from its own records plus open groups containing them |
| Partition guarantee | Implicit, and false | **Explicit invariant**: enumeration partitions the eligible set; no record appears twice, none is orphaned |
| Domain-wide multiplicity | Described as narrow | **Stated as an authorised audit disclosure** — not zero-information (§5.3) |
| `world()` | Unchanged | **Unchanged.** It still selects all eligible records for its query; only enumeration stops borrowing its selection |

No rule, requirement, decision or tally changes.

---

## 1. Owner finding and instruction, recorded verbatim

> "L2 does not preserve proposition identity … `_truth_projection("world", "king.status")` returns both
> records. … `propositions(domain, target)` is supposed to enumerate identities and return the records
> belonging to each identity. The proposed projection instead returns every active truth record sharing
> the applicability, regardless of identity."

> "Consequences include: Records appearing under the wrong proposition. Proposition states being
> calculated from unrelated records. A caller being shown a relationship between records that the engine
> has not established. The same records potentially appearing under multiple proposition entries."

> "Recommendation: Separate eligibility from identity selection. Keep a shared access/eligibility
> predicate where appropriate, but make proposition enumeration additionally require the exact
> `proposition_id`. `world()` can continue retrieving all eligible records for its intended query. The
> two surfaces should share the access rules without pretending their selection criteria are identical."

> "`possible_duplicates(domain, target=None)` … still reveals that multiple identities exist somewhere in
> the domain. The proposal should explicitly acknowledge that this is an authorized audit disclosure to
> promotion-capable actors, not a zero-information response. The existing refusal-before-registry-read
> rule for unauthorized callers remains essential."

**Scope of this revision:** the projection mismatch and the audit-disclosure statement. No architecture,
no ratifying text, no implementation is touched.

---

## 2. The defect, stated exactly

REV 4 gave one function two jobs. An **access rule** answers "may this actor see this record at all?".
A **selection rule** answers "does this record belong to the thing being asked about?". They coincide for
`world()`, whose subject is an applicability, and they do **not** coincide for enumeration, whose subject
is an identity.

With `prop-001` and `prop-002` both on `world / king.status`:

| Consequence named by the owner | Mechanism in REV 4 |
|---|---|
| Record under the wrong proposition | `_truth_projection` ignores `proposition_id`, so Truth B appears under `prop-001` |
| State from unrelated records | REV 4 §5.1 computed `state` from the same unfiltered list; two settled identities would both read `unresolved_relation` |
| Unestablished relationship | Showing A and B together asserts a shared subject the engine never established — and it is exactly the assertion U-1 forbids it to make |
| Duplicate appearance | Every identity on the target would list every record on the target |

The last one is the tell: a correct projection over a partition cannot yield duplicates. REV 4's could,
which is why §4.2 below replaces the shape rather than patching the filter.

---

## 3. Eligibility separated from selection

### 3.1 The shared access rule — unchanged semantics

```
_eligible_truth(domain, target) ->
    [r for r in self._records.values()
     if r.truth_apt and r.owner_character is None
     and r.domain == domain and r.target == target and r.status == ACTIVE]
```

| Question answered | "May an authorised caller see this authoritative record?" |
|---|---|
| Callers | `authoritative()` (and therefore `world()`) **and** enumeration |
| Selection performed | **None.** It is a visibility predicate over domain records |
| Private classes | Cannot pass: `truth_apt` excludes claim, account and belief; `owner_character is None` excludes anything owned |

This is REV 4's predicate with a name that states its job. Nothing about access widens or narrows.

### 3.2 The identity selection — the addition

```
_proposition_records(prop_id, domain, target) ->
    [r for r in _eligible_truth(domain, target) if r.proposition_id == prop_id]
```

| Property | Consequence |
|---|---|
| Access is inherited **through** `_eligible_truth`, never re-implemented | One visibility rule; no second filter to drift |
| `proposition_id` is a single required field on truth records (REV 4 §4.3) | A record belongs to exactly one identity, so it can appear under exactly one entry |
| No fallback when `prop_id` matches nothing | Empty list; the entry reports `state: none`, never a widened query |
| Domain and target are still checked | A record attached to a foreign-domain identity cannot be pulled in by id alone |

### 3.3 Why `world()` is not changed

`world()` asks *what does the authoritative record say about this applicability*. Its subject is the
applicability, so selecting all eligible records is correct there, and a conflict or unresolved relation
is the honest answer when several apply. Enumeration asks *what identities exist on this applicability
and what does each one hold*. Different question, different selection, same access rule.

The two surfaces now share **access** and differ **only** in selection — which is what the owner asked
for, and what "share the access rules without pretending their selection criteria are identical" means
mechanically.

### 3.4 Proposition state is computed per identity

```
_proposition_state(prop_id, domain, target):
    records = _proposition_records(prop_id, domain, target)
    open_groups = groups with status open containing any record in records
    if open_groups:        "conflict"
    elif len(records) >= 2: "unresolved_relation"
    elif len(records) == 1: "settled"
    else:                   "none"
```

| Rule | Reason |
|---|---|
| Precedence: `conflict` > `unresolved_relation` > `settled` > `none` | A record inside an open group is contested regardless of how many records its identity holds |
| `unresolved_relation` is per identity | Two live records on **one** identity; two settled identities are two settled identities (RAT-2) |
| A group spanning identities marks **each** identity that holds a member as `conflict` | The assertion did establish that relationship; it is not an inference (REV 3 §5.8) |
| No state is derived from another identity's records | The REV 4 defect, closed by construction |

**Relationship to the world surface, stated so it is not mistaken for a new rule.** `world()` aggregates
across identities and, with two or more eligible active records, reports an unresolved relation rather
than a settled answer (REV 3 §7.1 item 6, still unapplied; RAT-2). So a target can hold two identities
that are each `settled` while the applicability-level answer is `unresolved_relation`. That is not a new
decision and this revision adds no rule for it: it follows from U-1 (distinct propositions do not
automatically conflict) and RAT-2 (multiple live records are not automatically settled). REV 5's
contribution is that enumeration no longer presents the applicability-level aggregate as if it were each
identity's state.

---

## 4. The partition invariant

The corrected shape is checkable without judgement, which is what makes the original bug detectable
rather than merely fixed.

**For any `(domain, target)` with authorised enumeration:**

1. Every record in every entry's `records` satisfies `r.proposition_id == entry.id`.
2. No record id appears in two entries — guaranteed by (1) and single-valued `proposition_id`.
3. The union of all entries' `records` equals `_eligible_truth(domain, target)`.
4. No entry has a record whose `domain` or `target` differs from the query.
5. `state` for each entry is computed from that entry's `records` only.

A violation of (1), (2) or (4) is an **integrity failure reported as such**, not silently filtered.
A violation of (3) — an eligible truth record attached to no enumerated identity — is likewise an
integrity failure, because `proposition_id` is required: an orphan means the registry and the record
store have diverged. The response must not quietly omit it, or enumeration would under-report canon
while looking complete.

This invariant is an implementation self-check on the proposed mechanism. It is **not** a new map
requirement, and it is not executed here.

---

## 5. API contracts, corrected

### 5.1 `propositions(domain, target, *, actor)`

Same envelope, corrected contents:

```
engine.propositions(domain, target="king.status", actor="operator")
→ {
    "status": "ok",
    "viewer": "authority",
    "knowledge_conferred": false,
    "propositions": [
      { "id": "prop-00003", "label": "king.status.alive", "kind": "claim_about_world",
        "about_character": null, "created_by": "operator", "created_turn": 4,
        "state": "settled",
        "records": [ { "record_id": "truth-00012", "text": "the king is alive",
                       "created_turn": 4, "promotion_id": "promo-00004" } ] },
      { "id": "prop-00007", "label": "king.status.death", "kind": "claim_about_world",
        "about_character": null, "created_by": "model:canon", "created_turn": 6,
        "state": "settled",
        "records": [ { "record_id": "truth-00019", "text": "the king is dead",
                       "created_turn": 6, "promotion_id": "promo-00009" } ] } ] }
```

| Field | Source |
|---|---|
| `state` | `_proposition_state(entry.id, domain, target)` |
| `records` | `_proposition_records(entry.id, domain, target)` |
| Everything else | L1, unchanged from REV 4 §4.1 |

Each entry holds only its own records. Truth A cannot appear under `prop-00007`, and neither state is
computed from the other's records.

`target` remains **required** (REV 4 §5.1): no subject discovery on this surface.

### 5.2 `possible_duplicates(domain, *, actor, target=None)`

| Form | Returns | Withheld |
|---|---|---|
| `target=<t>` | Groups on `<t>` with ≥2 identities: member `id`s and `count` | labels, texts, `about_character` |
| `target=None` | The same shape, across the domain | labels, texts, `about_character`, targets, single-identity groups |

Grouping is by applicability **internally**; the domain-wide response does not name the target. Response
envelope gains `"disclosure": "authorised_audit"` so the caller is told what kind of answer it received.
The REV 3 §5.5 disclaimer remains verbatim in every response.

### 5.3 Domain-wide multiplicity is an authorised disclosure

REV 4 called this mode "deliberately narrow". Narrow is not the same as empty, and the difference should
be stated rather than left to inference.

| What the caller learns | What it does not learn |
|---|---|
| That at least one applicability in this domain holds two or more declared identities | Which applicability |
| How many such groups exist, and how many identities each holds | Any label, target, text or character |
| The relative mint order of the listed ids, which are sequential | Anything about a private record, which is unlinked (REV 4 §4.3) |
| That canon on some subject is identity-fragmented — the audit signal itself | Whether the identities denote the same proposition (the disclaimer) |

**This is an audit disclosure to a promotion-capable actor, not a zero-information response.** Its
justification is the one REV 3 §5.6 gave: without it, split identities are invisible and conflict
detection can silently halve. The information is disclosed to exactly the actors who can act on it —
those holding a live `promotion` grant in that domain — and it is coarse precisely so that it cannot be
walked back into a subject index.

The refusal rule is the precondition, not a formality: an unauthorised caller receives **no groups, no
counts and no ids**, because a count is itself information (D-02 behaviour 8). Authority is checked
before the registry is read (§7).

---

## 6. Worked example — the two identities REV 4 conflated

State: `prop-00003` holds Truth A ("the king is alive", turn 4). `prop-00007` holds Truth B ("the king is
dead", turn 6). Both are settled identities on `world / king.status`.

| Turn | Call | REV 4 result | REV 5 result |
|---|---|---|---|
| 7 | `propositions(world, "king.status", actor="operator")` | both entries list **both** records; both entries `unresolved_relation` | `prop-00003` → one record, `settled`; `prop-00007` → one record, `settled` |
| 8 | `world(world, "king.status")` | unchanged | unchanged: two eligible records → `unresolved_relation`, never a settled pick |
| 9 | `assert_incompatible(truth-00012, truth-00019, …)` | both entries `conflict` | both entries `conflict` — each holds a group member |
| 10 | `possible_duplicates(world, actor="operator")` | ids and counts | ids and counts, `disclosure: "authorised_audit"` |
| 11 | `possible_duplicates(world, actor="sable")` (no grant) | refused | refused; **no count** |

Turn 7 is the correction: the identity-level answer and the applicability-level answer are now visibly
different, and neither is manufactured from the other's records.

---

## 7. Refusals — unchanged, and load-bearing

| Situation | Result |
|---|---|
| No live grant in `domain` | `refused` / `no_domain_authority`; no identities, counts or ids |
| Grant in a different domain (rule 35) | `refused` / `no_domain_authority` |
| Character view requested | `refused` / `not_a_perspective_surface` |
| `target` omitted on §5.1 | `refused` / `target_required` |
| Unknown target, valid authority | `ok`, empty `propositions` |
| Identity with no active records | entry present, `state: none`, `records: []` |

---

## 8. Conformance of the correction

Not conformance evidence; a coherence check on the corrected items only.

| Property | v9 rules | Map requirements | Verdict |
|---|---|---|---|
| Access rule shared; selection differs | 5, 7 | H2, J14, K12 | **C** — one visibility predicate, two subjects |
| No private record reachable | 5, 6, 47 | H2, G3, D7, J7, J8 | **C** — unchanged from REV 4 §4.3 |
| Per-identity state | 6, 41 | J12, J13, J14 | **C** — state computed from that identity's records |
| No relationship asserted by enumeration | 8, 34 | J12, K15 | **C** — grouping by declared identity, no comparison |
| Partition, no duplicate or orphan | 7 | F1b, F3, F7 | **C** — stated invariant with integrity failure on violation |
| Audit disclosure bounded | 6, 32 | K2, K9, F4 | **C** — authorised, coarse, refusal precedes read |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded and audience-constrained; R-B,
R-C, R-D open.

---

## 9. Implementation list — deltas only

REV 4 §7 and §10 stand, with these rows replaced. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17′ | `engine.py` | `_eligible_truth(domain, target)` | REV 4's `_truth_projection`, renamed and **restricted to eligibility**; called by `authoritative()` and by enumeration |
| 17a | `engine.py` | `_proposition_records(prop_id, domain, target)` | `_eligible_truth` filtered by exact `proposition_id`; no fallback, no widening |
| 17b | `engine.py` | `_proposition_state(prop_id, domain, target)` | Per-identity state; `conflict` if a member is in an open group, else by count |
| 9″ | `engine.py` | `propositions(domain, target, *, actor)` | `state` and `records` from 17b/17a; partition invariant checked; `viewer`/`knowledge_conferred` |
| 10″ | `engine.py` | `possible_duplicates(domain, *, actor, target=None)` | Same gate; ids/counts only; `disclosure: "authorised_audit"`; no target in domain-wide mode |

`authoritative()` keeps its semantics; it calls `_eligible_truth` instead of an inline comprehension. That
is a refactor of one predicate, not a behaviour change.

**Untouched:** `eval/mem1/model.py` and the G1 harness; production Howl; r32; plan v3; Amendment 1; v9;
the map; the ledger; REV 1–4.

---

## 10. What this record does and does not claim

**Open:** U-6 / B3 — untouched. F7 replay — no facility, determinism unproven. R-B, R-C, R-D as before.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing was imported or compiled.
- No behavioural verification. Current-behaviour statements remain source inspection.
- No conformance proof. §8 checks coherence; it does not demonstrate conformance.
- No new rule, requirement, field or probe added to the ratified set.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9
  (`698254ca…`), the map (`b4ba33cc…`), the ledger (`da3e7239…`), REV 1 (`98447423…`), REV 2
  (`74ac6f77…`), REV 3 (`de4f68b6…`) or REV 4 (`1da08f2c…`).
- No transfer, sharing, destination interface, or cross-chat carryover.

**Sequence from here:** owner assessment of the complete proposal for ratification → separate
implementation authorization → implementation → conformance review → separate testing authorization.
