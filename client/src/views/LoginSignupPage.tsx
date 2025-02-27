import { Outlet } from 'react-router-dom';
import MainBackground from '../components/AuthBackground';

function LoginSignupPage() {
  return (
    <>
      <MainBackground />
      <Outlet /> {/* This will render the modals */}
    </>
  );
}

export default LoginSignupPage;