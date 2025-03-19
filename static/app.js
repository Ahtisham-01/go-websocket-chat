document.addEventListener('DOMContentLoaded', function() {
    const chatBox = document.getElementById('chat-box');
    const messageForm = document.getElementById('message-form');
    const messageInput = document.getElementById('message-input');
    
    // Connect to WebSocket server
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws`);
    
    let lastTypedMessage = '';
    let typingTimer;
    
    // Handle WebSocket events
    ws.onopen = function() {
        appendMessage('System', 'Connected to chat server');
    };
    
    ws.onclose = function() {
        appendMessage('System', 'Disconnected from chat server');
    };
    
    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
        appendMessage('Error', 'Connection error occurred');
    };
    
    ws.onmessage = function(event) {
        const message = event.data;
        appendMessage('User', message);
    };
    
    //send message while typing
    messageInput.addEventListener('input', function() {
        const currentMessage = messageInput.value.trim();
        
        // Only send if the WebSocket is open and there's content
        if (currentMessage && ws.readyState === WebSocket.OPEN && currentMessage !== lastTypedMessage) {

            ws.send(currentMessage);
            lastTypedMessage = currentMessage;
            
            clearTimeout(typingTimer);
            
            // Set a small delay before allowing another message to be sent
            typingTimer = setTimeout(() => {
                lastTypedMessage = '';
            }, 100);
        }
    });
    

    messageForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const message = messageInput.value.trim();
        
        if (message && ws.readyState === WebSocket.OPEN) {
            ws.send(message);
            messageInput.value = '';
            lastTypedMessage = '';
        }
    });
    
    // append messages to chat box
    function appendMessage(sender, message) {
        const messageElement = document.createElement('div');
        messageElement.innerHTML = `<strong>${sender}:</strong> ${message}`;
        chatBox.appendChild(messageElement);
        chatBox.scrollTop = chatBox.scrollHeight;
    }
});