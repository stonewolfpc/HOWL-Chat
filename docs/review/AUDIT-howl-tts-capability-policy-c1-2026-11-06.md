# Adversarial Review — C-1 Capability Policy Decision (Policy B with Floor)

**Date:** 2026-11-06
**Reviewer:** Muse (independent adversarial pass)
**Subject:** `docs/review/DECISION-howl-tts-capability-policy-c1-2026-11-06.md`
**Subject SHA-256 (verified):** `c25a0d3fd38e0767a4b7bf278cb1d6617e2e53caf11053dabb71449578c4c8ad`
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`
**Contract SHA-256 (verified):** `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393`
**Prior reviews:** `AUDIT-howl-tts-helper-stage2-adversarial-2026-11-06.md`,
`AUDIT-howl-tts-helper-capability-correction-2026-11-06.md` (which raised C-1)

**Status:** INDEPENDENT REVIEW ARTIFACT. Read-only. No contract, decision report, code, test, fixture,
or accepted Stage 2 artifact was modified. No Policy B implementation, no caller wiring, no Stage 3,
no native provider, no commit, no push.

**Verdict: PASS WITH FINDINGS.** The policy choice is correct and better supported by the frozen text
than the report claims. The proposed *wording* cannot be adopted as drafted: two blocking gaps would
produce a contract that contradicts itself on arrival.

---

## 1. Verification of the submitted artifacts

| Item | Result |
|---|---|
| Decision report SHA-256 | **Exact match** (302 lines) |
| Architecture contract SHA-256 | **Exact match** (2,416 lines) |
| Cited §17.3.2 range (1244–1281) | Correct; five service IDs sit at **1266–1270** (report says 1265–1271 — bracket-inclusive, harmless) |
| Cited §17.6.1 range (1683–1699) | Correct; registry table 1685–1695, the contested prose at 1698–1699 |
| Cited §17.6.3 range (1719–1732) | Correct; `Capability false` row at 1725 |
| Cited §17.7.2 range (1762–1790) | Correct |
| Cited §17.9 rows | `protocol_incompatible` at **1881**, `unsupported_capability` at **1900** — both exact |
| Cited §17.11.4 steps (2069–2073) | **Exact** — steps 4, 5, 6 |
| Cited §17.15 O-3 (line 2325) | **Exact** |
| Cited `protocol.go:527-531` | **Exact** — the `protocol_violation` + `detail.missing` block |
| Remaining ranges (§17.3.4, 17.4.1, 17.4.2, 17.4.4, 17.4.6, 17.4.7, 17.7.1, 17.11.1–.3, §5.2, §5.5, §3.2, §9) | All correct within ±2 lines |

Line discipline is unusually good for a report of this size. The problems below are omissions of
*substance*, not citation errors.

---

## 2. Independent reconstruction of the ambiguity

The report describes the contract as **silent** on C-1 and then argues from silence. That understates
the situation in one direction and overstates it in another.

### 2.1 The contract is internally inconsistent, not merely silent

Three normative sentences cannot all be true at once:

1. §17.3.2, line 1277–1279: “A provider descriptor's `capabilities` **MAY include false capabilities**;
   `hello` describes what the binary can expose, while `catalog` describes what the selected
   provider/asset set actually supports.”
2. §17.5.1, line 1429: `capabilities` is “**Zero or more** `Capability` objects; no duplicate IDs”, and
   the `Capability` DTO is “exactly” `{id, version}` with “no other fields permitted” (1436–1439).
3. §17.6.1, line 1698–1699: “A capability with `Supported:false` is false for that provider; **it is not
   silently omitted**.”

(1) asserts the wire can carry a *false* capability. (2) removes every field that could express falseness
and permits an empty array. (3) forbids the only encoding that remains — absence. There is no reading
under which all three hold. `Supported` exists only in the application projection (§3.2, line 195), which
is what (3) must be describing, and (1) must be read as “may *name* capabilities the current asset set
does not satisfy” — but neither is stated.

This matters for the decision: the owner is not filling a gap, they are resolving a contradiction, and
the amendment must repair **§17.3.2's impossible sentence too**. The report's exposure of ambiguity
points 2 and 4 is correct but files this as “not stated” instead of “self-contradictory.”

### 2.2 The report under-claims the strength of its own case

Four clauses decide more of C-1 than the report credits:

- **§5.3, line 484.** The `hello` row of the discovery/loading table allocates exactly two failure
  states: `protocol_incompatible` and `helper_start_failed`. `protocol_violation` is absent. Under the
  current implementation, a helper missing one capability is rejected with a code the contract never
  assigns to `hello`.
- **§17.11.5, line 2089.** “No restart for `protocol_incompatible`, `helper_missing`, …” — the
  “no automatic retry” property the report asserts in §4.5 is *derivable*, not merely desirable, and only
  through this row.
- **§17.18.1 A1-D7.** “`protocol_violation` remains the precise internal classification” for the faults
  enumerated in A1-D3/D5 (publication after terminal, unknown ID, second terminal). The most recent
  normative amendment to the contract already narrows this code. That is the strongest available evidence
  that reusing it for a capability omission is wrong.
- **§4.6, line 397 and §17.9, line 1900.** `unsupported_capability` is defined as “Requested optional
  capability unavailable” — the word *optional* is in the frozen text twice.

Consequence: the report presents a genuinely strong argument as a balanced one. The owner should know
that Policy A is not merely under-specified — **its current error classification violates §5.3 today**,
independent of which policy is chosen.

### 2.3 What is genuinely undecided

Only one thing: whether the *set* `{catalog, load, unload, synthesize, cancel}` carries a mandatory
minimum, and if so which members. §17.4.1 establishes that these eight operations are the complete v1
vocabulary, and that every operation's precondition is “Handshake complete.” It never says a helper must
implement each one, and the only remedy the contract names for an unavailable capability is per-request
`unsupported_capability` (§17.6.3 line 1725). So the default reading favors B, and the floor is a
*justified addition*, not an interpretation.

---

## 3. Clause-by-clause assessment of the proposed wording (§6)

| § | Proposal | Assessment |
|---|---|---|
| 6.1 | Rewrite §17.6.1 prose: presence means implemented; require `load`+`synthesize`; missing required is `protocol_incompatible` with bounded `detail.missing`; optional set returns `unsupported_capability` and writes no frame; absence is the wire encoding; no no-op substitute | **Correct in substance, incomplete in placement.** It repairs the §17.6.1 contradiction but leaves §17.3.2 line 1277–1279 asserting “MAY include false capabilities” on a DTO with no such field. Must be repaired in the same amendment. |
| 6.2 | Add to §17.3.2 that five IDs are the vocabulary, not a requirement | **Necessary and correct.** This is the sentence that actually fixes the source of the ambiguity. But it must also retract “MAY include false capabilities,” which it currently leaves standing. |
| 6.3 | §17.4.4: absent `unload` returns `unsupported_capability`; process stop is not unload; `stopping`→`absent`/`crashed`; no-op unload non-conformant | **Correct and well-targeted.** It closes the loophole where a helper could satisfy Policy A with a lying `unload`. Two gaps: the transition it asserts has no row in §17.7.2 (B1), and it does not say what the parent does when it must change provider/identity with no `unload` available (B4). |
| 6.4 | §17.11.4 steps 4–6: record `cancel_requested`; send `cancel` only if advertised; skip otherwise; grace then `KillTree`/`Wait` | **Correct, and preserves precedence** — step 4's “first records `cancel_requested`” is retained, so §17.11.6 rule 1 still governs. Two gaps: §17.11.6 rule 4 independently mandates “before sending `cancel`” and is not amended (B3); step 5 becomes conditional, so the cooperative-stop path disappears and every cancellation of an absent-`cancel` helper is kill-bound (B8). |
| 6.5 | §17.9 `protocol_incompatible` row gains “a hello missing a required service capability”; `detail.missing` bounded | **Right classification, wrong blast radius.** The sentence “`protocol_violation` remains reserved for unmatched IDs, duplicate terminals, and publication after terminal” contradicts the five other handshake sites that still emit `protocol_violation` (`protocol.go:516, 521, 530, 535, 539`) — see B2. |
| 6.6 | Status: subset helper remains `Available`; list marks absent IDs unsupported | **Ambiguous and definitionally unsafe.** `TTSProviderStatus.Capabilities []TTSCapability` is *provider*-scoped (§3.2, line 187), while service capabilities are process-scoped; and §17.7.1 defines `Available` in terms of “required assets are known present,” knowledge that only `catalog` supplies. See B5, B6. |
| 6.7 | No `supported` field on the wire DTO | **Correct and non-negotiable.** It also matches §17.5.1's “no other fields are permitted.” |

---

## 4. Findings

### B1 — BLOCKING (amendment completeness): the state machine cannot reach the state the amendment mandates

§6.5 assigns `incompatible` status and `ErrCodeModelLoadFailed`, and §17.11.5's no-restart rule keys off
`protocol_incompatible`. But §17.7.2 contains exactly one row producing `incompatible`:

> `handshaking` | **no common major** | `incompatible` | No automatic retry

A hello that is well-framed, version-compatible, and missing `howl.tts.load` has no legal transition. It
is not “no common major”; it is not `crashed`; it cannot remain `handshaking`. The amendment would make a
normative state unreachable — precisely the contradictory-rules outcome this review is meant to prevent.

**Required:** add a transition row (`handshaking` | hello missing a required service capability |
`incompatible` | bounded `detail.missing`; no automatic retry), and widen the §17.11.5 sentence to name
it, so restart suppression is authorized rather than merely implied.

**B1b (moderate, pre-existing, load-bearing under B).** §17.8 freezes “Idle close 120 s — Stop an idle
helper instance,” and §5.5 says an idle instance may be stopped, but §17.7.2 has no idle-close row
(`ready`/`degraded`/`unloaded` → `stopping` occurs only on “`shutdown` accepted”). Under B this becomes
the *primary* resource-release path for a helper without `unload`, so the gap stops being cosmetic. The
amendment should add the row or state explicitly that idle close is implemented as `shutdown`.

### B2 — BLOCKING (self-contradiction): §6.5 reserves `protocol_violation` while five handshake sites keep emitting it

The proposed sentence narrows `protocol_violation` to correlation and publication faults. In the current
parent, `validateHello` still returns `protocol_violation` for:

| Site | Condition |
|---|---|
| `protocol.go:516` | duplicate provider ID |
| `protocol.go:521`, `539` | invalid or duplicate extension descriptor |
| `protocol.go:518` | unregistered or wrong-version capability (`validateCapabilities`) |
| `protocol.go:530-531` | missing mandatory capability (the one being reclassified) |

Adopting §6.5 verbatim leaves four conditions in direct violation of the amended registry. The contract
also never states which code an **unregistered capability ID or wrong version** in `hello` should
produce — §17.6.1 says only that versions are checked, and §17.6.2's `unsupported_extension` covers
extensions, not core capabilities.

**Required:** one sentence classifying the whole hello-validation surface. Recommendation: an
unregistered ID, wrong version, duplicate descriptor, or malformed handshake descriptor is
`protocol_incompatible` (the only code §5.3 allocates to `hello`; wire origin, non-retryable, no restart),
distinct from `invalid_frame` for unparseable bytes (§17.2.7–§17.2.8). Without this, the amendment
trades one contradiction for another.

### B3 — MODERATE: §17.11.6 rule 4 still mandates an unconditional `cancel`

§17.11.6 rule 4 reads: “the adapter marks `cancel_requested`/`timeout` **before sending `cancel`**; if the
helper is alive, the supervisor kills and observes process exit before deleting the exact final/`.part`
paths.” §6.4 amends §17.11.4 only. Two normative rules would then disagree about whether a `cancel` frame
is always sent.

**Required:** amend §17.11.6 rule 4 to “before sending `cancel`, if the helper advertised
`howl.tts.cancel`,” and state that cancellation precedence and cleanup ordering are unchanged.

### B4 — MODERATE: no parent rule for provider/identity change when `unload` is absent

§17.8 freezes “One loaded provider per process — `provider_already_loaded` until unload.” §17.4.4 says a
`load` for a different provider/model returns `provider_already_loaded` and “the caller must `unload`
first.” If `unload` is absent, that instruction is unsatisfiable: the parent would deterministically fail
every identity change against a live helper.

**Required:** state that when `unload` is absent, a model/provider change is implemented as supervised
termination plus a fresh instance for the new key (§5.5 keys instances by
`provider_id + engine + model_identity`), and that the parent MUST NOT send a differing-identity `load` to
a live helper. §4.2 gestures at “a later request starts a new process and reloads” but never turns it into
a rule or names the trap.

### B5 — MODERATE: `Available` is defined in terms of `catalog` knowledge

§17.7.1: “`Available` means the helper is present, protocol-compatible, the provider is compiled in, and
**the required assets are known present**.” Asset presence is only knowable via `catalog`, which §5.3
permits without a load. With `catalog` optional, a catalog-less helper cannot be classified by that
definition at all — it is neither provably available nor unavailable before `load`.

**Required:** redefine availability for catalog-less helpers, e.g. “`Available` requires handshake
compatibility and a compiled provider; asset presence is established at `load`, and a failed load moves
to `unavailable` per the existing row.” Otherwise §6.6's “remains `Available`” is an assertion the state
vocabulary does not support.

### B6 — MODERATE: the status projection conflates process-scoped and provider-scoped capabilities

`TTSCapability` lives in `TTSProviderStatus.Capabilities` (§3.2, line 187) — provider-scoped. Service
capabilities describe the *helper process and its protocol surface*. §6.6 says “The status capability list
marks absent IDs as unsupported” without saying which list, inviting either a mixed array (erasing the
service/provider split §17.5.1 enforces on the wire) or silent omission of service capability state.

**Required:** name the surface. Recommendation: service capability state belongs in a helper/process
section of the status (or an explicitly separate field); provider capabilities stay in
`TTSProviderStatus.Capabilities`. Note in passing that §3.2's `State` comment already lists nine of the
thirteen §17.7.1 states — if the amendment touches status semantics, reconcile that too.

### B7 — MINOR: §4.6 line 389 is a second, earlier definition of `protocol_incompatible`

“`protocol_incompatible` | No common protocol major | `ErrCodeModelLoadFailed`.” §6.5 amends §17.9 only.
Both tables are normative-looking; the earlier one would remain narrower. Add the missing clause there
as well, or mark §4.6 as superseded by §17.9.

### B8 — MINOR: the floor's honest-provider argument has a hole for lazy runtimes

§4.1 argues B avoids pressuring providers into “no-op implementations.” That is right for
`unload`/`cancel`/`catalog`, but not for `load`. §17.4.4 requires `load` to validate assets *and* that
“the runtime can initialize,” and §17.11.3 step 8 says “On `load`, validate assets and initialize
provider/runtime; readiness becomes true only after the load result,” while §5.3 forbids any load during
`synthesize`. A runtime that can only initialize lazily therefore has no conformant option under either
policy; it must implement an honest eager `load` or be excluded.

This does not weaken the recommendation — mandatory explicit `load` is genuinely necessary, since without
it §17.11.4 step 2 (“ensures handshake and load”) cannot succeed and §17.4.1's synthesis precondition is
unreachable. But the report should state the limitation instead of implying the floor is free: **B
preserves neutrality for `unload`, `cancel`, and `catalog`, not for initialization timing.**

### B9 — MINOR: absent `cancel` makes every cancellation kill-bound

With §6.4's skip, step 5 (“the helper reader acknowledges and signals the provider context”) no longer
applies, so no cooperative stop can occur. Every cancellation of such a helper waits the full grace and
ends in `KillTree`. That is correct and bounded, but it is a latency floor and an orderly-shutdown
regression the owner should adopt knowingly. One sentence in §4.3.

### B10 — INFORMATIONAL: the amendment mechanism contradicts the contract's own convention

§17.18 exists as an **additive** clarification that “resolves the two open review findings without
renumbering or rewriting §§17.1–17.17,” and §17.18.3 states “The original §17 text and earlier proposal
history remain in place.” §6.1 proposes to *replace* lines 1697–1699 in place. Adoption should follow the
established A1 pattern — a new §17.19 with decision entries — so the corpus keeps one amendment
convention and the audit trail survives.

### B11 — INFORMATIONAL: four missed citations that favor the recommendation

§5.3 line 484; §17.11.5 line 2089; §17.18.1 A1-D7; §4.6 line 397. See §2.2. None changes the outcome;
all strengthen it, and their absence makes Policy A look more defensible than the text supports.

### B12 — INFORMATIONAL (pre-existing, out of scope, referred back)

§17.10.3 mandates publication “with an **atomic rename** to the final filename.” The fake helper publishes
by hard link (my C-3 in the correction review). The no-overwrite effect is preserved and arguably
stronger, but the corpus already diverges from its own wording. Not caused by this proposal; objective 9
prompted the check. Separately, §17.6.2 line 415 says extensions are rejected with
`unsupported_capability` while the same section's body says `unsupported_extension` — do not propagate
that into §6.5's registry edit.

---

## 5. Answers to the review objectives

1. **Clause recheck** — all cited ranges accurate; §2 above corrects the framing from “silence” to
   “contradiction” and adds four clauses the report omits.
2. **Internal consistency with the state machine and provider interface** — the policy is consistent;
   the wording is not, until B1/B1b/B4/B5/B6 are applied. The provider interface (§5.2) needs no change:
   `Unload` and `Close` are already distinct, and the helper owns process shutdown.
3. **Is mandatory `load` necessary?** Yes, and more strongly than argued: §17.4.1, §5.3, §17.4.4,
   §17.7.2, §17.11.3, and §17.11.4 step 2 all assume an explicit load. Preserving a lazy provider would
   require amending five sections and the readiness definition — a far larger change than C-1. The
   tradeoff the report must disclose: initialization timing is *not* neutralized by the floor (B8).
4. **`unload` absent** — process termination is categorically not an unload (§17.4.4 forbids it as a
   forced unload), so omission is the only honest advertisement; resource release is guaranteed by job
   termination; status must go `stopping`→`absent`/`crashed` per §17.7.2 and must never report
   `unloaded`; provider replacement then requires a new process, not a reload (B4).
5. **`cancel` absent** — the frame must be skipped (§17.4.1 gives `cancel` a live synthesis precondition
   and the helper would answer `unsupported_capability`, a terminal the cancellation path would have to
   absorb); boundedness is unaffected because §17.11.4 step 6's grace/`KillTree`/`Wait` sequence never
   depended on the acknowledgement; precedence is preserved provided step 4's `cancel_requested`
   recording stays first, which §6.4 keeps. Cost: kill-bound cancellations (B9).
6. **Error classification** — `protocol_incompatible` for missing required is right and is in fact
   *required* by §5.3; `unsupported_capability` for missing optional is right per §17.6.3/§4.6/§17.9;
   `protocol_violation` must be narrowed and the remaining handshake sites reclassified together (B2).
7. **Sections that must be amended together** — §17.3.2 (1277–1279), §17.4.4, §17.6.1 (1697–1699),
   §17.7.1 (`Available`), §17.7.2 (transition rows), §17.9 (1881), §17.11.4, §17.11.5 (2089),
   §17.11.6 rule 4, §4.6 (389), §5.3 (484). The report covers six of eleven.
8. **Future providers / compatibility / status / security** — §9's independent-update rule enumerates
   only “protocol major, provider ID, and extension contracts,” so Policy A adds an unlisted
   compatibility requirement and contradicts §9; that argument should be made explicitly. No security
   consequence either way (§17.11 ownership, path validation, and the two-process cap are orthogonal).
9. **Frozen guarantees touched** — none by intent. Ownership, correlation, acceptance barrier, and
   publication are untouched; §6.4 does not weaken step 6. The only cross-reference is B12, which is
   pre-existing.

---

## 6. Required corrections to the proposal

1. Add the `handshaking` → `incompatible` transition and the §17.11.5 no-restart clause (B1); add or
   explicitly deny an idle-close transition row (B1b).
2. Classify the entire hello-validation surface in one sentence and reclassify the four remaining
   `protocol_violation` sites, or §6.5 contradicts the code it governs (B2).
3. Amend §17.11.6 rule 4 alongside §17.11.4 (B3).
4. State the stop-then-start rule for identity change without `unload` (B4).
5. Redefine `Available` for catalog-less helpers (B5).
6. Name the status surface that carries service-capability state (B6).
7. Add the missing clause to §4.6 line 389 (B7).
8. State the initialization-timing limitation of the floor (B8) and the kill-bound cancellation cost (B9).
9. Adopt as an additive §17.19, not an in-place rewrite (B10); cite §5.3, §17.11.5, §17.18 A1-D7, and
   §4.6 (B11).

---

## 7. Verdict

**PASS WITH FINDINGS.** Policy B with the `load`+`synthesize` floor is the correct decision: it follows
the contract's own per-request remedy, it is supported by §5.3 and §17.18 more strongly than the report
argues, and the floor is justified by the operation table rather than by taste. The owner's selection is
sound and should be carried forward.

The **wording** must not be adopted as drafted. B1 and B2 are blocking: one makes a mandated state
unreachable, the other reserves `protocol_violation` while five handshake sites keep emitting it. Both are
mechanical to fix and neither reopens the policy question. B3–B6 are moderate and should be settled in
the same amendment rather than discovered during Stage 3 wiring.

A recommendation is not permission to implement. The current implementation remains Policy A and must not
be changed until the owner approves an architecture amendment.

---

## 8. Questions requiring owner decisions

1. **Floor membership.** Confirm `load`+`synthesize` (report's recommendation) rather than `synthesize`
   alone. §17.11.4 step 2 makes `load` operationally necessary, but it is still a choice.
2. **Handshake error code.** `protocol_incompatible` for all hello-validation failures (recommended, per
   §5.3), or a new registry code? Choosing a new code means touching §17.9, §5.3, §4.6, and §17.11.5.
3. **`Available` for catalog-less helpers.** Accept “handshake-compatible + compiled provider, assets
   proven at load” as the new definition, or make `catalog` part of the floor after all?
4. **Status surface.** Where does service-capability state live — a helper-level status section, or inside
   the provider-scoped `Capabilities` list?
5. **Idle close under B.** Is idle close a `shutdown` (graceful, §17.4.7) or an independent termination
   path needing its own §17.7.2 row?
6. **Cancellation cost.** Accept that a `cancel`-less helper is always kill-bound after the grace, or
   require `cancel` in the floor as well?
7. **Amendment mechanism.** Additive §17.19 in the A1 style, or an in-place rewrite of §17.6.1?

---

## 9. Boundaries of this review

**Performed:** hash verification of both artifacts; full read of the decision report; direct reading of
§§3.2, 4.6, 5.1–5.5, 9, 17.3.2–17.3.4, 17.4.1–17.4.7, 17.5.1, 17.6.1–17.6.3, 17.7.1–17.7.2, 17.8,
17.9, 17.10.3, 17.11.3–17.11.6, 17.15, 17.18; spot verification of the cited code sites in
`internal/audio/tts/helper/protocol.go`.

**Not performed:** no contract or report edit, no code, test, fixture, or golden change, no Policy B
implementation, no `Unload`/`Cancel` caller wiring, no Stage 3, no native provider, no build, no commit,
no push.
