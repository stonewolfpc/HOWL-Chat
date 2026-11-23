/**
 * Sidebar Context — World Section
 *
 * World display update, world selector modal, world selection,
 * and world auto-switch toggle.
 *
 * Depends on: sidebar-state.js, sidebar-session.js
 *
 * @module sidebar-context-world
 * @version 1.0.0
 */

async function updateWorldDisplay() {
    const worldNameEl      = document.getElementById('world-name-display');
    const worldDescEl      = document.getElementById('world-description-display');
    const worldImageEl     = document.getElementById('world-image');
    const worldThumbnailEl = document.getElementById('world-thumbnail');
    if (!worldNameEl) return;
    if (!window.sidebarState) return;

    const cd          = window.sidebarState.currentChatData;
    const placeholder = worldThumbnailEl ? worldThumbnailEl.querySelector('.thumbnail-placeholder') : null;

    if (cd.world) {
        worldNameEl.textContent = cd.world.title || cd.world.name;
        if (worldDescEl) {
            worldDescEl.textContent = cd.world.description || cd.world.details || '';
        }

        if (worldImageEl) {
            // New multi-image array format: use thumbnail from images array
            let imageSrc = cd.world._resolvedImage || '';

            if (!imageSrc && cd.world.images && Array.isArray(cd.world.images) && cd.world.images.length > 0) {
                const thumbIdx = cd.world.thumbnailIndex || 0;
                const thumb = cd.world.images[thumbIdx] || cd.world.images[0];
                imageSrc = thumb.data || thumb.thumbnail || '';
            }

            // Legacy single-image format fallback
            if (!imageSrc && cd.world.image) {
                imageSrc = window.sidebarState.resolveImageSrc(cd.world.image) || '';
            }

            // Async fallback: resolve {id, stored:true} reference from IndexedDB (legacy)
            if (!imageSrc && cd.world.image && window.HOWLStorage) {
                try {
                    const ref     = cd.world.image;
                    const imageId = (typeof ref === 'object' && ref.id) ? ref.id : null;
                    if (imageId) {
                        const imageData = await window.HOWLStorage.getImage(imageId, 'world');
                        if (imageData && imageData.data) {
                            imageSrc = imageData.data;
                            cd.world._resolvedImage = imageSrc;
                        } else {
                            console.warn('World image not found in IndexedDB:', imageId);
                        }
                    }
                } catch (error) {
                    console.warn('Failed to restore world image:', error);
                }
            }

            if (imageSrc) {
                worldImageEl.src = imageSrc;
                worldImageEl.style.display = 'block';
                if (placeholder) placeholder.style.display = 'none';
            } else {
                worldImageEl.style.display = 'none';
                if (placeholder) placeholder.style.display = '';
            }
        }
    } else {
        worldNameEl.textContent = 'No World Selected';
        if (worldDescEl) worldDescEl.textContent = '';
        if (worldImageEl) worldImageEl.style.display = 'none';
        if (placeholder) placeholder.style.display = '';
    }
}

function openWorldSelector() {
    const worlds = window.WorldRepository.all();
    const cd     = window.sidebarState.currentChatData;

    const modal = document.createElement('div');
    modal.className   = 'selector-modal';
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);`;

    let content = `
        <div style="background:var(--shell-panel-bg);border-radius:12px;padding:24px;max-width:400px;border:1px solid var(--surface-hover-strong);">
            <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">Select a World</h2>
            <div style="max-height:300px;overflow-y:auto;margin-bottom:16px;">
    `;

    if (worlds.length === 0) {
        content += `<div style="color:var(--surface-hover);text-align:center;padding:20px;">No worlds available</div>`;
    } else {
        worlds.forEach(world => {
            const isActive = cd.world && cd.world.id === world.id;
            content += `
                <div onclick="selectWorld('${world.id}'); document.querySelector('.selector-modal').remove();" style="
                    padding:12px;background:${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'};
                    margin-bottom:8px;border-radius:6px;cursor:pointer;
                    border:1px solid ${isActive ? 'var(--success)' : 'rgba(255,255,255,0.1)'};color:#fff;transition:all 0.2s ease;"
                    onmouseover="this.style.background='var(--surface-hover-strong)'"
                    onmouseout="this.style.background='${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'}'">
                    <strong style="display:block;margin-bottom:4px;">${world.title || world.name}</strong>
                    <small style="color:var(--surface-hover);">${world.description || 'No description'}</small>
                </div>
            `;
        });
    }

    content += `
            </div>
            <div style="display:flex;gap:8px;">
                <button onclick="document.querySelector('.selector-modal').remove();" style="flex:1;padding:10px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
                <button onclick="navigateTo('world.html'); document.querySelector('.selector-modal').remove();" style="flex:1;padding:10px;background:var(--success-bg);border:1px solid var(--success);color:#00b873;border-radius:6px;cursor:pointer;">+ New World</button>
            </div>
        </div>
    `;

    modal.innerHTML = content;
    document.body.appendChild(modal);
    modal.addEventListener('click', (e) => { if (e.target === modal) modal.remove(); });
}

async function selectWorld(worldId) {
    const worlds = window.WorldRepository.all();
    const world  = worlds.find(w => w.id === worldId);
    if (!world) { console.error('World not found:', worldId); return; }

    if (window.chatBridge) {
        await window.chatBridge.applyWorldChange(world);
    } else if (typeof contextDisplay !== 'undefined' && contextDisplay.selectWorld) {
        await contextDisplay.selectWorld(worldId);
    }

    if (window.go?.main?.App?.AddWorldToChat) {
        try { await window.go.main.App.AddWorldToChat(world); }
        catch (error) { console.warn('Backend world sync failed:', error); }
    }
}

function toggleWorldAutoSwitch() {
    const checkbox = document.getElementById('world-auto-switch');
    console.log('World auto-switch:', checkbox.checked);
}
