// app-pages.js - Page loading functions for Forum SPA

// Load the home page with all posts
async function loadHomePage() {
  try {
    const response = await fetch("/?api=true", {
      headers: {
        Accept: "application/json",
        "X-Requested-With": "XMLHttpRequest",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to load posts");
    }

    const data = await response.json();
    window.state.posts = data.Posts || [];

    document.getElementById("content").innerHTML = window.templates.homePage(
      window.state.posts
    );
    setupReactionButtons();
  } catch (error) {
    console.error("Error:", error);
    document.getElementById("content").innerHTML = window.templates.error(
      "Failed to load posts"
    );
  }
}

// Load a single post page
async function loadPostPage(postId) {
  try {
    const response = await fetch(`/post?id=${postId}`, {
      headers: { Accept: "application/json" },
    });

    if (!response.ok) {
      throw new Error("Failed to load post");
    }

    const data = await response.json();
    window.state.currentPost = data.post || data.Post;

    document.getElementById("content").innerHTML = window.templates.postDetail(
      window.state.currentPost
    );

    // Setup reaction buttons
    setupReactionButtons();

    // Setup comment form submission
    const commentForm = document.getElementById("comment-form");
    if (commentForm) {
      commentForm.addEventListener("submit", function (event) {
        if (window.appForms && window.appForms.submitComment) {
          window.appForms.submitComment(event);
        } else {
          console.error("appForms.submitComment not loaded");
          event.preventDefault();
        }
      });
    }
  } catch (error) {
    console.error("Error:", error);
    document.getElementById("content").innerHTML = window.templates.error(
      "Failed to load post"
    );
  }
}

// Load filtered posts
async function loadFilteredPosts(queryString) {
  try {
    const response = await fetch(`/filter${queryString}`, {
      headers: { Accept: "application/json" },
    });

    if (!response.ok) {
      throw new Error("Failed to load filtered posts");
    }

    const data = await response.json();
    window.state.posts = data.posts || [];

    // Determine filter title
    let title = "Filtered Posts";
    if (data.category) {
      title = `${data.category} Posts`;
    } else if (queryString.includes("user_created=true")) {
      title = "My Posts";
    } else if (queryString.includes("liked=true")) {
      title = "Liked Posts";
    }

    document.getElementById("content").innerHTML =
      window.templates.filteredPosts(title, window.state.posts);
  } catch (error) {
    console.error("Error:", error);
    document.getElementById("content").innerHTML = window.templates.error(
      "Failed to load filtered posts"
    );
  }
}

// Show login page
function showLoginPage() {
  document.getElementById("content").innerHTML = window.templates.loginForm();

  // Make sure we properly attach the event handler using the appForms module
  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    loginForm.addEventListener("submit", function (event) {
      if (window.appForms && window.appForms.submitLogin) {
        window.appForms.submitLogin(event);
      } else {
        console.error("appForms.submitLogin not loaded");
        event.preventDefault();
        alert(
          "Login functionality is not available. Please try refreshing the page."
        );
      }
    });
  }
}

// Show register page
function showRegisterPage() {
  document.getElementById("content").innerHTML =
    window.templates.registerForm();

  const registerForm = document.getElementById("register-form");
  if (registerForm) {
    registerForm.addEventListener("submit", function (event) {
      if (window.appForms && window.appForms.submitRegister) {
        window.appForms.submitRegister(event);
      } else {
        console.error("appForms.submitRegister not loaded");
        event.preventDefault();
      }
    });
  }
}

// Show create post page
function showCreatePostPage() {
  document.getElementById("content").innerHTML =
    window.templates.createPostForm();

  const createPostForm = document.getElementById("create-post-form");
  if (createPostForm) {
    createPostForm.addEventListener("submit", function (event) {
      if (window.appForms && window.appForms.submitPost) {
        window.appForms.submitPost(event);
      } else {
        console.error("appForms.submitPost not loaded");
        event.preventDefault();
      }
    });
  }
}

// Show error page
function showErrorPage(message) {
  document.getElementById("content").innerHTML =
    window.templates.error(message);
}

// Set up reaction buttons (likes/dislikes)
function setupReactionButtons() {
  const reactionButtons = document.querySelectorAll(
    ".like-button, .dislike-button"
  );

  reactionButtons.forEach((button) => {
    button.addEventListener("click", function () {
      if (window.appForms && window.appForms.submitReaction) {
        window.appForms.submitReaction(this);
      } else {
        console.error("appForms module not loaded correctly");
      }
    });
  });
}

// Export page functions
window.appPages = {
  loadHomePage,
  loadPostPage,
  loadFilteredPosts,
  showLoginPage,
  showRegisterPage,
  showCreatePostPage,
  showErrorPage,
  setupReactionButtons,
};
