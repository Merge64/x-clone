import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { ArrowLeft, Calendar } from 'lucide-react';
import { format } from 'date-fns';
import Layout from './Layout';
import PostList from '../components/posts/PostList';
import { getPostsByUsername, getUserInfo, getUserProfile } from '../utils/api';

function ProfilePage() {
  const { username } = useParams<{ username: string }>();
  const [posts, setPosts] = useState<any[]>([]);
  const [userInfo, setUserInfo] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [notFound, setNotFound] = useState(false);
  const [activeTab, setActiveTab] = useState('posts');
  const [isCurrentUser, setIsCurrentUser] = useState(false);

  const fetchUserPosts = async () => {
    if (!username) return;
    
    setIsLoading(true);
    setError('');
    setNotFound(false);
    
    try {
      // Fetch the current user's info to check if this is the current user's profile
      const currentUserInfo = await getUserInfo();
      const isOwnProfile = currentUserInfo.username === username;
      setIsCurrentUser(isOwnProfile);
      
      // Fetch the profile user's info using the new API endpoint
      try {
        const profileData = await getUserProfile(username);
        
        // Format the profile data to match our component's expectations
        setUserInfo({
          username: profileData.Username,
          displayName: profileData.Username, // You can use a different field if available
          bio: 'This is a bio placeholder', // Add this if available in the API
          location: profileData.Location || '',
          joinedAt: profileData.CreatedAt || new Date().toISOString(),
        });
        
        const fetchedPosts = await getPostsByUsername(username);
        // Ensure posts is always an array
        setPosts(Array.isArray(fetchedPosts) ? fetchedPosts : []);
      } catch (profileError) {
        console.error(`Error fetching profile for ${username}:`, profileError);
        setNotFound(true);
        setError(`The account @${username} doesn't exist`);
        setPosts([]);
      }
    } catch (error) {
      console.error('Error fetching profile data:', error);
      setError('Failed to load profile data. Please try again.');
      setPosts([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchUserPosts();
  }, [username]);

  if (isLoading) {
    return (
      <Layout>
        <div className="flex justify-center p-6">
          <div className="w-8 h-8 border-t-2 border-b-2 border-blue-500 rounded-full animate-spin"></div>
        </div>
      </Layout>
    );
  }

  if (notFound) {
    return (
      <Layout>
        <div className="flex flex-col items-center justify-center p-8 text-center mt-16">
          <p className="text-gray-500 mb-6">Hmm...this page doesn't exist. Try searching for something else.</p>
          
          <div className="w-full max-w-xs">
            <Link to="/home">
              <button className="py-2 px-6 bg-[#1A8CD8] text-white rounded-full hover:bg-blue-600">
                Search
              </button>
            </Link>
          </div>
        </div>
      </Layout>
    );
  }

  if (error && !notFound) {
    return (
      <Layout>
        <div className="p-6 text-center text-red-500">
          <p>{error}</p>
          <button 
            onClick={fetchUserPosts}
            className="mt-2 px-4 py-2 bg-[#1A8CD8] text-white rounded-full hover:bg-blue-600"
          >
            Try Again
          </button>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="border-b border-gray-800">
        <div className="p-4 flex items-center">
          <Link to="/home" className="mr-4">
            <ArrowLeft size={20} />
          </Link>
          <div>
            <h1 className="text-xl font-bold">{userInfo?.displayName || userInfo?.username}</h1>
            <p className="text-gray-500 text-sm">{posts.length} posts</p>
          </div>
        </div>
      </div>
      
      <div className="bg-gray-800 h-32"></div>
      
      <div className="p-4 border-b border-gray-800">
        <div className="flex justify-between">
          <div className="mt-[-48px]">
            <div className="w-24 h-24 rounded-full bg-black border-4 border-black flex items-center justify-center text-3xl font-bold">
              {userInfo?.username?.charAt(0).toUpperCase() || '?'}
            </div>
          </div>
          
          {isCurrentUser && (
            <button className="px-4 py-2 border border-gray-600 rounded-full font-bold hover:bg-gray-900">
              Edit profile
            </button>
          )}
        </div>
        
        <div className="mt-4">
          <h2 className="text-xl font-bold">{userInfo?.displayName || userInfo?.username}</h2>
          <p className="text-gray-500">@{userInfo?.username}</p>
          
          {userInfo?.bio && (
            <p className="mt-3">{userInfo.bio}</p>
          )}
          
          {userInfo?.location && (
            <p className="mt-2 text-gray-500">{userInfo.location}</p>
          )}
          
          <div className="flex items-center mt-3 text-gray-500">
            <Calendar size={16} className="mr-1" />
            <span>Joined {userInfo?.joinedAt ? format(new Date(userInfo.joinedAt), 'MMMM yyyy') : 'recently'}</span>
          </div>
          
          <div className="flex mt-3">
            <div className="mr-4">
              <span className="font-bold">0</span>{' '}
              <span className="text-gray-500">Following</span>
            </div>
            <div>
              <span className="font-bold">0</span>{' '}
              <span className="text-gray-500">Followers</span>
            </div>
          </div>
        </div>
      </div>
      
      <div className="flex border-b border-gray-800">
        <button
          className={`flex-1 py-4 text-center font-bold ${
            activeTab === 'posts' 
              ? 'text-white border-b-4 border-blue-500' 
              : 'text-gray-500 hover:bg-gray-900'
          }`}
          onClick={() => setActiveTab('posts')}
        >
          Posts
        </button>
        <button
          className={`flex-1 py-4 text-center font-bold ${
            activeTab === 'replies' 
              ? 'text-white border-b-4 border-blue-500' 
              : 'text-gray-500 hover:bg-gray-900'
          }`}
          onClick={() => setActiveTab('replies')}
        >
          Replies
        </button>
        <button
          className={`flex-1 py-4 text-center font-bold ${
            activeTab === 'media' 
              ? 'text-white border-b-4 border-blue-500' 
              : 'text-gray-500 hover:bg-gray-900'
          }`}
          onClick={() => setActiveTab('media')}
        >
          Media
        </button>
        <button
          className={`flex-1 py-4 text-center font-bold ${
            activeTab === 'likes' 
              ? 'text-white border-b-4 border-blue-500' 
              : 'text-gray-500 hover:bg-gray-900'
          }`}
          onClick={() => setActiveTab('likes')}
        >
          Likes
        </button>
      </div>
      
      <PostList 
        posts={posts} 
        onRepost={fetchUserPosts}
        emptyMessage={
          activeTab === 'posts' 
            ? `@${userInfo?.username} hasn't posted yet` 
            : `No ${activeTab} to display`
        }
      />
    </Layout>
  );
}

export default ProfilePage;