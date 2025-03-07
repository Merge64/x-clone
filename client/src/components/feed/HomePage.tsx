import { useState, useEffect } from 'react';
import { getAllPosts } from '../../utils/api';
import Layout from '../../views/Layout';
import CreatePost from '../posts/CreatePost';
import PostList from '../posts/PostList';
import UsernamePopup from './UsernamePopup';

function HomePage() {
  const [posts, setPosts] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState('for-you');
  const [isUsernamePopupOpen, setIsUsernamePopupOpen] = useState(false);

  const fetchPosts = async () => {
    setIsLoading(true);
    setError('');
    try {
      const fetchedPosts = await getAllPosts();
      // Ensure posts is always an array
      setPosts(Array.isArray(fetchedPosts) ? fetchedPosts : []);
    } catch (error) {
      console.error('Error fetching posts:', error);
      setError('Failed to load posts. Please try again.');
      setPosts([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handlePostCreated = () => {
    fetchPosts();
  };


  const closeUsernamePopup = () => {
    setIsUsernamePopupOpen(false);
  };

  return (
    <>
      <Layout>
        <div className="border-b border-gray-800">
          <div className="flex">
            <button
              className={`flex-1 py-4 text-center font-bold ${
                activeTab === 'for-you' 
                  ? 'text-white border-b-4 border-blue-500' 
                  : 'text-gray-500 hover:bg-gray-900'
              }`}
              onClick={() => setActiveTab('for-you')}
            >
              For You
            </button>
            <button
              className={`flex-1 py-4 text-center font-bold ${
                activeTab === 'following' 
                  ? 'text-white border-b-4 border-blue-500' 
                  : 'text-gray-500 hover:bg-gray-900'
              }`}
              onClick={() => setActiveTab('following')}
            >
              Following
            </button>
          </div>
        </div>

        <CreatePost onPostCreated={handlePostCreated} />

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
        ) : (
          <PostList 
            posts={posts} 
            onRepost={fetchPosts}
            emptyMessage={
              activeTab === 'for-you' 
                ? "No posts to display. Be the first to post something!" 
                : "You're not following anyone yet, or they haven't posted."
            }
          />
        )}
      </Layout>

      {/* Username Popup */}
      <UsernamePopup 
        isOpen={isUsernamePopupOpen} 
        onClose={closeUsernamePopup} 
      />
    </>
  );
}

export default HomePage;