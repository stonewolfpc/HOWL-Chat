/**
 * ig-pose.js — OpenPose skeleton editor for ControlNet.
 *
 * Draws the COCO-18 keypoint layout that OpenPose ControlNet models were trained
 * on: coloured limbs on a black background, joints as dots. Drag a joint to pose
 * the figure, pick a preset, then send the render straight into the ControlNet
 * image slot. No external pose library, no network.
 *
 * Coordinates are normalised 0..1 so the same pose renders at any output size.
 */
(function () {
    'use strict';

    // COCO-18 order used by every OpenPose ControlNet checkpoint.
    const NOSE = 0, NECK = 1, RSHO = 2, RELB = 3, RWRI = 4, LSHO = 5, LELB = 6, LWRI = 7,
          RHIP = 8, RKNEE = 9, RANK = 10, LHIP = 11, LKNEE = 12, LANK = 13,
          REYE = 14, LEYE = 15, REAR = 16, LEAR = 17;
    const JOINT_COUNT = 18;

    // Limb pairs and the colours OpenPose uses; matching them matters more than
    // they look, because the ControlNet was trained on this exact palette.
    const LIMBS = [
        [NECK, RSHO, [255, 0, 0]],
        [NECK, LSHO, [255, 85, 0]],
        [RSHO, RELB, [255, 170, 0]],
        [RELB, RWRI, [255, 255, 0]],
        [LSHO, LELB, [170, 255, 0]],
        [LELB, LWRI, [85, 255, 0]],
        [NECK, RHIP, [0, 255, 0]],
        [RHIP, RKNEE, [0, 255, 85]],
        [RKNEE, RANK, [0, 255, 170]],
        [NECK, LHIP, [0, 255, 170]],
        [LHIP, LKNEE, [0, 170, 255]],
        [LKNEE, LANK, [0, 85, 255]],
        [NECK, NOSE, [0, 0, 255]],
        [NOSE, REYE, [85, 0, 255]],
        [REYE, REAR, [170, 0, 255]],
        [NOSE, LEYE, [255, 0, 255]],
        [LEYE, LEAR, [255, 0, 170]]
    ];

    const PRESETS = {
        standing: {
            label: 'Standing',
            joints: [
                [0.500, 0.105], [0.500, 0.185], [0.443, 0.205], [0.428, 0.350], [0.422, 0.485],
                [0.557, 0.205], [0.572, 0.350], [0.578, 0.485], [0.470, 0.520], [0.462, 0.700],
                [0.458, 0.880], [0.530, 0.520], [0.538, 0.700], [0.542, 0.880], [0.489, 0.094],
                [0.511, 0.094], [0.478, 0.100], [0.522, 0.100]
            ]
        },
        tpose: {
            label: 'T-pose',
            joints: [
                [0.500, 0.105], [0.500, 0.190], [0.330, 0.190], [0.185, 0.192], [0.045, 0.194],
                [0.670, 0.190], [0.815, 0.192], [0.955, 0.194], [0.470, 0.520], [0.466, 0.700],
                [0.462, 0.880], [0.530, 0.520], [0.534, 0.700], [0.538, 0.880], [0.489, 0.094],
                [0.511, 0.094], [0.478, 0.100], [0.522, 0.100]
            ]
        },
        sitting: {
            label: 'Sitting',
            joints: [
                [0.500, 0.215], [0.500, 0.290], [0.445, 0.305], [0.430, 0.420], [0.440, 0.520],
                [0.555, 0.305], [0.570, 0.420], [0.560, 0.520], [0.472, 0.560], [0.300, 0.585],
                [0.300, 0.760], [0.528, 0.560], [0.700, 0.585], [0.700, 0.760], [0.489, 0.204],
                [0.511, 0.204], [0.478, 0.210], [0.522, 0.210]
            ]
        },
        walking: {
            label: 'Walking',
            joints: [
                [0.530, 0.110], [0.512, 0.190], [0.455, 0.210], [0.405, 0.345], [0.375, 0.470],
                [0.570, 0.205], [0.625, 0.340], [0.665, 0.465], [0.480, 0.520], [0.410, 0.680],
                [0.375, 0.850], [0.545, 0.520], [0.600, 0.700], [0.635, 0.880], [0.520, 0.098],
                [0.542, 0.098], [0.508, 0.104], [0.554, 0.104]
            ]
        },
        armsup: {
            label: 'Arms up',
            joints: [
                [0.500, 0.105], [0.500, 0.185], [0.440, 0.200], [0.405, 0.095], [0.385, 0.010],
                [0.560, 0.200], [0.595, 0.095], [0.615, 0.010], [0.470, 0.520], [0.462, 0.700],
                [0.458, 0.880], [0.530, 0.520], [0.538, 0.700], [0.542, 0.880], [0.489, 0.094],
                [0.511, 0.094], [0.478, 0.100], [0.522, 0.100]
            ]
        },
        running: {
            label: 'Running',
            joints: [
                [0.545, 0.115], [0.520, 0.195], [0.470, 0.205], [0.545, 0.255], [0.610, 0.190],
                [0.575, 0.210], [0.640, 0.350], [0.700, 0.465], [0.490, 0.520], [0.560, 0.650],
                [0.615, 0.800], [0.550, 0.530], [0.420, 0.640], [0.330, 0.700], [0.535, 0.102],
                [0.557, 0.102], [0.523, 0.108], [0.570, 0.108]
            ]
        }
    };

    const STORAGE_KEY = 'howl-pose-joints';
    const STORAGE_SIZE = 'howl-pose-size';

    let joints = null;
    let canvas = null;
    let ctx = null;
    let dragging = -1;
    let outputSize = 768;
    let hover = -1;

    function clonePreset(key) {
        const preset = PRESETS[key] || PRESETS.standing;
        return preset.joints.map(p => [p[0], p[1]]);
    }

    function loadSaved() {
        try {
            const raw = localStorage.getItem(STORAGE_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            if (!Array.isArray(parsed) || parsed.length !== JOINT_COUNT) return null;
            if (!parsed.every(p => Array.isArray(p) && p.length === 2 && isFinite(p[0]) && isFinite(p[1]))) return null;
            return parsed.map(p => [Number(p[0]), Number(p[1])]);
        } catch (_) {
            return null;
        }
    }

    function persist() {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(joints));
            localStorage.setItem(STORAGE_SIZE, String(outputSize));
        } catch (_) { /* storage may be full or disabled; the pose still works */ }
    }

    /** Draw the skeleton at any pixel size. */
    function draw(target, size, pose, opts) {
        const g = target.getContext('2d');
        const options = opts || {};
        g.save();
        g.fillStyle = '#000000';
        g.fillRect(0, 0, size, size);
        g.lineCap = 'round';
        g.lineJoin = 'round';

        const unit = size / 512;
        const limbWidth = Math.max(3, Math.round(7 * unit));
        const dotRadius = Math.max(3, Math.round(6 * unit));

        LIMBS.forEach(([a, b, color]) => {
            const pa = pose[a], pb = pose[b];
            if (!pa || !pb) return;
            g.strokeStyle = 'rgb(' + color.join(',') + ')';
            g.lineWidth = limbWidth;
            g.beginPath();
            g.moveTo(pa[0] * size, pa[1] * size);
            g.lineTo(pb[0] * size, pb[1] * size);
            g.stroke();
        });

        // Joints: OpenPose draws small dots so the model can locate keypoints.
        pose.forEach((p, i) => {
            if (!p) return;
            const x = p[0] * size;
            const y = p[1] * size;
            g.beginPath();
            g.arc(x, y, (options.highlight === i ? dotRadius * 1.6 : dotRadius), 0, Math.PI * 2);
            g.fillStyle = options.highlight === i ? '#ffffff' : 'rgb(255,255,255)';
            g.fill();
        });
        g.restore();
    }

    function previewSize() {
        return canvas ? canvas.width / (window.devicePixelRatio || 1) : 420;
    }

    function repaint() {
        if (!ctx) return;
        const size = previewSize();
        // Draw into an offscreen buffer at the preview size, then scale onto the
        // hi-DPI canvas so lines stay crisp.
        const buf = document.createElement('canvas');
        buf.width = Math.round(size);
        buf.height = Math.round(size);
        draw(buf, buf.width, joints, { highlight: hover >= 0 ? hover : dragging });
        ctx.setTransform(1, 0, 0, 1, 0, 0);
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.drawImage(buf, 0, 0, canvas.width, canvas.height);
    }

    function eventToUnit(e) {
        const rect = canvas.getBoundingClientRect();
        const x = (e.clientX - rect.left) / rect.width;
        const y = (e.clientY - rect.top) / rect.height;
        return [Math.min(1, Math.max(0, x)), Math.min(1, Math.max(0, y))];
    }

    function nearestJoint(unitX, unitY) {
        const threshold = 0.055;
        let best = -1;
        let bestDist = Infinity;
        joints.forEach((p, i) => {
            const dx = p[0] - unitX;
            const dy = p[1] - unitY;
            const d = Math.sqrt(dx * dx + dy * dy);
            if (d < bestDist && d < threshold) { bestDist = d; best = i; }
        });
        return best;
    }

    function setJointLabel(i) {
        const el = document.getElementById('ig-pose-selection');
        if (!el) return;
        const names = ['nose', 'neck', 'right shoulder', 'right elbow', 'right wrist',
            'left shoulder', 'left elbow', 'left wrist', 'right hip', 'right knee',
            'right ankle', 'left hip', 'left knee', 'left ankle', 'right eye', 'left eye',
            'right ear', 'left ear'];
        el.textContent = i >= 0 ? 'Moving ' + names[i] : 'Drag a joint to pose';
    }

    function renderToDataURL(size) {
        const buf = document.createElement('canvas');
        buf.width = size;
        buf.height = size;
        draw(buf, size, joints, {});
        return buf.toDataURL('image/png');
    }

    function sendToControlNet() {
        const size = outputSize;
        const dataUrl = renderToDataURL(size);
        window.igPoseJoints = joints.map(p => [p[0], p[1]]);
        window.igPoseSize = size;
        window.controlNetImageData = dataUrl;
        persist();

        const zone = document.getElementById('ig-upload-controlnet');
        if (zone) {
            zone.innerHTML = '<img alt="Pose control image" src="' + dataUrl +
                '" style="width:100%;height:100%;object-fit:contain;border-radius:6px;background:#000;">';
            zone.dataset.filled = '1';
        }
        const status = document.getElementById('ig-controlnet-status');
        if (status) status.textContent = 'Pose control image ready · ' + size + ' × ' + size;

        close();
        if (typeof setStatus === 'function') {
            setStatus('Pose control image ready — generate to apply it.');
        }
    }

    /** Public: restore a pose from a workflow and rebuild the control image. */
    function applyPose(savedJoints, size) {
        if (!Array.isArray(savedJoints) || savedJoints.length !== JOINT_COUNT) return false;
        joints = savedJoints.map(p => [Number(p[0]) || 0, Number(p[1]) || 0]);
        if (size) outputSize = Number(size) || outputSize;
        const sizeSelect = document.getElementById('ig-pose-size');
        if (sizeSelect) sizeSelect.value = String(outputSize);
        persist();
        const dataUrl = renderToDataURL(outputSize);
        window.igPoseJoints = joints.map(p => [p[0], p[1]]);
        window.igPoseSize = outputSize;
        window.controlNetImageData = dataUrl;
        const zone = document.getElementById('ig-upload-controlnet');
        if (zone) {
            zone.innerHTML = '<img alt="Pose control image" src="' + dataUrl +
                '" style="width:100%;height:100%;object-fit:contain;border-radius:6px;background:#000;">';
            zone.dataset.filled = '1';
        }
        return true;
    }

    function buildModal() {
        if (document.getElementById('ig-pose-overlay')) return;
        const overlay = document.createElement('div');
        overlay.id = 'ig-pose-overlay';
        overlay.className = 'ig-pose-overlay';
        overlay.innerHTML = [
            '<div class="ig-pose-dialog" role="dialog" aria-modal="true" aria-labelledby="ig-pose-title">',
            '  <header class="ig-pose-header">',
            '    <div>',
            '      <h2 id="ig-pose-title">Pose Editor</h2>',
            '      <p>OpenPose skeleton for ControlNet. Drag any joint.</p>',
            '    </div>',
            '    <button type="button" class="ig-btn ig-btn-secondary" id="ig-pose-close" aria-label="Close pose editor">Close</button>',
            '  </header>',
            '  <div class="ig-pose-body">',
            '    <div class="ig-pose-stage"><canvas id="ig-pose-canvas" width="420" height="420"></canvas></div>',
            '    <aside class="ig-pose-side">',
            '      <label class="ig-label" for="ig-pose-preset">Preset</label>',
            '      <select id="ig-pose-preset" class="ig-select">',
            Object.keys(PRESETS).map(k => '<option value="' + k + '">' + PRESETS[k].label + '</option>').join(''),
            '      </select>',
            '      <label class="ig-label" for="ig-pose-size">Control image size</label>',
            '      <select id="ig-pose-size" class="ig-select">',
            '        <option value="512">512 × 512</option>',
            '        <option value="768" selected>768 × 768</option>',
            '        <option value="1024">1024 × 1024</option>',
            '      </select>',
            '      <div class="ig-pose-actions">',
            '        <button type="button" class="ig-btn ig-btn-secondary" id="ig-pose-flip">Flip</button>',
            '        <button type="button" class="ig-btn ig-btn-secondary" id="ig-pose-reset">Reset</button>',
            '      </div>',
            '      <p class="ig-pose-hint" id="ig-pose-selection">Drag a joint to pose</p>',
            '      <p class="ig-pose-hint">OpenPose ControlNet reads this exact skeleton. Use an SDXL OpenPose ControlNet with an SDXL checkpoint.</p>',
            '    </aside>',
            '  </div>',
            '  <footer class="ig-pose-footer">',
            '    <span class="ig-pose-hint">The figure is drawn on black, as the model expects.</span>',
            '    <button type="button" class="ig-btn ig-generate" id="ig-pose-apply">Use as Control Image</button>',
            '  </footer>',
            '</div>'
        ].join('\n');
        document.body.appendChild(overlay);

        overlay.addEventListener('click', e => {
            if (e.target === overlay) close();
        });
        document.addEventListener('keydown', e => {
            if (e.key === 'Escape' && overlay.classList.contains('open')) close();
        });

        document.getElementById('ig-pose-close').addEventListener('click', close);
        document.getElementById('ig-pose-apply').addEventListener('click', sendToControlNet);

        document.getElementById('ig-pose-preset').addEventListener('change', e => {
            joints = clonePreset(e.target.value);
            dragging = -1;
            setJointLabel(-1);
            persist();
            repaint();
        });

        document.getElementById('ig-pose-size').addEventListener('change', e => {
            outputSize = parseInt(e.target.value, 10) || 768;
            persist();
        });

        document.getElementById('ig-pose-flip').addEventListener('click', () => {
            joints = joints.map(([x, y]) => [1 - x, y]);
            persist();
            repaint();
        });

        document.getElementById('ig-pose-reset').addEventListener('click', () => {
            joints = clonePreset(document.getElementById('ig-pose-preset').value);
            persist();
            repaint();
        });

        canvas = document.getElementById('ig-pose-canvas');
        ctx = canvas.getContext('2d');

        const onDown = e => {
            const [ux, uy] = eventToUnit(e);
            dragging = nearestJoint(ux, uy);
            if (dragging >= 0) {
                canvas.setPointerCapture(e.pointerId);
                setJointLabel(dragging);
                repaint();
            }
        };
        const onMove = e => {
            const [ux, uy] = eventToUnit(e);
            if (dragging >= 0) {
                joints[dragging] = [ux, uy];
                repaint();
                return;
            }
            const near = nearestJoint(ux, uy);
            if (near !== hover) {
                hover = near;
                canvas.style.cursor = near >= 0 ? 'grab' : 'default';
                repaint();
            }
        };
        const onUp = e => {
            if (dragging >= 0) {
                persist();
                if (canvas.hasPointerCapture(e.pointerId)) canvas.releasePointerCapture(e.pointerId);
                dragging = -1;
                repaint();
            }
        };

        canvas.addEventListener('pointerdown', onDown);
        canvas.addEventListener('pointermove', onMove);
        canvas.addEventListener('pointerup', onUp);
        canvas.addEventListener('pointercancel', onUp);
    }

    function sizeCanvas() {
        if (!canvas) return;
        const dpr = window.devicePixelRatio || 1;
        const target = canvas.clientWidth || 420;
        canvas.width = Math.round(target * dpr);
        canvas.height = Math.round(target * dpr);
        repaint();
    }

    function open() {
        buildModal();
        if (!joints) joints = loadSaved() || clonePreset('standing');
        try {
            const savedSize = parseInt(localStorage.getItem(STORAGE_SIZE) || '', 10);
            if (savedSize >= 256 && savedSize <= 2048) outputSize = savedSize;
        } catch (_) { /* ignore */ }
        const select = document.getElementById('ig-pose-size');
        if (select) select.value = String(outputSize);
        const overlay = document.getElementById('ig-pose-overlay');
        overlay.classList.add('open');
        // Size synchronously. requestAnimationFrame is throttled in background
        // tabs, which would leave the canvas blank until the tab is focused.
        sizeCanvas();
        setJointLabel(-1);
        requestAnimationFrame(() => {
            if (overlay.classList.contains('open')) sizeCanvas();
        });
    }

    function close() {
        const overlay = document.getElementById('ig-pose-overlay');
        if (overlay) overlay.classList.remove('open');
    }

    window.igOpenPoseEditor = open;
    window.igApplyPose = applyPose;

    function bind() {
        const btn = document.getElementById('ig-open-pose-editor');
        if (btn) btn.addEventListener('click', open);

        const clear = document.getElementById('ig-clear-controlnet');
        if (clear) {
            clear.addEventListener('click', () => {
                window.controlNetImageData = null;
                window.igPoseJoints = null;
                const zone = document.getElementById('ig-upload-controlnet');
                if (zone) {
                    zone.dataset.filled = '';
                    zone.innerHTML = '<span>Drop control image or click</span>';
                }
                const status = document.getElementById('ig-controlnet-status');
                if (status) status.textContent = '';
            });
        }
        window.addEventListener('resize', () => {
            const overlay = document.getElementById('ig-pose-overlay');
            if (overlay && overlay.classList.contains('open')) sizeCanvas();
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', bind);
    } else {
        bind();
    }
})();
