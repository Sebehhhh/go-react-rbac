import { useState, useEffect } from 'react';
import api from '../services/api';
import { Permission } from '../types/rbac';

export const usePermissions = () => {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPermissions = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await api.get<{ success: boolean; message: string; data: Permission[] }>('/permissions');
        setPermissions(response.data.data);
      } catch (err) {
        setError('Failed to fetch permissions');
        console.error(err);
      }
      setIsLoading(false);
    };

    fetchPermissions();
  }, []);

  return { permissions, isLoading, error };
};
