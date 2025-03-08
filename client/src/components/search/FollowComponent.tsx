import { createContext, useContext, useState } from "react";

type FollowContextType = {
  followingIds: Set<string>;
  toggleFollow: (userId: string, isFollowing: boolean) => void;
};

const FollowContext = createContext<FollowContextType>({
  followingIds: new Set(),
  toggleFollow: () => {},
});

export function FollowProvider({ children }: { children: React.ReactNode }) {
  const [followingIds, setFollowingIds] = useState<Set<string>>(new Set());

  const toggleFollow = (userId: string, isFollowing: boolean) => {
    setFollowingIds((prev) => {
      const newSet = new Set(prev);
      isFollowing ? newSet.delete(userId) : newSet.add(userId);
      return newSet;
    });
  };

  return (
    <FollowContext.Provider value={{ followingIds, toggleFollow }}>
      {children}
    </FollowContext.Provider>
  );
}

export const useFollow = () => useContext(FollowContext);
