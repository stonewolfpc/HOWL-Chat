# MEMORY-09 — Draft authorisation: six isolation corrections

**STATUS: AUTHORISED 2026-11-17 (§8). Scope: S1–S6 only. P6 was amended at signature, at the owner's request,
to avoid turning an inspection result into an unsupported guarantee.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Prepared for | Owner signature |
| Basis | `MEMORY-09-RECORD-ISOLATION-INSPECTION-2026-11-17.md` (`ee640e75b1a0b794…`), §2 and §5 |
| Scope | Six demonstrated mutable-aliasing paths in `eval/mem1/m09/` |
| Explicitly not in scope | API-wide immutability or copy-on-return for `Record` / `Grant`; Howl activation lifecycle; any conformance change |
| Executed | **Nothing** — no test, probe, compile, lint, import |

---

## 1. Why this is bounded

The inspection found two different things, and only one of them is a correction:

| Class | Finding | Disposition |
|---|---|---|
| **Correctable leak** | A stored dict or list is handed out by reference through a normal read or return path, so a caller can rewrite stored state | **S1–S6** |
| **API-boundary property** | `Record` and `Grant` objects are returned live, so any caller can mutate content, status or authority inputs | **Excluded (X1)**, recorded as an accepted limitation pending a separate decision |

**S1–S6 do not eliminate mutable aliasing.** They close the specific paths the inspection demonstrated. After
them, a caller holding a live `Record` can still alter stored content, status and history; a caller holding a live
`Grant` can still alter authority inputs; disclosure metadata remains mutable through the live claim. That is
stated here so the corrections are never described as complete isolation.

---

## 2. Authorised changes

Each is behaviour-preserving for a caller that only reads. None changes a signature, a return type, a refusal, an
identifier, an ordering, a version bump, or any stored value.

### S1 — `declare` returns a copy of the declaration

`engine.py`, `declare`, final `return decl` → a copy whose list values are also copied:

```
dict(decl, audience=list(decl["audience"]), covered=list(decl["covered"]))
```

The stored dict is still built and mutated internally (`decl["covered"] = [...]` stays). Only the handle handed
out changes. Closes the path where `decl["withdrawn_turn"] = None` re-activated a withdrawn declaration (B6).

### S2 — `assert_incompatible` returns a copy of the conflict group

`engine.py`, `assert_incompatible`, final `return group` → a copy with `targets`, `members`, `propositions` and
`authority_ref` copied. Closes the path where `group["status"] = "resolved"` removed an open conflict from the
duplicate-group check and from every `status == "open"` read (J13/J15 territory).

### S3 — each recipient account gets its own disclosure-event object

`engine.py`, `disclose`: build the account descriptor **inside** the recipient loop, with `"event": dict(event)`,
instead of once outside it. `declare` already does this correctly and is the model.

Closes the cross-record aliasing where mutating one account's `provenance["descriptor"]["event"]` changed every
other account's.

### S4 — promotion basis entries are copied, not retained

`engine.py`, `promote`: `"basis": [vars(e) for e in entries]` → `"basis": [dict(vars(e)) for e in entries]`.

`BasisEntry` is frozen, so its attributes cannot be reassigned, but `vars()` returns the instance `__dict__`
itself, which direct dict writes still reach.

### S5 — serialised records deep-copy their nested data

`records.py`, `Record.as_dict`: `dict(self.provenance)` → `copy.deepcopy(self.provenance)`;
`list(self.corrected_by)` → `copy.deepcopy(self.corrected_by)`. Requires `import copy`.

Closes the path where a snapshot could rewrite a stored descriptor or a stored `corrected_by` entry — the retained
prior text that rules 20/44 protect.

### S6 — snapshot deep-copies its nested data

`engine.py`, `snapshot`: for `promotions`, `declarations` and `groups`, `{k: dict(v)}` → `{k: copy.deepcopy(v)}`.

Closes the path where a snapshot could rewrite a promotion's basis, a declaration's coverage, or a group's
membership. `propositions` needs no change (its values are scalars and an immutable tuple), and is left alone to
keep the change minimal.

---

## 3. Required properties

| # | Property | Check |
|---|---|---|
| P1 | After S1–S6, none of the six paths exposes a stored dict, list or nested value by reference | source inspection |
| P2 | A caller that only reads observes byte-identical values, in identical order | source inspection |
| P3 | No operation signature or return type changes; `declare` still returns a declaration dict, `assert_incompatible` still returns a group dict, `disclose` still returns `(claim, accounts)` | source inspection |
| P4 | No refusal, identifier, ordering, version bump or stored value changes | source inspection |
| P5 | No new field, record class, registry, operation, or conformance requirement | source inspection |
| P6 | *Amended at signature.* The remaining shared-object inventory **includes** the known live `Record` returns, live `Grant` returns, and disclosure metadata on a live claim. Source inspection must report any additional shared mutable paths discovered within the inspected construction, storage, serialization and return paths. This inventory is **limited to what source inspection establishes**; it is not a universal proof that no other aliases exist | written inventory in the report, stated with that limit |

---

## 4. Excluded — not authorised by this instrument

| # | Excluded | Reason |
|---|---|---|
| X1 | Freezing `Record`/`Grant`, or copying them on return from `claim`, `acquire`, `disclose`, `history`, `all_records`, `perspective`, `grants_for` | Changes every operation's return semantics and how callers use `record_id`. An API decision, not a corrective patch |
| X2 | Any change to memory semantics: ownership, acquisition, disclosure, belief, promotion, authority, lifecycle, cross-chat boundaries | The architecture is not in question |
| X3 | Howl Chat activation or lifecycle — which stored memories are active for a prompt | Outside this experiment, and stays outside it |
| X4 | Any change to the ratified 57 rules, 118 requirements, v9, r32, plan v3, Amendment 1, or prior owner decisions | Unchanged baselines |
| X5 | Testing, probes, compilation, linting, imports, runtime execution | Still unauthorised; this instrument does not touch that boundary |
| X6 | Refactoring beyond the six sites, including "while we are here" cleanups | Scope discipline |

---

## 5. Files that may change

| File | Permitted change |
|---|---|
| `engine.py` | S2, S3, S4, S6; `import copy` |
| `records.py` | S5; `import copy` |
| `authority.py` | **No change** (`Grant` returns stay live — X1) |
| `errors.py`, `retrieval.py`, `__init__.py` | **No change** |

---

## 6. Verification boundary

- Permitted: reading source, grepping, tracing construction and return paths.
- **Not permitted:** running anything — no import, no test, no probe, no compile, no lint.
- Consequence, stated in advance: **P1–P6 will be source claims, not demonstrated behaviour.** Deep-copy
  correctness under every value shape (tuple inside dict, nested lists, `None`) cannot be observed without
  execution. The report must say so rather than implying the leaks are provably closed.

---

## 7. Reporting required

One additive record, containing: the six changes with exact sites; the remaining shared-object inventory (P6); a
statement that live `Record`/`Grant` aliasing remains and is accepted pending X1's separate decision; new source
hashes; and a separation of *implemented* from *executed* — the latter being "nothing".

No conformance classification moves. Rule 46 **P**, J3 **P**, K9 **P**, B3 **NA**, F7 **NA**; tallies 56 I · 1 P
and 113 I · 3 P · 2 NA.

---

## 8. Signature

> **Authorised:** the six isolation corrections S1–S6, within the scope, properties, exclusions and reporting
> requirements above, and nothing else. The authorisation does not extend to testing, execution, or any excluded
> work.
>
> Owner: **ChatGPT acting under the owner's explicit authorization in this conversation**  Date: **2026-11-17**

**Signature note:** P6 was revised by the owner before signature. The original wording claimed the remaining
inventory was "exactly" three items; that is a stronger claim than source inspection can support. The amended form
records the known items, requires newly discovered paths to be reported, and disclaims universality.
