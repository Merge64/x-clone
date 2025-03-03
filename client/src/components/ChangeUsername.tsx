import React, { useState, useEffect, useCallback } from 'react';
import { X } from 'lucide-react';
import { getUserInfo, updateUsername } from '../utils/api';
import { checkExists } from '../utils/auth';
import { useNavigate } from 'react-router-dom';

interface ChangeUsernameProps {
  onSkip?: () => void;
}

function ChangeUsername({ onSkip }: ChangeUsernameProps) {
const [username, setUsername] = useState('');
const [originalUsername, setOriginalUsername] = useState(''); // Store the original username
const [isLoading, setIsLoading] = useState(true);
const [error, setError] = useState('');
const [isValid, setIsValid] = useState(true); // Start with true to avoid initial error
const [isAvailable, setIsAvailable] = useState(true);
const [isChecking, setIsChecking] = useState(false);
const [initialLoad, setInitialLoad] = useState(true); // Track initial load
const navigate = useNavigate();

useEffect(() => {
// Fetch user info including the randomly generated username
const fetchUserInfo = async () => {
try {
setIsLoading(true);
const userInfo = await getUserInfo();
setUsername(userInfo.username);
setOriginalUsername(userInfo.username); // Store the original username
// Set initial validation state based on the fetched username
setIsValid(/^[a-zA-Z0-9_]+$/.test(userInfo.username));
setInitialLoad(false); // Mark initial load as complete
} catch (err) {
console.error('Error fetching user info:', err);
setError('Failed to load user information');
setInitialLoad(false); // Mark initial load as complete even on error
} finally {
setIsLoading(false);
}
};


fetchUserInfo();
}, []);

// Use useCallback to memoize the checkUsernameAvailability function
const checkUsernameAvailability = useCallback(async () => {
// If username is the same as original, no need to check availability
if (username === originalUsername) {
setIsAvailable(true);
return;
}


if (!username || username.length > 15) {
  setIsAvailable(false);
  return;
}

setIsChecking(true);
try {
  // Using the checkExists function to see if username already exists
  const exists = await checkExists(username, 'user');
  setIsAvailable(!exists);
} catch (err) {
  console.error('Error checking username availability:', err);
  setIsAvailable(false);
} finally {
  setIsChecking(false);
}
}, [username, originalUsername]);

useEffect(() => {
// Check username availability with debounce
if (!username || initialLoad) return;


// Set isChecking to true immediately when typing starts
if (username !== originalUsername) {
  setIsChecking(true);
}

const debounceTimer = setTimeout(checkUsernameAvailability, 500);
return () => clearTimeout(debounceTimer);
}, [username, checkUsernameAvailability, initialLoad, originalUsername]);

const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
const newUsername = e.target.value;


// Only allow letters, numbers, and underscores
if (newUsername && !/^[a-zA-Z0-9_]*$/.test(newUsername)) {
  return;
}

// Enforce maximum length of 15 characters
if (newUsername.length > 15) {
  return;
}

setUsername(newUsername);

// Simple validation for format - availability is checked in the useEffect
// Username can be any length up to 15 characters, but must contain only valid characters
const isFormatValid = newUsername.length === 0 || /^[a-zA-Z0-9_]+$/.test(newUsername);
setIsValid(isFormatValid);

// Set checking state to true immediately when the username changes
if (newUsername !== originalUsername && newUsername.length > 0) {
  setIsChecking(true);
}
};

const handleSkip = () => {
  // Use the onSkip prop if provided, otherwise navigate
  if (onSkip) {
    onSkip();
  } else {
    // Navigate to the home feed
    navigate('/home');
  }
};

const handleSubmit = async () => {
if (!isValid || !isAvailable || isChecking) return;


try {
  setIsLoading(true);
  await updateUsername(username);
  // If onSkip is provided, call it after successful update
  if (onSkip) {
    onSkip();
  } else {
    // Redirect to home feed after successful update
    navigate('/home');
  }
} catch (err) {
  console.error('Error updating username:', err);
  setError('Failed to update username');
} finally {
  setIsLoading(false);
}
};

if (isLoading) {
return (
<div className="flex items-center justify-center p-6">
<div className="w-8 h-8 border-t-2 border-blue-500 rounded-full animate-spin"></div>
</div>
);
}

// Determine if the Next button should be visible and enabled
const isSameAsOriginal = username === originalUsername;
const canProceed = isValid && isAvailable && !isChecking && username.length > 0 && !isSameAsOriginal;

// Determine what status message to show
let statusMessage = null;
let statusColor = "";

if (error) {
statusMessage = error;
statusColor = "text-red-500";
} else if (!isValid && username.length > 0) {
statusMessage = "Username can only contain letters, numbers, and underscores.";
statusColor = "text-red-500";
} else if (isValid && !isAvailable && !isChecking && username.length > 0 && !isSameAsOriginal) {
statusMessage = "This username is already taken. Please choose another one.";
statusColor = "text-red-500";
} else if (isChecking && username !== originalUsername && username.length > 0) {
statusMessage = "Checking username availability...";
statusColor = "text-blue-500";
}

return (
<div className="bg-black text-white p-6 rounded-lg">
<div className="flex justify-between mb-6">
<button onClick={handleSkip} className="text-gray-400 hover:text-white">
<X size={24} />
</button>
</div>


    <h1 className="text-3xl font-bold mb-2">What should we call you?</h1>
    <p className="text-gray-400 mb-6">Your @username is unique. You can always change it later.</p>
    
    <div className="relative mb-4">
      <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none text-gray-500">
        @
      </div>
      <input
        type="text"
        value={username}
        onChange={handleUsernameChange}
        className="w-full bg-black border border-gray-700 rounded-lg py-3 px-10 text-white focus:outline-none focus:border-blue-500"
        placeholder="Username"
        maxLength={15}
      />
      <div className="absolute inset-y-0 right-3 flex items-center pointer-events-none text-gray-500">
        {username.length}/15
      </div>
    </div>
    
    {/* Status message area with fixed height to prevent layout shifts */}
    <div className="h-6 mb-4">
      {statusMessage && <p className={statusColor}>{statusMessage}</p>}
    </div>
    
    {/* Fixed height container for Next button to prevent layout shifts */}
    <div className="h-14 mb-3">
      {/* Next button - only visible when username is different from original */}
      {!isSameAsOriginal && (
        <button
          onClick={handleSubmit}
          disabled={!canProceed}
          className={`w-full py-3 rounded-full font-bold transition-all duration-200 ${
            isChecking 
              ? 'bg-gray-500 text-black cursor-not-allowed' 
              : canProceed 
                ? 'bg-white text-black hover:bg-[#D7DBDC]' 
                : 'bg-white text-black opacity-50 cursor-not-allowed'
          }`}
        >
            Next
        </button>
      )}
    </div>
    
    <button
      onClick={handleSkip}
      className="w-full py-3 rounded-full border border-gray-700 text-white font-bold hover:bg-gray-900"
    >
      Skip for now
    </button>
  </div>
);
}

export default ChangeUsername;