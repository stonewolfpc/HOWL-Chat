# MEMORY-09 — S-2 Representation Investigation

## A positive logical representation is justified; P-1 is recommended, not adopted

**Status:** S-2 investigation under the owner's **B** authorization. **Logical only. No carrier is adopted. No schema,
code, migration, test, requirement, map change or canonical capture.** The recommendation is for the owner's decision;
it is not an implementation and not a final answer.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | how can a post-formation `scope` or `knowledge_class` change be represented in the lifecycle stream so that O-KR12, K-R14 and §6.4 are satisfied? |
| Candidate space | P-1 dedicated change entries; P-2 generic change entry; P-6 piggyback; P-4 external join; P-3 null; P-5 storage |
| Finding | a positive representation is justified; **P-4 excluded**, **P-6 inadequate for pure changes**, **P-3 is a legitimate decline but leaves O-KR12 unmet**; **P-1 and P-2 are viable** |
| Recommendation | **P-1** — dedicated `scope`/`class` change entries with a local per-memory/per-dimension version; **not adopted** |
| Not decided | exact kind names, fields, physical schema, implementation, concurrency, S-1 fold, N-c, requirement capture, canonical capture |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Question and baseline

**S2I-01 (Observed — the authorization).** The owner selected **B**: A-6 does not carry into S-2. S-2 may investigate
logical representations, including new kinds or fields, and may conclude that no admissible representation is
currently justified (S2B-01…05).

**S2I-02 (Calculated — the obligation).** O-KR12(d): for every actual class/scope change, a lifecycle-stream record
determines the new value, and the current value is a projection of the stream (KRO-06). The stream must be
independently rebuildable (K-R14) and the projection checkable (§6.4).

**S2I-03 (Observed — the inherited fences).** K-R12, K-R14, §6.4, K-R2's integrity principle, K-R13, K-R15,
§11/D1/D5/K-R8/K-R9/K-R10/I13, the r33/v9 axis separation, and S2H-12 (S2B-06). No authority, validity, cause,
N-b, N-c, N-a, state-machine or UX decision is in scope (S2B-03/10).

---

## 2. Candidate space under B

**S2I-04 (Calculated — the logical candidates).**

| Candidate | Form | Status under B |
|---|---|---|
| **P-1** | dedicated lifecycle change entries — one per dimension, or a pair (`scope_changed`, `class_changed`) | logically admissible; needs a local order |
| **P-2** | one generic change entry with a dimension field | logically admissible; needs a local order |
| **P-6** | piggyback a class/scope value on an existing status entry | admissible in form, but **inadequate for pure changes** |
| **P-4** | an external operation joined at projection | **excluded by K-R14** |
| **P-3** | null carrier — formation values only | admissible, but **declines** O-KR12 on a change |
| **P-5** | stored projection/cache | **storage, not a carrier** |

**S2I-05 (Calculated — what B changed).** Under A-6, P-1, P-2 and P-6 were barred as new kinds/fields. B lifts that
bar for **logical consideration**. K-R14 is unchanged and still excludes P-4; K-R15 is unchanged and still supplies no
lifecycle order (S2B-07).

---

## 3. Evaluation criteria

**S2I-06 (Calculated — the tests).**

| # | Criterion | Source |
|---|---|---|
| 1 | the carrier is in the lifecycle stream | K-R12 locus |
| 2 | it covers both `scope` and `knowledge_class` changes | K-R12 domain |
| 3 | the current value is derivable from the stream | K-R12 projection |
| 4 | the value is not a mutable authoritative field | K-R12 non-mutation |
| 5 | the stream is independently rebuildable | K-R14 |
| 6 | the check is deterministic and detects mismatch | §6.4 |
| 7 | no stored cross-stream join | K-R13 |
| 8 | ordering is local and does not settle N-c | K-R15; S2H |
| 9 | no authority, validity or cause is decided | §11/D6/D5/K-R9/K-R10/I13 |
| 10 | context inactivity is not lifecycle status | S2H-12 |

---

## 4. Candidate evaluations

### 4.1 P-4 — excluded

**S2I-07 (Calculated).** An external operation joined at projection would require the memory projection to consult
another stream. K-R14 requires independent rebuild and says the cross-stream join has no stored form; §6.4 says
neither projection may consult the other. **P-4 fails before any schema question** (C3D-13; KRO-25). B does not change
this.

### 4.2 P-6 — inadequate for pure changes

**S2I-08 (Calculated).** P-6 attaches a class/scope value to an existing status entry. It can represent a class/scope
change only when a status change also occurs. A **pure** scope or class change — one with no status change — would
have no entry to piggyback on. Creating a synthetic status entry to carry it would falsely change lifecycle status
(e.g., `corrected` when nothing was corrected).

**S2I-09 (Calculated — why that is fatal to P-6 as a general carrier).** O-KR12 is universal over actual changes. A
carrier that cannot represent a pure change fails the domain condition (criterion 2). P-6 could be a **partial**
mechanism for coincident changes, but not the general representation. It also risks conflating status and class/scope,
which K-R12 lists separately.

### 4.3 P-3 — a legitimate decline

**S2I-10 (Calculated).** P-3 (formation values only) is buildable and complete **if no change occurs**. It fails
O-KR12 when a change occurs. Under B, it is not the only option; it is the **decline** option. A process-level decline
is legitimate and must remain available (S2H-15), but it leaves the acknowledged obligation unmet.

### 4.4 P-5 — storage, not a carrier

**S2I-11 (Calculated).** P-5 caches the current value for the §6.4 check. It needs an authoritative input and a
rebuild rule; it does not supply either. It can accompany P-1 or P-2, but cannot substitute for them (C3D-05…09).

### 4.5 P-1 — dedicated change entries

**S2I-12 (Calculated — the form).** The lifecycle stream gains change entries for the two dimensions — either one kind
per dimension (`scope_changed`, `class_changed`) or a single kind with a dimension field (which would be P-2). Each
entry carries the memory identity, the new value, and a local version. The status kinds remain unchanged.

**S2I-13 (Calculated — evaluation).**

| Criterion | P-1 |
|---|---|
| 1 in-stream | **yes** |
| 2 both dimensions | **yes** |
| 3 projection | **yes** — formation value, then latest change |
| 4 non-mutation | **yes** — stored value is a cache |
| 5 independent rebuild | **yes** — from the lifecycle stream alone |
| 6 check | **yes** — fold and compare; mismatch is an integrity error |
| 7 no stored join | **yes** |
| 8 local order | **yes** — see §5 |
| 9 no authority/validity/cause | **yes, if fenced** — see S2I-20 |
| 10 context ≠ status | **yes, if fenced** — see S2I-21 |

### 4.6 P-2 — generic change entry

**S2I-14 (Calculated — the form).** One change kind with a `dimension` field (`scope`/`class`) and a new value, plus
a local version. It is more compact and extensible than P-1, but polymorphic: the dimension and value must be
validated together.

**S2I-15 (Calculated — evaluation).** P-2 satisfies criteria 1–10 in the same way as P-1, provided the dimension field
is validated and the value is checked against the existing `scope`/`class` vocabulary. Its advantage is fewer kinds;
its cost is a polymorphic entry and a broader validation surface.

---

## 5. The ordering problem

**S2I-16 (Calculated — the gap).** K-R12 requires a projection, but the lifecycle stream has no ordering rule
(DS-09; ND-1). K-R15 fixes `(turn, seq)` for the **acquisition** fold only; it does not supply a lifecycle order.
A positive carrier must not borrow K-R15, and must not silently settle N-c.

**S2I-17 (Calculated — the local solution).** The change entry carries a **per-memory, per-dimension version**. The
projection applies the entry with the highest version for each dimension. The version:
- is local to the dimension, not a stream-wide status order;
- does not order the four status kinds;
- does not decide whether the status stream needs an order key (N-c);
- can detect a gap if versions are sequential (a missing version is an integrity error).

**S2I-18 (Calculated — why this is not N-c spillover).** N-c asks whether the **status** stream needs an order key.
The version orders only the class/scope change entries. The status projection remains governed by its own unresolved
ordering question (S-1/N-c). The two orders are separate. If a future decision unifies them, that is a separate
decision, not an implication of S-2.

**S2I-19 (Unestablished — alternatives).** A predecessor link (each change points to the prior change) is an
alternative to a numeric version. It avoids a counter but adds a link field and a chain-integrity check. The choice
between version and predecessor is not settled here; both are local and both avoid N-c.

---

## 6. Projection, check and rebuild

**S2I-20 (Calculated — the projection rule).** For each memory and dimension:
1. start from the formation value (the value written at formation);
2. find the change entry with the highest local version for that dimension;
3. if none, the current value is the formation value;
4. otherwise, the current value is that entry's new value.

**S2I-21 (Calculated — the check).** Fold the lifecycle entries for the memory; apply the rule in S2I-20; compare the
result to the stored current value. A mismatch is an integrity error surfaced to `inspection`, never silently resolved
(§6.4). If versions are sequential, a gap is also an integrity error.

**S2I-22 (Calculated — rebuild).** The projection is reconstructable from the lifecycle stream alone; no acquisition
event, external operation or stored join is consulted (K-R14; K-R13). The stored current value is a cache, not an
authority.

**S2I-23 (Calculated — the completeness limit).** The check verifies that the projection matches the **recorded**
stream. It cannot detect a change that was never recorded. A sequential version can reveal an internal gap, but not a
missing **latest** change. Recording is a write-path obligation; the representation improves auditability, not
infallibility. This limit must be stated in the decision record.

---

## 7. Fences on authority, validity, causes and context

**S2I-24 (Calculated — authority).** A change entry may record provenance (who or what recorded the change) for audit,
but it must not define who is **permitted** to change class/scope. Permission remains §11/C-1, outside S-2.

**S2I-25 (Calculated — validity).** The entry records a new value; it does not assert that the change is warranted or
true. Validity remains D6 and the unresolved predicate, outside S-2.

**S2I-26 (Calculated — causes).** The entry may record a cause for audit, but it must not define which causes are
permitted. D5/K-R9/K-R10/I13 continue to constrain causes; S-2 does not relax or extend them.

**S2I-27 (Calculated — context).** The entry represents a lifecycle class/scope value, not context activation. A memory
being inactive in a scene is not a lifecycle status (S2H-12). The state machine decides participation; S-2 does not.

---

## 8. Recommendation

**S2I-28 (Calculated — the recommendation).** **P-1**: dedicated change entries, one per dimension, with a local
per-memory/per-dimension version. This is recommended because:
- it keeps `scope` and `class` distinct, with no polymorphic dimension field;
- the projection and check are straightforward: filter by kind, take the highest version;
- validation is simpler: each kind has one value domain;
- the four status kinds remain untouched, avoiding status/class conflation;
- it satisfies all ten criteria in S2I-06.

**S2I-29 (Calculated — P-2 remains viable).** A generic change entry is not rejected. It has fewer kinds and is more
extensible, at the cost of a polymorphic field and a broader validation surface. If extensibility is prioritized, P-2
is a legitimate alternative. **The choice between P-1 and P-2 is the owner's, not this record's.**

**S2I-30 (Observed — not adopted).** P-1 is **recommended, not adopted**. No carrier has been selected. The
recommendation is a decision input, not an implementation and not a final answer.

---

## 9. Decline assessment

**S2I-31 (Calculated — decline is legitimate).** A process-level decline is a completed substantive decision: it
records O-KR12 as unmet, preserves Q2 = Yes, and leaves the representation question open (S2H-15). It is not a failure
to decide.

**S2I-32 (Calculated — decline is not required here).** A positive representation is possible under B and satisfies
K-R12/K-R14/§6.4. Decline would leave an acknowledged obligation unmet even though an adequate representation is
available. The investigation therefore does not recommend decline, while preserving it as a valid owner outcome.

**S2I-33 (Calculated — the state machine does not force the choice).** The recommendation is grounded in the memory
contract's own obligation, not in an assumed state-machine design (S2H-07/08). Howl's future needs may inform the
contract's sufficiency, but they do not require a positive carrier here.

---

## 10. What remains open

**S2I-34 (Calculated — the open set).**

| Open | Why |
|---|---|
| P-1 vs P-2 | a design trade-off; owner decision |
| version vs predecessor link | both avoid N-c; not settled |
| one kind per dimension vs a shared change kind | P-1 vs P-2 |
| exact field names and physical schema | implementation, not S-2 |
| concurrency and serialization of versions | write-path design |
| whether to carry prior values for audit | optional; not required by the projection |
| S-1 fold scope | separate downstream decision |
| N-c | status ordering; untouched |
| requirement capture / map | separate gate |
| implementation / tests | separate gate |

---

## 11. Acceptance test

**S2I-35 (Calculated — the logical walk-through, not code).**

1. A memory is formed with `scope=scenario`, `knowledge_class=experiential`.
2. A `scope_changed` entry with version 1 records `scope=world`.
3. The projection yields `world`; the stored cache should be `world`; the check passes.
4. Rebuilding from the lifecycle stream alone yields `world` (K-R14).
5. If the cache says `scenario` while the stream yields `world`, the check reports an integrity error (§6.4).
6. If version 1 is missing while version 2 exists, the sequential check reports a gap.
7. No acquisition event or external operation is consulted (K-R13).
8. The entry records the value; it does not authorize, validate or define the cause of the change.

**S2I-36 (Observed — no test executed).** The walk-through is logical. No code, schema, probe or model call was run.

---

## 12. What this does not do

**S2I-37 (Observed).**

| Untouched | State |
|---|---|
| **carrier adoption** | **0** — P-1 recommended only |
| **schema / kind / field / key** | none created or implemented |
| **authority / validity / causes** | fenced, not decided |
| **N-b / N-c / N-a** | not decided |
| **§11 / O-8a / v9 domain** | not decided |
| **state machine / UX** | not designed |
| **implementation / tests** | **0** |
| **requirement capture** | **0**; map **118** |
| **canonical capture** | **0** |

---

## 13. Traceability

**S2I-38 (Observed).**

| What | Where |
|---|---|
| B authorization; baseline | S2B-01…S2B-17 |
| Request; addendum | S2R-01…S2R-35; S2H-01…S2H-22 |
| O-KR12 | KRO-06 |
| K-R12 / K-R14 / §6.4 / K-R2 / K-R13 / K-R15 | r33 `:83`, `:111`–`:113`, `:366`–`:367` |
| P-1/P-2/P-6/P-4/P-3/P-5 | C3D-10/11/13/14; C3P-10 |
| P-4 excluded by K-R14 | C3D-13; KRO-25 |
| P-6 nearest positive under A-6, now evaluated | C3D-14; S2I-08 |
| Local order; N-c separation | S2I-16…S2I-19 |
| Context inactivity ≠ lifecycle status | S2H-12 |
| Map unchanged | 118 rows |

---

## 14. Record status

**S2I-39 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier adopted | **0** — P-1 recommended only |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **investigation** — no canonical capture |

**End of investigation.** A positive logical representation is justified (S2I-28…S2I-33). **P-4 is excluded by
K-R14; P-6 cannot represent pure changes; P-3 is a legitimate decline but leaves O-KR12 unmet; P-5 is storage.** **P-1
and P-2 are viable; P-1 is recommended** — dedicated `scope`/`class` change entries with a local per-memory/per-
dimension version, which avoids N-c and satisfies K-R12/K-R14/§6.4 (S2I-12…S2I-29). **P-1 is not adopted; the owner
decides** (S2I-30). **No schema, code, requirement, map change or canonical capture follows** (S2I-37/39).
