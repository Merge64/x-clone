// components/LoginModPl.tsx
import React, { useState } from 'react';
import { X } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export function LoginModal() {
  const [email, setEmail] = useState('');
  const navigate = useNavigate();

  function handleClose() {
    navigate('/');
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
      <div className="bg-black w-full max-w-md rounded-2xl p-8 relative">
        <button
          onClick={handleClose}
          className="absolute left-4 top-4 text-gray-400 hover:text-white"
        >
          <X size={20} />
        </button>
        
        <div className="flex justify-center mb-8">
        <X size={50} className="text-white" />
        </div>

        <h1 className="text-2xl font-bold text-white text-center mb-8">
          Sign in to X
        </h1>

        <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 flex items-center justify-center gap-2 mb-3 hover:bg-gray-200 transition-colors">
          <img src="https://www.google.com/favicon.ico" 
               alt="Google" 
               className="w-5 h-5" />
          Sign in with Google
        </button>

        <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 flex items-center justify-center gap-2 mb-3 hover:bg-gray-200 transition-colors">
          <svg viewBox="0 0 24 24" className="w-5 h-5">
            <path d="M18.71 19.5C17.88 20.74 17 21.95 15.66 21.97C14.32 22 13.89 21.18 12.37 21.18C10.84 21.18 10.37 21.95 9.09997 22C7.78997 22.05 6.79997 20.68 5.95997 19.47C4.24997 17 2.93997 12.45 4.69997 9.39C5.56997 7.87 7.12997 6.91 8.81997 6.88C10.1 6.86 11.32 7.75 12.11 7.75C12.89 7.75 14.37 6.68 15.92 6.84C16.57 6.87 18.39 7.1 19.56 8.82C19.47 8.88 17.39 10.1 17.41 12.63C17.44 15.65 20.06 16.66 20.09 16.67C20.06 16.74 19.67 18.11 18.71 19.5ZM13 3.5C13.73 2.67 14.94 2.04 15.94 2C16.07 3.17 15.6 4.35 14.9 5.19C14.21 6.04 13.07 6.7 11.95 6.61C11.8 5.46 12.36 4.26 13 3.5Z" fill="currentColor"/>
          </svg>
          Sign in with Apple
        </button>

        <div className="flex items-center gap-3 my-4">
          <div className="h-px bg-gray-700 flex-1"></div>
          <span className="text-gray-500">or</span>
          <div className="h-px bg-gray-700 flex-1"></div>
        </div>

        <input
          type="text"
          placeholder="Phone, email, or username"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-3 focus:outline-none focus:border-blue-500"
        />

        <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 mb-3 hover:bg-gray-200 transition-colors">
          Next
        </button>

        <button className="w-full border border-gray-700 text-white font-semibold rounded-full py-2.5 px-4 mb-6 hover:bg-white/10 transition-colors">
          Forgot password?
        </button>

        <p className="text-gray-500 text-center">
          Don't have an account?{' '}
          <a href="/i/flow/signup" className="text-blue-500 hover:underline">
            Sign up
          </a>
        </p>
      </div>
    </div>
  );
}
