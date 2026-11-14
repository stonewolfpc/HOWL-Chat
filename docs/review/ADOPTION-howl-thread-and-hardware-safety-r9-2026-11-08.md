# Adoption Record — Howl Thread/Hardware/Performance Requirements, r9

**Prepared:** 2026-11-07 CST (UTC−06:00); session UTC date 2026-11-08  
**Status:** ADOPTED — REQUIREMENTS AND ACCEPTANCE CRITERIA ONLY — EFFECTIVE 2026-11-08  
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r9-2026-11-08.md`  
**Subject SHA-256:** `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`  
**Subject size:** 769 lines / 112,270 bytes, CRLF (769 CRLF pairs; 769 CR bytes)  
**Independent review:** `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r9-2026-11-08.md`, SHA-256 `606acaec2e05a3c384c8dd6a028f0ac2d84a2e39c853e2d9bbe145ea415e942a`, 161 lines / 18,530 bytes; verdict **Adoptable as written.**  
**Scope of this instrument:** This instrument adopts r9 as the requirements and acceptance-criteria baseline only. It does not authorize C1–C8, implementation, runtime execution, benchmarking, or deployment, and it does not resolve A-1 or any Stage 3 decision. The proposed stage plan remains unamended and non-controlling.

## 1. Artifact identity and discrepancy reconciliation

The binding artifact is the CRLF file identified above, verified from its actual bytes:

- `wc -l -c`: `769 112270`.
- SHA-256: `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`.
- Line endings: CRLF, 769 CR bytes, 769 LF bytes, 769 CRLF pairs.
- The audit's `112,270 bytes` and `0f71cee5…` are correct and bind the same byte sequence.

The earlier `111,501` figure was the length of the LF-normalized text after universal-newline translation, not the size of the file. It cannot be paired with the actual-byte SHA-256. The LF-normalized copy is a different byte sequence and must not be adopted. No hash collision is involved; the earlier byte-count statement is corrected, and the SHA-256 remains `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` for the actual artifact.

**Prepared-draft SHA-256 (historical, not current file identity):** `d4917b5407f1f58d4a7559a7d7beda3c61e0d9a80bc794a05955bcbf3f63a6a3` identified the prior draft before these corrections. It is not the identity of this revised draft or of the final adoption instrument.

### 1.1 Date basis and chronology

Preparation and verification occurred on **2026-11-07 local time (CST, UTC−06:00)**. The session's UTC date was **2026-11-08**. The r9 and audit headers and filenames use `2026-11-08` because that is the session/UTC date; their local file timestamps are 2026-11-07 (r9 at 20:01 CST; audit at 20:16 CST). They are not silently future-dated relative to UTC, and their bytes and hashes are unchanged.

This record was prepared on 2026-11-07 local time and retains the subject's existing `2026-11-08` path date to preserve artifact identity. No date is changed in r9 or the audit. The owner adopted the record on 2026-11-08 (UTC; 2026-11-07 CST).

## 2. Governing-input table — exact reproduction (M1 treatment)

This table is reproduced in full from r9 §2.2. It is not incorporated by reference or by attachment path. All seventeen rows were independently recomputed from the actual files on **2026-11-07 CST (2026-11-08 UTC)**.

| Input | Path | SHA-256 | Lines | Bytes |
|---|---|---|---:|---:|
| Original recorded requirements | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` | `ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5` | 344 | 25,933 |
| Audited r2 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md` | `c9bddd10608993ae73d86caa01a7adb6dd83220f98da8dab17e46c1aa3deb8ab` | 447 | 36,025 |
| Mimo r2 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r2-2026-11-07.md` | `567bc86c3d8b2a78b9687131d5937d77a9736393bed6a383dd0aba5ff948994b` | 221 | 29,130 |
| Audited r3 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r3-2026-11-07.md` | `e45f1d898ab3d66016d2714472fd1023bcc93926a5e250a94bbb9072d9082d16` | 513 | 47,710 |
| Mimo r3 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r3-2026-11-07.md` | `c9804763e1cded31697a1ec05f3b9136a5dfab33dce5b5a4ebfb63f6035379a3` | 226 | 30,143 |
| Audited r4 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r4-2026-11-07.md` | `52a550853291e6d4b0db8473e7f4c78185f1a7414d06ae207c5363b12e754429` | 603 | 68,813 |
| Mimo r4 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r4-2026-11-07.md` | `eafaf6cc515fadaf80b34966ab92a8e5e5703db010f3e9bed1cc1517372aad41` | 258 | 30,900 |
| Pre-review assessed r5 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r5-2026-11-07.md` | `93cfe1de5f362fba6840056aec70aa7e29763ce38827a0a16b66666cbdb4b4b4` | 686 | 88,436 |
| Pre-review assessed r6 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r6-2026-11-07.md` | `8398802b35fecedae60f49f92dcf9f70cd51a9107cb49e432b4e4e4e9252eaf9` | 700 | 95,092 |
| Independently reviewed r7 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r7-2026-11-07.md` | `58cd4656c35407691f4b590d095c28340eea6710cc28debb34abde36f387e772` | 715 | 97,790 |
| Mimo r7 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r7-2026-11-07.md` | `8a1956bc07d863da1d26bc803314457d061e9aacee1936f90f67f2472bc705b4` | 159 | 18,600 |
| Independently reviewed r8 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r8-2026-11-07.md` | `b14dc4684605bed61eae83283578a37c34c979d88a362cdd99ff07821c24b4cf` | 724 | 102,800 |
| Mimo r8 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r8-2026-11-07.md` | `71ef9b725769d6539c09f9c9bf3bcaa79f66bbc181cffa1db76415da3c867905` | 157 | 17,863 |
| Governing contract | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 | 182,904 |
| Proposed stage plan | `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 | 69,168 |
| Correction record | `docs/review/CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` | `8fe9356cc304f3135657d19a7f2ea6badd9cf2233585a03bd0bd5c4f6355e076` | 172 | 16,190 |
| Owner-decision instrument | `docs/review/OWNER-DECISIONS-howl-tts-stage3-2026-11-07.md` | `24d3344554e7ed6a465c8b7bdbca4461cda23f78353aa7626e3cf4ba7df122ba` | 610 | 45,891 |

The r9 row label **“Pre-review assessed r6 predecessor”** is imprecise: r6 was the subject of the post-r6 restart assessment, not a pre-review. The row is reproduced verbatim to preserve the exact table; its path, SHA-256, line count, and byte count are correct and are the binding identification. This is the M3 note permitted in lieu of an r10. No requirement, hash, or artifact is changed.

## 3. Adoption fields — completed and executed

**Owner identity or accountable role:** Stone Wolf Systems (accountable organization).  
**Stage acceptance authority (§7.2):** Stone Wolf Systems, as owner. The owner did not author the substantive r9 changes, so no separate §7.2 accountable acceptor is required. Engineering produces technical evidence and assesses conformance; stage acceptance remains the owner's act.  
**Change provenance:** DeepSeek prepared the substantive r9 changes under ChatGPT's direction; the owner did not author them. No human author, signature, or separate approval requirement is inferred or created.  
**Adoption date:** 2026-11-08 (UTC; 2026-11-07 CST).  
**Artifact path:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r9-2026-11-08.md`  
**Artifact SHA-256:** `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`  
**Adopted scope — requirements coverage only:** The requirements and acceptance criteria applicable to C1–C8 and eventual implementation. This establishes requirements coverage only; it does not adopt the proposed stage plan, authorize any stage, or authorize implementation. A-1 is open, so PR-1–PR-4 remain inactive and this record claims no performance-acceptance effect.  
**Governing-input table (§2.2):** Reproduced in full in §2 above; no attachment is used.  
**Governing-input hashes — preparation-time verification (not adoption):** [x] Yes  [ ] No — 17/17 rows rechecked on 2026-11-07 CST (2026-11-08 UTC); this is a historical preparation check, not an adoption attestation.  
**Governing-input hashes — adoption-time recheck:** [x] Yes  [ ] No — all 17 rows were rechecked against §2 on 2026-11-08 (UTC; 2026-11-07 CST); every SHA-256, line count, and byte count matched.  
**Requirements/design incorporation:** [x] Incorporated as requirements and acceptance criteria only; proposed stage plan remains non-controlling.  [ ] Not incorporated; criteria remain review inputs only. — Selected by the owner on adoption.  
**Proposed-stage-plan status:** [x] Not amended  [ ] Amended under §6.5 Option B — instrument path: / exact scope: / SHA-256: / design incorporation is Yes: [ ]  
**Governing contract status:** [x] Matches §2.2; unamended; the contract has no C8 stage or C8 performance clause; Option B is not a contract amendment.  
**Correction-record status:** [x] Matches §2.2; unamended.  
**Contract-level C8 performance requirement (if elected):** [x] None  [ ] Elected — instrument recorded in the separate §17 field below; this revision does not authorize it.  
**Separate §17 contract-amendment instrument (if elected):** None. No contract-level C8 performance requirement has been elected, and no §17 instrument is part of this record.  
**Independent r9 review record / verified hash:** `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r9-2026-11-08.md`; SHA-256 `606acaec2e05a3c384c8dd6a028f0ac2d84a2e39c853e2d9bbe145ea415e942a`; 161 lines / 18,530 bytes; verdict **Adoptable as written.**  
**Requirements-authority matrix (§7.1) checked:** [x] Checked; no unresolved identifier collision or authority leakage; responsibilities separated: [x] evidence production  [x] conformance assessment  [x] stage acceptance  [x] requirement/architecture change  
**Budget-conformance role (§10):** [x] Bounded to existing authority; [x] no owner authority delegated; [x] determination is not stage acceptance  
**A-1 owner scope decision (if a performance gate has acceptance effect):** [ ] Decided — decision record/hash:  [x] Open — PR-1–PR-4 inactive  
**Approved exceptions (scope, duration, safety evidence):** No exceptions are approved by this instrument. Any later exception requires a separate post-adoption owner decision under §12, with scope, duration, and safety evidence; this field does not pre-approve any future decision.  
**Authorization boundary:** [x] Adoption does not authorize any stage C1–C8, implementation, runtime execution, benchmarking, or deployment.  
**Decision:** [x] Adopt  [ ] Reject  [ ] Return for correction — Stone Wolf Systems adopts r9 under the scope stated above.  
**Adoption method:** Recorded by the owner's declaration. No signature or separate approval field is required by this instrument.

## 4. M1–M3 and evidence-preparation notes

- **M1:** The exact governing-input table is reproduced in §2. No attachment path is offered as a substitute.  
- **M2:** The mandatory governing-contract and correction-record confirmations appear in separate fields in §3. The optional contract-level C8 §17 instrument appears only in its own field, and it is **None** here.  
- **M3:** The r6 label is noted in §2. The r6 artifact remains the post-r6 assessed input, identified by its unchanged hash and dimensions; no r10 is required for this record-level distinction.  
- **Future live-`unavailable` verification evidence (not a prerequisite to adoption):** When the applicable stage is authorized and its verification is performed, the test evidence must capture the complete supervisor tuple `ProcessState=alive`, `State=unavailable`, `Available=true`, `ServiceCapabilitiesKnown=true`, **and** the separate helper-side `State=unloaded` response. The helper-side observation does not substitute for any supervisor field. A helper-side `unloaded` reading is not a supervisor recovery; no recovery is recorded unless the explicit same-identity `load` required by A2-D12 occurs.  
- **Supervisor-level health rows:** The supervisor-level rows in r9 §8.1 are r9 behavioral requirements. They are not attributed to the contract's wire-health rule in §17.7.2.  
- **Approved exceptions:** No exceptions are approved by this instrument. Post-adoption exceptions are separate owner instruments and do not retroactively alter this record.

## 5. Review and verification record

| Check | Result |
|---|---|
| r9 actual-byte identity | 769 lines / 112,270 bytes; 769 CRLF pairs; SHA-256 `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` |
| r9 unchanged by this record | Yes — no r9 edit, no r10, no replacement artifact |
| Governing-input rows | Preparation-time 17/17 and adoption-time 17/17 recomputed on 2026-11-08 (UTC; 2026-11-07 CST); every SHA-256, line count, and byte count matches r9 §2.2. |
| Independent review | Audit SHA-256 `606acaec2e05a3c384c8dd6a028f0ac2d84a2e39c853e2d9bbe145ea415e942a`, 161 / 18,530; verdict Adoptable as written |
| §7.1 matrix | 27 data rows, six columns; four responsibilities separated |
| §7.3 original criteria | AC-1–AC-7 present and unchanged |
| §8.1 | Thirteen contract states; live/exited `unavailable` split preserved; two wire-health exclusions grounded in §17.4.1 and A2-D4 |
| Scope evidence | Documentation and hash verification only; no build, test, replay, benchmark, runtime execution, or deployment performed |
| Prepared-draft identity | Prior draft SHA-256 `d4917b5407f1f58d4a7559a7d7beda3c61e0d9a80bc794a05955bcbf3f63a6a3`; historical only, not the final adoption-record hash |
| Final adoption-record hash | Computed only after completion and recorded in the separate finalization note; it is not copied into this record |

## 6. Owner completion record

The owner completed the adoption on **2026-11-08 (UTC; 2026-11-07 CST)**:

1. **Owner/accountable role:** Stone Wolf Systems.
2. **Stage acceptance authority:** Stone Wolf Systems, as owner. The owner did not author the substantive r9 changes; no separate §7.2 acceptor is required. Engineering produces evidence and assesses conformance; the owner accepts stages.
3. **Incorporation:** r9 requirements and acceptance criteria only; the proposed corrective stage plan remains unamended and non-controlling.
4. **A-1:** remains open; PR-1–PR-4 remain inactive; no performance-acceptance effect is claimed.
5. **Governing inputs:** all 17 rows rechecked at adoption; every SHA-256, line count, and byte count matched.
6. **Exceptions:** none approved by this instrument.
7. **Decision:** **Adopt**.

**Finalization rule:** the final adoption-record SHA-256 is calculated only after this record is complete. It is recorded in a separate finalization note; it is not copied into this record. The `d4917b54…` value in §1.1 and §5 is historical provenance only.

**Effect:** r9 is adopted as the requirements and acceptance-criteria baseline. A-1 remains open; PR-1–PR-4 remain inactive; the eight Stage 3 decisions remain independent and pending; and C1–C8 and all implementation/runtime work remain unauthorized.
