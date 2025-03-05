import React, { useState } from 'react';
import { X, Image, Smile, Calendar, MapPin } from 'lucide-react';
import { PostData } from '../../types/post';

interface QuoteModalProps {
  post: PostData;
  currentUser: any;
  onClose: () => void;
  onQuoteSubmit: (quoteText: string) => void;
  isReposting: boolean;
  formatDate: (dateString: string) => string;
  displayName: string;
}

const QuoteModal: React.FC<QuoteModalProps> = ({
  post,
  currentUser,
  onClose,
  onQuoteSubmit,
  isReposting,
  formatDate,
  displayName
}) => {
  const [quoteText, setQuoteText] = useState('');
  const MAX_CHARS = 280;

  const handleQuoteTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newText = e.target.value;
    if (newText.length <= MAX_CHARS) {
      setQuoteText(newText);
    }
  };

  const handleSubmit = () => {
    if (quoteText.trim()) {
      onQuoteSubmit(quoteText);
    }
  };

  return (
    <div className="fixed inset-0 bg-[#242D34]/60 flex items-start justify-center z-50 p-4 pt-64 overflow-y-auto">
      <div className="bg-black w-full max-w-lg rounded-xl border border-gray-800 shadow-xl ml-32">
        <div className="flex items-center p-3">
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white p-1"
            aria-label="Close"
          >
            <X size={20} />
          </button>
        </div>
        <div>
          <div className="flex-1 overflow-y-auto p-4">
            <div className="flex mb-4">
              <div className="mr-3">
                <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center text-lg font-bold">
                  {currentUser?.nickname?.charAt(0).toUpperCase() || currentUser?.username?.charAt(0).toUpperCase() || 'U'}
                </div>
              </div>

              <div className="w-full">
                <textarea
                  className="w-full bg-transparent border-none outline-none text-white resize-none mt-3 placeholder-gray-500"
                  placeholder="Add a comment"
                  value={quoteText}
                  onChange={handleQuoteTextChange}
                  maxLength={MAX_CHARS}
                  autoFocus
                ></textarea>
              </div>
            </div>

            <div className="ml-12 flex-1 mb-6">
              <div className="border border-gray-800 rounded-xl p-4 mt-2">
                <div className="flex items-center">
                  <div className="w-8 h-8 rounded-full bg-gray-700 flex items-center justify-center text-xs font-bold mr-2">
                    {displayName.charAt(0).toUpperCase()}
                  </div>
                  <span className="font-bold mr-1">{displayName}</span>
                  <span className="text-gray-500">@{post.username}</span>
                  <span className="text-gray-500 ml-2">·</span>
                  <span className="text-gray-500 ml-2 text-sm">
                    {formatDate(post.created_at)}
                  </span>
                </div>
                <p className="mt-2 text-gray-300 whitespace-pre-wrap">{post.body}</p>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-6 border-t border-gray-800 pt-3 flex items-center justify-between p-4">
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
            {quoteText.length > 0 && (
              <div className={`text-sm ${quoteText.length > MAX_CHARS * 0.8 ? 'text-yellow-500' : 'text-gray-500'}`}>
                {quoteText.length}/{MAX_CHARS}
              </div>
            )}

            <button
              onClick={handleSubmit}
              className={`px-4 py-1.5 rounded-full ${!quoteText.trim() || isReposting
                ? 'bg-[#787A7A] text-black cursor-not-allowed'
                : 'bg-[#EFF3F4] hover:bg-[#D7DBDC] text-black'
                } font-bold text-sm`}
              disabled={!quoteText.trim() || isReposting}
            >
              {isReposting ? 'Posting...' : 'Post'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default QuoteModal;