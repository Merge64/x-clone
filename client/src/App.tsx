import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LoginSignupPage from './views/LoginSignupPage';
import LoginModal from './components/auth/LoginModal';
import HomePage from './components/feed/HomePage';
import AuthChecker from './components/auth/AuthChecker';
import SignupModal from './components/auth/SignupModal';
import ProfilePage from './views/ProfilePage';
import UsernamePopup from './components/feed/UsernamePopup';
import PostDetailPage from './views/PostDetailPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public routes */}
        <Route path="/" element={<LoginSignupPage />}>
          <Route path="i/flow/login" element={<LoginModal />} />
          <Route path="i/flow/signup" element={<SignupModal />} />
        </Route>

        {/* Protected routes */}
        <Route path="/home" element={
          <AuthChecker>
            <HomePage />
          </AuthChecker>
        } />
        
        {/* Change username route as overlay */}
        <Route path="/change-username" element={
          <AuthChecker>
            <HomePage />
            <UsernamePopup isOpen={true} />
          </AuthChecker>
        } />
        
        {/* Direct username access route */}
        <Route 
          path="/:username" 
          element={
            <AuthChecker>
              <ProfilePage />
            </AuthChecker>
          } 
        />
        
        <Route 
          path="/:username/:postId" 
          element={
              <PostDetailPage />
          } 
        />

        {/* Redirect any unknown routes to home */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;