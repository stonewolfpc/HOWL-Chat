# MEMORY-09 — Phase 1A-1: Retraction and Knowledge Persistence Conformance Review

**Status:** research record. Read-only. **No implementation, no test, no probe, no model call, no
corpus, no canonical edit.** Nothing here adopts a model, retains or retires a class, adds a register
entry, or decides PH1A-02.

| Field | Value |
|---|---|
| Created | **2026-11-19** |
| Authority | None beyond the owner's work order for Phase 1A-1. Creates no decision. |
| Scope | What must happen to a character's knowledge when the information they acquired is retracted — in r32, v9, the owner ledger, the 118 requirements, and the two implementations in `eval/mem1` |
| Predecessor | `MEMORY-09-PHASE-1A-CONCEPTUAL-OBJECT-INVESTIGATION.md` (CO-01, CO-03, CO-06, PH1A-02) |
| Inputs verified | Hashes in §3, read from disk this session |
| Tests executed | **0.** Every behavioural statement below is a static reading of source, not a run. |
| Files changed | This file only. Nothing staged, committed, or registered. |

---

## 1. What this review adds, and what it deliberately does not

The Phase 1A investigation established that `eval/mem1` contains two object models: the harness drives
`model.py`, the engine is `m09/`, and no bridge exists between them (CO-01). It also established that
r32 and v9 each hold one half of the knowledge guarantee, and neither holds both.

This review takes the single behaviour where that split is observable end to end — **retraction** — and
asks what each source actually requires, what each implementation actually does, and where the
difference is a gap rather than a choice.

Three results were not visible in Phase 1A and are the reason this document exists:

1. **The two implementations do not disagree about retraction propagation. Only one of them has
   propagation to disagree about.** The harness never creates knowledge for a listener at all, so its
   apparent conformance with recipient-persistence requirements is vacuous, not demonstrated
   (RP-03, RP-04).
2. **r32 §8.4 and conformance row E1c conflict for the recipient**, and the conflict is invisible when
   either is read alone (RP-09). This is the substantive question underneath PH1A-02, and it cannot be
   answered by choosing r32 or v9 — both readings survive either choice.
3. **The hybrid relation the owner proposed is not excluded by any tier-1 decision.** It is excluded by
   one tier-3 design choice in v9 whose stated rationale is canon-enumeration safety, a rationale the
   hybrid need not violate (RP-11, RA-3).

This review does **not** recommend an architecture. It separates what is already required from what is
genuinely undecided, and it frames the decision rather than taking it.

---

## 2. Method and evidence labels

Every claim carries one label and, where it is a claim about a document, a tier under O-6a.

| Label | Meaning here |
|---|---|
| **Observed** | Read directly in the working tree at the hash in §3. A statement about code as written, never about code as run. |
| **Source-derived** | Entailed by quoted text. The quotation is given. |
| **Calculated** | Arithmetic over counted items (e.g. probe counts, line counts). |
| **Inferred** | A judgement this review makes and invites challenge on. Never presented as a requirement. |
| **Unestablished** | Not decidable from the sources read. Recorded as a question, not filled with a plausible answer. |
| **Superseded** | Replaced by a later record; cited only as history. |

**Tiers (O-6a, T1 ledger `da3e7239…`):** T1 owner-confirmed decisions and r32 D1–D7 · T2 r32's `K-R`
mechanisms and invariants · T3 v9's 57 rules, ratified as design baseline under O-6b.

**Precedence rule applied, not invented:** a T3 rule contradicting a T1 decision is **void**, not a
divergence (O-6a case 1); a T3 rule contradicting an applicable T2 mechanism is a **divergence**, and
neither side is implemented until the owner decides (case 2). This review applies that rule to
precedence only. It does not use it to settle substance, and it does not treat a later document as
authoritative because it is later.

**ID namespaces** are new and collision-free with Phase 1A (`CO-`, `PH1A-`, `TR-`, `RX-`, `HC-`):

- `RP-nn` — findings
- `RQ-nn` — open questions requiring an owner decision
- `RA-n` — relationship alternatives
- `RT-nn` — implementation traces (Appendix A)

---

## 3. Sources read, with hashes and tiers

| Source | SHA-256 (first 16) | Bytes/lines | Tier | Role here |
|---|---|---|---|---|
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` (r32) | `e88ef2d1697da1d9…` | 91,311 B / 1,028 L | T2 | K-R1–K-R19, R1b/R1c, §8.4, S34–S38 |
| `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | `698254ca262c0874…` | 236 L | T3 | rules 12, 15, 20, 29, 37–38, 42, 46–57 |
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | `da3e72398ecbaa0a…` | 636 L | T1 | D-01, O-1, O-3, O-4, O-5, O-8a-4, O-10b, O-13 |
| `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | `b4ba33ccbfe8f8c1…` | 211 L | T1 (ratified requirements) | A2a, A2b, B6, C1, E1a–E1c, E2, E3, H10, K16, K24 |
| `MEMORY-09-PHASE-1A-CONCEPTUAL-OBJECT-INVESTIGATION.md` | `82173a6f78ac84c4…` | 1,071 L | research | Predecessor; CO/PH1A IDs reused, not redefined |
| `eval/mem1/m09/engine.py` | `7505e9acca8e3229…` | 1,209 L | Observed | `retract`, `disclose`, `acquire`, `revoke_belief` |
| `eval/mem1/m09/records.py` | `ad2d4d11501c357b…` | 252 L | Observed | `Record`; `proposition_id` refusal on private records |
| `eval/mem1/m09/retrieval.py` | `cd2d380bd5d9d1f7…` | 101 L | Observed | `perspective` pool construction |
| `eval/mem1/model.py` | `813cdbff56c01ae4…` | 398 L | Observed | `MemoryRecord`, `AcquisitionEvent`, `KnowledgeEdge`, `retract`/`revoke` |
| `eval/mem1/retrieve.py` | `84f90b5ce0b3619f…` | 235 L | Observed | `eligible`, `retrieve` |
| `eval/mem1/form.py` | `f472052738388bab…` | 242 L | Observed | `_retract`, `_admit` (`via="witnessed"`) |
| `eval/mem1/conversation.py` | `fa1d9ce630462fa9…` | 414 L | Observed | driver; per-turn retrieval; no listener formation |

`MEMORY-09-ARCH-03-R3-INDEPENDENT-REVIEW-RECORD-2026-11-17.md` F5 is cited once, as a prior
observation about the same §4/§8.4 tension. It is not treated as authority.

---

## 4. The question, decomposed

### 4.1 Five operations that are routinely spoken of as one

The work order asks for these to be separated. They are separated here by **what changes**, **what is
preserved**, and **whether a character's knowledge changes**.

| Operation | What changes | What is preserved | Character knowledge |
|---|---|---|---|
| **Forgetting** | whether this character currently knows | the fact they once knew (event log) | decreases — the character no longer knows |
| **Source retraction** | whether the source record is asserted | the record, its provenance, every knower's account | **does not change** (T1 O-3; T2 K-R13; T3 rule 49) |
| **Account retraction** | whether one character's account is held | the record row itself, under status | decreases — this owner loses this copy |
| **Correction** | a record's content | the prior text, in history | does not change |
| **Supersession** | which authoritative record governs | both records, the chain | does not change (T3 rule 38 / K24) |
| **Promotion revocation** | whether a promotion remains authoritative | the promotion and its provenance | does not change (T3 rule 38 / K24) |

Two consequences follow immediately, and both are requirements rather than opinions:

- **Forgetting and source retraction are opposites in their effect on knowledge** and must not share an
  operation. One removes knowledge; the other must not (RP-01).
- **Account retraction is not source retraction.** It is the only operation in this list that removes a
  character's copy while leaving the source intact — and it is the one operation with no owner decision
  and no r32 equivalent (RQ-03).

### 4.2 Three acquisition circumstances

The work order asks whether one mechanism must serve all three. The sources do not answer uniformly.

| Circumstance | How r32 represents it | How v9 represents it |
|---|---|---|
| **Publicly available** | `ambient_authored`; audience derived from class + scope; often no event at all (§8.2) | declaration coverage; `narrator_grant` accounts created per named recipient (B1–B9) |
| **Private disclosure** | `told` event → edge; `via` precedence chain | `disclose()` → per-recipient `acquired` account with descriptor |
| **Independent discovery** | one shared memory, two `witnessed` edges | two independent accounts, no relation between them |

The third row is where the models diverge most sharply, and it is the row Phase 1A's CO-03 did not
reach: under v9 two characters who witness the same event hold **two records that the system cannot
recognise as being about the same thing** (RP-06).

### 4.3 Five guarantees that the alternatives trade against each other

These are the properties that cannot all be held at once by the current models. They are stated
separately so that §9 can score each alternative against them rather than against a preference.

- **G1 — Exclusive private ownership.** A's original record belongs to A alone (T1 O-1).
- **G2 — Recipient persistence.** B's account survives retraction of A's source (T1 O-3; T3 rule 49;
  requirements A2a, E1b, H10).
- **G3 — Retraction visibility.** A knower can be shown that what they hold has been retracted
  (T2 §8.4 / S34; requirement A2b read narrowly).
- **G4 — No resolvable pointer.** B's account never resolves to A's record, its contents or its history
  (T3 rule 12; T1 O-1; requirement J8).
- **G5 — Content sameness.** The system can tell that two accounts concern the same underlying claim
  (needed for G3 for recipients, and for conflict detection across knowers; **no requirement states it**).

---

## 5. What each source requires

### 5.1 Tier 1 — owner decisions

O-3 verbatim, on the point at issue:

> "Once disclosure occurs, subsequent retraction of the source does not erase the recipient's acquired
> account."

and, immediately after, the qualification:

> "Irreversibility applies to the **historical disclosure event and the recipient's acquired account**.
> It does **not** mean the recipient retains access forever regardless of other legitimate lifecycle
> operations."

O-1 establishes G1 and G4 together: B receives "a separate acquired record containing only what was
disclosed" and "gains no access to A's original, its unrelated contents, or its derivation history."

O-5 establishes that a data update never revises a belief, and that only the character's own model may
declare one. This constrains what any retraction notice may assert (RP-10).

### 5.2 Tier 2 — r32

The mechanism is stated in four places, and they are consistent with each other.

| Rule | Text (abridged) | What it fixes |
|---|---|---|
| **K-R3** | "Revocation is a new event (`via=revoked`), never a deleted history row. Forgetting changes current knowledge; it never rewrites what happened." | forgetting is an event, not a status write |
| **K-R11** | the edge fold depends **only** on acquisition-family events for (character, memory) | truth status may not enter the knowledge projection |
| **K-R12** | memory lifecycle is its **own authoritative stream** — `established \| superseded \| corrected \| retracted` — and status is a **projection**, not a mutable field | retraction is an event on a second stream |
| **K-R13** | "'This character's knowledge includes a retracted claim' is a **read-time join** of edge × memory projection. It is never stored on the edge." | G3 is computed, never stored |
| **K-R19** | "A revoked pair has **no edge**, not an edge with a 'revoked' flag." | forgetting removes the relation |

And the worked case, §8.4, verbatim:

> "`memory(retracted)` + A's event unchanged + projection live ⇒ A's channel: *'you were told X; it has
> since been retracted; you still believe it'*. Another character who knows the correction behaves
> differently. No mind is corrected by a data update."

The consequence paragraph under K-R13 is the load-bearing sentence for this review:

> "a retraction changes the memory stream only. The edge's `derived_from_event_ids` is unchanged, the
> character still knows the (now retracted) claim, and whether they have *learned of the retraction* is
> itself an acquisition event."

So r32 holds **G3 for any knower**, and pays for it with a resolvable memory id on every edge — which
is exactly **G4**'s prohibition under v9's rule 12.

### 5.3 Tier 3 — v9

| Rule | Text (abridged) | Effect on retraction |
|---|---|---|
| 12 | "recipient records carry a restricted origin descriptor, never access to the source record" | severs the join that G3 needs |
| 38 | "revocation or supersession alters no memory, knowledge or belief" | G2 holds for authority operations |
| 49 | "the disclosure event and the recipient's acquired account persist after any retraction of the source, and this does not preclude other legitimate lifecycle operations" | G2 holds for source retraction |
| 15 / 55 | belief is not revised by acquisition; revision is append-only | a retraction cannot reach a belief |
| 20 | "nothing is erased" | status changes are not deletion |

v9 states G2 and G4. It states nothing that provides G3 for a recipient, and nothing that provides G5.
That silence is not a contradiction of r32 — it is the gap-fill boundary O-6a case 3 describes, except
that here the gap-fill removes a behaviour r32 required rather than adding one r32 lacked.

### 5.4 The ratified requirements

| ID | Requirement | What it fixes |
|---|---|---|
| **A2a** | account exists; memory later retracted → character still holds what it acquired | G2 |
| **A2b** | same character, same turn → **no automatic knowledge of the retraction** | the retractor is not credited with knowledge they did not acquire |
| **E1a** | retraction → "only the memory's truth status changes" | no side effects |
| **E1b** | accounts held by knowers → unchanged | G2 |
| **E1c** | retraction turn → **"no notice account created"**; fails if "retraction knowledge appears unaccounted" | forbids a *stored* notice |
| **E2** | supersession → both retained, older knowledge historical | persistence under supersession |
| **E3** | private record retracted → only that character's record changes | no domain-truth side effect |
| **H10** | source retracted after disclosure → recipient account and event persist | G2 for disclosure |
| **B6** | declaration withdrawn → future coverage stops, existing accounts intact | G2 for coverage |
| **K16** | basis retracted after promotion → promotion stands | no automatic undo |
| **K24** | authority change leaves characters untouched | no character-state side effect |

Note what is **absent** from the 118: no row requires a recipient to be *shown* that their account's
source was retracted. G3 appears only through A2b, which concerns the retracting character, and through
E1c, which forbids unaccounted retraction knowledge. This absence is the finding in §7.

---

## 6. Observed behaviour of the two implementations

### 6.1 The harness (`model.py`, `form.py`, `retrieve.py`, `conversation.py`)

**RT-01 — A retraction is formed only from the retracting speaker's own span.**
`Determiner.run` computes `owner_id = owner_of.get(span.speaker)`; retraction is detected only when
`span.speaker != "Narrator"`, and `_retract` searches `self.store.current(owner_id=owner_id)`
(`form.py:174, 187–204`). The retraction therefore revokes **the speaker's own** active records.
No other character's records are visited, because no other character has any.

**RT-02 — The harness creates no knowledge for a listener, ever.**
`_admit` is called only for the speaker's span, and records `via="witnessed"` unconditionally
(`form.py:206–240`). There is no acquisition path where B acquires from A's utterance: B's turns form
B's own records from B's own speech. Consequences, all Observed:

- requirements whose precondition is "B was told" (C1, H4, H10, J5) cannot be exercised;
- the harness's apparent satisfaction of A2a/E1b/H10 is **vacuous** — there is no recipient account to
  remove or preserve;
- `via` does not record how anyone learned anything; it records that a record was formed.

This corrects a reading available from Phase 1A: the harness is not a weaker implementation of
r32's relation. It is an implementation of a **single-knower** relation, in which the question "what
happens to B when A retracts" has no representation at all.

**RT-03 — Retraction writes status on the target and appends two `revoked` events.**
`MemoryStore.retract` calls `revoke()` (sets `status=REVOKED`, `revoked_at_turn`, appends an event with
`via="revoked"`), then adds a second record whose `predicate` is `retraction:<target>` and whose
`retracts` names the target (`model.py:282–318`). The target keeps its span and is never deleted.

**RT-04 — The edge projection contradicts K-R19.**
`MemoryStore.edges()` emits one edge per **event**, including revocations
(`model.py:347–353`). A revoked pair therefore still yields an edge, with `via="revoked"`, and the
retraction event yields a second one. K-R19 requires "no edge", and K-R3's "revocation is a new event"
is satisfied — so the harness satisfies the event rule and violates the projection rule. Observed;
Source-derived from K-R19.

**RT-05 — The retraction is visible to exactly one character: the one who said it.**
Because the retraction record is owned by the retracting speaker, is `ACTIVE`, and inherits the target's
scope (`"owner"`), `eligible()` admits it in `character_perspective` for that owner only
(`retrieve.py:100–143`). The frozen fixture's single retraction probe is asked of Alaric, who is the
speaker of the retraction (`fixtures.py:189, 218, 242`). **No probe asks a recipient**, so no current or
planned measurement covers G3 for a second character. Observed; Calculated (one retraction probe in the
fixture).

**RT-06 — A world/narrator record can never be retracted.**
`_retract` searches only `current(owner_id=owner_id)` for the speaker's id, and the narrator has no
`owner_of` entry; narrator records are owned by the world or `chat`. A retraction sentence in narration
is skipped by `span.speaker != _NARRATOR`. Observed.

### 6.2 The engine (`m09/`)

**RT-07 — `retract()` has no actor, no owner check, and no authority check.**
`Engine.retract(self, record_id: str, turn: int)` refuses beliefs and truth records, then sets
`rec.status = REVOKED` (`engine.py:542–559`). Any caller holding the id may retract any claim or account
in the chat, including an account belonging to another character. Observed. This is the Phase 1A CO-06
finding, restated here because it is the mechanism an account-retraction decision would have to govern.

**RT-08 — Recipient accounts survive a source retraction, and the owner keeps nothing visible.**
`disclose()` creates independent accounts with descriptor `{speaker, record_class, turn, event}` — no
source id (`engine.py:575–700`). Retracting the source claim sets the **claim** to `REVOKED`; B's
account row is untouched. But `retrieval.perspective` excludes `status == "revoked"` from the only
character surface (`retrieval.py:46–52`), so if the retracted record is the character's own account, the
row survives in `_records` and in `history()` while becoming **unreachable to its owner through
retrieval**. G2 is therefore true at the storage layer and qualified at the channel layer (RP-02).

**RT-09 — Re-disclosure of a retracted claim is permitted.**
`disclose()` checks `src.chat_id`, `src.owner_character == discloser`, and
`src.record_class in NON_TRUTH_CLASSES`. It does **not** check `src.status` (`engine.py:575–700`).
A character may therefore disclose a claim they have just retracted, and each new recipient receives a
fresh active account. Observed. Whether this is a defect or a legitimate consequence of "history is not
erased" is **not established** by any source (RQ-05).

**RT-10 — There is no knowledge-side event stream.**
The engine stores account rows; it appends no acquisition or lifecycle events. `retract` changes a
mutable `status` field with no accompanying event. K-R12's "status is a projection of an append-only
stream" has no counterpart here — consistent with v9 never adopting K-R12, and inconsistent with r32 if
r32 governs. This is Phase 1A's PH1A-04, confirmed by inspection rather than inherited.

**RT-11 — No surface can answer "what does A know about a retracted thing".**
There is no cross-stream join: `authoritative()` reads only truth records; `perspective()` reads only
the character's own rows; `accounts_of()` returns rows without consulting any source's status. Observed.
This is the structural absence behind RP-09.

### 6.3 What both implementations share

Both express a change in *knowledge* by writing a status on the *record*. Neither implements forgetting
as an event distinct from a record-status write. r32 requires the two to be different streams
(K-R3, K-R11, K-R19). Observed; Source-derived.

---

## 7. The conflict that appears only when the sources are composed

This is the central finding of the review, and it is not a disagreement between r32 and v9. It is a
disagreement between two requirements that were each written for a different character.

**r32 §8.4** requires the knower's channel to display: *"you were told X; it has since been retracted"*
— computed by a read-time join (K-R13), with no stored notice.

**Requirement E1c** requires: *"No notice account created"*, and fails if *"retraction knowledge appears
unaccounted"*.

**Requirement A2b** requires: *"No automatic knowledge of the retraction"*, and fails if *"retraction
knowledge appears unaccounted"*.

Read A2b and E1c as written, and the character who is *shown* a retraction they were never told about has
knowledge that is unaccounted for — which is precisely what both rows forbid. Read §8.4 alone, and that
display is mandatory for every live projection.

Three readings are available, and **none is selected by any source**:

| Reading | What it permits | What it forbids | Text it rests on |
|---|---|---|---|
| **R-narrow** | A derived join may display the retraction, because it is computed from the source's status, not acquired | any *stored* notice account | E1c's word "account"; K-R13's "never stored on the edge" |
| **R-broad** | nothing automatic; a knower learns of a retraction only by being told | the §8.4 display, for any knower with no acquisition event | A2b/E1c's "appears unaccounted" |
| **R-split** | the retractor may be shown (they performed the act); other knowers may not | the §8.4 display for recipients | A2b's "same character" as the scope of the row |

Under R-broad and R-split, **§8.4's recipient display is void for anyone who was not separately told**,
and v9's current behaviour is correct — not because v9 designed it, but because it never built the join.
Under R-narrow, **v9 is non-conformant** and cannot be repaired without G5, because the join needs an
identity that survives the disclosure copy (RP-11).

This is why PH1A-02 cannot be answered as a binary. Choosing "shared reference" or "copied account" does
not decide whether the notice is owed; both readings remain open under either choice. The reading must be
decided separately, and it is the decision that determines whether G5 is required at all (RQ-01).

**A second, independent conflict sits inside §8.4's own sentence.** The display text ends *"you still
believe it"*. Under T1 O-5 and T3 rule 38, the memory system may not assert a character's belief; only
that character's model may declare one. §8.4 is T2; O-5 is T1. Under O-6a case 1, the belief clause
cannot bind the memory system. The legitimate, precedence-surviving part of §8.4 is *"you were told X;
it has since been retracted"*. This is an application of the recorded precedence rule, not a new
decision, and it is reported rather than repaired (RP-10).

---

## 8. Behaviour comparison

### 8.1 By operation, across the four sources

| Behaviour | r32 (T2) | v9 (T3) | Harness | Engine |
|---|---|---|---|---|
| Source record after retraction | status projected from a `retracted` lifecycle event (K-R12) | `status = REVOKED`, mutable field, no event | `status = REVOKED` + `revoked_at_turn`; two appended events | `status = REVOKED`, no event |
| Knower's account after source retraction | unchanged — edge intact (K-R13) | unchanged — separate row | **no knower exists** | unchanged — separate row |
| Knower shown the retraction | yes, read-time join (§8.4) | **no mechanism** | only the retractor, via an owned event record | **no mechanism** |
| Notice stored? | no — computed (K-R13) | n/a | yes — a stored `retraction:*` record | n/a |
| Forgetting | new `revoked` event; edge disappears (K-R3, K-R19) | no distinct operation; `retract()` on the account | `revoke()` + `retraction:*` record; edge still emitted (RT-04) | `retract()` on the account row |
| Correction | lifecycle event; content changes | `correct()` — truth records only, human-gated | `state_change` record supersedes prior | `correct()` — truth only |
| Supersession | lifecycle event; both retained | truth only (`supersede_authority`), human-gated | any record, same predicate | truth only |
| Promotion revocation | explicit, human-controlled (O-8a-4, O-10b) | `revoke_promotion`, `_require_human` | not represented | `revoke_promotion`, `_require_human` |
| Belief affected by retraction | no (§8.4; S34) | no (rule 15/38; beliefs not retractable) | no belief class exists | no — `retract` refuses beliefs |
| Who may retract | not stated for private records | **any caller** (RT-07) | formation-layer speaker confinement only | **any caller** |

### 8.2 By acquisition circumstance

| Circumstance | r32 | v9 | Harness | Same behaviour? |
|---|---|---|---|---|
| **Public / ambient** | audience derived from class+scope; no edge for many knowers; retraction removes it from every derived audience at once | coverage creates real per-recipient accounts; retraction of the source leaves every account active and no one can be told | world records are `world`-scoped and cannot be retracted at all (RT-06) | **No — opposite persistence** |
| **Private disclosure** | edge per recipient; retraction visible by join; account persists | copied account; retraction invisible; account persists | not representable (RT-02) | **No — G3 differs** |
| **Independent discovery** | one memory, two edges; both knowers see the retraction | two unrelated accounts; sameness unrecognisable | each speaker forms their own record; no relation | **No — G5 absent in v9** |

The first row is the one the work order asked about directly, and the answer is that the sources **do**
already treat the circumstances differently — but inconsistently with each other, and in the public case
the difference is not a design choice anyone recorded. v9's B6 governs *declaration withdrawal*, not
*retraction of a declared record*; there is no requirement for the latter (RP-07).

---

## 9. Relationship alternatives

Three alternatives are stated, including the hybrid the owner proposed. Each is scored against G1–G5.
"Pass" means the alternative can express the guarantee without further invention; "qualified" means it
can, but with a stated cost; "fail" means it cannot.

### RA-1 — r32: one shared memory, one edge per knower

- Knowledge is a reference: `edge(character, memory_id)`.
- Retraction is a lifecycle event on the shared memory; the join at read time produces the notice.

| G1 ownership | G2 persistence | G3 visibility | G4 no pointer | G5 sameness |
|---|---|---|---|---|
| **fail** — the memory is chat/world-owned; no character owns the original (Phase 1A CO-04) | pass — edge intact | **pass** | **fail** — the edge carries `memory_id` | **pass** — one id |

### RA-2 — v9: one copied account per knower

- Knowledge is a copy: each account has its own identity and an unresolvable descriptor.
- Retraction of the source cannot reach the copy, because nothing links them.

| G1 | G2 | G3 | G4 | G5 |
|---|---|---|---|---|
| **pass** | pass | **fail** | **pass** | **fail** |

### RA-3 — the hybrid: shared content identity, private copies

This is the owner's proposal, stated here as a candidate and explicitly **not** as a recommendation.

- Each account keeps its own identity, ownership, text and descriptor (G1, G2, G4 unchanged from RA-2).
- Each account additionally carries a **content identity** — a shared name for "the claim/content this
  account is an account of" — which is **not** a record id and does not resolve to any record.
- Retraction is recorded against the content identity (or against a source that declares one).
- A knower's channel may then compute: "the account you hold concerns content C; C is no longer
  asserted by its source" — without naming the source, its owner, or its unrelated contents.

| G1 | G2 | G3 | G4 | G5 |
|---|---|---|---|---|
| **pass** — accounts stay private | pass — copies unchanged | **pass, if the notice is impersonal** | **pass, if the identity is non-resolvable** | **pass** |

**What RA-3 costs, stated plainly.** G3 and G4 are both "pass" only under conditions that are not yet
defined:

1. **The identity must be established, not inferred.** Three sub-options exist and each has a different
   failure mode: (a) identity from exact communicated text — loses paraphrase and quotation drift;
   (b) identity declared by the conversation layer, alongside `content` and `recipients`, which matches
   the existing pattern that the layer establishes disclosure facts (B2) and never the memory system
   (rule 26); (c) identity as a v9 `proposition_id` — **excluded today**, because `records.py` refuses
   `proposition_id` on any private record.
2. **The notice must be impersonal or it becomes a disclosure.** "Content C is no longer asserted" tells
   B nothing about A. "A retracted what they told you" discloses A's act to B, which under v9 rule 11
   requires an acquisition event — and an acquisition event would satisfy E1c while creating exactly the
   notice account E1c forbids. The attributed form is therefore blocked twice over; the impersonal form
   is not blocked, and is also not required by any text.
3. **It adds a third identity type.** Records have identities; v9 has propositions; RA-3 adds content
   identities scoped to private records. That is new architecture under O-6a case 3 and would have to be
   declared as such, not attributed to r32.

**What RA-3 does not solve.** It does not decide RQ-01 (whether a notice is owed at all). It does not
give the *system* any authority to tell B that A retracted — only to report the status of a content
identity B already holds. And it does not make "who knows C" answerable across characters unless the
identity is enumerable, which would reintroduce the G4 problem the design was avoiding.

### 9.1 The guarantee matrix, together

| | G1 ownership | G2 persistence | G3 visibility | G4 no pointer | G5 sameness |
|---|---|---|---|---|---|
| **RA-1** r32 shared | fail | pass | pass | fail | pass |
| **RA-2** v9 copied | pass | pass | fail | pass | fail |
| **RA-3** hybrid | pass | pass | pass* | pass* | pass* |

`*` conditional on the three undefined points above. No alternative passes all five unconditionally, and
no source requires all five. **That is the decision.** It is recorded as RQ-01/RQ-02 and not resolved.

---

## 10. Findings register

| ID | Finding | Label | Source | Consequence |
|---|---|---|---|---|
| **RP-01** | Forgetting and source retraction have opposite effects on knowledge, and neither implementation separates them; both write record status instead of an event on a knowledge stream | Observed + Source-derived | K-R3, K-R11, K-R19 vs `model.py:282–318`, `engine.py:542–559` | PH1A-03/04 confirmed; any design contract must name which stream a status write belongs to |
| **RP-02** | In the engine, a retracted record survives in storage and in `history()` but is excluded from the only character retrieval surface, so "still holds" is true of the data and false of the channel | Observed | `retrieval.py:46–52`; `engine.py:1181` | A2a holds at rest, is qualified in retrieval; needs an owner reading of what "holds" means |
| **RP-03** | The harness has no listener acquisition: every record is formed for its speaker with `via="witnessed"` | Observed | `form.py:206–240`; `conversation.py:191–201` | Its conformance with C1/H4/H10/J5 is vacuous, not demonstrated |
| **RP-04** | The harness's single retraction probe is asked of the retracting character; no probe covers a recipient | Observed + Calculated | `fixtures.py:189, 218, 242` | G3 for recipients is unmeasurable under the frozen fixture |
| **RP-05** | The harness's `edges()` emits an edge for a revoked pair, and a second for the retraction record, contradicting K-R19 | Observed | `model.py:347–353`; K-R19 | If r32 governs, the projection is non-conformant; if v9 governs, the stream is unused |
| **RP-06** | v9 cannot represent that two accounts concern the same content; independent discovery produces unrelated copies | Observed | `records.py` `proposition_id` refusal; `engine.py` `acquire` | G5 absent; retraction can never reach co-knowers; cross-knower conflict detection is impossible |
| **RP-07** | No requirement covers retraction of a *declared* record; B6 covers withdrawal of the declaration only | Source-derived | Conformance map B6; no other row | Public-information retraction has no defined behaviour (RQ-04) |
| **RP-08** | The engine's `retract()` accepts any record id with no actor, ownership or authority check | Observed | `engine.py:542–559` | Account retraction is ungoverned; O-10b's rationale is not applied to accounts (RQ-03) |
| **RP-09** | §8.4's recipient notice and A2b/E1c's "no unaccounted retraction knowledge" conflict for the recipient, and three readings are each defensible | Source-derived | r32 §8.4, K-R13; conformance A2b, E1c | **The substantive decision beneath PH1A-02** (RQ-01) |
| **RP-10** | §8.4's "you still believe it" cannot bind the memory system: O-5 (T1) forbids the system asserting a belief | Source-derived + Inferred | r32 §8.4; v9 rules 15/38; O-5 | Precedence resolved by O-6a case 1; substance reported, not repaired |
| **RP-11** | The hybrid relation is excluded today by a tier-3 choice whose stated rationale is canon-enumeration safety, not by any tier-1 decision | Observed + Inferred | `records.py` module docstring; O-1; rule 12 | RA-3 is not forbidden by T1; it would be new architecture under O-6a case 3 |
| **RP-12** | A retracted claim can still be disclosed; `disclose()` never checks source status | Observed | `engine.py:575–700` | Retraction does not inhibit further propagation (RQ-05) |
| **RP-13** | The engine has no cross-stream join, so "A knows a retracted thing" is not answerable on any surface | Observed | `engine.py` `authoritative`, `accounts_of`; `retrieval.py` | K-R13 has no implementation anywhere |
| **RP-14** | The harness's `via` is asserted by the formation rule, not evidenced by an acquisition event | Observed | `form.py:206–240` | C1/C9's "origin told, teller cited" cannot be satisfied by this field |

---

## 11. Established requirements vs open decisions

### 11.1 Established — no owner decision needed

These follow from T1/T2/T3 text and are not reopened here.

1. A source retraction does not remove a recipient's account (O-3; rule 49; A2a, E1b, H10).
2. A retraction changes truth status only (E1a), and does not touch beliefs (O-5; rule 38; S34).
3. No stored notice account is created (E1c, under the narrow reading) — the notice, if any, is computed.
4. A character who learns of a retraction has acquired that knowledge; it does not arrive by itself
   (K-R13 consequence; A2b).
5. Retracting a promotion's basis does not undo the promotion (K16; rule 29).
6. Authority changes never rewrite character state (rule 38; K24).
7. History is preserved under every operation in §4.1 (rule 20; K-R3).

### 11.2 Open — requires the owner

| ID | Question | Why it cannot be settled from source | What it blocks |
|---|---|---|---|
| **RQ-01** | Is a knower owed a retraction notice at all, and if so is it derived (R-narrow) or acquired (R-broad)? | §8.4 and E1c cannot both bind for the recipient; no text ranks them | Whether G3/G5 are required; whether v9 is non-conformant |
| **RQ-02** | If a notice is owed, is it impersonal ("content C is no longer asserted") or attributed ("A retracted it")? | The attributed form discloses A's act and requires an acquisition event, which E1c forbids storing | Whether RA-3 can stay within rule 12 |
| **RQ-03** | Who may retract a character's *account*? | No source states it; `retract()` has no gate; O-10b governs promotions, not accounts | The only operation that can remove a character's knowledge |
| **RQ-04** | What happens to coverage accounts when the *declared record* is retracted? | B6 covers withdrawing the declaration, not retracting the record | Public-information retraction has no defined behaviour |
| **RQ-05** | Does retraction bind future re-disclosure of the retracted claim? | No requirement; `disclose()` does not check status | Whether RP-12 is a defect or a legitimate historical act |
| **RQ-06** | Must "same content" be representable across characters at all? | No requirement states it; needed only if RQ-01 is answered "yes" | Whether a content identity is required, and which of the three sub-options |

### 11.3 Decision brief for PH1A-02 (framing only — no option is recommended)

Phase 1A posed PH1A-02 as "shared reference (r32) or copied account (v9)". This review shows that framing
is necessary but not sufficient, because the answer does not determine the notice behaviour. The brief
should therefore put **three** questions, in this order, and stop:

**Q-A. The relation.** Which of RA-1, RA-2, RA-3 is the intended relation between a character's knowledge
and the thing known? The material consequence of each is already tabulated in §9.1. The question the
owner should be asked is not "which model wins" but: *is it acceptable that the system cannot say that B
and C know the same thing?* If yes, RA-2 stands and RQ-01 is moot. If no, RA-1 or RA-3 is required.

**Q-B. The notice.** If a knower's source is retracted, is the knower's channel permitted to say so
without having been told? This is RQ-01. It is the question §8.4 answers "yes" and E1c answers "no", and
it must be answered explicitly rather than inferred from the architecture choice, because either answer
is compatible with either relation.

**Q-C. The attribution.** If the answer to Q-B is yes, does the notice name the retracting actor? This is
RQ-02, and it is the difference between a status report and a disclosure, with the acquisition and
privacy consequences that follow.

**Q-D. The conversation-state relation.** Added by the follow-up recorded in §13: does MEMORY-09 track
conversation states at all, and where does applicability live? See `RQS-01`/`RQS-02` there. Q-D precedes
Q-B in order, because a retraction is itself an event in a state and the notice join is therefore
state-relative (RPS-09). It does not displace Q-A, Q-B or Q-C.

**What must not be put to the owner yet:** which implementation is the experimental subject (PH1A-12).
That question is downstream of Q-A, because the answer determines whether `m09/` is the subject with a
bridge, or `model.py` is the subject with an acquisition model that does not exist.

---

## 12. What follows, and what this review does not decide

**Consequences for the sequence the owner proposed.**

- **1A-2** can proceed as designed, but its brief must carry Q-A/Q-B/Q-C rather than a binary. The
  four-class decomposition (PH1A-01) is unaffected by this review and remains independent.
- **1A-3** (experimental subject) is now demonstrably downstream of Q-A: the harness cannot express the
  relation at all (RP-03), so "bridge the harness to `m09/`" and "add acquisition to the harness" are
  different projects with different evidence value.
- **1B** (design contract) should not fix the relation before Q-B, because a contract that fixed RA-2
  would silently decide RQ-01 by omission.

**Explicitly not decided here.** No model is adopted. No class is retained or retired. No requirement is
added, amended or relaxed. No register entry is proposed. No implementation change is proposed. No probe
is designed. No threshold, category, or pass criterion is invented. The three readings in §7 are recorded
as readings, not ranked.

**Limits of this review.** All behavioural claims are static readings at the hashes in §3; nothing was
executed, so no claim here is a demonstration of runtime behaviour. §8.4's display is quoted, not
rendered. The hybrid in RA-3 is a conceptual sketch, not a design: its three undefined points are stated
as undefined rather than filled. r32's revision changelog is self-divergent (Phase 1A Appendix A-9), so
r32 is cited by rule id and section, never by change history.

---

## 13. Addendum — conversation-state lifecycle

The owner's follow-up identified a lifecycle operation this review did not separate: a **conversation
state transition** — regenerate, next alternative, undo, revert — can make a memory inapplicable
without any change to what anyone asserted or knows. It is investigated as its own subject in
`MEMORY-09-PHASE-1A-1A-CHAT-STATE-LIFECYCLE-2026-11-19.md` (`RPS-`/`RQS-`/`RS-`/`RST-` IDs). It is
not folded into §4–§11 above, because it changes no record's content and no character's knowledge.

What the addendum establishes, in summary, and why this document's brief changes:

1. **No source governs the operation.** r32 and v9 say nothing about regenerate, branch, revert or
   undo; chat deletion is governed in detail (K-R35, K-R40–K-R44), chat state is silent. The 118
   requirements have no state row.
2. **`status` cannot carry applicability.** It is 3-valued and already means retraction/forgetting;
   a fourth value for "inapplicable here" would be the enum collapse the owner ruled out, and a
   reverted state could not restore activity without yet another value.
3. **The one coupling to this review is RPS-09:** a retraction occurs in a state, so the derived
   notice of §7 (R-narrow) must join on `(state, status)`. Q-B's substance is unchanged and no
   reading in §7 is withdrawn; only what a correct answer must be computed from changes.
4. **Howl cannot supply state provenance today.** Its history is one mutable list with positional
   message ids; undo rewrites the log atomically and destroys the previous state's messages;
   regenerate does not; Prev/Next only scrolls despite the documentation. This is stated as a
   condition on any design, not as a defect report.

Q-D is therefore added to §11.3 and placed before Q-B. PH1A-12 remains downstream of Q-A, and is now
also downstream of Q-D.

---

## Appendix A — Traces

**RT-01** `form.py:174` retraction gate → `_retract` → `store.current(owner_id=speaker)`.
**RT-02** `form.py:206–240` `_admit` → `add_event(via="witnessed")`; `conversation.py` forms per speaker only.
**RT-03** `model.py:282–302` `retract()` → `revoke()` (`:303–318`) → `add_event(predicate="retraction:…")`.
**RT-04** `model.py:347–353` `edges()` maps every event, including `via="revoked"`.
**RT-05** `retrieve.py:100–143` `eligible` → `store.current()`; `fixtures.py:189, 218, 242` p48 asked of Alaric.
**RT-06** `form.py:174` narrator excluded from retraction; `_retract` scoped to `owner_id`.
**RT-07** `engine.py:542–559` `retract(record_id, turn)` — no actor parameter.
**RT-08** `engine.py:575–700` `disclose` descriptors; `retrieval.py:46–52` perspective excludes `revoked`.
**RT-09** `engine.py:575–700` disclose checks class and owner, never status.
**RT-10** `engine.py` — no event store; `status` assigned directly at `:337, :370, :539, :558`.
**RT-11** `engine.py` `authoritative` → `_eligible_truth` (`:167`); `accounts_of` (`:1181`); no join.
**RT-12** `records.py:156–158` private records refuse `proposition_id`; docstring gives enumeration as the reason.
**RT-13** `model.py:347–353` + `K-R19` — projection rule violated while event rule is met.
**RT-14** `engine.py:561–574` `_require_human` — human gate exists for authority, absent for `retract`.

## Appendix B — Text locations for the load-bearing quotations

| Quotation | Location |
|---|---|
| "a retraction changes the memory stream only… whether they have *learned of the retraction* is itself an acquisition event" | r32, consequence under K-R13, lines 113–115 |
| "you were told X; it has since been retracted; you still believe it" | r32 §8.4, line 421 |
| "A revoked pair has **no edge**…" | r32 K-R19, line 145 |
| "the disclosure event and the recipient's acquired account persist after any retraction of the source" | v9 rule 49 |
| "Once disclosure occurs, subsequent retraction of the source does not erase the recipient's acquired account" | O-3, T1 ledger |
| "No notice account created" / "Retraction knowledge appears unaccounted" | Conformance map E1c |
| "No automatic knowledge of the retraction" | Conformance map A2b |
| "a private record cannot be attached to a proposition" | `m09/records.py` `__post_init__` |

---

**End of review.** No canonical artifact was modified. The next action requires an owner response to
Q-D (RQS-01) and then to RQ-01 (Q-B), which is the question the architecture choice does not answer.
