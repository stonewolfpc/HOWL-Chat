# MEMORY-09 — Owner Decision OD-a: Separate Event Structures (Option B)

**Status:** owner-decision record. **OD-a settled: Option B — separate event structures.** Documentation and
read-only analysis only. **No schema change, implementation, test, probe, benchmark, model call, canonical
edit or registry change is authorised by this record.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Decision | **OD-a = Option B**, separate event structures |
| Decision scope | **Event representation only** |
| Authority | Owner, stated in session and recorded here; as with every prior confirmation, **no persisted signature exists** |
| Canonical edits | **None.** r33, v9, the owner ledger and the conformance map are unchanged |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Exact owner wording — preserved verbatim

> For MEMORY-09 OD-a, my decision is Option B — separate event structures. This is an owner decision about
> representation only; it does not authorize schema changes or implementation.

> Acquisition events and memory-lifecycle/status events remain separate structures.
>
> Each structure owns its own event identity and occurrence-ordering responsibility, consistent with A2.
>
> Cross-stream references, where needed, must be typed rather than implying shared identity or chronology.
>
> Separate structures do not establish cross-stream chronological ordering.
>
> Cause, evidence, occurrence, and authority remain distinct responsibilities.

> This does not authorize schema changes, implementation, tests, canonical edits, or registry changes.
> OD-a is settled. OD-b remains a separate owner decision, and OD-g remains open. A2 and D-3b remain
> unchanged.

---

## 2. What is decided

**ODA-01 (Source-derived — owner).** **OD-a = Option B.** Two separately defined event structures are
retained, not one shared envelope with a stream discriminator (which was Option A).

**ODA-02 (Source-derived).** The two structures are **acquisition events** and **memory-lifecycle/status
events**.

**ODA-03 (Source-derived).** Each structure **owns its own event identity**. Identity is not shared, not
derived from a common sequence, and not transferred between structures.

**ODA-04 (Source-derived).** Each structure **owns its own occurrence-ordering responsibility**, explicitly
"**consistent with A2**" — a persisted, independent, universal occurrence order *within* that stream. A2 was
not reopened to reach this; it is applied.

**ODA-05 (Source-derived).** Cross-stream references, **where needed**, must be **typed** — they must state
what kind of relation they express. They may not imply shared identity; they may not imply shared chronology.

**ODA-06 (Source-derived).** **Separation does not establish cross-stream chronological ordering.** Option B
answers "how many structures", not "who came first across structures".

**ODA-07 (Source-derived).** Five responsibilities stay distinct: **cause, evidence, occurrence, authority**
— and, per ODA-03, **identity**. No field may perform two of them by implication.

**ODA-08 (Inferred).** ODA-07 is the decision's strongest constraint. It forbids the shortcut that would have
made Option A tempting: one `event_id` space, one `occurred_at`, one `cause`, one `evidence_span` serving
both streams and therefore silently conflating them.

---

## 3. What is not decided, and not authorised

**ODA-09 (Source-derived — owner).** Authorised by this decision: **nothing** beyond the representation
choice. No schema change, no implementation, no tests, no canonical edits, no registry changes.

**ODA-10 (Observed).** Still open, unchanged by OD-a:

| Item | State |
|---|---|
| **OD-b** — is `provisional` a truth/status state or a formation/review state? | **OPEN** — still gating |
| **OD-g** — cross-stream comparandum for turn-less events | **OPEN** |
| **DN-6a / DN-6b** | **OPEN** |
| **ND-1** | **OPEN**, widened |
| **A2** | **SETTLED** — unchanged |
| **D-3b** | **SETTLED** — unchanged |
| The **118** ratified requirements | **UNCHANGED** |

**ODA-11 (Observed).** OD-a is a **gating** decision: it was one of the pair identified as necessarily
preceding OD-c/d/e. Its resolution unblocks the *shape* question, not the *content* question — OD-b still
determines what the lifecycle structure must represent.

---

## 4. What Option B implies for the structures that exist today

**ODA-12 (Observed — one structure exists, one does not).** r33 §6 defines exactly two record structures:

| Structure | Status |
|---|---|
| **§6.1 `MemoryRecord`** (chat-owned) | exists — but this is a *state record*, not an event stream |
| **§6.2 `AcquisitionEvent`** (authoritative, append-only) | exists — 7 fields: `id`, `character_id`, `memory_id`, `turn`, `via`, `from_character_id` (nullable), `source_memory_id` (nullable), `confidence_in_source`, `evidence_span` |
| **memory-lifecycle/status event** | **does not exist in any schema** — 4 schemas mention `evidence_span`, none defines the structure |

**ODA-13 (Inferred).** Option B therefore lands on an **asymmetric** pair: one side is specified, one side
is absent. The decision tells us the missing side must be **its own structure** — it does not tell us its
fields. **That cannot be drafted before OD-b**, because a lifecycle event schema is a statement about what
lifecycle represents.

**ODA-14 (Inferred — a consequence worth naming).** Under Option A, the absent structure could have been
minimally specified as "an envelope with `stream=status`". Option B forecloses that shortcut: the lifecycle
structure must be defined on its own terms, which makes OD-b's answer more consequential, not less.

---

## 5. The architectural rationale — preserved, and labelled

**ODA-15 (Source-derived — owner, verbatim).**

> Separate event structures support independently movable world, scenario, and character state domains.
> Howl's state machine must be able to resolve and access the appropriate memories and events as these
> domains are switched and recombined, maintaining continuity without conflating their respective histories
> or lifecycle responsibilities.

**ODA-16 (Source-derived — owner's stated context).** The rationale is that world, scenario and character can
each move independently, "into any combination", and the state machine "has to have access to the memories
and events to maintain fluidity" across transitions.

**ODA-17 (Source-derived — owner's boundary on the rationale).** "This is rationale for the representation
decision, **not yet a specification** of how transitions, memory activation, or event lookup will be
implemented."

**ODA-18 (Inferred).** The rationale is **forward-looking**, and must be recorded as such: it is a statement
about Howl's eventual composability, not a derivation from observed behaviour of the assessed system. It is
preserved here as owner-stated rationale — the reason the choice was obvious — and is **not** a new
normative requirement on MEMORY-09.

**ODA-19 (Inferred — why the rationale favours B).** If three domains can move independently, a single shared
event lifecycle forces every domain to be described by the same state vocabulary and the same sequence — so
switching one domain would either rewrite others' histories or require a discriminator that means different
things per domain. Separate structures let each domain's events carry their own identity and order, and let
composition happen at **resolution** time rather than inside the event log.

---

## 6. What the rationale exposes in the current canon

These are observations, not requirements. They are recorded because the rationale surfaced them.

**ODA-20 (Observed — `scope` is a label, not a domain reference).** r33 §6.1 records `scope`; the values in
use are `world`, `scenario`, `character_public`, `character_private` (r33 `:69`). A repository-wide search
for `world_id` or `scenario_id` across all markdown returns **zero files**. So canon records **which class of
applicability** a memory has, and **not which world or scenario** it belongs to.

**ODA-21 (Inferred).** Domain-switching resolution therefore has no domain **identity** to resolve against in
the current record model. This is a missing-structure finding (**P-1** shape), not a defect in OD-a, and it is
**not** something OD-a authorises fixing.

**ODA-22 (Source-derived).** r33 **D1** (`:69`): `world`/`scenario`/`character_public`/`character_private`
are **scopes, not owners**; experiential memories belong to the originating **chat**. **K-R39**: ownership is
scope-independent. So composability must **not** be read as multi-owner memory: switching domains does not
change who owns the memory.

**ODA-23 (Source-derived — the rationale's capability is outside the assessed experiment).** v9 (`:228–229`)
lists **world/scenario switching** as outside the experiment, with "the active world and scenario held
constant". V9 rule 31 and **K18** keep promotion contained within its chat; rule 36 makes world scope an
**authority domain, not a repository**.

**ODA-24 (Inferred).** Consequence for how the rationale should be read: switching a domain resolves **this
chat's** records for that domain. It is not domain merging, not cross-domain import, and **K-R40** ("no
automatic transfer") forbids reading it as migration. Continuity across a switch is therefore a *resolution*
property, not a *transfer* property.

**ODA-25 (Observed — selection is already context-relative; status is not).** r33 `:280`: the eligibility
filter excludes "by ownership, scope, privacy, status, time, character knowledge, conflict state — **for the
named context**". So the architecture already resolves *what applies* per context.

**ODA-26 (Inferred — narrows a concern I would otherwise have raised).** Because selection is
context-relative and `status` is absolute, a domain switch changes **which records are selected**, not what
any record's status **is**. Option B does not require `status` to become context-relative, and OD-a does not
decide that it should. If a future need appears for context-relative status, that is a separate question.

---

## 7. Consistency with prior decisions

**ODA-27 (Inferred).** Option B is consistent with, and does not extend, the settled record:

| Prior decision | Interaction |
|---|---|
| **A2** (per-stream independent ordering) | B **implements** A2's shape: two streams, each with its own order |
| **D-3b** (statement-less events permitted) | B is what makes D-3b tractable — a statement-less status event lives in a structure that does not require an utterance |
| **DN-1** (prohibition covers the source record) | unaffected; B is about event representation |
| **DN-4** (forgetting = distinct op + event) | forgetting's **event** now has a home: the lifecycle/status structure. Which structure, and with what fields, is still open |
| **D-3a / SD-02 / ND-1** | untouched |
| **XT-03…XT-05** (A2 answers Q1 only) | B makes the two-sequence fact structural rather than a convention: no shared sequence exists to be misread |

**ODA-28 (Inferred).** A2's sufficiency claim is **not** strengthened by OD-a. B separates the two orders; it
does not supply the cross-stream comparandum that **OD-g** addresses. The two decisions are independent and
both remain as stated.

---

## 8. What is not established

**ODA-29 (Unestablished).** The lifecycle/status structure's fields, name, and whether it carries `turn` at
all — blocked on OD-b, and shaped by OD-g.

**ODA-30 (Unestablished).** Whether domain identity (a world or scenario reference) belongs in the record
model, in the application layer, or nowhere. Canon is silent (ODA-20); OD-a does not decide it, and it is
**not** raised here as a new owner decision — it is recorded so the question is not lost.

**ODA-31 (Unestablished).** How a transition resolves records at all (ODA-17's reserved boundary). No
transition, activation or lookup mechanism is specified, and none is authorised.

**ODA-32 (Observed).** No test, schema, probe or experiment was run to support this record. It records a
decision and what canon says around it.

**ODA-33 (Observed — method limit).** The `world_id`/`scenario_id` search was term-based across markdown
only; Go or Python sources were not searched for an equivalent domain identifier in this pass, and the
earlier finding that `chat_id` was **relocated** to `owner_chat_id` rather than removed is unaffected.

---

## 9. Record status

**ODA-34 (Observed).** OD-a is **SETTLED — Option B**. OD-b is the remaining gate. OD-g is open. DN-6a/DN-6b
open. ND-1 open and widened. A2 and D-3b unchanged. The 118 are unchanged. No canonical file was edited to
record this decision, deliberately: the decision lives in this dated, untracked record until the owner
authorises canonical capture.

**End of record.**
