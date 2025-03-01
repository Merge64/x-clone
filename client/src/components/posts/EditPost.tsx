import React, { useState } from 'react';
import { editPost, editQuote } from '../../utils/api';

interface EditPostProps {
  postId: string | number;
  initialContent: string;
  isQuoteRepost?: boolean;
  onCancel: () => void;
  onSave: (postId: string, newContent: string) => void;
}

function EditPost({ postId, initialContent, isQuoteRepost = false, onCancel, onSave }: EditPostProps) {
  const [editContent, setEditContent] = useState(initialContent);
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const MAX_CHARS = 280;
  const characterCount = editContent.length;

  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    if (newContent.length <= MAX_CHARS) {
      setEditContent(newContent);
    }
  };

  const handleSave = async () => {
    if (!editContent.trim()) {
      setError('Post cannot be empty');
      return;
    }
    
    setIsSubmitting(true);
    setError('');
    
    try {
      // Use different API endpoints based on whether it's a quote repost or regular post
      if (isQuoteRepost) {
        await editQuote(postId.toString(), editContent);
      } else {
        await editPost(postId.toString(), editContent);
      }
      
      onSave(postId.toString(), editContent);
    } catch (error) {
      console.error(`Error editing ${isQuoteRepost ? 'quote' : 'post'}:`, error);
      setError(`Failed to edit ${isQuoteRepost ? 'quote' : 'post'}. Please try again.`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
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
        <div className={`text-sm ${characterCount > MAX_CHARS * 0.8 ? 'text-yellow-500' : 'text-gray-500'}`}>
          {characterCount}/{MAX_CHARS}
        </div>
        
        <div className="flex space-x-2">
          <button
            onClick={onCancel}
            className="px-3 py-1 rounded-full border border-gray-600 text-gray-300 hover:bg-gray-800"
            disabled={isSubmitting}
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
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

export default EditPost;