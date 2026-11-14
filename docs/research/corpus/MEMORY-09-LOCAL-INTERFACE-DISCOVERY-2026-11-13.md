# `MEMORY-09` — Local Interface Discovery (pre-L2)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-LID-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **DISCOVERY — filesystem only. No binary executed, no server contacted, no prompt sent.**
**Authority:** owner statement, 2026-11-13: *"kobold is on a different drive … in a directory called ai apps … so is silly tavern."*
**Supersedes:** the "no binary found" finding in row 487 §1, whose search scope was PATH and the project tree only.

---

## 1. What was found

Location: `H:\AI apps\` (owner-stated directory, confirmed on disk).

| Item | Evidence | Version |
|---|---|---|
| `koboldcpp.exe` | present, **607 MB** single binary | **unknown — see §3** |
| `llama.cpp-b10076-cuda12\` | build tree + `llama-b10076-bin-win-cuda-12.4-x64.zip` | **b10076, CUDA 12.4** |
| `llama.cpp-prism-b9596-cuda12\` | build tree | **b9596, CUDA 12.4** |
| `SillyTavern\` | directory with `package.json` | **1.18.0** (read directly) |
| `*.kcpps` presets | 5 files (`gemma4`, `glm4`, `mag mel r1`, `qwen 9b`, `riconante v2`) | prior saved configurations |
| `KoboldCpp_Ternary-Bonsai-27B_Settings.txt` | setup notes | — |
| Model files | `G:\AI GGUF\GGUF\` via `config/model-registry.json` | GGUF, per-model sizes recorded |

**Method:** directory listing, PE metadata query, `package.json` read, bounded string search. **Nothing was executed.**

## 2. What this settles

- The earlier "no binary found" was **a scope limit, not an absence**. PATH and project-tree searches were correct; the binaries were elsewhere. Row 487 is corrected by a revision note, not rewritten.
- **KoboldCpp-first is now a supported selection**, not a provisional one: the binary exists, and it is the only candidate that presents both a native generate surface and an OpenAI-compatible one from a single process.
- `llama.cpp` direct remains viable and *better versioned* than KoboldCpp — two build tags are visible in directory names.

## 3. The gap that files cannot close

**KoboldCpp's version is not recoverable from the filesystem.**

| Attempted | Result |
|---|---|
| PE `VersionInfo` (`FileVersion`, `ProductVersion`) | **empty** — typical of a frozen Python build |
| Plain version strings in the binary | **none found** by a bounded search |

Absence of a plain string is **not** proof of absence: a frozen bundle compresses its payload, so the string is likely present but not scannable. This is recorded as *unavailable*, not as *absent*.

**Two ways to resolve it, both outside today's authority:**
1. `koboldcpp.exe --version` — executes the binary (no model, no generation). Still a process start, so it needs authorisation.
2. The server's own identity endpoints once L2 is authorised: `/api/extra/version`, `/api/v1/model`, `/api/extra/true_max_context_length`.

**Route 2 is better:** it reports the version *and* confirms the surface in one authorised step, and it is the same endpoint the harness will rely on later.

## 4. Context the binaries give, without a model loaded

- **CUDA 12.4 builds for both `llama.cpp` trees**, matching a Pascal-generation card (GTX 1080 8 GB, owner-reported, unverified). CUDA support for Pascal is a version-compatibility question that L2/L3 will answer, not this note.
- **A 27B-class preset** (`Ternary-Bonsai-27B`) on an 8 GB card implies **CPU/GPU split offload**, not full residency. That is a *performance-stream* fact, and it does not enter behavioural evidence — but it does mean generation will be slow, which is expected and acceptable.
- **Five saved `.kcpps` presets** exist. They are existing configurations, not harness conditions; the harness will define its own frozen params and record them (row 486 §9). Presets are useful only as evidence that these model/sampler combinations have run before.

## 5. What L2 would do, precisely

Authorisation required. On approval, a single authorised discovery step would:

1. `GET /api/extra/version` → server name and version.
2. `GET /api/v1/model` → loaded model identity, if any.
3. `GET /api/extra/true_max_context_length` → actual `n_ctx` (not the requested value).
4. Record all three in an identity block, with any missing field marked `unavailable`.

**No prompt, no generation, no model load triggered by the harness.** If the server is not already running, the probe reports "not listening" and stops — the harness does not start servers.

## 6. Not authorised, and not done

Starting `koboldcpp.exe`; running `--version`; any HTTP request to `127.0.0.1`; loading a model; generating; reading the `.kcpps` presets in detail; measuring hardware; any hosted route.

## 7. L2 executed — result and what it means

**Authorised:** owner, 2026-11-13 (*"you may run it"*). **Executed:** 2026-11-14T07:42:18Z.

**Instrument:** `eval/harness/discover.py`, five new checks (75–79). Fixed in code: loopback only, **GET on three identity paths and nothing else**, no retries, no server start, missing fields marked `unavailable`.

**Pre-check before the probe:** no `kobold` or `llama` process running; no answer on ports 5001, 5000, 8000 or 8080.

**Probe result:**

| Field | Value |
|---|---|
| listening | **false** |
| server version | `unavailable` |
| loaded model | `unavailable` |
| max context length | `unavailable` |
| identity GETs | three, all `http=0` (no response) |

**Recorded at:** `eval/runs/L2-2026-11-14T074218Z/record.json` (+ `record.sha256` sidecar), write-once.

**Interpretation.** This is an **availability** result, not a failure and not evidence about anything. The server is not running; the harness does not start servers, so it stopped and recorded that. Consequences:

- KoboldCpp's **version remains unknown** — the one fact L2 existed to obtain.
- The documented reply shapes §5 of row 487 remain **unconfirmed against this build**.
- No prompt was sent, no model was loaded, no server was started.

**To complete L2** requires a running server, which means a **model load** — a separate authorisation, and a model choice (`gemma4`, `glm4`, `mag mel r1`, `qwen 9b`, `riconante v2` presets exist, or a GGUF from `G:\AI GGUF`). Until then L2 is **attempted, not established**.

### 7.1 The conditions, as machine checks

The owner restated L2's boundaries after the run. Each is now a check, not a convention:

| Condition | Enforcement | Check |
|---|---|---|
| Probe only `127.0.0.1` | `LocalEndpoint` loopback guard; non-loopback raises | 75 |
| Bounded timeouts | `TIMEOUT_SECONDS = 3.0`, single attempt, no retry | 76 |
| No POST, no prompt, no model load, no generation | Method pinned to `GET`; path allowlist of three; source may contain no `POST`, `generate(`, `dispatch` | 80 |
| Record status, bodies, version/model/context | Raw bodies preserved verbatim in the block | 79 |
| Mark unavailable, never infer | Empty responses yield `unavailable`; version comes from the endpoint only, never from the executable | 77 |
| Silent endpoint → stop and report | Recorded as `listening: false` with a conclusion; nothing started | 78 |
| Discovery ≠ L3 authorisation | `discover.py` references no grant, no dispatcher, no `P2-GRANT`; discovery works with the grant **absent** | 81 |

**One honest discrepancy:** the instruction mentions *four* authorised GETs; the probe issues **three** — `/api/extra/version`, `/api/v1/model`, `/api/extra/true_max_context_length` — because only three identity facts are needed and each extra endpoint widens the surface. If a fourth identity path is intended, naming it is a one-line addition; guessing one would defeat the allowlist.

Suite: **81/81, exit 0.** A second probe after the checks were added again returned `listening: false`, and `eval/runs/P2-GRANT.json` remains **absent** — discovery created no authorisation.

## 7.2 L2 established — model-free startup, real identity

**Authority:** `WO-L2-START`, 2026-11-13. **Mode:** model-free startup + read-only discovery.

### Startup

| Item | Value |
|---|---|
| Executable | `H:\AI apps\koboldcpp.exe` |
| Arguments | `--nomodel --skiplauncher --host 127.0.0.1 --port 5001 --quiet` |
| Binding | `http://127.0.0.1:5001` — loopback only |
| Model loaded | **no** — `/api/extra/version` reports `"llm": false`; `/api/v1/model` reports `"inactive"` |
| Preset loaded | no |
| GUI | suppressed by `--skiplauncher`; none opened |

**Model-free server startup is supported.** `--nomodel` is documented as *"launch the GUI alone, without selecting any model"*; empirically, with `--skiplauncher`, the **API serves with no model loaded** — the ambiguity was real and only a test could settle it.

### Identity, as returned

| Endpoint | HTTP | Body |
|---|---|---|
| `/api/extra/version` | 200 | `{"result":"KoboldCpp","version":"1.120","protected":false,"llm":false,…}` |
| `/api/v1/model` | 200 | `{"result":"inactive"}` |
| `/api/extra/true_max_context_length` | 200 | `{"value":12288}` |

**Server: KoboldCpp 1.120.** The version that was unrecoverable from the filesystem came back from the server itself, which is why the endpoint route was preferred over executing the binary.

**Caveat on the context figure:** 12288 is the server's configured default with **no model resident**, not a model's trained context. It becomes a model fact only when a model is loaded.

**Recorded:** `eval/runs/L2-2026-11-14T075057Z/record.json`, block hash `09d9a60c…` (superseding the earlier silent-endpoint attempt at `…T074218Z`, kept as history).

### A real body corrected the parser — the rule working as designed

First run parsed `server_version` as **"KoboldCpp"**: the parser took `result`, assuming the fixture shape where `result` *was* the version. The real body puts the name in `result` and the version in a separate `version` field.

**Fix:** prefer `version`; treat a non-numeric `result` as the **server name** and record it separately; still accept a bare numeric `result` for older/other shapes. Regression check **82**. Suite: **82/82, exit 0.**

This is exactly the case the fixtures existed to survive: a documented-looking assumption met a real body and lost, and the parser changed rather than the interpretation.

### Preflight detail worth recording

The first launch attempt blocked the calling shell for the full 120 s because the child process inherited stdout, even though the server came up normally. It was **not** killed; a follow-up check found port 5001 answering `200`. A launcher should redirect the child's streams and return, rather than wait on an inherited pipe.

### Server left running

Per `WO-L2-START` §4, the server was **not stopped**. Two `koboldcpp.exe` processes are alive (bootloader + child) on 127.0.0.1:5001. Stopping it requires separate authorisation.

### Boundary

No prompt, no generation, no model load, no CUDA compatibility test, no behavioural evaluation. **L2 success authorises nothing further** — L3 needs its own decision, and with no model resident the next step would necessarily involve a model choice.

## 8. Claim ceiling

A filesystem inventory with three version facts and one explicit gap. It establishes **what is installed and where**, not what any of it returns. The distinction that matters: the adapter is still written against documented shapes (row 487 §5), and only a real body at L3 can confirm them.
