# DECISION — F3-C Root Authority and Found-Copy Precedence

**Date:** 2026-11-10
**Status:** resolved; store behavior implemented and tested
**Scope:** how `worldstore` treats a found copy when more than one root is registered
**Relates to:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` §3 (availability),
A10 (single writer), A11 (availability); `IMPLEMENTATION-F3C-STORE-FOUNDATION-2026-11-10.md`

## 1. The question

Suppose root A is the usual canonical location, root A is disconnected, and root B contains an
item with the requested identity. A lookup can prove existence, but not that root B's copy is
current. Does the store return B as an ordinary found item?

## 2. What the code actually had

- `library.Registration` has no priority, rank, or authority field. `DefaultSave` selects a
  **create** destination; it does not rank reads.
- `Registry.List` sorts roots by category and then path. That order is deterministic and
  reproducible, but it is not a statement of authority.
- `Lookup` returned the first readable copy and reported `Found` with no qualification. On a
  disposable single-root setup that is harmless; with a disconnected sibling it silently claims
  more than it knows.

So before this decision the honest answer was: precedence was not established at all, and the
result did not say so.

## 3. Decision

1. **Roots are equivalent search locations.** None is preferred by registration order, origin,
   role, or save-default status. Path order is retrieval order only, chosen so a result is
   reproducible.
2. **A found copy is authoritative only when both conditions hold:**
   - every registered root for the category was readable and searched without error; and
   - exactly one root holds the identity.
3. Otherwise the result is **found, authority indeterminate**: `Found=true`,
   `AuthorityIndeterminate=true`, with `Copies` and `UnavailableRoots`/`Errors` explaining why.
   `Item.AuthorityIndeterminate` carries the same fact through `Load`.
4. **An indeterminate found copy is readable for display but not editable.** `Editable` is false,
   and `Update` returns `authority_indeterminate` without touching any file. Editing a copy that
   may be shadowed would fork the item invisibly.
5. **Indeterminate is never absence.** It must not trigger creation, import, overwrite, duplicate
   cleanup, reference reassignment, or fallback to a secondary store.

### Behavior matrix

| Roots readable | Copies found | Result |
| --- | --- | --- |
| all | exactly 1 | authoritative found; editable if the owning root is writable |
| all | 2 or more | found, authority indeterminate (`copies: n`) |
| one or more unavailable, unreadable, or errored | 1 | found, authority indeterminate |
| one or more unavailable | 0 | incomplete search: `ErrUnavailable` |
| all | 0 | proven absence: `ErrNotFound` |
| no roots registered | 0 | incomplete search: `ErrUnavailable` |

## 4. Direct answers

- **Can a lower-priority root's copy be selected when a higher-priority root is unavailable?**
  There is no lower or higher priority. Any readable copy may be returned for display, and it is
  reported as authority-indeterminate unless every root was searched and it is the only copy.
- **Are roots ordered by authority?** No. They are sorted by path for reproducible retrieval.
  `DefaultSave` affects where new items are created, nothing else.
- **Does the result distinguish indeterminate authority from an ordinary find?** Yes:
  `Search.AuthorityIndeterminate`, `Search.Copies`, `Item.AuthorityIndeterminate`, and the
  `authority_indeterminate` update status.

## 5. Consequences

- A disconnected sibling can no longer be mistaken for "the canonical copy is here"; consumers see
  the item but must present it as not-current and must not edit it.
- Duplicate copies across two readable roots are treated the same way: ambiguous, not silently
  resolved by path order. A dedicated duplicate resolver is later work.
- No availability index is added, and no root preference is invented to avoid the ambiguity.
- Retrieval order remains deterministic so repeated reads do not flip between copies.

## 6. Compatibility

This clarifies the frozen availability contract; it does not change envelope, identity, revision,
reference, or ownership semantics. No data is moved or rewritten. The single-writer rule is
unchanged: an indeterminate state refuses the write rather than choosing a winner.

## 7. Evidence

Tests in `internal/worldstore`:

- unique copy across readable roots → authoritative, editable, update succeeds;
- duplicate copies across two readable roots → `Copies=2`, indeterminate, update refused,
  both files unchanged;
- found copy with a disconnected sibling root → indeterminate, `Load` preserves the flag,
  update refused, file byte-identical;
- unique copy with all roots readable → `Complete=true`, `AuthorityIndeterminate=false`.

The A13 regression suites remain the gate during consumer integration.
