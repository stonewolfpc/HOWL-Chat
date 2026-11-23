/**
 * Navigation Module
 * 
 * This module handles page navigation between the different sections
 * of the HOWL Chat application (Chat, World, Characters, Scenarios, Lorebooks).
 * 
 * @module navigation
 */

// Navigate to a specific page.
//
// Leaving Chat is a save boundary: DataSync owns that write, so it is awaited
// here before the document unloads. Firing it during unload instead left the
// request aborted mid-flight ("Failed to save chat data") or racing the
// conversation boundary. One completed save, then navigate.
async function navigateTo(page) {
    const dataSync = window.dataSync;
    if (dataSync?.currentChatId && !window.howlNavigationSaved) {
        try {
            await dataSync.saveCurrentChatData();
            window.howlNavigationSaved = true;
        } catch (error) {
            console.error('[navigation] save before leaving failed:', error);
        }
    }
    window.location.href = page;
}

// Set active navigation item based on current page
function setActiveNavItem() {
    const currentPage = window.location.pathname.split('/').pop() || 'index.html';
    
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
    });
    
    // Find and activate the matching nav item
    document.querySelectorAll('.nav-item').forEach(item => {
        const onclick = item.getAttribute('onclick');
        if (onclick && onclick.includes(currentPage)) {
            item.classList.add('active');
        }
    });
}

// Initialize navigation on page load
document.addEventListener('DOMContentLoaded', () => {
    setActiveNavItem();
});
