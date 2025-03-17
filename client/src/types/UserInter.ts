export interface UserInfo {
  id: number;
  created_at: string;
  username: string;
  nickname?: string | null;
  follower_count: number;
}
