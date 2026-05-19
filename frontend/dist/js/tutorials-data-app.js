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
