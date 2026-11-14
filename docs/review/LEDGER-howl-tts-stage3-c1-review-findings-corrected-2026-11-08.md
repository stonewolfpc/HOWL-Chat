# C1 — Corrected Finding Ledger

**Date:** 2026-11-08
**Status:** AUTHOR LEDGER. Not closure. **C1 REMAINS UNAUTHORIZED.**
**Supersedes the ledger section of:** reconciliation `5b72cc97…` (which is frozen; this file is the correction).
**Attestation:** separate artifact, unsigned.

## 1. The ambiguity being corrected

The previous ledger marked some rows "conceded" in prose while marking their IDs `TEXT REQUIRED`. That reads as a contradiction. The resolution is to separate two different things:

- **Mechanism concession** — the source code exhibits the failure the claim describes, so the author accepts the *mechanism*, whether or not the reviewer's exact wording is available.
- **Finding-level acceptance** — the author accepts the reviewer's scope, severity, and falsifier. This requires the reviewer's text.

A row can therefore be "mechanism conceded, finding-level text pending" without ambiguity.

## 2. Corrected ledger

| ID | Author response | Basis | Missing text | Consequence |
|---|---|---|---|---|
| E1-05 | mechanism conceded | source: `supervisor.go:172` authenticates by pointer | reviewer wording | keep pointer token as authentication; epoch is a label |
| E1 intent gap | mechanism conceded | source: `SetState` `:233–242` ignores intent | ID mapping | latch `retiring` atomically |
| E1 handshake conflict | mechanism conceded | source: `client.go:248–253` publishes after terminate | reviewer scope | observer owns typed terminal (option a) |
| E1-01…E1-04 | **cannot respond** | none | full finding text | reserved; no correction asserted |
| L-7 check-then-act | mechanism conceded | source: `client.go:860–861` | reviewer wording | `TerminateIf(token, intent)` |
| L-7 kill paths | mechanism conceded | source: `client.go:305`, `:312`, `:366` | reviewer scope | `KillIf(p, intent)` |
| L-7-03 | mechanism conceded | source: one liveness test gates both process and file `:860–863` | reviewer severity | split resource classes |
| L-7-01/02/04–06 | **cannot respond** | none | full finding text | reserved |
| L-8 | open | source: no intent field | reviewer semantics | owner decision on precedence |
| L-9 Start race | mechanism conceded | source: `client.go:213–216` then `:226` | reviewer wording | admission inside `s.mu` install |
| L-9 shutdown frame | mechanism conceded | source: `Shutdown` → `call` → `Start` | reviewer scope | one-shot control permit |
| Q6 | accepted with correction | source: edge list omitted release points | none | bounded wording retained |

## 3. Rules applied

1. A mechanism concession cites a line; without a line, the row is `cannot respond`.
2. A finding may be conceded in mechanism and still be open at finding level; the two columns are independent.
3. No row is closed by this ledger. Closure requires reviewer confirmation in the attestation.
4. Reviewer-confirmed ≠ author-agreed ≠ accepted amendment.

## 4. Status

Reserved IDs remain unanswered pending text. Frozen artifacts unchanged. No source, test, build, or r9 change.
