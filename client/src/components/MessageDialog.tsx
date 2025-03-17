import { useState } from 'react';
import { Send, X } from 'lucide-react';
import { sendMessage } from '../utils/api';
import { useNavigate } from 'react-router-dom';

interface MessageDialogProps {
  isOpen: boolean;
  onClose: () => void;
  recipient: string;
}

export default function MessageDialog({ isOpen, onClose, recipient }: MessageDialogProps) {
  const [message, setMessage] = useState('');
  const navigate = useNavigate();   
  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await sendMessage(recipient, message);
      setMessage('');
      onClose();
      // Navigate to messages page with the recipient info
      navigate('/messages', { 
        state: { 
          openConversation: {
            username: recipient,
            content: message,
            timestamp: new Date().toISOString()
          }
        }
      });
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex ml-64 items-center justify-center z-50">
      <div className="bg-black border border-gray-800 rounded-xl w-full max-w-lg mx-4">
        <div className="flex items-center justify-between p-4 border-b border-gray-800">
          <h2 className="text-xl font-bold">Message @{recipient}</h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-900 rounded-full"
          >
            <X size={20} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-4">
          <textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Start a new message"
            className="w-full h-32 bg-transparent border border-gray-800 rounded-lg p-3 focus:outline-none focus:border-gray-600 resize-none"
          />
          <div className="flex justify-end mt-4">
            <button
              type="submit"
              disabled={!message.trim()}
              className="px-4 py-2 bg-[#1A8CD8] text-white rounded-full font-bold disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-600 flex items-center gap-2"
            >
              <Send size={18} />
              Send
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}