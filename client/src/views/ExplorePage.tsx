import { useState, useEffect } from "react";
import { getAllPosts, getSearchedPosts } from "../utils/api";
import Layout from "../components/Layout";
import PostList from "../components/posts/PostList";
import { Search } from "lucide-react";

function ExplorePage() {
  const [posts, setPosts] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [activeTab, setActiveTab] = useState("for-you");
  const [keyword, setKeyword] = useState("");

  const fetchPosts = async () => {
    setIsLoading(true);
    setError("");
    try {
      const fetchedPosts = await getAllPosts();
      // Ensure posts is always an array
      setPosts(Array.isArray(fetchedPosts) ? fetchedPosts : []);
    } catch (error) {
      console.error("Error fetching posts:", error);
      setError("Failed to load posts. Please try again.");
      setPosts([]);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchSearchedPosts = async () => {
    setIsLoading(true);
    setError("");

    try {
      let fetchedPosts = null;

      if (keyword == "") {
        fetchedPosts = await getAllPosts(); // Buscar todos si está vacío
      } else {
        fetchedPosts = await getSearchedPosts(keyword); // Buscar literal
      }

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
  }, []);

  useEffect(() => {
    fetchSearchedPosts();
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchSearchedPosts();
  };

  return (
    <Layout>
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
          <button
            className={`flex-1 py-4 text-center font-bold ${
              activeTab === "following"
                ? "text-white border-b-4 border-blue-500"
                : "text-gray-500 hover:bg-gray-900"
            }`}
            onClick={() => setActiveTab("following")}
          >
            Following
          </button>
        </div>
      </div>

      {/* Search Filters (Optional) */}
      <div className="px-4 py-2 bg-black border-b border-gray-800">
        <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
          <button className="px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700">
            Top
          </button>
          <button className="px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700">
            Latest
          </button>
          <button className="px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700">
            People
          </button>
          <button className="px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700">
            Media
          </button>
          <button className="px-4 py-1 bg-gray-800 rounded-full text-sm font-medium text-white hover:bg-gray-700">
            Lists
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
            onClick={fetchSearchedPosts}
            className="mt-2 px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600"
          >
            Try Again
          </button>
        </div>
      ) : (
        <PostList
          posts={posts}
          onRepost={fetchSearchedPosts}
          emptyMessage={
            activeTab === "for-you"
              ? "No posts to display. Be the first to post something!"
              : "You're not following anyone yet, or they haven't posted."
          }
        />
      )}

      {/* Right Sidebar (Optional) */}
      <div className="hidden lg:block fixed right-0 top-0 h-screen w-80 border-l border-gray-800 p-4 bg-black overflow-y-auto">
        <div className="mb-6">
          <h2 className="text-xl font-bold mb-4">Search filters</h2>
          <div className="bg-gray-900 rounded-xl p-4">
            <h3 className="font-bold mb-2">People</h3>
            <div className="flex items-center justify-between mb-2">
              <span>From anyone</span>
              <div className="h-5 w-5 rounded-full bg-blue-500 flex items-center justify-center">
                <span className="text-white text-xs">✓</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span>People you follow</span>
              <div className="h-5 w-5 rounded-full border border-gray-600"></div>
            </div>

            <h3 className="font-bold mt-4 mb-2">Location</h3>
            <div className="flex items-center justify-between mb-2">
              <span>Anywhere</span>
              <div className="h-5 w-5 rounded-full bg-blue-500 flex items-center justify-center">
                <span className="text-white text-xs">✓</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span>Near you</span>
              <div className="h-5 w-5 rounded-full border border-gray-600"></div>
            </div>

            <button className="text-blue-500 mt-4">Advanced search</button>
          </div>
        </div>

        <div>
          <h2 className="text-xl font-bold mb-4">What's happening</h2>
          <div className="bg-gray-900 rounded-xl p-4 mb-4">
            <div className="mb-4">
              <p className="text-xs text-gray-500">Trending in your area</p>
              <p className="font-bold">#TrendingTopic</p>
              <p className="text-xs text-gray-500">10.4K posts</p>
            </div>
            <div>
              <p className="text-xs text-gray-500">Technology · Trending</p>
              <p className="font-bold">#WebDevelopment</p>
              <p className="text-xs text-gray-500">5.2K posts</p>
            </div>
          </div>

          <button className="text-blue-500">Show more</button>
        </div>
      </div>
    </Layout>
  );
}

export default ExplorePage;
