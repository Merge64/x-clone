import { useState, useEffect } from "react";
import { X, Eye, EyeOff } from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import { checkExists, login, validateToken } from "../../utils/auth";

function FloatingLabelInput({
  label,
  type,
  value,
  onChange,
  showToggle = false,
  error = "",
  onToggleVisibility = () => { },
  showPassword = false,
}: {
  label: string;
  type: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  showToggle?: boolean;
  error?: string;
  onToggleVisibility?: () => void;
  showPassword?: boolean;
}) {
  const [isFocused, setIsFocused] = useState(false);

  const inputType = type === 'password' && showPassword ? 'text' : type;

  return (
    <div className="relative">
      <div className={`relative border ${error ? 'border-red-500' : 'border-gray-700'} rounded-md focus-within:border-blue-500`}>
        <label
          className={`absolute transition-all duration-200 pointer-events-none ${isFocused || value
            ? 'text-xs text-gray-400 top-2 left-3'
            : 'text-base text-gray-500 top-1/2 -translate-y-1/2 left-3'
            }`}
        >
          {label}
        </label>

        <input
          type={inputType}
          value={value}
          onChange={onChange}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          className={`w-full bg-black text-white rounded-md px-3 pt-6 pb-2 focus:outline-none ${showToggle ? 'pr-10' : ''
            }`}
        />

        {showToggle && (
          <button
            type="button"
            onClick={onToggleVisibility}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-white"
          >
            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
          </button>
        )}
      </div>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
}

function LoginModal() {
  const [emailOrUsername, setEmailOrUsername] = useState("");
  const [password, setPassword] = useState("");
  const [passwordError, ] = useState('');
  const [loginError, setLoginError] = useState('');
  const [step, setStep] = useState(1);
  const [, setUserExists] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isCheckingToken, setIsCheckingToken] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  // Check for valid token on component mount
  useEffect(() => {
    const checkAuthToken = async () => {
      try {
        const isValid = await validateToken();
        if (isValid) {
          // If token is valid, redirect to home or the intended destination
          const from = location.state?.from?.pathname || "/home";
          navigate(from, { replace: true });
        }
      } catch (error) {
        console.error("Error validating token:", error);
        // Token validation failed, continue with login flow
      } finally {
        setIsCheckingToken(false);
      }
    };

    checkAuthToken();
  }, [navigate, location.state?.from?.pathname]);

  const handleClose = () => navigate("/");
  
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
  };

  const checkUserExists = async () => {
    if (!emailOrUsername.trim()) {
      setLoginError("Please enter an email or username");
      return;
    }

    setIsLoading(true);
    try {
      // Check if username or email exists
      const isUsernameExists = await checkExists(emailOrUsername, 'user');
      const isEmailExists = await checkExists(emailOrUsername, 'mail');

      if (isUsernameExists || isEmailExists) {
        setUserExists(true);
        setStep(2);
        setLoginError('');
      } else {
        setLoginError("No account found with that username or email");
      }
    } catch (error) {
      console.error("Error checking user:", error);
      setLoginError("Something went wrong. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async () => {
    if (step === 1) {
      await checkUserExists();
      return;
    }

    setIsLoading(true);
    try {
      const result = await login(emailOrUsername, password);
      
      if (!result.success) {
        setLoginError(result.error || "Login failed");
        return;
      }
      
      // Get the redirect path from location state or default to /home
      const from = location.state?.from?.pathname || "/home";
      navigate(from, { replace: true });
    } catch (error) {
      console.error("Login failed:", error);
      setLoginError("Something went wrong during login. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  // Show loading indicator while checking token
  if (isCheckingToken) {
    return (
      <div className="fixed inset-0 bg-[#242D34]/60 bg-opacity-50 flex items-center justify-center p-4">
        <div className="bg-black w-full max-w-md rounded-2xl p-8 flex flex-col items-center justify-center">
          <div className="w-8 h-8 border-t-2 border-blue-500 rounded-full animate-spin mb-4"></div>
          <p className="text-white">Checking login status...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-[#242D34]/60 bg-opacity-50 flex items-center justify-center p-4">
      {step === 1 && (
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
            Sign in to X-clone
          </h1>

          {loginError && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
              {loginError}
            </div>
          )}

          {/* Social login buttons */}
          <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 flex items-center justify-center gap-2 mb-3 hover:bg-gray-200 transition-colors">
            <img src="https://www.google.com/favicon.ico" alt="Google" className="w-5 h-5" />
            Sign in with Google
          </button>

          <button className="w-full bg-white text-black font-semibold rounded-full py-2.5 px-4 flex items-center justify-center gap-2 mb-3 hover:bg-gray-200 transition-colors">
            <svg viewBox="0 0 24 24" className="w-5 h-5">
              <path d="M18.71 19.5C17.88 20.74 17 21.95 15.66 21.97C14.32 22 13.89 21.18 12.37 21.18C10.84 21.18 10.37 21.95 9.09997 22C7.78997 22.05 6.79997 20.68 5.95997 19.47C4.24997 17 2.93997 12.45 4.69997 9.39C5.56997 7.87 7.12997 6.91 8.81997 6.88C10.1 6.86 11.32 7.75 12.11 7.75C12.89 7.75 14.37 6.68 15.92 6.84C16.57 6.87 18.39 7.1 19.56 8.82C19.47 8.88 17.39 10.1 17.41 12.63C17.44 15.65 20.06 16.66 20.09 16.67C20.06 16.74 19.67 18.11 18.71 19.5ZM13 3.5C13.73 2.67 14.94 2.04 15.94 2C16.07 3.17 15.6 4.35 14.9 5.19C14.21 6.04 13.07 6.7 11.95 6.61C11.8 5.46 12.36 4.26 13 3.5Z" fill="currentColor" />
            </svg>
            Sign in with Apple
          </button>

          <div className="flex items-center gap-3 my-4">
            <div className="h-px bg-gray-700 flex-1"></div>
            <span className="text-gray-500">or</span>
            <div className="h-px bg-gray-700 flex-1"></div>
          </div>

          {/* Step 1: Email/Username input */}
          <div className="mb-6">
            <FloatingLabelInput
              label="Email or username"
              type="text"
              value={emailOrUsername}
              onChange={(e) => setEmailOrUsername(e.target.value)}
            />
          </div>

          <button
            onClick={handleSubmit}
            disabled={isLoading}
            className={`w-full ${isLoading ? 'bg-gray-400' : 'bg-white hover:bg-gray-200'} text-black font-semibold rounded-full py-2.5 px-4 mb-6 transition-colors flex items-center justify-center`}
          >
            {isLoading ? (
              <>
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-black" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Checking...
              </>
            ) : (
              'Next'
            )}
          </button>

          <p className="text-gray-500 text-center">
            Don't have an account?{" "}
            <a href="/i/flow/signup" className="text-blue-500 hover:underline">
              Sign up
            </a>
          </p>
        </div>
      )}

      {/* Step 2: Password input */}
      {step === 2 && (
        <div className="bg-black w-full max-w-md rounded-2xl p-8 relative flex flex-col items-center">
          <button
            onClick={handleClose}
            className="absolute left-4 top-4 text-gray-400 hover:text-white"
          >
            <X size={20} />
          </button>

          <div className="flex justify-center mb-6">
            <X size={30} className="text-white" />
          </div>

          <h1 className="text-2xl font-bold text-white text-center mb-8">
            Enter your password
          </h1>

          {loginError && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md w-full">
              {loginError}
            </div>
          )}

          <div className="w-full mb-6 bg-c">
            <div className="w-full text-white border bg-[#101214] border-black rounded-md px-4 py-2 pr-10 focus:outline-none focus:border-blue-500 mb-4">
              <p className="text-sm text-white text-opacity-15">Username</p>
              <p className="text-white text-opacity-15">{emailOrUsername}</p>
            </div>

            <div className="relative w-full">
              <FloatingLabelInput
                label="Password"
                type="password"
                value={password}
                onChange={handlePasswordChange}
                showToggle={true}
                error={passwordError}
                onToggleVisibility={togglePasswordVisibility}
                showPassword={showPassword}
              />
            </div>

            <a href="/i/flow/signup" className="text-blue-500 hover:underline text-sm block mt-2">
              Forgot password?
            </a>
          </div>

          <button
            onClick={handleSubmit}
            disabled={isLoading || passwordError !== ''}
            className={`w-full ${isLoading || passwordError !== ''
                ? 'bg-[#787A7A] cursor-not-allowed text-gray-500'
                : 'bg-white hover:bg-[#D7DBDC] text-black'
              } font-semibold rounded-full py-2.5 px-4 mb-6 transition-colors flex items-center justify-center`}
          >
            {isLoading ? (
              <>
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Logging in...
              </>
            ) : (
              'Log in'
            )}
          </button>

          <p className="text-gray-500 text-center">
            Don't have an account?{" "}
            <a href="/i/flow/signup" className="text-blue-500 hover:underline">
              Sign up
            </a>
          </p>
        </div>
      )}
    </div>
  );
}

export default LoginModal;