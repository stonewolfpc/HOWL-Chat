/**
 * Menu Interaction Module
 * 
 * This module handles dropdown menu toggling and interactions
 * for the top menu bar in the HOWL Chat application.
 * 
 * @module menu
 */

// Menu toggling
function toggleMenu(menuId) {
    // Close all other menus first
    const allMenus = document.querySelectorAll('.dropdown');
    allMenus.forEach(menu => {
        if (menu.id !== menuId) {
            menu.classList.remove('show');
        }
    });
    
    // Toggle the requested menu
    const menu = document.getElementById(menuId);
    if (menu) {
        menu.classList.toggle('show');
    }
}

// Flyout toggling
function toggleFlyout(flyoutId) {
    const flyout = document.getElementById(flyoutId);
    const overlay = document.getElementById('flyout-overlay');
    
    if (flyout && overlay) {
        flyout.classList.toggle('show');
        overlay.classList.toggle('show');
    }
    
    // Close all dropdown menus when opening flyout
    const allMenus = document.querySelectorAll('.dropdown');
    allMenus.forEach(menu => {
        menu.classList.remove('show');
    });
}

// Close menus when clicking outside
document.addEventListener('click', function(event) {
    if (!event.target.closest('.menu-item')) {
        const allMenus = document.querySelectorAll('.dropdown');
        allMenus.forEach(menu => {
            menu.classList.remove('show');
        });
    }
});

// Handle menu dropdowns
document.querySelectorAll('.dropdown-item').forEach(item => {
    item.addEventListener('click', () => {
        console.log('Menu item clicked:', item.textContent);
    });
});

// Exit application - save chat and quit
async function exitApp() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return;
    }
    
    try {
        // Use backend ExitApp which handles final save and shutdown
        await window.go.main.App.ExitApp();
    } catch (error) {
        console.error('Error during application exit:', error);
        // Fallback: try to quit anyway
        if (window.runtime) {
            window.runtime.Quit();
        }
    }
}
