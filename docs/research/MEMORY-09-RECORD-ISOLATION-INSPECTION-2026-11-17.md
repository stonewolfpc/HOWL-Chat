# MEMORY-09 — Source inspection: record independence and mutable aliasing

**Source inspection only. No file changed, nothing executed, no import, no probe. Every finding below is derived by
reading the construction, storage, serialisation and return paths; none is an observed runtime result.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Code inspected | `records.py 778684e2a63d82d9…`, `engine.py fcbc7eb256287ca1…`, `authority.py f3445823…`, `retrieval.py cd2d380b…` |
| Code changed | **None** |
| Executed | **Nothing** |

---

## 0. Verdicts

| Question | Verdict |
|---|---|
| **Independent records** (each record its own identity, never a view/alias/projection) | **CONFIRMED** |
| **Snapshot isolation** (`snapshot()` / `as_dict()` cannot reach stored state) | **NOT CONFIRMED** — four demonstrated paths |
| **Cross-record mutable aliasing** (one record's mutation reaching another) | **CONFIRMED** — one demonstrated path, narrow |
| **Disclosure metadata isolation** | **CONFIRMED** through `snapshot()`/`as_dict()`; **NOT CONFIRMED** through the live claim record |
| **Other fields affected** | `provenance`, `corrected_by`, `promotions.basis`, `declarations.audience`/`covered`, `groups.members`/`targets`/`status` |
| **Violation of a ratified requirement by engine behaviour** | **None found** — every path below requires a caller to mutate a returned object |
| **Defect** | Four leaked mutable handles, of which two can change *authorisation and conflict state*, not merely metadata |

---

## 1. Trace: construction → storage → serialisation → return

### 1.1 What is copied, and what is only referenced

| Site | Statement | Copy depth |
|---|---|---|
| `acquire` (engine.py:196) | `desc = dict(descriptor or {})` | **shallow** — nested `event` object is shared with the caller and with every other account built from the same descriptor |
| `acquire` (engine.py:214) | `provenance={"descriptor": desc, "span": span, "turn": turn}` | stores `desc` by reference |
| `claim` (engine.py:186-188) | `provenance={...}` fresh literal | isolated |
| `disclose` (engine.py:606-608) | `event = {...}`; `descriptor = {..., "event": event}` built **once, outside the loop** | the same `event` object reaches every recipient account |
| `declare` (engine.py:681-682) | `"audience": list(audience)`, `"covered": []` | **copied** — caller's list not shared |
| `declare` (engine.py:684-688) | `descriptor={..., "event": {"kind": "declaration", ...}}` built **inside** the comprehension | a fresh descriptor and event per account — the contrast with `disclose` |
| `promote` (engine.py:418) | `"basis": [vars(e) for e in entries]` | **none** — `vars()` returns the instance `__dict__` itself |
| `assert_incompatible` (engine.py:790-796) | `group = {... "targets": list(targets) ...}` | copies its inputs |
| `Record.as_dict` (records.py:222) | `"provenance": dict(self.provenance)` | **shallow** — nested `descriptor`/`event` are the stored objects |
| `Record.as_dict` (records.py:230) | `"corrected_by": list(self.corrected_by)` | **shallow** — list copied, the dicts inside are the stored ones |
| `Record.as_dict` (records.py:225-227) | `"disclosure": {"recipients": list(...), "source_class": ...}` | **fresh dict + fresh list** — genuinely isolated |
| `Record.as_dict` (records.py:221) | `"domain": list(self.domain)` | isolated (and the tuple is immutable anyway) |
| `snapshot` (engine.py:1113-1118) | `{k: dict(v) for k, v in self._propositions/…}` | **shallow** for promotions, declarations, groups |

### 1.2 What the return paths hand out

| Surface | Returns | Live or copied |
|---|---|---|
| `_record_view` → `authoritative`, `world`, `belief_facts` | `{record_id, text, created_turn, promotion_id}` | **copied scalars** — isolated |
| `retrieval.perspective` | `Record` objects | **live** |
| `Engine.history`, `all_records` | `Record` objects | **live** |
| `claim`, `acquire`, `disclose` | `Record` objects | **live** |
| `declare` | `decl` — the stored dict | **live** |
| `assert_incompatible` | `group` — the stored dict | **live** |
| `AuthorityRegistry.grants_for` | `Grant` objects | **live** (`Grant` is a plain `@dataclass`, not frozen) |
| `narration_report`, `_refusal` | fresh dicts of scalars | isolated |

---

## 2. Demonstrated aliasing — each with its cause

### 2.1 Cross-record: the shared event token (Answer 2 — **confirmed**)

`disclose` builds `event` once and one `descriptor` once, then calls `acquire` per recipient. `acquire` shallow-
copies the descriptor, so each account gets its own `provenance["descriptor"]` dict — but `descriptor["event"]` is
the **same object** in every one. Therefore:

```
accounts[0].provenance["descriptor"]["event"]["turn"] = 999
  -> visible in accounts[1].provenance["descriptor"]["event"]
```

**Scope of the damage:** that object is read in exactly one way — truthiness, at `_validate_evidence` lines 237,
251, 258 (`if not descriptor.get("event")`). Setting it to `{}` would make a *later* evidence check fail for a
record already written; no stored record's meaning changes otherwise. No code compares tokens (no `descriptor ==`),
so no rule outcome flips. **Real aliasing, low semantic blast radius.**

### 2.2 Snapshot: `provenance` and `corrected_by` (Answer 1 — **not confirmed**)

`as_dict` shallow-copies both. Therefore:

```
snap["records"][i]["provenance"]["descriptor"]["event"]["turn"] = 999   # reaches the stored account
snap["records"][i]["corrected_by"][0]["prior_text"] = "rewritten"       # reaches the stored history entry
```

The second is the more serious: `corrected_by` is the retained prior text that rule 20 / rule 44 protect.

### 2.3 Snapshot: promotions, declarations, groups (Answer 1 / 5)

`snapshot` applies `dict(v)` — one level only. So:

| Key | Shared through the snapshot | Consequence |
|---|---|---|
| `promotions[pid]["basis"]` | the list and every dict inside (`vars(e)`) | a stored basis entry can be rewritten; rule 28 evidence changes retroactively |
| `declarations[did]["audience"]`, `["covered"]` | both lists | coverage can be rewritten; `audience` is written by `declare` only, `covered` is unreferenced elsewhere |
| `declarations[did]["withdrawn_turn"]` | — | **scalar**, safe |
| `groups[gid]["members"]`, `["targets"]`, `["propositions"]` | the lists | a conflict group's membership can be edited |

### 2.4 Live handles that can change state, not just metadata (Answer 5)

| Handle | Mutation | Effect |
|---|---|---|
| `declare(...)` → stored `decl` | `decl["withdrawn_turn"] = None` | **re-activates a withdrawn declaration**: `_validate_evidence`'s `narrator_grant` branch refuses on `withdrawn_turn is not None` (B6) |
| `assert_incompatible(...)` → stored `group` | `group["status"] = "resolved"` | removes an open conflict from the duplicate-group check at engine.py:786 and from any `status == "open"` read — J13/J15 territory |
| `grants_for(...)` → live `Grant` | reassign `operations`/`revoked_*` | alters the authority gate's input |
| any returned `Record` | `rec.text = …`, `rec.status = …`, `rec.corrected_by.append(…)` | stored record mutated; `__setattr__` guards **only** `proposition_id` |

These are the highest-severity findings, because they are reachable through **ordinary return values** of rated
operations — not through a hostile cast of `all_records()`.

### 2.5 Disclosure metadata (Answer 3)

| Path | Isolated? | Why |
|---|---|---|
| `snapshot()["records"][i]["disclosure"]` | **Yes** | `as_dict` builds a fresh dict and `list(recipients)`; `recipients` is itself a tuple |
| `claim.disclosure["source_class"] = "world_truth"` | **No** | `disclosure` is a plain dict on a mutable record |
| `claim.disclosure = {...}` | **No** | no `__setattr__` guard; only `proposition_id` is protected |
| `claim.disclosure["recipients"] += ("x",)` | **No** | tuple is immutable, but the key can be reassigned wholesale |

So the *serialised* view is clean and the *live object* is not. Since the live claim is what `disclose` returns,
this matters for Answer 4.

### 2.6 Records themselves (Answer 4)

**Confirmed: records are independently owned as identities.** Each `Record` is a distinct object in
`_records`; nothing returns a view, alias or projection of another record; each account has its own
`provenance` mapping and its own `record_id`. O-1's separateness and §8's "never a view, alias or projection" are
satisfied at that granularity.

**Not confirmed: records are independently owned as *states*.** A returned `Record` is the stored object, and its
nested values can be shared (§2.1). "Independent identity" holds; "no shared mutable value anywhere in the object
graph" does not.

---

## 3. What is genuinely isolated (stated for fairness)

- `_record_view` and everything built on it — `authoritative`, `world`, `belief_facts` — fresh dicts of scalars.
- `as_dict()["disclosure"]`, `as_dict()["domain"]`.
- `declare`'s `list(audience)` at construction; `assert_incompatible`'s `list(targets)` and `list(first.domain)`.
- `Grant.operations` is a `frozenset` (the attribute can be reassigned on a live Grant, but the set itself cannot
  be mutated).
- `BasisEntry` is `@dataclass(frozen=True)` — the instance cannot be reassigned. It does **not** protect
  `entry.__dict__`, which is what `vars(e)` stores; direct dict writes still reach it.

---

## 4. Does any of this violate a ratified requirement?

| Requirement | Verdict |
|---|---|
| **O-1 / §8** — each recipient record has its own identity and provenance, never a view, alias or projection | **Satisfied.** Distinct objects, distinct provenance mappings. §2.1 shares a value *inside* provenance, not the provenance itself or the record |
| **Rule 20 / 44** — retained history cannot be erased | **Satisfied by every engine operation.** Reachable only by a caller writing into `corrected_by` through a snapshot; the engine never rewrites prior text |
| **J13 / J15** — conflict state cannot be silently resolved | **Satisfied by every engine operation**; reachable through the live group dict returned by `assert_incompatible` |
| **B6** — a withdrawn declaration confers nothing | **Satisfied by the operation**; reachable through the live `decl` dict returned by `declare` |
| **Rule 12 / §8** — recipients cannot resolve the source | **Satisfied**; no pointer exists to mutate |
| **F7** — replay determinism | **NA** — no replay facility. Caller mutation would defeat any future replay, but nothing claims replay today |

**Conclusion: no semantic violation exists in the engine's own behaviour.** Every path requires a caller to mutate
an object the engine handed back. The defect is a **missing isolation boundary**, and two of its leaks happen to
pass through authorisation and conflict state rather than through inert metadata.

---

## 5. Minimum correction

Ordered by severity. None of these changes memory semantics, adds a field, or alters any operation's meaning.

| # | Change | Site | Size |
|---|---|---|---|
| 1 | Return a copied declaration: `dict(decl, audience=list(decl["audience"]), covered=list(decl["covered"]))` | `declare` return | 1 line |
| 2 | Return a copied group, copying its three lists | `assert_incompatible` return | 1-2 lines |
| 3 | Give each account its own event object — build the descriptor inside the loop with `"event": dict(event)`, as `declare` already does | `disclose` | 1 line moved |
| 4 | `[dict(vars(e)) for e in entries]` | `promote` | 1 line |
| 5 | `copy.deepcopy` for `provenance` and `corrected_by` in `as_dict` | `records.py` | 1 import, 2 expressions |
| 6 | `copy.deepcopy` for promotions, declarations and groups in `snapshot` | `engine.py` | 1 line |
| 7 | *(separate decision)* live `Record` returns from `claim`/`acquire`/`disclose`/`history`/`all_records`/`perspective` | API-wide | NOT part of the minimum |

**Why 1–4 are the minimum that matters:** they close the handles that can change *state*, not just metadata, and
each is behaviour-preserving — a caller who only reads sees exactly what it saw before.

**Why 7 is excluded:** making Records immutable or returning copies would change every operation's return type and
the way callers use ids (`record.record_id`). No ratified requirement demands an in-process tamper boundary, and
introducing one is an architectural decision, not a corrective patch. Recommend recording it as an accepted
property, or deciding separately.

**Impact of applying 1–6:** no read result changes for a non-mutating caller; `deepcopy` costs scale with record
count in `snapshot` (already the whole store) and are negligible per call; `vars()`→`dict(vars())` changes nothing
observable except that the caller's entry no longer doubles as the stored one.

---

## 6. Status, kept separate

| Item | Status |
|---|---|
| Implementation of B2/D-E1 | Present in source; **unchanged by this inspection** |
| This inspection | Source reading only; findings are claims about code, not observations of behaviour |
| Executed | **Nothing** — no import, no test, no probe, no compile |
| Classifications | Unchanged: rule 46 **P**, J3 **P**, K9 **P**, B3 **NA**, F7 **NA**; 56 I · 1 P and 113 I · 3 P · 2 NA |
| Testing | Still unauthorised |
