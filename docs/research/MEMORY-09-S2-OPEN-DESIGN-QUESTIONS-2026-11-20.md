# MEMORY-09 — S-2 Open Design Questions

## Completeness, recovery, concurrency and formation before any implementation authorization

**Status:** design-question analysis. **S-3 remains passed; these are the remaining details, not acceptance failures.**
**No option is selected. No schema, code, migration, test, deployment, requirement, map change or canonical capture.**
Howl state-machine design and UX remain outside.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Questions | completeness marker · recovery policy · concurrency guarantee · formation baseline |
| Finding | each has viable options and trade-offs; none invalidates S-3; none is selected |
| What this is | the examination needed before an implementation authorization request |
| What this is not | implementation authorization, schema design or a requirement change |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Task and boundary

**S2O-01 (Observed — the owner's recommendation).** Before drafting an implementation authorization request, examine
the four unresolved design questions: how completeness is known, what happens on incomplete history, what guarantees
same-dimension serialization, and how the initial baseline is represented.

**S2O-02 (Observed — S-3 stands).** P-1 passed the logical acceptance criteria (S3A-27). These questions are the
details needed to turn the accepted representation into an implementable design; they do not reopen the acceptance.

**S2O-03 (Observed — the boundary).** No schema, field definition, code, migration, test, probe, model call,
deployment, requirement or map change, or canonical capture. The map stays at **118**. Nothing here authorizes
implementation.

---

## 2. Completeness marker

**S2O-04 (Observed — why it matters).** The sequence detects a missing **internal** entry because a higher version
survives. It cannot detect a missing **latest** entry because no later version exists (P1D-17; S2F-10/11). A
completeness marker is the only way to distinguish "complete" from "stale but undetectable."

**S2O-05 (Calculated — the options).**

| Option | What it provides | What it costs | What it cannot do |
|---|---|---|---|
| **No marker** | nothing beyond the current gap detection | none | cannot detect a missing latest |
| **Expected-max counter** | a durable per-memory/per-dimension maximum; a rebuild with a lower max detects a missing latest | a mutable counter that must be updated atomically with the version assignment; if lost, the guarantee weakens | cannot prove completeness; the counter itself can be stale |
| **Periodic checkpoint** | a snapshot of the sequence (max version and/or hash) at a point | storage and a trust anchor; a missing latest after the checkpoint is undetected until the next checkpoint | cannot cover changes after the checkpoint |
| **Terminal marker** | a marker at close/deletion | only applies at the end of a sequence | not useful for live sequences |
| **External manifest** | an independent record of expected versions | an external dependency; must not become a projection input (K-R14) | cannot prove the stream is complete if the manifest is also incomplete |

**S2O-06 (Calculated — the trust problem).** A marker is a **detection aid**, not a proof. Any marker can be stale,
lost or itself incomplete. The strongest option is a durable expected-max counter updated atomically with the version
assignment; the weakest is no marker. A periodic checkpoint gives bounded confidence, not a guarantee.

**S2O-07 (Calculated — the boundary).** A marker must not become an authority for the **value**; it is an integrity
aid. If it is consulted by the projection, it must not violate K-R14's independent rebuild. A separate integrity check
that compares the stream to the marker is safer than feeding the marker into the projection.

**S2O-08 (Observed — no selection).** The marker is an open design decision. It is not required by K-R12/K-R14/§6.4;
it strengthens the completeness story. The implementation authorization request would need to state whether a marker
exists, its form, and its trust assumptions.

---

## 3. Recovery policy

**S2O-09 (Observed — the scenarios).** Recovery differs by failure:

| Failure | What is known | What is not known |
|---|---|---|
| **internal gap** | a higher version survives; the gap is detectable | the missing value |
| **missing latest** | nothing, unless a marker detects it | whether the projection is stale |
| **duplicate version** | a conflict exists | which entry is authoritative |
| **out-of-order arrival** | the higher version is present; the lower is pending | whether the lower will arrive |

**S2O-10 (Calculated — the policy options).**

| Option | Behavior | Trade-off |
|---|---|---|
| **Fail closed** | the dimension is unavailable; reads surface an integrity state | safest; may block use |
| **Last known + flag** | serve the last contiguous value as "last known, incomplete"; never as authoritative | transparent; consumers must handle the flag |
| **Checkpoint restore** | restore from a checkpoint/backup, then verify | can recover lost entries; depends on a trustworthy checkpoint |
| **Explicit re-record** | re-record the missing change from an authoritative source | can repair; requires an authority and a new entry |
| **Escalate to inspection** | no automatic repair; an operator investigates | safe; may be slow |

**S2O-11 (Calculated — the non-negotiable boundaries).**

- **No silent repair.** A missing value is not invented or guessed.
- **No untrusted projection as authoritative.** A value derived across a gap is not presented as current truth.
- **No lifecycle-status change for context.** An incomplete projection is an integrity state, not a new lifecycle
  `status`; context inactivity remains separate (S2H-12).
- **No authority laundering.** Recovery does not decide who may change class/scope; it repairs or surfaces history.
- **No validity decision.** Recovery does not decide whether a missing change was warranted.

**S2O-12 (Calculated — the shape of a coherent policy).** For an **internal gap**, the last contiguous value may be
served with an explicit incomplete flag; the higher versions are held. For a **missing latest**, without a marker the
system cannot know; with a marker, treat it as a gap. For a **duplicate**, reject/conflict and untrust the dimension.
For **out-of-order**, wait for the predecessor; do not repair. Any escalation is explicit.

**S2O-13 (Observed — no selection).** The policy is an open design decision. It is operational; it does not change the
logical representation. The implementation authorization request would need to choose a default and state the
consumer-visible behavior.

---

## 4. Concurrency guarantee

**S2O-14 (Observed — the invariant).** Same-dimension writes are serialized; versions are contiguous and unique; the
version assignment is atomic with the entry write (P1D-07/08; S2C-05/07).

**S2O-15 (Calculated — the mechanisms).**

| Mechanism | How it guarantees the invariant | Trade-off |
|---|---|---|
| **Single writer/queue per (memory, dimension)** | one writer assigns the next version in order | simple; may bottleneck; single point of failure |
| **Compare-and-swap (CAS)** | read the max, write with an expected max; conflict ⇒ retry | good throughput; needs atomic conditional writes; retries |
| **Pessimistic lock** | lock the pair during assignment | predictable; contention/deadlock risk; timeout needed |
| **Central sequencer** | one component assigns versions; writers append | strong ordering; central dependency |

**S2O-16 (Calculated — what the mechanism must provide).** Atomicity of **read-max + assign + write**; a conflict is
rejected or retried, never silently overwritten; a retry does not create a duplicate; the accepted version is
contiguous. The logical invariant is fixed; the mechanism is implementation.

**S2O-17 (Calculated — cross-dimension independence).** Scope and class have separate sequences, so a scope write and
a class write do not contend. A mechanism that serializes the whole memory would violate the independence property
unnecessarily; a per-dimension guard is sufficient.

**S2O-18 (Observed — no selection).** The mechanism is an implementation decision. It does not change P-1. The
implementation authorization request would need to name the guarantee and the failure behavior.

---

## 5. Formation baseline

**S2O-19 (Observed — why it matters).** P-1 uses the formation value as the baseline (P1D-06). If formation is not in
the lifecycle stream, the projection's base is an initial condition, not a stream entry. K-R12 says current values
are a projection of the stream; the cleaner the baseline's provenance, the stronger the acceptance story.

**S2O-20 (Calculated — the options).**

| Option | What it provides | What it costs |
|---|---|---|
| **Baseline in the record** | the memory's initial scope/class; simple | the base is not a stream entry; K-R12's "projection of the stream" is weaker |
| **Existing formation/provisional entry** | if the entry already carries scope/class, the baseline is in the stream | depends on the existing entry's payload; must verify |
| **New formation event** | an explicit `memory_formed` entry with initial scope/class | a new kind; not covered by P-1; needs a separate decision |
| **Initial checkpoint** | a snapshot of the baseline | adds a marker; overkill if the formation entry exists |

**S2O-21 (Source-derived — the existing architecture).** §11's diagram begins with `provisional (status, class,
attribution set)`; §8.1 says formation writes `scope`, `knowledge_class`, `status` (`:404`). Whether the existing
`provisional` entry carries `scope` is an existing representation question (A2-11/A2-12), not a P-1 creation.

**S2O-22 (Calculated — the preference, not a decision).** If the existing formation/provisional entry already carries
scope and class, use it as the baseline; that is the cleanest and adds no kind. If it does not, the baseline remains
the record's initial condition, or a separate formation-event decision is needed. **P-1 established only
`scope_changed` and `class_changed`; it did not authorize a new formation kind.**

---

## 6. Cross-cutting: logical versus implementation

**S2O-23 (Calculated — the classification).**

| Question | Kind | Where it belongs |
|---|---|---|
| completeness marker | representation/operational | a design decision; may add a marker element |
| recovery policy | operational | an implementation authorization detail |
| concurrency mechanism | implementation | an implementation authorization detail; the invariant is logical |
| formation baseline | representation | may be resolved by the existing formation entry; a new event needs a separate decision |

**S2O-24 (Calculated — what an implementation authorization request would need to answer).**

| Item | Needed |
|---|---|
| completeness | marker: none / counter / checkpoint / manifest; trust assumptions |
| recovery | default behavior; last-known flag; escalation; duplicate handling |
| concurrency | mechanism; atomic assignment; conflict/retry behavior |
| formation | baseline source; whether the existing provisional entry carries scope/class |
| schema/fields | physical names, types, storage — separate authorization |
| tests/deployment | separate authorization |
| requirement capture | separate step; map stays 118 |
| canonical capture | separate step |

**S2O-25 (Observed — no request is drafted).** This record examines the questions; it does not draft the implementation
authorization request and does not authorize implementation.

---

## 7. What this record does not do

**S2O-26 (Observed).**

| Untouched | State |
|---|---|
| **completeness marker** | no option selected |
| **recovery policy** | no option selected |
| **concurrency mechanism** | no option selected |
| **formation baseline** | preference noted; no decision |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |
| **authority / validity / causes** | not decided |
| **N-b / N-c / N-a / S-1** | not decided |

---

## 8. Traceability

**S2O-27 (Observed).**

| What | Where |
|---|---|
| P-1 decision | P1D-01…P1D-10 |
| S-3 acceptance | S3A-27 |
| Missing-latest limitation | P1D-17; S2F-10/11 |
| Gap taxonomy | S2F-06/07 |
| Concurrency invariant | S2C-05/07; S2F-04 |
| Formation writes | r33 `:279`, `:404`, `:604`; §11 `:604` |
| Context ≠ lifecycle status | S2H-12 |
| K-R12 / K-R14 / §6.4 | r33 `:111`–`:113`, `:366`–`:367` |
| Map unchanged | 118 rows |

---

## 9. Record status

**S2O-28 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Option selected | **0** |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **design-question analysis** — no canonical capture |

**End of analysis.** Four questions remain, and none invalidates S-3 (S2O-02). **Completeness** is a detection aid, not
a proof: no marker, an expected-max counter, a checkpoint or an external manifest, each with a trust cost
(S2O-04…S2O-08). **Recovery** must choose between fail-closed, last-known-with-flag, checkpoint restore, explicit
re-record or escalation, under the non-negotiables: no silent repair, no untrusted authority, no lifecycle-status
change for context (S2O-09…S2O-13). **Concurrency** needs atomic read-max/assign/write, with per-dimension
independence; the mechanism is implementation (S2O-14…S2O-18). **Formation** is cleanest if the existing
provisional/formation entry already carries scope/class; P-1 did not authorize a new formation kind (S2O-19…S2O-22).
**No option is selected, no implementation is authorized, and no requirement or canonical change follows** (S2O-25/26).
