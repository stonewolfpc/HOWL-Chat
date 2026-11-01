# Image & Video Generation — Audit

Same method as the sampler pass: enumerate every control, prove what is wired to what, then fix the ones that
were decorative. 137 controls scanned in `image-generator.html`; each id cross-checked against every
`getElementById` in the JS bundle.

## Fixed

| # | Problem | Was | Now |
|---|---------|-----|-----|
| 1 | **Negative prompt never applied on reuse** | `historyMenuAction('use-prompt')` wrote to `ig-negative-prompt`, an id that has never existed in the page (the field is `ig-negative`) | writes to `ig-negative`, so "Use Prompt" restores the negative too |
| 2 | **Phantom tiling control** | generation read `ig-tiling`, which does not exist; the actual control is `ig-vae-tiling` (already sent as `vae_tiling`) | removed the phantom read and the stale restore line; `vae_tiling` remains the single path |
| 3 | **LoRA restore to a dead multi-select** | restore code targeted `ig-lora`; the page has `ig-lora-selector` plus a `window.activeLoRAs` array | restores into `activeLoRAs` (accepting the new `{name,strength,enabled}` shape) and re-renders the chips |
| 4 | **Upscale seed lock button did nothing** | `ig-upscale-lock` had no handler and no backing state; the upscale seed field was never used for generation | real `toggleUpscaleSeedLock()`: locked state shows on the button, and the upscale seed becomes the generation seed while locked |
| 5 | **Board search did nothing** | `ig-board-search` was in the page and no JS ever read it | filters history cards by filename or prompt/metadata text as you type |
| 6 | **Generated images were cropped** | `.ig-batch-thumb` and `.ig-history-item img` used `object-fit: cover` in square tiles, cutting off portraits and landscapes | `object-fit: contain` on a dark tile, so the whole frame is visible at any aspect ratio |

## Verified after the fix

Rendered `image-generator.html` at 1440×900 with four generated images of different aspect ratios
(1:1, 2:3, 3:2, 3:8) injected through the real `.ig-batch-grid` / `.ig-batch-thumb` markup, plus three history
cards in the real `.ig-history-item` markup:

- every thumbnail renders as a 1:1 tile with `object-fit: contain` and the full image visible;
- no horizontal overflow in the grid or the preview area;
- history cards keep square tiles without cropping;
- typing `knight` into `#ig-board-search` leaves exactly the knight card visible and hides the other two —
  the filter is live.

## Control inventory results

- Dead controls found and now all handled: `ig-upscale-lock` (wired), `ig-board-search` (wired).
- JS ids that had no element: `ig-negative-prompt` (fixed), `ig-tiling` (removed), `ig-lora` (fixed).
  `ig-init-image` and `ig-board-context-menu` are created dynamically at runtime — not defects.
- Remaining `ig-*` ids in the page are all referenced by JS.

## Not covered

- A real diffusion run (no image model was loaded during the audit) — card sizing and control wiring were verified
  with representative images injected into the live DOM, not with sampler output.
- The Wails bridge in the packaged window, as before.
