#!/usr/bin/env node
/**
 * Static HTML Audit — dev tool, not shipped
 *
 * Checks all HTML files in frontend/dist for:
 *   1. <label for="X"> with no matching id="X" in the same file
 *   2. <img src="..."> pointing to a path that doesn't exist on disk
 *   3. <script src="..."> / <link href="..."> pointing to missing files
 *   4. Inline base64 data URLs stored directly in HTML (bloat indicator)
 *   5. Empty panel divs (id contains "panel" or "list", no children)
 *   6. TODO/placeholder text in visible elements
 *
 * Run: node tests/audit-html.js
 */

'use strict';

const fs   = require('fs');
const path = require('path');

const DIST = path.resolve(__dirname, '../frontend/dist');

const RED    = '\x1b[31m';
const YELLOW = '\x1b[33m';
const GREEN  = '\x1b[32m';
const RESET  = '\x1b[0m';
const BOLD   = '\x1b[1m';

let totalErrors   = 0;
let totalWarnings = 0;

function error(file, msg)   { console.log(`${RED}  ✗ ERROR${RESET}   ${file}: ${msg}`);   totalErrors++;   }
function warn(file, msg)    { console.log(`${YELLOW}  ⚠ WARN${RESET}    ${file}: ${msg}`);  totalWarnings++; }
function ok(msg)            { console.log(`${GREEN}  ✓${RESET} ${msg}`); }

function getHtmlFiles(dir) {
    const results = [];
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
        if (entry.isDirectory()) {
            results.push(...getHtmlFiles(path.join(dir, entry.name)));
        } else if (entry.name.endsWith('.html')) {
            results.push(path.join(dir, entry.name));
        }
    }
    return results;
}

// Minimal regex-based HTML parser — good enough for audit purposes
function auditHtmlFile(htmlPath) {
    const rel     = path.relative(DIST, htmlPath);
    const content = fs.readFileSync(htmlPath, 'utf8');
    let fileErrors = 0;

    // 1. Collect all id="..." values in the file
    const idSet = new Set();
    for (const m of content.matchAll(/\sid="([^"]+)"/g)) {
        idSet.add(m[1]);
    }

    // 2. Check <label for="X">
    for (const m of content.matchAll(/<label[^>]+\bfor="([^"]+)"/g)) {
        const forVal = m[1];
        if (!idSet.has(forVal)) {
            error(rel, `<label for="${forVal}"> — no element with id="${forVal}" found in this file`);
            fileErrors++;
        }
    }

    // 3. Check <img src="..."> — skip data URIs and empty
    for (const m of content.matchAll(/<img[^>]+\bsrc="([^"]+)"/g)) {
        const src = m[1];
        if (!src || src.startsWith('data:') || src.startsWith('http') || src.startsWith('{{')) continue;
        const resolved = path.resolve(path.dirname(htmlPath), src);
        if (!fs.existsSync(resolved)) {
            error(rel, `<img src="${src}"> — file not found on disk`);
            fileErrors++;
        }
    }

    // 4. Check <script src="...">
    for (const m of content.matchAll(/<script[^>]+\bsrc="([^"]+)"/g)) {
        const src = m[1];
        if (!src || src.startsWith('http') || src.startsWith('//')) continue;
        const resolved = path.resolve(path.dirname(htmlPath), src);
        if (!fs.existsSync(resolved)) {
            error(rel, `<script src="${src}"> — file not found on disk`);
            fileErrors++;
        }
    }

    // 5. Check <link href="..."> CSS
    for (const m of content.matchAll(/<link[^>]+\bhref="([^"]+)"/g)) {
        const href = m[1];
        if (!href || href.startsWith('http') || href.startsWith('//') || href.startsWith('#')) continue;
        const resolved = path.resolve(path.dirname(htmlPath), href);
        if (!fs.existsSync(resolved)) {
            error(rel, `<link href="${href}"> — file not found on disk`);
            fileErrors++;
        }
    }

    // 6. Inline base64 blobs in src attributes (not data: URIs in attributes — those are fine for tiny icons)
    const b64Matches = [...content.matchAll(/src="data:image\/[^;]+;base64,([^"]{200,})"/g)];
    if (b64Matches.length > 0) {
        warn(rel, `${b64Matches.length} large inline base64 image(s) hardcoded in HTML — move to CSS or JS`);
    }

    // 7. Duplicate nav-item images — same src appearing more than once inside .nav-sidebar
    const navSidebarMatch = content.match(/<div class="nav-sidebar">([\s\S]*?)<\/div>\s*\n?\s*\n?\s*<!--/);
    if (navSidebarMatch) {
        const navBlock = navSidebarMatch[1];
        const navSrcs = [...navBlock.matchAll(/src="([^"]+)"/g)].map(m => m[1]);
        const seen = {};
        for (const src of navSrcs) {
            seen[src] = (seen[src] || 0) + 1;
            if (seen[src] === 2) {
                error(rel, `Duplicate nav icon: src="${src}" appears ${seen[src]}x in nav-sidebar — extra orphaned nav-item`);
                fileErrors++;
            }
        }
    }

    // 8. Placeholder / TODO text in visible content
    const placeholders = ['TODO', 'PLACEHOLDER', 'Coming Soon', 'Not Implemented', 'TBD'];
    for (const p of placeholders) {
        const re = new RegExp(`>([^<]*${p}[^<]*)<`, 'i');
        if (re.test(content)) {
            warn(rel, `Contains placeholder text matching "${p}" in element content`);
        }
    }

    if (fileErrors === 0) ok(rel);
}

// ---------------------------------------------------------------------------
// JS file audit: check for unregistered function calls (pattern-based)
// ---------------------------------------------------------------------------

function auditJsFile(jsPath) {
    const rel     = path.relative(DIST, jsPath);
    const content = fs.readFileSync(jsPath, 'utf8');

    // Flag any remaining console.log with raw base64 data (accidental blob logging)
    if (/console\.log\([^)]*data:image/.test(content)) {
        warn(rel, 'console.log() contains base64 image data — will flood DevTools');
    }

    // Flag calling window.go.main.App.X where X is a known-removed function
    const removed = ['SetCharacters', 'UpdateScenarioContext', 'UpdateWorldContext'];
    for (const fn of removed) {
        if (content.includes(`App.${fn}`)) {
            error(rel, `Calls removed binding App.${fn}() — this will fail at runtime`);
        }
    }
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

console.log(`\n${BOLD}HOWL Chat — Static HTML/JS Audit${RESET}`);
console.log(`Scanning: ${DIST}\n`);

const htmlFiles = getHtmlFiles(DIST);
console.log(`${BOLD}HTML Files (${htmlFiles.length}):${RESET}`);
for (const f of htmlFiles) auditHtmlFile(f);

const jsDir   = path.join(DIST, 'js');
const jsFiles = fs.existsSync(jsDir)
    ? fs.readdirSync(jsDir).filter(f => f.endsWith('.js')).map(f => path.join(jsDir, f))
    : [];

console.log(`\n${BOLD}JS Files (${jsFiles.length}):${RESET}`);
for (const f of jsFiles) auditJsFile(f);

console.log(`\n${BOLD}Summary:${RESET}`);
console.log(`  ${RED}Errors:${RESET}   ${totalErrors}`);
console.log(`  ${YELLOW}Warnings:${RESET} ${totalWarnings}`);

if (totalErrors > 0) {
    console.log(`\n${RED}${BOLD}AUDIT FAILED — fix errors before shipping${RESET}\n`);
    process.exit(1);
} else {
    console.log(`\n${GREEN}${BOLD}AUDIT PASSED${RESET}\n`);
    process.exit(0);
}
