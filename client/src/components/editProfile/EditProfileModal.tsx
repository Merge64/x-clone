import React, { useEffect, useState } from 'react';
import Modal from './Modal';

interface EditProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialData: {
    nickname: string;
    bio?: string;
    location?: string;
  };
  onSave: (data: {
    nickname: string;
    bio?: string;
    location?: string;
  }) => void;
}

const EditProfileModal: React.FC<EditProfileModalProps> = ({
  isOpen,
  onClose,
  initialData,
  onSave,
}) => {
  const [formData, setFormData] = useState(initialData);

  useEffect(() => {
    setFormData(initialData);
  }, [initialData]);

  const handleSave = () => {
    onSave(formData);
    onClose();
  };
  console.log(initialData)
  return (
    <Modal isOpen={isOpen} onClose={onClose} onSave={handleSave} title="Edit profile">
      <div className="space-y-4">
        <div>
          <label className="block text-gray-500 text-sm mb-2" htmlFor="name">
            Name
          </label>
          <input
            type="text"
            id="name"
            value={formData.nickname}
            onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
            className="w-full bg-transparent border border-gray-800 rounded-md p-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            maxLength={50}
          />
        </div>

        <div>
          <label className="block text-gray-500 text-sm mb-2" htmlFor="bio">
            Bio
          </label>
          <textarea
            id="bio"
            value={formData.bio || ''}
            onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
            className="w-full bg-transparent border border-gray-800 rounded-md p-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            rows={3}
            maxLength={160}
          />
        </div>

        <div>
          <label className="block text-gray-500 text-sm mb-2" htmlFor="location">
            Location
          </label>
          <input
            type="text"
            id="location"
            value={formData.location || ''}
            onChange={(e) => setFormData({ ...formData, location: e.target.value })}
            className="w-full bg-transparent border border-gray-800 rounded-md p-2 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            maxLength={30}
          />
        </div>
      </div>
    </Modal>
  );
};

export default EditProfileModal;