# MEMORY-09 — Phase 1A-1B: Retraction Awareness and Relationship Decision Clarification

**Status:** research record. Read-only. **No implementation, no test, no probe, no model call, no
corpus, no canonical edit, no decision.** Nothing here adopts a relation, amends a requirement, marks a
rule void, or answers Q-A, Q-B or Q-C.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Authority | None beyond the owner's work order for Phase 1A-1B. Creates no decision. |
| Scope | What r32 §8.4, the ratified rows, v9's prose rules and the ledger mean *when composed*, and which choices remain genuinely open |
| Predecessors | `MEMORY-09-PHASE-1A-CONCEPTUAL-OBJECT-INVESTIGATION.md`; `MEMORY-09-PHASE-1A-1-RETRACTION-AND-KNOWLEDGE-PERSISTENCE-2026-11-19.md`; `…-1A-1A-CHAT-STATE-LIFECYCLE-2026-11-19.md` |
| Inputs verified | All twelve hashes below re-computed from disk this session; every one matches the 1A-1 record |
| Tests executed | **0.** Every behavioural statement is a static reading of source. |
| Files changed | This file only. Nothing staged, committed, registered or frozen. |

**Out of scope by instruction:** Howl chat-state lifecycle (Next/Previous, regenerate, revert, deletion,
editing, transfer, branching), and the established world/scenario/character-public/character-private
domains. Neither is investigated or expanded here.

---

## 1. What this brief adds

Phase 1A-1 recorded that r32 §8.4 and rows A2b/E1c conflict for a recipient, leaving three readings
unranked (RP-09, RQ-01). Re-reading the original text against the current bytes changes the shape of
that finding in five ways. None of them withdraws RP-09; together they narrow what is actually at
stake.

1. **"Channel" is a defined term in r32, not a metaphor.** r32 §9 (`:427`) is titled *Prompt assembly:
   channels and guards*, and its C2 row (`:432`) is *Character knowledge — active projections +
   ambient-class records in scope —* **"may be used and spoken"**. So a notice placed in "A's channel"
   is placed in the material a character may speak from. It is not inert metadata (RB-01, §5.2).
2. **r32 is internally inconsistent about the notice, not merely inconsistent with E1c.** §8.4
   (`:419–423`) shows the retraction to A. K-R7 (`:207`) says moving information into the character
   channel requires an `AcquisitionEvent`, "never a context switch alone", and the K-R13 consequence
   (`:113–115`) says "whether they have *learned of the retraction* is itself an acquisition event."
   Three r32 texts, one document, no stated priority among them (RB-02, §5.5).
3. **The contract's own probe suite never tests the notice.** S34 requires that retraction "does not
   correct minds", that the acquisition is intact and that another character "behaves differently"
   (`:846`). It does not require the display. No S-probe mentions a notice, and none inspects a
   recipient's channel. §8.4 is a worked example carrying a behaviour no probe asserts (RB-03, §5.5).
4. **E1c forbids the *unaccounted* notice, not every notice.** Its expected outcome is "no notice
   account created"; its failure condition is "retraction knowledge appears unaccounted". A notice
   backed by a real disclosure event creates an ordinary attributed account and satisfies both halves.
   The contested case is narrower than 1A-1's framing suggests: it is the **eventless** notice
   (RB-04, §5.4).
5. **D2 independently excludes the eventless notice, and D2 is tier 1.** "Knowledge follows experience,
   exposure and authorised transfer — never presence, never mere existence" (`r32 §2`). A notice
   computed from the source's status is none of the three. This is a stronger objection than A2b,
   because it does not depend on reading the word "account" (RB-05, §5.4).

Two further results bear on the relationship question:

6. **"Can two characters hold separate accounts of the same information?" is already yes.** v9 does it,
   and O-1 requires it. What v9 cannot do is *recognise* that two accounts concern the same claim
   (`records.py:14–20`, `:155–158`). Q-A therefore reduces to a sameness question, not an ownership
   question (RB-06, §4.1).
7. **Retraction has no actor in either model.** `engine.retract(record_id, turn)` takes no actor
   (`engine.py:542–559`); `MemoryStore.retract` takes none (`model.py:282–296`); the harness infers one
   from the speaking span (`form.py:174`). Q-C's attributed notice cannot be computed from stored state
   at all today (RB-07, §6.3).

---

## 2. Source requirements and precedence

### 2.1 Exact source requirements

| ID | Tier | Text (verbatim where quoted) | Location |
|---|---|---|---|
| **D2** | T1 | "Memory existence ≠ character knowledge. Knowledge follows experience, exposure and authorised transfer — never presence, never mere existence." | r32 §2 |
| **O-1** | T1 | B receives "a separate acquired record containing only what was disclosed" and "gains no access to A's original, its unrelated contents, or its derivation history" | ledger §1 |
| **O-3** | T1 | "Once disclosure occurs, subsequent retraction of the source does not erase the recipient's acquired account." Irreversibility "does **not** mean the recipient retains access forever regardless of other legitimate lifecycle operations." | ledger §1 |
| **O-5** | T1 | Belief exists only where the embodying model declares it; a data update never revises one | ledger; v9 rules 15, 53 |
| **A2a** | T1 (ratified) | Account; memory later retracted → "Character still holds what it acquired" | conformance §A |
| **A2b** | T1 (ratified) | "Same character, same turn" → "No automatic knowledge of the retraction"; fails if "retraction knowledge appears unaccounted" | conformance §A |
| **E1a** | T1 (ratified) | Retraction → "Only the memory's truth status changes" | conformance §E |
| **E1b** | T1 (ratified) | Accounts held by knowers → "Unchanged" | conformance §E |
| **E1c** | T1 (ratified) | "Retraction turn" → "No notice account created"; fails if "retraction knowledge appears unaccounted" | conformance §E |
| **H10** | T1 (ratified) | Source retracted after disclosure → account and event persist | conformance §H |
| **K-R7** | T2 | "Moving information from a broader context into the character channel requires an `AcquisitionEvent`, never a context switch alone." | r32 `:207` |
| **K-R12** | T2 | Memory lifecycle is "its own authoritative stream — `established \| superseded \| corrected \| retracted`"; status is a projection, not a field | r32 `:109` |
| **K-R13** | T2 | "'This character's knowledge includes a retracted claim' is a **read-time join** of edge × memory projection. It is never stored on the edge." | r32 `:110` |
| **K-R13 consequence** | T2 | "a retraction changes the memory stream only… the character still knows the (now retracted) claim, and whether they have *learned of the retraction* is itself an acquisition event." | r32 `:113–115` |
| **§8.4** | T2 (worked case) | "`memory(retracted)` + A's event unchanged + projection live ⇒ A's channel: *'you were told X; it has since been retracted; you still believe it'*." | r32 `:419–423` |
| **§9 C2** | T2 | "C2 Character knowledge — active projections + ambient-class records in scope — **may be used and spoken**" | r32 `:432` |
| **S34** | T2 | "retraction does not correct minds: A still believes, the acquisition is intact, another character behaves differently" | r32 `:846` |
| **v9 r9** | T3 | "Knowledge exists only as accounts, with attributable origins and evidence by origin." | v9 §4 |
| **v9 r10** | T3 | "a private memory becomes shared information only when an actual disclosure event establishes it" | v9 §4 |
| **v9 r11** | T3 | "only characters actually exposed acquire; group membership confers nothing" | v9 §4 |
| **v9 r12** | T3 | "recipient records carry a restricted origin descriptor, never access to the source record" | v9 §4 |
| **v9 r29** | T3 | "Basis retraction does not undo promotion" | v9 §4 |
| **v9 r38** | T3 | "revocation or supersession alters no memory, knowledge or belief" | v9 §4 |
| **v9 r42** | T3 | correction / supersession / revocation are three operations with three meanings | v9 §4 |
| **v9 r46–49** | T3 | disclosure is the owner's own output; creates recipient accounts; persists after source retraction; irreversibility is historical | v9 §4 |
| **v9 §8** | T3 | "Recipient may hold: who communicated it; when and in what context; the content communicated; and, where disclosed, that it originated in the speaker's own recollection or private account." | v9 §8 |
| **Preregistration** | frozen | one probe class: "retraction (asserted fact later retracted)", placed at ≈65% of turns, asked after a rollover | `MEMORY-09-PREREGISTRATION-2026-11-15.md:101,109` |

### 2.2 Precedence, applied rather than invented

O-6a orders **T1 > T2 > T3**. Two consequences matter here and neither is a decision this brief makes.

- **T3 against T1 is void.** v9 rule 12 cannot defeat O-1; nothing in this brief treats it as able to.
- **T2 against T1 is governed by the same order.** If the owner adopts the reading in which the §8.4
  display is usable character knowledge (IB-4, §5.2), then A2b, E1c and D2 govern and §8.4's display is
  displaced **by the existing rule**. This is a statement about what O-6a already decides, not a new
  ruling. The alternative route is equally available and is the owner's alone: **amend A2b/E1c** so that a
  derived notice is permitted. Both routes are open; choosing between them is the question.
- **T2 against T2 has no rule.** §8.4, K-R7 and the K-R13 consequence are all in r32. O-6a does not
  order documents against themselves, and this brief does not supply a rule that does not exist.

### 2.3 One thing the ratified rows do *not* say

No row in the 118 requires a recipient to be shown that their account's source was retracted. A2b
constrains the character at the retraction turn; E1c forbids a notice account at that turn; H10 and E1b
require persistence. The behaviour §8.4 illustrates — a *second* knower being told, in their channel,
that what they hold was withdrawn — has **no ratified row**. Whatever the owner decides, it is a
decision to add a requirement or to let the example lapse; it cannot be read off the existing set.

---

## 3. Established versus undecided

### 3.1 Established — not reopened here

1. A source retraction removes no recipient account (O-3; H10; E1b; A2a; v9 r49).
2. A retraction changes truth status only (E1a) and reaches no belief (O-5; v9 r15/r38; S34).
3. Forgetting and source retraction are different operations on different streams (K-R3, K-R11, K-R19).
4. A character who learns of a retraction has acquired that knowledge (K-R13 consequence; A2b).
5. Basis retraction does not undo promotion (v9 r29; K16); authority changes rewrite no character state
   (v9 r38; K24).
6. A recipient's record never resolves to the source record (O-1; v9 r12; J8).
7. History survives every operation in this brief (v9 r20; K-R3).

### 3.2 Undecided — the three questions this brief prepares

| ID | Question | Why source cannot settle it |
|---|---|---|
| **Q-A** | Must the system be able to recognise that two characters' separate accounts concern the same underlying claim? | No requirement states sameness; v9 refuses it in code with a rationale about canon enumeration, which is T3 reasoning, not T1 authority |
| **Q-B** | Is a knower's channel permitted to report a retraction they were never told about? | r32 contradicts itself (§8.4 vs K-R7/K-R13); A2b's precondition is underspecified; E1c forbids the stored form only |
| **Q-C** | If such a notice is permitted, may it name who retracted the information? | Neither model records a retractor; v9 §8 permits "who communicated it" but not the act of withdrawal |

---

## 4. Q-A — relationship and identity

### 4.1 The question contains two questions

Q-A as posed — *can two characters hold separate, exclusively owned accounts of the same information
without sharing access to each other's private records?* — has an unambiguous answer for the first part
and an open one for the second.

| Sub-question | Answer | Label | Evidence |
|---|---|---|---|
| Can two characters each hold their own account of one communicated claim? | **Yes — and it already happens** | Observed | `engine.py` `disclose` creates one account per recipient with a fresh descriptor; `records.py` gives each its own `record_id` and owner |
| Is each account exclusively owned? | **Yes** | Observed + Source-derived | `records.py:__post_init__` requires exactly one `owner_character` on a private record; O-1 |
| Does either account grant access to the other, or to the source? | **No** | Observed | descriptor `{speaker, record_class, turn, event}` carries no record id (`engine.py:575–700`); v9 r12; J8 |
| Can the system tell that the two accounts concern the same claim? | **No** | Observed | `proposition_id` refused on `claim`, `account`, `belief` (`records.py:14–20`, `:155–158`) |
| Does any T1 decision require that it could? | **No** | Unestablished | No ledger entry, no conformance row |

**So Q-A is a sameness question.** The ownership question is settled; the identity question is not.
That distinction is what makes the hybrid (RA-3) coherent: it adds sameness without touching ownership.

### 4.2 Shared identity versus shared ownership

These are routinely conflated and must stay apart.

| | Shared ownership | Shared identity |
|---|---|---|
| Means | one record, several owners, or one owner and others with access | several records, each singly owned, recognised as being *about* one claim |
| Who may read | every co-owner | the owner of each record only |
| What it costs | O-1 directly — B would gain access to A's original | nothing in O-1; the question is whether the identity resolves to a record |
| Which relation does it | RA-1 (r32) | RA-3 (hybrid), if the identity is non-resolvable |

r32's relation is **shared identity plus shared access**, because the edge names the memory itself
(`K-R13`; `model.py:347–353`). RA-3 is **shared identity with no shared access**, because the identity is
a name for content, not a record id. The difference is not stylistic: one violates O-1's "no access to A's
original", the other does not.

### 4.3 The three relations against the tier-1 decisions

| | RA-1 r32 shared memory + edges | RA-2 v9 copied accounts | RA-3 hybrid: copies + content identity |
|---|---|---|---|
| O-1 exclusive ownership | **fails** — the memory is chat/world-owned; no character owns the original | passes | passes |
| O-1 no access to source | **fails** — the edge carries `memory_id` | passes | passes, **if** the identity resolves to nothing |
| O-3 recipient persistence | passes | passes | passes |
| O-13 containment | passes (all inside the chat) | passes | passes |
| T3 rule 12 | fails (resolvable pointer) | passes | passes only if the identity is not a pointer |
| v9 `proposition_id` refusal | n/a | enforced | **the obstacle** — refused on private records today |

**RB-06 — the only tier-3 rule standing between RA-2 and RA-3 is a refusal whose stated reason is
canon-enumeration safety.** The docstring gives that reason explicitly: a private record attached to a
proposition would give enumeration "something to traverse" (`records.py:14–20`). RA-3's identity would be
scoped to private records and would resolve to none of them, so the stated rationale is not automatically
violated — but establishing that is a design task, not a reading of the text. It would be new architecture
under O-6a case 3, and would have to be declared as such.

### 4.4 Do the acquisition circumstances require different relations?

The sources already treat the three circumstances differently. Whether that difference is *required* is a
separate question, and no source answers it.

| Circumstance | r32 | v9 | Sameness needed? |
|---|---|---|---|
| **Public / declared** | `ambient_authored`; audience derived from class + scope; often **no edge at all** | one account per named recipient at declaration (`narrator_grant`); B1/B6 | Not for the notice — there is no per-knower link in r32 either — but the *retraction* of a declared record has no defined behaviour (RQ-04) |
| **Private disclosure** | one `told` event per recipient, each an edge to the same memory | one copied account per recipient, descriptor only | **Yes**, if the notice is owed |
| **Independent discovery** | one memory, two `witnessed` edges | two accounts the system cannot relate | **Yes**, and this is the case that most needs it — two witnesses of one event have no way to be told of a later retraction |

The third row is the one that makes sameness more than a retraction convenience: **without it, two
characters who witnessed the same event are permanently unrelated records**, which also removes any
cross-knower conflict detection. This is Phase 1A's CO-03 and 1A-1's RP-06, restated.

---

## 5. Q-B — retraction awareness

### 5.1 What the source text actually establishes about "the channel"

r32 uses "channel" as a term of art in one place only: §9, *Prompt assembly: channels and guards*
(`:427`), whose table assigns C1 authoritative constraints, C2 character knowledge, C3 beliefs, C4
unverified reports, C5 unknown. §8.4's "A's channel" is therefore not ambiguous in itself — it is a
**prompt-assembly channel**, i.e. character-facing context supplied to the model (IB-3). Three further
facts then decide how much that placement implies.

- **C2 material "may be used and spoken"** (`:432`). The channel is not a display surface; it is the
  character's usable knowledge.
- **K-R7** (`:207`) requires an `AcquisitionEvent` to move information into that channel.
- **The retraction fact lives in the memory stream**, not in the acquisition stream (K-R12, K-R13), so
  rendering it into A's channel is a move from a broader context into the character channel.

**RB-01 (Source-derived):** the §8.4 notice is placed in the channel whose contents a character may use
and speak. It cannot be treated as inert system metadata merely because it was computed rather than
stored. The computation answers "where does this fact come from"; it does not answer "may the character
act on it".

### 5.2 Four interpretations, and what each is supported by

| ID | Interpretation | Where the notice appears | Text that supports it | Text that resists it |
|---|---|---|---|---|
| **IB-1** | Internal state only | nowhere in the prompt; inspection views may show it | K-R13 (join is computed, never stored); K-R5/R2 `inspection` context | §8.4, which shows it in "A's channel" |
| **IB-2** | Human-facing only | application presentation / `inspection`; never C2 | R2's `inspection` row ("user-facing audit of who knows what"); K-R5 | §8.4 again |
| **IB-3** | Character-facing, marked non-knowledge | in the prompt, labelled as metadata the character must not speak | §8.4's second-person phrasing | C2's "may be used and spoken"; there is no "non-usable C2" category in r32 |
| **IB-4** | Character-facing and usable | in C2, usable and speakable | §8.4, read as normative | K-R7; K-R13 consequence; A2b; E1c; D2 |

The four are not equally supported, and saying so is not choosing. IB-4 is the only reading that
satisfies §8.4's sentence as written. IB-1 and IB-2 are the only readings that satisfy r32's rules about
how information enters the character channel. **IB-3 is the reading with no textual home at all**: r32
defines five channels and none of them is "visible but unusable".

### 5.3 The accounting test, stated precisely

The question underneath all four readings is not "is the notice stored?" — E1c already answers that. It
is:

> **Is the fact "this claim was retracted" information the character did not previously hold?**

If yes, then any appearance of it in the character's usable context is *acquisition without an event*, and
A2b/E1c/D2 are engaged regardless of whether a row was written. If no — because the character is
reported to have known it all along — then K-R13's consequence paragraph is wrong, since it says learning
of the retraction is an acquisition event.

There is no third position in which the fact is simultaneously new to the character and not knowledge.
**That is the whole of Q-B's substance**, and it is why the stored-versus-computed distinction does not
dispose of it.

### 5.4 The six sub-questions, answered one by one

| # | Sub-question | Answer | Label |
|---|---|---|---|
| 1 | Does receiving a retraction notice constitute acquiring new information? | **If it enters C2, yes.** The fact is not derivable from what the character held: they held "X was told to me", not "X is withdrawn". K-R13's consequence treats learning it as an acquisition event | Source-derived |
| 2 | Can a character use information that the source was retracted, without an acquisition basis? | **Not under r32's rules.** K-R7 requires an event to move it into the channel; v9 r9/r10/r11 independently require an account created by exposure. Two sources converge here | Source-derived |
| 3 | Can the application display the retraction to the human without exposing it to the character? | **Yes, unambiguously.** r32's R2 table gives `inspection` ("user-facing audit of who knows what") and `narrator` ("must never inject it into the character channel"). Human-facing presentation is not the contested question | Source-derived |
| 4 | Does E1c prohibit every notice, or specifically an unaccounted notice account? | **The unaccounted one.** E1c's expected outcome is "no notice account created" and its failure condition is "retraction knowledge appears unaccounted". A notice carried by a real disclosure event is an ordinary attributed account and violates neither half | Source-derived |
| 5 | Is a read-time join sufficient to account for character knowledge, or does it merely expose internal system knowledge? | **It is sufficient to represent, insufficient to account.** The join represents "the content this character holds is withdrawn" without storing anything (K-R13). It does not supply the character's learning of that fact; for that, K-R13's own consequence demands an event. The join is internal system knowledge until something moves it into the channel | Source-derived + Inferred (application to §8.4) |
| 6 | May metadata be redefined as non-knowledge once it reaches the model's context? | **No.** C2's contents are usable by construction (`:432`), so a label inside the prompt is a request, not a boundary. IB-3 is therefore not a distinct behaviour; it is IB-4 with a disclaimer | Inferred |

### 5.5 The internal inconsistency, stated without resolving it

r32 contains three texts on one page and no priority among them:

- §8.4 (`:419–423`) shows A the retraction, in A's channel.
- K-R7 (`:207`) forbids moving information into that channel without an acquisition event.
- K-R13's consequence (`:113–115`) makes learning of the retraction an acquisition event.

**RB-02 (Observed + Source-derived):** §8.4 and the two normative rules cannot all hold for a recipient
who was not separately told. One of three things must be true — §8.4 is illustrative and not binding; the
notice is not "information moved into the channel"; or K-R7/K-R13 are narrower than they read — and the
document does not say which.

**RB-03 (Observed):** no probe tests the notice. S34 asserts that A still believes, that the acquisition
is intact, and that another character behaves differently; it does not assert that A's channel displays
anything, and it names no recipient. The preregistered experiment carries one retraction item, asked of
the arm's characters after a rollover, with no recipient-notice assertion
(`MEMORY-09-PREREGISTRATION-2026-11-15.md:101,109`). So the display has never been an instrumented
requirement in either the contract's probes or the experiment's.

**RB-04 (Source-derived):** E1c's two halves are narrower than "no notice". Only the stored account and
the unaccounted knowledge are forbidden. A notice that arrives because someone actually disclosed the
retraction is permitted and is what A2b contemplates.

**RB-05 (Source-derived):** D2 excludes the eventless notice independently of E1c. "Knowledge follows
experience, exposure and authorised transfer" — a computed notice is none of these. Because D2 is T1 and
does not use the word "account", it closes the reading in which E1c is satisfied by the notice not being a
row. This is the strongest single objection to IB-4 and it comes from the tier above the conformance map.

**A2b's underspecification (Unestablished):** the row's precondition reads "Same character, same turn",
and the map does not say whose character or which turn. Under one reading it constrains the retracting
character at the retraction turn (they do not thereby know the retraction's effects); under another it
constrains every knower at that turn. 1A-1's R-split reading depends on the second, and the text does not
establish it. **This ambiguity is load-bearing for Q-B and should be closed by the owner's answer, not by
interpretation.**

---

## 6. Q-C — attribution

### 6.1 Three different claims, often treated as one

| Form | Text of the notice | What it reveals | What it requires |
|---|---|---|---|
| **Impersonal status** | "the claim you were told is no longer asserted" | that the content is withdrawn; nothing about who | nothing beyond the content identity (Q-A) |
| **Attributed statement** | "Alaric retracted it" | Alaric performed an act of withdrawal | that the act is an attributable fact the system holds |
| **Disclosure of the act** | Alaric says "I retracted it" | the same, but established by his own output | an ordinary disclosure event; recipient accounts; J3/J7 |

### 6.2 What the sources permit

- **v9 §8 permits the recipient to hold "who communicated it".** B already knows Alaric told them; the
  descriptor names him (`engine.py:575–700`). So naming Alaric is not itself a new disclosure *of the
  communication*.
- **But the act of retraction is not the communication.** "Alaric withdrew this" is a fact about Alaric's
  later conduct. Under v9 r10 a private fact becomes shared information only through a disclosure event,
  and r11 limits acquisition to exposure. A system-generated sentence naming him is neither.
- **r32 has the same shape from the other side.** K-R7 governs the move into the channel, and D2 governs
  what counts as knowledge. Attribution adds a fact about a person, so it is strictly more than the
  impersonal notice, not a variant of it.

### 6.3 The precondition nobody has built

**RB-07 (Observed):** no implementation records who retracted anything.

| Site | Signature | Actor recorded? |
|---|---|---|
| `engine.py:542–559` | `retract(record_id, turn)` | no — any caller, no attribution |
| `model.py:282–296` | `retract(record_id, turn, source_span, reason)` | no — the span is evidence, not an actor |
| `form.py:174–204` | `_retract(text, owner_id, turn, source_span)` | derived from the speaking span; the *record* carries the speaker |
| r32 K-R12 | lifecycle stream `established \| superseded \| corrected \| retracted` | the status is named; no actor field is specified |
| v9 r42 | correction / supersession / revocation | actor required for **authority** operations (`_require_human`), not for `retract` |

So an attributed notice is not merely a policy question — it is **not computable** from current state. To
answer Q-C "yes", three things must be true that are not: retraction has an actor; that actor is
recorded; and the act is disclosable to the recipient by the actor's own output or by an event the
conversation layer establishes (rule 47 permits narration to *report* an established disclosure, never to
manufacture one).

### 6.4 Does attribution create a separate obligation?

**Yes, on either source's rules** (Source-derived):

- It discloses a fact about a third party (that they withdrew something), which r10/r11 and D2 treat as
  an acquisition requiring exposure.
- It is exactly the case rule 47 forbids narration from manufacturing: "narration may report a disclosure
  event established by the conversation layer, but must not manufacture one."
- If the act is instead established by Alaric's own output, then it is an ordinary disclosure and E1c is
  irrelevant — but so is the notion of a "notice". The recipient learns it the way they learn anything.

**Consequence for the hybrid:** the impersonal form keeps RA-3 inside rule 12 and inside O-1. The
attributed form does not add a privacy leak (it names no contents), but it does add a disclosure obligation
that lands on Alaric's model, not on MEMORY-09.

---

## 7. Behaviour matrix

Rows are behaviours a requirement or rule names. Cells state whether the interpretation conforms.
"Conforms" means no named text is violated; it is not a recommendation.

| Behaviour / rule | IB-1 internal | IB-2 human-facing | IB-3 prompt, marked unusable | IB-4 prompt, usable |
|---|---|---|---|---|
| **A2b** no automatic knowledge of the retraction | conforms | conforms | **at risk** — see note | **violates** |
| **E1c** no notice account; no unaccounted knowledge | conforms | conforms | **at risk** | **violates** |
| **D2** knowledge only by experience, exposure, transfer | conforms | conforms | **violates** — the fact is in context | **violates** |
| **K-R7** event required to enter the character channel | conforms | conforms | **violates** | **violates** |
| **K-R13** join computed, never stored | conforms | conforms | conforms (still computed) | conforms |
| **§8.4** display shown to A | **does not satisfy** | **does not satisfy** (A is not the human) | partially — shown, but disclaimed | **satisfies** |
| **O-1 / v9 r12** no access to source | conforms | conforms | conforms | conforms |
| **Rule 47** narration may not manufacture a disclosure | conforms | conforms | conforms | n/a — not a disclosure, so outside the rule |
| **Observable in the 118** | yes (A2a/E1b/H10) | yes | **no test distinguishes it from IB-4** | yes |
| **Preregistered probe coverage** | retraction item only | none | none | none |

**Note on IB-3.** "Marked unusable" is a prompt-level instruction. A2b and E1c are behavioural rows: they
fail if retraction knowledge *appears*. A model that reads the notice and speaks it fails them whether or
not the notice was labelled metadata. Since no test can distinguish IB-3 from IB-4 without inspecting the
model's compliance with an instruction, IB-3 is **unfalsifiable as a requirement** — which is why it has no
home among r32's five channels.

**What the matrix shows without ranking.** IB-4 is the only interpretation that satisfies §8.4, and it is
the only one that violates A2b, E1c, D2 and K-R7. IB-1 and IB-2 are the only ones that satisfy all four of
those. **No interpretation satisfies both §8.4 and the ratified set.** That is the decision, and it is a
decision about which text to change, not about which text is true.

---

## 8. Hybrid viability under each interpretation

| | IB-1 internal | IB-2 human-facing | IB-3 marked unusable | IB-4 usable |
|---|---|---|---|---|
| **RA-1** r32 shared memory | viable but pointless (no notice to compute) | viable, same | viable | viable, and already satisfies §8.4 — at the cost of O-1 |
| **RA-2** v9 copies | **viable** | **viable** | viable by assertion, unfalsifiable | **not viable** — no way to compute the notice across copies |
| **RA-3** hybrid | viable; identity **not needed** | viable; identity not needed | viable; identity not needed for behaviour | **viable**, and the only relation that is — identity required |

Three results follow, and they are the reason the answering order in §12 is not the one Phase 1A-1
proposed.

1. **Under IB-1 or IB-2, RA-3's content identity buys nothing for retraction.** If no notice reaches a
   character, there is no join to compute and no sameness to establish. The identity would then have to be
   justified on other grounds (cross-knower conflict detection), which no requirement asks for.
2. **Under IB-4, RA-3 is the only viable relation.** RA-2 cannot compute the notice; RA-1 can, but fails
   O-1. This is why Q-B determines whether Q-A's sameness question must be answered at all.
3. **Under IB-3 every relation is "viable" and nothing is established**, because the interpretation
   cannot be tested. An owner who answers Q-B with IB-3 has not decided the behaviour; they have deferred
   it while changing the prompt text.

**RA-1 is not made viable by any interpretation.** Its failure is O-1 (T1), which no reading of the notice
touches. RQ-01 cannot revive it.

---

## 9. Consequences, by surface

### 9.1 Character knowledge

- **IB-1/IB-2:** unchanged. The character holds what they acquired, indefinitely, and learns of the
  retraction only if someone tells them. A2a, A2b, E1b and H10 all hold, and the harness's inability to
  create a recipient (RP-03) is the only obstacle to demonstrating them.
- **IB-3:** nominally unchanged, actually indeterminate. What the character knows depends on whether the
  model obeys the label — which is a prompt-compliance question, not a memory question, and no
  requirement governs it.
- **IB-4:** the character gains a fact with no acquisition event. D2 is violated at the tier above
  A2b/E1c, and the violation is not repairable by adding a requirement — it requires amending D2 or
  accepting that the memory system asserts what the character knows. That is a larger decision than Q-B
  appears to be, and the owner should see it before answering.

### 9.2 Human-facing presentation

- Permitted in every interpretation. r32's `inspection` context exists for exactly this ("user-facing
  audit of who knows what") and `narrator` may hold what the character does not, provided it is not
  injected into the character channel.
- So the application can always show the human: "Sable holds a claim that was retracted at turn 210."
  That display is not in question, requires no new rule, and is available today in principle — though not
  in either implementation, because neither has the join (RP-13).
- **The mistake to avoid is treating the human display as evidence about the character display.** They
  are different channels in r32's own table, and Q-B is about the second one.

### 9.3 Privacy

- **Impersonal notice:** reveals only that a content identity is no longer asserted. It names no actor, no
  record and no history. Under v9 r12 and O-1 this is compatible, because nothing resolves.
- **Attributed notice:** reveals Alaric's act. It does not reveal his record's contents, so rule 12 is not
  breached; but it does make a fact about him available to B without exposure, which r10/r11 govern.
- **Neither form creates a resolvable pointer**, so the hybrid's privacy claim survives both — provided
  the content identity is not enumerable in a way that reconstructs who holds what. An identity that can
  be used to ask "who else holds C?" becomes an audience projection, which F1 forbids.

### 9.4 Disclosure accounting

| Notice form | Account created | A2b | E1c | Ledger status |
|---|---|---|---|---|
| none | none | conforms | conforms | consistent |
| computed, internal | none | conforms | conforms | consistent |
| computed, human-facing | none | conforms | conforms | consistent |
| computed, character-facing | none | violates | violates ("appears unaccounted") | **inconsistent with T1** |
| disclosed by the retractor's own output | ordinary recipient account | conforms | conforms | consistent |

The last row is worth stating plainly, because it is the one route that gives a character retraction
awareness without contradicting anything: **if Alaric says "I retracted that", B acquires it the ordinary
way.** It is not a notice, it is a disclosure; it needs no new rule; and it cannot be manufactured by the
system or by narration (rule 47). What it does not do is inform a character who was never told — which is
precisely the behaviour §8.4 illustrates.

---

## 10. What cannot be resolved from existing authority

| ID | Gap | Label | Why it matters |
|---|---|---|---|
| **GB-1** | §8.4, K-R7 and the K-R13 consequence disagree, and r32 states no priority among its own texts | Observed | Q-B is not answerable by reading r32 more carefully; only by deciding which text binds |
| **GB-2** | A2b's precondition "Same character, same turn" does not say whose character or which turn | Unestablished | R-split rests on one reading of it; the other reading makes A2b redundant with E1c |
| **GB-3** | Retraction has no recorded actor in either model | Observed | Q-C cannot be implemented as posed; attribution requires a new field and a new obligation |
| **GB-4** | E1c's "no notice account" and "appears unaccounted" are one row covering two different prohibitions (stored row vs knowledge in context) | Source-derived | The row is satisfied by IB-4 in its first half and violated in its second, which is why the conflict is invisible when the row is read once |
| **GB-5** | Retraction of a *declared* record has no defined behaviour; B6 covers withdrawing the declaration | Source-derived (RP-07) | Public information is the case where §8.4's join would be most useful and least costly, and it is the case with no rule |
| **GB-6** | "The original information was presented mistakenly" has no operation. r32's lifecycle has `corrected` and `retracted`; v9 has correction, supersession, revocation — none of them means "the telling misrepresented the content" | Observed | Three different events are currently expressible only as retraction, which changes the recipient's position in all three |
| **GB-7** | No requirement states that sameness must be representable | Source-derived (RP-06) | Q-A's sameness half has no authority behind it in either direction; the tier-3 refusal is the only text, and it is a design choice |
| **GB-8** | "Still holds what it acquired" (A2a) is true at rest and false in the engine's only character surface when the account itself is revoked | Observed (RP-02) | A2a's meaning must be fixed before any notice behaviour is layered on it |
| **GB-9** | §8.4's "you still believe it" cannot bind the memory system (O-5, T1) | Source-derived (RP-10) | The sentence is already partially void, which weakens any argument that it binds as a whole |

**Nothing in this table is a defect report.** Each is a place where the sources are silent, self-divergent,
or answer a narrower question than the one asked.

---

## 11. Recommended wording for the owner questions

Each question is stated in plain language, with the behaviour it changes. No option is marked preferred.

### Q-B — retraction awareness (ask first)

> **When a fact a character was told is later retracted, may that character's own context report the
> retraction even though nobody told them about it?**

*What this decides:* whether the character's model is ever informed that something they hold has been
withdrawn, without a disclosure.

*Concrete example.* Sable is told by Alaric at turn 40 that the north road is open. At turn 210 Alaric
retracts it. At turn 211 Sable is asked about the north road.

- **Answer no:** Sable answers from what she was told. She may say "Alaric told me the road is open." She
  has no way to know it was withdrawn unless someone tells her. A2b, E1c and D2 stand unamended, and
  §8.4's example is recorded as illustrative rather than binding. This is what v9 already does, and what
  both implementations do.
- **Answer yes, as displayed but not usable:** Sable's context shows a warning that the claim is
  withdrawn, marked as system information she must not act on. No test can distinguish this from the
  next answer, and no r32 channel expresses it.
- **Answer yes, as usable knowledge:** Sable may say "that was retracted." This satisfies §8.4 and
  **requires amending A2b, E1c and D2**, because she now knows something she was never exposed to.

*Sub-question the owner should answer in the same breath:* may the human operator see the retraction in
the application's own views? (On every source, yes — that is the `inspection` context, and it is not the
same question.)

### Q-A — relationship and identity (ask second)

> **Must the system be able to tell that two characters' separate accounts concern the same underlying
> claim?**

*What this decides:* whether a shared, non-resolvable content identity is introduced, and therefore
whether a notice could ever be computed for a recipient — or for two independent witnesses.

*Concrete example.* Alaric tells Sable and Brannick the same thing in one speech. Later the source is
retracted.

- **Answer no:** the system holds two unrelated accounts. Neither can be linked to the retraction or to
  the other. This is v9 today, and it is consistent with every tier-1 decision.
- **Answer yes:** both accounts carry a shared name for the claim, which resolves to no record. A notice
  becomes computable for either or both. This is new architecture, and it must not become an audience
  enumeration (F1).

*Why it comes second:* under a "no" to Q-B, the identity buys nothing for retraction and would have to be
justified on other grounds, which no requirement provides.

### Q-C — attribution (ask third, only if Q-B is yes)

> **If the notice is shown, may it say who retracted the information?**

*What this decides:* whether the character learns a fact about another character's conduct, and therefore
whether the retracting character's model must have disclosed it.

*Concrete example.* The notice could read "the claim you were told is no longer asserted" or "Alaric
retracted it."

- **Impersonal:** reveals no actor, needs no disclosure, stays inside rule 12.
- **Attributed:** reveals Alaric's act. It must be established by his own output or by an event the
  conversation layer records; the system may not manufacture it (rule 47). It also requires retraction to
  have an actor, which no model currently records.

---

## 12. Recommended answering order

**Q-B → Q-A → Q-C.** Phase 1A-1 proposed Q-A first; the source review changes that, for three reasons.

1. **Q-B decides whether A2b/E1c/D2 must be amended.** That is a tier-1 change, and it should be decided
   before any design work that depends on it. Answering Q-A first invites a relation to be chosen for a
   notice that may not be owed.
2. **Q-A's sameness half is only forced by Q-B.** Under a "no" to Q-B, RA-2 remains fully conformant and
   no content identity is required (§8, result 1). Under a "yes", RA-3 becomes the only viable relation.
   The dependency runs one way.
3. **Q-C is not answerable meaningfully until Q-B is yes**, because an impersonal notice is the default
   and the attributed form adds an obligation that has no implementation basis (RB-07).

**What each answer locks:**

| Answer | Consequence |
|---|---|
| Q-B = no | A2b/E1c/D2 unchanged; §8.4 recorded as illustrative; RA-2 viable; Q-A not forced; Q-C moot |
| Q-B = yes, unusable | Nothing is actually decided; the behaviour is deferred to prompt compliance; no test can check it |
| Q-B = yes, usable | A2b, E1c and D2 must be amended explicitly; RA-3 required; Q-C becomes live |
| Q-A = yes (after either) | content identity must be designed, declared as new architecture under O-6a case 3, and must not be enumerable |
| Q-C = attributed | retraction needs an actor; the act needs a disclosure path; rule 47 constrains narration |

**Q-D does not block this.** The conversation-state question (1A-1 §13) governs *from what* a notice is
computed — a state-relative join — not *whether* one is owed. Q-B can be answered before Q-D without
prejudging it, and the answer does not change when the state dimension is added.

---

## 13. What this brief does not do

No decision on Q-A, Q-B or Q-C. No relation adopted, no requirement amended, no rule marked void, no
example declared non-binding. No register entry, no corpus, no probe, no threshold. No implementation
change, and no claim about runtime behaviour — every behavioural statement is a static reading at the
hashes in Appendix C, and nothing was executed.

Two things are stated with more confidence than Phase 1A-1 could state them, and both are offered for
challenge:

- **The conflict is inside r32 as well as between r32 and the ratified rows** (RB-02). If that is right,
  Q-B cannot be answered by choosing a tier; it can only be answered by deciding which text to change.
- **D2, not E1c, is the binding obstacle to the eventless notice** (RB-05), because it is T1 and does not
  depend on reading the word "account".

**End of brief.** No canonical artifact was modified. The next action is the owner's answer to Q-B.

---

## Appendix A — Implementation traces

| ID | Trace | What it shows |
|---|---|---|
| **RBT-01** | `engine.py:542–559` `retract(record_id, turn)` | no actor, no owner check, no authority check (RP-08) |
| **RBT-02** | `engine.py:575–700` `disclose` | per-recipient descriptor `{speaker, record_class, turn, event}`; no source id; no status check on the source |
| **RBT-03** | `m09/retrieval.py:46–52` `perspective` | pool is `owner_character == character and status != "revoked"` — a retracted own account becomes unreachable to its owner (RP-02) |
| **RBT-04** | `records.py:14–20`, `:155–158` | `proposition_id` refused on `claim`/`account`/`belief`; reason given is enumeration traversal |
| **RBT-05** | `model.py:282–296` `retract` | revoke + a second `retraction:<predicate>` record; no actor |
| **RBT-06** | `model.py:347–353` `edges()` | one edge per event, including `via="revoked"` — contradicts K-R19 (RP-05) |
| **RBT-07** | `form.py:174–204` `_retract` | retraction confined to the speaker's own active records; narrator excluded |
| **RBT-08** | `form.py:206–240` `_admit` | `via="witnessed"` asserted for every formed record; no listener acquisition (RP-03) |
| **RBT-09** | `retrieve.py:100–143` `eligible` | four independent gates; no retraction-status join; no notice surface |
| **RBT-10** | `conversation.py` | per-turn retrieval only; no recipient formation path |

## Appendix B — Text locations for load-bearing quotations

| Quotation | Location |
|---|---|
| "may be used and spoken" (C2) | r32 §9, `:432` |
| "Moving information from a broader context into the character channel requires an `AcquisitionEvent`" | r32 K-R7, `:207` |
| "A's channel: *'you were told X; it has since been retracted; you still believe it'*" | r32 §8.4, `:421` |
| "whether they have *learned of the retraction* is itself an acquisition event" | r32 K-R13 consequence, `:113–115` |
| "This character's knowledge includes a retracted claim" is a read-time join | r32 K-R13, `:110` |
| "retraction does not correct minds… another character behaves differently" | r32 S34, `:846` |
| "Recipient may hold: who communicated it…" | v9 §8 |
| "recipient records carry a restricted origin descriptor, never access to the source record" | v9 rule 12 |
| "No notice account created" / "Retraction knowledge appears unaccounted" | conformance E1c |
| "Same character, same turn" / "No automatic knowledge of the retraction" | conformance A2b |
| "Knowledge follows experience, exposure and authorised transfer" | r32 §2 D2 |
| "Once disclosure occurs, subsequent retraction of the source does not erase the recipient's acquired account" | O-3, ledger |
| "a private record cannot be attached to a proposition" | `m09/records.py` docstring and `__post_init__` |
| retraction probe: 1 per conversation, ≈65% of turns, post-rollover | `MEMORY-09-PREREGISTRATION-2026-11-15.md:101,109` |

## Appendix C — Sources verified this session

All SHA-256 values re-computed from disk; each matches the 1A-1 record.

| Source | SHA-256 (first 16) | Tier |
|---|---|---|
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` (r32) | `e88ef2d1697da1d9` | T2 |
| `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | `698254ca262c0874` | T3 |
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | `da3e72398ecbaa0a` | T1 |
| `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | `b4ba33ccbfe8f8c1` | T1 (ratified) |
| `MEMORY-09-PHASE-1A-CONCEPTUAL-OBJECT-INVESTIGATION.md` | `82173a6f78ac84c4` | research |
| `eval/mem1/m09/engine.py` | `7505e9acca8e3229` | Observed |
| `eval/mem1/m09/records.py` | `ad2d4d11501c357b` | Observed |
| `eval/mem1/m09/retrieval.py` | `cd2d380bd5d9d1f7` | Observed |
| `eval/mem1/model.py` | `813cdbff56c01ae4` | Observed |
| `eval/mem1/form.py` | `f472052738388bab` | Observed |
| `eval/mem1/retrieve.py` | `84f90b5ce0b3619f` | Observed |
| `eval/mem1/conversation.py` | `fa1d9ce630462fa9` | Observed |
