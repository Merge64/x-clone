// components/MainBackground.tsx
import React from 'react';
import { X } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

function MainBackground() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-black text-white flex">
      {/* Left side - Logo */}
      <div className="flex flex-1 items-center justify-center max-w-[50vw]">
        <X className="w-96 h-96 text-white" strokeWidth={1} />
      </div>

      {/* Right side - Login / Sign Up */}
      <div className="flex-0 p-10 flex flex-col justify-center ml-8 max-w-[50vw]">
        <h1 className="text-6xl font-bold mb-2">Happening now</h1>
        <h2 className="text-3xl font-bold mt-8 mb-8">Join today.</h2>

        {/* Sign up buttons */}
        <div className="w-full max-w-[300px]">
          <button className="w-full bg-white text-black rounded-full py-2 px-4 font-bold flex items-center justify-center gap-2 mb-3 min-w-[300px]">
            <img src="https://www.google.com/favicon.ico" alt="Google" className="w-5 h-5" />
            Sign up with Google
          </button>
          <button className="w-full bg-white text-black rounded-full py-2 px-4 font-bold flex items-center justify-center gap-2 min-w-[300px]">
          <svg viewBox="0 0 24 24" className="w-5 h-5">
              <path d="M16.365 1.43c0 1.14-.493 2.27-1.177 3.08-.744.9-1.99 1.57-2.987 1.57-.12 0-.23-.02-.3-.03-.01-.06-.04-.22-.04-.39 0-1.15.572-2.27 1.206-2.98.804-.94 2.142-1.64 3.248-1.68.03.13.05.28.05.43zm4.565 15.71c-.03.07-.463 1.58-1.518 3.12-.945 1.34-1.94 2.71-3.43 2.71-1.517 0-1.9-.88-3.63-.88-1.698 0-2.302.91-3.67.91-1.377 0-2.332-1.26-3.428-2.8-1.287-1.82-2.323-4.63-2.323-7.28 0-4.28 2.797-6.55 5.552-6.55 1.448 0 2.675.95 3.6.95.865 0 2.222-1.01 3.902-1.01.613 0 2.886.06 4.374 2.19-.13.09-2.383 1.37-2.383 4.19 0 3.26 2.854 4.42 2.955 4.45z" fill="currentColor"/>
            </svg>
            Sign up with Apple
          </button>
        </div>

        {/* Divider */}
        <div className="flex items-center max-w-[300px]">
          <div className="flex-1 border-t border-[#2F3336]" />
          <span className="px-4 my-2">or</span>
          <div className="flex-1 border-t border-[#2F3336]" />
        </div>

        {/* Create account button */}
        <button
          onClick={() => navigate('/i/flow/signup')}
          className="bg-[#1d9bf0] text-white rounded-full py-2 px-4 font-bold mb-3 max-w-[300px] min-w-[300px]"
        >
          Create account
        </button>

        <p className="text-xs text-gray-500 mb-8 max-w-[300px]">
          By signing up, you agree to the <a href="#" className="text-[#1d9bf0]">Terms of Service</a> and{' '}
          <a href="#" className="text-[#1d9bf0]">Privacy Policy</a>, including{' '}
          <a href="#" className="text-[#1d9bf0]">Cookie Use</a>.
        </p>

        {/* Sign in button */}
        <div className="mt-10 max-w-[300px]">
          <p className="font-bold mb-5">Already have an account?</p>
          <button
            onClick={() => navigate('/i/flow/login')}
            className="w-full border border-gray-600 text-[#1d9bf0] rounded-full py-2 px-4 font-bold hover:bg-[#1d9bf0]/10 min-w-[300px]"
          >
            Sign in
          </button>
        </div>
      </div>

      {/* Footer */}
      <footer className="fixed bottom-0 left-0 right-0 p-4 flex flex-wrap justify-center gap-x-4 text-sm text-gray-500">
        <a href="#" className="hover:underline">About</a>
        <a href="#" className="hover:underline">Download the X app</a>
        <a href="#" className="hover:underline">Help Center</a>
        <a href="#" className="hover:underline">Terms of Service</a>
        <a href="#" className="hover:underline">Privacy Policy</a>
        <a href="#" className="hover:underline">Cookie Policy</a>
        <a href="#" className="hover:underline">Accessibility</a>
        <a href="#" className="hover:underline">Ads info</a>
        <a href="#" className="hover:underline">Blog</a>
        <a href="#" className="hover:underline">Careers</a>
        <a href="#" className="hover:underline">Brand Resources</a>
        <a href="#" className="hover:underline">Advertising</a>
        <a href="#" className="hover:underline">Marketing</a>
        <a href="#" className="hover:underline">X for Business</a>
        <a href="#" className="hover:underline">Developers</a>
        <a href="#" className="hover:underline">Directory</a>
        <a href="#" className="hover:underline">Settings</a>
        <span>© 2025 X-clone Corp.</span>
      </footer>
    </div>
  );
}

export default MainBackground;
