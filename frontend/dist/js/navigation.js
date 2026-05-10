/**
 * Navigation Module
 * 
 * This module handles page navigation between the different sections
 * of the HOWL Chat application (Chat, World, Characters, Scenarios, Lorebooks).
 * 
 * @module navigation
 */

// Navigate to a specific page
function navigateTo(page) {
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
