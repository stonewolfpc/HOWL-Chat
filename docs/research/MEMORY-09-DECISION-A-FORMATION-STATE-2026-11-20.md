# MEMORY-09 — Decision A

## Does the authoritative stream have to contain the formation state?

**Status:** decision brief and recommended interpretation. **P-1 and S-3 remain closed. No implementation, schema,
code, migration, test, requirement, map change or canonical capture.** A “yes” answer would settle a representation
**obligation**, not a formation-event design.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | does K-R12, with K-R14 and §6.4, require the initial `scope`, `knowledge_class` and `status` to be reconstructible from the authoritative event stream, or may formation-time record fields be the baseline? |
| Recommended answer | **Yes** — the initial state must be reconstructible from the authoritative stream; record fields may be the stored projection, not the baseline authority |
| What that would decide | the representation **obligation** only |
| What it would not decide | the formation event’s shape, payload, kind name, authority, or implementation |
| Status of P-1 / S-3 | **closed**; not reopened |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The exact question

**DA-01 (Observed — the owner’s framing).** Does K-R12, read together with K-R14 and §6.4, require the initial
`scope`, `knowledge_class` and `status` to be **reconstructible from the authoritative event stream**, or may
formation-time **record fields** serve as the baseline from which subsequent lifecycle events are projected?

**DA-02 (Observed — the caution).** A decision that the initial state must be stream-derived would settle the
representation **obligation**. It would **not** authorize a particular formation-event design — its shape, payload,
kind name, authority or implementation.

**DA-03 (Observed — the boundary).** No schema, code, tests, requirements, map change or canonical capture. P-1 and S-3
remain closed; this record does not reopen them.

---

## 2. The text

### 2.1 K-R12

**DA-04 (Source-derived — `:111`).**

> “Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus
> `scope` and `knowledge_class` changes — and a memory’s **current** status, class and scope are a **projection of that
> stream**, not mutable authoritative fields.”

**DA-05 (Calculated — what is unqualified).** The sentence says **current** status, class and scope — not “post-formation
changes” or “values after the first change.” At formation, the current values are the initial values. If those values
are not in the stream, they are not a projection of it; they are a record-field baseline.

**DA-06 (Calculated — the stream’s declared contents).** The stream’s contents are the four status kinds **plus** scope
and class changes. A formation state is not one of the named kinds. The sentence therefore does not say **how** the
initial state enters the stream — but it does say the current values are a projection of it.

### 2.2 K-R14

**DA-07 (Source-derived — `:113`).**

> “Each stream is rebuilt and checked independently; a mismatch in either is an integrity failure. The cross-stream
> join has no stored form, so it cannot drift.”

**DA-08 (Calculated — what “independently” requires).** A memory-stream rebuild that must read the record’s formation
fields to obtain the baseline is not independent of the record. The record is the **projection**, not the stream. For
the memory stream to be rebuilt from its own history, the initial state must be part of that history.

**DA-09 (Calculated — the check).** The §6.4 check compares the computed projection to the stored projection. If the
baseline comes from the stored record, the check is not purely stream-derived; it cannot distinguish a missing
formation entry from a valid initial condition.

### 2.3 §6.4

**DA-10 (Source-derived — `:366`).**

> “Acquisition events project to edges; **memory events project to the memory’s status/scope/class**. Neither
> projection may consult the other.”

**DA-11 (Calculated — the projection’s input).** The sentence names **memory events** as the input to the
status/scope/class projection. A formation-time record field is not a memory event. A projection of “record field +
memory events” is not the projection §6.4 describes.

**DA-12 (Calculated — the independence purpose).** §6.4’s “neither projection may consult the other” is about the
acquisition/memory separation, not about the record. But the same independence logic applies within the memory stream:
the stream’s own history should supply its baseline.

---

## 3. The record-field design

**DA-13 (Source-derived — §4 `:279`).** Formation “assign[s] owner chat, **scope**, **knowledge class**, attribution,
time, status” and outputs “Records + events (+ edge projection).”

**DA-14 (Source-derived — §6.1 `:340`).** `MemoryRecord` has `scope`, `knowledge_class` and `status` fields.

**DA-15 (Source-derived — §8.1 `:404`).** “Formation writes `scope=world`, `knowledge_class=experiential`,
`status=active`, `established_turn=40`.”

**DA-16 (Calculated — what the design shows).** The current architecture stores the initial values as **record
fields**. That is evidence of the current design; it is not a normative statement that the initial state may remain
outside the stream. The record fields can be the **stored projection/cache**; the question is whether they are also
the **baseline authority**.

**DA-17 (Calculated — the distinction).** “Record field as cache” is compatible with K-R12. “Record field as the
authority from which the stream’s projection starts” is not, if K-R12’s “projection of that stream” includes the
initial state.

---

## 4. The gap

**DA-18 (Source-derived — A2-43).** `provisional` has a producer — formation at §11 `:604` — but **no typed event
representation**. This is the exact gap: the initial state has a producer, but not a stream entry.

**DA-19 (Source-derived — ECI-17).** “The initial `provisional` **must be representable on the lifecycle side**. It
cannot be inferred from the acquisition structure.” This is an **Inferred** consequence of K-R12/K-R14/§6.4; it
corroborates the stronger reading but is **not** the basis by itself.

**DA-20 (Calculated — why ECI-17 is not canon).** ECI-17 is a research inference, not an owner decision or a canonical
requirement. The owner’s caution is correct: Decision A should rest on K-R12/K-R14/§6.4, with ECI-17 as corroboration,
not on ECI-17 alone.

---

## 5. The two readings

### 5.1 Yes — the initial state must be stream-derived

**DA-21 (Calculated — the case for Yes).**

1. K-R12 says **current** values are a projection of the stream; the initial values are current at formation.
2. K-R14 requires independent rebuild; a baseline outside the stream makes the memory stream dependent on the record.
3. §6.4 names **memory events** as the projection’s input; a record field is not a memory event.
4. The stream is “its own authoritative stream”; a record-field baseline is a second authority for the initial value.
5. A2-43/ECI-17 identify the missing typed event as the gap the stronger reading predicts.

### 5.2 No — the record baseline is sufficient

**DA-22 (Calculated — the case for No).**

1. The stream is a **change log**; the formation record is the initial condition, not a change.
2. §4/§6.1/§8.1 explicitly design formation as writing record fields; K-R12 was written with that design in view.
3. K-R14’s “independently” can mean independent of the **acquisition** stream, not independent of the record’s initial
   condition.
4. The formation fields are set once and are not **mutable** authoritative fields; the “not mutable” clause is
   satisfied.
5. §6.4’s “memory events project” can be read as “changes project onto the initial state.”

### 5.3 Why Yes is the better reading

**DA-23 (Calculated — the unqualified text).** The No reading requires reading “stream” as “change log” and
“projection” as “initial state + changes.” Neither qualification appears in K-R12. The Yes reading takes the sentence
as written: current values, including the initial ones, are a projection of the stream.

**DA-24 (Calculated — independent rebuild).** The No reading makes the memory stream’s rebuild depend on the record.
K-R14’s “independently” is not limited to the acquisition/memory separation; it is a property of each stream.

**DA-25 (Calculated — authority, not mutability).** K-R12’s concern is **authority**, not mutability alone. A
formation field that is set once is not mutable, but it is still an authority outside the stream. The stronger reading
preserves the stream as the sole authority.

**DA-26 (Calculated — the design is not the requirement).** §4/§6.1/§8.1 describe the current design. They do not
amend K-R12. The gap A2-43 observes is exactly what the stronger reading predicts.

---

## 6. Recommended answer

**DA-27 (Calculated — recommended interpretation).** **Yes.** K-R12, read with K-R14 and §6.4, requires the initial
`scope`, `knowledge_class` and `status` to be **reconstructible from the authoritative event stream**. Formation-time
record fields may remain as the **stored projection/cache**, but they cannot be the **baseline authority**.

**DA-28 (Observed — it is a recommendation, not an owner decision).** The owner may adopt Yes or No. This record does
not make the determination; it presents the recommended interpretation and the counter-reading (DA-21…DA-26).

**DA-29 (Calculated — what Yes would not decide).** It would not decide:

| Not decided by Yes | Why |
|---|---|
| the formation entry’s kind name | a design choice |
| its payload (scope/class/status/attribution/turn) | a design choice |
| whether it is a new kind or an existing entry | a design choice |
| who may form / formation authority | an authority question, separate |
| whether formation can immediately advance to `active` | the §8.1 collision; N-b |
| its relation to `provisional` | N-b remains separate |
| the completeness marker | Decision B |
| schema, code, tests | implementation |
| requirement capture / canonical capture | separate gates |

**DA-30 (Calculated — the representation obligation, separated from design).** Yes settles only that the initial state
must be stream-reconstructible. It does not authorize a particular event shape, payload, authority or implementation.
A separate **formation-representation design decision** would be required.

---

## 7. What follows

**DA-31 (Calculated — if the owner adopts Yes).** The current record-field baseline is **insufficient as the baseline
authority**. A formation representation is required. The formation design — kind, payload, authority — is a **separate
decision** and is not authorized here. P-1’s post-formation representation is unaffected; S-3 remains passed.

**DA-32 (Calculated — if the owner adopts No).** The record baseline is sufficient; the stream begins with the first
change; the initial state is an immutable formation field. The current design is compliant; no formation event is
required. The gap is documented as a design choice, not a violation.

**DA-33 (Observed — no implementation either way).** Neither answer authorizes schema, code, tests, requirement capture
or canonical capture. P-1 and S-3 remain closed.

---

## 8. What this record does not do

**DA-34 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **formation-event design** | not selected |
| **N-b / `provisional`** | not decided |
| **completeness marker** | Decision B; not addressed |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 9. Traceability

**DA-35 (Observed).**

| What | Where |
|---|---|
| K-R12 | `:111` |
| K-R14 | `:113` |
| §6.4 | `:366`–`:367` |
| Formation writes record fields | §4 `:279`; §6.1 `:340`; §8.1 `:404` |
| `provisional` producer, no typed event | A2-43; §11 `:604` |
| Initial state must be lifecycle-side (inferred) | ECI-17 |
| P-1 formation resolution | P1F-09…P1F-17 |
| Decision scopes | FCD-04…FCD-17 |
| P-1 / S-3 closed | P1D-01; S3A-27 |
| Map unchanged | 118 rows |

---

## 10. Record status

**DA-36 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Recommended answer | **Yes** — a recommendation, not an owner decision |
| Formation-event design selected | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **decision brief** — no canonical capture |

**End of brief.** The recommended interpretation is **Yes** (DA-27): K-R12’s unqualified “current” values, K-R14’s
independent rebuild, and §6.4’s memory-event projection together require the initial `scope`/`knowledge_class`/`status`
to be **reconstructible from the authoritative stream**; record fields may be the stored projection, not the baseline
authority. **This settles the representation obligation only** (DA-29/30). The formation event’s shape, payload,
authority and implementation are a **separate design decision**, not authorized here. **P-1 and S-3 remain closed; no
schema, code, requirement, map change or canonical capture follows** (DA-33/34).
