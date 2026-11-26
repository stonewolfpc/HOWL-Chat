# Howl Foundation — Contract Paper Register

## Register `HF-CONTRACTS-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract-paper index and dependency register |
| Authority packet | `HOWL-FOUNDATION-CONTRACT-PACK-DEEPSEEK-INTAKE-2026-11-20.md` (`HF-CONTRACT-INTAKE-2026-11-20-01`) |
| Owner decisions | D1, D2, D3, D5, D6 selected; D7 requirement clarified; D4 owner direction established; D8 delegated |
| Status | **Review drafts, review-complete.** FC-1–FC-6 accepted as reviewed proposed drafts (2026-11-21); none is approved. FC-5 carries an owner-approved additive amendment (revision 6, 2026-11-21) applied after that acceptance; the amendment approves no contract and changes no other paper's status |
| Review baseline | FC-1–FC-6 frozen at the six pinned hashes in §6 |
| M0A | **Separate and unauthorized.** No M0A work is included here |
| Code/schema/API/migration/test/model/process changes | **0** |
| Implementation authorization | **None** |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Blocked** · **Observed** · **Unestablished**.

---

## 1. Papers

| Paper | File | Gate | Status | Dependency |
|---|---|---|---|---|
| **FC-1** | [Conversation State Contract](FC-1-CONVERSATION-STATE-CONTRACT-2026-11-21.md) | M1 | Accepted review draft | D1, D2, D3, D7 |
| **FC-2** | [Composition and Context Contract](FC-2-COMPOSITION-AND-CONTEXT-CONTRACT-2026-11-21.md) | M2 | Accepted review draft | D3, D4; FC-1 |
| **FC-3** | [Multi-Character Isolation Contract](FC-3-MULTI-CHARACTER-ISOLATION-CONTRACT-2026-11-21.md) | M3 | Accepted review draft | D5, D7; FC-1, FC-2 |
| **FC-4** | [Content Ownership and Synchronization Contract](FC-4-CONTENT-OWNERSHIP-AND-SYNCHRONIZATION-CONTRACT-2026-11-21.md) | M0B/M5 | Accepted review draft | D3; FC-1 |
| **FC-5** | [Capability, Resource, and Handoff Contract](FC-5-CAPABILITY-RESOURCE-AND-HANDOFF-CONTRACT-2026-11-21.md) | M4 | Accepted review draft | D6; HRA-DB |
| **FC-6** | [Evidence and Acceptance Contract](FC-6-EVIDENCE-AND-ACCEPTANCE-CONTRACT-2026-11-21.md) | All gates | Accepted review draft | FC-1–FC-5; D7 acceptance |
| **FC-7** | **Not filed** | M8 | **Blocked by D8** | D8 delegated decision record |

**FC-7 is deliberately absent.** D8 delegates the Howl ↔ MEMORY-09 lifecycle-fact decision to the GPT and DeepSeek agent/researcher process. Filing FC-7 now would either pre-empt that decision or create a placeholder. The contract register records the block rather than manufacturing a paper.

---

## 2. Dependency order

1. **FC-1** fixes conversation identity and event/state semantics, records D2's authored-object retention direction, and defines snapshot references. Conversation-event and state-history retention remain unresolved.
2. **FC-2** composes a reproducible prompt plan from FC-1 state and FC-4 canonical references.
3. **FC-3** applies the D5 hard-flush policy and preserves D7 voice identity across speaker changes.
4. **FC-4** defines canonical authored-object authority and revision binding used by FC-1 and FC-2.
5. **FC-5** defines capability admission, reservation and verified handoff under HRA-DB and D6. An owner-approved additive amendment (revision 6, 2026-11-21) adds two configuration records — the configuration decision record (FC5-27–FC5-30) and the launch construction record (FC5-31) — together with the identifier note, the §4 lead-in and the FC5-26a boundary pointer. It approves no contract, resolves no `DQ-HRA-*` question and settles no configuration-retention question.
6. **FC-6** defines receipts, claim levels, fault evidence and the D7 end-to-end acceptance test.
7. **FC-7** waits for the D8 delegated decision record.

FC-2 and FC-4 may be reviewed in parallel with FC-1, but their references to state identity and snapshots must not diverge from FC-1. FC-3 depends on both state and prompt composition. FC-6 cannot be finalized until the other papers expose their acceptance surfaces.

---

## 3. Owner decisions reflected

| Decision | Contract consequence |
|---|---|
| **D1** | Durable append-only event ledger; state/branch views are derived and reconstructible. |
| **D2** | Retain user data by default: authored worlds, scenarios, characters, lore, images and other library objects remain stored unless manually deleted from inside Howl. Orphanhood is a relationship or ownership condition, not a deletion authorization; no cascading deletion of authored objects is authorized. No automatic orphan cleanup is authorized as a general background sweep; the earlier blanket shorthand is narrowed by FC-1's directed cleanup, under which a parent's own original memory records are removed as part of that parent's explicit deletion (FC1-22e) while independently owned memory copies survive the source parent's deletion (FC1-22g, FC1-22i). FC-1 records three directed lifecycles: authored content is deleted only through explicit user deletion; memory originals are cleaned up with their owning parent's explicit deletion; conversation reproducibility snapshots are deleted with their owning conversation (FC1-22c). `orphaned` remains undefined for legacy, unknown-origin and interrupted-lifecycle records, destination ownership of a copy remains unresolved (FC1-22g, `FC1-MEMORY-COPY-DESTINATION`), and D2 does not establish conversation-event or state-history retention. |
| **D3** | Canonical library owns authored objects and revisions; conversations hold immutable references and reproducibility snapshots; browser storage becomes compatibility/import only after a bounded transition. |
| **D4** | Prompt construction must be preflight, model-compatible, reproducible and fail-closed. Technology selection remains research, not assumed. |
| **D5** | Hard flush on speaker change; keyed reuse only after adapter-specific evidence. |
| **D6** | One supported heavy-work instance initially, temporary until B or C passes recovery and concurrency acceptance. This is not a one-resident-model rule and not cross-instance reservation safety. |
| **D7** | Local-only TTS; unique, user-adjustable per-character voices preserved across switching and restoration; actual multi-character acceptance required. |
| **D8** | Delegated decision process; disagreement returns to the owner; no option selected in these papers. |

---

## 4. Cross-paper safeguards

1. **D1 + D3:** canonical references and reproducibility snapshots must remain meaningful across revisions and restoration. A snapshot is evidence, not authored-object authority.
2. **D4:** tokenizer selection is not the whole contract. The complete model-specific prompt-construction pipeline is evaluated.
3. **D5 + D7:** a hard flush clears transient prompt context, not character voice assignment or character configuration.
4. **D6 + HRA:** the temporary admission policy preserves supported multi-model residency and is never treated as T0.5 process proof or global reservation accounting.
5. **D8:** the delegated agents must expose disagreement. They do not manufacture consensus.
6. **D7:** two voice selectors in a UI are not acceptance. FC-6 defines the two-character, distinct-voice, speaker-switch, interruption and restoration test.
7. **Review acceptance ≠ approval:** accepting a review draft records that its text was reviewed; it does not approve the contract or authorize implementation.
8. **FC-1 + FC-4 (memory copies):** FC-4's read-only authored-content display copies are not FC-1's independently owned memory copies; FC-4 establishes no memory-record authority and settles no destination ownership.

---

## 5. What these papers do not do

- They do not approve FC-1–FC-7.
- They do not authorize M0A, Phase 1, M4, or any implementation slice.
- They do not select a tokenizer or prompt-library.
- They do not define `orphaned` by assumption.
- They do not resolve `DQ-HRA-01`–`DQ-HRA-12`.
- They do not select a D8 option.
- They do not change the canonical r33/v9, owner ledger, conformance map, HFS, HFR, HCDM, or Phase 1 records.
- They do not claim test, model, runtime or performance evidence.
- They do not convert review-draft acceptance into contract approval.

---

## 6. Review records (2026-11-21)

**FC-1 through FC-6 completed their review sequence and are accepted as reviewed proposed drafts.** This is a review record, not an approval, and not a statement that any implementation conforms to these papers.

| Paper | Review status | sha256 | Bytes |
|---|---|---|---|
| FC-1 | Accepted review draft — memory-copy semantics (OD-D2-COPY) and subsequent correction passes | `1381cae9bee8d8c39797b0a302160771f2dcd511c12210ddf0465541e89406e1` | 37,949 |
| FC-2 | Accepted review draft — composition-policy identity, immutable-plan and dispatch-attempt separation, content-hash boundary | `39eed6f576fae956e26c56da454f5be5aa91cbae9b72b7abfc41f033a78bbcef` | 14,883 |
| FC-3 | Accepted review draft | `38836b6da73e20a2e5166f400f1754aae77f5aced0cc30a36fe243db35ec2bb8` | 10,465 |
| FC-4 | Accepted review draft — D2 wording reconciled; FC-1 copy-semantics consistency pass | `585bd8a81535e5de440ed5991511b9415b54828f1797e7b5491682ff69ca9ccc` | 11,942 |
| FC-5 | Accepted review draft; **amended after review** by the owner-approved additive amendment (revision 6, 2026-11-21) — configuration decision record and launch construction record | `626ab4bede017da6f3e1431511e229b370635c277acd075915e85aa99b117665` | 23,056 |
| FC-6 | Accepted review draft | `22eaa0cdaf0589fc0d19aaf2399515967a2b26ef14ce554f995d2a5d77c04780` | 9,339 |

These hashes are the frozen review baseline. A later edit to a paper requires its own authorization and breaks the corresponding hash.

FC-4 is not the only paper revised after the baseline commit `0bc3b091b5023ac867f6d060bd56c481d6f84a3a`. FC-4 was revised under the authorized 2026-11-21 D2 wording reconciliation, and later under an authorized FC-1 copy-semantics consistency pass (FC4-13, FC4-14). FC-1 was revised under an owner-authorized memory-copy pass that recorded OD-D2-COPY (2026-11-21), followed by two authorized correction passes. FC-2 was revised under an authorized composition review covering its authority labels, plan identity versus dispatch-attempt identity, the content-hash boundary, privacy-policy identity and budget validity for estimated or unknown segment costs. That commit still contains the earlier text of all three papers; the hashes above are the current files. FC-3 and FC-6 are unchanged from that commit. FC-5 was revised under the owner-approved additive amendment (revision 6, 2026-11-21): four insertion points, 91 lines added and 0 removed, and removing those 91 lines reproduces the baseline text byte for byte, so the FC-5 row above records the amended bytes rather than the frozen review baseline for FC-5.

**Acceptance does not:** approve a paper; authorize implementation, M0A, M4 or Phase 1; select a tokenizer, broker or storage technology; resolve `DQ-HRA-01`–`DQ-HRA-12`; define `orphaned`; select a D8 option; or establish implementation conformance.

**Open threads after this review:** D4 research; D8 delegated decision (FC-7 remains unfiled); D2's unresolved questions — the `orphaned` status of legacy, unknown-origin and interrupted-lifecycle records, chat-scoped records, conversation-event and state-history retention, and privacy deletion — whose wording has been reconciled but whose decisions remain open; FC-1's memory-copy questions — whether the pre-deletion notice is blocking or advisory (`FC1-MEMORY-COPY-NOTIFY`), whether a world destination owns a copy (`FC1-MEMORY-COPY-DESTINATION`), whether copying is available only at deletion (`FC1-MEMORY-COPY-TIMING`), whether the notice also covers snapshots (`FC1-MEMORY-NOTICE-SCOPE`), and destination-side merge and conflict (`FC1-MEMORY-COPY-MERGE`); the `DQ-HRA-*` decisions; M0A's own release-truth gate; M4 authorization.

**D2 wording reconciliation:** performed 2026-11-21 across §3 of this register, the intake packet's D2 records and FC-4 FC4-14. This was a documentation pass only. D2's questions remain unresolved: `orphaned` is still undefined for legacy, unknown-origin and interrupted-lifecycle records, chat-scoped records and conversation-event/state-history retention remain open, and privacy deletion remains a separate unresolved policy. `docs/DOCUMENT-REGISTRY.md` remains unrevised.

**D2 scope transition (2026-11-21, OD-D2-COPY).** The owner directed that memory transfer is a copy operation: the source is preserved and remains parent-owned, each destination copy is independently owned, provenance survives deletion of the source, and no deletion cascades into copies. That direction narrowed the blanket “no automatic orphan cleanup” shorthand previously carried by §3 of this register, by FC-4 FC4-14 and by the intake packet's D2 records. Under FC-1, a parent's own original memory records are cleaned up as part of that parent's explicit deletion (FC1-22e), independently owned copies survive it (FC1-22g, FC1-22i), and conversation-owned reproducibility snapshots are deleted with their conversation (FC1-22c). The record covers four Nov 15–17 sources: `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` (K-R40, K-R45, K-R47); `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` (same rules, S54, S58); `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` (rules 31 and 45, O-13); `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` (Finding 1, D-1, D-4 Option B). Earlier move semantics are superseded by OD-D2-COPY and preserved in FC-1 FC1-22f as superseded history. FC-1's D-4 Option B boundary is unchanged: MEMORY-09 remains absent from the transfer path. This transition establishes no memory-record storage authority in FC-4, settles no destination ownership (FC1-22g, `FC1-MEMORY-COPY-DESTINATION`), resolves no `orphaned` status for legacy, unknown-origin or interrupted-lifecycle records, and authorizes no implementation. The intake packet's D2 wording was subsequently reconciled and accepted.

**End of register.** The contracts are a review surface. Review acceptance, contract approval and implementation authority are three separate owner decisions.
