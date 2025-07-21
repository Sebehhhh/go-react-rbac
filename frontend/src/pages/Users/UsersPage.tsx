import { useState } from 'react';
import { useUsers } from '../../hooks/useUsers';
import UserTable from '../../components/users/UserTable';
import UserFormModal from '../../components/users/UserFormModal';
import { User } from '../../types/user';
import api from '../../services/api';

const UsersPage = () => {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const { users, total, isLoading, error, mutate } = useUsers(page, 10, search);
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const handleOpenModal = (user: User | null) => {
    setEditingUser(user);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setEditingUser(null);
    setIsModalOpen(false);
  };

  const handleSave = async (data: any, userId: number | null) => {
    try {
      if (userId) {
        await api.put(`/users/${userId}`, data);
      } else {
        await api.post('/users', data);
      }
      mutate(); // Re-fetch users
      handleCloseModal();
    } catch (error) {
      console.error('Failed to save user:', error);
    }
  };

  const handleDelete = async (user: User) => {
    if (window.confirm(`Are you sure you want to delete ${user.first_name}?`)) {
      try {
        await api.delete(`/users/${user.id}`);
        mutate();
      } catch (error) {
        console.error('Failed to delete user:', error);
      }
    }
  };

  const handleToggleActivate = async (user: User) => {
    const action = user.is_active ? 'deactivate' : 'activate';
    if (window.confirm(`Are you sure you want to ${action} ${user.first_name}?`)) {
      try {
        await api.put(`/users/${user.id}/${action}`);
        mutate();
      } catch (error) {
        console.error(`Failed to ${action} user:`, error);
      }
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
        />
      )}
      {/* Add pagination controls here */}
      <UserFormModal 
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        user={editingUser}
        onSave={handleSave}
      />
    </div>
  );
};

export default UsersPage;
