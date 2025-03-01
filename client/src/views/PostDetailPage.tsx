import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { ArrowLeft, MoreHorizontal, MessageSquare, Repeat, Heart, Share, Trash2, Edit } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import Layout from '../components/Layout';
import { getPostById, deletePost, getUserInfo } from '../utils/api';

function PostDetailPage() {
  const { username, postId } = useParams<{ username: string; postId: string }>();
  const navigate = useNavigate();
  const [post, setPost] = useState<any>(null);
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showOptions, setShowOptions] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const fetchPost = async () => {
    if (!username || !postId) return;
    
    setIsLoading(true);
    setError('');
    
    try {
      const fetchedPost = await getPostById(username, postId);
      setPost(fetchedPost);
      
      const userInfo = await getUserInfo();
      setCurrentUser(userInfo);
    } catch (error) {
      console.error('Error fetching post:', error);
      setError('Failed to load post. It may have been deleted or is unavailable.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPost();
  }, [username, postId]);

  const handleDeletePost = async () => {
    if (!post || !window.confirm('Are you sure you want to delete this post?')) {
      return;
    }
    
    setIsDeleting(true);
    try {
      await deletePost(post.id);
      navigate('/home');
    } catch (error) {
      console.error('Error deleting post:', error);
      setError('Failed to delete post. Please try again.');
      setIsDeleting(false);
    }
  };

  const isOwnPost = currentUser?.username === post?.username;

  if (isLoading) {
    return (
      <Layout>
        <div className="flex justify-center p-6">
          <div className="w-8 h-8 border-t-2 border-b-2 border-blue-500 rounded-full animate-spin"></div>
        </div>
      </Layout>
    );
  }

  if (error || !post) {
    return (
      <Layout>
        <div className="p-6 text-center">
          <p className="text-red-500">{error || 'Post not found'}</p>
          <Link 
            to="/home"
            className="mt-4 inline-block px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600"
          >
            Back to Home
          </Link>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="border-b border-gray-800 p-4 flex items-center">
        <Link to="/home" className="mr-4">
          <ArrowLeft size={20} />
        </Link>
        <h1 className="text-xl font-bold">Post</h1>
      </div>
      
      <div className="p-4 border-b border-gray-800">
        <div className="flex">
          <div className="mr-3">
            <Link to={`/profile/${post.username}`}>
              <div className="w-12 h-12 rounded-full bg-gray-700 flex items-center justify-center text-xl font-bold">
                {post.username?.charAt(0).toUpperCase() || '?'}
              </div>
            </Link>
          </div>
          
          <div className="flex-1">
            <div className="flex items-center justify-between">
              <div>
                <Link 
                  to={`/profile/${post.username}`} 
                  className="font-bold hover:underline"
                >
                  {post.username}
                </Link>
                <p className="text-gray-500">@{post.username}</p>
              </div>
              
              <div className="relative">
                <button 
                  onClick={() => setShowOptions(!showOptions)}
                  className="p-2 rounded-full hover:bg-gray-800"
                >
                  <MoreHorizontal size={20} />
                </button>
                
                {showOptions && isOwnPost && (
                  <div className="absolute right-0 mt-2 w-48 bg-black border border-gray-800 rounded-lg shadow-lg z-10">
                    <button 
                      onClick={() => {
                        setShowOptions(false);
                        navigate(`/edit-post/${post.id}`);
                      }}
                      className="w-full text-left px-4 py-2 hover:bg-gray-800 flex items-center"
                      disabled={isDeleting}
                    >
                      <Edit size={16} className="mr-2" />
                      <span>Edit Post</span>
                    </button>
                    <button 
                      onClick={() => {
                        setShowOptions(false);
                        handleDeletePost();
                      }}
                      className="w-full text-left px-4 py-2 text-red-500 hover:bg-gray-800 flex items-center"
                      disabled={isDeleting}
                    >
                      <Trash2 size={16} className="mr-2" />
                      <span>{isDeleting ? 'Deleting...' : 'Delete Post'}</span>
                    </button>
                  </div>
                )}
              </div>
            </div>
            
            <div className="mt-3">
              <p className="text-xl whitespace-pre-wrap">{post.content}</p>
            </div>
            
            <div className="mt-4 text-gray-500">
              <time dateTime={post.created_at}>
                {formatDistanceToNow(new Date(post.created_at), { addSuffix: true })}
              </time>
            </div>
            
            <div className="flex justify-between mt-4 py-4 border-t border-gray-800 text-gray-500">
              <button className="flex items-center hover:text-blue-400">
                <MessageSquare size={18} className="mr-2" />
                <span>0</span>
              </button>
              
              <button className="flex items-center hover:text-green-500">
                <Repeat size={18} className="mr-2" />
                <span>0</span>
              </button>
              
              <button className="flex items-center hover:text-red-500">
                <Heart size={18} className="mr-2" />
                <span>0</span>
              </button>
              
              <button className="flex items-center hover:text-blue-400">
                <Share size={18} />
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div className="p-4 border-b border-gray-800">
        <div className="flex">
          <div className="mr-3">
            <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center text-lg font-bold">
              {currentUser?.username?.charAt(0).toUpperCase() || '?'}
            </div>
          </div>
          
          <div className="flex-1">
            <textarea
              className="w-full bg-transparent border-none outline-none resize-none placeholder-gray-500"
              placeholder="Post your reply"
              rows={3}
            />
            
            <div className="flex justify-end mt-2">
              <button className="px-4 py-2 bg-blue-500 text-white rounded-full font-bold hover:bg-blue-600">
                Reply
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div className="p-6 text-center text-gray-500">
        <p>No replies yet</p>
      </div>
    </Layout>
  );
}

export default PostDetailPage;