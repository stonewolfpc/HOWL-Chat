# TASK-03 — OQ-7 clarification 1: capability lifetime and write containment

| Field | Value |
|---|---|
| Purpose | Make the boundary's authority claim precise enough to be decided, not believed |
| Status | **Contract-level clarification. Design only. V3.1 unchanged. Nothing authorized.** |
| Subject | row 563 (OQ-7), clauses O7-03, O7-04, O7-06, O7-10 |
| Recovery token | `TASK-03-OQ7-LIFETIME-CONTAINMENT-2026-11-22` |

Continues OQ-7's clause numbering (`O7-26`…). V3.1 remains the baseline; row 563 remains historical evidence.

---

## 1. Why "unforgeable capability" was imprecise

**O7-26 — Go cannot express the property the phrase implied.** A capability token can be stored, copied, retained, or reused after release, and no type-level construction in Go can prove that the *calling* code currently holds a lock. There is no goroutine-local state to check against, and `go vet`, the race detector and the compiler say nothing about lock ownership across a call boundary.

**O7-27 — So the claim must be split into four properties, each decidable:**

| Property | Question | Decided by |
|---|---|---|
| **non-constructible** | can code outside the boundary produce one? | unexported type + no constructor + no exported field |
| **non-retainable** | can a caller keep it past the call? | it must never be returned, stored, or accepted as a parameter |
| **non-reusable** | can it be used for a second operation? | it must not exist as a value at all |
| **containment** | can a caller reach the writer without one? | the writer is unexported and no API returns or accepts it |

A token satisfies the first property and struggles with the other three. **An API shape satisfies all four**, which is why the token should not be the proof.

---

## 2. The shape that is provable

**O7-28 — Containment by construction:** the raw history writer stays unexported inside the boundary package, and the package exposes **one** operation that performs the entire transition:

```
boundary.Mutate(chatID, command) (Result, error)
    acquire lock (fail → typed refusal)
    read + validate marker
    evaluate the command's precondition against (state, generation, settled_generation)
    refuse if the transition is unresolved
    perform the content mutation
    record the marker CAS
    release
```

The caller never receives a writer, a handle, an open file, a path, or a token. **Lifetime is trivially satisfied because no capability value exists to outlive the call.**

**O7-29 — The capability token may remain as an internal detail, never as the argument for trust.** If a token appears in the public API, properties 2–4 return as open questions, and O7-26 says they cannot be closed.

**O7-30 — `command` is a closed, typed set** — append, rewrite, clear, duplicate, delete — with no function-valued fields. A callback parameter would let caller code run inside the exclusive hold, where it could block, re-enter, panic, or retain what it was given; a function parameter would let a writer escape (§1, containment).

**O7-31 — No callbacks and no caller-supplied predicates inside the hold.** The precondition must be evaluated by the boundary from its own knowledge of the command type and the marker it just read. A caller-supplied "expected pair" reintroduces exactly the read-then-write race of V3-28, moved one layer up.

---

## 3. The five requirements, made decidable

**O7-32 — Construction.** Code outside the boundary package cannot obtain anything that authorises a history write. *Decided by:* the writer type, the lock handle, and the file handle are all unexported; no exported function returns them; no exported struct embeds them.

**O7-33 — Lifetime.** No authority object outlives its operation. *Decided by:* the public API returns only data (`Result`, `error`), never a capability, handle or closure that captures one.

**O7-34 — Write containment.** Every mutation of `history.ndjson` or `history.json` requires passing through the boundary. *Decided by:* the negative assertion of §4, plus the writer being unreachable from outside the package.

**O7-35 — Atomic composition.** Lock acquisition, marker comparison, content mutation and marker recording occur under **one continuous exclusive hold**, with no release-and-reacquire between them. *Decided by:* a single `acquire`/`defer release` frame spanning the whole operation — a structural fact about one function, not a discipline across several.

**O7-36 — Failure visibility.** Lock unavailability, precondition refusal, and a failed marker write all reach the caller as typed errors. *Decided by:* no discard site exists at any entry point. **Silent non-persistence is release-blocking**, not a degradation (O7-12).

**O7-37 — Partial success must never be reported as success.** If the content mutation succeeds and the marker CAS fails, the caller receives a failure and the conversation is left in a fail-closed row:

| Content | Marker | Reported | Row | State |
|---|---|---|---|---|
| written | CAS succeeded | success | 14–17 | Healthy |
| written | CAS refused or failed | **failure** | 12, 18 or 19 | Pending |
| written, torn | unchanged | failure | 7 | Damaged `truncated` |
| not written | unchanged | failure | unchanged | unchanged |
| settled-clear written | settlement CAS failed | **failure** | 13 | Pending `unattributed` |

No ordering of these operations may produce a return value that claims success while the pair and the content disagree.

---

## 4. The negative build assertion, and its honest limit

**O7-38 — The assertion must be negative, not a checklist.** "The known writers all use the boundary" establishes nothing about unknown ones. The assertable form is: *the set of code sites that open `history.ndjson` or `history.json` for writing is a subset of the boundary package* — a whole-program property, failing the build otherwise.

**O7-39 — Its completeness depends on path construction being centralised.** If a file name can be assembled at runtime — concatenation, configuration, a value passed in — a static assertion over literals and constructors cannot see the write. Completeness therefore requires that all conversation paths come from one constructor, which is itself a containment property. **Where that does not hold, the assertion is a heuristic and must be described as one** — not as proof. This is Unestablished today.

**O7-40 — The assertion runs as a build gate, not a lint.** A check that can be disabled, skipped or locally suppressed is not build-bound evidence.

---

## 5. What this changes in OQ-7's verdict

**O7-41 — O7-04 is restated:** authority is **containment by API shape** plus the negative build assertion; the capability type is at most an internal safeguard, and never the unit of proof.

**O7-42 — O7-22's five conditions now read:**

1. one boundary operation that owns lock, comparison, mutation and settlement in a single hold (O7-28, O7-35);
2. a *negative* build assertion that no other code opens a history representation for writing, valid only if path construction is centralised (O7-38, O7-39);
3. OS-lock exclusivity held continuously across the transition (O7-16);
4. settlement written only under that hold, never by observation (O7-20);
5. typed refusals surfaced at every entry point, with partial success never reported as success (O7-36, O7-37).

**O7-43 — Nothing about OQ-6 improves as a result.** The settlement field's trustworthiness is now *more precisely bounded*, not raised: without O7-42's items 1–2 the field records what an unverified writer asserted.

---

## 6. Still open on this branch

| # | Open question |
|---|---|
| OQ-7a | Whether conversation paths can be centralised enough for the negative assertion to be proof rather than heuristic (O7-39) |
| OQ-7b | Lock scope: one lock per conversation directory, or per document with a documented order (O7-11) |
| OQ-7c | Re-entrancy: whether any legitimate operation needs a nested boundary call, which the single-hold shape forbids |
| OQ-7d | Where non-history writers (`memories.json`, configs) belong — outside the claim, but writing into the same directory (O7-06) |

**O7-44 — Nothing here is authorized.** No boundary operation, no unexported writer, no build assertion, no token, no lock composition, no marker change, no migration, no backfill, no self-test, no deletion enablement. This clarification narrows what the eventual implementation must prove; it does not begin it.
