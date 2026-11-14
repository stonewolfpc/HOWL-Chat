# `RCN-1-E1` — Erratum to `MEMORY-09-RCN-1` (count terminology)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `RCN-1-E1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **EFFECTIVE — terminological only.**
**Erratum to:** `docs/research/corpus/MEMORY-09-RECONCILIATION-1-ACCOUNTING-R1-WITHHOLDING-2026-11-13.md` (`MEMORY-09-RCN-1`, v1, `b05796e4…`, unchanged).
**Authority:** owner approval, 2026-11-13, of `MEMORY-09-SDP-1` v2 Appendix A (decision `D-6`). Companion to the `CORPUS-01` amendment made the same day.

---

## 1. What is corrected

`MEMORY-09-RCN-1` §2.1 described **46** as *"every level-2/3 item heading"*. That is a **raw heading-line count** and conflicts with the convention already stated in `CORPUS-01` ("Accounting basis": **42 item headings**, with the four `R6-03a`–`d` variants folded under the `R6-03` family heading).

## 2. Corrected terminology

| Count | Term | Composition |
|---:|---|---|
| 42 | **item headings** | One per family-level item; includes the `R6-03` container |
| 45 | **scored units** | 42 − 1 unscored container + 4 scored variants |
| 46 | **heading lines** (equivalently, **tagged headings**) | 42 + 4 variant sub-headings; each line carries a `Restraint:` tag |

## 3. What is not affected

- No count changes value.
- No item, tag, membership, split, score or threshold changes.
- The findings in `MEMORY-09-RCN-1` §2.2–§2.3 stand: the container has no scored body, the 45/46/42 relation holds, and the six separate `R-6` counts are unchanged.
- `MEMORY-09-RCN-1` itself remains **byte-identical**; this erratum is a separate, additive record.

## 4. Effect

The two artifacts now agree on terminology, and "46 headings" no longer appears anywhere as a convention name. The original reconciliation remains citable in full, with this erratum read alongside it.

## 5. Provenance

| Layer | Fact |
|---|---|
| Authorization | Owner, 2026-11-13 (`D-6`) |
| Applies to | `MEMORY-09-RCN-1` §2.1 terminology only |
| Filed by | DeepSeek (per work order), repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Companion change | `CORPUS-01` accounting basis, coverage line, operational membership block, `R-6` restatement (same date) |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
