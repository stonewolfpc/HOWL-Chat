# Adversarial Review — C-1 TTS Capability Policy Amendment, Revision 5

**Review date:** 2026-11-07
**Reviewer:** Mimo (independent, read-only)
**Subject:** `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`, revision 5
**Governing contract:** `DESIGN-howl-tts-helper-architecture-2026-11-06.md` (frozen)
**Policy decision:** `DECISION-howl-tts-capability-policy-c1-2026-11-06.md`
**Status of subject:** DRAFT FOR OWNER REVIEW — NOT ADOPTED. This review adopts nothing.

---

## 1. Identity of the artifacts

| Artifact | SHA-256 | Size |
|---|---|---|
| Amendment revision 5 | `8c728f2b2bbaa25533db4775e2fb55f850e6c3325977b76b31ac210c55c7cdbd` | 302 lines / 60,121 bytes |
| Frozen contract | `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` | 2,416 lines |
| C-1 decision report | `c25a0d3fd38e0767a4b7bf278cb1d6617e2e53caf11053dabb71449578c4c8ad` | 302 lines |
| Revision 4 audit | `772707bb5cb85b2002c668b602daeae0a330e3b2f3e06baaa8a6dc0edc550404` | 310 lines |
| Revision 3 audit | `9ceaab3038ac8855a87c53ef5c84d2e26a081ff98f3d12fb1e948635f417075b` | 412 lines |
| Revision 2 audit | `db6beb0a604cb65af4fe82a97d29db42395bc1312eabce051398b48ef4163bb3` | 330 lines |
| Revision 1 audit | `258fd1c1c981c2000ddcca4e87faea1e667d2b24114b68a48c8dea021d9c870e` | 358 lines |

Revision 4 was 278 lines / 53,108 bytes, so revision 5 adds 24 lines. The contract and
decision report are byte-identical to the values revision 4 carried and to the values asserted
in the amendment's own header (lines 6 and 8), so no comparison artifact moved under this
review.

No revision-5 hash was supplied for comparison. The value above is computed from the file under
review; the owner should confirm it matches the copy they intended to submit. A mismatch would
invalidate this review, not the amendment.

Prior audits are unchanged: revision 1 `258fd1c1…`, revision 2 `db6beb0a…`, revision 3
`9ceaab30…`, revision 4 `772707bb…`. The last two are cited in the amendment header at lines
11–14.

### Method

Every citation in the amendment is treated as a claim, not as evidence. Each cited line range
was read in the frozen contract and checked against the sentence or table row it is said to
name. Every disposition row (B1–B12b, F1–F13, N1–N8) is checked against the normative clause it
claims to resolve, not against its own description. Every owner decision is checked against the
policy it says it implements. Nothing in the amendment, contract, decision report, prior audits,
code, tests, fixtures, or Stage 2 artifacts was modified.

---

## 2. Verdict

**PASS WITH FINDINGS — recommend adoption after one one-line scoping correction (N9).**

All six focus areas the owner nominated pass. The revision 4 blocking defect is gone, the N1
identity collision is closed by an interlocking pair of clauses that works in both directions,
the budget-exhausted state is now fully representable, `detail["missing"]` is specified without
touching the wire DTO, provider features are projected through the catalog, and every N and F
disposition checked out against the clause it cites.

One defect found on this pass, N9, is not a policy error and not a regression introduced by
revision 5: it is a scoping defect inherited from revision 4, introduced by the F3 correction
that revision 4 made and that the revision 4 audit — this reviewer's — endorsed. It is a
one-word fix. It is recorded as moderate rather than blocking because the amendment's own
precedence rule (§17.19.2, §17.19.3) prevents the over-broad sentence from superseding the
frozen clause it collides with, and because the intent of the sentence is not in doubt. It
should still be fixed before adoption, because a literal implementation of the sentence as
written cannot synthesize at all.

Three editorial notes (E1–E4) and four immaterial citation nits are recorded below. None
requires a revision.

---

## 3. The six focus areas

### 3.1 Requirement 1 — same-identity recovery versus different-identity replacement

**Pass.** No path sends a differing-identity `load` to a helper that holds, or has failed on,
that identity.

The mechanism is a two-part interlock, and both parts are now present and mutually consistent:

- A2-D7 defines an identity change as "any difference in `provider_id`, `engine`, or
  `model_identity` from the currently loaded **or failed** identity" (line 76). The "or failed"
  phrase is the part that makes recovery work, and it was added in this revision.
- A2-D12 limits in-place recovery to "the same `provider_id`, `engine`, and `model_identity` as
  the failed attempt" and then states the complement explicitly: "If the requested identity
  differs in any of those three fields, `A2-D7`'s stop-and-replace procedure governs; the
  supervisor MUST NOT send a differing-identity `load` to the old helper" (line 81).

Checking the four reachable cases rather than the sentence:

| Situation | Governing text | Result |
|---|---|---|
| `unavailable` after failed load for A, helper live `unloaded`, explicit load for A | A2-D12 same-identity branch | one in-place `load`; `unavailable` → `loading` |
| `unavailable` after failed load for A, explicit load for B | A2-D12 → A2-D7 (B differs from the **failed** identity) | stop-and-replace; no `load` to the old process |
| `ready`/`busy`/`degraded` on A, explicit load for B | A2-D7 (B differs from the **loaded** identity) | stop-and-replace; helper's `provider_already_loaded` never reached |
| cleanly `unloaded` after a successful `unload`, explicit load for B | §17.7.3 `unloaded` → `loading`; §17.4.4 | in-place `load` for B is correct |

The fourth row is not a hole. After a successful `unload` no identity is loaded and none has
failed, so A2-D7's definition does not classify the load as an identity change, and §17.4.4's
`provider_already_loaded` precondition (§17.8 line 1840: "until unload") no longer holds. Sending
`load` for B to a helper that holds nothing is what the frozen protocol requires; forcing a
replacement there would be pure process churn. The guarantee is therefore exactly as strong as
it needs to be: no differing-identity load while an identity is held or has failed.

A2-D7's replacement sequence is intact and correctly ordered — assert no active synthesis under
the lock, never pre-empt (violation is `internal_error`), `shutdown`, grace, close stdin, await
exit, `KillTree` if still alive, observe `Wait()`, and only then start a fresh process. The
`shutdown`-error and transport-failure branches match A2-D8's terminal outcomes exactly, and
"on a transport failure or nonzero exit, this path starts no replacement" closes the obvious
replacement loop.

### 3.2 Requirement 2 — budget-exhausted availability

**Pass.** `Available=true` with `ProcessAlive=false` and `ServiceCapabilitiesKnown=false` is
coherent in every place it appears, and each of the three fields is derived from a different
authority, which is what makes the combination stable rather than accidental:

- `Available` is defined by A2-D5 as a configured-provider flag and explicitly "MUST NOT be read
  as proof that the helper process is alive". Case (b) then spells the combination out:
  `ProcessAlive=false`, `Loaded=false`, `Ready=false`, no automatic restart.
- `ProcessAlive` is derived only from `TTSProcessState` (A2-D4), so it says nothing about
  provider configuration.
- `ServiceCapabilitiesKnown` is false "after any exit" (A2-D6), so a exited process necessarily
  invalidates the snapshot without touching the other two fields.

The combination appears in four places and they agree: A2-D5 (the definition of case (b)),
§17.19.4 item 12 (the state rule), §17.19.4 table row 2 (line 128), and owner decision 13. The
new row is keyed `ProcessState=exited`, `State=unavailable`, which is disjoint from the frozen
`exited`+`crashed` row at line 1817 and from the `absent`+`unavailable` row at 1818, so the six
rows of §17.7.4 remain mutually exclusive on `(ProcessState, State)`.

Two supporting points that the revision gets right and that a weaker draft would have missed:

- A2-D5 keeps `Available` historical-compatibility evidence ("a compatible handshake
  established") and says budget exhaustion "does not redefine provider compatibility". That is
  the only reading under which `Available=true` after an exit is not a contradiction.
- `Loaded` and `Known` are deliberately *not* historical: both follow the current instance. The
  draft states this distinction in A2-D5 and A2-D6 rather than leaving it to inference.

The one unenumerated write-in is E1 below.

### 3.3 Requirement 3 — `loading` → `unloaded` as a table amendment

**Pass.** The edge is now an amendment to the legal-transition table, not explanatory prose, and
the labeling is consistent in all three instruments:

- A2-D12 (line 81): "This is an additive §17.7.3 transition edge: `loading` → `unloaded` on a
  typed non-fatal load failure, with `accepting_work=true`; no new helper state is introduced."
- Supersession map (row 100): "**Amended by A2-D12:** add `unloaded` to the `loading` row's legal
  transitions … fatal faults remain `failed`."
- Impact row 29: "Amend the legal-transition row: add `loading` → `unloaded` …"

The instrument is the right one. §17.7.3's table (1797–1804) is keyed `Helper state | Legal
transitions | accepting_work`, and its `loading` row currently reads `ready`, `failed`,
`stopping` with `false`. Adding `unloaded` to that cell is a genuine table edit, and the previous
"Clarified" label — which revision 4 used while adding an edge — is gone everywhere.

The edge also repairs an internal inconsistency in the frozen text rather than inventing policy:
line 1806–1808 says a missing model is "a typed error and the helper remains usable", which is
impossible if a typed failure leaves the helper stuck in `loading` with `accepting_work=false`.
The amendment supplies the missing half of a sentence the freeze already contains.

E2 below records one precision note about where `accepting_work=true` attaches.

### 3.4 Requirement 4 — `detail["missing"]` encoding

**Pass.** The representation matches the frozen envelope and cannot be parsed ambiguously.

- The carrier is unchanged: §17.2.4's `detail` is a `map[string]string` with at most 16 entries,
  keys of 1–64 bytes and values of 0–256 bytes (1146–1147; §17.2.5 table row at 1160). The
  application mirror `TTSFailure.Detail` is the same type (2222). No array, no new container, no
  wire change, and §17.5.1's `Capability` stays `{id, version integer}`.
- Ambiguity is impossible for structural rather than stylistic reasons: the v1 service
  vocabulary is the closed five-ID set of §17.6.1 (1687–1691), no ID contains a comma, the
  encoding forbids spaces and duplicates, and the order is fixed to table order. The result is a
  canonical byte string, so even string comparison is well defined.
- It cannot overflow: only the two mandatory IDs can appear, and
  `howl.tts.load,howl.tts.synthesize` is 31 bytes against a 256-byte value limit.
- It does not create a behavioural dependency on `detail`. The outcome is decided by the code
  alone, so the encoded list stays within §17.2.4's "small diagnostic facts only" and does not
  promote an optional diagnostic field into protocol logic.
- Style is consistent with the one existing example: §17.5.7 (1674) uses scalar string values
  under named keys (`{"min": "0.5", "max": "2.0"}`), not arrays.

The only thing the encoding does not define is what a *malformed* producer should do — e.g. an
unordered or duplicate-bearing list. That is a conformance question for the emitter, not an
ambiguity for the reader, so it needs no text.

### 3.5 Requirement 5 — capability projection and registry purity

**Pass.** Provider features are projected through catalog fields that already exist, and the
service/provider split is stated in both directions.

- §17.12's `TTSProviderCatalog` (2186–2198) already carries `Languages []string`,
  `Speed *TTSSpeedRange` (with `Supported`, `Min`, `Max`), `SSMLSupported`, and `StreamSupported`.
  A2-D2 names exactly these four as the provider-feature surface, so no new field is invented and
  no range data is duplicated into a boolean list.
- The registry is not polluted. A2-D6 keeps `Capabilities` provider-only and adds a separate
  process-scoped `ServiceCapabilities`; §17.19.1 forbids mixing the two sets in either the wire
  descriptors or the status projection; the four feature IDs are classified provider-scoped and
  barred from `hello.result.capabilities`; `howl.tts.synthesize` is named as the one intentional
  shared ID. §17.6.1's table (1685–1695) already separates the five service rows from the four
  feature rows, so the classification matches the frozen data.
- `Supported:false` is confined to `ServiceCapabilities`, and line 59 says so. The type
  `TTSCapability` carries the field in both arrays, but an advertised provider capability is
  true by presence, so the false case exists only on the service side. This is inherent to
  reusing one application type — which was F1's fix — and is not a registry intrusion.
- No behavioural gate rides on `Available`, so the A2-D5 redefinition cannot loosen
  authorization: synthesis is gated by `Loaded`/`Ready` (`Ready` is defined via §17.7.1's
  `ready`/`busy`/`degraded`), and the renderer-facing catalog keeps its own separate
  `TTSProviderCatalog.Available` (§17.12 line 2247–2249 maps `ProviderCatalog` into
  `TTSVoiceCatalog`, not `TTSProviderStatus`). The rename hazard — a UI reading "available" while
  assets are missing — is real for a *status* panel and is the owner's decision 3; it is not a
  propagation of the catalog's meaning.

### 3.6 Requirement 6 — traceability of N1–N8 and F1–F13

**Pass.** Every disposition row checked out against the clause it claims to resolve. The detail is
in §4 and §5. Two rows deserve comment here because they are the ones a reviewer should not take
on trust:

- **N7** claims both audit hashes are "in the header". They are: revision 3 at line 12, revision
  4 at line 14, each with its filename. The revision 2 hash is still absent, which is fine; it was
  never required.
- **F12** is now claimed "Resolved" rather than half-resolved, and the claim is defensible
  because the instrument changed: a new supersession row (98) scopes "or dependency" to
  startup-preventing dependencies. That scoping is not a stretch, because row 1818 is itself keyed
  `ProcessState=absent` — a live helper cannot match the row at all, so the narrow reading is the
  only consistent one. The map agrees with impact row 17, which narrows it "to the no-helper case".

### 3.7 The 13 owner decisions and the open carry-forwards

**Pass.** Thirteen decisions are recorded (§4 rows 1–13: seven original, four from revision 3, two
from revision 4). None closes an open question:

| Open item | Why it stays open | Evidence |
|---|---|---|
| O-1 — may the picker start a helper for `catalog`? | A2-D12 bars `catalog` from triggering a *load* and adds no restart path, but §17.7.2 line 1766 still lets an `absent` helper start on a catalog request. Starting is untouched. | 1766; A2-D12 last-but-three sentence |
| O-2 — WAV retention | A2-D9 deletes "only the request's exact output paths", which is §17.11.4 step 7's existing rule; published-file retention is not addressed. | 2074–2076 |
| O-3 — `auto` fallback | Named and refused: "O-3 remains open and is not decided here." | A2-D10; impact row 26 |
| O-4 — kokoro bridge timing | Untouched; `model_identity` is only a key field. | A2-D7 |
| O-5 — ORT version/DLL/CRT | Untouched. | — |
| O-6 — asset layout and manifest | A2-D5 requires `LastError.Detail` to name identity, not layout. | A2-D5 |
| O-7 — one binary versus many | A2-D1's floor constrains capabilities per helper, not packaging. | A2-D1 |
| B12a — hard link versus atomic rename | Explicitly "Not amended"; recorded in the map and again in §5. | rows 113, §5 |
| B12b — extension error wording | Explicitly "Not amended"; the amendment deliberately uses neither wording. | rows 114, §5 |

One editorial observation: §17.18.3 (line 2410) carries A1's sentence "O-1 through O-7 remain
exactly as listed in §17.15; this clarification decides none of them." §17.19.5 has no equivalent.
The substance is correct — nothing is closed — but a reader comparing A1 with A2 will notice the
missing sentence. See E3.

---

## 4. N1–N8 verification

| Finding | Claim | Verified against |
|---|---|---|
| **N1 — identity-specific recovery** | A2-D12 reuses a live helper only for the same identity; a difference follows A2-D7 and is never sent to the old helper; replacement stays budget-exempt | **Yes.** A2-D12 line 81 carries all three clauses; A2-D7 line 76 carries the "or failed identity" definition. Both directions checked in §3.1. |
| **N2 — §17.7.3 citation** | Map and impact row 29 now cite 1795–1808 | **Yes.** Row 100 and row 29 both read 1795–1808, which is the heading (1795), table (1797–1804), and prose (1806–1808). No "1792" remains anywhere in the file. |
| **N3 — §4 cleanup summary** | Decision 11 now says a terminal does not by itself require a kill | **Yes.** §4 row 11 and §17.19.4 item 9 both say "does not by itself authorize `KillTree`" and reserve the kill for when cleanup must proceed on a live helper, matching A2-D9 and step 7 (2074–2076). |
| **N4 — legal transition label** | `loading` → `unloaded` is now an amendment, not a clarification | **Yes.** See §3.3. |
| **N5 — provider-feature projection** | `Supported:false` scoped to service capabilities; provider features via catalog | **Yes.** §3.5. |
| **N6 — budget-exhausted availability** | `Available=true`, `ProcessAlive=false`, no automatic restart | **Yes.** §3.2. |
| **N7 — `detail.missing` and header hash** | Single `missing` key, comma-separated ordered; revision 3 and 4 hashes present | **Yes.** A2-D3 line 72; header lines 12 and 14. |
| **N8 — line 1818 wording** | Row 98 corrects 1818 to no-helper/startup-prevention without rewriting frozen text | **Yes.** Row 98, and the correction is the only one consistent with the row's `ProcessState=absent` key. |

---

## 5. F1–F13 verification

| Finding | Claim | Verified against |
|---|---|---|
| **F1 — capability type collision** | §17.12 authoritative; single `TTSCapability` `{ID, Version int, Supported, Detail}`; wire unchanged; row 3 records it | **Yes.** §17.12 (2169–2223) is named authoritative; line 45–51 and 167 record the supersession; §3.2's `Version string` (line 194) and the nine-value `State` comment (line 182) are superseded, not edited. |
| **F2 — undefined recovery transition** | Typed non-fatal load failure returns `unloaded` with `accepting_work=true`; fatal stays `failed` | **Yes.** A2-D12; row 29; consistent with 1806–1808. |
| **F3 — over-broad recovery trigger** | Only explicit `load` triggers recovery; five request kinds never do; update-gated row has precedence | **Partly.** The list and the precedence are right, and item 11 states the scope correctly ("do not trigger it"). But A2-D12's own wording, "never trigger a load", is broader than the rule it means, and collides with the frozen ensure-load. See **N9**. |
| **F4 — cancellation terminal branch** | Both branches stated; step 7 confirmed in the map | **Yes.** A2-D9 and row 107; rule 4 superseded by row 110; `internal_error` named in A2-D7. |
| **F5 — status shape unrepresentable** | §17.12 authoritative and already carries `ProcessState`/`ProcessAlive` | **Yes.** 2171–2172. |
| **F6–F11, F13 — minor traceability** | `busy` restored; 1753–1755; 2089–2090; one §17.12 range; four feature IDs classified; `ServiceCapabilitiesKnown` in the new row; `Loaded` during `stopping`; duplicate step removed; `internal_error` named | **Yes**, all eight. §17.19.4 item 6 includes `busy`; row 96 and 108 cite exactly 1753–1755 and 2089–2090; §17.12 is cited as 2140–2228 throughout; rows 127–128 carry `ServiceCapabilitiesKnown`; A2-D5 states `Loaded` "remains true during `stopping`". |
| **F12 — line 1818 wording** | Now scoped to no-helper/startup-prevention | **Yes.** See §3.6. |

---

## 6. Findings

### N9 — MODERATE (inherited from revision 4; should be fixed before adoption)

**A2-D12's synthesize clause is unscoped and collides with the only reachable load path.**

The text (line 81): "Recovery is triggered only by a new explicit `load` request; `synthesize`,
`catalog`, voice-list, `health`, and status requests never trigger a load. A `synthesize` request
while `Loaded=false` remains `provider_not_loaded` and does not initiate the retry."

Taken literally, the second sentence applies to every `Loaded=false` state, including the
ordinary `unloaded` state immediately after a successful handshake. That is not implementable
alongside the frozen interface. §17.12.1 (line 2235) defines `Synthesize` as "Ensure
handshake/load, synthesize, validate output, return typed result"; §17.11.4 step 2 (2067) says
"It ensures handshake and load". Worse, `Synthesize` is the *only* place the adapter can load at
all: `Initialize` explicitly does not spawn or load (2234), `ProviderStatus` "must not load a
model solely for status" (2241), `ProviderCatalog` "must not load a model" (2242), `GetVoices`
"never spawn or block" (2237), and no new load API is authorized. So if synthesize never triggers
a load, nothing does, and no provider can ever become `Loaded`.

Why this is moderate rather than blocking: the amendment's precedence rule saves it. §17.19.2
says the section controls "only for the matter identified in §17.19.3", and the map's row for
§17.11.4 lists steps 4–6 and step 7 only — step 2 is deliberately not superseded, and §17.12.1 is
not in the map at all. An implementer applying the precedence rule therefore keeps ensure-load
and reads A2-D12 as governing recovery, which is plainly its intent. That is also how §17.19.4
item 11 reads it ("do not trigger **it**"), so the summary is already better scoped than the
decision it summarizes.

Recommended fix, one sentence:

> A `synthesize` request while `State=unavailable` returns `provider_not_loaded` and does not
> initiate the retry; the ordinary ensure-load of §17.11.4 step 2 and §17.12.1 is unchanged.

Lineage, stated plainly because it affects how much weight to give this: revision 4 added the
"five non-triggering request kinds" list while resolving F3, and the revision 4 audit endorsed
that row as complete. Completeness was checked; compatibility with the ensure-load path was not.
The defect therefore predates revision 5 and was not introduced by it, but it originated in a
correction this reviewer approved.

The same over-broad "explicit" language appears twice more, at lower risk: A2-D5 ("it becomes
`Loaded`/`Ready` only after an explicit load succeeds") and §17.19.4 item 3 ("explicit `load` is
required before synthesis"). Both are true in the sense that a load must precede readiness;
neither is safe to read as "the caller must load before calling `Synthesize`". One definition of
"explicit load", used in all three places, would close the whole class.

### N10 — MINOR

**§17.19.4 item 11 restates the recovery rule with an exhaustive "only" that omits the identity
condition.**

Item 11 says an explicit `load` "may re-enter `loading` **only on** a live helper whose
helper-side state is `unloaded` with `accepting_work=true` and whose configuration is valid". For
an explicit load for a *different* identity after a failed attempt, every listed condition holds,
so the item permits something A2-D12 forbids with a MUST NOT and A2-D7 routes to replacement.
The decision text controls — §17.19.4 is a consolidated statement, not a decision entry — but the
summary should not permit by omission what the decision forbids by prohibition.

Fix: append "and only for the same `provider_id`/`engine`/`model_identity`" to item 11.

### E1 — EDITORIAL

A2-D5's `Available` enumeration does not name the `unavailable` case entered from `starting` +
`helper_start_failed` (§17.7.2 line 1769). The decision's closed-world lists — true for five
states plus `unavailable` cases (a) and (b), false for seven states plus "missing helper or
invalid provider configuration" — leave that entry unassigned. The definition resolves it
(`Available` requires "a compatible handshake established"; a first-attempt start failure has
none, so false), and false is the safe value, so no action is required. A reader may pause
nonetheless; naming the case would cost one clause.

### E2 — EDITORIAL

A2-D12's "with `accepting_work=true`" attaches a state property to an edge. §17.7.3's third
column is keyed by helper state, and the `unloaded` row already reads `true`, so the clause is
consistent and redundant — but an implementer could misread it as raising the `loading` row's
`false`. Suggested parenthetical: "(`accepting_work=true` is the `unloaded` row's existing
value)".

### E3 — EDITORIAL

§17.19.5 should carry A1's sentence that O-1 through O-7 remain unresolved (§17.18.3 line 2410).
Nothing in §17.19 closes any of them (§3.7), so this is symmetry and future-proofing, not a
correction.

### E4 — EDITORIAL

A2-D12 routes a live-but-`failed` helper through "the existing failed/crash path and `A2-D7`
stop-and-replace sequence". The budget treatment is split correctly by noun — the exemption in
A2-D7 attaches to "identity replacement", which a same-identity replacement is not, so §17.11.5's
budget governs — but the sentence invites the opposite reading by borrowing A2-D7's procedure
without saying the exemption does not travel with it. Five words would remove the doubt: "…the
`A2-D7` stop-and-replace sequence, without its budget exemption, govern".

Worth recording as a design observation rather than a defect: after a typed non-fatal failure the
helper is `unloaded` with `accepting_work=true`, so it could in principle take a new identity in
place; A2-D12 deliberately sends a differing identity through replacement instead. That is
conservative, explicit, budget-exempt, and costs one process start. Allowing in-place switching
after a clean failure would be a policy change, not a correction.

### Immaterial citation nits (no action)

- §3.2 is cited as "lines 179–198" (map row 93, impact rows 2 and 3); the block ends at 197 and
  198 is blank.
- §17.17 is cited as 2354–2368 (impact row 27); the table's last row, the A1 row the new A2 row
  will follow, is at 2369.
- §17.11.5's budget is attributed to "line 2088" (map row 109); the sentence spans 2087–2088.
- Map row 97 says "existing rows unchanged" while row 98 narrows row 1818. The two are consistent
  once read together, but they sit adjacent.

None of these misdescribes the clause it points at.

---

## 7. Representability of all thirteen states

Fields: process state, `ProcessAlive`, `Available`, `Loaded`, `Ready`,
`ServiceCapabilitiesKnown`, `ServiceCapabilities`. Every row is now derivable from a stated
rule; no cell requires an unstated decision.

| `State` | ProcessState | Alive | Available | Loaded | Ready | Known | ServiceCapabilities |
|---|---|---|---|---|---|---|---|
| `absent` | `absent` | false | false | false | false | false | empty |
| `starting` | `starting` | true | false | false | false | false | empty |
| `handshaking` | `alive` | true | false | false | false | false | empty (before handshake) |
| `unloaded` | `alive` | true | true | false | false | true | five IDs |
| `loading` | `alive` | true | true | false | false | true | five IDs |
| `ready` | `alive` | true | true | true | true | true | five IDs |
| `busy` | `alive` | true | true | true | true | true | five IDs |
| `degraded` | `alive` | true | true | true | true | true | five IDs |
| `stopping` | `alive`/`stopping` | true | false | true until exit | false | false | empty |
| `crashed` | `exited` | false | false | false | false | false | empty |
| `incompatible` | `alive`/`stopping`/`exited` | true while terminating | false | false | false | false | empty |
| `unavailable` (asset-specific, live) | `alive` | true | **true** | false | false | **true** | five IDs |
| `unavailable` (budget exhausted) | `exited` | **false** | **true** | false | false | **false** | empty |
| `unavailable` (helper/config) | `absent` | false | false | false | false | false | empty |
| `backoff` | `exited` | false | false | false | false | false | empty |

Checks:

- No state requires `Available=true` to imply a live process (A2-D5 forbids the reading).
- No state requires `Known=true` without a live, handshaken, non-stopping instance (A2-D6) —
  including the two `unavailable` variants, which differ exactly in `Known` and in `ProcessAlive`.
- `Loaded=true` with `Known=false` in `stopping` is the only mixed row, and it is stated:
  `Loaded` survives until exit is observed (A2-D5), the snapshot dies at `stopping` (A2-D6),
  `Ready` is false, so no new work can be dispatched against a stale snapshot.
- `incompatible` with a live process is stated by A2-D4 and matches frozen row 1819 ("process
  must be stopped").
- No state can be entered with a `Known=true` snapshot carried over from a replaced instance:
  A2-D6 clears it on replacement, and §17.19.4 item 7 says the old snapshot is invalid "as soon as
  replacement begins".

---

## 8. What this review did not do

- It did not verify the amendment against the implementation. Policy B is not implemented; the
  current implementation remains Policy A by design.
- It did not re-derive the C-1 policy decision itself. The decision to accept Policy B is the
  owner's and is out of scope here; this review asks only whether the text implements that
  decision without contradicting the frozen architecture.
- It did not adopt anything, authorize implementation, wire `Unload` or `Cancel`, authorize
  Stage 3, or authorize a commit or push.
- It did not modify the amendment, contract, decision report, prior audits, code, tests,
  fixtures, or Stage 2 artifacts.

The one point on which a reader should not take this review's word: the computed revision-5 hash.
If the owner's copy differs, re-review from the owner's bytes.

---

## 9. Conclusion

Revision 5 is adoptable after one scoping word. The six items the owner nominated all hold, the
N1 collision is closed in both directions, the budget-exhausted state is representable and
internally consistent, the `detail.missing` encoding fits the frozen envelope even in its
canonical-byte-string property, provider features stay in the catalog, and every disposition row
was checked against the clause it names rather than against its own description.

The remaining defect, N9, is inherited from revision 4 and was missed by the revision 4 review,
including by this reviewer. It does not change the policy, does not require a new owner decision,
and does not reopen any of F1–F13 or N1–N8. It should be corrected before adoption because the
sentence as written would strand the only load path in the adapter.

Adoption remains a separate formal decision. This review does not constitute it, and nothing here
authorizes Stage 3.

**End — revision 5 review.**
