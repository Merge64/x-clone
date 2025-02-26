// views/LoginSignupPage.tsx
import React from 'react';
import { Outlet } from 'react-router-dom';
import MainBackground from '../components/MainBackground';

function LoginSignupPage() {
  return (
    <>
      {/* Show your main background layout */}
      <MainBackground />

      {/* Where nested routes (the modals) get injected */}
      <Outlet />
    </>
  );
}

export default LoginSignupPage;
