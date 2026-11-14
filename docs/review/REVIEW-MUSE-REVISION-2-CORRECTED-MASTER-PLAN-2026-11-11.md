# MUSE — REVISION 2 SECTION-SCOPED RE-REVIEW
## Howl Chat Corrected Master Plan

**Review ID:** `REVIEW-MUSE-REVISION-2-CORRECTED-MASTER-PLAN-2026-11-11`
**Date:** 2026-11-11
**Scope:** Section-scoped re-review of Revision 2 against MUSE-01–MUSE-12.
**Authorization:** Read-only. No plan edit, implementation, test, benchmark, model load, migration, network experiment, or gate change was performed.

---

## 0. Artifact identity

| Artifact | SHA-256 | Size / lines | Result |
|---|---|---|---|
| `CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md` | `5c4c8e5604789ff16fd66503f353b49f098aae297aa72308cc35b511dc2a6d36` | 130,529 B / 1,223 | **Matches the reported Revision 2 hash** |
| `REVIEW-MUSE-CORRECTED-MASTER-PLAN-2026-11-11.md` | `56aa7b715a6d962a2ad2f08db04725c4c9fb328c8275daf827b8b160a80ef942` | 34,301 B / 328 | Prior review, unmodified |
| `MUSE-CORRECTION-DISPOSITION-HOWL-2026-11-11.md` | `6a2b8f026b6617077a66cf681e265280d23d4a275da4af4452c74b72089b6b40` | 5,366 B / 37 | Read; not treated as proof |
| `VERIFICATION-HOWL-CORRECTED-MASTER-PLAN-2026-11-11.md` | `1018aee260b6d268a4275ee4999e46c643e544b3c45cd3df6bbb5e8d6837b9c9` | 4,939 B / 68 | Revision 2 verification |
| Recovery addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | 222,234 B / 2,418 | **Unchanged** |
| Reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | 39,612 B / 421 | **Unchanged** |

All five protected F3-C artifact hashes (`85e1968b…`, `ed8bf2e2…`, `343104b5…`, `494539b9…`, `ce32ac0b…`) still verify byte-identical.

The plan grew 107,967 → 130,529 bytes (+22,562) and 1,072 → 1,223 lines (+151), consistent with added registers rather than a rewrite.

---

## 1. Verdict

**READY WITH NON-BLOCKING FINDINGS.**

All twelve findings are corrected in substance, not merely named in the disposition table. The two architectural errors that mattered most — the rollover dependency inversion and the hot-swap/resource cycle — are genuinely fixed. The network mission now has the per-identity accounting it previously only claimed. RAG accountability exists without inventing an unauthorized requirement. Migration safety, import hardening, and investigation authorization are all present as testable contract text.

Four precision defects remain, one of them a decision-level circularity inside the new hot-swap sequence. None blocks owner review; all should be fixed before the plan is used as an execution order.

---

## 2. Finding-by-finding re-review

| Finding | Previous defect | Revision 2 correction (actual text) | Complete? | Residual | Status |
|---|---|---|---|---|---|
| **MUSE-01** | `CWP-05` depended on `CWP-04` | `CWP-05` prerequisites: "`CWP-02`, `CWP-03`; canonical chat identity. Multi-character flush is explicitly **not** a prerequisite." `CWP-04` unchanged as the flush package. | **Yes** | None. The separation now matches `C-05`/`ODR-11`. | **Corrected** |
| **MUSE-02** | `CWP-08` ↔ `CWP-15` mutual prerequisite | New §8.4 four-step handoff `15a → 08a → 15b → 08b`; `CWP-08` prerequisites now "`CWP-15a` … `CWP-15b` before `CWP-08b` per §8.4"; `CWP-15` prerequisites "Phase handoff only … There is no reverse package-level dependency"; `AD-05` and `RG-05` updated. | **Package level: yes. Decision level: no.** | **N1 (Medium):** §8.4 step 1 says `CWP-15a` "May start when `INV-06` and `AD-05` evidence", but `AD-05`'s own "Required evidence before decision" is the `15a → 08a → 15b → 08b` sequence. `AD-05` is therefore both input to and output of step 1. | **Partially corrected** |
| **MUSE-03** | `CWP-01` blocked by three unsequenced gates | New §8.3 table: `UD-01` (Project Owner) → F3-C Phase 0/D-ROOT (owner through existing record) → D-STORAGE → `CWP-01` design review, each with release evidence and permitted pre-release work; "`INV-03` is evidence input … does not resolve F3-C, grant R6, or authorize phase-0 work"; "A local root/alias decision inside `CWP-01` is prohibited." | **Yes** | The mechanism for closing Phase 0 (owner disposition vs R6 machinery) is still not named. That is arguably correct — naming it would be inventing authority — but the owner should be told a `Phase-0` mechanism decision is part of gate 2, not just a "disposition". | **Corrected** |
| **MUSE-04** | Path rows and 28 tests absent at ID level | §10.3 (14 `N-*`), §10.4 (5 `L-*`), §10.5 (8 `P-*`), §10.6 (all 28 `V-NET` with owner, prerequisites, required evidence, state, negative flag, app-boundary flag), §10.7 (`RULE-NET-01`–`04`), §10.8 completion criterion. | **Yes** | None material. Every `V-NET` row is present and all 28 read "**Not run**". `N-14` is correctly retained as display-only, not egress. | **Corrected** |
| **MUSE-05** | RAG had a contradiction but no accountable disposition | New §6.4: no requirement identity invented; provisional `CWP-02` prompt-path ownership only; full lifecycle blocked by `OA-10`; embedding identity, index versioning, deletion/retention, retrieval correctness, and resource cost each assigned conditionally; `RG-13` conditional; `CWP-02` deliverable limits RAG to an ownership/decision brief; `CWP-15` includes embedding consumer "only if `OA-10` authorizes it". | **Yes — accountability complete; scope correctly deferred** | Two unconditional references remain: §9.3 "Prompt path" regression row lists RAG injection, and the addendum-level RAG question has no `UD-` identity. Both are cosmetic relative to `OA-10`. | **Corrected (scope requires owner adjudication)** |
| **MUSE-06** | No backup/atomicity/crash/free-space/single-writer responsibility | `CWP-01` deliverable now: "migration-safety contract: free-space preflight, backup/snapshot with named restore, atomic staging/commit, post-migration integrity verification, resumable or cleanly abortable interruption, single-writer enforcement"; validation adds "fault fixtures for power loss, partial write, insufficient space, and concurrent writer". `CWP-10` carries the same contract. `RG-02` and §9.3 "Durable state" updated. | **Yes** | Memory (`S6`) and character/world (`S3`) stores are dependencies of `CWP-01` but are not named in the migration-safety scope or `RG-02`, which say "chat/root/media". | **Corrected** |
| **MUSE-07** | §23.8 controls absent from deliverables | `CWP-20` adds in-flight cancellation/refusal and renderer/webview navigation guard; `CWP-23` adds `.part`, length, checksum/signature, truncation rejection, no mirror, browser-handoff separation; `CWP-24` adds redirect and oversized branches; `RULE-NET-01`–`04` trace each control to a package and to `V-NET` rows. | **Yes** | **N2 (Low):** `RULE-NET-01` requires checksum/signature verification unconditionally; addendum §23.8 says "verify a published checksum or signature **where available**". The plan is stricter than its authority. | **Corrected** |
| **MUSE-08** | Hot-swap was process-safe, not conversation-coherent | `CWP-08` deliverable adds "conversation-coherence contract: canonical transcript, re-render under the new template, re-tokenize, re-fit to the new budget, preserve/report what changed, and explicitly refuse an incoherent swap"; validation adds "cross-template and cross-context-size fixtures; a visibly reported truncation/re-fit"; investigation adds "architecture/tokenizer/template compatibility and re-fit behavior"; `RG-05` updated. | **Yes** | None. Refusal is explicit rather than silent. | **Corrected** |
| **MUSE-09** | No untrusted-import hardening | `CWP-27` adds "bounded read sizes, decompression limits, path containment for extracted entries, safe embedded-metadata handling, no in-place rewrite …, malformed-input failure that preserves the source"; validation adds "hostile fixtures for malformed, oversized, traversal, archive-bomb, and unsafe-reference inputs fail closed without source mutation"; risk correctly states hardening "is not a validated security guarantee until tested". | **Yes** | None. No unearned security claim. | **Corrected** |
| **MUSE-10** | Plan self-authorized investigations | §8.5: "Every investigation in this plan is **proposed**. This unadopted plan does not authorize its own investigations." `OA-09` added for model/benchmark/network/user-content/runtime investigations; package lines changed to "Investigation is proposed; it requires `OA-09`"; `CWP-13` and `CWP-24` list `OA-09`. | **Yes** | The phrase "authorized by this plan" survives only in negations ("No row is authorized by this plan", "No automatic enablement is authorized by this plan"). | **Corrected** |
| **MUSE-11** | `NET-02` mis-traced to Civitai package | `CWP-21` sources now "`NET-03`, `NET-10`, `NET-11`, `NET-13`"; risk states "`NET-02` is owned by `CWP-22`; Civitai metadata lookup is not cloud inference"; §12.1 `NET-02` → `CWP-22` sole work owner with `CWP-20` governance; §12.2 matches. | **Yes** | Link count legitimately fell 95 → 94 as a result; forward and reverse both 94 with zero mismatch. | **Corrected** |
| **MUSE-12** | Satisfied constraint rows said "Action required" | `REC-M-09`: "Preserve constraint; `CWP-25` is a dependency reference, not a work owner." `NET-20`: "Preserve absence; `CWP-20` is an absence-maintenance check, not a product action." | **Yes** | None. | **Corrected** |

**Tally:** 11 corrected, 1 partially corrected (MUSE-02, package-level cycle resolved; decision-level self-reference remains).

---

## 3. Mandatory deeper challenges

### A. Dependency graph — `15a → 08a → 15b → 08b`

- **Distinct deliverables:** yes. `15a` = resource-authority contract (states, consumers, admission, release, timeout, failure). `08a` = hot-swap consumer/state contract (drain, stop, spawn, readiness, commit, rollback). `15b` = enforced admission/release (design). `08b` = hot-swap implementation design. Each is separately named and separately owned.
- **Explicit handoff:** partially. The table gives "May start when" for each step, which is an input condition, but no step names its **output artifact consumed by the next step** beyond the column label. `15b` says it consumes `08a`; `08b` consumes `15b`. `08a` consumes `15a`. The artifact-to-artifact handoff is inferable but not stated as "input: X contract, output: Y contract".
- **Indirect later dependency:** **yes, one.** `AD-05` is listed as a start condition for `15a` while `AD-05`'s required evidence is the completed four-step sequence (see N1). This is a decision-level edge back from step 1 to step 4.
- **Package-level acyclicity:** yes. `CWP-15` has no dependency on `CWP-08`; `CWP-08` depends on `CWP-15a` and `CWP-15b`. No package-level cycle remains.
- **Implementation hidden in a preparatory phase:** no. `15b` is labeled "enforced implementation **design**" in the package deliverable and §8.4 states "No implementation is authorized at any step." Note **N3** below on the table's "enforced admission/release" phrasing, which reads more like implementation than the deliverable does.
- **`CWP-01` gates bypassed:** no. §8.3 keeps all three gates and explicitly prohibits a local root/alias decision inside `CWP-01`; `INV-03` is demoted to evidence input. Note that `08a` requires `CWP-01`, so hot-swap now correctly waits behind the keystone rather than around it.

### B. RAG ownership (`OA-10`)

- **Does `OA-10` prevent unauthorized scope?** Yes. `OA-10` names "RAG committed scope and full retrieval lifecycle" as not authorized, and §6.4 blocks embedding model, index, deletion, retention, rebuild, and resource claims behind it.
- **Does `CWP-02` accidentally own RAG architecture?** No, but the wording is close. `CWP-02`'s deliverable says "RAG is a **provisional, blocked sub-scope** … only an ownership/decision brief is permitted; no retrieval lifecycle implementation," and §6.4 says "`CWP-02` owns prompt-path integration only while scope is unresolved." "Owns … only" is bounded by the same sentence, so the risk is low; a stricter phrasing would be "`CWP-02` may investigate RAG integration ownership; it does not own RAG architecture."
- **Is `RG-13` conditional?** Yes — "If RAG is authorized by `OA-10` … Blocked — no requirement identity or owner authorization."
- **Are the sub-responsibilities blocked rather than quietly assigned?** Yes for embedding selection, index lifecycle, retrieval validation, deletion/retention, and resource cost: each row begins "Not authorized" and uses conditional "would own" phrasing.
- **Residual:** §9.3 lists RAG unconditionally in the prompt-path regression row, and `CWP-28` lists "RAG" as a dependency. Neither creates authorization; both should be marked conditional for consistency (N4).

### C. Durable-state integrity

Measurable acceptance now exists for: backup/restore (named restore in deliverable), atomic migration (staging/commit), crash consistency (power-loss fixture), free-space preflight, single-writer/concurrency (concurrent-writer fixture), partial-write recovery (partial-write fixture), and post-migration integrity verification. These are contract-level and testable.

Coverage: **chat and media are explicit** (`CWP-01`, `CWP-10`, `RG-02`, §9.3 "Durable state"). **Memory and character/world stores are not named**, though `CWP-01` lists `S3` and `S6` as dependencies and its schema must survive them. That is the one gap worth closing (N5).

### D. Network and download integrity

- `.part`, length/checksum, truncation rejection, no mirror, offline cancellation, and webview navigation guard are each preserved with an owning package (`CWP-20`, `CWP-23`, `CWP-24`) and a verification mapping (`V-NET-10/18/21/23/24/25`, `RULE-NET-01`–`04`).
- All 28 `V-NET` rows are individually identifiable and all read "**Not run**" (verified by count: 28 rows, 28 "Not run").
- **Grounding:** three of four `RULE-NET` entries match §23.8 wording (cancellation, navigation guard, no-silent-substitution). `RULE-NET-01` tightens "published checksum or signature where available" into an unconditional requirement (**N2**). All other new network text traces to `NET-*` or `V-NET-*`; no new requirement was invented.

### E. Model hot-swap coherence

Tokenizers, templates, architecture, re-fit, active-stream drain, readiness, commit, rollback, unload, resource release, and visible change reporting are all now named. Sampler and generation-parameter compatibility are covered indirectly through `CWP-09` provenance and `CWP-11` identity rather than in `CWP-08` itself — acceptable, since a swap does not change samplers but may change which are valid.

Incompatible swaps now **fail explicitly**: "explicitly refuse an incoherent swap" plus "a visibly reported truncation/re-fit." That is the correct behavior; silent incoherent output is no longer an acceptable outcome.

### F. Import hardening

Bounded reads, decompression limits, path containment, safe metadata handling, no in-place rewrite, malformed-input source preservation, partial-import failure behavior ("fail closed without source mutation"), and hostile fixtures are all present and testable. No security guarantee is claimed without evidence — the risk row says hardening "is not a validated security guarantee until tested."

### G. Governance

`OA-09` prevents self-authorization: investigations are "proposed," and any model/benchmark/network/user-content/runtime investigation requires it. `OA-10` prevents RAG scope creep. `§9.4` still forbids a "done" status and self-acceptance. No implementation, enablement, gate passage, or package acceptance is implied anywhere; §12.3 and §12.4 retain "not adopted / no implementation authorized." The two surviving "authorized by this plan" strings are both negations.

---

## 4. Traceability spot-check

| Check | Result |
|---|---|
| 66 original rows | **66**, each once |
| `NET-01`–`NET-20` | **20**, each once |
| Forward rows | **86 unique**, no duplicates |
| Outcome totals | **30 P / 30 M / 13 U / 5 C / 2 E / 6 S** — unchanged from reconciliation |
| 28 `V-NET` individual dispositions | **28**, all `Not run` |
| `N`/`L`/`P` identities | **14 / 5 / 8 = 27**, all present |
| 37 capability rows | **37**, seven-state vector preserved, no promotions |
| Forward/reverse consistency | **94 links each, 0 mismatches, 0 orphan packages** |
| Link-count change 95 → 94 | Expected consequence of removing `NET-02` from `CWP-21`; both directions updated |
| Owner decisions invented | **None**; `UD-01`–`UD-06` still unresolved |
| New requirement invented | **None**; `RULE-NET-*` are labeled existing addendum controls |

---

## 5. New findings introduced in Revision 2

These did not exist in Revision 1.

**N1 — Medium — `AD-05` is both prerequisite and product of the hot-swap sequence.**
§8.4 step 1: "`CWP-15a` … May start when `INV-06` and `AD-05` evidence." §7 `AD-05` "Required evidence before decision": "`CWP-15a` → `CWP-08a` → `CWP-15b` → `CWP-08b`." Step 1 therefore needs a decision whose evidence is produced by steps 1–4.
*Why it matters:* it is the last remaining cycle, and it sits on the riskiest runtime change. It also invites a "blocked on AD-05" claim that stalls step 1 indefinitely.
*Correction:* `15a` starts on `INV-06` alone; `AD-05` is recorded at the `15a`→`08a` boundary (or immediately after `15a`) as the decision adopting the resource contract. One line in both places.

**N2 — Low — `RULE-NET-01` is stricter than its controlling source.**
Addendum §23.8: "verify a published checksum or signature **where available**." `RULE-NET-01` omits the qualifier, making verification unconditional.
*Why it matters:* a release gate that cannot be satisfied for sources publishing no checksum becomes an unpassable acceptance criterion, or gets waived informally.
*Correction:* restore "where available," or state explicitly that unverified downloads are permitted only with a visible "unverified" label and remain non-authoritative.

**N3 — Low — §8.4 step 3 artifact label implies implementation.**
Table: "`CWP-15b` — enforced admission/release"; package deliverable: "`CWP-15b` enforced **implementation design**"; §8.4: "No implementation is authorized at any step." The three statements are reconcilable but the table reads as authorization-adjacent.
*Correction:* rename the step to "enforced admission/release **design**" to match the deliverable.

**N4 — Low — RAG is unconditional in two regression/dependency places.**
§9.3 "Prompt path" row lists "RAG" injection unconditionally; `CWP-28` lists "RAG" among dependencies. Every other RAG reference is `OA-10`-conditional, and there is a correctly conditional "RAG (conditional)" row.
*Correction:* mark both as "RAG (if `OA-10`)".

**N5 — Low — memory and character/world stores are outside the named migration-safety scope.**
`CWP-01` depends on `S3` and `S6`; `RG-02` and §9.3 "Durable state" name chat/root/media only.
*Correction:* add memory/character/world state to `RG-02`'s scope sentence, or state explicitly that they are not migrated in this phase and why.

**N6 — Observation — §9.3 now has two `Network` rows.**
Lines 885 and 889. The second (with cancellation, `.part`/checksum, webview guard) should be merged into the first; the section declares one row per domain.
*Correction:* merge.

**N7 — Observation — RG-12 says "47 path identities".**
§2.1, §10, and §12.3 all correctly say **27**. `47` appears to be 27 paths + 20 `NET` rows. *Correction:* change to "27 path identities and 20 network requirements."

---

## 6. Required before adoption (small, mechanical)

1. Fix the `AD-05` self-reference (N1) — the only substantive item.
2. Restore "where available" to `RULE-NET-01` (N2).
3. Rename §8.4 step 3 to a design artifact (N3).
4. Mark the two RAG references conditional (N4).
5. Add memory/character/world to `RG-02` scope (N5).
6. Merge the duplicate `Network` regression row (N6) and correct `47` → `27` (N7).

No further architectural re-review is required after these; they are consistency repairs inside already-correct sections.

---

## 7. What this review did not do

- Did not edit the plan, disposition, verification record, addendum, reconciliation, ledger, stage record, or any protected F3-C artifact.
- Did not accept the disposition record's claims; each correction was located in the revised text and measured against the finding it claims to close.
- Did not implement, configure, migrate, load a model, run a benchmark, exercise the UI, or perform any network experiment.
- Did not mark any `V-NET` test run, any capability operational, or any gate passed.
- Did not resolve `UD-01`–`UD-06`, the RAG scope question, the remote ASR/TTS question, or the Phase-0 closure mechanism.
- Did not authorize implementation, adoption, or enablement of anything.

**Remaining blockers to owner adoption:** none architectural. Six documentation corrections (§6) should land first so the plan's own counts and decision references are internally consistent.
