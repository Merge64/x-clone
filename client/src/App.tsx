import React from 'react';
import { X } from 'lucide-react';

function App() {
  return (
    <div className="min-h-screen bg-black text-white flex">
      {/* Left side - Logo */}
      <div className="flex flex-1 items-center justify-center max-w-[50vw]">
        <X className="w-96 h-96 text-white" strokeWidth={1} />
      </div>

      {/* Right side - Login Form */}
      <div className="flex-0 p-10 flex flex-col justify-center ml-8 max-w-[50vw]">
        <h1 className="text-6xl font-bold mb-2">Happening now</h1>
        <h2 className="text-3xl font-bold mt-8 mb-8">Join today.</h2>

        {/* Sign up buttons */}
        <div className="w-full max-w-[300px]">
          <button className="w-full bg-white text-black rounded-full py-2 px-4 font-bold flex items-center justify-center gap-2 mb-3 min-w-[300px]">
            <img src="https://www.google.com/favicon.ico" alt="" className="w-5 h-5" />
            Sign up with Google
          </button>
          <button className="w-full bg-white text-black rounded-full py-2 px-4 font-bold flex items-center justify-center gap-2 mb-3 min-w-[300px]">
            <img src="https://www.apple.com/favicon.ico" alt="" className="w-5 h-5" />
            Sign up with Apple
          </button>
        </div>

        {/* Divider */}
        <div className="flex items-center my-4 max-w-[300px]">
          <div className="flex-1 border-t border-gray-600"></div>
          <span className="px-4">or</span>
          <div className="flex-1 border-t border-gray-600"></div>
        </div>

        <button className="bg-[#1d9bf0] text-white rounded-full py-2 px-4 font-bold mb-3 max-w-[300px] min-w-[300px]">
          Create account
        </button>

        <p className="text-xs text-gray-500 mb-8 max-w-[300px]">
          By signing up, you agree to the{' '}
          <a href="#" className="text-[#1d9bf0]">Terms of Service</a> and{' '}
          <a href="#" className="text-[#1d9bf0]">Privacy Policy</a>, including{' '}
          <a href="#" className="text-[#1d9bf0]">Cookie Use</a>.
        </p>

        <div className="mt-10 max-w-[300px]">
          <p className="font-bold mb-5">Already have an account?</p>
          <button className="w-full border border-gray-600 text-[#1d9bf0] rounded-full py-2 px-4 font-bold hover:bg-[#1d9bf0]/10 min-w-[300px]">
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

export default App;