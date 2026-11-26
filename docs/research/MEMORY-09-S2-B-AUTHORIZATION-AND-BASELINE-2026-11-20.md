# MEMORY-09 — S-2 B Authorization and Decision Baseline

## The logical representation process is open; implementation and capture are not

**Status:** authorization record and decision baseline. **The owner selected B on 2026-11-20; S-2 is open.** S-2 is
authorized to reach a **logical** representation decision under the request (S2R-01…S2R-35) and the binding Howl
separation addendum (S2H-01…S2H-22). **No carrier is adopted. No schema, code, migration, test, requirement, map
change or canonical capture is authorized.**

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Owner selection | **B** — A-6 does not carry into S-2 |
| Authorizes | the S-2 logical representation decision process |
| Does not authorize | Howl state-machine/UX design; authority or validity; N-b/N-c/N-a; implementation; requirement capture; canonical capture |
| Acceptance test | **S-3** — O-KR12, K-R14, §6.4 |
| Carrier state | **none adopted**; the investigation may recommend, not adopt |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The authorization

**S2B-01 (Observed — the owner's selection).** The owner selected **B**: S-2 is authorized on the basis that **A-6 does
not carry into S-2**. A-6 remains a term of the **N-a decision** and is not repealed (S2R-20; S2A-28).

**S2B-02 (Observed — what B authorizes).** Under the request, S-2 may:

| Authorized | Scope |
|---|---|
| investigate how `scope` and `knowledge_class` changes are represented | logical representation |
| consider logical representations requiring new kinds or fields | no physical schema |
| determine the projection and reconstruction rules | O-KR12 / K-R14 |
| determine the integrity check | §6.4 |
| examine the entry representation | shape, not fold |
| examine any S-2-specific ordering need | local; must not settle N-c |
| conclude that no admissible representation is currently justified | process-level decline (S2H-15) |

**S2B-03 (Observed — what B does not authorize).** Designing Howl's state machine or UX; deciding transition
authority or validity; implementing code, schema, migration or tests; adding requirements to the 118-row map;
canonical capture or modification (S2R-25; S2H-16…S2H-19).

**S2B-04 (Observed — the addendum binds).** S2H-12 remains binding: **context inactivity is not a memory-lifecycle
event.** A memory removed from the active context is not thereby retracted, superseded, corrected or otherwise
changed in lifecycle `status`.

**S2B-05 (Observed — no carrier adopted).** B authorizes the **investigation**, not its answer. No carrier has been
selected. A recommendation from the investigation is not an adoption.

---

## 2. Inherited constraints

**S2B-06 (Calculated — what S-2 must satisfy).**

| Constraint | Requirement |
|---|---|
| **K-R12** | the carrier belongs to the lifecycle stream; class/scope changes are its subject; current values are a projection; no mutable authoritative field |
| **K-R14** | each stream is rebuilt and checked independently; no cross-stream input; no stored join |
| **§6.4** | memory events project to status/scope/class; the check is deterministic; mismatch is an integrity error |
| **K-R2** | the **integrity principle** — an unreconstructable projection is an integrity failure; its literal `derived_from_event_ids` clause is edge-specific and is **not** imported |
| **K-R13** | knowledge × memory is a read-time join; never stored |
| **K-R15** | acquisition ordering only; no automatic lifecycle ordering rule |
| **§11** | promotion authority is separate; S-2 does not decide who may promote or whether a change is permitted |
| **D1 / D5** | `world`/`scenario`/`character_public`/`character_private` are scopes, not owners; decay never affects scope |
| **K-R8 / K-R9 / K-R10 / I13** | class is authored, not inferred; ambient only for authored definitions; no downgrade by repetition, age or decay; popularity is not authorship |
| **r33 `scope` vs v9 `domain`** | separate axes; S-2 uses r33 scope values only; O-8a authority-promotion is distinct |
| **S2H-12** | context activation/deactivation is not lifecycle status |

**S2B-07 (Calculated — the A-6 change).** C3D-11's "no admissible positive carrier" was a finding **under A-6**. B
lifts that bar for **logical consideration**: P-1, P-2 and P-6 become logically admissible. **K-R14 still excludes an
external join (P-4)**; K-R15 still does not supply a lifecycle order; P-5 remains storage, not a carrier.

---

## 3. Exact responsibilities

**S2B-08 (Calculated — S-2's responsibilities).**

| Responsibility | In scope | Boundary |
|---|---|---|
| represent post-formation class/scope changes | yes | logical only |
| choose a logical carrier form (or decline) | yes | no physical schema |
| define the projection/reconstruction rule | yes | K-R12 / K-R14 |
| define the integrity check | yes | §6.4 |
| define the entry representation | yes | shape, not fold |
| define an S-2-specific ordering rule, if needed | yes | local; must not settle N-c |
| conclude no admissible representation is justified | yes | records O-KR12 as unmet |

**S2B-09 (Calculated — S-2 is not responsible for).** The N-a fold; N-b status vocabulary; N-c status ordering;
§11 promotion; O-8a; v9 authority domain; state-machine selection/switching/combination/activation/deactivation/
restore; UX; authority; validity; causes; implementation; requirement capture; canonical capture (S2R-17; S2H-11…14).

---

## 4. Exclusions

**S2B-10 (Calculated — the exclusion set).**

| Excluded | Gate |
|---|---|
| Howl state machine, context orchestration, UX | separate Howl responsibility; S2H |
| authority / permission / sufficiency | §11 / C-1 R-E; separate |
| validity / warrant / truth | D6; separate |
| causes of class/scope change | D5 / K-R9 / K-R10 / I13; not S-2's |
| N-b / N-c / N-a fold | S-1 |
| §11 promotion / O-8a / v9 `domain` | separate axes |
| physical schema / DDL / code / migration / tests | separate implementation authorization |
| requirement capture / map change | separate step; map stays **118** |
| canonical edits (r33, v9, ledger, registry) | separate canonical-capture authorization |

---

## 5. Acceptance conditions

**S2B-11 (Calculated — when an S-2 representation is adequate).**

| Condition | Test |
|---|---|
| **recording** | every actual class/scope change has a lifecycle-stream record that determines the new value |
| **projection** | the current value is derived from the stream; it is not a mutable authority |
| **independent rebuild** | the value can be reconstructed from the lifecycle stream alone (K-R14) |
| **check** | a deterministic fold-and-compare detects a stored/computed mismatch as an integrity error (§6.4) |
| **no stored join** | no cross-stream join is persisted (K-R13) |
| **no borrowed order** | any S-2 ordering rule is local and does not settle N-c |
| **no lifecycle/context conflation** | context inactivity is not a lifecycle status (S2H-12) |
| **decline is honest** | a decline records O-KR12 as unmet and does not deny Q2 = Yes |

**S2B-12 (Calculated — the completeness limit).** The check verifies that the projection matches the **recorded**
stream; it cannot detect a change that was never recorded. A positive representation improves auditability but does
not make recording infallible. This limit belongs in the decision record.

---

## 6. Decision procedure

**S2B-13 (Calculated — the sequence).**

1. Enumerate the logical representation candidates under B.
2. Test each against S2B-06 and S2B-11.
3. Consider the process-level decline as a substantive outcome.
4. Produce a decision brief with a recommendation or a decline, for the owner.
5. The owner decides; implementation, requirement capture and canonical capture remain separate gates.

**S2B-14 (Observed — no presumptions).** The procedure does not presume a positive carrier. It also does not presume
a decline. The investigation must be able to conclude either.

---

## 7. What remains unauthorized

**S2B-15 (Observed).**

| Untouched | State |
|---|---|
| **carrier** | none adopted |
| **schema / kind / field / key** | none created or implemented |
| **state machine / UX** | not designed |
| **authority / validity / causes** | not decided |
| **N-b / N-c / N-a** | not decided |
| **implementation / tests / probes** | **0** |
| **requirement capture** | **0**; map **118** |
| **canonical capture** | **0** |

---

## 8. Traceability

**S2B-16 (Observed).**

| What | Where |
|---|---|
| Owner selection B | owner, 2026-11-20 |
| Authorization request | S2R-01…S2R-35 |
| Howl separation addendum | S2H-01…S2H-22 |
| A-6 binds N-a, not S-2 by default | NAE-31/32; KRO-17; S2R-19…22 |
| C3D-11 under A-6; B changes the bar | C3D-11; S2B-07 |
| K-R12 / K-R14 / §6.4 / K-R2 / K-R13 / K-R15 | r33 `:83`, `:111`–`:113`, `:366`–`:367` |
| Non-absorption | S2R-17; S2A-38 |
| Separate gates | S2R-25 |
| Context inactivity ≠ lifecycle status | S2H-12 |

---

## 9. Record status

**S2B-17 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier adopted | **0** |
| S-2 open | **yes** — logical process only |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **authorization and baseline record** |

**End of baseline.** **B is the owner's selection; S-2 is open** (S2B-01). It is authorized to investigate the
**logical** representation of post-formation `scope`/`knowledge_class` changes, including new kinds or fields,
projection/reconstruction, the integrity check, entry representation, a local ordering rule, and a process-level
decline (S2B-02). It inherits K-R12, K-R14, §6.4, K-R2's integrity principle, K-R13, K-R15, §11/D1/D5/K-R8/K-R9/
K-R10/I13, the r33/v9 axis separation, and S2H-12 (S2B-06). It does not decide authority, validity, causes, N-b, N-c,
N-a, §11, O-8a, the state machine, UX, implementation, requirements or canon (S2B-03/10). **No carrier is adopted**
(S2B-05).
