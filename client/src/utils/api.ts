export function getToken(): string | null {
  return localStorage.getItem('Authentication');
}

export function decodeToken(token: string): any {
  try {
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

export async function getUserInfo(): Promise<any> {
  try {
    const response = await fetch('http://localhost:8080/api/user/info', {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to fetch user info');

    return await response.json();
  } catch (error) {
    console.error('Error fetching user info:', error);
    const token = getToken();
    if (token) {
      const decodedToken = decodeToken(token);
      if (decodedToken) {
        return {
          username: decodedToken.username || decodedToken.preferred_username || 'user' + Math.floor(Math.random() * 10000),
        };
      }
    }
    throw error;
  }
}

export async function updateUsername(username: string): Promise<any> {
  const response = await fetch('http://localhost:8080/api/user/update-username', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ username }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to update username');
  }

  return await response.json();
}

export async function checkUsernameAvailability(username: string): Promise<boolean> {
  try {
    const response = await fetch(`http://localhost:8080/api/search?q=${encodeURIComponent(username)}&f=unique-user`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;

    const data = await response.json();
    return !data.exists;
  } catch (error) {
    console.error('Error checking username availability:', error);
    return false;
  }
}

export async function getAllPosts(): Promise<any[]> {
  try {
    const response = await fetch('http://localhost:8080/api/posts', {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to fetch posts');

    const data = await response.json();
    return ensurePostsFormat(data);
  } catch (error) {
    console.error('Error fetching posts:', error);
    return [];
  }
}

function ensurePostsFormat(data: any): any[] {
  if (data && data.posts && Array.isArray(data.posts)) return data.posts.map(processPost);
  if (Array.isArray(data)) return data.map(processPost);
  if (data && typeof data === 'object') return Object.values(data).map(processPost);
  return [];
}

function processPost(post: any): any {
  if (!post) return null;
  return {
    ...post,
    id: post.id,
    created_at: post.created_at,
    userid: post.userid,
    username: post.username || 'unknown',
    nickname: post.nickname || post.username || 'unknown',
    body: post.body || '',
    is_repost: !!post.is_repost,
    likes_count: post.likes_count || 0,
    parentid: post.parent_id,
    parent_id: post.parent_id,
  };
}

export async function createPost(body: string): Promise<any> {
  try {
    const response = await fetch('http://localhost:8080/api/posts/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ body }),
    });

    if (!response.ok) throw new Error('Failed to create post');

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

    if (!response.ok) throw new Error(`Failed to fetch posts for user ${username}`);

    const data = await response.json();
    return ensurePostsFormat(data);
  } catch (error) {
    console.error(`Error fetching posts for user ${username}:`, error);
    return [];
  }
}

export async function getPostById(username: string, postId: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${username}/${postId}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      const fallbackResponse = await fetch(`http://localhost:8080/api/posts/post/${postId}`, {
        method: 'GET',
        credentials: 'include',
      });

      if (!fallbackResponse.ok) throw new Error(`Failed to fetch post ${postId}`);

      return processPost(await fallbackResponse.json());
    }

    return processPost(await response.json());
  } catch (error) {
    console.error(`Error fetching post ${postId}:`, error);
    try {
      const allPosts = await getAllPosts();
      return allPosts.find(p => p.id.toString() === postId.toString()) || null;
    } catch (fallbackError) {
      console.error('Fallback error:', fallbackError);
    }
    throw error;
  }
}

export async function editPost(postId: string, body: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/edit`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ body }),
    });

    if (!response.ok) throw new Error('Failed to edit post');

    return await response.json();
  } catch (error) {
    console.error('Error editing post:', error);
    throw error;
  }
}

export async function editQuote(postId: string, quote: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/edit`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ quote }),
    });

    if (!response.ok) throw new Error('Failed to edit quote');

    return await response.json();
  } catch (error) {
    console.error('Error editing quote:', error);
    throw error;
  }
}

export async function deletePost(postId: string): Promise<void> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/delete`, {
      method: 'DELETE',
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to delete post');
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

    if (!response.ok) throw new Error('Failed to repost');

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
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ quote }),
    });

    if (!response.ok) throw new Error('Failed to quote repost');

    return await response.json();
  } catch (error) {
    console.error('Error quote reposting:', error);
    throw error;
  }
}