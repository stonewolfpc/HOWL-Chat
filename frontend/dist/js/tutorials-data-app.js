/**
 * HOWL Chat — Tutorial Data: App Functions & UI
 *
 * @module tutorials-data-app
 */

'use strict';

window.TUTORIALS_DATA_APP = {
    id: 'app',
    title: 'App Functions',
    topics: [
        {
            id: 'menu-file',
            title: 'File Menu',
            content: `# File Menu

Access via the **File** button in the top menu bar.

## Save

**File → Save**

Saves the current chat session to the backend with its current name.

## Save As

**File → Save As**

Saves the current chat with a new name. The original session remains unchanged.

## Load

**File → Load**

Opens a list of saved chat sessions. Select one to restore:
- Message history
- Loaded world, scenario, characters
- Sampler and model settings

## Delete

**File → Delete**

Permanently removes a saved chat session. Requires confirmation.

## New Chat

**File → New Chat**

Creates a fresh chat session. Does not affect saved sessions.

## Settings

**File → Settings**

Opens the full settings modal (General, Storage, Network, Privacy, Audio).

## Advanced

**File → Advanced**

Opens the Sampler Settings flyout directly.

## Exit

**File → Exit**

Closes the application. All unsaved data may be lost — save first.

<info>
The File menu is your primary hub for session management. Get familiar with Save, Load, and New Chat — they are the most frequently used options.
</info>`
        },
        {
            id: 'menu-edit',
            title: 'Edit Menu',
            content: `# Edit Menu

Standard text editing operations.

## Available Actions

| Action | Description |
|--------|-------------|
| **Undo** | Reverses the last text edit in the input field. |
| **Redo** | Re-applies an undone edit. |
| **Cut** | Removes selected text and copies it to the clipboard. |
| **Copy** | Copies selected text to the clipboard. |
| **Paste** | Inserts clipboard text at the cursor position. |

## Usage

These actions work in any text field:
- Chat input bar
- World/character/scenario description fields
- Lorebook content editor
- Settings text inputs

## Keyboard Shortcuts

| Action | Windows/Linux | macOS |
|--------|---------------|-------|
| Undo | <key>Ctrl</key> + <key>Z</key> | <key>Cmd</key> + <key>Z</key> |
| Redo | <key>Ctrl</key> + <key>Shift</key> + <key>Z</key> | <key>Cmd</key> + <key>Shift</key> + <key>Z</key> |
| Cut | <key>Ctrl</key> + <key>X</key> | <key>Cmd</key> + <key>X</key> |
| Copy | <key>Ctrl</key> + <key>C</key> | <key>Cmd</key> + <key>C</key> |
| Paste | <key>Ctrl</key> + <key>V</key> | <key>Cmd</key> + <key>V</key> |

<info>
Undo/Redo is limited to the current text field. There is no global undo for deleted messages — use message actions (Copy, Edit) instead.
</info>`
        },
        {
            id: 'menu-view',
            title: 'View Menu',
            content: `# View Menu

Controls the visual presentation of the application.

## Accessibility Modes

| Mode | Effect |
|------|--------|
| **Colorblind Mode** | Adjusts colors for protanopia, deuteranopia, and tritanopia. |
| **High Contrast Mode** | Maximizes contrast between text and background. |

## Themes

| Theme | Description |
|-------|-------------|
| **Light** | Bright background, dark text. |
| **Dark** | Dark background, light text. Default. |
| **Custom** | Load a custom theme CSS file. |

## UI Scale

Adjust the interface size:
- **Increase UI Scale** — Makes everything larger.
- **Decrease UI Scale** — Makes everything smaller.
- **Reset UI Scale** — Returns to 1.0x.

Range: 0.5x to 2.0x.

## Panel Toggles

| Toggle | Shows/Hides |
|--------|-------------|
| **Toggle Sidebar** | Left sidebar (sessions, characters, context). |
| **Toggle Character Panel** | Character list in the sidebar. |
| **Toggle World Panel** | World/scenario display in the sidebar. |

## Fullscreen

**View → Fullscreen**

Expands the chat to fill the entire screen. Press <key>F11</key> or <key>Escape</key> to exit.

<info>
UI Scale and Theme settings persist across app restarts. They are saved in localStorage.
</info>`
        },
        {
            id: 'menu-settings',
            title: 'Settings Menu',
            content: `# Settings Menu

Quick access to all configuration panels.

## Sampler Settings

**Settings → Sampler Settings**

Opens the Sampler Settings flyout. Controls token selection behavior.

## Model Settings

**Settings → Model Settings**

Opens the Model Settings flyout. Controls model loading and performance.

## Story Tools

**Settings → Story Tools**

Opens the Story Tools flyout. Access memories, lorebooks, and knowledge documents.

## Advanced / Experimental

**Settings → Advanced / Experimental**

Placeholder for upcoming features.

## Keyboard Shortcuts

**Settings → Keyboard Shortcuts**

Shows a reference table of all keyboard shortcuts.

## Audio Settings

**Settings → Audio Settings**

Opens the Audio Settings modal. Configure TTS, ASR, VAD, and hardware devices.

## TTS Settings

**Settings → TTS Settings**

Quick access to TTS voice and speed controls.

<info>
All settings menus can also be accessed via the **File → Settings** modal, which provides a unified settings interface.
</info>`
        },
        {
            id: 'tts-voice-registry',
            title: 'TTS Voice Model Registry',
            content: `# TTS Voice Model Registry

HOWL's unified model registry also manages TTS voice models (Piper .onnx files). Scan once, remember forever — voices from any directory appear in the TTS dropdown automatically.

## How It Works

The registry stores voice models the same way it stores LLMs and image models:
- Scan a directory → all \`.onnx\` voice files are registered with their absolute paths
- Registry is persisted → saved to \`config/model-registry.json\`
- Directory is remembered → saved to \`config/model-scan-dirs.json\`
- Auto-rescan on startup → newly added voices appear automatically

## Where to Find It

**Settings → Audio Settings → TTS tab**

The TTS voice dropdown shows:
- Built-in voice (en_US-lessac-medium)
- All registered voices from your scanned directories

## How to Add Voices

**Scan a Directory:**
1. Open **Settings → Audio Settings**
2. Click the **TTS** tab
3. Click **+ Dir** next to the TTS Voice Model dropdown
4. Select a directory containing Piper .onnx voice models
5. All voices are scanned and registered
6. The directory is remembered for auto-rescan on startup
7. Select a voice from the refreshed dropdown

## Example Workflow

**New User (First Time Setup):**
1. Download Piper voices from the [Piper repository](https://github.com/rhasspy/piper)
2. Extract to \`F:\\Voices\\Piper\`
3. Open HOWL → Settings → Audio Settings → TTS
4. Click **+ Dir** → select \`F:\\Voices\\Piper\`
5. All voices appear in the dropdown
6. Done. Future HOWL sessions show all voices automatically.

**Power User (Multiple Voice Collections):**
1. Scan \`F:\\Voices\\Piper-English\` (50 English voices)
2. Scan \`F:\\Voices\\Piper-Multilingual\` (30 other languages)
3. Scan \`E:\\Work\\Voices\` (10 custom voices)
4. All 90 voices appear in the dropdown
5. Next week, download 5 new voices to \`F:\\Voices\\Piper-English\`
6. Restart HOWL → the 5 new voices appear automatically

## Registry vs Manual Selection

| Method | When to Use |
|--------|-------------|
| **Registry dropdown** | Your permanent voice collection. Scan once, use forever. |
| **Built-in only** | If you only need the default voice and don't want to manage a registry. |

<success>
The TTS registry works exactly like the LLM registry. Add your voice directories once, and all voices appear automatically on every startup. New voices added to those directories are picked up automatically.
</success>`
        },
        {
            id: 'image-model-registry',
            title: 'Image Model Registry',
            content: `# Image Model Registry

HOWL's unified model registry manages all image generation assets — diffusion models, VAEs, LoRAs, ControlNets, embeddings, and upscalers. Scan once, remember forever — assets from any directory appear in the Image Generator dropdowns automatically.

## How It Works

The registry stores image assets alongside LLMs and TTS voices:
- Scan a directory → all supported files (.safetensors, .ckpt, .gguf, .bundle, .onnx, .pt) are registered
- Registry is persisted → saved to \`config/model-registry.json\`
- Directory is remembered → saved to \`config/model-scan-dirs.json\`
- Auto-rescan on startup → newly added assets appear automatically
- Cross-directory access → assets from multiple directories appear together in all dropdowns

## Where to Find It

**Image Generator page** (access via main navigation or File → Image Generator)

The Image Generator has multiple dropdowns, all populated from the registry:
- **Model** — Main diffusion models (SDXL, Flux, SD1.5, etc.)
- **VAE** — Variational autoencoders
- **LoRA** — Style adapters
- **ControlNet** — Structure/pose control models
- **Embeddings** — Textual inversions
- **Upscalers** — Image upscaling models

## Supported Model Families

HOWL supports the following diffusion model families via stable-diffusion.cpp:

| Model Family | CLI Flag | Companion Models Required | Special Flags | Download Links |
|--------------|----------|---------------------------|---------------|----------------|
| **SD1.x / SD2.x** | \`-m\` | clip_l + vae | None | [CivitAI SD1.5](https://civitai.com/models/3/realistic-vision-v51) |
| **SDXL** | \`-m\` | clip_l + vae | None | [CivitAI SDXL](https://civitai.com/models/164278/juggernaut-xl-v9) |
| **SD3 / SD3.5** | \`-m\` | clip_l + clip_g + t5xxl + vae | None | [HuggingFace SD3](https://huggingface.co/stabilityai/stable-diffusion-3-medium) |
| **Flux.1 / Flux.2** | \`--diffusion-model\` | clip_l + t5xxl + vae | None | [HuggingFace Flux](https://huggingface.co/black-forest-labs/FLUX.1-dev) |
| **Chroma** | \`--diffusion-model\` | clip_l + t5xxl + vae | \`--chroma-disable-dit-mask\` | [HuggingFace Chroma](https://huggingface.co/filipstrand/chroma-v10) |
| **Anima** | \`--diffusion-model\` | \`--llm\` (Qwen3-0.6B-Base) + qwen_image_vae | None | [HuggingFace Anima](https://huggingface.co/circlestone-labs/Anima) |
| **HiDream-O1** | \`-m\` | None (self-contained) | None | [HuggingFace HiDream](https://huggingface.co/Comfy-Org/HiDream-O1-Image) |
| **Qwen Image** | \`-m\` | None (self-contained) | None | [HuggingFace Qwen Image](https://huggingface.co/Qwen/Qwen2-VL-72B-Instruct) |
| **AuraFlow** | ❌ NOT SUPPORTED | ❌ | ❌ | — |

## Companion Model Requirements

Some models require specific companion files. HOWL auto-detects these if they're in the same directory or scanned directories.

### For Flux / Chroma
- **clip_l**: CLIP-L text encoder (usually named \`clip_l.safetensors\` or \`CLIP-L.safetensors\`)
- **t5xxl**: T5-XXL text encoder (usually named \`t5xxl_fp16.safetensors\` or \`t5xxl_fp32.safetensors\`)
- **vae**: VAE model (usually named \`ae.safetensors\` or \`vae.safetensors\`)

### For Anima
- **llm**: Qwen3-0.6B-Base text encoder (NOT clip_l or t5xxl!)
  - Download: [Qwen3-0.6B-Base](https://huggingface.co/circlestone-labs/Anima/tree/main/split_files/text_encoders)
- **vae**: qwen_image_vae (NOT standard ae.safetensors!)
  - Download: [qwen_image_vae](https://huggingface.co/circlestone-labs/Anima/tree/main/split_files/vae)

### For SD3 / SD3.5
- **clip_l**: CLIP-L text encoder
- **clip_g**: CLIP-G text encoder
- **t5xxl**: T5-XXL text encoder
- **vae**: VAE model

### For SDXL
- **clip_l**: CLIP-L text encoder
- **vae**: VAE model

### For SD1.x / SD2.x
- **clip_l**: CLIP-L text encoder
- **vae**: VAE model

<warning>
Anima requires Qwen3-0.6B-Base as its text encoder, NOT clip_l or t5xxl. Using the wrong companion model will cause "unknown tensor" errors.
</warning>

## How to Add Assets

**Scan a Directory:**
1. Open the **Image Generator** page
2. Click **+ Dir** next to the Model dropdown
3. Select a directory containing diffusion models
4. All supported files are scanned and registered
5. The directory is remembered for auto-rescan on startup
6. Repeat for other asset types if stored in separate directories

## Example Workflow

**New User (First Time Setup):**
1. Download SDXL models from CivitAI or HuggingFace
2. Extract to \`D:\\AI GGUF\\Image\\SDXL\`
3. Download VAEs to \`D:\\AI GGUF\\Image\\VAEs\`
4. Download LoRAs to \`D:\\AI GGUF\\Image\\LoRAs\`
5. Open Image Generator → click **+ Dir** for Model → select SDXL directory
6. Click **+ Dir** for VAE → select VAEs directory
7. Click **+ Dir** for LoRA → select LoRAs directory
8. All assets appear in their respective dropdowns
9. Done. Future sessions show all assets automatically.

**Power User (Multiple Collections):**
1. Scan \`D:\\AI GGUF\\Image\\SDXL\` (20 SDXL models)
2. Scan \`E:\\Work\\Models\\Flux\` (15 Flux models)
3. Scan \`F:\\Experimental\\New\` (10 experimental models)
4. All 45 models appear in the Model dropdown together
5. Next week, download 5 new Flux models to \`E:\\Work\\Models\\Flux\`
6. Restart HOWL → the 5 new models appear automatically in the dropdown

## Cross-Directory Generation

The registry enables a powerful workflow: **generate with models from any directory without copying files**.

**Example:**
- Your configured \`ModelsDir\` is \`D:\\AI GGUF\\Image\\SDXL\`
- You also scanned \`E:\\Work\\Models\\Flux\`
- In the Model dropdown, you see models from BOTH directories
- Select a Flux model from \`E:\\Work\\Models\\Flux\`
- Click Generate → HOWL loads it from its original path
- No file copying required — the engine loads directly from wherever the file lives

## Registry vs Manual Selection

| Method | When to Use |
|--------|-------------|
| **Registry dropdown** | Your permanent asset collection. Scan once, use forever. |
| **Browse (manual)** | Testing a one-off model, or a file you don't want remembered. |
| **+ Dir** | Adding a new asset directory to your permanent collection. |

## Supported File Types by Asset Type

| Asset Type | Extensions |
|------------|------------|
| **Diffusion Models** | .safetensors, .ckpt, .gguf, .bundle |
| **VAEs** | .safetensors, .ckpt |
| **LoRAs** | .safetensors |
| **ControlNets** | .safetensors |
| **Embeddings** | .pt, .safetensors |
| **Upscalers** | .safetensors |
| **TTS Voices** | .onnx |
| **LLMs** | .gguf, .bin |

<success>
The image registry is the most powerful part of the unified system. You can have models scattered across multiple drives and directories, and they all appear together in the dropdowns. No need to consolidate your collection — just scan each directory once.
</success>`
        },
        {
            id: 'sidebar-left',
            title: 'Left Sidebar',
            content: `# Left Sidebar

The left sidebar provides access to chat sessions, characters, and context information.

## Session List

Shows all recent chat sessions. Click a session to load it.

- **New Session** button — Creates a new chat.
- **Session names** — Based on first message or manual save name.
- **Timestamp** — When the session was last active.

## Character Panel

Lists all characters loaded into the current chat.

- **Add Character** — Load a character from your library.
- **Character avatar** — Visual identifier.
- **Character name** — Click to expand details.

### Character Context Menu

Right-click a character to:
- **Replace Image** — Change the avatar.
- **View Memories** — See character-specific memories.
- **Remove from Chat** — Unload the character.

## Context Panel

Shows the currently loaded world, scenario, and lorebooks.

- **World** — Active world name and thumbnail.
- **Scenario** — Active scenario name.
- **Lorebooks** — Count of active lorebook entries.

<info>
The left sidebar can be toggled via **View → Toggle Sidebar** or by dragging the edge.
</info>`
        },
        {
            id: 'sidebar-right',
            title: 'Right Sidebar',
            content: `# Right Sidebar

The right sidebar displays supplementary information about the current chat.

## Info Panel

- **Current Chat Name** — Active session name.
- **Current World Name** — Loaded world.
- **Loaded Model** — Currently running LLM.
- **Context Usage** — How much of the context window is filled.

## World Panel

A larger view of the active world:
- World thumbnail
- World description
- Attached scenarios and characters

## Panel Toggle

Click the panel header to collapse or expand sections.

<info>
The right sidebar is optional. Hide it for more chat screen space via **View → Toggle World Panel**.
</info>`
        },
        {
            id: 'keyboard-shortcuts',
            title: 'Keyboard Shortcuts',
            content: `# Keyboard Shortcuts

Master these shortcuts to navigate HOWL efficiently.

## Chat

| Shortcut | Action |
|----------|--------|
| <key>Enter</key> | Send message |
| <key>Shift</key> + <key>Enter</key> | New line in input |
| <key>Escape</key> | Abort generation |
| <key>Up</key> | Recall previous message (input field) |

## Navigation

| Shortcut | Action |
|----------|--------|
| <key>Ctrl</key> + <key>N</key> | New chat |
| <key>Ctrl</key> + <key>S</key> | Save chat |
| <key>Ctrl</key> + <key>O</key> | Load chat |
| <key>F11</key> | Toggle fullscreen |

## Settings

| Shortcut | Action |
|----------|--------|
| <key>Ctrl</key> + <key>,</key> | Open settings |
| <key>Ctrl</key> + <key>Shift</key> + <key>S</key> | Open sampler settings |

## Text Editing

| Shortcut | Action |
|----------|--------|
| <key>Ctrl</key> + <key>Z</key> | Undo |
| <key>Ctrl</key> + <key>Shift</key> + <key>Z</key> | Redo |
| <key>Ctrl</key> + <key>X</key> | Cut |
| <key>Ctrl</key> + <key>C</key> | Copy |
| <key>Ctrl</key> + <key>V</key> | Paste |

## Panels

| Shortcut | Action |
|----------|--------|
| <key>Ctrl</key> + <key>B</key> | Toggle sidebar |
| <key>Ctrl</key> + <key>Shift</key> + <key>B</key> | Toggle world panel |

<info>
Some shortcuts may conflict with browser defaults when running in web mode. Use the in-app menu as a fallback.
</info>`
        },
        {
            id: 'themes-accessibility',
            title: 'Themes & Accessibility',
            content: `# Themes & Accessibility

HOWL supports multiple visual modes to suit your preferences and needs.

## Built-in Themes

| Theme | Description |
|-------|-------------|
| **Dark** | Default. Low eye strain in dim environments. |
| **Light** | High contrast for bright rooms. |
| **Midnight** | Deep blue tones. |
| **Forest** | Green-tinted dark mode. |
| **Solarized** | Classic programmer theme. |

## Switching Themes

1. Click **View** in the top menu.
2. Select **Light Theme** or **Dark Theme**.
3. Or use **File → Settings** for more options.

## UI Scale

Adjust the entire interface size:

1. **View → Increase UI Scale** (or **Decrease** / **Reset**).
2. Or use **Ctrl + Plus** / **Ctrl + Minus**.

## Accessibility Features

| Feature | Purpose |
|---------|---------|
| **Colorblind Mode** | Adjusts colors for color vision deficiencies. |
| **High Contrast Mode** | Maximizes text-background contrast. |
| **Reduced Motion** | Disables animations for vestibular disorders. |
| **Screen Reader Support** | ARIA labels on all interactive elements. |

<info>
Theme and UI Scale settings are saved to localStorage and persist across sessions. They are applied immediately without restart.
</info>`
        }
    ]
};
