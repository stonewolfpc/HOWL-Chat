# MEMORY-09 — Draft authorisation: W1, `withdraw_declaration` returns a copy

**STATUS: AUTHORISED 2026-11-17 (§6) and applied. W1 only. Source inspection only; nothing executed.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Basis | `MEMORY-09-ISOLATION-HARDENING-REPORT-2026-11-17.md` §3.2 — the seventh path, found during S1–S6 and deliberately not fixed |
| Scope | **One return expression**, one function |
| Executed | **Nothing** |

---

## 1. What is authorised

`engine.py`, `withdraw_declaration`, final `return decl` (l.717) →

```
return dict(decl, audience=list(decl["audience"]), covered=list(decl["covered"]))
```

The stored dict remains the stored dict: `decl["withdrawn_turn"] = turn` still mutates it, so withdrawal continues
to take effect for `_validate_evidence`'s `narrator_grant` branch (B6). Only the handle handed to the caller changes.

**Why it matters:** today a caller can call `withdraw_declaration(...)`, receive the stored dict, write
`result["withdrawn_turn"] = None`, and re-activate a declaration whose withdrawal it just performed — the same
defect class S1 closed in `declare`.

---

## 2. Required properties

| # | Property |
|---|---|
| W-P1 | The returned dict shares no dict or list with the stored declaration |
| W-P2 | A read-only caller observes identical values |
| W-P3 | Signature and return type unchanged (`-> dict`) |
| W-P4 | The stored declaration still mutates internally; withdrawal still takes effect; B6 unchanged |
| W-P5 | No new field, class, operation, requirement, refusal, identifier or ordering change |
| W-P6 | Inventory reported as **known paths**, not as proof of complete isolation |

---

## 3. Excluded

| # | Excluded |
|---|---|
| X1 | Freezing or copy-on-return for `Record`/`Grant` — still a separate API decision |
| X2 | Any other `return` in the module, including a sweep for further live handles. If one is spotted it is **reported, not fixed** |
| X3 | Howl activation/lifecycle |
| X4 | Any ratified baseline or conformance classification |
| X5 | Testing, execution, compilation, linting, imports |
| X6 | Refactoring, renaming, comment reflow beyond the one expression |

---

## 4. Files that may change

`engine.py` only. `records.py`, `authority.py`, `errors.py`, `retrieval.py`, `__init__.py`: **no change**.

---

## 5. Verification and reporting

Source inspection only: no import, no compile, no test. The report must state that W-P1/W-P2 are source claims, add
the new `engine.py` hash, and keep the known-path inventory language — a seventh correction does not make the
eighth impossible.

---

## 6. Signature

> **Authorised:** W1 only — `withdraw_declaration` returns a copy of its declaration, within the properties and
> exclusions above. Does not extend to testing, execution, or any excluded work.
>
> Owner: **owner, by explicit authorization in this conversation**  Date: **2026-11-17**

**Signature note:** the owner's approval restated the protections explicitly — the stored declaration still
receives `withdrawn_turn = turn`; only the returned handle changes; no other return path may be swept in;
`engine.py` is the only permitted file; nothing may run; and the report must keep the qualified known-path
inventory rather than claim complete isolation. Those restatements are recorded as part of the authorisation, not
as commentary on it.
