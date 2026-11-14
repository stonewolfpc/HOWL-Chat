# GATE — F3-C Implementation Verification v2

**Date:** 2026-11-11  
**Status:** **ADOPTED — VERIFICATION PROTOCOL ONLY.**  
**Authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Supersedes:** the original six-step implementation gate.  
**Does not grant:** P1a implementation authorization, F3-C acceptance, verification, deployment, or F3-D authorization.

Gate v2 is the implementation-verification protocol for F3-C. It tests both the implementation and the frozen contract. It exists because an implementation can pass every test, and faithfully follow the contract, while remaining unsafe if the contract or its acceptance criteria are defective.

The architecture freeze remains in force. P1a implementation is still not authorized. F3-C acceptance is still not granted. F3-D remains blocked. A9/A10/A12/A13 remain closed.

---

## 1. Authority binding

1. The only implementation authority is the consolidated contract at the exact path above.
2. The recorded SHA-256 must match the contract presented in the review packet. A mismatch stops the gate; no review result is valid against a different document.
3. Master-1 through Master-3c are decision history, not competing authorities.
4. Frozen F3-C decisions, LF-1, A1–A13, and accepted native evidence remain supporting authority, subject to the limits recorded in the consolidated contract.
5. A contract change after this gate record requires a new contract hash and a new Gate v2 authority check. The previous verification cannot be reused.

---

## 2. Review packet and independent evidence access

The implementer supplies one frozen packet. The reviewer must receive the complete implementation, not files selected by the implementer.

The packet includes:

- the exact contract path and SHA-256;
- the packet commit, branch, build identifier, and dates;
- a complete changed-file manifest, including added, modified, deleted, renamed, and generated files;
- the dependency closure needed to build, test, and run the changed code;
- all source, tests, fixtures, schemas, migrations, build scripts, and native harnesses;
- all test commands, toolchain versions, environment details, and raw command logs with exit statuses;
- all fault-injection records, race-test records, filesystem-recovery records, and native Wails/WebView2 captures;
- the evidence ledger described in §10;
- a list of every redaction, with its reason; secrets are never included in the packet.

The reviewer may inspect any repository file, not only the changed-file manifest. The manifest defines completeness; it does not limit access.

Evidence is invalid if:

- commands are filtered so a failure or exit status is hidden;
- logs are summaries without raw output;
- a test name is offered without its source, assertion, input, and observed result;
- a result depends on an unpinned environment or an unavailable artifact;
- the implementer controls the only copy of the evidence and refuses independent inspection.

---

## 3. Finding taxonomy

Every finding receives exactly one class. A finding may expose more than one class, in which case separate findings are recorded.

| Class | Definition | Required treatment |
|---|---|---|
| **Contract defect** | The frozen requirement is ambiguous, contradictory, incomplete, or unsafe. | Do not implement around it. Record the exact clause, a counterexample, and the impact. The contract must be corrected by a new freeze or formally dispositioned before verification can pass. |
| **Implementation defect** | The code violates a sound requirement. | Fix the code, add the missing positive and negative coverage, and rerun the affected and regression checks. |
| **Evidence defect** | The implementation may be correct, but the proof is missing, incomplete, irreproducible, or does not cover the claim. | Produce the missing evidence or record the item as unresolved. An evidence defect is never converted into a pass by assertion. |
| **Test defect** | The test does not meaningfully establish the claimed behavior. | Repair the test, add the adversarial case, or replace the test. A passing weak test is not evidence. |

A contract defect is not dismissed as an implementation deviation. A weak test is not treated as proof of correctness.

Findings also carry:

- a stable finding ID;
- severity: blocking, major, or minor;
- the contract clause or source location;
- the claim under review;
- the evidence and the adversarial counterexample;
- the observed result and impact;
- disposition: **confirmed**, **rejected**, or **unresolved**;
- the owner and the required regression.

Confirmed, rejected, and unresolved findings remain distinguishable in the final record. Unresolved findings never count as passed.

---

## 4. Status ladder

The four statuses are separate. No status is inferred from another.

| Status | Asserted by | Requires | Does not imply |
|---|---|---|---|
| **Implementation complete** | Implementer | The implementer claims all planned work is done and supplies the frozen packet. | Correctness, verification, or acceptance. |
| **Verification passed** | Gate v2 reviewer | Every §11 row is resolved; no blocking contract, implementation, evidence, or test defect remains; regression is complete; the staged Mimo challenge is complete. | F3-C acceptance or F3-D authorization. |
| **F3-C accepted** | Accountable owner | Verification passed and an explicit owner acceptance decision. | F3-D authorization or P1a implementation authorization. |
| **F3-D authorized** | Separate owner authorization | F3-C acceptance and an explicit F3-D scope authorization. | Any other downstream work. |

P1a implementation authorization is a separate planning and authorization decision. It is not granted by implementation completeness, verification, acceptance, or this protocol.

---

## 5. Procedure

### Stage 0 — Authority and packet freeze

1. Record the contract path and SHA-256.
2. Record the packet commit, build identifier, toolchain, OS, and test environment.
3. Confirm reviewer independence and read-only evidence access.
4. Confirm the packet includes the complete changed-file manifest and dependency closure.
5. Stop on authority mismatch, incomplete packet, or uncontrolled evidence.

### Stage 1 — Contract challenge

Challenge the contract before reviewing the implementation against it.

The reviewer asks whether the requirement itself is ambiguous, contradictory, incomplete, or unsafe. Required challenge areas include:

- identity, root incarnation, and revision semantics;
- hash construction, canonical bytes, and any circularity;
- target set and ownership;
- receipt and commit inventories, reconstruction, and forks;
- lock order, publication order, and crash windows;
- creation, Save As, append, autosave, switch, delete, and cast boundaries;
- ordinary and cast journal recovery;
- memory and history migration;
- tombstones and deletion evidence;
- wipe staging, tree hashing, gallery boundary, and the single root-identity commit point;
- acceptance criteria and the evidence they demand.

A contract defect is recorded as such. The implementation cannot cure it by choosing an interpretation. If the defect changes architecture, Gate v2 pauses until a new contract freeze exists.

### Stage 2 — Implementation review

1. Review the complete packet and dependency closure.
2. Trace every authoritative mutation caller, not merely the public API list.
3. Verify the lock order and the leaf-lock rule.
4. Verify target serialization, exact bytes, readback, and hashes.
5. Verify commit records, receipts, inventories, successor links, and reconstruction.
6. Verify creation, Save As, append, autosave, switch, delete, cast, migration, wipe, and gallery paths.
7. Verify error paths, poison states, and recovery outcomes.
8. Verify that tests and fixtures do not use a weaker path than production.
9. Record implementation defects with exact source and clause references.

### Stage 3 — Test adequacy

Inspect what the tests assert, not only whether they pass.

Required review:

- each assertion and the behavior it establishes;
- positive and negative cases for critical invariants;
- boundary values, malformed input, and conflicting operations;
- fault injection at journal boundaries;
- deterministic race coverage and scheduler assumptions;
- filesystem cases: missing, partial, duplicate, symlink, reparse point, hard link, and unexpected object;
- native versus mocked or shimmed execution;
- skipped, retried, flaky, or conditionally disabled tests;
- whether a test can pass while the claimed invariant is false;
- whether fixtures use real user data or only synthetic data.

A test that asserts only “no error” is insufficient for a semantic invariant. A test that passes through a mock cannot establish native Wails/WebView2 behavior.

### Stage 4 — Adversarial and destructive verification

Apply the §11 matrix with both the valid operation and the dangerous operation.

- Prove the valid operation succeeds.
- Prove the dangerous operation is rejected, blocked, or safely recovered.
- Attempt the crash window, stale base, changed session, duplicate operation, fork, malformed receipt, ambiguous journal, missing tombstone, bad tree hash, and interrupted migration.
- Exercise the wipe, delete, rollback, and recovery paths only in an isolated destructive environment.

### Stage 5 — Full regression after corrections

Any correction invalidates prior results for the affected area.

- A changed lock, lock order, serializer, target path, receipt format, recovery path, or journal boundary requires the full §11 matrix.
- A narrower change requires the affected rows plus the full regression set identified in the finding.
- No result is carried across a different packet hash, contract hash, or environment.
- The regression result is recorded, not assumed.

### Stage 6 — Adjudication

1. Reconcile every finding as confirmed, rejected, or unresolved.
2. Confirm that every §11 row has a result.
3. Confirm that all blocking findings are closed or explicitly unresolved.
4. Confirm the staged Mimo challenge is complete.
5. Issue the Gate v2 verdict in §9.
6. Do not issue an acceptance decision. That is a separate owner action.

### Stage 7 — Owner decision

The accountable owner may accept F3-C only after a verification-passed record. Acceptance is recorded separately from the Gate v2 verdict. F3-D authorization remains a further separate decision.

---

## 6. Required §11 matrix

Every item in the contract’s §11 must have a row. A missing row is `NOT TESTED`, not a pass.

Each row uses this chain:

```text
claim -> evidence -> adversarial counterexample -> result
```

Required row fields:

- contract §11 item and invariant;
- explicit claim, including the expected success and the expected rejection;
- evidence artifact, command, source location, hash, and observed output;
- adversarial counterexample and why it is the strongest practical attack;
- result: **PASS**, **FAIL**, **UNRESOLVED**, or **BLOCKED**;
- defect class and finding IDs, if not PASS;
- regression scope and rerun record.

A row is PASS only when the positive case and the negative case both behave as required, the evidence is reproducible, and the counterexample has been attempted and rejected.

### 6.1 Row inventory

| ID | Contract §11 item | Required focus |
|---|---|---|
| V1 | 1 | Canonicalization, request fingerprints, Base64, hash chain. |
| V2 | 2 | Receipt inventory and reconstruction. |
| V3 | 3 | Commit fork, orphan, and missing inventory. |
| V4 | 4 | Exact target bytes, `tgt1-` manifest, memory snapshot, history identity. |
| V5 | 5 | Append, dirty state, autosave, session invalidation, stale snapshots. |
| V6 | 6 | Locks, publication, autosave vs switch, wipe vs save. |
| V7 | 7 | Save As source binding and operation-ID conflicts. |
| V8 | 8 | Tombstones and deletion recovery. |
| V9 | 9 | Creation journal boundaries and identity revalidation. |
| V10 | 10 | Memory migration conflict, commit, and receipt/journal agreement. |
| V11 | 11 | History migration order and ID preservation. |
| V12 | 12 | Wipe tree hashes, gallery boundary, staging failure, root-identity commit. |
| V13 | 13 | Fault injection at every journal boundary. |
| V14 | 14 | Filesystem recovery: symlinks, reparse points, hard links, partial staging. |
| V15 | 15 | Race testing. |
| V16 | 16 | Native Wails/WebView2 paths. |
| V17 | 17 | Independent Mimo Flash Fast challenge before final approval. |

No row may be satisfied by a test name alone. No row may be satisfied by a positive case alone when the invariant is protective.

---

## 7. Isolated destructive testing

Wipe, delete, rollback, recovery, and root-incarnation tests must never risk real user data.

Required isolation:

- a dedicated disposable root, VM, container, or equivalent sandbox;
- synthetic fixtures only;
- an explicit check that the harness root is disposable before any destructive operation;
- no production state, no live user directory, and no shared application data;
- no destructive command outside the disposable sandbox;
- recorded cleanup and proof that the test did not reach a real root;
- fail-closed behavior if isolation cannot be proven.

An isolated destructive test that cannot prove its boundary is `BLOCKED`, not passed.

---

## 8. Staged Mimo challenge

Mimo receives the complete frozen contract and the complete implementation packet. The challenge is staged so each category gets an independent pass.

### Mimo Stage 1 — Contract and acceptance criteria

Mimo challenges the contract itself: ambiguity, contradiction, missing failure outcome, unsafe requirement, and defective acceptance criteria. Findings are recorded as contract defects or rejected with reasons.

### Mimo Stage 2 — Implementation against the contract

Mimo challenges the implementation: deviations, missing callers, lock violations, unsafe recovery, incorrect serialization, and unsupported assumptions. Findings are recorded as implementation defects or rejected with reasons.

### Mimo Stage 3 — Test evidence and claimed results

Mimo challenges the proof: weak assertions, missing negative cases, unreproducible commands, mocked native paths, hidden skips, and claims not supported by the evidence. Findings are recorded as evidence or test defects or rejected with reasons.

Mimo does not receive a single enormous prompt that combines all three stages. Mimo does not grant acceptance. Mimo findings enter the same taxonomy and adjudication process as every other finding.

---

## 9. Gate verdict

The Gate v2 record uses one verdict.

| Verdict | Meaning |
|---|---|
| **VERIFICATION PASSED** | All 17 §11 rows are PASS; no blocking contract, implementation, evidence, or test defect remains; all findings are adjudicated; regression is complete; the three Mimo stages are complete. |
| **VERIFICATION FAILED** | At least one confirmed blocking defect or failed §11 row remains. |
| **VERIFICATION INCOMPLETE** | At least one §11 row, evidence item, test-adequacy question, or Mimo stage remains unresolved or missing. This is not a pass. |
| **BLOCKED** | Authority, packet completeness, reviewer independence, or destructive-test isolation cannot be established. |

`VERIFICATION PASSED` is an engineering verdict. It is not F3-C acceptance. It does not authorize P1a implementation. It does not authorize F3-D.

---

## 10. Evidence ledger and final record

The final Gate v2 record contains:

- contract path and SHA-256;
- packet commit, build identifier, and environment;
- changed-file manifest hash and dependency-closure record;
- the complete §11 matrix;
- every finding with class, severity, evidence, counterexample, disposition, owner, and regression;
- confirmed/rejected/unresolved totals;
- the isolation proof for destructive tests;
- the three Mimo stage records;
- the regression record;
- the verdict and the reviewer identity;
- the date and the packet hash to which the verdict applies.

The record must be sufficient for a different reviewer to reproduce the result without asking the implementer to select new evidence.

---

## 11. Adoption and boundaries

Gate v2 is adopted as the F3-C implementation-verification protocol. It replaces the original six-step gate. It does not replace the consolidated contract, the architecture freeze, or the F3-C acceptance decision.

Current boundaries remain:

- **Architecture freeze:** granted.
- **Gate v2:** adopted.
- **P1a implementation authorization:** not granted.
- **F3-C acceptance:** not granted.
- **F3-D:** blocked.
- **A9/A10/A12/A13:** closed.

The next permissible step is a separately authorized P1a planning exercise. This protocol does not grant implementation permission.
