import { useState, useEffect, useCallback } from 'react';
import api from '../services/api';
import { Role } from '../types/rbac';

interface RolesResponse {
  success: boolean;
  message: string;
  data: Role[];
}

export const useRoles = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRoles = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await api.get<RolesResponse>('/roles');
      if (response.data && response.data.data && Array.isArray(response.data.data)) {
        setRoles(response.data.data);
      } else {
        setRoles([]); // Ensure roles is always an array
      }
    } catch (err) {
      setError('Failed to fetch roles');
      console.error(err);
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    fetchRoles();
  }, [fetchRoles]);

  return { roles, isLoading, error, mutate: fetchRoles };
};
