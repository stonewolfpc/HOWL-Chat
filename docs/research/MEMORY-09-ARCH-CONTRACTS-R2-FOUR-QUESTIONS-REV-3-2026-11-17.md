# MEMORY-09 — Architectural contracts, revision 3

**Status: PROPOSAL REV 3 — R-A confirmed, RAT-1/RAT-2 recorded as precedence clarifications, R-E resolved
into a contract. Not implemented, not tested.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 2 (`74ac6f77…`) **for §7, §8 R-A, §8 RAT-1, §8 RAT-2 and §8 R-E only** |
| Stands unamended | REV 2 in all other respects; REV 1 (`98447423…`) and REV 2 (`74ac6f77…`) preserved byte-identical |
| Authority | Owner confirmations on R-A, RAT-1, RAT-2, and the R-E requirement (recorded verbatim, §1) |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2 |
| Tests / probes | **0.** Nothing executed |
| Amends | **Nothing.** No edit to r32, v9, the 118-requirement map, the ledger, or any prior record |

---

## 0. What changed

| Item | REV 2 | REV 3 |
|---|---|---|
| **R-A** | Residual risk: two readings of O-8a-5 | **Confirmed** — shared structure, separate capabilities. Mechanism specified, including separate representation of default and delegated authority |
| **RAT-1** | Tension noted, resolved in prose | **Precedence clarification recorded**, dated and non-amending |
| **RAT-2** | Tension noted, resolved in prose | **Precedence clarification recorded**, dated and non-amending |
| **R-E** | Accepted false negative | **Contract**: creation, lookup, reuse, and multiplicity reporting. The false negative remains, but is bounded, surfaced and no longer unbounded |
| **D-A, D-B** | Specified | **Confirmed for correction**; still not applied — implementation remains unauthorised |

---

## 1. Owner confirmations, recorded verbatim

> "**R-A** — Recommend: confirm … O-8a-5 means shared authority structure and scoping, not identical
> capabilities for every operation. … One shared authority framework. The same domain-scoping and
> delegation infrastructure. Explicit operation-specific capabilities. Authoring permission does not
> imply promotion permission. Promotion permission does not imply authoring permission. Neither
> operation automatically grants the other."

> "One detail I'd make explicit: the human operator's default authority and delegated model capabilities
> must be represented separately. A model receiving promotion authority must not silently receive
> authorship authority. If authorship is delegated, that delegation must be explicit and recorded."

> "**RAT-1** — Recommend: record a precedence clarification … K10's latest-governs language cannot
> authorize automatic displacement when O-10b explicitly requires human-controlled authority changes.
> … We should document how that hierarchy applies here, not rewrite K10."

> "**RAT-2** — Recommend: record a precedence clarification … A single applicable promotion can settle
> authority. Multiple competing promotions without established incompatibility remain unresolved. An
> established equal-authority conflict remains a conflict. No automatic preference is introduced."

> "**R-E** — … Proposition identities must be persistent and domain-scoped. An actor should reuse an
> existing proposition identity when contributing to that same proposition. Creating a new identity must
> be explicit. The engine must never silently merge two identities based on text similarity. Potential
> duplicate identities should remain identifiable as an unresolved limitation, not be represented as
> proof that the propositions differ."

> "Neither clarification should modify the ratified r32, v9, or 118-requirement map. They should be
> recorded as dated, non-amending precedence clarifications under the existing hierarchy."

**Owner's caution, recorded because it governs what this file can claim:**

> "the reported 57/118 enumeration checks establish that the tables are complete. They do not establish
> that every proposed mechanism actually conforms. We still need the implementation conformance review
> later."

---

## 2. R-A — confirmed interpretation, and what it changes mechanically

### 2.1 Recorded confirmation

O-8a-5's "one shared authoring-authority regime" means **shared framework, scoping, delegation and
revocation infrastructure**, with **capability named per operation**. It does not mean that holding one
capability implies the other. This is an interpretation of an existing decision, not a new one, and it
amends nothing.

### 2.2 Mechanism, revised for the owner's detail

`Grant` carries two additions, not one:

| Field | Values | Why it must exist |
|---|---|---|
| `operations` | subset of `{"promotion", "authorship"}` | Capability is per operation; no implication either way |
| `origin` | `"domain_default"` or `"delegation"` | The owner required default authority and delegated capability to be **represented separately**; without this field a domain default and a delegated grant are indistinguishable in the record |

| Grant | `holder_kind` | `origin` | `operations` | Created by |
|---|---|---|---|---|
| Domain default | `human` | `domain_default` | `{promotion, authorship}` | declaring the domain |
| Delegated to a model | `model` | `delegation` | explicit at grant time; **defaults to `{promotion}` only** | the domain authority |

Rules that follow, each checkable without judgement:

1. A model grant that names only `promotion` confers no authorship. There is no fallback.
2. Delegating authorship is a distinct act: `grant(..., operations={"authorship"})`, recorded with
   `origin="delegation"`, `granted_by`, `granted_turn`.
3. Revocation is per grant, so revoking a promotion grant leaves any separate authorship grant intact.
4. `require(domain, actor, operation)` — the operation is now a required argument at every call site, so
   no operation can be reached through a grant that never named it.

**What this does not do:** it does not add a permission system. It is one field on an existing object,
one on its provenance, and a parameter on an existing check.

---

## 3. RAT-1 — precedence clarification: K10 versus O-10b

**Recorded 2026-11-17.** Non-amending. Does not edit K10, does not edit the map, does not create a
requirement. It documents how O-6a's hierarchy applies to a conflict that already existed.

**Subject.** K10 ("later promotions enter a chain; latest governs") and O-10b ("no actor may revoke or
supersede another actor's promotion, or their own prior promotion, merely by disagreeing with it;
resolution requires an explicit human edit").

**Determination.** Tier 1 governs. K10's "latest governs" describes the **supersession chain**, which is
created by the supersession operation. It does not authorise displacement by a bare promotion.

**Operational rule.**

| Step | Behaviour |
|---|---|
| 1 | A promotion establishes authority for **its own proposition** |
| 2 | A later promotion does **not** silently supersede an earlier one |
| 3 | Supersession requires the authorised **human** edit (`supersede_authority`) |
| 4 | The earlier promotion and its history remain preserved (rules 20, 44) |
| 5 | Two live promotions on one proposition are `unresolved_relation` until an assertion or a human edit changes that |

**Effect on conformance.** K10 remains satisfied through the supersession path; it is no longer
satisfiable by promotion alone. Recorded as a qualification of reading, not as a defect in K10 — and not
as a change to any file.

---

## 4. RAT-2 — precedence clarification: K7 versus O-10

**Recorded 2026-11-17.** Non-amending, same form as RAT-1.

**Subject.** K7 ("promotion settles authority; authority is no longer ambiguous") and O-10 (the
equal-authority conflict state, which is by definition unresolved).

**Determination.** Both hold. K7 describes the normal result of a promotion; O-10 describes the
exceptional state in which two legitimately promoted records coexist. Neither is a universal claim about
every query after every promotion.

**Operational rule.**

| Situation | Answer |
|---|---|
| One applicable promotion, no competing live record, no open group | `settled` |
| Several live promotions, incompatibility not established | `unresolved_relation` |
| Open conflict group | `conflict` |
| Nothing live | `none` |

No automatic preference is introduced in any branch.

---

## 5. R-E — proposition identity: creation, lookup, reuse

This is the section the owner asked for: **how identity actually works**, not a risk note.

### 5.1 The problem stated exactly

Two actors promoting claims about the same proposition may declare two identities. Automatic text
comparison is forbidden (rule 8, rule 34, §9). So the engine cannot detect the duplication, and if
identity were the *only* key to conflict detection, split identities would silently halve detection.

The contract therefore has to do three separate things: make **reuse the path of least resistance**, make
**creation explicit**, and make **multiplicity visible without asserting difference**.

### 5.2 Creation — explicit, and the only path

A proposition is minted as part of a promotion, never implicitly:

```
promote(..., proposition={"new": {"label": "king.status", "kind": "claim_about_world"}})
```

| Rule | Consequence |
|---|---|
| `proposition` is a **required** parameter | There is no default identity; omission is a refusal, not an inference |
| Minting requires a `label` | The actor's own identity string, recorded verbatim. It is **never compared**, never normalised, never used for similarity |
| `kind` and `about_character` are fixed at mint | A proposition cannot change kind later; a belief-ascription cannot become a world claim |
| Ids are monotonic and never reused | `prop-00007` remains `prop-00007`; there is no delete, no re-mint, no id recycling |
| Domain is part of the identity | A proposition belongs to one world or scenario; referencing it from another domain is refused |

### 5.3 Lookup — enumerate, then the actor chooses

The engine offers no similarity search and no ranking by text. It offers **enumeration by
applicability**:

```
engine.propositions(domain, target="king.status")
→ [ {id: "prop-00003", label: "king.status.alive", kind: "claim_about_world",
     created_by: "operator", created_turn: 4,
     state: "settled", live_records: ["truth-00012"], texts: [...]},
    {id: "prop-00007", label: "king.status.death",  kind: "claim_about_world",
     created_by: "model:canon", created_turn: 6,
     state: "unresolved_relation", live_records: ["truth-00019","truth-00024"], texts: [...]} ]
```

The actor reads the labels and the live texts and decides which identity is the one they mean. That
decision is a **declaration**, which is exactly what identity is; the engine does not make it, rank it,
or hint at it. Returning record text is presentation for a human or model reader — not comparison.

### 5.4 Reuse — by reference, and it is checked

```
promote(..., proposition={"id": "prop-00007"})
```

Refused when: the id is unknown; the id belongs to another chat (rule 31); the id belongs to another
domain (rule 35); or the promotion's declared applicability disagrees with the proposition's `target`.

Reuse is the cheap path by construction: the actor already has the id from a previous promotion or from
lookup, and passing it costs one parameter.

### 5.5 Multiplicity — visible, never interpreted

Two mechanisms, both non-semantic:

**At mint time.** The promotion result reports `coexisting_propositions: [ids]` — what else already
exists on that applicability. This is a *warning surface at the moment of creation*, not a block:
multiple propositions on one subject is the normal case under U-1 ("different propositions about the
same subject must not automatically conflict"), so blocking would forbid legitimate work.

**On review.** `engine.possible_duplicates(domain, target=None)` returns applicability groups holding
more than one proposition, with the count and the member ids. Its contract states, in the response
itself:

> **identity multiplicity — unresolved.** These identities share an applicability. The engine cannot
> determine whether they denote the same proposition. This report is **not** evidence that the
> propositions differ, and **not** evidence that a claim is missing.

Nothing is merged, scored, ordered by similarity, or auto-linked.

### 5.6 Why split identities do not halve conflict detection

This is the part that answers the owner's concern directly.

| Mechanism | Keyed on | Effect of a split identity |
|---|---|---|
| `unresolved_relation` | one `proposition_id` | **Lost** — two live records on two identities do not form one candidacy |
| **Assertion** (U-7) | records **and** applicability | **Not lost** — an assertion may name records from different propositions that share an applicability, so a real incompatibility remains expressible |
| `possible_duplicates` | applicability | **Not lost** — the multiplicity is reported |

So a split identity costs the *automatic* path to candidacy, not the *detectable* path. A domain
authority looking at `king.status` sees both identities, both live texts, and the multiplicity report,
and can assert incompatibility across them. What it costs is that nobody is told automatically that the
two are "the same claim" — which is the correct cost, because the engine cannot know that.

### 5.7 Deliberate non-capability: no merge, no relink

Merging two identities would mean re-pointing records from one id to another — rewriting authority
history (rule 20, F3, J16). It is **not proposed**, not stubbed, and not implied. If it is ever wanted,
it is an owner decision that must confront history rewriting explicitly.

The affordable alternative already exists: a human supersession can settle the subject, and the
assertion can record the disagreement; both preserve every original identity.

### 5.8 Worked example

| Turn | Actor | Action | Result |
|---|---|---|---|
| 4 | operator | mint `prop-00003` `king.status.alive`, promote "the king is alive" | settled |
| 6 | model:canon (promotion grant) | does **not** look up; mints `prop-00007` `king.status.death`, promotes "the king is dead" | both live; two identities; `coexisting_propositions` reported at mint |
| 7 | operator | `propositions(domain, "king.status")` | sees both, with both texts and both states |
| 8 | operator | `assert_incompatible(record_a="truth-00012", record_b="truth-00019", applicability=("king.status",), statement="…")` | one `conflict` group spanning **both** propositions; both preserved |
| 9 | human | `supersede_authority` on the older record | governing authority changes; both records and both identities retained |

Had the model reused `prop-00003` at turn 6, the state at turn 7 would have been
`unresolved_relation` on one identity — the automatic path. Both routes end at the same place once a
human edits.

### 5.9 R-E after this contract

| Before | After |
|---|---|
| "Accepted false negative" | **Bounded, surfaced, and non-silent**: multiplicity is reported at mint and on review; assertions span identities; nothing is merged or assumed |
| Undetectable duplication | Duplication is *undecidable by design* and *visible in practice*; the limitation is stated in the query contract itself |

R-E is not eliminated — it cannot be without semantic comparison, which is forbidden. It is now a
**declared limitation with an operating procedure**, which is the strongest honest form available.

---

## 6. Conformance of the new mechanisms — the check the owner asked for

Enumeration proved completeness; this checks coherence of the *new* items only.

| Mechanism | Rules | Requirements | Verdict |
|---|---|---|---|
| `proposition` required, explicit mint | 6 (applicability untouched), 26 (nothing inferred), 34 (no assessment) | K2, K9, L4 | **C** — identity is declared, not derived |
| Enumerated lookup, no similarity | 8, 34 | K2, G2, F7 (deterministic order by id) | **C** — ordering is by id, never by score |
| `coexisting_propositions` at mint | 34 | K2, K9 | **C** — reported, not acted upon |
| `possible_duplicates`, multiplicity text | 8, 13, 34 | J12, J13, F1b | **C** — a report, not a verdict; no record changes |
| Assertion spanning propositions | 41, 34, 21 | J13, J14, J16 | **C** — enables the state; no preference |
| No merge / no relink | 20, 44, 45 | F3, J16, K27 | **C** — absence, not capability |
| `Grant.operations` + `origin` | 24, 35, 39 | L7, L10, K21 | **C** — capability and provenance, per domain |

**No new contradictions.** The three prior residuals change status: R-A **closed** (confirmed), RAT-1 and
RAT-2 **recorded** as clarifications, R-E **bounded**. Remaining residuals: R-B (definitions/B3), R-C
(belief voiced in dialogue), R-D (F7 replay).

**Conformance tallies are unchanged from REV 2** — this revision adds no mechanism that shifts a
requirement's code: rules **39 U · 13 M · 5 N · 0 R · 0 X**, requirements **86 U · 22 M · 8 N · 1 R · 1
X**. F7 remains the only **R**; B3 the only **X**.

---

## 7. Updated implementation list (still unapplied)

### 7.1 Corrections confirmed, awaiting implementation authorisation

| # | File | Symbol | Change |
|---|---|---|---|
| 1 | `engine.py` | `promote` | Delete `supersede` and its displacement block (D-A, RAT-1) |
| 2 | `engine.py` | `retract` | Refuse belief records |
| 3 | `engine.py` | new `revoke_belief` | Model-gated, append-only revocation (D-B, D5) |
| 4 | `engine.py` | `promote` | Refuse a belief citation unless the proposition is a `belief_ascription` of that owner |
| 5 | `engine.py` | `acquire` | Evidence table; own-account exempt; inference basis owned; no truth reads |
| 6 | `retrieval.py` | `world` | `unresolved_relation` for ≥2 live records; never `settled` |

### 7.2 Structural work, now with owner direction

| # | File | Symbol | Change | Direction |
|---|---|---|---|---|
| 7 | `engine.py` | `_propositions` registry | Identity, `kind`, `label`, `about_character`, domain | U-1, U-4, R-E |
| 8 | `records.py` | `Record` | `proposition_id`, required on truth classes | U-1 |
| 9 | `engine.py` | `propositions(domain, target)` | Enumerated lookup; id-ordered | R-E |
| 10 | `engine.py` | `possible_duplicates(domain, target)` | Multiplicity report with the non-evidence disclaimer | R-E |
| 11 | `engine.py` | `promote` | Required `proposition` (id or mint); `coexisting_propositions` in the result | R-E |
| 12 | `engine.py` | new `assert_incompatible` | Assertion with applicability, spanning propositions | U-7 |
| 13 | `engine.py` | `_promotions`, `snapshot` | Persist groups, propositions, authorships | U-2, R-E |
| 14 | `retrieval.py` | new `belief_facts` | Belief-fact view; world excludes ascriptions | U-4 |
| 15 | `authority.py` | `Grant`, `require` | `operations` + `origin`; per-capability revocation | R-A |
| 16 | `engine.py` | new `author_ambient` | Authoring capability; `claim` + `declaration` | U-5, U-8 |
| 17 | `__init__.py` | `__all__` | Export `assert_incompatible`, `author_ambient`, `revoke_belief` | 3, 12, 16 |

**Untouched:** `eval/mem1/model.py` and the G1 harness; production Howl; r32; plan v3; Amendment 1; v9;
the map; the ledger; REV 1; REV 2.

---

## 8. What remains, and what this record does not claim

**Open:** U-6 / B3 (event versus definition) — deliberately unresolved, untouched. F7 replay — no facility
exists, determinism unproven. R-B, R-C, R-D as described.

**Not claimed:**

- No implementation. No source file changed; D-A and D-B are confirmed **and still uncorrected**.
- No execution. No test, probe or suite ran.
- No behavioural verification. Current-behaviour statements remain source inspection.
- No conformance proof — the enumeration checks prove the impact tables are **complete**, not that any
  mechanism conforms. That is the later implementation conformance review's job, and the owner's caution
  is recorded in §1.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9, the map, the
  ledger, REV 1 (`98447423…`) or REV 2 (`74ac6f77…`).
- No transfer, sharing, destination interface, or cross-chat carryover.

**Sequence from here:** owner review of REV 3 → implementation authorization → implementation →
conformance review → separate testing authorization.
