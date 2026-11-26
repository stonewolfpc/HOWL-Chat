# MEMORY-09 — Architecture v9

**Status:** **RATIFIED (design baseline) — O-6b.** Ratified as tier 3 under O-6a, subject to tiers 1
and 2. Owner confirmations were recorded in conversation; **no persisted signature exists and none is
implied by this file.** Nothing here is implemented or executed. Ratification establishes design
coherence, **not** demonstrated correctness: of 118 defined tests, **0 implemented, 0 executed**.
Any rule later found to contradict a tier-1 decision is **void**, recorded as a defect in this
document rather than as a divergence.

| Field | Value |
|---|---|
| Record created | **2026-11-16 23:47:21 CST (−0600)** = **2026-11-17 05:47:21 UTC** |
| Authority created | **Design baseline only** (O-6b). No implementation authority — **Gate 2** (`eval/mem1` implementation) remains unanswered and unauthorised |
| Scope | Experimental engine semantics (`eval/mem1` target). Production Howl was **not** inspected |
| Confirmed inputs | D-01, D-02, O-1, O-11, O-8a (all six), O-10, O-10b, O-12, O-13, O-3, O-4, O-5 (see `MEMORY-09-OWNER-DECISIONS-2026-11-17.md`) |
| Open inputs | **Gate 2** (implementation — unanswered), O-2 |
| Conformance permission | **O-7 / Amendment 1** — probes for the 118 requirements only; prospective; no prior-intent claim for plan §10 |
| Tests | **118** defined, **0 executed** (`MEMORY-09-CONFORMANCE-MAP-2026-11-17.md`) |
| Relationship to r32 | **Settled — O-6a.** Layered precedence: tier 1 owner decisions > tier 2 r32 mechanisms > tier 3 this document. Register entries **D-1** (transfer, O-13), **D-2** (D6 vs §10.1 rev 2), **D-3** (disclosure gap-fill), **D-4** (open: destination interface) |
| Registry | No registry row added (not authorised). See §10 |

Claim labels: **[D-01 conf]** · **[D-02 conf]** · **[OC]** owner clarification · **[P]** proposal ·
**[SO]** static observation · **[UA]** unverified assumption.

---

## 1. Responsibility boundary

The system answers different questions in different places. Collapsing them is what produced the
`owner_id` defect recorded in `MEMORY-09-SESSION-CORRECTIONS-2026-11-17.md` §4.

| Question | Owner |
|---|---|
| Whether a character speaks | **Model** |
| What the character chooses to disclose | **Model** |
| What a character believes, concludes, remembers | **Model** (the memory system records the result; it does not determine it) |
| Whether a character changes their belief | **Model** |
| Whether disclosure actually occurred | **Conversation / event layer** |
| What information was communicated | **Conversation / event layer** |
| Which characters were exposed | **Conversation / event layer** |
| What was actually said | **Conversation / event layer** |
| Recording memories and acquisitions | **Memory system** |
| Preserving provenance | **Memory system** |
| Identifying and representing conflicts | **Memory system** |
| Eligibility per retrieval context | **Memory system** |
| Applying authorised conflict-handling rules | **Memory system**, under authorised rules |
| Deciding which character is correct | **Nobody.** Not the model, not the memory system |
| Authoritative promotion | **Separately authorised arrangement** (O-8a) |

Prohibitions: no second AI model; no AI-based memory arbitration; the memory system must never
initiate disclosure, infer disclosure from proximity, or manufacture acquisition.

---

## 2. Four concepts

| Concept | Question it answers |
|---|---|
| **Memory** | what exists, in which domain, with what provenance |
| **Declaration** | what has been deliberately made generally available, by whom, over an immutable set |
| **Account** | what information a character has acquired, and how |
| **Belief** | what a character currently accepts, suspects or considers true |

The withdrawn four-class architecture (`ambient_authored` / `experiential` / `restricted` /
`character_specific` as a stored classification) is **not revived**. Behaviour it was trying to
express is carried by declarations, accounts, beliefs and domains.

---

## 3. Domains and truth-aptness

| Record class | Truth-apt? | Records | May change domain truth? |
|---|---|---|---|
| World truth | **Yes — authoritative** | what is true in the world | it is the authority |
| Scenario truth | **Yes — authoritative within that scenario** | what is true in that scenario | within its domain only |
| Public statement / claim | **No** | that a character said something, and what was said | never |
| Private recollection / interpretation | **No** | what a character remembers or makes of something | never |
| Acquired knowledge (account) | **No** | what information a character holds | never |
| Belief | **No** | what a character currently accepts | never |

`character_public` is **not a truth domain**. It is a visibility marker on a fact whose truth
resides in the world or scenario domain.

---

## 4. Rules 1–27

1. **Responsibility separation** — §1 boundaries hold.
2. **Containment** — every record is contained by exactly one chat. **[D-01 conf]**
3. **Character ownership** — every private record is owned by exactly one character. **[O-1 conf]**
4. **Association** — a private record names its owning character's perspective; association is not containment.
5. **Access** — a private record is retrievable only by its owner, unless a disclosure created a separate record for someone else.
6. **Applicability** — domain plus target says where something applies; never who owns or knows it.
7. **Truth-aptness** — only world and scenario records hold truth; nothing else may.
8. **No adjudication** — the system records claims, recollections, inferences and beliefs without marking them true or false against domain truth.
9. **Acquisition** — knowledge exists only as accounts, with attributable origins and evidence by origin.
10. **Disclosure** — a private memory becomes shared information only when an actual disclosure event establishes it; the model decides whether and what to disclose.
11. **Exposure limits acquisition** — only characters actually exposed acquire; group membership confers nothing.
12. **Provenance without bypass** — recipient records carry a restricted origin descriptor, never access to the source record.
13. **Conflict detection ≠ resolution** — §9.
14. **No truth mutation by non-authoritative records** — only an explicit promotion (§6) changes which record governs.
15. **Belief independence** — belief is not created, replaced or revoked by contradictory information, disclosure, promotion or acquisition, nor inferred from speech or behaviour. **[O-5 conf]**
16. **Promotion is explicit** — never from disclosure, sharing, repetition, persuasion, agreement or conflict detection.
17. **Promotion changes authority, not ownership** — it creates a world/scenario record with declared applicability and provenance; the source is unchanged.
18. **Promotion confers no knowledge** — no accounts are created; content equivalence never implies knowledge.
19. **Promotion compels no belief** — rejection, doubt, misunderstanding or indifference are valid outcomes, not system failures.
20. **History is preserved** — conflicting records remain; a later promotion supersedes the earlier *authoritative* record; nothing is erased.
21. **Authority, not history** — promotion settles which record governs truth queries; it never means the disagreement did not happen. **[O-11 conf]**
22. **Creation is open** — any participant, human or model, may create conversational information.
23. **Recording is system-governed** — the system preserves who contributed what, in which domain, from which perspective, by what route.
24. **Promotion authority is separate** — generation never confers it. A world or scenario may grant a model **standing promotion authority**; the grant must be explicit, domain-scoped, recorded and revocable. Promotion occurs only under the declared arrangement, never as a side effect of generation. **[O-8a-6 conf]**
25. **A statement is a claim** — a character's utterance records that they said it; it never establishes that they experienced what they described.
26. **Kind is declared, not judged** — contribution kind is declared by the model and recorded by the conversation layer; the memory system must not infer it and must not silently reclassify it. A missing declaration is recorded as **`Unspecified claim`**; the system never supplies a kind on the model's behalf. A declared recollection records how the model presented its contribution — never that it is verified truth. **[O-12 conf]**
27. **Creation is not promotion** — a record created by conversation is authoritative only through an explicit promotion or a standing authorised arrangement.
28. **Basis requirement** — every promotion records a **basis: a set of one or more typed entries**, each of kind `record_cited` or `actor_declaration`. Mixed bases are allowed; an **empty set is rejected**. The system validates structure and references only. **[O-8a-2 conf]**
29. **Basis retraction does not undo promotion** — authority stands until a separate authorised operation changes it. A retraction of the cited source is an input to a future promotion or revocation, never an automatic undo. **[O-8a-2 conf]**
30. **Declarations preserve content and author** — an `actor_declaration` entry records **what was declared**, plus the **identity and authority of the declaring actor**. It is a promotion property, not a contribution kind, and does not interact with rule 26. **[O-8a-2 conf]**
31. **Promotion never extends reach** — a promoted world/scenario record remains contained within the chat where it was created; its authority applies there and nowhere else. Cross-chat transfer, sharing and carryover are **not implemented by MEMORY-09**: they belong to the application and require the user's decision. **[O-13 conf]**
32. **Reach is still not knowledge** — even inside the originating chat, authority grants no character knowledge. Knowledge requires declaration coverage or an account. **[D-02 conf]**
33. **Citations must resolve within the chat** — each `record_cited` entry identifies **specific records** that exist, are **retrievable by the promotion operation within the promoting chat**, and comply with existing access rules. A citation to another chat's record cannot resolve, by construction (rule 31). **[O-8a-2 conf]**
34. **No relevance or coverage test** — the authorised actor judges whether the basis is sufficient. The system never assesses persuasiveness, relevance or coverage. **[O-8a-2 conf]**
35. **Authority is per domain** — promotion authority is held separately for each world and each scenario, never globally across the application. A grant in one domain permits nothing in another. **[O-8a-3 conf]**
36. **World scope is an authority domain, not a repository** — two chats using the same world may hold divergent, individually authoritative canon; neither is wrong (rule 31). **[O-13 conf]**
37. **Revocation and supersession are explicit** — both exist as operations within a domain; original records and provenance remain preserved, and supersession creates a new authoritative state without erasing history. **Execution is human-controlled (rule 43).** **[O-8a-4 conf, qualified by O-10b]**
38. **Authority changes never rewrite character state** — revocation or supersession alters no memory, knowledge or belief. **[O-8a-4 conf]**
39. **One authoring authority, two operations** — promotion and Q7 ambient authorship share one regime scoped per world/scenario; ambient authorship does not promote. **[O-8a-5 conf]**
40. **Shared governance is not shared content** — no shared content, no automatic promotion, no cross-chat transfer follows from a shared regime. **[O-8a-5 conf]**
41. **Equal-authority conflict is a state, not an error** — two legitimately promoted, incompatible, authoritative records in one domain with overlapping applicability are both preserved, with provenance and promotion history; neither receives precedence; the conflict is explicitly represented, and retrieval must not present either as settled truth. **[O-10 conf]**
42. **Three operations, three meanings** — **correction** changes a record's content; **supersession** changes which authoritative record governs; **revocation** changes whether a promotion remains authoritative. A correction never silently resolves an authority conflict. **[O-10 conf]**
43. **Authority changes are human-controlled** — no actor may unilaterally revoke or supersede another actor's promotion, **or their own prior promotion**, merely by disagreeing with it. Resolution of an authority conflict requires an explicit human edit. Authority to promote does not confer power to rewrite the shared authority state. **[O-10b conf]**
44. **A human edit changes authority, never history** — the original records, the conflicting accounts and their provenance survive the resolution. **[O-10b conf]**
45. **Transfer is outside MEMORY-09** — no transfer, share or carryover operation is exposed; orchestration belongs to the application, controlled by the user. **[O-13 conf / Finding 1]**
46. **Disclosure needs no authorisation** — a character may disclose what it has acquired or remembers through its own conversational output. The system records the event and its attribution, creates recipient-specific acquired accounts, and does not expose the source record. It grants no permission to speak. **[O-3 conf]**
47. **Narration cannot disclose for a character** — narration may report a disclosure event established by the conversation layer, but must not manufacture one, reproduce a private record, or speak private content on the owner's behalf. **[O-3 conf]**
48. **Disclosure is a claim, not truth** — it creates an attributable claim and a recipient-specific account; the source record is never exposed. Re-disclosure may occur from the discloser's own acquired account. **[O-3 conf]**
49. **Irreversibility is historical, not perpetual access** — the disclosure event and the recipient's acquired account persist after any retraction of the source, and this does not preclude other legitimate lifecycle operations. **[O-3 conf]**
50. **Two questions, two fields** — `via` answers how information entered the character's knowledge; `origin_class` answers whether the account is backed by an acquisition event. They never compete for one field. **[O-4 conf]**
51. **Own accounts create no events** — a character's own account is recorded `origin_class = own_account`, `via = none`, and creates or implies no acquisition event. This does not establish that no historical acquisition occurred, only that none is represented. **[O-4 conf]**
52. **A stated source is content, not history** — a source named inside a recollection is content of that recollection, never a verified acquisition event. The system does not manufacture an event to explain an account. **[O-4 conf]**
53. **Beliefs are declared, not inferred** — a belief exists only where the **model embodying that character** explicitly declares it. Not another character, not ambient narration, not the system. An utterance is not automatically a belief record. **[O-5 conf]**
54. **Narration may report a belief, never install one** — ambient narration may report a belief established by the conversation layer; it cannot invent a character's interior state. **[O-5 conf]**
55. **Belief revision is append-only** — a changed belief creates a new record; the prior belief remains retrievable by its owner. **[O-5 conf]**
56. **A belief is never promoted as truth** — a claim about what a character believes may become authoritative as a **fact about that character**; it never becomes authority over the belief's content. **[O-5 conf]**
57. **Disclosure carries the belief, not its content** — the recipient acquires that the character believes it, never that it is true, and gains no access to the original private record. **[O-5 conf]**

---

## 5. Six separate dimensions

| # | Dimension | Governs | Consequences |
|---|---|---|---|
| 1 | **Chat containment** | lifecycle | exactly one containing chat **[D-01 conf]**; deletion/transfer/restoration are **outside this experiment** |
| 2 | **Character ownership** | exclusive access among characters | only the owner retrieves and may disclose; ownership is never transferred, only copied by disclosure into a new record **[O-1 conf]** |
| 3 | **Association / attribution** | whose perspective or voice | makes a private record someone's; grants neither access nor knowledge |
| 4 | **Access / visibility** | who may retrieve | owner only for private records; declarations and exposure otherwise |
| 5 | **Acquisition** | what information a character holds | accounts with origin, evidence, source, turn |
| 6 | **Applicability** | where something applies | domain + target only |

Containment and character ownership coexist because they quantify over different sets (chats vs
characters), answer different questions (lifecycle vs exclusive access), and have disjoint
consequences. **D-01 is unamended by this distinction.**

---

## 6. Promotion lifecycle

1. **Source** — a character-originated record exists (recollection, claim, disclosure).
2. **Act** — an authorised actor performs an explicit promotion, recorded with actor, authority reference, timestamp with offset, source, target domain, declared applicability, basis, optional rationale.
3. **Result** — a world/scenario record is created, truth-apt, citing the source.
4. **Effect on truth queries** — the promoted record governs within its applicability; conflicting character-level records no longer govern.
5. **Effect on characters** — none automatic: no account, no belief change.
6. **Effect on history** — source, prior authoritative records, and the disagreement all remain.

**Bases:** any attributable basis may be cited, including a disclosure. Disclosure is evidence, not
authority; authority is O-8a.

---

## 7. Disclosure lifecycle

1. **Model** decides the character speaks and what to disclose.
2. **Conversation layer** establishes that disclosure occurred, its content and recipients — recorded as a claim event with spans.
3. **Exposure** follows from that event alone.
4. **Memory system** creates, per exposed character, their own acquired account: source, event/turn, communicated content, origin descriptor.
5. **Unexposed characters acquire nothing.**
6. **Source private record unchanged**, still exclusively owned.
7. **Disclosure establishes neither truth nor belief.**

A character may disclose only what they hold: an account, a private record, or a belief they voice.
Disclosing a belief produces a claim — the recipient acquires *that the character believes it*, not
the belief.

---

## 8. Privacy and provenance

**Recipient may hold:** who communicated it; when and in what context; the content communicated;
and, where disclosed, that it originated in the speaker's own recollection or private account.

**Recipient must not obtain:** the source record, its unrelated contents, its derivation history,
or any pointer resolving to them.

**Mechanism:** the recipient's origin reference is a **descriptor** (speaker, record type, turn),
not a resolvable identifier. The resolvable link exists only on the source side. Each recipient
record has its own identity and acquisition provenance; it is never a view, alias or projection.

---

## 9. Conflict: detection versus resolution

**May:** identify incompatible records; preserve competing accounts; preserve origins, owners,
acquisition history and truth status; decide eligibility per context; apply explicit authorised
rules; abstain.

**Must not:** prefer a belief because it is newer or repeated; promote a private recollection to
world truth; overwrite one character's memory with another's; force belief revision; treat
disclosure as verification; treat acquisition as acceptance; declare a character correct because
information was disclosed to them.

**[SO]** r32 D6 already assigns conflict detection to the application and forbids model arbitration;
its materiality predicate remains open (Q1; ARCH-03 finding F3). This architecture invents no new
conflict authority and assumes no policy; where D6 is silent, the system abstains.

---

## 10. Boundaries, exclusions, and declared defects

**Outside this experiment:** chat lifecycle, deletion, transfer, restoration, world/scenario
switching, reusable asset lifecycle. The active world and scenario are held constant.

**Not included by design:** temporal intervals, materiality, decay/surfacing, asset versioning,
group knowledge, transfer, deletion.

**Declared, not repaired:** `docs/DOCUMENT-REGISTRY.md` claims to cover `docs/` exhaustively. This
file adds one more unregistered file; no registry edit was authorised. The understatement now stands
at four known files (stage contract F18; ARCH-03 review record; R3 recovery closure; this set).
