import { useState, useEffect, useCallback } from 'react';
import api from '../services/api';
import { User } from '../types/user';

interface UsersResponse {
  data: {
    users: User[];
    total: number;
    page: number;
    limit: number;
  }
}

export const useUsers = (page = 1, limit = 10, search = '') => {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await api.get<UsersResponse>('/users', {
        params: { page, limit, search },
      });
      setUsers(response.data.data.users);
      setTotal(response.data.data.total);
    } catch (err) {
      setError('Failed to fetch users');
      console.error(err);
    }
    setIsLoading(false);
  }, [page, limit, search]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  return { users, total, isLoading, error, mutate: fetchUsers };
};
