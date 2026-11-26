# MEMORY-09 — Formation-Representation Decision Scope

## How the initial state enters the authoritative stream — scope only

**Status:** **decision scope opened.** No design selected. No schema, code, migration, test, requirement, map change
or canonical capture. P-1 and S-3 remain closed; Decision B remains untouched.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision | which representation makes the initial `scope` / `knowledge_class` / `status` reconstructible from the stream? |
| Constraint | Decision A = Yes; K-R12/K-R14/§6.4; record fields are caches, not authority |
| Design selected | **none** |
| What it would decide | the form/kind and payload of the initial-state entry, and how reconstruction works |
| What it would not decide | N-b, §8.1 collision, formation authority, completeness, implementation, canonical |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The question

**FRS-01 (Observed — the scope).** Decision A establishes that the initial `scope`, `knowledge_class` and `status`
must be reconstructible from the authoritative stream (DAO-01/05). This decision asks: **what representation makes
that true?**

**FRS-02 (Observed — the boundary).** The decision is **separately bounded**: it settles the initial-state
representation. It does not reopen P-1/S-3, does not decide N-b or the §8.1 collision wholesale, does not decide
formation authority or warrant, and does not authorize schema, code, tests, requirements or canonical capture.

**FRS-03 (Observed — Decision B untouched).** Completeness assurance remains its own decision. This scope does not
analyze or select a completeness marker.

---

## 2. Constraints from Decision A

**FRS-04 (Calculated — the hard constraints).**

| Constraint | Source |
|---|---|
| the initial values must be derivable from the stream | DA-27; DAO-05 |
| the record fields may be caches, not the independent source | DA-27; DAO-05 |
| the memory stream must remain independently rebuildable/checkable | K-R14 |
| memory events must project the initial values | §6.4 |
| P-1’s post-formation changes are unaffected | P1D-01…P1D-10 |
| no new requirement or canonical capture | DAO-02/12 |

**FRS-05 (Calculated — the compatibility test).** Any option must let a rebuild obtain the initial values **without
reading the formation-time record fields as the authority**. A representation that only points at the record, or
that requires the record to supply the values, fails the test.

---

## 3. Design questions the decision would need to answer

**FRS-06 (Calculated — the questions).**

| # | Question | Why it matters |
|---|---|---|
| 1 | **Form/kind**: a new typed formation entry, an existing entry, or an equivalent initial-state entry? | determines the stream’s shape |
| 2 | **Payload**: which of `scope`, `knowledge_class`, `status`, attribution, `established_turn` does it carry? | reconstruction needs at least the three Decision A fields |
| 3 | **Status semantics**: does it carry `provisional`, `active`, or both? | the §8.1 collision; N-b |
| 4 | **Relationship to `provisional`**: is it the initial `provisional`; does it later transition? | N-b overlap |
| 5 | **Authority**: who/what may form; is formation a user or system action? | authority is separate from representation |
| 6 | **Projection/check**: how the memory stream projects the initial values; how the record fields are validated as caches | K-R14 independence |
| 7 | **Ordering/identity**: where the entry sits; its version/order relative to the first change | sequence integrity |
| 8 | **Migration**: how existing records without a formation entry are treated | implementation/operational; not decided here |

**FRS-07 (Calculated — the minimum).** Decision A requires **questions 1, 2, 6 and 7** at minimum. Questions 3–5 and 8
may be constrained by the answers but are **not** automatically decided by this decision; N-b and §8.1 remain separate.

---

## 4. Options

**FRS-08 (Calculated — options, none selected).**

| Option | What it is | Strength | Weakness |
|---|---|---|---|
| **New typed formation entry** | a `memory_formed` (or equivalent) entry carrying the initial values | explicit; self-contained; clear projection base | adds a kind; overlaps N-b; needs its own validation |
| **Existing `established` entry** | the formation state is an `established` entry carrying initial values | reuses an existing kind | `established` may mean post-provisional; N-b collision |
| **Existing `provisional` entry** | the formation state is the initial `provisional` entry | natural fit with §11; may resolve ECI-17 | depends on N-b; §8.1 `active` collision |
| **Synthetic initial entry** | an entry generated at rebuild/checkpoint time | no new stored kind | not authoritative if generated; may not survive rebuild |
| **Pointer/hash to the formation record** | the stream points at the record’s formation values | minimal stream change | the record remains the authority; **fails FRS-05** |

**FRS-09 (Calculated — the pointer option is excluded).** A pointer or hash that requires reading the record to obtain
the values does not make the initial state reconstructible from the stream. It relocates the external baseline; it
does not remove it.

**FRS-10 (Calculated — the new-kind option is not automatically preferred).** A new kind is the clearest, but P-1’s
minimality rationale and the N-b overlap argue for checking whether an existing entry can carry the values first. No
option is selected.

---

## 5. Evaluation criteria

**FRS-11 (Calculated — how to compare options).**

| Criterion | Test |
|---|---|
| **reconstructibility** | can a rebuild obtain the initial `scope`/`class`/`status` from the stream alone? |
| **independence** | does the rebuild avoid treating the record fields as the authority? |
| **minimality** | does it add the fewest kinds/fields consistent with Decision A? |
| **P-1 consistency** | does it leave `scope_changed`/`class_changed` unchanged? |
| **N-b/§8.1 compatibility** | can it avoid deciding them, or does it force a choice? |
| **authority separation** | does it conflate representation with who may form? |
| **implementation neutrality** | does it avoid schema/code commitments? |

**FRS-12 (Calculated — no weighted ranking).** The criteria are a checklist, not a scoring model. A design that fails
reconstructibility or independence is excluded regardless of minimality.

---

## 6. Dependencies

**FRS-13 (Calculated).**

| Dependency | Relationship |
|---|---|
| **N-b / `provisional`** | the formation entry may be the initial `provisional`; N-b remains separate |
| **§8.1 collision** | `active` vs `provisional` at formation; the formation entry’s status may force or defer it |
| **P-1** | post-formation changes unchanged; the formation entry is the baseline |
| **Decision B** | completeness assurance; independent, untouched |
| **Migration** | existing records without a formation entry; implementation/operational, not decided here |
| **Authority / warrant** | who may form; separate from representation |

**FRS-14 (Calculated — coordination without subsumption).** The formation-representation decision may need to
**coordinate** with N-b and §8.1, but it should not **subsume** them. If the initial status value cannot be chosen
without N-b, that part is deferred rather than decided by implication.

---

## 7. What the decision would and would not decide

**FRS-15 (Calculated — would decide).**

- the form/kind of the initial-state entry;
- its payload, at least `scope`, `knowledge_class` and `status`;
- how the memory stream reconstructs the initial values;
- how the record fields are validated as caches;
- the entry’s ordering/identity relative to the first change.

**FRS-16 (Calculated — would not decide).**

| Not decided | Where it belongs |
|---|---|
| N-b’s `provisional` vocabulary | N-b |
| the §8.1 `active`/`provisional` collision | N-b/§8.1 |
| formation authority / warrant | authority decision |
| completeness marker | Decision B |
| schema / code / tests | implementation |
| requirement capture / canonical capture | separate gates |
| migration/backfill | implementation/operational |

**FRS-17 (Calculated — the authorization boundary).** This is a **representation decision**. Adopting an option would
not authorize schema, code, migration, tests, requirement capture or canonical capture. Those need separate
authorization.

---

## 8. Sequence

**FRS-18 (Calculated).**

1. Decision A — **Yes**, owner-adopted (DAO-01).
2. **This scope** — open the formation-representation decision; no design selected.
3. Formation-representation decision — choose the form/payload and reconstruction rule.
4. Decision B — completeness assurance, independently, on its own failure target.
5. Implementation authorization — only after the representation decisions and separate gates.

**FRS-19 (Calculated — no forced sequencing of N-b).** N-b and §8.1 may be decided before, alongside, or after the
formation representation. This scope does not order them; it only flags the overlap.

---

## 9. What this record does not do

**FRS-20 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **Decision A** | resolved Yes; not reopened |
| **formation design** | not selected |
| **N-b / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 10. Traceability

**FRS-21 (Observed).**

| What | Where |
|---|---|
| Decision A owner determination | DAO-01…DAO-14 |
| Decision A brief | DA-01…DA-36 |
| K-R12 / K-R14 / §6.4 | `:111`–`:113`, `:366`–`:367` |
| Formation record fields | §4 `:279`; §6.1 `:340`; §8.1 `:404` |
| `provisional` producer, no typed event | A2-43; §11 `:604` |
| Initial state must be lifecycle-side | ECI-17 |
| P-1 formation resolution | P1F-09…P1F-17 |
| P-1 / S-3 closed | P1D-01; S3A-27 |
| Map unchanged | 118 rows |

---

## 11. Record status

**FRS-22 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Formation design selected | **0** — scope only |
| N-b / §8.1 decided | **0** |
| Decision B touched | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **decision scope** — no canonical capture |

**End of scope.** Decision A is **Yes** at the interpretation level (DAO-01). This scope opens the **separately bounded
formation-representation decision**: how the initial `scope`/`knowledge_class`/`status` enters the authoritative
stream so that it is reconstructible, without selecting a design (FRS-01…FRS-19). The pointer/hash option is excluded
because it leaves the record as the authority (FRS-09). **N-b, §8.1, formation authority, Decision B and implementation
remain separate** (FRS-13…FRS-17). **No schema, code, requirement, map change or canonical capture follows**
(FRS-20/22).
