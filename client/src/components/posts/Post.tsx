import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { MessageSquare, Repeat, Heart, Share, MoreVertical } from 'lucide-react';
import {
  repost,
  toggleLike,
  deletePost,
  getUserInfo,
  editPost,
  editQuote,
  checkIfLiked,
  checkIfReposted,
} from '../../utils/api';
import PostOptions from './PostOptions';
import PostContent from './PostContent';
import CommentModal from './CommentModal';
import QuoteModal from './QuoteModal';
import { PostData } from '../../types/post';

interface PostProps {
  post: PostData;
  onRepost?: () => void;
  onEdit?: (postId: string, newContent: string) => void;
  onDelete?: (postId: string) => void;
}

function Post({ post, onRepost, onEdit, onDelete }: PostProps) {
  const navigate = useNavigate();
  const [isReposting, setIsReposting] = useState(false);
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [showOptions, setShowOptions] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [likesCount, setLikesCount] = useState(post.likes_count);
  const [repostsCount, setRepostsCount] = useState(post.reposts_count);
  const [showRepostOptions, setShowRepostOptions] = useState(false);
  const [showQuoteModal, setShowQuoteModal] = useState(false);
  const [showCommentModal, setShowCommentModal] = useState(false);
  const [commentsCount, setCommentsCount] = useState(post.comments_count || 0);
  const [isLiked, setIsLiked] = useState(false);
  const [isReposted, setIsReposted] = useState(false);
  const [isUndoingRepost, setIsUndoingRepost] = useState(false);

  const MAX_CHARS = 280;
  const isQuoteRepost = Boolean(post.is_repost && post.parent_id !== null && post.quote);
  const isSimpleRepost = Boolean(post.is_repost && post.parent_id !== null && !post.quote);

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

  useEffect(() => {
    setLikesCount(post.likes_count);
    setRepostsCount(post.reposts_count);
    setCommentsCount(post.comments_count || 0);
  }, [post.likes_count, post.reposts_count, post.comments_count]);

  useEffect(() => {
    const checkInteractions = async () => {
      try {
        const [likedStatus, repostedStatus] = await Promise.all([
          checkIfLiked(Number(post.id)),
          checkIfReposted(Number(post.id))
        ]);
        setIsLiked(likedStatus);
        setIsReposted(repostedStatus);
      } catch (error) {
        console.error('Error checking post interactions:', error);
      }
    };

    if (currentUser) {
      checkInteractions();
    }
  }, [post.id, currentUser]);

  useEffect(() => {
    const handleClickOutside = () => {
      if (showRepostOptions) {
        setShowRepostOptions(false);
      }
    };

    document.addEventListener('click', handleClickOutside);
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  }, [showRepostOptions]);

  useEffect(() => {
    if (showCommentModal || showQuoteModal) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'auto';
    }

    return () => {
      document.body.style.overflow = 'auto';
    };
  }, [showCommentModal, showQuoteModal]);

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

  const handleLikeToggle = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    const newLikedState = !isLiked;
    setIsLiked(newLikedState);
    setLikesCount(prevCount => newLikedState ? prevCount + 1 : prevCount - 1);

    try {
      await toggleLike(Number(post.id));

      if (onRepost) {
        onRepost();
      }
    } catch (error) {
      console.error('Error toggling like:', error);
      setIsLiked(!newLikedState);
      setLikesCount(post.likes_count);
    }
  };

  const toggleRepostOptions = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    // Toggle repost options regardless of current state
    setShowRepostOptions(!showRepostOptions);
  };

  const handleSimpleRepost = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isReposting) return;

    setIsReposting(true);
    setIsReposted(true);
    setRepostsCount(prevCount => prevCount + 1);

    try {
      await repost(Number(post.id), "");

      if (onRepost) {
        onRepost();
      }

      setShowRepostOptions(false);
    } catch (error) {
      console.error('Error reposting:', error);
      setIsReposted(false);
      setRepostsCount(post.reposts_count);
    } finally {
      setIsReposting(false);
    }
  };

  const handleUndoRepost = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isUndoingRepost) return;

    setIsUndoingRepost(true);
    setIsReposted(false);
    setRepostsCount(prevCount => Math.max(0, prevCount - 1));

    try {
      await repost(Number(post.id));

      if (onRepost) {
        onRepost();
      }
      
      setShowRepostOptions(false);
    } catch (error) {
      console.error('Error undoing repost:', error);
      setIsReposted(true);
      setRepostsCount(post.reposts_count);
    } finally {
      setIsUndoingRepost(false);
    }
  };

  const openQuoteModal = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowQuoteModal(true);
    setShowRepostOptions(false);
  };

  const closeQuoteModal = () => {
    setShowQuoteModal(false);
  };

  const handleQuoteRepost = async (quoteText: string) => {
    if (isReposting || !quoteText.trim()) return;

    setIsReposting(true);
    setIsReposted(true);
    setRepostsCount(prevCount => prevCount + 1);

    try {
      await repost(Number(post.id), quoteText);

      if (onRepost) {
        onRepost();
      }

      closeQuoteModal();
    } catch (error) {
      console.error('Error quote reposting:', error);
      setIsReposted(false);
      setRepostsCount(post.reposts_count);
    } finally {
      setIsReposting(false);
    }
  };

  const openCommentModal = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowCommentModal(true);
  };

  const closeCommentModal = () => {
    setShowCommentModal(false);
  };

  const handleCommentAdded = () => {
    setCommentsCount(prev => prev + 1);
    if (onRepost) {
      onRepost();
    }
  };

  const handlePostClick = () => {
    navigate(`/${post.username}/${post.id}`);
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

  // Determine what to show for the repost button
  const renderRepostButton = () => {
    return (
      <button
        className={`flex items-center ${isReposted ? 'text-green-500' : 'hover:text-green-500'}`}
        onClick={toggleRepostOptions}
        disabled={isReposting || isUndoingRepost}
      >
        <Repeat
          size={18}
          className="mr-1"
          fill={isReposted ? "currentColor" : "none"}
        />
        <span>{repostsCount}</span>
      </button>
    );
  };

  return (
    <div className="border-b border-gray-800 p-4 hover:bg-gray-900/50 transition-colors">
        {post.is_repost && (
          <div className="flex items-center text-gray-500 text-sm mb-2">
            <Repeat size={14} className="mr-2" />
            <span>{displayName} reposted</span>
          </div>
        )}

      <div className="block" onClick={handlePostClick}>
        <div className="flex space-x-3">
          <div className="w-10 h-10 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-lg font-bold">
            {displayName.charAt(0).toUpperCase()}
          </div>
          
          <div className="flex-1">
            <div className="flex items-center">
              <div className="flex-1">
                <Link
                  to={`/${post.username}`}
                  className="font-bold hover:underline"
                  onClick={(e) => e.stopPropagation()}
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
                    <MoreVertical size={18} />
                  </button>

                  {showOptions && (
                    <PostOptions
                      onEdit={handleEdit}
                      onDelete={handleDelete}
                      isDeleting={isDeleting}
                    />
                  )}
                </div>
              )}
            </div>

            <PostContent
              post={post}
              isEditing={isEditing}
              isQuoteRepost={isQuoteRepost}
              isSimpleRepost={isSimpleRepost}
              editContent={editContent}
              handleContentChange={handleContentChange}
              handleCancelEdit={handleCancelEdit}
              handleSaveEdit={handleSaveEdit}
              isSubmitting={isSubmitting}
              error={error}
              MAX_CHARS={MAX_CHARS}
              formatDate={formatDate}
            />

            {!isEditing && (
              <div className="flex justify-between mt-3 text-gray-500">
                <button
                  className="flex items-center hover:text-blue-400"
                  onClick={openCommentModal}
                >
                  <MessageSquare size={18} className="mr-1" />
                  <span>{commentsCount}</span>
                </button>

                <div className="relative">
                  {renderRepostButton()}

                  {showRepostOptions && (
                    <div className="absolute bottom-full mb-2 right-0 bg-black rounded-lg shadow-lg border border-gray-800 z-10 overflow-hidden w-40">
                      {isReposted ? (
                        <button
                          onClick={handleUndoRepost}
                          className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white whitespace-nowrap"
                        >
                          <Repeat size={16} className="mr-2" />
                          Undo repost
                        </button>
                      ) : (
                        <>
                          <button
                            onClick={handleSimpleRepost}
                            className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                          >
                            <Repeat size={16} className="mr-2" />
                            Repost
                          </button>
                          <button
                            onClick={openQuoteModal}
                            className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                          >
                            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"></path>
                            </svg>
                            Quote
                          </button>
                        </>
                      )}
                    </div>
                  )}
                </div>

                <button
                  className={`flex items-center ${isLiked ? 'text-red-500' : 'hover:text-red-500'}`}
                  onClick={handleLikeToggle}
                >
                  <Heart
                    size={18}
                    className="mr-1"
                    fill={isLiked ? "currentColor" : "none"}
                  />
                  <span>{likesCount}</span>
                </button>

                <button className="flex items-center hover:text-blue-400">
                  <Share size={18} />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {showCommentModal && (
        <CommentModal
          post={post}
          currentUser={currentUser}
          onClose={closeCommentModal}
          onCommentAdded={handleCommentAdded}
          formatDate={formatDate}
          displayName={displayName}
        />
      )}

      {showQuoteModal && (
        <QuoteModal
          post={post}
          currentUser={currentUser}
          onClose={closeQuoteModal}
          onQuoteSubmit={handleQuoteRepost}
          isReposting={isReposting}
          formatDate={formatDate}
          displayName={displayName}
        />
      )}
    </div>
  );
}

export default Post;