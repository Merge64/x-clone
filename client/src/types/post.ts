export interface PostData {
    id: number | string;
    created_at: string;
    userid?: string | number;
    username: string;
    nickname?: string;
    body: string;
    is_repost?: boolean;
    parent_id?: string | number;
    parentid?: string | number;
    likes_count: number;
    reposts_count: number;
    comments_count: number;
    quote?: string;
    parent_post?: PostData;
  }
  
  export interface CommentData extends PostData {
    in_reply_to_username?: string;
    in_reply_to_post_id?: string | number;
  }
  
  // Alias types for backward compatibility
  export type Post = PostData;
  export type Comment = CommentData;
  
  export interface User {
    username: string;
    nickname?: string;
  }