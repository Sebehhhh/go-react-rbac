import { Edit, PowerOff, Power, Trash2, Lock } from 'lucide-react';
import { User } from '../../types/user';

interface UserTableProps {
  users: User[];
  onEdit: (user: User) => void;
  onDelete: (user: User) => void;
  onToggleActivate: (user: User) => void;
  onUpdatePassword: (user: User) => void;
}

const UserTable: React.FC<UserTableProps> = ({ users, onEdit, onDelete, onToggleActivate, onUpdatePassword }) => {
  return (
    <div className="overflow-x-auto bg-white rounded-lg shadow">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Name</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Username</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Email</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Role</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Status</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-right text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {users.map((user) => (
            <tr key={user.id}>
              <td className="px-6 py-4 whitespace-nowrap">{user.first_name} {user.last_name}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">@{user.username}</td>
              <td className="px-6 py-4 whitespace-nowrap">{user.email}</td>
              <td className="px-6 py-4 whitespace-nowrap">{user.role.name}</td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                  {user.is_active ? 'Active' : 'Inactive'}
                </span>
              </td>
              <td className="px-6 py-4 text-sm font-medium text-right whitespace-nowrap">
                <button onClick={() => onEdit(user)} className="text-indigo-600 hover:text-indigo-900">
                  <Edit className="w-5 h-5" />
                </button>
                <button onClick={() => onUpdatePassword(user)} className="ml-4 text-yellow-600 hover:text-yellow-900">
                  <Lock className="w-5 h-5" />
                </button>
                <button onClick={() => onToggleActivate(user)} className="ml-4 text-gray-600 hover:text-gray-900">
                  {user.is_active ? <PowerOff className="w-5 h-5 text-red-500" /> : <Power className="w-5 h-5 text-green-500" />}
                </button>
                <button onClick={() => onDelete(user)} className="ml-4 text-red-600 hover:text-red-900">
                  <Trash2 className="w-5 h-5" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default UserTable;
