# R-07 — Reconciled Evidence Register
## Howl Chat / Amendment 02 Closure Follow-Through
### Date: 2026-11-11
### Status: Reconciles original R-07 with Muse independent review; no architecture selected; no implementation authorized

## 0. Inputs

- Original R-07: `docs/research/R-07-EVIDENCE-REPORT-2026-11-11.md`
  - SHA-256: `e3860e3d5d6ca71724d31d0021897819c968afc1b8e1aad09fc011f8d9a67134`
  - Size: `46,107 B / 876 lines`
- R-07 erratum and supplement: `docs/research/R-07-EVIDENCE-REPORT-2026-11-11-ERRATUM-AND-SUPPLEMENT.md`
- Independent review: `docs/review/REVIEW-R-07-EVIDENCE-MUSE-INDEPENDENT-2026-11-11.md`
  - SHA-256: `4bdb399cfcb56c75e32d3103a97fd57b06db8895ea3fd8dfef16feba5a98d792`
  - Size: `15,423 B / 180 lines`

This register reconciles findings; it does not re-execute tests, validate architecture, or select a mechanism. — **Observed**

## 1. Status vocabulary

- **Confirmed:** independently reproduced from source by the review.
- **Confirmed-strengthened:** independently reproduced and found more severe than the original wording.
- **Confirmed-understated:** independently reproduced but assigned too narrow a requirement framing.
- **Correction required:** identity/provenance or conclusion text must be corrected.
- **Open:** needs bounded research, source inspection, or owner decision.
- **Unestablished:** insufficient evidence.

## 2. EO-01 findings

| ID | Finding | Origin | Class | Status | Evidence | Next action |
|---|---|---|---|---|---|---|
| F-01-01 | Live chat files are written sequentially, not transactionally. | R-07 | Observed | Confirmed | `SaveChat` at `internal/backend/chat/persistence.go:191-251` | Preserve; use in candidate evaluation |
| F-01-02 | Individual JSON writes use temp-file rename, but the multi-file operation has no atomicity. | R-07 | Observed | Confirmed | `saveJSON` at `:569-581`; sequential call order in `SaveChat` | Preserve |
| F-01-03 | NDJSON is preferred history source; normal `SaveChat` only rewrites `history.json`. | R-07 | Observed | Confirmed-strengthened (`M-03`) | `loadHistoryFile` at `:583-629`; `SaveChat` at `:208-213` | Resolve representation authority before commit design |
| F-01-04 | `chat.Service`/memory state and live AppData chat state use different roots and formats. | R-07 | Observed | Confirmed | `persistence.go:61-75`; `memory/service.go:345-371`; `chat/service.go:147-244` | Enumerate roots and migration authority |
| F-01-05 | No journal, commit marker, operation ID, revision marker, or cross-store recovery coordinator was found. | R-07 | Observed | Confirmed | `SaveChat`; `ChatMeta`/`ChatHistory` | Candidate evaluation |
| F-01-06 | Best-effort load can skip corrupt/optional files and truncate NDJSON at the first bad line. | R-07 | Observed | Confirmed | `loadHistoryFile`; `LoadChat` optional loads | Fault-injection research |
| F-01-07 | Crash between files can leave newer metadata with older context, or persisted but inactive chat files. | R-07 | Inferred | Confirmed as reasoning, not executed | Sequential writes; config last | Phase R-2 fault tests |
| F-01-08 | Recovery cannot deterministically finish or reject a partial multi-file operation from current metadata alone. | R-07 | Inferred | Confirmed as reasoning, not executed | No commit/recovery marker | Phase R-2 fault tests |
| F-01-09 | Clearing world/scenario memory objects does not delete stale world/scenario files. | R-07 | Observed | Not independently challenged | `SaveChat` writes only non-nil objects | Isolated deletion test |
| F-01-10 | `memories_backup.json` is overwritten with `{}` before every memory write. | R-07 | Observed | Confirmed; no candidate assigned (`C-02`) | `memory/service.go:459-461` | Assign candidate and test |
| M-03 | A newer `history.json` snapshot is unreachable while NDJSON is nonempty. | Muse | Observed | Confirmed-strengthened | `loadHistoryFile` preference; `SaveChat` snapshot write | Separate representation fix |
| C-01 | NDJSON/history duality survives every EO-01 candidate unaddressed. | Muse | Inferred | Open | R-07 §1.8 candidates | Add representation workstream |
| C-02 | Backup overwrite has no candidate approach. | Muse | Observed | Open | `F-01-10` | Add backup candidate |

## 3. EO-02 findings

| ID | Finding | Origin | Class | Status | Evidence | Next action |
|---|---|---|---|---|---|---|
| F-02-01 | No single shared authorization decision point across lore, frontend memory, service memory, and RAG. | R-07 | Observed | Confirmed | Prompt paths and providers | Candidate evaluation |
| F-02-02 | World/scenario lore routed without active world/scenario comparison. | R-07 | Observed | Confirmed | `lorebook.go:344-376`; unused `scenarioName`, `worldName` | Fix requirement before implementation |
| F-02-03 | Character lore requires stable owner-ID match or explicit shareability. | R-07 | Observed | Confirmed | `lorebook.go:358-370` | Preserve correct behavior |
| F-02-04 | Hidden/stale lore excluded; unowned character lore fail-closed. | R-07 | Observed | Confirmed | `canConsider`; `shareableWithoutOwner` | Preserve |
| F-02-05 | Frontend memory injection bypasses backend authorization as composed user text. | R-07 | Observed | Confirmed-strengthened (`M-01`) | `chat-stream.js:100-180`; `GenerateCharacterResponse` | Re-scope as authorization + continuity |
| F-02-06 | Backend memory records lack durable source-message provenance. | R-07 | Observed | Not independently challenged | `types/memory.go:21-38`; DTO at `app-memory.go:30-52` | Source inspection |
| F-02-07 | RAG retrieval is similarity/budget-scoped, not audience-authorized. | R-07 | Observed | Confirmed by review scope; live reachability open | `rag/service.go:119-169` | Reachability test |
| F-02-08 | Prompt-pool authority depends on last sync/merge path. | R-07 | Observed | Not independently challenged | `mergeCanonicalLore`; `GetAllLorebookEntries` | Source trace |
| F-02-09 | Scope changes/promotion/revocation effects on caches and derived artifacts not centrally enforced. | R-07 | Inferred | Confirmed as reasoning gap | No central gate; snapshot/cache paths | Phase R-3 |
| F-02-10 | Live reachability of secondary `chat.Service` prompt/RAG path is not established. | R-07 | Unestablished | Open | `App.buildPromptMessages` vs `chat.Service.buildPrompt` | Static call graph + isolated test |
| M-01 | Frontend character memory crosses conversation boundaries without promotion; violates `OD-04`/`AC-C70`. | Muse | Observed | Confirmed; materially understated in R-07 | `memory-store.js:203-208`; no chat scope | New continuity requirement |
| M-02 | Chat-scoped lifecycle facts (`[RECENT EVENTS]`, inactive notes) travel through composed frontend text. | Muse | Observed | Confirmed | `chat-stream.js` composer | EO-03 payload anchoring |
| C-03 / R-01 | Authorization-gate latency is named as a risk but not measured by the research plan. | Muse | Inferred | Open | R-07 §2.8 candidate A; Phase R-5 | Add measurement gate |
| R-02 | Frontend composition contract change is not tested by any research phase. | Muse | Inferred | Open | EO-02 candidate D; Phase R-3 | Add contract-cost phase |

## 4. EO-03 findings

| ID | Finding | Origin | Class | Status | Evidence | Next action |
|---|---|---|---|---|---|---|
| F-03-01 | Inspected boundary is HTTP POST `/v1/chat/completions` with JSON-marshaled request. | R-07 | Observed | Confirmed | `http_client.go:144-180`; `:594-680` | Backend inventory |
| F-03-02 | Production retains no hash, receipt, or byte copy of transmitted request. | R-07 | Observed | Confirmed | `GenerateChatCompletionStreaming` | Phase R-4 |
| F-03-03 | Existing capture records message objects, not serialized transport bytes. | R-07 | Observed | Confirmed | `inference_capture.go:47-77` | Phase R-4 |
| F-03-04 | Existing capture is test-only and absent from production builds. | R-07 | Observed | Confirmed | `inference_capture_off.go:1-6` | Preserve |
| F-03-05 | Image requests transform last user message after generic assembly. | R-07 | Observed | Confirmed | `http_client.go:655-680` | Modality-specific capture |
| F-03-06 | Character requests collapse multiple frontend sources into one user string. | R-07 | Observed | Confirmed | `GenerateCharacterResponse`; composer | Provenance design |
| F-03-07 | Exact transport-byte identity is not independently established. | R-07 | Inferred | Confirmed | No transport receipt | Phase R-4 |
| F-03-08 | Other production model backends/alternate boundaries are not established. | R-07 | Unestablished | Open | Inspected path only | Backend inventory |

## 5. Provenance findings

| ID | Finding | Origin | Class | Status | Correction |
|---|---|---|---|---|---|
| P-01 | Revision 7 digest wrong in R-07. | Muse | Observed | Correction required | Use `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` |
| P-02 | Revision 8 digest malformed (62 hex chars). | Muse | Observed | Correction required | Use `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` |
| P-03 | `aa48…` artifact mislabeled as closure verification record. | Muse | Observed | Correction required | `aa48…` is the verification record; closure record is `f828…` |
| P-04 | R-07 completion checklist inaccurately marked hash verification passed. | Muse | Observed | Correction required | Replace with corrected checklist in erratum |

## 6. Reconciled conclusions

1. **EO-01 is confirmed.** The current implementation has no cross-store commit or deterministic recovery, and the NDJSON/history snapshot duality is an independent data-loss path. — **Observed/Inferred**
2. **EO-02 is confirmed and broader than authorization.** The active-world/scenario filter is inert, and frontend memory injection is both an authorization bypass and a conversation-scope violation. — **Observed**
3. **EO-03 is confirmed.** Message-object capture is not transport-byte evidence; image transformation makes the boundary modality-specific. — **Observed**
4. **R-07’s substantive findings survive review.** Its identity block, two archive digests, artifact label, checklist, and one conclusion sentence require correction. — **Observed**
5. **No candidate is selected.** Candidate approaches still need representation, backup, latency, and frontend-composition coverage. — **Observed/Inferred**
6. **The review does not validate architecture or authorize implementation.** — **Observed**

## 7. Open questions requiring bounded work

| ID | Question | Required work | Authorization |
|---|---|---|---|
| Q-01 | Can frontend composition carry memory IDs/owners/conversation IDs/promotion state? | Contract design + isolated test | Implementation planning later |
| Q-02 | What active-world/scenario filter is correct, and where is authority read? | Source trace + requirements decision | Read-only then design |
| Q-03 | Which history representation is authoritative, and how is newer snapshot reconciled? | Design + fault tests | Runtime research |
| Q-04 | What backup guarantee should `memories_backup.json` provide? | Design + fault tests | Runtime research |
| Q-05 | What capture proves bytes for text/image/generic/character paths? | Transport capture test | Runtime/network/model authorization |
| Q-06 | Which five source files are modified vs HEAD, and were they pre-existing? | Git baseline audit | Read-only |
| Q-07 | Is secondary `chat.Service` prompt/RAG live? | Static call graph + isolated path test | Runtime research |
| Q-08 | What latency budget applies to a shared authorization gate? | Benchmark on target hardware | Benchmark authorization |
| Q-09 | What is the cost of the frontend composition contract change? | Isolated UX/compat test | Runtime research |
| Q-10 | How should chat-scoped lifecycle facts be anchored for `AC-C71`? | Payload/provenance design | Design + test |

## 8. Baseline caveat

The owner reports five modified source files relative to Git HEAD. The review states those modifications predate its read-only work. The review artifact available here does not enumerate the five files. This register therefore does **not** claim a clean working-tree baseline. — **Unestablished**

## 9. Status

R-07 is preserved but is not a clean identity record without its erratum. The findings are confirmed; the evidence register is reconciled; architecture remains unvalidated; no implementation is authorized.

**Requirements baseline closed. Architecture unvalidated. Evidence reconciled. Next: bounded questions and owner decisions.**
