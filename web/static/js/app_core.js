// app-core.js - Core application functionality for Forum SPA

// Global application state
const state = {
  sessionID: null,
  username: null,
  posts: [],
  currentPost: null,
};

// Make state available globally
window.state = state;

// Initialize the application
document.addEventListener("DOMContentLoaded", async function () {
  // When state updates, check WebSocket connection
  window.addEventListener("stateUpdated", function () {
    if (
      window.chatConnection &&
      window.chatConnection.checkAndConnectWebSocket
    ) {
      window.chatConnection.checkAndConnectWebSocket();
    }
  });

  // Check login status and set up page
  await checkLogin();
  setupNavigationEvents();
  loadCurrentPage();
});

// Check if user is logged in
async function checkLogin() {
  try {
    const response = await fetch("/user/status");
    const data = await response.json();

    // Check if login status changed
    const wasLoggedIn = state.sessionID > 0;
    const isLoggedIn = data.sessionID > 0;

    // Update state
    state.sessionID = data.sessionID;
    state.username = data.username;
    updateUI();

    // Handle WebSocket connection based on login status change
    if (!wasLoggedIn && isLoggedIn) {
      // User just logged in - start WebSocket
      if (window.chatConnection) {
        window.chatConnection.connect();
        window.chatConnection.startChecking();
      }
    } else if (wasLoggedIn && !isLoggedIn) {
      // User just logged out - stop WebSocket
      if (window.chatConnection) {
        window.chatConnection.stopChecking();
      }
    }
  } catch (error) {
    console.error("Login check failed:", error);
  }
}

// Update UI elements after login status change
function updateUI() {
  // Update auth box
  document.getElementById("auth-box").innerHTML = window.templates.authBox(
    state.sessionID,
    state.username
  );

  // Add logout event listener if logged in
  if (state.sessionID) {
    document
      .getElementById("logout-link")
      .addEventListener("click", handleLogout);
  }

  // Update sidebar
  document.getElementById("sidebar").innerHTML = window.templates.sidebar(
    state.sessionID
  );

  // Update chat sidebar
  const chatSidebar = document.getElementById("chat-sidebar");
  if (chatSidebar) {
    if (state.sessionID) {
      chatSidebar.innerHTML = window.templates.chatSidebar();

      // Fetch users when UI is updated
      if (window.chatUI && window.chatUI.fetchAllUsers) {
        setTimeout(window.chatUI.fetchAllUsers, 100);
      }
    } else {
      chatSidebar.innerHTML = ""; // Empty for non-logged-in users
    }
  }
}

// Handle logout
async function handleLogout(event) {
  event.preventDefault();

  try {
    const response = await fetch("/logout", { method: "POST" });
    if (response.ok) {
      state.sessionID = null;
      state.username = null;
      updateUI();
      navigate("/");
    }
  } catch (error) {
    console.error("Logout failed:", error);
  }
}

// Set up navigation event listeners
function setupNavigationEvents() {
  // Handle link clicks with data-navigate attribute
  document.addEventListener("click", function (event) {
    const link = event.target.closest("[data-navigate]");
    if (link) {
      event.preventDefault();
      currentChatUser = null;
      navigate(link.getAttribute("href"));
    }
  });

  // Handle browser back/forward
  window.addEventListener("popstate", loadCurrentPage);
}

// Navigate to a new page
function navigate(url) {
  history.pushState(null, "", url);
  loadCurrentPage();
}

// Determine what to show based on current URL
function loadCurrentPage() {
  const path = window.location.pathname;
  const searchParams = new URLSearchParams(window.location.search);
  const content = document.getElementById("content");

  // Show loading indicator
  content.innerHTML = window.templates.loading();

  // Check if user is logged in first - redirect to login for protected pages
  if (
    !state.sessionID &&
    (path === "/" ||
      path === "/index.html" ||
      path === "/filter" ||
      path === "/post")
  ) {
    if (window.appPages) {
      window.appPages.showLoginPage();
    } else {
      // Fallback if appPages is not yet loaded
      console.error("appPages module not loaded");
      content.innerHTML =
        "<div>Error: Application modules not loaded correctly.</div>";
    }
    return;
  }

  // Route to correct page handler
  if (path === "/" || path === "/index.html") {
    if (window.appPages) window.appPages.loadHomePage();
  } else if (path === "/post") {
    const postId = searchParams.get("id");
    if (postId && postId !== "undefined") {
      if (window.appPages) window.appPages.loadPostPage(postId);
    } else {
      if (window.appPages)
        window.appPages.showErrorPage("Post ID is missing or invalid");
    }
  } else if (path === "/filter") {
    if (window.appPages)
      window.appPages.loadFilteredPosts(window.location.search);
  } else if (path === "/login") {
    if (window.appPages) window.appPages.showLoginPage();
  } else if (path === "/register") {
    if (window.appPages) window.appPages.showRegisterPage();
  } else if (path === "/createPost") {
    if (state.sessionID) {
      if (window.appPages) window.appPages.showCreatePostPage();
    } else {
      navigate("/login");
    }
  } else {
    if (window.appPages) window.appPages.showErrorPage("Page not found");
  }

  // Update chat sidebar if needed
  const chatSidebar = document.getElementById("chat-sidebar");
  if (chatSidebar && state.sessionID) {
    chatSidebar.innerHTML = window.templates.chatSidebar();
    if (window.chatUI && window.chatUI.fetchAllUsers) {
      setTimeout(window.chatUI.fetchAllUsers, 100);
    }
  }
}

// Export core functions
window.appCore = {
  state,
  checkLogin,
  updateUI,
  navigate,
  handleLogout,
  setupNavigationEvents,
  loadCurrentPage,
};
