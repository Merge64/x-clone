import React, { useState } from 'react';
import { X, Image, Smile, Calendar, MapPin } from 'lucide-react';
import { PostData } from '../../types/post';
import { addComment } from '../../utils/api';

interface CommentModalProps {
  post: PostData;
  currentUser: any;
  onClose: () => void;
  onCommentAdded: () => void;
  formatDate: (dateString: string) => string;
  displayName: string;
}

const CommentModal: React.FC<CommentModalProps> = ({
  post,
  currentUser,
  onClose,
  onCommentAdded,
  formatDate,
  displayName
}) => {
  const [newComment, setNewComment] = useState('');
  const [isAddingComment, setIsAddingComment] = useState(false);

  const MAX_CHARS = 280;

  const handleCommentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = e.target.value;
    if (text.length <= MAX_CHARS) {
      setNewComment(text);
    }
  };

  const submitComment = async () => {
    if (!newComment.trim() || isAddingComment) return;

    setIsAddingComment(true);
    try {
      const targetPostId = post.is_repost && !post.quote ? post.parent_id : post.id;
      await addComment(Number(targetPostId), newComment);
      setNewComment('');
      onCommentAdded();
      onClose();
    } catch (error) {
      console.error('Error adding comment:', error);
    } finally {
      setIsAddingComment(false);
    }
  };

  // Determine which content to display based on whether it's a repost
  const postContent = post.is_repost && post.quote ? post.quote : post.body;

  return (
    <div 
      className="fixed inset-0 bg-[#242D34]/60 flex items-start justify-center z-50 p-4 pt-64 overflow-y-auto"
      onClick={onClose} // This handles clicking on the overlay
    >
      <div 
        className="bg-black w-full max-w-lg rounded-xl border border-gray-800 shadow-xl"
        onClick={(e) => e.stopPropagation()} // Prevent click propagation inside modal
      >
        <div className="flex items-center p-3">
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white p-1 rounded-full hover:bg-gray-800"
            aria-label="Close"
          >
            <X size={20} />
          </button>
        </div>

        <div className="p-4">
          <div className="flex space-x-3">
            <div className="flex flex-col items-center">
              <div className="w-10 h-10 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-lg font-bold">
                {displayName.charAt(0).toUpperCase()}
              </div>
              <div className="w-0.5 bg-gray-700 h-full mt-2 mb-2"></div>
            </div>

            <div className="flex-1">
              <div className="flex items-center">
                <span className="font-bold mr-1">{displayName}</span>
                <span className="text-gray-500">@{post.username}</span>
                <span className="text-gray-500 ml-2">·</span>
                <span className="text-gray-500 ml-2 text-sm">
                  {formatDate(post.created_at)}
                </span>
              </div>
              <p className="mt-1 whitespace-pre-wrap text-gray-300">{postContent}</p>

              <div className="mt-3 text-gray-500">
                <span>Replying to </span>
                <span className="text-[#1D9BF0] hover:underline">
                  <a href={`/${post.username}`}>@{post.username}</a>
                </span>
              </div>
            </div>
          </div>
        </div>

        <div className="p-4">
          <div className="flex space-x-3">
            <div className="w-10 h-10 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-lg font-bold">
              {currentUser?.nickname?.charAt(0).toUpperCase() || 
               currentUser?.username?.charAt(0).toUpperCase() || 'U'}
            </div>
            <div className="flex-1">
              <textarea
                className="w-full bg-transparent border-none p-2 text-white resize-none focus:outline-none placeholder-gray-500 min-h-[80px]"
                placeholder="Post your reply"
                value={newComment}
                onChange={handleCommentChange}
                maxLength={MAX_CHARS}
                autoFocus
              ></textarea>
            </div>
          </div>
          
          <div className="pt-3 flex items-center justify-between">
            <div className="flex space-x-2 text-blue-400">
              <button type="button" className="p-2 rounded-full hover:bg-blue-500/10">
                <Image size={18} />
              </button>
              <button type="button" className="p-2 rounded-full hover:bg-blue-500/10">
                <Smile size={18} />
              </button>
              <button type="button" className="p-2 rounded-full hover:bg-blue-500/10">
                <Calendar size={18} />
              </button>
              <button type="button" className="p-2 rounded-full hover:bg-blue-500/10">
                <MapPin size={18} />
              </button>
            </div>

            <div className="flex items-center space-x-3">
              {newComment.length > 0 && (
                <div className={`text-sm ${newComment.length > MAX_CHARS * 0.8 ? 'text-yellow-500' : 'text-gray-500'}`}>
                  {newComment.length}/{MAX_CHARS}
                </div>
              )}

              <button
                onClick={submitComment}
                className={`px-4 py-1.5 rounded-full ${
                  !newComment.trim() || isAddingComment
                    ? 'bg-[#787A7A] text-black cursor-not-allowed'
                    : 'bg-[#EFF3F4] hover:bg-[#D7DBDC] text-black'
                } font-bold text-sm`}
                disabled={!newComment.trim() || isAddingComment}
              >
                {isAddingComment ? 'Posting...' : 'Reply'}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CommentModal;