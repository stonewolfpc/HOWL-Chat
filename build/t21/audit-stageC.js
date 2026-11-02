// Stage C theme audit — the exact probe used on the remaining pages (2026-11-02).
//
// Run inside the harness page (build/bin/howl-ui-harness.exe, http://127.0.0.1:8799)
// via the browser's evaluate channel:
//
//   (async () => { <this file's body> })()
//
// Method, and why each part exists:
//   * transitions/animations are disabled and the caller waits >= 1200 ms before
//     reading. theme.js drives a 4s custom-property transition; 600 ms returned
//     the PREVIOUS theme's colours and produced one false "stuck light in dark"
//     finding on the memories flyout. Anything below 1s is not trustworthy.
//   * the theme is flipped through theme.js's own setThemeMode(), never by
//     toggling the body class, so a missing/renamed entry point fails loudly
//     instead of silently passing (the stage's recorded audit lesson).
//   * text nodes are measured against the RESOLVED ancestor background
//     (alpha-composited down the tree), which also covers labels inside filled
//     child elements - the case that hid white-on-teal in the world page.
//   * offender = luminance delta < 60 (0-255 scale) for text at 16px or larger.
//   * "stuck dark" = an opaque area > 6000 px^2 whose resolved background
//     luminance is < 25; expected in dark mode, must be 0 in light mode.
//   * background resolution is memoised per element: recomputing the ancestor
//     chain per node blew the 10s evaluation budget on index/story-tools.
//
// One tab only: every harness page holds an open /__events SSE stream, so six
// tabs exhaust the browser's six-connection budget and later pages never finish
// loading. That produced two false "page never loads" results before it was
// identified.
(function auditScript() {
  const wait = ms => new Promise(r => setTimeout(r, ms));
  if (!document.getElementById('__audit_css')) {
    const s = document.createElement('style');
    s.id = '__audit_css';
    s.textContent = '*,*::before,*::after{transition:none!important;animation:none!important}';
    document.head.appendChild(s);
  }
  const parse = c => {
    const m = /rgba?\(([^)]+)\)/.exec(c || '');
    if (!m) return null;
    const p = m[1].split(',').map(parseFloat);
    return { r: p[0], g: p[1], b: p[2], a: p[3] === undefined ? 1 : p[3] };
  };
  const rel = c => {
    const f = v => { v /= 255; return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4); };
    return 0.2126 * f(c.r) + 0.7152 * f(c.g) + 0.0722 * f(c.b);
  };
  const over = (t, u) => t.a >= 1 ? t : {
    r: t.r * t.a + u.r * (1 - t.a),
    g: t.g * t.a + u.g * (1 - t.a),
    b: t.b * t.a + u.b * (1 - t.a),
    a: 1
  };
  const measure = () => {
    const cache = new Map();
    const bgOf = e => {
      if (cache.has(e)) return cache.get(e);
      const stack = [];
      let n = e, res = null;
      while (n && n.nodeType === 1) {
        if (cache.has(n)) { res = cache.get(n); break; }
        stack.push(n); n = n.parentElement;
      }
      if (!res) res = { r: 255, g: 255, b: 255, a: 1 };
      for (let i = stack.length - 1; i >= 0; i--) {
        const c = parse(getComputedStyle(stack[i]).backgroundColor);
        if (c && c.a > 0.001) res = over(c, res);
        cache.set(stack[i], res);
      }
      return cache.get(e) || res;
    };
    const low = [], dark = [];
    for (const e of document.body.querySelectorAll('*')) {
      const cs = getComputedStyle(e);
      if (cs.display === 'none' || cs.visibility === 'hidden' || parseFloat(cs.opacity) < 0.3) continue;
      const r = e.getBoundingClientRect();
      if (r.width < 4 || r.height < 4) continue;
      const bg = bgOf(e);
      if (rel(bg) * 255 < 25 && r.width * r.height > 6000 && parse(parse(cs.backgroundColor) ? cs.backgroundColor : '')?.a > 0.5) {
        dark.push((e.tagName + (e.id ? '#' + e.id : e.className ? '.' + String(e.className).split(' ')[0] : '')).slice(0, 40));
      }
      const hasText = [...e.childNodes].some(n => n.nodeType === 3 && n.textContent.trim().length > 1);
      if (!hasText) continue;
      if (parseFloat(cs.fontSize) < 16) continue;
      const fg = parse(cs.color);
      if (!fg || fg.a < 0.3) continue;
      const delta = Math.abs(rel(fg) - rel(bg)) * 255;
      if (delta < 60) {
        low.push({
          sel: e.tagName.toLowerCase() + (e.id ? '#' + e.id : ''),
          text: e.textContent.trim().slice(0, 40),
          fg: cs.color,
          bg: 'rgb(' + [bg.r, bg.g, bg.b].map(v => Math.round(v)).join(',') + ')',
          delta: +delta.toFixed(1)
        });
      }
    }
    return {
      bodyBg: getComputedStyle(document.body).backgroundColor,
      bodyClass: document.body.className,
      lowContrastCount: low.length,
      lowContrast: low.slice(0, 10),
      darkAreaCount: dark.length,
      darkAreas: dark.slice(0, 8)
    };
  };
  window.setThemeMode('light'); return wait(1200).then(() => {
    const light = measure();
    window.setThemeMode('dark');
    return wait(1200).then(() => {
      const dark = measure();
      window.setThemeMode('light');
      return { page: location.pathname, light, dark };
    });
  });
})
