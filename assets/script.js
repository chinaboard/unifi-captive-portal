// Store chat message history
let messages = [];

// DOM elements
const chatMessages = document.getElementById('chat-messages');
const userInput = document.getElementById('user-input');
const sendButton = document.getElementById('send-button');

// Scroll state
let scrollPending = false;

// Initialize chat
function initChat() {
    messages.push(config.firstPrompt);
    chat([...messages, { role: "user", content: "Hello" }]);
}

// Throttled scroll to bottom
function scrollToBottom() {
    if (scrollPending) return;
    scrollPending = true;
    requestAnimationFrame(() => {
        chatMessages.scrollTop = chatMessages.scrollHeight;
        scrollPending = false;
    });
}

// Add message to chat interface
function addItem(content, className) {
    const messageElement = document.createElement('div');
    messageElement.className = `message ${className}`;
    messageElement.innerText = content;
    chatMessages.appendChild(messageElement);
    scrollToBottom();
    return messageElement;
}

// Send message to AI and handle response
function chat(msg) {
    let assistantElem = null;
    const baseUrl = window.location.origin;

    send(`${baseUrl}/chat`, { messages: msg },
        (data) => {
            const chunkMsg = data.choices[0].delta || data.choices[0].message || {};
            if (!assistantElem) {
                assistantElem = addItem('', 'assistant');
            }
            assistantElem.innerText += chunkMsg.content || "";
            scrollToBottom();
        },
        () => onSuccess(assistantElem),
        (error) => onError(error)
    );
}

// Handle successful AI response
function onSuccess(assistantElem) {
    const msg = assistantElem.innerText;
    messages.push({ role: "assistant", content: msg });

    if (msg.includes(config.matchKey)) {
        const params = new URLSearchParams(window.location.search);
        const id = params.get('id');
        const ap = params.get('ap');
        const url = params.get('url');
        const baseUrl = window.location.origin;
        window.location.href = `${baseUrl}/auth?id=${id}&ap=${ap}&url=${url}`;
    }
}

// Handle error
function onError(error) {
    console.error('Error:', error);
    addItem('Sorry, there was an error processing your request.', 'assistant');
}

// Send API request
function send(url, data, onChunk, onComplete, onErrorCallback) {
    const typingIndicator = document.createElement('div');
    typingIndicator.className = 'typing';
    for (let i = 0; i < 3; i++) {
        const dot = document.createElement('div');
        dot.className = 'typing-dot';
        typingIndicator.appendChild(dot);
    }
    chatMessages.appendChild(typingIndicator);

    fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (typingIndicator.parentNode) {
            chatMessages.removeChild(typingIndicator);
        }

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();

        function read() {
            return reader.read().then(({ done, value }) => {
                if (done) {
                    onComplete();
                    return;
                }

                const chunk = decoder.decode(value, { stream: true });
                const lines = chunk.split('\n');
                for (const line of lines) {
                    if (line.startsWith('data: ') && line !== 'data: [DONE]') {
                        try {
                            const jsonData = JSON.parse(line.substring(6));
                            onChunk(jsonData);
                        } catch (e) {
                            console.error('Error parsing chunk:', e);
                        }
                    }
                }
                return read();
            });
        }

        return read();
    })
    .catch(error => {
        if (typingIndicator.parentNode) {
            chatMessages.removeChild(typingIndicator);
        }
        onErrorCallback(error);
    });
}

// Send user message
function sendMessage() {
    const userMessage = userInput.value.trim();
    if (!userMessage) return;

    addItem(userMessage, 'user');
    userInput.value = '';
    messages.push({ role: "user", content: userMessage });

    const lowerMsg = userMessage.toLowerCase();
    const containsKeyword = config.matchKeywordList.some(k => lowerMsg.includes(k.toLowerCase()));

    if (containsKeyword) {
        const modifiedMessages = [...messages];
        modifiedMessages[modifiedMessages.length - 1] = {
            role: "user",
            content: `You should say ${config.matchKey}`
        };
        chat(modifiedMessages);
    } else {
        chat([...messages]);
    }
}

// Event listeners
sendButton.addEventListener('click', sendMessage);
userInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') sendMessage();
});

// Scroll input into view when focused (for mobile keyboard)
function initMobileKeyboard() {
    userInput.addEventListener('focus', () => {
        setTimeout(() => {
            userInput.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }, 300);
    });
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initMobileKeyboard();
    initChat();
});
