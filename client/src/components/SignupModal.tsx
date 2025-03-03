import React, { useState, useEffect } from 'react';
import { X, ChevronDown } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export function SignupModal() {
  const [isLoading, setIsLoading] = useState(true);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [nameCharCount, setNameCharCount] = useState(0);
  const [month, setMonth] = useState('Month');
  const [day, setDay] = useState('Day');
  const [year, setYear] = useState('Year');
  const [emailError, setEmailError] = useState('');
  const [birthDateError, setBirthDateError] = useState('');
  const [isBirthDateValid, setIsBirthDateValid] = useState(false);
  const [showErrorAlert, setShowErrorAlert] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    // Simulate a delay for loading
    const timer = setTimeout(() => setIsLoading(false), 200);
    return () => clearTimeout(timer);
  }, []);

  // Auto-hide error alert after 5 seconds
  useEffect(() => {
    if (showErrorAlert) {
      const timer = setTimeout(() => {
        setShowErrorAlert(false);
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [showErrorAlert]);

  const handleClose = () => navigate('/');

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setName(value);
    setNameCharCount(value.length);
  };

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const handleEmailChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setEmail(value);

    if (!validateEmail(value)) {
      setEmailError('Please enter a valid email address');
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/search?q=${value}&f=unique-mail`);
      const data = await response.json();

      if (data.exists) {
        setEmailError('Email is already in use');
      } else {
        setEmailError('');
      }
    } catch (error) {
      console.error('Error checking email uniqueness:', error);
      setEmailError('An error occurred. Please try again later.');
    }
  };

  const validatePassword = (password: string) => {
    if (password.length < 8) {
      setPasswordError('Password must be at least 8 characters');
      return false;
    }

    if (password.length > 72) {
      setPasswordError('Password must not be greater than 72 characters');
      return false;
    }

    if (!/[A-Z]/.test(password)) {
      setPasswordError('Password must contain at least one uppercase letter');
      return false;
    }

    if (!/[a-z]/.test(password)) {
      setPasswordError('Password must contain at least one lowercase letter');
      return false;
    }

    if (!/[0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
      setPasswordError('Password must contain at least one number or special character');
      return false;
    }

    if (password === password.toUpperCase()) {
      setPasswordError('Password must contain at least one non-uppercase character');
      return false;
    }

    setPasswordError('');
    return true;
  };
  
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    validatePassword(value);
  };

  const validateBirthDate = () => {
    if (month === 'Month' || day === 'Day' || year === 'Year') {
      setBirthDateError('Please select a valid birth date');
      setIsBirthDateValid(false);
      return false;
    }

    const birthDate = new Date(`${year}-${month}-${day}`);
    const age = new Date().getFullYear() - birthDate.getFullYear();

    if (age < 13) {
      setBirthDateError('You must be at least 13 years old');
      setIsBirthDateValid(false);
      return false;
    }

    setBirthDateError('');
    setIsBirthDateValid(true);
    return true;
  };

  // Generate a username that fits within the 15 character limit
  const generateUsername = (baseName: string) => {
    // Clean the base name to only include valid characters (letters, numbers, underscores)
    const cleanedName = baseName.replace(/[^a-zA-Z0-9_]/g, '');
    
    // Generate exactly 7 random digits
    const randomDigits = Math.floor(Math.random() * 10000000).toString().padStart(7, '0');
    
    // Calculate how much space we have left for the base name
    const maxBaseNameLength = 15 - randomDigits.length;
    
    // Trim the base name if needed to fit within the limit
    const trimmedBaseName = cleanedName.substring(0, maxBaseNameLength);
    
    // Combine the trimmed base name with the random digits
    return trimmedBaseName + randomDigits;
  };

  // Check if username is unique and keep trying until it is
  const getUniqueUsername = async (baseName: string) => {
    let isUnique = false;
    let username = '';
    
    while (!isUnique) {
      username = generateUsername(baseName);
      try {
        const response = await fetch(`http://localhost:8080/api/search?q=${username}&f=unique-user`);
        const data = await response.json();
        
        if (!data.exists) {
          isUnique = true;
        }
      } catch (error) {
        console.error('Error checking username uniqueness:', error);
        // If there's an error, we'll just use the generated username
        isUnique = true;
      }
    }
    
    return username;
  };

  const handleSubmit = async () => {
    if (!validateEmail(email) || !validatePassword(password) || !validateBirthDate()) return;

    try {
      // Generate a unique username
      const uniqueUsername = await getUniqueUsername(name);
      
      const signupResponse = await fetch("http://localhost:8080/api/i/flow/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
          nickname: name, 
          username: uniqueUsername, 
          mail: email, 
          password 
        }),
        credentials: "include",
      });

      if (!signupResponse.ok) {
        const errorData = await signupResponse.json();
        throw new Error(errorData.error || "Signup failed");
      }

      navigate("/i/flow/login");
    } catch (error) {
      console.error("Signup failed:", error);
    }
  };

  if (isLoading) {
    return (
      <div className="fixed inset-0 bg-black flex items-center justify-center z-50">
        <div className="w-8 h-8 border-t-2 border-blue-500 rounded-full animate-spin"></div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-[#242D34]/60 flex items-center justify-center z-50">
      <div className="bg-black w-full max-w-xl rounded-2xl relative p-2">
        <button onClick={handleClose} className="absolute left-4 top-4 text-gray-400 hover:text-white">
          <X size={20} />
        </button>

        <div className="flex justify-center mb-8">
          <X size={50} className="text-white" />
        </div>
        <div className="px-24">
          <h1 className="text-4xl font-bold text-white text-left mb-8">Create your account</h1>

          <div>
            <div className="relative mb-4">
              <input
                type="text"
                placeholder="Name"
                value={name}
                onChange={handleNameChange}
                maxLength={50}
                className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 focus:outline-none focus:border-blue-500"
              />
              <span className="absolute right-3 top-3 text-gray-500 text-sm">{nameCharCount} / 50</span>
            </div>

            <input
              type="email"
              placeholder="Email"
              value={email}
              onChange={handleEmailChange}
              className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-4 focus:outline-none focus:border-blue-500"
            />
            {emailError && <p className="text-red-500 text-sm mb-4">{emailError}</p>}

            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={handlePasswordChange}
              className="w-full bg-black text-white border border-gray-700 rounded-md px-4 py-3 mb-4 focus:outline-none focus:border-blue-500"
            />
            {passwordError && <p className="text-red-500 text-sm mb-4">{passwordError}</p>}

            <div className="mb-12">
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
                    onBlur={validateBirthDate}
                  >
                    <option>Month</option>
                    {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
                      <option key={m} value={m}>{new Date(0, m - 1).toLocaleString('default', { month: 'long' })}</option>
                    ))}
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
                    onBlur={validateBirthDate}
                  >
                    <option>Day</option>
                    {Array.from({ length: 31 }, (_, i) => i + 1).map((d) => (
                      <option key={d} value={d}>{d}</option>
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
                    onBlur={validateBirthDate}
                  >
                    <option>Year</option>
                    {Array.from({ length: 100 }, (_, i) => new Date().getFullYear() - i).map((y) => (
                      <option key={y} value={y}>{y}</option>
                    ))}
                  </select>
                  <div className="absolute right-3 top-3 pointer-events-none text-gray-500">
                    <ChevronDown size={18} />
                  </div>
                </div>
              </div>
              {birthDateError && <p className="text-red-500 text-sm mb-4">{birthDateError}</p>}
            </div>
          </div>

          <button
            onClick={handleSubmit}
            disabled={!!emailError || !validateEmail(email) || !password || !isBirthDateValid}
            className={`w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 mb-6 transition-colors ${!!emailError || !validateEmail(email) || !password || !isBirthDateValid ? "opacity-50 cursor-not-allowed" : "hover:bg-gray-200"}`}
          >
            Sign up
          </button>
        </div>
      </div>
    </div>
  );
}

export default SignupModal;