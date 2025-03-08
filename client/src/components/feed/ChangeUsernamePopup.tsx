import ChangeUsername from '../auth/ChangeUsername';
import { useNavigate } from 'react-router-dom';

interface UsernamePopupProps {
  isOpen: boolean;
  onClose?: () => void;
}

function UsernamePopup({ isOpen, onClose }: UsernamePopupProps) {
  const navigate = useNavigate();
  
  if (!isOpen) return null;
  
  const handleClose = () => {
    if (onClose) {
      onClose();
    } else {
      // If no onClose provided, navigate back to home
      navigate('/home');
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-[#242D34]/60">
      <div className="w-full max-w-md">
        <ChangeUsername onSkip={handleClose} />
      </div>
    </div>
  );
}

export default UsernamePopup;