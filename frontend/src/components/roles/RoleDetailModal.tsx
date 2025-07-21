import { Dialog, Transition } from '@headlessui/react';
import { Fragment } from 'react';
import { Role } from '../../types/rbac';
import { X, Shield, Users, Key } from 'lucide-react';

interface RoleDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  role: Role | null;
}

const RoleDetailModal: React.FC<RoleDetailModalProps> = ({ isOpen, onClose, role }) => {
  if (!role) return null;

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={onClose}>
        <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex items-center justify-center min-h-full p-4 text-center">
            <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
              <Dialog.Panel className="w-full max-w-2xl p-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
                <div className="flex items-center justify-between mb-6">
                  <Dialog.Title as="h3" className="text-xl font-semibold leading-6 text-gray-900 flex items-center">
                    <Shield className="w-6 h-6 mr-2 text-indigo-600" />
                    Role Details: {role.name}
                  </Dialog.Title>
                  <button 
                    onClick={onClose}
                    className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
                  >
                    <X className="w-5 h-5" />
                  </button>
                </div>

                <div className="space-y-6">
                  {/* Basic Info */}
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <h4 className="font-medium text-gray-900 mb-3 flex items-center">
                      <Shield className="w-4 h-4 mr-2" />
                      Basic Information
                    </h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700">Role Name</label>
                        <div className="mt-1 flex items-center">
                          <span className="text-sm text-gray-900">{role.name}</span>
                          {role.is_system_role && (
                            <span className="ml-2 px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full">System Role</span>
                          )}
                        </div>
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700">Type</label>
                        <p className="mt-1 text-sm text-gray-900">
                          {role.is_system_role ? 'System Role' : 'Custom Role'}
                        </p>
                      </div>
                    </div>
                    <div className="mt-4">
                      <label className="block text-sm font-medium text-gray-700">Description</label>
                      <p className="mt-1 text-sm text-gray-900">{role.description || 'No description provided'}</p>
                    </div>
                    <div className="mt-4">
                      <label className="block text-sm font-medium text-gray-700">Created</label>
                      <p className="mt-1 text-sm text-gray-900">
                        {new Date(role.created_at).toLocaleDateString('en-US', { 
                          year: 'numeric', 
                          month: 'long', 
                          day: 'numeric' 
                        })}
                      </p>
                    </div>
                  </div>

                  {/* Permissions */}
                  <div className="bg-indigo-50 p-4 rounded-lg">
                    <h4 className="font-medium text-gray-900 mb-3 flex items-center">
                      <Key className="w-4 h-4 mr-2" />
                      Permissions ({role.permissions?.length || 0})
                    </h4>
                    {role.permissions && role.permissions.length > 0 ? (
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        {role.permissions.map((permission) => (
                          <div 
                            key={permission.id} 
                            className="bg-white p-3 rounded-md border border-indigo-200 hover:border-indigo-300 transition-colors"
                          >
                            <div className="flex items-start">
                              <Key className="w-4 h-4 mt-0.5 mr-2 text-indigo-600 flex-shrink-0" />
                              <div className="flex-1 min-w-0">
                                <h5 className="text-sm font-medium text-gray-900 truncate">
                                  {permission.name}
                                </h5>
                                <p className="text-xs text-gray-500 mt-1">
                                  {permission.description}
                                </p>
                                <div className="flex items-center mt-2">
                                  <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-gray-100 text-gray-800">
                                    {permission.resource}
                                  </span>
                                  <span className="ml-2 inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-indigo-100 text-indigo-800">
                                    {permission.action}
                                  </span>
                                </div>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <div className="text-center py-4">
                        <Key className="mx-auto h-12 w-12 text-gray-400" />
                        <p className="mt-2 text-sm text-gray-500">No permissions assigned to this role.</p>
                      </div>
                    )}
                  </div>

                  {/* Stats */}
                  <div className="bg-green-50 p-4 rounded-lg">
                    <h4 className="font-medium text-gray-900 mb-3 flex items-center">
                      <Users className="w-4 h-4 mr-2" />
                      Usage Statistics
                    </h4>
                    <div className="grid grid-cols-2 gap-4">
                      <div className="text-center">
                        <div className="text-2xl font-bold text-green-600">
                          {role.permissions?.length || 0}
                        </div>
                        <div className="text-sm text-gray-600">Total Permissions</div>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-blue-600">
                          --
                        </div>
                        <div className="text-sm text-gray-600">Users Assigned</div>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="mt-6 flex justify-end">
                  <button 
                    onClick={onClose}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-transparent rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 transition-colors"
                  >
                    Close
                  </button>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
};

export default RoleDetailModal;