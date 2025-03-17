import { PostData } from "../types/post";

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

export async function getUserProfile(username: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/profile/${username}`, {
      method: "GET",
      credentials: "include",
    });

    if (!response.ok) throw new Error(`Failed to fetch profile for ${username}`);

    const data = await response.json();

    return {
      username: data.profile.username,
      nickname: data.profile.nickname,
      location: data.profile.location || "",
      bio: data.profile.bio || "",
      followerCount: data.profile.follower_count || 0, 
      followingCount: data.profile.following_count || 0, 
    };
  } catch (error) {
    console.error(`Error fetching profile for ${username}:`, error);
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

export async function checkIfLiked(postId: number): Promise<boolean> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/check/${postId}/liked`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;
    const data = await response.json();
    return data.liked || false;
  } catch (error) {
    console.error('Error checking if post is liked:', error);
    return false;
  }
}

export async function checkIfReposted(postId: number): Promise<boolean> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/check/${postId}/reposted`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;
    const data = await response.json();
    return data.reposted || false;
  } catch (error) {
    console.error('Error checking if post is reposted:', error);
    return false;
  }
}

export async function getRepostsCount(postId: number): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/count/${postId}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;
    const data = await response.json();
    return data.reposts_count
  } catch (error) {
    console.error('Error checking count of reposts:', error);
    return;
  }
}

export async function getLikesCount(postId: number): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/count/${postId}/likes`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;
    const data = await response.json();
    return data.likes_count
  } catch (error) {
    console.error('Error checking count of reposts:', error);
    return;
  }
}


export async function getCommentsCount(postId: number): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/count/${postId}/comments`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;
    const data = await response.json();
    return data.comments_count
  } catch (error) {
    console.error('Error checking count of reposts:', error);
    return;
  }
}

export async function getComments(postId: number): Promise<any[]> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/comments/${postId}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to fetch comments');
    const data = await response.json();
    return Array.isArray(data.comments) ? data.comments.map(processComment) : [];
  } catch (error) {
    console.error('Error fetching comments:', error);
    return [

    ];
  }
}

function processComment(comment: any): any {
  if (!comment) return null;
  return {
    ...comment,
    id: comment.id,
    created_at: comment.created_at,
    userid: comment.userid,
    username: comment.username || 'unknown',
    nickname: comment.nickname || comment.username || 'unknown',
    body: comment.body || '',
    likes_count: comment.likes_count || 0,
    reposts_count: comment.reposts_count || 0,
    comments_count: comment.comments_count || 0,
    in_reply_to_username: comment.in_reply_to_username || comment.parent_username,
    in_reply_to_post_id: comment.in_reply_to_post_id || comment.parent_id
  };
}

export async function addComment(parentId: number, body: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/comments/${parentId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ body }),
    });

    if (!response.ok) throw new Error('Failed to add comment');
    return await response.json();
  } catch (error) {
    console.error('Error adding comment:', error);
    return {
      success: true,
      comment: {
        body: body
      }
    };
  }
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

export async function getSearchedPosts(keyword: string, filter: string) {
  let endpointURl = `http://localhost:8080/api/search?q=${keyword}`;
  if (filter) {
    endpointURl += `&f=${filter}`;
  }

  try {
    const response = await fetch(endpointURl, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) throw new Error('Failed to fetch data');

    const data = await response.json();

    return filter === 'user' ? data.users || data : ensurePostsFormat(data);
  } catch (error) {
    console.error('Error fetching data:', error);
    return [];
  }
}

export async function getPrivateSearchedPosts(keyword: string, filter: string) {
  let endpointURl = `http://localhost:8080/api/private/search?q=${keyword}`;
  if (filter) {
    endpointURl += `&f=${filter}`;
  }
  try {
    const response = await fetch(endpointURl, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) throw new Error('Failed to fetch data');

    const data = await response.json();

    return filter === 'user' ? data.users || data : ensurePostsFormat(data);
  } catch (error) {
    console.error('Error fetching data:', error);
    return [];
  }
}

function ensurePostsFormat(data: any): any[] {
  if (data && data.posts && Array.isArray(data.posts)) return data.posts.map(processPost);
  if (Array.isArray(data)) return data.map(processPost);
  if (data && typeof data === 'object') return Object.values(data).map(processPost);
  return [];
}

export async function FollowUser(username: string) {
  try {
    const response = await fetch(`http://localhost:8080/api/profile/follow/${username}`, {
      method: 'POST',
      credentials: 'include',
    });
    if (!response.ok) throw new Error('Failed to follow user');

    const data = await response.json();
    return ensurePostsFormat(data);

  } catch (error) {
    console.error('Error to following user:', error);
    return [];
  }
}

export async function UnfollowUser(username: string) {
  try {
    const response = await fetch(`http://localhost:8080/api/profile/unfollow/${username}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    if (!response.ok) throw new Error('Failed to unfollow user');
    const data = await response.json();
    return ensurePostsFormat(data);

  } catch (error) {
    console.error('Error to unfollowing user:', error);
    return [];
  }
}

export async function IsAlreadyFollowing(username: string): Promise<boolean> {
  try {
    const response = await fetch(`http://localhost:8080/api/profile/is-following/${username}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) return false;

    const data = await response.json();
    return data.isFollowing === true;

  } catch (error) {
    console.error('Error checking follow status:', error);
    return false;
  }
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
    reposts_count: post.reposts_count || 0,
    comments_count: post.comments_count || 0,
    parentid: post.parent_id,
    parent_id: post.parent_id,
  };
}

export async function toggleLike(id: number): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${id}/like`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to like post');
    
    return response.json();
  } catch (error) {
    console.error('Error liking post:', error);
    throw error;
  }
}

export async function repost(id: number, quote?: string): Promise<any> {
  try {
    // Get current user info first
    const userInfo = await getUserInfo();
    
    const response = await fetch(`http://localhost:8080/api/posts/${id}/repost`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      credentials: 'include',
      body: JSON.stringify({ 
        quote,
        username: userInfo.username,
        nickname: userInfo.nickname
      }),
    });

    if (!response.ok) throw new Error('Failed to repost');
    
    const data = await response.json();
    
    // Ensure the response includes user information
    return {
      ...data,
      username: userInfo.username,
      nickname: userInfo.nickname
    };
  } catch (error) {
    console.error('Error reposting:', error);
    throw error;
  }
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

export async function getPostsByUsername(username: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/user/${username}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to create post');

    return await response.json();
  } catch (error) {
    console.error('Error creating post:', error);
    throw error;
  }
}

export async function getRepliesByUsername(username: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/replies/user/${username}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to create post');

    return await response.json();
  } catch (error) {
    console.error('Error creating post:', error);
    throw error;
  }
}

export async function getLikesByUsername(username: string): Promise<any> {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/likes/user/${username}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (!response.ok) throw new Error('Failed to create post');

    return await response.json();
  } catch (error) {
    console.error('Error creating post:', error);
    throw error;
  }
}

export const getPostById = async ( postId: string): Promise<PostData | null> => {
  try {
    const response = await fetch(`http://localhost:8080/api/posts/${postId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch post');
    }
    const data = await response.json();
    const postData = data.post;

    return {
      id: postData.id,
      created_at: postData.created_at,
      userid: postData.userid,
      username: postData.username,
      nickname: postData.nickname || '',
      body: postData.body || postData.quote || '',
      likes_count: postData.likes_count,
      reposts_count: postData.reposts_count,
      comments_count: postData.comments_count || 0,
      is_repost: postData.is_repost || false,
      parent_id: postData.parent_id,
      quote: postData.quote,
      parent_post: postData.parent_post,
    };
  } catch (error) {
    console.error('Error fetching post:', error);
    return null;
  }
};

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
    // Get current user info first
    const userInfo = await getUserInfo();
    
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/repost`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      credentials: 'include',
      body: JSON.stringify({ 
        username: userInfo.username,
        nickname: userInfo.nickname
      }),
    });

    if (!response.ok) throw new Error('Failed to repost');
    
    const data = await response.json();
    
    // Ensure the response includes user information
    return {
      ...data,
      username: userInfo.username,
      nickname: userInfo.nickname
    };
  } catch (error) {
    console.error('Error reposting:', error);
    throw error;
  }
}

export async function quoteRepost(postId: string, quote: string): Promise<any> {
  try {
    // Get current user info first
    const userInfo = await getUserInfo();
    
    const response = await fetch(`http://localhost:8080/api/posts/${postId}/quote`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      credentials: 'include',
      body: JSON.stringify({ 
        quote,
        username: userInfo.username,
        nickname: userInfo.nickname
      }),
    });

    if (!response.ok) throw new Error('Failed to quote repost');

    const data = await response.json();
    
    // Ensure the response includes user information
    return { 
      ...data,
      username: userInfo.username,
      nickname: userInfo.nickname
    };
  } catch (error) {
    console.error('Error quote reposting:', error);
    throw error;
  }
}

export async function getFollows(username: string, followType: string) {
  try {
    const response = await fetch(`http://localhost:8080/api/profile/${username}/${followType}`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) throw new Error(`Failed to get ${followType} user`);
    const data = await response.json();
    return data;

  } catch (error) {
    console.error(`Error to get ${followType} user:`, error);
    return [];
  }
}

export async function updateUserProfile(profileData: { nickname: string; bio?: string; location?: string; birthdate?: string }): Promise<any> {
  try {
    const response = await fetch('http://localhost:8080/api/profile/edit', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(profileData),
    });

    if (!response.ok) throw new Error('Failed to update profile');

    return await response.json();
  } catch (error) {
    console.error('Error updating profile:', error);
    throw error;
  }
}

export async function listConversations() {
  try {
    const response = await fetch(`http://localhost:8080/api/messages`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) throw new Error(`Failed to get active conversations`);
    const data = await response.json();
    console.log(data)
    return data;

  } catch (error) {
    console.error(`Error:`, error);
    return [];
  }
}

export async function getMessagesConversation(currUsername:string, secondUsername:string) {
  try {
    const response = await fetch(`http://localhost:8080/api/messages/${currUsername}/${secondUsername}`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) throw new Error(`Failed to get active conversations`);
    const data = await response.json();
    console.log(data)
    return data;

  } catch (error) {
    console.error(`Error:`, error);
    return [];
  }
}

export async function sendMessage(receiverUsername: string, message: string) {
  try {
    const response = await fetch(`http://localhost:8080/api/messages/dm/${receiverUsername}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ message }),
    });

    if (!response.ok) throw new Error('Failed to send message');

    return await response.json();
  } catch (error) {
    console.error('Error sending message:', error);
    throw error;
  }
}
