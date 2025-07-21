import { useState, useEffect, useCallback } from 'react';
import api from '../services/api';
import { Role } from '../types/rbac';

interface RolesResponse {
  data: {
    roles: Role[];
  }
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
      setRoles(response.data.data.roles);
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
