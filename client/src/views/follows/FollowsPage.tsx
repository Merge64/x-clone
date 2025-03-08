import { useState, useEffect } from "react";
import { FollowCard } from "../../components/search/FollowCard";
import { UserInfo } from "../../types/UserInter";
import { getUserInfo, getFollows } from "../../utils/api";
import Navbar from "../navbar/Navbar";
import { useParams } from "react-router-dom";

interface FollowingResponse {
  following_count: number;
  users: UserInfo[];
}

function FollowsPage() {
  const { typeFollowsURL } = useParams<{ typeFollowsURL: string }>();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [apiResponse, setFollowing] = useState<FollowingResponse | null>(null);

  // Validate and get the correct type
  const getTypeFollows = () => {
    if (typeFollowsURL === "following" || typeFollowsURL === "followers") {
      return typeFollowsURL;
    }
    throw new Error("Invalid follow type");
  };

  const fetchFollowing = async () => {
    setIsLoading(true);
    setError("");

    try {
      const typeFollows = getTypeFollows();
      const info = await getUserInfo();
      const fetchedFollowing = await getFollows(info.username, typeFollows);
      setFollowing(fetchedFollowing);
    } catch (error) {
      console.error("Error fetching follows:", error);
      setError("Failed to load follows list. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchFollowing();
  }, [typeFollowsURL]); // Re-fetch when the URL parameter changes

  // Handle invalid follow type
  if (!["following", "followers"].includes(typeFollowsURL || "")) {
    return (
      <Navbar>
        <div className="min-h-screen bg-black text-white">
          <div className="p-4 text-center text-red-400">
            Invalid follow type. Please use 'following' or 'followers'.
          </div>
        </div>
      </Navbar>
    );
  }

  return (
    <Navbar>
      <div className="min-h-screen bg-black text-white">
        <div className="max-w-2xl mx-auto">
          <header className="sticky top-0 z-10 bg-black/80 backdrop-blur-sm border-b border-gray-800">
            <div className="px-4 py-3">
              <h1 className="text-xl font-bold capitalize">{typeFollowsURL}</h1>
              {apiResponse && (
                <p className="text-sm text-gray-400">
                  {apiResponse.following_count} {typeFollowsURL}
                </p>
              )}
            </div>
          </header>

          <main>
            {isLoading ? (
              <div className="flex justify-center p-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
              </div>
            ) : error ? (
              <div className="p-4 text-center text-red-400">{error}</div>
            ) : apiResponse?.users.length === 0 ? (
              <div className="p-4 text-center text-gray-400">
                No {typeFollowsURL} found.
              </div>
            ) : (
              <div className="divide-y divide-gray-800">
                {apiResponse?.users.map((user) => (
                  <FollowCard key={user.id} {...user} />
                ))}
              </div>
            )}
          </main>
        </div>
      </div>
    </Navbar>
  );
}

export default FollowsPage;