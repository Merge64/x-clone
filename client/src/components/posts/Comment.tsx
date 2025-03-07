import React, { useState, useEffect } from 'react';
import { Heart, MessageCircle, Repeat, Share, MoreVertical } from 'lucide-react';
import { repost, toggleLike, deletePost, getUserInfo, editPost, checkIfLiked, checkIfReposted } from '../../utils/api';
import PostOptions from './PostOptions';
import CommentModal from './CommentModal';
import QuoteModal from './QuoteModal';
import { PostData } from '../../types/post';

interface CommentProps {
  comment: PostData;
  formatDate: (date: string | number | Date) => string;
  onCommentAdded?: () => void;
}

function Comment({ comment, formatDate, onCommentAdded }: CommentProps) {
  const [isLiked, setIsLiked] = useState(false);
  const [isReposted, setIsReposted] = useState(false);
  const [likesCount, setLikesCount] = useState(comment.likes_count);
  const [repostsCount, setRepostsCount] = useState(comment.reposts_count);
  const [commentsCount, setCommentsCount] = useState(comment.comments_count);
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [showOptions, setShowOptions] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(comment.body);
  const [showRepostOptions, setShowRepostOptions] = useState(false);
  const [showQuoteModal, setShowQuoteModal] = useState(false);
  const [showCommentModal, setShowCommentModal] = useState(false);
  const [isReposting, setIsReposting] = useState(false);
  const [isUndoingRepost, setIsUndoingRepost] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState('');

  const MAX_CHARS = 280;
  const isCurrentUserAuthor = comment.username === currentUser?.username;

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
    const checkInteractions = async () => {
      try {
        const [likedStatus, repostedStatus] = await Promise.all([
          checkIfLiked(Number(comment.id)),
          checkIfReposted(Number(comment.id))
        ]);
        setIsLiked(likedStatus);
        setIsReposted(repostedStatus);
      } catch (error) {
        console.error('Error checking comment interactions:', error);
      }
    };

    if (currentUser) {
      checkInteractions();
    }
  }, [comment.id, currentUser]);

  const handleLikeToggle = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    const newLikedState = !isLiked;
    setIsLiked(newLikedState);
    setLikesCount(prev => newLikedState ? prev + 1 : prev - 1);

    try {
      await toggleLike(Number(comment.id));
    } catch (error) {
      console.error('Error toggling like:', error);
      setIsLiked(!newLikedState);
      setLikesCount(comment.likes_count);
    }
  };

  const toggleRepostOptions = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowRepostOptions(!showRepostOptions);
  };

  const handleSimpleRepost = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isReposting) return;

    setIsReposting(true);
    setIsReposted(true);
    setRepostsCount(prev => prev + 1);

    try {
      await repost(Number(comment.id), "");
      setShowRepostOptions(false);
    } catch (error) {
      console.error('Error reposting:', error);
      setIsReposted(false);
      setRepostsCount(comment.reposts_count);
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
    setRepostsCount(prev => Math.max(0, prev - 1));

    try {
      await repost(Number(comment.id));
      setShowRepostOptions(false);
    } catch (error) {
      console.error('Error undoing repost:', error);
      setIsReposted(true);
      setRepostsCount(comment.reposts_count);
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
    setRepostsCount(prev => prev + 1);

    try {
      await repost(Number(comment.id), quoteText);
      closeQuoteModal();
    } catch (error) {
      console.error('Error quote reposting:', error);
      setIsReposted(false);
      setRepostsCount(comment.reposts_count);
    } finally {
      setIsReposting(false);
    }
  };

  const openCommentModal = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowCommentModal(true);
  };

  const closeCommentModal = () => {
    setShowCommentModal(false);
  };

  const handleCommentAdded = () => {
    setCommentsCount(prev => prev + 1);
    if (onCommentAdded) {
      onCommentAdded();
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

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditContent(comment.body);
    setError('');
  };

  const handleSaveEdit = async () => {
    if (!editContent.trim()) {
      setError('Comment cannot be empty');
      return;
    }

    try {
      await editPost(comment.id.toString(), editContent);
      setIsEditing(false);
    } catch (error) {
      console.error('Error editing comment:', error);
      setError('Failed to edit comment. Please try again.');
    }
  };

  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isDeleting) return;

    if (window.confirm('Are you sure you want to delete this comment? This action cannot be undone.')) {
      setIsDeleting(true);
      try {
        await deletePost(comment.id.toString());
        // You might want to add a callback here to refresh the parent component
      } catch (error) {
        console.error('Error deleting comment:', error);
        alert('Failed to delete comment. Please try again.');
      } finally {
        setIsDeleting(false);
        setShowOptions(false);
      }
    } else {
      setShowOptions(false);
    }
  };

  return (
    <div className="border-b border-gray-800">
      <div className="p-4 hover:bg-gray-900/50 transition-colors">
        <div className="flex gap-3">
          <div className="w-10 h-10 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-lg font-bold">
            {comment.nickname?.charAt(0).toUpperCase() || comment.username?.charAt(0).toUpperCase()}
          </div>
          <div className="flex-1">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1 mb-0.5">
                <span className="font-bold hover:underline">
                  {comment.nickname || comment.username}
                </span>
                <span className="text-gray-500">
                  @{comment.username} · {formatDate(comment.created_at)}
                </span>
              </div>
              {isCurrentUserAuthor && (
                <div className="relative">
                  <button onClick={toggleOptions} className="p-2 hover:bg-gray-800 rounded-full">
                    <MoreVertical className="w-5 h-5" />
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

            {isEditing ? (
              <div className="mt-2">
                <textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  className="w-full bg-transparent border border-gray-700 rounded-lg p-2 text-white"
                  rows={3}
                  maxLength={MAX_CHARS}
                />
                {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
                <div className="flex justify-end gap-2 mt-2">
                  <button
                    onClick={handleCancelEdit}
                    className="px-4 py-1 rounded-full border border-gray-500 hover:bg-gray-800"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleSaveEdit}
                    className="px-4 py-1 rounded-full bg-blue-500 hover:bg-blue-600"
                  >
                    Save
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-white mb-3">{comment.body}</div>
            )}

            <div className="flex items-center justify-between max-w-md text-gray-500">
              <button className="flex items-center group" onClick={openCommentModal}>
                <div className="p-2 rounded-full group-hover:bg-blue-500/10 group-hover:text-blue-500">
                  <MessageCircle size={18} />
                </div>
                <span className="text-sm group-hover:text-blue-500">{commentsCount}</span>
              </button>

              <div className="relative">
                <button className="flex items-center group" onClick={toggleRepostOptions}>
                  <div className={`p-2 rounded-full ${isReposted ? 'text-green-500' : 'group-hover:bg-green-500/10 group-hover:text-green-500'}`}>
                    <Repeat size={18} fill={isReposted ? "currentColor" : "none"} />
                  </div>
                  <span className={`text-sm ${isReposted ? 'text-green-500' : 'group-hover:text-green-500'}`}>
                    {repostsCount}
                  </span>
                </button>

                {showRepostOptions && (
                  <div className="absolute bottom-full mb-2 right-0 bg-black rounded-lg shadow-lg border border-gray-800 z-10 overflow-hidden w-40">
                    {isReposted ? (
                      <button
                        onClick={handleUndoRepost}
                        className="w-full text-left px-4 py-3 hover:bg-gray-900 flex items-center text-white"
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
                          <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                          </svg>
                          Quote
                        </button>
                      </>
                    )}
                  </div>
                )}
              </div>

              <button className="flex items-center group" onClick={handleLikeToggle}>
                <div className={`p-2 rounded-full ${isLiked ? 'text-red-500' : 'group-hover:bg-red-500/10 group-hover:text-red-500'}`}>
                  <Heart size={18} fill={isLiked ? "currentColor" : "none"} />
                </div>
                <span className={`text-sm ${isLiked ? 'text-red-500' : 'group-hover:text-red-500'}`}>
                  {likesCount}
                </span>
              </button>

              <button className="group">
                <div className="p-2 rounded-full group-hover:bg-blue-500/10 group-hover:text-blue-500">
                  <Share size={18} />
                </div>
              </button>
            </div>
          </div>
        </div>
      </div>

      {showCommentModal && (
        <CommentModal
          post={comment}
          currentUser={currentUser}
          onClose={closeCommentModal}
          onCommentAdded={handleCommentAdded}
          formatDate={formatDate}
          displayName={comment.nickname || comment.username}
        />
      )}

      {showQuoteModal && (
        <QuoteModal
          post={comment}
          currentUser={currentUser}
          onClose={closeQuoteModal}
          onQuoteSubmit={handleQuoteRepost}
          isReposting={isReposting}
          formatDate={formatDate}
          displayName={comment.nickname || comment.username}
        />
      )}
    </div>
  );
}

export default Comment;