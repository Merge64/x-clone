import { useState, useEffect } from "react";
import { getSearchedPosts } from "../utils/api";
import Navbar from "./navbar/Navbar";
import PostList from "../components/posts/PostList";
import { Search } from "lucide-react";
import { UserList } from "../components/search/UserList";

function ExplorePage() {
  const [posts, setPosts] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [activeTab, setActiveTab] = useState("for-you");
  const [keyword, setKeyword] = useState("");
  const [order, setOrder] = useState("");

  const fetchPosts = async () => {
    setIsLoading(true);
    setError("");

    try {
      let fetchedPosts = null;
      fetchedPosts = await getSearchedPosts(keyword, order);
      setPosts(Array.isArray(fetchedPosts) ? fetchedPosts : []);
    } catch (error) {
      console.error("Error fetching posts:", error);
      setError("Failed to load posts. Please try again.");
      setPosts([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, [order, keyword]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchPosts();
  };

  const handleMostLiked = () => {
    setOrder("");
  };

  const handleLatest = () => {
    setOrder("latest");
  };

  const handleUser = () => {
    setOrder("user");
  };

  return (
    <Navbar>
      {/* Search Bar */}
      <div className="sticky top-0 z-10 bg-black bg-opacity-90 backdrop-blur-sm border-b border-gray-800 px-4 py-2">
        <form onSubmit={handleSearch} className="relative">
          <div className="relative flex items-center">
            <div className="absolute left-3 text-gray-500">
              <Search size={18} />
            </div>
            <input
              type="text"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              placeholder="Search posts"
              className="w-full bg-gray-900 rounded-full py-2 pl-10 pr-4 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </form>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-800">
        <div className="flex">
          <button
            className={`flex-1 py-4 text-center font-bold ${
              activeTab === "for-you"
                ? "text-white border-b-4 border-blue-500"
                : "text-gray-500 hover:bg-gray-900"
            }`}
            onClick={() => setActiveTab("for-you")}
          >
            For You
          </button>
        </div>
      </div>

      {/* Search Filters (Optional) */}
      <div className="px-4 py-2 bg-black border-b border-gray-800">
        <div className="flex justify-center gap-2 overflow-x-auto py-2 scrollbar-hide">
          <button
            className="mx-1.5 px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700"
            onClick={handleMostLiked}
          >
            Top
          </button>
          <button
            className="mx-1.5 px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700"
            onClick={handleLatest}
          >
            Latest
          </button>
          <button
            className="mx-1.5 px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700"
            onClick={handleUser}
          >
            People
          </button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-6">
          <div className="w-8 h-8 border-t-2 border-b-2 border-blue-500 rounded-full animate-spin"></div>
        </div>
      ) : error ? (
        <div className="p-6 text-center text-red-500">
          <p>{error}</p>
          <button
            onClick={fetchPosts}
            className="mt-2 px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600"
          >
            Try Again
          </button>
        </div>
      ) : order === "user" ? (
        keyword === "" ? ( // Nueva condición para campo vacío
          <div className="p-6 text-center text-gray-500">
            <p>Enter something to search</p>
          </div>
        ) : (
          <UserList
            users={posts}
            emptyMessage="No users found."
            className="space-y-4 p-4"
          />
        )
      ) : (
        <PostList
                posts={posts}
                onRepost={fetchPosts}
                emptyMessage={activeTab === "for-you"
                  ? "No posts to display. Be the first to post something!"
                  : "You're not following anyone yet, or they haven't posted."} isLoading={false}        />
      )}
    </Navbar>
  );
}

export default ExplorePage;
