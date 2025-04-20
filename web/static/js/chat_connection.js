// chat-connection.js - WebSocket connection management for chat

// Global variables for WebSocket connection
let socket = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 1000; // 1 second delay
let wsCheckInterval = null;

// Connect to WebSocket when user is logged in
function connect() {
  // Only connect if user is logged in and no active connection exists
  if (
    !window.state ||
    !window.state.sessionID ||
    (socket !== null && socket.readyState === WebSocket.OPEN)
  ) {
    return;
  }

  // Stop trying if max reconnect attempts reached
  if (reconnectAttempts >= maxReconnectAttempts) {
    const statusElement = document.getElementById("chat-status");
    if (statusElement) {
      statusElement.textContent = "Disconnected";
      statusElement.className = "disconnected";
    }
    return;
  }

  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsUrl = `${protocol}//${window.location.host}/ws`;

  // Update connection status indicator
  const statusElement = document.getElementById("chat-status");
  if (statusElement) {
    statusElement.textContent = "Connected";
    statusElement.className = "connected";
  }

  try {
    socket = new WebSocket(wsUrl);

    socket.onopen = function () {
      if (statusElement) {
        statusElement.textContent = "Connected";
        statusElement.className = "connected";
      }
      reconnectAttempts = 0; // Reset attempts on successful connection
      if (window.chatUI && window.chatUI.fetchAllUsers) {
        window.chatUI.fetchAllUsers();
      }
    };

    socket.onclose = function () {
      if (statusElement) {
        statusElement.textContent = "Disconnected";
        statusElement.className = "disconnected";
      }
      socket = null;

      // Try to reconnect if user is still logged in
      if (window.state && window.state.sessionID > 0) {
        reconnectAttempts++;
        const delay = reconnectDelay * Math.pow(2, reconnectAttempts - 1);
        setTimeout(connect, delay);
      }
    };

    socket.onerror = function () {
      if (statusElement) {
        statusElement.textContent = "Error";
        statusElement.className = "disconnected";
      }
    };

    socket.onmessage = function (event) {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "user_list" && window.chatMessages) {
          window.chatMessages.handleUserList(data.users);
        } else if (data.type === "message" && window.chatMessages) {
          window.chatMessages.handleMessage(data);
        } else if (data.type === "typing") {
          // Show typing indicator
          const typingIndicator = document.getElementById("typing-indicator");
          if (typingIndicator) {
            typingIndicator.innerHTML = `${
              data.username || "User"
            } is typing<span class="typing-dots"><span>.</span><span>.</span><span>.</span></span>`;
            typingIndicator.style.display = "block";
          }
        } else if (data.type === "typing_stopped") {
          // Hide typing indicator
          const typingIndicator = document.getElementById("typing-indicator");
          if (typingIndicator) {
            typingIndicator.style.display = "none";
          }
        } else if (data.type === "history" && window.chatUI) {
          if (Array.isArray(data.messages)) {
            window.chatUI.displayMessageHistory(data.messages);
          } else {
            window.chatUI.displayMessageHistory([]);
          }
        } else if (data.type === "more_history" && window.chatUI) {
          if (Array.isArray(data.messages)) {
            window.chatUI.displayMoreMessageHistory(data.messages);
          }
        }
      } catch (e) {
        console.log("Error processing WebSocket message:", e);
      }
    };
  } catch (e) {
    console.log("Error creating WebSocket:", e);
    if (statusElement) {
      statusElement.textContent = "Failed";
      statusElement.className = "disconnected";
    }

    // Try to reconnect if user is still logged in
    if (window.state && window.state.sessionID > 0) {
      reconnectAttempts++;
      const delay = reconnectDelay * Math.pow(2, reconnectAttempts - 1);
      setTimeout(connect, delay);
    }
  }
}

// Check WebSocket connection status
function checkAndConnectWebSocket() {
  // If user isn't logged in, stop checking and clean up
  if (!window.state || !window.state.sessionID || window.state.sessionID <= 0) {
    stopChecking();

    // Reset connection
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.close();
    }
    socket = null;
    reconnectAttempts = 0;

    // Update status indicator
    const statusElement = document.getElementById("chat-status");
    if (statusElement) {
      statusElement.textContent = "Not logged in";
      statusElement.className = "disconnected";
    }
    return;
  }

  // Only attempt to connect if needed
  if (
    (!socket || socket.readyState !== WebSocket.OPEN) &&
    reconnectAttempts < maxReconnectAttempts
  ) {
    connect();
  }
}

// Function to start WebSocket checking interval
function startChecking() {
  // Clear any existing interval first
  if (wsCheckInterval) {
    clearInterval(wsCheckInterval);
  }

  // Only set interval if user is logged in
  if (window.state && window.state.sessionID > 0) {
    wsCheckInterval = setInterval(checkAndConnectWebSocket, 10000);
  }
}

// Function to stop WebSocket checking interval
function stopChecking() {
  if (wsCheckInterval) {
    clearInterval(wsCheckInterval);
    wsCheckInterval = null;
  }

  // Also close any existing connection
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close();
    socket = null;
  }
}

// Check if WebSocket is connected
function isConnected() {
  return socket && socket.readyState === WebSocket.OPEN;
}

// Get the socket (used by other modules)
function getSocket() {
  return socket;
}

// Add cleanup for when page unloads
window.addEventListener("beforeunload", function () {
  stopChecking();
});

// Export chat connection module functions
window.chatConnection = {
  socket: getSocket,
  connect,
  checkAndConnectWebSocket,
  startChecking,
  stopChecking,
  isConnected,
};
