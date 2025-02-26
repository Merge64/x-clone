// Main.tsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import LoginSignupPage from './views/LoginSignupPage';
import { LoginModal } from './components/LoginModal';
import { SignupModal } from './components/SignupModal'; // You’ll create this similarly to LoginModal

import './index.css';

const rootElement = document.getElementById('root');
if (!rootElement) {
  throw new Error("Failed to find the root element");
}

ReactDOM.createRoot(rootElement).render(
  <BrowserRouter>
    <Routes>
      {/* Top-level route always shows the main background page */}
      <Route path="/" element={<LoginSignupPage />}>
        {/* Nested routes show the modals on top */}
        <Route path="i/flow/login" element={<LoginModal />} />
        <Route path="i/flow/signup" element={< SignupModal />} />
      </Route>
    </Routes>
  </BrowserRouter>
);
