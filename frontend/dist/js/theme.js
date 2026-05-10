/**
 * Theme Management Module
 * 
 * This module handles theme cycling through the Artherius color palette
 * with smooth transitions and animations.
 * 
 * @module theme
 */

// Theme definitions for Artherius color palette
const themes = [
    {
        color: '#14b8a6',
        glow: 'rgba(20, 184, 166, 0.5)',
        border: 'rgba(20, 184, 166, 0.3)',
        textGlow: 'rgba(20, 184, 166, 0.4)',
        secondary: '#06b6d4'
    },
    {
        color: '#06b6d4',
        glow: 'rgba(6, 182, 212, 0.5)',
        border: 'rgba(6, 182, 212, 0.3)',
        textGlow: 'rgba(6, 182, 212, 0.4)',
        secondary: '#14b8a6'
    },
    {
        color: '#a855f7',
        glow: 'rgba(168, 85, 247, 0.5)',
        border: 'rgba(168, 85, 247, 0.3)',
        textGlow: 'rgba(168, 85, 247, 0.4)',
        secondary: '#ec4899'
    }
];

let currentThemeIndex = 2; // Start with cyan
const targetFPS = 30; // Target 30fps for smooth animation
const frameInterval = 1000 / targetFPS; // ~33.33ms per frame
let animationStartTime = 0;
let isAnimating = false;

// Easing function for smooth transitions
function easeInOutCubic(t) {
    return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

// Smoother easing function for borders
function easeInOutQuad(t) {
    return t < 0.5 ? 2 * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

// Fade color interpolation
function fadeColor(color1, color2, progress) {
    const r1 = parseInt(color1.slice(1, 3), 16);
    const g1 = parseInt(color1.slice(3, 5), 16);
    const b1 = parseInt(color1.slice(5, 7), 16);
    
    const r2 = parseInt(color2.slice(1, 3), 16);
    const g2 = parseInt(color2.slice(3, 5), 16);
    const b2 = parseInt(color2.slice(5, 7), 16);
    
    const easedProgress = easeInOutCubic(progress);
    
    const r = Math.round(r1 + (r2 - r1) * easedProgress);
    const g = Math.round(g1 + (g2 - g1) * easedProgress);
    const b = Math.round(b1 + (b2 - b1) * easedProgress);
    
    return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
}

// Fade RGBA interpolation
function fadeRGBA(rgba1, rgba2, progress) {
    const r1 = parseInt(rgba1.slice(5, 8), 16);
    const g1 = parseInt(rgba1.slice(9, 12), 16);
    const b1 = parseInt(rgba1.slice(13, 16), 16);
    const a1 = parseFloat(rgba1.slice(17, 20));
    
    const r2 = parseInt(rgba2.slice(5, 8), 16);
    const g2 = parseInt(rgba2.slice(9, 12), 16);
    const b2 = parseInt(rgba2.slice(13, 16), 16);
    const a2 = parseFloat(rgba2.slice(17, 20));
    
    const easedProgress = easeInOutCubic(progress);
    
    const r = Math.round(r1 + (r2 - r1) * easedProgress);
    const g = Math.round(g1 + (g2 - g1) * easedProgress);
    const b = Math.round(b1 + (b2 - b1) * easedProgress);
    const a = (a1 + (a2 - a1) * easedProgress).toFixed(2);
    
    return `rgba(${r}, ${g}, ${b}, ${a})`;
}

// Background color for fading to
const backgroundColor = '#0a0a0f';

// Theme cycling function
function cycleTheme() {
    if (isAnimating) return;
    
    isAnimating = true;
    animationStartTime = performance.now();
    
    const currentTheme = themes[currentThemeIndex];
    const nextThemeIndex = (currentThemeIndex + 1) % themes.length;
    const nextTheme = themes[nextThemeIndex];
    
    const fadeDuration = 4000; // 4 seconds for fade
    const pauseDuration = 1000; // 1 second pause
    
    function animate(timestamp) {
        const elapsed = timestamp - animationStartTime;
        const progress = Math.min(elapsed / fadeDuration, 1);
        
        // Fade to background first
        if (progress < 0.5) {
            const fadeProgress = progress * 2;
            document.documentElement.style.setProperty('--theme-color', fadeColor(currentTheme.color, backgroundColor, fadeProgress));
            document.documentElement.style.setProperty('--theme-glow', fadeRGBA(currentTheme.glow, 'rgba(15, 21, 31, 0.5)', fadeProgress));
            // Keep border static - don't change to prevent pulsing
            document.documentElement.style.setProperty('--theme-text-glow', fadeRGBA(currentTheme.textGlow, 'rgba(100, 100, 100, 0.4)', fadeProgress));
        } else {
            // Fade to new theme
            const fadeProgress = (progress - 0.5) * 2;
            document.documentElement.style.setProperty('--theme-color', fadeColor(backgroundColor, nextTheme.color, fadeProgress));
            document.documentElement.style.setProperty('--theme-glow', fadeRGBA('rgba(15, 21, 31, 0.5)', nextTheme.glow, fadeProgress));
            // Keep border static - don't change to prevent pulsing
            document.documentElement.style.setProperty('--theme-text-glow', fadeRGBA('rgba(100, 100, 100, 0.4)', nextTheme.textGlow, fadeProgress));
        }

        if (progress < 1) {
            requestAnimationFrame(animate);
        } else {
            // Set final theme (keep border static)
            document.documentElement.style.setProperty('--theme-color', nextTheme.color);
            document.documentElement.style.setProperty('--theme-glow', nextTheme.glow);
            // Don't change border - keep it static to prevent pulsing
            document.documentElement.style.setProperty('--theme-text-glow', nextTheme.textGlow);
            document.documentElement.style.setProperty('--theme-secondary', nextTheme.secondary);

            currentThemeIndex = nextThemeIndex;
            isAnimating = false;

            // Schedule next cycle
            setTimeout(cycleTheme, pauseDuration);
        }
    }
    
    requestAnimationFrame(animate);
}

// Start theme cycling after initial delay
setTimeout(cycleTheme, 2000);
