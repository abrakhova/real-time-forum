// templates.js - Contains all HTML templates for the forum SPA

const templates = {
  // Home page posts
  postCard: (post) => `
      <div class="post">
        <h2><a href="/post?id=${post.ID}" data-navigate>${post.Title}</a></h2>
        <p>${post.Content}</p>
        <div class="post-meta">
          <div class="left">
            <span class="username">Posted by: <strong>${
              post.Username
            }</strong></span>
            <span class="date">Date: ${post.Date || "Unknown"}</span>
          </div>
          <div class="right">
            <button class="like-button" data-id="${
              post.ID
            }" data-type="like" data-for="post">👍 ${post.Likes || 0}</button>
            <button class="dislike-button" data-id="${
              post.ID
            }" data-type="dislike" data-for="post">👎 ${
    post.Dislikes || 0
  }</button>
          </div>
        </div>
      </div>
    `,

  homePage: (posts) => `
      <h2>All posts</h2>
      ${
        posts.length > 0
          ? posts.map((post) => templates.postCard(post)).join("")
          : "<p>No posts available.</p>"
      }
    `,

  // Single post view
  postDetail: (post) => `
      <div class="post" id="post-${post.ID}">
        <h2>Category: ${post.Category}</h2>
        <h3>${post.Title}</h3>
        <p class="author">Posted by: <strong>${post.Username}</strong> on ${
    post.Date || "Unknown date"
  }</p>
        <div class="post-content">${post.Content}</div>
        
        <div class="post-actions">
          <button class="like-button" data-id="${
            post.ID
          }" data-type="like" data-for="post">
            👍 <span>${post.Likes || 0}</span>
          </button>
          <button class="dislike-button" data-id="${
            post.ID
          }" data-type="dislike" data-for="post">
            👎 <span>${post.Dislikes || 0}</span>
          </button>
        </div>
        
        <h3>Comments (${post.Comments ? post.Comments.length : 0})</h3>
        <div id="comments-container">
          ${templates.comments(post.Comments)}
        </div>
        
        ${templates.commentForm(post.ID)}
      </div>
    `,

  // Comments section
  comments: (comments) => {
    if (!comments || comments.length === 0) {
      return "<p>No comments yet. Be the first to comment!</p>";
    }

    return comments
      .map(
        (comment) => `
        <div class="comment">
          <p class="comment-author"><strong>${
            comment.Username
          }</strong> commented:</p>
          <p class="comment-content">${comment.Content}</p>
          <div class="comment-actions">
            <button class="like-button" data-id="${
              comment.ID
            }" data-type="like" data-for="comment">
              👍 <span>${comment.Likes || 0}</span>
            </button>
            <button class="dislike-button" data-id="${
              comment.ID
            }" data-type="dislike" data-for="comment">
              👎 <span>${comment.Dislikes || 0}</span>
            </button>
          </div>
        </div>
      `
      )
      .join("");
  },

  // Comment form
  commentForm: (postId) => {
    // This uses a global state.sessionID to determine if user is logged in
    if (window.state && window.state.sessionID) {
      return `
          <form id="comment-form">
            <h4>Add a Comment</h4>
            <input type="hidden" name="post_id" value="${postId}">
            <textarea name="content" class="comment-box" placeholder="Write your comment here..." required minlength="1" maxlength="500"></textarea>
            <button type="submit">Submit Comment</button>
          </form>
        `;
    } else {
      return `
          <div class="login-prompt">
            <p>You must be logged in to add a comment. 
                <a href="/login" data-navigate>Login</a> or 
                <a href="/register" data-navigate>Register</a>
            </p>
          </div>
        `;
    }
  },

  // Filtered posts view
  filteredPosts: (title, posts) => `
      <h2>${title}</h2>
      ${
        posts.length > 0
          ? posts
              .map(
                (post) => `
            <div class="post">
              <h2><a href="/post?id=${post.id || post.ID}" data-navigate>${
                  post.title || post.Title
                }</a></h2>
              <p>Category: ${post.category || post.Category}</p>
            </div>
          `
              )
              .join("")
          : "<p>No posts found.</p>"
      }
    `,

  // Auth forms
  loginForm: () => `
      <div class="auth-form">
        <h1>Login</h1>
        <div id="login-error" class="error" style="display: none;"></div>
        <form id="login-form">
          <label for="identifier">Username or Email:</label>
          <input type="text" id="identifier" name="identifier" required placeholder="Enter your username or email">
          
          <label for="password">Password:</label>
          <input type="password" id="password" name="password" required placeholder="Enter your password">
          
          <button type="submit">Login</button>
        </form>
        <p>Don't have an account? <a href="/register" data-navigate>Register here</a></p>
        <div class="category-button">
          <a href="/" data-navigate class="back-to-home">Back to Home</a>
        </div>
      </div>
    `,

  registerForm: () => `
      <div class="auth-form">
        <h1>Create Account</h1>
        <div id="register-error" class="error" style="display: none;"></div>
        <form id="register-form">
          <label for="username">Username:</label>
          <input type="text" id="username" name="username" required placeholder="Choose a username">
          
          <label for="first_name">First Name:</label>
          <input type="text" id="first_name" name="first_name" required placeholder="Enter your first name">
          
          <label for="last_name">Last Name:</label>
          <input type="text" id="last_name" name="last_name" required placeholder="Enter your last name">
          
          <label for="age">Age:</label>
          <input type="number" id="age" name="age" required min="13" placeholder="Enter your age">
          
          <label for="gender">Gender:</label>
          <select id="gender" name="gender" required>
            <option value="">Select gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
            <option value="prefer_not_to_say">Prefer not to say</option>
          </select>
          
          <label for="email">Email Address:</label>
          <input type="email" id="email" name="email" required placeholder="Enter your email">
          
          <label for="password">Password:</label>
          <input type="password" id="password" name="password" required placeholder="Create a password">
          
          <button type="submit">Register</button>
        </form>
        <p>Already have an account? <a href="/login" data-navigate>Login here</a></p>
        <div class="category-button">
          <a href="/" data-navigate class="back-to-home">Back to Home</a>
        </div>
      </div>
    `,

  // Create post form
  createPostForm: () => `
      <div class="form-container">
        <h2>Create a New Post</h2>
        <form id="create-post-form">
          <label for="title">Title:</label>
          <input type="text" id="title" name="title" required minlength="5" maxlength="50">
          
          <label for="content">Content:</label>
          <textarea id="content" name="content" required placeholder="Write your post content here..."></textarea>
          
          <label>Categories:</label>
          <div class="checkbox-container">
            <label><input type="checkbox" name="categories" value="General"> General</label>
            <label><input type="checkbox" name="categories" value="Local News & Events"> Local News & Events</label>
            <label><input type="checkbox" name="categories" value="Viking line"> Viking Line</label>
            <label><input type="checkbox" name="categories" value="Travel"> Travel</label>
            <label><input type="checkbox" name="categories" value="Sailing"> Sailing</label>
            <label><input type="checkbox" name="categories" value="Cuisine & food"> Cuisine & Food</label>
            <label><input type="checkbox" name="categories" value="Politics"> Politics</label>
          </div>
          
          <button type="submit">Publish Post</button>
        </form>
        <div class="category-button" style="margin-top: 20px;">
          <a href="/" data-navigate class="back-to-home">Cancel</a>
        </div>
      </div>
    `,

  // Chat interface template
  chatInterface: (username) => `
    <div class="chat-interface">
    <div class="chat-title">
      <h2>Chat with ${username}</h2>
      <button id="back-to-posts" class="back-button">Back to Posts</button>
    </div>
    
    <div id="messages-container" class="messages-container">
      <div class="chat-empty-state">
        <h3>Start a conversation</h3>
        <p>No messages yet. Send a message to start the conversation.</p>
      </div>
    </div>
    
    <div id="typing-indicator" class="typing-indicator"></div>
    
    <div class="chat-footer">
      <textarea id="message-input" placeholder="Type a message..." rows="2"></textarea>
      <button id="send-message-button">Send</button>
    </div>
  </div>
  `,

  // Helper templates
  loading: () => '<div class="loading">Loading...</div>',

  error: (message) => `
      <div class="error-container">
        <h1>Error</h1>
        <p>${message}</p>
        <div class="category-button">
          <a href="/" data-navigate class="back-to-home">Back to Home</a>
        </div>
      </div>
    `,

  // Auth box in header
  authBox: (sessionID, username) => {
    if (sessionID && username) {
      return `
          <p>Logged in as: <strong>${username}</strong></p>
          <a href="/createPost" data-navigate>Create Post</a>
          <a href="#" id="logout-link">Logout</a>
        `;
    } else {
      return `
          <a href="/login" data-navigate>Login</a>
          <a href="/register" data-navigate>Register</a>
        `;
    }
  },

  // Sidebar with categories
  sidebar: (sessionID) => {
    if (!sessionID) {
      return `
        <aside id="category-sidebar" style="display: none;">
          <!-- Sidebar hidden for non-logged-in users -->
        </aside>
      `;
    }

    // If the user is logged in, display the sidebar with categories
    return `
      <aside id="category-sidebar">
        <h2>Categories:</h2>
        <p><a href="/filter?category=General" data-navigate>General</a></p>
        <p><a href="/filter?category=Local%20News%20%26%20Events" data-navigate>Local News & Events</a></p>
        <p><a href="/filter?category=Viking%20line" data-navigate>Viking line</a></p>
        <p><a href="/filter?category=Travel" data-navigate>Travel</a></p>
        <p><a href="/filter?category=Sailing" data-navigate>Sailing</a></p>
        <p><a href="/filter?category=Cuisine%20%26%20food" data-navigate>Cuisine & food</a></p>
        <p><a href="/filter?category=Politics" data-navigate>Politics</a></p>

        <h2>Filters:</h2>
      <p><a href="/filter?user_created=true" data-navigate>My Posts</a></p>
      <p><a href="/filter?liked=true" data-navigate>Liked Posts</a></p>
      
      <div class="category-button">
        <a href="/" data-navigate class="back-to-home">Back to Home</a>
      </div>
    </aside>
    `;
  },

  // Chat sidebar template
  chatSidebar: () => `
    <div class="chat-header">
      <h3>Chat</h3>
      <span id="chat-status" class="connected">Connected</span>
    </div>
    <div id="users-list" class="users-list">
      <p class="empty-users-message">Loading users...</p>
    </div>
  `,
};

// Make templates available globally
window.templates = templates;
