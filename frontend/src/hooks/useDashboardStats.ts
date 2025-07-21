import { useState, useEffect } from 'react';
import api from '../services/api';

interface DashboardStats {
  total_users: number;
  active_users: number;
  inactive_users: number;
  new_users_today: number;
  new_users_this_week: number;
  total_roles: number;
}

interface RoleDistribution {
  role_name: string;
  user_count: number;
}

interface ActivityLog {
  id: number;
  action: string;
  resource: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
  user: {
    id: number;
    username: string;
    email: string;
    first_name: string;
    last_name: string;
  };
}

interface UserAnalytics {
  date: string;
  user_count: number;
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

        console.log('Dashboard Stats Response:', statsRes.data);
        console.log('Role Distribution Response:', roleDistRes.data);
        console.log('Recent Activity Response:', recentActivityRes.data);
        console.log('User Analytics Response:', userAnalyticsRes.data);

        setData({
          stats: statsRes.data?.data || null,
          roleDistribution: Array.isArray(roleDistRes.data?.data) ? roleDistRes.data.data : [],
          recentActivity: Array.isArray(recentActivityRes.data?.data) ? recentActivityRes.data.data : [],
          userAnalytics: Array.isArray(userAnalyticsRes.data?.data) ? userAnalyticsRes.data.data : [],
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
