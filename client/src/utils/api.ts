// Function to get the JWT token from cookies or localStorage
export function getToken(): string | null {
  // If you're using cookies, the token is automatically sent with requests
  // If you're using localStorage:
  return localStorage.getItem('Authentication');
}

// Function to decode JWT token and extract user information
export function decodeToken(token: string): any {
  try {
    // JWT tokens are split into three parts: header.payload.signature
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (error) {
    console.error('Error decoding token:', error);
    return null;
  }
}

// Function to get user information from the backend
export async function getUserInfo(): Promise<any> {
  try {
    const response = await fetch('http://localhost:8080/api/user/info', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Include cookies in the request
    });

    if (!response.ok) {
      throw new Error('Failed to fetch user info');
    }

    return await response.json();
  } catch (error) {
    console.error('Error fetching user info:', error);
    
    // Fallback: If the API call fails, try to get info from the JWT token
    const token = getToken();
    if (token) {
      const decodedToken = decodeToken(token);
      if (decodedToken) {
        return {
          username: decodedToken.username || decodedToken.preferred_username || 'user' + Math.floor(Math.random() * 10000),
          // Add other user properties from the token if available
        };
      }
    }
    
    throw error;
  }
}

// Function to update the username
export async function updateUsername(username: string): Promise<any> {
  const response = await fetch('http://localhost:8080/api/user/update-username', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Include cookies in the request
    body: JSON.stringify({ username }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to update username');
  }

  return await response.json();
}

// Function to validate if a username is available
export async function checkUsernameAvailability(username: string): Promise<boolean> {
  try {
    const response = await fetch(`http://localhost:8080/api/search?q=${encodeURIComponent(username)}&f=unique-user`, {
      method: 'GET',
      credentials: 'include',
    });
    
    if (!response.ok) {
      return false;
    }
    
    const data = await response.json();
    return !data.exists; // Return true if username doesn't exist (is available)
  } catch (error) {
    console.error('Error checking username availability:', error);
    return false;
  }
}

// Posts API functions
export async function getAllPosts(): Promise<any[]> {
  try {
    const response = await fetch('http://localhost:8080/api/posts', {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch posts');
    }

    const data = await response.json();
    console.log("Raw API response:", data);
    
    // Process the posts to ensure they have the correct structure
    const processedPosts = ensurePostsFormat(data);
    console.log("Processed posts:", processedPosts);
    return processedPosts;
  } catch (error) {
    console.error('Error fetching posts:', error);
    return [];
  }
}

// Helper function to ensure posts have the correct format
function ensurePostsFormat(data: any): any[] {
  // If data is an object with a posts property that is an array (most common case)
  if (data && data.posts && Array.isArray(data.posts)) {
    console.log("Processing posts array from data.posts");
    return data.posts.map(processPost);
  }
  
  // If data is already an array
  if (Array.isArray(data)) {
    console.log("Processing direct posts array");
    return data.map(processPost);
  }
  
  // If data is an object with values that should be treated as an array
  if (data && typeof data === 'object') {
    console.log("Processing posts from object values");
    return Object.values(data).map(processPost);
  }
  
  // If all else fails, return an empty array
  console.log("No valid posts data found, returning empty array");
  return [];
}

// Process a single post to ensure it has the correct structure
function processPost(post: any): any {
  if (!post) return null;
  
  console.log("Processing post:", post);
  
  // Ensure the post has all required fields and map field names to match our component expectations
  const processedPost = {
    ...post,
    id: post.id,
    created_at: post.created_at,
    userid: post.userid,
    username: post.username || 'unknown',
    nickname: post.nickname || post.username || 'unknown',
    body: post.body || '',
    is_repost: !!post.is_repost,
    likes_count: post.likes_count || 0,
    // Map parent_id to parentid for consistency with our component
    parentid: post.parent_id,
    parent_id: post.parent_id
  };
  
  console.log("Processed post:", processedPost);
  return processedPost;
}

export async function createPost(body: string): Promise<any> {
  try {
    const response = await fetch('http://localhost:8080/api/posts/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ body }),
    });

    if (!response.ok) {
      throw new Error('Failed to create post');
    }

    return await response.json();
  } catch (error) {
    console.error('Error creating post:', error);
    throw error;
  }
}

export async function getPostsByUsername(username: string): Promise<any[]> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${username}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch posts for user ${username}`);
    }

    const data = await response.json();
    return ensurePostsFormat(data);
  } catch (error) {
    console.error(`Error fetching posts for user ${username}:`, error);
    return [];
  }
}

export async function getPostById(username: string, postId: string): Promise<any> {
  try {
    console.log(`Fetching post with ID ${postId} by user ${username}`);
    
    // First try to get the post directly
    const response = await fetch(`http://localhost:8080/api/posts/${username}/${postId}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      // If that fails, try to get the post without specifying a username
      // This is a fallback in case the API supports this format
      const fallbackResponse = await fetch(`http://localhost:8080/api/posts/post/${postId}`, {
        method: 'GET',
        credentials: 'include',
      });
      
      if (!fallbackResponse.ok) {
        throw new Error(`Failed to fetch post ${postId}`);
      }
      
      const post = await fallbackResponse.json();
      return processPost(post);
    }

    const post = await response.json();
    return processPost(post);
  } catch (error) {
    console.error(`Error fetching post ${postId}:`, error);
    
    // As a last resort, try to find the post in the global posts list
    try {
      const allPosts = await getAllPosts();
      const foundPost = allPosts.find(p => p.id.toString() === postId.toString());
      if (foundPost) {
        return foundPost;
      }
    } catch (fallbackError) {
      console.error('Fallback error:', fallbackError);
    }
    
    throw error;
  }
}

export async function editPost(postId: string, content: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/edit`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ content }),
    });

    if (!response.ok) {
      throw new Error('Failed to edit post');
    }

    return await response.json();
  } catch (error) {
    console.error('Error editing post:', error);
    throw error;
  }
}

export async function deletePost(postId: string): Promise<void> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/delete`, {
      method: 'DELETE',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to delete post');
    }
  } catch (error) {
    console.error('Error deleting post:', error);
    throw error;
  }
}

export async function repostPost(postId: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/repost`, {
      method: 'POST',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to repost');
    }

    return await response.json();
  } catch (error) {
    console.error('Error reposting:', error);
    throw error;
  }
}

export async function quoteRepost(postId: string, quote: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/quote`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ quote }),
    });

    if (!response.ok) {
      throw new Error('Failed to quote repost');
    }

    return await response.json();
  } catch (error) {
    console.error('Error quote reposting:', error);
    throw error;
  }
}