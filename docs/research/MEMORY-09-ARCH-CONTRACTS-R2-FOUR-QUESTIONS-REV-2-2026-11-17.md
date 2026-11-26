# MEMORY-09 — Architectural contracts, revision 2

**Status: PROPOSAL REV 2 — revised around the owner's eight recommendations. Not implemented, not tested.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | `MEMORY-09-ARCH-CONTRACTS-R2-FOUR-QUESTIONS-2026-11-17.md` (REV 1, `98447423…`) **for the eight decided points only** |
| REV 1 status | **Preserved unamended** — 44 152 B, `98447423cbab06cc…`. Not edited, not withdrawn |
| Authority | Owner **recommendations** on U-1…U-8 (recorded verbatim, §1). This file decides nothing further |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2; no code written |
| Tests / probes | **0.** Nothing executed |
| Reports | Two implementation defects (owner-confirmed) specified for correction — **not corrected** |

---

## 0. Delta from REV 1

| Point | REV 1 | REV 2 |
|---|---|---|
| **U-1** | Read `target` as a proposition identifier (no new data) | **Explicit proposition identity** — a separate `Proposition` entity; `target` keeps its rule-6 meaning untouched |
| **U-2** | K7/K10 read proposed as a question | **Human-only supersession**; `promote()` cannot displace anything; K10 read as the supersession chain |
| **U-3** | `experienced` flagged as unresolved | **Canonical names preserved**; `experienced` is a concept label for `witnessed`; `narrator_grant` determined to be an acquisition mechanism (r32 K-R11) |
| **U-4** | Optional `about` parameter | **No `about` parameter.** Proposition `kind` (closed two-value enum) carries the distinction, with four demonstrations |
| **U-5** | `claim` + `declaration` reuse proposed, flagged | Reuse **tested against the semantics**; adopted because ambient material is assertoric, with the failure condition named (definitions — B3) |
| **U-6** | B3 unresolved | **Unchanged — still unresolved.** Nothing added |
| **U-7** | "who may assert" left open | **Domain authority**, with mandatory attribution and applicability identification |
| **U-8** | Single authority gate for authorship and promotion | **Capabilities split**: `Grant.operations`; authoring ≠ promotion in both directions |
| Defects | Found, not fixed | **Specified for correction** under existing authority; still not applied |

---

## 1. Owner recommendations, recorded verbatim

The owner stated these as **recommendations**, not as registered decisions. They are recorded in that
form and used as the governing direction for this revision.

> "I would not authorize implementation yet."

> "**U-1** — Recommend: explicit proposition identity … A proposition needs an explicit, stable
> identity. Records referencing that identity concern the same proposition. Applicability determines
> whether their claims overlap. Different propositions about the same subject must not automatically
> conflict. The engine does not compare natural-language content to decide truth or incompatibility."

> "**U-2** — Recommend: human-only supersession … A new incompatible promotion creates or preserves an
> unresolved authority situation. Supersession requires an explicit human edit."

> "**U-3** — Recommend: preserve canonical names; resolve mapping … Do not silently rename `witnessed`
> to `experienced`. Also, `narrator_grant` must not disappear simply because it wasn't in the
> five-item breakdown."

> "**U-4** — Recommend: yes, if it is the minimum structural safeguard … One caveat: don't let the
> `about` field become a general-purpose mechanism for inventing new truth categories."

> "**U-5** — Recommend: reuse only if semantics genuinely match … If they cannot, propose the smallest
> necessary structural distinction and explain why."

> "**U-6** — Recommend: leave unresolved. Do not solve this by quietly expanding the contract."

> "**U-7** — Recommend: domain authority, with explicit attribution … The engine records the authorized
> assertion; it does not invent one."

> "**U-8** — Recommend: yes, but keep it non-promotional … authoring authority is not promotion
> authority … All contributions remain attributable."

**The standard this is measured against (owner, verbatim):**

> "The most important thing here is that we're not trying to build the cleverest memory engine
> possible. We're trying to build one that behaves correctly when characters disagree, when someone
> lies, when beliefs change, when canon changes, and when the model has incomplete information."

---

## 2. Defect corrections (specified, not applied)

Both are owner-confirmed implementation defects. Neither needs new architecture, and **neither is
implemented by this file**.

### D-A — `promote(supersede=True)` bypasses the human gate

**Defect.** `promote` accepts `supersede=True`, marks prior authoritative records `SUPERSEDED` and sets
`superseded_by`, without calling `_require_human`. A model holding a standing promotion grant (8a-6)
could therefore displace an existing promotion — including its own — with no human edit. O-10b says
exactly the opposite, and it is tier 1.

**Correction.** Remove `supersede` from `promote` entirely. Displacement is available only through
`supersede_authority`, which already requires a human holder. No new permission system, no new
operation, no new field.

**Consequence.** Every promotion creates a record and displaces nothing. Two live promotions on one
proposition produce `unresolved_relation` (§3), which is the state O-10 describes as unresolved.

### D-B — `retract()` can revoke another character's belief

**Defect.** `retract(record_id, turn)` sets `status = REVOKED` on any non-truth record with no ownership
or role check. A belief can therefore be revoked by anyone, contradicting O-5 ("only that character's
own model may declare, revise, or revoke") and rule 53.

**Correction.** `retract` refuses belief records outright; belief lifecycle moves to model-gated
operations:

| Operation | Gate | Effect |
|---|---|---|
| `declare_belief` | `role == "character_model"` and `declarer == character` | new record; prior beliefs `superseded`, retained (rule 55) |
| `revoke_belief` | same gate | **new** belief record carrying the revocation; prior `superseded_by`; nothing deleted (D5) |

The second row also closes a smaller gap: D5 expects explicit belief revocation to create a new record
with the prior retained, while `retract` only mutated status. Both defects are in the same code path.

**Explicitly refused:** an operator or narrator acting "with authority" through `retract`. Authority in
a domain is promotion authority; it confers nothing over a character's interior state (O-5).

---

## 3. Contract C-1 (revised) — explicit propositions, and conflict

### 3.1 Proposition identity

A **Proposition** is a claim identity, not a record and not a memory. It is the thing two promotions can
share, and the reason two statements can be known to concern the same claim without comparing their
wording.

| Field | Meaning |
|---|---|
| `proposition_id` | stable identity, minted once, never reused |
| `domain` | the world or scenario it applies in |
| `target` | **applicability**, unchanged in meaning from rule 6 |
| `kind` | `claim_about_world` or `belief_ascription` — **closed enum, two values** |
| `about_character` | required iff `kind == belief_ascription`; otherwise absent |
| `created_by`, `created_turn` | attribution of the identity itself |

Propositions live in a `_propositions` registry alongside the engine's existing `_promotions` and
`_declarations` registries — **not** in the record set, so **no new record class and no §3 amendment**.
Truth records gain one field, `proposition_id`.

A promoting actor either **references** an existing proposition or **mints** one as part of the
promotion. Minting is an act of declaration — "these are the same claim" is the actor's statement, never
the engine's inference. The engine never compares text.

### 3.2 Candidacy, overlap, and the four states

| Concept | Rule |
|---|---|
| **Same proposition** | candidacy — two live authoritative records on one `proposition_id` |
| **Overlapping applicability** | same `domain` + same `target`, different propositions — *never* a conflict by itself |
| **Incompatibility** | asserted by a domain authority (§3.3) |
| **Content comparison** | never performed |

| Live candidacy on a proposition | State | Authority |
|---|---|---|
| 0 records | `none` | — |
| 1 record | `settled` | K7 |
| ≥2 records, no assertion | `unresolved_relation` | O-10 (unresolved), U-2 |
| open conflict group | `conflict` | rule 41, J13, J14 |

Different propositions with the same applicability are answered separately. This is the owner's U-1
requirement — "different propositions about the same subject must not automatically conflict" — made
structural rather than conventional.

### 3.3 The assertion (U-7)

An assertion is made by an actor holding authority in the applicable domain and must carry:

| Requirement (owner) | Field |
|---|---|
| "Identify the records or proposition involved" | `propositions[]`, `records[]` |
| "Identify the overlapping applicability" | `domain`, `targets[]` (the shared applicability) |
| "Be attributable to the asserting actor" | `asserted_by`, `authority_ref`, `turn` |
| "Preserve the assertion as an assertion, not as an independently verified truth" | `statement` verbatim; stored in the group, never in a truth field |
| "Leave the underlying records intact" | no record touched; membership only |

The group persists with `status` `open`/`resolved`, and closure only by human supersession or revocation
(rule 43, O-10b). The engine records the assertion; it never invents one, and never promotes the
assertion into a verdict.

**Refusals:** no domain authority · `record_a == record_b` · a member is not truth-apt · a member is in
another chat · the pair already shares an open group · a member is not `ACTIVE`.

### 3.4 K7 / K10 (U-2)

- `promote()` displaces nothing (D-A). Its own proposition becomes governed by the new record.
- Displacement is `supersede_authority`, human-gated, which creates the K10 chain.
- Therefore K10's "latest governs" is read as **the supersession chain governs**, not "any later
  promotion wins". This is a **precedence application**, not a new rule: tier 1 (O-10b) outranks tier 3
  (rule 20/K10). See §6, finding **RAT-1**.

---

## 4. Contract C-2 (revised) — belief and world truth

### 4.1 Mechanism

`about` is gone as a free parameter. The distinction lives on the **proposition**, where it cannot be
omitted, defaulted, or repurposed:

- `kind = claim_about_world` — answers world/scenario truth queries for its applicability.
- `kind = belief_ascription` + `about_character` — answers **only** "what does <character> believe"
  queries. It is never returned by a world-truth query.

A promotion whose basis cites a belief record must reference a `belief_ascription` proposition whose
`about_character` equals that belief's owner. Anything else is refused at write time.

### 4.2 The four demonstrations the owner required

| Required demonstration | How it holds |
|---|---|
| **Cannot be omitted** when promoting a belief-related record | A truth record cannot be constructed without `proposition_id` (record validation); a promotion citing a belief must reference an ascription proposition, or it is refused |
| **Cannot be used to promote the underlying belief as world truth** | The ascription proposition's `kind` excludes it from world answers; a `claim_about_world` proposition carrying a belief citation is refused at write |
| **"Sable believes P" can still be promoted as a fact about Sable** | An ascription proposition with `about_character="sable"` is authoritative, and answered by the belief-fact view |
| **World-truth retrieval excludes belief records** | `world()` filters on `kind`; only `claim_about_world` propositions are eligible |

**The caveat honoured.** `kind` is closed at two values with no minting path for a third: an unknown kind
is refused at construction. It cannot become a general-purpose mechanism for inventing truth
categories, because there is nothing to extend — any new kind is an owner decision, not a parameter.

**Residual, stated not hidden.** An actor promoting `P` as world truth citing a *claim* in which a
character voiced a belief is undetectable: the claim is not a belief record, and inspecting its text is
the adjudication rule 8 forbids. The actor is then asserting `P` on their own authority, which is a
different act from deriving it from a belief. This residual shrinks only if the conversation layer ever
marks such claims as belief-reports — a change outside MEMORY-09's scope.

---

## 5. Contract C-3 (revised) — evidence, with canonical names preserved (U-3)

### 5.1 Naming

| Owner's concept | Canonical kind | Treatment |
|---|---|---|
| "experienced" | `witnessed` | **Concept label only.** No enum change; `experienced` never appears as a `via` value in the contract |
| told | `told` | unchanged |
| read | `read` | unchanged |
| inferred | `inferred` | unchanged |
| verified | `verified` | unchanged |
| *(not listed)* | `narrator_grant` | **Retained** — see §5.2 |
| *(exempt)* | `none` | own accounts only (rule 51) |

`witnessed` is preserved because `via` is a **fixed total order** in tier 2 (r32 K-R26:
`witnessed > verified > read > told > inferred > narrator_grant`). Adding `experienced` as a value would
insert a term into that order, which no one is authorised to edit. Mapping the concept onto the
canonical name costs nothing and changes no tier-2 byte.

### 5.2 `narrator_grant` — determination

The owner asked whether it is "an acquisition mechanism, a separate grant mechanism, or something else
under the existing rules". Source inspection of r32 answers it:

| Evidence | Text |
|---|---|
| **K-R11** | "The edge fold depends **only** on acquisition-family events for (character, memory): `witnessed \| told \| read \| inferred \| verified \| narrator_grant \| revoked`" |
| **K-R26** | it is one of the six ranked `via` kinds |
| **K-R46** | destination knowledge is "a **new** event with its own `via` (`read \| told \| narrator_grant`)" |

**Determination:** `narrator_grant` is an **acquisition mechanism** — an acquisition-family event kind
whose evidence is the grant or declaration act, not a character's own experience. It is **not** the
authority grant of 8a-6; that is a different object entirely (who may promote or author). The two share
a word and nothing else.

**Name separation (following the O-4 precedent for `Recollection` vs `own_account`):**

| Term | Meaning |
|---|---|
| **authority grant** (`Grant`, 8a-1/8a-6) | who may promote or author in a domain |
| **`narrator_grant`** (`via`) | how information entered a character's knowledge, by an authorised act |

### 5.3 Evidence table

| Kind | Evidence required at write | Refused when |
|---|---|---|
| `witnessed` | acquisition event reference (`span`) attributable to the character | no span; character absent from the event |
| `told` | disclosure event reference and `descriptor.speaker` ≠ owner | no speaker; self-told; no event |
| `read` | `descriptor.source` **and** a reading event | no source; no event |
| `inferred` | `descriptor.basis` = ids of **the same character's** existing accounts | empty basis; any id not owned by that character; any id that does not exist |
| `verified` | `descriptor.verification` event **and** its basis | no verification event; no basis |
| `narrator_grant` | `descriptor.declaration` resolving to a live declaration in this chat | unknown or withdrawn declaration |
| `none` (`own_account`) | **nothing** — the absence is the record's meaning | never |

**Structural invariant.** The validator may read only: the record being written, records already owned
by that character, the declaration registry, and the event reference. It may never call
`authoritative()`. Evidence establishes *acquisition*, never *correctness*.

---

## 6. Contract C-4 (revised) — ambient authorship (U-5, U-8)

### 6.1 Reuse test (U-5)

The owner asked for the test, not the assumption. Applied:

| Question | Finding | Verdict |
|---|---|---|
| Can `claim` represent an attributed authored contribution? | Rule 22 opens creation to "any participant, human or model"; rule 23 requires the system to record who contributed what, in which domain, by what route; §3 describes the class as a statement and its author | **Yes** |
| Can `declaration` carry it? | O-12's declared kinds include `declaration`; the *operation* declares it at the call boundary, so nothing is inferred from content | **Yes** |
| Does anything depend on the contributor being a *character*? | Rule 3 governs private records; a domain-scoped authored claim is not private and is not owned by a character | **No conflict** |
| Is ambient material **assertoric**? | It states something about the world ("the capital is Vareth"), which is what a claim is | **Yes** |

**Adopted:** `claim` + `declaration` + provenance block `{"authorship": {actor, domain, grant_id, turn}}`.

**The failure condition, named:** reuse holds **because** ambient material is assertoric. If authored
material is later required to be a **definition** — a category statement rather than an assertion — the
claim class cannot represent it without collapsing the distinction B3 is waiting on. That is exactly
the boundary at which reuse would stop being valid, and it belongs to U-6/B3, not to this contract.

### 6.2 Authority (U-8)

| Rule | Consequence |
|---|---|
| Authoring authority is not promotion authority | separate capability on the grant |
| Standing promotion authority does not grant authorship | a model grant for promotion confers no authoring |
| Authorship does not grant promotion | authoring creates a non-truth record; promotion stays its own operation |
| All contributions attributable | author, domain and grant recorded on every authored record |

**Reconciliation with O-8a-5.** 8a-5 says one **shared authoring-authority regime** applies to both
operations, "the operations remain distinct". The reading adopted: the regime is shared in **structure**
— one holder source, one per-domain scoping, one grant record form, one revocation path — while
**capability is per operation**. That is what makes U-8's separation and 8a-5's sharing compatible.
Residual **R-A** records the alternative reading, which would contradict U-8.

### 6.3 Operation

`author_ambient(actor, domain, target, text, turn, rationale=None, revises=None)`

- requires `authorship` capability in that domain;
- creates a `claim` with kind `declaration`, contributor the actor;
- never creates truth; never appears in a world answer;
- `revises` appends a new record and marks the prior `superseded`, preserving text and history;
- is a lawful `record_cited` basis for a later promotion, which remains the only path to authority.

**Refusals:** no authorship capability · unknown domain · empty text · `revises` naming a record from
another chat or a truth record · attempting to author in a domain where the actor holds only promotion
capability.

---

## 7. Authority model revision (U-8) — one added field

| Element | Change |
|---|---|
| `Grant` | gains `operations: frozenset[str]` over `{"promotion", "authorship"}` |
| Domain declaration | installs the human operator with **both** capabilities |
| Model grant (8a-6) | defaults to `{"promotion"}`; authorship requires an explicit grant naming it |
| `require(domain, actor)` | becomes `require(domain, actor, operation)` |
| Revocation | per grant; revoking promotion capability leaves authorship capability untouched, and vice versa |

This is **one field on an existing object** plus a parameter on an existing check.

---

## 8. Contradiction analysis — the demonstration the owner asked for

Nine mechanisms against the ratified set. Verdicts: **C** compatible · **P** compatible via a recorded
precedence application · **R** residual risk, named.

| # | Mechanism | Rules examined | Requirements examined | Verdict |
|---|---|---|---|---|
| M1 | Explicit proposition identity | 6 (applicability unchanged), 41 (candidacy key), 14 (promotion governs) | K2, K9 (additive fields), J13 (overlap ≠ conflict), A4 | **C** — applicability keeps its rule-6 meaning; the proposition is a separate dimension, so rule 6 is not reinterpreted |
| M2 | No implicit supersession | 20, 21, 42, 43 | K7, K10, K22, K23, K26 | **P** — RAT-1 below |
| M3 | Proposition `kind`, belief exclusion | 7 (no new truth category), 15, 17 | D10, H3, A4, K2 | **C** — `kind` describes the *subject* of a truth record, not a new truth class; exactly two values, closed |
| M4 | Evidence validation | 9, 11, 50, 51, 52 | C5, C7, C8, F6 | **C** — refusals are structural; no truth read anywhere in the path |
| M5 | Ambient via `claim`/`declaration` | 22, 23, 26, 27, 31, 45 | K25, K27, L6, L8 | **R-B** — §3's claim-class description is stretched for a non-character author; semantics hold while material is assertoric |
| M6 | Assertion + conflict group | 8, 13, 20, 21, 34, 37, 41 | J12, J13, J14, J15, J16 | **C** — the assertion is recorded, never assessed; no record touched; closure human-gated |
| M7 | Split authority capability | 24, 35, 39, 40 | K21, L7, L10 | **P** — R-A below |
| M8 | Append-only belief revocation | 53, 55 | D5, D9, O-5 | **C** — new record, prior retained; strictly closer to rule 55 and D5 than the code it replaces |
| M9 | Canonical `via` names | 50, 51 | C5, C9 | **C** — tier-2 order untouched; no enum value added |

### Pre-existing tensions exposed (not created)

| ID | Tension | Status |
|---|---|---|
| **RAT-1** | **K10** says "Latest governs; predecessors retained as history". **O-10b** (tier 1) says no actor may displace another's or their own prior promotion without a human edit. For the same-actor case these cannot both hold literally | **Pre-existing.** Not caused by any mechanism here. Resolved by precedence (tier 1 wins), read as the supersession chain. K10's text is a ratified requirement and **cannot be edited** — a qualification record would need an owner instruction |
| **RAT-2** | **K7** "authority no longer ambiguous" versus **O-10**'s equal-authority state, which is by definition unresolved | **Pre-existing.** Read as: a promotion's own proposition is unambiguous unless a competing live promotion or an assertion exists |

**No mechanism above creates a new contradiction with any of the 57 rules or 118 requirements.** Two
*existing* tensions between ratified items are surfaced, given a precedence reading, and left for the
owner to qualify or leave standing.

### Residual risks

| ID | Risk | Consequence |
|---|---|---|
| **R-A** | If O-8a-5 is read as *shared capability* rather than shared structure, U-8's separation contradicts it | The owner's recommendation (U-8) and 8a-5 need one reading chosen; this file adopts shared structure |
| **R-B** | Ambient material that is a definition rather than an assertion breaks the `claim` reuse | Bounded by B3; not repaired here, by instruction |
| **R-C** | The claim-path gap in C-2 (§4.2) — a belief voiced in dialogue, promoted as fact from a claim | Undetectable without text inspection; stated, not closed |
| **R-D** | F7 — no replay facility; determinism remains unproven | Unchanged from R2 |
| **R-E** | Proposition minting is an actor's declaration, so two actors may mint two identities for "the same" claim | Accepted: this is a false negative, never a false conflict — the safer direction under U-1 |

---

## 9. Conformance impact (recomputed) — all 57 rules

Codes: **U** unchanged · **M** mechanism modified/added · **N** newly satisfiable · **R** owner reading
required · **X** still not implementable.

| Rule | | Rule | | Rule | |
|---|---|---|---|---|---|
| 1 | U | 2 | U | 3 | U |
| 4 | U | 5 | U | 6 | U |
| 7 | U | 8 | **M** | 9 | U |
| 10 | U | 11 | U | 12 | U |
| 13 | **N** | 14 | **M** | 15 | U |
| 16 | U | 17 | U | 18 | U |
| 19 | U | 20 | **M** | 21 | **M** |
| 22 | U | 23 | **M** | 24 | **M** |
| 25 | U | 26 | U | 27 | **N** |
| 28 | U | 29 | U | 30 | U |
| 31 | U | 32 | U | 33 | U |
| 34 | **M** | 35 | U | 36 | U |
| 37 | **M** | 38 | U | 39 | **N** |
| 40 | U | 41 | **N** | 42 | **M** |
| 43 | **M** | 44 | U | 45 | U |
| 46 | U | 47 | U | 48 | U |
| 49 | U | 50 | **M** | 51 | U |
| 52 | U | 53 | **M** | 54 | U |
| 55 | **M** | 56 | **N** | 57 | U |

**Changes from REV 1:** rule 6 moves **R → U** (applicability is no longer reinterpreted), rule 24 moves
**U → M** (grants carry capabilities). Non-U rules and notes:

| Rule | Note |
|---|---|
| 8 | Incompatibility is an assertion, never a verdict |
| 13 | Detection becomes representable |
| 14 | Supersession governs; a promotion governs only its own proposition |
| 20 | Groups and supersession chains preserved; nothing erased |
| 21 | Group membership is authority-state history, not precedence |
| 23 | The authorship route becomes distinguishable from promotion |
| 24 | A grant now names its operations |
| 27 | "Creation is not promotion" becomes demonstrable |
| 34 | The authority judges; the engine validates structure only |
| 37 | Group closure is revocation/supersession |
| 39 | The second operation exists and is defined |
| 41 | Reachable without false positives and without a pick |
| 42 | Correction/supersession/revocation kept distinct; the `promote(supersede=True)` path is removed |
| 43 | Human gate is the only route to displacement |
| 50 | Evidence table enforces the two-field distinction |
| 53 | Belief lifecycle gated on the character's own model |
| 55 | Belief revocation is append-only |

**Tally: 39 U · 13 M · 5 N · 0 R · 0 X = 57.**

---

## 10. Conformance impact (recomputed) — all 118 requirements

| Req | | Req | | Req | | Req | |
|---|---|---|---|---|---|---|---|
| A1 | U | A2a | U | A2b | U | A3 | U |
| A4 | **M** | B1 | U | B2 | U | B3 | **X** |
| B4 | U | B5 | U | B6 | U | B7a | U |
| B7b | U | B7c | U | B8 | U | B9 | U |
| C1 | U | C2 | **M** | C3 | U | C4 | **M** |
| C5 | **N** | C6 | U | C7 | **M** | C8 | **M** |
| C9 | U | D1 | U | D2 | U | D3 | U |
| D4 | U | D5 | U | D6 | U | D7 | U |
| D8 | U | D9 | **M** | D10 | **N** | E1a | U |
| E1b | U | E1c | U | E2 | U | E3 | U |
| F1 | **M** | F1b | **M** | F2 | U | F3 | **M** |
| F4 | U | F5 | U | F6 | **N** | F7 | **R** |
| F8 | **M** | G1 | U | G2 | U | G3 | U |
| H1 | U | H2 | U | H3 | U | H4 | U |
| H5 | U | H6 | U | H7 | U | H8 | U |
| H9 | U | H10 | U | I1 | U | I2 | U |
| I3 | U | J1 | U | J2 | U | J3 | U |
| J4 | U | J5 | U | J6 | U | J7 | U |
| J8 | U | J9 | U | J10 | U | J11 | U |
| J12 | **N** | J13 | **N** | J14 | **N** | J15 | **M** |
| J16 | **M** | K1 | U | K2 | **M** | K3 | U |
| K4 | U | K5 | U | K6 | U | K7 | **M** |
| K8 | U | K9 | **M** | K10 | **M** | K11 | **M** |
| K12 | U | K13 | U | K14 | U | K15 | U |
| K16 | U | K17 | U | K18 | U | K19 | U |
| K20 | U | K21 | U | K22 | **M** | K23 | **M** |
| K24 | U | K25 | **N** | K26 | **M** | K27 | U |
| L1 | U | L2 | U | L3 | U | L4 | U |
| L5 | U | L6 | **N** | L7 | **M** | L8 | U |
| L9 | U | L10 | **M** | — | — | — | — |

**Changes from REV 1:** K7 moves **R → M** (U-2 settles the reading), F7 stays **R** (no replay exists),
L7 and L10 move **U → M** (grants now name capabilities).

| Requirement | Note |
|---|---|
| A4 | Belief-fact exclusion sharpens the world answer |
| C2 | Inference basis validated against the character's own accounts |
| C4 | Read requires a source **and** an event |
| C5 | Wrong evidence refused at write — newly satisfiable |
| C7 | Own-account exemption checked before the table |
| C8 | Validator forbidden from parsing text |
| D9 | Belief lifecycle gated on the character's own model |
| D10 | Enforced structurally via proposition `kind` |
| F1, F1b, F8 | Proposition registry, groups and authorship enter the snapshot |
| F3 | Revision appends; text never rewritten |
| F6 | Evidence resolves per origin by construction |
| F7 | Canonical ordering fixed; **no replay facility** — determinism unproven |
| J12–J14 | Flag-and-abstain, group representation, unresolved retrieval |
| J15, J16 | Correction cannot close; human closure recorded, members retained |
| K2, K9 | Proposition identity recorded with domain, applicability, provenance |
| K7 | Settled per proposition; competing live records are unresolved (RAT-2) |
| K10 | Chain is created by supersession (RAT-1) |
| K11 | Group membership joins preserved history |
| K22, K23 | Revocation and supersession close groups explicitly |
| K25 | Ambient authorship exists and does not promote |
| K26 | The `promote(supersede=True)` bypass is removed |
| L6 | Creation vs authority demonstrable |
| L7, L10 | Grant capability and per-capability revocation |

**Tally: 86 U · 22 M · 8 N · 1 R · 1 X = 118.**

---

## 11. Exact proposed changes to the implementation (still unapplied)

`engine.py` 392 lines · `records.py` 154 · `authority.py` 156 · `retrieval.py` 85.

### 11.1 Corrections under existing authority — owner-recommended, not applied

| # | File | Symbol | Change |
|---|---|---|---|
| 1 | `engine.py` | `promote` | Delete the `supersede` parameter and its displacement block (D-A) |
| 2 | `engine.py` | `retract` | Refuse belief records; route to belief operations (D-B) |
| 3 | `engine.py` | new `revoke_belief` | Model-gated, append-only revocation record (D5, O-5) |
| 4 | `engine.py` | `promote` | Refuse a belief citation unless the proposition is a `belief_ascription` of that owner |
| 5 | `engine.py` | `acquire` | Evidence table; own-account exempt; inference basis owned by the character; no truth reads |
| 6 | `retrieval.py` | `world` | Emit `unresolved_relation` for ≥2 live records; never `settled` |

### 11.2 Requires explicit owner approval [PA]

| # | File | Symbol | Change | Owner basis |
|---|---|---|---|---|
| 7 | `engine.py` | new `_propositions` registry | Proposition identity, `kind`, `about_character` | U-1, U-4 |
| 8 | `records.py` | `Record` | Add `proposition_id`; validate presence on truth classes | U-1, U-4 |
| 9 | `engine.py` | new `assert_incompatible` | Domain-authority assertion with applicability | U-7 |
| 10 | `engine.py` | `_promotions`, `snapshot` | Persist conflict groups and closure | U-2 |
| 11 | `retrieval.py` | new `belief_facts` | Belief-fact view; world answers exclude ascriptions | U-4 |
| 12 | `authority.py` | `Grant`, `require` | `operations` capability set; per-capability revocation | U-8 |
| 13 | `engine.py` | new `author_ambient` | Authoring capability; `claim` + `declaration`; non-promoting revision | U-5, U-8 |
| 14 | `__init__.py` | `__all__` | Export `assert_incompatible`, `author_ambient`, `revoke_belief` | 3, 9, 13 |

**Untouched:** `eval/mem1/model.py` and the G1 harness; production Howl; r32; plan v3; Amendment 1; v9;
the map; the ledger; REV 1.

---

## 12. What this record does not claim

- **No implementation.** No source file changed. The two defects are **specified, not corrected**.
- **No execution.** No test, probe or suite ran; `py_compile` was not needed.
- **No behavioural verification.** Current-behaviour statements are source inspection ([SO]).
- **No ratification.** The eight items are the owner's **recommendations**, recorded verbatim. If they
  are to become registered decisions, that is a separate instruction.
- **No resolution of B3.** U-6 is honoured: nothing was added to make B3 look solved.
- **No amendment** to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9, the map,
  the ledger or REV 1 (`98447423…`).
- **No transfer, sharing or carryover**, and no destination interface.

**Two items need one line from the owner before implementation authorisation:** the O-8a-5 reading in
R-A, and whether RAT-1/RAT-2 should be recorded as qualifications against K10/K7 (the map cannot be
edited without instruction).
