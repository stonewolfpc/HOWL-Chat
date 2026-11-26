# MEMORY-09 — Review supplement: event-layer ownership and K9 time

**Status: narrowly scoped supplement to `MEMORY-09-CONFORMANCE-SOURCE-REVIEW-2026-11-17.md`. It addresses
two review issues only — the rule 46 / J3 disclosure-event finding, and the K9 timestamp claim. It changes no
implementation file, no review document, and no contract. Nothing was executed.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Scope | §4.1 and §4.2 of the consolidated review, and one imprecise sentence of my own |
| Code changed | **None** |
| Documents changed | **None** — this file is additive; the review, report, supplement and clarification stand unedited |
| Executed | **Nothing** — no test, probe, compile, lint or runtime check |
| Trigger | Owner correction request, two items |

---

## Part 1 — Rule 46 / J3: whose event is it?

### 1.1 Exact clauses

**v9 §1, responsibility boundary** (rows relevant to a disclosure):

> | Question | Owner |
> |---|---|
> | Whether a character speaks | **Model** |
> | What the character chooses to disclose | **Model** |
> | Whether disclosure actually occurred | **Conversation / event layer** |
> | What information was communicated | **Conversation / event layer** |
> | Which characters were exposed | **Conversation / event layer** |
> | What was actually said | **Conversation / event layer** |
> | Recording memories and acquisitions | **Memory system** |
> | Preserving provenance | **Memory system** |

> Prohibitions: no second AI model; no AI-based memory arbitration; the memory system must never initiate
> disclosure, infer disclosure from proximity, or manufacture acquisition.

**v9 §7, disclosure lifecycle** (complete):

> 1. **Model** decides the character speaks and what to disclose.
> 2. **Conversation layer** establishes that disclosure occurred, its content and recipients — recorded as a
>    claim event with spans.
> 3. **Exposure** follows from that event alone.
> 4. **Memory system** creates, per exposed character, their own acquired account: source, event/turn,
>    communicated content, origin descriptor.
> 5. **Unexposed characters acquire nothing.**
> 6. **Source private record unchanged**, still exclusively owned.
> 7. **Disclosure establishes neither truth nor belief.**

**v9 §8, privacy and provenance — mechanism:**

> **Mechanism:** the recipient's origin reference is a **descriptor** (speaker, record type, turn), not a
> resolvable identifier. The resolvable link exists only on the source side. Each recipient record has its own
> identity and acquisition provenance; it is never a view, alias or projection.

**Rules (v9 §4), verbatim:**

> 10. **Disclosure** — a private memory becomes shared information only when an actual disclosure event establishes it; the model decides whether and what to disclose.
> 11. **Exposure limits acquisition** — only characters actually exposed acquire; group membership confers nothing.
> 22. **Creation is open** — any participant, human or model, may create conversational information.
> 23. **Recording is system-governed** — the system preserves who contributed what, in which domain, from which perspective, by what route.
> 25. **A statement is a claim** — a character's utterance records that they said it; it never establishes that they experienced what they described.
> 26. **Kind is declared, not judged** — … A missing declaration is recorded as **`Unspecified claim`**; the system never supplies a kind on the model's behalf.
> 46. **Disclosure needs no authorisation** — a character may disclose what it has acquired or remembers through its own conversational output. **The system records the event and its attribution**, creates recipient-specific acquired accounts, and does not expose the source record. It grants no permission to speak.
> 47. **Narration cannot disclose for a character** — narration may report a disclosure event established by the conversation layer, but must not manufacture one, reproduce a private record, or speak private content on the owner's behalf.
> 48. **Disclosure is a claim, not truth** — it creates an attributable claim and a recipient-specific account; the source record is never exposed.

**K9 (map):** "Operation recorded | Actor, authority, time, domain, applicability, source".
**J3 (map):** "Voluntary disclosure → Claim event with content and recipients".

### 1.2 Current implementation evidence

`engine.py:552` — `disclose(self, *, discloser, recipients, source_record_id, turn) -> list[Record]`

| Step | Lines | What happens |
|---|---|---|
| Validate recipients, ownership, containment | 560–571 | before mutation |
| Derive content by source class | 573–581 | `belief` → `"<name> believes: …"`; `account` → text; `claim` → `"<name> said: …"`; truth → refused |
| Build an event token | 583 | `event = {"kind": "disclosure", "turn": turn}` |
| Create one account per recipient | 584–590 | `acquire(via="told", descriptor={"speaker", "record_class", "turn", "event"})` |
| Return | 591 | the accounts only |

**No record is created for the event.** No `claim` is written. The disclosure exists as N accounts plus an
in-line dict token whose identity is `(kind="disclosure", turn=turn)`.

Interface slots that already exist and are unused by this path:

| Slot | Signature | Used by `disclose`? |
|---|---|---|
| `claim(..., span=...)` | `engine.py:172-173`; span stored in provenance (`:188`) | **No** |
| `acquire(..., span=...)` | `engine.py:193-194`; span stored in provenance (`:214`) | **No** |
| `descriptor: dict` — unconstrained | `engine.py:194` | yes, but only with the synthesised token |

### 1.3 The five questions

**Q1 — Is an external event reference expected to be supplied to MEMORY-09?**

**Established: a *span* may be, and a *resolvable pointer* must not be.** v9 §7 step 2 has the conversation
layer establish the event and says it is "recorded as a claim event **with spans**"; §8 fixes the recipient's
origin reference as "a **descriptor** (speaker, record type, turn), **not a resolvable identifier**", and adds
that "the resolvable link exists only on the source side". Rule 12 states the same prohibition.

So the contract expects MEMORY-09 to receive *facts about the event* (who, what, which recipients, and a span
when one exists) and forbids it from receiving a pointer that would resolve back to the event or the source.
There is no clause anywhere requiring an external event **id**.

**Q2 — Does the current experiment have an interface capable of consuming such a reference?**

**Established: partly.** Three slots exist (`claim(span=)`, `acquire(span=)`, and the free-form `descriptor`),
so a span could be carried today. The gap is that `disclose` takes no span and synthesises its own event token
instead of forwarding one, and it returns accounts without an event record. The *interface* is adequate; the
*operation* does not use it.

**Q3 — Is the absence a missing interface, a missing external dependency, or a genuine missing record?**

**Established: a genuine missing record — and it is MEMORY-09's own.** The three candidates separate cleanly:

| Candidate | Test | Verdict |
|---|---|---|
| Missing external dependency | Would the event layer have to supply something it does not already supply? `disclose` already receives discloser, recipients, source, turn. | **No** — the inputs are present |
| Missing interface | Would a slot have to be invented? `span` and `descriptor` exist; only the signature of `disclose` omits them. | **Minor, not the defect** |
| Genuine missing record | Does rule 46 require a record that does not exist? "**The system records the event and its attribution**"; §7 step 2 "recorded as a claim event with spans"; J3 "Claim event with content and recipients". | **Yes** |

The boundary the owner asked me to respect cuts the other way from my first reading. v9 §1 assigns the *facts*
to the conversation layer; rule 46 assigns the *recording* to the memory system. MEMORY-09 is not being asked
to decide whether a disclosure occurred — the caller asserts that by invoking `disclose` with named recipients.
It is being asked to write down what it was told. It currently writes down the consequences and not the event.

**Not manufacturing event semantics** would mean not inventing exposure, not inferring disclosure from
proximity, not deciding who was present. Creating one record that says "at turn T, X disclosed C to [A, B]",
from arguments the caller supplied, is none of those; it is exactly what rule 46 names.

**Q4 — Would the proposed `_events` registry duplicate or violate another layer's responsibility?**

**Established: it would duplicate MEMORY-09's own representation, and is therefore rejected.** §7 step 2 does
not say "recorded in an events table"; it says "recorded as a **claim event** with spans". The canonical carrier
already exists: the `claim` class, with contributor attribution, declared kind, span, and provenance. A second
registry would give one event two representations — a claim-shaped one and a registry-shaped one — with two
places where "what was said" lives, two histories to keep consistent (rule 20), and two attribution paths
(rule 23). That is the duplication defect the architecture spends REV 4–9 removing elsewhere.

My §4.1 option (a) was therefore wrong in its shape, though right in its diagnosis. The correction is recorded
in §1.5.

**Q5 — Can rule 46 and J3 be satisfied without a new MEMORY-09 event registry?**

**Established: yes — one `claim` per disclosure, created by `disclose`.** No new class, no new field, no new
registry, no change to the record model.

### 1.4 Established / inferred / unresolved

| Kind | Statement |
|---|---|
| **Established** | Rule 46 requires MEMORY-09 to record the event; it does not require a new structure to do so |
| **Established** | The canonical carrier is a `claim` event with spans (v9 §7 step 2), and the `claim` class already carries contributor, kind, span and provenance |
| **Established** | A resolvable event pointer is forbidden into recipient provenance (§8, rule 12); an external span is permitted |
| **Established** | The current code creates no event record; the P classification is correct |
| **Inferred** | Two disclosures in one turn are indistinguishable by the current token, because `(kind, turn)` is not unique within a turn |
| **Unresolved** | Whether the event claim's `contribution_kind` should be supplied by the event layer or default to `Unspecified claim` (rule 26 forbids the system choosing one) — an owner/implementation choice |
| **Unresolved** | Whether the event claim is the discloser's private record (owner = discloser, as `claim` requires) or something else; the `claim` class has no other shape today |

### 1.5 Smallest permissible resolution

**Replace the synthesised token with one recorded claim; forward, do not invent.**

| Step | Change | Size |
|---|---|---|
| 1 | `disclose` gains optional `span` and optional `kind` parameters | two keyword parameters, defaulted |
| 2 | `disclose` creates **one** `claim` (contributor = discloser, kind = the caller's declaration or `Unspecified claim`, text = the disclosed content, span forwarded, provenance records recipients, source class and turn) — the event | one existing operation call |
| 3 | Each recipient account's descriptor carries the event claim's **record id? No** — §8 forbids a resolvable reference into recipient provenance | nothing; the descriptor keeps `speaker`, `record_class`, `event`, `turn` |
| 4 | Return value reports the event claim alongside the accounts | one key in the return |

Constraints that keep it contract-consistent:

| Constraint | Reason |
|---|---|
| The event claim is a `claim`, never truth | rule 48, `truth_apt` |
| Its `kind` comes from the caller or is `Unspecified claim` | rule 26 / L9: the system must not supply a kind |
| Recipients are recorded on the **event**, not as a resolvable link on the accounts | §8: descriptor is speaker/class/turn only |
| No `_events` registry | §7 step 2 names the claim as the record |
| No change to rule 11's exposure logic | exposure still follows from the event alone |

**Behavioural consequences to expect, none of them violations:** `perspective(discloser)` will now return the
event claim (the discloser's own utterance, rule 5); `history()` grows by one record per disclosure; and the
disclosure becomes retrievable as an event rather than only as N consequences.

**Not implemented.** This supplement is a reconciliation, and the instruction forbids implementation. The change
also needs a delta row, because no REV 3–13 row authorises it — which is precisely why the earlier report was
wrong to call the row complete rather than reporting it.

### 1.6 Classification

| Row | Classification | Change |
|---|---|---|
| Rule 46 | **P** | none — remains partial until the event claim exists |
| J3 | **P** | none |
| Rule 10 | I | none — "shared information only when an actual disclosure event establishes it" is satisfied by the operation, and is unaffected by whether the event is recorded |
| J7, H4, J8 | I | none — recipient accounts and their provenance are correct |

---

## Part 2 — K9: timestamp, determinism, and one imprecise claim of mine

### 2.1 The claim under review, and its correction

The consolidated review says:

> v9 §6 step 2 requires the promotion to record "…, **timestamp with offset**, …". `_promotions` records
> `turn`. … The engine has no clock, by design: REV 10 S-2 forbids deriving any response field from a clock …

**That attribution is imprecise, and I withdraw the wording.** Three separate things were bundled:

| Statement | Precise position |
|---|---|
| A clock read during retrieval is forbidden | **True**, but the prohibition is phrased in REV 10 §3.1's definitional sentence ("from a clock"), not in S-2's rule text |
| S-2 forbids it | **Partly** — S-2 forbids any field "computed from a read made after `S`"; a clock read is the case §3.1 enumerates |
| A **persisted** write-time timestamp is forbidden | **False.** Nothing in S-1…S-8 addresses persistence at all |

The distinction matters because it changes the conclusion from "the contract forbids timestamps" to "the
contract constrains *where* a time value may come from relative to a response".

### 2.2 Exact clauses

**v9 §6, promotion lifecycle, step 2:**

> **Act** — an authorised actor performs an explicit promotion, recorded with actor, authority reference,
> **timestamp with offset**, source, target domain, declared applicability, basis, optional rationale.

**REV 10 §3.1, the definitional sentence:**

> `S` is a value, not a view. After it is taken, no field of the response may be derived from a later read of
> the record store, the group store or the registry, **from a clock**, or from a second pass over
> `_eligible_truth`. Everything the response contains — entry `state`s, record counts, `relation`,
> `open_groups`, `establishing_groups`, `integrity` — is a projection of `S`.

**REV 10 §3.2, S-1…S-6 (verbatim):**

> | S-1 | **One take.** `S` is taken once per call, before the first read that contributes to the response. |
> | S-2 | **No re-read.** Every response field is a projection of `S`. No field is computed from a post-`S` read. |
> | S-3 | **One map, two projections.** … |
> | S-4 | **Staleness permitted, mixing not.** A concurrent write may make `S` stale; the response must still be a function of exactly one `S`. No read-your-writes guarantee is offered within a call, and none is required: mutating calls (`assert_incompatible`, `supersede`, `retract`) are separate calls whose effects appear in later snapshots. |
> | S-5 | **Refusals take no snapshot.** … |
> | S-6 | **Assembled means assembled from `S`.** … |

Note what S-1…S-6 quantify over: **response fields**, **`S`**, and **reads**. None of them mentions a write, a
record's persisted provenance, or what a mutating call may store.

**REV 11 S-7 and S-8:**

> | **S-7** | **Coherent capture.** Every field of `S` must have co-existed at one committed state. `S` is admissible only if there is a single committed version `v` with `accepted = accepted@v`, `returned_identities = identities@v`, `resolution = resolution@v` and `groups = groups@v`. Atomic snapshot, transactionally consistent read, and version-check-and-retry all satisfy this; the contract prescribes none of them. |
> | **S-8** | **Frozen identity set.** `returned_identities` is captured inside the same coherent read as the rest of `S` … none may consult the registry afterwards. |

S-7 constrains the *capture*, not the *write*. A value written at promotion time and read inside the capture
co-exists with everything else by construction, because it is part of the committed state.

**REV 11 §3.4, what capture coherence is not:**

> | Settled by S-7 | Not settled |
> |---|---|
> | `S` describes a state that once existed | Whether it is the newest — liveness remains the caller's concern |
> | No response can present a blended state | … |
> | A write landing mid-capture cannot produce a phantom `unresolved_relation` | Replay. `S` is still not recorded or persisted; F7 stays **R** |

**REV 12 and REV 13 amendments.** S-9 (capture refusals publish nothing derived), S-10 (ordering), S-11
(content independence), S-12 (metadata boundary) — quoted in the verbatim-sources file §13–§14 — introduce no
wording about time, clocks or persistence. They extend the *refusal* and *metadata* boundaries only. **No
amendment in REV 12 or REV 13 restricts a persisted timestamp.**

### 2.3 The four-way distinction

| Case | What it is | Contract verdict | Why |
|---|---|---|---|
| **(a) Captured at write, persisted** | A promotion stores `at` with offset as part of its own record | **Permitted** | It is written before publication, inside the mutating call, and is part of the committed state that S-7 captures. No clause forbids it. K9 and v9 §6 step 2 require its content |
| **(b) Obtained during retrieval** | A query calls `now()` to stamp or age the answer | **Forbidden** | REV 10 §3.1: no response field derived "from a clock"; S-2: no field computed from a read made after `S` |
| **(c) Regenerated during replay** | A replay re-stamps records rather than reading stored bytes | **Forbidden by F7's premise, not by S-1…S-8** | F7 asks that replaying the same ordered records yield identical answers; re-stamping makes the *records* differ. A replay must read the persisted value |
| **(d) Turn counter as ordering** | `created_turn` orders records and breaks ties | **Permitted, and not a timestamp** | It is deterministic data. It answers "in what order", never "when" |

Case (a) and case (b) are different events in different phases. My earlier sentence collapsed them, which is
what made K9 look like a contract conflict rather than a missing field.

### 2.4 Current implementation evidence

| Where "time" appears | Source | Semantics |
|---|---|---|
| `Record.created_turn: int` | `records.py:74`, exposed in `as_dict` `:191` | ordinal turn |
| Promotion operation record | `_promotions[promo_id]` in `promote`: `promotion_id, actor, domain, target, turn, record_id, proposition_id, rationale, basis, superseded` | ordinal turn; **no timestamp, no offset** |
| Proposition identity | `_propositions[pid]`: `created_by`, `created_turn` | ordinal turn |
| Ordering | `_take` sorts identities by `proposition_id`; `_assemble` sorts records by `record_id`; `history` sorts by `(created_turn, record_id)` | ids and turn |
| `snapshot()` | `version`, records, propositions, groups, promotions, declarations | no timestamp |

So the implementation records **case (d)** and neither (a) nor (b). K9's "time" is satisfied only in the
ordinal sense, which is why it is **P** and not **I**.

### 2.5 Can a persisted creation timestamp with offset satisfy K9 while preserving the contract?

**Yes, under five conditions. Established, with one inference named.**

| # | Condition | Source |
|---|---|---|
| C1 | The clock is read **once, inside the mutating call**, before publication and before the version bump, so the value is part of the committed state | S-7 (admissibility is about capture, not about how a committed value was produced) |
| C2 | The value is **persisted**, and any response that shows it reads it from `S` like any other field | S-1, S-2 |
| C3 | **No ordering, state or comparison** may depend on it. Ordering stays `(created_turn, record_id)` | REV 11 §3.4 (liveness is not settled by S); F7's determinism claim is ordering by id |
| C4 | **No retrieval-time clock read**, and no derived field such as age or recency | §3.1 "from a clock"; S-2 |
| C5 | The offset is recorded as written, so the value is unambiguous across zones | v9 §6 step 2 "timestamp with offset"; K9 "time" |

**Inference, stated as such:** because K9 concerns a recorded *operation* rather than a *response*, the smallest
compliant placement is inside the promotion operation record — not in the three view envelopes. If it stays
there, **the snapshot contract is untouched**: no new `S` field, no new capture step, no A-step change, no
envelope change.

### 2.6 What adding it would and would not disturb

| Effect | Assessment |
|---|---|
| S-1…S-8 | Unchanged, if the value is not projected into the three views |
| F7 "replay same ordered records → identical answers" | Preserved — the timestamp becomes part of the records being replayed |
| Two fresh runs of the same scenario | **Would differ** in the timestamp. F7 does not require run-to-run identity of values, and the determinism claim in the code is about ordering; but any future test that compares two fresh runs byte-for-byte would now fail. This is a real consequence to state, not a defect |
| Tally | None; K9's classification is about content, not about counts |
| Authorisation | A **new persisted field**, therefore outside the current authorisation ("do not add … fields") — hence reported, not implemented |

### 2.7 Classification and remaining owner decision

| Row | Classification | Change |
|---|---|---|
| K9 | **P** | none — "time" is ordinal only |
| v9 rules 20, 21, 42, 43, 44 | I | none — none of them mentions a timestamp |
| F7 | **NA** | none — no replay facility |

**Owner decision required:** authorise one persisted field (promotion record `at`, ISO-8601 with offset) under
conditions C1–C5, or ratify that a turn counter is an acceptable reading of "time" for the experiment. The
second option is the smaller one and requires no code change; it does change what K9 means, so it is yours.

---

## 3. What this supplement does not do

| Not done | Note |
|---|---|
| No implementation | No file in `eval/mem1/` changed; neither resolution is applied |
| No test, probe, compile, lint or runtime check | Nothing executed |
| No rewrite of any review document | The consolidated review, the report, the supplement and the clarification stand byte-identical |
| No contract change | v9, the map, REV 1–13, r32, plan v3, Amendment 1 and the acceptance instrument untouched |
| No new architecture | The event resolution uses the existing `claim` class; the time resolution names one field, subject to your decision |

---

## 4. Verification for this supplement

| Item | Value |
|---|---|
| Implementation | `authority.py f3445823…`, `engine.py 706d5f6c…`, `records.py f211e6ce…`, `errors.py 48125d7e…`, `retrieval.py cd2d380b…`, `__init__.py 06ccf8f6…` — unchanged by this supplement |
| Contract and sources | REV 13 `5b745628…`; v9 `698254ca…`; map `b4ba33cc…`; ledger `da3e7239…`; acceptance `9a25ded0…` |
| Evidence read | `disclose` (engine.py:552-591), `claim`/`acquire` signatures (:172-196), `acquire` provenance (:214), `_promotions` construction in `promote`, `Record.created_turn` (records.py:74, :191) |
| Executed | **nothing** |

**Sequence:** this supplement → owner decision on §1.5 (implement the event claim) and §2.7 (timestamp or turn)
→ implementation authorisation, if any → conformance re-review → separate testing authorization. This document
requests none of these.
