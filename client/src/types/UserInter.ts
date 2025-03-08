export interface UserInfo {
  id: number;
  created_at: string;
  username: string;
  mail: string;
  location?: string | null;
  nickname?: string | null;
  follower_count: number;
}
