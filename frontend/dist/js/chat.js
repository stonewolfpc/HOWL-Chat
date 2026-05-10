/**
 * Chat Module
 * 
 * This module handles chat message interactions, including adding messages,
 * streaming responses, and message actions (copy, edit, regenerate, continue).
 * 
 * @module chat
 */

// Add user message to chat
function addUserMessage(text) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message user';
    messageDiv.innerHTML = `
        <div class="chat-bubble">
            <div class="chat-sender">You</div>
            <div class="chat-text">${text}</div>
            <div class="chat-bubble-buttons">
                <button class="bubble-button" onclick="copyMessage(this)">Copy</button>
                <button class="bubble-button" onclick="editMessage(this)">Edit</button>
            </div>
        </div>
    `;
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

// Add assistant message to chat
function addAssistantMessage(sender, text) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant';
    messageDiv.innerHTML = `
        <div class="chat-bubble">
            <div class="chat-sender">${sender}</div>
            <div class="chat-text">${text}</div>
            <div class="chat-bubble-buttons">
                <button class="bubble-button" onclick="copyMessage(this)">Copy</button>
                <button class="bubble-button" onclick="regenerateMessage(this)">Regenerate</button>
                <button class="bubble-button" onclick="continueMessage(this)">Continue</button>
            </div>
        </div>
    `;
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

// Stream assistant response (simulated)
function streamAssistantResponse(sender, fullText) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant';
    
    const bubble = document.createElement('div');
    bubble.className = 'chat-bubble';
    bubble.innerHTML = `
        <div class="chat-sender"><span class="streaming-indicator"></span>${sender}</div>
        <div class="chat-text"></div>
        <div class="chat-bubble-buttons">
            <button class="bubble-button" onclick="copyMessage(this)">Copy</button>
            <button class="bubble-button" onclick="regenerateMessage(this)">Regenerate</button>
            <button class="bubble-button" onclick="continueMessage(this)">Continue</button>
        </div>
    `;
    
    messageDiv.appendChild(bubble);
    chatHistory.appendChild(messageDiv);
    
    const textElement = bubble.querySelector('.chat-text');
    let currentIndex = 0;
    
    // Simulate streaming by adding characters one at a time
    const streamInterval = setInterval(() => {
        if (currentIndex < fullText.length) {
            textElement.textContent += fullText[currentIndex];
            currentIndex++;
            chatHistory.scrollTop = chatHistory.scrollHeight;
        } else {
            clearInterval(streamInterval);
            // Remove streaming indicator
            const indicator = bubble.querySelector('.streaming-indicator');
            if (indicator) {
                indicator.remove();
            }
        }
    }, 30); // 30ms per character for smooth streaming
}

// Copy message text
function copyMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const text = bubble.querySelector('.chat-text').textContent;
    navigator.clipboard.writeText(text).then(() => {
        const originalText = button.textContent;
        button.textContent = 'Copied!';
        setTimeout(() => {
            button.textContent = originalText;
        }, 1500);
    });
}

// Edit message (user messages only)
function editMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const textElement = bubble.querySelector('.chat-text');
    const currentText = textElement.textContent;
    
    const input = document.querySelector('.input-field');
    input.value = currentText;
    input.focus();
    
    // Remove the message for editing
    const messageDiv = button.closest('.chat-message');
    messageDiv.remove();
}

// Regenerate message (assistant messages only)
function regenerateMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const sender = bubble.querySelector('.chat-sender').textContent;
    
    // Remove current message
    const messageDiv = button.closest('.chat-message');
    messageDiv.remove();
    
    // Stream new response
    const responses = [
        "Let me reconsider that from a different perspective.",
        "Here's an alternative approach to your request.",
        "I'll provide a different response based on the context.",
        "Allow me to rephrase that for better clarity."
    ];
    const randomResponse = responses[Math.floor(Math.random() * responses.length)];
    streamAssistantResponse(sender, randomResponse);
}

// Continue message (assistant messages only)
function continueMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const sender = bubble.querySelector('.chat-sender').textContent;
    
    const continuations = [
        " Furthermore, I should mention that this approach has additional benefits.",
        " Additionally, there are several other factors to consider in this context.",
        " Moreover, the implications of this decision extend beyond the immediate scope.",
        " In addition to what I've mentioned, there are further details worth exploring."
    ];
    const randomContinuation = continuations[Math.floor(Math.random() * continuations.length)];
    streamAssistantResponse(sender, randomContinuation);
}

// Handle send button
document.querySelector('.send-button').addEventListener('click', () => {
    const input = document.querySelector('.input-field');
    const message = input.value.trim();
    if (message) {
        addUserMessage(message);
        input.value = '';
        
        // Simulate AI response with streaming
        setTimeout(() => {
            streamAssistantResponse("Knight", "I understand your request. Let me assist you with that matter.");
        }, 500);
    }
});

// Handle Enter key in input field
document.querySelector('.input-field').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        document.querySelector('.send-button').click();
    }
});

// Load Model - opens file picker
function loadModel() {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.gguf,.bin,.safetensors';
    
    input.onchange = (e) => {
        const file = e.target.files[0];
        if (file) {
            console.log('Loading model:', file.name);
            // TODO: Backend integration for model loading
            // This will be implemented when backend is ready
            alert(`Model selected: ${file.name}\n\nBackend integration pending.`);
        }
    };
    
    input.click();
}

// Unload Model - aggressively unloads current model
function unloadModel() {
    if (confirm('⚠️ AGGRESSIVE UNLOAD ⚠️\n\nThis will forcefully unload the current model from memory.\nAll unsaved context will be lost.\n\nAre you sure you want to proceed?')) {
        console.log('AGGRESSIVE UNLOAD initiated');
        // TODO: Backend integration for aggressive model unloading
        // This will implement:
        // - Force memory cleanup
        // - Clear GPU cache
        // - Terminate model process
        // - Reset all model state
        // - Clear conversation context
        alert('🔴 MODEL AGGRESSIVELY UNLOADED\n\nGPU cache cleared.\nMemory freed.\nModel process terminated.\n\nBackend integration pending.');
    }
}

// Placeholder functions for new button handlers
function undoMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const messageDiv = button.closest('.chat-message');
    messageDiv.remove();
    console.log('Message undone');
}

function previousMessage(button) {
    console.log('Previous message');
    // TODO: Implement message history navigation
}

function nextMessage(button) {
    console.log('Next message');
    // TODO: Implement message history navigation
}

function deleteMessage(button) {
    const messageDiv = button.closest('.chat-message');
    messageDiv.remove();
    console.log('Message deleted');
}

function addImage() {
    console.log('Add image');
    // TODO: Implement image upload
}

function addAudio() {
    console.log('Add audio');
    // TODO: Implement audio recording
}

function textToSpeech() {
    console.log('Text to speech');
    // TODO: Implement TTS
}

function addFile() {
    console.log('Add file');
    // TODO: Implement document upload for RAG 2
}
