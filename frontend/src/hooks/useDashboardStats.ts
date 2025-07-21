import { useState, useEffect } from 'react';
import api from '../services/api';

interface DashboardStats {
  totalUsers: number;
  activeUsers: number;
  newRegistrations: number;
  totalRoles: number;
}

interface RoleDistribution {
  roleName: string;
  userCount: number;
}

interface ActivityLog {
  id: number;
  action: string;
  resource: string;
  ipAddress: string;
  userAgent: string;
  createdAt: string;
  user: { // Assuming user details are nested
    id: number;
    username: string;
    email: string;
  };
}

interface UserAnalytics {
  date: string;
  newUsers: number;
}

interface DashboardData {
  stats: DashboardStats | null;
  roleDistribution: RoleDistribution[] | null;
  recentActivity: ActivityLog[] | null;
  userAnalytics: UserAnalytics[] | null;
  isLoading: boolean;
  error: string | null;
}

export const useDashboardData = () => {
  const [data, setData] = useState<DashboardData>({
    stats: null,
    roleDistribution: null,
    recentActivity: null,
    userAnalytics: null,
    isLoading: true,
    error: null,
  });

  useEffect(() => {
    const fetchData = async () => {
      setData((prev) => ({ ...prev, isLoading: true, error: null }));
      try {
        const [statsRes, roleDistRes, recentActivityRes, userAnalyticsRes] = await Promise.all([
          api.get('/dashboard/stats'),
          api.get('/dashboard/role-distribution'),
          api.get('/dashboard/recent-activity'),
          api.get('/dashboard/user-analytics'),
        ]);

        setData({
          stats: statsRes.data.data,
          roleDistribution: roleDistRes.data.data,
          recentActivity: recentActivityRes.data.data,
          userAnalytics: userAnalyticsRes.data.data,
          isLoading: false,
          error: null,
        });
      } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
        setData((prev) => ({ ...prev, isLoading: false, error: 'Failed to load dashboard data.' }));
      }
    };

    fetchData();
  }, []);

  return data;
};
