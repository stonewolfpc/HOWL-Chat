# REVIEW — Independent Review of the F3C-TARGET-2 R5 Integration Review

**Status:** **INDEPENDENT REVIEW ONLY — NOT APPROVAL, NOT FREEZE, NOT R6, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Reviewed artifact:** `docs/review/INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md`
**Integration review SHA-256:** `b29bb562c2a5b596baa3bb5689d6aaaf0adf9c580ac18ab72a99eae8c13f03f6`
**Integration review size:** 411 lines; 29,097 bytes.
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
**R5 SHA-256:** `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`
**R5 size:** 1475 lines; 77,233 bytes.
**Original matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
**Original matrix SHA-256:** `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`
**First precision addendum:** `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
**First precision addendum SHA-256:** `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`
**Second precision addendum:** `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md`
**Second precision addendum SHA-256:** `be399e5bfda3fec62e4356ad8e0dadde11a643881bd74fc314ba2bf9b116e403`
**Latest precision review:** `docs/review/REVIEW-INDEPENDENT-ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md`
**Latest precision review SHA-256:** `d97870088f74e8c9b2e79580225f5bbb56a08cca08e1d9cf9f66e9cfdab9361d`
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
**Gate v2 SHA-256:** `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`
**Scope:** documentation-only independent review of the integration artifact. No source, test, build, migration, model-load, destructive-test, application-data, authority-record, freeze, acceptance, P1a, R6, or F3-D action was performed or authorized. No existing file was modified. No repository source inspection was performed; this review tests the integration document's internal coherence and its stated proof obligations.

---

## 1. Verdict

**The integration review is the right artifact, and its central judgment is correct: the architecture is conditionally coherent, not proven.**

Preserve the integration review as the authoritative composition of R5, the original matrix, and both precision addenda. Accept its five integration findings (I1–I5), its state machine, its entry-point inventory, its retention lifecycle, its conditional lock analysis, and its twelve-entry contradiction register as requirements and boundaries — not as proof.

**Do not authorize R6.** The three attack areas below each contain at least one gap that would let an implementation be built to the letter of the current documents yet violate the safety property the documents claim.

---

## 2. The integration review's load-bearing claims

The review stands or falls on four compositions:

1. **OD-18 (I2/CR-04):** one authority process per storage root, recommended as a storage-level exclusive lock plus a bounded deployment assumption for split-view cases.
2. **Commit/durability/release separation (§3.2):** visibility, power-loss durability, and gate release are three distinct points; no generic `success` state.
3. **Storage-access inventory (§5):** chat, creation, cast, delete, wipe, migrations, plans, startup recovery, gate release, gallery, attachments, import, export, diagnostics, background writes, world/scenario snapshots, content library, browser storage, and operation-map paths each get a gate class and lock boundary.
4. **Contradiction register (§8, CR-01–CR-12):** genuine conflicts, compatible boundaries, and scope mismatches are distinguished rather than blurred.

All four are the right structures. Each is attacked below.

---

## 3. Attack 1 — OD-18 implementation feasibility

### 3.1 What the integration review gets right

- Incarnation numbers are correctly demoted from exclusion mechanism to identity.
- Go mutexes are correctly confined to in-process synchronization.
- The lock is correctly placed as an admission lock before any F3-C lock, never nested inside `chatWriteMu`.
- Split views are correctly named as the defeat condition for any lock that lives inside one view.

### 3.2 Feasibility gaps R6 must close

1. **Supported-storage enumeration is missing.** The review requires proof "across supported filesystems" but does not enumerate them. R6 must name the exact supported Windows configurations: local NTFS, and explicitly in-or-out for ReFS, SMB/network shares, removable USB, OneDrive- or cloud-synced folders, containers, and snapshotted/duplicated volumes. Every unlisted configuration must fail closed or be declared unsupported — never silently treated as locked.

2. **Stale-handle behavior after process crash is unspecified.** If the authority process dies while holding the lock, R6 must define: who may observe staleness, how staleness is distinguished from a live holder, what recovery may and may not do before the lock is re-acquired, and what happens if the stale handle cannot be cleared. Crash recovery that proceeds without re-acquired exclusion is a second writer.

3. **Lock-file deletion and recreation races.** If the lock is a lock file, deletion by cleanup, antivirus, sync tools, or a second process that believes the lock is stale can fork authority. R6 must define the lock path's ownership (under OD-02b/OD-17), forbid cleanup tools from touching it by path, and make deletion-then-recreation resolve to exactly one owner or fail closed.

4. **Network and removable filesystems may not provide exclusion.** Opportunistic locks, byte-range locks, and exclusive-create semantics differ across SMB versions and removable media, and may silently degrade. R6 must either prove the chosen primitive on each supported configuration or exclude that configuration. A primitive that returns success without exclusion is worse than no lock: it manufactures false confidence.

5. **Split-view containment needs a boundary, not a hope.** Option A plus a bounded deployment assumption is a coherent combination only if the assumption is enforceable: how does the system detect that two processes hold the same root through different mounts, and what does it do when it cannot tell? "Outside the supported deployment" must map to a concrete refusal, not documentation prose. If detection is impossible, the residual risk must be owned explicitly, and the architecture must not be described as single-writer safe in those deployments.

6. **Lock acquisition failure must be fail-closed and diagnosable.** Every F3-C storage access denied by lock failure needs a deterministic error that blocks ordinary access and recovery alike, distinguishes "lock held by live authority" from "lock subsystem broken," and never degrades into proceeding without the lock after a timeout or retry budget.

**Disposition:** OD-18 remains correctly unresolved. The recommendation direction (storage-level exclusion plus bounded assumption) is accepted; the feasibility proof is not. R6 must carry a per-configuration lock matrix with crash, stale-handle, deletion, network, removable, and split-view columns, each cell either proved or excluded with a fail-closed behavior.

---

## 4. Attack 2 — The state machine

### 4.1 What the integration review gets right

- Visibility, durability, and release are correctly separated; `S2_RECORD_INDETERMINATE`, `S5_VERIFIED_PENDING_RELEASE`, and `S6_RELEASED` enforce the separation structurally.
- Torn bytes are correctly barred from identifying a package.
- Matching target bytes are correctly barred from implying commitment.
- Orphan packages and orphan markers are correctly denied commit status.

### 4.2 Completeness challenges

1. **Every transition needs a total event mapping.** §3.3 lists the required transitions but does not prove totality. R6 must answer, for each state, every possible next observation: new commit arrival, power loss, package disappearance, marker loss, duplicate/conflicting record discovery, wipe intent, tombstone arrival, and retention expiry. Any unmapped (state, event) pair is a future accidental release.

2. **Resolution exits for `S7_BLOCKED` and `D1_DELETING`/`W1_WIPING` are underdetermined.** `S7` is a hard block, but the machine must define the only legal exits (operator/owner disposition, wipe, or destructive path — never silent repair-by-guess) and what evidence each exit requires. `D1` and `W1` need defined completion and failure exits, including crash mid-wipe and crash mid-delete, with no path that strands ordinary access half-open.

3. **`S5` is the most dangerous state and needs the strongest guard.** "Targets match, verified state durable, release not yet complete" sits one step from ordinary access. R6 must prove: what serializes release against a successor commit that begins projection concurrently; what happens if power is lost after release-state write but before gate observation flips; and what prevents a reader from observing `S5` bytes as released. Release must be a durable, serialized, observable flip — not a race between a writer and a reader.

4. **`S2` re-entry from `S6` must not expose mixed successor state.** The transition `S6 → S2` on new-commit projection is correct, but R6 must define the exact instant ordinary reads stop: readers admitted just before the transition must not observe a half-projected successor, and readers arriving during `S2` must block even if some target bytes still look like the old release.

5. **Marker/package loss during `S4` must have deterministic outcomes.** If the repair marker or the referenced package becomes unreadable mid-repair, the machine must route to exactly one of retry-from-record, `S7` block, or abort-to-prior-verified-state — never to package substitution by content guess.

6. **Duplicate/conflicting records need the exclusive-create proof.** §3.4 routes duplicates to `S7` via exclusive-create and the operation map. R6 must show the create primitive is actually exclusive on every supported configuration (same feasibility matrix as OD-18) and that "never last-writer inference" holds under crash between the two competing creates.

**Disposition:** accept the state machine as the required contract shape. Require R6 to publish the total transition table with one deterministic recovery outcome per (state, event) pair and a proof that no path reaches ordinary access except through `S5 → S6` release.

---

## 5. Attack 3 — The lock proof

### 5.1 What the integration review gets right

- The frozen order is restated without silent modification.
- `chatStateMu` is correctly confined to a brief leaf with no disk I/O.
- No new F3-C lock may be inserted; OD-18 is correctly placed outside the order.
- Startup, multi-chat, root-journal, verification, gate-release, and operation-map paths are correctly marked conditional or gapped rather than claimed proven.

### 5.2 Proof gaps confirmed and extended

1. **Startup recovery exact sequence.** The required order (exclusion → incarnation → authority/lease/barrier reconciliation → package/record/marker/target reconciliation → gate release) is the right skeleton. R6 must bind each step to exact lock acquisitions and releases, state where root journals are written, and prove no step publishes anything a gate or reader can observe before release. Ambiguous incarnation state must fail closed before any target is touched.

2. **Multi-chat recovery ordering.** Stable `(rootId, chatId)` order is necessary but not sufficient: R6 must prove no path holds a lower-ordered chat lock while waiting on a higher-ordered one across recovery, repair, marker write, and release, including error and retry paths. One reverse wait in an error handler is an inversion.

3. **Wipe exclusion under crash.** Tree-Lock excluding recovery's tree-R is correct while both run, but crash mid-wipe leaves wipe intent, partial new-root state, and old roots behind. R6 must prove wipe intent is durable, wipe recovery resumes under the same exclusion, and startup recovery cannot begin repairing the old root while wipe recovery is pending. Wipe and recovery must be mutually exclusive across restarts, not just within one process lifetime.

4. **Gate release serialization.** Release must be serialized with successor projection and ordinary readers. R6 must give the exact reader lock sequence (what shared lock or sequence a gate observation plus target read holds), the exact release write sequence, and the proof that neither a reader nor a successor projection can interleave between verification and observable release.

5. **Verification purity and operation-map discipline.** Pure-function verification and chatState-leaf operation-map use are correctly stated as required rules. R6 must enforce them structurally (no disk I/O under `chatStateMu`, no upward lock from verification under `chatWriteMu`) rather than by convention, with a named owner per caller.

**Disposition:** the conditional-correctness framing is correct and must be preserved. The proof is not complete, and the five items above are the minimum proof target. No implementation may claim the frozen lock order covers recovery, wipe, or release until R6 discharges them.

---

## 6. Attack 4 — Contradiction register ownership

The twelve-entry register is accepted as a boundary artifact with one correction in emphasis: **none of the twelve may be closed by inference, compatible labeling, or a future-R6 promise — only by a named owner decision or formal amendment.**

- **CR-01 (frozen §4.1 vs new roots):** owner OD-17 ruling plus formal amendment if §4.1 was exhaustive. OD-02b enumeration alone cannot overrule a frozen exhaustive layout.
- **CR-02 (create-only records vs repeated target replacement):** owner AM-1/OD-02a ruling. No reinterpretation of "create-only" in R6 prose.
- **CR-03 (same-volume atomic move vs eight targets plus deletions):** per-file protocol proof; multi-file atomicity must never be claimed.
- **CR-04 (in-process locks vs cross-process writer):** OD-18 decision; bounded assumption if no shared primitive. This is the critical path.
- **CR-05/CR-06 (frozen target/record schemas vs v2 additions):** formal compatibility/version ruling. V2 fields must not be described as frozen-compliant before that ruling.
- **CR-07 (retained evidence vs cleanup):** positive retention rule with reference-aware deletion. Age- or prefix-based cleanup of referenced evidence must be forbidden structurally.
- **CR-08/CR-09 (V18–V24/A14/V24):** correctly held as compatible-but-blocked. Require a new Gate record before any PASS claim; keep V24 blocked pending the multi-character baseline.
- **CR-10 (world/scenario audit vs chat-store authority):** correctly marked scope mismatch. Forbid citation of the world/scenario audit as chat-store evidence in R6 and implementation review.
- **CR-11 (wipe inventory vs new roots):** every new OD-02b path, marker, and projection state must appear in the wipe inventory with owner and crash behavior.
- **CR-12 (tmp cleanup vs durable markers):** path-prefix or age deletion of referenced evidence must be forbidden; any tmp-like hosting of durable state needs explicit retention linkage.

Each contradiction needs an explicit owner and disposition before R6 is drafted. The register itself must be carried into R6 as an acceptance checklist.

---

## 7. Inventory assessment

The storage-access inventory is the right boundary expansion: a chat-message-only safety case would be insufficient for Howl, and the review correctly refuses that narrowing. Gallery copies, attachments, snapshots, imports, exports, diagnostics copies, browser storage, and operation-map paths are all real exfiltration or bypass routes for unverified state.

Two reinforcements are required:

1. **Source-level enumeration remains unverified by the review's own admission.** Accept that limit explicitly: the contract inventory defines the classes and required per-path obligations (gate check, lock, access class, error behavior, source-level owner), but R6 must bind every concrete caller to exactly one row. Any caller without a row fails closed.
2. **Copy and inspection paths need publication-grade rules.** Gallery staging, diagnostics support copies, exports, and background autosave must follow the same verification-before-publication rule as ordinary reads. "Read-only" or "background" is not a bypass class; R6 must assign each such path either an ordinary gate or an explicit privileged class with owner and audit.

---

## 8. Authorization boundary

- **Integration review:** preserved as composition and proof-obligation authority. Not R6, not a contract revision, not proof.
- **R5:** proposal only. Not frozen, not accepted, not implemented, not verified.
- **Matrices and addenda:** preserved; findings accepted as requirements, not approvals.
- **OD-18:** unresolved. Blocks any cross-process single-writer claim.
- **OD-02a / OD-02b / OD-17:** awaiting owner; OD-02a is inert without OD-02b.
- **R6:** not authorized, not written.
- **F3-C phase 0:** open. **P1a:** unauthorized. **F3-C acceptance:** not granted. **F3-D:** blocked.
- **A9/A10/A12/A13:** closed; not reopened. **V20:** blocked. **Gate v2:** V1–V17 only. **A14/V24:** open / blocked.

**Required before R6:** owner rules on OD-17 and OD-18; separate OD-02a/OD-02b decisions with joint activation confirmed; acceptance or replacement of the state machine; acceptance of the entry-point inventory and retention lifecycle; explicit owners and dispositions for all twelve contradictions; and a lock-order proof discharging §7.2–§7.4 plus the extensions in §5 above.

No R6 drafting, freeze, acceptance, or implementation follows from this review.
