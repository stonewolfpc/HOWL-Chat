# MEMORY-09 — Architectural contracts, revision 4

**Status: PROPOSAL REV 4 — proposition enumeration bounded by character access. Not implemented, not
tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 3 (`de4f68b6…`) **for §5.3, §5.5, and the affected rows of §6 and §7 only** |
| Stands unamended | REV 3 in every other respect; REV 1 (`98447423…`), REV 2 (`74ac6f77…`) and REV 3 preserved byte-identical |
| Authority | Owner's enumeration/access-control instruction, recorded verbatim (§1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed |
| Amends | **Nothing.** No edit to r32, v9, the 118-requirement map, the ledger, or any prior record |
| Ratification | **Not ratified.** This remains a proposal; implementation authorisation is a separate gate |

---

## 0. What changed

| Item | REV 3 | REV 4 |
|---|---|---|
| `propositions(domain, target)` | Returned every identity with `live_records` and `texts`, to any caller | **Authority-gated, target-required, projection-only.** Identities and applicability metadata for the caller; only records the caller may already retrieve; no private record, ever |
| `possible_duplicates` | Applicability groups with counts and member ids | Same purpose, **same gate**, and in domain-wide mode **ids and counts only** — no labels, no texts, no `about_character` |
| `coexisting_propositions` | Ids reported at mint | Ids **and count** only; no labels, no texts |
| Private records and propositions | Not addressed | **Structurally unlinkable.** A non-truth record cannot carry `proposition_id`; the edge does not exist |
| Record projection | Implicit | One shared predicate with `world()`, named, so the two surfaces cannot drift apart |

No rule, requirement, decision or tally changes. This revision constrains a proposed mechanism; it does
not add to the ratified set.

---

## 1. Owner instruction, recorded verbatim

> "Proposition enumeration must preserve character access boundaries.
>
> REV 3 says `propositions(domain, target)` returns every identity with labels, live texts, and states.
> That must not become a route around O-1 or D-02.
>
> An actor must not gain access to another character's private record merely because they can enumerate
> propositions in the same domain. The implementation contract needs to distinguish:
>
> - proposition identity and applicability metadata;
> - records the requesting actor is authorized to access;
> - private content that must remain inaccessible.
>
> The identity can be visible without making every associated record's text visible. This is a necessary
> access-control clarification, not a reason to redesign R-E."

**Scope of this revision:** the enumeration boundary only. R-E's creation, reuse, multiplicity and
no-merge provisions stand as written in REV 3 §5.

---

## 2. The defect in REV 3 §5.3, stated exactly

REV 3 returned, for every identity: `label`, `state`, `live_records` and `texts`. Three distinct leaks
were possible from that one shape, and they need three distinct answers.

| # | Leak | Why it violates the ratified set |
|---|---|---|
| L-1 | **Private text.** If any record reachable from an identity were private, its text would be returned | Rule 5 and O-1: a private record is retrievable only by its owner. H2: abstains in every context. G3 and D7: no leak into a non-owner's context |
| L-2 | **Private identity.** A record id is a handle; a handle is an access path if anything resolves it | J8: an origin descriptor must resolve to nothing. Handing out an id invites a resolver |
| L-3 | **Private existence.** A label or target authored by one actor can describe content another actor holds privately; enumeration would confirm that the subject exists and is being worked on | D-02 behaviour 8 (existence, who knows, and how are separable), rule 32 (reach is not knowledge), rule 6 (applicability is never ownership or knowledge) |

L-3 is the one most easily dismissed, and it is the one that makes "just return metadata" insufficient
on its own. A proposition registry is not a private record, so nothing in v9 forbids reading it — but
nothing in v9 entitles a *character* to it either. Rule 6 is explicit that applicability says where
something applies, **never who owns or knows it**; the converse also holds, and an enumeration surface
must not become a knowledge channel by accident.

---

## 3. The audience decision: enumeration is not a perspective surface

**Decision.** Proposition enumeration is an **authoring and adjudication surface**, not a character
retrieval surface. It is available to an actor holding a live `promotion` capability in that domain, and
to no one else.

| Caller | Result |
|---|---|
| Actor with a live `promotion` grant in `domain` | L1 identities + L2 authorised records (§4) |
| Actor with only `authorship`, or no grant | **Refused** — `no_domain_authority` (§7) |
| A character perspective (`perspective` / `beliefs_of` / `accounts_of`) | **Unchanged** — those surfaces never call enumeration, and enumeration gains no character entry point |
| Narrator | No path; narration may not reproduce or derive content (rule 47, H7) |

**Why `promotion` and not `authorship`:** minting and reuse are properties of promotion (REV 3 §5.2,
§5.4). Ambient authorship creates claims and declarations, not propositions. Requiring the operation that
actually uses the registry keeps capability and purpose aligned, and it is checkable at one call site.

**Why authority is not access.** `require(domain, actor, "promotion")` proves capability **in a domain**.
It proves nothing about any character's private memory. O-1's exclusivity is a property of the record,
not of the reader's rank, so a promotion grant must not widen record eligibility by one row. §4.3 makes
that structural rather than aspirational.

**Why this is not a redesign of R-E.** Creation, reuse, multiplicity reporting and no-merge are
untouched. Only the audience of the lookup changes, and the lookup was never the part the owner asked to
preserve — reuse by id remains the cheap path for every authorised actor.

---

## 4. Three layers, three answers

The owner named three things that must be distinguished. They are three layers of one response, and each
has its own rule.

### 4.1 L1 — identity and applicability metadata

| Field | Returned | Notes |
|---|---|---|
| `id` | yes | `prop-00007`; opaque, never derived from text |
| `label` | yes | The minter's identity string, verbatim. Actor-authored text; see §4.4 |
| `kind` | yes | `claim_about_world` / `belief_ascription` / … (closed set, REV 2) |
| `about_character` | yes, when the proposition is an ascription | The subject of an **authoritative** record the caller may read anyway |
| `domain`, `target` | yes | Domain is part of the identity (REV 3 §5.2) |
| `created_by`, `created_turn` | yes | Provenance already required of promotions (K9) |
| `state` | yes | `settled` / `unresolved_relation` / `conflict` / `none`, computed as today |
| **any record id** | **no** | L2 only |
| **any record text** | **no** | L2 only |

L1 is metadata about **declared identity**. It contains no record handle and no record content.

### 4.2 L2 — records the caller may access

L2 is a **projection through an existing eligibility predicate**, not a new query over the record store:

```
_truth_projection(domain, target) ->
    [r for r in self._records.values()
     if r.truth_apt and r.owner_character is None
     and r.domain == domain and r.target == target and r.status == ACTIVE]
```

| Property | Consequence |
|---|---|
| `world()` and enumeration call **the same function** | A change to one cannot create a second, divergent visibility rule |
| `owner_character is None` **and** `truth_apt` | Only domain records; no account, claim or belief can pass |
| Active only | Superseded and revoked records stay in history, not in enumeration |
| Texts returned are the same bytes `world()` returns | Enumeration adds **no text surface**; it only groups what an authorised truth query already returns |

L2 therefore cannot be wider than `world()`. It is narrower in one respect: it is grouped by proposition
identity, which `world()` does not do.

### 4.3 L3 — private content: absent by construction, not by filter

The strongest available guarantee is that the edge does not exist:

| Rule | Effect |
|---|---|
| `Record.proposition_id` is **required** on `world_truth` / `scenario_truth` | Every authoritative record is attached to one identity (REV 3 item 8) |
| `Record.proposition_id` is **refused** on `claim`, `account` and `belief` | A private record cannot be attached to a proposition, so enumeration has nothing to traverse |
| No operation sets `proposition_id` after construction | No late attachment, no backfill, no migration path |
| The proposition registry stores `label` and `kind`, never a record reference | The registry cannot be used as an index into private records |

This is why L3 needs no filter. A filter can be forgotten at a call site; an absent field cannot. It also
means a proposition is **not** a container: it is a declared identity that authoritative records point
at, and only authoritative records are in the graph.

**Consequence for belief ascriptions.** A `belief_ascription` proposition is attached to authoritative
records of the form "Sable believes P" — records a promotion created and an authorised actor signed. It
is **not** attached to Sable's belief record, which stays private, unlinked and unenumerable. O-5's
distinction survives mechanically: the fact *about* Sable is domain canon; Sable's mind is not.

### 4.4 The one residual honesty: labels are actor-authored

A label is declared by the actor who minted the identity. The engine must not and cannot judge whether a
label describes private content (rule 8, rule 34). Therefore:

- labels are returned **only** to the domain-authority audience, never to a character context;
- a label is never propagated into a record, a disclosure, a coverage declaration or narration;
- the contract does **not** claim label text is safe to treat as content. What it claims is that the
  audience is an authority holder in that domain, who already reads that domain's canon.

This is a stated limitation, not a guarantee. Pretending otherwise would require semantic judgement the
architecture forbids.

---

## 5. API contracts

### 5.1 `propositions(domain, target, *, actor)`

`target` is **required**. There is no subject-discovery mode on this surface.

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
        "state": "unresolved_relation",
        "records": [ { "record_id": "truth-00019", "text": "the king is dead", … },
                     { "record_id": "truth-00024", "text": "the king still lives", … } ] } ] }
```

| Field | Meaning |
|---|---|
| `viewer` | Constant `"authority"`. The response declares its own audience so a caller cannot mistake it for a perspective read |
| `knowledge_conferred` | Constant `false`. Rules 18 and 32: reading canon is not acquiring it. No account is created, and the result must not be injected into a character's context |
| `records` | L2 only. Empty when no authoritative record is active on that target |

**Target required, and why.** A wildcard would turn enumeration into subject discovery: a list of every
target in the domain, including subjects the caller never thought to query. R-E needs reuse, not a
catalogue. Duplicate *review* is served by §5.2 instead, which returns no subject inventory.

### 5.2 `possible_duplicates(domain, *, actor, target=None)`

| Form | Returns | Withheld |
|---|---|---|
| `target=<t>` | Groups on `<t>` with ≥2 identities: member `id`s and `count` | labels, texts, `about_character` |
| `target=None` | The same shape, across the domain | labels, texts, `about_character`, targets with a single identity |

Every response carries the REV 3 §5.5 disclaimer verbatim:

> **identity multiplicity — unresolved.** These identities share an applicability. The engine cannot
> determine whether they denote the same proposition. This report is **not** evidence that the
> propositions differ, and **not** evidence that a claim is missing.

**The domain-wide mode is deliberately narrow.** It names applicability groups that already hold two or
more declared identities — that is, subjects with more than one authored proposition. It returns no
labels, no texts and no character names, so it cannot be read as a content index. It reveals only that
multiplicity exists somewhere in the domain, which is precisely the audit signal the owner asked for in
REV 3 §5.6. An actor who wants detail must name the target, which is itself an applicability it could
already query through `world()`.

### 5.3 `coexisting_propositions` at mint

Returned to the promoting actor only, as `{"count": n, "ids": [...]}`. No labels, no texts. It is a
warning at the moment of creation; the actor can enumerate the target under §5.1 if it wants detail.
Nothing is blocked: multiple propositions on one subject is the normal case under U-1.

### 5.4 Reuse — unchanged, and still checked

`promote(..., proposition={"id": "prop-00007"})` remains refused for an unknown id, a foreign chat
(rule 31), a foreign domain (rule 35), or an applicability that disagrees with the proposition's
`target`. Enumeration is a convenience for finding the id, never a substitute for these checks: reuse is
validated against the registry, not against what the caller claims to have seen.

---

## 6. Why this does not weaken R-E

The owner's concern in REV 2 was that split identities could halve conflict detection. Nothing here
changes that analysis; it changes who may look.

| Mechanism | Keyed on | Effect of split identity | Gate |
|---|---|---|---|
| `unresolved_relation` | one `proposition_id` | Detected automatically only within one identity | none (engine state) |
| Assertion (U-7) | records **and** applicability | Spans identities; incompatibility remains expressible | actor's existing authority to assert |
| `possible_duplicates` | applicability | Multiplicity reported | domain `promotion` grant |
| `propositions` | applicability | Identity and canon visible | domain `promotion` grant |

A domain authority still sees both identities, both live texts and the multiplicity report, and can
assert across them. What changes is that a **character** can no longer obtain any of it — which is the
point, and which costs R-E nothing, because characters were never the audience for adjudication.

---

## 7. Refusals, fail-closed

| Situation | Result |
|---|---|
| No live grant in `domain` | `refused` / `no_domain_authority`; **no identities, no counts, no ids** |
| Grant in a different domain (rule 35) | `refused` / `no_domain_authority` |
| Caller asks for a character view of enumeration | `refused` / `not_a_perspective_surface` |
| `target` omitted on §5.1 | `refused` / `target_required` |
| Unknown target, valid authority | `ok`, empty `propositions` (abstention is not an error) |
| Proposition whose records are all superseded/revoked | L1 identity present, `state: none`, `records: []` |

A refusal returns nothing derived from the registry — not even a count. A count is information under
D-02 behaviour 8, and returning one before checking authority would leak by anticipation.

---

## 8. Worked example — the leak REV 3 allowed, and the answer

State: Sable privately holds an account that the queen poisoned the king. No one has promoted anything
about it. Canon holds two competing promotions on `king.status`.

| Turn | Actor | Call | REV 3 result | REV 4 result |
|---|---|---|---|---|
| 7 | Sable (character, no grant) | `propositions(world, "king.status")` | both identities **and both texts** — Sable is handed canon it never acquired, and, if any linked record were private, that too | `refused / no_domain_authority`; Sable's knowledge is unchanged and its private account is not in the response |
| 8 | model:canon (promotion grant) | `propositions(world, "king.status")` | same list | two identities, two truth texts, no private record, `knowledge_conferred: false` |
| 9 | model:canon | `possible_duplicates(world)` | targets + labels across the domain | groups with ≥2 ids and counts; no labels, texts or character names |
| 10 | model:canon | `promote(..., proposition={"id":"prop-00007"})` | accepted | accepted — reuse is validated against the registry |
| 11 | operator | `assert_incompatible(record_a="truth-00012", record_b="truth-00019", applicability=("king.status",), …)` | one `conflict` group spanning both | unchanged |

**Sable's private account never appears at any turn.** It has no `proposition_id` (§4.3), so it is not in
the graph; and Sable has no enumeration path (§3). Two independent reasons, either sufficient.

---

## 9. Conformance of the boundary itself

This checks coherence with existing requirements; it is **not** conformance evidence, and no row here is
a new requirement.

| Boundary property | v9 rules | Map requirements | Verdict |
|---|---|---|---|
| No private record in any enumeration result | 5, 6, 47 | H2, G3, D7, J7, J8 | **C** — structurally absent, not filtered (§4.3) |
| Authority gates the surface; capability ≠ access | 24, 32, 35 | K8, K21, L7 | **C** — grant is domain-scoped and confers no record entitlement |
| No account created by reading canon | 18, 32 | K4, D-02 behaviour 1 | **C** — `knowledge_conferred: false`, no acquisition path |
| Texts identical to `world()` | 7, 41 | J14, H3, K12 | **C** — shared predicate; no new text surface |
| Multiplicity report is not a verdict | 8, 34 | J12, J13 | **C** — disclaimer in the response; no record changes |
| Consumption order independent of the caller | 34 | F7 | **C** — ordering by proposition id, then record id |
| Unknown / unauthorised refused without widening | 5 | F4, F5, B9 | **C** — refusal precedes any registry read |

**Tallies unchanged from REV 3:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M ·
8 N · 1 R · 1 X**. F7 remains the only **R**; B3 the only **X**. This revision adds no mechanism that
shifts a requirement's code.

**Residuals unchanged:** R-A closed; RAT-1 and RAT-2 recorded; R-E bounded and now audience-constrained;
R-B (definitions/B3), R-C (belief voiced in dialogue), R-D (F7 replay) open.

---

## 10. Implementation list — deltas only

REV 3 §7 stands, with these four rows amended. Nothing is applied.

| # | File | Symbol | Change | Direction |
|---|---|---|---|---|
| 8′ | `records.py` | `Record` | `proposition_id` required on truth classes, **refused on `claim`/`account`/`belief`**, never assignable after construction | O-1, rule 5, §4.3 |
| 9′ | `engine.py` | `propositions(domain, target, *, actor)` | `target` required; `authority.require(domain, actor, "promotion")` first; L1 + `_truth_projection`; `viewer`/`knowledge_conferred` in the result | §3, §5.1 |
| 10′ | `engine.py` | `possible_duplicates(domain, *, actor, target=None)` | Same gate; ids and counts only; disclaimer in every response | §5.2 |
| 11′ | `engine.py` | `promote` | `coexisting_propositions` as `{count, ids}`; no labels or texts | §5.3 |

One new internal symbol, shared with the existing surface:

| # | File | Symbol | Change |
|---|---|---|---|
| 17 | `engine.py` | `_truth_projection(domain, target)` | The single active-truth predicate, called by both `world()`'s backing `authoritative()` path and enumeration |

**Untouched:** `eval/mem1/model.py` and the G1 harness; production Howl; r32; plan v3; Amendment 1; v9;
the map; the ledger; REV 1; REV 2; REV 3.

---

## 11. What this record does and does not claim

**Open:** U-6 / B3 (event versus definition) — untouched. F7 replay — no facility exists, determinism
unproven. R-B, R-C, R-D as before.

**Not claimed:**

- No implementation. No source file changed. D-A and D-B remain confirmed **and uncorrected**.
- No execution. No test, probe or suite ran; nothing was compiled or imported.
- No behavioural verification. Statements about current behaviour are source inspection.
- No conformance proof. §9 checks coherence of a proposal; it does not demonstrate that any mechanism
  conforms. That is the later conformance review's job.
- No new rule, requirement, field or probe added to the ratified set. The boundary is an implementation
  constraint on a proposed mechanism.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9
  (`698254ca…`), the map (`b4ba33cc…`), the ledger (`da3e7239…`), REV 1 (`98447423…`), REV 2
  (`74ac6f77…`) or REV 3 (`de4f68b6…`).
- No transfer, sharing, destination interface, or cross-chat carryover.

**Sequence from here:** owner ratification decision on the architectural proposal → separate
implementation authorization → implementation → conformance review → separate testing authorization.
