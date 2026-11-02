# E: taken offline — 2026-11-01

## Why

Two unclean reboots the same day, both ending in `KERNEL_DATA_INPAGE_ERROR (0x7A)`:

| Time | Event | Device |
|---|---|---|
| 14:12:02 | 198 × `disk` Event 51, "error detected during a paging operation" | `\Device\Harddisk5\DR5` |
| 14:13:24 | Kernel-Power 41 (unclean shutdown) | — |
| 15:38:49 | Kernel-Power 41 (no storage errors logged) | — |

Disk 5 is the 2 TB NVMe that held E: (FriendlyName "PCIe SSD", BusType NVMe,
serial `6479_A73C_4030_01FF`, volume label "Inland"). Windows reported it
`Healthy`, which is not a useful signal: the 0x7A bugcheck plus a burst of paging
errors is the signature of a failing NVMe controller.

No minidump was written for the 15:38 crash (newest dump on disk is 10/29), so
there is no kernel dump to attribute that one definitively. The verifier's
`--hash` pass read E: roughly a minute before it, which is a plausible trigger
and is why `build/verify-move.py` was rewritten to never open E:.

## What was verified before offlining

- 46/46 manifest rows present on G: (29 byte-exact files + 17 app-written
  sidecars, valid JSON, zero `E:` references inside).
- Both orphan files at exact size in `G:\AI GGUF\Image models`.
- 0 `E:` references in `config/*.json`; registry 133 entries, 0 E: paths.
- Model data on E: is duplicated on G: (`G:\AI GGUF` = 179 GB GGUF + 57 GB image
  models + 39 GB modpack).
- User confirmed everything else on the drive is backed up or not needed.

## Action

```
powershell -NoProfile -Command "Set-Disk -Number 5 -IsOffline 1"
```

`-IsOffline 1`, not `$true`: through this shell `$true` arrives as the string
"True" and the cmdlet rejects it (that is exactly how the first attempt failed
and changed nothing).

## Reversing it

```
powershell -NoProfile -Command "Set-Disk -Number 5 -IsOffline 0"
```

Nothing was deleted or formatted. The volume and its data are intact; it is
simply not mounted. Physically removing the drive is the durable fix once the
screwdriver turns up.
