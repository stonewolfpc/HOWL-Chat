# MUSE correction disposition — Howl corrected master plan

**Plan:** `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md`  
**Plan SHA-256:** `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874`  
**Plan size:** 131,755 bytes / 1,223 lines  
**Review:** `docs/review/REVIEW-MUSE-CORRECTED-MASTER-PLAN-2026-11-11.md`  
**Review SHA-256:** `56aa7b715a6d962a2ad2f08db04725c4c9fb328c8275daf827b8b160a80ef942`  
**Status:** Revision 3 proposed plan — **not adopted; no implementation authorization**  
**Revision identity:** Revision 2 recorded the MUSE-01–MUSE-12 correction. Revision 3 is the separate N1–N7 closure identity and O-1–O-3 documentary cleanup; technical scope, work packages, and dependencies are frozen except for those corrections.

## Dispositions

| Finding | Disposition | Sections / packages changed |
|---|---|---|
| **MUSE-01** — rollover incorrectly depends on multi-character flush | **Accepted and corrected.** `CWP-05` now depends on `CWP-02`, `CWP-03`, and canonical identity only; `CWP-04` is explicitly not a prerequisite. | `CWP-05`; rollover/flush separation remains in `C-05` and `ODR-11`. |
| **MUSE-02** — CWP-08/CWP-15 cycle | **Accepted and corrected.** Added an acyclic four-step phase handoff: `CWP-15a` resource contract → `CWP-08a` hot-swap contract → `CWP-15b` enforcement → `CWP-08b` implementation. No reverse package-level edge. | §8.4; `AD-05`; `CWP-08`; `CWP-15`; `RG-05`. |
| **MUSE-03** — CWP-01 has three unsequenced gates | **Accepted and corrected.** Added named order, authority, release evidence, and permitted preparatory work for `UD-01`, F3-C Phase 0/D-ROOT, and D-STORAGE. `INV-03` is evidence input, not gate release. | §8.3; `CWP-01`; `RG-02`. |
| **MUSE-04** — network identities only referenced by range | **Accepted and corrected.** Added individual `N-01`–`N-14`, `L-01`–`L-05`, `P-01`–`P-08`, and all 28 `V-NET` rows with proof, owner, prerequisites, evidence, state, negative flag, and app-boundary flag. `V-NET-23`, `V-NET-26`, and `V-NET-27` are explicit. | §10.3–§10.8; `CWP-20`–`CWP-25`; §12.2 note. |
| **MUSE-05** — RAG has no accountable disposition | **Partially corrected; scope deferred to owner.** No requirement or `UD-` identity is fabricated. §6.4 records that no RAG requirement exists; `CWP-02` is provisional prompt-path owner only; `OA-10` blocks the full lifecycle; `RG-13` is conditional. Embedding identity, index versioning, retrieval correctness, deletion/retention, and resource coordination are assigned as conditional responsibilities. | §6.4; §7.1 `OA-10`; `CWP-02`; `CWP-15`; `RG-13`; §9.3. |
| **MUSE-06** — durable-state integrity absent | **Accepted and corrected.** Added backup/restore, atomic staging/commit, crash/interruption behavior, free-space preflight, post-migration integrity, single-writer enforcement, and partial-write tests to chat and media migration ownership. | `CWP-01`; `CWP-10`; §9.3; `RG-02`. |
| **MUSE-07** — addendum §23.8 rules not carried into deliverables | **Accepted and corrected.** Added individual control trace for `.part`, length/checksum/signature, truncation rejection, no mirror substitution, offline cancellation, and webview navigation guard. | `CWP-20`; `CWP-23`; `CWP-24`; §10.7; `V-NET-10/18/23/24/25/28`; `RULE-NET-01`–`04`. |
| **MUSE-08** — hot-swap preserves process, not conversation | **Accepted and corrected.** Added canonical transcript, new-template re-render, re-tokenization, re-fit, visible change reporting, refusal of incoherent swaps, and cross-template/context-size tests. | `CWP-08`; `RG-05`. |
| **MUSE-09** — no untrusted-import hardening | **Accepted and corrected.** Added bounded reads, decompression limits, path containment, safe embedded metadata, no in-place rewrite, malformed-input source preservation, and hostile fixtures. | `CWP-27`; §9.3 import/export row. |
| **MUSE-10** — plan self-authorizes investigations | **Accepted and corrected.** Removed self-authorizing phrases; added `OA-09`; all runtime/model/benchmark/network/user-content investigation now requires explicit owner authorization. | §8.1; §8.5; §7.1 `OA-09`; authorization lines in `CWP-01`, `CWP-02`, `CWP-06`, `CWP-07`, `CWP-09`; `CWP-13`; `CWP-24`. |
| **MUSE-11** — NET-02 mis-traced to Civitai package | **Accepted and corrected.** `NET-02` is owned solely by `CWP-22`; `CWP-21` no longer lists it. `CWP-20` is named as governance dependency, not work owner. | §12.1 `NET-02`; §12.2 `CWP-21`/`CWP-22`; `CWP-21` sources and risks. |
| **MUSE-12** — satisfied constraint rows say “Action required” | **Accepted and corrected.** `REC-M-09` now says preserve constraint with `CWP-25` as dependency reference; `NET-20` says preserve absence with `CWP-20` as absence-maintenance check. | §12.1. |

## Structural result

- **28 work packages remain; no new package identity was added.** `CWP-15a/15b` and `CWP-08a/08b` are phases inside existing packages.
- **66 original rows and 20 network rows remain exactly once.**
- **No new requirement, owner decision, exclusion, or `UD-` question was invented.**
- **RAG remains without a requirement identity.** Its scope decision is explicit and deferred.
- **The dependency graph is acyclic:** the only hot-swap/resource ordering is `15a → 08a → 15b → 08b`; package-level `CWP-15` has no dependency on `CWP-08`.
- **All 28 V-NET tests have individual rows and are not run.**
- **No authorization boundary was crossed.**

MUSE-13–MUSE-16 are review observations, not required findings. They remain deferred and were not used to expand the plan.

## Revision 3 consistency closure (N1–N7)

| Finding | Disposition | Exact correction |
|---|---|---|
| **N1** — AD-05 is both prerequisite and product | **Accepted and corrected.** | AD-05 now records `CWP-15a` as its initial evidence and is recorded at the `15a → 08a` boundary. `CWP-15a` starts on `INV-06` plus authorization to prepare the contract; AD-05 is not an input. |
| **N2** — checksum requirement stronger than source | **Accepted and corrected.** | `RULE-NET-01` and `CWP-23` now say “checksum/signature **where available**”; absent verification must be visibly marked unverified and non-authoritative. |
| **N3** — enforcement terminology implies implementation | **Accepted and corrected.** | §8.4 step 3 and `CWP-15` now call `CWP-15b` an admission/release **enforcement design**, not enforced implementation. |
| **N4** — RAG references unconditional | **Accepted and corrected.** | §9.3 prompt-path row and `CWP-28` dependencies now mark RAG as conditional on `OA-10`. |
| **N5** — memory/character/world migration scope omitted | **Accepted and corrected.** | `CWP-01` now covers chat, memory, character, and world stores; media remains `CWP-10`; §9.3 and `RG-02` name the same scope. |
| **N6** — duplicate Network regression rows | **Accepted and corrected.** | §9.3 has one merged Network row containing cancellation, `.part`/checksum/truncation, and webview guard. |
| **N7** — incorrect path count | **Accepted and corrected.** | `RG-12` now says 28 V-NET tests, 27 path identities, and 20 network requirements. |

## Documentary cleanups (O-1–O-3)

| Observation | Disposition |
|---|---|
| **O-1 — Revision 2 label covered two contents** | **Resolved.** The N1–N7 closure is now Revision 3 in the plan, this disposition, and the verification record. Revision 2 remains the MUSE-01–MUSE-12 correction identity. |
| **O-2 — “AD-05 is adopted” wording** | **Resolved.** The plan and records say `AD-05` is **recorded** at the `15a → 08a` boundary. No dependency or phase-order change was made. |
| **O-3 — outcome distribution absent from verification** | **Resolved.** The verification record now states both scopes: 6 Satisfied rows in the 86-row forward trace, with 5 distinct qualifiers; 19 Satisfied rows in the 103-row full reconciliation, with 13 additional protected-register rows. The retired variants shorthand remains removed. |
