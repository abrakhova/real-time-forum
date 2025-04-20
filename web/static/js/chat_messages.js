// chat-messages.js - Chat message handling functionality

// Global variables for messages
const onlineUsers = [];
const usersWithUnreadMessages = new Set(); // Store user IDs who have unread messages
const lastMessagesData = {}; // Store last messages data by user ID

// Handle user list updates from the server
function handleUserList(users) {
  onlineUsers.length = 0;
  
  if (Array.isArray(users)) {
    users.forEach((user) => {
      onlineUsers.push(user.id);

      // Store last messages data if available
      if (user.lastMessages) {
        lastMessagesData[user.id] = user.lastMessages;
      }
    });
  }

  // Update the UI immediately
  if (window.chatUI && window.chatUI.updateUsersList) {
    window.chatUI.updateUsersList();
  }
}

// Handle incoming messages
function handleMessage(message) {
  displayMessage(message);

  // If a message is from someone else, mark it as unread
  if (message.sender_id !== window.state.sessionID) {
    // Only mark as unread if we're not currently viewing that user's chat
    let currentChatUser = null;
    if (window.chatUI && window.chatUI.currentUser) {
      currentChatUser = window.chatUI.currentUser();
    }

    if (!currentChatUser || message.sender_id !== currentChatUser.id) {
      usersWithUnreadMessages.add(message.sender_id);
      showMessageNotification(message);
    }
  }
}

// Update the displayMessage function in chat_messages.js
function displayMessage(message) {
    let currentChatUser = null;
    if (window.chatUI && window.chatUI.currentUser) {
      currentChatUser = window.chatUI.currentUser();
    }
  
    // Only display if it's part of the current chat
    if (
      currentChatUser &&
      ((message.sender_id === currentChatUser.id &&
        message.receiverID === window.state.sessionID) ||
        (message.sender_id === window.state.sessionID &&
          message.receiverID === currentChatUser.id))
    ) {
      const messagesContainer = document.getElementById("messages-container");
      if (!messagesContainer) {
        return;
      }
  
      // Clear empty state if it exists
      const emptyState = messagesContainer.querySelector(".chat-empty-state");
      if (emptyState) {
        emptyState.remove();
      }
  
      // Find existing "sending..." message if this is a confirmation of our sent message
      if (message.sender_id === window.state.sessionID) {
        const pendingMessages =
          messagesContainer.querySelectorAll(".message.outgoing");
        for (const pending of pendingMessages) {
          const timeElem = pending.querySelector(".message-time");
          if (timeElem && timeElem.textContent.includes("Sending...")) {
            // Update this message instead of creating a new one
            const messageDate = new Date(message.timestamp);
            const dateStr = messageDate.toLocaleDateString();
            const timeStr = messageDate.toLocaleTimeString();
            timeElem.textContent = `${dateStr} ${timeStr}`;
            // Add timestamp as data attribute
            timeElem.setAttribute("data-timestamp", message.timestamp);
  
            // If we have new message data, update the sorting of users
            if (!lastMessagesData[message.receiverID]) {
              lastMessagesData[message.receiverID] = {};
            }
            lastMessagesData[message.receiverID][message.id || Date.now()] = {
              timestamp: message.timestamp,
              content: message.content,
              sender_id: message.sender_id,
            };
  
            if (window.chatUI && window.chatUI.updateUsersList) {
              window.chatUI.updateUsersList();
            }
  
            return;
          }
        }
      }
  
      // Create new message element
      const messageElem = document.createElement("div");
      messageElem.className = "message";
  
      // Get sender name
      let senderName = window.state.username;
      if (message.sender_id !== window.state.sessionID) {
        senderName = message.username || currentChatUser.name;
      }
  
      if (message.sender_id === window.state.sessionID) {
        messageElem.classList.add("outgoing");
      } else {
        messageElem.classList.add("incoming");
      }
  
      const messageDate = new Date(message.timestamp);
      const dateStr = messageDate.toLocaleDateString();
      const timeStr = messageDate.toLocaleTimeString();
      const escapeHTML =
        window.chatUI && window.chatUI.escapeHTML
          ? window.chatUI.escapeHTML
          : (text) =>
              text
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;");
  
      messageElem.innerHTML = `
              <div class="message-sender">${escapeHTML(senderName)}</div>
              <div class="message-text">${escapeHTML(message.content)}</div>
              <div class="message-time" data-timestamp="${
                message.timestamp
              }">${dateStr} ${timeStr}</div>
          `;
  
      messagesContainer.appendChild(messageElem);
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
  
      // Update last messages data for sorting
      const otherUser =
        message.sender_id === window.state.sessionID
          ? message.receiverID
          : message.sender_id;
      if (!lastMessagesData[otherUser]) {
        lastMessagesData[otherUser] = {};
      }
      lastMessagesData[otherUser][message.id || Date.now()] = {
        timestamp: message.timestamp,
        content: message.content,
        sender_id: message.sender_id,
      };
  
      if (window.chatUI && window.chatUI.updateUsersList) {
        window.chatUI.updateUsersList();
      }
    } else {
      // Still update the last messages data for user sorting
      if (message.sender_id !== window.state.sessionID) {
        const senderID = message.sender_id;
        if (!lastMessagesData[senderID]) {
          lastMessagesData[senderID] = {};
        }
        lastMessagesData[senderID][message.id || Date.now()] = {
          timestamp: message.timestamp,
          content: message.content,
          sender_id: message.sender_id,
        };
  
        if (window.chatUI && window.chatUI.updateUsersList) {
          window.chatUI.updateUsersList();
        }
      }
    }
  }

// Show notification for new message
function showMessageNotification(message) {
  // Find the username from users list
  let senderName = "User";
  const userItem = document.querySelector(
    `.user-item[data-user-id="${message.sender_id}"]`
  );

  if (userItem) {
    senderName = userItem.dataset.username;

    // Add notification indicator to this user
    userItem.classList.add("has-new-message");
  }

  // Create notification if browser supports it
  if ("Notification" in window) {
    if (Notification.permission === "granted") {
      new Notification(`New message from ${senderName}`, {
        body:
          message.content.substring(0, 50) +
          (message.content.length > 50 ? "..." : ""),
      });
    } else if (Notification.permission !== "denied") {
      Notification.requestPermission();
    }
  }
}

// Send a message
function sendMessage() {
  const socket = window.chatConnection ? window.chatConnection.socket() : null;
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    alert("WebSocket not connected. Please refresh the page.");
    return;
  }

  let currentChatUser = null;
  if (window.chatUI && window.chatUI.currentUser) {
    currentChatUser = window.chatUI.currentUser();
  }

  if (!currentChatUser) {
    alert("Please select a user to chat with");
    return;
  }

  const messageInput = document.getElementById("message-input");
  if (!messageInput) {
    return;
  }

  const content = messageInput.value.trim();

  if (!content) {
    return;
  }

  // Make sure receiverID is explicitly an integer
  const receiverID = parseInt(currentChatUser.id, 10);

  const message = {
    type: "message",
    receiverID: receiverID,
    content: content,
  };

  // Add message locally immediately for better UX
  if (window.chatUI && window.chatUI.displayLocalMessage) {
    window.chatUI.displayLocalMessage(content);
  }

  // Send via websocket
  try {
    socket.send(JSON.stringify(message));
    messageInput.value = "";
    messageInput.focus();
  } catch (e) {
    console.log("Error sending message:", e);
    alert("Failed to send message: " + e.message);
  }
}

// Check if user has unread messages
function hasUnreadMessages(userId) {
  return usersWithUnreadMessages.has(userId);
}

// Clear unread messages for a user
function clearUnreadMessages(userId) {
  usersWithUnreadMessages.delete(userId);
}

// Check if user is online
function isUserOnline(userId) {
  return onlineUsers.includes(userId);
}

// Get last messages data for a user
function getLastMessagesData(userId) {
  return lastMessagesData[userId] || null;
}

// Export the chat messages module functions
window.chatMessages = {
  handleUserList,
  handleMessage,
  displayMessage,
  showMessageNotification,
  sendMessage,
  hasUnreadMessages,
  clearUnreadMessages,
  isUserOnline,
  getLastMessagesData,
};

// Initialize when DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  // Check if we need to connect WebSocket
  if (window.state && window.state.sessionID > 0) {
    setTimeout(() => {
      if (window.chatConnection && window.chatConnection.connect) {
        window.chatConnection.connect();
      }
    }, 1000);
  }
});

// Only start the interval if user is already logged in
if (window.state && window.state.sessionID > 0) {
  if (window.chatConnection && window.chatConnection.startChecking) {
    window.chatConnection.startChecking();
  }
}
