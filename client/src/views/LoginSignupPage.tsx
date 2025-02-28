import { Outlet } from 'react-router-dom';
import AuthBackground from '../components/AuthBackground';

function LoginSignupPage() {
  return (
    <>
      <AuthBackground />
      <Outlet /> {/* This will render the modals */}
    </>
  );
}

export default LoginSignupPage;