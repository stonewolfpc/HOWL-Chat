# Adversarial Review — C-1 Amendment Revision 3 (§17.19 draft)

**Date:** 2026-11-07
**Reviewer:** Mimo (independent adversarial pass)
**Subject:** `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`
**Status:** INDEPENDENT REVIEW ARTIFACT — READ-ONLY; the amendment remains unadopted

**Verdict: PASS WITH FINDINGS — conditional.** Revision 3 resolves both previous blockers and twelve of
fourteen dispositions are now executed in the normative text. One blocking defect remains: the status
surface A2-D6 mandates cannot be expressed by the type the amendment tells implementers to use for it.
Four moderate gaps follow from that and from A2-D12's dependency on an unstated helper transition. The
policy and the architecture both stand; §17.19 should not be adopted until F1 is corrected.

---

## 1. Artifact hashes (verified on disk)

| Artifact | SHA-256 | State |
|---|---|---|
| Amendment draft, revision 3 | `a13fb198696622a1d5663c774a1e81935b78cf755e997f67bbbf9ebc10161607` | **matches expected**; 246 lines, 43,144 bytes |
| Frozen architecture | `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` | unchanged |
| C-1 decision report | `c25a0d3fd38e0767a4b7bf278cb1d6617e2e53caf11053dabb71449578c4c8ad` | unchanged |
| Prior audit (rev 2) | `db6beb0a604cb65af4fe82a97d29db42395bc1312eabce051398b48ef4163bb3` | unchanged |
| Original C-1 audit | `258fd1c1c981c2000ddcca4e87faea1e667d2b24114b68a48c8dea021d9c870e` | unchanged |

The contract was re-hashed after this review and is byte-identical (§9). Revision 3 grew from 229 to 246
lines and from 35,330 to 43,144 bytes; every added block was read, not inferred from the delta.

---

## 2. Verification method and coverage

Read the full 229→246-line draft, then read each cited clause in the frozen contract directly. Hashes and
line numbers were checked independently; the draft's own citations were treated as claims, not evidence.
Contract re-hashed at the end.

| Check | Result |
|---|---|
| A2 decisions present | **12** (A2-D1 … A2-D12) |
| Impact rows | **29**, numbered 1–29 with no gaps |
| B1–B12 dispositions | **14**, arithmetic re-derived: 2 + 10 + 2 = 14 over 12 findings |
| §3.1 count check | Correct, including the B12 split rationale |
| Forbidden edit verbs in §2 | **None** — no `Replace`/`Redefine`/`Rewrite`/`Delete` |
| Line ranges re-verified | 30 of 32 exact; two border errors (§17.7.1, §17.11.5) |
| Impact map vs supersession map | No contradictions; asymmetries are granularity only (see F1 for the one substantive exception) |
| Two-process cap, cancellation precedence, acceptance barrier | Preserved |

Contract clauses read directly for this review: §3.2 (165–197), §4.6 (386–400), §4.7 (412–417),
§5.3 (480–488), §5.5 (via §2.5 cross-reference), §17.1 (1065–1075), §17.2.4–17.2.8 headings and
§17.2.5, §17.4.1 (1318–1335), §17.4.4 (1360–1380), §17.4.5 (1381–1387), §17.4.7 (1397–1402),
§17.6.1 (1670–1700), §17.6.2 (1701–1718), §17.6.3 (1720–1732), §17.7.1 (1736–1760),
§17.7.2 (1762–1790), §17.7.3 (**1792–1808**), §17.7.4 (1810–1822), §17.8, §17.9 (1858–1912),
§17.11.3 (2049–2062), §17.11.4 (2064–2080), §17.11.5 (2082–2098), §17.11.6 (2120–2138),
§17.12 (2169–2226), §17.15, §17.17, §17.18.

---

## 3. Disposition of the prior review (R1–R14)

| Prior finding | Status | Evidence |
|---|---|---|
| **R1** — additive method vs in-place rewrites | **RESOLVED** | §2 preamble now defines Add / Supersede / No textual change / Confirm; rows 4, 7, 12, 15, 19 read "No textual change; preserve the original." Grep for `Replace|Redefine|Rewrite|Delete` returns nothing. Supersession is expressed only through §17.19.3; §17.19's precedence clause limits effect to matters listed there. |
| **R2** — §17.7.4 could not represent the new combination | **RESOLVED** | Row added in both §17.19.3 and §17.19.4: `ProcessState=alive`, `State=unavailable`, `Available=true`, `Loaded=false`, `Ready=false`; existing rows untouched; impact row 17 changed from "Confirm" to "Add one row." |
| **R3** — `unavailable` had no retry transition | **ADDRESSED, new gaps** (F2, F3) | A2-D12 adds `unavailable` + explicit request → `loading`, bounded to one attempt, no automatic retry. |
| **R4** — `ProcessAlive` supersession too narrow | **RESOLVED** | Map now names "entire `ProcessAlive` sentence"; A2-D4 gives a complete five-value `TTSProcessState` projection; §17.19.4 repeats it. Border citation wrong (F7). |
| **R5** — duplicate idle-close transition row | **RESOLVED** | A2-D8 states idle close emits `shutdown` and "does **not** add a duplicate transition row"; §17.19.3 and row 16 repeat it. |
| **R6** — rejected shutdown mis-mapped; close-stdin missing | **RESOLVED** | Both A2-D7 and A2-D8 now classify a `shutdown` error response as "already stopping → close stdin, await exit, classify by exit code 0/nonzero"; transport failure separately. Close stdin is A2-D7 step (3) and appears in A2-D8's sequence. Consistent with §17.9 line 1906 and §17.11.5 exit codes. |
| **R7** — A2-D7 pre-empted active synthesis | **RESOLVED** | "An identity change executes under the existing per-instance serialization lock. On entry the adapter MUST assert that no synthesis is active … this path MUST NOT cancel or pre-empt it." |
| **R8** — `ServiceCapabilitiesKnown` invalidation incomplete | **RESOLVED** | A2-D6 now states an iff-invariant covering `stopping`, `exited`, `crashed`, `incompatible`, `backoff`, replacement, and pre-handshake; asset-specific `unavailable` on a live instance retains, non-authoritative clears. |
| **R9** — wrong `detail` bound citation | **RESOLVED** | A2-D3 says §17.2.5. Verified: §17.2.5 "Frame and field limits" (1149) and the 16-entry/64/256 bound at 1160. §17.2.4 is the error envelope, §17.2.7 malformed frames, §17.2.8 input disposition — all correctly cited alongside it. |
| **R10** — termination vs exit observer | **RESOLVED** | A2-D4: "The exit observer is intent-aware … MUST NOT be reclassified as `crashed`; only an unexpected exit is `crashed`." |
| **R11** — B11 over-claimed citations | **RESOLVED** | All five now in §1: §4.6 line 397 (A2-D1), §5.3 line 484 (A2-D3), §17.18.1 A1-D7 (A2-D3), §17.11.5 line 2089 (A2-D4), §17.9 lines 1867–1872 (A2-D3). All five also appear in §17.19.3. |
| **R12** — B8 asserted but not written | **RESOLVED** | A2-D1 now contains the substance: "`load` is mandatory and MUST initialize … not initialization timing. A lazy-only runtime that cannot initialize on `load` is non-conformant." |
| **R13** — off-by-two citations; B12b boundary | **MOSTLY RESOLVED** | §17.8 → 1846 ✓; §17.11.6 rule 4 → 2127–2130 ✓; B12b boundary now explicit in A2-D3. Two *new* border errors (F7, F8) and one range inconsistency (F9). |
| **R14** — `LastError` must name the failed model identity | **RESOLVED** | A2-D5 requires `LastError.Detail` to name `provider_id` and `model_identity` (and `engine` when present) — exactly §17.4.1's `load` parameter names (line 1328), and type-compatible with `TTSFailure.Detail map[string]string` (§17.12 line 2222). |

Every B1–B12 disposition is now executed in §1, including B8 and B11, so §3's traceability table matches the
normative text for the first time.

---

## 4. Findings

### F1 — BLOCKING: the mandated service-capability list cannot carry `Supported` in the implementation-ready shape

A2-D6 requires: "The status projection adds a distinct `ServiceCapabilities []TTSCapability` field … When
true, the list reports all five IDs with `Supported:true` or `Supported:false`."

The contract has **two different types with this name**:

| Location | Definition |
|---|---|
| §3.2 line 192–197 | `TTSCapability{ ID string; Version string; Supported bool; Detail string }` |
| §17.12 line 2212–2215 | `TTSCapability{ ID string; Version int }` — **no `Supported`** |

§17.12 is headed "Implementation-ready application interface summary" (line 2140), and its
`TTSProviderStatus` (2169–2184) is the authoritative status shape. Impact row 25 instructs: "Add `Loaded`,
`ServiceCapabilities`, and `ServiceCapabilitiesKnown` to the implementation-ready status shape." Adding
`ServiceCapabilities []TTSCapability` there yields elements of `{ID, Version int}` — there is no field in
which to report `Supported:true` or `Supported:false`, which is the operative requirement of A2-D6.

Impact row 3 compounds it: "§3.2, `TTSCapability` (lines 192–198) … **No type change**; document that
`Supported:false` is application-only." That is true of §3.2 and false of §17.12, where no `Supported`
exists to be application-only.

The same field addition also creates a *new* inconsistency: `Version` is `string` in §3.2 and `int` in
§17.12, so "mirror the §3.2 status fields" (the map's §17.12 entry) and "add `ServiceCapabilities
[]TTSCapability`" to both produces the identical new field with two different element types. §17.6.2's
extension versions are integers (line 1705), and §3.2's `string` is the outlier — but §17.19 never says
which one the new field uses.

**Consequence.** At Stage 3 an implementer cannot satisfy A2-D6's "reports all five IDs with
`Supported:true`/`Supported:false`" using the §17.12 shape, and is told by impact row 3 that no type change
is needed. Either the projection is unimplementable, or an unlisted type change is required — which is
precisely the failure mode R1 and R2 were raised to eliminate. `ServiceCapabilitiesKnown` is otherwise
correct; the defect is in what the list is required to contain.

**Required:** state which `TTSCapability` the projected list uses and reconcile `Version`; if §17.12's
shape is authoritative, add `Supported bool` there (and `Detail` if the projection needs it) and change
impact row 3 from "No type change" to the specific addition.

Verified not a problem: `ServiceCapabilitiesKnown=true` implies `Available=true` under A2-D5/D6, so no
stale-capability authorization path exists; the empty-list state is fail-safe and unreachable for dispatch
because every `Known=false` state has no live handshaked process.

### F2 — MODERATE: A2-D12's retry precondition depends on a helper transition §17.7.3 does not define

A2-D12 conditionally requires a retry: "if the helper is still live, **its helper-side state still permits
a new `load`**, and the provider/engine/config remains valid, the supervisor MUST attempt one new `load`."

§17.7.3 (lines 1795–1808) is the helper-side state table:

```
unloaded     → loading, stopping            accepting_work: true
loading      → ready, failed, stopping      accepting_work: false
ready        → synthesizing, loading (only after unload), stopping
failed       → stopping
```

and the prose at 1806–1808: "`failed` is a terminal-in-process state for a fatal provider/protocol fault …
**A normal missing model or unsupported request is not `failed`; it is a typed error and the helper remains
usable.**"

So a typed load failure is explicitly *not* `failed` — but the table gives `loading` no transition back to
`unloaded`, which is the only state from which `load` is legal. The helper's post-failure state is
therefore undefined, while the prose asserts the helper "remains usable." That gap is pre-existing, but
A2-D12 is the first clause to make it load-bearing, and **§17.7.3 appears nowhere** — not in the 29 impact
rows, not in §17.19.3, not in §2.1's confirmed-unchanged list.

**Consequence.** The retry guard cannot be evaluated as written, and the two plausible readings diverge:
read the table strictly and no retry is ever legal (the recovery path is dead code); read the prose and the
helper returns to `unloaded` (recovery works, but on an unstated transition). This is the same class of
defect as rev 2's R2 — a normative rule resting on an unrepresentable state.

**Required:** either cite §17.7.3 and state that a typed load failure returns the helper to `unloaded` with
`accepting_work=true` (an additive clarification in §17.19.3), or replace the guard with something
observable from the parent — the helper answered the failed `load` with a typed asset error, the process is
alive, and `health`'s `accepting_work` is true (§17.4.5 line 1385, result shape line 1636).

### F3 — MODERATE: "the next explicit TTS operation" is broad enough to load a model on a voice-list query

A2-D12 triggers on "the next explicit TTS operation after an asset-specific load failure." That term is not
defined, and the contract's operation set includes `catalog` and `health`. §17.4.5 line 1383 says `health`
"MUST NOT load a model, initialize ONNX Runtime, discover new …", and §5.3 line 485 permits `catalog` to
read assets but not load a model — with line 490 noting this "prevents a voice-list request from
initializing ONNX Runtime."

As written, a voice-list request could satisfy A2-D12's trigger and force a 120 s `load` (§17.8 line 1842),
inverting §5.3's separation and colliding directly with open decision **O-1** — "may the Character voice
picker start a helper for `catalog`…?" (§17.15 line 2323) — which A2-D10 and A2-D11 both say remains open.

A second ambiguity: §17.7.2 line 1790 leaves `unavailable` only via "helper/assets/config **updated** →
`absent`", while A2-D12 leaves it via explicit request → `loading`. Both now exist with different
conditions, and A2-D12 does not say which governs when configuration changed *and* a request arrives.

**Required:** name the trigger operation (`load` or `synthesize` only), state that `catalog`, voice-list,
and `health` requests never satisfy it, and give the update-gated row precedence when configuration has
changed.

### F4 — MODERATE: A2-D9 leaves the observed-terminal cancellation branch without a stated cleanup order

A2-D9 states one branch: "The adapter then observes the normal cancel grace; if no terminal response
arrives, it calls `KillTree`, waits for exit, and deletes only the request's exact output paths."

The other branch — a terminal response *does* arrive after `cancel_requested` is recorded — is left to
"the existing precedence rules." Under §17.11.6 rule 1 that success "is rejected and cleaned," and cleaning
means deleting a published WAV while the helper is still alive, which §17.11.4 **step 7** forbids: "It
deletes exact output leftovers only after any live helper has been terminated and observed, or after the
exit observer has already completed" (lines 2074–2076).

The order is recoverable only from a secondary artifact: §17.19.3 says rule 4 is superseded but that
"Cancellation precedence and cleanup ordering are unchanged," which preserves rule 4's "if the helper is
alive, the supervisor kills and observes process exit before deleting." A2-D9's own normative text does not
restate it, and §17.11.4 **step 7 is not cited in the supersession map at all** — the map covers "steps 4–6
(lines 2069–2073)", while step 7 sits at 2074–2076 and is the clause that actually governs deletion order.

**Consequence.** An implementer working from §17.19 alone, which A2-D11 makes authoritative, reads "kill
only if no terminal arrives" and will delete output of a cancelled request while the helper is alive.

**Required:** state both branches explicitly — kill, observe `Wait()`, then delete on the terminal-received
branch too — and add §17.11.4 step 7 (2074–2076) to the map as confirmed. This is not a barrier failure:
rules 1–3 and 5–7 remain intact and no success can revive a request.

### F5 — MODERATE: the §3.2 inspector cannot represent the distinction the new health row draws

§3.2's `TTSProviderStatus` (179–190) has no `ProcessState` and no `ProcessAlive`. `ProcessState` is defined
in §17.7.1 (1738–1742) and the derived `ProcessAlive` becomes normative in A2-D4.

The new §17.7.4 row's *only* discriminator is `ProcessState`: `alive` + `unavailable` + `Available=true`
(retryable under A2-D12) versus `absent` + `unavailable` + `Available=false` (update-gated). Through the
§3.2 interface both present as `State=unavailable`, and impact row 2 adds `Loaded`,
`ServiceCapabilities`, and `ServiceCapabilitiesKnown` but not the two fields the health table branches on.

§17.12 carries `ProcessState` and `ProcessAlive` (2171–2172), so the facts are representable *somewhere* —
but the amendment never says which shape is authoritative, and it edits the one that cannot express the
distinction. Combined with F1 this is the same root question: which status type is normative.

**Required:** either add `ProcessState`/`ProcessAlive` to §3.2 or state in §17.19.3 that §17.12's shape
supersedes §3.2's for status interpretation.

### F6 — MINOR: the consolidated unload row drops `busy`

§17.19.4 item 6 reads "`ready`/`degraded` + `unload` while idle → `unloaded`". The frozen row it
summarizes, §17.7.2 line 1781, is "`ready`/`busy`/`degraded` | `unload` while idle → `unloaded`". The
consolidation silently omits `busy`. Since A2-D11 preserves line 1781, item 6 is an imprecise summary
rather than a supersession — but §17.19.4 is presented as "a consolidated statement of the amended
behavior," so the mismatch is worth one word. (The frozen row's inclusion of `busy` alongside "while idle"
is itself odd; the amendment is the right place to say which reading governs.)

### F7 — MINOR: the superseded `ProcessAlive` sentence spans three lines, not one

§17.19.3 cites "§17.7.1, entire `ProcessAlive` sentence (line 1755)". The sentence begins at line 1753
("`ProcessAlive` is true for `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`,
and `stopping` while the process is alive; it is false for …") and ends at 1755. The effect statement in
both the map and §2 row 15 is correct and complete; only the anchor is wrong. Behaviourally harmless —
the conflicting values (`incompatible`, `unavailable`) are on line 1755 — but "entire sentence (line 1755)"
is self-contradictory.

### F8 — MINOR: the no-restart sentence is 2089–2090, not 2088–2090

§17.19.3 cites "§17.11.5, no-restart sentence (lines 2088–2090)". Line 2088 is the tail of the preceding
backoff bullet ("restarts per key in 10 minutes; reset after 10 minutes of stable readiness"); the
no-restart sentence is 2089–2090. Note this is a small regression: the rev 2 draft cited line 2089
correctly, and A2-D4 in this same revision still says "§17.11.5 line 2089", so the map now disagrees with
its own decision entry.

### F9 — MINOR: §17.12 is cited with two different ranges

The map says §17.12 "status interface (lines 2169–2222)"; impact row 25 says "(lines 2140–2228)". Both are
defensible — §17.12 runs 2140–2228 and the status struct occupies 2169–2184 — but the map's range stops at
2222, one line before `TTSFailure` closes at 2223, and omits the `TTSProviderState`/`TTSProcessState` note
at 2226. Pick one range per artifact and make it cover the status shape.

### F10 — MINOR: four capability IDs are left unclassified

§17.19.1 and A2-D1 declare the service vocabulary to be exactly five IDs. §17.6.1's table (1685–1695) also
lists `howl.tts.speed`, `howl.tts.language`, `howl.tts.ssml`, and `howl.tts.streaming`. The amendment never
states that these four are provider-scoped. A2-D3 then makes "a provider capability in the service array"
`protocol_incompatible`, so an implementer must decide whether a helper advertising `howl.tts.speed` in
`hello` is incompatible — and the text does not say. The decision report asserts the split
(§6.1: "Service capabilities use the five core service IDs"), so this is faithful to the owner's decision;
it is simply not restated where the classification rule now lives.

### F11 — MINOR: the new health row omits `ServiceCapabilitiesKnown`

§17.19.4's additive row lists `ProcessState`, `State`, `Available`, `Loaded`, `Ready`, and `LastError`. For
a "live, compatible helper" that completed the handshake, A2-D6 makes `ServiceCapabilitiesKnown=true`. It
is derivable, but this row exists precisely to make the legal combination explicit, and the owner's
representability question names this field. One clause closes it.

### F12 — EDITORIAL

- `Loaded` during `stopping`: A2-D5 clears `Loaded` on "process exit/crash … or identity/process
  replacement", not on `stopping`, so `State=stopping`, `Available=false`, `Loaded=true` is legal. That is
  coherent (resources are still resident until exit) and required by nothing, but it is unstated.
- Existing §17.7.4 row 1818 says "No helper **or dependency**". Under A2-D5 a dependency failure on a live
  helper now maps to the new row, so the phrase over-describes. The row is still reached only when
  `ProcessState=absent`, so no state is unreachable — wording only.

### F13 — EDITORIAL

- A2-D7 states step (7) twice ("only after exit is observed, start a fresh process for the new key if the
  transition is still permitted", then "After an observed exit, step 7 may start the fresh process …").
- A2-D7's "the path reports the invariant failure under existing rules" names no code. The registry has
  `internal_error` ("Unexpected provider/helper fault", line 1907); naming it would make the sentence
  actionable.

---

## 5. Representability check (all thirteen states)

Every combination is legal and noncontradictory as stated; the two unstated-but-derived cells are marked.

| `State` | `ProcessState` | `ProcessAlive` | `Available` | `Loaded` | `Ready` | `Known` | Source |
|---|---|---:|---:|---:|---:|---:|---|
| `absent` | `absent` | false | false | false | false | false | A2-D4/D5/D6 |
| `starting` | `starting` | true | false | false | false | false | A2-D4/D5/D6 |
| `handshaking` | `alive` | true | false | false | false | false | A2-D4/D5 |
| `unloaded` | `alive` | true | true | false | false | true | §17.19.4 item 3 |
| `loading` | `alive` | true | true | false | false | true | items 3–4 |
| `ready` | `alive` | true | true | true | true | true | item 4 |
| `busy` | `alive` | true | true | true | true | true | A2-D5 |
| `degraded` | `alive` | true | true | true | true | true | A2-D5 |
| `stopping` | `stopping` (or `alive`) | true | false | true * | false | false | A2-D5/D6, item 8 |
| `crashed` | `exited` | false | false | false | false | false | A2-D5/D6 |
| `incompatible` (terminating) | `alive`/`stopping` | true | false | false | false | false | A2-D4, §17.7.4 line 1819 |
| `unavailable`, asset-specific, live | `alive` | true | true | false | false | true * | new row; A2-D6 |
| `unavailable`, helper/config | `absent` | false | false | false | false | false | §17.7.4 line 1818 |
| `backoff` | `absent` | false | false | false | false | false | A2-D4 |

\* derivable, not stated (F11, F12). No unreachable state and no contradictory pair was found. The
per-identity scope question raised in the rev 2 review is resolved by §17.4.4's "one loaded provider per
process": only one identity can be loaded at a time, so a provider-scoped `State` with the identity carried
in `LastError.Detail` is coherent.

---

## 6. Answers to the eight review requirements

1. **Hashes and full read** — done; §1 artifacts table.
2. **Additive mechanism** — correct and internally consistent. Original clauses stand; supersession is
   expressed only through §17.19.3; no §2 row uses a rewrite verb; the §3.2 preservation claim about the
   `State` comment is accurate. The one contradiction is F1, where an impact row asserts "no type change"
   for a type that must change.
3. **Process-state consistency** — `ProcessAlive` now derives solely from `TTSProcessState` and all five
   values are covered (F7 is the citation only). All thirteen `State` values are representable and
   noncontradictory (§5). The live-helper `unavailable` row exists and is complete except for `Known`
   (F11). Recovery is explicit for asset-specific failures and update-gated for helper/config failures,
   with no automatic retry — but it rests on an undefined helper transition (F2) and an over-broad trigger
   (F3).
4. **Capability validity** — the invariant is correct: true iff the current instance completed a compatible
   handshake and has not entered `stopping`, `exited`, `crashed`, `incompatible`, `backoff`, or been
   replaced. Every exit, stop, incompatible result, and replacement clears it; asset-specific `unavailable`
   on a live instance retains it. No stale snapshot can authorize operations, because every `Known=false`
   state has no live handshaked process. The defect is in the list's *content* (F1), not its validity.
5. **Shutdown and identity replacement** — normal and error-response shutdown paths now preserve
   close-stdin, grace, `KillTree`, `Wait()`, and cleanup ordering (R6 resolved, verified against §17.9
   line 1906 and §17.11.5 exit codes). Identity replacement runs under the serialization lock and cannot
   pre-empt synthesis (R7 resolved). Cancellation precedence and the terminal-publication barrier remain
   intact; F4 is a missing branch statement, not a barrier failure.
6. **Tables and arithmetic** — 12 decisions, 29 rows, 14 dispositions; arithmetic re-derived and correct.
   The maps no longer contradict each other; the asymmetries are granularity (rows 10, 11, 21, 26, 28 have
   no map entry because nothing is superseded; the map carries three confirmations that are folded into
   broader rows). Two border citations are wrong (F7, F8) and one range is doubled (F9).
7. **Newly introduced issues** — F1 (unimplementable projection), F3 (retry trigger scope and O-1
   collision), F4 (missing cleanup branch), F5 (status-shape mismatch), F6 (dropped `busy`), F10
   (unclassified capability IDs). No behavior changes the accepted Policy B decision: A2-D1 matches
   §5/§6.1 of the decision report, A2-D3 matches §6.5, A2-D6 matches §6.6, A2-D2 matches §6.7, and A2-D7/D9
   match §6.3/§6.4.
8. **Defects vs preferences** — F1–F5 are defects with concrete consequences (unimplementable requirement,
   undefined state, wrong trigger, unsafe deletion order, unrepresentable distinction). F6–F11 are real but
   small inaccuracies. F12–F13 are editorial.

---

## 7. What revision 3 gets right

The corrections were made in the right places, not papered over: the supersession map was rewritten so
that "superseded" never means "deleted"; `ProcessAlive` became a derived property instead of a 13-state
table; `ServiceCapabilitiesKnown` became an invariant instead of an event list; the shutdown error path was
split into registry-consistent branches; and A2-D7 was scoped to the serialization lock rather than
inventing pre-emption. A2-D12 is the right shape for recovery — bounded, explicit-only, and it correctly
separates "failed attempt" from "provider verdict". The document's scope discipline is intact throughout:
nothing authorizes implementation, and §17.19.5 restates the boundaries.

---

## 8. Remaining owner decisions

1. **Which status shape is authoritative** (§3.2 inspection type vs §17.12 summary) and whether
   `TTSCapability.Version` stays `string` or becomes `int`. F1 and F5 both reduce to this.
2. **Whether a voice-list/`catalog` request may trigger the A2-D12 retry `load`** — this is the O-1
   boundary, and answering it decides F3's trigger wording.
3. **Whether an identity-change replacement spends restart budget** (§17.11.5's 3-per-10-minutes rule).
   A2-D7 starts a replacement process without saying whether it counts. Not previously raised; unstated in
   both directions.
4. **Whether the supervisor kills a live helper after a cooperative cancellation terminal** (F4). The map
   says ordering is unchanged, which implies yes; the §1 text implies no.

Unchanged and still open: O-1 – O-7, B12a (§17.10.3 hard-link vs atomic rename), B12b (§4.7 vs §17.6.2
extension error wording). None is decided — correctly — by this amendment.

---

## 9. Verdict and boundaries

**PASS WITH FINDINGS — conditional.** Revision 3 resolves R1 (method) and R2 (health-table row) and
executes all fourteen dispositions in normative text; the method, the state model, and the capability
invariant are now coherent. It is **not adoptable as written** because of F1: a normative requirement
(`Supported:true`/`Supported:false` for all five service IDs) has no field in the type the amendment tells
implementers to carry it in, and an impact row states that no type change is required. F2–F5 should be
settled in the same pass, since each is a one- or two-sentence clarification that would otherwise be
rediscovered at Stage 3. No blocking issue reopens a policy or architectural question; the accepted Policy
B decision is unchanged.

Minimum changes before approval:

1. Reconcile `TTSCapability` and state the authoritative status shape; fix impact row 3 (F1, F5).
2. Make A2-D12's helper-state guard observable or cite §17.7.3 (F2), and narrow "explicit TTS operation"
   to `load`/`synthesize` (F3).
3. State the observed-terminal cancellation branch and add §17.11.4 step 7 to the map (F4).
4. Restore `busy` in §17.19.4 item 6; correct §17.7.1 to 1753–1755 and §17.11.5 to 2089–2090; pick one
   §17.12 range; classify the four feature capability IDs; add `ServiceCapabilitiesKnown` to the new row
   (F6–F11).

**Boundaries of this review.** Read-only throughout. No edit to the amendment, contract, decision report,
prior audits, code, tests, fixtures, or Stage 2 artifacts. No implementation, no caller wiring, no Stage 3,
no dependency download, no packaging, no commit, no push. The contract was re-hashed after this review and
remains `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393`.

**The amendment remains unadopted unless the owner explicitly accepts it. Stage 3 remains unauthorized
regardless of this verdict.**
