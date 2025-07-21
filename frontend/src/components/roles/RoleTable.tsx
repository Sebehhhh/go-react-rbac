import { Role } from '../../types/rbac';
import { Edit, Trash2, Eye } from 'lucide-react';

interface RoleTableProps {
  roles: Role[];
  onEdit: (role: Role) => void;
  onDelete: (role: Role) => void;
  onViewDetails: (role: Role) => void;
}

const RoleTable: React.FC<RoleTableProps> = ({ roles, onEdit, onDelete, onViewDetails }) => {
  return (
    <div className="overflow-x-auto bg-white rounded-lg shadow">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Name</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Description</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Type</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">Permissions</th>
            <th className="px-6 py-3 text-xs font-medium tracking-wider text-right text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {roles && roles.length > 0 ? (
            roles.map((role) => (
              <tr key={role.id}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <span className="font-medium">{role.name}</span>
                    {role.is_system_role && (
                      <span className="ml-2 px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full">System</span>
                    )}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{role.description}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs rounded-full ${role.is_system_role ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800'}`}>
                    {role.is_system_role ? 'System Role' : 'Custom Role'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                    {role.permissions?.length || 0} permissions
                  </span>
                </td>
                <td className="px-6 py-4 text-sm font-medium text-right whitespace-nowrap">
                  <button onClick={() => onViewDetails(role)} className="text-blue-600 hover:text-blue-900">
                    <Eye className="w-5 h-5" />
                  </button>
                  <button onClick={() => onEdit(role)} className="ml-4 text-indigo-600 hover:text-indigo-900">
                    <Edit className="w-5 h-5" />
                  </button>
                  {!role.is_system_role && (
                    <button onClick={() => onDelete(role)} className="ml-4 text-red-600 hover:text-red-900">
                      <Trash2 className="w-5 h-5" />
                    </button>
                  )}
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                No roles found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default RoleTable;
