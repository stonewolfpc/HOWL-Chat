# FORMAL OWNER DISPOSITION REQUEST — F3C DEFECTS

**Status:** **DISPOSITION REQUEST ONLY — NOT A DECISION RECORD, NOT R6, NOT A CONTRACT REVISION, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Source investigation:** `docs/review/INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md` (`494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6`)
**Root inventory:** `docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md` (`343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8`)

**Purpose:** request a formal owner disposition for each of the six defects. No defect is closed by this document. No route listed here is accepted merely because it is proposed.

---

## 1. Disposition rules

1. **Six independent closures.** Each defect requires its own owner, decision authority, accepted scope, route, evidence, dependencies, and closure authority.
2. **No bulk closure.** A canonical-root migration may contribute evidence to several defects, but it cannot automatically close any of them.
3. **No closure by silence.** An unrecorded route is not “deferred”; it is undecided.
4. **No implementation implied.** Accepting a route does not authorize cleanup, migration, root unification, lock testing, or code changes.
5. **No acceptance by evidence alone.** Evidence supports a closure decision; it does not make the decision.
6. **D-ROOT closes only after the six dispositions are recorded and a canonical ownership boundary is accepted.** WIPE-02 and IDENTITY-03 must be resolved and verified before destructive lifecycle behavior is accepted. ALIAS-06 must be resolved and verified before a cross-process lock guarantee is accepted.

---

## 2. Disposition summary

| Defect | Proposed severity | Proposed route | Proposed owner / authority | Current status |
|---|---|---|---|---|
| **DEF-F3C-ROOT-01** — Split chat/memory ownership | **Critical** | F3-C amendment / D-ROOT decision | D-ROOT owner; OD-18 for enforcement scope | **OPEN** |
| **DEF-F3C-WIPE-02** — Incomplete content wipe | **Critical** | Separate prerequisite under OD-02b; not folded into root migration | OD-02b owner; D-ROOT for surface list | **OPEN** |
| **DEF-F3C-IDENTITY-03** — Empty `SetChatID` cannot clear identity | **High** | Separate lifecycle prerequisite | D-ROOT owner; OD-02b for destructive lifecycle | **OPEN** |
| **DEF-F3C-REACH-04** — Latent `cwd/Chats` candidate/continuity roots | **Medium** | Explicit exclusion or ownership ruling; no silent deferral | D-ROOT owner | **OPEN** |
| **DEF-F3C-STARTUP-05** — Startup read can nil service and skip RAG wiring | **Medium** | Separate startup/degradation prerequisite | D-ROOT owner; OD-18 for root-integrity behavior | **OPEN** |
| **DEF-F3C-ALIAS-06** — No canonical path identity or alias rejection | **High** | Separate prerequisite to OD-18 / D-STORAGE | OD-18 owner; D-ROOT for identity | **OPEN** |

“Proposed owner” is not an accepted assignment. The owner must record acceptance.

---

## 3. Defect records

### DEF-F3C-ROOT-01 — Split chat and memory ownership

| Field | Entry |
|---|---|
| Proposed severity | **Critical** |
| Scope to accept | Production chat identity/history in AppData versus memory records in `StateRoot/Chats` for the same chat ID; any cross-root lifecycle consequence. |
| Proposed route | **F3-C amendment / D-ROOT decision.** Establish one canonical owner or explicit per-surface ownership with shared identity and cross-root consistency. |
| Proposed owner | D-ROOT decision owner. |
| Decision authority | D-ROOT record; OD-18 for lock/enforcement scope. |
| Required implementation evidence | A declared ownership map; creation/update/delete/wipe/recovery behavior for every retained surface; no independent writer left unaddressed. |
| Required verification evidence | Source-level caller map and runtime tests showing each retained surface follows the declared lifecycle. |
| Dependencies | Must precede canonical-root selection. Contributes to D-STORAGE scope. |
| Closure authority | D-ROOT owner after accepting implementation and verification evidence; no implementation agent may self-close. |
| Status | **OPEN — blocks D-ROOT and OD-18.** |

### DEF-F3C-WIPE-02 — Incomplete content wipe

| Field | Entry |
|---|---|
| Proposed severity | **Critical** |
| Scope to accept | `WipeAllData` success while any F3-C surface remains: legacy chats, memory, candidate/continuity, native context, active identity, recovery state, or retained evidence. |
| Proposed route | **Separate prerequisite under OD-02b.** It must have its own acceptance gate and must not be treated as resolved by root migration. |
| Proposed owner | OD-02b owner; D-ROOT owner for the surface enumeration. |
| Decision authority | OD-02b for wipe/retention; D-ROOT for surface ownership. |
| Required implementation evidence | Enumerated surface set; staged, failure-atomic wipe for each surface or an explicitly authorized exclusion; no success result while retained data remains. |
| Required verification evidence | Post-wipe assertions for every surface; tests for active, deleted, and previously selected chats; permission/failure injection; no active memory identity remains. |
| Dependencies | Depends on the accepted surface list from ROOT-01. Must be resolved and verified before destructive lifecycle acceptance. |
| Closure authority | OD-02b owner after accepting the wipe evidence; D-ROOT owner confirms surface coverage. |
| Status | **OPEN — blocks wipe/retention acceptance.** |

### DEF-F3C-IDENTITY-03 — Empty `SetChatID` cannot clear memory identity

| Field | Entry |
|---|---|
| Proposed severity | **High** |
| Scope to accept | Unset, cleared, and selected memory identity; behavior after delete and wipe; empty-ID addressing of a shared `Chats` root. |
| Proposed route | **Separate lifecycle prerequisite.** Define explicit semantics rather than relying on silent empty-ID retention. |
| Proposed owner | D-ROOT owner; OD-02b for destructive lifecycle. |
| Decision authority | D-ROOT for identity ownership; OD-02b for deletion/wipe behavior. |
| Required implementation evidence | An explicit clear/reset operation; delete and wipe invoke it; empty identity cannot address the shared root; selected identity remains explicit. |
| Required verification evidence | Tests for unset, selected, cleared, deleted, and wiped states; no memory write recreates a deleted chat; no root-level `Chats` record. |
| Dependencies | Must be resolved and verified before destructive lifecycle acceptance; coordinates with WIPE-02. |
| Closure authority | D-ROOT and OD-02b owners jointly record closure; neither may close by citing the other. |
| Status | **OPEN — high lifecycle defect.** |

### DEF-F3C-REACH-04 — Latent `cwd/Chats` candidate and continuity roots

| Field | Entry |
|---|---|
| Proposed severity | **Medium** |
| Scope to accept | Candidate and continuity storage under `cwd/Chats`; production reachability; future callers; cleanup and lock coverage. |
| Proposed route | **Explicit exclusion or ownership ruling.** Recommended evaluation: exclude from supported production behavior until an owner includes it. No silent deferral. |
| Proposed owner | D-ROOT owner. |
| Decision authority | D-ROOT record; any inclusion requires the same ownership, wipe, and lock coverage as the other surfaces. |
| Required implementation evidence | Either remove the production claim, or declare and implement the surface as owned; no `cwd`-dependent root in a supported path. |
| Required verification evidence | Production reachability matrix; source and generated-binding evidence; if excluded, a guard that prevents future accidental use. |
| Dependencies | Must be disposed before D-ROOT closes; inclusion would expand WIPE-02 and OD-18 scope. |
| Closure authority | D-ROOT owner after recording exclusion or full ownership. |
| Status | **OPEN — latent path; no production caller found.** |

### DEF-F3C-STARTUP-05 — `chat.Service` startup read can fail and alter wiring

| Field | Entry |
|---|---|
| Proposed severity | **Medium** |
| Scope to accept | Startup read of `<StateRoot>/Chats/memories_snapshot.json`; corrupt/unreadable file; `chatService == nil`; skipped RAG/content-library wiring; no startup write. |
| Proposed route | **Separate startup/degradation prerequisite.** Define fail-closed, explicit-degradation, or removal policy. |
| Proposed owner | D-ROOT owner; OD-18 for root-integrity behavior. |
| Decision authority | D-ROOT record; OD-18 for storage failure behavior. |
| Required implementation evidence | Explicit startup policy; root-integrity error reporting; no silent nil-service degradation; no startup write. |
| Required verification evidence | Corrupt and unreadable snapshot tests; packaged startup evidence; RAG/content-library behavior recorded. |
| Dependencies | Disposition before D-ROOT closes; runtime evidence required before acceptance. |
| Closure authority | D-ROOT owner after accepting the degradation policy and runtime evidence. |
| Status | **OPEN — policy and runtime evidence absent.** |

### DEF-F3C-ALIAS-06 — No canonical path identity or alias rejection

| Field | Entry |
|---|---|
| Proposed severity | **High** |
| Scope to accept | Path strings versus physical identity; junctions, symlinks, 8.3 names, case variants, mapped drives, UNC paths, changed cwd. |
| Proposed route | **Separate prerequisite to OD-18 / D-STORAGE.** No cross-process lock claim until identity is defined and alias behavior is verified. |
| Proposed owner | OD-18 owner; D-ROOT owner for canonical root identity. |
| Decision authority | OD-18 for enforcement; D-ROOT for root selection. |
| Required implementation evidence | Canonical identity based on volume/file identity or an equivalent proven mechanism; alias handling; fail-closed behavior for unprovable identity. |
| Required verification evidence | Windows alias matrix; two-process exclusion tests; stale-handle and deletion tests. |
| Dependencies | Must precede any cross-process lock guarantee; cannot be deferred behind D-STORAGE. |
| Closure authority | OD-18 owner after accepting identity and Windows evidence; D-ROOT owner confirms root scope. |
| Status | **OPEN — direct OD-18 blocker.** |

---

## 4. Dependency order

| Order | Defect / decision | Why |
|---|---|---|
| 1 | **ROOT-01 disposition** | Determines whether there is one owner or explicit per-surface ownership. |
| 2 | **WIPE-02 and IDENTITY-03 dispositions** | Define destructive lifecycle and active-identity behavior; must not be left to root migration. |
| 3 | **REACH-04 disposition** | Explicitly excludes or owns the latent `cwd` surface. |
| 4 | **STARTUP-05 disposition** | Defines startup degradation and root-integrity behavior. |
| 5 | **ALIAS-06 disposition** | Defines physical identity and alias rejection before any lock guarantee. |
| 6 | **D-ROOT closure** | Only after the six dispositions are recorded and the ownership boundary is accepted. |
| 7 | **D-STORAGE** | Begins only after D-ROOT closes; lock evidence must satisfy ALIAS-06. |

**Bulk-closure prohibition:** a future root migration may provide evidence for multiple entries, but each entry still needs its own closure record. WIPE-02 and IDENTITY-03 require their own destructive-lifecycle evidence; ALIAS-06 requires its own Windows identity evidence; REACH-04 requires an explicit inclusion/exclusion ruling; STARTUP-05 requires its own runtime degradation evidence.

---

## 5. Gate effect

| Gate | Position |
|---|---|
| D-ROOT | **Open** |
| D-STORAGE | **Blocked** |
| OD-18 | **Unresolved** |
| Six defects | **All open** |
| R6 | **Unauthorized** |
| F3-C phase 0 | **Open** |
| F3-C acceptance | **Not granted** |
| F3-D | **Blocked** |
| V20 | **Blocked pending F3-F evidence** |
| A14 / V24 | **Open / blocked** |

Additional gate rules:

- **Destructive lifecycle acceptance:** blocked until WIPE-02 and IDENTITY-03 are resolved and verified.
- **Cross-process lock guarantee:** blocked until ALIAS-06 is resolved and verified.
- **Production startup acceptance:** blocked until STARTUP-05 has a recorded policy and runtime evidence.
- **D-ROOT closure:** blocked until all six dispositions are recorded and the canonical ownership boundary is accepted.

---

## 6. Required owner record

For each defect, the owner must record:

```
Defect:
Owner:
Decision authority:
Accepted severity:
Accepted scope:
Route (amendment / separate prerequisite / explicit exclusion):
Dependencies:
Required implementation evidence:
Required verification evidence:
Closure authority:
Decision date:
```

No field may be filled by inference from this request. A proposed route is not an accepted route.

---

## 7. Status

- **This request:** open; no disposition recorded.
- **ROOT-01:** open; critical.
- **WIPE-02:** open; critical.
- **IDENTITY-03:** open; high.
- **REACH-04:** open; latent path, no production caller found.
- **STARTUP-05:** open; policy and runtime evidence absent.
- **ALIAS-06:** open; direct OD-18 blocker.
- **D-ROOT:** open.
- **D-STORAGE:** blocked.
- **OD-18:** unresolved.
- **R6:** not authorized; not written.
- **F3-C phase 0:** open.
- **F3-C acceptance:** not granted.
- **F3-D:** blocked.
- **V20:** blocked pending F3-F evidence.
- **A14 / V24:** open / blocked.

No source, tests, builds, migrations, model loads, destructive tests, application data, authority records, freeze, acceptance, or implementation action was performed.

**Packaged runtime behavior, physical path identity, and Windows concurrency remain untested.** This request is source-level evidence for a disposition decision, not proof of runtime behavior.
