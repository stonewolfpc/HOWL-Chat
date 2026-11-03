'use strict';

// ---------------------------------------------------------------------------
// State Persistence
// ---------------------------------------------------------------------------
function saveState() {
    const state = {
        _version: 2,
        model:     document.getElementById('ig-model').value,
        vae:       document.getElementById('ig-vae-model').value,
        clipL:     document.getElementById('ig-clip-l').value,
        clipG:     document.getElementById('ig-clip-g').value,
        t5xxl:     document.getElementById('ig-t5xxl').value,
        taesd:     document.getElementById('ig-taesd').checked,
        prompt:    document.getElementById('ig-prompt').value,
        negative:  document.getElementById('ig-negative').value,
        width:     document.getElementById('ig-width').value,
        height:    document.getElementById('ig-height').value,
        steps:     document.getElementById('ig-steps').value,
        seed:      document.getElementById('ig-seed').value,
        sampler:   document.getElementById('ig-sampler').value,
        scheduler: document.getElementById('ig-scheduler').value,
        cfg:       document.getElementById('ig-cfg').value,
        guidance:  document.getElementById('ig-guidance').value,
        slgScale:  document.getElementById('ig-slg-scale').value,
        eta:       document.getElementById('ig-eta').value,
        batch:     document.getElementById('ig-batch').value,
        clipSkip:  document.getElementById('ig-clip-skip').value,
        // LoRA: new strength-based system (array of {name, strength, enabled})
        loras:     window.activeLoRAs || [],
        hires:     document.getElementById('ig-hires-enable').checked,
        hrScale:   document.getElementById('ig-hr-scale')?.value || '2',
        hrUpscaler: document.getElementById('ig-hr-upscaler')?.value || '',
        strength:  document.getElementById('ig-strength')?.value || '0.75',
        strengthInpaint: document.getElementById('ig-strength-inpaint')?.value || '0.75',
        outpaintStrength: document.getElementById('ig-outpaint-strength')?.value || '0.8',
        outpaintMargin: document.getElementById('ig-outpaint-margin')?.value || '128',
        outpaintDirection: document.querySelector('.ig-dir-btn.active')?.dataset?.dir || 'all',
        upscaleFactor: document.getElementById('ig-upscale-factor')?.value || '2',
        upscaleDenoise: document.getElementById('ig-upscale-denoise')?.value || '0.4',
        resizeMode: document.getElementById('ig-resize-mode')?.value || 'just_resize',
        styleStrength: document.getElementById('ig-style-strength')?.value || '0.5',
        preserveColors: document.getElementById('ig-preserve-colors')?.checked || false,
        matchHistogram: document.getElementById('ig-match-histogram')?.checked || false,
        latentNoise: document.getElementById('ig-latent-noise')?.value || '0',
        maskBlur: document.getElementById('ig-mask-blur')?.value || '4',
        maskInvert: document.getElementById('ig-mask-invert')?.checked || false,
        // Inpaint-specific
        brushHardness: document.getElementById('ig-brush-hardness')?.value || '100',
        maskBlurInpaint: document.getElementById('ig-mask-blur-inpaint')?.value || '4',
        maskFeather: document.getElementById('ig-mask-feather')?.value || '0',
        maskedContent: document.getElementById('ig-masked-content')?.value || 'fill',
        resizeModeInpaint: document.getElementById('ig-resize-mode-inpaint')?.value || 'just_resize',
        styleStrengthInpaint: document.getElementById('ig-style-strength-inpaint')?.value || '0.5',
        preserveColorsInpaint: document.getElementById('ig-preserve-colors-inpaint')?.checked || false,
        matchHistogramInpaint: document.getElementById('ig-match-histogram-inpaint')?.checked || false,
        latentNoiseInpaint: document.getElementById('ig-latent-noise-inpaint')?.value || '0',
        // Outpaint-specific
        edgeBehavior: document.getElementById('ig-edge-behavior')?.value || 'mirror',
        resizeModeOutpaint: document.getElementById('ig-resize-mode-outpaint')?.value || 'just_resize',
        styleStrengthOutpaint: document.getElementById('ig-style-strength-outpaint')?.value || '0.5',
        preserveColorsOutpaint: document.getElementById('ig-preserve-colors-outpaint')?.checked || false,
        matchHistogramOutpaint: document.getElementById('ig-match-histogram-outpaint')?.checked || false,
        latentNoiseOutpaint: document.getElementById('ig-latent-noise-outpaint')?.value || '0',
        // Upscale-specific
        upscaleMethod: document.getElementById('ig-upscale-method')?.value || 'latent_sd',
        upscaleSharpen: document.getElementById('ig-upscale-sharpen')?.value || 'medium',
        upscaleAntiArtifact: document.getElementById('ig-upscale-anti-artifact')?.value || 'off',
        upscalePrompt: document.getElementById('ig-upscale-prompt')?.value || '',
        upscaleNegative: document.getElementById('ig-upscale-negative')?.value || '',
        upscaleSampler: document.getElementById('ig-upscale-sampler')?.value || 'euler_a',
        upscaleSteps: document.getElementById('ig-upscale-steps')?.value || '20',
        upscaleCFG: document.getElementById('ig-upscale-cfg')?.value || '7',
        upscaleSeed: document.getElementById('ig-upscale-seed')?.value || '-1',
        ctrlModel: document.getElementById('ig-control-model').value,
        ctrlStrength: document.getElementById('ig-control-strength').value,
        videoFrames: document.getElementById('ig-video-frames')?.value,
        videoFPS:    document.getElementById('ig-video-fps')?.value,
        prediction:  document.getElementById('ig-prediction')?.value,
        rng:         document.getElementById('ig-rng')?.value,
        samplerRng:  document.getElementById('ig-sampler-rng')?.value,
        flowShift:   document.getElementById('ig-flow-shift')?.value,
        vaeTileOverlap: document.getElementById('ig-vae-tile-overlap')?.value,
        disableMetadata: document.getElementById('ig-disable-metadata')?.checked,
        skipLayerStart: document.getElementById('ig-skip-layer-start')?.value,
        skipLayerEnd:   document.getElementById('ig-skip-layer-end')?.value,
        highNoiseSteps: document.getElementById('ig-high-noise-steps')?.value,
        highNoiseCfg:   document.getElementById('ig-high-noise-cfg')?.value,
        highNoiseGuidance: document.getElementById('ig-high-noise-guidance')?.value,
        highNoiseSlg:   document.getElementById('ig-high-noise-slg')?.value,
        highNoiseEta:   document.getElementById('ig-high-noise-eta')?.value,
        highNoiseSkipStart: document.getElementById('ig-high-noise-skip-start')?.value,
        highNoiseSkipEnd:   document.getElementById('ig-high-noise-skip-end')?.value,
        backend:      document.getElementById('ig-backend')?.value,
        paramsBackend: document.getElementById('ig-params-backend')?.value,
        vaeTiling:    document.getElementById('ig-vae-tiling')?.checked,
        temporalTiling: document.getElementById('ig-temporal-tiling')?.checked,
        circular:     document.getElementById('ig-circular')?.checked,
        flashAttention: document.getElementById('ig-fa')?.checked,
        offloadToCPU: document.getElementById('ig-offload-cpu')?.checked,
        // Video-specific
        videoFormat: document.getElementById('ig-video-format')?.value || 'mp4',
        videoSeedMode: document.getElementById('ig-video-seed-mode')?.value || 'fixed',
        videoMotion: document.getElementById('ig-video-motion')?.value || '0.5',
        videoDenoise: document.getElementById('ig-video-denoise')?.value || '0.6',
        videoLoop: document.getElementById('ig-video-loop')?.checked || false,
        // Global expert
        cfgRescale: document.getElementById('ig-cfg-rescale')?.value || '0',
        faceRestore: document.getElementById('ig-face-restore')?.value || '0',
        sChurn: document.getElementById('ig-s-churn')?.value || '0',
        sTmin: document.getElementById('ig-s-tmin')?.value || '0',
        sTmax: document.getElementById('ig-s-tmax')?.value || '50',
        sNoise: document.getElementById('ig-s-noise')?.value || '1',
        // UI state
        currentMode: typeof currentMode !== 'undefined' ? currentMode : 'txt2img',
        seedLocked: typeof seedLocked !== 'undefined' ? seedLocked : false,
        aspectLocked: typeof _aspectLocked !== 'undefined' ? _aspectLocked : true,
        // Note: lastUsedSeed is NOT persisted to avoid reusing bad seeds across reloads
    };
    localStorage.setItem(IG_STATE_KEY, JSON.stringify(state));
}

function restoreState() {
    const raw = localStorage.getItem(IG_STATE_KEY);
    if (!raw) return;
    try {
        const s = JSON.parse(raw);
        if (s.prompt)    document.getElementById('ig-prompt').value    = s.prompt;
        if (s.negative)  document.getElementById('ig-negative').value  = s.negative;
        if (s.model)     document.getElementById('ig-model').value     = s.model;
        if (s.width)     document.getElementById('ig-width').value    = s.width;
        if (s.height)    document.getElementById('ig-height').value   = s.height;
        if (s.steps)     document.getElementById('ig-steps').value     = s.steps;
        if (s.seed)      document.getElementById('ig-seed').value      = s.seed;
        if (s.sampler)   document.getElementById('ig-sampler').value    = s.sampler;
        if (s.scheduler) document.getElementById('ig-scheduler').value  = s.scheduler;
        if (s.cfg)       document.getElementById('ig-cfg').value       = s.cfg;
        if (s.guidance)  document.getElementById('ig-guidance').value  = s.guidance;
        if (s.slgScale)  document.getElementById('ig-slg-scale').value  = s.slgScale;
        if (s.eta)       document.getElementById('ig-eta').value       = s.eta;
        if (s.batch)     document.getElementById('ig-batch').value       = s.batch;
        if (s.clipSkip)  document.getElementById('ig-clip-skip').value   = s.clipSkip;
        if (s.vae)       document.getElementById('ig-vae-model').value = s.vae;
        if (s.clipL)     document.getElementById('ig-clip-l').value    = s.clipL;
        if (s.clipG)     document.getElementById('ig-clip-g').value    = s.clipG;
        if (s.t5xxl)     document.getElementById('ig-t5xxl').value    = s.t5xxl;
        if (s.taesd !== undefined) document.getElementById('ig-taesd').checked = s.taesd;
        // LoRA: restore new format (loras) with backward compatibility (lora)
        if (s.loras && Array.isArray(s.loras)) {
            window.activeLoRAs = s.loras;
            renderActiveLoRAs();
        } else if (s.lora && Array.isArray(s.lora)) {
            // Migrate old format: array of names -> array of {name, strength: 1.0, enabled: true}
            window.activeLoRAs = s.lora.map(name => ({ name, strength: 1.0, enabled: true }));
            renderActiveLoRAs();
        }
        if (s.hires !== undefined)  document.getElementById('ig-hires-enable').checked = s.hires;
        if (s.hrScale !== undefined && document.getElementById('ig-hr-scale')) document.getElementById('ig-hr-scale').value = s.hrScale;
        if (s.hrUpscaler !== undefined && document.getElementById('ig-hr-upscaler')) document.getElementById('ig-hr-upscaler').value = s.hrUpscaler;
        if (s.strength !== undefined && document.getElementById('ig-strength')) document.getElementById('ig-strength').value = s.strength;
        if (s.strengthInpaint !== undefined && document.getElementById('ig-strength-inpaint')) document.getElementById('ig-strength-inpaint').value = s.strengthInpaint;
        if (s.outpaintStrength !== undefined && document.getElementById('ig-outpaint-strength')) document.getElementById('ig-outpaint-strength').value = s.outpaintStrength;
        if (s.outpaintMargin !== undefined && document.getElementById('ig-outpaint-margin')) document.getElementById('ig-outpaint-margin').value = s.outpaintMargin;
        if (s.outpaintDirection && document.querySelectorAll('.ig-dir-btn')) {
            document.querySelectorAll('.ig-dir-btn').forEach(btn => { btn.classList.toggle('active', btn.dataset.dir === s.outpaintDirection); });
        }
        if (s.upscaleFactor !== undefined && document.getElementById('ig-upscale-factor')) document.getElementById('ig-upscale-factor').value = s.upscaleFactor;
        if (s.upscaleDenoise !== undefined && document.getElementById('ig-upscale-denoise')) document.getElementById('ig-upscale-denoise').value = s.upscaleDenoise;
        if (s.resizeMode !== undefined && document.getElementById('ig-resize-mode')) document.getElementById('ig-resize-mode').value = s.resizeMode;
        if (s.styleStrength !== undefined && document.getElementById('ig-style-strength')) document.getElementById('ig-style-strength').value = s.styleStrength;
        if (s.preserveColors !== undefined && document.getElementById('ig-preserve-colors')) document.getElementById('ig-preserve-colors').checked = s.preserveColors;
        if (s.matchHistogram !== undefined && document.getElementById('ig-match-histogram')) document.getElementById('ig-match-histogram').checked = s.matchHistogram;
        if (s.latentNoise !== undefined && document.getElementById('ig-latent-noise')) document.getElementById('ig-latent-noise').value = s.latentNoise;
        if (s.maskBlur !== undefined && document.getElementById('ig-mask-blur')) document.getElementById('ig-mask-blur').value = s.maskBlur;
        if (s.maskInvert !== undefined && document.getElementById('ig-mask-invert')) document.getElementById('ig-mask-invert').checked = s.maskInvert;
        // Inpaint restore
        if (s.brushHardness !== undefined && document.getElementById('ig-brush-hardness')) document.getElementById('ig-brush-hardness').value = s.brushHardness;
        if (s.maskBlurInpaint !== undefined && document.getElementById('ig-mask-blur-inpaint')) document.getElementById('ig-mask-blur-inpaint').value = s.maskBlurInpaint;
        if (s.maskFeather !== undefined && document.getElementById('ig-mask-feather')) document.getElementById('ig-mask-feather').value = s.maskFeather;
        if (s.maskedContent !== undefined && document.getElementById('ig-masked-content')) document.getElementById('ig-masked-content').value = s.maskedContent;
        if (s.resizeModeInpaint !== undefined && document.getElementById('ig-resize-mode-inpaint')) document.getElementById('ig-resize-mode-inpaint').value = s.resizeModeInpaint;
        if (s.styleStrengthInpaint !== undefined && document.getElementById('ig-style-strength-inpaint')) document.getElementById('ig-style-strength-inpaint').value = s.styleStrengthInpaint;
        if (s.preserveColorsInpaint !== undefined && document.getElementById('ig-preserve-colors-inpaint')) document.getElementById('ig-preserve-colors-inpaint').checked = s.preserveColorsInpaint;
        if (s.matchHistogramInpaint !== undefined && document.getElementById('ig-match-histogram-inpaint')) document.getElementById('ig-match-histogram-inpaint').checked = s.matchHistogramInpaint;
        if (s.latentNoiseInpaint !== undefined && document.getElementById('ig-latent-noise-inpaint')) document.getElementById('ig-latent-noise-inpaint').value = s.latentNoiseInpaint;
        // Outpaint restore
        if (s.edgeBehavior !== undefined && document.getElementById('ig-edge-behavior')) document.getElementById('ig-edge-behavior').value = s.edgeBehavior;
        if (s.resizeModeOutpaint !== undefined && document.getElementById('ig-resize-mode-outpaint')) document.getElementById('ig-resize-mode-outpaint').value = s.resizeModeOutpaint;
        if (s.styleStrengthOutpaint !== undefined && document.getElementById('ig-style-strength-outpaint')) document.getElementById('ig-style-strength-outpaint').value = s.styleStrengthOutpaint;
        if (s.preserveColorsOutpaint !== undefined && document.getElementById('ig-preserve-colors-outpaint')) document.getElementById('ig-preserve-colors-outpaint').checked = s.preserveColorsOutpaint;
        if (s.matchHistogramOutpaint !== undefined && document.getElementById('ig-match-histogram-outpaint')) document.getElementById('ig-match-histogram-outpaint').checked = s.matchHistogramOutpaint;
        if (s.latentNoiseOutpaint !== undefined && document.getElementById('ig-latent-noise-outpaint')) document.getElementById('ig-latent-noise-outpaint').value = s.latentNoiseOutpaint;
        // Upscale restore
        if (s.upscaleMethod !== undefined && document.getElementById('ig-upscale-method')) document.getElementById('ig-upscale-method').value = s.upscaleMethod;
        if (s.upscaleSharpen !== undefined && document.getElementById('ig-upscale-sharpen')) document.getElementById('ig-upscale-sharpen').value = s.upscaleSharpen;
        if (s.upscaleAntiArtifact !== undefined && document.getElementById('ig-upscale-anti-artifact')) document.getElementById('ig-upscale-anti-artifact').value = s.upscaleAntiArtifact;
        if (s.upscalePrompt !== undefined && document.getElementById('ig-upscale-prompt')) document.getElementById('ig-upscale-prompt').value = s.upscalePrompt;
        if (s.upscaleNegative !== undefined && document.getElementById('ig-upscale-negative')) document.getElementById('ig-upscale-negative').value = s.upscaleNegative;
        if (s.upscaleSampler !== undefined && document.getElementById('ig-upscale-sampler')) document.getElementById('ig-upscale-sampler').value = s.upscaleSampler;
        if (s.upscaleSteps !== undefined && document.getElementById('ig-upscale-steps')) document.getElementById('ig-upscale-steps').value = s.upscaleSteps;
        if (s.upscaleCFG !== undefined && document.getElementById('ig-upscale-cfg')) document.getElementById('ig-upscale-cfg').value = s.upscaleCFG;
        if (s.upscaleSeed !== undefined && document.getElementById('ig-upscale-seed')) document.getElementById('ig-upscale-seed').value = s.upscaleSeed;
        if (s.ctrlModel) document.getElementById('ig-control-model').value = s.ctrlModel;
        if (s.ctrlStrength) document.getElementById('ig-control-strength').value = s.ctrlStrength;
        if (s.videoFrames && document.getElementById('ig-video-frames')) document.getElementById('ig-video-frames').value = s.videoFrames;
        if (s.videoFPS && document.getElementById('ig-video-fps')) document.getElementById('ig-video-fps').value = s.videoFPS;
        if (s.prediction !== undefined && document.getElementById('ig-prediction')) document.getElementById('ig-prediction').value = s.prediction;
        if (s.rng !== undefined && document.getElementById('ig-rng')) document.getElementById('ig-rng').value = s.rng;
        if (s.samplerRng !== undefined && document.getElementById('ig-sampler-rng')) document.getElementById('ig-sampler-rng').value = s.samplerRng;
        if (s.flowShift !== undefined && document.getElementById('ig-flow-shift')) document.getElementById('ig-flow-shift').value = s.flowShift;
        if (s.vaeTileOverlap !== undefined && document.getElementById('ig-vae-tile-overlap')) document.getElementById('ig-vae-tile-overlap').value = s.vaeTileOverlap;
        if (s.disableMetadata !== undefined && document.getElementById('ig-disable-metadata')) document.getElementById('ig-disable-metadata').checked = s.disableMetadata;
        if (s.skipLayerStart !== undefined && document.getElementById('ig-skip-layer-start')) document.getElementById('ig-skip-layer-start').value = s.skipLayerStart;
        if (s.skipLayerEnd !== undefined && document.getElementById('ig-skip-layer-end')) document.getElementById('ig-skip-layer-end').value = s.skipLayerEnd;
        if (s.highNoiseSteps !== undefined && document.getElementById('ig-high-noise-steps')) document.getElementById('ig-high-noise-steps').value = s.highNoiseSteps;
        if (s.highNoiseCfg !== undefined && document.getElementById('ig-high-noise-cfg')) document.getElementById('ig-high-noise-cfg').value = s.highNoiseCfg;
        if (s.highNoiseGuidance !== undefined && document.getElementById('ig-high-noise-guidance')) document.getElementById('ig-high-noise-guidance').value = s.highNoiseGuidance;
        if (s.highNoiseSlg !== undefined && document.getElementById('ig-high-noise-slg')) document.getElementById('ig-high-noise-slg').value = s.highNoiseSlg;
        if (s.highNoiseEta !== undefined && document.getElementById('ig-high-noise-eta')) document.getElementById('ig-high-noise-eta').value = s.highNoiseEta;
        if (s.highNoiseSkipStart !== undefined && document.getElementById('ig-high-noise-skip-start')) document.getElementById('ig-high-noise-skip-start').value = s.highNoiseSkipStart;
        if (s.highNoiseSkipEnd !== undefined && document.getElementById('ig-high-noise-skip-end')) document.getElementById('ig-high-noise-skip-end').value = s.highNoiseSkipEnd;
        if (s.backend !== undefined && document.getElementById('ig-backend')) document.getElementById('ig-backend').value = s.backend;
        if (s.paramsBackend !== undefined && document.getElementById('ig-params-backend')) document.getElementById('ig-params-backend').value = s.paramsBackend;
        if (s.vaeTiling !== undefined && document.getElementById('ig-vae-tiling')) document.getElementById('ig-vae-tiling').checked = s.vaeTiling;
        if (s.temporalTiling !== undefined && document.getElementById('ig-temporal-tiling')) document.getElementById('ig-temporal-tiling').checked = s.temporalTiling;
        if (s.circular !== undefined && document.getElementById('ig-circular')) document.getElementById('ig-circular').checked = s.circular;
        if (s.flashAttention !== undefined && document.getElementById('ig-fa')) document.getElementById('ig-fa').checked = s.flashAttention;
        if (s.offloadToCPU !== undefined && document.getElementById('ig-offload-cpu')) document.getElementById('ig-offload-cpu').checked = s.offloadToCPU;
        // Video restore
        if (s.videoFormat !== undefined && document.getElementById('ig-video-format')) document.getElementById('ig-video-format').value = s.videoFormat;
        if (s.videoSeedMode !== undefined && document.getElementById('ig-video-seed-mode')) document.getElementById('ig-video-seed-mode').value = s.videoSeedMode;
        if (s.videoMotion !== undefined && document.getElementById('ig-video-motion')) document.getElementById('ig-video-motion').value = s.videoMotion;
        if (s.videoDenoise !== undefined && document.getElementById('ig-video-denoise')) document.getElementById('ig-video-denoise').value = s.videoDenoise;
        if (s.videoLoop !== undefined && document.getElementById('ig-video-loop')) document.getElementById('ig-video-loop').checked = s.videoLoop;
        // Global expert restore
        if (s.cfgRescale !== undefined && document.getElementById('ig-cfg-rescale')) document.getElementById('ig-cfg-rescale').value = s.cfgRescale;
        if (s.faceRestore !== undefined && document.getElementById('ig-face-restore')) document.getElementById('ig-face-restore').value = s.faceRestore;
        if (s.sChurn !== undefined && document.getElementById('ig-s-churn')) document.getElementById('ig-s-churn').value = s.sChurn;
        if (s.sTmin !== undefined && document.getElementById('ig-s-tmin')) document.getElementById('ig-s-tmin').value = s.sTmin;
        if (s.sTmax !== undefined && document.getElementById('ig-s-tmax')) document.getElementById('ig-s-tmax').value = s.sTmax;
        if (s.sNoise !== undefined && document.getElementById('ig-s-noise')) document.getElementById('ig-s-noise').value = s.sNoise;

        // Restore UI state
        if (s.currentMode && typeof setWorkflowMode === 'function') {
            setWorkflowMode(s.currentMode);
        }
        if (s.seedLocked !== undefined) {
            seedLocked = s.seedLocked;
            const lockBtn = document.getElementById('ig-seed-lock');
            if (lockBtn) lockBtn.classList.toggle('ig-locked', seedLocked);
        }
        if (s.aspectLocked !== undefined && typeof _aspectLocked !== 'undefined') {
            _aspectLocked = s.aspectLocked;
            const aspectBtn = document.getElementById('ig-aspect-lock');
            if (aspectBtn) {
                aspectBtn.classList.toggle('active', _aspectLocked);
                aspectBtn.innerHTML = _aspectLocked ? '&#x1F512; Lock' : '&#x1F513; Free';
            }
        }
        // Note: lastUsedSeed is NOT restored to avoid reusing bad seeds across reloads

        updateLabels();
    } catch (e) { /* ignore */ }
}

function updateLabels() {
    const ids = ['ig-steps', 'ig-cfg', 'ig-strength', 'ig-strength-inpaint', 'ig-outpaint-strength',
                 'ig-upscale-denoise', 'ig-guidance', 'ig-slg-scale', 'ig-eta',
                 'ig-control-strength', 'ig-video-frames', 'ig-video-fps',
                 'ig-style-strength', 'ig-latent-noise', 'ig-mask-blur',
                 'ig-brush-hardness', 'ig-mask-blur-inpaint', 'ig-mask-feather',
                 'ig-style-strength-inpaint', 'ig-latent-noise-inpaint',
                 'ig-style-strength-outpaint', 'ig-latent-noise-outpaint',
                 'ig-flow-shift', 'ig-skip-layer-start', 'ig-skip-layer-end',
                 'ig-vae-tile-overlap', 'ig-high-noise-steps', 'ig-high-noise-cfg',
                 'ig-high-noise-guidance', 'ig-high-noise-slg', 'ig-high-noise-eta',
                 'ig-high-noise-skip-start', 'ig-high-noise-skip-end',
                 'ig-upscale-steps', 'ig-upscale-cfg',
                 'ig-video-motion', 'ig-video-denoise',
                 'ig-cfg-rescale', 'ig-face-restore',
                 'ig-s-churn', 'ig-s-tmin', 'ig-s-tmax', 'ig-s-noise'];
    ids.forEach(id => {
        const slider = document.getElementById(id);
        const valEl = document.getElementById(id + '-val');
        if (slider && valEl) {
            let val = slider.value;
            if (id === 'ig-flow-shift' && parseFloat(val) === 0) val = 'Auto';
            if (id === 'ig-high-noise-steps' && parseInt(val, 10) === -1) val = 'Auto';
            valEl.textContent = val;
        }
    });
}

// ---------------------------------------------------------------------------
// Event Binding
// ---------------------------------------------------------------------------
function bindEvents() {
    // Sliders
    const sliderIds = ['ig-steps', 'ig-cfg', 'ig-strength', 'ig-strength-inpaint', 'ig-outpaint-strength',
                       'ig-upscale-denoise', 'ig-guidance', 'ig-slg-scale', 'ig-eta',
                       'ig-control-strength', 'ig-video-frames', 'ig-video-fps', 'ig-brush-size',
                       'ig-style-strength', 'ig-latent-noise', 'ig-mask-blur',
                       'ig-brush-hardness', 'ig-mask-blur-inpaint', 'ig-mask-feather',
                       'ig-style-strength-inpaint', 'ig-latent-noise-inpaint',
                       'ig-style-strength-outpaint', 'ig-latent-noise-outpaint',
                       'ig-flow-shift', 'ig-skip-layer-start', 'ig-skip-layer-end',
                       'ig-vae-tile-overlap', 'ig-high-noise-steps', 'ig-high-noise-cfg',
                       'ig-high-noise-guidance', 'ig-high-noise-slg', 'ig-high-noise-eta',
                       'ig-high-noise-skip-start', 'ig-high-noise-skip-end',
                       'ig-upscale-steps', 'ig-upscale-cfg',
                       'ig-video-motion', 'ig-video-denoise',
                       'ig-cfg-rescale', 'ig-face-restore',
                       'ig-s-churn', 'ig-s-tmin', 'ig-s-tmax', 'ig-s-noise'];
    sliderIds.forEach(id => {
        const slider = document.getElementById(id);
        if (slider) {
            slider.addEventListener('input', () => {
                const valEl = document.getElementById(id + '-val');
                if (valEl) {
                    let val = slider.value;
                    if (id === 'ig-flow-shift' && parseFloat(val) === 0) val = 'Auto';
                    if (id === 'ig-high-noise-steps' && parseInt(val, 10) === -1) val = 'Auto';
                    valEl.textContent = val;
                }
                saveState();
            });
        }
    });

    // ── Aspect ratio & resolution controls ────────────────────────────────────
    let _aspectLocked = true;  // lock on by default

    const _lockBtn     = document.getElementById('ig-aspect-lock');
    const _optimizeBtn = document.getElementById('ig-optimize-res');
    const _wInput      = document.getElementById('ig-width');
    const _hInput      = document.getElementById('ig-height');

    function _snap64(n) { return Math.max(64, Math.round(n / 64) * 64); }

    function _currentRatio() {
        const w = parseInt(_wInput.value, 10) || 512;
        const h = parseInt(_hInput.value, 10) || 512;
        return w / h;
    }

    function _syncAspectButtons() {
        const w = parseInt(_wInput.value, 10);
        const h = parseInt(_hInput.value, 10);
        document.querySelectorAll('.ig-aspect-btn').forEach(b => {
            b.classList.toggle('active', parseInt(b.dataset.w) === w && parseInt(b.dataset.h) === h);
        });
    }

    function _setDimensions(w, h) {
        _wInput.value = w;
        _hInput.value = h;
        _syncAspectButtons();
        saveState();
    }

    // Aspect ratio buttons — set exact preset dimensions
    document.querySelectorAll('.ig-aspect-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            _setDimensions(parseInt(btn.dataset.w), parseInt(btn.dataset.h));
        });
    });

    // Lock toggle
    if (_lockBtn) {
        _lockBtn.addEventListener('click', () => {
            _aspectLocked = !_aspectLocked;
            _lockBtn.classList.toggle('active', _aspectLocked);
            _lockBtn.innerHTML = _aspectLocked ? '&#x1F512; Lock' : '&#x1F513; Free';
        });
    }

    // Width input — adjust height to preserve ratio when locked
    if (_wInput) {
        _wInput.addEventListener('change', () => {
            if (_aspectLocked) {
                const ratio = _currentRatio(); // ratio before this change was applied...
                // ratio was computed from current (new) w / old h — recalculate from active btn
                const activeBtn = document.querySelector('.ig-aspect-btn.active');
                const r = activeBtn
                    ? parseInt(activeBtn.dataset.w) / parseInt(activeBtn.dataset.h)
                    : _currentRatio();
                const newW = parseInt(_wInput.value, 10);
                _hInput.value = _snap64(newW / r);
            }
            _syncAspectButtons();
            saveState();
        });
    }

    // Height input — adjust width to preserve ratio when locked
    if (_hInput) {
        _hInput.addEventListener('change', () => {
            if (_aspectLocked) {
                const activeBtn = document.querySelector('.ig-aspect-btn.active');
                const r = activeBtn
                    ? parseInt(activeBtn.dataset.w) / parseInt(activeBtn.dataset.h)
                    : _currentRatio();
                const newH = parseInt(_hInput.value, 10);
                _wInput.value = _snap64(newH * r);
            }
            _syncAspectButtons();
            saveState();
        });
    }

    // Optimize — apply model-aware resolution based on optimal ranges and VRAM constraints
    if (_optimizeBtn) {
        _optimizeBtn.addEventListener('click', () => {
            const family = (window.currentModelFamily || '').toLowerCase();
            
            // Model-specific optimal resolution ranges (min, max, base)
            const modelRanges = {
                'flux':     { min: 512,  max: 2048, base: 1024 },
                'sdxl':     { min: 512,  max: 2048, base: 1024 },
                'sd3':      { min: 512,  max: 2048, base: 1024 },
                'sd3.5':    { min: 512,  max: 2048, base: 1024 },
                'chroma':   { min: 512,  max: 2048, base: 1024 },
                'kolors':   { min: 512,  max: 2048, base: 1024 },
                'auraflow': { min: 512,  max: 2048, base: 1024 },
                'hidream':  { min: 512,  max: 2048, base: 1024 },
                'sd2.1':    { min: 384,  max: 1024, base: 768 },
                'sd21':     { min: 384,  max: 1024, base: 768 },
                'cascade':  { min: 384,  max: 1024, base: 768 },
                'sdxl-turbo': { min: 256, max: 1024, base: 512 },
                'turbo':    { min: 256, max: 1024, base: 512 },
                'default':  { min: 256, max: 1024, base: 512 } // SD1.5 and unknown
            };
            
            const range = modelRanges[family] || modelRanges['default'];
            
            // Try to get VRAM info to adjust for available memory
            let vramGB = 8; // Conservative default
            if (window.go?.main?.App?.GetGPUInfo) {
                try {
                    const gpuInfo = window.go.main.App.GetGPUInfo();
                    if (gpuInfo && gpuInfo.TotalMemory) {
                        vramGB = Math.floor(gpuInfo.TotalMemory / 1024); // Convert MB to GB
                    }
                } catch (e) {
                    console.warn('Failed to get GPU info, using default VRAM');
                }
            }
            
            // Adjust max resolution based on VRAM
            const vramAdjustments = {
                4:  0.5,  // 4GB VRAM - very conservative
                6:  0.6,  // 6GB VRAM
                8:  0.75, // 8GB VRAM
                12: 1.0, // 12GB VRAM - full capability
                16: 1.2, // 16GB VRAM - can go higher
                24: 1.5  // 24GB+ VRAM - high end
            };
            
            const adjustment = vramAdjustments[vramGB] || 0.75;
            const adjustedMax = Math.floor(range.max * adjustment);
            const adjustedBase = Math.floor(range.base * Math.min(1, adjustment));
            
            // Preserve current aspect ratio
            const activeBtn = document.querySelector('.ig-aspect-btn.active');
            const r = activeBtn
                ? parseInt(activeBtn.dataset.w) / parseInt(activeBtn.dataset.h)
                : _currentRatio();
            
            // Calculate optimal resolution within adjusted bounds
            const budget = adjustedBase * adjustedBase;
            const w = _snap64(Math.sqrt(budget * r));
            const h = _snap64(w / r);
            
            // Ensure within min/max bounds
            const finalW = Math.max(range.min, Math.min(adjustedMax, w));
            const finalH = Math.max(range.min, Math.min(adjustedMax, h));
            
            _setDimensions(finalW, finalH);
            
            // Also apply optimal sampler, scheduler, steps, cfg, guidance for this family
            if (typeof applyModelPreset === 'function') {
                applyModelPreset(family);
            }
            
            console.log(`[Optimize] Model: ${family}, VRAM: ${vramGB}GB, Resolution: ${finalW}x${finalH}`);
        });
    }

    // Auto-expand textareas (JS fallback for browsers without field-sizing)
    function autoExpand(textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = Math.min(textarea.scrollHeight, 200) + 'px';
    }
    const promptTa = document.getElementById('ig-prompt');
    const negativeTa = document.getElementById('ig-negative');
    if (promptTa) {
        promptTa.addEventListener('input', () => autoExpand(promptTa));
        // Initial expand
        autoExpand(promptTa);
    }
    if (negativeTa) {
        negativeTa.addEventListener('input', () => autoExpand(negativeTa));
        autoExpand(negativeTa);
    }

    // Inputs that trigger save
    ['ig-model', 'ig-vae-model', 'ig-clip-l', 'ig-clip-g', 'ig-t5xxl', 'ig-taesd',
     'ig-prompt', 'ig-negative', 'ig-width', 'ig-height', 'ig-seed',
     'ig-sampler', 'ig-scheduler', 'ig-batch', 'ig-clip-skip', 'ig-lora-selector',
     'ig-hires-enable', 'ig-control-model',
     'ig-prediction', 'ig-rng', 'ig-sampler-rng', 'ig-disable-metadata',
     'ig-backend', 'ig-params-backend', 'ig-vae-tiling', 'ig-temporal-tiling',
     'ig-circular', 'ig-fa', 'ig-offload-cpu',
     'ig-outpaint-margin', 'ig-outpaint-strength', 'ig-upscale-factor', 'ig-upscale-denoise',
     'ig-hr-scale', 'ig-hr-upscaler',
     'ig-resize-mode', 'ig-preserve-colors', 'ig-match-histogram', 'ig-mask-invert',
     'ig-masked-content', 'ig-resize-mode-inpaint', 'ig-preserve-colors-inpaint', 'ig-match-histogram-inpaint',
     'ig-edge-behavior', 'ig-resize-mode-outpaint', 'ig-preserve-colors-outpaint', 'ig-match-histogram-outpaint',
     'ig-upscale-method', 'ig-upscale-sharpen', 'ig-upscale-anti-artifact',
     'ig-upscale-prompt', 'ig-upscale-negative', 'ig-upscale-sampler', 'ig-upscale-seed',
     'ig-video-frames', 'ig-video-fps', 'ig-video-format', 'ig-video-seed-mode', 'ig-video-loop',
     'ig-cfg-rescale', 'ig-face-restore',
     'ig-s-churn', 'ig-s-tmin', 'ig-s-tmax', 'ig-s-noise'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener('change', saveState);
    });

    // Upscale factor change recalculates dimensions and target display
    const upscaleFactorSel = document.getElementById('ig-upscale-factor');
    if (upscaleFactorSel) {
        upscaleFactorSel.addEventListener('change', () => {
            if (currentMode === 'upscale' && initImageData) {
                const img = document.getElementById('ig-init-image');
                if (img) {
                    const factor = parseInt(upscaleFactorSel.value || '2', 10);
                    const w = Math.round(img.naturalWidth * factor / 64) * 64;
                    const h = Math.round(img.naturalHeight * factor / 64) * 64;
                    document.getElementById('ig-width').value = w;
                    document.getElementById('ig-height').value = h;
                    updateUpscaleTarget(img.naturalWidth, img.naturalHeight, factor);
                    saveState();
                }
            }
        });
    }

    // File upload for img2img
    const uploadInput = document.getElementById('ig-upload-input');
    const uploadZone = document.getElementById('ig-upload-zone');
    if (uploadInput && uploadZone) {
        uploadZone.addEventListener('click', () => {
            uploadInput.dataset.context = 'ig-upload-zone';
            uploadInput.click();
        });
        uploadInput.addEventListener('change', (e) => {
            if (e.target.files?.[0]) {
                const context = uploadInput.dataset.context;
                if (context === 'ig-upload-controlnet') {
                    handleControlNetUpload(e.target.files[0]);
                } else if (context === 'ig-upload-mini-video-end') {
                    handleVideoEndUpload(e.target.files[0]);
                } else if (context === 'ig-upload-mini-img2img-mask') {
                    handleImg2imgMaskUpload(e.target.files[0]);
                } else {
                    handleImageUpload(e.target.files[0]);
                }
            }
        });
        // Drag and drop
        uploadZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadZone.classList.add('dragover');
        });
        uploadZone.addEventListener('dragleave', () => uploadZone.classList.remove('dragover'));
        uploadZone.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadZone.classList.remove('dragover');
            if (e.dataTransfer.files?.[0]) handleImageUpload(e.dataTransfer.files[0]);
        });
    }

    // Mini upload zones in Mode Tools panel
    document.querySelectorAll('.ig-upload-mini').forEach(zone => {
        zone.addEventListener('click', () => {
            if (uploadInput) {
                uploadInput.dataset.context = zone.id;
                uploadInput.click();
            }
        });
    });

    // Brush controls for inpaint (now in left panel Mode Tools)
    const brushSize = document.getElementById('ig-brush-size');
    if (brushSize) {
        brushSize.addEventListener('input', () => {
            const valEl = document.getElementById('ig-brush-val');
            if (valEl) valEl.textContent = brushSize.value;
        });
    }
    const clearMaskBtn = document.getElementById('ig-clear-mask');
    if (clearMaskBtn) clearMaskBtn.addEventListener('click', clearMask);
    initMaskDrawing();

    // Outpaint direction buttons
    document.querySelectorAll('.ig-dir-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.ig-dir-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            saveState();
        });
    });

    // Generate
    document.getElementById('ig-generate').addEventListener('click', onGenerate);

    // Model dropdown change (updates presets, no backend load)
    const modelSel = document.getElementById('ig-model');
    if (modelSel) modelSel.addEventListener('change', onModelChange);

    // Unload Model button
    const unloadBtn = document.getElementById('ig-unload-model');
    if (unloadBtn) unloadBtn.addEventListener('click', onUnloadModel);

    // Reload Model button
    const reloadBtn = document.getElementById('ig-reload-model');
    if (reloadBtn) reloadBtn.addEventListener('click', onReloadModel);

    // Panel toggles
    document.querySelectorAll('.ig-panel-header').forEach(header => {
        header.addEventListener('click', () => togglePanel(header.closest('.ig-panel')));
    });

    // Column resizers
    document.querySelectorAll('.ig-resizer').forEach(resizer => {
        resizer.addEventListener('mousedown', e => {
            const targetAttr = resizer.dataset.target;
            let prev;
            if (targetAttr === 'sidebar') {
                prev = document.getElementById('ig-sidebar');
            } else if (targetAttr === 'history') {
                prev = resizer.previousElementSibling;
                while (prev && !prev.classList.contains('ig-main')) prev = prev.previousElementSibling;
                if (!prev) return;
                const historyCol = document.querySelector('.ig-history-col');
                if (!historyCol) return;
                const startX = e.clientX;
                const startWidth = historyCol.offsetWidth;

                function onMove(ev) {
                    const delta = startX - ev.clientX;
                    const newW = Math.min(Math.max(startWidth + delta, 180), 600);
                    historyCol.style.width = newW + 'px';
                    historyCol.style.flexShrink = '0';
                }
                function onUp() {
                    window.removeEventListener('mousemove', onMove);
                    window.removeEventListener('mouseup', onUp);
                }
                window.addEventListener('mousemove', onMove);
                window.addEventListener('mouseup', onUp);
                return;
            }
            if (!prev) return;
            const startX = e.clientX;
            const startWidth = prev.offsetWidth;

            function onMove(ev) {
                const delta = ev.clientX - startX;
                const newW = Math.min(Math.max(startWidth + delta, 200), 600);
                prev.style.width = newW + 'px';
                prev.style.flexShrink = '0';
            }
            function onUp() {
                window.removeEventListener('mousemove', onMove);
                window.removeEventListener('mouseup', onUp);
            }
            window.addEventListener('mousemove', onMove);
            window.addEventListener('mouseup', onUp);
        });
    });

    // Workflow mode buttons
    document.querySelectorAll('.ig-workflow-btn').forEach(btn => {
        btn.addEventListener('click', () => setWorkflowMode(btn.dataset.mode));
    });

    // Seed lock
    const lockBtn = document.getElementById('ig-seed-lock');
    if (lockBtn) lockBtn.addEventListener('click', toggleSeedLock);

    // Seed reuse
    const reuseBtn = document.getElementById('ig-seed-reuse');
    if (reuseBtn) reuseBtn.addEventListener('click', reuseLastSeed);

    // img2img Clear button
    const clearImg2Img = document.getElementById('ig-img2img-clear');
    if (clearImg2Img) clearImg2Img.addEventListener('click', clearInitImage);

    // img2img Use Metadata button
    const useMetaBtn = document.getElementById('ig-img2img-use-meta');
    if (useMetaBtn) {
        useMetaBtn.addEventListener('click', async () => {
            await handleUseMetadata('img2img');
        });
    }

    // Inpaint Clear button
    const clearInpaint = document.getElementById('ig-inpaint-clear');
    if (clearInpaint) clearInpaint.addEventListener('click', clearInitImage);

    // Inpaint Use Metadata button
    const useMetaInpaint = document.getElementById('ig-inpaint-use-meta');
    if (useMetaInpaint) {
        useMetaInpaint.addEventListener('click', async () => {
            await handleUseMetadata('inpaint');
        });
    }

    // Outpaint Clear button
    const clearOutpaint = document.getElementById('ig-outpaint-clear');
    if (clearOutpaint) clearOutpaint.addEventListener('click', clearInitImage);

    // Outpaint Use Metadata button
    const useMetaOutpaint = document.getElementById('ig-outpaint-use-meta');
    if (useMetaOutpaint) {
        useMetaOutpaint.addEventListener('click', async () => {
            await handleUseMetadata('outpaint');
        });
    }

    // Inpaint mask tool buttons
    const eraserToggle = document.getElementById('ig-eraser-toggle');
    if (eraserToggle) {
        eraserToggle.addEventListener('click', () => {
            eraserToggle.classList.toggle('active');
            setStatus(eraserToggle.classList.contains('active') ? 'Eraser active.' : 'Brush active.');
        });
    }
    const previewToggle = document.getElementById('ig-mask-preview-toggle');
    if (previewToggle) {
        previewToggle.addEventListener('click', () => {
            previewToggle.classList.toggle('active');
            setStatus(previewToggle.classList.contains('active') ? 'Mask preview on.' : 'Mask preview off.');
        });
    }
    const invertMaskBtn = document.getElementById('ig-invert-mask');
    if (invertMaskBtn) {
        invertMaskBtn.addEventListener('click', () => {
            setStatus('Mask invert toggled.');
        });
    }

    // Upscale method change shows/hides conditional controls
    const upscaleMethodSel = document.getElementById('ig-upscale-method');
    if (upscaleMethodSel) {
        function updateUpscaleVisibility() {
            const method = upscaleMethodSel.value;
            const isLatent = method === 'latent_sd' || method === 'latent_sdxl' || method === 'sdxl_refiner';
            const isEsrgan = method === 'esrgan' || method === 'realesrgan' || method === 'swinir' ||
                             method === 'ultrasharp' || method === 'animesharp' || method === 'remacri';
            document.querySelectorAll('.ig-upscale-latent-only').forEach(el => {
                el.style.display = isLatent ? '' : 'none';
            });
            document.querySelectorAll('.ig-upscale-esrgan-only').forEach(el => {
                el.style.display = isEsrgan ? '' : 'none';
            });
        }
        upscaleMethodSel.addEventListener('change', () => {
            updateUpscaleVisibility();
            saveState();
        });
        // Initial visibility
        updateUpscaleVisibility();
    }

    // Upscale Clear button
    const clearUpscale = document.getElementById('ig-upscale-clear');
    if (clearUpscale) clearUpscale.addEventListener('click', clearInitImage);

    // Upscale Use Metadata button
    const useMetaUpscale = document.getElementById('ig-upscale-use-meta');
    if (useMetaUpscale) {
        useMetaUpscale.addEventListener('click', async () => {
            await handleUseMetadata('upscale');
        });
    }

    // Upscale random seed button
    const upscaleRandom = document.getElementById('ig-upscale-random');
    if (upscaleRandom) {
        upscaleRandom.addEventListener('click', () => {
            const seed = Math.floor(Math.random() * 2147483647);
            const seedInput = document.getElementById('ig-upscale-seed');
            if (seedInput) seedInput.value = seed;
            saveState();
        });
    }

    // Video Clear button
    const clearVideo = document.getElementById('ig-video-clear');
    if (clearVideo) clearVideo.addEventListener('click', clearInitImage);

    // Video Use Metadata button
    const useMetaVideo = document.getElementById('ig-video-use-meta');
    if (useMetaVideo) {
        useMetaVideo.addEventListener('click', async () => {
            await handleUseMetadata('video');
        });
    }

    // LoRA selector: add LoRA to active list
    const loraSelector = document.getElementById('ig-lora-selector');
    if (loraSelector) {
        loraSelector.addEventListener('change', () => {
            const name = loraSelector.value;
            if (!name) return;
            // Prevent duplicates
            if (!window.activeLoRAs.find(l => l.name === name)) {
                window.activeLoRAs.push({ name, strength: 1.0, enabled: true });
                renderActiveLoRAs();
                saveState();
            }
            loraSelector.value = ''; // Reset to placeholder
        });
    }

    // Initialize activeLoRAs array if not present
    if (!window.activeLoRAs) window.activeLoRAs = [];

    // Keyboard shortcuts
    document.addEventListener('keydown', handleKeydown);
}

// ---------------------------------------------------------------------------
// LoRA Active Panel Rendering
// ---------------------------------------------------------------------------
function renderActiveLoRAs() {
    const panel = document.getElementById('ig-lora-active-panel');
    if (!panel) return;

    if (!window.activeLoRAs || window.activeLoRAs.length === 0) {
        panel.innerHTML = '<div class="ig-lora-empty">No LoRAs selected. Choose from dropdown above.</div>';
        return;
    }

    panel.innerHTML = window.activeLoRAs.map((lora, index) => `
        <div class="ig-lora-item ${lora.enabled ? '' : 'disabled'}" data-index="${index}">
            <div class="ig-lora-header">
                <span class="ig-lora-name" title="${lora.name}">${lora.name}</span>
                <label class="ig-lora-toggle">
                    <input type="checkbox" ${lora.enabled ? 'checked' : ''} onchange="toggleLoRAEnabled(${index})">
                    <span>On</span>
                </label>
                <button class="ig-lora-remove" onclick="removeLoRA(${index})" title="Remove">&#x2715;</button>
            </div>
            <div class="ig-lora-strength">
                <span class="ig-lora-strength-label">Strength:</span>
                <input type="range" class="ig-slider" min="-2" max="2" step="0.05" value="${lora.strength}"
                       oninput="updateLoRAStrength(${index}, this.value)">
                <span class="ig-lora-strength-val">${lora.strength.toFixed(2)}</span>
            </div>
        </div>
    `).join('');
}

function toggleLoRAEnabled(index) {
    if (window.activeLoRAs[index]) {
        window.activeLoRAs[index].enabled = !window.activeLoRAs[index].enabled;
        renderActiveLoRAs();
        saveState();
    }
}

function removeLoRA(index) {
    window.activeLoRAs.splice(index, 1);
    renderActiveLoRAs();
    saveState();
}

function updateLoRAStrength(index, value) {
    if (window.activeLoRAs[index]) {
        window.activeLoRAs[index].strength = parseFloat(value);
        // Update display value
        const panel = document.getElementById('ig-lora-active-panel');
        const item = panel?.querySelector(`[data-index="${index}"]`);
        const valEl = item?.querySelector('.ig-lora-strength-val');
        if (valEl) valEl.textContent = parseFloat(value).toFixed(2);
        saveState();
    }
}

// ---------------------------------------------------------------------------
// Metadata Extraction
// ---------------------------------------------------------------------------
async function handleUseMetadata(mode) {
    if (!initImageData) {
        setStatus('Please upload an image first to extract metadata.');
        return;
    }

    setStatus('Extracting metadata from image...');

    try {
        const meta = await window.go.main.App.ExtractImageMetadataFromBase64(initImageData);

        if (!meta || (!meta.prompt && !meta.steps && !meta.seed)) {
            setStatus('No generation metadata found in this image.');
            return;
        }

        // Populate fields if metadata exists
        if (meta.prompt) {
            const promptEl = document.getElementById('ig-prompt');
            if (promptEl) promptEl.value = meta.prompt;
        }

        if (meta.negative_prompt) {
            const negEl = document.getElementById('ig-negative');
            if (negEl) negEl.value = meta.negative_prompt;
        }

        if (meta.seed !== undefined && meta.seed !== 0) {
            const seedEl = document.getElementById('ig-seed');
            if (seedEl) seedEl.value = meta.seed;
        }

        if (meta.steps && meta.steps > 0) {
            const stepsEl = document.getElementById('ig-steps');
            if (stepsEl) {
                stepsEl.value = meta.steps;
                updateLabels();
            }
        }

        if (meta.cfg_scale && meta.cfg_scale > 0) {
            const cfgEl = document.getElementById('ig-cfg');
            if (cfgEl) {
                cfgEl.value = meta.cfg_scale;
                updateLabels();
            }
        }

        if (meta.sampler) {
            const samplerEl = document.getElementById('ig-sampler');
            if (samplerEl) {
                // Try to find matching sampler (handle variations)
                const samplerLower = meta.sampler.toLowerCase();
                const options = Array.from(samplerEl.options);
                const match = options.find(o => o.value.toLowerCase() === samplerLower ||
                    o.textContent.toLowerCase().includes(samplerLower));
                if (match) samplerEl.value = match.value;
            }
        }

        if (meta.scheduler) {
            const schedEl = document.getElementById('ig-scheduler');
            if (schedEl) {
                const schedLower = meta.scheduler.toLowerCase();
                const options = Array.from(schedEl.options);
                const match = options.find(o => o.value.toLowerCase() === schedLower);
                if (match) schedEl.value = match.value;
            }
        }

        if (meta.width && meta.height && meta.width > 0 && meta.height > 0) {
            const widthEl = document.getElementById('ig-width');
            const heightEl = document.getElementById('ig-height');
            if (widthEl) widthEl.value = meta.width;
            if (heightEl) heightEl.value = meta.height;
        }

        // Mode-specific fields
        if (mode === 'img2img') {
            // Note: We don't auto-set denoise strength from metadata
            // as it depends on the workflow context
        }

        saveState();
        setStatus(`Metadata extracted: ${meta.prompt ? meta.prompt.substring(0, 30) + '...' : 'no prompt'}, Seed: ${meta.seed || 'N/A'}, Steps: ${meta.steps || 'N/A'}`);

        // Auto-expand textareas if content is long
        const promptTa = document.getElementById('ig-prompt');
        const negTa = document.getElementById('ig-negative');
        if (promptTa && window.autoExpand) window.autoExpand(promptTa);
        if (negTa && window.autoExpand) window.autoExpand(negTa);

    } catch (err) {
        console.error('Metadata extraction failed:', err);
        setStatus('Failed to extract metadata: ' + (err.message || err));
    }
}

// ---------------------------------------------------------------------------
// Panel Toggles
// ---------------------------------------------------------------------------
function togglePanel(panel) {
    if (!panel) return;
    const body = panel.querySelector('.ig-panel-body');
    const toggle = panel.querySelector('.ig-panel-toggle');
    if (!body) return;
    const isCollapsed = panel.classList.contains('ig-panel-collapsed');
    if (isCollapsed) {
        panel.classList.remove('ig-panel-collapsed');
        body.style.display = '';
        if (toggle) toggle.innerHTML = '&#x25BC;';
    } else {
        panel.classList.add('ig-panel-collapsed');
        body.style.display = 'none';
        if (toggle) toggle.innerHTML = '&#x25B6;';
    }
}

// ---------------------------------------------------------------------------
// Workflow Modes
// ---------------------------------------------------------------------------
let currentMode = 'txt2img';
let initImageData = null; // base64 string of uploaded image for img2img

function setWorkflowMode(mode) {
    currentMode = mode;
    document.querySelectorAll('.ig-workflow-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.mode === mode);
    });

    const isImg2Img = mode === 'img2img';
    const isInpaint = mode === 'inpaint';
    const isTxt2Img = mode === 'txt2img';
    const isOutpaint = mode === 'outpaint';
    const isUpscale = mode === 'upscale';
    const isVideo = mode === 'video';
    const needsImage = isImg2Img || isInpaint || isOutpaint || isUpscale;

    // Mode Tools panel — show correct section
    document.querySelectorAll('.ig-mode-section').forEach(el => el.style.display = 'none');
    const activeSection = document.getElementById('ig-mode-' + mode);
    if (activeSection) activeSection.style.display = '';

    // Mode-scoped options.
    //
    // ig-controls.css has always declared body.mode-txt2img .ig-txt2img-only,
    // body.mode-inpaint .ig-inpaint-only and friends, but nothing ever added the
    // body class, so every mode showed every control: options that cannot work
    // in the current mode looked available and then did nothing.
    document.body.classList.remove('mode-txt2img', 'mode-img2img', 'mode-inpaint',
                                   'mode-outpaint', 'mode-upscale', 'mode-video');
    document.body.classList.add('mode-' + mode);

    // Whole panels that only make sense in some modes. Anything not listed here
    // is common to every mode (model, sub-models, performance, DiT settings).
    const PANEL_MODES = {
        'Prompt': ['txt2img', 'img2img', 'inpaint', 'outpaint', 'video'],
        'Parameters': ['txt2img', 'img2img', 'inpaint', 'outpaint', 'video'],
        'Sampler Options': ['txt2img', 'img2img', 'inpaint', 'outpaint', 'video'],
        'Sampler Shaping (DPM++/SDE)': ['txt2img', 'img2img', 'inpaint', 'outpaint', 'video'],
        'Presets': ['txt2img', 'img2img', 'inpaint', 'outpaint', 'video'],
        'ControlNet': ['txt2img', 'img2img', 'inpaint', 'outpaint'],
        'High-Noise Fix Settings': ['txt2img']
    };
    document.querySelectorAll('.ig-panel').forEach(panel => {
        const titleEl = panel.querySelector('.ig-panel-title');
        if (!titleEl) return;
        const allowed = PANEL_MODES[titleEl.textContent.trim()];
        if (!allowed) return;
        panel.style.display = allowed.includes(mode) ? '' : 'none';
    });

    // A few individual controls live inside shared panels but only mean
    // something in some modes (batch count produces one image per run, so it is
    // meaningless for upscale and video). Hide the whole field row, not just the
    // input, so no orphaned label is left behind.
    const FIELD_MODES = {
        'ig-batch': ['txt2img', 'img2img', 'inpaint', 'outpaint']
    };
    Object.keys(FIELD_MODES).forEach(id => {
        const el = document.getElementById(id);
        if (!el) return;
        const row = el.closest('.ig-field, .ig-checkbox, label') || el;
        row.style.display = FIELD_MODES[id].includes(mode) ? '' : 'none';
    });

    // The model list itself is mode-scoped: video families only make sense in
    // Video mode, and still-image checkpoints cannot render video.
    if (typeof _populateDiffusionDropdown === 'function' && typeof models !== 'undefined' && models.length) {
        const modelEl = document.getElementById('ig-model');
        const previous = modelEl ? modelEl.value : '';
        _populateDiffusionDropdown();
        if (modelEl && previous && [...modelEl.options].some(o => o.value === previous)) {
            modelEl.value = previous;
        }
    }

    // Canvas placeholder / upload zone
    const placeholder = document.getElementById('ig-canvas-placeholder');
    const uploadZone = document.getElementById('ig-upload-zone');
    const previewArea = document.getElementById('ig-preview-area');
    const maskCanvas = document.getElementById('ig-mask-canvas');

    if (placeholder) {
        if (needsImage) {
            placeholder.style.display = 'none';
            if (uploadZone) uploadZone.style.display = 'flex';
        } else {
            placeholder.style.display = '';
            if (uploadZone) uploadZone.style.display = 'none';
            if (!isVideo) {
                initImageData = null;
                clearMask();
            }
        }
        // Reset mini upload zones (preserve ControlNet always; preserve video end only when entering video)
        document.querySelectorAll('.ig-upload-mini').forEach(z => {
            if (z.id === 'ig-upload-controlnet') return;
            if (isVideo && z.id === 'ig-upload-mini-video-end') return;
            z.classList.remove('ig-uploaded');
            const span = z.querySelector('span');
            if (span) {
                const defaultText = z.id === 'ig-upload-mini-video-end'
                    ? 'Drop end frame or click to browse'
                    : 'Drop image or click to browse';
                span.textContent = defaultText;
            }
        });
        if (!isVideo) {
            videoEndImageData = null;
            window.videoEndImageData = null;
        }
    }

    if (maskCanvas) {
        maskCanvas.style.display = isInpaint ? 'block' : 'none';
        maskCanvas.classList.toggle('active', isInpaint);
        if (isInpaint) resizeMaskCanvas();
    }

    if ((isTxt2Img || isVideo) && currentResult) {
        if (currentResult.images?.length > 1) {
            let gridHtml = '<div class="ig-batch-grid">';
            currentResult.images.forEach((b64, idx) => {
                gridHtml += `<img src="data:image/png;base64,${b64}" alt="Generated ${idx+1}" class="ig-batch-thumb" data-idx="${idx}" onclick="window.viewBatchImage(${idx})">`;
            });
            gridHtml += '</div>';
            previewArea.innerHTML = gridHtml;
        } else if (currentResult.image_data) {
            previewArea.innerHTML = `<img src="data:image/png;base64,${currentResult.image_data}" alt="Generated">`;
        }
    } else if (!needsImage && !initImageData) {
        previewArea.innerHTML = '';
        if (placeholder) placeholder.style.display = '';
    }

    // Upscale mode: show target resolution if image already loaded
    const targetEl = document.getElementById('ig-upscale-target');
    if (isUpscale && initImageData) {
        const img = document.getElementById('ig-init-image');
        if (img && targetEl) {
            const factor = parseInt(document.getElementById('ig-upscale-factor')?.value || '2', 10);
            updateUpscaleTarget(img.naturalWidth || parseInt(document.getElementById('ig-width').value, 10) || 512,
                                img.naturalHeight || parseInt(document.getElementById('ig-height').value, 10) || 512,
                                factor);
        }
    } else if (targetEl) {
        targetEl.style.display = 'none';
    }

    // Update labels for new mode sliders
    updateLabels();
}

// --- Image-to-Image Mask Handler ---
let img2imgMaskData = null;
window.img2imgMaskData = img2imgMaskData;
function handleImg2imgMaskUpload(file) {
    if (!file.type.startsWith('image/')) {
        setStatus('Please upload an image file for the mask.');
        return;
    }
    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        const commaIdx = dataUrl.indexOf(',');
        img2imgMaskData = commaIdx >= 0 ? dataUrl.slice(commaIdx + 1) : dataUrl;
        window.img2imgMaskData = img2imgMaskData;
        const zone = document.getElementById('ig-upload-mini-img2img-mask');
        if (zone) {
            zone.classList.add('ig-uploaded');
            const span = zone.querySelector('span');
            if (span) span.textContent = 'Mask loaded — click to replace';
        }
        setStatus('Mask loaded for img2img.');
    };
    reader.readAsDataURL(file);
}

// --- Video End Image Handler ---
let videoEndImageData = null;
window.videoEndImageData = videoEndImageData;
function handleVideoEndUpload(file) {
    if (!file.type.startsWith('image/')) {
        setStatus('Please upload an image file for the end frame.');
        return;
    }
    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        const commaIdx = dataUrl.indexOf(',');
        videoEndImageData = commaIdx >= 0 ? dataUrl.slice(commaIdx + 1) : dataUrl;
        window.videoEndImageData = videoEndImageData;
        const zone = document.getElementById('ig-upload-mini-video-end');
        if (zone) {
            zone.classList.add('ig-uploaded');
            const span = zone.querySelector('span');
            if (span) span.textContent = 'End frame loaded — click to replace';
        }
        setStatus('Video end frame loaded.');
    };
    reader.readAsDataURL(file);
}

// --- ControlNet Upload Handler ---
let controlNetImageData = null;
function handleControlNetUpload(file) {
    if (!file.type.startsWith('image/')) {
        setStatus('Please upload an image file for ControlNet.');
        return;
    }
    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        const commaIdx = dataUrl.indexOf(',');
        controlNetImageData = commaIdx >= 0 ? dataUrl.slice(commaIdx + 1) : dataUrl;

        const zone = document.getElementById('ig-upload-controlnet');
        if (zone) {
            zone.classList.add('ig-uploaded');
            zone.querySelector('span').textContent = 'Control image loaded';
        }
        setStatus('ControlNet image ready.');
    };
    reader.readAsDataURL(file);
}

// --- File Upload Handler ---
function handleImageUpload(file) {
    if (!file.type.startsWith('image/')) {
        setStatus('Please upload an image file (PNG, JPG, WEBP).');
        return;
    }
    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        const commaIdx = dataUrl.indexOf(',');
        initImageData = commaIdx >= 0 ? dataUrl.slice(commaIdx + 1) : dataUrl;

        const previewArea = document.getElementById('ig-preview-area');
        const uploadZone = document.getElementById('ig-upload-zone');
        previewArea.innerHTML = `<img id="ig-init-image" src="${dataUrl}" alt="Init Image" style="max-width:100%;max-height:100%;object-fit:contain;border-radius:8px;">`;
        if (uploadZone) uploadZone.style.display = 'none';

        // Mark the specific mini zone that triggered the upload
        const uploadInput = document.getElementById('ig-upload-input');
        const context = uploadInput?.dataset?.context;
        const targetZone = context ? document.getElementById(context) : null;
        document.querySelectorAll('.ig-upload-mini').forEach(z => {
            if (z.id === 'ig-upload-controlnet') return;
            z.classList.add('ig-uploaded');
            const span = z.querySelector('span');
            if (span) span.textContent = 'Image loaded — click to replace';
        });
        if (targetZone) {
            targetZone.classList.add('ig-uploaded');
            const span = targetZone.querySelector('span');
            if (span) span.textContent = 'Image loaded — click to replace';
        }

        // For upscale mode, calculate and show target dimensions
        if (currentMode === 'upscale') {
            const img = new Image();
            img.onload = () => {
                const factor = parseInt(document.getElementById('ig-upscale-factor')?.value || '2', 10);
                const w = Math.round(img.naturalWidth * factor / 64) * 64;
                const h = Math.round(img.naturalHeight * factor / 64) * 64;
                document.getElementById('ig-width').value = w;
                document.getElementById('ig-height').value = h;
                updateUpscaleTarget(img.naturalWidth, img.naturalHeight, factor);
                saveState();
            };
            img.src = dataUrl;
        }

        // Size mask canvas to match the new image
        requestAnimationFrame(() => {
            resizeMaskCanvas();
            clearMask();
        });

        setStatus('Image loaded. Adjust strength and click Generate.');
    };
    reader.readAsDataURL(file);
}

function clearInitImage() {
    initImageData = null;
    const previewArea = document.getElementById('ig-preview-area');
    const placeholder = document.getElementById('ig-canvas-placeholder');
    const uploadZone = document.getElementById('ig-upload-zone');
    if (previewArea) previewArea.innerHTML = '';
    if (placeholder) placeholder.style.display = '';
    if (uploadZone) uploadZone.style.display = 'flex';
    document.querySelectorAll('.ig-upload-mini').forEach(z => {
        if (z.id === 'ig-upload-controlnet') return;
        z.classList.remove('ig-uploaded');
        const span = z.querySelector('span');
        if (span) span.textContent = 'Drop image or click to browse';
    });
    const targetEl = document.getElementById('ig-upscale-target');
    if (targetEl) { targetEl.style.display = 'none'; targetEl.textContent = ''; }
    img2imgMaskData = null;
    window.img2imgMaskData = null;
    const maskZone = document.getElementById('ig-upload-mini-img2img-mask');
    if (maskZone) {
        maskZone.classList.remove('ig-uploaded');
        const span = maskZone.querySelector('span');
        if (span) span.textContent = 'Drop mask image or click';
    }
    clearMask();
    setStatus('Image cleared.');
}

function updateUpscaleTarget(origW, origH, factor) {
    const targetEl = document.getElementById('ig-upscale-target');
    if (!targetEl) return;
    const w = Math.round(origW * factor / 64) * 64;
    const h = Math.round(origH * factor / 64) * 64;
    targetEl.textContent = `Target: ${w} x ${h} px (${factor}x)`;
    targetEl.style.display = '';
}

// --- Mask Drawing ---
let maskCtx = null;
let isDrawing = false;
let lastPos = null;

function resizeMaskCanvas() {
    const canvas = document.getElementById('ig-mask-canvas');
    const img = document.getElementById('ig-init-image');
    if (!canvas || !img) return;

    // Match canvas size to the displayed image size
    const rect = img.getBoundingClientRect();
    const wrapper = document.querySelector('.ig-canvas-wrapper');
    if (!wrapper) return;
    const wrapRect = wrapper.getBoundingClientRect();

    canvas.width = rect.width;
    canvas.height = rect.height;
    canvas.style.width = rect.width + 'px';
    canvas.style.height = rect.height + 'px';

    maskCtx = canvas.getContext('2d');
    clearMask();
}

function clearMask() {
    if (!maskCtx) return;
    maskCtx.fillStyle = '#000000';
    maskCtx.fillRect(0, 0, maskCtx.canvas.width, maskCtx.canvas.height);
}

function getMaskBase64() {
    if (!maskCtx) return '';
    return maskCtx.canvas.toDataURL('image/png').split(',')[1];
}

function initMaskDrawing() {
    const canvas = document.getElementById('ig-mask-canvas');
    if (!canvas) return;

    canvas.addEventListener('mousedown', (e) => {
        if (currentMode !== 'inpaint') return;
        isDrawing = true;
        lastPos = { x: e.offsetX, y: e.offsetY };
    });
    canvas.addEventListener('mousemove', (e) => {
        if (!isDrawing || currentMode !== 'inpaint') return;
        drawStroke(lastPos.x, lastPos.y, e.offsetX, e.offsetY);
        lastPos = { x: e.offsetX, y: e.offsetY };
    });
    window.addEventListener('mouseup', () => { isDrawing = false; lastPos = null; });

    // Touch support
    canvas.addEventListener('touchstart', (e) => {
        if (currentMode !== 'inpaint') return;
        e.preventDefault();
        const t = e.touches[0];
        const rect = canvas.getBoundingClientRect();
        isDrawing = true;
        lastPos = { x: t.clientX - rect.left, y: t.clientY - rect.top };
    }, { passive: false });
    canvas.addEventListener('touchmove', (e) => {
        if (!isDrawing || currentMode !== 'inpaint') return;
        e.preventDefault();
        const t = e.touches[0];
        const rect = canvas.getBoundingClientRect();
        const pos = { x: t.clientX - rect.left, y: t.clientY - rect.top };
        drawStroke(lastPos.x, lastPos.y, pos.x, pos.y);
        lastPos = pos;
    }, { passive: false });
    window.addEventListener('touchend', () => { isDrawing = false; lastPos = null; });
}

function drawStroke(x0, y0, x1, y1) {
    if (!maskCtx) return;
    const size = parseInt(document.getElementById('ig-brush-size')?.value || '16', 10);
    maskCtx.globalCompositeOperation = 'source-over';
    maskCtx.strokeStyle = '#ffffff';
    maskCtx.lineWidth = size;
    maskCtx.lineCap = 'round';
    maskCtx.lineJoin = 'round';
    maskCtx.beginPath();
    maskCtx.moveTo(x0, y0);
    maskCtx.lineTo(x1, y1);
    maskCtx.stroke();
}

// ---------------------------------------------------------------------------
// Model Control
// ---------------------------------------------------------------------------
function onUnloadModel() {
    setStatus('Unloading model...');
    if (window.go?.main?.App?.UnloadImageModel) {
        window.go.main.App.UnloadImageModel()
            .then(() => {
                setStatus('Model unloaded');
                const statusEl = document.getElementById('ig-model-status');
                if (statusEl) { statusEl.textContent = 'No model loaded'; statusEl.className = 'ig-panel-status'; }
                updateLiveModelInfo(null);
            })
            .catch(err => { setStatus('Unload failed: ' + err.message); });
    } else {
        setStatus('Unload not available in this backend.');
    }
}

function onReloadModel() {
    const name = document.getElementById('ig-model').value;
    if (!name) { setStatus('No model selected to reload.'); return; }
    setModel(name);
}

// Ctrl+L / "Load Model" button handler.
// Explicitly triggers a backend model load for the currently selected model.
async function onLoadModelClick() {
    const name = document.getElementById('ig-model').value;
    if (!name) { setStatus('No model selected to load.'); return; }
    await setModel(name);
}

// ---------------------------------------------------------------------------
// Presets
// ---------------------------------------------------------------------------
const PRESETS = {
    default:       { prompt: '', negative: '', steps: 20, cfg: 7,  sampler: 'euler_a', scheduler: 'karras', width: 512, height: 512 },
    realistic:     { prompt: 'masterpiece, best quality, highly detailed, realistic, 8k uhd, professional photography', negative: 'cartoon, anime, 3d render, blurry, low quality, deformed', steps: 25, cfg: 7, sampler: 'dpm++2m', scheduler: 'karras', width: 512, height: 768 },
    cinematic:     { prompt: 'cinematic shot, dramatic lighting, film grain, color graded, movie still', negative: 'overexposed, underexposed, video game, cartoon', steps: 30, cfg: 7, sampler: 'dpm++2m', scheduler: 'karras', width: 1024, height: 576 },
    anime:         { prompt: 'masterpiece, best quality, anime style, detailed background, vibrant colors', negative: 'realistic, 3d, western cartoon, photo, blurry', steps: 25, cfg: 7, sampler: 'euler_a', scheduler: 'karras', width: 512, height: 768 },
    illustration:  { prompt: 'detailed illustration, artistic, stylized, clean linework, concept art', negative: 'photograph, realistic, 3d render, noisy', steps: 25, cfg: 7, sampler: 'euler_a', scheduler: 'normal', width: 768, height: 768 },
    'low-light':   { prompt: 'night scene, moody lighting, dark atmosphere, cinematic shadows', negative: 'bright, overexposed, daylight, flat lighting', steps: 30, cfg: 8, sampler: 'dpm++2m', scheduler: 'karras', width: 768, height: 512 },
    'high-detail': { prompt: 'hyperdetailed, intricate, 8k resolution, sharp focus, ultra detailed', negative: 'simple, low detail, blurry, soft focus', steps: 40, cfg: 7, sampler: 'dpm++2m', scheduler: 'exponential', width: 1024, height: 1024 },
};

function applyPreset() {
    const key = document.getElementById('ig-preset').value;
    if (!key || !PRESETS[key]) { setStatus('Select a preset first.'); return; }
    const p = PRESETS[key];
    if (p.prompt) document.getElementById('ig-prompt').value = p.prompt;
    if (p.negative) document.getElementById('ig-negative').value = p.negative;
    document.getElementById('ig-steps').value = p.steps;
    document.getElementById('ig-cfg').value = p.cfg;
    document.getElementById('ig-sampler').value = p.sampler;
    document.getElementById('ig-scheduler').value = p.scheduler;
    document.getElementById('ig-width').value = p.width;
    document.getElementById('ig-height').value = p.height;
    updateLabels();
    saveState();
    setStatus('Preset applied: ' + key);
}

function saveCurrentPreset() {
    const name = prompt('Enter preset name (no spaces):');
    if (!name) return;
    const key = name.toLowerCase().replace(/\s+/g, '-');
    PRESETS[key] = {
        prompt: document.getElementById('ig-prompt').value,
        negative: document.getElementById('ig-negative').value,
        steps: parseInt(document.getElementById('ig-steps').value, 10),
        cfg: parseFloat(document.getElementById('ig-cfg').value),
        sampler: document.getElementById('ig-sampler').value,
        scheduler: document.getElementById('ig-scheduler').value,
        width: parseInt(document.getElementById('ig-width').value, 10),
        height: parseInt(document.getElementById('ig-height').value, 10),
    };
    const sel = document.getElementById('ig-preset');
    const opt = document.createElement('option');
    opt.value = key;
    opt.textContent = name;
    sel.appendChild(opt);
    sel.value = key;
    setStatus('Preset saved: ' + name);
}

// ---------------------------------------------------------------------------
// Seed Tools
// ---------------------------------------------------------------------------
function randomizeSeed() {
    const seed = Math.floor(Math.random() * 2147483647);
    document.getElementById('ig-seed').value = seed;
    saveState();
}

let seedLocked = false;
function toggleSeedLock() {
    seedLocked = !seedLocked;
    const btn = document.getElementById('ig-seed-lock');
    if (btn) btn.classList.toggle('ig-locked', seedLocked);
}

let upscaleSeedLocked = false;
function toggleUpscaleSeedLock() {
    upscaleSeedLocked = !upscaleSeedLocked;
    const btn = document.getElementById('ig-upscale-lock');
    if (btn) btn.classList.toggle('ig-locked', upscaleSeedLocked);
    setStatus(upscaleSeedLocked ? 'Upscale seed locked' : 'Upscale seed will be randomized');
}

function reuseLastSeed() {
    if (lastUsedSeed !== -1) {
        document.getElementById('ig-seed').value = lastUsedSeed;
        saveState();
    } else {
        setStatus('No previous seed to reuse.');
    }
}

// ---------------------------------------------------------------------------
// Prompt Weighting
// ---------------------------------------------------------------------------
function wrapPromptWeight() {
    const ta = document.getElementById('ig-prompt');
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const text = ta.value;
    const before = text.slice(0, start);
    const selected = text.slice(start, end);
    const after = text.slice(end);
    const newText = before + '(' + (selected || 'word') + ')' + after;
    ta.value = newText;
    ta.selectionStart = start + 1;
    ta.selectionEnd = start + 1 + (selected ? selected.length : 4);
    ta.focus();
    saveState();
}

function adjustWeight(delta) {
    const ta = document.getElementById('ig-prompt');
    const text = ta.value;
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const selected = text.slice(start, end);
    if (!selected) { setStatus('Select text to weight first.'); return; }
    let match = selected.match(/^\(([\s\S]+)\)$/);
    if (match) {
        const inner = match[1];
        const numMatch = inner.match(/^(.*):(\d+(\.\d+)?)$/);
        if (numMatch) {
            const newVal = Math.max(0, Math.min(5, parseFloat(numMatch[2]) + delta * 0.1));
            const newText = text.slice(0, start) + '(' + numMatch[1] + ':' + newVal.toFixed(1) + ')' + text.slice(end);
            ta.value = newText;
        } else {
            const newVal = (1 + delta * 0.1).toFixed(1);
            const newText = text.slice(0, start) + '(' + inner + ':' + newVal + ')' + text.slice(end);
            ta.value = newText;
        }
    } else {
        const newVal = (1 + delta * 0.1).toFixed(1);
        const newText = text.slice(0, start) + '(' + selected + ':' + newVal + ')' + text.slice(end);
        ta.value = newText;
    }
    ta.focus();
    saveState();
}

// ---------------------------------------------------------------------------
// Keyboard Shortcuts
// ---------------------------------------------------------------------------
function handleKeydown(e) {
    // Ctrl+Enter = Generate
    if (e.ctrlKey && e.key === 'Enter') {
        e.preventDefault();
        onGenerate();
        return;
    }
    // Ctrl+S = Save Image As
    if (e.ctrlKey && e.key === 's' && !e.shiftKey) {
        e.preventDefault();
        saveImageAs();
        return;
    }
    // Ctrl+L = Load Model
    if (e.ctrlKey && e.key === 'l') {
        e.preventDefault();
        onLoadModelClick();
        return;
    }
    // Ctrl+U = Unload Model
    if (e.ctrlKey && e.key === 'u') {
        e.preventDefault();
        onUnloadModel();
        return;
    }
    // Ctrl+R = Reload Model
    if (e.ctrlKey && e.key === 'r') {
        e.preventDefault();
        onReloadModel();
        return;
    }
    // R = Random Seed (when not in an input)
    if (e.key === 'r' && !e.ctrlKey && !e.altKey && !e.metaKey) {
        const active = document.activeElement;
        if (active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA')) return;
        e.preventDefault();
        randomizeSeed();
    }
    // Workflow mode shortcuts
    if (e.ctrlKey) {
        const modeMap = { '1': 'txt2img', '2': 'img2img', '3': 'inpaint', '4': 'outpaint', '5': 'upscale', '6': 'video' };
        if (modeMap[e.key]) {
            e.preventDefault();
            setWorkflowMode(modeMap[e.key]);
        }
    }
}
