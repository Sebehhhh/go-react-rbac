import { useState } from 'react';
import { useRoles } from '../../hooks/useRoles';
import RoleTable from '../../components/roles/RoleTable';
import RoleFormModal from '../../components/roles/RoleFormModal';
import RoleDetailModal from '../../components/roles/RoleDetailModal';
import { Role } from '../../types/rbac';
import api from '../../services/api';
import { showSuccessAlert, showErrorAlert, showConfirmAlert } from '../../utils/alerts';

const RolesPage = () => {
  const { roles, isLoading, error, mutate } = useRoles();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [detailRole, setDetailRole] = useState<Role | null>(null);

  const handleOpenModal = (role: Role | null) => {
    setEditingRole(role);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setEditingRole(null);
    setIsModalOpen(false);
  };

  const handleViewDetails = (role: Role) => {
    setDetailRole(role);
    setIsDetailModalOpen(true);
  };

  const handleCloseDetailModal = () => {
    setDetailRole(null);
    setIsDetailModalOpen(false);
  };

  const handleSave = async (data: any, roleId: number | null) => {
    try {
      if (roleId) {
        await api.put(`/roles/${roleId}`, { name: data.name, description: data.description });
        await api.put(`/roles/${roleId}/permissions`, { permission_ids: data.permission_ids });
        showSuccessAlert('Role updated successfully!');
      } else {
        const response = await api.post('/roles', { name: data.name, description: data.description });
        if (response.data.data && data.permission_ids && data.permission_ids.length > 0) {
          await api.put(`/roles/${response.data.data.id}/permissions`, { permission_ids: data.permission_ids });
        }
        showSuccessAlert('Role created successfully!');
      }
      mutate(); // Re-fetch roles
      handleCloseModal();
    } catch (error: any) {
      showErrorAlert(error.response?.data?.message || 'Failed to save role.');
    }
  };

  const handleDelete = async (role: Role) => {
    const confirmed = await showConfirmAlert('Delete Role', `Are you sure you want to delete role ${role.name}?`);
    if (confirmed) {
      try {
        await api.delete(`/roles/${role.id}`);
        mutate();
        showSuccessAlert('Role deleted successfully!');
      } catch (error: any) {
        showErrorAlert(error.response?.data?.message || 'Failed to delete role.');
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
          onViewDetails={handleViewDetails}
        />
      )}
      <RoleFormModal 
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        role={editingRole}
        onSave={handleSave}
      />
      <RoleDetailModal 
        isOpen={isDetailModalOpen}
        onClose={handleCloseDetailModal}
        role={detailRole}
      />
    </div>
  );
};

export default RolesPage;
