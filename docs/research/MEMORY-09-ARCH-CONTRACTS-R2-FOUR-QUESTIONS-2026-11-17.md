# MEMORY-09 — Architectural contracts for the four R2 questions

**Status: PROPOSAL FOR OWNER REVIEW. Not ratified. Not implemented. Not tested.**
Recorded **2026-11-17 09:17:41 CST (−0600)** = **2026-11-17 15:17:41 UTC**.

| Field | Value |
|---|---|
| Authority | None. This file proposes contracts; it decides nothing and amends nothing |
| Ratified inputs | v9 (57 rules, tier 3) · Conformance map (118 requirements) · r32 (tier 2, `e88ef2d1…`) · O-1…O-13, D-01, D-02, O-8a-1…6 |
| Owner-approved this round | **Four directions**, quoted in §11.2. **Not** their technical mechanisms |
| Implementation | **Unmodified.** No file under `eval/mem1/` was changed for this record |
| Tests / probes run | **0.** No suite, probe or self-test executed; no code changed, so none was needed |
| Scope | Contract definition and reconciliation only |

Claim labels: **[RC]** existing ratified requirement/rule · **[OD]** owner-approved direction (this round) ·
**[PA]** proposed architecture (new, needs approval) · **[SO]** static observation from source inspection.

---

## 0. What was approved — and what this file proposes

The four directions were approved **as directions**. This table separates them from the mechanisms that
would implement them, because a direction is not a mechanism.

| Layer | Status | Examples in this file |
|---|---|---|
| **Ratified** [RC] | Existing authority, unedited | 57 rules, 118 requirements, O-1…O-13, D-01, D-02, O-8a |
| **Owner-approved direction** [OD] | Approved this round | "do not infer incompatibility from target equality alone"; "beliefs must never silently become world truth"; the five acquisition kinds; "ambient authorship as a distinct authoring operation" |
| **Proposed architecture** [PA] | **Needs explicit approval** | conflict-group record · `assert_incompatible` · `unresolved_relation` status · required `about` parameter · belief target namespace · descriptor key schema · `author_ambient` |
| **Proposed requirements** | **Not added** — the set stays at 118 | NR-1…NR-4 in §11.3, listed for decision, not insertion |

**No mechanism here is authorised by the direction that motivated it.** Where a proposal needs a new
structural element it is labelled [PA], and listed in §8 and §11.3.

---

## 1. Contract C-1 — Incompatibility and conflict detection

### 1.1 Intended observable behaviour

| # | Situation | Required observation |
|---|---|---|
| 1 | Two promotions on **different** propositions about one subject | No candidacy, no conflict; each proposition answers for itself |
| 2 | Two promotions on the **same** declared proposition, no incompatibility established | Both retained, neither presented as settled, **not** labelled a conflict |
| 3 | An authority **establishes** incompatibility between two live promotions | Explicit conflict group; both preserved with provenance and promotion history; neither preferred |
| 4 | No information establishing incompatibility exists | No conflict is manufactured |
| 5 | A conflict group is open | The truth query returns the conflict, not a pick |
| 6 | Correction of one member | Group stays open |
| 7 | Human supersession or revocation | Group closes; both records and the disagreement survive |

### 1.2 Minimum mechanism

**Proposition identity is declared, never inferred.** `Record.target` already keys applicability
(rule 6). The only semantic move is to state what that key *is*: a **declared proposition identifier**,
not a subject name. Two statements about one subject are promoted under two different targets — by the
actors who promote them — and never share a candidacy slot unless the actors declare the same
proposition.

That gives a mechanical candidacy rule with no content comparison and no semantic judge:

```
candidacy(domain, target) = promotions P where P.domain == domain, P.target == target,
                            P.status == ACTIVE
```

| Live candidacy | State | Authority |
|---|---|---|
| 0 records | `none` | — |
| 1 record | `settled` | [RC] K7 |
| ≥2 records, no assertion, no group | `unresolved_relation` | [PA] — abstention, not a conflict label |
| live records + open group | `conflict` | [RC] rule 41, J13, J14 |

`unresolved_relation` is the honest state: the engine has established nothing about incompatibility (it
cannot), so it does not call the records a conflict — and equally will not hand back one as settled. It
returns every live record and says the relation is undeclared.

**Incompatibility is established only by an explicit assertion** from an actor holding authority in that
domain, recorded verbatim with actor, authority, turn and statement. The engine never derives it from
text, target equality, recency, repetition or acquisition.

**The conflicted group is persisted**, because J13 requires it to be *explicitly represented*, J15
requires correction not to dissolve it, and J16 requires a human edit to close it. Shape [PA]:

| Field | Meaning |
|---|---|
| `group_id`, `domain`, `target` | identity and scope |
| `members` | record ids and their `promotion_id`s (J13 provenance) |
| `asserted_by`, `authority`, `turn`, `statement` | who established incompatibility, under what authority, in what words |
| `status` | `open` / `resolved` |
| `resolved_by`, `resolution_turn`, `resolution_kind` | `supersession` or `revocation` only |

### 1.3 Ratified elements satisfied

| Element | How |
|---|---|
| **Rule 41** [RC] | Candidacy retained; incompatibility becomes an established state rather than a content guess; both records and their promotion history persist; retrieval cannot present either as settled |
| **Rule 34** [RC] | The authority judges; the engine validates structure only. The assertion is recorded, never assessed |
| **Rule 21** [RC] | The group changes presentation, not the record of disagreement |
| **Rule 13** [RC] | Detection becomes representable; resolution stays separate |
| **Rule 42** [RC] | Correction still cannot close a group; closure is supersession or revocation |
| **Rules 43–44, O-10b** [RC] | Closure is human-gated; original records and the fact of disagreement survive |
| **J12–J16** [RC] | J12 flag without verdict · J13 both retained with provenance · J14 unresolved in retrieval · J15 correction leaves it open · J16 human edit resolves, history survives |
| **O-10** [RC] | The equal-authority state: both authoritative, both preserved, neither wins |

### 1.4 Possible conflicts with existing requirements

| Element | Risk | Reading proposed |
|---|---|---|
| **K7** — "promotion settles authority" | With two live promotions and no assertion, the answer is not settled | K7 describes the ordinary case: a proposition is governed by the promotion that created it. It does not say every later promotion displaces predecessors — that is K10. **Needs owner confirmation (U-2)** |
| **K10** — "latest governs" | Under C-1 a bare promotion does not displace the earlier one | Tier 1 (O-10b) outranks tier 3: an actor may not displace even *their own* prior promotion without a human edit. K10's chain is therefore read as the **supersession** chain. **Needs owner confirmation (U-2)** |
| **D-2** — materiality predicate | Conflating representation with materiality would pre-empt D-2 | Explicitly separated: C-1 defines *representation*. Whether a conflict is **material to a character's current turn** is D-2's question, remains open, and is neither decided nor anticipated here |
| **Rule 6** [RC] | Re-reading `target` as a proposition key narrows "applicability" | Rule 6 is untouched in form; the granularity of the key is fixed. **Needs owner confirmation (U-1)** |

### 1.5 New structure required

| Element | Type | Why it cannot be avoided |
|---|---|---|
| `assert_incompatible(actor, record_a, record_b, statement, turn)` | **Operation** [PA] | Rule 41's state is unreachable if incompatibility can never be established, and the system may not establish it |
| Conflict group | **Persisted object** [PA] | J13 "explicitly represented", J15 "must not resolve", J16 "history survives" — a transient return value cannot carry these |
| `unresolved_relation` | **New retrieval status** [PA] | Otherwise the choice is "conflict" (manufactured) or "settled" (silent winner) |
| Target granularity | **Contract reading, no field** [PA] | Distinguishes same-proposition from same-subject without adding data |

**No new record class. No new field on `Record`.** If the owner prefers data-level separation, a
`proposition` field is the alternative (U-1): it costs a field and buys protection against two actors
choosing different targets for the same proposition.

### 1.6 Provenance, ownership, authority, history

- **Provenance** — group members referenced by record id *and* `promotion_id`; the assertion carries its
  own actor, authority, turn and verbatim statement.
- **Ownership** — untouched: groups reference authoritative records, which are never character-owned
  (rule 17).
- **Authority** — asserting requires a live grant in that domain (8a-1/8a-3) and confers no power to
  resolve; resolution stays human (O-10b).
- **History** — no record edited or deleted; group creation and closure appear in the snapshot (F1).

### 1.7 Execution paths

**No false conflict.** `promote(target="king.status", text="the king is dead")` then
`promote(target="king.death.cause", text="the king died of poison")` → two propositions; each answers
for itself; neither `unresolved`, neither `conflict`.

**Unresolved relation.** A second `promote` on `"king.status"` with `text="the king is alive"`, no
assertion → both records returned, `status="unresolved_relation"`. Not `settled`, and **not** `conflict`.

**Established conflict.** `assert_incompatible(actor="operator", record_a="truth-00012",
record_b="truth-00019", statement="these describe the same event differently", turn=9)` → group
`cg-0001` `open`; `authoritative("king.status")` → `status="conflict"`, `groups=["cg-0001"]`, both
records.

**Refusals.** Assertion refused when: the actor holds no domain authority; `record_a == record_b`; a
record is not truth-apt (rule 7); a record is outside this chat (rule 31); the pair already shares an
open group; a member is not `ACTIVE` (no live conflict).

**Correction.** `correct(record_id="truth-00012", …)` with `cg-0001` open → text changes, group stays
`open` (J15). `supersede_authority`/`revoke_promotion` by a human holder → group `resolved`,
`resolution_kind` recorded, members retained (J16, rule 44).

### 1.8 Implementable vs requiring approval

| Part | Status |
|---|---|
| Replace "multiple active records ⇒ conflict" with the four-state answer | **Implementable under existing authority** — it removes an inference the direction forbids |
| Conflict group, assertion operation, new status value | **Requires explicit owner approval** [PA] |
| Target-as-proposition reading | **Requires owner confirmation** (U-1) |

**Defect found while specifying ([SO], no test run):** `promote(supersede=True)` changes authority state
while bypassing `_require_human`. An actor holding a standing grant could supersede another actor's
record — or its own — without the human edit O-10b requires. Conformance gap against rules 43 / K26,
independent of C-1; its fix implements existing rules and needs no new architecture.

---

## 2. Contract C-2 — Structural separation of belief and world truth

### 2.1 Intended observable behaviour

| # | Situation | Required observation |
|---|---|---|
| 1 | Belief declared by the character's own model | Belief record, owned by that character, append-only revision |
| 2 | Promotion citing that belief | Authority over **a fact about the character**, never over the belief's content |
| 3 | Truth query for the belief's proposition | Belief-derived records ineligible; answer comes from world/scenario authority or abstention |
| 4 | Disclosure of a belief | Recipient acquires *that the character believes it* (rule 57), not its content |
| 5 | Revision or revocation of a belief | Only by that character's own model |
| 6 | Prior beliefs after revision | Still retrievable by their owner |

### 2.2 Minimum mechanism

Two rules and one convention:

1. **Every promotion declares `about`** ∈ {`world`, `belief:<character>`} [PA], recorded in the
   promotion metadata and never derived. There is no default, because a default is an assumption about
   what the actor meant.
2. **A promotion whose basis cites a belief must declare `about="belief:<owner>"`** and target the
   belief namespace (`belief:<character>:<proposition>`), never the belief's own proposition. The
   engine refuses the alternative — this is rule 56's guard, able to fire.
3. **Truth queries filter by construction**: a belief-fact promotion carries a namespaced target, so it
   is never returned for the world proposition. The separation is spatial, not semantic.

Together these mean a belief cannot reach world truth *through the belief path*. They do not cover an
actor promoting the belief's content as world truth citing a **claim** in which a character voiced it:
the engine cannot determine that the claim's content *is* `P`, and making it try would be the semantic
adjudication rule 8 forbids. That residual is stated rather than papered over — the actor is then
asserting `P` on their own authority, which is a different act from deriving it from a belief.

### 2.3 Ratified elements satisfied

| Element | How |
|---|---|
| **Rule 56 / D10** [RC] | Authority over "Sable believes P" is a fact about Sable, structurally marked and spatially excluded from P |
| **Rules 53–55** [RC] | Declarer gate retained; narration cannot install; revision stays append-only |
| **O-5** [RC] | Declaration-only creation, own-model ownership, no promotion of content, append-only history |
| **O-4, O-11** [RC] | Belief remains its own record class; promotion changes authority, not ownership |
| **Rules 3, 5, 7** [RC] | Belief owned by one character; never truth-apt; never lifted to a truth class |
| **Rule 57** [RC] | Disclosure already carries the belief, not the content |

### 2.4 Possible conflicts

| Element | Risk | Reading proposed |
|---|---|---|
| **K2** [RC] | A belief-fact's applicability is a character, not a world proposition | Applicability is still declared and recorded; the domain stays world/scenario; the *subject* is a character. **Needs confirmation (U-4)** |
| **K9** [RC] recorded fields | `about` is an additional recorded parameter | K9 is a minimum, not a closed list. Recorded as [PA] regardless |
| **Rule 6** [RC] | The namespace stretches "applicability" | Same reading as C-1; same confirmation (U-1) |
| **F2** [RC] | A belief-fact *is* a truth record about a belief | F2 constrains belief records, not authoritative records about beliefs. Unchanged |
| **D-2** | none | Not touched |

### 2.5 New structure required

| Element | Type | Why |
|---|---|---|
| `about` (required promotion parameter) | **New recorded parameter** [PA] | Without it a belief-derived promotion is indistinguishable from a world promotion at the point of writing |
| Belief namespace for targets | **Convention** [PA] | Makes exclusion spatial and cheap |
| Belief `target` populated | **Field reuse** — `Record.target` exists and is empty on beliefs | Lets a belief-fact state which proposition it concerns without adding a field |

**No new record class. No new field.**

### 2.6 Provenance, ownership, authority, history

Belief records keep `owner_character`, `attribution` (declarer + role), and their append-only chain
(`supersedes`/`superseded_by`). A belief-fact promotion records `about="belief:<character>"`, its basis
(including the cited belief), actor, authority and turn. Ownership is never transferred: the truth
record is domain-owned, the belief stays the character's. Prior beliefs remain retrievable after
revision (status `superseded`, never deleted).

### 2.7 Execution paths

**Declared.** `declare_belief(character="sable", declarer="sable", role="character_model",
text="the king was poisoned", turn=4)` → `belief-00003` owned by `sable`. A claim by another character,
or narration, cannot reach this method (rules 53/54).

**Correct promotion.** `promote(actor="operator", domain=("world","w1"),
target="belief:sable:king_poisoned", text="Sable believes the king was poisoned",
basis=[BasisEntry("record_cited", record_id="belief-00003")], about="belief:sable", turn=7)` →
authoritative for the belief-fact only.

**Refusal (rule 56).** Same basis with `target="king.poisoned"`, `about="world"` → refused: *a promotion
citing a belief cannot govern the belief's content*.

**Query separation.** `authoritative(("world","w1"), "king.poisoned")` → absent; the belief-fact is
ineligible. `authoritative(("world","w1"), "belief:sable:king_poisoned")` → settled.

**Revision.** A second `declare_belief` → new record, prior `superseded`; `perspective("sable", …)` still
returns the prior belief to its owner.

**Revocation refusal.** `retract("belief-00003")` by anyone other than Sable's model → **refused**
(defect fixed: current `retract` accepts any non-truth record with no ownership or role check).

### 2.8 Implementable vs requiring approval

| Part | Status |
|---|---|
| Refuse a promotion that cites a belief and targets its content | **Implementable now** — enforces rule 56 with existing data |
| Gate belief revocation on declarer identity and role | **Implementable now** — defect fix against O-5/rule 53 |
| Required `about` parameter + belief namespace | **Requires owner approval** [PA] |
| Populating `target` on belief records | **Requires owner approval** as a data-shape change (field exists; its use on beliefs is new) |

---

## 3. Contract C-3 — Acquisition evidence and write-time validation

### 3.1 Intended observable behaviour

Every acquired account is accepted only if the evidence required by its own acquisition kind is present
in the record. Every refusal happens at write time. No validator ever consults domain truth.

### 3.2 Minimum mechanism

One validation table consulted by `acquire`, plus an exemption that is itself checked:

| Kind | Evidence required at write | Refused when |
|---|---|---|
| `experienced` see §3.4 | an acquisition event reference (`span`) attributable to that character, plus content | no span; character absent from the event |
| `told` | a disclosure event reference and `descriptor.speaker` ≠ owner | no speaker; self-told; no event |
| `read` | `descriptor.source` (document/asset identity) **and** a reading event | no source; no event |
| `inferred` | `descriptor.basis` = ids of **the same character's** existing accounts | empty basis; any id not owned by that character; any id that does not exist |
| `verified` | `descriptor.verification` event **and** its basis | no verification event; no basis |
| `narrator_grant` | `descriptor.declaration` resolving to a declaration in this chat | unknown or withdrawn declaration |
| `own_account` (`via="none"`) | **nothing** — the absence is the record's meaning (rule 51) | never, by this validator |

**The invariant, enforced structurally.** The validator may read only: the record being written, records
already owned by that character, the declaration registry, and the event reference. It must never call
`authoritative()`. That restriction — not a comment — is what keeps provenance validation from becoming
truth adjudication: evidence establishes *acquisition*, never *correctness*.

`inferred` is where this bites hardest: refusing a basis the character does not already hold prevents
inference from manufacturing access, while accepting any basis the character does hold keeps the system
out of judging whether the inference was good.

### 3.3 Ratified elements satisfied

| Element | How |
|---|---|
| **C5** [RC] | Wrong or missing evidence refused at write; previously absent |
| **Rule 50** [RC] | `via` and `origin_class` keep their separate questions |
| **Rule 9** [RC] | Knowledge exists only as attributable accounts with evidence by origin |
| **Rule 11** [RC] | Exposure governs acquisition; `told` requires the disclosure event |
| **Rule 51** [RC] | `own_account`/`via=none` exempt and never treated as missing evidence |
| **Rule 52** [RC] | A source named inside a recollection stays content; the validator must not read text |
| **O-1, O-4, O-12, D-02** [RC] | Provenance without bypass; two fields; declared kinds; existence is not knowledge |

### 3.4 Name reconciliation — unresolved, not assumed

The approved direction lists `experienced`, `told`, `read`, `inferred`, `verified`. The ratified `via`
vocabulary is `witnessed > verified > read > told > inferred > narrator_grant` (r32 K-R26, tier 2).

| Difference | Consequence |
|---|---|
| `experienced` vs `witnessed` | Adding `experienced` as a **new value** inserts a term into a fixed total order in tier 2, which nobody is authorised to edit. Treating it as a **plain-language synonym** changes nothing. **Owner must choose (U-3)** |
| `narrator_grant` absent from the five | It exists in the ratified order and is produced by declaration coverage (B1). Either a sixth acquisition kind or out of scope. **Owner must state which (U-3)** |

This file validates six kinds plus the exempt case, marking `experienced` as the label for `witnessed`
**pending U-3**. It renames nothing.

### 3.5 Possible conflicts

| Element | Risk | Reading proposed |
|---|---|---|
| **C7** [RC] | A validator demanding evidence would break own accounts | Explicit exemption checked before the table |
| **C8** [RC] | Requiring verification for a recollection's stated source would manufacture an event | The validator never parses text; a named source is not evidence |
| **F6** [RC] | Every record's evidence must resolve per origin | This is the mechanism that makes F6 satisfiable |
| **B1/B7** [RC] | Declarations create coverage accounts | `narrator_grant` validation resolves the declaration |
| **K-R1** (tier 2) | Edge creation depends on explicit acquisition events | Consistent: three kinds now require the event reference instead of assuming it |

### 3.6 New structure required

**No new field, no new class, no new operation.** Two shapes change:

1. `descriptor` — an unconstrained dict — acquires a **documented key contract** per via (`speaker`,
   `source`, `basis`, `verification`, `declaration`, `event`). A schema decision, not a field addition,
   and [PA].
2. `span` stops being optional for three kinds.

### 3.7 Execution paths

`acquire(character="sable", via="told", text="the king was poisoned", descriptor={"speaker":"alaric"},
turn=6)` → accepted.

`acquire(character="sable", via="told", text="…", descriptor={}, …)` → **refused**: a told account with
no speaker is not evidence.

`acquire(character="sable", via="inferred", text="the poisoner knew the court",
descriptor={"basis":["acct-00021"]}, …)` where `acct-00021` belongs to `alaric` → **refused**: inference
cannot grant access the character did not already hold.

`acquire(character="sable", via="read", text="…", descriptor={"source":"chronicle-3"}, span=None)` →
**refused**: reading needs the event, not just the book.

`acquire(character="sable", origin_class="own_account", text="I remember the king falling", …)` →
accepted with `via="none"`, no event, no missing-evidence error.

**Truth isolation.** None of these paths reads `authoritative()`; a refusal is structural, and an
accepted account may be false — being told something establishes communication, not correctness.

### 3.8 Implementable vs requiring approval

| Part | Status |
|---|---|
| Validation table, typed refusals, own-account exemption, inference-basis ownership check | **Implementable now** — C5 is an existing ratified requirement; this is its mechanism |
| Descriptor key contract | **Requires owner approval** as a schema decision [PA] |
| `experienced` naming and `narrator_grant` status | **Requires owner decision (U-3)** — cannot be resolved by implementation |

---

## 4. Contract C-4 — Q7 ambient authorship

### 4.1 Intended observable behaviour

| # | Situation | Required observation |
|---|---|---|
| 1 | An actor with domain authority authors world/scenario material | Record exists, attributed to that actor and authority, in that domain |
| 2 | Truth query for that material | Nothing authoritative; not truth-apt |
| 3 | The material promoted later | Promotion is a separate operation with its own basis and authority |
| 4 | Editing authored content | A new version; the prior retained; nothing promoted |
| 5 | Audit | Authorship and promotion separately attributable, different identities and turns |
| 6 | Authoring across domains, or without authority | Refused |

### 4.2 Minimum operational meaning of Q7

> **Q7 ambient authorship** is the operation by which an actor holding the domain's authoring authority
> contributes world or scenario *material* — non-truth-apt, attributed, domain-scoped — which becomes
> authoritative only through a separate promotion.

That is the smallest meaning satisfying O-8a-5 ("one shared authoring-authority regime, two distinct
operations") without inventing a broader authority regime: the same per-domain grant permits promotion
and authorship, and authorship's only product is a non-authoritative record.

**Record class — deliberately not a new one.** v9 §2 states the withdrawn four-class architecture
(including `ambient_authored`) is **not revived**; a new class would also amend §3's table. So ambient
material is recorded as a **claim** with the existing declared kind `declaration`, contributor the
authoring actor, and a provenance block `{"authorship": {actor, domain, authority, turn}}` [PA]. This
stretches §3's description of the claim class ("that a character said something"): the contributor here
is an authoring authority, not a character speaker. **Needs owner confirmation (U-5).**

**Authority.** `author_ambient` requires a live grant in that domain — the same requirement as
`promote`. No basis is required: 8a-2 governs *promotion* bases, and importing a basis requirement into
authorship would invent a second evidence regime.

### 4.3 Ratified elements satisfied

| Element | How |
|---|---|
| **Rule 39** [RC] | Two operations, one regime: same grant, different products |
| **K25** [RC] | Ambient authorship creates no authoritative record; promotion remains explicit |
| **Rule 27** [RC] | "Creation is not promotion" becomes demonstrable rather than vacuous |
| **Rule 40** [RC] | No shared content, no automatic promotion, no cross-chat reach |
| **O-8a-1, 8a-3, 8a-5, 8a-6** [RC] | Holder default, per-domain scope, shared regime with distinct operations, grants unchanged |
| **O-10, O-11** [RC] | No authority change by authorship; no character state touched |
| **Rules 22–23** [RC] | Creation stays open at participant level; recording now distinguishes the authorship route |

### 4.4 Possible conflicts

| Element | Risk | Reading proposed |
|---|---|---|
| **§2/§3** [RC] | `ambient_authored` is explicitly not revived | No new class; reuse `claim` + `declaration`. A distinct class would be a §3 amendment and a new approval (U-5) |
| **O-12** [RC] | Kind is declared by the model, never supplied by the system | The *operation* declares `declaration` at the call boundary; nothing is inferred from content or reclassified |
| **B3** [RC] | Ambient material is definition-like; events are not general knowledge | B3 still has no event/definition marker. C-4 does not create one and does not resolve B3 |
| **F3** [RC] | No prior record edited | Revision appends and marks the prior `superseded`; text never rewritten |
| **K27** [RC] | No transfer/sharing | Authorship is chat-contained; nothing crosses |

### 4.5 New structure required

| Element | Type | Why |
|---|---|---|
| `author_ambient(actor, domain, target, text, turn, rationale=None, revises=None)` | **Operation** [PA] | K25 needs an operation that is neither promotion nor a character utterance |
| `provenance["authorship"]` | **Provenance key** [PA] | Makes authorship and promotion separately auditable (direction, point 5) |
| Revision linkage (`revises` → `superseded_by`) | **Parameter + status transition** [PA] | "Editing must never promote" needs a defined non-promoting edit path |

**No new record class, no new contribution kind, no new field.**

### 4.6 Provenance, ownership, authority, history

The authored record names its author, domain and the grant relied on — the same authority object
promotion uses, recorded separately. No character owns it (not truth-apt, not private): it is domain
material, owned by nobody, scoped to the chat (rule 2). Characters learn it only through declaration
coverage (rules 10–11), never because it exists. Revision leaves the prior record intact and marks it
superseded; the chain is the history.

### 4.7 Execution paths

`author_ambient(actor="operator", domain=("world","w1"), target="city.capital",
text="the capital is Vareth", turn=2)` → claim record, kind `declaration`, authorship provenance.

`authoritative(("world","w1"), "city.capital")` → `none`. The material exists; nothing is authoritative.

`promote(actor="operator", domain=("world","w1"), target="city.capital",
text="the capital is Vareth", basis=[BasisEntry("record_cited", record_id="claim-00007")],
about="world", turn=3)` → authoritative; the claim is unchanged; authorship and promotion have distinct
ids and turns.

`author_ambient(actor="model:canon", …)` with no grant → **refused** (8a-1/8a-3). With a grant → allowed,
and the record still does not become truth.

`author_ambient(..., turn=8, revises="claim-00007")` → new record; prior `superseded`; **no promotion**,
and truth queries are unaffected until someone promotes.

### 4.8 Implementable vs requiring approval

| Part | Status |
|---|---|
| An authority-gated operation creating a non-truth record that never promotes | **Requires owner approval** as to shape: reuse of `claim`/`declaration` (U-5) |
| Separately auditable authorship provenance | **Requires approval** of the provenance key [PA] |
| Non-promoting revision path | **Requires approval** of the status transition [PA] |

---

## 5. Cross-contract interactions

| Pair | Interaction | Consequence |
|---|---|---|
| **C-1 × C-2** | Belief-fact promotions live in the `belief:<character>:` namespace, so two contradictory claims about what Sable believes can conflict without ever colliding with world propositions | Clean separation; groups behave identically inside the namespace |
| **C-1 × C-3** | A group changes presentation, not acquisition; characters may acquire contested content | No interaction with evidence validation |
| **C-2 × C-3** | Disclosing a belief creates a `told` account whose content is "X believes: P" | Evidence is the disclosure event; **no belief is created for the recipient** (rule 57) |
| **C-3 × C-4** | Ambient material reaches characters only through declaration coverage → `narrator_grant` | Validator must resolve `descriptor.declaration`; ambient records are citation targets for promotion |
| **C-4 × C-1** | Ambient material is non-truth-apt | It can never be a conflict member; membership is restricted to truth records |
| **All × B3** | None of the four introduces an event/definition marker | **B3 remains absent** and unresolved |

---

## 6. R2 question → proposed resolution

| R2 item | Class in R2 | Resolution proposed now | Residual |
|---|---|---|---|
| **Q-R2-1** conflict | **B** | C-1: declared proposition identity; incompatibility only by assertion; four-state answer; persisted group | U-1, U-2, U-7 + [PA] approvals |
| **Q-R2-2** belief vs truth | **C** | C-2: required `about`, namespace exclusion, live rule-56 guard, belief revocation gate | U-4; claim-path residual stated |
| **Q-R2-3** evidence | **B** | C-3: per-kind validation table, own-account exemption, no-truth-reads rule | U-3 + descriptor schema |
| **Q-R2-4** ambient authorship | **B** | C-4: authority-gated non-truth authoring, separate audit, non-promoting revision | U-5, U-8 |
| **B3** | **C** | **Unchanged** — no representation exists, none invented | U-6 |
| **C5** | **B** | **Now implementable** via C-3 | descriptor schema |
| **K25** | **B** | **Now implementable** via C-4 | operation approval |

Three of R2's B/C items move to implementable-pending-approval; B3 stays C; every move depends on a
named approval rather than an assumption.

---

## 7. Remaining unresolved architectural questions

| # | Question | Blocks |
|---|---|---|
| **U-1** | Is `target` confirmed as a **proposition** identifier (one target = one proposition), or is a separate `proposition` field required? | C-1, C-2 (rule 6 reading) |
| **U-2** | Are K7/K10 read as: a bare promotion does not displace a predecessor, and displacement requires human supersession (O-10b)? | C-1 states vs K7/J13 |
| **U-3** | Is `experienced` the plain-language name for ratified `witnessed`, or a new value? Is `narrator_grant` a sixth acquisition kind? | C-3 completeness |
| **U-4** | Is a **required** `about` parameter on promotion acceptable (K9 as minimum, not closed list)? | C-2 rule-56 guard |
| **U-5** | Is ambient material carried by `claim` + `declaration`, or does the owner want a distinct non-truth class (amending §3)? | C-4 |
| **U-6** | B3 — event/definition distinction: still no representation. Confirm it stays out of scope until an owner decision defines it | B3 |
| **U-7** | Who may **assert** incompatibility: any domain authority, human only, or any participant as a non-truth claim? | C-1 authority surface |
| **U-8** | Does the held-authority default extend to ambient authorship without a new grant, given 8a-5's shared regime? | C-4 gating |

**None is resolved by assumption.** All eight are open.

---

## 8. Proposed changes to the current implementation (exact, unexecuted)

Nothing below has been applied. Symbols are from the current build (`engine.py` 392 lines,
`records.py` 154, `authority.py` 156, `retrieval.py` 85).

### 8.1 Implementable under existing authority (defect fixes; no new architecture)

| # | File | Symbol | Change |
|---|---|---|---|
| 1 | `engine.py` | `promote(..., supersede=True)` | Close the O-10b hole: refuse it, or require the human gate as `supersede_authority` does |
| 2 | `engine.py` | `retract` | Refuse to revoke a belief unless the caller is that character's model (O-5, rule 53) |
| 3 | `engine.py` | `promote` | Replace "≥2 active ⇒ conflict" with the four-state answer; retain every record |
| 4 | `engine.py` | `promote` | Add the rule-56 refusal: a basis citing a belief may not target that belief's content |
| 5 | `engine.py` | `acquire` | Per-via evidence table; own-account exempt; inference basis must be owned by the character; no truth reads |
| 6 | `retrieval.py` | `world` | Report `unresolved_relation` distinctly from `conflict`; never `settled` for ≥2 live records |

### 8.2 Requires explicit owner approval [PA]

| # | File | Symbol | Change | Depends on |
|---|---|---|---|---|
| 7 | `engine.py` | new `assert_incompatible` | Authority-gated assertion creating a persisted group | U-7 |
| 8 | `engine.py` | `_promotions`, `snapshot` | Carry groups and their closure | U-2 |
| 9 | `engine.py` | `promote` | Required `about`; belief-fact targets namespaced | U-4 |
| 10 | `records.py` | `Record` | No new field; populate `target` on beliefs | U-4 |
| 11 | `engine.py` | new `author_ambient` | Authority-gated non-truth authoring with provenance and non-promoting revision | U-5, U-8 |
| 12 | `retrieval.py` | `world` | Exclude the belief-fact namespace from proposition answers | U-4 |
| 13 | `__init__.py` | `__all__` | Export the two new operations | 7, 11 |

**Not touched by any of this:** `eval/mem1/model.py` and the rest of the G1 harness; production Howl;
r32; plan v3; Amendment 1; v9; the map; the ledger.

---

## 9. Conformance impact — all 57 v9 rules

Codes: **U** unchanged · **M** mechanism modified/added · **N** newly satisfiable ·
**R** owner reading required · **X** still not implementable.

| Rule | | Rule | | Rule | |
|---|---|---|---|---|---|
| 1 | U | 2 | U | 3 | U |
| 4 | U | 5 | U | 6 | **R** |
| 7 | U | 8 | **M** | 9 | U |
| 10 | U | 11 | U | 12 | U |
| 13 | **N** | 14 | **M** | 15 | U |
| 16 | U | 17 | U | 18 | U |
| 19 | U | 20 | **M** | 21 | **M** |
| 22 | U | 23 | **M** | 24 | U |
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

**Notes on the 18 non-U rules**

| Rule | Note |
|---|---|
| 6 | `target` read as proposition identity (U-1) |
| 8 | Incompatibility is an assertion, never a verdict |
| 13 | Detection becomes representable |
| 14 | The assertion changes presentation, never which record governs |
| 20 | Groups add a representation; nothing erased |
| 21 | Group membership is authority-state history, not precedence |
| 23 | Authorship route becomes distinguishable from promotion |
| 27 | "Creation is not promotion" becomes demonstrable |
| 34 | The authority judges; the engine validates structure only |
| 37 | Group closure is revocation/supersession |
| 39 | The second operation exists and is defined |
| 41 | Reachable without false positives and without a pick |
| 42 | Correction/supersession/revocation kept distinct; `promote(supersede=True)` flagged |
| 43 | Human gate extended to closure; the bypass is named |
| 50 | Evidence table makes the two-field distinction enforceable |
| 53 | Belief revocation gate added |
| 55 | Belief history preserved; target populated |
| 56 | Guard exists and can fire |

**Tally: 39 U · 12 M · 5 N · 1 R · 0 X = 57.**

---

## 10. Conformance impact — all 118 ratified requirements

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
| K4 | U | K5 | U | K6 | U | K7 | **R** |
| K8 | U | K9 | **M** | K10 | **M** | K11 | **M** |
| K12 | U | K13 | U | K14 | U | K15 | U |
| K16 | U | K17 | U | K18 | U | K19 | U |
| K20 | U | K21 | U | K22 | **M** | K23 | **M** |
| K24 | U | K25 | **N** | K26 | **M** | K27 | U |
| L1 | U | L2 | U | L3 | U | L4 | U |
| L5 | U | L6 | **N** | L7 | U | L8 | U |
| L9 | U | L10 | U | — | — | — | — |

**Notes on the 30 non-U requirements**

| Requirement | Note |
|---|---|
| A4 | Belief-fact exclusion sharpens the truth answer |
| C2 | Inference basis validated against the character's own accounts |
| C4 | Read requires a source and an event |
| C5 | Wrong evidence refused at write — **newly satisfiable** |
| C7 | Own-account exemption made explicit and checked first |
| C8 | Validator forbidden from parsing text |
| D9 | Revocation gated on the declaring model |
| D10 | Enforced structurally, not by wording |
| F1, F1b, F8 | Conflict groups and authorship enter the snapshot and reconstruction |
| F3 | Status transitions defined for revision; text never rewritten |
| F6 | Evidence resolves per origin by construction |
| F7 | Canonical ordering fixed; **no replay facility exists** — determinism still unproven |
| J12–J14 | Flag-and-abstain, group representation, unresolved retrieval |
| J15 | Correction cannot close a group |
| J16 | Human closure recorded; members retained |
| K2 | Promotion records `about` alongside applicability |
| K7 | Settling read per proposition (U-2) |
| K9 | `about` added to recorded fields (K9 is a minimum) |
| K10 | Chain defined as supersession, not bare promotion |
| K11 | Group membership joins preserved history |
| K22, K23 | Revocation and supersession close groups explicitly |
| K25 | Ambient authorship exists and does not promote |
| K26 | The `promote(supersede=True)` bypass is named and closed |
| L6 | Creation vs authority demonstrable |

**Tally: 88 U · 19 M · 8 N · 2 R · 1 X = 118.**

---

## 11. Three-way classification

### 11.1 Existing ratified requirements (unchanged, not re-opened)

57 v9 rules · 118 map requirements · D-01, D-02 · O-1, O-3, O-4, O-5, O-6a, O-6b, O-7, O-8a-1…6,
O-10, O-10b, O-11, O-12, O-13 · r32 as tier 2 · D-4 (Option B) · D-1/D-2/D-3 register entries.

Nothing in this file edits, extends or settles them. The two interpretive readings it relies on most
(U-1, U-2) are submitted **as questions**, not applied.

### 11.2 Owner-approved direction (this round, quoted)

1. "Conflict detection must distinguish actual incompatibility from mere shared subject, target, or domain."
2. "Beliefs must never silently become world truth through retrieval, disclosure, promotion, or association."
3. The five acquisition kinds, with the invariant "evidence that information was acquired does not establish that the information is true."
4. "Implement ambient authorship as a distinct authoring operation, separate from promotion."

### 11.3 Genuinely new proposals — **not added to the ratified set**

| # | Proposed requirement | Status |
|---|---|---|
| **NR-1** | Where incompatibility is not established, the answer for a proposition is `unresolved_relation`, never `settled` and never `conflict` | Proposed; set stays at 118 |
| **NR-2** | A promotion citing a belief may govern only a fact about that character, and is ineligible for the belief's own proposition | Proposed |
| **NR-3** | Every acquired account is refused at write unless its kind's evidence is present; evidence is never truth | Proposed mechanism for existing C5 |
| **NR-4** | Ambient-authored records are non-truth-apt, separately attributable from promotions, and promotable only by a separate operation | Proposed mechanism for existing K25 |

**Structural elements proposed, each needing approval:** conflict-group object · `assert_incompatible` ·
`unresolved_relation` status · required `about` parameter · belief target namespace · descriptor key
schema · `author_ambient` · authorship provenance key · non-promoting revision linkage.

**Nothing above is authorised. The count remains 57 rules and 118 requirements.**

---

## 12. What this record does not claim

- **No implementation.** No file under `eval/mem1/` was modified.
- **No execution.** No test, probe, suite or self-test ran.
- **No behavioural verification.** Every statement about current behaviour is source inspection ([SO]),
  not observed execution.
- **No ratification.** The four contracts are proposals; the owner approved directions, not mechanisms.
- **No resolution by assumption.** U-1…U-8 remain open, including the two readings this file depends on
  most heavily.
- **No amendment** to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9, the map,
  or the ledger.
- **No transfer, sharing or carryover** is proposed anywhere; K27 and O-13 are untouched.

**Next step in the required sequence:** owner review of these four contracts → decisions on U-1…U-8 →
implementation authorization → implementation → conformance review → separate testing authorization.
