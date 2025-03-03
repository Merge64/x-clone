import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { MessageSquare, Repeat, Heart, Share, MoreHorizontal } from 'lucide-react';
import { repostPost, deletePost, getUserInfo, editPost, editQuote } from '../../utils/api';

interface PostData {
  id: string | number;
  created_at: string;
  userid: number;
  username: string;
  nickname?: string;
  parent_id?: number | string | null;
  quote?: string | null;
  body: string;
  is_repost?: boolean;
  likes_count: number;
  parent_post?: {
    id: string | number;
    created_at: string;
    username: string;
    nickname?: string;
    body: string;
  } | null;
}

interface PostProps {
  post: PostData;
  onRepost?: () => void;
  onEdit?: (postId: string, newContent: string) => void;
  onDelete?: (postId: string) => void;
}

function Post({ post, onRepost, onEdit, onDelete }: PostProps) {
  const [isReposting, setIsReposting] = useState(false);
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [showOptions, setShowOptions] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  
  const MAX_CHARS = 280;
  const isQuoteRepost = post.is_repost && post.parent_id !== null && post.quote;
  const isSimpleRepost = post.is_repost && post.parent_id !== null && !post.quote;
  
  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const info = await getUserInfo();
        setCurrentUser(info);
      } catch (error) {
        console.error('Error fetching user info:', error);
      }
    };

    fetchUserInfo();
  }, []);

  const isCurrentUserAuthor = post.username === currentUser?.username;

  const toggleOptions = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowOptions(!showOptions);
  };
  
  const handleEdit = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    const contentToEdit = isQuoteRepost ? (post.quote || '') : post.body;
    setEditContent(contentToEdit);
    setIsEditing(true);
    setShowOptions(false);
  };
  
  const handleCancelEdit = () => {
    setIsEditing(false);
    setError('');
  };
  
  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    if (newContent.length <= MAX_CHARS) {
      setEditContent(newContent);
    }
  };
  
  const handleSaveEdit = async () => {
    if (!editContent.trim()) {
      setError('Post cannot be empty');
      return;
    }
    
    setIsSubmitting(true);
    setError('');
    
    try {
      // Use different API endpoints based on whether it's a quote repost or regular post
      if (isQuoteRepost) {
        await editQuote(post.id.toString(), editContent);
      } else {
        await editPost(post.id.toString(), editContent);
      }
      
      if (onEdit) {
        onEdit(post.id.toString(), editContent);
      }
      
      setIsEditing(false);
    } catch (error) {
      console.error(`Error editing ${isQuoteRepost ? 'quote' : 'post'}:`, error);
      setError(`Failed to edit ${isQuoteRepost ? 'quote' : 'post'}. Please try again.`);
    } finally {
      setIsSubmitting(false);
    }
  };
  
  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isDeleting) return;
    
    if (window.confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
      setIsDeleting(true);
      try {
        await deletePost(post.id.toString());
        
        if (onDelete) {
          onDelete(post.id.toString());
        }
      } catch (error) {
        console.error('Error deleting post:', error);
        alert('Failed to delete post. Please try again.');
      } finally {
        setIsDeleting(false);
        setShowOptions(false);
      }
    } else {
      setShowOptions(false);
    }
  };
  
  const handleRepost = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isReposting) return;
    
    setIsReposting(true);
    try {
      await repostPost(post.id.toString());
      if (onRepost) {
        onRepost();
      }
    } catch (error) {
      console.error('Error reposting:', error);
    } finally {
      setIsReposting(false);
    }
  };
  
  const displayName = post.nickname || post.username;
  
  const formatDate = (dateString: string) => {
    try {
      const cleanedDateString = dateString.replace(/(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})(\.\d+)? ([-+]\d{4}) [-+]\d{2}/, '$1');
      const date = new Date(cleanedDateString);
      if (isNaN(date.getTime())) {
        console.error('Invalid date:', dateString);
        return 'some time ago';
      }
      return formatDistanceToNow(date, { addSuffix: true });
    } catch (error) {
      console.error('Error formatting date:', error);
      return 'some time ago';
    }
  };
  
  return (
    <div className="border-b border-gray-800 p-4 hover:bg-gray-900/50 transition-colors">
      {post.is_repost && (
        <div className="flex items-center text-gray-500 text-sm mb-2">
          <Repeat size={14} className="mr-2" />
          <span>{displayName} reposted</span>
        </div>
      )}
      
      <div className="block">
        <div className="flex">
          <div className="mr-3">
            <Link to={`/profile/${post.username}`} className="block">
              <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center text-lg font-bold">
                {displayName.charAt(0).toUpperCase() || '?'}
              </div>
            </Link>
          </div>
          
          <div className="flex-1">
            <div className="flex items-center">
              <div className="flex-1">
                <Link 
                  to={`/profile/${post.username}`} 
                  className="font-bold hover:underline"
                >
                  {displayName}
                </Link>
                <span className="text-gray-500 ml-1">@{post.username}</span>
                <span className="text-gray-500 ml-2">·</span>
                <span className="text-gray-500 ml-2">
                  {formatDate(post.created_at)}
                </span>
              </div>
              
              {isCurrentUserAuthor && (
                <div className="relative">
                  <button 
                    onClick={toggleOptions}
                    className="text-gray-500 hover:text-blue-400 transition-colors"
                    aria-label="More options"
                  >
                    <MoreHorizontal size={18} />
                  </button>
                  
                  {showOptions && (
                    <div className="absolute right-0 mt-2 w-32 bg-black rounded-lg shadow-lg border border-gray-800 z-10 overflow-hidden">
                      <button 
                        onClick={handleEdit}
                        className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                      >
                        <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                        </svg>
                        Edit
                      </button>
                      <button 
                        onClick={handleDelete}
                        className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-red-500"
                        disabled={isDeleting}
                      >
                        <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                        </svg>
                        {isDeleting ? 'Deleting...' : 'Delete'}
                      </button>
                    </div>
                  )}
                </div>
              )}
            </div>
            
            {isQuoteRepost && !isEditing && (
              <div className="mt-2 mb-3">
                <p className="whitespace-pre-wrap">{post.quote}</p>
              </div>
            )}
            
            {!isEditing && !isQuoteRepost && !isSimpleRepost && (
              <div className="mt-1">
                <p className="whitespace-pre-wrap">{post.body}</p>
              </div>
            )}
            
            {isEditing && (
              <div className="mt-3">
                <textarea
                  className="w-full bg-gray-800 border border-gray-700 rounded-lg p-3 text-white resize-none min-h-[100px]"
                  value={editContent}
                  onChange={handleContentChange}
                  maxLength={MAX_CHARS}
                  disabled={isSubmitting}
                  placeholder={isQuoteRepost ? "Edit your quote..." : "Edit your post..."}
                />
                
                <div className="flex items-center justify-between mt-2">
                  <div className={`text-sm ${editContent.length > MAX_CHARS * 0.8 ? 'text-yellow-500' : 'text-gray-500'}`}>
                    {editContent.length}/{MAX_CHARS}
                  </div>
                  
                  <div className="flex space-x-2">
                    <button
                      onClick={handleCancelEdit}
                      className="px-3 py-1 rounded-full border border-gray-600 text-gray-300 hover:bg-gray-800"
                      disabled={isSubmitting}
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleSaveEdit}
                      className={`px-3 py-1 rounded-full ${
                        isSubmitting || !editContent.trim()
                          ? 'bg-blue-800 text-gray-300 cursor-not-allowed'
                          : 'bg-blue-500 hover:bg-blue-600 text-white'
                      }`}
                      disabled={isSubmitting || !editContent.trim()}
                    >
                      {isSubmitting ? 'Saving...' : 'Save'}
                    </button>
                  </div>
                </div>
                
                {error && <p className="text-red-500 mt-1 text-sm">{error}</p>}
              </div>
            )}
            
            {(isQuoteRepost || isSimpleRepost) && post.parent_post && (
              <div className="mt-3 border border-gray-700 rounded-lg p-4 hover:bg-gray-800/50">
                <div className="flex items-center">
                  <div className="w-5 h-5 rounded-full bg-gray-700 flex items-center justify-center text-xs font-bold mr-2">
                    {(post.parent_post.nickname || post.parent_post.username).charAt(0).toUpperCase()}
                  </div>
                  <span className="font-bold mr-1">
                    {post.parent_post.nickname || post.parent_post.username}
                  </span>
                  <span className="text-gray-500">@{post.parent_post.username}</span>
                  <span className="text-gray-500 ml-2">·</span>
                  <span className="text-gray-500 ml-2 text-sm">
                    {formatDate(post.parent_post.created_at)}
                  </span>
                </div>
                <div className="mt-2">
                  <p className="whitespace-pre-wrap text-gray-300">{post.parent_post.body}</p>
                </div>
              </div>
            )}
            
            {!isEditing && (
              <div className="flex justify-between mt-3 text-gray-500">
                <button className="flex items-center hover:text-blue-400">
                  <MessageSquare size={18} className="mr-1" />
                  <span>0</span>
                </button>
                
                <button 
                  className={`flex items-center ${isReposting ? 'text-green-500' : 'hover:text-green-500'}`}
                  onClick={handleRepost}
                >
                  <Repeat size={18} className="mr-1" />
                  <span>0</span>
                </button>
                
                <button className="flex items-center hover:text-red-500">
                  <Heart size={18} className="mr-1" />
                  <span>{post.likes_count}</span>
                </button>
                
                <button className="flex items-center hover:text-blue-400">
                  <Share size={18} />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Post;