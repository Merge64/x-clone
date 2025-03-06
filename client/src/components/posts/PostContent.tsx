import React from 'react';
import { useNavigate } from 'react-router-dom';
import { PostData } from '../../types/post';

interface PostContentProps {
  post: PostData;
  isEditing: boolean;
  isQuoteRepost: boolean;
  isSimpleRepost: boolean;
  editContent: string;
  handleContentChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  handleCancelEdit: () => void;
  handleSaveEdit: () => void;
  isSubmitting: boolean;
  error: string;
  MAX_CHARS: number;
  formatDate: (dateString: string) => string;
}

const PostContent: React.FC<PostContentProps> = ({
  post,
  isEditing,
  isQuoteRepost,
  isSimpleRepost,
  editContent,
  handleContentChange,
  handleCancelEdit,
  handleSaveEdit,
  isSubmitting,
  error,
  MAX_CHARS,
  formatDate
}) => {
  const navigate = useNavigate();

  const handleReferencedPostClick = (e: React.MouseEvent, username: string, postId: number) => {
    e.preventDefault();
    e.stopPropagation();
    navigate(`/${username}/${postId}`);
  };

  if (isEditing) {
    return (
      <div className="mt-3">
        <textarea
          className="w-full bg-gray-800 border border-gray-700 rounded-lg p-3 text-white resize-none min-h-[100px]"
          value={editContent}
          onChange={handleContentChange}
          maxLength={MAX_CHARS}
          disabled={isSubmitting}
          placeholder={isQuoteRepost ? "Edit your quote..." : "Edit your post..."}
          autoFocus
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
    );
  }

  return (
    <>
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

      {(isQuoteRepost || isSimpleRepost) && post.parent_post && (
        <div 
          className="mt-3 border border-gray-700 rounded-lg p-4 hover:bg-gray-800/50 cursor-pointer"
          onClick={(e) => handleReferencedPostClick(e, post.parent_post!.username, Number(post.parent_post!.id))}
        >
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
    </>
  );
};

export default PostContent;