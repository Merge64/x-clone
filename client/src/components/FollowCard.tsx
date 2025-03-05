import { useEffect, useState } from "react";
import { FollowUser, IsAlreadyFollowing, UnfollowUser } from "../utils/api";
import { UserInfo } from "./user/UserLayout";

export function FollowCard(user: UserInfo) {
  const [isFollowing, setIsFollowing] = useState<boolean | null>(null); // Estado inicial indefinido

  useEffect(() => {
    async function checkFollowingStatus() {
      const stateIsFollowing = await IsAlreadyFollowing(String(user.id));
      setIsFollowing(stateIsFollowing);
    }

    checkFollowingStatus();
  }, [user.id]);

  if (isFollowing === null) {
    return <p>Loading...</p>;
  }

  const selectedAction = isFollowing ? UnfollowUser : FollowUser;
  const selectedText = isFollowing ? "Unfollow" : "Follow";

  const handleButton = () => {
    setIsFollowing(!isFollowing);
    selectedAction(String(user.id));
  };

  return (
    <div className="flex items-center justify-between w-full p-3 hover:bg-gray-800/50 transition-colors rounded-xl">
      <div className="flex items-center gap-2">
        <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center text-white">
          {user.username.charAt(0).toUpperCase() || "?"}
        </div>
        <div className="flex flex-col">
          <div className="flex items-center gap-1">
            <strong>{user.nickname}</strong>
          </div>
          <span className="opacity-60">@{user.username}</span>
        </div>
      </div>

      <aside>
        <button
          className="px-4 py-1.5 bg-white text-black font-bold rounded-full hover:bg-gray-200 transition-colors"
          onClick={handleButton}
        >
          {selectedText}
        </button>
      </aside>
    </div>
  );
}
