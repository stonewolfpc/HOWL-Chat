# Muse — Independent review of R-07 evidence report

**Date:** 2026-11-11
**Reviewer:** Muse (independent instance)
**Review type:** Evidence verification and challenge; read-only
**Subject:** `docs/research/R-07-EVIDENCE-REPORT-2026-11-11.md`
**Subject identity:** `e3860e3d5d6ca71724d31d0021897819c968afc1b8e1aad09fc011f8d9a67134`, `46,107 B / 876 lines` — **Observed**
**Status:** Complete; no document or source file modified

## 1. Verdict

**Findings confirmed. Report provenance requires correction.**

`EO-01`, `EO-02`, and `EO-03` are all supported by the source I traced independently. The central claims are not interpretations of ambiguous code; they are direct readings of functions whose behavior contradicts their own comments and parameter names. In one case the evidence is stronger than the report states.

The report itself, however, fails its own identity standard. Two of its six recorded digests do not reproduce, one of which is not even a valid SHA-256 length, and one cited artifact is mislabeled. The completion checklist marks those verifications as passed. This does not change any finding, and it requires no new investigation to fix — but it must be corrected before the report is cited as an evidence base, because a reader who runs verification against it will get failures that look like archive tampering.

This review does not validate the architecture, approve any candidate, or authorize implementation or research execution.

## 2. What I verified, and how

I did not accept the report's reading of any cited path. For each finding I opened the source and traced the execution path myself:

| Area | Method | Result |
|---|---|---|
| `EO-01` persistence | Read `SaveChat`, `loadHistoryFile`, call structure | **Confirmed**, one detail strengthened |
| `EO-02` bypass 1 | Read `RouteResolvedEntries` and all call sites | **Confirmed** |
| `EO-02` bypass 2 | Read `buildCharacterSpecificMessage`, `getCharacterMemories`, `GenerateCharacterResponse` | **Confirmed**, scope implication understated |
| `EO-03` boundary | Read character request assembly and the transport path | **Confirmed** |
| Report identity | Recomputed all six digests, resolved each cited artifact | **Four defects** |

## 3. Provenance defects

### `P-01` — Revision 7 digest is wrong

Report records `4cd760643d57f039ecb7fc11503608c8041b8e524b1aa5258e1ed8729fb42dd1`.
Actual: `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1`.

One character at position 38: `8` recorded as `5`. Correct length, wrong value.

### `P-02` — Revision 8 digest is malformed

Report records `7f97d3594f03cc88dc2ca9148345ad7b7336cdf4ce61f6831abe2478911890`.
This is **62 hexadecimal characters**, not 64. It cannot be a SHA-256 digest and could not have been produced by `sha256sum`. Actual: `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f`.

### `P-03` — One cited artifact is mislabeled

Report cites "Closure verification record: `aa48d8a529ba0d6dd9d76639136e9e9b72d81ca670307598ec144676d0c3e1d3`."

That digest is correct, but it resolves to `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md`, the Amendment 02 verification record. No closure verification record exists under that name. The closure document is `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md`, `f828fd0469cef0508e1e2c39482875ae90665bcc07d96821df293a9746863050`, which the report also cites correctly.

### `P-04` — The completion checklist is inaccurate

§9 records `[x] All four immutable reference hashes verified` and `[x] Revision 6-9 archive identities verified for the checked files`. The second is false for revisions 7 and 8.

The important detail: **the correct values were available in a file the report itself cites.** The closure document's §2 table records revision 7 as `…b8e824b1…` and revision 8 as `…ad7b7433b6cdf…1890f` — both correct. The errors are transcription, not fabrication, and the source of truth the report depended on was accurate.

**Why this matters beyond tidiness:** every prior artifact in this chain was accepted on reproducible digests. A report whose identity block fails `sha256sum -c` produces a false tamper signal, which is worse than an absent hash, because it will be read as evidence that the archives changed.

**Required:** republish the identity block from computed values, correct the §9 checkboxes, and relabel `P-03`. No new investigation is needed.

## 4. `EO-01` — cross-store commit and recovery

**Confirmed.** `SaveChat` writes `meta.json`, `history.json`, optional `world.json`/`scenario.json`, `characters.json`, `lorebooks.json`, `memories.json`, then `config.json`, returning on the first failure. Each file is written individually; there is no spanning transaction, no commit marker, and no recovery coordinator.

**One detail stronger than reported.** `SaveChat`'s own comment states: "The append log (history.ndjson) is the true source of truth and must NOT be deleted," immediately above its write of `history.json`, which it labels "atomic write … for fallback/backup only." `loadHistoryFile` then prefers non-empty NDJSON. The consequence is sharper than "a normal save does not reconcile the preferred source": **the newer snapshot is unreachable while NDJSON is non-empty.** A save that captures newer state in `history.json` cannot be observed by the loader at all until NDJSON is empty or unreadable.

The report's other claims reproduce: corrupt NDJSON returns the prefix before the bad line, scanner errors return partial history, and missing or corrupt history yields an empty list rather than an error, so load is structurally unable to distinguish "empty" from "lost."

The report is careful here in a way worth noting: `F-01-07` and `F-01-08` are classed **Inferred**, not Observed, because the crash behavior was not executed. That is correct discipline.

## 5. `EO-02` — the two bypasses, traced independently

### 5.1 Lore routing does not filter by active world or scenario

`RouteResolvedEntries(entries, characterID, scenarioName, worldName)`:

- `case ScopeWorld:` appends unconditionally.
- `case ScopeScenario:` appends unconditionally.
- `case ScopeCharacter:` requires `entry.OwnerID == characterID`, or `shareableWithoutOwner`.
- `default:` appends unconditionally.

`worldName` and `scenarioName` are **never read in the function body**. Go does not reject unused parameters, so this compiles silently. The doc comment above the function states "Scenario entries → all characters in the active scenario," which is not what the code does.

This is not a subtle bug or a race. It is a provably inert authorization check, and the caller passes the values as though it were live:

- `app-prompt.go:922` — `RouteResolvedEntries(resolved, characterID, scenarioName, worldName)`
- `internal/backend/chat/service.go:668` and `:750` — `RouteResolvedEntries(resolvedEntries, activeCharID, s.scenarioID, s.worldID)`

**Confirmed: `F-02-02`.** The report's reading is correct and its evidence class of Observed is appropriate.

### 5.2 Frontend memory arrives as unauthorized user text

`buildCharacterSpecificMessage` composes one string containing character description, personality, background, recent system events, prior responses, inactive-character notes, the backend lore block, and `[MEMORIES]` drawn from `window.MemoryStore.getCharacterMemories(character.id)` — up to four entries.

`GenerateCharacterResponse(message string, characterID string)` then does:

```go
// The frontend already built a complete character-specific message including
// character info, lore, memories, group context, and the user message.
messages := []llama.ChatMessage{
    {Role: "user", Content: message},
}
```

The backend's own comment states the design: the frontend composed the memory. The backend therefore receives memory prose as undifferentiated user content and cannot recover scope, provenance, or revocation state from it. **Confirmed: `F-02-05`.**

### 5.3 A third finding the report understates

`getCharacterMemories` filters on `scope === 'character' && ownerKind === 'character_id' && ownerId === characterId`. There is **no conversation or chat scope in the predicate.**

This means the injection is not only an authorization bypass; it is a direct violation of the conversation-scoping requirement that `AC-C70` states — "not retrievable in another chat unless explicitly promoted." A memory recorded in conversation X is injected into the frontend-composed turn in conversation Y for the same character, with no promotion event, and the backend cannot detect that a promotion was never made.

The report records the mechanism (`F-02-05`) but frames it as an authorization gap. It is both an authorization gap and a continuity-scope violation, and the second framing connects it to `OD-04` and `AC-C70`, which the report does not make. **`M-01` below.**

### 5.4 Credit where the code is correct

`F-02-04` accurately credits fail-closed behavior: unowned character entries require `shareableWithoutOwner`, which admits only public or unlabelled legacy entries. The report is not a fault-finding exercise, and this entry is evidence of that.

## 6. `EO-03` — the transport boundary

**Confirmed.** `GenerateCharacterResponse` sends exactly one user message containing the whole frontend composition. There is no hash, no byte copy, and no independent record at the point of transmission. Image requests replace the last user message with vision content after generic assembly, so a receipt over pre-image messages would describe something other than what was sent.

The report's distinction between recording message objects at the inference interface and capturing serialized transport bytes is correct, and its treatment of the test-only capture mechanism as insufficient is proportionate.

## 7. Paths the report missed or underweighted

| ID | Item | Severity |
|---|---|---|
| `M-01` | Frontend memory injection is a conversation-scope violation under `OD-04`/`AC-C70`, not only an authorization gap. `getCharacterMemories` has no chat scope, and the backend cannot observe a missing promotion. | Material — changes which requirement the fix must satisfy |
| `M-02` | `buildCharacterSpecificMessage` also injects `[RECENT EVENTS]` (death, step-out) and inactive-character notes into the same composed string. These are chat-scoped lifecycle facts reaching the model through the unauthorizable path, which bears directly on how `AC-C71`'s payload test must be anchored. | Material for `EO-03` |
| `M-03` | `SaveChat` writes `history.json` while declaring NDJSON the source of truth; because load prefers non-empty NDJSON, a newer snapshot is unreachable rather than merely unreconciled. | Strengthens `EO-01` |

None of these contradict the report. All three increase the severity of what it found.

## 8. Candidate approaches

Each candidate set carries advantages, risks, failure modes, and detection methods, and none is selected. That is the right posture. Three gaps:

**`C-01` — the NDJSON/`history.json` duality survives every `EO-01` candidate unaddressed.** Atomicity is orthogonal to the dual-representation problem. Generation-swap (B) would need history unified inside a generation to fix it; journaling (C) would preserve the defect entirely while appearing to solve durability. The report's candidates address commit, not representation, and the representation bug is a data-loss path on its own.

**`C-02` — `F-01-10` has no candidate.** `memories_backup.json` is overwritten with `{}` before every memory write, so the backup is not a backup. No approach section mentions it. It is a data-loss vector independent of commit atomicity and should be assigned to one.

**`C-03` — latency risk is named but unassigned.** `EO-02` candidate A lists "added latency and cache-invalidation complexity" as a risk; nothing in the research plan measures it. See `R-01`.

## 9. Research plan

Phases R-1 through R-5 are sufficient to distinguish the `EO-01` candidates, and the authorization gates on each phase are stated correctly. Two gaps:

**`R-01` — hardware evaluation omits the authorization gate.** R-5 measures "candidate commit/capture mechanisms on representative constrained hardware." The authorization decision point is a candidate mechanism whose principal risk is cost, and it is not in R-5's scope. A shared gate that adds unacceptable latency on target hardware would be rejected by measurement, not by argument, and the plan currently cannot produce that measurement.

**`R-02` — no phase tests the frontend composition contract change.** `EO-02` candidate D changes the multi-character composition contract (memory IDs and owners instead of prose). R-3 tests injection behavior, but no phase estimates the cost of that contract change or verifies that the composed path can carry verifiable references. Given `M-01` and `M-02`, both bypasses live in that composer, so the plan's coverage of the highest-consequence path is thinner than its coverage of the lore path.

## 10. Conclusions versus evidence

One sentence exceeds what the evidence supports, and only in its target:

> "the evidence contradicts any assumption that Revision 9's proposed lore routing already enforces active world/scenario authority"

Revision 9 is a requirements document whose acceptance criteria are explicitly labeled unvalidated design targets. It never asserted that routing enforces anything. The correct statement is that the evidence contradicts any assumption that the **current** lore routing enforces active world/scenario authority — which is true, verified, and consequential.

The remainder of §10 is properly bounded: "does not establish," "evidence gaps, not merely documentation gaps," and "no candidate architecture is selected or validated."

The report also holds its own authorization line throughout — §0.2, §8, §9's unchecked boxes. That is consistent with how the earlier chain behaved and should be preserved.

## 11. Summary judgment

| Question asked | Answer |
|---|---|
| Is each finding supported by the cited source and execution path? | **Yes.** All three, traced independently. |
| Were relevant paths missed? | **Two underweighted** (`M-01`, `M-02`), **one strengthened** (`M-03`). No contradicting path found. |
| Do the candidate approaches address the demonstrated failure modes? | **Mostly.** `C-01`, `C-02` are unaddressed failure modes. |
| Is the research plan sufficient to distinguish viable architectures? | **Largely.** `R-01`, `R-02` are gaps. |
| Does any conclusion exceed the evidence? | **One sentence**, misattributing the contradiction to Revision 9 rather than current code. |

The three obligations were worth running. Each produced a real defect, and the two authorization bypasses are more serious than "documentation deficiency": one is an inert permission check that callers believe is active, and the other routes conversation-scoped memory into other conversations in a way the backend structurally cannot see.

## 12. Authorization boundary

No implementation, migration, runtime investigation, model load, benchmark, network test, Android pairing, or research execution is authorized by this review. `OD-01`–`OD-04` remain settled. `O-01`–`O-10` remain open. The report's findings are confirmed as findings; no candidate is validated, and the architecture remains unapproved.
