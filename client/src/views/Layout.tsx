import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { X, Home, User, LogOut, Bell } from 'lucide-react';
import { logout } from '../utils/auth';
import { getUserInfo } from '../utils/api';

interface LayoutProps {
  children: React.ReactNode;
}

function Layout({ children }: LayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const [currentUsername, setCurrentUsername] = useState<string | null>(null);

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const userInfo = await getUserInfo();
        setCurrentUsername(userInfo.username);
      } catch (error) {
        console.error('Error fetching user info:', error);
      }
    };

    fetchUserInfo();
  }, []);

  return (
    <div className="min-h-screen bg-black text-white flex">
      {/* Sidebar */}
      <div className="w-64 border-r border-gray-800 p-4 fixed h-full">
        <div className="flex flex-col h-full">
          <div className="mb-8">
            <X size={30} className="text-white" />
          </div>

          <nav className="flex-1">
            <ul className="space-y-4">
              <li>
                <Link 
                  to="/home" 
                  className={`flex items-center p-2 rounded-full hover:bg-gray-800 ${location.pathname === '/home' ? 'font-bold' : ''}`}
                >
                  <Home size={24} className="mr-4" />
                  <span className="text-xl">Home</span>
                </Link>
              </li>
              <li>
                <Link 
                  to="/notifications" 
                  className={`flex items-center p-2 rounded-full hover:bg-gray-800 ${location.pathname === '/notifications' ? 'font-bold' : ''}`}
                >
                  <Bell size={24} className="mr-4" />
                  <span className="text-xl">Notifications</span>
                </Link>
              </li>
              <li>
                <Link 
                  to={`/${currentUsername ?? ''}`} 
                  className={`flex items-center p-2 rounded-full hover:bg-gray-800 `}
                >
                  <User size={24} className="mr-4" />
                  <span className="text-xl">Profile</span>
                </Link>
              </li>
            </ul>
          </nav>

          <button 
            onClick={handleLogout}
            className="flex items-center p-2 rounded-full hover:bg-gray-800 mt-auto"
          >
            <LogOut size={24} className="mr-4" />
            <span className="text-xl">Logout</span>
          </button>
        </div>
      </div>

      {/* Main content */}
      <div className="ml-64 flex-1">
        <main className="max-w-2xl mx-auto">
          {children}
        </main>
      </div>
    </div>
  );
}

export default Layout;
