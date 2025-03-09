import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { FollowUser, getUserInfo, IsAlreadyFollowing, UnfollowUser } from "../../utils/api";
import { UserInfo } from "../../types/UserInter";

export function FollowCard(user: UserInfo) {
  const [isFollowing, setIsFollowing] = useState<boolean | null>(null);
  const [isCurrentUser, setIsCurrentUser] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const info = await getUserInfo();
        const isOwnProfile = info.username === user.username;
        setIsCurrentUser(isOwnProfile);

        if (!isOwnProfile) {
          const following = await IsAlreadyFollowing(user.username);
          setIsFollowing(following);
        }
      } catch (error) {
        console.error("Error fetching user info:", error);
      }
    };
    fetchUserInfo();
  }, [user.username]);

  const handleFollowButton = async (event: React.MouseEvent) => {
    event.stopPropagation(); // Prevents clicking the button from triggering navigation

    setIsFollowing(!isFollowing);
    if (isFollowing) {
      await UnfollowUser(user.username);
    } else {
      await FollowUser(user.username);
    }
  };

  const goToProfile = () => {
    navigate(`/${user.username}`);
  };

  return (
    <div
      className="flex items-center justify-between w-full p-3 hover:bg-gray-800/50 transition-colors rounded-xl cursor-pointer"
      onClick={goToProfile}
    >
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
      {!isCurrentUser && (
        <button
          className={`px-4 py-1.5 rounded-full font-bold transition-colors ${
            isFollowing
              ? "bg-transparent border border-gray-600 text-white hover:border-red-500 hover:text-red-500 hover:bg-red-500/10"
              : "bg-white text-black hover:bg-gray-200"
          }`}
          onClick={handleFollowButton} 
        >
          {isFollowing ? "Unfollow" : "Follow"}
        </button>
      )}
    </div>
  );
}
