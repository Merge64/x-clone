import React, { useState, useEffect } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { validateToken } from './auth';

interface AuthCheckerProps {
  children: React.ReactNode;
}

function AuthChecker({ children }: AuthCheckerProps) {
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const location = useLocation();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const isValid = await validateToken();
        setIsAuthenticated(isValid);
      } catch (error) {
        console.error('Error validating token:', error);
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  if (isLoading) {
    return (
      <div className="fixed inset-0 bg-black flex items-center justify-center z-50">
        <div className="w-8 h-8 border-t-2 border-blue-500 rounded-full animate-spin"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    // Redirect to login page with the return url
    return <Navigate to="/" state={{ from: location }} replace />;
  }

  return <>{children}</>;
}

export default AuthChecker;