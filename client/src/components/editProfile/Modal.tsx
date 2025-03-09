import React from 'react';
import { X } from 'lucide-react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave?: () => void;
  title: string;
  children: React.ReactNode;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, onSave, title, children }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className="bg-black w-full max-w-xl rounded-xl border border-gray-800">
        <div className="flex items-center justify-between p-4 border-b border-gray-800">
          <div className="flex items-center gap-8">
            <button onClick={onClose} className="hover:bg-gray-900 p-2 rounded-full">
              <X size={20} />
            </button>
            <h2 className="text-xl font-bold">{title}</h2>
          </div>
          {onSave && (
            <button
              onClick={onSave}
              className="px-4 py-1 bg-white text-black rounded-full font-bold hover:bg-gray-200"
            >
              Save
            </button>
          )}
        </div>
        <div className="p-4">{children}</div>
      </div>
    </div>
  );
}

export default Modal