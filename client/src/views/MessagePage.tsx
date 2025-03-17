import { useEffect, useState } from "react";
import { Search, Settings, ArrowLeft } from "lucide-react";
import Navbar from "./navbar/Navbar";
import { getMessagesConversation, getUserInfo, listConversations, sendMessage } from "../utils/api";

interface Message {
  id: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  conversation_id: number;
  sender_username: string;
  content: string;
  username: string;
  nickname: string;
  timestamp: string;
}

interface ConversationResponse {
  id: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  sender_username: string;
  sender_nickname: string;
  receiver_username: string;
  receiver_nickname: string;
  messages: Message[];
}

function MessagesPage() {
  const [isLoading, setIsLoading] = useState(true);
  const [conversationList, setConversations] = useState<Message[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<Message | null>(null);
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [secondUser, setSecondUser] = useState<string | null>(null); // Ensure secondUser is a string
  const [conversationResponse, setConversationResponse] = useState<ConversationResponse | null>(null);
  const [newMessage, setNewMessage] = useState("");

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const info = await getUserInfo();
        setCurrentUser(info);
      } catch (error) {
        console.error("Error fetching user info:", error);
      }
    };

    fetchUserInfo();
  }, []);

  const fetchConversationsList = async () => {
    try {
      const fetchedPosts = await listConversations();
      if (Array.isArray(fetchedPosts)) {
        setConversations(fetchedPosts);
      } else {
        console.error("Invalid response format:", fetchedPosts);
        setConversations([]);
      }
    } catch (error) {
      console.error("Error fetching conversation list:", error);
      setConversations([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchConversationsList();
  }, []);

  const handleConversationClick = (conversation: Message) => {
    setSelectedConversation(conversation);
    setSecondUser(conversation.username); // Set secondUser correctly
  };

  useEffect(() => {
    const fetchMessagesConversation = async () => {
      if (!currentUser || !secondUser) {
        console.warn("Skipping fetch: currentUser or secondUser is missing");
        return;
      }

      try {
        const info = await getMessagesConversation(currentUser.username, secondUser);
        setConversationResponse(info);
      } catch (error) {
        console.error("Error fetching conversation:", error);
      }
    };

    fetchMessagesConversation();
  }, [currentUser, secondUser]);

  // Function to send a message
  const handleSendMessage = async () => {
    if (!newMessage.trim() || !secondUser) return;

    try {
      await sendMessage(secondUser, newMessage); // Ensure correct username is used
      setNewMessage(""); // Clear input after sending

      // Refresh messages after sending
      const updatedConversation = await getMessagesConversation(currentUser.username, secondUser);
      setConversationResponse(updatedConversation);
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

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

          {/* Messages List */}
          <div className="overflow-y-auto">
            {isLoading ? (
              <p className="text-gray-500 text-center py-4">Loading conversations...</p>
            ) : conversationList.length === 0 ? (
              <p className="text-gray-500 text-center py-4">No conversations found</p>
            ) : (
              conversationList.map((conversation) => (
                <div
                  key={conversation.id} // Unique key added here
                  className="px-4 py-3 hover:bg-gray-900 cursor-pointer flex items-start gap-3"
                  onClick={() => handleConversationClick(conversation)}
                >
                  <div className="w-10 h-10 text-sm rounded-full bg-gray-800 flex-shrink-0 flex items-center justify-center">
                    {conversation.nickname[0].toUpperCase()}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between">
                      <span className="font-bold truncate">{conversation.nickname}</span>
                      <span className="text-gray-500">{conversation.timestamp}</span>
                    </div>
                    <p className="text-gray-500 truncate">{conversation.content}</p>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Right side - Message content */}
        <div className="flex-1 flex flex-col">
          {selectedConversation ? (
            <>
              <div className="p-3 border-b border-gray-800">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gray-800 flex items-center justify-center">
                    {selectedConversation.nickname[0]}
                  </div>
                  <div>
                    <h2 className="font-bold">{selectedConversation.nickname}</h2>
                    <p className="text-sm text-gray-500">@{selectedConversation.username}</p>
                  </div>
                </div>
              </div>

              <div className="flex-1 overflow-y-auto p-4">
                <div className="text-white space-y-2">
                  {conversationResponse?.messages?.map((message) => (
                    <div
                      key={message.id} // Unique key added here
                      className={`bg-gray-800 p-3 rounded-lg w-fit max-w-xs ${
                        message.sender_username === currentUser?.username ? "ml-auto" : "mr-auto"
                      }`}
                    >
                      <p className="text-sm text-gray-400">{message.sender_username}</p>
                      <p>{message.content}</p>
                    </div>
                  ))}
                </div>
              </div>

              <div className="p-4 border-t border-gray-800">
                <div className="flex items-end gap-2">
                  <input
                    type="text"
                    placeholder="Start a new message"
                    className="flex-1 bg-transparent border border-gray-800 rounded-2xl px-4 py-3 focus:outline-none focus:border-gray-700"
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                  />
                  <button
                    className="text-blue-400 px-4 py-2 rounded-full disabled:opacity-50"
                    onClick={handleSendMessage}
                    disabled={!newMessage.trim()}
                  >
                    Send
                  </button>
                </div>
              </div>
            </>
          ) : (
            <div className="flex items-center justify-center flex-1 text-gray-500">
              Select a conversation to start messaging
            </div>
          )}
        </div>
      </div>
    </Navbar>
  );
}

export default MessagesPage;
