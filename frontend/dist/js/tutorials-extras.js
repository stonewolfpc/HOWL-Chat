/**
 * HOWL Chat — Tutorial Extras
 *
 * Everything the base renderer does not cover: screenshots, paired
 * good/bad examples, section anchors, an on-this-page strip, prev/next
 * paging and click-to-zoom.
 *
 * Loads after tutorials-ui.js. The hooks are optional — tutorials-ui.js
 * calls extendTutorialMarkdown() and decorateTutorialTopic() only when
 * they are defined, so this file can be removed and nothing breaks.
 *
 * Content syntax this adds:
 *   ![alt text](img/tutorials/file.png "Caption shown under the shot")
 *   <good>one short block of example text</good>
 *   <bad>one short block of example text</bad>
 *
 * @module tutorials-extras
 */

'use strict';

(function () {

    var IMG = /!\[([^\]]*)\]\(([^)\s]+)(?:\s+"([^"]*)")?\)/g;

    function exampleBlock(label, cls, icon, body) {
        // Blank lines inside a wrapper would be split into separate <p> blocks
        // by the paragraph pass, so collapse them to keep the card intact.
        return '<div class="tutorials-example ' + cls + '">' +
            '<div class="tutorials-example-label"><span class="icon ' + icon + '"></span>' + label + '</div>' +
            '<div class="tutorials-example-body">' + body.replace(/\n{2,}/g, '\n') + '</div>' +
            '</div>';
    }

    window.extendTutorialMarkdown = function (html) {
        html = html.replace(/&lt;good&gt;([\s\S]*?)&lt;\/good&gt;/g, function (m, body) {
            return exampleBlock('Good', 'good', 'icon-check', body);
        });
        html = html.replace(/&lt;bad&gt;([\s\S]*?)&lt;\/bad&gt;/g, function (m, body) {
            return exampleBlock('Bad', 'bad', 'icon-alert', body);
        });
        html = html.replace(IMG, function (m, alt, src, caption) {
            var cap = caption || alt || '';
            return '<figure class="tutorials-figure" onclick="_openTutorialShot(this)">' +
                '<img src="' + src + '" alt="' + alt + '" loading="lazy">' +
                '<figcaption>' + cap + '</figcaption></figure>';
        });
        return html;
    };

    window.decorateTutorialTopic = function (contentEl, topic) {
        var heads = Array.prototype.slice.call(contentEl.querySelectorAll('h2'));
        var used = {};

        heads.forEach(function (h, i) {
            var base = (h.textContent || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '') || ('section-' + (i + 1));
            var id = base, n = 2;
            while (used[id]) { id = base + '-' + n; n++; }
            used[id] = true;
            h.id = id;
        });

        // Only long topics earn a jump strip.
        if (heads.length >= 3) {
            var toc = document.createElement('nav');
            toc.className = 'tutorials-toc';
            toc.innerHTML = '<div class="tutorials-toc-title">On this page</div>' + heads.map(function (h) {
                return '<a class="tutorials-toc-link" href="#' + h.id + '" onclick="_scrollToTutorialSection(event, \'' + h.id + '\')">' + h.textContent + '</a>';
            }).join('');
            var h1 = contentEl.querySelector('h1');
            if (h1 && h1.nextSibling) contentEl.insertBefore(toc, h1.nextSibling);
            else contentEl.insertBefore(toc, contentEl.firstChild);
        }

        var flat = [];
        ((window.TUTORIALS_DATA || {}).sections || []).forEach(function (section) {
            (section.topics || []).forEach(function (t) { flat.push(t); });
        });
        var idx = flat.findIndex(function (t) { return t.id === topic.id; });
        if (idx < 0) return;

        var prev = flat[idx - 1];
        var next = flat[idx + 1];
        var pager = document.createElement('div');
        pager.className = 'tutorials-pager';
        pager.innerHTML =
            (prev
                ? '<button class="tutorials-pager-btn prev" onclick="navigateToTutorial(\'' + prev.id + '\')"><span class="tutorials-pager-dir">Previous</span><span class="tutorials-pager-name">' + prev.title + '</span></button>'
                : '<span></span>') +
            (next
                ? '<button class="tutorials-pager-btn next" onclick="navigateToTutorial(\'' + next.id + '\')"><span class="tutorials-pager-dir">Next</span><span class="tutorials-pager-name">' + next.title + '</span></button>'
                : '<span></span>');
        contentEl.appendChild(pager);
    };

    window._scrollToTutorialSection = function (e, id) {
        e.preventDefault();
        var el = document.getElementById(id);
        if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    };

    window._openTutorialShot = function (figure) {
        var img = figure && figure.querySelector ? figure.querySelector('img') : null;
        if (!img) return;
        var capEl = figure.querySelector('figcaption');
        var caption = (capEl && capEl.textContent) || img.alt || '';

        var box = document.getElementById('tutorials-lightbox');
        if (!box) {
            box = document.createElement('div');
            box.id = 'tutorials-lightbox';
            box.className = 'tutorials-lightbox';
            box.innerHTML = '<figure><img alt=""><figcaption></figcaption></figure>';
            box.addEventListener('click', function () { box.classList.remove('show'); });
            document.body.appendChild(box);
        }
        box.querySelector('img').src = img.src;
        box.querySelector('img').alt = img.alt || '';
        box.querySelector('figcaption').textContent = caption;
        box.classList.add('show');
    };

    // Escape closes the zoom first, before the tutorials overlay handler runs.
    document.addEventListener('keydown', function (e) {
        if (e.key !== 'Escape') return;
        var box = document.getElementById('tutorials-lightbox');
        if (box && box.classList.contains('show')) {
            box.classList.remove('show');
            e.stopPropagation();
        }
    }, true);
})();
