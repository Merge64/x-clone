import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { ArrowLeft, Calendar } from "lucide-react";
import { format } from "date-fns";
import Navbar from "./navbar/Navbar";
import PostList from "../components/posts/PostList";
import { getLikesByUsername, getPostsByUsername, getRepliesByUsername, getUserInfo, getUserProfile } from "../utils/api";
import { IsAlreadyFollowing, UnfollowUser, FollowUser } from "../utils/api";
import { useNavigate } from "react-router-dom";

function ProfilePage() {
    const { username } = useParams<{ username: string }>();
    const [posts, setPosts] = useState<any[]>([]);
    const [userInfo, setUserInfo] = useState<any>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");
    const [notFound, setNotFound] = useState(false);
    const [activeTab, setActiveTab] = useState("posts");
    const [isCurrentUser, setIsCurrentUser] = useState(false);
    const [isFollowing, setIsFollowing] = useState<boolean>(false);
    const [currentUserInfo, setCurrentUserInfo] = useState<any>(null);
    const navigate = useNavigate();

    const fetchUserPosts = async () => {
        if (!username) return;

        setIsLoading(true);
        setError("");
        setNotFound(false);

        try {
            // Fetch user profile
            const profileData = await getUserProfile(username);
            setUserInfo({
                username: profileData.username,
                nickname: profileData.nickname,
                location: profileData.location,
                joinedAt: profileData.CreatedAt || new Date().toISOString(),
                followersCount: profileData.followerCount,
                followingCount: profileData.followingCount,
            });

            let fetchedPosts = [];

            if (activeTab === "posts") {
                fetchedPosts = await getPostsByUsername(username); // Fetch user’s posts
            } else if (activeTab === "replies") {
                fetchedPosts = await getRepliesByUsername(username);
            } else if (activeTab === "likes") {
                fetchedPosts = await getLikesByUsername(username);
            }

            setPosts(Array.isArray(fetchedPosts) ? fetchedPosts : []);
        } catch (error) {
            console.error(`Error fetching profile or posts for ${username}:`, error);
            setNotFound(true);
            setError(`The account @${username} doesn't exist`);
            setPosts([]);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        const fetchUserInfo = async () => {
            try {
                const info = await getUserInfo();
                setCurrentUserInfo(info);
                const isOwnProfile = info.username === username;
                setIsCurrentUser(isOwnProfile);

                if (!isOwnProfile && username) {
                    const following = await IsAlreadyFollowing(username);
                    setIsFollowing(following);
                }
            } catch (error) {
                console.error("Error fetching user info:", error);
            }
        };
        fetchUserInfo();
    }, [username]);

    const handleFollowAction = async () => {
        if (!username || !currentUserInfo) return;

        try {
            if (isFollowing) {
                await UnfollowUser(username);
            } else {
                await FollowUser(username);
            }
            setIsFollowing(!isFollowing);

            // Refresh user profile to update followers count
            const updatedProfile = await getUserProfile(username);
            setUserInfo((prev: any) => ({
                ...prev,
                followersCount: updatedProfile.FollowersCount || 0,
            }));
        } catch (error) {
            console.error("Error updating follow status:", error);
        }
    };

    useEffect(() => {
        fetchUserPosts();
    }, [activeTab, username]);

    if (notFound) {
        return (
            <Navbar>
                <div className="flex flex-col items-center justify-center p-8 text-center mt-16">
                    <p className="text-gray-500 mb-6">
                        Hmm...this page doesn't exist. Try searching for something else.
                    </p>
                    <div className="w-full max-w-xs">
                        <Link to="/home">
                            <button className="py-2 px-6 bg-[#1A8CD8] text-white rounded-full hover:bg-blue-600">
                                Search
                            </button>
                        </Link>
                    </div>
                </div>
            </Navbar>
        );
    }

    if (error && !notFound) {
        return (
            <Navbar>
                <div className="p-6 text-center text-red-500">
                    <p>{error}</p>
                    <button
                        onClick={fetchUserPosts}
                        className="mt-2 px-4 py-2 bg-[#1A8CD8] text-white rounded-full hover:bg-blue-600"
                    >
                        Try Again
                    </button>
                </div>
            </Navbar>
        );
    }

    return (
        <Navbar>
            <div className="border-b border-gray-800">
                <div className="p-4 flex items-center">
                    <button onClick={() => navigate(-1)} className="mr-4">
                        <ArrowLeft size={20} />
                    </button>

                    <div>
                        <h1 className="text-xl font-bold">
                            {userInfo?.nickname || "Loading..."}
                        </h1>
                        <p className="text-gray-500 text-sm">{posts.length} posts</p>
                    </div>
                </div>
            </div>

            <div className="bg-gray-800 h-32"></div>

            <div className="p-4 border-b border-gray-800">
                <div className="flex justify-between">
                    <div className="mt-[-48px]">
                        <div className="w-24 h-24 rounded-full bg-black border-4 border-black flex items-center justify-center text-3xl font-bold">
                            {(userInfo?.nickname || userInfo?.username || "?").charAt(0).toUpperCase()}
                        </div>
                    </div>

                    {isCurrentUser ? (
                        <button className="px-4 py-2 border border-gray-600 rounded-full font-bold hover:bg-gray-900">
                            Edit profile
                        </button>
                    ) : (
                        <button
                            className={`px-4 py-1.5 rounded-full font-bold transition-colors ${isFollowing
                                ? "bg-transparent border border-gray-600 text-white hover:border-red-500 hover:text-red-500 hover:bg-red-500/10"
                                : "bg-white text-black hover:bg-gray-200"
                            }`}
                            onClick={handleFollowAction}
                        >
                            {isFollowing ? "Unfollow" : "Follow"}
                        </button>
                    )}
                </div>

                <div className="mt-4">
                    <h2 className="text-xl font-bold">
                        {userInfo?.nickname || userInfo?.username || "Loading..."}
                    </h2>
                    <p className="text-gray-500">@{userInfo?.username}</p>

                    {userInfo?.bio && <p className="mt-3">{userInfo.bio}</p>}

                    {userInfo?.location && (
                        <p className="mt-2 text-gray-500">{userInfo.location}</p>
                    )}

                    <div className="flex items-center mt-3 text-gray-500">
                        <Calendar size={16} className="mr-1" />
                        <span>
              Joined{" "}
                            {userInfo?.joinedAt
                                ? format(new Date(userInfo.joinedAt), "MMMM yyyy")
                                : "recently"}
            </span>
                    </div>

                    <div className="flex mt-3">
                        <div className="mr-4">
              <span
                  className="font-bold cursor-pointer hover:underline"
                  onClick={() => navigate("./following")}
              >
                {userInfo?.followingCount || 0}
              </span>{" "}
                            <span className="text-gray-500">Following</span>
                        </div>
                        <div>
              <span
                  className="font-bold cursor-pointer hover:underline"
                  onClick={() => navigate("./followers")}
              >
                {userInfo?.followersCount || 0}
              </span>{" "}
                            <span className="text-gray-500">Followers</span>
                        </div>
                    </div>

                </div>
            </div>

            <div className="flex border-b border-gray-800">
                <button
                    className={`flex-1 py-4 text-center font-bold ${activeTab === "posts"
                        ? "text-white border-b-4 border-blue-500"
                        : "text-gray-500 hover:bg-gray-900"
                    }`}
                    onClick={() => setActiveTab("posts")}
                >
                    Posts
                </button>
                <button
                    className={`flex-1 py-4 text-center font-bold ${activeTab === "replies"
                        ? "text-white border-b-4 border-blue-500"
                        : "text-gray-500 hover:bg-gray-900"
                    }`}
                    onClick={() => setActiveTab("replies")}
                >
                    Replies
                </button>
                <button
                    className={`flex-1 py-4 text-center font-bold ${activeTab === "likes"
                        ? "text-white border-b-4 border-blue-500"
                        : "text-gray-500 hover:bg-gray-900"
                    }`}
                    onClick={() => setActiveTab("likes")}
                >
                    Likes
                </button>
            </div>

            <PostList posts={posts} isLoading={isLoading} />
        </Navbar>
    );
}

export default ProfilePage;
