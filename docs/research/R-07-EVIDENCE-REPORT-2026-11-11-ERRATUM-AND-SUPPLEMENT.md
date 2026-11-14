# R-07 — Erratum and Supplemental Findings Record
## Howl Chat / Amendment 02 Closure Follow-Through
### Date: 2026-11-11
### Status: Additive correction; original R-07 remains immutable; no architecture selected; no implementation authorized

## 0. Purpose and authority

This record corrects the identity block and completion checklist of:

- `docs/research/R-07-EVIDENCE-REPORT-2026-11-11.md`
- Original SHA-256: `e3860e3d5d6ca71724d31d0021897819c968afc1b8e1aad09fc011f8d9a67134`
- Original size: `46,107 B / 876 lines`

The original report is **not edited or replaced** by this record. This erratum is separately hashed and must be read with the original. — **Observed**

Independent review preserved unchanged at:

- `docs/review/REVIEW-R-07-EVIDENCE-MUSE-INDEPENDENT-2026-11-11.md`
- SHA-256: `4bdb399cfcb56c75e32d3103a97fd57b06db8895ea3fd8dfef16feba5a98d792`
- Size: `15,423 B / 180 lines`
- Reviewer: Muse (independent instance)
- Review status: complete; read-only; no document or source file modified. — **Observed**

This erratum is not a Revision 10, does not amend Amendment 02, does not approve any architecture, and does not authorize implementation or research execution. — **Observed**

## 1. Identity corrections

### P-01 — Revision 7 digest

| Field | Value |
|---|---|
| Recorded in R-07 | `4cd760643d57f039ecb7fc11503608c8041b8e524b1aa5258e1ed8729fb42dd1` |
| Correct value | `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` |
| Error | One character at position 38: `8` recorded as `5` |
| Classification | **Correction required — Observed** |

Correct file: `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-7-2026-11-11.md`. — **Observed**

### P-02 — Revision 8 digest

| Field | Value |
|---|---|
| Recorded in R-07 | `7f97d3594f03cc88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890` |
| Correct value | `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` |
| Error | Recorded value is 62 hexadecimal characters and cannot be a SHA-256 digest |
| Classification | **Correction required — Observed** |

Correct file: `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md`. — **Observed**

### P-03 — Artifact label correction

R-07 recorded:

> Closure verification record: `aa48d8a529ba0d6dd9d76639136e9e9b72d81ca670307598ec144676d0c3e1d3`

Correction:

| Artifact | Correct label | Path | SHA-256 |
|---|---|---|---|
| Amendment 02 revision-cycle closure | Closure record | `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md` | `f828fd0469cef0508e1e2c39482875ae90665bcc07d96821df293a9746863050` |
| Amendment 02 verification record | Verification record | `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md` | `aa48d8a529ba0d6dd9d76639136e9e9b72d81ca670307598ec144676d0c3e1d3` |

The `aa48…` digest is correct for the verification record; it is not a “closure verification record.” No artifact with that name was found. — **Observed**

### P-04 — Completion-checklist correction

R-07 §9 marked these as passed:

- `[x] All four immutable reference hashes verified`
- `[x] Revision 6-9 archive identities verified for the checked files`

The second statement was false for Revisions 7 and 8 because of P-01 and P-02. The first statement remains supportable for the four required inputs when their paths are labeled correctly, but the original checklist did not distinguish the closure record from the verification record. — **Observed**

Corrected status:

- [x] Four required input digests reproduce when paired with the correct paths.
- [ ] Original R-07 identity block reproduces without correction.
- [x] Revision 6 digest verified.
- [ ] Revision 7 digest in original R-07 verified.
- [ ] Revision 8 digest in original R-07 verified.
- [x] Correct Revision 7 and Revision 8 digests established by independent recomputation.
- [x] No Revision 10 found.
- [x] Original R-07 preserved byte-for-byte.
- [ ] Full prior archive/sidecar enumeration complete.
- [ ] Independent review reconciliation complete.

## 2. Independent reproduction of the three obligations

### EO-01 — confirmed

Muse independently read `SaveChat`, `loadHistoryFile`, and the call structure. It confirmed:

- sequential per-file saves
- no spanning transaction
- no commit marker
- no recovery coordinator
- corrupt NDJSON returns the prefix before the bad line
- scanner errors return partial history
- missing/corrupt history yields empty rather than an error
- load cannot distinguish “empty” from “lost”

Classification: **Confirmed — Observed**. — **Observed**

Muse strengthened the report’s history finding as `M-03`: because load prefers nonempty NDJSON, a newer `history.json` snapshot is not merely unreconciled; it is unreachable while NDJSON remains nonempty. Classification: **Confirmed and strengthened — Observed**.

### EO-02 — both bypasses confirmed

Muse independently traced both:

1. `RouteResolvedEntries` appends world and scenario entries unconditionally and never reads `scenarioName` or `worldName` in the function body. This is an inert authorization check that callers pass values to as though it were active. Classification: **Confirmed — Observed**.
2. `buildCharacterSpecificMessage` composes character context, lore, memories, recent events, prior responses, and inactive-character notes into one string; `GenerateCharacterResponse` sends that string as one user message. The backend cannot recover memory scope, provenance, or revocation state. Classification: **Confirmed — Observed**.

Muse found the second problem is stronger than R-07 stated: `getCharacterMemories` has no conversation/chat scope, so same-character memories can cross conversation boundaries without promotion. Classification: **Confirmed and materially understated — Observed**.

### EO-03 — confirmed

Muse confirmed:

- character requests send exactly one user message containing the frontend composition
- no hash or byte copy exists at transmission
- image requests transform the last user message after generic assembly
- recording message objects is not proof of transmitted transport bytes
- the test-only capture mechanism is insufficient as a production receipt

Classification: **Confirmed — Observed**.

### One conclusion requiring correction

R-07 §10 said:

> “the evidence contradicts any assumption that Revision 9’s proposed lore routing already enforces active world/scenario authority”

Correction:

> The evidence contradicts any assumption that the **current** lore routing enforces active world/scenario authority.

Revision 9 is a requirements document with unvalidated design targets; it did not assert that its proposed routing was already implemented. Classification: **Correction required — Inferred**.

## 3. Additional findings from independent review

### M-01 — Frontend memory is a conversation-scope violation, not only an authorization gap

`getCharacterMemories` filters only on:

- `scope === 'character'`
- `ownerKind === 'character_id'`
- `ownerId === characterId`

It has no conversation/chat scope. A memory recorded in conversation X is injected into conversation Y for the same character without a promotion event. The backend cannot observe that the promotion never occurred.

This is both an authorization gap and a continuity-scope violation under `OD-04` / `AC-C70`. Classification: **Confirmed — Observed**.

### M-02 — Lifecycle material also travels through the composed frontend message

`buildCharacterSpecificMessage` also injects `[RECENT EVENTS]` such as death and step-out, plus inactive-character notes, into the same composed string. These are chat-scoped lifecycle facts reaching the model through the path that cannot be independently authorized or attributed.

This matters for anchoring `AC-C71` payload tests. Classification: **Confirmed — Observed**.

### M-03 — Newer history snapshot is unreachable while NDJSON is nonempty

`SaveChat` writes `history.json` and calls it a fallback/backup, while `loadHistoryFile` prefers nonempty `history.ndjson`. A newer snapshot captured only in `history.json` cannot be observed until NDJSON is empty or unreadable.

Classification: **Confirmed and strengthened — Observed**.

## 4. Review gaps and omissions now tracked

### C-01 — NDJSON/history duality is not addressed by any EO-01 candidate

R-07’s candidate approaches address commit atomicity, not representation. Generation-swap and journaling could still preserve the dual-representation defect; journaling in particular could appear to solve durability while leaving the data-loss path intact. Classification: **Review finding — Inferred**.

### C-02 — `memories_backup.json` has no candidate

`MemoryStorage.saveMemories` overwrites `memories_backup.json` with `{}` before each write. It is therefore not a backup. R-07 recorded this as `F-01-10` but assigned no candidate approach. It is a data-loss vector independent of commit atomicity. Classification: **Review finding — Observed**.

### C-03 / R-01 — Authorization-gate latency is unmeasured

EO-02 candidate A names added latency and cache-invalidation complexity as a risk, but the research plan does not measure the authorization decision point on target hardware. A shared gate that adds unacceptable latency would need to be rejected by measurement, not argument. Classification: **Review finding — Inferred**.

### R-02 — Frontend composition contract change is untested by the plan

EO-02 candidate D changes the multi-character composition contract from prose to verifiable memory IDs/owners. The research plan tests injection behavior but does not estimate the contract-change cost or verify that the composed path can carry verifiable references. Both demonstrated bypasses live in that composer, so this is the highest-consequence path with the thinnest plan coverage. Classification: **Review finding — Inferred**.

## 5. Supplemental evidence questions

These are questions, not findings:

- Q-01: Can the frontend composer carry memory IDs, owner IDs, conversation IDs, and promotion state without changing user-visible behavior?
- Q-02: What exact filter should active-world/active-scenario lore use, and where is current authority read?
- Q-03: Which history representation is authoritative, and how is a newer snapshot reconciled with NDJSON?
- Q-04: What recovery guarantee should `memories_backup.json` provide, if any?
- Q-05: What transport capture is sufficient for text, image, generic, and frontend-composed character paths?
- Q-06: Which five source files are modified relative to Git HEAD, and were those modifications present before R-07 and the review?
- Q-07: Is the secondary `chat.Service` prompt/RAG path live in production?
- Q-08: What latency budget applies to a shared authorization gate on target hardware?
- Q-09: What is the cost of changing the frontend composition contract?
- Q-10: How should chat-scoped lifecycle facts be represented so `AC-C71` payload tests can anchor them?

## 6. Baseline caveat

The owner reports five modified source files relative to Git HEAD. The independent review states that its own work was read-only and that those modifications predate its work.

This erratum does **not** independently establish a clean baseline. The exact five-file list is not enumerated in the review artifact available here. Classification: **Unestablished**.

Any future evidence artifact must state whether it was produced against the modified working tree or against a clean commit, and must not treat “reviewer did not modify files” as proof that the tree matched HEAD.

## 7. Preserved artifacts and status

| Artifact | SHA-256 | Status |
|---|---|---|
| R-07 original report | `e3860e3d5d6ca71724d31d0021897819c968afc1b8e1aad09fc011f8d9a67134` | Preserved, not edited |
| Muse independent review | `4bdb399cfcb56c75e32d3103a97fd57b06db8895ea3fd8dfef16feba5a98d792` | Preserved, not edited |
| Revision 9 | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` | Preserved |
| Closure record | `f828fd0469cef0508e1e2c39482875ae90665bcc07d96821df293a9746863050` | Preserved |
| Verification record | `aa48d8a529ba0d6dd9d76639136e9e9b72d81ca670307598ec144676d0c3e1d3` | Preserved |
| Revision 6 | `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` | Preserved |
| Revision 7 | `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` | Preserved |
| Revision 8 | `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` | Preserved |

The original R-07 must not be cited as a clean identity record without this erratum. Its substantive findings remain supported; its provenance block and one conclusion sentence require correction. — **Observed**

## 8. Corrected completion status

- [x] Three evidence obligations independently reproduced.
- [x] EO-01 confirmed; representation defect strengthened.
- [x] EO-02 both bypasses confirmed; conversation-scope violation added.
- [x] EO-03 transport-receipt gap confirmed.
- [x] Original R-07 preserved byte-for-byte.
- [x] Independent review preserved byte-for-byte.
- [x] Identity errors identified with correct replacement values.
- [x] Additional findings M-01–M-03 recorded.
- [x] Omissions C-01–C-03 and R-01–R-02 recorded.
- [ ] R-07 and review reconciled into one evidence register.
- [ ] Exact five-file baseline enumerated.
- [ ] Any architecture decision prepared.

## 9. Authorization boundary

This erratum authorizes nothing beyond documentation. It does not authorize source changes, schema or migration changes, runtime experiments, model loading, inference tests, benchmarks, network or packet capture, Android pairing, user-data inspection beyond existing read-only authorization, database mutation, new dependencies, implementation, deployment, or research execution.

**Requirements baseline closed. Architecture unvalidated. Evidence corrected and reconciled next.**
