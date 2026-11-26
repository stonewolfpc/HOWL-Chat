# MEMORY-09 — J3 scope determination

**Status: narrowly scoped answer to one question — the scope and allocation of requirement J3. Additive; no
code, no map edit, no schema proposal, nothing executed.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Question | Does J3 require MEMORY-09 to persist the recipient set, or only the system boundary to preserve it? |
| Answer | **J3 is engine-scoped and requires the claim event, with content and recipients, to exist inside the engine's scope.** Option C fails J3 as written |
| Code changed | **None** — no schema proposed either |
| Executed | **Nothing** |

---

## 1. The requirement, verbatim, with its allocation

**Map, §J, line 137 — the complete row, all four columns:**

> | J3 | Voluntary disclosure | Claim event with content and recipients | Disclosure inferred; event missing |

**Allocation:**

| Property | Value | Source |
|---|---|---|
| ID | `J3` | map line 137 |
| Section | `## J — Claims, disclosure, provenance, conflict (16)` | map line 131 |
| Ordinal | 3rd of 16 in J | map rows J1–J16, lines 135–150 |
| Set | 118 ratified requirements (A 5 · B 11 · C 9 · D 10 · E 5 · F 9 · G 3 · H 10 · I 3 · J 16 · K 27 · L 10) | map line 12 |
| Requirement column | "Voluntary disclosure" — the scenario, not an obligation | map line 137 |
| Expected column | "Claim event with content and recipients" | map line 137 |
| Failure column | "Disclosure inferred; event missing" | map line 137 |

The row carries no qualifier naming an author, a layer, or an external input. That absence is the whole
question, and it is answered by scope rather than by the row's own words.

---

## 2. Scope: whose conformance is J3?

Three ratified statements, each explicit:

**Map, Scope row (line 11):**

> | Scope | Ratified requirements (**O-6b**) for **experimental engine semantics**. **Production Howl not inspected** |

**Ledger, O-6b (CONFIRMED), lines 470–479:**

> **Decision (owner, this session): v9 and its 118-test conformance map are ratified as the design
> baseline for MEMORY-09's experimental engine, subject to tiers 1 and 2 of O-6a.**
> …
> - **The 118-test conformance map**, as **defined requirements**. Ratification is of the
>   *requirements*, not of any result: **0 tests implemented, 0 executed, 0 passing.**
> - **Scope:** experimental engine semantics (`eval/mem1` target). Production Howl was not inspected
>   and is not covered.

**Ledger, Gate 2 authorisation, lines 558–559:**

> **Implementation must conform to:** the O-6a precedence hierarchy · all **57** ratified architecture
> rules · the **118** O-6b conformance requirements as the governing requirements baseline …

**Established: J3 is a requirement on the experimental engine.** It is not a statement about the
conversation-to-memory boundary in the abstract; it is one of the 118 rows whose subject is `eval/mem1`.
An engine that produces no claim event does not satisfy J3, whatever the layer upstream does.

**Check for contrary evidence.** The map does name the other layer once — `H8` (line 119): "Narration reports
an established disclosure | **Reporting an event the conversation layer established is permitted**; no access to
unrelated private content follows | Reporting refused, or reporting grants access". So the map is willing to say
"the conversation layer" where the fact belongs there. J3 does not say it, and H8's expectation still tests the
engine's own behaviour. **No clause in the map or in O-6b places J3 outside the engine.**

---

## 3. What the row's failure column forces

The map states its own convention (line 13): "Every row states requirement · expected outcome · failure
condition." The **failure column is the operative test** — a row fails if its condition is met, and J3's
condition is:

> **"Disclosure inferred; event missing"**

Two distinct failure conditions joined by a semicolon, and the second is unconditional: **an engine in which the
event is missing fails J3.** There is no reading in which "event missing" is acceptable so long as some other
component holds it, because the row is evaluated against the engine.

The expected column adds the substance: **"Claim event with content and recipients"** — the event must carry
both the content and the recipient set. Recipients are not implied by the accounts, not derivable by grouping,
and not optional: they are named as event content.

**Established: J3 requires the claim event, including its recipient set, to be represented inside the engine's
scope.**

**Not established by J3:** that the engine *decides* who was exposed. "Persist" and "determine" are different
obligations, and J3 speaks only to what must exist. The negative side — that the engine must never infer
exposure — is carried by `J4` ("No actual disclosure event → No accounts created; fails if: Accounts from
intent, proximity or silence") and `J5` ("Exposure governs acquisition"), both engine-side and unaffected by
who authors the event.

---

## 4. Consequence for A, B and C

| Option | J3 as written | Amendment needed? |
|---|---|---|
| **A** — `disclose()` creates the claim event | Satisfies the expected column: an event with content and recipients exists in the engine | No |
| **B** — the conversation layer establishes and supplies the event; the engine persists the claim and the recipient set | **Satisfies the expected column identically.** The row names no author, so caller-supplied establishment is admissible; the engine holds the event, decides nothing, and infers nothing | **No** |
| **C** — the event stays entirely outside MEMORY-09 | **Fails.** The engine's store contains no claim event, which is precisely J3's second failure condition | **Yes** — and a scope amendment is not a conformance result |

**This is the substantive answer to your question: J3 requires MEMORY-09 to hold the recipient set, and B is the
arrangement that does so without giving MEMORY-09 authority over exposure.** C is not available without changing
a ratified requirement's applicability, because J3 is engine-scoped by the map's own Scope row and by O-6b's
"baseline for MEMORY-09's experimental engine".

Note also that B keeps A and B distinct in exactly the way you asked: the conversation layer remains the
authority for occurrence, content and recipients; the engine is the persistence mechanism. Calling a persistence
interface does not make the engine the authority that determined the event — and J3 does not require it to be.
**A versus B therefore remains undecided by J3**; it is decided by rule 46, §1 and §7 step 2, which the previous
supplement recorded as unresolved. J3 narrows the *shape* without settling the *author*.

---

## 5. If C were nonetheless chosen, the instrument required

There is a precedent, and it shows both the mechanism and its limits — **D-4, Option B, scope determination**
(ledger 590–625):

> **Decision (owner, 2026-11-17): MEMORY-09 is absent from the transfer path.** No destination-write
> interface exists. … This is a **scope determination for the MEMORY-09 experiment**, **not** a declaration
> that the Howl application has no transfer responsibility.

> **Owner's wording, verbatim:** These guarantees are not enforceable by MEMORY-09 under the ratified
> experiment scope. Their implementation and enforcement within Howl Chat remain outside this decision and
> unresolved.

> **Untestable within the ratified set:** the 118 requirements cannot test whether Howl Chat fulfils these
> guarantees. That would require separate **application-level** requirements and tests, which O-6b closed.
> These are not defects in MEMORY-09; they are guarantees whose enforcement falls outside its authorised scope.

> - **No new conformance requirements or probes** — the set stays at **118**.

Four properties of that precedent matter here:

| Property | D-4 | If applied to J3 |
|---|---|---|
| Authority | An explicit **owner** scope determination | Would be required; an implementer cannot reallocate |
| Effect on the requirement | Marked unenforceable within MEMORY-09's scope, **gap recorded** | J3 would be recorded as out of scope — not as satisfied |
| Effect on the count | "the set stays at **118**" | Same; no new requirement could be added to cover it |
| Effect on classification | The guarantees stay unresolved | J3 would remain **P**, and would become a boundary gap rather than a pass |

So C is *admissible only as an owner scope determination* — never as a reinterpretation of J3 — and even then it
does not produce a pass. That is the precise reason your instruction not to "call it satisfied" is correct.

---

## 6. Classifications, unchanged

| Row | Classification | Note |
|---|---|---|
| **J3** | **P** | Unchanged. Now with a sharper reason: the event and its recipient set must exist in the engine, and they do not |
| **Rule 46** | **P** | Unchanged; authorship unresolved |
| Rules 56 I · 1 P (46) | unchanged | — |
| Requirements 113 I · 3 P (J3, K9, K10) · 2 NA (B3, F7) | unchanged | — |
| K9 | **P** | Unchanged; ordinal, not a timestamp |
| F7 | **NA** | No replay conformance claimed |

The consolidated review remains a **source-conformance assessment**. It is not, and is not to be read as, a
passing conformance result: nothing has executed, so no row has a result at all.

---

## 7. What this document does not do

| Not done | Note |
|---|---|
| No code change | `eval/mem1/m09/` untouched |
| No schema proposal | No field, key, parameter or registry is designed here — per instruction |
| No implementation of A, B or C | Not implemented, not recommended for implementation yet |
| No map edit | J3's row, section and counts are unmodified; the map is not amended to accommodate any reading |
| No reallocation performed | The J3 scope finding is an interpretation of the map's existing scope, not a change to it |
| No test, probe, compile, lint or runtime check | Nothing executed |

---

## 8. Verification

| Item | Value |
|---|---|
| Evidence read | map lines 11–13 (scope, counts, convention), 131–150 (§J), 119 (H8), 201–212 (review answers); ledger 468–480 (O-6b), 550–559 (Gate 2 baseline), 590–625 (D-4) |
| Count check | A 5 · B 11 · C 9 · D 10 · E 5 · F 9 · G 3 · H 10 · I 3 · J 16 · K 27 · L 10 = **118** |
| Sources | map `b4ba33cc…`, ledger `da3e7239…`, v9 `698254ca…`, REV 13 `5b745628…` — unchanged |
| Implementation | `authority.py f3445823…`, `engine.py 706d5f6c…`, `records.py f211e6ce…`, `errors.py 48125d7e…`, `retrieval.py cd2d380b…`, `__init__.py 06ccf8f6…` — unchanged |
| Executed | **nothing** — every statement above is a reading of ratified text |
