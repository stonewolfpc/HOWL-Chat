# MEMORY-09 — B2 implementation review R1: refusals, event token, atomicity

**Scope: source inspection only. No code changed, nothing executed, no probe written. These are answers to three
review points, not a conformance review and not a test result.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Code reviewed | `engine.py fcbc7eb256287ca1…`, `records.py 778684e2a63d82d9…`, `authority.py f3445823…`, `errors.py 48125d7e…`, `retrieval.py cd2d380b…`, `__init__.py 06ccf8f6…` |
| Code changed by this review | **None** |
| Executed | **Nothing** — no test, probe, compile, lint or import |
| Verdict | Architecture unchanged; two pure safeguards identified; event token has no functional effect; atomicity holds by exhaustion of raise sites, with two named future-edit hazards |

---

## 1. Which refusals are ratified, and which are engineering safeguards

The implementation raises twelve refusal conditions. They are **not** one class of decision, and the distinction
matters because refusal *is* memory policy when it constrains what can be recorded about a character or domain.

### 1.1 Classification

| # | Refusal | Anchor | Anchor type | New site? | Restricts memory content? |
|---|---|---|---|---|---|
| 1 | `content` blank or non-string | J4 ("accounts from intent, proximity or **silence**"); §5.3 "content explicit" | Ratified + owner decision | **New** | No — restricts the *call*, not what may be stored about anyone |
| 2 | `recipients` empty | Rule 11 ("only characters actually exposed acquire"); J5; §5.3 "absent or empty refused" | Ratified + owner decision | Preserved | No |
| 3 | Recipient not a non-empty string | **None** | **Safeguard** | New | No |
| 4 | Recipient named twice | **None** — §5.3's "no reinterpretation" supports refusing over deduping, but no clause requires either | **Safeguard** | New | No |
| 5 | Discloser among recipients | C-3 `told` evidence ("a character cannot be told by itself"); §7's exposure semantics | Ratified (pre-existing behaviour) | Preserved | No |
| 6 | Source record absent | Rule 33 (a citation must resolve) | Ratified | Preserved | No |
| 7 | Source in another chat | Rule 31 (containment) | Ratified | Preserved | No |
| 8 | Discloser does not own the source | Rule 46 ("what it has **acquired or remembers** through its own conversational output") | Ratified | Preserved | No |
| 9 | Source is domain truth | Rule 46 / v9 §7 step 1 (`claim || account || belief`) | Ratified | Preserved | No |
| 10 | Unknown contribution kind | Rule 26 / L9 (kind is declared, never inferred) — the same check `claim()` already makes | Ratified | New input, old rule | No |
| 11 | A recipient cannot be `told` (pre-flight) | C-3 (`told` requires speaker and event) | Ratified | New site, same refusal as `acquire` already made | No |
| 12 | Malformed `disclosure` metadata, including direct construction | J3 + rule 46 | Ratified | New | No |

**Count: 10 ratified or owner-authorised, 2 pure safeguards (3, 4).**

### 1.2 The specific point raised — discloser as recipient

Not new policy. The pre-B2 `disclose` refused it explicitly with the same error class and message, and C-3 would
refuse it independently through the `told` evidence check (`speaker == character`). The B2 rewrite preserved the
check and moved it earlier; it neither added nor removed a permission.

### 1.3 What the two safeguards actually do

| Safeguard | Effect | Why it is not policy |
|---|---|---|
| Non-string recipient | Rejects `None`, `""`, `1` | Types only; no character or domain gains or loses a capability |
| Duplicate recipient | Rejects `["a", "a"]` | Constrains call shape. The alternative — silently deduping — would construct a *second* set, which §5.3 forbids. Refusing keeps the supplied set verbatim |

Neither can make a character know less or more; both can only make a malformed call fail loudly instead of
silently producing a differently-shaped event. **Neither was required by the ratified text, and both are recorded
here so that discretion is not mistaken for policy.** If the owner prefers leniency on either, the change is one
condition each and touches no memory semantics.

### 1.4 One genuine behaviour change worth naming

Refusals 1–12 now occur **before** `_mutation`. Previously the checks ran inside it, so a refused call advanced
`version`. Now it does not. This is a property of the refusal path, not of the architecture, and it makes refused
calls invisible to the version counter — which is closer to S-5 ("refusals take no snapshot") than the old
behaviour was. Recorded so it is not discovered later as a surprise.

---

## 2. The recipient event token: what it can and cannot affect

The token is `{"kind": "disclosure", "turn": turn}`, identical for two disclosures in the same turn.

### 2.1 What the code does with it

| Consumer | Behaviour |
|---|---|
| `_validate_evidence`, `witnessed` (l.237) | truthiness only |
| `_validate_evidence`, `told` (l.251) | truthiness only |
| `_validate_evidence`, `read` (l.258) | truthiness only |
| `_validate_evidence`, `verified` (l.277) | truthiness of a different key |
| Anywhere else | **nothing** — `grep 'descriptor =='` returns no equality test, and no path compares one record's descriptor with another's |

So the token is *evidence that an event exists* (C-3), never an identity.

### 2.2 Consequence per surface

| Surface | Effect of the collision |
|---|---|
| **Acquisition provenance** | Two same-turn accounts by one recipient carry identical descriptor content. They remain distinct records with distinct ids and provenance dicts. No rule requires distinguishable provenance between two accounts, and no code path treats identical descriptors as duplicates (`grep` for dedupe logic: none) |
| **Replay** | None. Ordering is by `record_id` (`_assemble`, `history` tie-break), and ids are assigned by call order, so the token is never an ordering or comparison input. F7 remains **NA**: no replay facility exists, so this is a statement about the code, not about replay behaviour |
| **Retrieval** | None. `perspective` filters by owner and status and scores on `r.text`; `history` filters on owner/domain/target and sorts on `(created_turn, record_id)`. The descriptor is never read |
| **Snapshot** | The token is visible, so a reader sees two accounts whose descriptors agree. That reader also cannot tell *which* claim produced either account — which is rule 12 / §8 working as designed, not a defect |

### 2.3 Residual risk, stated rather than fixed

If a future requirement ever needs per-event identity **on the recipient side**, it cannot be satisfied by storing
the claim id — that is the resolvable link rule 12 forbids. It would need a non-resolvable opaque token supplied by
the layer, which is a new decision, not a defect in this one. Nothing today needs it.

### 2.4 An aliasing note, pre-existing and not token-specific

`acquire` stores `dict(descriptor)` — a shallow copy — so every account of one disclosure shares the same nested
`event` object, and `as_dict()` shallow-copies `provenance` into snapshots. Nothing mutates it today. A caller who
mutated a returned snapshot could therefore alter stored provenance. This predates B2 and affects every `acquire`
call with a nested descriptor; it is unchanged by this implementation.

---

## 3. Atomicity: can anything raise after the first write?

### 3.1 The mechanism has no rollback

`_mutation` is:

```python
self._depth += 1
try:
    yield
finally:
    self._depth -= 1
    if self._depth == 0:
        self._version += 1
```

There is no journal and no undo. Records added before an exception **remain**, and the version still advances —
because the bump is in `finally`. A partial write would therefore be *published* as a normal version step and
would not be detectable by version comparison. This is why the ordering guarantee has to be true by construction
rather than by hope.

### 3.2 Enumeration of everything after the first write

Inside `disclose`'s `with self._mutation():` block:

| Step | Can it raise? | Why not |
|---|---|---|
| `self._id("claim")` | No | integer increment and string format |
| `Record(...)` for the claim | No | every branch in `__post_init__` is satisfied by validated inputs: owner is a non-empty character; chat is the engine's; class is `claim`; `proposition_id` is `None`; `contribution_kind` is a member of `CONTRIBUTION_KINDS`; `disclosure` passed the new record-level check in the same construction |
| `self._add(...)` | No | dict assignment |
| `self.acquire(...)` per recipient | **No raise statements exist in `acquire`** — `grep -n raise` over its body returns nothing | its only refusal is `_validate_evidence`, already run per recipient in pre-flight with the same character, `via`, `origin_class`, descriptor and turn |
| `Record(...)` inside `acquire` | No | `origin_class="acquired"`, `via="told"` (a member of `VIAS`, not `none`), owner present, chat present, `proposition_id=None`, `disclosure=None` |
| `tuple(...)`, `return` | No | pure |

**Conclusion: under the current source, no exception is reachable after the first write.** The guarantee is real
but *conditional* — it is a property of the code as written, not of a mechanism that enforces it.

### 3.3 The two conditions the guarantee depends on

| Condition | Breaks if |
|---|---|
| C-a | `acquire` gains a refusal *outside* `_validate_evidence` (a new check in its body) | pre-flight would not cover it, and it would fire mid-loop |
| C-b | `Record.__post_init__` gains a validation applicable to `claim`/`account` shapes | construction inside the write block could raise after the claim landed |

Both are future-edit hazards, not current defects. The source carries a comment marking the pre-flight as
load-bearing; that comment is the only thing protecting C-a and C-b, and a reviewer who removes it removes the
warning, not the hazard.

### 3.4 A related pre-existing property

`acquire` used to be called by `disclose` **inside** its own `_mutation` and still is; the depth counter makes the
version bump once. That part is unchanged and correct.

---

## 4. Status after this review

| Item | Status |
|---|---|
| D-E1 representation | Implemented; **this review finds no contradiction with §5.3** |
| B2 interface | Implemented; content supplied, recipients supplied once, source id not stored |
| Refusals | 10 ratified/authorised, 2 safeguards — now named individually |
| Event token | No functional effect on provenance, ordering, retrieval or version; F7 still **NA** |
| Atomicity | Holds by exhaustion of raise sites; conditional on C-a and C-b |
| Source review (full, all 118 + 57) | **Not done** — this memo reviews three points, not the row set |
| Tests and execution | **Still unauthorised; nothing has run** |
| Classifications | Unchanged: rule 46 **P**, J3 **P**, K9 **P**, B3 **NA**, F7 **NA**; 56 I · 1 P and 113 I · 3 P · 2 NA |

**No classification moves because of this memo.** It reports that the engineering decisions did not quietly
change the architecture; it does not report that anything works.

---

## 5. What was inspected

`grep` for every `descriptor.get(...)` consumer, every `event` occurrence, and any `descriptor ==` / dedupe
comparison (none); `_mutation` and `_add` bodies; the full `acquire` body for `raise` (none); the `disclose` write
block line by line; `retrieval.perspective` and `Engine.history` for descriptor use (none); `Record.__post_init__`
branch by branch against the pre-validated inputs.

**Nothing executed.** Every statement above is a claim about source, not an observation of behaviour.
