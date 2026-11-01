/**
 * HTML Utilities — single owner of escapeHtml
 *
 * Six files used to define their own escapeHtml: two different algorithms
 * (a textContent round-trip and a regex replace) across four variants. Which
 * one a page got depended on script order, and the div-based version does not
 * escape quotes, so the same call was safe in a text node and unsafe inside an
 * attribute. There is one implementation now, and it escapes both contexts.
 *
 * Load this before any script that renders user or model text.
 *
 * @module html-utils
 * @version 1.0.0
 */

'use strict';

var HTML_ESCAPES = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;'
};

function escapeHtml(value) {
    if (value === null || value === undefined) return '';
    return String(value).replace(/[&<>"']/g, function (ch) { return HTML_ESCAPES[ch]; });
}
