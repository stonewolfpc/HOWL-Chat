# MEMORY-09 — Ordering Key Families: Direct Comparison

**Status:** comparative assessment. Read-only. **No selection, no schema change, no O-16 entry, no
implementation, no test, no probe, no model call, no corpus, no Git operation.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Question | Does a universal occurrence order satisfy the whole event model — including status events — better than a composite `(turn, sub)`? |
| Baseline | r33 `dfb28b6003ca356f` · harness `eval/mem1/` |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Four axes, not three

**OK-01 (Inferred).** The owner's three-way separation is adopted, with a fourth kept explicitly apart
because conflating it with ordering is what produced the `event_id` defect:

| Axis | Question | Answers today |
|---|---|---|
| **Evidence** | what supports this event? | `evidence_span` — nullable, never fabricated (D-3b) |
| **Cause** | what produced this event? | **no field**; `via` names the kind, not the cause |
| **Occurrence** | where does this event belong in the history? | **no field** |
| **Identity** | which event is this? | `id` — unique (K-R16), but currently positional and unstable |

**OK-02 (Source-derived).** Identity and ordering are **related but not identical**, as the owner said: a
stable id answers "which event", an ordering key answers "which came first". **Replacing `event_id` with a
stable id would fix identity and leave ordering untouched**, and vice versa.

---

## 2. The two families, defined precisely

**OK-03 (Inferred).** **Family A — universal occurrence order.** Every event carries a persisted ordinal,
assigned at write, from a single sequence per scope. Two variants must be distinguished, because the
criterion set does not decide between them:

| | Scope of the sequence | Property |
|---|---|---|
| **A1** | one global sequence across both streams | a single total order over everything |
| **A2** | one sequence **per stream** | a total order within acquisition, and within memory |

**OK-04 (Inferred).** **Family B — composite `(turn, sub)`.** The turn plus a sub-position that is
statement-derived when a statement exists and self-assigned otherwise.

---

## 3. Criterion by criterion

| # | Criterion | **A1** | **A2** | **B** |
|---|---|---|---|---|
| 1 | Total ordering across all relevant event types | yes, globally | yes, **per stream** | partial |
| 2 | Stable identity across reconstruction | **orthogonal** | **orthogonal** | **orthogonal** |
| 3 | No dependence on insertion order | yes | yes | conditional |
| 4 | No fabricated conversational positions | yes | yes | **no, or** |
| 5 | Consistent ordering, statement-bearing and -less | yes | yes | **no** |
| 6 | Multiple events within one turn | yes | yes | yes for distinct spans, **no** for one span |
| 7 | Chat scope recoverable via `owner_chat_id` | yes | yes | yes |
| 8 | No competing sources of ordering truth | yes | yes | **no** |

### 3.1 Criterion 1 — what "across all event types" can mean

**OK-05 (Source-derived).** K-R14 and §6.4: "**Two streams, two checks.** … Neither projection may consult
the other." Nothing in T2 requires comparing a retraction to a learning; the join is read-time and never
stored (K-R13).

**OK-06 (Inferred).** So **A2 is sufficient** and A1 is stronger than anything requires. A1 is not free:
a single global allocator makes every write depend on the other stream's high-water mark, which is a shared
dependency the two-stream design exists to avoid. **The criterion as worded can be read as demanding A1;
no requirement demands it.** That reading question is the owner's.

### 3.2 Criterion 2 — not discriminating

**OK-07 (Inferred).** Neither family supplies stable identity. Both require ids **stored, never computed**.
The harness's `f"ev-{len(self._events) + 1:05d}"` fails this for a reason unrelated to either key family:
it is derived from list length at append time. **A stable id is a prerequisite for both families, not an
advantage of either.**

### 3.3 Criteria 4 and 8 — the decisive pair

**OK-08 (Inferred).** Family B cannot satisfy criteria 4 and 8 **at the same time**, and the reason is
structural rather than stylistic:

| B's choice | Then |
|---|---|
| one numeric space for both statement-derived and assigned sub-positions | an assigned value is **indistinguishable from a statement position** in the log — the "inferred statement position" D-3b forbids (criterion 4) |
| tag the two kinds inside the field | the field now has **two generators**, so two events' order cannot be compared without a cross-kind rule — two ordering authorities (criterion 8) |

**OK-09 (Inferred).** There is no third option: the ambiguity and the dual authority are the same problem
viewed from either end. This is the strongest single result in the comparison.

### 3.4 Criterion 5 — B's first component is also unsafe

**OK-10 (Inferred).** B's key is `(turn, sub)`, and a statement-less event may have **no turn either** — a
model-initiated loss or a lifecycle action need not occur *in* a turn. §6.2 requires `turn`, but that
requirement was written for events that follow utterances. So B fails on **both** components, not just the
second. **A does not presuppose a turn for ordering**, which is why it survives where B does not.

**OK-11 (Unestablished).** Whether statement-less events nonetheless receive the turn in which the system
acted. If they do, B's first component survives; the second still does not (OK-08).

### 3.5 Criterion 6 — the known granularity limit

**OK-12 (Observed).** B orders distinct spans correctly and ties when one span yields two events for one
pair. The deterministic extractor cannot produce that today (ten rules, ten distinct predicates, one
`search` per span); model-assisted formation — RQ-MF-01, **open and untested** — can. A has no such
limit: each event takes its own ordinal.

### 3.6 Criterion 7 — satisfied, with one caveat

**OK-13 (Source-derived).** Chat scope is recoverable for any event: `event → memory_id →
MemoryRecord.owner_chat_id`. True for acquisition, forgetting and retraction alike, because each concerns a
memory.

**OK-14 (Inferred).** The caveat is timing, not correctness: allocating a **per-chat** ordinal requires the
join **at write**, which is available (the memory is in hand). A direct carrier would be more robust for
reasoning after deletion, but K-R43 removes records with the chat, so no orphan case arises.

---

## 4. Does the universal mechanism cover status events?

This was the owner's condition for authorising anything. **The answer is no — not yet, and the reason is
larger than ordering.**

**OK-15 (Observed).** §5's dimension table states:

> **Memory lifecycle history** | append-only **memory events** | authoritative for status/scope/class;
> status is projected, never a mutable field (K-R12)

**OK-16 (Observed — the gap).** **§6 defines no memory-event schema.** §6.1 is `MemoryRecord`, §6.2 is
`AcquisitionEvent`, §6.3 is `KnowledgeEdge`, §6.4 is consistency rules. The memory stream is **referenced
as existing and authoritative** and is **never specified**.

**OK-17 (Inferred).** Consequences, in order of importance:

1. **No ordering mechanism can be "carried by every event" of a stream that has no events to carry it.**
   Universal ordering is necessary but **not yet applicable** to status.
2. **ND-1 widens.** It is no longer only "does the status stream need an ordering rule?" — the stream has
   **no schema at all**, so ordering cannot even be specified for it. That is a bigger gap than the one
   recorded.
3. **This is the second undefined-but-referenced structure found**, after `evidence_span` (ES-02). The
   pattern is now established: r32's §5 declares structures that §6 does not define.

**OK-18 (Inferred).** So a decision for A2 would be **correct for acquisition and forgetting** and
**silent for status** — which is exactly the boundary the owner asked not to cross. The honest position is
that A satisfies the acquisition stream completely, and the status stream cannot be assessed until it has
a schema.

---

## 5. The coverage obligation implied by SD-02

**OK-19 (Source-derived, stated as an obligation).** D-3b implies one obligation, in three clauses:

> **CO-1 (cause)** — every acquisition-family event records the cause that produced it.
> **CO-2 (absence)** — an event with no source statement records that absence as absence, never
> substituting evidence or a position.
> **CO-3 (identity)** — such an event is independently identifiable as an occurrence.

**OK-20 (Calculated — coverage of each clause in the 118).**

| Clause | Nearest row | Covered? |
|---|---|---|
| CO-1 | none. **K-R90** answers "how did they learn this" from events — learning only, not cause | **no** |
| CO-2 | none. **I15** ("absence is never evidence about absence") is a different claim | **no** |
| CO-3 | **K-R16** requires unique ids | **partly** — uniqueness, not stability across reconstruction |

**OK-21 (Inferred — species).** This is an **in-scope obligation with no ratified test** — the **NEW-1
species**, not the D-4 species (it is not outside MEMORY-09's authority; it is inside it and untested).
Recorded as a coverage obligation. **No row is added**, the set stays at **118**, and the map is unchanged.

---

## 6. Result

**OK-22 (Inferred).** On the criteria as given:

- **Family B fails 4 ∧ 8 (structurally, OK-08), fails 5 (OK-10) and fails 6 for model-assisted formation
  (OK-12).**
- **Family A satisfies 3, 4, 5, 6, 7, 8** on the acquisition stream; **criterion 1 depends on whether
  "across all event types" means globally (A1) or per stream (A2)**, and **criterion 2 is orthogonal**.
- **Neither family covers status events**, because there is no status-event schema (OK-16).

**OK-23 (Inferred).** So Family A is the only family that satisfies the set — with two caveats that are
decisions, not findings: A2 rather than A1 unless a global order is wanted; and status ordering is
undefined territory regardless of the choice.

**OK-24 (Observed).** This record selects nothing. The owner's lean toward a universal order is consistent
with the analysis, and the decision remains theirs.

---

## 7. What a choice would still leave open

**OK-25 (Inferred).** None of these are settled by picking a family:

| Item | Why it survives the choice |
|---|---|
| the **cause descriptor** (CO-1) | a shape for "what produced this" — independent of how events are ordered |
| **identity** (criterion 2) | stored ids, still needed |
| **memory-event schema** (OK-16) | the status stream has nothing to attach a key to |
| **allocation scope** (A1 vs A2) | per-chat, per-stream, or global |
| **DN-6a / DN-6b** | who may cause forgetting or retraction — untouched throughout |
| **ND-1** | widened, not closed |

---

## 8. Limits

**OK-26 (Observed).** Static reading only. No fold was run, no key was implemented, and the memory stream's
absence of a schema is established from §§5–6 of r33 — a document that may be silent by design about a
stream it expects elsewhere. If a memory-event schema exists in another artifact, OK-16 is an evidence gap
rather than a finding, and this record should be corrected.

**End of comparison.**
