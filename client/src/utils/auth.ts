export async function validateToken(): Promise<boolean> {
    try {
      const response = await fetch('http://localhost:8080/api/auth/validate', {
        method: 'GET',
        credentials: 'include',
      });
      
      return response.ok;
    } catch (error) {
      console.error('Error validating token:', error);
      return false;
    }
  }
  

  export async function logout(): Promise<void> {
    try {
      await fetch('http://localhost:8080/api/auth/logout', {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Error logging out:', error);
    }
  }
  
  export async function checkExists(query: string, type: 'user' | 'mail'): Promise<boolean> {
    try {
      const filter = type === 'user' ? 'unique-user' : 'unique-mail';
      const response = await fetch(`http://localhost:8080/api/search?q=${encodeURIComponent(query)}&f=${filter}`, {
        method: 'GET',
        credentials: 'include',
      });
      
      if (!response.ok) {
        return false;
      }
      
      const data = await response.json();
      return data.exists;
    } catch (error) {
      console.error(`Error checking if ${type} exists:`, error);
      return false;
    }
  }
  
  export async function login(usernameOrEmail: string, password: string): Promise<{success: boolean, error?: string}> {
    try {
      const response = await fetch("http://localhost:8080/api/i/flow/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username_or_email: usernameOrEmail,
          password: password
        }),
        credentials: "include",
      });
  
      if (!response.ok) {
        const errorData = await response.json();
        return {
          success: false,
          error: errorData.error === "Invalid credentials"
            ? "Wrong password"
            : errorData.error || "Login failed"
        };
      }
  
      return { success: true };
    } catch (error) {
      console.error("Login failed:", error);
      return { 
        success: false, 
        error: "Something went wrong during login. Please try again."
      };
    }
  }