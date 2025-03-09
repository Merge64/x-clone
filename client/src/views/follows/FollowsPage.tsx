import { useState, useEffect } from "react";
import { FollowCard } from "../../components/search/FollowCard";
import { UserInfo } from "../../types/UserInter";
import { getUserInfo, getFollows } from "../../utils/api";
import Navbar from "../navbar/Navbar";
import { useParams } from "react-router-dom";

interface FollowsResponse {
  following_count?: number;
  follower_count?: number;
  users: UserInfo[];
}

function FollowsPage() {
  const { username, typeFollowsURL } = useParams<{ username: string, typeFollowsURL: string }>();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [apiResponse, setFollows] = useState<FollowsResponse | null>(null);

  // Validate and get the correct type
  const getTypeFollows = () => {
    if (typeFollowsURL === "following" || typeFollowsURL === "followers") {
      return typeFollowsURL;
    }
    throw new Error("Invalid follow type");
  };

  const getCount = () => {
    if (!apiResponse) return 0;
    return typeFollowsURL === "following" 
      ? apiResponse.following_count || 0 
      : apiResponse.follower_count || 0;
  };

  const fetchFollows = async () => {
    setIsLoading(true);
    setError("");

    try {
      const typeFollows = getTypeFollows();
      // const info = await getUserInfo();
      const fetchedFollows = await getFollows(String(username), typeFollows);
      setFollows(fetchedFollows);
    } catch (error) {
      console.error("Error fetching follows:", error);
      setError("Failed to load follows list. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchFollows();
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
              <p className="text-sm text-gray-400">
                {getCount()} {typeFollowsURL}
              </p>
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