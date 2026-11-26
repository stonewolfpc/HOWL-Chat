# MEMORY-09 — Isolation hardening report: S1–S6

**Implemented under `MEMORY-09-ISOLATION-HARDENING-AUTHORIZATION-DRAFT-2026-11-17.md` (signed 2026-11-17, P6
amended). Nothing executed: no test, probe, compile, lint or import. Every statement below is a source claim.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Authorisation | S1–S6 only, with X1–X6 excluded |
| Implemented | **Yes** — six sites |
| Executed | **Nothing** |
| Conformance classifications | **Unchanged** |
| New finding | **A seventh live-dict return exists** (`withdraw_declaration`); reported, **not fixed** — outside S1–S6 |

---

## 1. What changed

| # | File | Site | Change |
|---|---|---|---|
| **S1** | `engine.py` | `declare`, final return (l.700-704) | `return dict(decl, audience=list(...), covered=list(...))` — the stored `decl` is still built and mutated internally, so withdrawal and coverage keep working |
| **S2** | `engine.py` | `assert_incompatible`, final return (l.811-817) | copy of the group with `targets`, `members`, `propositions`, `authority_ref` copied — the stored group is what every `status == "open"` read consults |
| **S3** | `engine.py` | `disclose` (l.607-624, l.639-642) | a nested `_account_descriptor()` returns `{..., "event": dict(event)}` per recipient, used in both pre-flight and the write loop — no two accounts share an event object |
| **S4** | `engine.py` | `promote` (l.419) | `"basis": [dict(vars(e)) for e in entries]` |
| **S5** | `records.py` | `Record.as_dict` (l.225, l.234) | `copy.deepcopy` for `provenance` and `corrected_by` |
| **S6** | `engine.py` | `snapshot` (l.1175-1177) | `copy.deepcopy` for `groups`, `promotions`, `declarations`; `propositions` left as `dict(v)` |

Both modules gained `import copy` (`engine.py:30`, `records.py:32`).

**W1 (2026-11-17, separate authorisation)** — `withdraw_declaration` now returns a copy as well (l.717-720). It was
not part of S1–S6 and was not folded into them; see §3.2.

**Files not changed:** `authority.py`, `errors.py`, `retrieval.py`, `__init__.py` — as §5 of the authorisation required.

---

## 2. Required properties, as satisfied

| # | Property | Finding |
|---|---|---|
| P1 | None of the six paths exposes a stored dict, list or nested value by reference | **Holds by inspection.** S1/S2 copy the outer dict *and* every list value; S3 copies the event per account; S4 copies each basis dict; S5 deep-copies both nested fields; S6 deep-copies the three registries |
| P2 | Read-only callers observe identical values, in identical order | **Holds by inspection.** Every change is a copy of the same content: dict order preserved by `dict(...)` construction, list order preserved, no sorting or filtering introduced |
| P3 | No signature or return-type change | **Holds.** `declare` → `dict`; `assert_incompatible` → `dict`; `disclose` → `(Record, tuple[Record, ...])`; `snapshot` → `dict`; `as_dict` → `dict` |
| P4 | No refusal, identifier, ordering, version bump or stored value changes | **Holds.** No `raise` added or removed; `_id`, `_seq*`, `_version` untouched; all stored structures still built exactly as before. S1's internal `decl["covered"] = [...]` still mutates the stored dict |
| P5 | No new field, record class, registry, operation or requirement | **Holds.** Two `import copy` statements and six call-site changes; nothing exported changed |
| P6 | *Amended:* known inventory reported, newly discovered paths reported, no universality claimed | **See §3 — and the amendment earned its keep immediately** |

**Limits of this table:** P1 and P2 are claims about code, not observations. Deep-copy behaviour across every value
shape (a tuple inside a dict, nested lists, `None`, shared sub-objects) cannot be confirmed without execution, and
none was run. If `deepcopy` were to fail on some value shape, it would fail at runtime — undetected here.

---

## 3. Known remaining shared-object paths (P6, amended)

**This is the inventory source inspection establishes, not a proof that no other alias exists.** Even after the
seventh correction in §3.2, the inventory will still be described this way: known paths, reported as found, with
the source-only limitation intact. Corrections close demonstrated leaks; they do not certify isolation, and no
report of this kind should be read as doing so.

### 3.1 Known before this change, still present by design (X1)

| Path | What stays live | Note |
|---|---|---|
| `claim`, `acquire`, `disclose`, `supersede_authority`, `correct`, `retract`, `declare_belief`, … | `Record` objects | Content, status, `corrected_by` and provenance remain writable through a returned record; `__setattr__` guards only `proposition_id` |
| `history`, `all_records`, `retrieval.perspective` | `Record` objects | same |
| `AuthorityRegistry.grants_for` | `Grant` objects | authority inputs remain writable |
| a live claim's `disclosure` dict | the metadata itself | `claim.disclosure["source_class"] = …` still writes through (the *serialised* view is isolated — S5 and the existing fresh-dict construction in `as_dict`) |

### 3.2 Newly discovered during this work — **corrected under W1**

| Path | Finding | Disposition |
|---|---|---|
| `withdraw_declaration` (engine.py:717) | Returned the **stored** `decl` dict, exactly the defect class S1 corrected in `declare`. A caller could undo the withdrawal it had just performed (`decl["withdrawn_turn"] = None`), re-activating a declaration that B6 says confers nothing | **Fixed under W1** (2026-11-17, its own authorisation: `MEMORY-09-WITHDRAWAL-COPY-AUTHORIZATION-DRAFT-2026-11-17.md`). The stored dict still receives `withdrawn_turn = turn`, so B6 holds; only the returned handle changed. Source claim, not observed |

This is the amended P6 working as intended: the original "exactly three items" wording would have asserted
completeness that the code does not support.

### 3.3 Latent, not currently shared

| Site | Why safe today | What would change that |
|---|---|---|
| `as_dict()["attribution"]` — `dict(self.attribution or {})` | every value is a string | an `attribution` value becoming a dict or list |
| `snapshot()["propositions"]` — `dict(v)` | values are scalars and one immutable tuple | a proposition field gaining a mutable value |
| `Record.domain` | a tuple | it becoming a list |

Recorded so the next inspection does not have to rediscover them; none is a leak today.

---

## 4. What was not changed

| # | Excluded item | Status |
|---|---|---|
| X1 | Freezing `Record`/`Grant` or copying on return | **Not done.** The API decision remains open and is the reason §3.1 exists |
| X2 | Memory semantics: ownership, acquisition, disclosure, belief, promotion, authority, lifecycle, containment | **Untouched** |
| X3 | Howl activation/lifecycle | **Untouched** — still outside this experiment |
| X4 | v9, 57 rules, 118 requirements, r32, plan v3, Amendment 1, prior decisions | **Untouched** |
| X5 | Testing, probes, compilation, linting, imports, runtime | **Nothing run** |
| X6 | Refactoring beyond the six sites | **Not done** — including §3.2, which is reported rather than patched |

---

## 5. Status, kept separate

| Question | Answer |
|---|---|
| Are S1–S6 implemented? | **Yes** — six sites in `engine.py` and `records.py` |
| Is the isolation *demonstrated*? | **No.** P1/P2 are source claims. Nothing has been imported, so even a syntax error is undetected |
| Did any classification move? | **No** — rule 46 **P**, J3 **P**, K9 **P**, B3 **NA**, F7 **NA**; 56 I · 1 P and 113 I · 3 P · 2 NA |
| Is testing authorised? | **No** |

### Source hashes after the change

| File | Hash | Status |
|---|---|---|
| `engine.py` | `c05a6a83974118cb…` | changed (was `780a9cd760175062…` after S1–S6; W1 applied on top) |
| `records.py` | `28c37c547c87d84a…` | changed (was `778684e2a63d82d9…`) |
| `authority.py` | `f344582350af88f1…` | unchanged |
| `errors.py` | `48125d7ebae8c013…` | unchanged |
| `retrieval.py` | `cd2d380bd5d9d1f7…` | unchanged |
| `__init__.py` | `06ccf8f60625ff2c…` | unchanged |

---

## 6. What remains before this is trustworthy

1. **One-line follow-up for §3.2** — **done**: W1 applied 2026-11-17 under its own authorisation. What remains is not another copy fix, but the X1 decision.
2. **X1 decision** — whether `Record`/`Grant` returns are accepted as an in-process trust boundary or isolated at
   the API. This is the only remaining way stored state can be rewritten through an interface.
3. **Execution** — the only thing that can turn P1/P2 from readings into findings. Still unauthorised.
