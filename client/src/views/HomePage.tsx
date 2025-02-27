import { X, LogOut } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { logout } from '../utils/auth';

function HomePage() {
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-black text-white">
      <header className="border-b border-gray-800 p-4">
        <div className="flex items-center justify-between">
          <X size={30} className="text-white" />
          <button 
            onClick={handleLogout}
            className="flex items-center text-gray-400 hover:text-white"
          >
            <LogOut size={20} className="mr-2" />
            <span>Logout</span>
          </button>
        </div>
      </header>
      <main className="max-w-screen-xl mx-auto p-4">
        <h1 className="text-2xl font-bold mb-4">Home</h1>
        <p>Welcome to your home page! You are successfully authenticated.</p>
      </main>
    </div>
  );
}

export default HomePage;