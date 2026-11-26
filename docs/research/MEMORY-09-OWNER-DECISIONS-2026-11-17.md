# MEMORY-09 — Owner Decision Ledger

**Status:** decision ledger. Confirmations below were recorded **in conversation**; no persisted
signature exists. Nothing here authorises implementation, test execution, or contract edits.

| Field | Value |
|---|---|
| Created | **2026-11-16 23:47:21 CST (−0600)** = **2026-11-17 05:47:21 UTC** |
| Authority | Records owner confirmations made in session; creates no new authority |
| Companion records | `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md`, `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md`, `MEMORY-09-SESSION-CORRECTIONS-2026-11-17.md` |

Status vocabulary: **CONFIRMED** (owner stated in session) · **NARROWLY ANSWERED** · **OPEN** ·
**DEFERRED**.

---

## 1. Confirmed

### D-01 — CONFIRMED (pre-existing, retained verbatim)

> Experiential memories belong to the originating **chat**; `world`/`scenario`/`character_public`/
> `character_private` are **scopes**, not owners. Chat deletion removes memories, not reusable
> assets. No automatic cross-chat carryover.

**Recorded defect:** the approval instrument for D1–D7 is not locatable in this workspace
(ARCH-03 findings A2/A3). Recorded, not re-litigated.

**Rationale attached (non-normative, outside experiment scope):** chat saved → memories saved with
it; chat retained → memories retained; chat deleted → orphaned memories explicitly deleted;
exceptions are explicit keep and explicit transfer. Cross-reference: agrees with K-R39/K-R40;
**differs** in adding a "keep" action where K-R40 says manual transfer is the only preservation
mechanism. Referred to D8, not resolved here.

### D-02 — CONFIRMED (eight behaviours)

Existence is not knowledge; behaviour summary as approved:

1. Separation — existence implies no knowledge; belief implies no truth.
2. Acquisition — knowledge only by acquiring; presence, proximity, relevance, retrieval, storage confer nothing.
3. Deliberate generality — general availability established by an explicit attributable act, never inferred from scope, age, popularity or access.
4. Restriction — withholding from unauthorised characters, including later arrivals.
5. Individual knowledge and belief — including belief contradicting a general fact, without altering the fact.
6. Events are not general knowledge.
7. Late arrival — no prior knowledge except what is explicitly declared general.
8. Distinguishability — existence, who knows, and how each came to know are separable.

### O-1 — CONFIRMED

Character ownership means **exclusive ownership and control of the original private record among
characters**, independent of chat containment. Operative example:

- A owns and retains exclusive access to the original private memory.
- A may disclose from it; B receives a separate acquired record containing only what was disclosed.
- B gains no access to A's original, its unrelated contents, or its derivation history.
- A's original remains unchanged.
- Both records remain contained by their originating chat under D-01.

**D-01 is unamended.** This is an additional ownership distinction, not a replacement.

### O-11 — CONFIRMED

- Promotion creates a separate world/scenario record.
- The source character record remains unchanged in identity, ownership, content and history.
- For authoritative world/scenario truth queries, the promoted record takes precedence within its declared applicability.
- For character memory and recollection queries, the original record remains available to its owner.
- Conflict/history remains represented; promotion changes which record governs authoritative truth, not whether the disagreement occurred.

Operative framing: *what does the authoritative record say* and *what does this character remember*
are two legitimate questions with legitimately different answers.

### O-8b — NARROWLY ANSWERED

Promotion changes **authority**, not ownership, knowledge or belief. Characters may accept, reject,
doubt, misunderstand or ignore promoted information; non-acceptance is not a system failure.

### O-9 — NARROWLY ANSWERED

Disclosure **may** be cited as the **basis** for an authorised promotion. Disclosure never promotes
by itself. Evidence and authority are different things; sufficiency of a basis folds into O-8a.

---

## 2. Open

### O-8a — promotion authority (ALL SIX CONFIRMED — series closed)

Reframed: creation is open to any participant; recording is system-governed; **authority to promote
is separate and explicit.**

**8a-1 is CONFIRMED (owner, this session):** promotion authority is held by the world/scenario
authoring authority, defaulting to the human operator, declared per domain. This is the
default-and-delegation arrangement approved by the owner.

**Correction recorded:** an earlier status report in this session listed 8a-1 as open after that
approval. The ledger error is corrected here and logged in
`MEMORY-09-SESSION-CORRECTIONS-2026-11-17.md` (correction 19).

| # | Question | Recommendation (not decision) |
|---|---|---|
| 8a-1 | Who holds promotion authority? | **CONFIRMED** — world/scenario authoring authority, default the human operator, per domain |
| 8a-2 | What bases are sufficient? | **CONFIRMED** — basis kinds `record_cited` and `actor_declaration`; a basis is a set of ≥1 typed entries (mixed allowed, empty rejected); citations must resolve in-chat; declarations preserve content and author; no relevance or coverage test |
| 8a-3 | Per world/scenario or global? | **CONFIRMED** — per domain: held separately for each world and each scenario, never globally |
| 8a-4 | May the same authority revoke or supersede? | **CONFIRMED, qualified** — explicit revocation or supersession exist within the domain; history and provenance preserved. **Execution is human-controlled (O-10b):** no actor may unilaterally revoke or supersede another's promotion, or their own prior one |
| 8a-5 | Shared with Q7 ambient authorship? | **CONFIRMED** — one shared authoring-authority regime, scoped per world/scenario; the operations stay distinct |
| 8a-6 | May an authorised arrangement delegate promotion to the model (e.g. canonised narration)? | **CONFIRMED** — permitted: a world or scenario may grant a model standing promotion authority, explicit, domain-scoped, recorded and revocable |

**8a-2 CONFIRMED — governing contract (owner, this session):**

- **Basis types:** `record_cited` and `actor_declaration`.
- **Basis structure:** a set of one or more typed entries. Mixed bases are allowed; empty sets are rejected.
- **Record citations:** must identify specific records that exist, are retrievable by the promotion
  operation within the promoting chat, and comply with existing access rules.
- **Actor declarations:** must preserve the actual declared content and the identity and authority of
  the declaring actor.
- **No relevance or coverage test:** the authorised actor judges whether the basis is sufficient. The
  memory system validates structure and references, not persuasiveness.
- **O-13 boundary:** promotion cannot cite records from another chat, nor extend authority beyond the
  originating chat.

**Architectural consequence (stated explicitly):** world scope defines an **authority domain**, not a
shared cross-chat repository.

**8a-3 CONFIRMED (owner, this session):** promotion authority is **per domain** — held separately for
each world and each scenario, never globally across Howl Chat. Consistent with the explicit
delegation boundary established in 8a-1 and 8a-6.

**Correction 21 is resolved by this confirmation** and is marked as such in the corrections record.

**8a-6 CONFIRMED text (owner, this session):**

> "A world or scenario may grant an AI model standing promotion authority. The grant must be explicit,
> scoped to that domain, recorded, and revocable. The model's ordinary generation does not itself
> grant authority; promotion occurs only under the declared arrangement."

**Consequence:** AI-driven canonisation is an **intentional supported capability**, not an accidental
side effect. It is a grant of authority, not a change in what generation means.

**8a-4 CONFIRMED — revocation and supersession (owner, this session):**

The authority holder may explicitly revoke or supersede a prior promotion within their authorised
domain. Safeguards, as stated:

- Revocation and supersession are **explicit operations**.
- Original promotion records and provenance **remain preserved**.
- Supersession establishes a **new authoritative state without erasing history**.
- Revoking or retracting a cited basis **does not automatically revoke** its promotion.
- Authority changes **do not rewrite** character memories, knowledge, or beliefs.

**8a-4 QUALIFICATION (owner, later this session — O-10b):** the operations above remain meaningful,
but **execution is human-controlled**. An actor's authority to promote does **not** grant unilateral
power to rewrite the shared authority state. The confirmed wording above is retained unedited; this
paragraph records the later qualification rather than replacing it.

**Consequence for delegated authority:** a model holding a standing grant (8a-6) may promote, but may
not revoke or supersede — not another actor's promotion and not its own — without an explicit human
edit.

**8a-5 CONFIRMED — shared regime with Q7 ambient authorship (owner, this session):**

One shared authoring-authority regime applies to both promotion and Q7 ambient authorship, scoped per
world/scenario. The **operations remain distinct**: authority to author ambient content does not mean
every contribution is promoted to authoritative truth.

### O-8a — CLOSED: consolidated record

| # | Decision |
|---|---|
| 8a-1 | Authority holder — human operator by default; explicit delegation to authorised actors supported |
| 8a-2 | Promotion basis — typed basis sets, mechanical validation, no relevance or coverage judgement |
| 8a-3 | Authority scope — independently scoped per world and scenario |
| 8a-4 | Revocation and supersession — explicit operations; history and provenance preserved; **execution human-controlled (O-10b)** |
| 8a-5 | Q7 ambient authorship — shared authoring-authority regime, distinct operations |
| 8a-6 | Standing model authority — explicit, scoped, recorded, revocable grants |

**Invariant surviving the whole series:** shared authority governance does **not** mean shared content,
automatic promotion, or cross-chat transfer. Authority is shared; content, promotion and reach are not.

### O-13 — cross-chat reach of promoted world authority (CONFIRMED — originating chat only)

**Decision (owner, this session):** a promoted world/scenario record **remains contained within the
chat where it was created.** Its authority applies within that chat. Sharing it with another chat is
**not** an automatic consequence of promotion — and not a consequence of promotion at all.

**Responsibility boundary established by this decision:**

| Layer | Responsibility |
|---|---|
| MEMORY-09 | Records, provenance, promotion and authority **within the originating chat** |
| Howl Chat application | Controls any transfer between chats |
| User | Decides whether a memory is transferred |
| Promotion | Never transfers, copies or exposes a memory to another chat by itself |

**Consequence for implementation:** the memory system must not implement cross-chat transfer,
sharing or carryover. That capability is out of scope **by decision**, not by omission.

**Question as originally posed (owner, this session):** whether a promoted record's authority applies
only within its originating chat, or to every chat using the same world. **Answered: originating chat
only.** The owner declined automatic cross-chat availability on the ground that world-level authority
does not by itself establish availability elsewhere, and that turning world scope into a de facto
transfer mechanism would break D-01.

**Related settled points:** basis kinds (8a-2) are confirmed and independent of this question; reach
is not knowledge (rules 31–32).

### O-14 — retraction awareness and character acquisition (CONFIRMED — Q-B, Branch 2)

**Decision (owner, 2026-11-19). Verbatim:**

> Characters do not automatically become aware of retractions because application status changes.
>
> Characters retain what they acquired. They learn about a retraction through legitimate acquisition:
> witnessing, disclosure, reading, inference, verification, or another valid route.
>
> MEMORY-09 records and handles retraction independently of character awareness.

**Answer given: Q-B = Branch 2.** The question was posed in Phase 1A-1B, prepared by 1A-1C and
narrowed by 1A-1D. The owner's answer settles it; the open status recorded at 1A-1B §3.2 is closed.

**What this decides.**

| | |
|---|---|
| **Retraction handling** | MEMORY-09 records status, reason and history, and stops presenting the claim as current truth |
| **Retraction awareness** | never produced by MEMORY-09 from status; only by an acquisition event |
| **Application state / human inspection** | unaffected — the application may hold the retraction and show it to the operator |
| **Character acquisition routes** | `witnessed`, `told`, `read`, `inferred`, `verified`, `narrator_grant` (K-R11), or another valid route |

**"Acquisition event" is read broadly.** The routes above are K-R11's event kinds, not a disclosure
requirement. An inference or a verification is an acquisition event **without being a disclosure
event**: the character's model concludes or establishes it, and the event records that. What the
system may never do is derive an acquisition from application status and present it as the
character's — that is the move this decision forbids.

**What this does not change.** A2b, E1c and D2 need **no amendment**: the decision agrees with them.
O-1, O-3, O-5 and H10 stand. No relation (RA-1, RA-2, RA-3) is adopted, and no content identity is
introduced. Q-A, Q-C and Q-D remain open. Inference remains a **model act recorded as an event**, never
a system derivation from status.

**Consequential reconciliation, performed under this decision.** r32 §8.4's display sentence is
classified as illustrative and non-binding; r32 is revised to **r33** for that classification alone,
with no normative rule changed. K-R7 and the K-R13 consequence are untouched. Recording the
classification is an explicit act required by this decision, not an inference from it.

**Carried gaps, unchanged by this decision.** No `inspection` context exists in either implementation
(1A-1C P-5), so operator inspection remains **unbuilt**. The disclosure path is implemented but
**uncalled** (1A-1C P-7). The v9 harness's active retraction-event record and m09's status-only
retraction remain divergent representations (V-5), and that record's semantic status is
**Unestablished**. Under O-14 these gaps matter more, because acquisition is now the only route to
awareness.

**Not authorised:** implementation, testing, probe execution, corpus work, or any change to the 118
requirements. This is a decision about behaviour, not a build order.

### O-12 — contribution kind (CONFIRMED)

**Decision (owner, this session):** contribution kind is **declared by the model and recorded by the
conversation layer**. The memory system must not infer it and must not silently reclassify it.

**Missing declaration is not invented.** A model output that supplies no kind is recorded as
**`Unspecified claim`**. The system never supplies a kind on the model's behalf.

| Model output | Recorded kind |
|---|---|
| Explicitly declares recollection | Recollection |
| Explicitly declares hypothesis | Hypothesis |
| Explicitly declares jest | Jest |
| No declaration supplied | **Unspecified claim** |

**A declared recollection is not independently verified truth.** It records how the model presented
its contribution, with provenance — never that the contribution is correct.

**Rationale (recorded as given):** the AI supplies experience and personality; the memory system
preserves what happened without pretending to understand the AI better than the AI does.

### O-10 — conflict handling (CONFIRMED)

**D6 is retained for its original purpose:** deciding whether a conflict is **material** to a
character's current turn. Its guards stand — the system detects and never arbitrates, unequal
authority is not a conflict, and only an authorised operation clears a conflict group.

**Added: the equal-authority conflict state.** Two records are authoritative, both were promoted by
legitimately authorised actors, both apply to the same domain, their applicability overlaps, and
their content is incompatible. Neither automatically wins. The state is **unresolved**, and:

- both authoritative records remain preserved;
- both retain their provenance and promotion history;
- neither receives automatic precedence;
- the conflict is explicitly represented;
- retrieval must not silently select one as settled truth;
- only an authorised supersession or revocation changes the authority state.

**Three operations, three meanings:**

| Operation | What it changes |
|---|---|
| **Correction** | The content of a record |
| **Supersession** | Which authoritative record governs |
| **Revocation** | Whether a promotion remains authoritative |

A correction to one record **must not** silently resolve an authority conflict. Supersession and
revocation must be explicit.

### O-10b — who may change another actor's authority (CONFIRMED)

**Decision (owner, this session):** a single actor **does not** remove a record for everyone. **A human
edit is required** — including when the actor disagrees with their own earlier promotion.

**Rationale as given:** the same principle as human interaction, which is the reason this system
exists. Even when it produces confusion or conflict, it is not realistic to let one actor's response
override another's memory.

**Stated properties:**

- One actor cannot unilaterally remove or supersede another actor's authoritative record.
- An actor cannot automatically erase or replace even their **own** previously promoted record merely
  because they now disagree with it.
- Conflicting authoritative records remain represented.
- Resolution requires an **explicit human edit**.
- The edit changes the authoritative state; it does not erase the historical record or the fact that
  disagreement occurred.

**This is a safeguard, not a limitation to be engineered around:** it protects the ability of
characters to disagree, including about their own memories. The system must not manufacture
consensus by letting one participant erase another's account.

**Reconciliation with 8a-4:** the explicit revocation/supersession operations remain; execution is now
human-controlled. 8a-4's confirmed wording is retained unedited — see the qualification above.

### Findings carried from the O-10 review

**Finding 1 — transfer rules (scope correction, not a reopening of O-13).** K-R45–K-R47 in r32 specify
transfer as a MEMORY-09 responsibility, contradicting O-13. Resolution, as recommended and agreed:

- Transfer remains a **Howl Chat application** responsibility.
- MEMORY-09 exposes **no** transfer operation.
- Howl may orchestrate an authorised transfer using its own lifecycle and persistence mechanisms.
- MEMORY-09 must not silently implement transfer, sharing or carryover.

The contract text has **not** been altered; the discrepancy is reported and its ownership recorded
here.

**Finding 2 — D6 versus §10.1 revision 2.** Three distinct things exist: D6's established materiality
rule, §10.1 revision 1 (five conditions), and §10.1 revision 2 (three conditions plus E1–E3).
Revision 2 remains **[P] proposed** until its conditions and exposure obligations are reviewed against
D6. **No approval by implication**, and no silent substitution.

**Finding 3 — r32 has no concept of disclosure.** The term does not occur in r32 at all
(0 occurrences of "disclos", case-insensitive). What the contract has instead is **`told` as an
evidence kind** — one value in the `via` precedence chain (`witnessed > verified > read > told >
inferred > narrator_grant`), recorded as an acquisition event with source attribution and never
promoted by repetition (G2). r32 treats telling as **evidence**, not as a governed act requiring
authority. The v9 disclosure lifecycle therefore has **no contract basis**: it is new architecture,
and the contract's silence is not endorsement. Reported, not reconciled.

### O-2 — lifecycle consequences (OPEN)

If any lifecycle consequence is intended for character ownership (survival or transfer on chat
deletion), D-01 §2–§3 must be reopened explicitly. Out of experiment scope.

### O-3 — disclosure authorisation (CONFIRMED)

**The original framing is withdrawn.** Gate-based disclosure ("sharing requires authority") is
withdrawn: it would have the memory system decide whether a character may speak its own mind.

**Decision (owner, this session), verbatim:**

> Disclosure requires no MEMORY-09 authorisation. A character may disclose information it has
> acquired or remembers through its own conversational output. The system records the disclosure and
> its attribution, creates recipient-specific acquired accounts, and does not expose the source record.
>
> Narration cannot independently disclose, reproduce, or expose a character's private record on that
> character's behalf. It may report an established disclosure event without gaining access to
> unrelated private content.
>
> Disclosure creates an attributable claim, not authoritative truth. Re-disclosure may occur from the
> disclosing character's own acquired account. Once disclosure occurs, subsequent retraction of the
> source does not erase the recipient's acquired account.

**Narration distinctions, as tightened by the owner:**

| Narration may | Narration may not |
|---|---|
| Report a disclosure event established by the conversation layer | Manufacture a disclosure that did not occur |
| Describe what is publicly observable | Reproduce or expose a character's private record |
| — | Speak private content on the owner's behalf |

**Irreversibility, clarified by the owner:** irreversibility applies to the **historical disclosure
event and the recipient's acquired account**. It does **not** mean the recipient retains access
forever regardless of other legitimate lifecycle operations.

**Effect:** character agency is protected without turning private memory into a system-enforced gag
order.

### O-4 — origin of the character's own account (CONFIRMED)

**Decision (owner, this session), verbatim:**

> A character's own account is a distinct origin from witnessed acquisition. It is recorded with
> `origin_class = own_account` and `via = none`; no acquisition event is created or implied by that
> record. This does not establish that no historical acquisition occurred; it establishes only that
> no acquisition event is **represented** by the record. `via` remains a fixed order over evidence
> kinds and is never a statement about the character's memory of how they learned something. A source
> the character mentions within a recollection is content of that recollection, not a verified
> acquisition event. The system records such accounts without endorsing them and without requiring
> the model to supply provenance.

**Two separate questions, two separate fields:**

| Question | Field |
|---|---|
| How did this information enter the character's knowledge? | `via` |
| Is this account backed by an acquisition event, or is it the character's own account? | `origin_class` |

They must not compete for the same field.

| Character's account | `origin_class` | `via` |
|---|---|---|
| "I witnessed the dragon attack." | `event_acquired` | `witnessed` |
| "Bob told me about the attack." | `event_acquired` | `told` |
| "I remember the dragon attacking." | `own_account` | `none` |

The third establishes what the character **presents** as their own recollection — not that they
witnessed the attack.

**Qualification added by the owner:** `own_account` does **not** mean no historical acquisition ever
occurred. It means **no acquisition event is established by this record.** The character may genuinely
have learned the information years ago; unless that acquisition is independently represented, the
memory system cannot invent it.

**Consequence for a stated source:** "My mother told me the king was murdered when I was young"
remains an account of the character's recollection. The system must **not** manufacture an acquisition
event involving the mother merely because the model supplied that explanation.

**Why `recollected` is not a `via` kind:** r32's `via` is a **fixed total order** over evidence kinds
(K-R26: `witnessed > verified > read > told > inferred > narrator_grant`). Inserting a non-evidence
kind would force a contract amendment no one is authorised to make, and would ask a malformed
question — whether a recollection outranks an inference. K-R26 is untouched.

**Name separation:** `Recollection` remains the O-12 contribution kind (what the model declares it is
doing); `own_account` is the origin class (what the record is). One word no longer carries two
meanings.

### O-5 — belief independence (CONFIRMED — ratified as an architecture constraint)

**Ratified by the owner, this session.** It remains an **architecture constraint, not an amendment to
D-02**. No prior decision is reopened and no implementation is authorised by it.

**Central principle:** the model owns the character's mental state. MEMORY-09 records it; MEMORY-09
does not manufacture it. **A character's mind is not a database field for the system to edit.**

**Confirmed elements:**

1. **Declaration-only creation.** An utterance such as "I think the king is dead" is **not** a belief
   record. The model must explicitly declare the character's belief. The memory system does not
   interpret dialogue as psychology.
2. **Whose declaration counts (owner's clarification).** Only the **model embodying that character**.
   Not another character, not an ambient narrator, not MEMORY-09 itself.
3. **Narration may report, never install.** Ambient narration may report a belief already established
   by the conversation layer; it cannot invent a character's interior state.
4. **Beliefs are never promoted as world truth.** A claim that "Sable believes P" may become
   authoritative as a **fact about Sable**, provided the promotion has authority. That does **not**
   establish P. Even an authoritative record about someone's belief is never authority over the
   **content** of that belief.
5. **Append-only revision.** A changed belief creates a new record; the previous belief remains
   retrievable. The character's history is preserved, not rewritten.
6. **Disclosure carries the belief, not its content.** The recipient learns that Sable believes P —
   not that P is true — and gains no access to Sable's original private record.
7. **Acquisition compels no psychological response.** Rejection, doubt, misunderstanding and
   indifference are valid outcomes, never system failures.

**Rule 15 of v9 is ratified by this decision and loses its `[P]` marker.**

### O-6a — relationship between v9 and r32 (CONFIRMED)

**Decision (owner, this session): layered precedence, with an explicit divergence register.**
This establishes the relationship between the architecture and the contract. It does **not** ratify
v9 (that is O-6 part 2), does not amend r32, does not resolve any divergence's substance, and does
not authorise implementation.

**Tier 1 — Owner-confirmed decisions.** D-01, D-02, r32's D1–D7, and every confirmed `O-*` decision.
Highest authority. A rule in either document that contradicts a tier-1 decision is void.

**Tier 2 — r32's normative mechanism rules** (90 `K-R` rules, 18 invariants). Governs existing
mechanisms wherever those mechanisms remain within MEMORY-09's scope and do not conflict with tier 1.
**Tier 2 is not permission to override tier 1.**

**Tier 3 — v9's architecture rules** (57 rules). Governs newly defined architecture semantics where
tiers 1 and 2 do not already govern. **Tier 3 is not permission to rewrite tier 2.**

**Three cases:**

1. **v9 contradicts a tier-1 decision → void.** Not a divergence; a defect in v9.
2. **v9 contradicts an applicable r32 mechanism → divergence.** Recorded in the register. Neither
   rule is implemented until explicitly decided.
3. **r32 is silent, v9 speaks → gap-fill.** v9 governs, provided the new rule contradicts no confirmed
   owner decision and no applicable r32 mechanism. Silence is not contradiction. A gap-fill must be
   identified as **new architecture**, never retroactively attributed to r32.

**The register separates three things that must not be conflated:** resolved precedence (which
document governs) · unresolved substance (what the conflicting rules actually require) ·
implementation status (whether either rule is authorised for implementation). Resolving precedence
does **not** resolve the underlying conflict and does **not** authorise implementation.

| # | Affected rules | Controlling decision | Resolved precedence | Unresolved substance | Implementation status |
|---|---|---|---|---|---|
| **D-1** | r32 **K-R45, K-R46, K-R47** | **O-13** | Tier 1 governs: their **transfer responsibilities are outside MEMORY-09's authority under O-13** | How cross-chat transfer is actually performed, if at all | **Not authorised.** r32 unedited; no transfer capability in MEMORY-09 |
| **D-2** | r32 **D6** vs **§10.1 revision 2** | none — internal to r32 | Not a v9/r32 divergence: D6 is the established rule, revision 2 is `[P]` proposed | Which materiality predicate applies, and whether rev 2's E1–E3 are adopted | **Not authorised.** No substitution, no approval by implication |
| **D-3** | v9 disclosure rules **46–49** | **O-3** | Tier 3 governs: r32 is silent (0 occurrences of "disclos") | None at present — gap, not conflict | **Not authorised.** Identified as new architecture, not attributed to r32 |
| **D-4 — CLOSED (Option B; see the D-4 decision)** | Application-initiated cross-chat transfer and MEMORY-09 destination interface | None | O-13 excludes transfer responsibility from MEMORY-09 | Whether MEMORY-09 accepts a destination-record write initiated by the application, or is entirely absent from the transfer path | **Not authorised.** No interface or write capability implied. *Closed 2026-11-17 by the D-4 decision — Option B. Precedence, substance and status cells preserved as recorded at O-6a.* |

**D-4 remains an unanswered interface-boundary question.** It is not an implied permission for
MEMORY-09 to accept destination writes.

**Scope of this confirmation:** the document relationship only. It does not ratify v9, amend r32,
resolve D-4, or authorise implementation.

### O-6b — ratification of architecture v9 and the conformance map (CONFIRMED)

**Decision (owner, this session): v9 and its 118-test conformance map are ratified as the design
baseline for MEMORY-09's experimental engine, subject to tiers 1 and 2 of O-6a.**

**What is ratified:**

- **v9's 57 rules**, as tier 3 — newly defined architecture semantics, governing only where tiers 1 and
  2 do not already govern.
- **The 118-test conformance map**, as **defined requirements**. Ratification is of the
  *requirements*, not of any result: **0 tests implemented, 0 executed, 0 passing.**
- **Scope:** experimental engine semantics (`eval/mem1` target). Production Howl was not inspected
  and is not covered.

**What this does not do:**

- It does **not** amend r32. r32 is byte-identical to `e88ef2d1…`; its transfer rules remain outside
  MEMORY-09's authority under O-13 (D-1).
  *(Qualification recorded 2026-11-19: r32 was subsequently revised to **r33** under **O-14**, which
  classifies §8.4 as illustrative and changes no rule. The sentence above states the position **at
  ratification** and is preserved unamended.)*
- It does **not** resolve **D-1** (how transfer is performed), **D-2** (D6 vs §10.1 revision 2 —
  revision 2 stays `[P]`), or **D-4** (whether MEMORY-09 accepts an application-initiated destination
  write). **D-4 stays open.** *(Qualification recorded 2026-11-17: this sentence states the position **at ratification**. D-4 was subsequently closed by the D-4 decision — Option B, MEMORY-09 absent from the transfer path. The original wording is preserved unamended.)*
- It does **not** authorise implementation of any rule or any test. **O-7 / R-5a remains unresolved**,
  so conformance work still has no authorised channel.
- It does **not** add a registry row, inspect production code, or create any persisted signature. Like
  every prior confirmation, it is recorded in conversation and labelled as such.

**Ratification is of design coherence, not of correctness.** No claim is made that the architecture,
if implemented, produces the behaviours its tests describe — that is what the tests are for, and none
has been run.

**Owner's preserved subtlety:** ratification does **not** mean every future interpretation of these
rules is automatically correct. The **tier-1 contradiction safeguard remains in force** — any ratified
rule later found to contradict a tier-1 decision is void, recorded as a defect in v9 rather than as a
divergence.

**Ordering rationale (owner):** a design can be complete and ratified while implementation remains
blocked. Keeping those statuses separate is the discipline this workstream has been building.

**Standing state after this decision:** a ratified, **unimplementable** baseline — frozen design,
awaiting an authorised channel.

v9 and the 87-test map are proposals awaiting owner approval. Approval would not authorise
implementation or execution.

### O-7 — conformance-probe authorisation (CONFIRMED — Option B)

**Provenance correction (#24):** O-7 was originally titled "O-7 — R-5a". `R-5a` is **not a
source-defined identifier** — it appears only in four internally authored records and in no source
document. The authority for the restriction is **plan v3 §10, bullet 2**: *"No new architecture,
rule, field or probe in the contract while this plan executes."* Historical occurrences of `R-5a`
are preserved as historical references and are never independent authority. The phrase **"in the
contract"** was **not** written by the owner; it arose in collaborative assistant/DeepSeek work, and
its authorship and original intent are **not** attributed to the owner.

**Decision (owner, 2026-11-17): Option B — a numbered, dated, hashed Experiment Plan Amendment 1**,
using the preregistration amendment-register convention. It **prospectively permits only** conformance
probes corresponding to the **118 requirements ratified under O-6b**. Plan v3 is **not modified**; its
bytes and approved identity (`21,616 B`, `6c4078bd…`) are preserved.

**Exception wording (verbatim, as approved):**

> Notwithstanding plan §10, bullet 2, conformance probes that operationalize only the 118
> requirements already ratified under O-6b are permitted.

**Instrument:** `docs/research/MEMORY-09-EXPERIMENT-PLAN-AMENDMENT-1-CONFORMANCE-2026-11-17.md` —
amendment **1**, dated **2026-11-17**, SHA-256 `6afd2a0d8c19a5cd4be5d47da693192222f55842c6dacbcdb14901fdc4723656`, registered here. Per the freeze
record's convention, **an unregistered amendment does not bind**.

**Bounds of this decision:** prospective only · establishes **no** claim about §10's original intent ·
attaches **no** authorship of "in the contract" to the owner · does **not** authorise implementation
(Gate 2) · does **not** alter r32 (`e88ef2d1…`) · does **not** expand the 118 requirements.

### Gate 2 — implementation of `eval/mem1` (AUTHORISED — implementation only; testing **not** authorised)

**A separate board item, approved as an item and not as an authorisation.**

Gate 1 — O-7 plus Amendment 1 — permits narrowly scoped conformance work only. **Implementation of
any `eval/mem1` code remains unauthorised** until the owner makes a separate decision. The G0 receipt
says the same independently: *"A G0 pass authorizes nothing beyond scope: implementation still requires
separate authorization."*

Nothing in Gate 1 implies Gate 2. No decision exists for Gate 2; it has never been asked. Until it
is answered, the 118 requirements can be specified as probes but not executed against anything.

*(The paragraph above records the position **before** the Gate 2 decision; preserved unamended.)*

**Decision (owner, 2026-11-17): Gate 2 is authorised — implementation only.** The sequence is
**implementation first, architecture review second, testing not at all** until the owner explicitly
approves the next stage.

**Implementation must conform to:** the O-6a precedence hierarchy · all **57** ratified architecture
rules · the **118** O-6b conformance requirements as the governing requirements baseline · every
other confirmed decision, including **D-4** and **O-13**.

**Restrictions (exact):**

- `eval/mem1` only.
- No production Howl integration.
- No changes to r32, experiment plan v3, Amendment 1, or ratified decision records.
- No new architecture, rules, fields, requirements or probes.
- No MEMORY-09 transfer, destination-write, sharing or carryover interface.
- Do not resolve open questions through implementation assumptions.
- **Do not run conformance probes or other tests under this authorisation.**

**Mandatory review before any testing** — a separate architecture-conformance review that:

1. maps each of the 57 rules to implementation locations, or reports no evidence;
2. maps all 118 requirements to evidence as implemented, partial, absent, contradictory or not
   assessable;
3. identifies every implementation assumption and its source of authority;
4. identifies deviations, omissions, contradictions and unresolved questions;
5. confirms no prohibited scope or artifact changed;
6. distinguishes source inspection from executed verification.

**Testing remains separately unauthorised.** No probe, selftest or suite may run until the owner
reviews the implementation and explicitly approves the next stage.

**Review delivered:** `docs/research/MEMORY-09-IMPLEMENTATION-REVIEW-R1-2026-11-17.md`.
**Implementation:** `eval/mem1/m09/` (886 lines, 6 modules) — **0 tests executed.**

### D-4 — destination-write boundary (CONFIRMED — Option B, scope determination)

**Decision (owner, 2026-11-17): MEMORY-09 is absent from the transfer path.** No destination-write
interface exists. No transfer-specific memory-creation path exists. This is a **scope determination
for the MEMORY-09 experiment**, **not** a declaration that the Howl application has no transfer
responsibility.

**Reasoning, three layers:** (1) K-R1 / Q2 — edge creation depends on explicit acquisition events, and
Q2 fixes K-R1 as the only creation path; (2) **O-13** and **v9 rule 45** — transfer responsibility is
outside MEMORY-09; (3) **O-6b** ratified requirements **I2** and **K27** — transferable ownership and
transfer, share or carryover operations are each excluded. Option A would introduce a **second
creation route** outside the ratified acquisition mechanism, contradicting the approved scope.

**Recorded alongside this decision:**

- No destination-write interface.
- No transfer-specific memory creation path.
- No implied application-level transfer design.
- **K-R46, K-R47, K-R48** enforcement identified as a **known coverage gap**.
- **No new conformance requirements or probes** — the set stays at **118**.
- **No implementation authorisation** — Gate 2 untouched by this decision.

**Owner's wording, verbatim:**

> These guarantees are not enforceable by MEMORY-09 under the ratified experiment scope. Their
> implementation and enforcement within Howl Chat remain outside this decision and unresolved.

| Contract guarantee | Consequence under B |
|---|---|
| **K-R46** — no events or edges copied | No MEMORY-09 enforcement of this guarantee during transfer |
| **K-R47** — audit chain survives source deletion | No MEMORY-09 enforcement of this guarantee during transfer |
| **K-R48** — additivity | No MEMORY-09 enforcement of this guarantee during transfer |

**Untestable within the ratified set:** the 118 requirements cannot test whether Howl Chat fulfils
these guarantees. That would require separate **application-level** requirements and tests, which
O-6b closed. These are not defects in MEMORY-09; they are guarantees whose enforcement falls outside
its authorised scope.

**Boundary discipline:** closing the boundary is not solving the application-level problem. Transfer
remains an unresolved Howl Chat concern, and nothing here designs it.

---

## 3. What is not authorised by anything in this file

Implementation. Source or test edits. Test, probe or benchmark execution. Contract amendment.
Registry edits. Persistence of signatures (none exist). Promotion of any recommendation to a
decision. Inference of approval from implementation, tests, prior assistant statements, or contract
wording.
