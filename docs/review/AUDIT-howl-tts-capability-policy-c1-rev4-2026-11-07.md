# Adversarial Review — C-1 Amendment Revision 4 (§17.19 draft)

**Date:** 2026-11-07
**Reviewer:** Mimo (independent adversarial pass)
**Subject:** `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`
**Status:** INDEPENDENT REVIEW ARTIFACT — READ-ONLY; the amendment remains unadopted

**Verdict: PASS WITH FINDINGS — recommend adoption after one clarification.** Revision 4 resolves the
blocking type defect and the four moderate findings from revision 3. Every one of the six F-disposition
rows is executed in normative text, the authoritative status contract is now unambiguous, and the
recovery, cancellation, and replacement paths are coherent with the frozen contract. One moderate defect
remains: A2-D7 and A2-D12 both claim the identity-change-plus-explicit-load event and give opposite
instructions. Three minor citation/labelling items and four editorial residues should be fixed in the same
pass. No policy or architectural question is reopened.

---

## 1. Artifact hashes (verified on disk)

| Artifact | SHA-256 | State |
|---|---|---|
| Amendment draft, revision 4 | `49996de02e656fa45100dd057dbc2e124b2ec032d2b9ffdf76eb86bb1a99acdd` | **matches expected**; 278 lines, 53,108 bytes |
| Frozen architecture | `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` | unchanged |
| C-1 decision report | `c25a0d3fd38e0767a4b7bf278cb1d6617e2e53caf11053dabb71449578c4c8ad` | unchanged |
| Revision 3 audit | `9ceaab3038ac8855a87c53ef5c84d2e26a081ff98f3d12fb1e948635f417075b` | unchanged |
| Revision 2 audit | `db6beb0a604cb65af4fe82a97d29db42395bc1312eabce051398b48ef4163bb3` | unchanged |

The contract and the decision report were re-hashed after this review and are byte-identical (§9).
Revision 4 grew from 246 to 278 lines and from 43,144 to 53,108 bytes; every added block was read.

---

## 2. Verification method and coverage

Read the full 278-line draft, then read every cited clause in the frozen contract and in the decision
report directly. Citations were treated as claims, not evidence. Contract and decision re-hashed at the
end.

| Check | Result |
|---|---|
| A2 decisions | **12** (A2-D1 … A2-D12) |
| Impact rows | **31**, numbered 1–31 with no gaps |
| B1–B12 dispositions | **14** (2 + 10 + 2 over 12 numbered findings) |
| F-finding disposition rows | **6** (F1–F5 individually, F6–F13 grouped) |
| Forbidden rewrite verbs in §2 | **None** |
| "No type change" residue | **None** — the F1 contradiction is gone |
| Cited line ranges re-verified | 33 of 34 exact; §17.7.3 wrong by 3 lines (N2) |
| §3.2 mentions | All eight frame §3.2 as superseded or earlier proposal; none leaves it authoritative |
| Two-process cap, cancellation precedence, acceptance barrier | Preserved |

Contract clauses read directly for this review: §3.2 (177–198), §4.6 (386–403), §4.7 (412–417),
§5.3 (480–488), §17.1 (1064–1074), §17.2.5 (1149–1163), §17.3.2 (1277–1281), §17.4.1 (1323–1335),
§17.4.4 (1368–1379), §17.4.6 (1389–1395), §17.4.7 (1397–1402), §17.5.1 (1412–1424),
§17.6.1 (1683–1700), §17.6.2 (1701–1717), §17.6.3 (1719–1733), §17.7.1 (1751–1760),
§17.7.2 (1764–1790), §17.7.3 (**1795–1808**), §17.7.4 (1810–1822), §17.8 (1844–1848),
§17.9 (1876–1912), §17.10.3 (1945–1959), §17.11.3 (2049–2062), §17.11.4 (2062–2080),
§17.11.5 (2082–2098), §17.11.6 (2116–2138), §17.12 (2138–2228), §17.15 (2321–2328).
Decision report: §2.5 (72–84), §5 (214–231), §6.1–§6.7 (232–268).

---

## 3. Disposition of the revision 3 review

### 3.1 Blocking and moderate findings (F1–F5)

| Finding | Status | Evidence |
|---|---|---|
| **F1 — capability type collision (BLOCKING)** | **RESOLVED** | §17.19.1, A2-D6, and the map now state that §17.12's status shape is authoritative and that the single application type is `{ID string; Version int; Supported bool; Detail string}`. Impact row 3 no longer says "No type change"; it records the supersession, and row 25 records the concrete additions to §17.12. The phrase "No type change" no longer appears anywhere. `Version` is `int` in both the wire DTO (verified: `{"id": "howl.tts.synthesize", "version": 1}`, line 1418) and the application projection. |
| **F2 — undefined helper recovery transition (MODERATE)** | **RESOLVED** | A2-D12 states the rule normatively: a typed asset/model load failure that is not fatal "returns the helper to `unloaded` with `accepting_work=true`; `failed` remains the fatal-in-process state". §17.7.3 is now cited in the map and appears as impact row 29. Verified consistent with §17.7.3's `unloaded` → `loading` edge (line 1799) and line 1807–1808's "not `failed`; it is a typed error and the helper remains usable". Mechanism label imprecise (N4). |
| **F3 — over-broad recovery trigger (MODERATE)** | **RESOLVED** | A2-D12: "Recovery is triggered only by a new explicit `load` request; `synthesize`, `catalog`, voice-list, `health`, and status requests never trigger a load." The collision with O-1 and §5.3 line 490 is gone. Update-gated precedence is explicit: "If the helper/assets/config have changed, the §17.7.2 update-gated … row governs first." |
| **F4 — cancellation terminal branch (MODERATE)** | **RESOLVED** | A2-D9 states both branches and now cites step 7 by number: "if cleanup is required while the helper is live, the supervisor terminates and observes it before deletion, as step 7 requires". Impact row 22 records step 7 as confirmed and the map carries a dedicated step-7 entry (2074–2076). Verified against lines 2074–2076 and §17.11.6 rule 4 (2127–2130), whose own kill clause is likewise conditioned on the helper being alive. |
| **F5 — status shape unrepresentable (MODERATE)** | **RESOLVED** | §17.12 is authoritative and does carry `ProcessState` and `ProcessAlive` (verified lines 2171–2172). The new §17.7.4 row is expressible in it, and row 2 of the impact list correctly records that §3.2 needs no textual change because it is no longer the operative shape. |

### 3.2 Minor findings (F6–F13)

| Finding | Status | Evidence |
|---|---|---|
| **F6** — `busy` dropped from unload summary | **Resolved** | §17.19.4 item 6 now reads `ready`/`busy`/`degraded`, matching §17.7.2 line 1781. |
| **F7** — `ProcessAlive` sentence anchor | **Resolved** | Map cites lines 1753–1755. Verified: the sentence begins at 1753 and ends at 1755. |
| **F8** — no-restart anchor | **Resolved** | Map cites 2089–2090; A2-D4 cites line 2089. Verified: 2088 is the budget tail, the no-restart sentence is 2089–2090. The internal disagreement is gone. |
| **F9** — two §17.12 ranges | **Resolved** | Both the map and impact row 25 use 2140–2228. Verified: §17.12 heading at 2140, closing prose at 2228. |
| **F10** — unclassified capability IDs | **Resolved** | §17.19.1 classifies `speed`/`language`/`ssml`/`streaming` as provider-scoped and forbids them in `hello.result.capabilities`; A2-D3 makes misplacement `protocol_incompatible`; `howl.tts.synthesize` is named the intentional shared ID. Verified against §17.6.1 lines 1692–1695 and the descriptor example at 1418, which does advertise `howl.tts.synthesize`. |
| **F11** — `ServiceCapabilitiesKnown` missing from the new row | **Resolved** | The row now carries `ServiceCapabilitiesKnown=true` in both the map and §17.19.4. |
| **F12** — `Loaded` during `stopping`; line 1818 wording | **Half resolved** | A2-D5 now states `Loaded` "remains true during `stopping` until exit is observed". The "No helper **or dependency**" wording at 1818 is untouched, by design (N8). |
| **F13** — duplicate step, unnamed error code | **Resolved** | A2-D7 lists steps (1)–(7) once and names `internal_error`, verified present in the registry at line 1907. |

All six disposition rows in §3.2 of the draft describe what the normative text now does; none
overstates the change.

---

## 4. New findings

### N1 — MODERATE: A2-D7 and A2-D12 give opposite instructions for a differing-identity explicit load

A2-D7 opens: "An identity change executes under the existing per-instance serialization lock", then defines
the stop-and-replace sequence, and states "No success is reported and no differing-identity `load` is sent
to the old process."

A2-D12 states: "On a later explicit `load`, if the helper is live, its helper-side state is `unloaded` with
`accepting_work=true`, and the provider/engine/config remains valid, the supervisor MUST attempt one new
`load` on that helper and transition `unavailable` → `loading`."

A2-D12's retry conditions carry **no identity qualifier**. Consider the reachable case: a `load` for
identity A fails on a missing asset → supervisor `unavailable`, helper live in `unloaded` with
`accepting_work=true`. A caller then issues `load` for **identity B**. Every A2-D12 condition is satisfied,
so A2-D12 says send `load` on the existing helper; A2-D7 says an identity change is a stop-and-replace and
forbids a differing-identity `load` to the old process.

Both readings are defensible on their own — reusing a live `unloaded` helper is cheaper and the helper
holds no conflicting provider state, which is why `provider_already_loaded` (§17.4.4 line 1373) would not
fire — but the text does not pick one, and the two paths have materially different consequences for the
restart budget, the old process's lifetime, and the capability snapshot.

**Consequence.** Two conformant implementations diverge on a reachable state. One of them violates a MUST
in A2-D7. This is the only place in revision 4 where two normative decisions claim the same event.

**Required:** qualify A2-D12 with "for the same `provider_id`/`engine`/`model_identity`" and state that a
differing identity follows A2-D7. One clause; no other text changes.

### N2 — MINOR: §17.7.3 is cited as 1792–1808 in both places

The map row and impact row 29 both cite "§17.7.3, helper-side state (lines 1792–1808)". The §17.7.3
heading is at **1795**; lines 1792–1793 are still §17.7.2's closing prose about `backoff` and the exit
observer. Correct range: **1795–1808**. This is a newly introduced citation error (revision 3 did not cite
§17.7.3 at all). Behaviourally harmless — the table and the prose at 1806–1808 are inside the cited range
— but the anchor is wrong in two places.

### N3 — MINOR: the §4 summary of owner decision 11 contradicts A2-D9 and §17.11.6 rule 4

§4 states: "a terminal does not by itself authorize a kill; success is rejected and step 7 cleanup ordering
governs. **Only a missing terminal triggers the kill path.**"

A2-D9 states: "if cleanup is required while the helper is live, the supervisor terminates and observes it
before deletion, as step 7 requires". §17.11.6 rule 4 itself says "if the helper is alive, the supervisor
kills and observes process exit before deleting the exact final/`.part` paths" — conditioned on the helper
being alive and on deletion happening, not on a missing terminal.

So a terminal-received request whose output must be cleaned **does** terminate a live helper, and §4's
closing sentence denies it. §4 is non-normative, but it is presented as the record of the eleven owner
decisions, and the decision it misstates is the one that prevents deleting a file a live helper may still
be writing. Fix: replace the sentence with "a terminal does not by itself require a kill; a kill is
required only when cleanup must proceed on a live helper (§17.11.4 step 7)".

### N4 — MINOR: the §17.7.3 entry is labelled a clarification but adds a legal transition

The map says "Clarified by A2-D12 … No new state is introduced", and impact row 29's action is "Add
clarification". Neither is false, but the effect is a **new edge** in a table headed "Legal transitions":
`loading` currently reaches only `ready`, `failed`, and `stopping` (line 1800). A2-D12 adds
`loading` → `unloaded`.

Because §17.19 controls "only for the matter identified in §17.19.3", a careful reader will get the right
answer; a reader treating the amended table as exhaustive will not, and the retry guard becomes
unverifiable again. Recommend the row say what it does: add `unloaded` to the `loading` row's legal
transitions, with `accepting_work=true` after the typed failure.

### N5 — EDITORIAL: A2-D2's provider-side `Supported:false` projection has no vehicle

A2-D2 says "The application projects an absent capability as `TTSCapability.Supported=false`; `catalog`
remains the surface for provider/asset/voice availability." The only list that reports `false` is
`ServiceCapabilities`, which A2-D6 confines to the five service IDs. A provider feature the binary cannot
expose is absent from `ProviderDescriptor.capabilities` and has no entry in the provider `Capabilities`
array either, so nothing carries its `Supported:false`. The decision report's §6.1 asked for a
provider-capability projection; revision 4 moved that surface to `catalog`, which already has
`SSMLSupported`, `StreamSupported`, and `Speed *TTSSpeedRange` (lines 2196–2197) and is the better answer.
The sentence should be scoped to service capabilities so it describes the mechanism that exists.

### N6 — EDITORIAL: the `Available` enumeration omits budget-exhausted `unavailable`

§17.7.2 line 1789 defines a third route into `unavailable`: `crashed`/`backoff` + budget exhausted. A2-D5
enumerates `unavailable` twice — true for an asset-specific failure, false for missing helper or invalid
provider configuration — and that route is neither. The general definition ("the helper is present, the
handshake is compatible, and the provider is compiled/configured") resolves it to **true**, and A2-D12
correctly adds no retry path for it ("If no live helper remains, the existing `unavailable` → `absent`
update-gated row and restart-budget policy govern"), so the state is coherent. The enumeration simply reads
as exhaustive and is not. One clause closes it.

### N7 — EDITORIAL: `detail.missing` encoding; header hash

A2-D3 restricts `detail.missing` to the absent mandatory service IDs and bounds it by §17.2.5 (16 entries;
key 64; value 256 bytes — verified line 1160). At most two IDs, so the bound is trivial, but the field does
not exist in the frozen contract and its encoding (separator, ordering) is unstated while `TTSFailure.Detail`
is `map[string]string` (line 2222). Separately, the draft header cites the revision 3 audit without its
SHA-256, while every other referenced artifact has one.

### N8 — EDITORIAL: row 1818 still over-describes the helper/config case

The map says the §17.7.4 rows are unchanged, and A2-D5 now routes a dependency failure on a live helper to
the new row. Row 1818's "No helper **or dependency**" therefore names a case it no longer covers. It is
reached only when `ProcessState=absent`, so nothing is misclassified — wording only. (Carried over from
revision 3's F12; deliberately not fixed.)

---

## 5. Representability check (all thirteen states)

| `State` | `ProcessState` | `ProcessAlive` | `Available` | `Loaded` | `Ready` | `ServiceCapabilities` | `Known` |
|---|---|---:|---:|---:|---:|---|---:|
| `absent` | `absent` | false | false | false | false | empty | false |
| `starting` | `starting` | true | false | false | false | empty | false |
| `handshaking` | `alive` | true | false | false | false | empty | false |
| `unloaded` | `alive` | true | true | false | false | five IDs | true |
| `loading` | `alive` | true | true | false | false | five IDs | true |
| `ready` | `alive` | true | true | true | true | five IDs | true |
| `busy` | `alive` | true | true | true | true | five IDs | true |
| `degraded` | `alive` | true | true | true | true | five IDs | true |
| `stopping` | `stopping` | true | false | true (until exit) | false | empty | false |
| `crashed` | `exited` | false | false | false | false | empty | false |
| `incompatible` (terminating) | `alive`/`stopping` | true | false | false | false | empty | false |
| `unavailable`, asset-specific | `alive` | true | true | false | false | five IDs | true |
| `unavailable`, helper/config | `absent` | false | false | false | false | empty | false |
| `unavailable`, budget exhausted | `exited` | false | true * | false | false | empty | false |
| `backoff` | `absent` | false | false | false | false | empty | false |

\* derivable from A2-D5's general definition, not enumerated (N6). No unreachable state, no contradictory
pair, and no state in which a stale capability snapshot can authorize an operation: every `Known=false`
row has either no live process or a process that has not completed a compatible handshake.

---

## 6. Answers to the ten review requirements

1. **Hash and full read** — hash matches the expected value exactly; 278 lines; §1.
2. **Contract and decision unchanged** — both byte-identical; §1.
3. **F1–F5 and F6–F13** — F1–F5 all resolved with the remediation the revision 3 audit asked for; F6–F11
   and F13 resolved; F12 half resolved by design; §3.
4. **Authoritative status contract** — §17.12 is unambiguously authoritative, stated in §17.19.1, in both
   A2-D6 clauses, in the map, and in impact rows 2, 3, and 25. The wire `Capability` (`{id, version
   integer}`, verified at line 1418) and the application `TTSCapability` are distinct, and the wire gains no
   `supported` field. `Version` is `int` in both types; §3.2's `Version string` is explicitly superseded. All
   eight §3.2 mentions frame it as superseded or as the earlier proposal — none leaves it authoritative. The
   seven status fields are coherent across all thirteen states (§5).
5. **Helper recovery** — a typed non-fatal load failure returns `unloaded` with `accepting_work=true`;
   `failed` remains terminal-in-process; the helper can accept the next explicit `load` because §17.7.3
   already permits `unloaded` → `loading`; no automatic retry exists; `synthesize`, `catalog`, voice-list,
   `health`, and status requests are all excluded, and a `synthesize` while `Loaded=false` stays
   `provider_not_loaded`; the update-gated row has explicit precedence. Two precision items (N1, N4).
6. **Cancellation and cleanup** — both branches are distinct and stated; a cooperative terminal does not
   by itself authorize a kill; rejected success and exact-output cleanup preserve step 7, which is now
   cited by number; `KillTree` → observed `Wait()` → delete holds wherever termination is required; rule 4
   is the only rule narrowed, and the precedence rules and the terminal-publication barrier are untouched.
   One summary-level misstatement (N3).
7. **Identity replacement** — cannot pre-empt an active synthesis (lock invariant, violation is
   `internal_error`, no cancellation); fully serialized; the restart-budget exemption is explicit and does
   not create an unbounded path, because restart still occurs only on an explicit request (§17.11.5 line
   2087) and the budget for unexpected exits is expressly unchanged. The exemption's interaction with the
   retry path is the gap in N1.
8. **Tables and arithmetic** — 12 decisions, 31 impact rows, 14 dispositions, 6 F rows, all re-derived;
   every cited line reference re-verified, with one range wrong (N2).
9. **New contradictions** — N1 is the only normative one. No silent Policy B change: A2-D1 matches
   decision §5 and §6.1, A2-D2 and A2-D6 match §6.7, A2-D3 matches §6.2/§6.5, A2-D7/D8 match §6.3,
   A2-D9 matches §6.4, and A2-D10/A2-D6 match §6.6. O-1 is not silently closed: §17.7.2 line 1766 still
   permits a catalog request to start a helper from `absent`, and A2-D12 amends only the recovery trigger.
10. **Defects vs preferences** — N1 is a defect with a concrete divergence. N2–N4 are real inaccuracies in
    citations and labels. N5–N8 are editorial. Nothing in this review asks the contract to support a
    requirement it does not contain.

---

## 7. What revision 4 gets right

The blocking fix was made at the right level: rather than adding `Supported` to the wire DTO (which §6.7
and the frozen fixtures forbid), the draft names §17.12 as authoritative, defines one application
projection, and records the type change in both impact rows that touch it. The recovery decision is now
bounded in every direction — explicit trigger only, one attempt, no backoff, no restart path, update-gated
precedence — and it cites the helper-side table instead of inventing a mechanism. A2-D9's cleanup clause is
the correct reading of step 7 rather than a loosening of it. The replacement path's budget exemption is
justified in terms the frozen contract already uses ("restart only on the next explicit TTS request"), so it
does not create the unbounded recovery loop that a naive exemption would.

Scope discipline is intact: §17.19.5 restates every boundary, B12a/B12b stay deferred, O-1–O-7 stay open,
and nothing authorizes code, Stage 3, or a commit.

---

## 8. Remaining owner decisions

1. **Differing-identity `load` while `unavailable`** — reuse the live `unloaded` helper, or stop and
   replace? (N1.) This is a genuine behavioural choice with a cost/latency tradeoff, not a typo fix.
2. **`Available` for budget-exhausted `unavailable`** — confirm the definitional answer (true) or state the
   opposite. (N6.)
3. Carried forward and untouched: O-1 – O-7, B12a (§17.10.3 hard-link vs atomic rename), B12b (§4.7 vs
   §17.6.2 extension error wording).

---

## 9. Verdict and boundaries

**PASS WITH FINDINGS — recommend adoption after one clarification.** No blocking defect remains. F1 is
fixed in a way that respects the frozen wire shape, the four moderate findings from revision 3 are closed
in normative text, and the status contract, state projection, recovery rule, cancellation branches, and
replacement path are each individually coherent with §§17.1–17.18.

Before adoption:

1. **N1 (required):** qualify A2-D12's retry to the same provider/engine/model identity and route a
   differing identity to A2-D7.
2. **N2–N4 (should be fixed in the same pass):** correct §17.7.3 to 1795–1808, restate §4's decision-11
   sentence, and record the new `loading` → `unloaded` edge as a transition rather than a clarification.
3. **N5–N8 (editorial):** scope A2-D2's `Supported:false` sentence to service capabilities, add the
   budget-exhausted `unavailable` value to A2-D5's enumeration, define `detail.missing`'s encoding and add
   the rev3 audit hash to the header, and correct row 1818's wording.

**Boundaries of this review.** Read-only throughout. No edit to the amendment, contract, decision report,
prior audits, code, tests, fixtures, or Stage 2 artifacts. No implementation, no caller wiring, no Stage 3,
no dependency download, no packaging, no commit, no push. The contract was re-hashed after this review and
remains `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393`.

**The amendment remains unadopted unless the owner explicitly accepts it. Stage 3 remains unauthorized
regardless of this verdict.**
