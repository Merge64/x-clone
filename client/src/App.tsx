import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LoginSignupPage from './views/LoginSignupPage';
import LoginModal from './components/LoginModal';
import ChangeUsername from './components/ChangeUsername';
import HomePage from './views/HomePage';
import AuthChecker from './utils/AuthChecker';
import SignupModal from './components/SignupModal';
import ProfilePage from './views/ProfilePage';
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
          <Route path="/c" element={
          <AuthChecker>
            <ChangeUsername />
          </AuthChecker>
        } />
        <Route 
          path="/profile/:username" 
          element={
            <AuthChecker>
              <ProfilePage />
            </AuthChecker>
          } 
        />
        <Route 
          path="/post/:username/:postId" 
          element={
            <AuthChecker>
              <PostDetailPage />
            </AuthChecker>
          } 
        />

        {/* Redirect any unknown routes to home */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;