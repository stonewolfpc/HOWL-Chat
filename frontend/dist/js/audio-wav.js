/**
 * PCM-to-WAV encoding for microphone input.
 *
 * The renderer captures mono Float32 PCM through Web Audio. This module
 * resamples it to 16 kHz and emits a standard 16-bit PCM WAV, so the backend
 * never needs a global ffmpeg installation to accept a recording.
 */
(function (root, factory) {
    if (typeof module === 'object' && module.exports) {
        module.exports = factory();
    } else {
        root.AudioWav = factory();
    }
})(typeof globalThis !== 'undefined' ? globalThis : this, function () {
    'use strict';

    const TARGET_RATE = 16000;

    function resampleTo16k(samples, inputRate) {
        if (!samples || samples.length === 0) return new Float32Array(0);
        if (!inputRate || inputRate === TARGET_RATE) return new Float32Array(samples);

        const ratio = inputRate / TARGET_RATE;
        const length = Math.max(1, Math.round(samples.length / ratio));
        const output = new Float32Array(length);

        if (ratio >= 1.5) {
            // Box-filter the source window before decimating. This is not a
            // studio resampler, but it avoids the worst aliasing for speech.
            for (let i = 0; i < length; i++) {
                const start = Math.floor(i * ratio);
                const end = Math.min(samples.length, Math.max(start + 1, Math.floor((i + 1) * ratio)));
                let sum = 0;
                for (let j = start; j < end; j++) sum += samples[j];
                output[i] = sum / (end - start);
            }
            return output;
        }

        for (let i = 0; i < length; i++) {
            const position = i * ratio;
            const left = Math.floor(position);
            const right = Math.min(left + 1, samples.length - 1);
            const mix = position - left;
            output[i] = samples[left] * (1 - mix) + samples[right] * mix;
        }
        return output;
    }

    function writeAscii(view, offset, text) {
        for (let i = 0; i < text.length; i++) view.setUint8(offset + i, text.charCodeAt(i));
    }

    function encodeWav16(samples, sampleRate) {
        const pcm = resampleTo16k(samples, sampleRate);
        const dataSize = pcm.length * 2;
        const buffer = new ArrayBuffer(44 + dataSize);
        const view = new DataView(buffer);

        writeAscii(view, 0, 'RIFF');
        view.setUint32(4, 36 + dataSize, true);
        writeAscii(view, 8, 'WAVE');
        writeAscii(view, 12, 'fmt ');
        view.setUint32(16, 16, true);
        view.setUint16(20, 1, true);          // PCM
        view.setUint16(22, 1, true);          // mono
        view.setUint32(24, TARGET_RATE, true);
        view.setUint32(28, TARGET_RATE * 2, true);
        view.setUint16(32, 2, true);
        view.setUint16(34, 16, true);
        writeAscii(view, 36, 'data');
        view.setUint32(40, dataSize, true);

        let offset = 44;
        for (let i = 0; i < pcm.length; i++, offset += 2) {
            let sample = pcm[i];
            if (!Number.isFinite(sample)) sample = 0;
            sample = Math.max(-1, Math.min(1, sample));
            view.setInt16(offset, sample < 0 ? sample * 0x8000 : sample * 0x7fff, true);
        }
        return new Uint8Array(buffer);
    }

    return { TARGET_RATE, resampleTo16k, encodeWav16 };
});
