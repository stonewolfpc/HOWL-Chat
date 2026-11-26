# MEMORY-09 — Final event-ownership reconciliation

**Status: additive supplement. It supersedes §1 of `MEMORY-09-REVIEW-SUPPLEMENT-EVENTS-AND-TIME-2026-11-17.md`
only in its confidence, not in its evidence. No code, no existing review document, and no contract was changed.
Nothing was executed.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Scope | Event-record ownership under v9 §1 and §7; the standing of the one-claim proposal; K9 restated |
| Code changed | **None** |
| Existing documents changed | **None** — this file is additive |
| Executed | **Nothing** |
| Prior position withdrawn | Supplement §1.5 "smallest permissible resolution" is downgraded from *resolution* to *candidate conditional on reading A* |

---

## 1. The exact clauses, and what each one assigns

The question is not whether a disclosure event exists — all clauses agree it does. The question is **which
component writes it down**. The clauses divide into three kinds: those that name an owner explicitly, those
whose grammar implies one, and those that are silent.

| # | Clause (verbatim) | Source | Layer assignment |
|---|---|---|---|
| C1 | "Whether disclosure actually occurred \| **Conversation / event layer**" | v9 §1 | **Explicit** — determination |
| C2 | "What information was communicated \| **Conversation / event layer**" | v9 §1 | **Explicit** — determination |
| C3 | "Which characters were exposed \| **Conversation / event layer**" | v9 §1 | **Explicit** — determination |
| C4 | "What was actually said \| **Conversation / event layer**" | v9 §1 | **Explicit** — determination |
| C5 | "Recording memories and acquisitions \| **Memory system**" | v9 §1 | **Explicit** — persistence |
| C6 | "Preserving provenance \| **Memory system**" | v9 §1 | **Explicit** — persistence |
| C7 | "**Conversation layer** establishes that disclosure occurred, its content and recipients — **recorded as a claim event with spans**." | v9 §7 step 2 | **Implied by grammar** — the em-dash clause is predicated of the conversation layer |
| C8 | "**Memory system** creates, per exposed character, their own acquired account: source, event/turn, communicated content, origin descriptor." | v9 §7 step 4 | **Explicit** — persistence, accounts only |
| C9 | "**Exposure** follows from that event alone." | v9 §7 step 3 | **Silent** on persistence; the event is the authority for exposure |
| C10 | "The system records the event and its attribution, creates recipient-specific acquired accounts, and does not expose the source record." | rule 46 | **Ambiguous** — "the system" is not defined as MEMORY-09; the sentence bundles two acts |
| C11 | "Linking is not disclosure — disclosure is the establishment of a new **attributable claim**…" | rule 10 | **Silent** on author |
| C12 | "**A statement is a claim** — a character's utterance records that they said it" | rule 25 | **Implied** — the *utterance* is what records it |
| C13 | "it creates an attributable claim and a recipient-specific account" | rule 48 | **Ambiguous** — "it" is the disclosure act |
| C14 | "narration may report a disclosure event **established by the conversation layer**" | rule 47 | **Implied** — the event's establishment is the conversation layer's |
| C15 | "Public statement / claim \| **No** \| that a character said something, and what was said \| never" | v9 §3 | **Names MEMORY-09's class** — fixes the shape, not the author |
| C16 | "The edge fold depends **only** on acquisition-family events for (character, memory): `witnessed \| told \| read \| inferred \| verified \| narrator_grant \| revoked`" | r32 K-R11, as quoted in REV 2 §5.10 | **Implied** — MEMORY-09's canonical graph is acquisition edges; a disclosure event is not among them |
| C17 | "**Recipient may hold:** who communicated it; when and in what context; the content communicated… **Recipient must not obtain:** the source record… or any pointer resolving to them." | v9 §8 | **Explicit** — access boundary |

**Two clusters, one unresolved subject.**

- **Cluster A — recording belongs to the memory system:** C5, C6, C10 (first clause), C13, C15.
- **Cluster B — establishment and its record belong to the conversation layer:** C1–C4, C7, C14, C16.

Neither cluster contains a clause that says "the memory system creates the disclosure claim record" or "the
conversation layer persists it in MEMORY-09's store". **The sources do not settle ownership.** That is the
finding, and I am not going to manufacture the missing sentence.

---

## 2. The two readings, and where they actually differ

| | **Reading A** — MEMORY-09 creates the event claim | **Reading B** — the conversation layer creates it; MEMORY-09 receives facts |
|---|---|---|
| Who writes the record | `disclose()` inside MEMORY-09 | The event layer, via `claim()` or its own store |
| What "recorded as a claim event" means | The record's class is `claim`, written by the disclosure operation | The event is represented *in* MEMORY-09's claim vocabulary, but authored by the other layer |
| What MEMORY-09 receives | Discloser, recipients, content, span — and creates the claim from them | An existing event claim's identity, or only its facts |
| §7 step 2 grammar | Read as describing the *outcome* of the layer's establishing act, not its storage | Read literally: the layer records it |
| Rule 46 "the system records" | "The system" = the memory system | "The system" = the platform as a whole |
| K-R11 fold | Unaffected — a claim contributes no acquisition edge | Explains why: MEMORY-09's graph carries `told`, never the event |
| Where the disagreement lies | Not in *what* is recorded, nor in its class, nor in its access rules — only in **which component holds the write** | |

**What both readings agree on** (established, not interpretation):

1. The record, wherever written, is a **claim** — attributable, non-truth-apt, never truth, never belief (C15,
   rules 25, 48).
2. Its content is "what was said", not what it implies (rule 25; §3 class table).
3. It is not a knowledge source: reach is not knowledge (D-02), and no claim may enter `_eligible_truth`.
4. Recipients must not obtain it or any pointer to it (C17, §8).
5. Exposure is established by the event, not inferred by the memory system (C9, rule 11, §1 prohibition on
   "infer disclosure from proximity").

So the architecture settles the **class, the content, the truth status and the access rule**. It does not settle
the **writer**.

---

## 3. The five questions

### Q1 — Does the conversation layer create and persist the claim event, or pass established facts to MEMORY-09?

**Unresolved by the sources, and phrased as a false binary in one respect.** Three sub-questions hide here:

| Sub-question | Answer | Basis |
|---|---|---|
| Who *determines* that a disclosure occurred, its content, recipients, spans? | The conversation layer | C1–C4, explicit |
| Who *records* them? | Not stated for the event record | C7 grammar vs C10/C5 — conflicting |
| Who *persists* the result today? | MEMORY-09 is the only store in the experiment | Structural fact, not a contract clause |

That last row matters: the experiment has one record store. Under Reading B the event layer would either
(a) call MEMORY-09's `claim()`, i.e. use MEMORY-09 as its store, or (b) keep its own store, which no ratified
document describes. Reading B(a) collapses toward A at the persistence layer; Reading B(b) invents a component
the contract does not specify.

**Established: the determination is the conversation layer's. Unresolved: the write.**

### Q2 — What does "recorded as a claim event" mean in the ownership model?

**Interpretation, labelled:** it fixes the **shape and class** of the representation, not its author.

| It means | It does not mean |
|---|---|
| The event is represented as a claim — attributable, non-truth-apt (C15, rule 25) | That a new record class, registry or table is required |
| Its content is what was said, with spans where the event layer has them (C7) | That MEMORY-09 decides whether the disclosure happened |
| Its attribution names the speaker (rule 46's "and its attribution") | That the event is truth, knowledge, or an acquisition |

The phrase's only ownership content is negative: it forbids the event being recorded as truth (rule 48) or as a
belief (§7 step 7).

### Q3 — Does MEMORY-09 have responsibility for the event record, or only for recipient accounts?

**Established for accounts; unresolved for the event.**

| Object | Responsibility | Status |
|---|---|---|
| Recipient-specific acquired accounts, with descriptor, no resolvable link | **MEMORY-09** | Established — C5, C8, C17, rules 11, 12 |
| Source record's exclusivity and unchangedness | **MEMORY-09** | Established — §7 step 6, rule 5 |
| The disclosure **event** record | **Contested** — C7 assigns the act, C10 assigns the recording | Unresolved |

What can be said without forcing it: **MEMORY-09 has responsibility for the event record's class, attribution
and access, and contested responsibility for its creation.** Rule 46's first clause is the strongest single
piece of evidence for MEMORY-09 creation; §7 step 2's grammar and rule 47's "established by the conversation
layer" are the strongest against.

### Q4 — How do ownership and access interact with O-1, O-3, D-02 and narration?

This is where the ownership question stops being academic.

| Constraint | Exact text | Consequence for the event record |
|---|---|---|
| **O-1** | "A private record is retrievable only by its owner" (rule 5); v9's ownership model applies to records | If MEMORY-09 writes the claim, the claim has an owner. Today `claim()` sets owner = contributor, so the event would be **the discloser's private record** — retrievable by them alone, and by nobody else, including recipients |
| **§8** | "the recipient's origin reference is a **descriptor** (speaker, record type, turn), **not a resolvable identifier**" | Recipients keep the descriptor. They must not be able to resolve the event, so the event's ownership cannot be widened to them |
| **O-3 / rule 47** | "narration may report a disclosure event **established by the conversation layer**, but must not manufacture one, **reproduce a private record**, or speak private content on the owner's behalf" | If the event is the discloser's private record, narration reporting it must not reproduce it. Narration's licence is anchored to the *layer's* establishment (C14), not to a MEMORY-09 record — which is evidence that the event's canonical home is the layer, or that narration does not read the record |
| **D-02** | Reach is not knowledge | The event claim must never enter `_eligible_truth` or count as knowledge. Current code satisfies this structurally: `authoritative` is produced only from `_take` over truth-apt classes, and `claim` is not truth-apt |
| **rule 11 / §7 step 3** | "**Exposure** follows from that event alone" | The event is the authority for exposure. If MEMORY-09 creates it from caller arguments, it is recording an assertion, not deciding exposure — permitted; if it created the event without recipients supplied, it would be manufacturing exposure — forbidden |

**The consequence that follows from O-1 whatever the reading:** the event's access rule is *not* uniform. Under
Reading A it is the discloser's private record. Under Reading B(b), persisted outside MEMORY-09's store, it
escapes O-1 entirely — no owner, no `perspective()` membership, no rule 5 protection — which the architecture
nowhere authorises. Reading B(a) (the layer calls `claim()`) lands back on Reading A's access outcome.

**So access does not distinguish A from B(a). Only B(b) differs, and B(b) is unspecified.** That narrows the
question usefully: the choice is between *MEMORY-09 writes the claim* and *the caller writes it through
MEMORY-09's `claim()`* — both producing the same record with the same owner and the same access.

### Q5 — Does one claim per disclosure satisfy the full contract, or only the recording-shape clause?

**It satisfies the shape clauses and fails a different one — and the failure does not depend on ownership.**

| Requirement | Satisfied by "one claim per disclosure"? |
|---|---|
| Is a claim, non-truth-apt, never truth or belief (§7 step 7, rules 25, 48) | **Yes** |
| Content = what was said, with spans (§7 step 2, C7) | **Yes**, if the span is forwarded |
| Attribution names the speaker (rule 46) | **Yes** — contributor = discloser |
| Does not confer knowledge (D-02) | **Yes** — no acquisition edge, not authoritative |
| Recipients cannot resolve it (§8, rule 12) | **Yes** — descriptor unchanged |
| **J3: "Claim event with content and recipients"** | **No.** No path in the current API records the recipient set |
| Ownership (O-1) | **Ambiguous** — it produces the discloser's private record, which is one of the two readings, not a settled rule |
| Narration may report it (rule 47) | **Constrained** — reporting a private record is forbidden; only the event layer's establishment may be reported |

**The decisive finding, independent of ownership.** Recipients are currently recordable *only as N separate
accounts*. Neither interface can attach a recipient set to a claim event:

| Path | Recipients available? | Where they go |
|---|---|---|
| `disclose(discloser, recipients, source_record_id, turn)` | yes | into N `acquire()` calls; nowhere on any event |
| `claim(contributor, contributor_kind, text, turn, kind, span)` | **no parameter** | provenance is fixed at `{contributor, contributor_kind, span, turn}` (`engine.py:188`) |

Therefore J3 cannot be satisfied today **under either reading**: under A, `disclose` writes no claim; under B,
the caller can call `claim()` but has no way to record who was exposed. This is why J3 is **P** for a reason
that survives the ownership question entirely.

**Status of the proposal, stated precisely:**

| Question | Answer |
|---|---|
| Contractually established? | **No** — neither the author nor the recipients placement is settled |
| Merely plausible? | **Yes, under Reading A** — it is the shape §7 step 2 names, written by the operation rule 46 describes |
| Unsupported? | **No** — its class, content, truth status and access rule are all supported; only authorship and the recipient placement are open |

One-claim-per-disclosure is a **candidate**, not a resolution. My previous supplement called it the smallest
permissible resolution; that label was too strong and is withdrawn here.

---

## 4. What would settle the owners' question — and what remains

Three candidate settlements, none of them derivable from the current text:

| Option | What it requires | Cost |
|---|---|---|
| **A** — `disclose()` creates the event claim | Forward a recipient set onto the claim (a new provenance key), plus optional `span`/`kind` | A new provenance key; the claim becomes the discloser's private record by existing rules |
| **B** — the caller creates the claim, MEMORY-09 consumes it | `claim()` must accept recipients; `disclose()` must accept an existing event id | Same recipient key; plus an event-id parameter, and a check that the event names the same discloser |
| **C** — the event stays outside MEMORY-09 | Record that rule 46's first clause is satisfied by the conversation layer; MEMORY-09 records accounts only | No code change; rule 46/J3 would then be **layer-satisfied**, not MEMORY-09-satisfied, and the review would have to say so |

Under **A** and **B** the concrete delta is the same: **one provenance key carrying the recipient set**, which is
what J3 names and no current path can record. Under **C** no delta is needed, and the P classifications become a
statement about the boundary rather than about missing code.

**The unresolved item is a single question:** does the experiment's boundary mean MEMORY-09 must be able to
record *who was exposed*, or is that fact only ever the conversation layer's and never MEMORY-09's?

That question is yours. I am not resolving it, and I have not implemented any of A, B or C.

---

## 5. Classifications, preserved

| Row | Classification | Change from the consolidated review |
|---|---|---|
| **Rule 46** | **P** | **None.** The recording exists nowhere in MEMORY-09; whether it must is unresolved |
| **J3** | **P** | **None.** Independently of ownership, no path records the recipient set |
| Rule 10, 11, 12, 25, 47, 48 | I | None |
| J7, J8, H4, K5 | I | None |
| Rules 56 I · 1 P (46) | unchanged | — |
| Requirements 113 I · 3 P (J3, K9, K10) · 2 NA (B3, F7) | unchanged | — |

**No implementation.** No claim event, no `_events` registry, no replacement structure was written, and no
existing file was modified to accommodate either reading.

---

## 6. K9 — standing position, restated so it is not re-litigated

| Proposition | Status |
|---|---|
| A persisted write-time timestamp with offset is **not prohibited** by the snapshot contract (S-1…S-8, S-9…S-12) | Established — the prohibition in REV 10 §3.1 addresses response fields derived "from a clock" after `S` |
| A **retrieval-time** clock read is prohibited | Established — §3.1, with S-2's "no field computed from a read made after `S`" |
| The turn counter is an **ordinal**, not a timestamp | Established — `Record.created_turn` (`records.py:74`), `_promotions[...]["turn"]`; no clock, no offset |
| F7 conformance | **Not claimed.** F7 remains **NA**; no replay facility exists, and none was implied |
| K9 | **P** — unchanged |

**No timestamp field was added. The turn counter was not ratified as equivalent to a timestamp.** That decision
stays open with the owner, alongside the event-ownership question.

---

## 7. Preservation and verification

| Item | Value |
|---|---|
| This file | additive; nothing edited to accommodate it |
| Implementation | `authority.py f3445823…`, `engine.py 706d5f6c…`, `records.py f211e6ce…`, `errors.py 48125d7e…`, `retrieval.py cd2d380b…`, `__init__.py 06ccf8f6…` — unchanged |
| Review documents | consolidated review `9fb62d24…`, report `bd61da90…`, report supplement `f3484393…`, verbatim sources `4083c937…`, clarification `1ea28175…`, events-and-time supplement `2aaf61de…` — all unchanged |
| Contract and sources | REV 13 `5b745628…`, v9 `698254ca…`, map `b4ba33cc…`, ledger `da3e7239…`, acceptance `9a25ded0…` — unchanged |
| Evidence read for this supplement | v9 §1 table, §2 concepts, §3 class table, §7 steps 1–7, §8 provenance; rules 10, 11, 25, 46, 47, 48; r32 K-R11 as quoted in REV 2 §5.10; `claim`/`acquire`/`disclose` signatures and `acquire` provenance |
| Executed | **nothing** |
