import { useDashboardData } from '../hooks/useDashboardStats';
import { Users, UserPlus, Shield, Activity } from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, Legend, CartesianGrid } from 'recharts';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

const DashboardPage = () => {
  const { stats, roleDistribution, recentActivity, userAnalytics, isLoading, error } = useDashboardData();

  if (isLoading) {
    return <div className="p-6 text-center">Loading dashboard data...</div>;
  }

  if (error) {
    return <div className="p-6 text-center text-red-500">Error: {error}</div>;
  }

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">Dashboard Overview</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-md flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">Total Users</p>
            <p className="text-3xl font-bold text-gray-900">{stats?.total_users || 0}</p>
          </div>
          <div className="p-3 bg-blue-100 rounded-full text-blue-600">
            <Users className="w-6 h-6" />
          </div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">Active Users</p>
            <p className="text-3xl font-bold text-gray-900">{stats?.active_users || 0}</p>
          </div>
          <div className="p-3 bg-green-100 rounded-full text-green-600">
            <UserPlus className="w-6 h-6" />
          </div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">New Registrations</p>
            <p className="text-3xl font-bold text-gray-900">{stats?.new_users_today || 0}</p>
          </div>
          <div className="p-3 bg-yellow-100 rounded-full text-yellow-600">
            <UserPlus className="w-6 h-6" />
          </div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">Total Roles</p>
            <p className="text-3xl font-bold text-gray-900">{stats?.total_roles || 0}</p>
          </div>
          <div className="p-3 bg-red-100 rounded-full text-red-600">
            <Shield className="w-6 h-6" />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Role Distribution Chart */}
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">User Role Distribution</h2>
          {roleDistribution && roleDistribution.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={roleDistribution}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="user_count"
                  nameKey="role_name"
                  label={({ role_name, percent }) => `${role_name} (${(percent * 100).toFixed(0)}%)`}
                >
                  {roleDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-500">No role distribution data available.</p>
          )}
        </div>

        {/* User Analytics Chart */}
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">New User Analytics</h2>
          {userAnalytics && userAnalytics.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={userAnalytics}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="user_count" fill="#8884d8" name="New Users" />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-500">No user analytics data available.</p>
          )}
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Recent Activity</h2>
        {recentActivity && recentActivity.length > 0 ? (
          <ul className="divide-y divide-gray-200">
            {recentActivity.map((activity) => (
              <li key={activity.id} className="py-3 flex justify-between items-center">
                <span className="text-gray-700">
                  User '<span className="font-medium">{activity.user?.first_name} {activity.user?.last_name || activity.user?.username || activity.user?.email || 'Unknown'}</span>' {activity.action} {activity.resource} from {activity.ip_address}
                </span>
                <span className="text-sm text-gray-500">{new Date(activity.created_at).toLocaleString()}</span>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-500">No recent activity.</p>
        )}
      </div>
    </div>
  );
};

export default DashboardPage;
