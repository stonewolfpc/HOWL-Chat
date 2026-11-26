# MEMORY-09 — Entry Condition and Cross-Structure Integration

## How L1/L2/L3 establish the initial `provisional` under OD-a's separate structures

**Status:** analysis record. **No schema invented. No selection made. OD-g not resolved and not widened.**
Read-only: no canonical edit, schema change, implementation, test, probe, model call, registry change or
ledger entry.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Question | Under OD-a's separate acquisition and lifecycle structures, how does each of L1/L2/L3 establish the initial `provisional` produced by formation — especially when a record stays provisional permanently with no exit event? |
| Answers | cross-stream dependency · ordering assumption · fold requirement, per option |
| Decides | **nothing**; OD-b remains the owner's |
| OD-g | **untouched** — no comparandum proposed, none needed |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The constraints already in force

**ECI-01 (Observed — OD-a, owner wording preserved).** "Acquisition events and memory-lifecycle/status events
remain separate structures." "Each structure owns its own event identity and occurrence-ordering
responsibility, consistent with A2." "Cross-stream references, where needed, must be **typed** rather than
implying shared identity or chronology." "Separate structures do not establish cross-stream chronological
ordering."

**ECI-02 (Observed — the stream rules that bound any answer).**

| Rule | Text | Line |
|---|---|---|
| **K-R2** | "A projection that cannot be reconstructed from its events is an **integrity failure**, not a policy question." | `:93` |
| **K-R11** | The edge fold depends **only** on acquisition-family events for (character, memory). "Nothing else may influence it." | `:110` |
| **K-R12** | Memory lifecycle is its own authoritative stream; a memory's current status, class and scope are a **projection of that stream**, not mutable authoritative fields. | `:111` |
| **K-R13** | Retraction knowledge is a **read-time join** of edge × memory projection; never stored on the edge. | `:112` |
| **K-R14** | "Each stream is rebuilt and checked independently; a mismatch in either is an integrity failure. **The cross-stream join has no stored form**, so it cannot drift." | `:113` |
| **D4** | Provisional records are retrievable, attributed, **never promoted by retrieval**. | `:72` |
| **I15** | Absence of information is never information about absence. | `:301` |
| **I18** | The system may not use absence of evidence as evidence that nothing happened; each subsystem must name the evidence that would reveal its own silent omission. | `:306` |

**ECI-03 (Calculated — the contradiction inside the question).** The question as posed contains an assumption
that canon does not support: that the initial state comes **from** the acquisition structure. K-R11 confines the
edge fold to acquisition events; K-R12 makes status a projection of the **lifecycle** stream; K-R2 and K-R14
require each stream to be **independently reconstructable**. Composed, these forbid derivation of lifecycle
status from acquisition facts. The dependency the question anticipates is therefore **not a design gap but a
prohibited shape** — §4 states it exactly.

---

## 2. Three statements about what formation writes

**ECI-04 (Source-derived — §11 `:604`).** "evidence turn → **provisional** (status, class, attribution set)".
Formation writes the entry state and, with it, class and attribution.

**ECI-05 (Source-derived — §8.1 step 1, `:404`, untagged).** "Formation writes `scope=world`,
`knowledge_class=experiential`, **`status=active`**, `established_turn=40`."

**ECI-06 (Observed — the collision).** §11 says formation writes `provisional`. §8.1 says formation writes
`active`. Both are untagged design text in the same document. Neither is marked `[P]`; neither is a proposal.

**ECI-07 (Inferred — the charitable reconciliation, offered not adopted).** §11's exit conditions are **evidence
kinds**, and for world facts the condition is "authority or corroboration" (`:612`). Authority can be present
**at formation**. So formation may write `provisional` and advance in the same act. Under that reading §8.1 is
not contradicting §11; it records the **post-advance** value because the exit condition was satisfied at entry.
The two statements then agree on mechanism and differ only in whether the intermediate state is materialised.

**ECI-08 (Inferred — the alternative reading).** Authored world facts are **born `active`**, and §11's diagram
describes only the acquired/claimed path. This reading is also consistent with §8.1 and with §11's per-kind
rules; it differs from ECI-07 in whether entry is uniform across kinds.

**ECI-09 (Unestablished — which reading holds).** Canon does not say. This record does not choose, and the
choice is **not** OD-b: it is a question about whether entry is uniform, which belongs with the per-kind rules.

**ECI-10 (Inferred — the consequence that matters here).** Under **either** reading, a record's entry state is
**written by formation** and is **observable as `provisional` only while no exit condition holds**. The rumour is
the case where no exit condition **ever** holds, which is why it is the only case where `provisional` is
observable indefinitely. §8.1's `active` does not weaken the entry requirement; it shows entry and promotion
coinciding for one kind.

**ECI-11 (Observed — formation is a write, not an event, in every current text).** §11 and §8.1 both describe
formation **writing record fields** (`status`, `class`, `scope`, `established_turn`, attribution set). Neither
describes appending a **lifecycle event**. The acquisition side appends events (§6.2); the lifecycle side has no
schema at all (ODA-12). So "the initial condition" currently has **no event representation on either side**.

**ECI-12 (Observed — K-R17 confirms formation writes are the failure surface).** `:137`: a silent no-op "hides
formation bugs" — formation is already treated as a place where silent failure is a real risk, and K-R17's
remedy is *rejection at write*, not later reconciliation.

---

## 3. Four kinds of cross-structure dependency

**ECI-13 (Inferred — the distinction the question needs).** "Cross-stream dependency" is not one thing. Four
kinds behave differently, and only one is forbidden:

| Kind | What crosses | Status |
|---|---|---|
| **Content** | values written at formation (turn, class, attribution) originate acquisition-side | **Permitted** if copied at formation, not read live |
| **Status-conferral** | lifecycle status derived from acquisition facts | **Forbidden** — K-R11 + K-R12 |
| **Ordering** | comparing a lifecycle position against an acquisition position | **OD-g** — not introduced here, and not required |
| **Existence/closure** | the lifecycle stream's default assumes a record exists and that "no event" is complete | **Required to be addressed** — I15/I18 |

**ECI-14 (Inferred).** Distinguishing these dissolves most of the apparent problem. The question "where does
the initial state come from" mixes **content** (permitted) with **status-conferral** (forbidden). The answer is
that the *values* may cross and the *status* may not.

**ECI-15 (Calculated — what "typed cross-stream reference" buys).** OD-a permits typed references. A reference
from a lifecycle entry fact to the formation's acquisition event is **typed** (it says "entry grounded in this
acquisition") and does **not** imply shared identity or chronology. It is therefore within OD-a. It is also
**not** status-conferral: the reference records provenance, the status is written by the entry fact.

---

## 4. The prohibition, stated exactly

**ECI-16 (Source-derived + Inferred).** The shape canon forbids is: *"a record has acquisition events, therefore
it is provisional."* Reasons, each independent:

1. **K-R11** confines the edge fold to acquisition events; status is not an edge property.
2. **K-R12** makes status a projection of the **lifecycle** stream — so an acquisition-only derivation has no
   stream to project from.
3. **K-R2/K-R14** require each stream to be reconstructable **from its own events**; a lifecycle projection that
   requires acquisition input fails the independent-rebuild test.
4. **D4** separates retrievability from promotion; presence of an acquisition grants no status.

**ECI-17 (Inferred — the consequence).** The initial `provisional` **must be representable on the lifecycle
side**. It cannot be inferred from the acquisition structure. This holds for **all three** options, and it is
the single most important result of this record.

**ECI-18 (Inferred — and it resolves the owner's worry).** Because the entry fact is lifecycle-side, there is
**no required cross-stream dependency and no cross-stream ordering assumption** at entry. The integration
problem the question anticipates appears only in a variant that derives the default from acquisition presence —
which ECI-16 forbids. The dependency is not an L2 defect; it is a trap that L2 makes visible and that **all
three options must avoid identically**.

---

## 5. What each option does at the entry boundary

### L1 — state-naming

**ECI-19 (Inferred).** Entry is a **state assertion**: a lifecycle event names `provisional` (or the record's
creation carries it). Reconstructable from its own stream. No fold needed for entry — a state is read, not
derived.

**ECI-20 (Calculated — cost).** L1 needs a **record-creation** fact to exist, or the assertion has nothing to
attach to. If creation is not itself representable on the lifecycle stream, L1 has an existence dependency
(§3, kind 4) that it cannot close by naming states.

**ECI-21 (Inferred — benefit).** Under L1, `provisional` is asserted, so "still provisional" cannot be confused
with "transition lost". This is the one option where the permanent-provisional case needs **no closure
argument** at all.

### L2 — transition-only

**ECI-22 (Inferred — the sharp problem, stated plainly).** A transition-only representation has **nothing to
relate** at entry: `entered` needs a source, and at formation there is no prior state. So L2 must either (i)
record an **entry relation with a null/absent source** — a distinct relation kind, not a normal transition — or
(ii) treat entry as the record's creation, with the fold defaulting to `provisional` when the stream holds no
relevant relation.

**ECI-23 (Inferred — both are K-R12-consistent, and they are the a/b readings in miniature).** (i) is K-R12a's
shape: entry is on the stream. (ii) is K-R12c's shape: creation is the entry fact, and `established` is what a
transition produces. Either way the entry fact is lifecycle-side, satisfying ECI-17.

**ECI-24 (Inferred — permanent provisional).** Under either shape, a rumour is "entry recorded, no onward
relation". The state is correct. **But** the fold must distinguish:

- **complete, no successor** — the rumour: nothing further was ever warranted;
- **incomplete, successor expected** — a promotion that should exist but was not recorded.

**ECI-25 (Inferred — what that distinction requires).** A **closure** fact that is verifiable, not assumed:
either an explicit terminal marker, or a per-stream completeness guarantee that "no relation exists" is a
checked statement rather than an unobserved one. **K-R14 supplies that guarantee only within a stream** —
"each stream is rebuilt and checked independently" — and it does not extend across streams, which is why the
guarantee must be a lifecycle-stream property.

**ECI-26 (Inferred — the I15/I18 exposure, made concrete).** Without closure, the fold reports a state derived
from **absence of a relation**, which is exactly what I15 forbids and what I18 requires a subsystem to be able
to detect about itself. L2 is therefore **conditional on closure** — not defective, but incomplete until the
closure fact is chosen.

**ECI-27 (Inferred — phantom provisional).** A second failure mode specific to (ii): if formation's record write
succeeds but its entry fact does not, the fold's default reports a `provisional` record that **was never
formed** — a record whose existence is an artefact of the default. L1 does not have this failure mode, because
an absent assertion yields no state rather than a defaulted one.

### L3 — act-only

**ECI-28 (Inferred — the act-placement dilemma).** L3 derives state from acts. At entry the relevant act is
**formation**, whose evidence is an acquisition-side fact (`via`, `turn`, attribution). L3 must place that act
either:

- **lifecycle-side** — an entry act of its own. Then the lifecycle stream has an entry event, which is L2's
  shape with a wider payload; or
- **acquisition-side** — derive status from the acquisition event. Then status is conferred by acquisition, which
  **ECI-16 forbids**.

**ECI-29 (Calculated — L3 has no third position).** There is no other place for the entry act to live. So L3 is
either L2-with-payload on the lifecycle side, or canon-violating on the acquisition side. It is not a
structurally distinct integration.

**ECI-30 (Inferred — L3's absence exposure is worst at entry).** If entry is an act-only fact, a record whose
entry act is missing has **no derivable state at all**, and any default supplied to fix that is precisely the
"absence as evidence" move. L3 needs the same closure argument as L2, plus a total derivation to avoid an
undefined state.

---

## 6. The result: at entry, the options converge

**ECI-31 (Calculated — the constraint matrix).**

| | Entry fact lives | Reconstructable from own stream (K-R2/K-R14) | Status-conferral from acquisition | Needs closure for rumour | Cross-stream ordering |
|---|---|---|---|---|---|
| **L1** | lifecycle (state assertion) | **yes** | no | **no** — state is asserted | none |
| **L2 (a)** | lifecycle (entry relation) | yes | no | yes — closure or terminal marker | none |
| **L2 (c)** | lifecycle (creation = entry) | yes | no | yes — default + closure | none |
| **L3** | lifecycle (entry act) | yes | no | yes + total derivation | none |
| **any of them, acquisition-derived** | acquisition | **no** | **yes** | n/a | implied | 

**ECI-32 (Inferred — the finding).** **At the entry boundary, L1/L2/L3 collapse into one requirement:** a
lifecycle-side fact recording that the record exists in `provisional`, with values copied from formation. They
diverge only for **subsequent** transitions. OD-b's choice is therefore mostly about *post-entry* transitions,
not about entry.

**ECI-33 (Inferred — what that means for the decision).** The owner's integration concern is real and is closed
by constraint rather than by selection. It removes "how does entry work" from the discriminator set and leaves
the options judged on what they preserve for transition, forgetting, and ground — which is where the earlier
comparison already placed the trade-off.

**ECI-34 (Inferred — no cross-stream chronology is introduced).** Nothing above compares a lifecycle position
against an acquisition position. Values are **copied** at formation (same turn, not a comparison); references are
**typed** (OD-a permits); the fold reads only its own stream. Distinct from OD-g, which remains open and
unexercised.

**ECI-35 (Calculated — shared attribute ≠ cross-stream ordering).** A turn copied at formation is **one fact
recorded twice**, not two positions compared. Comparing `T_lifecycle` with `T_acquisition` to decide precedence
*would* be the OD-g question; nothing here does that, and nothing here needs it.

---

## 7. Fold requirements, extended

**ECI-36 (Inferred — what each option must satisfy to be complete, not merely representable).**

| Requirement | L1 | L2 | L3 |
|---|---|---|---|
| Entry fact on lifecycle stream | required | required | required |
| Per-stream ordering only (A2) | required | required | required |
| Total mapping from event to state | not applicable (read) | **required** | **required** |
| Unknown relation/act handling | not applicable | fail closed | fail closed |
| Closure: no-successor vs missing-successor | **not required** | **required** | **required** |
| Phantom-provisional avoidance | inherent | required | required |
| Ground retained without authorising status | payload choice | **payload choice** | inherent but must not invert |

**ECI-37 (Inferred — the one asymmetry that survives).** L1 is the only option that needs **no closure
argument** for the permanent-provisional case, because its state is asserted rather than inferred from what did
not happen. Every other difference between the options at entry is a cost both L2 and L3 pay.

**ECI-38 (Inferred — I18's demand, per option).** Each subsystem must name the evidence that would reveal its own
silent omission. For L1: an expected assertion that is absent. For L2: a stream that fails its completeness
check, or a relation whose source does not resolve. For L3: an act whose derived state is undefined. The rumour
case is the one where **nothing is expected**, so a representation that cannot say "nothing is expected here"
cannot distinguish correctness from loss.

---

## 8. Separations held

**ECI-39 (Observed).**

| Dimension | Status here |
|---|---|
| Verification (C′) | untouched; no option homes it |
| Cause | entry's cause is the formation act; referenced, never read as status (ECI-15) |
| Evidence | copied at formation; the §6.2 non-nullability problem is untouched and **not** resolved here |
| Occurrence | per-stream only; the entry fact carries formation's turn as a **value**, not a comparison |
| Identity | entry fact needs its own identity (ODA-03); no id scheme proposed |
| Authority | **ND-1** open; nothing here says who may form or advance a record |

**ECI-40 (Inferred).** Nothing in this record required OD-g, ND-1, C′ or `active`/`established` to be settled,
which confirms that the entry condition is decidable independently of them.

---

## 9. What this record does not do

**ECI-41 (Observed).** No schema, no fields, no relation vocabulary, no closure mechanism is proposed. The
"closure fact" is named as a **requirement**, not designed.

**ECI-42 (Observed).** No selection among L1/L2/L3. No reading among ECI-07/ECI-08 is adopted.

**ECI-43 (Observed).** OD-g is neither resolved nor widened: no comparandum is proposed, and the record shows
none is needed.

**ECI-44 (Observed).** No canonical file edited; 118 requirements unchanged; no ledger entry; Howl's state
machine untouched.

---

## 10. Limits

**ECI-45 (Observed).** Basis: r33 §6.1/§6.2/§8.1/§11 and K-R2/K-R11–R14/I15/I18, the OD-a decision record, and
the two prior OD-b records. **Not inspected:** any lifecycle-event schema (none exists in four schemas — ODA-12);
r32's own bytes; whether any implementation writes a record status at creation (`internal/backend/memory/*.go`
beyond the types file was not read). The ECI-07/ECI-08 reconciliation is **Inferred** and would need the
per-kind rules' intent to confirm — it is not established by the texts.

**End of record.**
