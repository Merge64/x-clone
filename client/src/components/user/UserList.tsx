import { FollowCard } from "../FollowCard";

interface UserListProps {
  users: any[];
  emptyMessage?: string;
  className?: string;
}

export function UserList({ users, emptyMessage, className }: UserListProps) {
  // Agregar validación adicional
  const validUsers = users?.filter((user) => user?.id && user?.username) || [];

  if (validUsers.length === 0) {
    return (
      <div className="p-6 text-center text-gray-500">
        <p>{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className={`${className}`}>
      {validUsers.map((user) => (
        <FollowCard
          key={user.id}
          {...user}
          nickname={user.nickname || user.displayName || user.username}
        />
      ))}
    </div>
  );
}
