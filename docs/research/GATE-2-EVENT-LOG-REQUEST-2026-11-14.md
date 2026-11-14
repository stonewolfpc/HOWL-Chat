# Gate 2 — Event-Log Query Request

**Identity:** `HOWL-GATE2-REQ-1` · **Date:** 2026-11-14 · **Mode:** request only, nothing executed
**Gate:** **2 — event-log enumeration on `C:`.** Closed. This document asks for it; it does not open it.
**Storage access:** none performed. No volume read, no model read, no probe.

---

## 1. Why this is a gate and not a free action

Reading the Windows event log is a query against `C:`, executed on the same machine and the same stack whose behaviour is under investigation. It is also the cheapest remaining source of positive evidence about **system stability** — whether the old failure pattern stopped, and whether anything new started.

Cheap is not the same as authorised. Gate 2 was listed closed for that reason, and a query that returns thousands of rows is not self-limiting the way a bounded file read is.

**What this query can establish:** whether WHEA, storage-reset, paging-error and bugcheck events stopped after the driver fix and after `E:` was removed, and whether their onset correlates with the failure window.

**What it cannot establish — and this is the whole boundary:** nothing about the digest disagreement. Event logs record that a fault was reported, not what any read returned. A clean log does not exonerate the storage path; a dirty one does not explain a digest. If this runs, it answers a **stability** question and sits beside the hash investigation, not inside it.

---

## 2. The query, with corrections

The plan supplied in the companion session proposes a 30-day System-log scan for IDs `17, 18, 19, 41, 51, 129, 153, 4101, 1001`. Two corrections matter for this investigation.

**Correction A — the storage-integrity events are missing.** `7`, `9` and `11` are the disk-level events (bad block, device not responding, controller error). `51` is the paging-operation error, which is directly relevant to the `E:` page-file hypothesis. A scan that omits `7/9/11` cannot speak to storage integrity, only to paging and resets.

**Correction B — uncorrected WHEA is missing.** `17/18/19` are the *corrected* hardware-error events. `46/47`, where present, are the ones that indicate uncorrected or fatal hardware faults and are what would point at CPU/IMC or PCIe.

```powershell
# Gate 2 — read-only. 30 days covers before and after E: removal (2026-11-13).
$since = (Get-Date).AddDays(-30)

Get-WinEvent -FilterHashtable @{ LogName = 'System'; StartTime = $since } |
Where-Object { $_.Id -in 7,9,11,17,18,19,41,46,47,51,1001,129,153,4101 } |
Select-Object TimeCreated, ProviderName, Id, LevelDisplayName,
              @{n='Bugcheck';e={ if ($_.Id -eq 1001) { ($_.Message -split "`n" | Select-String 'BugcheckCode|bugcheck was') -join ' ' } else { '' } }} |
Sort-Object TimeCreated |
Format-List
```

Notes on the query itself:

- **`-ErrorAction SilentlyContinue` is removed.** The supplied version would hide an access failure and return an empty list, which reads as "no events" — the one outcome that must never be fabricated.
- **`1001` messages carry the bugcheck code.** The projection extracts it so `0x1E` can be confirmed from the log rather than from memory.
- **`Sort-Object TimeCreated`** makes onset visible without a second command.

**Output handling:** written to `eval/runs/GATE2-EVENT-LOG-2026-11-14.txt`; counts reported, messages truncated per event. The file is evidence and is not edited afterwards.

---

## 3. What is asked

| # | Ask |
|---|---|
| 1 | Authorise **gate 2** for this single query, this window, this ID set |
| 2 | Confirm the window (30 days) or substitute one |
| 3 | Confirm the output path and that a truncated message column is acceptable |

On authorisation the query runs once, the counts are reported, and any interpretation distinguishes **positive evidence** (a fault type that stopped at a known moment) from **absence of a report**, which is not evidence of anything.

---

## 4. What is explicitly not part of this request

- No memory diagnostic, no reboot-based test.
- No drive benchmark, surface scan, SMART attribute pull, or stress test.
- No full-file hash, no model read, no S-9a execution, no USN query.
- No change to any driver, page-file setting, power plan, or registry value.
- Stages 2–4 of the companion plan (memory diagnostic, workload observation, hardware triage) are **separate** and would each need their own authorisation; none is implied by gate 2.

---

## 5. Claim ceiling

**Establishes:** the exact query proposed, the two defects found in the supplied version, and the boundary between a stability question and the digest question.

**Does not establish:** that the query will return anything, that absence of events means health, or anything at all about which layer returned disagreeing bytes.
