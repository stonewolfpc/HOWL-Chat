# TASK-03 — clarification 1j: the write ordering, and the handle-bound sequence

| Field | Value |
|---|---|
| Purpose | Close two ordering/specification checks raised against 1i |
| Status | **Contract text. Nothing adopted. No implementation authorized. Deletion unavailable** |
| Relationship | **Additive.** None of rows 528–538 is edited |
| Supersedes | 1i V-03 (understated what `removed` records); 1i V-17 (incomplete sequence); 1e Q-09 restated, not changed |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1J-2026-11-22` |

Prefixes: `W-` normative, `WT-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

**Scope note.** Both checks were completed before drafting. The first found a defect introduced by 1i; the second found an under-specification. Nothing else in 1i required change, and this document adds no new outcome, kind, state or owner decision.

---

## 1. Check one: does `removed` precede verification?

**W-01 — The ordering was already fixed, in 1e.** `TASK-03-CONTRACT-CLARIFY-1E-2026-11-22.md` states it twice:

> **Q-09** — `removed` is written only after the acquired identity is shown to be absent. Path absence alone never advances the state machine and never supports `complete`.
>
> `prepared → intent → pending → removing → [destructive removal] → verify identity absent → removed → terminal → cleanup`

So the contract **does not** permit writing `removed` before verification, and 1j does not change Q-09.

**W-02 — But 1i contradicted it.** 1i V-03 described the marker as "evidence that a removal was attempted and that its call returned" — which is what the marker would record if it were written before verification. It is not. `removed` records **both** the successful removal call **and** the absent observation that authorized it, in that order (Q-09). 1i's sentence is wrong and is corrected here.

**W-03 — The defect is not cosmetic.** Read literally, V-03 would allow a marker meaning "the call returned", which combined with a K1 in another process could be arranged into a completion claim. Under Q-09 the marker cannot mean that. The correction closes the reading, and V-05's context dimension remains what prevents the cross-process claim.

**W-04 — `complete` is a fresh claim, never a reinterpretation.** The terminal receipt is written in the same process, immediately after `removed`, from the observation that authorized `removed`. It is **not** derived from the marker by a later reader, and no reader may reconstruct `complete` from a marker plus an independent observation.

| Derivation | Permitted |
|---|---|
| K1 observed in `original` context authorizes `removed`, then terminal `complete` | yes |
| marker + later K1 in `reconciliation` | **no** — `completion_unconfirmed` (1i V-09) |
| marker alone, no observation | **no** — nothing is written |

**W-05 — The new crash point is already covered, and it is not window 11.** Fixing the ordering exposes a distinct crash position: removal returned, verification not yet done, so `removed` was never written. Reconciliation then finds state `removing` with K1 → **`absent_unverified`** (1i V-05), which claims nothing about completion. Window 11 remains exactly "marker written, terminal not yet written". No new outcome is needed, and none is added.

**W-06 — A failed or ambiguous removal writes nothing.** If the removal call returns an error, or its result cannot be established, `removed` is not written and the state remains `removing`. The operation never "optimistically" advances and never writes `removed` on the expectation that verification will later agree.

**W-07 — Verification is read-only.** It opens, identifies, compares and closes. It performs no write, no repair, no retry of removal, and no second destructive attempt. There is no re-verification loop that could upgrade an outcome.

---

## 2. Check two: is the handle-bound sequence specified tightly enough?

**W-08 — It was not.** 1i V-17 required that the directory be opened and identified through the held volume handle and that removal act on that opened object. It did not specify how children are reached or how verification is bound. That gap is closed here — as **required design**, still gated on evidence (W-14).

**W-09 — The sequence, as four required properties.** All four must hold for the guarantee to exist:

> 1. **Open.** The acquired directory identity is opened **by ID through the held volume handle**, and its identity is confirmed equal to the acquired identity from `FileIdInfo` on that same handle.
> 2. **Enumerate.** Children are enumerated **through the directory object obtained in step 1** — by handle-relative enumeration on that object, not by rebuilding a path.
> 3. **Remove.** Each child is removed through a handle derived from that same directory object. Removal never acts on a path string.
> 4. **Verify.** The acquired directory identity is shown absent **by ID open on the same held volume handle**, in the same process (Q-09, 1i V-06).

**W-10 — No path fallback at any step.** If any of the four requires path-based re-resolution — including "open the directory by ID, then enumerate by path", or "enumerate by handle, then remove by path" — the guarantee is **not established**. A mixed sequence is not a weaker guarantee; it is a different, unsupported design.

**W-11 — What the guarantee is, stated narrowly.** When all four hold, the deletion is bound to the object the operation acquired, and the completion claim rests on that object. It still does **not** claim:

- continuous physical device identity (1h U-22);
- exclusion of identity reuse across the interval, so `attribution` remains `unknown` (1h U-11);
- that a clone cannot exist elsewhere.

**W-12 — A partial demonstration enables nothing.** If steps 1–3 are demonstrated and step 4 is not, deletion is not "partly available". G-1 is a single gate: unmet until the whole sequence is established.

---

## 3. Evidence status

**W-13 — This is feasibility, not contract.** W-09 is a specification of what must be demonstrated. Whether Windows permits it is **Unestablished**.

**W-14 — E-27 is split, because it was carrying two questions.**

| # | Question | Status |
|---|---|---|
| **E-27** | Does a held volume handle keep the calls bound to that volume object across the operation? | Unestablished |
| **E-28** | Can the four-step sequence of W-09 be performed at all — in particular, can children be enumerated and removed through handles derived from the directory object, with no path re-resolution? | Unestablished |

E-25, E-26 and E-24 are unchanged.

**W-15 — The consequence of E-28 failing.** If children cannot be enumerated and removed without path re-resolution, the identity-bound design has no implementation, `complete` stays unreachable, and deletion remains unavailable. It does **not** license a path-based fallback (1i V-23).

---

## 4. Preserved

**W-16 — Carried unchanged.** Q-09's ordering. The 1h failure taxonomy and control open. The 1i context dimension and its matrix — including that only `removed` + K1 + `original` yields `complete`. `completion_unconfirmed`, immutable and never upgraded. No re-run from a marker. Count-only baseline, no manifest, `count_delta_unattributed` with its non-readings and the string-table lint. `absent_after_acquisition`, `absent_unverified`, `aborted`. FAT/exFAT/SMB refusal. D-10 unchanged; G-1 unchanged.

**W-17 — No new owner decision, no new outcome, no new evidence kind.** Five ledger kinds and two evidence kinds still hold.

---

## 5. Acceptance tests

- **WT-01** `removed` written before a K1 observation is rejected as an invalid transition (W-01, W-02).
- **WT-02** No code path derives `complete` from the marker plus an independent observation (W-04).
- **WT-03** Removal-call success followed by a crash, then restart with K1 → `absent_unverified`, never `completion_unconfirmed` and never `complete` (W-05).
- **WT-04** Removal call erroring → no `removed` write; state remains `removing` (W-06).
- **WT-05** Verification performs zero writes and cannot be invoked twice to change an outcome (W-07).
- **WT-06** Child enumeration and child removal are performed on handles derived from the directory object; a path string reaching either step fails the test (W-09, W-10).
- **WT-07** Verification opens by ID on the same held volume handle; verification by path fails the test.
- **WT-08** Steps 1–3 passing with step 4 failing leaves deletion unavailable; no partial enablement (W-12).
- **WT-09** The documented residual claims of W-11 are absent from any receipt or string produced in these paths.

---

## 6. Verification and authorization

Rows 528–538 remain byte-identical; 1i (`3f00fc6cfa6d…`) unchanged. 1j registered as row 539; counts updated; no historical entry altered.

**Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed. Nothing here was validated against Windows: E-24, E-25, E-26, E-27 and E-28 are all open, and W-09 is specification only.

**Contract drafting on this branch stops here.** The two raised checks are closed; no further speculative expansion is warranted, and the next evidence must come from a bounded feasibility investigation of W-09, not from more text.

TASK-03 remains in contract preparation. Not accepted, not closed. Deletion remains **unavailable**.
