# Adversarial Review — C-1 Amendment Revision 2 (§17.19 draft)

**Date:** 2026-11-07
**Reviewer:** Mimo (independent adversarial pass)
**Subject:** `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`
**Subject SHA-256 (verified):** `e7088adf1c261006d12a49dbfcbbb1ce26a17703699fb977afd6776b57ae3390` (229 lines)
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`
**Contract SHA-256 (verified):** `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` (2,416 lines)
**Prior audit reviewed against:** `AUDIT-howl-tts-capability-policy-c1-2026-11-06.md` (`258fd1c1…`)

**Status:** INDEPENDENT REVIEW ARTIFACT — READ-ONLY. This review confirms its own read-only boundary by
re-hashing the contract at the end (§8). No drafting, no implementation.

**Verdict: PASS WITH FINDINGS — conditional.** The policy content is correct and the five owner
conditions are met in substance. The draft is **not adoptable as written**: two blocking defects change
what the owner would actually be approving (R1, R2). No policy question needs reopening.

---

## 1. Verification performed

| Check | Result |
|---|---|
| Draft SHA-256 | Exact match, `e7088adf…`, 229 lines as stated |
| Contract unchanged | **Exact match** — `4ff6b72e…`, identical to the C-1 review |
| "Frozen architecture unchanged" claim | Verified true, not merely asserted |
| Line references checked | **34 of 36 exact**; two off by two (§17.8, §17.11.6 rule 4) |
| Every B1–B12 disposition traced into §1's normative text | 12 of 14 executed; B8 and half of B11 asserted only in §3 |
| Finding-count arithmetic (§3.1) | Correct: 2 + 10 + 2 = 14 dispositions over 12 numbered findings |
| Impact rows vs supersession map, clause by clause | **They do not agree** — 5 rows instruct in-place rewrites the method forbids (R1) |
| Contract re-hashed after review | Unchanged (§8) |

Independently confirmed accurate: §17.1 item 10 at line 1069; §17.3.2 at 1277–1279; §17.4.1 at
1323–1335; §17.4.4 "caller must `unload` first" at 1372–1374; §17.6.1 at 1697–1699; §17.6.3 at
1719–1732; §17.7.1 `Available` at 1757–1760; §17.7.4 `ProcessState=alive, State=incompatible` at **1819**;
§17.9 `protocol_violation` **1878**, `protocol_incompatible` **1881**, `unsupported_capability` **1900**;
§17.11.4 steps 4–6 at 2069–2073; §17.11.5 at 2089; §17.15 O-3 at 2325; §17.17 at 2354–2368; §4.6 at 389;
§4.7 at 415; §5.3 `hello` row at 484; §3.2 `State` comment listing 9 of 13 states.

Two additional verifications worth recording:

- **§17.4.1 has no capability column** (columns are Operation / `params` / Result / Preconditions). A2-D1's
  statement that `health` and `shutdown` are ungated control operations is therefore a genuine
  clarification, not a contradiction. Correct.
- **§3.2 is headed "Additive provider inspector … Proposed data shapes"** (line 165). Adding
  `Loaded`/`ServiceCapabilities`/`ServiceCapabilitiesKnown` is legitimate; it does not violate a frozen
  DTO. A2-D6 and impact rows #2/#25 are sound on this point.

---

## 2. Blocking findings

### R1 — BLOCKING: the amendment method is self-contradictory, and the impact list would rewrite the text it promises to preserve

A2-D11 states: "This section is additive and uses the A1 convention. It preserves the original text and
supersedes only the clauses listed in §17.19.3." §1's preamble repeats it: "On approval it is appended
after §17.18.3. The original §§17.1–17.18 remain in place."

But §2's required actions mutate the original text in five rows:

| Row | Clause | Required action as written |
|---:|---|---|
| 4 | §4.6 line 389 | "**Broaden** to include semantic hello validation" |
| 7 | §17.3.2 lines 1277–1279 | "**Replace** with the A2-D2 reading" |
| 12 | §17.6.1 lines 1697–1699 | "**Replace** with floor, optional set, absence semantics…" |
| 15 | §17.7.1 lines 1736–1760 | "**Redefine** `Available`…" |
| 19 | §17.9 lines 1858–1912 | "**Broaden** `protocol_incompatible`…" |

The same five clauses appear in §17.19.3 as *superseded by* decisions — a different mechanism. So the
owner is being asked to approve two incompatible procedures for the same clauses. Either §17.19 is
self-sufficient (the original sentences stand and stop controlling, and rows 4/7/12/15/19 are
descriptions of effect, not edits) or the original text is rewritten (and A2-D11's preservation claim,
on which B10's entire disposition rests, is false).

This is not pedantry. B10 was raised precisely because an in-place rewrite discards the audit trail the
A1 amendment deliberately kept. Adopting this draft under the "additive" label while executing five
rewrites would silently reintroduce the register the review rejected — with the label now asserting the
opposite.

**Required:** restate rows 4, 7, 12, 15, 19 as "addressed by §17.19 A2-Dx; original sentence stands and
is superseded for the listed matter" and reserve in-place verbs for the genuinely additive insertions
(new fields, new rows, item 11, checklist row). Then state explicitly which of the two mechanisms
applies to each row, and make §1's text carry every normative statement those clauses need — because
under the additive reading, §17.19 alone must be sufficient.

### R2 — BLOCKING: the "three separate facts" cannot be represented in §17.7.4, and the impact row that should fix it says "confirm"

The owner's coherence condition is: `Available=true` with `State=unavailable` is acceptable *provided*
`LastError` identifies the failed model instance. The draft encodes the first two facts and half the
third, but the health table has no row for the resulting combination:

- §17.7.4 line 1818 asserts `ProcessState=absent`, `State=unavailable`, `Ready=false` — "No helper or
  dependency; TTS unavailable." Under A2-D5 an asset-specific failure leaves `Available=true` and the
  process may be **alive**.
- §17.7.4 line 1819 is `ProcessState=alive`, `State=incompatible` — not this case.
- §17.7.2 line 1775 explicitly permits `loading` → `unavailable` with "**Process may remain alive until
  stopped**."

So the combination `ProcessState=alive`, `State=unavailable`, `Available=true`, `Loaded=false`,
`Ready=false` is legal under A2-D5 and has no row in §17.7.4. Impact row #17's required action is
"**Confirm** and cite A2-D4" — a no-op for the one table that must gain a row. Confirming §17.7.4
unchanged would leave the contract asserting `ProcessState=absent` for a state the amendment allows to
be alive.

Worse, this exposes a **pre-existing contract conflict** the amendment now depends on: §17.7.1 line 1755
lists `unavailable` among the states where `ProcessAlive` is false, while §17.7.2 line 1775 says the
process may be alive in that state. A2-D5 makes that contradiction load-bearing.

**Required:** add a §17.7.4 row (live process, `unavailable`, available-but-not-loaded, `LastError` naming
the failed model identity), and change impact row #17 from "Confirm" to "Add row". Also state that
`Available=true` + `unavailable` applies per failed attempt, not per provider, or explain how a
provider-level `State` and `LastError` remain meaningful when a single identity fails (see R3).

---

## 3. Moderate findings

### R3 — `Available=true` invites a retry the state machine cannot perform

§17.7.2 line 1790 offers exactly one exit from `unavailable`: "helper/assets/config **updated** → `absent`,
explicit new attempt only." Under A2-D5, `Available=true` now means the configured provider can be used
in principle — which is precisely an invitation to retry. If an asset is restored without a
configuration change, there is no transition: `unavailable` → (`nothing`) on an explicit request.

**Required:** either add an explicit-request transition from `unavailable` (subject to the restart
budget and the "no automatic retry" rule), or state that `Available=true` does **not** authorize a retry
and that only the update event clears the state. Leaving both readings open is the same
unreachable-transition defect as B1, inverted.

### R4 — The `ProcessAlive` supersession is narrower than the decision it introduces

The supersession map names only "§17.7.1, `ProcessAlive` sentence listing **`incompatible`** as false."
Line 1755 lists `incompatible` **and `unavailable`** (and `backoff`) as false. A2-D4's rule —
"`ProcessAlive` follows `ProcessState`, not the aggregate provider state" — logically supersedes the
whole sentence, and must, because §17.7.2 line 1775 keeps an `unavailable` process alive.

**Required:** widen the entry to the full sentence and say `ProcessAlive` is derived from
`TTSProcessState` alone.

### R5 — A2-D8 adds a transition row that already exists

A2-D8 specifies "`ready`/`degraded`/`unloaded` + idle deadline → `stopping` via `shutdown`", and impact
row #16 instructs "Add **both** rows." But §17.7.2 line 1782 already reads:

> `ready`/`degraded`/`unloaded` | `shutdown` accepted | `stopping` | Stop accepting work

Once idle close is defined as *emitting `shutdown`*, line 1782 covers it exactly. Adding a second row
with the same From/Event/To would duplicate a frozen transition and create two sources for one state
change.

This also corrects **my own B1b**, which said §17.7.2 has "no idle-close row." Precisely: there is no
idle-close *trigger*; the state change is already frozen. The missing piece is one definitional sentence,
not a second row.

**Required:** replace "add both rows" with "add the `handshaking` → `incompatible` row (B1); state that
idle close emits `shutdown`, so line 1782 governs the transition."

### R6 — "Rejected shutdown → `crashed`" contradicts the `shutdown` error code, and the sequence drops a mandated step

A2-D8: "a rejected shutdown request or a transport failure while sending it → failed shutdown, terminate
as needed, observe `Wait()`, and `crashed`." But §17.9 line 1906 defines the `shutdown` error code as
"Helper is stopping; new request may succeed later", retryable. The only plausible reason a `shutdown`
request is refused is that the helper is *already* stopping — in which case the correct outcome is to
await exit and classify by exit code (`absent` for 0, `crashed` otherwise), not to declare `crashed`
immediately. §17.4.1 also gives `shutdown` the precondition "any provider state" with result
`{"accepted":true}`, so a semantic refusal is not an expected shape.

Separately, A2-D7's replacement sequence (steps 3–5: shutdown grace, `KillTree`, `Wait()`) omits the
**close stdin** step that §17.11.5 line 2091 mandates in the graceful-shutdown path.

**Required:** distinguish a `shutdown` *error response* (await exit, classify by exit code) from a
transport write failure (kill and classify), and add "close stdin" to A2-D7's sequence.

### R7 — A2-D7 authorizes pre-empting an active synthesis, which the serialization lock forbids

A2-D7 step (1): "if synthesis is active, cancel under A2-D9 and wait for the terminal response or grace."
But §17.11.4 step 1 makes the adapter acquire the per-instance serialization lock before anything else,
and §17.7.2/§17.8 freeze "one active synthesis; no queue; `busy` for a second work request." An identity
change is itself a request that waits on that lock, so by the time the replacement sequence begins, no
synthesis can be active. Detecting one would mean the lock was bypassed.

As written, the step reads as permission to cancel a caller's in-flight synthesis in order to change
model. That is a new pre-emption semantic, not the existing lifecycle A2-D7 claims to reuse.

**Required:** state that the identity change executes under the per-instance lock, so step (1) is a
defensive assertion that no synthesis is active — not an authorization to cancel one — or, if
pre-emption is genuinely intended, say so explicitly and amend §17.7.2/§17.8 instead of claiming no new
mechanism.

### R8 — `ServiceCapabilitiesKnown` invalidation is incomplete on failure paths

A2-D6 invalidates the snapshot "before a validated compatible handshake" and "after an incompatible
handshake or process replacement". The owner asked whether it is invalidated on **every**
process-replacement and failure path. It is not: `crashed`, unexpected exit, `stopping`, `unavailable`,
and `backoff` are not named, so the snapshot of a process that no longer exists can remain marked known
and authoritative.

The rule also lacks a stated invariant. "Invalidate on replacement" and "retain across a crash" cannot
both be right, and neither follows from the other.

**Required:** replace the event list with one invariant — the snapshot is authoritative **iff the current
instance has completed a compatible handshake** — and state that crash/exit/stop clear it because that
condition ceases to hold.

---

## 4. Minor findings

**R9 — wrong bound citation.** A2-D3 says `detail.missing` "is bounded by §17.2.4". §17.2.4 is the error
*envelope* (lines 1126–1141). The bound is §17.2.5, line 1160: "Error detail | 16 entries; key 64; value
256 bytes." Five IDs fit, so the rule is sound; the reference is not.

**R10 — termination versus the exit observer.** A2-D4 says the supervisor terminates an `incompatible`
helper and observes `Wait()`. §17.7.2 line 1786 maps "any live state | unexpected `Wait()` exit →
`crashed`." The incompatible path's exit is *expected*, so it must not be reclassified. A2-D4 does state
"The aggregate `State` remains `incompatible`", but does not say the exit observer is intent-aware.
Add one sentence, or an implementation can race the two rules.

**R11 — B11's disposition over-claims, and the strongest support is uncited.** §3 says B11's resolution is
to "cite all four in the amendment's rationale and supersession map." In fact §5.3 line 484 and
§17.11.5 line 2089 appear in §17.19.3, while **§17.18.1 A1-D7** and **§4.6 line 397** appear only inside
the finding-description column, never in §1's normative text or the map. Additionally, the draft never
cites §17.9 lines **1867–1872**, which already state that `protocol_violation` is the classification for
"an unmatched ID, duplicate terminal, or publication after terminal" — the single strongest in-contract
support for A2-D3, stronger than the table row at 1878 the draft does cite.

**R12 — B8's disposition is asserted but not executed.** §3 claims A2-D1/A2-D5 state that "a lazy-only
runtime is non-conformant" and that neutrality excludes initialization timing. Neither sentence exists in
§1. §17.4.4's "the runtime can initialize" and §17.11.3 step 8 already imply it, so the substance is
defensible — but the traceability table records a resolution the normative text does not contain.

**R13 — two line references off by two, and one partial overlap with a deferred item.**

- §17.8 idle close is at line **1846** (1844 is "Cancel grace 2 s"). Impact row #18 and the supersession
  map both say 1844.
- §17.11.6 **rule 4** is at lines **2127–2130** (2125–2126 is rule 3). The supersession map says
  2125–2128.
- A2-D3 classifies a "wrong … extension version" at handshake as `protocol_incompatible`, while A2-D11
  defers the extension error-wording conflict to B12b. These are different surfaces (§17.6.2 line 1708
  governs *later requests* using an unadvertised extension; A2-D3 governs the hello descriptor), but the
  draft should say so rather than leave a deferred item apparently half-decided.

**R14 — `LastError` must name the model identity.** The owner conditioned coherence on `LastError`
identifying "the failed model instance rather than implying the provider itself disappeared." A2-D5 and
§17.19.4 item 5 say only "identifies the asset-specific failure." Add the identity requirement
explicitly, including `provider_id` + `model_identity`, consistent with §5.5's instance key.

---

## 5. Direct answers to the five requested checks

**1. Is `Available=true` with `State=unavailable` consistent across every status definition and
transition?** *Mostly, no.* The §17.7.1 mapping is superseded cleanly and §3.2's 9-of-13 state comment is
correctly fixed. But §17.7.4 cannot represent the state (R2); §17.7.1's `ProcessAlive` sentence still
contradicts §17.7.2 line 1775 (R4); `unavailable` has no retry transition (R3); and `LastError` lacks the
identity requirement (R14).

**2. Is `ServiceCapabilitiesKnown` invalidated on every process-replacement and failure path?**
*No* (R8). `crashed`, unexpected exit, `stopping`, `unavailable`, and `backoff` are unaddressed.

**3. Do shutdown failure, `KillTree`, `Wait()`, and output cleanup preserve the terminal-publication
barrier?** *Yes, with two defects.* The barrier itself is intact: A2-D9 keeps §17.11.6 rules 1–3 and 5–7,
and step 6's grace/`KillTree`/`Wait()` ordering is unchanged, so no success frame can revive a cancelled
request. But "rejected shutdown → `crashed`" mis-maps the `shutdown` error code (R6), and A2-D7 drops the
mandated "close stdin" (R6). R10 is a third, lower-risk ordering concern.

**4. Does identity replacement handle active synthesis and cancellation precedence correctly?**
*Precedence yes, replacement mechanism no.* Cancellation precedence is preserved because A2-D9 records
`cancel_requested` before anything else, so rules 1–2 still decide in favour of cancellation. But A2-D7's
"if synthesis is active, cancel" is inconsistent with the per-instance serialization lock and the frozen
no-queue rule (R7).

**5. Do the 29 impact rows, the supersession map, and B1–B12 dispositions agree?** *No.* Rows 4/7/12/15/19
prescribe rewrites the map calls supersessions (R1); row #17 says "confirm" where a row must be added
(R2); row #16 adds a duplicate transition (R5); B8 and part of B11 are claimed but not written (R11, R12).

---

## 6. What the draft gets right

Worth stating plainly, because the corrections above are mechanical and the substance is sound:

- **Policy B with the `load`+`synthesize` floor** is implemented without reopening it, and A2-D1's
  separation of service capabilities from `health`/`shutdown` is correct against §17.4.1's column set.
- **A2-D2** correctly identifies that §17.3.2's "MAY include false capabilities" is unencodable by a
  `{id, version}` DTO, and resolves it in the only coherent direction — absence is the wire encoding,
  `Supported:false` is application-only, catalog carries asset availability.
- **A2-D3** is the right repair for B2 and matches §5.3 line 484's allocation and §17.9 lines 1867–1872.
  R9/R13 are citation and scoping fixes, not policy errors.
- **A2-D6**'s separation of `ServiceCapabilities` from the provider-scoped array is exactly the fix B6
  needed, and `ServiceCapabilitiesKnown` is a genuinely better answer than five false values.
- **B12a/B12b are deferred explicitly rather than folded in**, and the count check in §3.1 is accurate.
- The document's own scope discipline is intact: nothing in it authorizes implementation.

---

## 7. Verdict and approval condition

**PASS WITH FINDINGS — conditional.** The amendment's substance is adoptable; its text is not. The owner
should not accept revision 2 from this review alone: **R1 and R2 are blocking** because both change what
approval means — R1 by making the amendment method ambiguous, R2 by asserting that a table requiring a
new row needs no change. R3–R8 are moderate and should be settled in the same pass so they are not
rediscovered during Stage 3 wiring. R9–R14 are mechanical.

Required before a third review or approval:

1. Resolve the additive-versus-edit contradiction and restate the five affected impact rows (R1).
2. Add the §17.7.4 row, fix impact row #17, and widen the `ProcessAlive` supersession (R2, R4).
3. Supply the `unavailable` retry rule (R3).
4. De-duplicate the idle-close transition and correct the rejected-shutdown mapping and the missing
   close-stdin step (R5, R6).
5. Scope A2-D7 to the serialization lock, or declare pre-emption (R7).
6. Replace A2-D6's event list with the live-handshake invariant (R8).
7. Correct §17.2.5, §17.8 line 1846, §17.11.6 lines 2127–2130; add `LastError` identity and the §17.9
   1867–1872 citation; execute B8 and complete B11 (R9–R14).

---

## 8. Boundaries of this review

**Performed:** hash verification of the draft and the contract; full read of all 229 lines; direct reading
of every cited contract clause plus §17.2.5, §17.6.2, §17.7.3, §17.9 lines 1867–1872, §17.11.5 exit codes,
§17.12, §17.17, and §17.15; independent check of all 34 line references; clause-by-clause comparison of the
29 impact rows against the supersession map and the 14 dispositions; verification that the draft's
"frozen architecture unchanged" claim is true.

**Not performed:** no edit to the draft, the contract, the decision report, or the prior audit; no code,
test, fixture, golden file, or Stage 2 artifact change; no implementation of Policy B; no `Unload`/`Cancel`
wiring; no Stage 3; no build, commit, or push. The contract was re-hashed after this review and is
unchanged.
