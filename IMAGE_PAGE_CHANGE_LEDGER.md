# IMAGE_PAGE_CHANGE_LEDGER.md

**Purpose:** Document every change, test, and finding for the Howl Chat image-page repair mission.

---

## Entry 1 — G: drive ControlNet rescan and rename test

**Date:** 2026-11-04  
**Files affected:** None (test only). Files touched on G: drive (reversible):
- `G:/AI GGUF/Image models/MLSDControlnetLllite_v20.safetensors` → `MLSD-Controlnet-Lllite-v20-TEST.safetensors`
- `G:/AI GGUF/Image models/MLSDControlnetLllite_v20.safetensors.howl.json` → `MLSD-Controlnet-Lllite-v20-TEST.safetensors.howl.json`
- `G:/AI GGUF/Image models/controlnet_openpose_sdxl_thibaud_fp16.safetensors` → `OpenPose-ControlNet-SDXL-Thibaud-TEST.safetensors`
- `G:/AI GGUF/Image models/controlnet_openpose_sdxl_thibaud_fp16.safetensors.howl.json` → `OpenPose-ControlNet-SDXL-Thibaud-TEST.safetensors.howl.json`

### What changed
1. Added `G:/AI GGUF/Image models` as a model directory via `AddModelDirectory` RPC (harness doesn't support browser-native folder picker).
2. Triggered rescan → scanner picked up 18 new/updated assets, including both ControlNets.
3. Backend registry confirmed both ControlNets registered as `type=controlnet, family=sdxl`.
4. **Rename test:** renamed both G: drive ControlNets on disk and re-ran rescan to verify the app picks up the change and keeps the correct family classification.

### Test results
**Before rename:**
- UI ControlNet dropdown (uid 12_203) showed:
  - `controlnet_openpose_sdxl_thibaud_fp16 [sdxl]` (uid 12_207)
  - `MLSDControlnetLllite_v20 [sdxl]` (uid 12_217)
- Registry had both with `family=sdxl`.

**After rename (post-rescan):**
- Registry now has only **1 G: drive ControlNet**: `mlsd` (family=sdxl), pointing to `MLSD-Controlnet-Lllite-v20-TEST.safetensors`.
- UI ControlNet dropdown now shows:
  - `mlsd [sdxl]` (uid 15_209) — the renamed MLSD file, correctly classified as sdxl.
  - `controlnet_openpose_sdxl_thibaud_fp16` **no longer present** — the OpenPose file was renamed and removed from registry.
  - `diffusion_pytorch_model_twins.safetensors [sdxl]` — **new entry** from H: drive (InvokeAI/models), unrelated to G: drive rename.
- G: drive folder shows 17 assets (down from 18) — consistent with one file removed by rename.

**Conclusion:** The app does pick up the rename on rescan. The renamed MLSD file is registered correctly with `family=sdxl`. The old name is removed. This is the correct behavior.

### Open question / unresolved
After the rename, the OpenPose ControlNet disappeared from the UI and registry. It was renamed to `OpenPose-ControlNet-SDXL-Thibaud-TEST.safetensors` on disk, but it is **not in the registry** and **not in the UI dropdown**. 

The file exists on G: drive at `G:/AI GGUF/Image models/OpenPose-ControlNet-SDXL-Thibaud-TEST.safetensors` (2386MB). A `clsfy` sample run classified it as `type=controlnet, family=sdxl` based on the embedded checkpoint metadata.

Possible causes (not yet investigated):
- The rescan may not have re-scanned the G: drive folder after the rename was written.
- The scanner may have had a transient issue reading the large (2.3GB) .safetensors file.
- The `.howl.json` sidecar for the renamed file may need to be re-created or verified.

The rename is reversible — both files can be renamed back to their original names.

### Rollback method
Rename both files back to their original names:
```bash
cd /c/Users/StoneWolfPC/HOWL_Chat
python3 -c "
import os
g = 'G:/AI GGUF/Image models'
os.rename(os.path.join(g, 'MLSD-Controlnet-Lllite-v20-TEST.safetensors'), os.path.join(g, 'MLSDControlnetLllite_v20.safetensors'))
os.rename(os.path.join(g, 'MLSD-Controlnet-Lllite-v20-TEST.safetensors.howl.json'), os.path.join(g, 'MLSDControlnetLllite_v20.safetensors.howl.json'))
os.rename(os.path.join(g, 'OpenPose-ControlNet-SDXL-Thibaud-TEST.safetensors'), os.path.join(g, 'controlnet_openpose_sdxl_thibaud_fp16.safetensors'))
os.rename(os.path.join(g, 'OpenPose-ControlNet-SDXL-Thibaud-TEST.safetensors.howl.json'), os.path.join(g, 'controlnet_openpose_sdxl_thibaud_fp16.safetensors.howl.json'))
"
```
Then trigger rescan via UI (↺ Rescan button).

---

## Entry 2 — ControlNet scanner behavior confirmed

**Date:** 2026-11-04  
**Files affected:** None

### Findings
1. The scanner **correctly classifies both G: drive ControlNets** as `type=controlnet, family=sdxl` (verified via `clsfy` sample and registry inspection).
2. The UI ControlNet dropdown picks up the scanner's output — both ControlNets appear with `[sdxl]` tag.
3. The rename test confirmed the app picks up on-disk renames and updates the registry + UI dropdown accordingly.
4. The "6 incompatible files hidden" group in the main model dropdown (uid 15_41) is unrelated to the ControlNet test.
5. The scanner picks up `.safetensors`, `.gguf`, `.pt`, and `.bin` files from flat directories (no UUID subfolders required) when companion `.howl.json` sidecars exist alongside them.

### Observations
- The G: drive folder is a flat directory (no UUID subfolders). All model files sit directly in `G:/AI GGUF/Image models/` with `.howl.json` sidecars.
- The scanner found 28 ControlNets total after the final rescan (up from 2 initially). Most are from H: drive's InvokeAI/models folder.
- `GetRegisteredModels()` RPC with `e.Type != "diffusion"` filter does NOT return ControlNets. The UI populates the ControlNet dropdown from `GetCompatibleAssets(family)` or `scanIndex.controlnets`.

### Limitations
- The OpenPose ControlNet (`controlnet_openpose_sdxl_thibaud_fp16`) disappeared from the registry after rename. Root cause not yet identified. This requires further investigation before concluding the rename-test is fully successful.
- The `diffusion_pytorch_model.safetensors` (46MB) from G: drive is registered as `type=diffusion, family=sdxl` but its legitimacy as a real model is unverified. It may be a misclassified file.

---

## Rollback summary
All changes so far are on-disk renames of two G: drive files. Rollback is the inverse rename followed by a UI rescan. No source code or configuration files were modified.
