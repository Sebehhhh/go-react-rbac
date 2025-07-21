import { Dialog, Transition } from '@headlessui/react';
import { Fragment, useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Role } from '../../types/rbac';
import { usePermissions } from '../../hooks/usePermissions';

const roleSchema = z.object({
  name: z.string().min(1, 'Role name is required'),
  description: z.string().optional(),
  permission_ids: z.array(z.number()).optional(),
});

type RoleFormInputs = z.infer<typeof roleSchema>;

interface RoleFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  role: Role | null;
  onSave: (data: RoleFormInputs, roleId: number | null) => void;
}

const RoleFormModal: React.FC<RoleFormModalProps> = ({ isOpen, onClose, role, onSave }) => {
  const { permissions } = usePermissions();
  const [selectedPermissions, setSelectedPermissions] = useState<number[]>([]);
  const { register, handleSubmit, formState: { errors }, reset } = useForm<RoleFormInputs>({
    resolver: zodResolver(roleSchema),
  });

  useEffect(() => {
    if (isOpen) {
      if (role) {
        const rolePermissionIds = role.permissions?.map(p => p.id) || [];
        reset({
          name: role.name,
          description: role.description,
        });
        setSelectedPermissions(rolePermissionIds);
      } else {
        reset({ name: '', description: '' });
        setSelectedPermissions([]);
      }
    }
  }, [role, isOpen, reset]);

  const handlePermissionChange = (permissionId: number, checked: boolean) => {
    if (checked) {
      setSelectedPermissions(prev => [...prev, permissionId]);
    } else {
      setSelectedPermissions(prev => prev.filter(id => id !== permissionId));
    }
  };

  const onSubmit = (data: RoleFormInputs) => {
    const formData = { ...data, permission_ids: selectedPermissions };
    onSave(formData, role ? role.id : null);
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={onClose}>
        <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex items-center justify-center min-h-full p-4 text-center">
            <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-95">
              <Dialog.Panel className="w-full max-w-md p-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
                <Dialog.Title as="h3" className="text-lg font-medium leading-6 text-gray-900">
                  {role ? 'Edit Role' : 'Create Role'}
                </Dialog.Title>
                <form onSubmit={handleSubmit(onSubmit)} className="mt-4 space-y-4">
                  <div>
                    <label htmlFor="name">Role Name</label>
                    <input id="name" {...register('name')} className="w-full px-3 py-2 mt-1 border rounded-md" />
                    {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
                  </div>
                  <div>
                    <label htmlFor="description">Description</label>
                    <textarea id="description" {...register('description')} className="w-full px-3 py-2 mt-1 border rounded-md"></textarea>
                  </div>
                  <div>
                    <h4 className="mb-2 text-md font-medium">Permissions</h4>
                    <div className="max-h-60 overflow-y-auto border border-gray-200 rounded-md p-3">
                      <div className="grid grid-cols-1 gap-2">
                        {permissions && permissions.map(perm => (
                          <div key={perm.id} className="flex items-center">
                            <input
                              type="checkbox"
                              id={`perm-${perm.id}`}
                              checked={selectedPermissions.includes(perm.id)}
                              onChange={(e) => handlePermissionChange(perm.id, e.target.checked)}
                              className="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
                            />
                            <label htmlFor={`perm-${perm.id}`} className="ml-2 text-sm text-gray-900">
                              <span className="font-medium">{perm.name}</span> - {perm.description}
                            </label>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                  <div className="mt-6 flex justify-end space-x-2">
                    <button type="button" onClick={onClose} className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-transparent rounded-md hover:bg-gray-200">Cancel</button>
                    <button type="submit" className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700">Save</button>
                  </div>
                </form>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
};

export default RoleFormModal;
