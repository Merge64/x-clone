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
  getRepostsCount,
  getCommentsCount,
  getLikesCount,
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
  disableNavigation?: boolean;
}

function Post({ post, onRepost, onEdit, onDelete, disableNavigation = false }: PostProps) {
  const navigate = useNavigate();
  const [isReposting, setIsReposting] = useState(false);
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [showOptions, setShowOptions] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(post.quote || post.body);
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
  const [isLikeProcessing, setIsLikeProcessing] = useState(false);

  const MAX_CHARS = 280;
  const isQuoteRepost = Boolean(post.is_repost && post.parent_id !== null && post.quote);
  const isSimpleRepost = Boolean(post.is_repost && post.parent_id !== null && !post.quote);
  const isCurrentUserAuthor = post.username === currentUser?.username;

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
    const fetchCounts = async () => {
      try {
        const targetId = isSimpleRepost ? post.parent_post?.id : post.id;
        if (!targetId) return;

        const [likes, reposts, comments] = await Promise.all([
          getLikesCount(Number(targetId)),
          getRepostsCount(Number(targetId)),
          getCommentsCount(Number(targetId))
        ]);

        setLikesCount(likes);
        setRepostsCount(reposts);
        setCommentsCount(comments);
      } catch (error) {
        console.error('Error fetching counts:', error);
      }
    };

    fetchCounts();
  }, [post.id, post.parent_post?.id, isSimpleRepost]);

  useEffect(() => {
    const checkInteractions = async () => {
      if (!currentUser) return;

      try {
        const targetId = isSimpleRepost ? post.parent_post?.id : post.id;
        if (!targetId) return;

        const [likedStatus, repostedStatus] = await Promise.all([
          checkIfLiked(Number(targetId)),
          checkIfReposted(Number(targetId))
        ]);

        setIsLiked(likedStatus);
        setIsReposted(repostedStatus);
      } catch (error) {
        console.error('Error checking interactions:', error);
      }
    };

    checkInteractions();
  }, [currentUser, post.id, post.parent_post?.id, isSimpleRepost]);

  const handlePostClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!disableNavigation) {
      navigate(`/${post.username}/status/${post.id}`);
    }
  };

  const toggleOptions = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowOptions(!showOptions);
  };

  const handleEdit = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsEditing(true);
    setShowOptions(false);
  };

  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isDeleting || !window.confirm('Are you sure you want to delete this post?')) return;

    setIsDeleting(true);
    try {
      await deletePost(post.id.toString());
      onDelete?.(post.id.toString());
    } catch (error) {
      console.error('Error deleting post:', error);
      alert('Failed to delete post');
    } finally {
      setIsDeleting(false);
      setShowOptions(false);
    }
  };

  const handleLikeToggle = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isLikeProcessing) return;

    setIsLikeProcessing(true);
    const prevState = isLiked;
    const prevCount = likesCount;

    try {
      const targetId = isSimpleRepost ? post.parent_post?.id : post.id;
      if (!targetId) return;

      setIsLiked(!prevState);
      setLikesCount(count => prevState ? count - 1 : count + 1);

      await toggleLike(Number(targetId));

      const [currentLikeStatus, currentCount] = await Promise.all([
        checkIfLiked(Number(targetId)),
        getLikesCount(Number(targetId))
      ]);

      setIsLiked(currentLikeStatus);
      setLikesCount(currentCount);
    } catch (error) {
      console.error('Error toggling like:', error);
      setIsLiked(prevState);
      setLikesCount(prevCount);
    } finally {
      setIsLikeProcessing(false);
    }
  };

  const handleSimpleRepost = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isReposting) return;

    setIsReposting(true);
    const prevState = isReposted;
    const prevCount = repostsCount;

    try {
      const targetId = isSimpleRepost ? post.parent_post?.id : post.id;
      if (!targetId) return;

      setIsReposted(true);
      setRepostsCount(count => count + 1);

      await repost(Number(targetId), "");
      onRepost?.();

      const currentCount = await getRepostsCount(Number(targetId));
      setRepostsCount(currentCount);
    } catch (error) {
      console.error('Error reposting:', error);
      setIsReposted(prevState);
      setRepostsCount(prevCount);
    } finally {
      setIsReposting(false);
      setShowRepostOptions(false);
    }
  };

  const handleUndoRepost = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isUndoingRepost) return;

    setIsUndoingRepost(true);
    const prevState = isReposted;
    const prevCount = repostsCount;

    try {
      const targetId = isSimpleRepost ? post.parent_post?.id : post.id;
      if (!targetId) return;

      setIsReposted(false);
      setRepostsCount(count => Math.max(0, count - 1));

      await repost(Number(targetId));
      onRepost?.();

      const currentCount = await getRepostsCount(Number(targetId));
      setRepostsCount(currentCount);
    } catch (error) {
      console.error('Error undoing repost:', error);
      setIsReposted(prevState);
      setRepostsCount(prevCount);
    } finally {
      setIsUndoingRepost(false);
    }
  };

  const handleQuoteRepost = async (quoteText: string) => {
    if (isReposting || !quoteText.trim()) return;

    setIsReposting(true);
    try {
      const targetId = isSimpleRepost ? post.parent_post?.id : post.id;
      if (!targetId) return;

      await repost(Number(targetId), quoteText);
      onRepost?.();
      setShowQuoteModal(false);
    } catch (error) {
      console.error('Error quote reposting:', error);
    } finally {
      setIsReposting(false);
    }
  };

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return formatDistanceToNow(date, { addSuffix: true });
    } catch (error) {
      console.error('Error formatting date:', error);
      return 'some time ago';
    }
  };

  return (
    <div
      className={`border-b border-gray-800 hover:bg-gray-900/50 transition-colors ${!disableNavigation ? 'cursor-pointer' : ''}`}
      onClick={handlePostClick}
    >
      {post.is_repost && (
        <div className="flex items-center text-gray-500 text-sm px-4 pt-3">
          <Repeat size={14} className="mr-2" />
          <span>
            <Link to={`/${post.username}`} className="hover:underline" onClick={e => e.stopPropagation()}>
              {post.nickname || post.username}
            </Link>
            &nbsp;reposted
          </span>
        </div>
      )}

      {!post.is_repost && post.parent_id !== null && post.parent_id !== undefined && (
        <div className="flex items-center text-gray-500 text-sm px-4 pt-3">
          <MessageSquare size={14} className="mr-2" />
          <span>
            replying to{" "}
            <Link to={`/${post.parent_post?.username}`} className="hover:underline" onClick={e => e.stopPropagation()}>
              {post.parent_post?.nickname || post.parent_post?.username}
            </Link>
          </span>
        </div>
      )}

      <div className="p-4">
        <div className="flex space-x-3">
          <div className="flex-shrink-0">
            <Link to={`/${post.username}`} onClick={e => e.stopPropagation()}>
              <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center text-lg font-bold">
                {(post.nickname || post.username)?.charAt(0).toUpperCase()}
              </div>

            </Link>
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-1 min-w-0">
                <Link
                  to={`/${post.username}`}
                  className="font-bold hover:underline truncate"
                  onClick={e => e.stopPropagation()}
                >
                  {post.nickname || post.username}
                </Link>
                <span className="text-gray-500">@{post.username}</span>
                <span className="text-gray-500">·</span>
                <span className="text-gray-500">{formatDate(post.created_at)}</span>
              </div>

              {isCurrentUserAuthor && (
                <div className="relative ml-2">
                  <button
                    onClick={toggleOptions}
                    className="p-2 hover:bg-gray-800 rounded-full"
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
              handleContentChange={(e) => setEditContent(e.target.value)}
              handleCancelEdit={() => setIsEditing(false)}
              handleSaveEdit={async () => {
                if (!editContent.trim()) {
                  setError('Post cannot be empty');
                  return;
                }

                setIsSubmitting(true);
                try {
                  if (isQuoteRepost) {
                    await editQuote(post.id.toString(), editContent);
                  } else {
                    await editPost(post.id.toString(), editContent);
                  }
                  onEdit?.(post.id.toString(), editContent);
                  setIsEditing(false);
                } catch (error) {
                  console.error('Error editing post:', error);
                  setError('Failed to edit post');
                } finally {
                  setIsSubmitting(false);
                }
              }}
              isSubmitting={isSubmitting}
              error={error}
              MAX_CHARS={MAX_CHARS}
              formatDate={formatDate}
            />

            <div className="flex justify-between mt-3 text-gray-500">
              <button
                className="flex items-center hover:text-blue-400"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  setShowCommentModal(true);
                }}
              >
                <MessageSquare size={18} className="mr-1" />
                <span>{commentsCount}</span>
              </button>

              <div className="relative">
                <button
                  className={`flex items-center ${isReposted ? 'text-green-500' : 'hover:text-green-500'}`}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    if (isReposted) {
                      handleUndoRepost(e);
                    } else {
                      setShowRepostOptions(!showRepostOptions);
                    }
                  }}
                  disabled={isReposting || isUndoingRepost}
                >
                  <Repeat
                    size={18}
                    className="mr-1"
                    fill={isReposted ? "currentColor" : "none"}
                  />
                  <span>{repostsCount}</span>
                </button>

                {showRepostOptions && !isReposted && (
                  <div className="absolute bottom-full mb-2 right-0 bg-black rounded-lg shadow-lg border border-gray-800 z-10 overflow-hidden w-40">
                    <button
                      onClick={handleSimpleRepost}
                      className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                    >
                      <Repeat size={16} className="mr-2" />
                      Repost
                    </button>
                    <button
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        setShowQuoteModal(true);
                        setShowRepostOptions(false);
                      }}
                      className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
                    >
                      <MessageSquare size={16} className="mr-2" />
                      Quote
                    </button>
                  </div>
                )}
              </div>

              <button
                className={`flex items-center ${isLiked ? 'text-red-500' : 'hover:text-red-500'}`}
                onClick={handleLikeToggle}
                disabled={isLikeProcessing}
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
          </div>
        </div>
      </div>

      {showCommentModal && (
        <CommentModal
          post={post}
          currentUser={currentUser}
          onClose={() => setShowCommentModal(false)}
          onCommentAdded={() => {
            setCommentsCount(prev => prev + 1);
            onRepost?.();
          }}
          formatDate={formatDate}
          displayName={post.nickname || post.username}
        />
      )}

      {showQuoteModal && (
        <QuoteModal
          post={post}
          currentUser={currentUser}
          onClose={() => setShowQuoteModal(false)}
          onQuoteSubmit={handleQuoteRepost}
          isReposting={isReposting}
          formatDate={formatDate}
          displayName={post.nickname || post.username}
        />
      )}
    </div>
  );
}

export default Post;