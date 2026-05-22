# HOWL Chat - Comprehensive CSS Style Audit

**Generated:** January 2026  
**Scope:** App-wide CSS elements, borders, text boxes, styles, colors, patterns, shapes, and fonts

---

## Table of Contents
1. [CSS Custom Properties (Variables)](#css-custom-properties)
2. [Dynamic Theme System (JavaScript)](#dynamic-theme-system)
3. [Color Palette Summary](#color-palette-summary)
4. [Background Colors](#background-colors)
5. [Border Colors](#border-colors)
6. [Text/Font Colors](#textfont-colors)
7. [Shadows & Glows](#shadows--glows)
8. [Gradients](#gradients)
9. [Font Properties](#font-properties)
10. [Inline Styles (HTML)](#inline-styles-html)
11. [JavaScript-Generated Styles](#javascript-generated-styles)
12. [Component-Specific Styles](#component-specific-styles)

---

## CSS Custom Properties (Variables)

### Primary Theme Variables
- `--theme-color` - Main accent color (dynamically set by theme.js)
- `--theme-glow` - Glow effect color (dynamically set by theme.js)
- `--theme-text-glow` - Text glow color (dynamically set by theme.js)
- `--theme-secondary` - Secondary accent color (dynamically set by theme.js)
- `--static-text` - Static text color (#e0e0e0 default)
- `--text-secondary` - Secondary text color (#a0a0b0 default)
- `--text-primary` - Primary text color
- `--text-tertiary` - Tertiary text color

### Story Tools Variables
- `--surface-primary` - Primary surface background
- `--surface-secondary` - Secondary surface background
- `--surface-tertiary` - Tertiary surface background
- `--surface-hover` - Hover state surface
- `--border-primary` - Primary border color
- `--border-secondary` - Secondary border color
- `--accent-primary` - Primary accent color
- `--accent-primary-hover` - Hover state primary accent
- `--accent-secondary` - Secondary accent color
- `--accent-secondary-hover` - Hover state secondary accent
- `--danger-primary` - Primary danger color
- `--danger-primary-hover` - Hover state danger color
- `--danger-color` - Danger color
- `--danger-color-hover` - Hover state danger color
- `--warning-color` - Warning color

---

## Dynamic Theme System (JavaScript)

**File:** `frontend/dist/js/theme.js`

### Theme Definitions (3 themes)

#### Theme 1: Teal
- **color:** `#14b8a6`
- **glow:** `rgba(20, 184, 166, 0.5)`
- **border:** `rgba(20, 184, 166, 0.3)`
- **textGlow:** `rgba(20, 184, 166, 0.4)`
- **secondary:** `#06b6d4`

#### Theme 2: Cyan
- **color:** `#06b6d4`
- **glow:** `rgba(6, 182, 212, 0.5)`
- **border:** `rgba(6, 182, 212, 0.3)`
- **textGlow:** `rgba(6, 182, 212, 0.4)`
- **secondary:** `#14b8a6`

#### Theme 3: Purple
- **color:** `#a855f7`
- **glow:** `rgba(168, 85, 247, 0.5)`
- **border:** `rgba(168, 85, 247, 0.3)`
- **textGlow:** `rgba(168, 85, 247, 0.4)`
- **secondary:** `#ec4899`

### Animation Details
- **Fade Duration:** 4000ms (4 seconds)
- **Pause Duration:** 1000ms (1 second)
- **Background Color:** `#0a0a0f` (fade transition target)
- **Easing Functions:**
  - `easeInOutCubic` - For color transitions
  - `easeInOutQuad` - For border transitions
- **Target FPS:** 30fps
- **Status:** Automatic cycling disabled (line 156 commented out to prevent flickering)

---

## Color Palette Summary

### Hex Colors
- `#14b8a6` - Teal (Theme 1 primary)
- `#06b6d4` - Cyan (Theme 2 primary)
- `#a855f7` - Purple (Theme 3 primary)
- `#ec4899` - Pink (secondary accent)
- `#0a0a0f` - Dark background (fade target)
- `#e0e0e0` - Static text
- `#a0a0b0` - Secondary text
- `#1a1a2e` - Dark input background
- `#0066ff` - Electric blue (loading progress)
- `#00ccff` - Cyan blue (loading progress)
- `#00ffff` - Bright cyan (loading progress)
- `#00ff66` - Green (success state)
- `#ff3333` - Red (error state)
- `#ff6b6b` - Light red (danger state)
- `#f87171` - Soft red (warning state)
- `#4ade80` - Green (success indicator)
- `#fbbf24` - Yellow (warning indicator)
- `#eab308` - Yellow (tutorial warning)
- `#7aa2f7` - Blue (entity action)
- `#8888a0` - Gray (disabled/placeholder)
- `#c0a0a0` - Brownish gray (negative prompt)
- `#4ade80` - Green (loaded status)

### RGBA Colors (Transparency)
- **Dark Backgrounds:**
  - `rgba(0, 0, 0, 0.35)` - Input background
  - `rgba(0, 0, 0, 0.25)` - Canvas wrapper
  - `rgba(0, 0, 0, 0.45)` - Upload zone
  - `rgba(0, 0, 0, 0.5)` - Various overlays
  - `rgba(0, 0, 0, 0.6)` - Status background
  - `rgba(0, 0, 0, 0.7)` - Modal overlay
  - `rgba(0, 0, 0, 0.85)` - Tutorial overlay
  - `rgba(0, 0, 0, 0.9)` - Loading progress track

- **Panel Backgrounds:**
  - `rgba(21, 21, 31, 0.98)` - Main modal/container
  - `rgba(21, 21, 31, 0.95)` - Form inputs
  - `rgba(21, 21, 31, 0.9)` - Metadata viewer
  - `rgba(21, 21, 31, 0.8)` - Session items
  - `rgba(21, 21, 31, 0.7)` - Panels
  - `rgba(21, 21, 31, 0.6)` - Cards
  - `rgba(15, 21, 31, 0.95)` - Lorebook creation
  - `rgba(15, 21, 31, 0.8)` - Type items
  - `rgba(15, 21, 31, 0.6)` - Type selector
  - `rgba(15, 21, 31, 0.52)` - Editor sections
  - `rgba(15, 21, 31, 0.4)` - Card content
  - `rgba(15, 15, 23, 0.6)` - Tutorial nav
  - `rgba(15, 15, 23, 0.8)` - Code block background
  - `rgba(10, 10, 15, 0.5)` - Headers/footers
  - `rgba(10, 10, 15, 0.35)` - Sidebar
  - `rgba(10, 10, 15, 0.3)` - Toggle slider

- **Borders:**
  - `rgba(100, 100, 100, 0.5)` - Strong border
  - `rgba(100, 100, 100, 0.4)` - Medium border
  - `rgba(100, 100, 100, 0.3)` - Standard border
  - `rgba(100, 100, 100, 0.25)` - Light border
  - `rgba(100, 100, 100, 0.2)` - Very light border
  - `rgba(100, 100, 100, 0.15)` - Subtle border
  - `rgba(100, 100, 100, 0.12)` - Minimal border
  - `rgba(100, 100, 100, 0.1)` - Faint border
  - `rgba(100, 100, 100, 0.08)` - Very faint border
  - `rgba(100, 100, 100, 0.05)` - Barely visible border
  - `rgba(100, 100, 100, 0.02)` - Almost invisible border

- **Text:**
  - `rgba(212, 212, 212, 0.9)` - High opacity text
  - `rgba(212, 212, 212, 0.8)` - Medium-high opacity text
  - `rgba(212, 212, 212, 0.7)` - Standard secondary text
  - `rgba(212, 212, 212, 0.68)` - Slightly muted text
  - `rgba(212, 212, 212, 0.65)` - Muted text
  - `rgba(212, 212, 212, 0.6)` - More muted text
  - `rgba(212, 212, 212, 0.55)` - Section header text
  - `rgba(212, 212, 212, 0.5)` - Placeholder text
  - `rgba(212, 212, 212, 0.45)` - Breadcrumb text
  - `rgba(212, 212, 212, 0.4)` - Very muted text
  - `rgba(212, 212, 212, 0.3)` - Extremely muted text
  - `rgba(212, 212, 212, 0.2)` - Almost invisible text
  - `rgba(212, 212, 212, 0.1)` - Barely visible text

- **Accents:**
  - `rgba(20, 184, 166, 0.5)` - Theme glow (teal)
  - `rgba(20, 184, 166, 0.4)` - Theme text glow
  - `rgba(20, 184, 166, 0.3)` - Theme border
  - `rgba(20, 184, 166, 0.2)` - Light theme accent
  - `rgba(20, 184, 166, 0.18)` - Hover theme accent
  - `rgba(20, 184, 166, 0.15)` - Light theme background
  - `rgba(20, 184, 166, 0.12)` - Very light theme
  - `rgba(20, 184, 166, 0.1)` - Minimal theme
  - `rgba(20, 184, 166, 0.08)` - Faint theme
  - `rgba(6, 182, 212, 0.5)` - Theme glow (cyan)
  - `rgba(6, 182, 212, 0.3)` - Theme border (cyan)
  - `rgba(6, 182, 212, 0.1)` - Ambient orb 2
  - `rgba(168, 85, 247, 0.5)` - Theme glow (purple)
  - `rgba(168, 85, 247, 0.3)` - Theme border (purple)
  - `rgba(100, 140, 220, 0.08)` - Entity action background
  - `rgba(100, 140, 220, 0.18)` - Entity action hover
  - `rgba(100, 140, 220, 0.25)` - Entity action border

- **Danger/Error:**
  - `rgba(239, 68, 68, 0.8)` - Danger text
  - `rgba(239, 68, 68, 0.6)` - Error border
  - `rgba(239, 68, 68, 0.5)` - Error text
  - `rgba(239, 68, 68, 0.4)` - Error box-shadow
  - `rgba(239, 68, 68, 0.3)` - Danger border hover
  - `rgba(239, 68, 68, 0.2)` - Danger background
  - `rgba(239, 68, 68, 0.12)` - Light danger
  - `rgba(220, 50, 50, 0.8)` - Delete button background
  - `rgba(220, 50, 50, 0.4)` - Delete button border
  - `rgba(220, 50, 50, 0.3)` - Warning border
  - `rgba(220, 50, 50, 0.2)` - Warning background
  - `rgba(220, 50, 50, 0.18)` - Warning hover
  - `rgba(255, 67, 54, 0.2)` - Delete hover
  - `rgba(255, 50, 50, 0.6)` - Error border (loading)
  - `rgba(255, 50, 50, 0.4)` - Error shadow (loading)

- **Success/Info:**
  - `rgba(0, 255, 100, 0.6)` - Success border
  - `rgba(0, 255, 100, 0.4)` - Success shadow
  - `rgba(0, 150, 255, 0.6)` - Loading border
  - `rgba(0, 150, 255, 0.4)` - Loading shadow
  - `rgba(0, 150, 255, 0.3)` - Loading border base
  - `rgba(0, 150, 255, 0.2)` - Loading shadow base
  - `rgba(234, 179, 8, 0.25)` - Warning border (tutorial)
  - `rgba(234, 179, 8, 0.08)` - Warning background (tutorial)
  - `rgba(30, 30, 35, 0.8)` - Dropdown image background

- **Miscellaneous:**
  - `rgba(255, 255, 255, 0.4)` - Shimmer animation
  - `rgba(255, 255, 255, 0.3)` - Scrollbar hover
  - `rgba(255, 255, 255, 0.2)` - Scrollbar thumb
  - `rgba(255, 255, 255, 0.1)` - Hover background
  - `rgba(255, 255, 255, 0.05)` - Light hover
  - `rgba(255, 255, 255, 0.02)` - Very light hover
  - `rgba(255, 255, 255, 0.01)` - Minimal hover

---

## Background Colors

### Solid Backgrounds
- `#0a0a0f` - Fade transition background
- `#1a1a2e` - Select option background
- `#e0e0e0` - Static text (used as background in some contexts)

### Transparent/Dark Backgrounds
- `rgba(0, 0, 0, 0.35)` - Input fields
- `rgba(0, 0, 0, 0.25)` - Canvas wrapper
- `rgba(0, 0, 0, 0.45)` - Upload zones
- `rgba(0, 0, 0, 0.5)` - Overlays
- `rgba(0, 0, 0, 0.6)` - Status backgrounds
- `rgba(0, 0, 0, 0.7)` - Modal overlays
- `rgba(0, 0, 0, 0.85)` - Tutorial overlay
- `rgba(0, 0, 0, 0.9)` - Loading progress track
- `rgba(10, 10, 15, 0.5)` - Headers/footers
- `rgba(10, 10, 15, 0.35)` - Sidebar
- `rgba(10, 10, 15, 0.3)` - Toggle slider
- `rgba(15, 15, 23, 0.6)` - Tutorial nav
- `rgba(15, 15, 23, 0.8)` - Code blocks
- `rgba(15, 21, 31, 0.95)` - Lorebook creation
- `rgba(15, 21, 31, 0.8)` - Type items
- `rgba(15, 21, 31, 0.6)` - Type selector
- `rgba(15, 21, 31, 0.52)` - Editor sections
- `rgba(15, 21, 31, 0.4)` - Card content
- `rgba(21, 21, 31, 0.98)` - Main modals
- `rgba(21, 21, 31, 0.95)` - Form inputs
- `rgba(21, 21, 31, 0.9)` - Metadata viewer
- `rgba(21, 21, 31, 0.8)` - Session items
- `rgba(21, 21, 31, 0.7)` - Panels
- `rgba(21, 21, 31, 0.6)` - Cards
- `rgba(30, 30, 35, 0.8)` - Dropdown images

### Hover Backgrounds
- `rgba(20, 184, 166, 0.15)` - Theme accent hover
- `rgba(20, 184, 166, 0.12)` - Theme accent light hover
- `rgba(20, 184, 166, 0.08)` - Theme accent faint hover
- `rgba(255, 255, 255, 0.1)` - Generic hover
- `rgba(255, 255, 255, 0.05)` - Light hover
- `rgba(255, 255, 255, 0.02)` - Very light hover
- `rgba(100, 100, 100, 0.3)` - Gray hover
- `rgba(100, 100, 100, 0.2)` - Light gray hover
- `rgba(100, 100, 100, 0.15)` - Very light gray hover
- `rgba(100, 100, 100, 0.12)` - Minimal gray hover
- `rgba(100, 100, 100, 0.08)` - Faint gray hover

### Special Backgrounds
- `rgba(239, 68, 68, 0.12)` - Generating button state
- `rgba(220, 50, 50, 0.2)` - Danger hover
- `rgba(220, 50, 50, 0.1)` - Close button hover
- `rgba(100, 140, 220, 0.08)` - Entity action
- `rgba(100, 140, 220, 0.18)` - Entity action hover

---

## Border Colors

### Standard Borders (Gray Scale)
- `rgba(100, 100, 100, 0.5)` - Strong border
- `rgba(100, 100, 100, 0.4)` - Medium border
- `rgba(100, 100, 100, 0.3)` - Standard border (most common)
- `rgba(100, 100, 100, 0.25)` - Light border
- `rgba(100, 100, 100, 0.2)` - Very light border
- `rgba(100, 100, 100, 0.15)` - Subtle border
- `rgba(100, 100, 100, 0.12)` - Minimal border
- `rgba(100, 100, 100, 0.1)` - Faint border
- `rgba(100, 100, 100, 0.08)` - Very faint border
- `rgba(100, 100, 100, 0.05)` - Barely visible border
- `rgba(100, 100, 100, 0.02)` - Almost invisible border

### Theme Borders (Dynamic)
- `rgba(20, 184, 166, 0.3)` - Teal theme border
- `rgba(20, 184, 166, 0.35)` - Teal theme border (slightly stronger)
- `rgba(20, 184, 166, 0.4)` - Teal theme border (strong)
- `rgba(6, 182, 212, 0.3)` - Cyan theme border
- `rgba(168, 85, 247, 0.3)` - Purple theme border
- `rgba(20, 184, 166, 0.25)` - Light theme border
- `rgba(100, 140, 220, 0.25)` - Entity action border

### Danger/Error Borders
- `rgba(239, 68, 68, 0.6)` - Error border
- `rgba(239, 68, 68, 0.3)` - Danger border
- `rgba(239, 68, 68, 0.5)` - Danger border hover
- `rgba(220, 50, 50, 0.4)` - Delete button border
- `rgba(220, 50, 50, 0.3)` - Warning border
- `rgba(255, 50, 50, 0.6)` - Error border (loading)

### Success/Info Borders
- `rgba(0, 255, 100, 0.6)` - Success border
- `rgba(0, 150, 255, 0.6)` - Loading border
- `rgba(0, 150, 255, 0.3)` - Loading border base
- `rgba(234, 179, 8, 0.25)` - Warning border (tutorial)
- `rgba(20, 184, 166, 0.2)` - Code border (tutorial)

### Border Styles
- `1px solid` - Standard solid border
- `2px solid` - Stronger solid border
- `3px solid` - Emphasis border (active states)
- `2px dashed` - Dashed border (placeholders)
- `1px dashed` - Light dashed border

---

## Text/Font Colors

### Primary Text Colors
- `#e0e0e0` - Static text (default)
- `var(--static-text)` - Dynamic static text
- `rgba(212, 212, 212, 1.0)` - Full opacity text
- `rgba(212, 212, 212, 0.9)` - High opacity text
- `rgba(212, 212, 212, 0.8)` - Medium-high opacity text

### Secondary Text Colors
- `#a0a0b0` - Secondary text (default)
- `var(--text-secondary)` - Dynamic secondary text
- `rgba(212, 212, 212, 0.7)` - Standard secondary text
- `rgba(212, 212, 212, 0.68)` - Slightly muted text
- `rgba(212, 212, 212, 0.65)` - Muted text
- `rgba(212, 212, 212, 0.6)` - More muted text
- `rgba(212, 212, 212, 0.55)` - Section header text
- `rgba(212, 212, 212, 0.5)` - Placeholder text
- `rgba(212, 212, 212, 0.45)` - Breadcrumb text
- `rgba(212, 212, 212, 0.4)` - Very muted text
- `rgba(212, 212, 212, 0.3)` - Extremely muted text
- `rgba(212, 212, 212, 0.2)` - Almost invisible text
- `rgba(212, 212, 212, 0.1)` - Barely visible text
- `rgba(100, 100, 100, 0.5)` - Gray placeholder text

### Theme Text Colors
- `#14b8a6` - Teal theme text
- `#06b6d4` - Cyan theme text
- `#a855f7` - Purple theme text
- `#ec4899` - Pink theme text
- `var(--theme-color)` - Dynamic theme color
- `var(--theme-secondary)` - Dynamic secondary theme

### Status Text Colors
- `#00ccff` - Loading status text
- `#00ff66` - Success status text
- `#ff3333` - Error status text
- `#f87171` - Warning status text
- `#fbbf24` - Yellow warning text
- `#eab308` - Tutorial warning text
- `#ff6b6b` - Danger text
- `#c0a0a0` - Negative prompt text
- `#7aa2f7` - Entity action text
- `#8888a0` - Disabled/placeholder text
- `#4ade80` - Green success text

### Font Weights
- `400` - Normal/Regular
- `500` - Medium
- `600` - Semi Bold
- `700` - Bold
- `bold` - Bold (keyword)

### Font Sizes
- `10px` - Key hints, chips
- `11px` - Meta info, tags, small labels
- `12px` - Small text, secondary labels
- `13px` - Standard text, labels
- `14px` - Body text, buttons
- `15px` - Slightly larger text
- `16px` - Section headers, small titles
- `18px` - Titles, headers
- `20px` - Large headers
- `24px` - Very large headers
- `28px` - Extra large headers
- `2.5rem` - Hero headers (40px)

### Font Families
- `'Segoe UI', Tahoma, Geneva, Verdana, sans-serif` - System sans-serif
- `'Courier New', monospace` - Monospace for code
- `'JetBrains Mono', 'Fira Code', 'Consolas', monospace` - Code fonts

---

## Shadows & Glows

### Box Shadows
- `0 0 8px var(--theme-glow)` - Standard theme glow
- `0 0 10px var(--theme-glow)` - Medium theme glow
- `0 0 12px var(--theme-glow)` - Strong theme glow
- `0 0 15px var(--theme-glow)` - Very strong theme glow
- `0 0 20px var(--theme-glow)` - Extra strong theme glow
- `0 0 30px rgba(0, 0, 0, 0.5)` - Modal shadow
- `0 0 40px rgba(0, 0, 0, 0.6)` - Large modal shadow
- `0 4px 12px rgba(0, 0, 0, 0.1)` - Card shadow
- `0 2px 8px rgba(20, 184, 166, 0.2)` - History item hover
- `0 4px 24px rgba(0, 0, 0, 0.3)` - Canvas image shadow
- `0 0 12px #00ccff` - Loading glow
- `0 0 24px #0066ff` - Loading glow outer
- `0 0 20px #00ccff` - Loading pulse mid
- `0 0 40px #0066ff` - Loading pulse outer
- `0 0 60px #00ffff` - Loading pulse extra
- `0 0 16px rgba(0, 150, 255, 0.4)` - Loading state
- `0 0 16px rgba(255, 50, 50, 0.4)` - Error state
- `0 0 16px rgba(0, 255, 100, 0.4)` - Completed state
- `0 2px 0 rgba(0, 0, 0, 0.3)` - Keyboard key shadow

### Text Shadows
- `0 0 8px var(--theme-glow)` - Standard text glow
- `0 0 10px var(--theme-glow)` - Medium text glow
- `0 0 8px #0066ff` - Loading text shadow
- `0 0 8px #00ff66` - Success text shadow
- `0 0 8px #ff3333` - Error text shadow

### Drop Shadows (Filter)
- `drop-shadow(0 0 8px var(--theme-glow))` - Active nav icon filter
- `filter: brightness(0) invert(1)` - Icon visibility on dark backgrounds

---

## Gradients

### Linear Gradients
- `linear-gradient(90deg, #0066ff, #00ccff, #00ffff)` - Loading progress fill
- `linear-gradient(90deg, #00ff66, #00ccff)` - Completed state
- `linear-gradient(90deg, #ff3333, #ff6666)` - Error state
- `linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent)` - Shimmer animation

---

## Font Properties

### Font Families
1. **System Sans-Serif:** `'Segoe UI', Tahoma, Geneva, Verdana, sans-serif`
2. **Monospace:** `'Courier New', monospace`
3. **Code Fonts:** `'JetBrains Mono', 'Fira Code', 'Consolas', monospace`

### Font Weights
- `400` (normal) - Regular text
- `500` (medium) - Medium weight
- `600` (semi-bold) - Semi-bold headers
- `700` (bold) - Bold headers
- `bold` - Bold keyword

### Font Sizes
- `10px` - Key hints, chips, badges
- `11px` - Meta info, tags, small labels
- `12px` - Small text, secondary labels, buttons
- `13px` - Standard text, labels
- `14px` - Body text, standard buttons
- `15px` - Slightly larger text
- `16px` - Section headers, small titles
- `18px` - Titles, modal headers
- `20px` - Large headers
- `22px` - Mobile headers
- `24px` - Very large headers
- `28px` - Extra large headers
- `2.5rem` (40px) - Hero headers

### Text Transform
- `uppercase` - Loading status, section headers
- `lowercase` - (not commonly used)
- `capitalize` - (not commonly used)

### Letter Spacing
- `1px` - Loading status
- `0.5px` - Section headers

### Line Height
- `1.0` - (rare)
- `1.3` - Standard text
- `1.4` - Body text
- `1.45` - Lorebook text
- `1.5` - Textareas
- `1.6` - Content text
- `1.7` - Tutorial text

---

## Inline Styles (HTML)

**File:** `frontend/dist/index.html`

### Inline Style Found (Line 140)
```html
<div style="display: flex; align-items: center; gap: 12px;">
```
- **display:** flex
- **align-items:** center
- **gap:** 12px

### Inline Style Found (Line 141)
```html
<img src="assets/HOWL.png" alt="HOWL Chat" style="width: 44px; height: 44px;">
```
- **width:** 44px
- **height:** 44px

---

## JavaScript-Generated Styles

### theme.js
**Dynamic CSS Property Setting:**
```javascript
document.documentElement.style.setProperty('--theme-color', nextTheme.color);
document.documentElement.style.setProperty('--theme-glow', nextTheme.glow);
document.documentElement.style.setProperty('--theme-text-glow', nextTheme.textGlow);
document.documentElement.style.setProperty('--theme-secondary', nextTheme.secondary);
```

### loading-progress.js
**Dynamic Style Setting:**
```javascript
progressFill.style.width = progress + '%';
```

### image-generator.js
**Dynamic Style Setting:**
```javascript
document.getElementById('ig-meta').style.display = 'block';
document.getElementById('ig-actions').style.display = 'flex';
item.innerHTML = `<img ... style="width:100%;height:100%;object-fit:cover;">`;
```

---

## Component-Specific Styles

### Ambient Effects (ambient.css)
- **Orb 1:** `rgba(20, 184, 166, 0.1)` with blur(96px)
- **Orb 2:** `rgba(6, 182, 212, 0.1)` with blur(96px)
- **Hexagon Grid:** SVG pattern with `#14b8a6` stroke
- **Opacity:** 0.3 for hexagon grid

### Loading Progress (loading-progress.css)
- **Background:** `rgba(10, 10, 15, 0.9)`
- **Border:** `rgba(0, 150, 255, 0.3)`
- **Fill Gradient:** `linear-gradient(90deg, #0066ff, #00ccff, #00ffff)`
- **Shadow:** `0 0 12px #00ccff, 0 0 24px #0066ff`
- **Text Color:** `#00ccff`
- **Text Shadow:** `0 0 8px #0066ff`

### Image Generator (image-generator.css)
- **Input Background:** `rgba(0, 0, 0, 0.35)`
- **Slider Thumb:** `var(--theme-color, #14b8a6)`
- **Canvas Wrapper:** `rgba(0, 0, 0, 0.25)`
- **Upload Zone:** `rgba(0, 0, 0, 0.45)` with dashed border
- **Action Buttons:** `rgba(20, 184, 166, 0.08)` background
- **Danger Actions:** `rgba(220, 50, 50, 0.08)` background
- **Entity Actions:** `rgba(100, 140, 220, 0.08)` background

### Lorebooks (lorebooks.css)
- **Card Background:** `rgba(15, 21, 31, 0.4)`
- **Chip Background:** `rgba(20, 184, 166, 0.14)`
- **Chip Muted:** `rgba(100, 100, 100, 0.15)`
- **Section Background:** `rgba(15, 21, 31, 0.52)`
- **Type Item Background:** `rgba(15, 21, 31, 0.8)`

### Story Tools (story-tools.css)
- Uses CSS custom properties for theming
- `--surface-primary`, `--surface-secondary`, `--surface-tertiary`
- `--surface-hover`
- `--border-primary`, `--border-secondary`
- `--accent-primary`, `--accent-secondary`
- `--danger-primary`, `--danger-color`
- `--warning-color`

### Tutorials (tutorials.css)
- **Overlay:** `rgba(0, 0, 0, 0.85)`
- **Container:** `rgba(21, 21, 31, 0.98)`
- **Nav Background:** `rgba(15, 15, 23, 0.6)`
- **Code Background:** `rgba(15, 15, 23, 0.8)`
- **Info Box:** `rgba(20, 184, 166, 0.08)`
- **Warning Box:** `rgba(234, 179, 8, 0.08)`
- **Breadcrumb:** `rgba(212, 212, 212, 0.45)`

### Modals/Popups
- **Overlay:** `rgba(0, 0, 0, 0.7)` with backdrop-filter blur(4px)
- **Container:** `rgba(21, 21, 31, 0.98)`
- **Header/Footer:** `rgba(10, 10, 15, 0.5)`
- **Border:** `rgba(100, 100, 100, 0.3)`
- **Shadow:** `0 0 30px rgba(0, 0, 0, 0.5)`

### Buttons
- **Primary:** `var(--theme-glow)` background
- **Secondary:** `rgba(100, 100, 100, 0.2)` background
- **Danger:** `rgba(239, 68, 68, 0.2)` background
- **Hover:** Increased opacity and glow

### Inputs/Forms
- **Background:** `rgba(21, 21, 31, 0.95)`
- **Border:** `rgba(100, 100, 100, 0.3)`
- **Focus:** Border changes to `var(--theme-color)` with glow
- **Placeholder:** `rgba(100, 100, 100, 0.5)`

### Scrollbars
- **Width:** 4px, 6px, or 8px depending on context
- **Track:** Transparent
- **Thumb:** `rgba(255, 255, 255, 0.2)`
- **Thumb Hover:** `rgba(255, 255, 255, 0.3)`

---

## Pattern & Shape Styles

### Hexagon Grid Pattern
- **SVG Data URL:** Embedded hexagon pattern
- **Stroke:** `#14b8a6`
- **Stroke Width:** 0.5
- **Size:** 20px × 20px
- **Animation:** 60s linear infinite rotation

### Border Radius
- `4px` - Small elements (inputs, buttons)
- `5px` - Medium elements
- `6px` - Standard elements
- `8px` - Larger elements (cards, modals)
- `10px` - Large elements
- `12px` - Extra large elements (main containers)
- `16px` - Very large elements (tutorial container)
- `50%` - Circular elements (orbs, avatars)

### Shapes
- **Circular:** Ambient orbs, avatar images
- **Rounded Rectangles:** Cards, buttons, inputs
- **Hexagonal:** Background pattern (SVG)

---

## Animation Styles

### Keyframe Animations
- `menubar-pulse` - Menu bar glow pulse
- `navbar-pulse` - Navbar glow pulse
- `panel-pulse` - Panel glow pulse
- `pulse` - General pulse animation
- `ambient-pulse-1` - Ambient orb 1 pulse
- `ambient-pulse-2` - Ambient orb 2 pulse
- `hexagon-rotate` - Hexagon grid rotation
- `shimmer` - Loading progress shimmer
- `pulse-glow` - Loading progress glow pulse
- `fadeIn` - Fade in animation
- `ig-spin` - Image generator spinner
- `tutorialsFadeIn` - Tutorial fade in

### Transitions
- `0.15s ease` - Fast transitions (buttons, hover)
- `0.2s ease` - Standard transitions
- `0.25s ease` - Medium transitions
- `0.3s ease` - Slow transitions
- `0.2s ease` - Transform transitions
- `width 0.3s linear` - Progress bar
- `transform 0.2s` - Transform only

---

## Summary Statistics

### Total CSS Files Analyzed: 19
1. animations.css
2. icons.css
3. character-header.css
4. world-header.css
5. scenario-header.css
6. world-form.css
7. world-grid.css
8. world-bubble.css
9. model-settings.css
10. sampler-settings.css
11. real-settings.css
12. image-generator.css
13. image-generator-panels.css
14. world-dropdown-item.css
15. world-creation-popup.css
16. chat-session-list.css
17. chat-session-modal.css
18. loading-progress.css
19. lorebooks.css
20. story-tools.css
21. tutorials.css
22. ambient.css
23. nav-icon-images.css

### Total HTML Files: 9
- index.html (with inline styles)
- audio-settings.html
- character.html
- image-generator.html
- lorebooks.html
- real-settings.html
- scenario.html
- story-tools.html
- world.html

### Total JavaScript Files with Style Generation: 3
- theme.js (dynamic theme system)
- loading-progress.js (progress bar styles)
- image-generator.js (inline image styles)

### Color Categories
- **Hex Colors:** 25 unique values
- **RGBA Colors:** 100+ unique combinations
- **CSS Variables:** 20+ custom properties
- **Dynamic Themes:** 3 theme presets

---

## Recommendations

1. **Consolidate Similar RGBA Values:** Many similar opacity levels could be standardized
2. **Create Color Tokens:** Consider using CSS custom properties for common colors
3. **Standardize Border Widths:** Mix of 1px, 2px, 3px could be reduced to 2 standards
4. **Unify Font Sizes:** Reduce from 10+ sizes to 5-6 standard sizes
5. **Shadow System:** Create a shadow scale (small, medium, large, x-large)
6. **Theme System:** The existing 3-theme system is well-implemented
7. **Accessibility:** Ensure sufficient contrast ratios for all text colors

---

**End of Audit**
