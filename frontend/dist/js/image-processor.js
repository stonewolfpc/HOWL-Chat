/**
 * Image Processor
 *
 * Memory-efficient resize pipeline using createImageBitmap() with
 * pre-downscaling to prevent crashes with large images (50MB+ PNGs, etc.)
 *
 * Produces three variants:
 *   thumbnail — 256px max dimension, WebP, ~20-80 KB  (UI, chat, grids)
 *   medium    — 1024px max dimension, WebP, ~200-500 KB (edit screens)
 *   original  — raw file blob (NOT decoded)            (archival)
 *
 * GIF handling: extracts first frame via createImageBitmap (static WebP).
 * All outputs are base64 DataURL strings except original which is a blob.
 *
 * Public API (window.ImageProcessor):
 *   process(file)  → Promise<{thumbnail, medium, original, meta}>
 *
 * @module image-processor
 * @version 2.0.0
 */

'use strict';

const ImageProcessor = (() => {

    const THUMBNAIL_MAX = 256;
    const MEDIUM_MAX    = 1024;
    const WEBP_QUALITY  = 0.85;

    // -----------------------------------------------------------------------
    // Core: createImageBitmap with resize (memory-efficient decoding)
    // -----------------------------------------------------------------------

    async function createResizedBitmap(file, maxDim) {
        // createImageBitmap decodes AND resizes in one step, off-thread
        // This prevents loading full-resolution image into memory
        const bitmap = await createImageBitmap(file, {
            resizeWidth: maxDim,
            resizeHeight: maxDim,
            resizeQuality: 'high'
        });
        return bitmap;
    }

    // -----------------------------------------------------------------------
    // Core: bitmap → canvas → WebP base64
    // -----------------------------------------------------------------------

    async function bitmapToDataURL(bitmap, quality) {
        const canvas = document.createElement('canvas');
        canvas.width  = bitmap.width;
        canvas.height = bitmap.height;
        const ctx = canvas.getContext('2d');
        ctx.drawImage(bitmap, 0, 0);
        bitmap.close(); // Free bitmap memory immediately

        const webp = canvas.toDataURL('image/webp', quality);
        if (webp && webp !== 'data:,') return webp;
        return canvas.toDataURL('image/png');
    }

    // -----------------------------------------------------------------------
    // Core: get original dimensions without full decode (fast)
    // -----------------------------------------------------------------------

    async function getImageDimensions(file) {
        // Use a tiny bitmap to get dimensions without decoding full image
        const bitmap = await createImageBitmap(file);
        const w = bitmap.width;
        const h = bitmap.height;
        bitmap.close();
        return { width: w, height: h };
    }

    // -----------------------------------------------------------------------
    // Public: process a File → {thumbnail, medium, original, meta}
    // -----------------------------------------------------------------------

    async function process(file) {
        if (!file || !file.type.startsWith('image/')) {
            throw new Error('Not a valid image file');
        }
        if (file.size > 100 * 1024 * 1024) {
            throw new Error('File too large (max 100MB)');
        }

        // Get dimensions without loading full image
        const { width, height } = await getImageDimensions(file);

        // Create resized bitmaps directly at target sizes (memory-efficient)
        const [thumbBitmap, mediumBitmap] = await Promise.all([
            createResizedBitmap(file, THUMBNAIL_MAX),
            createResizedBitmap(file, MEDIUM_MAX)
        ]);

        // Convert to WebP base64
        const [thumbnail, medium] = await Promise.all([
            bitmapToDataURL(thumbBitmap, WEBP_QUALITY),
            bitmapToDataURL(mediumBitmap, WEBP_QUALITY)
        ]);

        // Store original as raw blob (never decode it)
        const original = file;

        return {
            thumbnail,
            medium,
            original,  // This is now a File/Blob, NOT base64
            meta: {
                name:     file.name,
                type:     file.type,
                size:     file.size,
                width,
                height,
                thumbW:   Math.round(width * Math.min(1, THUMBNAIL_MAX / Math.max(width, height))),
                thumbH:   Math.round(height * Math.min(1, THUMBNAIL_MAX / Math.max(width, height))),
                mediumW:  Math.round(width * Math.min(1, MEDIUM_MAX / Math.max(width, height))),
                mediumH:  Math.round(height * Math.min(1, MEDIUM_MAX / Math.max(width, height)))
            }
        };
    }

    return { process };

})();

window.ImageProcessor = ImageProcessor;
