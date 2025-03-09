import { useState } from 'react';
import { Search, Settings, ArrowLeft } from 'lucide-react';
import Navbar from './navbar/Navbar';

interface Message {
  id: string;
  username: string;
  nickname: string;
  content: string;
  timestamp: string;
  avatar?: string;
}

function MessagesPage() {

  const [messages] = useState<Message[]>([
    {
      id: '1',
      username: 'Guizzomastuxo',
      nickname: 'Queeny',
      content: 'f',
      timestamp: 'Feb 13',
    },
    {
      id: '2',
      username: 'Queeny',
      nickname: 'Santi, Queeny',
      content: 'Queeny: Hola!',
      timestamp: 'Feb 13',
    }
  ]);

  return (
    <Navbar>
      <div className="flex h-screen bg-black text-white">
        {/* Left sidebar - Messages list */}
        <div className="w-[320px] border-r border-gray-800">
          <div className="p-3 flex items-center justify-between border-b border-gray-800">
            <h1 className="text-xl font-bold">Messages</h1>
            <div className="flex gap-4">
              <button className="hover:bg-gray-800 p-2 rounded-full">
                <Settings size={20} />
              </button>
              <button className="hover:bg-gray-800 p-2 rounded-full">
                <ArrowLeft size={20} />
              </button>
            </div>
          </div>
          
          <div className="p-2">
            <div className="relative">
              <input
                type="text"
                placeholder="Search Direct Messages"
                className="w-full bg-gray-900 rounded-full py-2 pl-10 pr-4 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
              />
              <Search className="absolute left-3 top-2.5 text-gray-500" size={16} />
            </div>
          </div>

          <div className="overflow-y-auto">
            {messages.map((message) => (
              <div
                key={message.id}
                className="px-4 py-3 hover:bg-gray-900 cursor-pointer flex items-start gap-3"
              >
                <div className="w-10 h-10 rounded-full bg-gray-800 flex-shrink-0 flex items-center justify-center">
                  {message.avatar || message.nickname[0]}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <span className="font-bold truncate">{message.nickname}</span>
                    <span className="text-sm text-gray-500">{message.timestamp}</span>
                  </div>
                  <p className="text-gray-500 truncate">{message.content}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Right side - Message content */}
        <div className="flex-1 flex flex-col">
          <div className="p-3 border-b border-gray-800">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-gray-800 flex items-center justify-center">
                Q
              </div>
              <div>
                <h2 className="font-bold">Queeny</h2>
                <p className="text-sm text-gray-500">@Guizzomastuxo</p>
              </div>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4">
            {/* Messages will go here */}
          </div>

          <div className="p-4 border-t border-gray-800">
            <div className="flex items-end gap-2">
              <input
                type="text"
                placeholder="Start a new message"
                className="flex-1 bg-transparent border border-gray-800 rounded-2xl px-4 py-3 focus:outline-none focus:border-gray-700"
              />
              <button className="text-blue-400 px-4 py-2 rounded-full disabled:opacity-50">
                Send
              </button>
            </div>
          </div>
        </div>
      </div>
    </Navbar>
  );
}

export default MessagesPage;