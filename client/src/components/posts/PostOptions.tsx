import React from 'react';

interface PostOptionsProps {
  onEdit?: (e: React.MouseEvent) => void;
  onDelete: (e: React.MouseEvent) => void;
  isDeleting: boolean;
}

const PostOptions: React.FC<PostOptionsProps> = ({ onEdit, onDelete, isDeleting }) => {
  return (
    <div className="absolute right-0 top-0 mt-8 bg-black rounded-lg shadow-lg border border-gray-800 z-10 overflow-hidden w-40">
      {onEdit && (
        <button
          onClick={onEdit}
          className="w-full text-left px-4 py-3 hover:bg-gray-900 text-white"
        >
          Edit
        </button>
      )}
      <button
        onClick={onDelete}
        className={`w-full text-left px-4 py-3 hover:bg-gray-900 ${
          isDeleting ? 'text-gray-500' : 'text-red-500'
        }`}
        disabled={isDeleting}
      >
        {isDeleting ? 'Deleting...' : 'Delete'}
      </button>
    </div>
  );
};

export default PostOptions;