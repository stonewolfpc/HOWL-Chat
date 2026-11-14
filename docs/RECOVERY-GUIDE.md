# HOWL Chat — Spare-Laptop Recovery Guide

**Backup commit:** `7a3a175` (plus this guide) · **Branch:** `master`
**Repository:** https://github.com/stonewolfpc/HOWL-Chat

> ## Read this first
> The repository's **default branch is `main`, which contains only an empty
> initial commit**. A plain `git clone` lands there and gives you one file.
> Always clone with `-b master`, or switch branches after cloning.

---

## 1. Clone

```bash
git clone -b master https://github.com/stonewolfpc/HOWL-Chat.git
cd HOWL-Chat
git checkout master
git log --oneline -3     # expect 7a3a175 or later
```

If you already cloned and see a single file, run `git fetch origin master:master && git checkout master`.

The clone is roughly 325 MB because earlier commits tracked 14 built `.exe` files. That is history, not the current tree; it is not repairable without rewriting history, which this backup deliberately did not do.

## 2. Prerequisites

| Need | Detail |
|---|---|
| OS | Windows 10/11 x64 |
| WebView2 runtime | built into Windows 11; on Windows 10 install the Evergreen runtime |
| **Go** | **1.23 or newer** (`go.mod` says `1.23.8`; the backup was built with 1.27.3) |
| Node / npm | **not needed** — `frontend/dist` is committed pre-built |
| Wails CLI | **not needed** to build — no cgo, no MSVC, no clang |
| MSVC / Build Tools | **not needed** (`CGO_ENABLED=0`) |

`llama.cpp` runs as an **external `llama-server.exe`**, not linked in. That is why there is no cgo toolchain requirement.

## 3. Build and run

```bat
build_howl_cgo.bat
```

or directly:

```bat
set CGO_ENABLED=0
set GOFLAGS=-buildvcs=false
go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o build\bin\howl-chat.exe .
```

Then `start-howl-chat.bat` to launch, or `start-howl-chat.bat dev` for live-reload mode.

`go build` fetches modules from `go.sum`; no vendored directory is present, so the machine needs internet on first build.

## 4. Restore separately — not in the repository

| Item | Why excluded | How to restore |
|---|---|---|
| Model weights (`.gguf`, `.safetensors`, …) | licence + size | download again; place them wherever you like, then point the app at the folder |
| `runtime/llama/llama-server.exe` and its DLLs | build artifact | download a llama.cpp Windows release |
| Root DLLs (`ggml-*.dll`, `cublas*_12.dll`, `cudart64_12.dll`, `stable-diffusion.dll`, …) | build artifact | same release, plus the stable-diffusion.cpp release for image generation |
| `sd-cli.exe`, `sd-server.exe` | build artifact | stable-diffusion.cpp release (image generation only) |
| `build/bin/` | build output | rebuild with step 3 |
| `backups/*.tar.gz` | heavyweight snapshot | whatever separate copy you keep of it |
| `models/`, `images/gallery/`, `runtime/`, `.gocache/` | local data / caches | created as needed |

Nothing here is irreplaceable **except** any local snapshot archives you keep outside Git — those must be copied by hand.

## 5. Configuration to re-check after cloning

Committed, but machine-specific — paths will point at this PC's drives:

- `config/model-scan-dirs.json` — folders the app scans for models
- `config/model-registry.json` — known models and their paths
- `config/settings.json` — general settings
- `config/image-settings.json` — image-page settings

Model paths referencing `G:\`, `H:\` or `E:\` will not exist on the spare laptop. Re-point them, or let the app re-scan.

## 6. Credentials

No API keys, tokens, certificates or `.env` files are committed. If the app needs a key for an online service, enter it again by hand on the new machine.

## 7. Verify the checkout works

```bat
go test ./... -count=1
```

Frontend behaviour tests live in `tests/frontend/*.test.js` and are plain JS, not wired to a package manager.

## 8. What was verified before this guide was written

- Remote `master` head matches local `HEAD` (`7a3a175`).
- A fresh `git clone -b master` into a temp directory produced 1,879 tracked files with `go.mod`, `wails.json`, `main.go`, `frontend/dist` and `config/` present, and **zero** `.exe`/`.dll` in the working tree.
- `go build -tags "desktop,production"` from that tree succeeded, producing a 23.5 MB `howl-chat.exe`.

Not verified: launching the app, loading a model, or generating a response — those need the external assets in §4.
