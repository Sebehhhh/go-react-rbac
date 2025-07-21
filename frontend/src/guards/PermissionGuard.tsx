import React from 'react';
import { useAuthStore } from '../hooks/useAuth';

interface PermissionGuardProps {
  permission: string; // e.g., "users.read", "roles.create"
  children: React.ReactNode;
}

const PermissionGuard: React.FC<PermissionGuardProps> = ({ permission, children }) => {
  const { user, permissions } = useAuthStore();

  // Super Admin always has access
  if (user?.role.name === 'Super Admin') {
    return <>{children}</>;
  }

  // Check if the user's permissions include the required permission
  if (permissions.includes(permission)) {
    return <>{children}</>;
  }

  return null; // Don't render children if permission is not met
};

export default PermissionGuard;
