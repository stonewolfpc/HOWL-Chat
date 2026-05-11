/**
 * Chat Module
 *
 * This module handles chat message interactions, including adding messages,
 * streaming responses, and message actions (copy, edit, delete, regenerate, continue, navigation).
 *
 * @module chat
 */

// SVG Icons for chat bubble buttons
const ICONS = {
    copy: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>`,
    edit: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>`,
    delete: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>`,
    regenerate: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>`,
    previous: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>`,
    next: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>`,
    continue: `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>`
};

// Wait for Wails runtime to be ready
let wailsReady = false;

// Check if Wails runtime is available
function checkWailsReady() {
    if (window.wails && window.wails._bindings) {
        wailsReady = true;
        console.log('Wails runtime is ready');
    } else {
        console.log('Waiting for Wails runtime...');
        setTimeout(checkWailsReady, 100);
    }
}

// Start checking for Wails runtime
checkWailsReady();

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
                <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
                <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
                <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
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
                <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
                <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
                <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
                <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
                <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
                <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
                <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
            </div>
        </div>
    `;
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;
    updateNavigationButtons();
}

// Stream assistant response (simulated)
function streamAssistantResponse(sender, fullText) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant';

    const bubble = document.createElement('div');
    bubble.className = 'chat-bubble';
    bubble.innerHTML = `
        <div class="chat-sender">${sender}</div>
        <div class="chat-text"></div>
        <div class="chat-bubble-buttons">
            <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
            <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
            <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
            <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
            <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
            <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
            <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
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
            updateNavigationButtons();
        }
    }, 30); // 30ms per character for smooth streaming
}

// Stream assistant response from backend with real token streaming
let currentStreamingBubble = null;
let thinkStartTime = null;
let thinkContent = '';
let mainContent = '';
let isInThinkMode = false;

// Regex patterns for think block detection
// Gemma 4 uses special tokens: <|channel>thought for start, <channel|> for end
const THINK_START_REGEX = /<\|channel>thought/i;
const THINK_END_REGEX = /<channel\|>/i;

function streamAssistantResponseFromBackend(message) {
    const chatHistory = document.querySelector('.chat-messages');

    // Create assistant message bubble with typing indicator
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant streaming';

    const bubble = document.createElement('div');
    bubble.className = 'chat-bubble';
    bubble.innerHTML = `
        <div class="chat-sender">AI<span class="typing-indicator"><span></span><span></span><span></span></span></div>
        <div class="chat-text"></div>
        <div class="chat-think"></div>
        <div class="chat-bubble-buttons">
            <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
            <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
            <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
            <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
            <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
            <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
            <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
        </div>
    `;

    messageDiv.appendChild(bubble);
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;

    // Set current streaming bubble reference so onChunk can find the elements
    currentStreamingBubble = bubble;

    const textElement = bubble.querySelector('.chat-text');
    const thinkArea = bubble.querySelector('.chat-think');

    // Reset think state
    thinkContent = '';
    mainContent = '';
    isInThinkMode = false;
    thinkStartTime = null;

    const onStart = () => {
        console.log('Stream started');
        // Switch send button to abort
        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Abort';
            sendButton.dataset.isAbort = 'true';
        }
    };

    const onChunk = (chunk) => {
        if (currentStreamingBubble) {
            // Check for start marker
            if (THINK_START_REGEX.test(chunk)) {
                isInThinkMode = true;
                thinkStartTime = Date.now();
                chunk = chunk.replace(THINK_START_REGEX, '');
            }

            // Check for end marker
            if (THINK_END_REGEX.test(chunk)) {
                isInThinkMode = false;
                chunk = chunk.replace(THINK_END_REGEX, '');
                // Update think area when done
                if (thinkContent) {
                    updateThinkArea(thinkArea, thinkContent, thinkStartTime);
                }
            }

            // Route tokens to separate fields
            if (isInThinkMode) {
                thinkContent += chunk;
            } else {
                mainContent += chunk;
                textElement.textContent = mainContent;
            }
            chatHistory.scrollTop = chatHistory.scrollHeight;
        }
    };

    const onComplete = () => {
        // Remove typing indicator
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) {
            indicator.remove();
        }
        messageDiv.classList.remove('streaming');

        // If there's think content, create collapsible area
        if (thinkContent) {
            updateThinkArea(thinkArea, thinkContent, thinkStartTime);
        }

        // Switch abort button back to send
        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Send';
            sendButton.dataset.isAbort = 'false';
        }

        // Clean up listeners
        if (window.runtime) {
            window.runtime.EventsOff('chat:chunk');
            window.runtime.EventsOff('chat:complete');
            window.runtime.EventsOff('chat:error');
            window.runtime.EventsOff('chat:aborted');
        }
        currentStreamingBubble = null;
        updateNavigationButtons();
    };

    const onError = (err) => {
        textElement.textContent += '\n[Error: ' + err + ']';
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) {
            indicator.remove();
        }
        messageDiv.classList.remove('streaming');

        // Switch abort button back to send
        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Send';
            sendButton.dataset.isAbort = 'false';
        }

        currentStreamingBubble = null;
    };

    const onAborted = () => {
        console.log('Stream aborted');
        // Remove typing indicator
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) {
            indicator.remove();
        }
        messageDiv.classList.remove('streaming');

        // Switch abort button back to send
        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Send';
            sendButton.dataset.isAbort = 'false';
        }

        // Clean up listeners
        if (window.runtime) {
            window.runtime.EventsOff('chat:chunk');
            window.runtime.EventsOff('chat:complete');
            window.runtime.EventsOff('chat:error');
            window.runtime.EventsOff('chat:aborted');
        }
        currentStreamingBubble = null;
    };

    // Subscribe to events
    if (window.runtime) {
        window.runtime.EventsOn('chat:start', onStart);
        window.runtime.EventsOn('chat:chunk', onChunk);
        window.runtime.EventsOn('chat:complete', onComplete);
        window.runtime.EventsOn('chat:error', onError);
        window.runtime.EventsOn('chat:aborted', onAborted);
    }

    // Call backend to start streaming
    window.go.main.App.SendMessageStream(message).catch((err) => {
        console.error('Error starting stream:', err);
        onError(err.message || err);
    });
}

// Update think area with content and timer (collapsible, hidden by default)
function updateThinkArea(thinkArea, content, startTime) {
    if (!content) {
        thinkArea.innerHTML = '';
        thinkArea.style.display = 'none';
        return;
    }

    const elapsed = startTime ? Date.now() - startTime : 0;
    const elapsedSeconds = (elapsed / 1000).toFixed(1);

    thinkArea.style.display = 'block';
    thinkArea.innerHTML = `
        <div class="thinking-toggle" onclick="toggleThink(this)">
            <span>Show reasoning (${elapsedSeconds}s)</span>
        </div>
        <div class="thinking-block">
            <pre>${content}</pre>
        </div>
    `;
}

// Toggle think area expansion
function toggleThink(toggleElement) {
    const thinkArea = toggleElement.closest('.chat-think');
    const block = thinkArea.querySelector('.thinking-block');
    const toggleText = toggleElement.querySelector('span');

    if (block.classList.contains('open')) {
        block.classList.remove('open');
        toggleText.textContent = toggleText.textContent.replace('Hide', 'Show');
    } else {
        block.classList.add('open');
        toggleText.textContent = toggleText.textContent.replace('Show', 'Hide');
    }
}

// Stream assistant response from backend with image
function streamAssistantResponseFromBackendWithImage(message, imageData) {
    const chatHistory = document.querySelector('.chat-messages');

    // Create assistant message bubble with typing indicator
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant streaming';

    const bubble = document.createElement('div');
    bubble.className = 'chat-bubble';
    bubble.innerHTML = `
        <div class="chat-sender">AI<span class="typing-indicator"><span></span><span></span><span></span></span></div>
        <div class="chat-text"></div>
        <div class="chat-think"></div>
        <div class="chat-bubble-buttons">
            <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
            <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
            <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
            <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
            <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
            <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
            <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
        </div>
    `;

    messageDiv.appendChild(bubble);
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;

    // Set current streaming bubble reference so onChunk can find the elements
    currentStreamingBubble = bubble;

    const textElement = bubble.querySelector('.chat-text');
    const thinkArea = bubble.querySelector('.chat-think');

    // Reset think state
    thinkContent = '';
    mainContent = '';
    isInThinkMode = false;
    thinkStartTime = null;

    const onStart = () => {
        console.log('Stream started with image');
        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Abort';
            sendButton.dataset.isAbort = 'true';
        }
    };

    const onChunk = (chunk) => {
        if (currentStreamingBubble) {
            // Check for start marker
            if (THINK_START_REGEX.test(chunk)) {
                isInThinkMode = true;
                thinkStartTime = Date.now();
                chunk = chunk.replace(THINK_START_REGEX, '');
            }

            // Check for end marker
            if (THINK_END_REGEX.test(chunk)) {
                isInThinkMode = false;
                chunk = chunk.replace(THINK_END_REGEX, '');
                if (thinkContent) {
                    updateThinkArea(thinkArea, thinkContent, thinkStartTime);
                }
            }

            // Route tokens to separate fields
            if (isInThinkMode) {
                thinkContent += chunk;
            } else {
                mainContent += chunk;
                textElement.textContent = mainContent;
            }
            chatHistory.scrollTop = chatHistory.scrollHeight;
        }
    };

    const onComplete = () => {
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) {
            indicator.remove();
        }
        messageDiv.classList.remove('streaming');

        if (thinkContent) {
            updateThinkArea(thinkArea, thinkContent, thinkStartTime);
        }

        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Send';
            sendButton.dataset.isAbort = 'false';
        }

        if (window.runtime) {
            window.runtime.EventsOff('chat:chunk');
            window.runtime.EventsOff('chat:complete');
            window.runtime.EventsOff('chat:error');
            window.runtime.EventsOff('chat:aborted');
        }
        currentStreamingBubble = null;
        updateNavigationButtons();
    };

    const onError = (err) => {
        textElement.textContent += '\n[Error: ' + err + ']';
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) {
            indicator.remove();
        }
        messageDiv.classList.remove('streaming');

        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Send';
            sendButton.dataset.isAbort = 'false';
        }

        currentStreamingBubble = null;
    };

    const onAborted = () => {
        console.log('Stream aborted');
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) {
            indicator.remove();
        }
        messageDiv.classList.remove('streaming');

        const sendButton = document.querySelector('.send-button');
        if (sendButton) {
            sendButton.textContent = 'Send';
            sendButton.dataset.isAbort = 'false';
        }

        if (window.runtime) {
            window.runtime.EventsOff('chat:chunk');
            window.runtime.EventsOff('chat:complete');
            window.runtime.EventsOff('chat:error');
            window.runtime.EventsOff('chat:aborted');
        }
        currentStreamingBubble = null;
    };

    if (window.runtime) {
        window.runtime.EventsOn('chat:start', onStart);
        window.runtime.EventsOn('chat:chunk', onChunk);
        window.runtime.EventsOn('chat:complete', onComplete);
        window.runtime.EventsOn('chat:error', onError);
        window.runtime.EventsOn('chat:aborted', onAborted);
    }

    // Call backend with image data
    if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.SendMessageStreamWithImage(message, imageData).catch((err) => {
            console.error('Error starting stream with image:', err);
            onError(err.message || err);
        });
    }
}
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
document.addEventListener('DOMContentLoaded', () => {
    const sendButton = document.querySelector('.send-button');
    if (sendButton) {
        sendButton.addEventListener('click', async () => {
            // Check if button is in abort mode
            if (sendButton.dataset.isAbort === 'true') {
                // Abort the current stream
                if (window.go && window.go.main && window.go.main.App) {
                    try {
                        await window.go.main.App.AbortStream();
                    } catch (error) {
                        console.error('Error aborting stream:', error);
                    }
                }
                return;
            }

            // Normal send mode
            const input = document.querySelector('.input-field');
            const message = input.value.trim();
            if (message) {
                addUserMessage(message);
                input.value = '';

                // Get image data if present
                let imageData = "";
                if (currentImage) {
                    // Extract base64 data from data URL (remove the data:image/...;base64, prefix)
                    imageData = currentImage.split(',')[1];
                    removeImage(); // Clear image after sending
                }

                // Check if Wails runtime is available
                if (!window.go || !window.go.main || !window.go.main.App) {
                    console.error('Wails runtime not available');
                    addAssistantMessage("AI", "Error: Backend not available. Please run this application through Wails (wails dev) to enable backend functionality.");
                    return;
                }

                // Start streaming response with or without image
                if (imageData) {
                    streamAssistantResponseFromBackendWithImage(message, imageData);
                } else {
                    streamAssistantResponseFromBackend(message);
                }
            }
        });
    }
});

// Handle Enter key in input field
document.addEventListener('DOMContentLoaded', () => {
    const inputField = document.querySelector('.input-field');
    if (inputField) {
        inputField.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                const sendButton = document.querySelector('.send-button');
                if (sendButton) {
                    sendButton.click();
                }
            }
        });
    }
});

// Check model status on page load
document.addEventListener('DOMContentLoaded', () => {
    checkModelStatus();
});

// Load Model - opens native OS file dialog via Wails
async function loadModel() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        alert('Error: Backend not available. Please run through Wails.');
        return;
    }

    try {
        // Use native OS file dialog via Wails backend binding
        const modelPath = await window.go.main.App.OpenModelFilePicker();
        if (!modelPath) {
            return; // user cancelled
        }

        console.log('Loading model:', modelPath);
        updateModelDisplay('Loading...');

        try {
            await window.go.main.App.LoadModel(modelPath);
            const modelName = modelPath.split(/[\\/]/).pop();
            updateModelDisplay(modelName);
            console.log('Model loaded successfully:', modelName);
        } catch (error) {
            console.error('Error loading model:', error);
            updateModelDisplay('Load failed');
            alert('Failed to load model: ' + (error.message || error));
        }
    } catch (error) {
        console.error('Error opening file picker:', error);
        alert('Failed to open file dialog: ' + (error.message || error));
    }
}

// Unload Model - aggressively unloads current model
async function unloadModel() {
    if (confirm('⚠️ AGGRESSIVE UNLOAD ⚠️\n\nThis will forcefully unload the current model from memory.\nAll unsaved context will be lost.\n\nAre you sure you want to proceed?')) {
        console.log('AGGRESSIVE UNLOAD initiated');
        try {
            await window.go.main.App.UnloadModel();
            updateModelDisplay("No model loaded");
            alert('🔴 MODEL AGGRESSIVELY UNLOADED\n\nGPU cache cleared.\nMemory freed.\nModel process terminated.');
        } catch (error) {
            console.error('Error unloading model:', error);
            alert(`Failed to unload model: ${error}`);
        }
    }
}

// Update model name display
function updateModelDisplay(modelName) {
    const display = document.getElementById('model-name-display');
    if (display) {
        display.textContent = modelName;
    }
}

// Check and update model display on page load
async function checkModelStatus() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.log('Wails runtime not available, skipping model status check');
        return;
    }

    try {
        const isLoaded = await window.go.main.App.IsModelLoaded();
        if (isLoaded) {
            const modelName = await window.go.main.App.GetLoadedModelName();
            updateModelDisplay(modelName);
        } else {
            updateModelDisplay("No model loaded");
        }
    } catch (error) {
        console.error('Error checking model status:', error);
        updateModelDisplay("No model loaded");
    }
}

// Regenerate AI response - removes last AI message and generates new one
async function regenerateMessage(button) {
    const messageDiv = button.closest('.chat-message');
    const chatHistory = document.querySelector('.chat-messages');

    // Find all assistant messages
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));

    if (assistantMessages.length === 0) {
        console.log('No AI messages to regenerate');
        return;
    }

    // Remove the last assistant message
    const lastAssistant = assistantMessages[assistantMessages.length - 1];
    lastAssistant.remove();

    // Re-stream the response using the last user message
    const userMessages = Array.from(chatHistory.querySelectorAll('.chat-message.user'));
    if (userMessages.length > 0) {
        const lastUserMessage = userMessages[userMessages.length - 1];
        const userText = lastUserMessage.querySelector('.chat-text').textContent;
        streamAssistantResponseFromBackend(userText);
    }

    updateNavigationButtons();
}

// Previous AI message - navigate to previous assistant message
function previousMessage(button) {
    const chatHistory = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));
    const currentIndex = assistantMessages.findIndex(msg => msg.contains(button.closest('.chat-bubble')));

    if (currentIndex > 0) {
        const previous = assistantMessages[currentIndex - 1];
        chatHistory.scrollTop = previous.offsetTop - chatHistory.offsetTop;
    }
}

// Next AI message - navigate to next assistant message
function nextMessage(button) {
    const chatHistory = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));
    const currentIndex = assistantMessages.findIndex(msg => msg.contains(button.closest('.chat-bubble')));

    if (currentIndex < assistantMessages.length - 1) {
        const next = assistantMessages[currentIndex + 1];
        chatHistory.scrollTop = next.offsetTop - chatHistory.offsetTop;
    }
}

// Continue from last AI response
async function continueMessage(button) {
    const chatHistory = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));

    if (assistantMessages.length === 0) {
        console.log('No AI messages to continue from');
        return;
    }

    const lastAssistant = assistantMessages[assistantMessages.length - 1];
    const lastText = lastAssistant.querySelector('.chat-text').textContent;

    // Stream continuation by appending to the last message
    streamAssistantResponseFromBackend(lastText);
}

// Edit message (works for both user and AI messages) - edits in-place within the bubble
function editMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const messageDiv = button.closest('.chat-message');
    const textElement = bubble.querySelector('.chat-text');
    const currentText = textElement.textContent;
    const isUserMessage = messageDiv.classList.contains('user');

    console.log('editMessage called, isUserMessage:', isUserMessage);

    // Replace text element with textarea for editing
    const textarea = document.createElement('textarea');
    textarea.className = 'chat-text chat-text-edit';
    textarea.value = currentText;
    textarea.rows = 4;

    // Replace the text element with textarea
    textElement.replaceWith(textarea);
    textarea.focus();

    // Add save/cancel buttons
    const buttonsContainer = bubble.querySelector('.chat-bubble-buttons');
    const originalButtons = buttonsContainer.innerHTML;

    // Create edit controls
    const editControls = document.createElement('div');
    editControls.className = 'chat-edit-controls';

    const saveButton = document.createElement('button');
    saveButton.className = 'bubble-button';
    saveButton.textContent = isUserMessage ? 'Confirm' : 'Save';
    saveButton.onclick = () => saveEdit(saveButton, isUserMessage);

    const cancelButton = document.createElement('button');
    cancelButton.className = 'bubble-button';
    cancelButton.textContent = 'Cancel';
    cancelButton.onclick = () => cancelEdit(cancelButton);

    editControls.appendChild(saveButton);
    editControls.appendChild(cancelButton);

    buttonsContainer.innerHTML = '';
    buttonsContainer.appendChild(editControls);

    // Store original buttons for cancel
    buttonsContainer.dataset.originalButtons = originalButtons;
}

// Save the edited message
async function saveEdit(button, isUserMessage) {
    const bubble = button.closest('.chat-bubble');
    const messageDiv = button.closest('.chat-message');
    const textarea = bubble.querySelector('.chat-text-edit');
    const editedText = textarea.value;

    console.log('saveEdit called, isUserMessage:', isUserMessage, 'editedText:', editedText);

    if (isUserMessage) {
        // For user messages: keep the edited message, remove only messages AFTER it
        const chatHistory = document.querySelector('.chat-messages');
        const messages = Array.from(chatHistory.querySelectorAll('.chat-message'));
        const messageIndex = messages.indexOf(messageDiv);

        console.log('User message edit, messageIndex:', messageIndex, 'total messages:', messages.length);

        // Update the text of the message being edited
        const newTextElement = document.createElement('div');
        newTextElement.className = 'chat-text';
        newTextElement.textContent = editedText;
        textarea.replaceWith(newTextElement);

        // Remove only messages AFTER the edited message (from end to start to avoid index issues)
        for (let i = messages.length - 1; i > messageIndex; i--) {
            messages[i].remove();
        }

        // Clear the input field to prevent old text from being sent with new messages
        const input = document.querySelector('.input-field');
        input.value = '';

        // Clear backend chat history to match frontend state
        if (window.go && window.go.main && window.go.main.App) {
            try {
                await window.go.main.App.ClearChatHistory();
                console.log('Backend chat history cleared');
            } catch (error) {
                console.error('Error clearing chat history:', error);
            }
        }

        // Send the edited message to generate new response
        streamAssistantResponseFromBackend(editedText);
    } else {
        // For AI messages: just update the text in place
        const newTextElement = document.createElement('div');
        newTextElement.className = 'chat-text';
        newTextElement.textContent = editedText;
        textarea.replaceWith(newTextElement);
    }

    // Restore original buttons
    const buttonsContainer = bubble.querySelector('.chat-bubble-buttons');
    buttonsContainer.innerHTML = buttonsContainer.dataset.originalButtons;

    updateNavigationButtons();
}

// Cancel editing and restore original
function cancelEdit(button) {
    const bubble = button.closest('.chat-bubble');
    const textarea = bubble.querySelector('.chat-text-edit');
    const originalText = textarea.dataset.originalText || textarea.value;

    // Replace textarea with original text element
    const newTextElement = document.createElement('div');
    newTextElement.className = 'chat-text';
    newTextElement.textContent = originalText;
    textarea.replaceWith(newTextElement);

    // Restore original buttons
    const buttonsContainer = bubble.querySelector('.chat-bubble-buttons');
    buttonsContainer.innerHTML = buttonsContainer.dataset.originalButtons;
}

// Delete message and show banner if AI message
function deleteMessage(button) {
    const messageDiv = button.closest('.chat-message');
    const isAssistant = messageDiv.classList.contains('assistant');

    messageDiv.remove();

    if (isAssistant) {
        showGenerateBanner();
    }

    updateNavigationButtons();
}

// Update navigation buttons (gray out prev/next when no messages exist)
function updateNavigationButtons() {
    const chatHistory = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));

    const prevButtons = chatHistory.querySelectorAll('#prev-btn');
    const nextButtons = chatHistory.querySelectorAll('#next-btn');

    const hasPrevious = assistantMessages.length > 1;
    const hasNext = false; // Simplified - could implement full history navigation

    prevButtons.forEach(btn => {
        btn.classList.toggle('disabled', !hasPrevious);
    });

    nextButtons.forEach(btn => {
        btn.classList.toggle('disabled', !hasNext);
    });
}

// Show "generate ai response" banner at bottom of chat
function showGenerateBanner() {
    const chatPanel = document.querySelector('.chat-panel');
    let banner = document.getElementById('generate-banner');

    if (!banner) {
        banner = document.createElement('div');
        banner.id = 'generate-banner';
        banner.className = 'generate-banner';
        banner.innerHTML = 'Generate AI Response';
        banner.onclick = () => {
            const input = document.querySelector('.input-field');
            const userMessages = Array.from(document.querySelectorAll('.chat-message.user'));
            if (userMessages.length > 0) {
                const lastUserMessage = userMessages[userMessages.length - 1];
                const userText = lastUserMessage.querySelector('.chat-text').textContent;
                streamAssistantResponseFromBackend(userText);
            }
            banner.remove();
        };
        chatPanel.appendChild(banner);
    }
}

// New chat - clears all messages and history
async function newChat() {
    // Clear frontend chat messages
    const chatHistory = document.querySelector('.chat-messages');
    chatHistory.innerHTML = '';

    // Clear input field
    const input = document.querySelector('.input-field');
    if (input) {
        input.value = '';
    }

    // Remove generate banner if present
    const banner = document.getElementById('generate-banner');
    if (banner) {
        banner.remove();
    }

    // Clear backend chat history
    if (window.go && window.go.main && window.go.main.App) {
        try {
            await window.go.main.App.ClearChatHistory();
            console.log('Backend chat history cleared for new chat');
        } catch (error) {
            console.error('Error clearing chat history:', error);
        }
    }
}

let currentImage = null; // Store current image data

function addImage() {
    // Create hidden file input
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.style.display = 'none';
    document.body.appendChild(input);

    input.onchange = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                currentImage = e.target.result; // Store base64 data
                // Show image preview in input bar
                showImagePreview(currentImage);
            };
            reader.readAsDataURL(file);
        }
        document.body.removeChild(input);
    };

    input.click();
}

function showImagePreview(imageData) {
    // Remove existing preview if any
    const existingPreview = document.querySelector('.image-preview');
    if (existingPreview) {
        existingPreview.remove();
    }

    // Create preview element
    const preview = document.createElement('div');
    preview.className = 'image-preview';
    preview.innerHTML = `
        <img src="${imageData}" alt="Preview">
        <button class="remove-image" onclick="removeImage()">×</button>
    `;

    // Insert before input bar
    const inputBar = document.querySelector('.input-bar');
    inputBar.parentNode.insertBefore(preview, inputBar);
}

function removeImage() {
    currentImage = null;
    const preview = document.querySelector('.image-preview');
    if (preview) {
        preview.remove();
    }
}

function addAudio() {
    console.log('Add audio');
    // TODO: Implement audio recording
}

let ttsAudio = null;

function textToSpeech() {
    // Create a file input for audio files
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = 'audio/*';
    fileInput.style.display = 'none';
    
    fileInput.onchange = function(e) {
        const file = e.target.files[0];
        if (!file) return;
        
        // Create object URL and play
        const audioUrl = URL.createObjectURL(file);
        
        // Stop any currently playing TTS audio
        if (ttsAudio) {
            ttsAudio.pause();
            ttsAudio = null;
        }
        
        ttsAudio = new Audio(audioUrl);
        
        // Visual feedback
        const ttsButton = document.querySelector('.tts-icon');
        if (ttsButton) {
            ttsButton.classList.add('active');
        }
        
        ttsAudio.onended = function() {
            URL.revokeObjectURL(audioUrl);
            if (ttsButton) {
                ttsButton.classList.remove('active');
            }
        };
        
        ttsAudio.onerror = function() {
            console.error('Error playing TTS audio');
            URL.revokeObjectURL(audioUrl);
            if (ttsButton) {
                ttsButton.classList.remove('active');
            }
        };
        
        ttsAudio.play().catch(err => {
            console.error('Failed to play TTS audio:', err);
            if (ttsButton) {
                ttsButton.classList.remove('active');
            }
        });
    };
    
    document.body.appendChild(fileInput);
    fileInput.click();
    document.body.removeChild(fileInput);
}

function addFile() {
    console.log('Add file');
    // TODO: Implement document upload for RAG 2
}
