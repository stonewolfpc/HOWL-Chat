# Adoption Record — §17.19 TTS Capability Policy Clarification A2

**Date:** 2026-11-07  
**Status:** ADOPTION COMPLETE — GOVERNING ARCHITECTURE CONTRACT  
**Scope:** document adoption only. This record does not authorize implementation, Stage 3, code,
capability negotiation, caller wiring, a recovery UI, a helper build, a native/provider dependency,
model or voice assets, fixtures, tests, or any change to accepted Stage 1 or Stage 2 artifacts.

## 1. Authority and verified inputs

| Artifact | Path | SHA-256 | Lines | Bytes |
|---|---|---|---:|---:|
| Accepted amendment (Revision 6) | `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` | `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784` | 320 | 65,274 |
| C-1 decision report | `docs/review/DECISION-howl-tts-capability-policy-c1-2026-11-06.md` | `c25a0d3fd38e0767a4b7bf278cb1d6617e2e53caf11053dabb71449578c4c8ad` | 302 | 26,197 |
| Independent Revision 6 audit | `docs/review/AUDIT-howl-tts-capability-policy-c1-rev6-2026-11-07.md` | `19bd7202a605600f23cc7556a1dc899fc97365eff13f245ba0483fd345ff91e2` | 329 | 21,506 |
| Architecture, before adoption | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` | 2,416 | 139,351 |
| Architecture, after adoption | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 | 182,904 |

The pre-adoption architecture hash was verified before the first adoption edit. The accepted
amendment, C-1 decision, and Revision 6 audit were verified before editing and are unchanged after
adoption. The Revision 6 audit verdict was **PASS — recommend adoption. No required fixes.**

## 2. Exact sections changed

| Location | Change | Authority |
|---|---|---|
| Document header | Status changed to governing contract; Stage 0 status records A2 adoption; adoption record and verified authority/audit hashes added. | A2-D11 |
| §17.1 item 11 | Added the A2 adoption entry after the preserved A1 entry. | A2-D11 |
| §17.7.2 | Added `handshaking` + semantic `hello` failure → `incompatible`; added `unavailable` + same-identity explicit `load` → `loading`. Existing rows preserved. | A2-D3/A2-D4/A2-D12 |
| §17.7.3 | Added `unloaded` to the legal `loading` transitions and recorded `accepting_work=true` for a typed non-fatal load failure. | A2-D12 |
| §17.7.4 | Added the two `unavailable` health rows and a note limiting the existing startup-dependency wording. Existing rows preserved. | A2-D5/A2-D12 |
| §17.12 | Added `Loaded`, `ServiceCapabilities`, and `ServiceCapabilitiesKnown`; extended application `TTSCapability`; recorded provider/service scope separation. | A2-D6 |
| §17.17 | Added the A2 checklist row and replaced the stale “acceptance pending” footer with the adoption status. | A2-D11 |
| §17.19 | Appended the adopted clarification, including its supersession map, state/status rules, and scope confirmation. | All A2 decisions |

No other sections were rewritten. The pre-adoption decisions remain in place as historical evidence;
§17.19 controls the matters listed in its supersession map. Line references inside §17.19 that point
to §§17.1–17.18 are explicitly identified as pre-adoption historical locators.

**Changed-file inventory:** the only files written by this adoption task are the governing
architecture document, `docs/review/README.md`, and this adoption record. The amendment, C-1 decision,
Revision 6 audit, and all other historical review artifacts were not modified.

## 3. Decision-to-contract mapping

| Decision | Adopted contract location |
|---|---|
| A2-D1 — Capability floor | §17.19.1, §17.19.2 A2-D1, §17.19.3; §17.4.1 and §4.6 confirmed/clarified. |
| A2-D2 — Wire and application semantics | §17.19.1, §17.19.2 A2-D2, §17.19.3; §17.3.2 superseded for false-capability wording; §17.5.1 wire DTO unchanged; §17.12 projection. |
| A2-D3 — Hello semantic validation | §17.19.2 A2-D3, §17.19.3; §17.7.2 failure row; §17.9, §5.3, §4.6, and §17.18.1 A1-D7 confirmed. |
| A2-D4 — `incompatible` transition | §17.19.2 A2-D4, §17.19.4 item 2; §17.7.2 row; §17.7.1 and §17.11.5 clarified. |
| A2-D5 — `Available`, `Loaded`, `Ready` | §17.19.2 A2-D5, §17.19.4 items 3/5; §17.7.4 rows; §17.7.1 supersession; §17.12 fields. |
| A2-D6 — Service-capability status | §17.19.2 A2-D6, §17.19.4 item 10; §17.12 authoritative status shape and `TTSCapability`; §3.2 superseded as authority. |
| A2-D7 — Absent `unload` / identity change | §17.19.2 A2-D7, §17.19.4 item 7; §17.4.4 superseded for the no-`unload` case; §17.11.5 budget rule clarified. |
| A2-D8 — Idle close uses `shutdown` | §17.19.2 A2-D8, §17.19.4 item 8; existing §17.7.2 row and §17.8 value confirmed; no duplicate row. |
| A2-D9 — Cancellation without `cancel` | §17.19.2 A2-D9, §17.19.4 item 9; §17.11.4 steps 4–7 and §17.11.6 rule 4 clarified. |
| A2-D10 — No silent downgrade or no-op | §17.19.1, §17.19.2 A2-D10; §4.6 optional-capability scope confirmed. |
| A2-D11 — Amendment method / deferred issues | §17.19.1, §17.19.2 A2-D11, §17.19.3; §17.1 item 11 and §17.17 updated; B12a/B12b preserved as deferred. |
| A2-D12 — `unavailable` recovery | §17.19.2 A2-D12, §17.19.4 items 5/11; §17.7.2 retry row; §17.7.3 edge; §17.7.4 rows; §17.12.1 and §17.11.5 confirmed. |

## 4. Verification commands and results

| Check | Command | Exit |
|---|---|---:|
| Authoritative hashes before editing | `sha256sum` on the accepted amendment, C-1 decision, and Revision 6 audit | 0 |
| Pre-adoption architecture hash | `sha256sum docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | 0 |
| Structural/policy checks | Python assertions over required/optional capabilities, absence semantics, fail-closed behavior, `protocol_incompatible`, ordinary `unloaded` ensure-load, `unavailable` recovery, identity triple, A2-D7 routing, no budget inheritance, `detail.missing`, 12 decisions, thirteen states, status fields, B12a/B12b, stale wording, and code fences | 0 |
| Final architecture hash | `sha256sum` and Python `hashlib.sha256` independently | 0 |
| Final counts | `wc -l -c` | 0 |
| Index markers | Python assertions for §21/§22 and governing-contract status | 0 |

Post-adoption hash was independently computed with both `sha256sum` and Python's `hashlib` and
matched byte-for-byte.

## 5. Deviations, omissions, and unresolved items

- No accepted policy decision was reopened, simplified, or omitted.
- The only editorial deviations are the E5 cleanup (plural “these additive rows” and removal of a
  stale line-number locator), the historical-locator note, and the explicit B12a/B12b deferral
  sentence. None changes policy.
- O-1 through O-7 remain explicitly unresolved.
- B12a (hard-link versus atomic-rename in §17.10.3) and B12b (`unsupported_extension` versus
  `unsupported_capability` wording) remain deferred.
- I1 remains a non-blocking integration observation: no frozen adapter method sends an explicit
  `load`, so in-session recovery depends on the update-gated reset or instance re-creation. It is not
  a contract defect and does not authorize implementation.
- No contradiction prevented faithful adoption.

## 6. Governing-contract confirmation

The post-adoption `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` is now the governing
architecture contract for the matters §17.19 identifies. It remains design-only. Stage 3 and all
implementation work remain unauthorized. No commit, push, or staging was performed.
