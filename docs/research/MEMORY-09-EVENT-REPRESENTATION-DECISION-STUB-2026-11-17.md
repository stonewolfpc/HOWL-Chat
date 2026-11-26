# MEMORY-09 — Event representation and interface: decision stub

**Status: all three decisions CLOSED. D-E2 = B2 (owner, §2b); D-E1 = the representation in §2c (owner authorisation
of 2026-11-17, implemented); D-E3's operative form = §5.3. Implementation exists in `eval/mem1/m09/`; **nothing has
been executed**, and this is not a conformance result.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Purpose | Fix the boundary before any implementation authorisation, so no field is added accidentally, no authority expands accidentally, and the engine's permitted behaviour is unambiguous |
| Authority held by this document | **None.** It records and separates; the owner decides |
| Code changed | **None** |
| Schema proposed | **None** — the representation's form is open decision D-E1 and is deliberately unspecified here |
| Executed | **Nothing** |
| Prior records | Unchanged; this stub is additive |
| Amendments | 2026-11-17 — §3 D-E3 and §5 recast to distinguish the integrity contract per interface. Then: **D-E2 closed by owner decision (B2)**; §2b added, §3/§5/§6 updated. Then: **D-E1 closed by owner authorisation and implemented**; §2c added, §3/§6/§7/§8 updated |

---

## 1. Governing principle

> **MEMORY-09 must preserve the conversation layer's established recipient set. It must not determine it.**

Everything in this stub follows from that sentence. It resolves what looked like an authority conflict: the
engine needs the recipient set *inside its own store* (J3, engine-scoped) without acquiring any authority over
*who was exposed* (rule 11, J4, J5).

**Settled facts, carried forward from the J3 scope determination and the ownership reconciliation:**

| # | Fact | Basis |
|---|---|---|
| 1 | The engine must hold the disclosure claim event | J3 expected column; engine-scoped by map Scope and O-6b |
| 2 | That event must contain the communicated content **and** the exposed-recipient set | J3 "Claim event with content and recipients" |
| 3 | The conversation layer establishes whether disclosure occurred, what was communicated, and who was exposed | v9 §1; §7 step 2 |
| 4 | MEMORY-09 must not infer exposure or independently decide that disclosure occurred | rule 11; J4; v9 §1 prohibitions |
| 5 | Option C leaves J3 unsatisfied; it would require an owner scope determination with J3 recorded as a boundary gap, never a pass | D-4 precedent; J3 failure condition "event missing" |
| 6 | J3 does not decide whether `disclose()` creates the event or the layer supplies an established event for persistence | rule 46 / §1 / §7 step 2 remain unresolved as to authorship |

---

## 2. The two candidate interfaces

Both satisfy J3's expected column. Neither is authorised, and neither is implemented.

| Responsibility | **B1 — established-event handoff** | **B2 — established-facts ingestion** |
|---|---|---|
| Establish disclosure | Conversation layer | Conversation layer |
| Supply content and recipients | Existing event | Explicit inputs |
| Write claim record | Conversation layer / event persistence path | MEMORY-09 |
| Validate recipient/account consistency | MEMORY-09 | MEMORY-09 |
| Create acquired accounts | MEMORY-09 | MEMORY-09 |
| Decide exposure | **Never MEMORY-09** | **Never MEMORY-09** |

**Current implementation, stated so the gap is not mistaken for a capability:** neither path carries the
recipient set onto a claim. `claim()` fixes provenance at `{contributor, contributor_kind, span, turn}`
(`engine.py:185-188`). `disclose()` derives content from the source record (`engine.py:573-581`) and builds
`event = {"kind": "disclosure", "turn": turn}` (`:583`), then writes N accounts and no claim (`:584-590`).

---

## 2b. Owner decision — B2 (APPROVED 2026-11-17)

**Decision text, verbatim. The four numbered statements are the decision; the flow beneath them is illustrative.**

> 1. Conversation layer establishes what happened. It knows who disclosed something, what was communicated, and
>    who received it.
> 2. MEMORY-09 records that established event. It stores one disclosure claim containing the communicated content
>    and the recipient set as structured metadata.
> 3. MEMORY-09 creates the recipient accounts. It uses the exact same supplied recipient set to create each
>    recipient's acquired account.
> 4. The recipient set is never reconstructed or inferred. The conversation layer is authoritative about who
>    received the disclosure.

```
Conversation layer  ->  establishes disclosure facts
                    ->  supplies established event facts
                            |
MEMORY-09           ->  stores disclosure claim + recipient metadata
                    ->  same supplied recipient set
                    ->  one acquired account per recipient
```

**Settled architecturally:**

| # | Settled | Effect on this stub |
|---|---|---|
| 1 | MEMORY-09 writes the claim record; the layer supplies the facts | D-E2 = **B2** |
| 2 | One supplied set serves both the claim's metadata and account creation | D-E3's operative form = **§5.3**, single source |
| 3 | No reconstruction or inference of the recipient set, ever | §5.3's prohibitions are binding |
| 4 | The recipient set is **structured metadata on the claim** | D-E1's architectural shape settled; its concrete representation still open |
| 5 | No MEMORY-09 determination of exposure | Unchanged from §1 fact 4 |

**Not settled, in the owner's words:** "We haven't approved field names, schema, API signatures, or implementation
work." Those belong to a separate implementation authorisation, and the no-new-fields restriction still stands
against them.

**Consequence recorded, not implemented:** under B2 the communicated content is an explicit input, so
`disclose()`'s current content-definition-by-derivation (`engine.py:573-581`, e.g. `f"{discloser} believes:
{src.text}"`) is superseded architecturally. It must be replaced when implementation is authorised; it has not
been changed.

---

## 2c. D-E1 CLOSED — the representation and the implementation (2026-11-17)

Authorised by the owner on 2026-11-17 to *determine and implement* D-E1 and B2, superseding the no-new-fields
restriction only as far as this architecture requires (§3 of that authorisation). Implemented; **not executed**.

### 1. Representation chosen

`Record.disclosure: dict | None`, **claims only**:

```
{"recipients": ("bram", "tamsin"), "source_class": "belief"}
```

| Property | Choice | Why |
|---|---|---|
| Placement | Field on `Record`, beside the other class-specific structures (`origin_class`/`via`, `contribution_kind`/`contributor`, `attribution`) | The record model already carries class-specific structured metadata; a registry or side table would have been a second store of one event |
| Marker | `disclosure is not None` | "No disclosure claim may silently omit its recipient metadata" becomes structural: presence *is* the marker, so no reader infers disclosure-ness from text or provenance |
| Recipients | `tuple[str, ...]`, exactly as supplied, duplicates refused | Immutable, and no canonicalisation that could be read as reinterpretation. Order is the layer's |
| Source | `source_class` only (`claim`/`account`/`belief`) | Records what kind of private record was disclosed without creating a resolvable link (rule 12 / §8). The source **id is not stored** |
| Visibility | `as_dict()` emits it, so `snapshot()` carries the event | J3 observability without adding a view, a count or a response field |
| New classes | None | No new record class, no event registry, no second store |

### 2. Interface — B2

```
Engine.disclose(*, discloser, content, recipients, source_record_id, turn,
                span=None, kind=None) -> (claim: Record, accounts: tuple[Record, ...])
```

| Input | Role | Notes |
|---|---|---|
| `content` | The communicated content, established by the conversation layer | **Required. Never derived.** The old derivation (`f"{discloser} believes: {src.text}"` and siblings) is removed |
| `recipients` | The established recipient set | Required. Used verbatim, **once**, for the claim metadata and for account creation |
| `discloser`, `turn` | Attribution of the event | Required |
| `source_record_id` | Ownership and containment check (rule 46, 31) plus `source_class` | Required. Read, never stored as a link |
| `span` | Claim provenance only | Optional; not copied onto recipient accounts |
| `kind` | Declared contribution kind | Optional; absent records `Unspecified claim` (rule 26) |

### 3. Refusals — all decided before the first write

| Condition | Error | Rule |
|---|---|---|
| `content` absent, non-string or blank | `DisclosureError` | 46 / B2 |
| `recipients` empty | `DisclosureError` | 11 / §5.3 |
| A recipient is not a non-empty string | `DisclosureError` | 11 |
| A recipient named twice | `DisclosureError` | §5.3 (no silent dedupe) |
| The discloser is among the recipients | `DisclosureError` | 46 / C-3 `told` |
| Source record absent | `CitationError` | 33 |
| Source in another chat | `DisclosureError` | 31 |
| Discloser does not own the source | `DisclosureError` | 46 |
| Source is domain truth | `DisclosureError` | 46 / §7 |
| Unknown contribution kind | `ContributionKindError` | 26 |
| A recipient cannot be `told` (evidence pre-flight) | `EvidenceError` | C-3 |
| Malformed `disclosure` metadata anywhere, including direct construction | `DisclosureError` | 46 / J3 |

### 4. Atomicity

1. Every refusal above is evaluated before the claim is constructed.
2. `_validate_evidence` is run for **every** recipient before any write, so the account loop cannot refuse halfway.
3. The claim and all accounts ride one enclosing `_mutation`, so the version advances once (S-7).

A refused disclosure therefore leaves **no claim and no account**. No rollback journal was added; the ordering is
the mechanism, and it is commented as load-bearing in the source. Observable consequence: refusals now occur
outside `_mutation`, so a *refused* call no longer advances `version`.

### 5. Files and functions changed

| File | Change |
|---|---|
| `records.py` | `Record.disclosure` field; validation block in `__post_init__`; `disclosure` emitted by `as_dict()`; `DisclosureError` imported |
| `engine.py` | `_DISCLOSER_CONTRIBUTOR_KIND`; `disclose` rewritten (B2 signature, pre-flight, claim + accounts, tuple return); `NON_TRUTH_CLASSES` imported |
| `authority.py`, `errors.py`, `retrieval.py`, `__init__.py` | **Unchanged** |

No caller of `disclose` exists anywhere in the repository, so the signature change breaks no call site.

### 6. Relationship to §5.3

| §5.3 requirement | Where it is satisfied |
|---|---|
| One authoritative supplied set | `supplied` → `recipients_set` (one tuple) → `disclosure["recipients"]` **and** the `acquire` loop. No second list exists |
| No second independently constructed set | `set(supplied)` appears **only** in a length comparison for duplicate detection; it is never stored, never returned and never drives account creation |
| No reconstruction or reinterpretation | Nothing in the account path reads recipients from any record, owner list, group, proximity or provenance |
| Sole source for content | `content` is an argument; the source record's `text` is not read by `disclose` at all |
| Absent or empty set refused | `DisclosureError` before any write |

### 7. Limitations and unresolved items

1. **Rule 57's expression moved.** The old text prefix ("… believes: …") is gone because content is now supplied.
   The rule is still satisfied structurally — the account is non-truth-apt and `descriptor.record_class` names the
   source class — but the engine no longer phrases the account, so the layer must. If the architecture wants the
   prefix, that is a further decision.
2. **Descriptor event token unchanged** (`{"kind": "disclosure", "turn": turn}`). It is deliberately
   non-resolvable (rule 12 / J8), so two disclosures in one turn remain indistinguishable *from the recipient
   side*. The claim disambiguates for the engine and the discloser.
3. **`contributor_kind` is fixed** at `"character"` by the operation. If a layer must record a different kind,
   that is an additional input and an additional decision.
4. **Atomicity is ordering-dependent.** Any future edit that adds a refusal after the claim is written would break
   it. Recorded as a constraint on future edits, not a guarantee of the type system.
5. **Nothing is executed.** No test, probe, compile, lint or import has run.

### 8. Status, kept separate

| Question | Answer |
|---|---|
| Implemented? | **Yes** — source present in `eval/mem1/m09/` |
| Executed? | **No** — no test, probe, compile, lint or import |
| J3 classification | Still **P**. The representation now exists; the row cannot move to **I** without a fresh source-conformance review, and cannot pass without execution |
| Rule 46 classification | Still **P** for the same reason |

---

## 3. The three decisions — status

Kept separate deliberately. They have different blockers and different owners of risk.

### D-E1 — Representation — **CLOSED: see §2c (owner authorisation, implemented)**

| | |
|---|---|
| **Decision** | `Record.disclosure` on claims: `{"recipients": tuple, "source_class": str}`, presence is the marker (§2c.1) |
| **Authorisation** | Owner, 2026-11-17: determine and implement D-E1 and B2, superseding the no-new-fields restriction *only* as far as this architecture requires |
| **Restriction still in force** | The general prohibition stands outside this scope; the authorisation named one field and its supporting changes, not a general licence |
| **What it must not become** | A general metadata bag, a second recipient source, or a resolvable link to the source record |
| **Effect** | J3's representation now exists in source; the classification stays **P** pending review and execution (§2c.8) |

### D-E2 — Interface — **CLOSED: B2 (owner, 2026-11-17)**

| | |
|---|---|
| **Decision** | **B2 — established-facts ingestion.** The conversation layer supplies the established facts; MEMORY-09 writes the claim and the accounts (§2b) |
| **Resolves** | Rule 46's authorship question in the architectural sense: the engine writes the record, the layer remains the authority for the facts |
| **Does not resolve** | Field names, schema, API signature, refusal specifics, implementation |
| **Retained** | B1 is not deleted; §5.2 records the non-selected alternative so its integrity contract remains auditable |

### D-E3 — Integrity

| | |
|---|---|
| **Question** | What equality and refusal conditions ensure the event and the resulting accounts cannot disagree? |
| **Depends on** | Nothing conceptually — it is a consequence of either interface, not a reason to prefer one. Its **operational expression** does depend on D-E2: two independently supplied sets under B1, one set used as the sole source under B2 (see §5) |
| **Must not** | Be satisfied by recomputation, reconstruction or inference of exposure |
| **Established now** | The prohibition. What differs per interface is only whether integrity is achieved by **comparison** (two sets) or by **construction** (one set) |

---

## 4. B1's unresolved dependency: the persistence path

B1's table row "Write claim record: conversation layer / event persistence path" is not yet a defined path.

| Fact | Consequence |
|---|---|
| The experiment has **one** record store | The layer cannot persist a claim anywhere except MEMORY-09's own store, via `claim()` |
| `claim()` cannot carry a recipient set | B1's write step depends on **D-E1** exactly as B2's does; B1 does not avoid the representation change |
| No second store exists or is described by any ratified document | Inventing one would be new architecture, prohibited by the same restriction that blocks D-E1 |

**Therefore:** B1, as stated, is really "the layer calls the existing claim interface to create the event, then
`disclose()` creates accounts against it". That keeps the authority split intact and is the only form of B1 the
current scope permits. If the owner intends a genuinely separate event store, that is a scope change, not an
interface choice, and it is not authorised by anything in this stub.

---

## 5. Integrity under each interface

> **B2 is approved (owner, 2026-11-17, §2b). §5.3 is therefore the operative contract. §5.2 and §5.4 are retained
> as the non-selected alternative — kept for audit, not for implementation.**

### 5.1 Two supplied sets, or one

An equality check is only meaningful where two sets can disagree. That condition holds for B1 and not for B2, so
the integrity contract splits — not in its purpose, but in its mechanism.

| | **B1 — established-event handoff** | **B2 — established-facts ingestion** |
|---|---|---|
| Recipient sets supplied | **Two, independently**: the event's established set, and the set whose accounts this call creates | **One**: the established set |
| Where integrity comes from | **Comparison** — the two can disagree, and the engine detects it | **Construction** — both the claim and the accounts are built from the same supplied value |
| Integrity contract | Compare the two supplied sets; refuse on mismatch, writing nothing | Use the established set as the **sole source** for both the claim and account creation; no reconstruction and no independent reinterpretation |
| What divergence is detectable | A layer or caller disagreeing with the event | Nothing — it cannot disagree with itself. The engine's only failure mode is reinterpretation, so the prohibition is the contract |

### 5.2 The comparison (B1)

The owner's qualification, recorded verbatim:

> B1's equality check is necessary, but it must compare the authoritative event recipient set against the
> recipients whose accounts are being created — not against a reconstructed or inferred set. Otherwise the engine
> could accidentally manufacture its own interpretation of exposure while supposedly validating the event.

| Property | Requirement |
|---|---|
| Left side (B1) | The **authoritative** recipient set as established and supplied by the conversation layer |
| Right side (B1) | The recipient set whose accounts are being created in this call |
| Operation | Set equality between two **supplied** values |
| Forbidden | Deriving either side from records, from group membership, from proximity, from the source's provenance, or from the accounts already written |
| On mismatch | Refusal. The engine records nothing rather than reconciling two authoritative claims about exposure |
| What it is not | Not a judgement that disclosure occurred; not a determination of who was exposed; not a repair path |

### 5.3 The single-source requirement (B2)

| Property | Requirement |
|---|---|
| Sole source | The established recipient set supplied by the conversation layer |
| Both uses | The claim event's recipient set **and** the accounts created in this call |
| Forbidden | Deriving the set again from records, group membership, proximity, the source's provenance, the source record's owner list, or the accounts as they are written |
| Forbidden, restated | Building the claim from inputs and the accounts from a second reading of the same inputs, where the two readings could differ |
| On absent or empty set | Refusal — an event with no recipients is not a disclosure (rule 11) |
| What it is not | Not a comparison, and not a weaker form of §5.2. It is the absence of a second source, which is a stronger guarantee about engine behaviour and a weaker one about caller behaviour |

**The honest asymmetry:** B1 can catch a caller who contradicts the event; B2 cannot catch a caller who supplies a
wrong set, because the first supply *is* the authority. B2's protection is that the engine never reinterprets;
its exposure is entirely upstream, where the conversation layer already owns the determination.

### 5.4 If B2 is to carry two sets

If B2 is intended to accept the event's set and a separate account-creation set, that must be stated explicitly as
part of **D-E2**, and §5.2's comparison then applies unchanged. The stub does not assume it, and the single-source
requirement in §5.3 is written for the one-set reading only.

**Standing of D-E3:** conceptually established; its operational expression is pending D-E2. It remains a
consequence of either interface rather than a reason to prefer one — but the reason is now specific rather than
blanket, since the two mechanisms protect against different failures.

---

## 6. What this stub does not authorise

| Not authorised | Note |
|---|---|
| Adding any further field, key or record class | The 2026-11-17 authorisation named `Record.disclosure` and its supporting changes, and is now **spent**. The plan v3 §10 / ledger 567 restriction stands outside that scope |
| Extending B2 beyond disclosure | The authorisation covers D-E1 and B2 only; no other operation gained inputs or metadata |
| Treating implementation as verification | §2c.8 keeps implemented, reviewed and executed apart |
| Inventing a second record store | New architecture; prohibited |
| Reallocating J3 to another layer | J3 is engine-scoped; only an owner scope determination could narrow it, and it would still not pass |
| Inference of exposure, in any form | Prohibited by rule 11, J4, J5 and fact 4 above |
| Testing of any kind | Unauthorised throughout |

**What a future implementation authorisation would need to name explicitly:** D-E1's concrete representation; the exact field(s) permitted; the permitted interface signature for B2; the refusal behaviour when the supplied set is absent or empty; and that J3's P may then be re-assessed as implementation progress — never as an executed result.

---

## 7. Classifications, unchanged

| Row | Classification | Reason held |
|---|---|---|
| **Rule 46** | **P** | Event authorship unresolved (fact 6) |
| **J3** | **P** | Representation now implemented (§2c) but **unreviewed and unexecuted**; a source review must re-inspect before the row can move, and only execution could make it pass |
| **K9** | **P** | Turn ordinal is not a timestamp; a persisted write-time timestamp remains permitted but unauthorised |
| **B3** | **NA** | Event-kind declaration boundary, not assessable as written |
| **F7** | **NA** | No replay facility; no replay conformance claimed |
| Rules | **56 I · 1 P (46)** | unchanged |
| Requirements | **113 I · 3 P (J3, K9, K10) · 2 NA (B3, F7)** | unchanged |

**Neither a decision nor a source change moves a classification.** B2 and §2c change what the engine is built to
do; they do not produce evidence that it does it. J3's representation now exists in source, but the row stays **P**
because no review has re-inspected the changed code and no probe has run. Rule 46 stays **P** likewise.

The consolidated review remains a **source-conformance assessment** of the code as it stood on its own date. It
is not a passing conformance result, and it now describes a revision that has since changed (`engine.py`, 
`records.py`). A fresh source review is required before any classification is restated.

---

## 8. Verification

| Item | Value |
|---|---|
| Evidence read for this stub | ledger 516 (plan v3 §10 bullet 2 authority), 567 (no new fields), 468–480 (O-6b), 590–625 (D-4); map line 11 (scope), line 137 (J3); `engine.py:185-188`, `:573-584` |
| Code, after the 2026-11-17 implementation | `authority.py f3445823…` (unchanged), `engine.py fcbc7eb256287ca1…` (was `706d5f6c…`), `records.py 778684e2a63d82d9…` (was `f211e6ce…`), `errors.py 48125d7e…` (unchanged), `retrieval.py cd2d380b…` (unchanged), `__init__.py 06ccf8f6…` (unchanged) |
| Line counts | `engine.py` 1118 → 1160; `records.py` 210 → 242 |
| Execution | **none** — no test, probe, compile, lint or import |
| Contract | v9 `698254ca…`, map `b4ba33cc…`, ledger `da3e7239…`, REV 13 `5b745628…` — unchanged |
| Executed | **nothing** |
