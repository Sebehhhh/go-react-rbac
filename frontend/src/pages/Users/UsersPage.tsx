import { useState } from 'react';
import { useUsers } from '../../hooks/useUsers';
import UserTable from '../../components/users/UserTable';
import UserFormModal from '../../components/users/UserFormModal';
import PasswordUpdateModal from '../../components/users/PasswordUpdateModal';
import { User } from '../../types/user';
import api from '../../services/api';
import { showSuccessAlert, showErrorAlert, showConfirmAlert } from '../../utils/alerts';

const UsersPage = () => {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const { users, total, isLoading, error, mutate } = useUsers(page, 10, search);
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false);
  const [passwordUpdateUser, setPasswordUpdateUser] = useState<User | null>(null);

  const handleOpenModal = (user: User | null) => {
    setEditingUser(user);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setEditingUser(null);
    setIsModalOpen(false);
  };

  const handleOpenPasswordModal = (user: User) => {
    setPasswordUpdateUser(user);
    setIsPasswordModalOpen(true);
  };

  const handleClosePasswordModal = () => {
    setPasswordUpdateUser(null);
    setIsPasswordModalOpen(false);
  };

  const handleSave = async (data: any, userId: number | null) => {
    try {
      if (userId) {
        await api.put(`/users/${userId}`, data);
        showSuccessAlert('User updated successfully!');
      } else {
        await api.post('/users', data);
        showSuccessAlert('User created successfully!');
      }
      mutate(); // Re-fetch users
      handleCloseModal();
    } catch (error: any) {
      showErrorAlert(error.response?.data?.message || 'Failed to save user.');
    }
  };

  const handleDelete = async (user: User) => {
    const confirmed = await showConfirmAlert('Delete User', `Are you sure you want to delete ${user.first_name}?`);
    if (confirmed) {
      try {
        await api.delete(`/users/${user.id}`);
        mutate();
        showSuccessAlert('User deleted successfully!');
      } catch (error: any) {
        showErrorAlert(error.response?.data?.message || 'Failed to delete user.');
      }
    }
  };

  const handleToggleActivate = async (user: User) => {
    const action = user.is_active ? 'deactivate' : 'activate';
    const confirmed = await showConfirmAlert(
      `${action === 'deactivate' ? 'Deactivate' : 'Activate'} User`,
      `Are you sure you want to ${action} ${user.first_name}?`
    );
    if (confirmed) {
      try {
        await api.put(`/users/${user.id}/${action}`);
        mutate();
        showSuccessAlert(`User ${action}d successfully!`);
      } catch (error: any) {
        showErrorAlert(error.response?.data?.message || `Failed to ${action} user.`);
      }
    }
  };

  const handleUpdatePassword = async (userId: number, passwordData: any) => {
    try {
      await api.put(`/users/${userId}/password`, passwordData);
      showSuccessAlert('Password updated successfully!');
      handleClosePasswordModal();
    } catch (error: any) {
      showErrorAlert(error.response?.data?.message || 'Failed to update password.');
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">Users Management</h1>
        <button onClick={() => handleOpenModal(null)} className="px-4 py-2 font-bold text-white bg-indigo-600 rounded-md hover:bg-indigo-700">
          Add User
        </button>
      </div>
      <div className="mb-4">
        <input 
          type="text"
          placeholder="Search users..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full px-3 py-2 border rounded-md"
        />
      </div>
      {isLoading && <p>Loading...</p>}
      {error && <p className="text-red-500">{error}</p>}
      {!isLoading && !error && (
        <UserTable 
          users={users}
          onEdit={handleOpenModal}
          onDelete={handleDelete}
          onToggleActivate={handleToggleActivate}
          onUpdatePassword={handleOpenPasswordModal}
        />
      )}
      {/* Add pagination controls here */}
      <UserFormModal 
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        user={editingUser}
        onSave={handleSave}
      />
      {passwordUpdateUser && (
        <PasswordUpdateModal 
          isOpen={isPasswordModalOpen}
          onClose={handleClosePasswordModal}
          userId={passwordUpdateUser.id}
          userFullName={`${passwordUpdateUser.first_name} ${passwordUpdateUser.last_name}`}
          onUpdatePassword={handleUpdatePassword}
        />
      )}
    </div>
  );
};

export default UsersPage;
