# MEMORY-09 — `evidence_span` Contract Investigation

**Status:** read-only investigation. No schema change, no implementation, no test, no probe, no model
call, no corpus, no Git operation. **No remedy selected.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Question | What does r32 intend `evidence_span` to identify, and can it carry the ordering O-G needs? |
| Baseline | r33 `dfb28b6003ca356f` · `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` · `eval/mem1/{model,form}.py` |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The three problems, kept apart

The owner's framing is adopted as this report's spine, because each problem needs a different remedy and
conflating them is how a documentation gap becomes a schema change:

| Problem | What it looks like | Remedy it needs |
|---|---|---|
| **P-1 Missing structure** | no field can carry the order | add a field |
| **P-2 Inadequate definition** | a field exists and is used, but what it *is* is unstated | define it |
| **P-3 Over-narrow requirement** | the requirement names a mechanism where it intends a property | reword the requirement |

**ES-01 (Observed).** P-1, P-2 and P-3 all fit something here. Which ones are *real* is what this report
establishes — and the answer is not the one the first investigation assumed.

---

## 2. What `evidence_span` actually is

**ES-02 (Calculated — exhaustive search).** `evidence_span` occurs in **four schemas** and in **zero
definitions**. Every occurrence is a field name in a block; no document says what it contains. It appears
in no prose, no rule, and no example value anywhere in the project.

| Location | Schema | Sibling fields |
|---|---|---|
| `DESIGN-MEMORY-CONTINUITY-2026-11-11.md:107` | **`Event`** | `event_id`, `continuity_id`, **`chat_id`**, **`turn_id`**, **`source_message_id`**, `event_type`, `participants[]`, `summary`, **`evidence_span`**, `observed_at`, `effective_at`, `origin`, `confidence`, `status` |
| `DESIGN-MEMORY-CONTINUITY-2026-11-11.md:169` | **`StateProposal.changes[]`** | `field`, `operation`, `value`, **`evidence_span`**, `confidence` |
| `MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md:240` | **`Fact`** (ground-truth ledger) | `fact_id`, `subject_character_id`, `predicate`, `object`, `valid_from`/`valid_to`, `status`, `source_event_id`, **`evidence_span`**, `visibility`, `derived_from` |
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md:349` | **`AcquisitionEvent`** (r32) | the §6.2 list — no `chat_id`, no `source_message_id` |

**ES-03 (Source-derived — the decisive inference).** In the `Event` schema, `evidence_span` sits
**beside** `source_message_id` and `turn_id`. It therefore **cannot be the message identifier**: the
design already has one. `evidence_span` is a **span within** a message — a location or segment, not the
message's identity.

**ES-04 (Unestablished).** Whether that span is an **identifier** (`span_id`) or a **text range**
(offsets / quotation) is stated nowhere. The harness uses an identifier (`source_span`, `sid`); the
contract names a span and never resolves the ambiguity. **This is P-2, and it is real.**

**ES-05 (Observed — a regression, not a gap).** The v1 design's `Event` carries **`chat_id`** and
**`source_message_id`**. r32's `AcquisitionEvent` carries **neither**. ARCH-03 F6 records "no `seq` (and
no chat field)" as if the chat field had never been designed. It *had* been — it was dropped between the
design layer and the contract. That changes the diagnosis: the per-chat counter's namespace existed and
was removed, which is a simplification of an existing structure, not an absence of one.

**ES-06 (Observed — a second drop).** The v1 `Event.status` enum includes **`forgotten`**, i.e. the
earlier design treated forgetting as a *status*. K-R19 forbids exactly that ("no edge with a 'revoked'
flag"). DN-4 = YES settles the mechanism in favour of K-R19; the design document is the older layer and
was not edited. Recorded so the conflict is not mistaken for a live contradiction.

---

## 3. Can a span carry an ordinal? — yes, under either reading

**ES-07 (Inferred).** This is the finding that changes the PH1A-05 picture. Both candidate readings of
`evidence_span` yield a **total order within a turn**, with no new counter:

| If `evidence_span` is… | The ordinal is… | Persisted? | Stable across rebuild? |
|---|---|---|---|
| an identifier | an intra-turn ordinal embedded in, or derivable from, the id | yes, on the event | if ids are stable (ES-13) |
| a text range | the span's **start offset** within the message | yes, on the event | **yes** — the message text does not change |

**ES-08 (Source-derived).** Either way the key is `(turn, span position)`, and both components are already
recorded. The message text is immutable in a log, so a start offset is content-derived order — **not**
insertion order, and therefore not the thing K-R15 forbids.

**ES-09 (Inferred).** So P-1 (**missing structure**) is probably **not real**. A carrier exists in the
message layer and was carried onto the event in the harness. What is missing is a **definition**, not a
field.

**ES-10 (Calculated — the tie that remains).** A span position distinguishes **distinct spans**. It
cannot distinguish **two events born of the same span**. That is §4, and it is the only granularity O-G
cannot supply by itself.

---

## 4. Event granularity: can one statement produce two events for one pair?

The owner's caution was right: this must be established, not assumed. It has two separate answers.

### 4.1 Normatively — permitted

**ES-11 (Source-derived).** K-R18: "A learn event is valid at any time. A learn while an episode is open
is a *continuation* (the character learns again / confirms), not a new episode." K-R23: "A weaker later
event … is appended to the log". Nothing in T2 restricts how many learn events one statement may
produce — the architecture explicitly contemplates several events in one open episode.

**ES-12 (Inferred).** A design that **requires** at most one event per pair per span would not be
implementing K-R15; it would be **adding a constraint the architecture does not have**. That is the
distinction the owner asked to preserve, and it is why "one event per pair per statement" cannot be
silently adopted as a precondition.

### 4.2 Empirically — one, in the deterministic extractor, by construction

**ES-13 (Observed — `form.py`).** The extractor has ten rules, each with a **distinct predicate**, each
applied once per span by `pattern.search`; `_admit` records `formed_spans` and counts duplicates. So at
most one candidate per (rule, span), and distinct predicates mean at most one event per (pair, span).
**The condition holds today — as an artefact of the rule table.**

**ES-14 (Observed).** The module docstring states the limit plainly: "`RQ-MF-01` (model-assisted
formation) stays open and untested: nothing here calls a model". Model-assisted formation is exactly the
path that would break ES-13's guarantee, because a model can emit two claims about one fact from one
sentence.

**ES-15 (Inferred).** Therefore: the one-event-per-span property is **true of the current instrument and
not guaranteed by the contract**. O-G is sufficient *in fact* for the deterministic extractor and
*not established* in general.

---

## 5. Stability across reconstruction

**ES-16 (Observed — harness).** `event_id = f"ev-{len(self._events) + 1:05d}"` is derived from append
position. Two consequences: it is **insertion order in a field named like an identity**, and a rebuilt
log is **renumbered**, so it cannot be compared to the original (K-R2/K-R14).

**ES-17 (Inferred).** An order carried by `event_id` would therefore fail the one property K-R15 exists
to protect — reproducibility. An order carried by a **span position** does not have this defect *if* the
span reference is stored and not recomputed.

**ES-18 (Unestablished).** Whether span **identifiers** are stable across a reload depends on how they are
constructed (the `EvidenceTurn.span_ids` source was not inspected). A **start offset** has no such
dependency — it is a property of immutable text. That is a mild argument for treating the span as a
range rather than an opaque id.

---

## 6. Consistency across event kinds — this fails today

The owner asked whether the identifier applies consistently to acquisition, forgetting and retraction
events.

| Event kind | Which span it records | Verdict |
|---|---|---|
| **Acquisition** | `source_span=span_id` passed from the forming span (`form.py:232` → `_append_event`) | **correct** |
| **Retraction record** | `retract(..., source_span, ...)` → the retraction's own span | **correct** |
| **Forgetting event** | `revoke(record_id, turn, reason)` has **no span parameter**; it builds its event with `source_span=rec.source_span` — **the acquisition statement's span** | **wrong** |

**ES-19 (Observed — `model.py:303-311`).** The revocation event is stamped with the provenance of the
event that *created* the record, not of the statement that caused the loss.

**ES-20 (Inferred — the consequence is precise).** Under O-G this is not cosmetic:

1. the forgetting event **ties** with the learn it is supposed to close — the exact ambiguity O-G was
   introduced to remove, reappearing at the one place it matters;
2. it sorts **at the moment of learning**, so a forget in turn 40 from a fact learned in turn 12 would
   order as if it happened in turn 12's span — the fold would reconstruct the wrong history;
3. K-R90 ("how did they learn this?" is answered from events) is corrupted: the forgetting event claims
   the wrong cause.

**ES-21 (Inferred).** This is the strongest argument yet that the **carrier** must be audited per event
kind before it is adopted. A correct carrier, wired inconsistently, reproduces the original defect with
more machinery.

---

## 7. K-R15's wording: mechanism or property?

**ES-22 (Source-derived).** The purpose is stated, separately from the rule. R1c's review concern:
"`revoked_turn` must be **entirely reproducible** from the acquisition stream, including ordering and
multiple learn/forget cycles." §6.4: "replaying the same ordered log produces the same edge, **on any
machine, in any insertion order**."

**ES-23 (Source-derived).** The rule itself is mechanism-level: "totally ordered by `(turn, seq)`, where
`seq` is a monotonic per-chat counter … Events sharing a turn are resolved by `seq` alone."

**ES-24 (Inferred).** So the requirement's **purpose is determinism and reproducibility**; the counter is
**one means** to it. That supports the functional reading (SM-08 (b)) — but the text names the mechanism,
so a contract that orders by span position instead is **not literally conformant** without a wording
change. The honest statement is: *the goal is met; the sentence is narrower than the goal.*

**ES-25 (Inferred).** That is **P-3 (over-narrow requirement)** — but P-3 here is mild. It is a one-clause
rewording, not a redesign, and it should be decided *after* ES-03/ES-04 are settled, not before.

---

## 8. The three problems, resolved as far as the evidence allows

| Problem | Verdict | Basis |
|---|---|---|
| **P-1 Missing structure** | **Probably not real** | the message layer carries a span, and the event carries a reference to it (ES-07/ES-09) |
| **P-2 Inadequate definition** | **Real, and the blocker** | four schemas use `evidence_span`; none says whether it is an id or a range (ES-02/ES-04) |
| **P-3 Over-narrow requirement** | **Real but mild** | K-R15 names a counter where its stated purpose is determinism (ES-22…ES-25) |
| **P-4 Inconsistent carrier use** | **Real, discovered here** | the forgetting event records the wrong span (ES-19/ES-20) |

**ES-26 (Inferred).** The owner's wager was correct: **the cheapest path is to define the existing
structure, not to add machinery.** P-2 costs one sentence. P-3 costs one clause. P-1 — the thing the
first investigation proposed — is likely unnecessary.

**ES-27 (Superseded).** The prior report's option table treated **O-A (`seq`)** as the reference remedy
and O-G as the cheap alternative. The ordering is now reversed: **O-A addresses a problem that probably
does not exist**, and O-G addresses P-2 + P-3 at lower cost — provided P-4 is fixed.

---

## 9. What this changes about O-G

**ES-28 (Inferred).** O-G is **necessary and not sufficient**, for two reasons that are now concrete
rather than hypothetical:

1. **Intra-span granularity** (ES-10/ES-15): two events from one span tie. Today's extractor cannot
   produce that; a model-assisted one can, and RQ-MF-01 is open.
2. **Forgetting provenance** (ES-19/ES-20): the carrier is wired wrong for the event class DN-4
   created, and O-G would inherit the bug.

**ES-29 (Inferred).** Neither reason is an argument for a counter. Both are arguments for **finishing the
definition and auditing the wiring** — which is what the owner proposed, and what this report confirms.

---

## 10. The three questions, as they now stand

| Question | Status after this investigation |
|---|---|
| **1. What does `evidence_span` identify?** | **Half answered.** It is a span *within* a message, not the message (ES-03). Whether id or range is **still unstated** (ES-04) — and the answer decides ES-18, not whether a carrier exists |
| **2. Can one statement produce two events for one pair?** | **Answered in two parts.** Normatively yes, and unrestricted (ES-11/ES-12). Empirically no, by rule-table construction (ES-13) — a guarantee the open model-assisted path would remove (ES-14/ES-15) |
| **3. Does K-R15 mandate a counter?** | **Evidence favours "no".** Its stated purpose is reproducibility (ES-22); the counter is a means (ES-23). But the text is mechanism-level, so ordering by span requires a one-clause rewording (ES-24/ES-25) |

**ES-30 (Observed).** No remedy is selected. No schema change is proposed. P-4 is an implementation
defect, not a contract gap, and fixing it is not authorised here.

---

## 11. Limits

**ES-31 (Observed).** Static reading only. `EvidenceTurn.span_ids` construction was **not** inspected
(ES-18 stays open). The v1 design document is an earlier layer and is **not** a normative source for
r32 — it is used here only as evidence of intent (ES-03/ES-05).

**End of investigation.**
