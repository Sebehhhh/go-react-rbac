import { Dialog, Transition } from '@headlessui/react';
import { Fragment, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { User } from '../../types/user';
import { useRoles } from '../../hooks/useRoles';

const userSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  email: z.string().email(),
  role_id: z.preprocess((val) => Number(val), z.number().positive('Role is required')),
});

type UserFormInputs = z.infer<typeof userSchema>;

interface UserFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: User | null;
  onSave: (data: UserFormInputs, userId: number | null) => void;
}

const UserFormModal: React.FC<UserFormModalProps> = ({ isOpen, onClose, user, onSave }) => {
  const { roles } = useRoles();
  const { register, handleSubmit, formState: { errors }, reset } = useForm<UserFormInputs>({
    resolver: zodResolver(userSchema),
  });

  useEffect(() => {
    if (isOpen) {
      if (user) {
        reset({ ...user, role_id: user.role.id });
      } else {
        reset({ first_name: '', last_name: '', email: '', role_id: 0 });
      }
    }
  }, [user, isOpen, reset]);

  const onSubmit = (data: UserFormInputs) => {
    onSave(data, user ? user.id : null);
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={onClose}>
        <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex items-center justify-center min-h-full p-4 text-center">
            <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
              <Dialog.Panel className="w-full max-w-md p-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
                <Dialog.Title as="h3" className="text-lg font-medium leading-6 text-gray-900">
                  {user ? 'Edit User' : 'Create User'}
                </Dialog.Title>
                <form onSubmit={handleSubmit(onSubmit)} className="mt-4 space-y-4">
                  <div>
                    <label htmlFor="first_name">First Name</label>
                    <input id="first_name" {...register('first_name')} className="w-full px-3 py-2 mt-1 border rounded-md" />
                    {errors.first_name && <p className="text-sm text-red-500">{errors.first_name.message}</p>}
                  </div>
                  <div>
                    <label htmlFor="last_name">Last Name</label>
                    <input id="last_name" {...register('last_name')} className="w-full px-3 py-2 mt-1 border rounded-md" />
                    {errors.last_name && <p className="text-sm text-red-500">{errors.last_name.message}</p>}
                  </div>
                  <div>
                    <label htmlFor="email">Email</label>
                    <input id="email" type="email" {...register('email')} className="w-full px-3 py-2 mt-1 border rounded-md" />
                    {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}
                  </div>
                  <div>
                    <label htmlFor="role_id">Role</label>
                    <select id="role_id" {...register('role_id')} className="w-full px-3 py-2 mt-1 border rounded-md">
                      <option value={0}>Select a role</option>
                      {roles && roles.map(role => <option key={role.id} value={role.id}>{role.name}</option>)}
                    </select>
                    {errors.role_id && <p className="text-sm text-red-500">{errors.role_id.message}</p>}
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

export default UserFormModal;
