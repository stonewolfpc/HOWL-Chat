# MEMORY-09 — Architectural contracts, revision 13

**Status: PROPOSAL REV 13 — the authorization boundary is defined in terms of what may be read, and `domain` /
`target` are specified as opaque labels rather than lookups. Not implemented, not tested, not ratified.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Supersedes | REV 12 (`0592d689…`) **for §2.2's "Data read: none" row, §2.3's S-5 and S-11, §3.1's A1 and A6, §3.2's framing paragraph, §4's invariant 19 and §7 rows 17h″/17m only** |
| Stands unamended | REV 12 in every other respect; REV 1–11 preserved byte-identical |
| Authority | Owner's two precision points, recorded verbatim (§1) |
| Ratification | **Not ratified.** On acceptance the operative text becomes REV 13 §2–§5 |
| Implementation | **Unmodified.** `eval/mem1/m09/` byte-identical to R2, including `authority.py` |
| Tests / probes | **0.** Nothing executed, nothing compiled |
| Amends | **Nothing** to r32, v9, the 118-requirement map or the ledger |

---

## 0. What changed

| Item | REV 12 | REV 13 |
|---|---|---|
| Class A "reads no data" | Literal, and false as soon as the gate consults a grant store | Restated over two named classes of data: **authorization metadata** (readable to evaluate the gate) and **protected memory data** (unreadable before the gate) |
| The gate's own knowledge | Unspecified | Named: `require(domain, actor, "promotion")` consults grants and delegations; that read is permitted, is the only pre-gate read, and is never echoed |
| A1 | "Validate `view`; resolve `(domain, target)`" — the second half ambiguous | Validate `view` only. `domain` and `target` are opaque labels, compared by equality, never looked up, never refused |
| A domain/target catalogue | Would have been an ordering hole | If an implementation keeps one, its lookup is an explicit capture step (A6′) after the gate, and an unknown label stays indistinguishable from an empty one |
| `unresolved(record_absent)` | Correctly three-valued, but its status was only implied | Stated explicitly: a capture-resolution outcome, not a memory state or a knowledge state; never published as a `state` |
| Invariants | 1–19 | 1–19, with 19 restated and 20–21 added (metadata boundary; no existence oracle) |

No rule, requirement, state value or tally changes. `authority.py` is not edited — the REV 7 decision not to
touch it stands.

---

## 1. Owner instruction, recorded verbatim

> "**A. The Class A 'no data read' guarantee needs a precise definition of the authorization boundary.** …
> `no_domain_authority` must establish whether the caller has promotion authority. Where does that authority
> check obtain its information? If it reads an authorization/delegation store, then the literal claim that
> Class A reads no data is too broad. That may be perfectly acceptable: authorization metadata is not the
> protected memory data whose access the gate is meant to control. But the contract should say so explicitly.
> I would distinguish: Authorization metadata: may be consulted to evaluate the gate. Protected memory data:
> records, groups, proposition registry and their contents; inaccessible until the gate passes. That
> preserves the actual security guarantee without making an impossible claim about how authorization works."

> "**B. A1 should distinguish argument validation from data resolution.** A1 currently says: 'Validate
> `view`; resolve `(domain, target)`'. But the refusal table says view validation happens before any store
> access. If resolving the domain or target requires looking up stored data, A1 could violate that ordering.
> Specify whether `(domain, target)` is simply validated as supplied input, or whether it is looked up in
> persistent state. If it is a lookup, its position relative to the authorization gate needs to be explicit."

> "Neither issue requires new memory semantics or a new requirement. They are precision corrections to the
> capture and refusal procedure."

> "One additional distinction is worth preserving: `unresolved(record_absent)` is a capture-resolution
> outcome, not a new kind of memory or a new knowledge state. The document should keep that distinction clear
> when saying no state values have changed."

> "My recommendation: one narrowly scoped REV 13, limited to those two precision points. Then review the
> resulting text for acceptance. Do not reopen REV 1–12, the 57 rules, the 118 requirements, or the broader
> architecture."

**Scope:** the two boundary definitions and the status of `unresolved(record_absent)`. Nothing else.

---

## 2. The authorization boundary

### 2.1 The overclaim

REV 12 S-5: "Authorization refusals take no snapshot and **read no data**… not the record store, not the
group store, not the registry." But `no_domain_authority` is produced by `require(domain, actor,
"promotion")`, and that call must consult something. The literal reading is either false or it implies the
gate operates on information it does not have. Neither is acceptable in a contract that will be implemented
by someone else.

### 2.2 Two classes of data

| Class | Contents | May Class A read it? |
|---|---|---|
| **Authorization metadata** | Grant records keyed by `(domain, actor, operation)`; delegation edges and their origins; the operation names a grant may carry | **Yes** — this is what the gate is for |
| **Protected memory data** | The record store and its record contents; the group store and its members; the proposition registry and its identities; and every value derived from them (`accepted`, `returned_identities`, `resolution`, `groups`, `establishes`) | **No** — unreachable until the gate passes |

Authorization metadata contains no proposition identity, no record identity, no belief, no text about any
character or event. It says which actor may perform which operation in which domain. Consulting it therefore
discloses nothing the gate exists to protect — and the gate cannot function without it.

### 2.3 Rules

| Rule | Statement |
|---|---|
| S-5 (restated) | **Authorization refusals read no protected memory data and take no snapshot.** They may consult authorization metadata, and only that, to evaluate the gate. Their envelope contains no protected-memory field and no authorization-metadata field. |
| S-12 (new) | **Metadata boundary.** Protected memory data is read only after the gate admits the call, and only inside a capture attempt. Authorization metadata is read only to evaluate the gate and to choose the refusal code; it is never published, summarised, counted or hinted at in any response, refused or served. |
| S-11 (extended) | **Content independence** now covers two repositories: neither the refusal class nor the code may vary with the content of the stores **or** with the content of the authorization metadata beyond the boolean outcome of the gate. |

### 2.4 The no-oracle consequences

| Question a caller might try to answer | Answer given | Why safe |
|---|---|---|
| "Does this domain exist?" | `no_domain_authority`, identical to every other unauthorized case | REV 7 §5.1 already collapses undeclared domain, absent grant and `actor=None`; the metadata boundary does not reopen it |
| "Who holds the grant?" | Nothing. No actor list, no grant count, no delegation edge | S-12 forbids publishing metadata; a refusal carries a code only |
| "Which operation was refused?" | Nothing beyond the gate named | `require` names one operation per call site; the code does not enumerate |
| "Does this target hold data?" | Nothing. Class A precedes every protected read; Class B is unreachable pre-gate (S-10) | §3's opaque labels keep an unknown target indistinguishable from an empty one |

Reading grants is thus permitted and bounded: the gate may *use* authorization metadata to decide, and may
never *utter* it.

---

## 3. `domain` and `target` are opaque labels

### 3.1 A1, corrected

| Step | REV 12 | REV 13 |
|---|---|---|
| A1 | "Validate `view`; resolve `(domain, target)`" | **Validate `view`** only — a name check against the three defined views, no store access. Class A: `unknown_view` |

`domain` and `target` require no validation and no resolution, because they are **labels, not references**:

- They are supplied as strings by the caller and compared by equality against the `domain` and `target`
  fields already carried by records.
- Nothing is looked up to interpret them. An unmatched label simply matches nothing.
- A label that matches nothing must be **indistinguishable** from a label that matches an identity whose
  records are all superseded: both yield an empty result, neither yields a refusal.

That last line is the reason no new refusal code is introduced. A `unknown_domain` code would be an existence
oracle over protected memory, which S-12 exists to prevent — and it would also break REV 7 §5.1's uniform
collapse.

`domain` is additionally the key the gate consults in authorization metadata. Passing it to `require(...)` is
therefore not a protected-memory read; it is the gate's own argument.

### 3.2 If an implementation keeps a catalogue

Some implementations may want a table of known domains or targets — for validation, display, or indexing.
The contract permits it and fixes its position:

| Step | Action | Ordering |
|---|---|---|
| A6′ | **Catalogue lookup** (optional): if a domain/target catalogue exists, look up `(domain, target)` here | Inside the capture, after the gate (A2), inside the version bracket (A10) |

| Consequence | Rule |
|---|---|
| A catalogue is protected memory data, not authorization metadata | It may not be read before the gate; the lookup is A6′, not A1 |
| An unknown label at A6′ is not a refusal | It yields the same empty result as a known label with no matching records |
| The catalogue must be part of the same coherent capture | Otherwise `returned_identities` and `accepted` could be drawn from different catalogue states — the REV 11 defect again |

A1's guarantee is therefore narrower and true: **input validation touches no store**, because the only input
that needs validating is the view name.

---

## 4. What `unresolved(record_absent)` is not

REV 12 §3.2 introduced a three-valued resolution map. Its status is now stated explicitly, because "no state
value changed" must not be read as though it had.

| It is | It is not |
|---|---|
| A value in `S.resolution`, describing what a member reference resolved to **during one capture** | A `state` value in an entry |
| A report that the group store names a record the record store does not contain | A memory kind, a knowledge state, or a belief |
| A divergence between two stores, surfaced as an integrity fault | A dispute, a settlement, or an absence of evidence |

| Field | Values |
|---|---|
| entry `state` | `none`, `settled`, `unresolved_relation`, `conflict`, `multiple` — **unchanged, five values** |
| `relation` | `asserted`, `not_established` — unchanged |
| `S.resolution[m]` | `resolved(pid)`, `unresolved(record_absent)` — internal, never published as a state |

`unresolved(record_absent)` reaches the envelope only as `establishes: false` with `cause:
"member_unresolved"` (authority) or an `integrity` code (canon), which are diagnostics of an envelope that
already existed. It adds no code to the map and no value to the state vocabulary.

---

## 5. Invariants 19–21

REV 9 §6 items 1–12 and REV 10–12 items 13–18 stand. Item 19 is restated; 20 and 21 are added.

19. **Refusal separation:** Class A refusals read no protected memory data and take no snapshot; Class B
    refusals are reachable only after the gate, discard any partial capture, and publish nothing derived from
    it. Both use one fixed envelope shape, carrying a code and nothing else.
20. **Metadata boundary:** protected memory data is read only post-gate and only inside a capture attempt;
    authorization metadata is read only to evaluate the gate and choose a code, and is never published,
    counted or hinted at. Neither a served nor a refused response increases in content with the size of the
    authorization metadata.
21. **No existence oracle:** `view` is the only input whose validity can produce a refusal. `domain` and
    `target` are opaque labels; an unmatched label is indistinguishable from an empty match, and no refusal
    class or code varies with whether a domain, target, identity or record exists.

These constrain the proposed mechanism. They add nothing to the ratified map and are not executed here.

---

## 6. Ratification status

| Item | State |
|---|---|
| REV 9 | Ratified within its declared sections; recorded in REV 10 §6, unchanged. |
| REV 10's snapshot contract | Not ratified; superseded. |
| REV 11 | Accepted in principle; superseded in part. |
| REV 12 | Accepted in principle. On acceptance of this record, the operative snapshot text is **REV 13 §2–§5** with REV 12 §2–§4 as amended here. |
| REV 8's unamended sections | Still flagged; unchanged by this revision. |
| Tallies | Unchanged: 57 rules, 118 requirements, same statuses. |
| Implementation, testing | Separate gates. Neither requested. |

---

## 7. Implementation list — deltas only

REV 12 §7 stands, with these rows replaced or added. Nothing is applied.

| # | File | Symbol | Change |
|---|---|---|---|
| 17h‴ (REV 12 17h″) | `engine.py` | `_capture(...)` | A1 now validates only `view`; an optional catalogue lookup is A6′; gate at A2 consults metadata and nothing else |
| 17n | `engine.py` | `_validate_inputs(view)` | Total on the view name; no store access; the only input validation in the call |
| 17o | `engine.py` | `_catalogue_lookup(...)` *(conditional)* | Only if a catalogue exists; called at A6′, inside the capture; empty result, never a refusal |
| 17m′ (REV 12 17m) | `errors.py` | refusal classes | Both classes carry `code` only; Class A additionally may consult metadata before constructing the code |
| 17p | `authority.py` | — | **No change.** `require(domain, actor, operation)` — REV 3's signature — is consulted as-is; the metadata it reads is not edited or extended |
| 6⁸ / 18⁷ / 9⁸ | `retrieval.py`, `engine.py` | `world`, `belief_facts`, `propositions` | Unchanged from REV 12 apart from propagating the A1/A6′ split |

**D-A and D-B remain confirmed and uncorrected.**

**Untouched:** `eval/mem1/m09/*.py` (byte-identical, including `authority.py`), `eval/mem1/model.py`, the G1
harness, production Howl Chat, r32, plan v3, Amendment 1, v9, the map, the ledger, REV 1–12.

---

## 8. Coherence

| Property | v9 rules | Map / contract | Verdict |
|---|---|---|---|
| Different questions do not borrow each other's verdicts | 6, 7, 41 | C-2; J12–J14 | **C** |
| Belief ascriptions never answer world truth | 7, 56 | C-2; D10; K12 | **C** |
| Conflict only where established | 13, 41 | O-10; J12, J13 | **C** |
| Underspecification is not conflict | 8, 34 | K15 | **C** — `unresolved(record_absent)` is a store divergence, not a dispute (§4) |
| No adjudication by the engine | 34 | K15, O-10 | **C** |
| History preserved, not re-reported as live | 20, 44 | D-02; K22 | **C** |
| Enumeration is authority-only | 21, 24, 35 | REV 4 §3, §7; K8, K21 | **C** — gate reads metadata only, then every protected read is post-gate |
| Canon reads add no text surface | 18, 32, 47 | REV 4 §4.2; L2; H7 | **C** |
| Fail closed | 8 | F4 | **C** — no admissible capture, no answer |
| Engine states what it did | 8 | K15 | **C** — A refused caller is told which boundary refused them, and nothing about the contents |

**Tallies unchanged:** rules **39 U · 13 M · 5 N · 0 R · 0 X**; requirements **86 U · 22 M · 8 N · 1 R ·
1 X**. F7 the only **R**; B3 the only **X**.

**Residuals unchanged:** R-A closed; RAT-1/RAT-2 recorded; R-E bounded; R-B, R-C, R-D open.

---

## 9. What this record does and does not claim

**Still open, carried forward unchanged:** B3; belief-voiced dialogue; replay / F7; R-B, R-C, R-D.

**Not claimed:**

- No implementation. No source file changed, `authority.py` included. D-A and D-B remain confirmed **and
  uncorrected**.
- No execution. No test, probe, suite or invariant check ran; nothing imported, compiled or linted.
- **No concurrency safety.** The current store is single-threaded and read-only on these paths, so A10's retry
  path and Class B are unreachable today. That is vacuity, not evidence: the contract states what a future
  implementation must guarantee and records that nothing has demonstrated it.
- No claim about how `supersede` publishes its two writes; nothing was executed to determine it.
- No new rule, requirement, state value or capability. `inconsistent_snapshot` remains a proposed refusal
  code; `unresolved(record_absent)` is an internal resolution value, not a state (§4); no new refusal code is
  introduced for domains or targets, deliberately, because it would be an existence oracle.
- No amendment to r32 (`e88ef2d1…`), plan v3 (`6c4078bd…`), Amendment 1 (`6afd2a0d…`), v9 (`698254ca…`), the
  map (`b4ba33cc…`), the ledger (`da3e7239…`), or REV 1–12 (`98447423…`, `74ac6f77…`, `de4f68b6…`,
  `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`, `5b3a9008…`, `7d703b77…`, `40850419…`,
  `0592d689…`).
- No transfer, destination interface, sharing, or cross-chat carryover.
- No change to the REV 10 §6 ratification record or its flagged reading.

**Sequence from here:** owner review of REV 13 for acceptance → acceptance of the snapshot contract as
REV 13 §2–§5 → separate implementation authorization → implementation → conformance review → separate
testing authorization. This record requests neither authorization.
