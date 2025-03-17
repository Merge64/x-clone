import React, { useState } from 'react';
import { createPost } from '../../utils/api';
import { getUserInfo } from '../../utils/api';

interface CreatePostProps {
  onPostCreated: () => void;
}

function CreatePost({ onPostCreated }: CreatePostProps) {
  const [body, setBody] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [userInfo, setUserInfo] = useState<any>(null);
  const [characterCount, setCharacterCount] = useState(0);
  const MAX_CHARS = 280;

  React.useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const info = await getUserInfo();
        setUserInfo(info);
      } catch (error) {
        console.error('Error fetching user info:', error);
      }
    };

    fetchUserInfo();
  }, []);

  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    if (newContent.length <= MAX_CHARS) {
      setBody(newContent);
      setCharacterCount(newContent.length);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!body.trim()) {
      setError('Post cannot be empty');
      return;
    }
    
    setIsSubmitting(true);
    setError('');
    
    try {
      await createPost(body);
      setBody('');
      setCharacterCount(0);
      onPostCreated();
    } catch (error) {
      console.error('Error creating post:', error);
      setError('Failed to create post. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="border-b border-gray-800 p-4">
      <form onSubmit={handleSubmit}>
        <div className="flex">
          <div className="mr-4">
            <div className="w-12 h-12 rounded-full bg-gray-700 flex items-center justify-center text-xl font-bold">
              {userInfo?.username?.charAt(0).toUpperCase() || '?'}
            </div>
          </div>
          <div className="flex-1">
            <textarea
              className="w-full bg-transparent border-none outline-none text-xl resize-none min-h-[100px]"
              placeholder="What's happening?"
              value={body}
              onChange={handleContentChange}
              disabled={isSubmitting}
            />
            
            <div className="flex items-center justify-between mt-2">
              <div className={`text-sm ${characterCount > MAX_CHARS * 0.8 ? 'text-yellow-500' : 'text-gray-500'}`}>
                {characterCount}/{MAX_CHARS}
              </div>
              
              <button
                type="submit"
                disabled={isSubmitting || !body.trim() || characterCount > MAX_CHARS}
                className={`px-4 py-2 rounded-full font-bold ${
                  isSubmitting || !body.trim() || characterCount > MAX_CHARS
                    ? 'bg-[#787A7A] text-black cursor-not-allowed'
                    : 'bg-[#EFF3F4] hover:bg-[#D7DBDC] text-black'
                }`}
              > 
                {isSubmitting ? 'Posting...' : 'Post'}
              </button>
            </div>
            
            {error && <p className="text-red-500 mt-2">{error}</p>}
          </div>
        </div>
      </form>
    </div>
  );
}

export default CreatePost;