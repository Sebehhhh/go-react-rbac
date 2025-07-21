import { useState } from 'react';
import { useRoles } from '../../hooks/useRoles';
import RoleTable from '../../components/roles/RoleTable';
import RoleFormModal from '../../components/roles/RoleFormModal';
import { Role } from '../../types/rbac';
import api from '../../services/api';

const RolesPage = () => {
  const { roles, isLoading, error, mutate } = useRoles();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);

  const handleOpenModal = (role: Role | null) => {
    setEditingRole(role);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setEditingRole(null);
    setIsModalOpen(false);
  };

  const handleSave = async (data: any, roleId: number | null) => {
    try {
      if (roleId) {
        await api.put(`/roles/${roleId}`, { name: data.name, description: data.description });
        await api.put(`/roles/${roleId}/permissions`, { permission_ids: data.permission_ids });
      } else {
        const response = await api.post('/roles', { name: data.name, description: data.description });
        if (response.data.data && data.permission_ids && data.permission_ids.length > 0) {
          await api.put(`/roles/${response.data.data.id}/permissions`, { permission_ids: data.permission_ids });
        }
      }
      mutate(); // Re-fetch roles
      handleCloseModal();
    } catch (error) {
      console.error('Failed to save role:', error);
    }
  };

  const handleDelete = async (role: Role) => {
    if (window.confirm(`Are you sure you want to delete role ${role.name}?`)) {
      try {
        await api.delete(`/roles/${role.id}`);
        mutate();
      } catch (error) {
        console.error('Failed to delete role:', error);
      }
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">Roles Management</h1>
        <button onClick={() => handleOpenModal(null)} className="px-4 py-2 font-bold text-white bg-indigo-600 rounded-md hover:bg-indigo-700">
          Add Role
        </button>
      </div>
      {isLoading && <p>Loading...</p>}
      {error && <p className="text-red-500">{error}</p>}
      {!isLoading && !error && (
        <RoleTable 
          roles={roles}
          onEdit={handleOpenModal}
          onDelete={handleDelete}
        />
      )}
      <RoleFormModal 
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        role={editingRole}
        onSave={handleSave}
      />
    </div>
  );
};

export default RolesPage;
