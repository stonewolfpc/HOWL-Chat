# MEMORY-09 — Entry Uniformity and the §8.1 Discrepancy

## Two questions kept apart, and the missing premise the reconciliation needs

**Status:** analysis record. **Nothing adopted. Nothing selected.** Kept **outside** OD-b as an independent
per-kind semantic question. No schema, no canonical edit, no test, no ledger entry.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Question A | Does formation **always** create a `provisional` lifecycle condition? |
| Question B | Can a particular formation operation **satisfy an exit condition immediately**, making `provisional` unobservable in the resulting state? |
| Relationship | **Not equivalent** (EKR-03). B may explain §8.1 without answering A. |
| Effect on OD-b | **none** — entry ownership is already constrained (ECI-17/18); this record does not reopen it |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Why the two questions must not merge

**EKR-01 (Observed — the three statements).**

| # | Statement | Source | Tag |
|---|---|---|---|
| 1 | "evidence turn → provisional (status, class, attribution set)", then four evidence-kind exits | §11 `:600–612` | untagged design text |
| 2 | "Formation writes `scope=world`, `knowledge_class=experiential`, **`status=active`**, `established_turn=40`" | §8.1 `:404` | untagged worked case |
| 3 | formation may enter `provisional` and advance immediately for a particular kind | ECI-07 (this workstream) | **Inferred** |

**EKR-02 (Inferred — statement 3's standing).** Statement 3 is **not established**. It was offered as a
reconciliation, and a reconciliation is a hypothesis about intent until grounded in text. This record tests
whether §11's own per-kind rules can ground it — and reports what is still missing.

**EKR-03 (Inferred — the non-equivalence).** A and B differ in what they quantify over:

- **A** is about **every** formation: is `provisional` always the entry condition?
- **B** is about **some** formation operations: can an exit be satisfied by the same act that enters?

B can be true while A is false (formations of that kind never enter provisional at all), and A can be true while
B is false (entry is always provisional and always observable, even briefly). §8.1 is evidence about **B** only if
§11 already says something about evaluation timing — which §11 does not (EKR-10).

**EKR-04 (Inferred — why merging them would be costly).** Merging turns §8.1 into evidence about A, and then
either (i) A becomes "no, not always" — a change to §11's diagram on the strength of a worked example, or (ii) B
is dismissed because A holds — leaving §8.1 unexplained. Both outcomes settle a semantic question by inference
from a case study.

---

## 2. The per-kind rules, verbatim

**EKR-05 (Source-derived — §11 `:612`).**

> "Promotion per kind: episodic — on witnessing; preference — immediate for holder; belief — never promotable to
> fact; rumour — **stays provisional**, repetition is not corroboration; world fact — authority or corroboration;
> absence — attributed. Rejected: blanket three-source rule, promotion by retrieval, decay altering truth,
> deleting history."

**EKR-06 (Source-derived — the diagram's exit labels, `:605–609`).** `active` is labelled "(episodic, preference,
validated fact, belief)" — so `active` is the exit for four of the six kinds; `superseded`, `corrected` and
`retracted` are cross-kind exits available to any record.

**EKR-07 (Observed — what §11 does *not* say).** §11 nowhere states **when** an exit condition is evaluated:
at formation, on the arrival of new evidence, or continuously. The rules name the **condition** for each kind,
not the **evaluation point**.

---

## 3. Question A — is entry always provisional?

**EKR-08 (Source-derived — for).** The diagram has exactly one entry edge, from "evidence turn" to
`provisional`, and labels it with the three things formation writes ("status, class, attribution set"). No second
entry edge exists. On the face of the text, **A holds**.

**EKR-09 (Observed — against).** §8.1 step 1 is a formation that writes `active`. It is a worked case, not a
rule, and it carries no `[P]`/`[D]` tag; but it is the only place in the corpus where a formation's status write
is shown, and it does not write `provisional`.

**EKR-10 (Calculated — the asymmetry).** One entry edge (structural) versus one worked case (illustrative). If
the question were decided on **structural** evidence, A holds. §8.1 is the only counter-evidence and it is
illustrative. This record does not weigh them into a verdict; it records that **A is Source-derived on the
diagram and contradicted by one worked example** (EKR-09), which is why it is a genuine open question rather than
an oversight.

---

## 4. Question B — can an exit be satisfied at formation?

**EKR-11 (Inferred — the missing premise).** B requires a premise §11 does not supply: **that exit conditions
are evaluated at the formation act**, not only when subsequent evidence arrives. Call it **P-EVAL**. Without
P-EVAL, statement 1 and §8.1 cannot both be right; with P-EVAL they can. P-EVAL is **not stated anywhere in
r33** — it is inferred from the per-kind wording only.

**EKR-12 (Inferred — the cost of importing P-EVAL).** P-EVAL is not a neutral convenience. If exits are evaluated
at formation, formation is an operation that can both create and advance, which means the **authority** to form
and the **authority** to advance are exercised in one act. ND-1 (who may advance) is open, so P-EVAL would make
entry depend on an open authority question. That is a reason to ground P-EVAL in text, not to assume it.

**EKR-13 (Calculated — kind-by-kind availability, grounded in §11's own wording).**

| Kind | §11 exit condition | Available **at the formation act**? | Grounding |
|---|---|---|---|
| **episodic** | "on witnessing" | **Plausible** | witnessing is an acquisition route (`via=witnessed`, §6.2 `:345`); the act that forms may be the witnessing |
| **preference** | "immediate for holder" | **Plausible — strongest** | "immediate" is itself an evaluation-time word; it admits no interval |
| **world fact** | "authority **or** corroboration" | **Plausible via authority**, not via corroboration | authority is present at formation if the attribution set records it; corroboration needs later events |
| **absence** | "attributed" | **Plausible** | the attribution set is written at formation (`:604`) |
| **belief** | "never promotable to fact" | **Unestablished** | constrains the **target**, not the timing; whether belief enters `active` or stays provisional is unstated |
| **rumour** | "**stays** provisional" | **No** | the only kind with an explicit no-exit; repetition is not corroboration |

**EKR-14 (Calculated — the result of the audit).** For **four of six** kinds, §11's own exit condition is
expressible at the formation act; **one** (rumour) has no exit; **one** (belief) is silent. So immediate exit is
**not a special case invented to explain §8.1** — on §11's wording it is the *common* case, and `provisional` is
the exception.

**EKR-15 (Inferred — what that does and does not settle).** It makes statement 3 **plausible and kind-grounded**
for §8.1's case, and it identifies exactly which premise it needs (P-EVAL, EKR-11). It does **not** establish
statement 3, because P-EVAL is inferred, and it does **not** answer A.

**EKR-16 (Inferred — §8.1 does not discriminate between explanations).** §8.1's record is `scope=world` **and**
`knowledge_class=experiential` — it sits in the overlap of the episodic and world-fact kinds. Its `status=active`
is equally consistent with "episodic, promoted on witnessing" and "world fact, promoted on authority". So §8.1
**corroborates** statement 3 without distinguishing which kind's rule did the work.

---

## 5. §8.1's second problem — storage, not just value

**EKR-17 (Source-derived — K-R12 `:111`).** A memory's "current status, class and scope are a **projection of that
stream**, not mutable authoritative fields."

**EKR-18 (Observed — §8.1 writes a field).** §8.1 step 1 says formation **writes `status=active`**. That is a
field write, exactly the shape K-R12 forbids for status.

**EKR-19 (Inferred — the discrepancy is therefore double).**

1. **Value level** — §8.1 writes `active`; §11 enters `provisional` (EKR-01 rows 1–2).
2. **Storage level** — §8.1 writes `status` as a field; K-R12 makes status a projection (EKR-17/18).

**EKR-20 (Inferred — which is more serious).** Storage level. The value question could be resolved by P-EVAL
(EKR-11); the storage question cannot, because no evaluation timing makes a **mutable authoritative field**
consistent with K-R12. §8.1 appears to predate, or to ignore, the projection rule — an **Inferred** chronology,
not established (the revision log never mentions §8.1 or K-R12, per CLV-02).

**EKR-21 (Inferred — why this matters to OD-b without reopening it).** Under L2/L3, the lifecycle state is
**derived**. A literal §8.1 would require the derivation to accept a **stored** starting value, which is neither
a relation nor an act. So §8.1 read literally is not merely a different entry value; it is a different **storage
model** for status. That is a second, independent reason it cannot be used to answer A or B.

---

## 6. The consequence that is load-bearing for any fold

**EKR-22 (Inferred — the fold-default failure).** Under L2(c), a record with no lifecycle relation folds to
`provisional`. If a formation legitimately produces `active` **and never records an entry relation** (§8.1 read
literally), the fold derives `provisional` for an `active` record — a **wrong state produced by a correct
default**. L1 is unaffected (its state is asserted); L3 is affected the same way as L2.

**EKR-23 (Inferred — the requirement this creates, stated without a mechanism).** Either (i) every formation
records an entry fact carrying the **resulting** state, or (ii) P-EVAL is rejected and formation always enters
`provisional`, so the default is always right. The first is a representation requirement; the second is a
semantic ruling. **Neither is chosen here.**

**EKR-24 (Calculated — so the §8.1 question is load-bearing, not cosmetic).** It determines whether a
fold-default representation can be correct for records that never observe `provisional`. That is why the owner's
instruction to keep it independent is the right sequencing: OD-b can be decided on its own discriminators, but a
**fold-based** selection inherits EKR-22 as a downstream obligation.

---

## 7. What would settle it — evidence, not inference

**EKR-25 (Inferred — the four candidate resolutions, none adopted).**

| # | Resolution | What it would require | Effect on A / B |
|---|---|---|---|
| R1 | §8.1 is illustrative shorthand | a statement that worked cases abbreviate | A holds; B unproven |
| R2 | Formation may promote immediately (P-EVAL) | a rule stating exits are evaluated at formation | A holds; **B true** |
| R3 | Authored world facts are born `active` | a per-kind entry rule | **A false** for that kind |
| R4 | §8.1 predates K-R12 and is stale | revision history (none exists — CLV-01) | neither; §8.1 withdrawn |

**EKR-26 (Observed — the evidence that does not exist).** There is no revision note for §8.1 or K-R12 (CLV-02),
no `[P]`/`[D]` tag on §8.1's step, and no statement of evaluation timing. So none of R1–R4 can be selected from
the corpus as it stands.

**EKR-27 (Inferred — the cheapest discriminating evidence).** A single sentence in the per-kind rules stating
**when** an exit condition is evaluated would separate R2 from R3, and would leave R1/R4 as bookkeeping. That is
the narrowest thing that would move this question.

---

## 8. What this record does not do

**EKR-28 (Observed).** It does **not** adopt statement 3, does **not** reject §8.1, and does **not** answer A or
B.

**EKR-29 (Observed).** It does **not** reopen entry ownership: ECI-16/17 (the lifecycle-side requirement) depend
on K-R11, K-R12, K-R2/K-R14 and D4, none of which touches the entry *value*. Whether entry is always
`provisional` or sometimes `active`, the entry **fact** must still be lifecycle-side.

**EKR-30 (Observed).** It does **not** enter OD-b's discriminator set and does **not** affect C′, OD-g, ND-1, or
`active`/`established`. `active` vs `established` remains a **vocabulary** question (K-R12), distinct from this
**value-at-entry** question.

---

## 9. Limits

**EKR-31 (Observed).** Basis: r33 §6.2, §8.1, §11, K-R12, and the two OD-b records. **Not inspected:** any
implementation's formation path (`internal/backend/memory/*.go` beyond the types file); r32's bytes; whether
`formation_method` (`:343`) distinguishes formation routes in a way that bears on P-EVAL. The P-EVAL premise is
**Inferred** and is the weakest link in statement 3.

**End of record.**
