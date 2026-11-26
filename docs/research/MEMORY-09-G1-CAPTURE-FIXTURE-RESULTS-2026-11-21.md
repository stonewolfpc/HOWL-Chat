# MEMORY-09 — G1 Capture Fixture Results

## The data-model contract passes; production capture wiring is still unproven

**Status:** **G1 fixture contract passed.** The test-only harness runs all F1–F10 cases and the character-tell case
without changing production behavior. The current production capture path still does not receive the envelope; that
remains a separate boundary.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Gate | **G1 — capture completeness** |
| Fixture result | **PASS** — 11 cases, 3 focused assertions |
| Harness | test-only; `internal/backend/memory/capture_contract_test.go` |
| Production files changed | **0** |
| Production API changed | **0** |
| L-2 / Decision A | **unchanged** |
| BGE / G2 | **not started** |
| Map | **118** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Commands and result

**G1R-01 (Observed — focused run).**

```
go test ./internal/backend/memory -run 'TestG1' -v
```

Result: **PASS**. `TestG1CaptureContract` passed all 11 subtests; `TestG1ExplicitNoneIsNotMissing`,
`TestG1WorldSpeechIsNotAmbient` and `TestG1CharacterTellRemainsUnclassified` passed.

**G1R-02 (Observed — package run).**

```
go test ./internal/backend/memory
```

Result: `ok  howl-chat/internal/backend/memory  0.804s`.

**G1R-03 (Observed — no production change).** The only new file is the test harness. No production type, function,
signature, storage path or binding was changed to make the fixture pass.

---

## 2. Case results

**G1R-04 (Observed — F1–F10 and character-tell).**

| Case | Result | Observed contract outcome |
|---|---|---|
| **F1** stable character message | **PASS** | complete; source `char-a`; eligible `experiential` / `character_specific` |
| **F2** user message | **PASS** | complete; source character **explicit none**; class eligibility deferred |
| **F3** narrator message | **PASS** | complete; no authored ref; ambient/restricted not eligible |
| **F4** empty `CharacterID` | **PASS** | **capture-incomplete**; missing `source_character_id` |
| **F5** duplicate display names | **PASS** | complete; stable `char-b` selected; no name matching |
| **F6** missing attribution | **PASS** | **capture-incomplete** for experience/belief; missing `attribution` |
| **F7** authored world definition | **PASS** | complete; `ambient` / `restricted` eligible |
| **F8** world-related character speech | **PASS** | complete; eligible `experiential` / `character_specific`; **not** ambient/restricted |
| **F9** legacy record without source | **PASS** | **capture-incomplete**; missing `source_message_id`, `speaker_role` |
| **F10** missing world/scenario | **PASS** | **capture-incomplete**; missing `world_id`, `scenario_id` |
| **character-tell** | **PASS** | complete; attribution `told`; no class assigned; non-class processing permitted |

---

## 3. Focused assertions

**G1R-05 (Observed — explicit none vs missing).** `TestG1ExplicitNoneIsNotMissing` passed: F2 carries a non-nil
explicit-none source character and is complete; F4 carries a nil source character and is incomplete. A user-originated
statement is not treated as a malformed character memory.

**G1R-06 (Observed — world speech is not authored truth).** `TestG1WorldSpeechIsNotAmbient` passed: F8 is eligible only
for `experiential` / `character_specific`; `ambient` and `restricted` are excluded without an authored asset reference.

**G1R-07 (Observed — character-tell remains unclassified).** `TestG1CharacterTellRemainsUnclassified` passed: capture
is complete, attribution is `told`, no class is assigned, and non-class processing remains permitted.

**G1R-08 (Observed — no reference class encoded).** Every fixture case has an empty `referenceClass`; the test fails if
one is populated. `EligibleClasses` is a permission set, not an expected label.

**G1R-09 (Observed — no class assigned).** Every envelope carries an empty assigned class; the harness does not test a
classifier against its own assumptions.

---

## 4. What this proves

**G1R-10 (Calculated — the data model can express the contract).** Existing `ChatHistory`, `ChatMeta` and `ChatData`
fields can carry source message ID, role, stable source character (or explicit none), world/scenario IDs and content.
The fixture models the missing attribution and authored-asset fields explicitly.

**G1R-11 (Calculated — the states are distinct).** Capture-incomplete is a separate outcome from classification-
unresolved; F4/F6/F9/F10 are incomplete, not unresolved.

**G1R-12 (Calculated — the fixture is source-grounded).** The fixture uses existing app types read-only and a test-only
envelope builder. It does not call the production memory API or persist anything.

---

## 5. What this does not prove

**G1R-13 (Observed — production wiring is not established).** The current `ProcessMessageForEvents` API still receives
only message text, history strings and `{ID, Name}` identities; the frontend still passes empty arrays and has no App
binding. The fixture proves the data model can express the contract, not that the production path passes it.

**G1R-14 (Calculated — G2 is still blocked on production capture).** A classifier experiment needs the production
capture path to supply the envelope. G2 has not started; no class labels, rules, BGE or corpus were used.

**G1R-15 (Calculated — G1's pass is scoped).** G1 passes as a **fixture contract demonstration**. It does not
authorize a production API change, L-2 amendment, BGE fetch or G2 run.

---

## 6. Files and verification

**G1R-16 (Observed — changed files).**

| File | Change |
|---|---|
| `internal/backend/memory/capture_contract_test.go` | **new** test-only harness |
| production files | **none** |

**G1R-17 (Observed — repository state).** The same **6** pre-existing modified files remain untouched; the new test
file is untracked. Map **118**; canon hashes unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`).

**G1R-18 (Observed — test scope).** Focused G1 tests and the full memory package passed. No browser, UI, model or
network test was run.

---

## 7. Next boundary

**G1R-19 (Calculated — the next problem).** The production capture boundary is the next problem: pass the envelope
from `ChatHistory`/`ChatMeta`/`ChatData` into the memory service, or explicitly reject formation as
capture-incomplete. That is a production API change and requires separate authorization.

**G1R-20 (Observed — not started).** No production API change, L-2/Decision A amendment, BGE fetch, G2 classifier or
corpus work has started.

---

## 8. Traceability

**G1R-21 (Observed).**

| What | Where |
|---|---|
| Fixture contract | GFC-01…31 |
| G1 audit | CCG-01…33 |
| Capture-incomplete vs unresolved | GFC-01/16/19; G1R-11 |
| F1–F10 | GFC-13; G1R-04 |
| Character-tell | GFC-17…19; G1R-07 |
| Ambient authored-only | K-R8; G1R-06 |
| Production gap | CCG-04…08; G1R-13 |
| No production change | G1R-03/16 |
| Map | **118** |

---

## 9. Record status

**G1R-22 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Production code/schema edited | **0** |
| Requirements changed | **0** — map **118** |
| G1 fixture | **PASS** |
| G1 production wiring | **not established** |
| L-2 / Decision A amended | **0** |
| BGE / G2 started | **0** |
| Tests executed | **G1 fixture + memory package pass** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **G1 result** — no canonical capture |

**End of results.** The fixture-only harness passes all F1–F10 and character-tell cases without a production change
(G1R-01…09). It proves the existing app data model can express the capture contract and that capture-incomplete is
distinct from unresolved (G1R-10…12). It does **not** prove the production memory API receives the envelope; that
remains the next boundary and requires separate authorization (G1R-13…15/19). **No L-2 amendment, BGE fetch, G2 start
or production change follows** (G1R-20/22).
