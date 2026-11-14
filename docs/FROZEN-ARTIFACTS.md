# HOWL Chat — Frozen Artifact Manifest

**Status:** Integrity and placement manifest only. It creates no decision, contract, acceptance, or implementation authorization.
**Last verified:** 2026-11-12
**Rule:** A frozen file is immutable. If it is wrong, add a new dated decision or erratum; do not edit the frozen bytes.

---

## 1. What is frozen

| Artifact | Frozen scope | Date | SHA-256 | Sidecar / record |
|---|---|---|---|---|
| F3-C data and ownership contracts | Document, identity, path, reference, update, import/rollback, snapshot, R1–R3 | 2026-11-10 | `aa11be0047c9410a72a37334d758dbd4d05dfb60cd318477cd5b103157df40ca` | No sidecar; freeze record is the artifact |
| F3-C master architecture and implementation contract | Chat identity/root incarnation/revisions; commits/receipts; canonicalization; targets; locks/publication; operations; recovery; tombstones; memory/history migration; wipe/root identity; Gate v2 | 2026-11-11 | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` | No sidecar; freeze record is `FREEZE-F3C-ARCHITECTURE-2026-11-11.md` |
| F3-C architecture freeze | The architecture only; explicit exclusions listed in the record | 2026-11-11 | `2d9220e8a9cfeb15e621836892d89b5bae23bec4b7ce7c7555d5af0d62fb315a` | No sidecar |
| T2.9.5B r14 register | Exact r14 decision-control revision, including CL2 Option A closure and preserved gates | 2026-11-04 | `bd10519bd5e3611fd984c0835f3eac74b04e67818caa8055675967ec5f94cc44` | `t2.9.5b-authorization-decision-register-r14.md.sha256` |
| T2.9.5B r14 checksum | Two-space checksum line for r14 | 2026-11-04 | `bfbbdc68a143c85f56fcb47feeb2ec1aa20db967cc46278921b94687ad7604cd` | Outside the checksummed pair |
| T2.9.5B r14 freeze record | Freeze state, scope, verification, and limitations | 2026-11-04 | `e1d1df8180ecf266799a1a50040b96d71bae4a5516c776b23c22b76d802206d0` | No self-referential hash |
| TTS Stage 0 protocol | Stage 0 contract inside the governing design, as amended by adopted Clarification A2 | 2026-11-06; A2 adopted 2026-11-07 | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` (whole governing design) | No sidecar; adoption record is the amendment record |

**Important distinction:** "Frozen" never means "accepted implementation." The F3-C architecture is frozen; F3-C implementation is not accepted. The T2.9.5B register is frozen; the T2.9.5B stage is not authorized. The TTS Stage 0 protocol is frozen; no TTS implementation is authorized.

---

## 2. F3-C frozen artifacts

### 2.1 Data and ownership contracts

`docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`

Expected:

```text
aa11be0047c9410a72a37334d758dbd4d05dfb60cd318477cd5b103157df40ca
```

This freeze covers the document envelope, identity, paths, references, update semantics, import/rollback, snapshot semantics, and the three cross-cutting resolutions:

- **R1:** refresh vs cast sync are separate explicit choices;
- **R2:** canonical lore belongs to the item; generated lore belongs to the chat; legacy `automated_*` files are `legacy_generated_unattributed`;
- **R3:** item folders contain no binaries.

### 2.2 Master architecture and implementation contract

`docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`

Expected:

```text
61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31
```

This is the consolidated implementation authority. It supersedes the isolated Master-1 through Master-3c sequence as the consolidated authority. It does not accept F3-C, authorize P1a, or claim implementation evidence.

### 2.3 Architecture freeze

`docs/review/FREEZE-F3C-ARCHITECTURE-2026-11-11.md`

Expected:

```text
2d9220e8a9cfeb15e621836892d89b5bae23bec4b7ce7c7555d5af0d62fb315a
```

Frozen: chat identity and root incarnation; revision/commit records; receipts and missing-receipt rule; successor links and fork inventories; canonicalization and target hashes; target sets and ownership; locks/publication/operation map; creation/Save As/append/autosave/switch/delete/cast; journal recovery; tombstones; memory/history migration; wipe staging/tree hashing/root-identity commit; Gate v2.

Not frozen: F3-C implementation completeness; P1a authorization; F3-D; cast-confirmation implementation; A9/A10/A12 evidence for the new persistence contract.

**Do not edit any of these three files in place.** A change requires a new dated decision and a new frozen revision with its own record.

---

## 3. T2.9.5B frozen lineage

### 3.1 Current controlling pair

| Artifact | SHA-256 |
|---|---|
| `docs/t2.9.5b-authorization-decision-register-r14.md` | `bd10519bd5e3611fd984c0835f3eac74b04e67818caa8055675967ec5f94cc44` |
| `docs/t2.9.5b-authorization-decision-register-r14.md.sha256` | `bfbbdc68a143c85f56fcb47feeb2ec1aa20db967cc46278921b94687ad7604cd` |
| `docs/t2.9.5b-authorization-decision-register-r14.freeze.md` | `e1d1df8180ecf266799a1a50040b96d71bae4a5516c776b23c22b76d802206d0` |

The freeze record controls the apparent header conflict:

- the r14 register body says `DRAFT — REVISION 14 — NOT FROZEN — NOT AUTHORIZED`;
- the separate freeze record says `FROZEN — accepted controlling decision-control revision`;
- the freeze record freezes the exact r14 register and checksum bytes.

**Do not edit the r14 header to "fix" the conflict.** The freeze record is the later, specific control. `docs/AMENDMENT-REGISTER.md` and `docs/README.md` record the resolution.

The stage state remains:

```text
T2.9.5B stage: NOT FROZEN — NOT AUTHORIZED
Implementation authorization: NOT GRANTED
Measurement authorization: NOT GRANTED
```

### 3.2 Earlier frozen revisions

These remain byte-for-byte preserved. They are historical frozen revisions, not current alternatives to r14.

| Revision | Register SHA-256 | Sidecar SHA-256 | Freeze record SHA-256 |
|---|---|---|---|
| Base register | `286458c996a15e627aa589f699806a1406b0702f644f973230de830da27d4882` | `ff6b11b52f63ba952e16b871ead3949cd90d3bc6eded0e389bacf89b5f37cd7e` | `79371b0176dc30ca8f317878fb228c3fb7229f544b20167d391807338aa8df73` |
| r7 | `acf0ca6297c7a6a611f23d07904d5cd50cfda60e06e3dc4cbd83de2a62f88122` | `aeff93646ae4e95ec0cc6b1be4350abd6e6662a95f56c5859e4f48f39ef667a2` | `06c688a182db1713822fc3c247feb09b5b9840ca20733ddd84b697ff34595998` |
| r8 | `ee4f4278f994ee3e6438adb0417fc1a2bb9915cf6b073071ddfb935a6bf55119` | `07a2f03c605ca260dd25af75847c2f43411a2fd427b7e5fcccbdbcb5709e1a4f` | `e4cfe3e9ce68d2ae19abe6cfb00f6bb6a8ccf7bccdf26f05f03a694a22ded1ba` |
| r9 | `85a568d9f45cd0eee6ab64f3f9eae33540aa8aac54133e0276ae890ab619678c` | `398547522be4c48995fe654d7d1f4ff7606f4d919f7ae9db28d885542da6d634` | `9f63b2444fa6a599dfdb211887d54ad507a8b10824bfd14c679a82bc87d1ae7b` |
| r11 | `38a21308e6a6a8be85005c542b24492294a49027fabff8eb7a7ba37ecbc04f0e` | `65e637d62008976bab66b0d30b4bc2fee7636cf4cb882c17f3d243a58ecac8a5` | `3dbbe4d45a356a9208461f3884d30d9c6a01aed172afd631a01d50582b28fea4` |
| r13 | `cc228d1b8298dfd190f3fe371a43ec0510fb4b874097b5d35593993ac1714659` | `096d5e9aae5ed04495da6958bd57cf35812502955802b8920947442a0e178ff3` | `a4b8976851be091c19e40081aa999e75c9e5ba34e35c6b91e9ce799870a72055` |
| r14 | `bd10519bd5e3611fd984c0835f3eac74b04e67818caa8055675967ec5f94cc44` | `bfbbdc68a143c85f56fcb47feeb2ec1aa20db967cc46278921b94687ad7604cd` | `e1d1df8180ecf266799a1a50040b96d71bae4a5516c776b23c22b76d802206d0` |

Preserved but **not frozen** by a separate freeze record: r2, r3, r4, r5, r6, r10, r12. They have checksum sidecars but no freeze record. Do not treat a sidecar as a freeze.

---

## 4. TTS Stage 0 frozen protocol

The frozen protocol is **inside** the governing design, not a separate file:

`docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`
SHA-256: `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f`

Stage 0 was authorized and frozen for review on 2026-11-06. Clarification A2 was adopted on 2026-11-07 from accepted Revision 6:

- adoption record: `docs/review/ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md`
- adoption SHA-256: `9da36691054df40d98c772f7ba9643bc959416e39a7f2d30681d9404cd16f043`
- standalone draft: `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`
- draft SHA-256: `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784`

The draft's `NOT ADOPTED` header is its pre-adoption state. The adoption record controls. The draft remains historical; §17.19 of the governing design is the governing text.

This freeze is **design/protocol authority only**. It authorizes no implementation, dependency addition, build, launch, Windows spike, inference, model/voice operation, production edit, commit, or push.

---

## 5. What is not frozen

Do not mistake any of these for a frozen artifact:

| Artifact | Actual status |
|---|---|
| Amendment 01 | Current planning amendment; not frozen |
| Amendment 02 revision 9 | Closed requirements baseline; not frozen; no rev 10 |
| Amendment 02 revisions 6/7/8 | Archived revisions |
| TTS §17.19 standalone draft | Historical draft; content adopted into the governing design |
| F3C-TARGET-2 R1–R4 | Historical proposals |
| F3C-TARGET-2 R5 | Current proposal only; not frozen |
| `GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` | Adopted protocol; not a freeze and not an authorization |
| `CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md` | Proposed revision 3; not adopted |
| F3-C implementation evidence | Evidence only; F3-C not accepted |
| ROOT-01 / Gate B / D-ROOT closure | Source evidence and open defects; not authorization |
| `docs/audio-architecture/*` | Implementation architecture; not T2.9.5B authorization |

The corrected master plan's internal phrase "frozen technical scope/dependency" describes its internal revision scope. It is **not** an owner adoption and does not create a frozen external contract.

---

## 6. Verification commands

Run from the repository root.

### 6.1 T2.9.5B current pair

```bash
sha256sum -c docs/t2.9.5b-authorization-decision-register-r14.md.sha256
# expected: docs/t2.9.5b-authorization-decision-register-r14.md: OK
```

Then:

```bash
sha256sum docs/t2.9.5b-authorization-decision-register-r14.md
sha256sum docs/t2.9.5b-authorization-decision-register-r14.md.sha256
sha256sum docs/t2.9.5b-authorization-decision-register-r14.freeze.md
```

Expected:

```text
bd10519bd5e3611fd984c0835f3eac74b04e67818caa8055675967ec5f94cc44  docs/t2.9.5b-authorization-decision-register-r14.md
bfbbdc68a143c85f56fcb47feeb2ec1aa20db967cc46278921b94687ad7604cd  docs/t2.9.5b-authorization-decision-register-r14.md.sha256
e1d1df8180ecf266799a1a50040b96d71bae4a5516c776b23c22b76d802206d0  docs/t2.9.5b-authorization-decision-register-r14.freeze.md
```

### 6.2 Earlier frozen revisions

```bash
for f in \
  docs/t2.9.5b-authorization-decision-register.md \
  docs/t2.9.5b-authorization-decision-register-r7.md \
  docs/t2.9.5b-authorization-decision-register-r8.md \
  docs/t2.9.5b-authorization-decision-register-r9.md \
  docs/t2.9.5b-authorization-decision-register-r11.md \
  docs/t2.9.5b-authorization-decision-register-r13.md; do
  sha256sum -c "$f.sha256" || exit 1
done
```

### 6.3 F3-C frozen artifacts

```bash
sha256sum \
  docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md \
  docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md \
  docs/review/FREEZE-F3C-ARCHITECTURE-2026-11-11.md
```

Expected:

```text
aa11be0047c9410a72a37334d758dbd4d05dfb60cd318477cd5b103157df40ca  docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md
61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31  docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md
2d9220e8a9cfeb15e621836892d89b5bae23bec4b7ce7c7555d5af0d62fb315a  docs/review/FREEZE-F3C-ARCHITECTURE-2026-11-11.md
```

### 6.4 TTS governing design and adoption

```bash
sha256sum \
  docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md \
  docs/review/ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md \
  docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md
```

Expected:

```text
ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f  docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md
9da36691054df40d98c772f7ba9643bc959416e39a7f2d30681d9404cd16f043  docs/review/ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md
dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784  docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md
```

---

## 7. Change control

A frozen artifact changes only by:

1. a new dated owner decision that states what is changing and why;
2. a new frozen revision, not an in-place edit of the old one;
3. its own freeze record and checksum;
4. a new row in `docs/DOCUMENT-REGISTRY.md` and `docs/AMENDMENT-REGISTER.md`;
5. an explicit supersession statement linking old and new.

A hash mismatch is an integrity incident, not a reason to regenerate the sidecar silently.

---

## 8. Bottom line

Frozen artifacts and open authority are now separate in one place:

- F3-C architecture: **frozen**; implementation: **not accepted**.
- F3-C data/ownership: **frozen**.
- T2.9.5B r14: **frozen**; stage: **not authorized**.
- TTS Stage 0: **frozen protocol**; implementation: **not authorized**.
- Amendment 01/02, TARGET-2 R5, corrected master plan: **not frozen**.

No frozen byte was changed by this reconciliation.
