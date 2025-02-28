// This file contains API functions to interact with your backend

// Function to get the JWT token from cookies or localStorage
export function getToken(): string | null {
    // If you're using cookies, the token is automatically sent with requests
    // If you're using localStorage:
    return localStorage.getItem('Authentication');
  }
  
  // Function to decode JWT token and extract user information
  export function decodeToken(token: string): any {
    try {
      // JWT tokens are split into three parts: header.payload.signature
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      );
      return JSON.parse(jsonPayload);
    } catch (error) {
      console.error('Error decoding token:', error);
      return null;
    }
  }
  
  // Function to get user information from the backend
  export async function getUserInfo(): Promise<any> {
    try {
      const response = await fetch('http://localhost:8080/api/user/info', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Include cookies in the request
      });
  
      if (!response.ok) {
        throw new Error('Failed to fetch user info');
      }
  
      return await response.json();
    } catch (error) {
      console.error('Error fetching user info:', error);
      
      // Fallback: If the API call fails, try to get info from the JWT token
      const token = getToken();
      if (token) {
        const decodedToken = decodeToken(token);
        if (decodedToken) {
          return {
            username: decodedToken.username || decodedToken.preferred_username || 'user' + Math.floor(Math.random() * 10000),
            // Add other user properties from the token if available
          };
        }
      }
      
      throw error;
    }
  }
  
  // Function to update the username
  export async function updateUsername(username: string): Promise<any> {
    const response = await fetch('http://localhost:8080/api/user/update-username', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify({ username }),
    });
  
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to update username');
    }
  
    return await response.json();
  }
  
  // Function to validate if a username is available
  export async function checkUsernameAvailability(username: string): Promise<boolean> {
    try {
      const response = await fetch(`http://localhost:8080/api/search?q=${encodeURIComponent(username)}&f=unique-user`, {
        method: 'GET',
        credentials: 'include',
      });
      
      if (!response.ok) {
        return false;
      }
      
      const data = await response.json();
      return !data.exists; // Return true if username doesn't exist (is available)
    } catch (error) {
      console.error('Error checking username availability:', error);
      return false;
    }
  }