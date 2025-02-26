// components/LoginModal.tsx
import React, { useState } from 'react';
import { X, ChevronDown } from 'lucide-react';
import { useNavigate } from 'react-router-dom';


export function SignupModal() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [nameCharCount, setNameCharCount] = useState(0);
  const [showDateFields, setShowDateFields] = useState(false);
  const [month, setMonth] = useState('Month');
  const [day, setDay] = useState('Day');
  const [year, setYear] = useState('Year');
  const navigate = useNavigate();

  function handleClose() {
    navigate('/');
  }

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setName(value);
    setNameCharCount(value.length);
  };

  const toggleDateFields = () => {
    setShowDateFields(!showDateFields);
  };

  return (
    <div className="fixed inset-0 bg-[#242D34]/60 flex items-center justify-center p-4 z-50">
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
          Create your account
        </h1>

        {/* First version of the signup form */}
        <div className={showDateFields ? "hidden" : "block"}>
          <div className="relative mb-4">
            <input
              type="text"
              placeholder="Name"
              value={name}
              onChange={handleNameChange}
              maxLength={50}
              className="w-full bg-black text-white border border-blue-500 rounded-md px-4 py-3 focus:outline-none focus:border-blue-600"
            />
            <span className="absolute right-3 top-3 text-gray-500 text-sm">
              {nameCharCount} / 50
            </span>
          </div>

          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-4 focus:outline-none focus:border-blue-500"
          />

          <div className="mb-6">
            <h3 className="text-white font-bold mb-1">Date of birth</h3>
            <p className="text-xs text-gray-500 mb-3">
              This will not be shown publicly. Confirm your own age, even if this account is for a business, a pet, or something else.
            </p>

            <div className="flex gap-2">
              <div className="relative flex-1">
                <select
                  className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 appearance-none focus:outline-none focus:border-blue-500"
                  value={month}
                  onChange={(e) => setMonth(e.target.value)}
                >
                  <option>Month</option>
                  <option value="1">January</option>
                  <option value="2">February</option>
                  <option value="3">March</option>
                  <option value="4">April</option>
                  <option value="5">May</option>
                  <option value="6">June</option>
                  <option value="7">July</option>
                  <option value="8">August</option>
                  <option value="9">September</option>
                  <option value="10">October</option>
                  <option value="11">November</option>
                  <option value="12">December</option>
                </select>
                <div className="absolute right-3 top-3 pointer-events-none text-gray-500">
                  <ChevronDown size={18} />
                </div>
              </div>

              <div className="relative w-24">
                <select
                  className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 appearance-none focus:outline-none focus:border-blue-500"
                  value={day}
                  onChange={(e) => setDay(e.target.value)}
                >
                  <option>Day</option>
                  {Array.from({ length: 31 }, (_, i) => i + 1).map(day => (
                    <option key={day} value={day}>{day}</option>
                  ))}
                </select>
                <div className="absolute right-3 top-3 pointer-events-none text-gray-500">
                  <ChevronDown size={18} />
                </div>
              </div>

              <div className="relative w-28">
                <select
                  className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 appearance-none focus:outline-none focus:border-blue-500"
                  value={year}
                  onChange={(e) => setYear(e.target.value)}
                >
                  <option>Year</option>
                  {Array.from({ length: 100 }, (_, i) => new Date().getFullYear() - i).map(year => (
                    <option key={year} value={year}>{year}</option>
                  ))}
                </select>
                <div className="absolute right-3 top-3 pointer-events-none text-gray-500">
                  <ChevronDown size={18} />
                </div>
              </div>
            </div>
          </div>

          <button
            onClick={toggleDateFields}
            className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 mb-6 hover:bg-gray-200 transition-colors"
          >
            Next
          </button>
        </div>

        {/* Second version of the signup form */}
        <div className={showDateFields ? "block" : "hidden"}>
          <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 flex items-center justify-center gap-2 mb-3 hover:bg-gray-200 transition-colors">
            <img src="https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_24x24dp.png"
              alt="Google"
              className="w-5 h-5" />
            Sign up with Google
          </button>

          <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 flex items-center justify-center gap-2 mb-3 hover:bg-gray-200 transition-colors">
            <svg viewBox="0 0 24 24" className="w-5 h-5">
              <path d="M18.71 19.5C17.88 20.74 17 21.95 15.66 21.97C14.32 22 13.89 21.18 12.37 21.18C10.84 21.18 10.37 21.95 9.09997 22C7.78997 22.05 6.79997 20.68 5.95997 19.47C4.24997 17 2.93997 12.45 4.69997 9.39C5.56997 7.87 7.12997 6.91 8.81997 6.88C10.1 6.86 11.32 7.75 12.11 7.75C12.89 7.75 14.37 6.68 15.92 6.84C16.57 6.87 18.39 7.1 19.56 8.82C19.47 8.88 17.39 10.1 17.41 12.63C17.44 15.65 20.06 16.66 20.09 16.67C20.06 16.74 19.67 18.11 18.71 19.5ZM13 3.5C13.73 2.67 14.94 2.04 15.94 2C16.07 3.17 15.6 4.35 14.9 5.19C14.21 6.04 13.07 6.7 11.95 6.61C11.8 5.46 12.36 4.26 13 3.5Z" fill="currentColor" />
            </svg>
            Sign up with Apple
          </button>

          <div className="flex items-center gap-3 my-4">
            <div className="h-px bg-gray-700 flex-1"></div>
            <span className="text-gray-500">or</span>
            <div className="h-px bg-gray-700 flex-1"></div>
          </div>

          <input
            type="text"
            placeholder="Name"
            value={name}
            onChange={handleNameChange}
            className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-3 focus:outline-none focus:border-blue-500"
          />

          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-3 focus:outline-none focus:border-blue-500"
          />

          <input
            type="password"
            placeholder="Password"
            className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-3 focus:outline-none focus:border-blue-500"
          />

          <div className="mb-6">
            <p className="text-xs text-gray-500 mb-2">
              By signing up, you agree to the Terms of Service and Privacy Policy, including Cookie Use. Others will be able to find you by email or phone number when provided.
            </p>
          </div>

          <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 mb-6 hover:bg-gray-200 transition-colors">
            Sign up
          </button>

          <p className="text-gray-500 text-center">
            Already have an account?{' '}
            <a href="#" className="text-blue-500 hover:underline">
              Sign in
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}



